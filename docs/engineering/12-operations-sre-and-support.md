# Operations, SRE, and support

_Phase 12 of 16 in the [complete engineering review](00-overview.md)._

Infrastructure, SLOs, observability, alerting, on-call, incidents, continuity, capacity, FinOps, support, and retirement.

Assess every control as **Pass**, **Fail**, **Blocked**, or **Not Applicable**. A Pass needs current evidence for the exact product, revision, artifact, configuration, environment, and data state under review.

## Infrastructure and Platform Readiness

_Consolidated from `quality standards/13-operations-sre/01-infrastructure-and-platform-readiness.md`; 21 non-duplicative controls._

### Universal controls

- [ ] **USEQ-B52808A9** — Infrastructure changes receive review, testing, and security analysis.
- [ ] **USEQ-2963750A** — Production accounts, subscriptions, projects, environments, and control planes are separated appropriately.
- [ ] **USEQ-A41434C9** — Unnecessary services, packages, ports, protocols, users, and capabilities are removed.
- [ ] **USEQ-64B98F44** — Workloads run with least operating-system and platform privilege.
- [ ] **USEQ-C18D735C** — Resource quotas, account limits, and provider limits are documented.
- [ ] **USEQ-FF663A86** — Monitoring covers platform control-plane, quota, capacity, and regional failures.
- [ ] **USEQ-92A27983** — Configuration drift, unauthorized resources, and shadow infrastructure are detected.
- [ ] **USEQ-7D17BE10** — Resource deletion and recreation have tested recovery procedures.
- [ ] **USEQ-3FBED686** — Multi-zone and multi-region architecture match approved availability and recovery objectives.
- [ ] **USEQ-57F42939** — Infrastructure capacity is validated under failure.
- [ ] **USEQ-7F1F8F9C** — Provider status and escalation paths are integrated into operations.
- [ ] **USEQ-A0E78A59** — Cost alerts and runaway-resource protections exist.
- [ ] **USEQ-C02FFCDE** — Infrastructure credentials and machine identities are inventoried, scoped, and rotated.

### Containers and orchestration, when used

- [ ] **USEQ-DA0AF8A9** — Base images are minimal, trusted, pinned, scanned, and supported.
- [ ] **USEQ-2B34F31F** — Unnecessary capabilities are removed.
- [ ] **USEQ-AE7EA05E** — Privileged execution and host access are prohibited unless explicitly approved.
- [ ] **USEQ-D99BEFB1** — Platform access follows least privilege.
- [ ] **USEQ-FF113370** — Node, cluster, and control-plane upgrades are rehearsed.
- [ ] **USEQ-4B41FF50** — Persistent-storage failure and restoration are tested.

### Serverless or managed execution, when used

- [ ] **USEQ-6EB9C058** — Service identities are scoped per workload where practical.
- [ ] **USEQ-3D0C3C21** — Regional availability and failover match approved objectives.

## Observability, Logging, Metrics, Traces, and Audit

_Consolidated from `quality standards/13-operations-sre/02-observability-logging-metrics-traces-and-audit.md`; 13 non-duplicative controls._

### Universal controls

- [ ] **USEQ-364C0C23** — Logs, metrics, traces, synthetic checks, and real-user signals cover critical journeys.
- [ ] **USEQ-15DAEE33** — Telemetry consistently identifies service, environment, release, instance, region, and tenant where appropriate.
- [ ] **USEQ-ACC724DF** — Trace context propagates across trusted service boundaries.
- [ ] **USEQ-9C8C69CA** — Latency uses meaningful percentiles.
- [ ] **USEQ-2CAC6BCD** — Traffic, workload, saturation, correctness, freshness, durability, queue lag, and reconciliation are measured where relevant.
- [ ] **USEQ-1172021F** — Logs are structured, searchable, and use consistent levels.
- [ ] **USEQ-ACC02B9C** — Passwords, reusable tokens, session identifiers, secret keys, payment credentials, and private keys are never logged.
- [ ] **USEQ-8605B669** — Log injection and unsafe control characters are handled.
- [ ] **USEQ-E25ED114** — Authentication, authorization changes, administrative actions, sensitive reads, exports, deletions, and configuration changes are audited.
- [ ] **USEQ-4CEDEF76** — Log and audit retention match operational, security, privacy, legal, and contractual requirements.
- [ ] **USEQ-B4212946** — Telemetry pipeline failure, backlog, cost, and quota exhaustion are monitored.
- [ ] **USEQ-5FEB8562** — Investigators can locate relevant evidence within the incident-response objective.
- [ ] **USEQ-CAC36B21** — Dashboards, alerts, and diagnostic queries are tested during simulated incidents.

