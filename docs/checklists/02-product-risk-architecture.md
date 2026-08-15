# Product, risk, and architecture

> Align business intent, accountable ownership, risk tolerance, and system understanding.

Sections 4–6 of the master checklist. For each applicable item, capture status, owner, evidence, release, environment, reviewer, and evidence age.

## 4. Product and business readiness

### 4.1 Product definition

- [ ] **PRC-04-001** — The problem being solved and intended users are documented.
- [ ] **PRC-04-002** — Release scope, non-goals, and deferred features are explicit.
- [ ] **PRC-04-003** — Every feature has testable acceptance criteria.
- [ ] **PRC-04-004** — Critical user journeys are ranked by business and user impact.
- [ ] **PRC-04-005** — Critical journeys include onboarding, normal use, error recovery, account recovery, cancellation, and deletion where applicable.
- [ ] **PRC-04-006** — Product behavior for unavailable dependencies and partial failure is specified.
- [ ] **PRC-04-007** — Product behavior for empty, malformed, incomplete, stale, or conflicting data is specified.
- [ ] **PRC-04-008** — User-visible terminology is consistent.
- [ ] **PRC-04-009** — Pricing, entitlements, quotas, billing rules, and plan limits are unambiguous.
- [ ] **PRC-04-010** — Cancellation, refund, reversal, dispute, and undo behavior is defined.
- [ ] **PRC-04-011** — Expected customer support obligations are documented.
- [ ] **PRC-04-012** — Service hours, maintenance expectations, SLOs, SLAs, and contractual commitments are aligned.
- [ ] **PRC-04-013** — The launch does not make claims that the implemented product cannot substantiate.
- [ ] **PRC-04-014** — Product management has accepted the final behavior of critical journeys.

### 4.2 Launch objectives

- [ ] **PRC-04-015** — Define measurable launch success criteria.
- [ ] **PRC-04-016** — Define measurable launch failure and rollback criteria.
- [ ] **PRC-04-017** — Define adoption, conversion, quality, reliability, support, and safety indicators.
- [ ] **PRC-04-018** — Define the observation period before declaring the launch complete.
- [ ] **PRC-04-019** — Define who has authority to pause, continue, expand, or roll back the launch.
- [ ] **PRC-04-020** — Establish expected traffic and customer cohorts.
- [ ] **PRC-04-021** — Establish whether the launch is internal, private, regional, percentage-based, invitation-only, or public.
- [ ] **PRC-04-022** — Ensure marketing, sales, support, and operations use the same release date and scope.
- [ ] **PRC-04-023** — Ensure help content, screenshots, pricing pages, documentation, and legal text match the released behavior.
- [ ] **PRC-04-024** — Ensure customers are notified of breaking changes, downtime, migrations, or changed data use where required.

---

## 5. Governance, ownership, and risk

### 5.1 Ownership

- [ ] **PRC-05-001** — Every production service has a named business owner.
- [ ] **PRC-05-002** — Every production service has a named engineering owner.
- [ ] **PRC-05-003** — Every data set has an accountable data owner.
- [ ] **PRC-05-004** — Every security-sensitive component has a security contact.
- [ ] **PRC-05-005** — Every third-party integration has an owner.
- [ ] **PRC-05-006** — Every operational dashboard, alert, and runbook has an owner.
- [ ] **PRC-05-007** — Every scheduled job and background process has an owner.
- [ ] **PRC-05-008** — Ownership remains valid during holidays, leave, and personnel changes.
- [ ] **PRC-05-009** — Escalation paths include primary, secondary, management, security, privacy/legal, and provider escalation contacts.
- [ ] **PRC-05-010** — Offboarding and role-change processes remove or adjust access promptly.

### 5.2 Risk assessment

- [ ] **PRC-05-011** — Classify confidentiality impact.
- [ ] **PRC-05-012** — Classify integrity impact.
- [ ] **PRC-05-013** — Classify availability impact.
- [ ] **PRC-05-014** — Classify privacy impact.
- [ ] **PRC-05-015** — Classify financial and fraud impact.
- [ ] **PRC-05-016** — Classify physical, psychological, societal, and user-safety impact.
- [ ] **PRC-05-017** — Classify regulatory and contractual impact.
- [ ] **PRC-05-018** — Classify reputational and customer-trust impact.
- [ ] **PRC-05-019** — Identify likely attackers, abusive users, insiders, compromised accounts, automated agents, and accidental failure sources.
- [ ] **PRC-05-020** — Identify foreseeable misuse and abuse cases, not only intended use.
- [ ] **PRC-05-021** — Maintain a release-specific risk register.
- [ ] **PRC-05-022** — Give each risk an owner, likelihood, impact, treatment, evidence, and review date.
- [ ] **PRC-05-023** — Record accepted residual risk and the person authorized to accept it.
- [ ] **PRC-05-024** — Review high-impact assumptions with independent subject-matter experts.
- [ ] **PRC-05-025** — Verify that launch scope remains inside the organization’s risk appetite.
- [ ] **PRC-05-026** — Ensure risk exceptions cannot remain open indefinitely.
- [ ] **PRC-05-027** — Ensure compensating controls are tested rather than merely described.

