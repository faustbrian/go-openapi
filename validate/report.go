// Package validate performs version-aware OpenAPI document validation and
// returns stable, machine-readable diagnostics. Document combines official
// schema checks with prose-specific rules, schema checks and supported reference
// diagnostics; SchemaInstance separately checks runtime values and direction.
// Decoding a typed view is not validation, and Report.Valid means no collected
// error diagnostics, not exhaustive external-resource or business correctness.
//
// External references/examples/codecs require explicit borrowed providers.
// Validation is synchronous and cooperative; a provider can block cancellation
// unless it honors context. FailFast and MaxDiagnostics truncate output after
// passes, not work or intermediate allocations. Validator caches can be reused
// concurrently when caller-supplied collaborators are also safe for that use.
package validate

import (
	"context"
	"errors"

	openapi "github.com/faustbrian/go-openapi/v2"
	openapischema "github.com/faustbrian/go-openapi/v2/jsonschema"
	"github.com/faustbrian/go-openapi/v2/jsonvalue"
	"github.com/faustbrian/go-openapi/v2/reference"
)

// ErrLimitExceeded reports a validation input or work bound exhaustion.
var ErrLimitExceeded = errors.New("OpenAPI validation limit exceeded")

// Source identifies the validation subsystem that produced a diagnostic.
type Source string

const (
	// SourceDocument identifies OpenAPI data-model validation.
	SourceDocument Source = "document"
	// SourceSchema identifies validation of an OpenAPI Schema Object.
	SourceSchema Source = "schema"
	// SourceReference identifies URI-reference target validation.
	SourceReference Source = "reference"
)

// Severity identifies diagnostic impact.
type Severity string

const (
	// SeverityError prevents a document from being valid.
	SeverityError Severity = "error"
	// SeverityWarning reports a non-fatal interoperability concern.
	SeverityWarning Severity = "warning"
)

// Diagnostic is one stable validation finding.
type Diagnostic struct {
	// Code identifies the rule or diagnostic class; callers need not parse Message.
	Code string
	// Message explains the finding and may contain document names or provider error text.
	Message string
	// Severity distinguishes validity-blocking errors from advisory warnings.
	Severity Severity
	// Source identifies the document, schema, or reference validation subsystem.
	Source Source
	// InstanceLocation is an escaped JSON Pointer in the validated input; empty selects the root.
	InstanceLocation string
	// KeywordLocation identifies a schema keyword when supplied; semantic rules may leave it empty.
	KeywordLocation string
	// AbsoluteKeywordLocation is the absolute schema keyword URI when the evaluator supplies one.
	AbsoluteKeywordLocation string
	// SpecificationVersion is the exact selected document revision.
	SpecificationVersion string
	// SpecificationSection names the relevant prose section when known.
	SpecificationSection string
}

// ExternalExampleResource is one explicitly retrieved serialized example.
type ExternalExampleResource struct {
	// RetrievalURI records the identity returned by the application resolver.
	RetrievalURI string
	// ContentType declares the serialized example media type for decoding; no sniffing is implied.
	ContentType string
	// Data contains borrowed serialized bytes; callers/providers must not mutate them during validation.
	Data []byte
}

// ExternalExampleResolver retrieves serialized example bytes through an
// application-owned authorization and transport policy.
type ExternalExampleResolver interface {
	// ResolveExternalExample retrieves an authorized example synchronously. The provider
	// owns transport policy, context cooperation, and safe lifetime of returned bytes.
	ResolveExternalExample(context.Context, string) (ExternalExampleResource, error)
}

// ExternalExampleResolverFunc adapts a function to ExternalExampleResolver.
type ExternalExampleResolverFunc func(
	context.Context,
	string,
) (ExternalExampleResource, error)

// ResolveExternalExample implements ExternalExampleResolver.
func (resolver ExternalExampleResolverFunc) ResolveExternalExample(
	ctx context.Context,
	identifier string,
) (ExternalExampleResource, error) {
	return resolver(ctx, identifier)
}

// MediaTypeExampleCodec serializes and parses one application-owned media
// representation. Implementations own format policy and perform no implicit
// I/O through the validator.
type MediaTypeExampleCodec interface {
	// Encode serializes an immutable value synchronously under application format policy.
	Encode(context.Context, jsonvalue.Value) ([]byte, error)
	// Decode parses borrowed bytes synchronously into immutable semantic data; it must
	// not retain mutable input aliases or perform unauthorized resource retrieval.
	Decode(context.Context, []byte) (jsonvalue.Value, error)
}

