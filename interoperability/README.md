# OpenAPI interoperability harness

This directory is a non-releasable engineering harness. It compares the root
`github.com/faustbrian/go-openapi/v2` module with pinned maintained OpenAPI
implementations while keeping peer dependencies outside the public module.
Applications must not import this module as a production library.

The root dependency's zero pseudo-version is an unpublished candidate
placeholder, not a release pin. The repository interoperability target
explicitly checks the local candidate using a replacement inside a disposable
copy; it does not qualify a public release:

```sh
make -f verification/package.mk interoperability
```

After publication, check the public module from the repository root:

```sh
OPENAPI_INTEROPERABILITY_VERSION=v2.0.0 ./scripts/check-interoperability.sh
```

The script defaults to public mode, uses the public Go proxy
and checksum database without a local replacement, and rejects the candidate
placeholder. Both modes leave tracked module files unchanged.

The root [interoperability guide](../docs/interoperability.md) documents the
pinned implementations, fixture policy, observed results, limitations, and
update procedure. The checked-in [`expected.tsv`](expected.tsv) matrix is
observational evidence only; OpenAPI specifications and recorded package
decisions remain authoritative when implementations differ.