## Alerting, On-Call, and Runbooks

_Consolidated from `quality standards/13-operations-sre/03-alerting-oncall-and-runbooks.md`; 15 non-duplicative controls._

### Universal controls

- [ ] **USEQ-2D2D4869** — Alerts are tied to user impact, SLO risk, data integrity, security, fraud, or impending capacity failure.
- [ ] **USEQ-768749B7** — Non-actionable information uses non-paging channels.
- [ ] **USEQ-402BD03D** — Every alert has an owner and a corresponding diagnostic or mitigation runbook.
- [ ] **USEQ-846BA894** — Alert routing and escalation are tested.
- [ ] **USEQ-467D6189** — Responders have required access and tools.
- [ ] **USEQ-57378001** — Fast and slow SLO consumption are detected.
- [ ] **USEQ-37471699** — Thresholds avoid obvious noise and missed incidents.
- [ ] **USEQ-06E3518C** — Failure of monitoring or paging generates an independent signal.
- [ ] **USEQ-4A5C9EF3** — Runbooks include symptoms, impact, diagnosis, containment, rollback or roll-forward, recovery, verification, and escalation.
- [ ] **USEQ-9B7E1084** — Runbooks do not rely on undocumented tribal knowledge.
- [ ] **USEQ-12D77066** — Runbook commands are reviewed and clearly identify destructive actions.
- [ ] **USEQ-11AF58B5** — Required dashboards, consoles, credentials, and documentation remain accessible during incidents.
- [ ] **USEQ-B43706EC** — A customer-facing status mechanism exists where appropriate.
- [ ] **USEQ-1D858AD2** — On-call load and alert volume are sustainable.
- [ ] **USEQ-9E774722** — On-call, alerting, and runbook drills have been completed.

## Backup, Disaster Recovery, and Business Continuity

_Consolidated from `quality standards/13-operations-sre/04-backup-disaster-recovery-and-business-continuity.md`; 17 non-duplicative controls._

### Universal controls

- [ ] **USEQ-004EF7C7** — Business-impact analysis identifies critical services, journeys, data, dependencies, and manual operations.
- [ ] **USEQ-F69270D3** — Recovery Point Objectives are defined for every stateful system.
- [ ] **USEQ-72E6B0CB** — Recovery Time Objectives are defined for every critical journey.
- [ ] **USEQ-0DF2FA43** — Backups run automatically at the required frequency and success is monitored.
- [ ] **USEQ-98D70D6C** — Backups are encrypted and access follows least privilege.
- [ ] **USEQ-76293243** — Immutable or tamper-resistant copies exist where risk warrants them.
- [ ] **USEQ-F6E2640E** — Full, partial, clean-environment, regional-loss, and provider-loss restoration are tested as applicable.
- [ ] **USEQ-8C6EF82F** — DNS, certificates, secrets, identity, networking, external integrations, and communication systems are included in recovery planning.
- [ ] **USEQ-7C161E16** — Restored data is checked for integrity and reconciled.
- [ ] **USEQ-D4E7515B** — Actual tested recovery time and recovery point meet approved objectives.
- [ ] **USEQ-C41BABC1** — Failover and failback are tested.
- [ ] **USEQ-0AAD5260** — Disaster-recovery capacity can carry required load.
- [ ] **USEQ-F4AB063C** — Recovery does not create duplicate notifications, charges, webhooks, messages, or jobs.
- [ ] **USEQ-FD45C5FC** — Restoration does not silently re-enable revoked users, revoked secrets, expired permissions, or deleted data.
- [ ] **USEQ-75188025** — Ransomware and destructive-administrator scenarios are exercised where relevant.
- [ ] **USEQ-DFA9E714** — Recovery does not depend on a single person, device, credential, or unavailable provider.
- [ ] **USEQ-AB371F8F** — Continuity procedures cover customer support, communications, essential manual operations, and decision authority.

