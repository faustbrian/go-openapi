# Security and trust-boundary model

This document describes the implemented security model reviewed by the
2026-07-22 release audit. It is a bounded engineering claim, not a guarantee
against every future input, dependency defect, or caller misconfiguration.

The labels below keep different kinds of evidence separate:

- **Specification requirement** is imposed by OpenAPI or an incorporated
  standard.
- **Observed fact** is directly established by current source and executable
  tests.
- **Package policy** is a deliberate restriction or guarantee chosen by this
  module.
- **Inference** is a security conclusion drawn from the preceding evidence.

## Protected properties

The protected properties are process availability; bounded CPU, memory,
goroutine, file, socket, and output use; the integrity of parsed semantics and
pinned conformance artifacts; filesystem and network authority; confidentiality
of paths, credentials, document contents, and resolver details; deterministic
results; and the absence of mutable cross-call aliasing.

Attackers may control JSON, YAML, Schema Objects, references, runtime
expressions, parameter values, media payloads, public API descriptions, reader
and writer failures, remote HTTP responses, and files beneath a caller-granted
root. Callbacks and explicitly supplied resolvers are caller code and therefore
belong to the caller's trust domain.

## Trust boundaries

| Boundary | Authority granted | Implemented control | Evidence |
| --- | --- | --- | --- |
| JSON and YAML readers | Bytes and reader behavior only | Strict single-document parsing, duplicate-key and invalid UTF-8 rejection, exact numbers, JSON-equivalent YAML, configurable independent limits, and context-aware reads | `parse` tests, fuzzers, race tests, and mutation gate |
| Semantic values and models | Caller-owned immutable data | Constructor copies, caller-owned returned slices, ordered exact semantics, and bounded high-level traversal | `jsonvalue`, `model`, generated model, composition, diff, validation, and serialization tests |
| File resolver | Only explicitly listed existing directories | Canonical roots opened with `os.OpenRoot`, post-symlink containment, accepted extensions, regular-file checks, cumulative document and byte limits, close lifecycle, and redacted access failures | `reference/file_resolver*_test.go` and resolver fuzzing |
| HTTP resolver | Exact schemes, hosts, ports, and optional CIDRs | No environment proxy, no transparent decompression, DNS address approval and pinned dialing, special-purpose address denial, redirect revalidation, credential stripping, bounded headers, bodies, documents, redirects, addresses, concurrency, and duration | `reference/http_resolver*_test.go` and resolver fuzzing |
| External references | Only the caller's explicit `Resolver` | URI and JSON Pointer validation, bounded graph traversal, cycle termination, cache identity, cancellation, and redacted resolver failures | `reference` tests, fuzzers, race tests, and mutation gate |
| JSON Schema compiler | Only explicitly configured dialect and resource loaders | Dialect separation, bounded traversal, single-flight construction, cancellable waiters, and redacted loader and compilation failures | `jsonschema` conformance, concurrency, fuzz, race, and mutation evidence |
| Writers and generated evidence | Caller-supplied writer or repository-local destination | Byte, node, and depth limits; deterministic ordering; atomic temporary-file replacement; bounded, single-value internal decoders | `serialize` and internal command tests and fuzzers |
| Official artifacts | Repository checkout only | Pinned source revision or retrieval date, SPDX license, HTTPS license source, SHA-256, regular-file and symlink checks, and offline verification | `specification/manifest.json` and `make -f verification/package.mk conformance` |

## Parser policy

**Specification requirement:** OpenAPI descriptions use JSON or YAML with the
JSON data model and string mapping keys. Version-specific OpenAPI and JSON
Schema rules are applied only after syntax parsing.

**Package policy:** Ambiguous syntax is rejected. This includes concatenated or
multi-document input, duplicate keys, non-string YAML keys, anchors, aliases,
merge keys, custom tags, non-JSON number forms, invalid UTF-8, and unpaired
JSON surrogate escapes. YAML is not a general-purpose YAML loader.

