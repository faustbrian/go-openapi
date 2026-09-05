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

The module is active, stable, and requires Go 1.26.6 or newer. It owns OpenAPI
document semantics; applications continue to own HTTP serving, routing,
business compatibility policy, and any authority granted to external resource
resolvers.

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
package main

import (
	"context"
	"fmt"
	"strings"

	openapi "github.com/faustbrian/go-openapi"
	"github.com/faustbrian/go-openapi/parse"
	"github.com/faustbrian/go-openapi/validate"
)

func main() {
	document, err := openapi.ParseYAML(
		context.Background(),
		strings.NewReader(`
openapi: 3.2.0
info:
  title: Example
  version: "1"
paths: {}
`),
		parse.DefaultLimits(),
	)
	if err != nil {
		panic(err)
	}

	report, err := validate.Document(context.Background(), document)
	if err != nil {
		panic(err)
	}

	fmt.Println(document.SpecificationVersion(), report.Valid())
}
```

This prints `3.2.0 true`. The tested
[package examples](example_test.go) also demonstrate canonical serialization
and document merging.

Strict YAML accepts only JSON-compatible scalar tags and string mapping keys.
Strict JSON rejects duplicate members and unpaired UTF-16 surrogate escapes.
Reference loading, external access, and resource budgets are caller-owned and
explicit.

## Package map

- `openapi`, `parse`, `specversion`, `jsonvalue`, `model`, `swagger20`, and
  `oas30` through `oas32` select and expose immutable document models.
- `validate`, `jsonschema`, `security`, `expression`, and `discriminator` apply
  document, Schema Object, security-requirement, runtime-expression, and schema
  selection semantics.
- `reference` and `implicit` resolve explicitly authorized references and
  connections under caller-selected bounds.
- `media`, `parameter`, `response`, `server`, and `xmlvalue` implement focused
  protocol value handling without owning an HTTP server.
- `compose`, `convert`, `diff`, and `serialize` transform, compare, and emit
  documents with explicit loss, conflict, and resource policies.
- `specification` exposes the embedded, provenance-checked specification
  resources used by the module.

The [API reference](https://pkg.go.dev/github.com/faustbrian/go-openapi) lists
every public package and exported identifier.

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

## Construction, errors, and lifecycle

Plain functions handle stateless value operations. Constructors validate
coherent options before returning immutable values or concurrency-safe
resolvers and validators; named `Default*` functions expose finite defaults.
Operations that may read, resolve, validate, or traverse caller-controlled data
accept a context and honor cancellation. Stable error categories support
`errors.Is`; the package does not retry automatically.

Parsed values own their immutable data. Callers own supplied readers, writers,
transports, callbacks, and resolver authority. A file resolver owns its opened
roots and must be closed; an HTTP resolver retains only idle connections, which
the caller releases with `CloseIdleConnections`. The package starts no
background lifecycle. See the [feature reference](docs/reference.md) and
[security and ownership model](docs/security.md) for the complete contracts.

## Documentation

Start with the [documentation index](docs/README.md). It separates usage,
specification behavior, security, performance, interoperability, and
maintainer references. The [specification decision register](docs/specification-decisions.md)
records intentional behavior where specification text or ecosystem practice
requires an explicit choice.
Shared construction, ownership, lifecycle, and composition expectations are in
the versioned [Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/README.md)
and its [Protocols and descriptions family](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/design-language.md#package-families-and-selection).

For adoption and maintenance, see the [executable examples](example_test.go),
[FAQ and troubleshooting](docs/faq.md), [compatibility policy](COMPATIBILITY.md),
[deprecation and migration policy](DEPRECATION.md),
[performance evidence](docs/performance.md), [support policy](SUPPORT.md),
[security guidance](docs/security.md), and
[private vulnerability-reporting process](SECURITY.md). Release changes are in
the [changelog](CHANGELOG.md).

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