## SLIs, SLOs, Error Budgets, and Reliability Governance

_Consolidated from `quality standards/13-operations-sre/05-sli-slo-error-budgets-and-reliability-governance.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-D1BD4768** — Identify the user journeys and service outcomes whose failure matters.
- [ ] **USEQ-B3E613B3** — Define service-level indicators from user-observable availability, latency, correctness, durability, freshness, and quality.
- [ ] **USEQ-0FC335A7** — Specify good and total events or equivalent valid measurement semantics.
- [ ] **USEQ-DD333F33** — Use windows and percentiles aligned to user and business consequences.
- [ ] **USEQ-14445940** — Set objectives from user need, dependency capability, cost, and risk rather than arbitrary perfection.
- [ ] **USEQ-B6C7F51A** — Measure from multiple points when server health can differ from user experience.
- [ ] **USEQ-A332E7D6** — Exclude events only through documented semantics that cannot hide provider-caused failure.
- [ ] **USEQ-918009AE** — Define an error-budget policy before budget exhaustion occurs.
- [ ] **USEQ-D85BEB1A** — Link budget status to release pace, change risk, reliability work, and escalation.
- [ ] **USEQ-E1B20661** — Use fast and slow burn alerts to detect acute and chronic objective risk.
- [ ] **USEQ-0616ED60** — Review whether objectives are too loose, too strict, or measuring the wrong outcome.
- [ ] **USEQ-77209F37** — Define dependencies' contribution to the end-to-end objective.
- [ ] **USEQ-4E57FA91** — Avoid multiplying optimistic component availability claims without shared-failure analysis.
- [ ] **USEQ-3F631903** — Publish objectives and current status to relevant product, engineering, support, and leadership stakeholders.
- [ ] **USEQ-C7355F4C** — Use SLO misses to prioritize systemic improvement, not punish individuals.
- [ ] **USEQ-A8A37E99** — Preserve separate integrity and safety gates that cannot be averaged into availability.
- [ ] **USEQ-5BB4035B** — Review objectives after material product, traffic, architecture, dependency, or contractual change.
- [ ] **USEQ-D72ACBC8** — Measure planned maintenance according to user and contractual semantics.
- [ ] **USEQ-F2601113** — Keep raw measurement and calculation reproducible.
- [ ] **USEQ-38C611A2** — Do not promise external service levels that internal evidence and capacity cannot support.

## Service Management and Customer Support

_Consolidated from `quality standards/13-operations-sre/06-service-management-and-customer-support.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-6D436E9F** — Maintain a service catalog with owner, purpose, users, dependencies, support hours, objectives, data, and lifecycle state.
- [ ] **USEQ-9AD25B90** — Define support channels, severity, response, escalation, entitlement, and communication expectations.
- [ ] **USEQ-48319180** — Provide support staff with current product behavior, known issues, diagnostics, privacy, security, and escalation guidance.
- [ ] **USEQ-7751FFCD** — Restrict and audit support access to customer data and administrative functions.
- [ ] **USEQ-EE8096B5** — Verify requester identity before disclosing or changing sensitive information.
- [ ] **USEQ-303A0806** — Separate service requests, defects, incidents, security events, privacy requests, and complaints while preserving linkage.
- [ ] **USEQ-8921C2A9** — Use consistent priority based on impact and urgency.
- [ ] **USEQ-8BD3AE1B** — Keep customers informed with accurate status, workaround, risk, and resolution information.
- [ ] **USEQ-52C655EE** — Avoid promising dates or causes before evidence supports them.
- [ ] **USEQ-7E1DF527** — Capture recurring requests and friction as product and documentation improvement signals.
- [ ] **USEQ-5CA932E3** — Maintain searchable knowledge articles with owner, effective version, review date, and retirement.
- [ ] **USEQ-320B9CC3** — Coordinate changes, releases, maintenance, suppliers, and support readiness.
- [ ] **USEQ-3847EF4D** — Measure resolution quality, reopen rate, escalation, user effort, satisfaction, and recurrence—not closure volume alone.
- [ ] **USEQ-17B8BFF4** — Provide accessible support channels and alternatives.
- [ ] **USEQ-6EEB4C4C** — Define handoff across time zones, teams, suppliers, and incident command.
- [ ] **USEQ-6509842E** — Protect support notes, attachments, recordings, and exports according to data sensitivity.
- [ ] **USEQ-23D14CF4** — Review complaints and vulnerable-user cases with appropriate expertise.
- [ ] **USEQ-BB51346B** — Define service acceptance before transfer from project to operation.
- [ ] **USEQ-DFA02FCC** — Retire support obligations only through communicated lifecycle policy.
- [ ] **USEQ-115E72C7** — Use continual improvement to reduce demand caused by product defects and confusing design.

