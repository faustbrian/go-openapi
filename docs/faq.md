# FAQ and troubleshooting

## Which package should an application import first?

Start with `github.com/faustbrian/go-openapi`. Its `ParseJSON` and `ParseYAML`
functions select the declared Swagger or OpenAPI version and return an immutable
document. Import a focused subpackage only for a specific operation such as
validation, reference resolution, composition, conversion, or serialization.
The root [package map](../README.md#package-map) and
[feature reference](reference.md) describe those boundaries.

## Does parsing fetch external references?

No. Parsing and validation perform no implicit filesystem or network access.
External references remain disabled until the caller supplies an explicitly
authorized resolver. File and HTTP resolvers deny all roots and destinations by
default and enforce caller-selected resource limits.

## Is this a server or code generator?

No. The module models and evaluates OpenAPI descriptions. It does not register
routes, generate handlers or clients, choose an HTTP stack, or decide whether a
wire-compatible API change is compatible with an application's business rules.

## How should errors be classified?

Use `errors.Is` with the sentinel categories exported by the package that owns
the failed operation. Do not classify failures by matching error strings.
Parser and resolver messages are bounded and redact untrusted content; wrapped
causes should be logged only when their caller-owned source is safe to disclose.
Context cancellation and deadlines remain distinguishable through the supplied
context.

## Which values can be reused concurrently?

Parsed semantic values and generated models are immutable. Validators and the
built-in resolvers document their concurrent-use guarantees. Caller-provided
callbacks and values behind caller-owned pointers retain their own
synchronization requirements. See the
[concurrency and ownership contract](security.md#concurrency-and-ownership).

## Which resources require cleanup?

Close each `reference.FileResolver` after its bounded operation to release its
opened roots. Call `CloseIdleConnections` on a `reference.HTTPResolver` when its
operation no longer needs retained keep-alive connections. The module starts no
background lifecycle, and callers continue to own supplied readers, writers,
callbacks, and transports.

## Does the module provide public testing helpers?

The module does not provide a separate public testing-helper package. Public
constructors accept ordinary readers, writers, contexts, callbacks, and
resolver interfaces, so applications can provide deterministic local fixtures
directly. The
[executable examples](../example_test.go) demonstrate supported entry points.
The `interoperability` module and specification corpora are maintainer evidence,
not production dependencies or a public testing API.

## Why is otherwise valid YAML rejected?

OpenAPI uses a JSON data model. This module accepts a single JSON-equivalent
YAML document and rejects aliases, merge keys, custom tags, non-string keys,
duplicate keys, and non-JSON scalars so producer and consumer semantics cannot
depend on a YAML implementation. See the
[specification decision](specification-decisions.md#openapi-dec-010-json-equivalent-yaml-input).

## Troubleshooting

### An external reference is disabled or denied

Supply a file or HTTP resolver only when external access is required, then grant
the narrowest roots, schemes, hosts, ports, and address ranges. The
[reference guide](reference.md#references) and
[network and filesystem policy](security.md#network-and-filesystem-policy)
describe the deny-by-default model.

### Parsing or traversal exceeds a limit

Start from the relevant `Default*` policy, measure the expected document, and
raise only the limiting dimension. Keep a caller deadline for untrusted input.
Do not remove bounds or use an unlimited process-wide resolver.

### A declared OpenAPI version is rejected

The module accepts only the exact versions listed in the
[supported specifications](../README.md#supported-specifications). Future
versions are rejected rather than interpreted as the latest known version.

### More help is needed

Use the [support process](../SUPPORT.md) for reproducible defects and adoption
questions. Report security vulnerabilities through the private process in
[SECURITY.md](../SECURITY.md).
