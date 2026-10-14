# Package map

[Documentation index](README.md) · [Repository entry](../README.md)

This map follows the current source modules, not a synchronized release.
Use the repository README for published installation versions; an unpublished
source major is not installable merely because its package appears here.
Internal packages and non-releasable fixtures are intentionally excluded.

## github.com/faustbrian/go-openapi/v2

`openapi` is a version-aware OpenAPI document toolkit for Go. It provides immutable lossless semantic values, bounded JSON and YAML parsing, generated typed views for Swagger 2.0 and OpenAPI 3.0 through 3.2, normative validation, OpenAPI JSON Schema dialects, explicit reference resolution, composition, conversion, compatibility diffing, and bounded deterministic serialization.

| Import path | Responsibility |
| --- | --- |
| [`github.com/faustbrian/go-openapi/v2`](../document.go) | Use openapi for version-aware document parsing and as the entry point to the package's bounded modeling, validation, reference, composition, conversion, and serialization surfaces. |
| [`github.com/faustbrian/go-openapi/v2/compose`](../compose/filter.go) | Public compose package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/convert`](../convert/convert.go) | Public convert package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/diff`](../diff/operations.go) | Public diff package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/discriminator`](../discriminator/select.go) | Public discriminator package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/expression`](../expression/evaluate.go) | Public expression package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/implicit`](../implicit/resolve.go) | Public implicit package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/jsonschema`](../jsonschema/annotated_enum.go) | Public jsonschema package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/jsonvalue`](../jsonvalue/value.go) | Public JSON value package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/media`](../media/encoding.go) | Public media package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/model`](../model/collection.go) | Public model package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/oas30`](../oas30/model_generated.go) | Public oas30 package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/oas31`](../oas31/model_generated.go) | Public oas31 package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/oas32`](../oas32/model_generated.go) | Public oas32 package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/parameter`](../parameter/decode.go) | Public parameter package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/parse`](../parse/json.go) | Public parse package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/reference`](../reference/bundle.go) | Public reference package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/response`](../response/select.go) | Public response package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/security`](../security/evaluate.go) | Public security package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/serialize`](../serialize/json.go) | Public serialize package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/server`](../server/expand.go) | Public server package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/specification`](../specification/resources.go) | Public specification package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/specversion`](../specversion/version.go) | Public specification version package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/swagger20`](../swagger20/model_generated.go) | Public swagger20 package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/validate`](../validate/components.go) | Public validate package; see the linked source and module guide. |
| [`github.com/faustbrian/go-openapi/v2/xmlvalue`](../xmlvalue/null.go) | Public XML value package; see the linked source and module guide. |

Choose the next guide in the [documentation index](README.md)
or return to the [repository entry](../README.md) for installation and ownership.