### 5.3 Control framework

- [ ] **PRC-05-028** — Map applicable security controls to a recognized verification standard.
- [ ] **PRC-05-029** — Map applicable privacy controls to legal and organizational requirements.
- [ ] **PRC-05-030** — Map accessibility requirements to the selected WCAG conformance target.
- [ ] **PRC-05-031** — Map operational controls to SLO, incident, continuity, and recovery requirements.
- [ ] **PRC-05-032** — Map contractual controls to customer, processor, provider, and partner agreements.
- [ ] **PRC-05-033** — Identify conflicts between product requirements, legal requirements, and technical controls.
- [ ] **PRC-05-034** — Document how conflicts were resolved and who approved the resolution.
- [ ] **PRC-05-035** — Apply separation of duties to high-impact changes and approvals.
- [ ] **PRC-05-036** — Periodically review whether the selected control level still matches the app’s risk.

---

## 6. Architecture and system understanding

- [ ] **PRC-06-001** — Maintain an up-to-date system-context diagram.
- [ ] **PRC-06-002** — Maintain component and deployment diagrams.
- [ ] **PRC-06-003** — Maintain end-to-end data-flow diagrams.
- [ ] **PRC-06-004** — Mark trust boundaries explicitly.
- [ ] **PRC-06-005** — Mark Internet-facing and externally reachable interfaces.
- [ ] **PRC-06-006** — Mark administrative and support interfaces.
- [ ] **PRC-06-007** — Mark where authentication and authorization decisions occur.
- [ ] **PRC-06-008** — Mark encryption and decryption boundaries.
- [ ] **PRC-06-009** — Mark where sensitive data is collected, transformed, stored, cached, logged, exported, backed up, and deleted.
- [ ] **PRC-06-010** — Mark all third-party data transfers.
- [ ] **PRC-06-011** — Identify every stateful component.
- [ ] **PRC-06-012** — Identify every asynchronous boundary and queue.
- [ ] **PRC-06-013** — Identify common dependencies and shared failure domains.
- [ ] **PRC-06-014** — Identify single points of failure.
- [ ] **PRC-06-015** — Identify control-plane dependencies needed during incidents or recovery.
- [ ] **PRC-06-016** — Identify circular and hidden dependency chains.
- [ ] **PRC-06-017** — Document required consistency, ordering, durability, and freshness properties.
- [ ] **PRC-06-018** — Document architectural invariants that must never be violated.
- [ ] **PRC-06-019** — Document expected failure behavior for every critical dependency.
- [ ] **PRC-06-020** — Record material architecture decisions and rejected alternatives.
- [ ] **PRC-06-021** — Update diagrams whenever the release changes data flow, trust, topology, or dependencies.
- [ ] **PRC-06-022** — Verify the documented architecture against the deployed environment rather than assuming documentation is accurate.

### 6.1 Threat and abuse modeling

- [ ] **PRC-06-023** — Perform a threat model for the actual release scope.
- [ ] **PRC-06-024** — Include all trust boundaries, actors, entry points, data stores, and external systems.
- [ ] **PRC-06-025** — Model spoofing and identity attacks.
- [ ] **PRC-06-026** — Model unauthorized modification and integrity attacks.
- [ ] **PRC-06-027** — Model information disclosure and privacy attacks.
- [ ] **PRC-06-028** — Model denial of service and resource exhaustion.
- [ ] **PRC-06-029** — Model privilege escalation.
- [ ] **PRC-06-030** — Model repudiation and inadequate accountability.
- [ ] **PRC-06-031** — Model business-logic abuse.
- [ ] **PRC-06-032** — Model insider and compromised-administrator threats.
- [ ] **PRC-06-033** — Model dependency and build-system compromise.
- [ ] **PRC-06-034** — Model tenant-escape and cross-customer attacks.
- [ ] **PRC-06-035** — Model automation, scraping, spam, fraud, and economic abuse.
- [ ] **PRC-06-036** — Model recovery-system and backup attacks.
- [ ] **PRC-06-037** — Model support-channel and account-recovery social engineering.
- [ ] **PRC-06-038** — Convert identified threats into requirements and tests.
- [ ] **PRC-06-039** — Revisit the threat model after material design or scope changes.
- [ ] **PRC-06-040** — Obtain independent review for high-value or high-risk systems.

---