**Observed fact:** Input bytes are read through the caller context under
`parse.Limits`. JSON checks cancellation during token traversal. YAML checks it
while the underlying decoder consumes syntax and at both document boundaries,
then during semantic traversal. Parser diagnostics do not include source text
or member names.

**Inference:** YAML alias expansion, parser differentials, and diagnostic
reflection cannot be used to smuggle alternative OpenAPI semantics through the
document parser. The caller still needs a deadline appropriate to its service.

## Network and filesystem policy

**Package policy:** Core parsing and validation never load a remote or local
resource implicitly. A nil resolver means no external retrieval. File and HTTP
resolvers grant no authority with their default options until the caller adds
roots or hosts.

**Observed fact:** HTTP identifiers reject user information, query strings, and
fragments. Redirects are re-authorized and lose `Authorization`, `Cookie`, and
`Proxy-Authorization`. Environment proxies and transparent decompression are
disabled. Responses with unsupported explicit media types or non-identity
content encodings are rejected. File identifiers must be absolute `file` URIs
inside a root after symlink evaluation.

**Inference:** A description cannot independently turn a `$ref` into ambient
SSRF or filesystem access. Explicitly allowing a private CIDR or broad root is
a transfer of authority by the caller and must be treated like any other
security-sensitive configuration.

The supplied HTTP policy makes an authorization decision on resolved addresses
and dials the approved address. It does not claim that an explicitly allowed
remote service is honest. Content authenticity, TLS trust-root policy, service
identity beyond normal HTTPS verification, and application credentials remain
caller responsibilities.

## Resource and lifecycle policy

### Repository automation authority

The internal provenance, specification-matrix, model-generation and dependency
audit commands run only when a developer explicitly selects a repository root.
The mutation-report command likewise reads an explicit operator-selected report,
which may be outside that checkout. These are maintainer tools, not public API
resolvers or sandboxes for hostile checkouts. Repository roots, manifest path
metadata and output locations belong to the invoking developer's trust domain;
parsed API payloads do not select those paths. Run them only in a reviewed,
stable checkout, without an untrusted process concurrently replacing files.

Manifest, field-inventory, dependency-evidence and mutation-report decoding
retains its existing byte limits. Provenance additionally rejects non-local
artifact paths, symlink components and non-regular files before checking pinned
digests. Generated specification inventories and Go source are public artifacts,
so their directories deliberately use mode 0755, subject to the process umask;
these output locations must not contain credentials or private payloads.

The maintainer owns this accepted tooling-authority boundary. The rationale is
that explicitly invoked code-generation and verification tools need checkout
access; narrowing the standalone report reader to that root would break valid
temporary-report use. Reviewed metadata, stable checkouts, bounded decoders and
provenance checks mitigate misuse. Reassess the narrow scanner annotations if
these commands become service entry points, receive untrusted path metadata,
operate in concurrently hostile directories or generate confidential output.

**Observed fact:** Independent limits exist for input and output bytes, syntax
tokens, semantic values, scalar size, object and array width, depth, references,
documents, redirects, addresses, concurrency, diagnostics, and operation-
specific graph work. Limits are checked before wide semantic copies on audited
paths. File handles and HTTP bodies are closed, temporary output is replaced
atomically, and resolver document budgets are cumulative.

**Package policy:** Zero-valued option fields use documented defaults only where
the API explicitly says so. Invalid negative or overflow-prone limits are
rejected. Callers should lower defaults to their expected document size and use
deadlines for every untrusted operation.

**Inference:** The implemented guards bound known amplification axes; aggregate
process safety still depends on callers bounding the number of simultaneous
top-level operations and the lifetime of objects they retain.

`security.Satisfied` checks array, requirement-object, and scope-array lengths
before copying their collections. Scheme and scope counts apply independently
to requirements and supplied credentials, with zero fields selecting defaults
and positive overrides retaining caller-selected finite limits. Requirement
scheme names and required scopes share an inclusive fixed 1 MiB raw-byte
budget; credential scheme names and granted scopes share a separate 1 MiB
budget. Every occurrence counts before deduplication, using remaining-budget
subtraction rather than overflow-prone addition. All requirements are admitted
and validated before any successful result; valid anonymous alternatives skip
credential admission because those credentials cannot affect the result.
Other evaluations admit all credentials before building one call-local index
reused across alternatives. Limit errors disclose no labels. Callers must
inventory formerly accepted oversized sets/labels before adopting this v2
behavior, and must not mutate borrowed credential maps/slices during the
synchronous call. Allocation of caller-owned semantic values/credentials and
the number of concurrent calls remain application responsibilities.

