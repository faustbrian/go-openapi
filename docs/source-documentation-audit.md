# Source documentation audit

## Pre-edit owner, scope and acceptance

Base `5e997b34a59739be97075bd3335edcda0fca5aa3`, clean detached main-source
worktree. This worker owns the full standalone OpenAPI source-doc pass; the
original unrelated checkout and existing reference worktree remain untouched.
Parent owns shared enforcement, publication and aggregate evidence. Ordinary
contracts/comments are Tier A; actual producer/output/directive trust changes
are Tier C and require independent complete-final-diff review before commit.
No runtime, dependency, fixture, normative-policy/history or release changes.

Outcome: accurate version selection, lossless values versus typed validation,
parser/resolver admission and authority, IO/callback/cancellation/ownership,
model presence/unknown fields, transformation loss and diagnostics/resource
contracts across root and interoperability modules. Direct consumers include
versioned model views, validation/JSON Schema, reference operations, media and
parameter/server/security helpers, the pinned independent interoperability module
and published JSON Schema v2 contracts. Inspect actual imported versions without
substituting sibling working directories.

Named risks: unsupported full-conformance claims, parse-versus-validation/authorization
confusion, absent/null/invalid fields, reference cycles/URI/network/file policy,
copy/borrow/immutability and shared dialect/cache concurrency, partial results and
writer failures, source-label/private-data diagnostic leakage, output/input/work
budgets and source-provenance/generation trust. Smallest checks: all-source AST/
signature/literal/token/directive identity, native GoDoc/package/API recipe and
existing applicable Examples, real producer regeneration if output comments change,
offline specification/local links/source checks/pinned changed+staged secrets and
one ordinary whole-diff truth/usefulness review. No race/fuzz/public-description/
service/performance/release campaigns for unchanged executable contracts. Every
Go/native compile command requires explicit parent slot and owned disposable
GOCACHE/GOMODCACHE/GOTMPDIR with success/failure cleanup.

## Reproducible full inventory

`git ls-files "*.go"`, modules.json/packages.json and native source-docs
inventory/check identify 251 tracked Go files, two modules, 35 source package
entries and eight internally generated model/model-test files. Initial structural
input: 308 exported comment omissions, six package-comment omissions and eight
manual generator reviews from the retained development checker. This is an
investigation input, not a prose-quality score. No existing source-audit file.