## Incident, Problem Management, and Postmortems

_Consolidated from `quality standards/13-operations-sre/07-incident-problem-management-and-postmortems.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-F8D01DA0** — Define incident criteria, severity, authority, roles, communication, and escalation before incidents occur.
- [ ] **USEQ-FD4D66CA** — Prioritize protection of people, data integrity, security, and essential service over preserving change plans.
- [ ] **USEQ-A3F2248F** — Establish incident command and one authoritative timeline.
- [ ] **USEQ-C36750C1** — Separate command, operations, investigation, communication, and documentation roles when scale requires it.
- [ ] **USEQ-D18DD1C3** — Use clear hypotheses, evidence, decisions, owners, and timestamps during response.
- [ ] **USEQ-232D5A01** — Contain blast radius and stop harmful automation or change.
- [ ] **USEQ-D3AB5A8E** — Preserve evidence while restoring service.
- [ ] **USEQ-C6CD577E** — Communicate known facts, uncertainty, impact, mitigation, and next update without speculation.
- [ ] **USEQ-84C1C16C** — Declare recovery only after user journeys, integrity, backlog, dependencies, and monitoring are verified.
- [ ] **USEQ-B5D09F0D** — Monitor for recurrence and delayed effects.
- [ ] **USEQ-D2A2831E** — Link related incidents to a problem record when systemic cause or recurrence exists.
- [ ] **USEQ-864FDC25** — Analyze technical, process, interface, organizational, supplier, and incentive contributors.
- [ ] **USEQ-C87A0E16** — Avoid blame, hindsight bias, and a single human-error root cause.
- [ ] **USEQ-887E41C0** — Identify why controls did not prevent, detect, contain, or recover sooner.
- [ ] **USEQ-B45BF06B** — Assign corrective actions that are specific, prioritized, owned, funded, and verifiable.
- [ ] **USEQ-A0EC83FE** — Prefer actions that change systems and incentives over reminders alone.
- [ ] **USEQ-4EE52A7C** — Track actions to effective completion and review recurrence.
- [ ] **USEQ-FEA91FF6** — Share lessons with other systems that have the same failure class.
- [ ] **USEQ-7617216D** — Update tests, monitors, runbooks, architecture, requirements, training, and risk records.
- [ ] **USEQ-2805E81B** — Measure restoration, impact, detection, recurrence, action effectiveness, and learning—not incident count alone.

## Operational Capacity, Performance, and Efficiency

_Consolidated from `quality standards/13-operations-sre/08-operational-capacity-performance-and-efficiency.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-88305F87** — Maintain forecasts for ordinary, peak, seasonal, launch, growth, recovery, and abuse demand.
- [ ] **USEQ-1012D643** — Track capacity limits and quotas across compute, memory, storage, network, connections, queues, databases, third parties, and control planes.
- [ ] **USEQ-874742D9** — Measure saturation and headroom at the actual bottleneck.
- [ ] **USEQ-CD9DAD72** — Define minimum safe headroom for failure, deployment, backup, recovery, and traffic shift.
- [ ] **USEQ-D9DC2BF5** — Review forecasts against actual demand and update assumptions.
- [ ] **USEQ-7E0F4822** — Scale based on leading indicators where reactive scaling would be too late.
- [ ] **USEQ-0C1180BC** — Test scale-up, scale-down, rebalancing, and capacity during a failure domain loss.
- [ ] **USEQ-D9D69015** — Prevent autoscaling from overwhelming a database or external dependency.
- [ ] **USEQ-E699A7C6** — Use per-tenant, per-user, and per-operation quotas to protect shared service.
- [ ] **USEQ-46661625** — Define overload admission, queuing, backpressure, prioritization, degradation, and shedding.
- [ ] **USEQ-49019E6A** — Protect critical operations before optional or expensive operations.
- [ ] **USEQ-C8DBF731** — Monitor latency distributions, throughput, errors, queue age, resource use, and unit cost together.
- [ ] **USEQ-5E82CAD5** — Detect leaks and gradual degradation through endurance signals.
- [ ] **USEQ-9197F66F** — Plan storage, index, log, backup, and retention growth.
- [ ] **USEQ-5FE14FFF** — Coordinate capacity with supplier limits and contract lead times.
- [ ] **USEQ-36E4D361** — Keep emergency capacity actions documented, authorized, and reversible.
- [ ] **USEQ-A73F088E** — Verify capacity after architecture, data-shape, feature, dependency, or workload change.
- [ ] **USEQ-F46A2D61** — Use performance incidents to update models and tests.
- [ ] **USEQ-C40A674E** — Avoid maintaining wasteful idle capacity when resilience can be met more efficiently.
- [ ] **USEQ-8006F4CA** — Do not reduce headroom below recovery and SLO needs solely to meet short-term cost targets.

