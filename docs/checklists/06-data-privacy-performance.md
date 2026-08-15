# Data, privacy, and performance

> Protect data integrity and privacy while proving user-experience and capacity targets.

Sections 23–26 of the master checklist. For each applicable item, capture status, owner, evidence, release, environment, reviewer, and evidence age.

## 23. Data model, databases, queues, caches, and integrity

- [ ] **PRC-23-001** — Every data set has an owner and classification.
- [ ] **PRC-23-002** — The data model enforces required constraints.
- [ ] **PRC-23-003** — Required fields cannot silently become absent.
- [ ] **PRC-23-004** — Uniqueness rules are enforced at the correct consistency boundary.
- [ ] **PRC-23-005** — Referential integrity is enforced or continuously checked.
- [ ] **PRC-23-006** — Transactions cover all operations requiring atomicity.
- [ ] **PRC-23-007** — Locking and concurrency strategies are documented.
- [ ] **PRC-23-008** — Lost updates and write races are prevented.
- [ ] **PRC-23-009** — Duplicate events and messages are handled safely.
- [ ] **PRC-23-010** — Queue delivery semantics are understood.
- [ ] **PRC-23-011** — Consumers are idempotent where messages may be redelivered.
- [ ] **PRC-23-012** — Ordering assumptions are explicit and tested.
- [ ] **PRC-23-013** — Poison messages cannot block an entire queue indefinitely.
- [ ] **PRC-23-014** — Dead-letter queues have monitoring, ownership, and safe replay procedures.
- [ ] **PRC-23-015** — Batch jobs checkpoint and resume safely.
- [ ] **PRC-23-016** — Reconciliation detects dropped, duplicated, or inconsistent records.
- [ ] **PRC-23-017** — Integrity checks or checksums exist where corruption risk warrants them.
- [ ] **PRC-23-018** — Replication lag is measured and accounted for.
- [ ] **PRC-23-019** — Read-after-write behavior matches user expectations.
- [ ] **PRC-23-020** — Failover does not create split-brain or conflicting writes.
- [ ] **PRC-23-021** — Cache keys include every required identity, tenant, locale, permission, and variant dimension.
- [ ] **PRC-23-022** — Cache invalidation is correct for security and correctness-sensitive data.
- [ ] **PRC-23-023** — Cache failure does not disclose data or violate invariants.
- [ ] **PRC-23-024** — Search indexes enforce access controls and deletion.
- [ ] **PRC-23-025** — Search-index freshness and reconciliation are monitored.
- [ ] **PRC-23-026** — Character encoding, collation, case, normalization, and sorting are intentional.
- [ ] **PRC-23-027** — Date and time values use an unambiguous storage representation.
- [ ] **PRC-23-028** — Data quality indicators cover completeness, validity, consistency, timeliness, and duplication.
- [ ] **PRC-23-029** — Orphaned data and references are detected.
- [ ] **PRC-23-030** — Data repair procedures are reviewed, reversible where possible, and audited.
- [ ] **PRC-23-031** — Manual production data changes require authorization and an audit trail.

### 23.1 Schema changes and migrations

- [ ] **PRC-23-032** — Every migration has a reviewed plan.
- [ ] **PRC-23-033** — Migration behavior is tested against representative data volume and shape.
- [ ] **PRC-23-034** — Migration duration and resource consumption are measured.
- [ ] **PRC-23-035** — Online migrations avoid unacceptable locking and downtime.
- [ ] **PRC-23-036** — Application and schema changes are backward- and forward-compatible during mixed-version deployment where required.
- [ ] **PRC-23-037** — Old and new application versions can coexist safely during rollout.
- [ ] **PRC-23-038** — A backup or recovery point exists before destructive changes.
- [ ] **PRC-23-039** — Migration preconditions are validated.
- [ ] **PRC-23-040** — Migration progress and failures are observable.
- [ ] **PRC-23-041** — Partial migration can resume or recover safely.
- [ ] **PRC-23-042** — Row counts, totals, constraints, checksums, and business invariants are verified afterward.
- [ ] **PRC-23-043** — Rollback feasibility is documented honestly.
- [ ] **PRC-23-044** — When rollback is unsafe, a tested roll-forward recovery exists.
- [ ] **PRC-23-045** — Deprecated columns, data, and indexes are removed only after all consumers have migrated.
- [ ] **PRC-23-046** — Data migrations are rehearsed using the same automation used for production.

