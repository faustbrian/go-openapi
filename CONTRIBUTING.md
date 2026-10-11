# Contributing

## Before Editing

1. Read [`AGENTS.md`](AGENTS.md) and the affected module's docs.
2. Run the standalone Go baseline below.
3. Identify owned dependencies and reverse dependants in `modules.json`.
4. Preserve unrelated work and generated/corpus provenance.

## Changes

Keep commits focused and conventional. Update every affected changelog with
the behavior and migration impact. Public API changes require compatibility
evidence and documentation. Specification behavior requires a decision record,
fixture coverage, and interoperability evidence.

New direct dependencies and dependency updates must follow the
[dependency governance policy](AGENTS.md#dependencies-and-supply-chain). Package-local
update bots are forbidden; the root policy owns every module and action update.

Specification-backed changes must follow the
[specification governance contract](AGENTS.md#design) and the
[specification decision register](docs/specification-decisions.md), update
the affected stable decision entries, and complete the Specification Decisions
section of the pull request template. An unresolved interpretation or stale
source pin is release-blocking; peer behavior cannot silently select policy.

Explicit full and release mutation gates must finish with zero surviving viable
mutants. Local mutation is optional, not a universal contribution prerequisite.

Do not fork shared verification workflows, add permanent replacements, machine-specific
paths, bypass flags, broad mutation exclusions, or aggregate quality metrics
that hide a failing package.

## Verification

Use the declared Go 1.27.2 producer; module minima remain Go 1.27.0. From the
affected module directory, test without sibling checkouts or the repository
workspace:

```bash
GOWORK=off go test ./...
```

Both declared modules have no services or test/build tags. The untagged command
runs deterministic fuzz seeds and saved corpus regressions, not active
exploration or the explicit `publicinterop` corpus. With the checksum-pinned
public Tools v2.1.1 verifier selected by `.golib.yaml`, run from the root:

```bash
make inventory
make check
make ci
```

Both Make routes use bounded local checks without building the verifier release
from source. One selected race verifies that concurrent callers receive one
shared dialect engine without duplicate construction; another verifies HTTP
resolver slot admission and cancellation while waiting. Each runs once within
90 seconds, three minutes total. Shared setup still checks out immutable Tools
source, but no sibling repository or goal ledger is needed.

Base-free checks run seeds/corpora; changed-source CI selects bounded active
smoke for causally affected targets. Short active exploration selects one
target; optional local mutation selects one package:

```bash
GOMAXPROCS=2 GOWORK=off go test -run '^$' -fuzz '^FuzzJSONParserDeterminism$' -fuzztime=2s -parallel=1 -timeout=2m ./parse
golib mutation --module . --package ./parse
```

`make ci` adds repository and offline specification/decision validation.
`make specification-check-online` selects upstream authority checks explicitly;
`make check-full` includes them and retains all 27 native 1,000-execution fuzz
targets, full coverage/race/mutation/API/security, generation/provenance
conformance, candidate and public-description interoperability, and named
performance allocation budgets with their original bounds. The isolated
interoperability module retains its pinned independent graph and public-proxy
requirements. Online authority and public-description qualification require
network access; a local pass does not prove public library publication.

PR, main and scheduled shared CI use bounded checks. Full campaigns are selected
explicitly with `make check-full` or `release_dry_run: true`; Required still
fails unless shared CI succeeds. Use `golib release check
--module .` and `golib release dry-run --module .` at the release boundary,
then verify the applicable consumer and publication claims separately. The
[current gate guide](docs/reference.md#development-gates) lists package-owned
operations and explicit CLI checks. Report every unavailable or failing
required command; do not describe partial results as release-ready.

## Selected Monthly Assurance

The owner-level [selected-assurance caller](.github/workflows/assurance.yml) is
configured for main only, on day 9 of each month at 04:31 UTC. It selects one
of the 27 current declared root-module fuzz targets using the shared UTC-week
rotation, with 30 seconds of active exploration. The selected risks are:
parser/value determinism and limits, reference resolution, transformations,
parameter decoding and owned manifest decoders.
HTTP resolver fuzzing uses an in-process fake transport and file resolution
uses owned temporary inputs; this is not live HTTP, public-description or
independent-consumer certification. The owned interoperability module remains
separate.

Base-free ordinary checks execute deterministic seeds, not active exploration;
committed-range CI selects affected active smoke. Manual assurance requires an
exact anchored selector from the root fuzz operation in `.golib.yaml`, for
example `^FuzzJSONParserDeterminism$`.
Empty, unknown and malformed selections refuse. The shared job bounds setup
and execution to ten minutes, active workers and CPU to two, owned memory to
2 GiB, native output to 1 MiB and private-cache lifetime. Failure, cancellation,
exhausted budgets and cleanup refusal remain attributable failures. This lane
does not run mutation or replace Required, original full fuzz budgets,
explicit `make check-full`, or relevant release qualification.

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.