## Cost Efficiency and FinOps

_Consolidated from `quality standards/13-operations-sre/09-cost-efficiency-and-finops.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-CD017015** — Assign ownership for material technology cost and its business or user value.
- [ ] **USEQ-2AA3663E** — Tag or otherwise attribute resources, services, environments, teams, products, and tenants sufficiently for analysis.
- [ ] **USEQ-E8247F73** — Separate fixed, variable, shared, support, data, network, license, supplier, and incident costs.
- [ ] **USEQ-D39687B3** — Define unit costs aligned to meaningful product outcomes and workloads.
- [ ] **USEQ-C6D9971F** — Forecast base, growth, launch, seasonal, failure, backup, and recovery cost.
- [ ] **USEQ-AB343848** — Set budgets and anomaly alerts without allowing cost controls to terminate critical service silently.
- [ ] **USEQ-8517773D** — Detect idle, orphaned, overprovisioned, underutilized, duplicate, and forgotten resources.
- [ ] **USEQ-5887E8D3** — Review data transfer, storage tiers, retention, logging, and high-cardinality telemetry cost.
- [ ] **USEQ-72822800** — Optimize architecture only when total lifecycle cost and quality improve.
- [ ] **USEQ-0B270727** — Account for engineering labor, complexity, support, migration, and lock-in—not provider invoice alone.
- [ ] **USEQ-3C56970B** — Use commitment and reservation mechanisms only with credible demand and exit analysis.
- [ ] **USEQ-5464C515** — Protect security, privacy, accessibility, reliability, and recovery requirements from harmful cost cutting.
- [ ] **USEQ-0A8057B2** — Make expensive user actions, queries, abuse, and retry amplification visible and bounded.
- [ ] **USEQ-D7F14772** — Review supplier pricing changes and quota interactions.
- [ ] **USEQ-5F69A20F** — Measure cost per successful outcome, active user, transaction, tenant, or workload where useful.
- [ ] **USEQ-0EC26FF1** — Include environmental resource efficiency where data is available.
- [ ] **USEQ-4CA042C8** — Verify savings after implementation and watch for shifted cost elsewhere.
- [ ] **USEQ-1E9FC44C** — Remove temporary launch and migration capacity after exit criteria are met.
- [ ] **USEQ-26C038B6** — Maintain cost evidence for major architecture and procurement decisions.
- [ ] **USEQ-8E6BBEEE** — Treat unexpected cost growth as an operational incident when material.

## Asset, Patch, and Technology Lifecycle Management

_Consolidated from `quality standards/13-operations-sre/10-asset-patch-and-technology-lifecycle-management.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-64F3DC86** — Maintain an authoritative inventory of software, services, infrastructure, devices, domains, certificates, licenses, accounts, and dependencies.
- [ ] **USEQ-BD4C9D3B** — Record owner, purpose, environment, criticality, data, exposure, version, support status, and retirement date.
- [ ] **USEQ-DB080B59** — Discover unmanaged or shadow assets continuously.
- [ ] **USEQ-E6C9769D** — Track vendor support, end-of-life, certificate, license, key, and contract dates.
- [ ] **USEQ-D4A7A5BA** — Classify patches by exploit activity, exposure, reachability, impact, and operational risk.
- [ ] **USEQ-9099195E** — Test patches in representative environments and define emergency paths.
- [ ] **USEQ-E6FF7190** — Apply patches within risk-based targets and verify deployed versions.
- [ ] **USEQ-542BC36E** — Use compensating controls only with owner, monitoring, expiry, and replacement plan.
- [ ] **USEQ-10EB8EC0** — Prevent unsupported assets from silently remaining production dependencies.
- [ ] **USEQ-B10C60D7** — Coordinate application, runtime, operating system, database, firmware, network, and supplier updates.
- [ ] **USEQ-95DCE399** — Test mixed versions, rollback, recovery, and data compatibility.
- [ ] **USEQ-EBD5EE42** — Maintain clean trusted installation media, artifacts, and configuration.
- [ ] **USEQ-6A7A184A** — Monitor failed, partial, or reverted patch deployments.
- [ ] **USEQ-E4979EC7** — Remove obsolete packages, services, accounts, ports, and capabilities.
- [ ] **USEQ-677D277F** — Review license use and renewal before it can interrupt service.
- [ ] **USEQ-25987E4B** — Protect management interfaces and inventory data.
- [ ] **USEQ-BDBE6B90** — Reassess capacity and performance after major updates.
- [ ] **USEQ-3021D362** — Include patch and asset state in incident investigation.
- [ ] **USEQ-82716355** — Retire duplicate and orphaned assets.
- [ ] **USEQ-A4575CA1** — Report lifecycle risk to decision-makers before it becomes an emergency.

