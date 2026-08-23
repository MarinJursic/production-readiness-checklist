# Security policy

This repository contains guidance and automation, not a deployed service. Security issues can still matter—for example, unsafe commands in prompts, workflow vulnerabilities, malicious links, or guidance that could create a serious security failure.

## Reporting a vulnerability

Please use GitHub’s private vulnerability reporting feature under **Security → Advisories → Report a vulnerability**. Do not open a public issue for an undisclosed vulnerability.

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
