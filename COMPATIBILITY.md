# Compatibility Policy

Each releasable directory is an independent Go module and follows semantic
versioning. The root module uses `v<version>` tags. A releasable nested module
uses `<module-directory>/v<version>` tags. The current `interoperability`
module is an internal, non-releasable harness and has no release tag.

The next root release is v2 and uses `github.com/faustbrian/go-openapi/v2`.
Move root and subpackage imports to that path and use Go 1.27.0 or newer.
There are no version-specific source directories or branches. The nominal
migration preserves exported declaration names, but adopts public JSON Schema
v2.0.0. Use `github.com/faustbrian/go-json-schema/v2` for `Schema`, `Result`,
`OutputUnit`, resource loaders, keyword compilers and the `Limits` passed to
`validate.NewValidatorWithDocumentSchemaLimits`; v1 named types are not
interchangeable with these v2 types. Standard Unicode property patterns such
as `\p{Script=Greek}` are supported, while nonstandard patterns such as
`\p{Greek}` are rejected. Review stored expressions before migrating.
Explicit canonical format assertions also reject invalid combined URI-template
prefix/explode modifiers such as `{var:1*}`. OpenAPI does not implicitly enable
format assertions. Existing v1 consumers can remain on the published v1
module with Go 1.26.6.

OpenAPI v2 is not published yet. The tracked non-releasable interoperability
module retains its public v1 dependencies until publication; disposable
candidate composition uses the planned OpenAPI `/v2` with public JSON Schema
`/v2`. Candidate checks do not prove public OpenAPI v2 availability or adoption.

Before `v1`, minor releases MAY contain reviewed breaking changes, but every
break MUST be documented with migration guidance. Patch releases MUST remain
backward compatible. At and after `v1`, incompatible exported API or documented
behavior changes require a new major version.

Compatibility includes exported Go APIs, error classification, serialization,
protocol behavior, persistence schemas, environment variables, command output,
resource ownership, ordering, retry/idempotency semantics, and documented
defaults. A compile-compatible change can still be behaviorally breaking.

Specification-backed modules MUST NOT diverge from their declared standards.
Ambiguities require documented decisions and stable tests. Deprecated APIs
follow [`DEPRECATION.md`](DEPRECATION.md).
The [specification decision register](docs/specification-decisions.md) is the
compatibility authority for observable specification interpretations.