## Decommissioning and Secure Retirement

_Consolidated from `quality standards/13-operations-sre/11-decommissioning-and-secure-retirement.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-04858636** — Define owner, scope, dependencies, customers, users, data, contracts, and completion criteria.
- [ ] **USEQ-A7C0D18F** — Confirm replacement or continuity plans before disabling essential capability.
- [ ] **USEQ-A8406215** — Inventory inbound and outbound interfaces, jobs, queues, domains, certificates, secrets, accounts, suppliers, and data copies.
- [ ] **USEQ-AA9D494C** — Notify affected users, customers, operators, support teams, and partners with sufficient notice.
- [ ] **USEQ-0E6C800B** — Provide migration, export, rollback, and support appropriate to impact.
- [ ] **USEQ-5C76742D** — Stop new use and dependency creation before final shutdown.
- [ ] **USEQ-949BE861** — Drain or reconcile in-flight transactions, messages, files, and background work.
- [ ] **USEQ-1276B959** — Export, archive, retain, delete, or legally hold data according to approved rules.
- [ ] **USEQ-87571159** — Verify deletion across primary stores, caches, indexes, replicas, backups, analytics, suppliers, and test copies as applicable.
- [ ] **USEQ-0CA1D171** — Revoke credentials, keys, tokens, certificates, service identities, network paths, and administrative access.
- [ ] **USEQ-C0FC8ACE** — Remove DNS, routing, callbacks, webhooks, firewall rules, monitoring, alerts, and status entries.
- [ ] **USEQ-58AA8637** — Terminate contracts, subscriptions, licenses, and provider resources only after data and continuity obligations are satisfied.
- [ ] **USEQ-FFBF3B61** — Preserve required source, artifacts, provenance, documentation, audit, incident, and decision records.
- [ ] **USEQ-1873E848** — Verify no traffic, processing, billing, data collection, or externally reachable surface remains.
- [ ] **USEQ-0173A672** — Monitor for unexpected calls and dependencies after shutdown.
- [ ] **USEQ-35351CBA** — Destroy media and hardware according to sensitivity.
- [ ] **USEQ-90749519** — Update asset, architecture, data-flow, service, and ownership inventories.
- [ ] **USEQ-D1519935** — Conduct a final security and privacy review.
- [ ] **USEQ-6E6CB31E** — Record residual obligations and long-term record owners.
- [ ] **USEQ-ADC1A37E** — Review lessons and update lifecycle planning.

## Operational Access and Production Change

_Consolidated from `quality standards/13-operations-sre/12-operational-access-and-production-change.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-BDDCDDAA** — Use individual identities and strong authentication for production access.
- [ ] **USEQ-A9A76931** — Grant the least privilege for the shortest practical period.
- [ ] **USEQ-0F159E5F** — Separate ordinary development access from privileged production roles.
- [ ] **USEQ-7812911D** — Require purpose, ticket or incident, approval, scope, and expiry for elevated access.
- [ ] **USEQ-731B3211** — Use controlled administrative entry points and hardened workstations where risk warrants.
- [ ] **USEQ-4C2D10C6** — Record commands, queries, changes, targets, time, actor, reason, and outcome proportionately.
- [ ] **USEQ-5259F535** — Protect recordings and audit data from tampering and excessive sensitive-content capture.
- [ ] **USEQ-D8F85AD0** — Avoid direct manual changes when reviewed reproducible automation can perform the work.
- [ ] **USEQ-13511541** — Use dry run, preview, row limits, transaction, backup, and peer verification for data or destructive operations.
- [ ] **USEQ-FF5D690D** — Prevent broad queries and exports without explicit authorization.
- [ ] **USEQ-C749F02A** — Make emergency access available during identity or control-plane failure and alert on use.
- [ ] **USEQ-E22FBA9C** — Review break-glass activity promptly.
- [ ] **USEQ-EC73D6ED** — Revoke access automatically on expiry, role change, departure, or incident containment.
- [ ] **USEQ-0A33DF21** — Prevent support impersonation from hiding the true operator.
- [ ] **USEQ-5C0FF165** — Verify user and tenant context before acting on customer data.
- [ ] **USEQ-BAD2E536** — Use two-person control for irreversible high-impact operations where appropriate.
- [ ] **USEQ-EB86FC6E** — Test access paths and required tools before an incident.
- [ ] **USEQ-0126D4F9** — Review dormant privileges and actual usage.
- [ ] **USEQ-D8EB3355** — Measure manual change, failed intervention, access exception, and resulting incident trends.
- [ ] **USEQ-A5591561** — Convert repeated manual operations into safe automation and runbooks.

