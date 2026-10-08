# Security Policy

## Reporting

Report suspected vulnerabilities through a
[private GitHub security advisory](https://github.com/faustbrian/go-openapi/security/advisories/new).
Do not open a public issue containing exploit details, credentials, private
fixtures, or affected deployment information.

Include the affected module and version, impact, reproduction, preconditions,
and any suggested mitigation. Reports are acknowledged as soon as practical;
timelines depend on severity and verification.

## Supported Versions

The current stable release line receives security fixes. Publishing v2 does
not end security support for the latest stable `v1` release line.
Support windows are documented per module and in
[`COMPATIBILITY.md`](COMPATIBILITY.md).

The published v1.0.0 security evaluator lacks credential count and label-byte
admission. It is not fixed by changing only its requirement limits. Upgrade
to `github.com/faustbrian/go-openapi/v2` v2.0.0 or newer. If migration is
temporarily unavailable, apply the caller-owned admission and validation
workaround in [the security model](docs/security.md#historical-v1-mitigation)
before evaluating untrusted requirements or credentials. Continued v1 report
triage does not imply that a patched v1 version has been published.

## Security Gates

Releases require isolated tests, race and hostile-input checks, exact coverage
and mutation results, `govulncheck`, secret scanning, license verification,
SBOM generation, provenance validation, and clean-consumer resolution. A
missing scanner or unavailable service is a failed gate, not a warning.

Security fixes MUST include a regression test that does not publish weaponized
details or real secrets. Credentials MUST be redacted from logs and evidence.

## Repository Assurance

The repository [safety and concurrency policy](AGENTS.md#safety-and-concurrency)
and [supply-chain policy](AGENTS.md#dependencies-and-supply-chain) define shared
trust boundaries and release requirements. Package-specific security guidance
refines those rules for its owned boundary.
