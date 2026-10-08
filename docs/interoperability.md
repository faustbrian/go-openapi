# Interoperability evidence

The decision-attributed review is
[`specification/interoperability.tsv`](../specification/interoperability.tsv).
The isolated harness owns the raw executable matrix in
[`interoperability/expected.tsv`](../interoperability/expected.tsv). Both record
observed behavior, not normative authority. The OpenAPI specifications and
their incorporated standards remain authoritative when implementations differ.

## Pinned implementations

| Implementation | Version | License | Version-specific source |
| --- | --- | --- | --- |
| `getkin/kin-openapi` | v0.149.0 | MIT | [source](https://github.com/getkin/kin-openapi/tree/v0.149.0), [license](https://github.com/getkin/kin-openapi/blob/v0.149.0/LICENSE) |
| `pb33f/libopenapi` | v0.41.2 | MIT | [source](https://github.com/pb33f/libopenapi/tree/v0.41.2), [license](https://github.com/pb33f/libopenapi/blob/v0.41.2/LICENSE) |
| `openapi/loads` | v0.25.3 | Apache-2.0 | [source](https://github.com/go-openapi/loads/tree/v0.25.3), [license](https://github.com/go-openapi/loads/blob/v0.25.3/LICENSE) |

The reviewed license-file SHA-256 values are, respectively,
`612a11e78c07e12765d9cf3866e3edc5f8212c541602721b14e465311afb3aa6`,
`7c119da422b6d9c6498f2df86d4fd6935f5b2b32d92df640ad8e9b60a8c3dac6`,
and `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`.

## Pinned public descriptions

| Description | Version | Revision | License |
| --- | --- | --- | --- |
| [Swagger Petstore](https://github.com/swagger-api/swagger-petstore/tree/8f0dd286987880b4af7bce552aca3813166f3049) | OpenAPI 3.0.4 | `8f0dd286987880b4af7bce552aca3813166f3049` | Apache-2.0 |
| [GitHub REST API](https://github.com/github/rest-api-description/tree/417c4fb368fc6a7162ce5f3eeeddce1a9a217747) | OpenAPI 3.1.0 | `417c4fb368fc6a7162ce5f3eeeddce1a9a217747` | MIT |

The description SHA-256 values are, respectively,
`0d810997f6409d5cff6f0cf2c1466814ba52250a784cd841cacb93514c7a8502`
and
`9d85f3a842c0215768f30f83ac7d1595430236fc51ce9c84e344b991a9f6b3da`.
Their license files are retained byte-for-byte with SHA-256 values
`b40930bbcf80744c86c46a12bc9da056641d722716c378f5659b9e555ef833e1`
and
`3243761cbac07e6d169a5a2f4e7c25cc544da85248e735df74c3672e055cc87b`.
The complete source paths, revisions, licenses, and update inputs are in
[`specification/manifest.json`](../specification/manifest.json).

The runner is compiled from the isolated module in
[`interoperability/go.mod`](../interoperability/go.mod). Its complete selected
graph and download checksums are committed, `go mod tidy -diff` and
`go mod verify` run before execution, and `go run -mod=readonly` prevents graph
updates. Exact direct versions are discovered again from Go build information
at runtime and emitted into the matrix. Candidate mode copies the pinned module to a temporary directory, adapts
the historical owned imports and dependencies to `/v2`, and reconciles shared
suppliers with the root graph. Its local root replacement qualifies only the
current candidate, not the separately maintained published-v1 harness or a
published-v2 release.
Competitor packages never enter `openapi`'s `go.mod`, production binary, or
core dependency graph. All tools run with external reference loading disabled
or with fixtures that contain no external reference.

## Fixture policy

The fixtures are small, locally authored audit cases. They are not imported
official artifacts and make no provenance claim beyond their reviewed source in
this repository.

- `oas30.json`, `oas31.json`, `oas32.json`, and `swagger20.json` exercise the
  minimal supported root shape independently for each version.
- `oas31-boolean-schema.json` exercises the normative OpenAPI 3.1 boolean
  Schema Object form.
- `oas32-tag-kind.json` exercises the OpenAPI 3.2 Tag Object `kind` field.
- `duplicate-key.json` and `yaml-alias.yaml` exercise deliberate strict parser
  policy. JSON duplicate members are ambiguous, and this package restricts
  YAML to one JSON-equivalent representation without aliases.

The public descriptions are a separate, unmodified interoperability corpus.
They test realistic size and ecosystem conventions; they never override
normative requirements. `openapi` uses a documented 20,000,000-operation
official-schema budget and 1,000,000-node/reference traversal budgets for this
reviewed corpus. Core defaults remain conservative.

`parse`, `model`, and `validate` are separate columns because accepting bytes,
constructing a typed model, and enforcing normative rules are distinct claims.
`roundtrip` reparses output rendered by the same implementation; it is recorded
only where that implementation exposes an applicable renderer. `na` means the
selected library has no applicable surface in this comparison and is never
treated as success.

## Current comparison on 2026-10-08

The candidate comparison also uses `go-openapi/loads` v0.25.3. Its ten
fixture rows retain their previous outcomes: the Swagger 2.0 row parses,
models, and round-trips; other rows remain `na` on this peer surface. A
separate non-empty Swagger case verifies operation, response, local-schema
reference, quoted-default, and extension data through load and reanalysis.

The candidate-v2 comparison uses `pb33f/libopenapi` v0.41.2. All ten fixture
rows retain their previous parse, model, validation, and round-trip outcomes;
only the peer version cells changed. Validation remains `na` on this peer
surface, and Swagger 2.0 rendering/reloading remains a recorded limitation.
The July findings below retain their original version attribution.

The peer now uses `pb33f/go-yaml` v0.1.1 (Apache-2.0) internally. Its pinned
[license](https://github.com/pb33f/go-yaml/blob/v0.1.1/LICENSE) has SHA-256
`4a4e095045f46875715234a884741acd2b61b88959d6f47a243ae03ef95c4e76`.
Upstream YAML node identities and low-level node APIs changed; this harness
uses document bytes, model construction, and rendering instead of those node
APIs. The peer migration does not change the root runtime dependency graph or
public parser API. Separate JSON/YAML semantic checks exercise retained
operation, response, schema-reference, boolean-schema, and quoted-string data.

## Kin 0.149.0 adoption

The current harness selects Kin 0.149.0 and its jsonschema/v6 6.0.3
validation dependency. The complete 0.147.0-to-0.149.0 interval includes
OpenAPI 3.1 format handling, rational `multipleOf` evaluation, positional
`prefixItems`, boolean schemas, and schema-cycle guards. These supplier
changes do not alter the root runtime graph or establish general specification
conformance beyond the exercised surfaces.

The refreshed matrix changes one outcome row: the valid OpenAPI 3.1 boolean
schema fixture now loads, validates, and round-trips successfully. Other
outcomes are unchanged, including the OpenAPI 3.2 tag-kind validation limit.
Separate JSON/YAML checks retain operation, response, local-reference,
required-property, extension, and quoted-string data through reload. Boolean
schema checks preserve both `true` and `false`; malformed input and external
references remain rejected by the selected runner boundary.

The dated findings below retain their original supplier versions. They are
historical observations, not descriptions of Kin 0.149.0.

## Findings on 2026-07-22

Observed facts:

- All applicable implementations accept and rebuild the basic OpenAPI 3.0,
  OpenAPI 3.1, OpenAPI 3.2, and Swagger 2.0 fixtures.
- `getkin/kin-openapi` v0.147.0 rejects the valid OpenAPI 3.1 boolean Schema
  Object during loading. This is classified as an implementation limitation,
  not evidence against the normative form.
- `getkin/kin-openapi` v0.147.0 models basic OpenAPI 3.2 but rejects the new Tag
  Object `kind` field during validation. This is classified as incomplete 3.2
  validation in the compared version.
- `pb33f/libopenapi` v0.38.7 parses and models every valid fixture, including
  both newer forms. Its compared surface does not perform a separate normative
  validation pass, so validation remains `na`.
- `pb33f/libopenapi` v0.38.7 cannot render and reload the Swagger 2.0 fixture
  through `Document.Render`. This is classified as a renderer-surface
  limitation; its Swagger model construction succeeds.
- The independent YAML loaders accept aliases and both independent JSON loaders
  accept a duplicate member in these fixtures. `openapi` rejects both by
  deliberate strict-input policy to prevent parser-dependent semantics.
- The pinned Swagger Petstore description parses, models, validates, and
  round-trips in `openapi` and `getkin/kin-openapi`; `pb33f/libopenapi`
  parses, models, and round-trips it on its applicable surface. `openapi`
  reports nine non-fatal XML interoperability warnings.
- The pinned GitHub REST API description parses, models, and round-trips in
  `openapi` and `pb33f/libopenapi`. `openapi` reports two ambiguous path
  templates and five discriminator targets whose property is not required,
  plus stable warning classes recorded by the tagged public-corpus test.
  `getkin/kin-openapi` also rejects validation and cannot complete its compared
  round trip. These outcomes are observations about the pinned description and
  tool versions, not authority for either interpretation.

These findings do not establish that any implementation is generally more
correct. The matrix tests matched documents and explicit surfaces only. It does
not compare implicit network loading, diagnostics, performance, code
generation, or unsupported feature synthesis.

## Update procedure

1. Change pinned tool versions in `interoperability/go.mod` and run `go mod tidy`
   from that directory to update the complete graph and `go.sum`.
2. For public descriptions, update `specification/manifest.json` and
   `scripts/sync-spec.sh`, then run the sync script and verify exact checksums.
3. Verify each version-specific source, license, release notes, and applicable
   API behavior; update this document and license hashes.
4. Run `INTEROP_UPDATE=true make -f verification/package.mk interoperability`
   to verify the pinned graph
   and regenerate the observed matrix in a temporary module.
5. Review every changed cell and diagnostic class against normative text and
   classify it as a true
   incompatibility, an implementation defect or limitation, or deliberate
   package policy.
6. Run `make -f verification/package.mk interoperability` again without update
   mode. The checked-in matrix
   must match exactly.

The scheduled CI job performs the final command. Version or behavior drift is a
failure requiring review; the expected matrix must never be loosened merely to
accept a new independent result.