// MediaTypeExampleCodecResolver selects an explicitly configured codec for a
// Media Type Object. Returning nil leaves that representation unchecked.
type MediaTypeExampleCodecResolver interface {
	// ResolveMediaTypeExampleCodec selects a codec synchronously; nil skips that
	// representation. Providers own context cooperation and concurrency safety.
	ResolveMediaTypeExampleCodec(
		context.Context,
		string,
		jsonvalue.Value,
	) (MediaTypeExampleCodec, error)
}

// MediaTypeExampleCodecResolverFunc adapts a function to a codec resolver.
type MediaTypeExampleCodecResolverFunc func(
	context.Context,
	string,
	jsonvalue.Value,
) (MediaTypeExampleCodec, error)

// ResolveMediaTypeExampleCodec implements MediaTypeExampleCodecResolver.
func (resolver MediaTypeExampleCodecResolverFunc) ResolveMediaTypeExampleCodec(
	ctx context.Context,
	mediaType string,
	mediaTypeObject jsonvalue.Value,
) (MediaTypeExampleCodec, error) {
	return resolver(ctx, mediaType, mediaTypeObject)
}

// Options controls diagnostics and explicitly authorized external work.
// Callback and resolver instances are borrowed and never closed; callers own
// their concurrency safety, cancellation cooperation, and side effects.
type Options struct {
	// FailFast returns only the first diagnostic after validation passes finish; it
	// is not an early-work cutoff and does not prevent callbacks or retrieval.
	FailFast bool
	// MaxDiagnostics caps returned findings after collection; it must be positive
	// and is not a bound on intermediate diagnostic allocation.
	MaxDiagnostics int
	// MaxDocumentNodes limits semantic values visited in the input document.
	MaxDocumentNodes int
	// MaxDocumentDepth limits input document nesting, counting the root as one.
	MaxDocumentDepth int
	// SchemaResourceLoader explicitly authorizes non-pinned Schema Object
	// dialects and resources. Nil keeps external schema loading disabled.
	SchemaResourceLoader openapischema.ResourceLoader
	// ReferenceResourceURI supplies the retrieval URI used as the document base.
	ReferenceResourceURI string
	// ReferenceResolver explicitly authorizes external OpenAPI references. Nil
	// leaves them unresolved without making an otherwise valid document fail.
	ReferenceResolver reference.Resolver
	// ReferenceLimits bound reference scanning and target resolution.
	ReferenceLimits reference.Limits
	// MaxReferences limits resolved reference occurrences in one document.
	MaxReferences int
	// ExternalExampleResolver explicitly authorizes serialized example
	// retrieval. Nil keeps retrieval disabled.
	ExternalExampleResolver ExternalExampleResolver
	// MediaTypeExampleCodecResolver explicitly supplies application-owned
	// serialization for non-JSON media examples.
	MediaTypeExampleCodecResolver MediaTypeExampleCodecResolver
	// MaxExternalExampleBytes bounds each retrieved serialized example.
	MaxExternalExampleBytes int
	schemaCompilerFactory   func(openapi.Document, ...openapischema.Option) (*openapischema.Compiler, error)
	schemaMarshaller        func(jsonvalue.Value) ([]byte, error)
	schemaValidator         func(*openapischema.Compiler, context.Context, jsonvalue.Value) (openapischema.OutputUnit, error)
}

// DefaultOptions returns bounded collect-all validation policy.
func DefaultOptions() Options {
	return Options{
		MaxDiagnostics:          10_000,
		MaxDocumentNodes:        1_000_000,
		MaxDocumentDepth:        256,
		ReferenceLimits:         reference.DefaultLimits(),
		MaxReferences:           100_000,
		MaxExternalExampleBytes: 16 << 20,
	}
}

// Report is an immutable validation result.
type Report struct {
	diagnostics []Diagnostic
}

// Valid reports whether no error diagnostics were produced.
func (report Report) Valid() bool {
	for _, diagnostic := range report.diagnostics {
		if diagnostic.Severity == SeverityError {
			return false
		}
	}
	return true
}

// Diagnostics returns a caller-owned copy in deterministic evaluation order.
func (report Report) Diagnostics() []Diagnostic {
	return append([]Diagnostic(nil), report.diagnostics...)
}
