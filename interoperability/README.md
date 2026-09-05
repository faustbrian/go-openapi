# OpenAPI interoperability harness

This directory is a non-releasable engineering harness. It compares the root
`github.com/faustbrian/go-openapi` module with pinned maintained OpenAPI
implementations while keeping peer dependencies outside the public module.
Applications must not import this module as a production library.

The root [interoperability guide](../docs/interoperability.md) documents the
pinned implementations, fixture policy, observed results, limitations, and
update procedure. The checked-in [`expected.tsv`](expected.tsv) matrix is
observational evidence only; OpenAPI specifications and recorded package
decisions remain authoritative when implementations differ.
