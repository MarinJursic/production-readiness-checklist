# Governance

Production Readiness Checklist is maintained by Marin Jursic. This document
explains how decisions are made while the project has one active maintainer.
It describes the current process; it does not imply a foundation, committee,
or independent certification body.

## Decision ownership

The maintainer has final responsibility for repository changes, releases,
security response, control semantics, and compatibility decisions. Community
feedback is welcome through issues and pull requests. A proposal is judged on
evidence, safety, technology neutrality, compatibility, maintenance cost, and
fit with the published product and trust contracts.

Material decisions should be recorded in a public issue, pull request, or
architecture document. Security-sensitive details remain private until safe
disclosure under [SECURITY.md](SECURITY.md).

## Changes that need extra review

Changes to control meaning, stable IDs, schemas, trust boundaries, evidence
authority, release workflows, package publication, scanner execution, AI
permissions, remediation permissions, or security policy require an explicit
maintainer review. Generated artifacts must be rebuilt from their declared
sources, and the pull request must explain any compatibility or evidence
invalidation effect.

No review may turn missing, stale, conflicting, incomplete, or lower-authority
evidence into Pass. AI output remains advisory unless a separately reviewed
deterministic contract and authoritative evidence establish the result.

## Releases

The maintainer approves release contents and versioning. Scanner releases must
use the documented tagged workflow, immutable inputs, public provenance, and
all required validation and security gates. A failed, compromised, or
materially incorrect release is revoked and replaced with a new version; tags
and version numbers are never reused.

## Conflicts of interest

Anyone reviewing a change should disclose a material financial, employment, or
personal interest that could reasonably affect the decision. The maintainer
records how that interest was handled. If the project gains additional
maintainers, a materially conflicted maintainer should not be the sole approver.

## Becoming a maintainer

Repeated, careful contributions that preserve the trust model can lead to an
invitation to maintain the project. The existing maintainer grants and revokes
access, records the change publicly, and applies least privilege. Repository,
npm, signing, and release access are separate capabilities and are not granted
automatically together.

## Continuity

If the maintainer can no longer operate the project, the preferred successor is
an established contributor with a public record of safe, compatible work. A
handover should rotate credentials, re-establish trusted publishing and signing
identities, review outstanding security reports, and publish the new ownership
clearly. If a safe transfer is not possible, the repository should be archived
and the last supported versions identified rather than silently abandoned.

## Process changes

Governance changes use the same public pull-request process as other material
changes. The reason, effect, and transition date must be recorded. Security
policy and license changes also follow their own stated requirements.