- `.`: `benchmark_test.go`, `document.go`, `document_internal_test.go`, `document_test.go`, `documentation_test.go`, `example_test.go`, `generate.go`, `interoperability_test.go`, `parse.go`, `public_interoperability_test.go`, `version.go`, `version_test.go`.
- `compose`: `filter.go`, `filter_internal_test.go`, `filter_test.go`, `fuzz_test.go`, `merge.go`, `merge_internal_test.go`, `merge_test.go`.
- `convert`: `convert.go`, `convert_test.go`, `coverage_internal_test.go`, `fuzz_test.go`, `oas30.go`, `oas31.go`, `oas31_internal_test.go`, `oas32.go`, `oas32_internal_test.go`, `swagger.go`, `swagger20_downgrade.go`, `swagger20_downgrade_internal_test.go`, `swagger_internal_test.go`.
- `diff`: `fuzz_test.go`, `operations.go`, `operations_internal_test.go`, `operations_test.go`.
- `discriminator`: `select.go`, `select_test.go`.
- `expression`: `evaluate.go`, `evaluate_internal_test.go`, `evaluate_test.go`, `expression.go`, `expression_test.go`, `fuzz_test.go`.
- `implicit`: `resolve.go`, `resolve_test.go`.
- `internal/modelaccess`: `access.go`, `access_test.go`.
- `internal/modelgen`: `generate.go`, `generate_test.go`.
- `internal/modelgen/cmd/modelgen`: `main.go`, `main_test.go`.
- `internal/quality/cmd/dependencyaudit`: `main.go`, `main_test.go`.
- `internal/quality/cmd/mutationgate`: `main.go`, `main_test.go`.
- `internal/specification`: `evidence_test.go`, `matrix.go`, `matrix_test.go`, `objects.go`, `objects_test.go`.
- `internal/specification/cmd/provenance`: `main.go`, `main_test.go`.
- `internal/specification/cmd/specmatrix`: `main.go`, `main_test.go`.
- `interoperability`: `kin_semantic_test.go`, `loads_semantic_test.go`, `runner.go`, `runner_test.go`, `semantic_roundtrip_test.go`.
- `jsonschema`: `annotated_enum.go`, `annotated_enum_test.go`, `compiler.go`, `compiler_internal_test.go`, `compiler_test.go`, `dialect_advisory.go`, `dialect_advisory_test.go`, `fuzz_test.go`, `public_dependency_test.go`.
- `jsonvalue`: `value.go`, `value_fuzz_test.go`, `value_internal_test.go`, `value_test.go`.
- `media`: `encoding.go`, `encoding_internal_test.go`, `encoding_test.go`, `fuzz_test.go`, `linkset.go`, `linkset_test.go`, `multipart.go`, `multipart_test.go`, `serialization_type.go`, `serialization_type_test.go`, `sse.go`, `sse_test.go`.
- `model`: `collection.go`, `collection_test.go`, `field.go`, `field_test.go`.
- `oas30`: `model_generated.go`, `model_generated_test.go`.
- `oas31`: `model_generated.go`, `model_generated_test.go`.
- `oas32`: `model_generated.go`, `model_generated_test.go`.
- `parameter`: `decode.go`, `decode_errors_test.go`, `decode_internal_test.go`, `decode_test.go`, `encode.go`, `encode_errors_test.go`, `encode_internal_test.go`, `encode_test.go`, `escape.go`, `escape_test.go`, `fuzz_test.go`, `header_strategy.go`, `header_strategy_internal_test.go`, `header_strategy_test.go`, `limits_test.go`, `options.go`, `options_internal_test.go`, `options_test.go`, `swagger20.go`, `swagger20_test.go`.
- `parse`: `boundary_test.go`, `fuzz_test.go`, `internal_test.go`, `json.go`, `json_test.go`, `yaml.go`, `yaml_test.go`.
- `reference`: `bundle.go`, `bundle_internal_test.go`, `bundle_test.go`, `chain.go`, `chain_internal_test.go`, `dereference.go`, `dereference_internal_test.go`, `dereference_test.go`, `errors_test.go`, `file_resolver.go`, `file_resolver_internal_test.go`, `file_resolver_test.go`, `fuzz_test.go`, `graph.go`, `graph_test.go`, `http_resolver.go`, `http_resolver_internal_test.go`, `http_resolver_test.go`, `object.go`, `object_test.go`, `pointer.go`, `pointer_internal_test.go`, `pointer_test.go`, `resolve.go`, `resolve_internal_test.go`, `resolve_test.go`, `resolver_internal_fuzz_test.go`, `resource_amplification_test.go`, `scan.go`, `scan_test.go`.
- `response`: `select.go`, `select_test.go`.
- `security`: `evaluate.go`, `evaluate_admission_test.go`, `evaluate_budget_test.go`, `evaluate_label_admission_test.go`, `evaluate_test.go`.
- `serialize`: `boundary_test.go`, `fuzz_test.go`, `internal_test.go`, `json.go`, `json_test.go`, `yaml.go`, `yaml_test.go`.
- `server`: `expand.go`, `expand_test.go`, `resolve.go`, `resolve_test.go`.
- `specification`: `resources.go`, `resources_test.go`.
- `specversion`: `version.go`, `version_test.go`.
- `swagger20`: `model_generated.go`, `model_generated_test.go`.
- `validate`: `cache_internal_test.go`, `components.go`, `components_test.go`, `coverage_contracts_internal_test.go`, `deprecated.go`, `deprecated_test.go`, `document.go`, `document_test.go`, `examples.go`, `examples_internal_test.go`, `examples_test.go`, `external_docs.go`, `external_docs_test.go`, `fuzz_test.go`, `headers.go`, `headers_test.go`, `helpers_internal_test.go`, `http_authentication_schemes_test.go`, `http_status_codes_test.go`, `instance.go`, `instance_internal_test.go`, `instance_test.go`, `links.go`, `links_internal_test.go`, `links_test.go`, `media_types.go`, `media_types_internal_test.go`, `media_types_test.go`, `metadata.go`, `metadata_test.go`, `operations.go`, `operations_internal_test.go`, `parameters.go`, `parameters_test.go`, `paths.go`, `paths_internal_test.go`, `paths_test.go`, `prose_authority_internal_test.go`, `references.go`, `references_internal_test.go`, `references_test.go`, `report.go`, `request_bodies.go`, `request_bodies_test.go`, `required_fields_test.go`, `responses.go`, `responses_internal_test.go`, `responses_test.go`, `root.go`, `root_test.go`, `schemas.go`, `schemas_test.go`, `security.go`, `security_internal_test.go`, `security_test.go`, `servers.go`, `servers_internal_test.go`, `servers_test.go`, `swagger.go`, `swagger_defaults.go`, `swagger_test.go`, `tags.go`, `tags_test.go`.
- `xmlvalue`: `null.go`, `null_test.go`.


