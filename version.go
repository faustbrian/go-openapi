// Package openapi parses JSON or restricted YAML into immutable typed views of
// pinned Swagger 2.0 and OpenAPI 3.0, 3.1, and 3.2 revisions. ParseJSON and
// ParseYAML select the version; Decode accepts an existing JSON semantic value.
// Unknown fields, exact number spelling, and member order survive decoding,
// but source whitespace, comments, and scalar styles do not.
//
// Decoding checks the root version marker, not document validity. Use validate
// for diagnostics, jsonschema for instance plans, and reference for explicitly
// authorized retrieval. Readers remain caller-owned. Cancellation is cooperative
// and cannot interrupt a caller's blocking Read. Returned built-in models share
// immutable values and are safe to read concurrently without a shutdown step.
package openapi

import "github.com/faustbrian/go-openapi/v2/specversion"

// ErrMalformedVersion reports a version string that is not syntactically an
// OpenAPI or Swagger version.
var ErrMalformedVersion = specversion.ErrMalformedVersion

// ErrUnsupportedVersion reports a syntactically valid version whose
// specification revision is not pinned and implemented by this package.
var ErrUnsupportedVersion = specversion.ErrUnsupportedVersion

// Dialect identifies a version line with distinct model and validation
// semantics.
type Dialect = specversion.Dialect

const (
	// DialectSwagger20 selects the pinned Swagger 2.0 model rules.
	DialectSwagger20 = specversion.DialectSwagger20
	// DialectOAS30 selects the pinned OpenAPI 3.0 model rules.
	DialectOAS30 = specversion.DialectOAS30
	// DialectOAS31 selects the pinned OpenAPI 3.1 model rules.
	DialectOAS31 = specversion.DialectOAS31
	// DialectOAS32 selects the pinned OpenAPI 3.2 model rules.
	DialectOAS32 = specversion.DialectOAS32
)

// Version preserves the exact supported patch revision selected by a
// document or caller. Its zero value is invalid.
type Version = specversion.Version

// ParseVersion parses an exact supported OpenAPI or Swagger version. Future
// patch releases are rejected until their authoritative specification inputs
// have been pinned and reviewed.
func ParseVersion(value string) (Version, error) {
	return specversion.Parse(value)
}
