# openapi

[![CI](https://github.com/faustbrian/go-openapi/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-openapi/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-openapi/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-openapi.svg)](https://pkg.go.dev/github.com/faustbrian/go-openapi)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-openapi?sort=semver)](https://github.com/faustbrian/go-openapi/releases)
[![Go](https://img.shields.io/badge/go-1.26.6-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`openapi` is a version-aware OpenAPI document toolkit for Go. It provides
immutable lossless semantic values, bounded JSON and YAML parsing, generated
typed views, normative validation, explicit reference resolution, composition,
conversion, compatibility diffing, and deterministic serialization.

## Supported specifications

- Swagger 2.0
- OpenAPI 3.0.0 through 3.0.4
- OpenAPI 3.1.0 through 3.1.2
- OpenAPI 3.2.0

Future versions are rejected rather than interpreted as the latest known
version. Unknown fields and `x-` extensions remain available through the
lossless semantic representation.

## Installation

```sh
go get github.com/faustbrian/go-openapi
```

## Quick start

Parsing requires an explicit context and finite limits. It performs no implicit
network or filesystem access.

```go
limits := parse.DefaultLimits()
document, err := openapi.ParseYAML(ctx, reader, limits)
if err != nil {
	return err
}

validator := validate.NewValidator()
report, err := validator.Document(ctx, document)
if err != nil {
	return err
}

for _, diagnostic := range report.Diagnostics() {
	fmt.Printf("%s %s: %s\n",
		diagnostic.Severity,
		diagnostic.InstanceLocation,
		diagnostic.Code,
	)
}
```

Strict YAML accepts only JSON-compatible scalar tags and string mapping keys.
Strict JSON rejects duplicate members and unpaired UTF-16 surrogate escapes.
Reference loading, external access, and resource budgets are caller-owned and
explicit.

## Capabilities

- Version-specific typed views preserve absent, null, valid, and invalid
  states without discarding the original document.
- Validation covers document structure, semantic rules, Schema Objects,
  request and response direction, media encodings, security requirements, and
  runtime expressions.
- Reference resolution supports internal and explicitly authorized external
  resources under depth, node, URI, and byte limits.
- Serialization covers parameters, headers, cookies, form data, multipart
  bodies, server variables, and canonical JSON output.
- Composition supports bundling, overlays, filtering, merging, and explicit
  version conversion with loss reporting.
- Compatibility diffing classifies operation, schema, parameter, response,
  security, callback, webhook, and component changes.

The package does not fetch references implicitly, choose an HTTP client,
generate application handlers, or infer business compatibility from syntax
alone.

## Documentation

Start with the [documentation index](docs/README.md). It separates usage,
specification behavior, security, performance, interoperability, and
maintainer references. The [specification decision register](docs/specification-decisions.md)
records intentional behavior where specification text or ecosystem practice
requires an explicit choice.
Shared construction, ownership, lifecycle, and composition expectations are in
the versioned [Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/README.md)
and its [Protocols and descriptions family](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/design-language.md#package-families-and-selection).

## Development

`make check` runs the repository contract. Conformance evidence lives under
`specification/conformance/`; see [CONTRIBUTING.md](CONTRIBUTING.md) before
changing generated views, normative requirements, or compatibility behavior.

## Related packages

- [JSON Schema](https://github.com/faustbrian/go-json-schema) provides the
  Schema Object dialect engine used by OpenAPI 3.1 and 3.2.
- [JSON-RPC](https://github.com/faustbrian/go-jsonrpc) and
  [OpenRPC](https://github.com/faustbrian/go-openrpc) cover RPC runtime and
  description contracts.

## License

MIT. See [LICENSE](LICENSE).