## Variants, generators and native boundaries

Root generate.go declares real local command directories for specmatrix and
modelgen. The producer owns all eight oas30/oas31/oas32/swagger20 model and test
outputs; they are not external source and must never be manually rewritten.
Trace actual command cwd/flags/inputs/output writes before changing templates or
headers. The retained checker marks the bare modelgen header for review; verify
with authenticated Tools2.4.1 command-directory support before interpreting that
as an unavailable declared generator. Keep every output symbol/token, frozen
prose/field inventory, corpus/license/source pin and generated test assertion.

Root native API mode is apidiff with the frozen api/v2.txt; no GoDoc-text baseline is configured.
Discover any native nested API/doc/example/generation recipe before selecting
checks. The interoperability module retains its independent published dependency
graph, public-proxy and optional external/public-description boundaries; source
auditing is not new public-consumer or conformance qualification. publicinterop
is an explicit test tag. The file resolver's os.Root containment/platform limits
must be traced, not converted into a fabricated platform guarantee.

## Missing documentation versus technical questions

Missing public fields/interfaces and command/interoperability package comments
are structural gaps. Canonical source/model package responsibility, immutable
semantic representation versus malformed typed access, provider defaults, nil/
zero/presence, borrowed streams/callbacks and mutation/concurrency, cancellation
and reference authorization, transformation loss and diagnostic/resource
claims require actual source/tests/specification evidence. Record unresolved
implementation limitations separately from unanswered policy/rationale; never
rewrite normative expectations or frozen decisions to endorse an existing defect.
Audit every private lifecycle/parser/cache/producer/assurance helper and every
existing test/example/fuzz/benchmark/marker, including generated output owners.

The pre-edit inventory and plan above are retained; completed semantic coverage
and native outcomes follow below.

## Semantic coverage and preserved variants

| Source family | Audited contracts and rationale |
| --- | --- |
| Root/specversion and four model editions | Exactly one root marker versus direct edition decoders; exact pinned revision, malformed versus unsupported; generated shape access retains unknown/invalid data without requiredness/enum validation. Zero views are not valid empty specification objects. |
| jsonvalue/model/modelaccess | Deeply immutable semantic children versus shallow generic T; outer slice copies and retained string backing, empty/absent/null/invalid states; synchronous decoder ownership and malformed reference-union branch selection. |
| parse | Buffered byte admission, decoded token/key/value/depth budgets, unpaired surrogate rejection, single JSON/YAML input, strict YAML tags/aliases, line-versus-byte offsets and unsanitized Cause; cancellation cannot preempt arbitrary IO or syntax work on buffered data. |
| reference | URI identity/$self, pointer versus anchor, duplicate anchors/cycles, supported bundling versus expansion, literal source-data exclusions, call-local caches, resource partial-failure effects; cumulative built-in reservations and explicit os.Root/transport cleanup. |
| validate/jsonschema | Official-schema versus prose rules and public JSON Schema v2 delegation; direction warnings versus Swagger errors, isolated single-flight construction/retry/publication, borrowed loader/codec/example policies, post-collection output caps, external-reference disabled outcomes versus invalid targets. |
| compose/convert/diff | Same exact-version merge versus same-dialect comparison; callbacks and source provenance, component rename/alias handling, source-ordered losses and selected downgrade/legacy variants. No full business equivalence or target-validity guarantee from successful transformation. |
| expression/implicit/discriminator | Already-decoded values versus code/IO execution, case-sensitive names/header ambiguity, callback templates and caller escaping, entry/current document scope and first-tag match versus ambiguous operation IDs, hints rather than schema validation. |
| parameter/media/response/server/XML | Isolated serialization/scalar coercion/non-JSON codec responsibilities, patch-specific empty/style choices, byte checks versus intermediate allocations, SSE UTF-8 replacement and EOF framing, exact/range/default response selection, literal server substitutions/no enum enforcement, caller-reported binary lengths and XML representation plans. |
| Maintainer commands/specification/interoperability | Authentic command cwd/flags/producer identities, bounded trusted-checkout decoders, per-file atomic writes/partial batches, historical evidence not a fresh conformance result, published-root nested dependency graph and contained fixtures. |

