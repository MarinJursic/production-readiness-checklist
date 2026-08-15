# Release readiness assessment: [release identifier]

> Copy this file for each release. Link detailed evidence records instead of pasting secrets, personal data, or large raw reports.

## Release identity

| Field | Value |
| --- | --- |
| Release identifier | |
| Source commit/tag | |
| Artifact/image digest | |
| Configuration version | |
| Database migrations | |
| Feature-flag state | |
| Target environment/regions | |
| Previous known-good release | |
| Planned deployment window | |
| Assessment owner | |
| Evidence snapshot time | |

## Scope

### Critical user journeys

- [Journey, user, expected outcome]

### Included components

- [Service, worker, frontend, scheduled job, infrastructure component]

### Data and third parties

- [Data classification, store, processor/provider]

### Exclusions

| Excluded area | Rationale | Reviewer |
| --- | --- | --- |
| | | |

## Track applicability

| Track | Required? | Owner | Rationale or evidence location |
| --- | --- | --- | --- |
| Release foundations | Yes | | |
| Product, risk, and architecture | | | |
| Source, build, and supply chain | | | |
| Environments, quality, and experience | | | |
| Application security | | | |
| Data, privacy, and performance | | | |
| Reliability and operations | | | |
| Maintenance, vendors, and compliance | | | |
| Conditional feature modules | | | |
| Evidence, sign-off, and decision | Yes | | |

## Immediate no-go screen

No-go controls describe dangerous conditions. A **Yes** result stops the release.

| Control | Condition present? | Evidence | Owner | Resolution |
| --- | --- | --- | --- | --- |
| PRC-02-001–PRC-02-020 | No / Yes / Unknown | | | |

## Control results

Duplicate rows as needed or link to a control-management system.

| Control ID | Status | Evidence | Owner | Reviewer | Evidence date/expiry | Next action |
| --- | --- | --- | --- | --- | --- | --- |
| PRC- | Pass / Fail / Blocked / N/A | | | | | |

## Open blockers and failures

| Control ID | Impact | Owner | Required action | Due |
| --- | --- | --- | --- | --- |
| | | | | |

## Risk exceptions

| Exception | Controls | Risk owner | Compensating control | Expiry |
| --- | --- | --- | --- | --- |
| | | | | |

## Evidence package

- Architecture and data-flow diagrams:
- Threat model:
- Test reports:
- Security and dependency reports:
- Accessibility evidence:
- Performance and capacity evidence:
- Resilience, restore, and rollback evidence:
- Monitoring, alerting, on-call, and runbooks:
- Legal, privacy, contractual, and third-party review:
- Deployment and post-deployment plan:

## Sign-offs

| Area | Signer | Decision | Date | Notes |
| --- | --- | --- | --- | --- |
| Product | | Approve / Reject / Conditional | | |
| Engineering | | | | |
| Quality | | | | |
| Security | | | | |
| Reliability/operations | | | | |
| Data | | | | |
| Privacy | | | | |
| Accessibility | | | | |
| Legal/compliance | | | | |
| Support | | | | |
| Business risk | | | | |
| Release authority | | | | |

## Decision

- Decision: **GO / CONDITIONAL GO / NO-GO**
- Decision-maker:
- Time:
- Rationale:
- Observation period:
- Stop/rollback criteria:
- Follow-up review:

Use the full [go/no-go decision record](go-no-go-decision.md) for material launches.
