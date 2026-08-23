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