## Historical v1 mitigation

Public v1.0.0 does not bound credential scheme/scope counts or scheme/scope
label bytes, and can rebuild granted-scope indexes across alternatives. Its
requirement limits alone do not bound that credential work. This matters when
applications expose attacker-influenced requirements or granted credentials
to `security.Satisfied` without equivalent upstream admission. No credential
authentication bypass or measured process exhaustion is established.

Upgrade to `github.com/faustbrian/go-openapi/v2` v2.0.0 or newer for package-owned
admission. No fixed version of the historical v1 module is published for this
issue. While migrating, applications retaining v1 must admit both input domains
before calling the evaluator: bound alternative, scheme and scope occurrence
counts, and total raw bytes of all scheme/scope labels. Count duplicate and
unused entries, not only matching entries; use finite service-appropriate
ceilings and reject rather than truncate inputs. Check byte allowances using
remaining-budget subtraction to avoid integer overflow. Bound upstream input
construction and concurrent calls separately.

Validate the shape of every requirement alternative and scope before allowing
any successful result. Do not rely on an earlier anonymous or matching
alternative to validate later malformed values. These are caller-owned
mitigations, not evidence that v1 has acquired the v2 guarantees. The maintainer
owns report triage and migration guidance; applications own any temporary
upstream admission. Review the workaround whenever accepted input shapes,
resource ceilings or the selected evaluator version change.

## Concurrency and ownership

**Observed fact:** Parsed semantic values and generated models are immutable.
Collection accessors return outer slices owned by the caller. Resolver counters,
compiler construction, and validation caches are synchronized. Callbacks are
not deliberately invoked while package locks are held. Race tests exercise all
production packages.

**Package policy:** Immutability of a generic collection does not recursively
freeze an arbitrary caller-defined element type. Callers must not mutate data
behind their own pointers concurrently unless that type documents safety.

## Content and secret handling

**Package policy:** Descriptions, examples, defaults, extensions, Markdown, and
other free-form strings are data. The module preserves them; it does not mark
them safe HTML, shell, source code, filesystem paths, or templates. Downstream
renderers and generators must apply contextual escaping.

**Observed fact:** The module emits no telemetry. Resolver access errors redact
paths, URLs, host details, anchor text, and underlying transport messages while
preserving stable error classification. Parser `Error()` strings are bounded
and content-free. A caller can deliberately inspect wrapped reader errors with
`errors.Unwrap`; those causes belong to the caller's reader trust boundary and
should not be logged blindly.

## Supply chain and update procedure

The historical dependency review, including graph-only modules, is recorded
in [`dependencies.tsv`](dependencies.tsv); current selected versions and
checksums are authoritative in `go.mod` and `go.sum`. Follow the current
[dependency update procedure](dependencies.md), verify the resolved graph, and
inspect the selected shared-tooling security results for vulnerability,
secret, license and SBOM checks. Official evidence updates follow
[`specification/README.md`](../specification/README.md) and must preserve exact
source bytes.

Security-sensitive updates must add a failing regression, retain explicit
limits and error classification, run the relevant fuzz and race targets, and
pass mutation testing before the release gates.

## Non-goals and residual responsibilities

The package does not authenticate API traffic, execute runtime expressions,
sanitize rich text, manage application secrets, establish a sandbox for caller
callbacks, or make an explicitly trusted remote server safe. Optional resolver
authority is disabled until configured.

Interoperability observations, fuzz campaigns, performance budgets, and
cross-platform tests remain separate evidence under
[interoperability](interoperability.md), [performance](performance.md), and
`specification/conformance/`. None may be inferred from this threat model,
aggregate coverage, official-schema success, or self round trips.
