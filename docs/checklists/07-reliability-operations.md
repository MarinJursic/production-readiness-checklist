# Reliability and operations

> Prepare infrastructure, observability, response, recovery, deployment, and post-launch verification.

Sections 27–35 of the master checklist. For each applicable item, capture status, owner, evidence, release, environment, reviewer, and evidence age.

## 27. Reliability, resilience, and failure engineering

SRE practice centers reliability decisions on user-facing SLIs, approved SLOs, and an error-budget policy. It also warns that apparently redundant systems frequently share hidden dependencies and failure domains. ([sre.google](https://sre.google/workbook/implementing-slos/))

- [ ] **PRC-27-001** — User-centric service-level indicators are defined.
- [ ] **PRC-27-002** — SLOs cover the most important availability, latency, correctness, durability, freshness, and quality outcomes.
- [ ] **PRC-27-003** — SLOs are approved by product and engineering stakeholders.
- [ ] **PRC-27-004** — SLO measurement represents actual user experience.
- [ ] **PRC-27-005** — An error-budget policy controls release and reliability decisions.
- [ ] **PRC-27-006** — The service is within its permitted error budget before a risky launch.
- [ ] **PRC-27-007** — Critical dependencies have reliability expectations compatible with the product SLO.
- [ ] **PRC-27-008** — Dependency reliability limitations are engineered around where necessary.
- [ ] **PRC-27-009** — Every network call has an appropriate timeout.
- [ ] **PRC-27-010** — Retries are bounded.
- [ ] **PRC-27-011** — Retries use exponential backoff and jitter where appropriate.
- [ ] **PRC-27-012** — Retried operations are safe or idempotent.
- [ ] **PRC-27-013** — Circuit breakers or equivalent controls prevent cascading failure.
- [ ] **PRC-27-014** — Bulkheads or isolation prevent one workload from exhausting all capacity.
- [ ] **PRC-27-015** — Fallback behavior is defined and tested.
- [ ] **PRC-27-016** — Graceful degradation prioritizes critical functionality.
- [ ] **PRC-27-017** — Failure in optional functionality does not unnecessarily fail the whole product.
- [ ] **PRC-27-018** — Health checks distinguish process liveness from readiness to serve.
- [ ] **PRC-27-019** — Instances stop receiving work before shutdown.
- [ ] **PRC-27-020** — In-flight work is completed, transferred, or safely retried during shutdown.
- [ ] **PRC-27-021** — Redundancy crosses the failure domains it is intended to survive.
- [ ] **PRC-27-022** — Shared control planes, identity systems, networks, storage, DNS, and providers are included in failure analysis.
- [ ] **PRC-27-023** — Single points of failure are removed or explicitly accepted.
- [ ] **PRC-27-024** — Network delay, packet loss, partitions, and DNS failure have been tested.
- [ ] **PRC-27-025** — Dependency timeout, malformed response, throttling, and outage have been tested.
- [ ] **PRC-27-026** — Cache outage and cache stampede behavior have been tested.
- [ ] **PRC-27-027** — Database failover and replication lag have been tested.
- [ ] **PRC-27-028** — Queue backlog and broker failure have been tested.
- [ ] **PRC-27-029** — Disk-full, memory pressure, CPU exhaustion, file-descriptor exhaustion, and connection exhaustion have been tested.
- [ ] **PRC-27-030** — Clock skew and time-service failure have been considered.
- [ ] **PRC-27-031** — Certificate and key expiry have been considered.
- [ ] **PRC-27-032** — Scheduled jobs handle overlapping or missed execution.
- [ ] **PRC-27-033** — Stateful failover preserves required integrity.
- [ ] **PRC-27-034** — Recovery does not create uncontrolled retry or traffic surges.
- [ ] **PRC-27-035** — Maintenance and degraded modes are documented.
- [ ] **PRC-27-036** — Controlled fault-injection or equivalent resilience tests have been completed for high-risk paths.
- [ ] **PRC-27-037** — Recovery time from tested failures meets the approved objective.

---

## 28. Infrastructure and platform readiness

- [ ] **PRC-28-001** — Infrastructure is defined as code or otherwise reproducible and auditable.
- [ ] **PRC-28-002** — Infrastructure changes receive review.
- [ ] **PRC-28-003** — Infrastructure code is tested and security-scanned.
- [ ] **PRC-28-004** — Production accounts, subscriptions, projects, and clusters are separated appropriately.
- [ ] **PRC-28-005** — Resource ownership and criticality are labeled.
- [ ] **PRC-28-006** — Operating systems, runtimes, database engines, and platform versions are supported.
- [ ] **PRC-28-007** — Security patches are current according to risk policy.
- [ ] **PRC-28-008** — Unnecessary services, packages, ports, and capabilities are removed.
- [ ] **PRC-28-009** — Workloads run with the least operating-system and platform privilege.
- [ ] **PRC-28-010** — Administrative access uses strong authentication and controlled entry points.
- [ ] **PRC-28-011** — Just-in-time or time-bounded administrative access is used where practical.
- [ ] **PRC-28-012** — Network segmentation and private connectivity are configured as designed.
- [ ] **PRC-28-013** — Storage encryption and access policies are correct.
- [ ] **PRC-28-014** — Critical resources have deletion protection or equivalent safeguards.
- [ ] **PRC-28-015** — Resource quotas and account limits are documented.
- [ ] **PRC-28-016** — Monitoring covers platform control-plane and quota failures.
- [ ] **PRC-28-017** — Configuration drift and unauthorized resources are detected.
- [ ] **PRC-28-018** — Asset inventory matches deployed infrastructure.
- [ ] **PRC-28-019** — Production resource deletion and recreation have tested recovery procedures.
- [ ] **PRC-28-020** — Multi-zone or multi-region architecture matches the required recovery and availability objectives.
- [ ] **PRC-28-021** — Infrastructure capacity has been validated under failure.
- [ ] **PRC-28-022** — Cloud or platform provider status and escalation paths are integrated into operations.
- [ ] **PRC-28-023** — Infrastructure cost alerts and runaway-resource protections exist.
- [ ] **PRC-28-024** — Break-glass access works during identity or network failure and is audited.
- [ ] **PRC-28-025** — Infrastructure credentials and machine identities are rotated and inventoried.

### 28.1 Containers and orchestration, when used

- [ ] **PRC-28-026** — Base images are minimal, trusted, pinned, and scanned.
- [ ] **PRC-28-027** — Workloads do not run as root unless justified.
- [ ] **PRC-28-028** — Filesystems are read-only where practical.
- [ ] **PRC-28-029** — Unnecessary Linux or equivalent capabilities are removed.
- [ ] **PRC-28-030** — Privileged containers and host access are prohibited unless explicitly approved.
- [ ] **PRC-28-031** — Resource requests and limits are configured.
- [ ] **PRC-28-032** — Runtime security policies are enforced.
- [ ] **PRC-28-033** — Platform role-based access follows least privilege.
- [ ] **PRC-28-034** — Network policies restrict unnecessary communication.
- [ ] **PRC-28-035** — Admission or policy checks prevent prohibited workloads.
- [ ] **PRC-28-036** — Image signatures or provenance are verified.
- [ ] **PRC-28-037** — Secrets are not embedded in images.
- [ ] **PRC-28-038** — Node, cluster, and control-plane upgrades have been rehearsed.
- [ ] **PRC-28-039** — Eviction, rescheduling, and rolling-update behavior are tested.
- [ ] **PRC-28-040** — Persistent-volume failure and restoration are tested.

### 28.2 Serverless or managed execution, when used

- [ ] **PRC-28-041** — Concurrency limits and scaling behavior are understood.
- [ ] **PRC-28-042** — Cold-start impact is acceptable.
- [ ] **PRC-28-043** — Execution duration, memory, payload, temporary-storage, and connection limits are handled.
- [ ] **PRC-28-044** — Event retries and duplicate delivery are safe.
- [ ] **PRC-28-045** — Provider lock-in and recovery dependencies are documented.
- [ ] **PRC-28-046** — Service identities are scoped per function or workload where practical.
- [ ] **PRC-28-047** — Provider logging and tracing are sufficient for incident investigation.
- [ ] **PRC-28-048** — Regional availability and failover match the SLO.

---

## 29. Observability, logging, metrics, and audit trails

- [ ] **PRC-29-001** — Logs, metrics, traces, synthetics, and real-user signals cover the critical journeys.
- [ ] **PRC-29-002** — Telemetry uses a consistent service, environment, version, and instance identity.
- [ ] **PRC-29-003** — Timestamps are accurate and use an unambiguous time basis.
- [ ] **PRC-29-004** — Requests and background operations have correlation identifiers.
- [ ] **PRC-29-005** — Distributed trace context propagates across trusted service boundaries.
- [ ] **PRC-29-006** — The release version appears in telemetry.
- [ ] **PRC-29-007** — Deployments and configuration changes are annotated on dashboards.
- [ ] **PRC-29-008** — Availability is measured from the user’s perspective.
- [ ] **PRC-29-009** — Latency is measured with meaningful percentiles.
- [ ] **PRC-29-010** — Errors distinguish user mistakes, expected business outcomes, dependency failures, and system defects.
- [ ] **PRC-29-011** — Traffic or workload volume is measured.
- [ ] **PRC-29-012** — Saturation is measured for relevant resources.
- [ ] **PRC-29-013** — Correctness, freshness, durability, queue lag, and reconciliation metrics exist where relevant.
- [ ] **PRC-29-014** — Business indicators can reveal silent technical failure.
- [ ] **PRC-29-015** — Dependency health is observable.
- [ ] **PRC-29-016** — Dashboards support both high-level triage and detailed diagnosis.
- [ ] **PRC-29-017** — Logs are structured and searchable.
- [ ] **PRC-29-018** — Log levels are consistent.
- [ ] **PRC-29-019** — Errors include enough context to diagnose without exposing sensitive data.
- [ ] **PRC-29-020** — Passwords, reusable tokens, session identifiers, secret keys, and payment credentials are never logged.
- [ ] **PRC-29-021** — Personal and confidential data is excluded or minimized.
- [ ] **PRC-29-022** — Log-injection and control-character attacks are handled.
- [ ] **PRC-29-023** — Audit records identify actor, tenant, operation, target, time, source, outcome, and reason where appropriate.
- [ ] **PRC-29-024** — Authentication, authorization changes, administrative actions, sensitive reads, exports, deletion, and configuration changes are audited.
- [ ] **PRC-29-025** — Audit data is protected against unauthorized modification and deletion.
- [ ] **PRC-29-026** — Audit access is restricted and audited.
- [ ] **PRC-29-027** — Log and audit retention matches security, privacy, legal, and operational requirements.
- [ ] **PRC-29-028** — Sampling does not discard critical security or correctness events.
- [ ] **PRC-29-029** — Telemetry pipeline failure is itself monitored.
- [ ] **PRC-29-030** — Telemetry backlog and quota exhaustion are monitored.
- [ ] **PRC-29-031** — Synthetic checks validate critical journeys externally.
- [ ] **PRC-29-032** — Real-user monitoring is privacy-reviewed.
- [ ] **PRC-29-033** — An investigator can locate relevant evidence within the incident-response objective.
- [ ] **PRC-29-034** — Observability cost and cardinality are controlled without losing essential signals.
- [ ] **PRC-29-035** — Dashboards and queries have been tested during a simulated incident.

---

## 30. Alerting, on-call, and runbooks

- [ ] **PRC-30-001** — Alerts are tied to user impact, SLO risk, data integrity, security, or impending capacity failure.
- [ ] **PRC-30-002** — Paging alerts are actionable.
- [ ] **PRC-30-003** — Non-actionable information is sent through non-paging channels.
- [ ] **PRC-30-004** — Alert severity definitions are documented.
- [ ] **PRC-30-005** — Every alert has an owner.
- [ ] **PRC-30-006** — Every page has a corresponding diagnostic or mitigation runbook.
- [ ] **PRC-30-007** — Alert routing and escalation have been tested.
- [ ] **PRC-30-008** — Primary and secondary coverage exist during required service hours.
- [ ] **PRC-30-009** — On-call responders have the required access and tools.
- [ ] **PRC-30-010** — No alert depends on one irreplaceable person.
- [ ] **PRC-30-011** — Burn-rate or equivalent alerts detect both fast and slow SLO consumption.
- [ ] **PRC-30-012** — Alert thresholds avoid obvious noise and missed incidents.
- [ ] **PRC-30-013** — Duplicate alerts are correlated or suppressed appropriately.
- [ ] **PRC-30-014** — Maintenance and deployments do not suppress unrelated incidents.
- [ ] **PRC-30-015** — Failure of the monitoring or paging system generates an independent signal.
- [ ] **PRC-30-016** — Runbooks include symptoms, impact, diagnosis, containment, rollback, recovery, verification, and escalation.
- [ ] **PRC-30-017** — Runbooks avoid undocumented tribal knowledge.
- [ ] **PRC-30-018** — Runbook commands are safe, reviewed, and clearly distinguish destructive actions.
- [ ] **PRC-30-019** — Access to required dashboards and consoles works during an incident.
- [ ] **PRC-30-020** — Provider support and escalation information is current.
- [ ] **PRC-30-021** — Customer support knows how to escalate production problems.
- [ ] **PRC-30-022** — A public or customer-facing status mechanism exists where appropriate.
- [ ] **PRC-30-023** — On-call load and alert volume are operationally sustainable.
- [ ] **PRC-30-024** — Handover procedures preserve active incident context.
- [ ] **PRC-30-025** — On-call and runbook drills have been completed.

---

## 31. Security validation and vulnerability management

- [ ] **PRC-31-001** — Security requirements are traceable to the threat model and recognized verification criteria.
- [ ] **PRC-31-002** — The release is assessed against the applicable OWASP ASVS 5.0.0 requirements.
- [ ] **PRC-31-003** — Source-code security analysis has run.
- [ ] **PRC-31-004** — Dependency and SBOM analysis has run.
- [ ] **PRC-31-005** — Secret scanning has run.
- [ ] **PRC-31-006** — Infrastructure and configuration analysis has run.
- [ ] **PRC-31-007** — Container or system-image analysis has run where applicable.
- [ ] **PRC-31-008** — Deployed application security testing has run.
- [ ] **PRC-31-009** — API security testing has run.
- [ ] **PRC-31-010** — Authentication and recovery testing has run.
- [ ] **PRC-31-011** — Authorization and tenant-isolation testing has run.
- [ ] **PRC-31-012** — Injection and unsafe-input testing has run.
- [ ] **PRC-31-013** — Browser and session security testing has run.
- [ ] **PRC-31-014** — File-processing security testing has run where applicable.
- [ ] **PRC-31-015** — Business-logic and abuse testing has run.
- [ ] **PRC-31-016** — Manual code review covers the highest-risk components.
- [ ] **PRC-31-017** — Independent penetration testing is complete where warranted.
- [ ] **PRC-31-018** — Penetration testing includes authenticated users, multiple roles, multiple tenants, APIs, administrative tools, and business logic as applicable.
- [ ] **PRC-31-019** — False positives are documented with evidence.
- [ ] **PRC-31-020** — Findings have owners and remediation targets.
- [ ] **PRC-31-021** — Fixed findings have been retested.
- [ ] **PRC-31-022** — Relevant known-exploited vulnerabilities are treated as priority blockers unless exposure and mitigation are convincingly demonstrated.
- [ ] **PRC-31-023** — Unsupported software and unpatchable components have explicit treatment plans.
- [ ] **PRC-31-024** — Compensating controls are validated in the deployed environment.
- [ ] **PRC-31-025** — A safe emergency-patching path exists.
- [ ] **PRC-31-026** — External attack-surface discovery is continuous.
- [ ] **PRC-31-027** — A vulnerability-disclosure policy and reporting channel exist.
- [ ] **PRC-31-028** — Security reports receive acknowledgment, triage, remediation, disclosure, and credit processes as appropriate.
- [ ] **PRC-31-029** — Customer security advisories can be produced promptly.
- [ ] **PRC-31-030** — Security evidence is retained with the release.
- [ ] **PRC-31-031** — Security regressions are represented in automated tests where practical.
- [ ] **PRC-31-032** — Production security scanning is configured not to damage service or data.
- [ ] **PRC-31-033** — Security monitoring covers account takeover, privilege changes, secret use, suspicious exports, abuse, and control tampering.

---

## 32. Incident response and crisis readiness

NIST SP 800-61 Revision 3, finalized in April 2025, treats incident response as part of broader risk management and emphasizes preparation, detection, response, and recovery. ([csrc.nist.gov](https://csrc.nist.gov/pubs/sp/800/61/r3/final))

- [ ] **PRC-32-001** — A current incident-response plan exists.
- [ ] **PRC-32-002** — Incident severity levels are defined.
- [ ] **PRC-32-003** — Authority to declare an incident is clear.
- [ ] **PRC-32-004** — Incident roles include command, technical response, communications, documentation, security, privacy/legal, and business ownership as needed.
- [ ] **PRC-32-005** — Primary and alternate communication channels exist.
- [ ] **PRC-32-006** — Contact details are current and available during identity, network, or provider failure.
- [ ] **PRC-32-007** — Procedures cover detection, triage, containment, eradication, recovery, and closure.
- [ ] **PRC-32-008** — Procedures distinguish outages, security incidents, privacy incidents, safety events, and fraud events.
- [ ] **PRC-32-009** — Evidence-preservation and forensic procedures exist.
- [ ] **PRC-32-010** — Logging and time synchronization support investigation.
- [ ] **PRC-32-011** — Credentials, sessions, certificates, and keys can be revoked or rotated rapidly.
- [ ] **PRC-32-012** — Compromised builds and dependencies can be identified and replaced.
- [ ] **PRC-32-013** — Customer and public communications templates exist.
- [ ] **PRC-32-014** — A status-page process exists where appropriate.
- [ ] **PRC-32-015** — Regulatory, contractual, insurer, partner, and law-enforcement notification requirements are mapped.
- [ ] **PRC-32-016** — Notification deadlines and decision owners are documented.
- [ ] **PRC-32-017** — Provider and subprocessor incident contacts are current.
- [ ] **PRC-32-018** — Alternate operations are possible during provider failure.
- [ ] **PRC-32-019** — Tabletop exercises cover realistic scenarios.
- [ ] **PRC-32-020** — Exercises include account takeover, data exfiltration, destructive access, ransomware, dependency compromise, regional outage, and data corruption as relevant.
- [ ] **PRC-32-021** — Exercise findings have owners and deadlines.
- [ ] **PRC-32-022** — Post-incident reviews are blameless but accountable.
- [ ] **PRC-32-023** — Corrective actions are tracked to completion.
- [ ] **PRC-32-024** — Lessons are converted into tests, controls, runbooks, and architecture changes.
- [ ] **PRC-32-025** — Incident metrics include detection, acknowledgment, containment, recovery, impact, and recurrence.
- [ ] **PRC-32-026** — Incident records preserve decisions and timelines.
- [ ] **PRC-32-027** — The plan can operate when the normal identity, chat, email, documentation, or cloud platform is unavailable.

---

## 33. Backups, disaster recovery, and continuity

- [ ] **PRC-33-001** — Business-impact analysis identifies critical services and data.
- [ ] **PRC-33-002** — Recovery Point Objective is defined for each stateful system.
- [ ] **PRC-33-003** — Recovery Time Objective is defined for each critical journey.
- [ ] **PRC-33-004** — Backup scope includes databases, files, queues where needed, configuration, infrastructure definitions, secrets or recovery mechanisms, certificates, and required artifacts.
- [ ] **PRC-33-005** — Backups run automatically at the required frequency.
- [ ] **PRC-33-006** — Backup success is monitored.
- [ ] **PRC-33-007** — Backups are encrypted.
- [ ] **PRC-33-008** — Backup access follows least privilege.
- [ ] **PRC-33-009** — Backup copies are isolated from ordinary production compromise.
- [ ] **PRC-33-010** — Immutable or otherwise tamper-resistant copies exist where risk warrants.
- [ ] **PRC-33-011** — Retention meets operational, legal, contractual, and privacy requirements.
- [ ] **PRC-33-012** — Backup keys and recovery credentials are protected.
- [ ] **PRC-33-013** — Point-in-time recovery is available where required.
- [ ] **PRC-33-014** — Full restoration has been tested.
- [ ] **PRC-33-015** — Partial restoration has been tested.
- [ ] **PRC-33-016** — Restoration to a clean environment has been tested.
- [ ] **PRC-33-017** — Regional or provider-loss recovery has been tested when part of the design.
- [ ] **PRC-33-018** — Restore order and dependency requirements are documented.
- [ ] **PRC-33-019** — DNS, certificates, secrets, identity, networking, and external integrations are included in the recovery procedure.
- [ ] **PRC-33-020** — Restored data is checked for integrity and reconciliation.
- [ ] **PRC-33-021** — Actual tested recovery time meets the objective.
- [ ] **PRC-33-022** — Failover and failback are both tested.
- [ ] **PRC-33-023** — Disaster-recovery capacity can carry the required load.
- [ ] **PRC-33-024** — Recovery does not accidentally send duplicate notifications, charges, webhooks, or jobs.
- [ ] **PRC-33-025** — Backup restoration does not re-enable revoked users, secrets, or deleted data without detection.
- [ ] **PRC-33-026** — Privacy deletion and backup-retention interactions are documented.
- [ ] **PRC-33-027** — A ransomware or destructive-administrator scenario has been exercised.
- [ ] **PRC-33-028** — Recovery does not depend on a single person or unavailable device.
- [ ] **PRC-33-029** — Continuity procedures cover customer support, communications, and essential manual operations.
- [ ] **PRC-33-030** — Recovery tests produce retained evidence and corrective actions.

---

## 34. Deployment and release engineering

- [ ] **PRC-34-001** — The deployment plan identifies exact artifacts and target environments.
- [ ] **PRC-34-002** — The plan identifies migrations, configuration changes, infrastructure changes, flags, and external dependencies.
- [ ] **PRC-34-003** — Required approvals are complete.
- [ ] **PRC-34-004** — Deployment permissions are restricted.
- [ ] **PRC-34-005** — The person deploying cannot silently substitute an unapproved artifact.
- [ ] **PRC-34-006** — Predeployment checks validate environment health and capacity.
- [ ] **PRC-34-007** — Required backups or recovery points exist.
- [ ] **PRC-34-008** — Database migration steps and ordering are explicit.
- [ ] **PRC-34-009** — Application, worker, scheduled-job, and schema compatibility is maintained during rollout.
- [ ] **PRC-34-010** — Mixed-version operation is tested where it will occur.
- [ ] **PRC-34-011** — Deployment uses a progressive, canary, blue-green, cohort, or similarly controlled approach where risk warrants.
- [ ] **PRC-34-012** — The smallest safe initial cohort is used.
- [ ] **PRC-34-013** — Automated health gates examine user-facing success, latency, errors, saturation, and business integrity.
- [ ] **PRC-34-014** — Canary and control populations are comparable.
- [ ] **PRC-34-015** — Rollout pauses automatically or operationally on defined signals.
- [ ] **PRC-34-016** — Stop and rollback criteria are documented before launch.
- [ ] **PRC-34-017** — Rollback can be initiated quickly.
- [ ] **PRC-34-018** — Rollback has been rehearsed.
- [ ] **PRC-34-019** — Rollback is compatible with database and event changes.
- [ ] **PRC-34-020** — When rollback is unsafe, roll-forward recovery is rehearsed.
- [ ] **PRC-34-021** — The previous artifact and configuration remain available.
- [ ] **PRC-34-022** — Feature flags provide an independent way to disable risky behavior where appropriate.
- [ ] **PRC-34-023** — Traffic shifting and load-balancer behavior are tested.
- [ ] **PRC-34-024** — Cache and CDN invalidation are planned.
- [ ] **PRC-34-025** — DNS changes and TTL timing are planned.
- [ ] **PRC-34-026** — Certificate changes are planned and monitored.
- [ ] **PRC-34-027** — External providers are ready for the expected change.
- [ ] **PRC-34-028** — Rate limits and quotas are increased where needed.
- [ ] **PRC-34-029** — Scheduled jobs and workers are deployed in a safe order.
- [ ] **PRC-34-030** — No critical deployment step exists only in someone’s memory.
- [ ] **PRC-34-031** — Manual steps are minimized, scripted where practical, and peer-checked.
- [ ] **PRC-34-032** — The complete procedure has been rehearsed in a production-like environment.
- [ ] **PRC-34-033** — On-call engineering, operations, security, data, and provider contacts are available as needed.
- [ ] **PRC-34-034** — Customer support and communications are ready.
- [ ] **PRC-34-035** — The change window leaves adequate observation and rollback time.
- [ ] **PRC-34-036** — Every deployment action is auditable.

---

## 35. Post-deployment verification

- [ ] **PRC-35-001** — Verify the deployed artifact digest.
- [ ] **PRC-35-002** — Verify the deployed configuration and feature-flag state.
- [ ] **PRC-35-003** — Verify database schema and migration completion.
- [ ] **PRC-35-004** — Verify migration row counts, totals, constraints, checksums, and business invariants.
- [ ] **PRC-35-005** — Run production-safe smoke tests for every critical journey.
- [ ] **PRC-35-006** — Test login, logout, recovery, and privileged authentication.
- [ ] **PRC-35-007** — Test authorization using multiple roles and tenants.
- [ ] **PRC-35-008** — Verify public, authenticated, administrative, API, worker, and scheduled-job paths.
- [ ] **PRC-35-009** — Verify third-party payments, webhooks, email, SMS, identity, storage, and other critical integrations.
- [ ] **PRC-35-010** — Verify TLS, certificates, redirects, CORS, cache controls, and security headers.
- [ ] **PRC-35-011** — Verify CDN and cache behavior.
- [ ] **PRC-35-012** — Review errors, latency, traffic, saturation, queues, database health, and dependency metrics.
- [ ] **PRC-35-013** — Review logs and traces for unexpected exceptions or data leakage.
- [ ] **PRC-35-014** — Verify alerts and synthetic monitoring.
- [ ] **PRC-35-015** — Compare canary and control outcomes.
- [ ] **PRC-35-016** — Verify analytics and business indicators.
- [ ] **PRC-35-017** — Perform an accessibility spot check on changed critical journeys.
- [ ] **PRC-35-018** — Confirm customer support receives expected telemetry and documentation.
- [ ] **PRC-35-019** — Confirm no unexpected increase in abuse, fraud, failed payments, duplicate operations, or support volume.
- [ ] **PRC-35-020** — Maintain heightened monitoring for the agreed observation period.
- [ ] **PRC-35-021** — Stop expansion when predefined thresholds are crossed.
- [ ] **PRC-35-022** — Roll back or declare an incident rather than rationalizing unexplained anomalies.
- [ ] **PRC-35-023** — Record the final go, hold, rollback, or incident decision.
- [ ] **PRC-35-024** — Schedule follow-up review for deferred observations and temporary controls.

---