All tracked test files were inventoried by actual declarations, fixture/seed
selection and assertion boundaries, with focused body inspection at the changed
contracts and private lifecycle/parser/cache/producer seams. The 27 fuzz targets
now state their actual semantic/determinism/round-trip or panic-only assertions;
23 benchmarks identify deterministic workloads, timed setup/cold/warm state and
comparison assumptions. Existing three Examples retain executable bodies/output.
No new campaigns, fixture modifications or assurance threshold changes are made.
The publicinterop tag, both go:generate routes, go:embed, scanner directives and
all frozen inputs retain their actual owners; no unfinished source debt marker
was found in the tracked Go inventory. No platform-specific source suffix or
external generated Go implementation is present.

## Existing limitations versus unanswered rationale

Confirmed by static implementation, not new runtime qualification: expression
interpolation discards an internal non-string marshal error and yields empty
text; linkset input admission discards its default marshal error, and linkset
serialization applies an anchor only after encountering its source member.
These are existing limitations, not deliberate normative interpretations or a
conformance success. No runtime repair or expected-protocol rewrite is included.
The API seams and current security guidance distinguish them from actual bounds.
Generic aliasing, unsanitized reader/decoder causes, post-collection diagnostic
caps and cooperative blocking-IO cancellation are contract limitations, not
blanket process-memory/CPU/privacy guarantees.

Unanswered owner questions: intended error policy for expression/linkset marshal
failures; desired order-independent linkset anchor handling; whether later parser
admission should constrain allocations before third-party decoding; and whether
validation should add an early collection/work diagnostic policy. Existing
normative decisions, typed history, pin manifests, dated reviews and evidence.tsv
remain unchanged; answering these needs a separate behavioral assignment.

## Selected native verification scope

Use Go 1.27.2 from modules.json with task-owned three-cache session. Native root
`go generate .` runs specmatrix then modelgen; repeat generation and compare all
eight outputs, preserving executable AST/signatures/tokens and the two TSVs,
evidence.tsv, manifest/corpus/licenses/history/dependencies/API baseline. The
producer's only permitted literal deltas are emitted comments/provenance; its
comment-stripped executable output and other producer logic remain identical.
Run the native selected internal/specification/modelgen/provenance checks and
actual command-output smoke tests, plus the four generated model smoke packages.
Inspect parsed/HTML comments and native GoDoc for all 35 source packages/modules.
Run the existing root Example selector with normal vet, and configured root
apidiff check; no GoDoc-text baseline is configured and no frozen API refresh is
needed. Native conformance's expected Git output guard compares the current
explicitly staged generated documentation. This preserves the real guard while
allowing intentional generated comment changes; the guard is not omitted.
Offline specification/local documentation/links, source checker and pinned
added-content/staged scans remain applicable. Parent owns the independent whole
final Tier C diff review before commit and all external delivery.

## Native outcomes and proof limits

