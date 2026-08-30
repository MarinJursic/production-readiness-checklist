# Security policy

This repository contains guidance and automation, not a deployed service. Security issues can still matter—for example, unsafe commands in prompts, workflow vulnerabilities, malicious links, or guidance that could create a serious security failure.

## Reporting a vulnerability

Please use [GitHub’s private vulnerability reporting form](https://github.com/MarinJursic/production-readiness-checklist/security/advisories/new). Do not open a public issue for an undisclosed vulnerability.

Include:

- the affected file, control, prompt, or workflow;
- the security impact and realistic abuse scenario;
- reproduction or evidence;
- a suggested correction, if available.

Avoid including real credentials, personal data, exploit payloads against systems you do not own, or other sensitive material.

## Repository security gates

Repository-owned workflows enforce complementary checks rather than treating any
single tool as proof of safety:

- CodeQL analyzes Go, Python, and GitHub Actions on pull requests, pushes to the
  default branch, and a weekly schedule with the `security-extended` query suite;
- Gitleaks scans the complete committed Git history on pull requests and default-
  branch pushes, with redaction enabled so a finding is not copied into logs;
- dependency review rejects pull requests that introduce a moderate-or-higher
  vulnerability in runtime, development, or unknown scope dependencies and also
  evaluates dependency license metadata; and
- validation and release workflows independently audit pinned Python packages,
  reachable Go vulnerabilities, static analysis findings, and the scanner's
  targeted filesystem-safety rules. Release tags repeat the dependency and
  complete-history secret gates before publishing.

Actions are pinned to immutable commit SHAs and tool invocations use exact module
versions. Dependabot proposes monthly updates so a pin does not silently become a
permanent version. A protected default branch should require the Validate,
CodeQL, Dependency review, and Secret scan checks and should enable GitHub secret
push protection. Those hosted repository settings are not encoded by workflow
files and must be verified in repository administration after the workflows land.

These checks reduce specific risks; they do not replace review, private
vulnerability reporting, incident response, or consumer verification.

### Historical upstream rules

The scanner's Gitleaks adapter is based on a pinned upstream rules file. An
older commit briefly stored that raw file at
`scanner/adapter/gitleaks-v8.30.0.toml`. Some of its upstream detection examples
look like Google API keys, so GitHub can report those examples as repository
secrets even though they are rule data rather than project credentials.

The raw historical path is the repository's only secret-scanning exclusion.
Current releases store the rules as `gitleaks-v8.30.0.toml.gz`; the scanner
checks the SHA-256 of both the compressed file and its expanded bytes before it
uses them. A repository test also requires the raw path to stay absent, the
exclusion to stay limited to that one path, and both hashes to match the Go
adapter. All other source, documentation, workflow, test, and adapter files
remain in GitHub secret scanning and the complete-history Gitleaks workflow.

## npm distribution security

The supported personal installation is global and stays outside scanned
projects: `npm install -g --ignore-scripts @marinjursic/prc@X.Y.Z`. The launcher
has no third-party JavaScript dependencies, lifecycle scripts, network fallback,
post-install download, or background updater. npm installs only the matching
OS/CPU package. The launcher verifies its release-bound manifest, native binary,
and every bundled runtime file before each execution.

Release packaging excludes website media and contributor-only documentation,
while retaining every control source, catalog contract, schema, adapter
manifest, pack, and benchmark fixture needed by scanner commands. Hard file and
byte budgets stop unexpected package growth instead of silently reducing scan
coverage. A normal scan does not install target dependencies or run target
scripts, and its default reports are stored outside the target with bounded
retention.

New npm versions must come from the exact tagged GitHub Actions workflow using
short-lived OIDC trusted publishing and npm provenance. Long-lived npm tokens,
local execution of the publisher, mismatched workflow identity, missing OIDC,
different immutable registry bytes, or an unbound tarball member stop the
release. npm provenance establishes origin, not freedom from vulnerabilities;
consumers should still use the latest supported patch and review security
advisories.

The release job's Python validation environment is binary-only and hash-bound
to the exact Linux wheels in `requirements-release.lock.txt`. The final
publication job verifies package structure with Python's standard library and
does not install extra Python packages before attestation or npm publication.

## Scope and response

Maintainers will assess reports affecting the repository’s content or automation. Application-specific findings discovered while using the checklist belong to the affected application’s security process, not this public repository.

The latest published checklist minor release receives documentation and workflow
security corrections. Before scanner 1.0, only the latest published scanner
prerelease is supported; older scanner prereleases should be upgraded rather
than patched in parallel. Development branch builds and failed release-workflow
artifacts are unsupported.

The maintainer targets acknowledgement within three business days and an initial
severity and scope assessment within seven business days. These are coordination
targets, not a contractual SLA. Resolution timing depends on exploitability,
affected users, upstream fixes, and safe disclosure coordination. If a report
indicates active harm, contact the affected provider or organization through its
published incident channel as well.

## Compromise and revocation

A compromised or materially incorrect scanner release, workflow dependency,
adapter, pack, catalog, or publisher identity is denied for future use rather
than silently replaced. The maintainer will identify affected versions and
digests, publish an advisory when safe, mark releases as revoked, remove or
revoke attestations where supported, update trust-store or registry revocations,
and publish a new version. A version, tag, or release-asset name is never reused.

See [scanner releases and verification](docs/scanner/releases.md) for the
consumer verification and revocation procedure.