---

## 24. Privacy and data protection

A universal checklist cannot enumerate every privacy law. The mandatory control is an applicability register based on users, jurisdictions, data, purposes, sector, and transfers. GDPR is one major example requiring data protection by design/default, processing security, rights handling, and breach governance. ([eur-lex.europa.eu](https://eur-lex.europa.eu/eli/reg/2016/679/oj/eng))

- [ ] **PRC-24-001** — Maintain an inventory of personal and sensitive data.
- [ ] **PRC-24-002** — Record the source, purpose, lawful basis or authorization, recipients, location, retention, and deletion method for each category.
- [ ] **PRC-24-003** — Collect only data necessary for the documented purpose.
- [ ] **PRC-24-004** — Privacy-protective defaults are used.
- [ ] **PRC-24-005** — Product behavior matches the published privacy notice.
- [ ] **PRC-24-006** — New data uses have undergone privacy review.
- [ ] **PRC-24-007** — Consent is used only where appropriate.
- [ ] **PRC-24-008** — Consent is specific, informed, granular, recorded, and withdrawable.
- [ ] **PRC-24-009** — Refusing nonessential consent does not improperly block the service.
- [ ] **PRC-24-010** — Nonessential trackers are prevented from operating before valid permission where required.
- [ ] **PRC-24-011** — Cookie and tracker inventories match actual production behavior.
- [ ] **PRC-24-012** — Analytics, advertising, attribution, session replay, and experimentation tools receive only approved data.
- [ ] **PRC-24-013** — Sensitive fields are redacted from logs, monitoring, support tools, recordings, and analytics.
- [ ] **PRC-24-014** — Access, correction, deletion, restriction, objection, and portability workflows exist where applicable.
- [ ] **PRC-24-015** — Rights requests cover downstream systems and processors.
- [ ] **PRC-24-016** — Requester identity verification is proportionate and does not collect excessive additional data.
- [ ] **PRC-24-017** — Rights-response deadlines are tracked.
- [ ] **PRC-24-018** — Retention periods are explicit.
- [ ] **PRC-24-019** — Deletion occurs across primary data, caches, indexes, derived data, and downstream systems.
- [ ] **PRC-24-020** — Backup deletion behavior is documented and legally acceptable.
- [ ] **PRC-24-021** — Legal holds override deletion only through a controlled process.
- [ ] **PRC-24-022** — Pseudonymization or anonymization claims have been validated.
- [ ] **PRC-24-023** — Reidentification risk has been assessed.
- [ ] **PRC-24-024** — Data export does not expose other users’ or tenants’ information.
- [ ] **PRC-24-025** — High-risk processing has an appropriate privacy-impact assessment.
- [ ] **PRC-24-026** — Children’s, biometric, health, precise-location, communications, financial, and other sensitive data trigger enhanced review.
- [ ] **PRC-24-027** — Cross-border transfers and data-residency commitments are documented.
- [ ] **PRC-24-028** — Processor and subprocessor contracts cover permitted use, security, deletion, assistance, and breach notification.
- [ ] **PRC-24-029** — Privacy contacts and escalation paths exist.
- [ ] **PRC-24-030** — A personal-data incident can be identified, scoped, and reported within applicable deadlines.
- [ ] **PRC-24-031** — Production support access to personal data is minimized and audited.
- [ ] **PRC-24-032** — Synthetic or masked data is used outside production.
- [ ] **PRC-24-033** — Test accounts and deleted users do not remain indefinitely in downstream tools.
- [ ] **PRC-24-034** — Interfaces avoid manipulative consent, cancellation, or privacy choices.

---

## 25. Performance and user experience efficiency

For public web experiences, current Core Web Vitals guidance uses LCP of at most 2.5 seconds, INP of at most 200 milliseconds, and CLS of at most 0.1 at the 75th percentile, segmented by mobile and desktop. These are useful defaults, not substitutes for product-specific performance objectives. ([web.dev](https://web.dev/articles/vitals))

- [ ] **PRC-25-001** — Performance objectives are defined per critical user journey.
- [ ] **PRC-25-002** — Objectives use relevant percentiles, not averages alone.
- [ ] **PRC-25-003** — Frontend and backend budgets are documented.
- [ ] **PRC-25-004** — Baseline performance is recorded before launch.
- [ ] **PRC-25-005** — Laboratory tests and real-user measurements are both used where applicable.
- [ ] **PRC-25-006** — Public pages meet approved loading, interaction, and visual-stability targets.
- [ ] **PRC-25-007** — API latency targets are defined per endpoint or journey.
- [ ] **PRC-25-008** — Database-query latency and volume are measured.
- [ ] **PRC-25-009** — Excessive query patterns and repeated fetches have been eliminated.
- [ ] **PRC-25-010** — Payload sizes are bounded.
- [ ] **PRC-25-011** — Compression is enabled appropriately.
- [ ] **PRC-25-012** — Images, fonts, scripts, and styles are delivered efficiently.
- [ ] **PRC-25-013** — Critical resources are prioritized.
- [ ] **PRC-25-014** — Noncritical work is deferred without harming accessibility or correctness.
- [ ] **PRC-25-015** — Cache and CDN behavior has been validated.
- [ ] **PRC-25-016** — Cold-start performance is acceptable.
- [ ] **PRC-25-017** — Connection and thread pools are sized and monitored.
- [ ] **PRC-25-018** — Memory, CPU, file-descriptor, connection, and storage leaks have been investigated.
- [ ] **PRC-25-019** — Long-running or resource-intensive requests have limits.
- [ ] **PRC-25-020** — Performance is tested with realistic authentication, authorization, data volume, and cache state.
- [ ] **PRC-25-021** — Performance is tested on representative mobile devices and slower networks where relevant.
- [ ] **PRC-25-022** — Third-party script and service impact is measured.
- [ ] **PRC-25-023** — Performance degradation during partial dependency failure is acceptable.
- [ ] **PRC-25-024** — Performance regressions have automated release thresholds.
- [ ] **PRC-25-025** — Client telemetry does not itself materially harm performance.
- [ ] **PRC-25-026** — Performance objectives are connected to user and business outcomes.

---

## 26. Capacity, scalability, and overload control

- [ ] **PRC-26-001** — Forecast normal, peak, burst, launch, seasonal, and abuse traffic.
- [ ] **PRC-26-002** — Forecast data, file, index, queue, log, and backup growth.
- [ ] **PRC-26-003** — Identify the maximum safe capacity of each critical component.
- [ ] **PRC-26-004** — Load tests use realistic request mixes and data distributions.
- [ ] **PRC-26-005** — Load tests include authentication and authorization overhead.
- [ ] **PRC-26-006** — Load tests include external dependency behavior.
- [ ] **PRC-26-007** — Stress testing identifies how the system fails beyond capacity.
- [ ] **PRC-26-008** — Spike testing covers sudden surges.
- [ ] **PRC-26-009** — Endurance testing identifies leaks and gradual degradation.
- [ ] **PRC-26-010** — Capacity testing includes the largest expected tenant or customer.
- [ ] **PRC-26-011** — Provider quotas and account limits have been reviewed.
- [ ] **PRC-26-012** — Database connections, locks, storage, replicas, and I/O have sufficient headroom.
- [ ] **PRC-26-013** — Queue, cache, load-balancer, network, and worker capacity have sufficient headroom.
- [ ] **PRC-26-014** — Autoscaling minimums, maximums, triggers, cooldowns, and scale-down behavior are validated.
- [ ] **PRC-26-015** — Scaling does not overload a database or dependency.
- [ ] **PRC-26-016** — Capacity remains sufficient during a zone, node, or region failure according to the design.
- [ ] **PRC-26-017** — Rate limits protect shared capacity fairly.
- [ ] **PRC-26-018** — Per-user and per-tenant quotas prevent noisy-neighbor failure.
- [ ] **PRC-26-019** — Backpressure propagates safely.
- [ ] **PRC-26-020** — Queues are bounded or have explicit overload behavior.
- [ ] **PRC-26-021** — Load shedding preserves the most important operations.
- [ ] **PRC-26-022** — Retry storms are prevented.
- [ ] **PRC-26-023** — Expensive optional functionality can degrade or disable safely.
- [ ] **PRC-26-024** — DDoS scenarios and provider protections have been reviewed.
- [ ] **PRC-26-025** — Capacity-related alert thresholds provide enough time to act.
- [ ] **PRC-26-026** — The cost of expected and peak operation has been modeled.
- [ ] **PRC-26-027** — Cost limits do not silently terminate essential operation.
- [ ] **PRC-26-028** — Capacity assumptions have owners and review dates.

---
