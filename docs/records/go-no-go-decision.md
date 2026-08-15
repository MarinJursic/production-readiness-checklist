# Go/no-go decision: [release identifier]

## Decision context

| Field | Value |
| --- | --- |
| Release/artifact/configuration | |
| Target environment and scope | |
| Evidence package | |
| Meeting/time | |
| Release authority | |
| Participants | |

## Gate summary

| Gate | Result | Evidence or exception |
| --- | --- | --- |
| No immediate no-go condition is present | | |
| Critical journeys meet approved objectives | | |
| Current evidence covers every critical requirement | | |
| No residual risk exceeds risk tolerance | | |
| Rollout and stop criteria are ready | | |
| Rollback or safe roll-forward is proven | | |
| Restore and incident response are proven | | |
| Monitoring, on-call, support, and communications are ready | | |
| Required sign-offs are complete | | |

## Open items

| Control/exception | Impact | Owner | Due/expiry | Monitoring |
| --- | --- | --- | --- | --- |
| | | | | |

## Decision

Select one:

- **GO** — Every release-blocking item passes and residual risk is within tolerance.
- **CONDITIONAL GO** — Remaining issues are not blockers, have tested compensating controls, named risk owners, monitoring, remediation dates, and automatic expiry; launch scope is bounded.
- **NO-GO** — Evidence is missing or stale, a blocker remains, recovery is unproven, required approval is absent, or remaining risk is not understood.

Decision:

Rationale:

Conditions or reduced scope:

Observation period:

Stop, rollback, and incident-declaration thresholds:

## Approval

| Role | Name | Approve/Reject | Time | Notes |
| --- | --- | --- | --- | --- |
| Release authority | | | | |
| Business risk owner | | | | |
| Engineering/operations | | | | |
| Required specialists | | | | |

## Production-readiness declaration

> Release [identifier], represented by artifact [digest] and configuration [version], has passed all applicable production-readiness gates documented in evidence package [reference]. No known unresolved risk exceeds the approved risk tolerance. Remaining accepted risks are recorded with owners, controls, monitoring, remediation dates, and expiry. Deployment, rollback or roll-forward, restoration, incident response, and customer communication have been tested and are ready.
