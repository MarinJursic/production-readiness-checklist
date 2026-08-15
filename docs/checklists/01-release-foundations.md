# Release foundations

> Define what readiness means, stop unsafe launches early, and identify the exact release under review.

Sections 1–3 of the master checklist. For each applicable item, capture status, owner, evidence, release, environment, reviewer, and evidence age.

## 1. How to operate this checklist

For **every checkbox**, record:

- **Status:** Pass, Fail, Blocked, or Not Applicable.
- **Applicability rationale:** especially for every Not Applicable item.
- **Owner:** one accountable person, not just a team name.
- **Evidence:** test report, log, screenshot, query, configuration export, code review, architecture decision, contract, drill report, or monitoring link.
- **Tested release:** commit, tag, artifact digest, image digest, configuration version, database migration version, and feature-flag state.
- **Environment:** exact environment in which the evidence was collected.
- **Reviewer:** someone other than the implementer for material controls.
- **Evidence date and expiry:** security scans, accessibility audits, capacity tests, and restore tests become stale.
- **Exception:** risk owner, justification, compensating control, monitoring, remediation date, and automatic expiry.

Rules:

- A Not Applicable result without a written reason is a failure.
- A release change made after evidence was captured invalidates every affected result.
- Do not calculate readiness by averaging scores. One serious failure can outweigh hundreds of passing items.
- Conditional sections become mandatory whenever their triggering feature, data type, user group, jurisdiction, or risk exists.
- Risk acceptance is not the same as fixing a problem. It must be explicit, time-bounded, visible, and owned.
- Evidence must cover the **actual production artifact and production configuration**, not merely a developer’s branch or an earlier build.

---

## 2. Immediate no-go conditions

The release must not proceed while any of these is true:

- [ ] **PRC-02-001** — A critical user journey has not been tested successfully end to end.
- [ ] **PRC-02-002** — Authentication, authorization, administrative access, or tenant isolation has not been independently verified.
- [ ] **PRC-02-003** — A known critical vulnerability or relevant known-exploited vulnerability remains unmitigated.
- [ ] **PRC-02-004** — Production secrets, signing keys, credentials, tokens, or personal data have been exposed.
- [ ] **PRC-02-005** — A destructive or irreversible data migration has not been rehearsed against representative data.
- [ ] **PRC-02-006** — Backup restoration has not been demonstrated within the required recovery objectives.
- [ ] **PRC-02-007** — Rollback or safe roll-forward has not been demonstrated.
- [ ] **PRC-02-008** — Expected peak load, burst load, or dependency-failure behavior is unknown.
- [ ] **PRC-02-009** — The release cannot meet its approved reliability objectives under expected operating conditions.
- [ ] **PRC-02-010** — Monitoring cannot detect failure of critical user journeys.
- [ ] **PRC-02-011** — No qualified person is available and authorized to respond during and immediately after launch.
- [ ] **PRC-02-012** — Required incident-response, escalation, customer-communication, or breach-notification procedures do not exist.
- [ ] **PRC-02-013** — A mandatory legal, privacy, payment, accessibility, safety, or contractual requirement is unmet.
- [ ] **PRC-02-014** — The production artifact cannot be traced to reviewed source, dependencies, build process, tests, and approval.
- [ ] **PRC-02-015** — An unsupported or end-of-life runtime, operating system, database, library, service, or protocol is required for operation.
- [ ] **PRC-02-016** — Known data corruption, cross-user disclosure, cross-tenant access, duplicate charging, or irreversible financial risk remains possible without a tested control.
- [ ] **PRC-02-017** — A single person, credential, machine, region, provider, or undocumented manual procedure is the only recovery path.
- [ ] **PRC-02-018** — Evidence was collected before material changes that have not been retested.
- [ ] **PRC-02-019** — A release blocker has been relabeled, suppressed, or accepted without an accountable risk owner.
- [ ] **PRC-02-020** — The team cannot describe concrete stop, rollback, and incident-declaration criteria before deployment begins.

Vulnerability priority should not be based on a raw severity score alone. Exploit activity, asset exposure, business impact, local controls, and environmental context must be considered; FIRST’s CVSS guidance specifically recommends enriching base scores with threat and environmental information. ([first.org](https://www.first.org/cvss/v4.0/implementation-guide))

---

## 3. Release control record

### 3.1 Release identity

- [ ] **PRC-03-001** — Assign a unique release identifier.
- [ ] **PRC-03-002** — Record the exact source commit or commits.
- [ ] **PRC-03-003** — Record the immutable application artifact digest.
- [ ] **PRC-03-004** — Record every deployable service, worker, scheduled job, frontend bundle, and infrastructure component included.
- [ ] **PRC-03-005** — Record database schema and migration versions.
- [ ] **PRC-03-006** — Record production configuration versions.
- [ ] **PRC-03-007** — Record feature flags and their intended launch states.
- [ ] **PRC-03-008** — Record third-party service and API versions relied upon.
- [ ] **PRC-03-009** — Record the intended production regions, accounts, projects, clusters, and domains.
- [ ] **PRC-03-010** — Record the target deployment date and change window.
- [ ] **PRC-03-011** — Record the previous known-good release and rollback artifact.
- [ ] **PRC-03-012** — Freeze or formally track changes after the readiness evidence snapshot.

### 3.2 Scope

- [ ] **PRC-03-013** — List every user-facing domain and subdomain.
- [ ] **PRC-03-014** — List every API, webhook, callback, streaming connection, and machine interface.
- [ ] **PRC-03-015** — List background workers, queues, scheduled tasks, batch processes, and administrative tools.
- [ ] **PRC-03-016** — List databases, object stores, search indexes, caches, message brokers, analytics stores, and backup locations.
- [ ] **PRC-03-017** — List external identity, payment, email, SMS, analytics, advertising, support, logging, CDN, and infrastructure providers.
- [ ] **PRC-03-018** — List all user types, administrative roles, service accounts, and external actors.
- [ ] **PRC-03-019** — List supported countries, languages, currencies, time zones, browsers, devices, and accessibility targets.
- [ ] **PRC-03-020** — List all data classifications and regulated data handled.
- [ ] **PRC-03-021** — List all excluded components and explain why they are outside the release assessment.
- [ ] **PRC-03-022** — Identify every shared platform or dependency that can affect the app even though another team operates it.

---