## Standards and source references

- [ISO/IEC 27001:2022 — Information security management systems](https://www.iso.org/standard/27001)
- [NIST Cybersecurity Framework 2.0](https://www.nist.gov/cyberframework)
- [ISO 22301:2019 — Business continuity management systems](https://www.iso.org/standard/75106.html)
- [OpenTelemetry Specifications](https://opentelemetry.io/docs/specs/)
- [Google SRE Workbook — Implementing SLOs](https://sre.google/workbook/implementing-slos/)
- [ISO/IEC 20000-1:2018 — Service management system requirements](https://www.iso.org/standard/70636.html)
- [NIST SP 800-61 Rev. 3 — Incident Response](https://csrc.nist.gov/pubs/sp/800/61/r3/final)
- [ISO/IEC 25010:2023 — Product quality model](https://www.iso.org/standard/78176.html)
- [ISO/IEC/IEEE 15939:2017 — Measurement process](https://www.iso.org/standard/71197.html)
- [ISO/IEC/IEEE 12207:2026 — Software life cycle processes](https://www.iso.org/standard/90219.html)
- [ISO 9001:2015 — Quality management systems](https://www.iso.org/standard/62085.html)
- [ISO/IEC 38500:2024 — Governance of IT](https://www.iso.org/standard/81684.html)
- [ISO/IEC 21031:2024 — Software Carbon Intensity](https://www.iso.org/standard/86612.html)
- [ISO/IEC 27701:2025 — Privacy information management systems](https://www.iso.org/standard/85819.html)
- [NIST SP 800-63-4 — Digital Identity Guidelines](https://pages.nist.gov/800-63-4/)

---

[Previous phase](11-developer-experience-platform-and-delivery.md) · [Next: Phase 13: Documentation and knowledge](13-documentation-and-knowledge.md)
