# Evidence, sign-off, and decision

> Assemble the evidence package, obtain accountable sign-off, and make the final go/no-go decision.

Sections 40–43 of the master checklist. For each applicable item, capture status, owner, evidence, release, environment, reviewer, and evidence age.

## 40. Final evidence package

Before approval, assemble one immutable or versioned evidence package containing:

- [ ] **PRC-40-001** — Release identifier, source commits, artifact digests, configuration versions, migrations, and flag states.
- [ ] **PRC-40-002** — Final scope and architecture diagrams.
- [ ] **PRC-40-003** — Data inventory and data-flow records.
- [ ] **PRC-40-004** — Threat model and abuse-case assessment.
- [ ] **PRC-40-005** — Applicable legal, privacy, accessibility, security, and contractual requirements.
- [ ] **PRC-40-006** — Requirements-to-test traceability.
- [ ] **PRC-40-007** — Functional and user-acceptance test results.
- [ ] **PRC-40-008** — Browser, device, localization, and accessibility results.
- [ ] **PRC-40-009** — Performance, load, stress, endurance, and capacity results.
- [ ] **PRC-40-010** — Resilience, failover, rollback, and recovery results.
- [ ] **PRC-40-011** — Security scan and penetration-test results.
- [ ] **PRC-40-012** — SBOM, license report, provenance, and artifact-signature evidence.
- [ ] **PRC-40-013** — Vulnerability disposition and known-exploitation review.
- [ ] **PRC-40-014** — Infrastructure and production-configuration review.
- [ ] **PRC-40-015** — Backup and restoration evidence.
- [ ] **PRC-40-016** — Incident-response and continuity drill evidence.
- [ ] **PRC-40-017** — Monitoring, alert, dashboard, on-call, and runbook review.
- [ ] **PRC-40-018** — Third-party risk and contract review.
- [ ] **PRC-40-019** — Open defect list.
- [ ] **PRC-40-020** — Open risk and exception register.
- [ ] **PRC-40-021** — Every exception’s owner, compensating control, expiry, and remediation plan.
- [ ] **PRC-40-022** — Deployment, rollback, and post-deployment validation plans.
- [ ] **PRC-40-023** — Approval record and decision rationale.

---

## 41. Required sign-offs

Sign-off does not transfer operational responsibility; it confirms that the signer reviewed the evidence within their area.

- [ ] **PRC-41-001** — Product owner: scope, behavior, acceptance criteria, launch objectives, and customer impact.
- [ ] **PRC-41-002** — Engineering owner: architecture, code, tests, migrations, maintainability, and technical risk.
- [ ] **PRC-41-003** — Quality owner: test completeness, evidence, defects, and regression status.
- [ ] **PRC-41-004** — Security owner: threat model, ASVS coverage, findings, supply chain, and residual security risk.
- [ ] **PRC-41-005** — Reliability or operations owner: SLOs, capacity, monitoring, on-call, rollout, rollback, and recovery.
- [ ] **PRC-41-006** — Data owner: integrity, migration, retention, reconciliation, and recovery.
- [ ] **PRC-41-007** — Privacy owner: data inventory, purpose, rights, notices, transfers, and privacy risk.
- [ ] **PRC-41-008** — Accessibility owner or qualified reviewer: target conformance and unresolved barriers.
- [ ] **PRC-41-009** — Legal or compliance owner: applicable obligations, terms, contracts, and regulatory risks.
- [ ] **PRC-41-010** — Support owner: support readiness, escalation, documentation, and customer communications.
- [ ] **PRC-41-011** — Business risk owner: explicit acceptance of remaining material risk.
- [ ] **PRC-41-012** — Executive or designated release authority: final go/no-go decision for high-impact launches.

A person should not sign an area for which no evidence exists or which they are not qualified or authorized to assess.

---

## 42. Final go/no-go decision rule

### GO

A release may be approved only when:

- [ ] **PRC-42-001** — Every applicable release-blocking item passes.
- [ ] **PRC-42-002** — Every critical requirement has current evidence tied to the exact production artifact.
- [ ] **PRC-42-003** — Critical user journeys meet their approved correctness, performance, security, accessibility, and reliability objectives.
- [ ] **PRC-42-004** — No known residual risk exceeds the organization’s risk appetite.
- [ ] **PRC-42-005** — Every accepted risk has an accountable owner and expiration.
- [ ] **PRC-42-006** — Rollout, rollback or roll-forward, restoration, incident response, and communication are demonstrably workable.
- [ ] **PRC-42-007** — The necessary operational and support personnel are ready.
- [ ] **PRC-42-008** — Required sign-offs are complete.

### CONDITIONAL GO

Conditional approval is appropriate only when:

- [ ] **PRC-42-009** — The remaining issue is not a release blocker.
- [ ] **PRC-42-010** — User, security, privacy, integrity, financial, legal, and safety impact is bounded and understood.
- [ ] **PRC-42-011** — A tested compensating control exists.
- [ ] **PRC-42-012** — Monitoring will detect deterioration.
- [ ] **PRC-42-013** — The feature can be disabled or rolled back safely.
- [ ] **PRC-42-014** — A named risk owner accepts the issue.
- [ ] **PRC-42-015** — A near-term remediation date and automatic expiry exist.
- [ ] **PRC-42-016** — Launch scope is reduced appropriately.

### NO-GO

The decision is no-go whenever:

- [ ] **PRC-42-017** — Evidence is missing or stale.
- [ ] **PRC-42-018** — A critical control is assumed rather than tested.
- [ ] **PRC-42-019** — A serious anomaly remains unexplained.
- [ ] **PRC-42-020** — Rollback, restoration, or incident response is unavailable.
- [ ] **PRC-42-021** — Required legal or contractual approval is absent.
- [ ] **PRC-42-022** — The team is relying on “we will monitor it” without a concrete, tested detection and mitigation path.
- [ ] **PRC-42-023** — Organizational pressure is the only reason to accept a risk.
- [ ] **PRC-42-024** — The decision-maker does not understand the remaining risk.

---

## 43. The production-readiness declaration

The strongest responsible declaration is:

> **Release [identifier], represented by artifact [digest] and configuration [version], has passed all applicable production-readiness gates documented in evidence package [reference]. No known unresolved risk exceeds the approved risk tolerance. Remaining accepted risks are recorded with owners, controls, monitoring, remediation dates, and expiry. Deployment, rollback or roll-forward, restoration, incident response, and customer communication have been tested and are ready.**

Do **not** declare:

> “This application has no issues,” “cannot fail,” “is completely secure,” or “has no holes.”

Those are not testable engineering conclusions. The checklist instead gives you an auditable basis for saying the release is **fit for production at a defined and accepted level of risk**.
