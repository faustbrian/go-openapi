# OpenAPI interoperability harness

This directory is a non-releasable engineering harness. Its maintained module
uses the published `github.com/faustbrian/go-openapi` v1.0.0 dependency so
ordinary module builds do not require an unpublished v2 producer. It compares
OpenAPI implementations while keeping peer dependencies outside the public
module. Permanent harness adoption of v2 remains pending its publication.
Applications must not import this module as a production library.

The repository interoperability target adapts the maintained runner's imports
and owned root dependency to v2 only inside a disposable copy. Candidate mode
uses a zero pseudo-version placeholder with a local replacement; it does not
qualify a public release:

```sh
make -f verification/package.mk interoperability
```

After publication, check the public module from the repository root:

```sh
OPENAPI_INTEROPERABILITY_VERSION=v2.0.0 ./scripts/check-interoperability.sh
```

The script defaults to public mode, requires an explicit published v2 version,
uses the public Go proxy and checksum database without any replacement, and
rejects the candidate placeholder. Both modes preserve peer dependencies and
leave tracked source and module files unchanged.

The root [interoperability guide](../docs/interoperability.md) documents the
pinned implementations, fixture policy, observed results, limitations, and
update procedure. The checked-in [`expected.tsv`](expected.tsv) matrix is
observational evidence only; OpenAPI specifications and recorded package
decisions remain authoritative when implementations differ.