Pinned Go 1.27.2 bounded sessions ran with GOWORK=off and owned GOCACHE,
GOMODCACHE and GOTMPDIR. All 251 tracked Go files have equivalent comment-free
executable AST/signatures and exact directive/security annotations. Raw tokens
and literal spellings are exact in 250 files; the sole exception is seven
modelgen emitted-comment string literals. Removing only emitted Go comment lines
and normalizing the equivalent indexed package-format placeholder makes those
producer strings identical; no producer call, expression or control-flow delta
exists. All eight generated outputs have exact runtime AST/signatures/body tokens.
No optional-semicolon normalization or Example identifier change was needed.
All 2,249 actual comment groups parsed and rendered as native GoDoc HTML/text.

| Applicable command/contract | Result |
| --- | --- |
| gofmt on tracked Go; final affected gofmt + stdlib parser/scanner/comment helper | Passed; source preservation as qualified above. |
| Root go generate . twice | Passed; eight output files byte-stable. Canonical comments/header provenance came only from the genuine producer. |
| make -f verification/package.mk conformance | Passed: native extraction/modelgen tests, generation, actual Git output guard, provenance tests and offline authenticated manifest verification. Only the eight generated docs were explicitly staged for its unstaged-output guard. This is selected producer evidence, not renewed whole-API conformance. |
| go test ./internal/modelgen/cmd/modelgen ./internal/specification/cmd/specmatrix ./oas30 ./oas31 ./oas32 ./swagger20 -count=1 -timeout=2m | Passed existing command/output and four generated smoke packages. Assertions/corpora unchanged. |
| go test -run '^Example' -count=1 -timeout=2m . and root Example/local-link selector | Passed with normal go-test vet; all three existing Examples/body/output unchanged. The source-traced Example naming concern did not produce a native failure and was not treated as a defect or changed. |
| go doc -all for all 35 catalog packages in their owning module | Passed and inspected; final modelgen/media comments refreshed, original other 33 renders reused for unchanged inputs. Nested interoperability retained published OpenAPI v2.0.0, never a sibling replacement. |
| Declared authenticated Tools v2.3.0 api check | Passed root apidiff; interoperability API gate explicitly not applicable. api/v2.txt remains byte-exact. |
| Declared Tools v2.3.0 docs check --local | Passed native local documentation validation. |
| Declared Tools v2.3.0 specification check | Passed offline ten current specification decisions; no online check or normative rewrite. |
| Authenticated Tools v2.4.1 source-docs check --json | Zero errors, eight retained availability/output reviews for genuine modelgen headers. Native command/output/identity evidence resolves this audit's review need without changing header identity or pretending the proxy has no reviews. |

The initial authenticated v2.4.1 API invocation was rejected by v2.3.0 version
admission before the API check ran; it was not an API compatibility failure.
The coordinator supplied matching authenticated v2.3.0 (source
d38bf31e7403501075442de9fde78f9479052996, binary SHA-256
0c01e8d75e12c3e4930c3f399660489cf682004e6b894dd2bb07d5399024b745,
checksum-set SHA-256
639af50cddd7dbcc45561a067dd0f361abbddaf1ba5810c680b4fc5b0ba2c10b).
No repository/tool admission policy was weakened. Published v2.4.1 source-comment
inspection is distinct from the repository-declared native v2.3.0 gates.

The two TSV inventories, existing evidence.tsv, normative/specification corpora,
licenses, manifests, API baseline, dependency graphs, catalog, scripts, CI,
.verification and typed decision history are unchanged from base. The native
run also checked 153 selected current material inputs for byte stability;
152 are frozen non-Go authority/configuration inputs byte-identical to base.
The remaining specification/resources.go carries intentional comments only,
with embedded directives and executable AST preserved. No race, active fuzz,
performance, live service, publicinterop corpus campaign, new public consumer or
release qualification was run. Existing unchanged runtime evidence is not
represented as newly executed. Limitations above are source observations, not
new behavioral probes, remedies or policy endorsements.

Both bounded sessions removed their exact three caches on all exits. Helpers,
saved base source, logs and rendered outputs are removed after capturing these
outcomes; coordinator-owned verifiers/scanner are preserved. Ordinary complete
source-to-base diff inspection found no remaining material issue. The pinned
added-content/new-file scan passed with zero findings; final staging retains
that changed-content boundary rather than rescanning immutable historical data.
The parent independent whole-final-diff disposition remains pending before the
coherent Intent commit. Publication, hosted CI and enforcement remain
coordinator-owned.
