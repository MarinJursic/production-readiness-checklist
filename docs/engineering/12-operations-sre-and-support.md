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

## Digital Sustainability, Environmental Impact, and Responsible Computing

_Consolidated from `final consolidated corpus/08-operations-sre-observability-continuity-cost-sustainability.md#Digital Sustainability, Environmental Impact, and Responsible Computing`; 237 non-duplicative controls._

### Scope, governance, accountability, and ethical boundaries

- [ ] **USEQ-5CF3177D** — Define which products, services, organizations, suppliers, users, regions, infrastructure, lifecycle stages, and environmental impact categories are inside the sustainability assessment boundary.
- [ ] **USEQ-13992319** — Assign an accountable sustainability owner with authority to influence product, architecture, procurement, operations, and retirement decisions.
- [ ] **USEQ-58F7D983** — Define decision rights for resolving conflicts among sustainability, accessibility, security, privacy, safety, reliability, performance, affordability, and business requirements.
- [ ] **USEQ-88774BC5** — Include affected users, workers, communities, suppliers, and environmental stakeholders in material sustainability decisions.
- [ ] **USEQ-888DCA36** — Treat environmental, social, and economic sustainability as connected concerns rather than optimizing carbon in isolation.
- [ ] **USEQ-21CC1B21** — Document the organization’s theory of change, material impacts, dependencies, assumptions, and limits of influence.
- [ ] **USEQ-D3F0C27B** — Identify relevant environmental laws, reporting rules, customer commitments, procurement requirements, and sector-specific obligations.
- [ ] **USEQ-6C5D6E42** — Prohibit sustainability claims that cannot be traced to a defined boundary, method, evidence set, time period, and responsible approver.
- [ ] **USEQ-7B0F07E8** — Distinguish mandatory requirements, voluntary goals, experiments, estimates, and aspirational targets in all plans and communications.
- [ ] **USEQ-7A460804** — Include sustainability requirements in product requirements, architecture decisions, supplier contracts, Definition of Done, release gates, and retirement plans.
- [ ] **USEQ-A060C72B** — Use independent review for material public claims, regulated disclosures, major lifecycle assessments, and high-impact trade-offs.
- [ ] **USEQ-1885C630** — Keep sustainability governance active after launch through operational ownership, reassessment triggers, and funded improvement work.

### Necessity, demand shaping, product strategy, and rebound effects

- [ ] **USEQ-CBFAABA8** — Assess whether the product, feature, automation, data collection, model, or recurring process is necessary before optimizing its implementation.
- [ ] **USEQ-22D56445** — Define the user or societal outcome that justifies the resource demand and identify lower-impact non-digital or simpler digital alternatives.
- [ ] **USEQ-97517034** — Reject features whose expected harm, waste, manipulation, exclusion, or maintenance burden outweighs evidenced benefit.
- [ ] **USEQ-1D0695BB** — Reduce unnecessary engagement loops, autoplay, infinite consumption, dark patterns, spam, and attention extraction that create avoidable resource use.
- [ ] **USEQ-5DC602F1** — Design pricing, quotas, defaults, incentives, and success metrics so they do not reward wasteful usage or artificial activity.
- [ ] **USEQ-4705259A** — Assess whether efficiency improvements could increase total demand and offset or reverse expected reductions.
- [ ] **USEQ-4FBD8481** — Measure absolute environmental impact as well as normalized efficiency so growth is not hidden by better per-unit metrics.
- [ ] **USEQ-F1113914** — Establish environmental and human budgets for features, campaigns, experiments, models, data products, and infrastructure.
- [ ] **USEQ-99D85387** — Prioritize reuse, repair, extension, consolidation, and retirement before creating additional systems or duplicate capabilities.
- [ ] **USEQ-6701A91D** — Define explicit sunset criteria for low-value features, experiments, data sets, models, environments, domains, and services.
- [ ] **USEQ-0EF58E65** — Include long-term maintenance, migration, support, and end-of-life costs in business cases.
- [ ] **USEQ-75B2BF2E** — Avoid shifting environmental cost to users through excessive device requirements, bandwidth, battery consumption, printing, travel, or repeated manual effort.

### Measurement model, boundaries, baselines, and uncertainty

- [ ] **USEQ-F4CC9192** — Define a reproducible measurement method before using a sustainability metric as a target or claim.
- [ ] **USEQ-2E6577FB** — Specify the functional unit used to normalize impact and prove that it represents meaningful delivered value rather than convenient activity.
- [ ] **USEQ-5F686CE1** — Define organizational, operational, temporal, geographic, lifecycle, and supplier boundaries for every measurement.
- [ ] **USEQ-3FFF3A34** — Separate operational energy and emissions from embodied impacts associated with hardware, facilities, networks, and manufacturing.
- [ ] **USEQ-1D668E78** — Measure or estimate greenhouse-gas emissions, energy, water, materials, electronic waste, land or biodiversity effects, and relevant pollution where material.
- [ ] **USEQ-9E6E23B4** — Distinguish direct measurement, provider data, modeled estimates, proxies, and assumptions.
- [ ] **USEQ-1C242EF6** — Record data sources, sampling, allocation rules, emission factors, grid factors, conversion methods, and versions.
- [ ] **USEQ-18A9F03F** — Report material uncertainty, data gaps, sensitivity, and confidence rather than presenting estimates as exact facts.
- [ ] **USEQ-8E40F7EA** — Use consistent methods across baselines and comparison periods or disclose and quantify method changes.
- [ ] **USEQ-D9D0B323** — Separate absolute totals, intensity metrics, avoided impact, enabled impact, and offsets; do not substitute one for another.
- [ ] **USEQ-605351E9** — Do not count hypothetical avoided emissions as equivalent to measured reductions in the product’s own footprint.
- [ ] **USEQ-9B9B9BBD** — Do not use offsets to avoid feasible direct reduction work, and disclose offset quality and boundaries when offsets are reported.
- [ ] **USEQ-D9A7B76F** — Protect measurement telemetry against tampering, accidental double counting, inappropriate aggregation, and privacy leakage.
- [ ] **USEQ-7F6B32D9** — Validate measurement pipelines and reconcile them against billing, capacity, asset, and provider records where possible.

### Targets, budgets, prioritization, and decision records

- [ ] **USEQ-55AD6D41** — Establish time-bounded, measurable targets linked to a documented baseline and product outcome.
- [ ] **USEQ-41D1322A** — Decompose organizational goals into product, service, workload, supplier, asset, and team responsibilities without double counting.
- [ ] **USEQ-006458B6** — Use budgets for transfer, requests, compute, memory, storage, retention, model inference, training, build work, and observability where material.
- [ ] **USEQ-1CDC439C** — Define thresholds that trigger investigation, rollback, redesign, capacity review, or retirement.
- [ ] **USEQ-4BBD1617** — Prioritize work using expected impact, confidence, cost, reversibility, time to benefit, and effects on people and resilience.
- [ ] **USEQ-9CE76496** — Record significant sustainability trade-offs in architecture and product decision records.
- [ ] **USEQ-63805953** — Test whether an optimization simply moves impact between client and server, region and region, present and future, or organization and user.
- [ ] **USEQ-81CEFB3A** — Prefer interventions that improve multiple quality attributes without creating hidden harms.
- [ ] **USEQ-00554C4D** — Require a compensating plan when a justified feature exceeds an approved budget.
- [ ] **USEQ-04D91E1E** — Do not average away severe local harm through favorable global totals.
- [ ] **USEQ-B5D5DF49** — Review targets after material traffic growth, architecture change, supplier change, hardware refresh, model change, incident, or new measurement evidence.
- [ ] **USEQ-5AE35A89** — Track improvement actions to verified outcomes rather than closing them when implementation activity is complete.

### Inclusive, efficient, and durable user experience

- [ ] **USEQ-CB209167** — Design the shortest understandable path to the user’s legitimate goal and remove unnecessary steps, reloads, retries, and repeated data entry.
- [ ] **USEQ-C45FED74** — Make essential journeys usable on lower-powered devices, slower or intermittent networks, small screens, assistive technologies, and constrained data plans where the intended population requires it.
- [ ] **USEQ-4DFA935A** — Preserve accessibility, readability, safety, and user control when reducing assets or functionality.
- [ ] **USEQ-8C18F591** — Use progressive enhancement and graceful degradation so core outcomes do not depend unnecessarily on high-cost capabilities.
- [ ] **USEQ-630FB495** — Let users control autoplay, animation, background activity, notifications, downloads, synchronization, quality levels, and data-heavy behavior where practical.
- [ ] **USEQ-01CB4937** — Default nonessential high-impact behavior to off or to the lowest reasonable resource level.
- [ ] **USEQ-F78DBB93** — Communicate the resource, cost, privacy, or quality consequences of meaningful user choices without coercion or blame.
- [ ] **USEQ-B2FF50FE** — Prevent accidental repeated submissions, downloads, payments, renders, uploads, and notifications.
- [ ] **USEQ-FDABAD0D** — Use clear status, validation, recovery, and offline messaging to reduce wasted attempts and support demand.
- [ ] **USEQ-9040A903** — Avoid forced device upgrades, unnecessary browser restrictions, or short compatibility windows without documented necessity.
- [ ] **USEQ-63E6B914** — Measure task success, abandonment, retries, support contacts, and time spent as sustainability signals as well as UX signals.
- [ ] **USEQ-CBB44B7C** — Involve representative users early, including people affected by bandwidth, energy, device, language, disability, or economic constraints.

### Content, media, typography, documents, and communications

- [ ] **USEQ-8A741146** — Publish only content that has an identified audience, owner, purpose, review date, retention rule, and retirement path.
- [ ] **USEQ-FE511AC4** — Remove or archive stale, duplicate, contradictory, unowned, and low-value content.
- [ ] **USEQ-A08EC0B1** — Use concise, structured, plain-language content that helps users complete tasks without repeated searches or support contacts.
- [ ] **USEQ-2C3E104A** — Provide text or lower-bandwidth alternatives for media where they preserve the intended outcome.
- [ ] **USEQ-C291BB52** — Select image, audio, video, animation, font, and document formats appropriate to actual quality needs and supported clients.
- [ ] **USEQ-97644690** — Compress, resize, crop, transcode, and deliver media according to display context rather than shipping source-quality assets by default.
- [ ] **USEQ-A6FE8EE0** — Avoid autoplay, decorative high-resolution media, unnecessary background video, and repeated asset downloads.
- [ ] **USEQ-816E7255** — Use responsive media and explicit dimensions to avoid unnecessary transfer and relayout.
- [ ] **USEQ-53953C3F** — Subset and preload fonts only when evidence shows benefit; minimize font families, weights, variants, and blocking behavior.
- [ ] **USEQ-6B5533CE** — Design downloadable documents for reuse, accessibility, efficient file size, and long-term readability.
- [ ] **USEQ-DD89FE18** — Avoid requiring printing when an accessible digital workflow can achieve the outcome.
- [ ] **USEQ-E4ABE085** — Prevent marketing, analytics, and notification campaigns from creating disproportionate transfer, storage, attention, or support waste.
- [ ] **USEQ-A60E4DA2** — Expire obsolete feeds, previews, cached assets, derivatives, and campaign resources safely.

### Client, frontend, application, and code efficiency

- [ ] **USEQ-7477455D** — Set measurable client-side budgets for transferred bytes, requests, execution time, memory, energy-sensitive work, long tasks, and background activity.
- [ ] **USEQ-481F4341** — Remove unused, duplicate, obsolete, unreachable, and unnecessarily generated code and assets.
- [ ] **USEQ-99D96F89** — Defer, lazy-load, stream, or conditionally load nonessential resources without harming accessibility, discoverability, security, or correctness.
- [ ] **USEQ-E2330712** — Prefer native platform capabilities and simpler representations when they meet requirements with less code and dependency overhead.
- [ ] **USEQ-49967306** — Modularize high-cost functionality so users who do not need it do not pay its resource cost.
- [ ] **USEQ-E0C0004A** — Treat first-party and third-party code under the same performance, privacy, security, and sustainability budgets.
- [ ] **USEQ-588808D0** — Minimize polling, timers, observers, background synchronization, speculative work, and network chatter.
- [ ] **USEQ-92BC5536** — Pause or reduce work when the application is hidden, inactive, offline, on a constrained connection, or under user power-saving preferences where safe.
- [ ] **USEQ-D0E5ADD2** — Use caching carefully to avoid repeated transfer while preventing stale, private, or permanently orphaned data.
- [ ] **USEQ-0555DF66** — Choose rendering and hydration strategies based on measured end-to-end cost across server, network, and client rather than fashion.
- [ ] **USEQ-A084DB55** — Prevent memory leaks, runaway event handlers, duplicate subscriptions, and retained resources that waste energy over long sessions.
- [ ] **USEQ-7E8877C1** — Include client resource regressions in continuous integration and production monitoring.
- [ ] **USEQ-7110753C** — Keep semantic, accessible, indexable output while optimizing implementation size and execution.

### Backend, API, algorithm, database, and data efficiency

- [ ] **USEQ-B2C8404E** — Define resource budgets and cost models for high-volume endpoints, jobs, event flows, searches, reports, and data transformations.
- [ ] **USEQ-523BAFB1** — Choose algorithms and data structures based on realistic scale, resource use, correctness, and maintainability.
- [ ] **USEQ-FB833D84** — Eliminate repeated computation, unnecessary serialization, avoidable fan-out, excessive round trips, and duplicate storage.
- [ ] **USEQ-E186F16E** — Bound requests, pagination, query complexity, response size, retries, recursion, batch size, and execution duration.
- [ ] **USEQ-90DB858D** — Use batching, caching, materialization, precomputation, and asynchronous processing only when freshness, consistency, authorization, and failure behavior remain correct.
- [ ] **USEQ-8FC65F49** — Optimize database access using measured plans, indexes, cardinality, locality, and workload behavior.
- [ ] **USEQ-B1F3AB1F** — Prevent unbounded scans, N-plus-one patterns, unused columns, excessive joins, and repeatedly derived values where they are material.
- [ ] **USEQ-8F0A9252** — Define data freshness and refresh frequencies based on user need rather than maximum technical frequency.
- [ ] **USEQ-10CDDC3D** — Use event-driven or scheduled work rather than continuous polling when it is reliably more efficient.
- [ ] **USEQ-21F0724F** — Deduplicate idempotent work and prevent retry storms, hot loops, poison messages, and unbounded queue growth.
- [ ] **USEQ-951A3A18** — Control observability volume, cardinality, sampling, and retention without losing required diagnostic, security, or audit evidence.
- [ ] **USEQ-4C7049BA** — Continuously identify and retire low-value jobs, indexes, replicas, caches, exports, reports, and derived data.
- [ ] **USEQ-365B130E** — Test resource behavior with representative volume, skew, cache state, concurrency, failure, and recovery.

### Network, transfer, caching, and distribution

- [ ] **USEQ-B2C4739E** — Minimize round trips and transfer volume while preserving correctness, freshness, privacy, security, and accessibility.
- [ ] **USEQ-24A45B5D** — Use compression suited to content type, client capability, latency, and compute trade-offs.
- [ ] **USEQ-5FD8938A** — Use durable caching directives and content-addressed assets where safe.
- [ ] **USEQ-DDD0B8A8** — Avoid invalidation strategies that repeatedly force full downloads when smaller or versioned updates suffice.
- [ ] **USEQ-8F05EC03** — Use edge delivery, content distribution, peer distribution, or regional replication only when measured total impact and resilience justify them.
- [ ] **USEQ-2914A322** — Keep data close to users and processing where it reduces transfer without violating residency, privacy, or resilience needs.
- [ ] **USEQ-81C3B6CD** — Reduce chatty protocols and choose efficient representations for high-volume paths.
- [ ] **USEQ-AB9AE435** — Support resumable and partial transfer for large resources where interrupted transfer is plausible.
- [ ] **USEQ-D517D2B1** — Prevent crawlers, bots, scrapers, abusive clients, and accidental loops from causing unbounded traffic.
- [ ] **USEQ-793660E3** — Use rate limits, quotas, backpressure, and request coalescing to control waste and protect shared resources.
- [ ] **USEQ-2205156F** — Monitor cache hit ratios, transfer per successful outcome, retransmission, aborted work, and regional paths.
- [ ] **USEQ-707C1649** — Include network equipment, intermediary services, and provider egress in impact boundaries where material.

### Hosting, infrastructure, capacity, and operational efficiency

- [ ] **USEQ-50405BB4** — Select hosting and infrastructure using verifiable energy, emissions, water, resource, resilience, data-location, and supplier-governance evidence.
- [ ] **USEQ-5DAA5FA9** — Right-size capacity using measured demand, required failure headroom, startup characteristics, and service objectives.
- [ ] **USEQ-0C287229** — Eliminate idle, abandoned, duplicate, oversized, and permanently overprovisioned resources.
- [ ] **USEQ-40D1DD20** — Use autoscaling, scheduling, consolidation, suspension, and shutdown controls that preserve availability and recovery requirements.
- [ ] **USEQ-6D96904B** — Schedule flexible workloads for lower-impact times or regions when privacy, latency, reliability, and contractual constraints permit.
- [ ] **USEQ-DF29852A** — Distinguish grid average, market-based claims, marginal effects, and actual renewable matching in infrastructure decisions.
- [ ] **USEQ-8AB2AF58** — Include backup power, cooling, storage, network, facility overhead, and embodied infrastructure in assessments where material.
- [ ] **USEQ-2738410D** — Monitor energy, utilization, carbon intensity, water, waste, capacity, and service outcomes at a granularity that supports action.
- [ ] **USEQ-F8EB7249** — Avoid maintaining unused environments, replicas, clusters, virtual machines, containers, snapshots, volumes, addresses, and reserved capacity.
- [ ] **USEQ-47DBC483** — Test shutdown, scale-to-zero, wake-up, failover, and recovery behavior before relying on them for reduction.
- [ ] **USEQ-6DF5A476** — Ensure resilience architecture survives required failures without keeping unjustified permanent duplication.
- [ ] **USEQ-C1FC62BA** — Include control planes, observability, build systems, development environments, and disaster-recovery infrastructure in the footprint.
- [ ] **USEQ-EE0EEBE1** — Negotiate access to provider sustainability data and document limitations when data is unavailable.

### Storage, retention, backups, archives, and data lifecycle

- [ ] **USEQ-CE56DAD4** — Collect and retain only data justified by a current user, legal, safety, security, operational, or research purpose.
- [ ] **USEQ-B2FD1EBB** — Classify hot, warm, cold, archival, backup, and disposable data and use storage suited to access and recovery needs.
- [ ] **USEQ-5AA75D77** — Define deletion, compaction, deduplication, compression, tiering, and archive policies.
- [ ] **USEQ-8EF05145** — Remove orphaned blobs, snapshots, logs, indexes, replicas, derivatives, test data, and abandoned exports.
- [ ] **USEQ-6CC15740** — Use retention periods based on evidenced need rather than indefinite default retention.
- [ ] **USEQ-7BC8164C** — Distinguish backups from archives and avoid using repeated full backups as an unmanaged preservation strategy.
- [ ] **USEQ-691E5191** — Optimize backup frequency, incrementality, retention, geographic copies, and testing against recovery objectives and impact.
- [ ] **USEQ-72E8E360** — Preserve legal holds and evidentiary requirements without retaining unrelated data.
- [ ] **USEQ-F8EFD0C1** — Design deletion to propagate across caches, indexes, replicas, analytical systems, derived assets, and suppliers.
- [ ] **USEQ-D6515E61** — Monitor storage growth by owner, purpose, class, age, duplication, access frequency, and deletion eligibility.
- [ ] **USEQ-BC43D21A** — Prevent compression, deduplication, or tiering from weakening encryption, access control, recovery, or integrity.
- [ ] **USEQ-5863DC60** — Include the energy and material cost of long-term preservation, migration, and media refresh in lifecycle plans.

### AI, machine learning, automation, and high-compute workloads

- [ ] **USEQ-267D27DD** — Justify the use of an AI or high-compute method against simpler alternatives that can achieve the required outcome.
- [ ] **USEQ-FD5481AF** — Define lifecycle boundaries for data preparation, training, fine-tuning, evaluation, inference, retrieval, monitoring, experimentation, and hardware.
- [ ] **USEQ-329FEE22** — Measure energy, emissions, water, hardware use, utilization, and functional outcomes for material AI workloads.
- [ ] **USEQ-4F4E1909** — Report model size and compute only alongside task quality, usage, lifetime, hardware, region, and uncertainty needed to interpret impact.
- [ ] **USEQ-5710DB96** — Use the smallest model, precision, context, retrieval scope, generation length, and frequency that reliably meet requirements.
- [ ] **USEQ-2809A00E** — Reuse appropriate existing models and artifacts when lifecycle impact, security, privacy, licensing, and quality are favorable.
- [ ] **USEQ-FB9693DF** — Control repeated experiments, hyperparameter searches, evaluation duplication, speculative agents, and unattended loops.
- [ ] **USEQ-BA68CA4F** — Cache or reuse safe deterministic results where it reduces work without creating stale, private, or incorrect outcomes.
- [ ] **USEQ-AE78BFC0** — Route requests by complexity only when routing quality, fairness, privacy, and failure behavior are validated.
- [ ] **USEQ-23A25A60** — Schedule flexible training and batch inference according to capacity and environmental signals where suitable.
- [ ] **USEQ-6169B64D** — Monitor resource and environmental regressions across model, prompt, data, provider, and traffic changes.
- [ ] **USEQ-A024A74B** — Include embodied accelerator impact and hardware replacement cycles in major procurement and architecture decisions.
- [ ] **USEQ-12E1E291** — Do not claim an AI system is sustainable from operational efficiency alone when training, hardware, rebound, or enabled effects are material.
- [ ] **USEQ-0B2665A3** — Provide kill switches and spending or resource limits for autonomous and agentic workloads.

### Hardware lifecycle, assets, circularity, and electronic waste

- [ ] **USEQ-62B94A8E** — Maintain an inventory of material hardware and virtualized assets with owner, age, location, utilization, support state, and retirement path.
- [ ] **USEQ-348B40AE** — Extend useful hardware life when security, reliability, compatibility, repairability, and efficiency remain acceptable.
- [ ] **USEQ-6C407ACB** — Avoid software requirements and update policies that force unjustified hardware replacement.
- [ ] **USEQ-3FE758C5** — Include embodied impact, repairability, modularity, spare parts, support duration, and take-back options in procurement.
- [ ] **USEQ-FE206CBE** — Prefer repair, redeployment, refurbishment, resale, donation, parts recovery, and certified recycling over disposal where safe and lawful.
- [ ] **USEQ-D6FAF009** — Erase data securely and verify custody during reuse, return, recycling, and disposal.
- [ ] **USEQ-04DF1172** — Track electronic waste by type, destination, treatment, recovery, and evidence of responsible handling.
- [ ] **USEQ-8A5BFC8A** — Prevent asset hoarding, forgotten stock, premature replacement, and perpetual retention of unsupported equipment.
- [ ] **USEQ-F19BB58C** — Coordinate capacity planning and hardware refresh so efficiency gains are not erased by excess inventory or reduced lifetime.
- [ ] **USEQ-CBABBF8C** — Require suppliers to disclose relevant lifecycle, materials, labor, repair, packaging, transport, and end-of-life information.
- [ ] **USEQ-016BEEB7** — Include batteries, displays, peripherals, networking equipment, storage media, and user devices where the product materially drives their use.
- [ ] **USEQ-217653AB** — Use asset identification and utilization records to support sustainability decisions and audits.

### Water, materials, waste, pollution, ecosystems, and location risk

- [ ] **USEQ-8CDDB1C3** — Assess water use and water-stress context for material computing, cooling, energy, and manufacturing dependencies.
- [ ] **USEQ-10A887ED** — Avoid optimizing energy or emissions in ways that create disproportionate water, material, pollution, biodiversity, or community harm.
- [ ] **USEQ-AFA6FD10** — Identify hazardous materials, chemicals, refrigerants, batteries, packaging, and waste streams associated with material infrastructure.
- [ ] **USEQ-51A9A772** — Assess location-specific climate, grid, water, ecosystem, disaster, and community risks rather than relying only on global averages.
- [ ] **USEQ-89709CE3** — Include upstream mining, manufacturing, transport, construction, and downstream disposal when they materially affect decisions.
- [ ] **USEQ-7808AE7C** — Use life-cycle assessment principles for major comparisons and document exclusions and data quality.
- [ ] **USEQ-0D63B3D7** — Prefer suppliers with transparent environmental management, credible targets, responsible sourcing, and independently supportable data.
- [ ] **USEQ-7872C26C** — Detect and prevent illegal, unsafe, or opaque waste export and disposal pathways.
- [ ] **USEQ-6DA0EC95** — Consider biodiversity and land-use effects in large facilities, energy sourcing, and hardware supply chains where material.
- [ ] **USEQ-F7780EDB** — Assess adaptation needs for heat, drought, flood, wildfire, grid instability, and other climate-related conditions.
- [ ] **USEQ-2B1C8F57** — Include environmental justice and unequal community impacts in site and supplier decisions.
- [ ] **USEQ-DBC553F9** — Document when reliable data is unavailable and avoid making precise claims from weak proxies.

### Suppliers, procurement, contracts, and shared responsibility

- [ ] **USEQ-5804B79B** — Include measurable sustainability requirements in supplier selection and renewal for material products and services.
- [ ] **USEQ-106A573E** — Require suppliers to define measurement boundaries, methods, data quality, assurance, and update cadence for reported claims.
- [ ] **USEQ-92042528** — Distinguish supplier-wide averages from the impact attributable to the purchased service.
- [ ] **USEQ-07CD6281** — Obtain data on energy, emissions, water, waste, hardware lifecycle, regions, and subcontractors where material.
- [ ] **USEQ-BEE4F085** — Require disclosure of significant method changes and restatements.
- [ ] **USEQ-CD3E24CA** — Assess whether contractual renewable-energy claims correspond to the locations and time periods of actual consumption.
- [ ] **USEQ-1129D00B** — Include rights to audit, request evidence, remediate, terminate, export data, and migrate when material commitments are not met.
- [ ] **USEQ-A80ADE7A** — Avoid procurement criteria that privilege polished reporting over actual verified reduction and responsible practice.
- [ ] **USEQ-A984E0E3** — Assess concentration, lock-in, data portability, and replacement impact before outsourcing critical capabilities.
- [ ] **USEQ-6C3BAA57** — Flow applicable requirements to subcontractors and supply-chain participants.
- [ ] **USEQ-C0E91E6F** — Use supplier scorecards carefully and prevent unverified composite scores from hiding severe weaknesses.
- [ ] **USEQ-1B6F08D8** — Revoke access, delete data, recover assets, and close environmental obligations when suppliers are retired.

### Longevity, compatibility, resilience, and low-resource operation

- [ ] **USEQ-FD386600** — Design stable interfaces, data formats, identifiers, and migration paths that reduce premature replacement and repeated redevelopment.
- [ ] **USEQ-E5E831CE** — Support reasonable client, operating-system, protocol, and assistive-technology lifetimes based on user context and security risk.
- [ ] **USEQ-788DE467** — Document compatibility windows and deprecation decisions with user impact, migration cost, and environmental consequences.
- [ ] **USEQ-A676EC10** — Use graceful degradation, offline capability, local recovery, and resumable work where they improve resilience and reduce waste.
- [ ] **USEQ-D82A953C** — Avoid brittle dependencies on one provider, region, proprietary format, or short-lived service without an exit path.
- [ ] **USEQ-DD4E26D9** — Preserve data portability and interoperable export so users can move without unnecessary recreation or loss.
- [ ] **USEQ-BD44967E** — Maintain disaster recovery and business continuity without uncontrolled permanent resource duplication.
- [ ] **USEQ-E6F6B0C4** — Design for climate-related infrastructure disruption, degraded connectivity, supply shortages, and power constraints where relevant.
- [ ] **USEQ-2D6A7EF2** — Keep recovery procedures, artifacts, keys, documentation, and skills viable for the intended service lifetime.
- [ ] **USEQ-18554C68** — Retire unsupported, insecure, or inefficient technology through planned migration rather than emergency replacement.
- [ ] **USEQ-61E768BD** — Measure product lifetime, upgrade frequency, migration work, and redevelopment as sustainability indicators.
- [ ] **USEQ-CE286A5D** — Ensure optimizations do not reduce diagnosability or make failures more expensive to detect and repair.

### Privacy, security, accessibility, and safety safeguards

- [ ] **USEQ-2C5DF06D** — Do not weaken accessibility, privacy, security, safety, reliability, or user rights to meet an environmental target.
- [ ] **USEQ-FE5C9BBC** — Collect sustainability telemetry at the least identifying granularity needed for the decision.
- [ ] **USEQ-B32833BD** — Aggregate, minimize, protect, retain, and delete energy or behavior telemetry according to documented purpose.
- [ ] **USEQ-C11D4863** — Prevent sustainability scores, device characteristics, resource patterns, or power signals from becoming tracking or fingerprinting mechanisms.
- [ ] **USEQ-13DA6246** — Do not use resource efficiency as justification for omitting necessary authentication, encryption, backups, audit records, or abuse controls.
- [ ] **USEQ-9CBBD1BC** — Optimize security controls through architecture and proportional configuration rather than disabling them.
- [ ] **USEQ-ABB8A79B** — Include malicious traffic, fraud, scraping, spam, cryptomining, denial of service, and compromised automation in waste-reduction controls.
- [ ] **USEQ-9B66981F** — Ensure low-resource or offline modes preserve authorization, integrity, privacy, and safe recovery.
- [ ] **USEQ-19B4881C** — Provide accessible alternatives when visual, audio, animation, or interaction reductions affect users differently.
- [ ] **USEQ-F55E8913** — Review unequal effects on users with older devices, disabilities, limited connectivity, low income, or dependence on shared devices.
- [ ] **USEQ-0C6BFCF8** — Document and approve unavoidable trade-offs with affected stakeholders and qualified reviewers.
- [ ] **USEQ-3A2A4DEE** — Monitor whether environmental controls create new exclusion, surveillance, manipulation, or safety risks.

### Reporting, claims, transparency, assurance, and anti-greenwashing

- [ ] **USEQ-346178E8** — Define every public or contractual sustainability claim in testable terms.
- [ ] **USEQ-17FA871A** — Identify the assessed product version, geography, time period, functional unit, lifecycle boundary, method, data sources, and exclusions.
- [ ] **USEQ-CB26CF79** — Disclose uncertainty, estimates, supplier dependence, material limitations, and known negative impacts.
- [ ] **USEQ-FEF4E43F** — Do not describe partial conformance, a single efficiency improvement, renewable-energy procurement, or offsets as total sustainability.
- [ ] **USEQ-808F4B6E** — Do not compare products unless the boundary, functional unit, quality level, lifetime, and methods are sufficiently aligned.
- [ ] **USEQ-7C8C50E7** — Keep evidence available for the life of the claim and any required retention period.
- [ ] **USEQ-3CB866B2** — Have qualified independent reviewers verify high-impact claims and regulated disclosures.
- [ ] **USEQ-34536507** — Correct or withdraw claims promptly when data, methods, products, suppliers, or assumptions change materially.
- [ ] **USEQ-4B9F5962** — Publish progress and setbacks using consistent historical baselines or clearly restate prior values.
- [ ] **USEQ-278F1D86** — Separate measured reductions from projections, commitments, avoided emissions, enabled impacts, and compensatory instruments.
- [ ] **USEQ-9771DA8A** — Prevent dashboards from using incomplete precision, favorable denominators, omitted growth, or selective time windows.
- [ ] **USEQ-68334A49** — Provide a contact and correction process for sustainability concerns and data errors.

### Continuous improvement, operations, incidents, and release gates

- [ ] **USEQ-1B038E56** — Monitor resource and environmental indicators alongside user outcomes, service objectives, cost, and quality.
- [ ] **USEQ-E58D1198** — Annotate deployments, model changes, campaigns, supplier changes, and traffic shifts on sustainability trends.
- [ ] **USEQ-D899C2D4** — Investigate unexpected increases in compute, storage, transfer, energy, water, waste, or hardware demand.
- [ ] **USEQ-8642C1BC** — Include sustainability in incident and problem reviews when waste, runaway activity, capacity, or supplier behavior contributed.
- [ ] **USEQ-374A20F0** — Create alerts or automated limits for uncontrolled jobs, agent loops, queue growth, data growth, traffic, builds, experiments, and scaling.
- [ ] **USEQ-D2B9385B** — Test that shutdown, throttling, degradation, and kill switches work without causing unsafe or irreversible outcomes.
- [ ] **USEQ-3ADC9F49** — Convert recurring waste into architectural, product, operational, or supplier corrective action rather than one-time cleanup.
- [ ] **USEQ-1C89DDCE** — Reassess after major changes to usage, regions, algorithms, models, devices, suppliers, regulations, or measurement methods.
- [ ] **USEQ-7C1E20A5** — Retain release evidence showing that applicable environmental budgets and no-go gates passed.
- [ ] **USEQ-624C88CF** — Block release when a material impact is unknown because required measurement or ownership is absent.
- [ ] **USEQ-5139B668** — Block release when the product exceeds an approved hard budget without an authorized, monitored, expiring exception.
- [ ] **USEQ-3820272D** — Block release when a public environmental claim is unsupported, misleading, or materially inconsistent with actual operation.
- [ ] **USEQ-BDA8D2C8** — Block release when optimization creates an unresolved accessibility, privacy, security, safety, integrity, or reliability failure.

## Final Gap Closure — Operational Acceptance, Service Transition, and Organizational Adoption

_Consolidated from `final consolidated corpus/08-operations-sre-observability-continuity-cost-sustainability.md#Final Gap Closure — Operational Acceptance, Service Transition, and Organizational Adoption`; 118 non-duplicative controls._

### Transition governance and accountability

- [ ] **USEQ-8E27F1AD** — Define the service, capability, users, organizations, suppliers, locations, data, and environments included in the transition.
- [ ] **USEQ-35A71BDE** — Assign an accountable transition owner with authority across product, delivery, operations, support, security, privacy, data, and suppliers.
- [ ] **USEQ-FDCEFCB5** — Identify the people authorized to accept operational responsibility and residual risk.
- [ ] **USEQ-FB298C6E** — Establish transition objectives, success measures, entry criteria, exit criteria, stop conditions, rollback criteria, and observation period.
- [ ] **USEQ-6FA91009** — Maintain an integrated transition plan covering technology, data, people, process, contracts, communications, training, support, and retirement.
- [ ] **USEQ-0697E472** — Identify dependencies, sequencing constraints, freeze periods, peak events, and competing organizational changes.
- [ ] **USEQ-0904B9FE** — Assess user, customer, worker, accessibility, safety, security, privacy, financial, and continuity impacts.
- [ ] **USEQ-6159AEF8** — Define governance for decisions during delivery–operations overlap.
- [ ] **USEQ-A4F57314** — Prevent schedule pressure from transferring an unaccepted or unsupported service into operation.
- [ ] **USEQ-1243DFDC** — Keep transition risks, assumptions, actions, owners, evidence, and decisions current.
- [ ] **USEQ-1CBA7858** — Replan when scope, readiness, staffing, suppliers, data, or operating conditions change materially.

### Operational acceptance criteria

- [ ] **USEQ-5FFA89EF** — Define operational acceptance criteria while requirements and architecture are being established.
- [ ] **USEQ-58BC0BE5** — Include service objectives, capacity, resilience, recoverability, security, privacy, accessibility, supportability, maintainability, cost, and sustainability.
- [ ] **USEQ-A015BF2C** — Include observability, alerting, on-call, escalation, incident, problem, change, and continuity readiness.
- [ ] **USEQ-B5D1F965** — Include data integrity, migration, reconciliation, retention, deletion, backup, restore, and legal-hold readiness.
- [ ] **USEQ-146FDE47** — Include supplier support, contractual obligations, licenses, quotas, service limits, and exit provisions.
- [ ] **USEQ-31DE5766** — Include production access, privileged procedures, break-glass access, and separation of duties.
- [ ] **USEQ-F5DC0C23** — Include runbooks for routine, degraded, emergency, and recovery operation.
- [ ] **USEQ-EB933861** — Include documentation, training, competence, staffing, handover, and support-channel readiness.
- [ ] **USEQ-DEB21FB0** — Include known defects, exceptions, technical debt, workarounds, and end-of-support commitments.
- [ ] **USEQ-FC3E6CE7** — Make criteria measurable and linked to evidence.
- [ ] **USEQ-DA324CE6** — Define which unmet criteria block acceptance and which may receive time-bounded conditional acceptance.
- [ ] **USEQ-01C02CC6** — Require the receiving organization to review and agree to the criteria before final testing.

### Service design for operability and supportability

- [ ] **USEQ-CFA62933** — Design the service so ordinary operation does not depend on undocumented expert intervention.
- [ ] **USEQ-40F75816** — Provide clear ownership for every service, component, data set, integration, job, alert, runbook, and supplier.
- [ ] **USEQ-212321D8** — Expose health, version, configuration, dependencies, queueing, saturation, and user-impact indicators.
- [ ] **USEQ-9A08AF02** — Provide safe diagnostics without exposing secrets or personal data.
- [ ] **USEQ-6C68E857** — Make routine maintenance, patching, scaling, backup, restore, failover, and certificate rotation executable and testable.
- [ ] **USEQ-4AEF1160** — Provide bounded and auditable administrative actions.
- [ ] **USEQ-F2960F7A** — Prevent support and operations tools from bypassing authorization, tenant isolation, or data-governance controls.
- [ ] **USEQ-D06B97D0** — Provide safe maintenance and degraded modes where required.
- [ ] **USEQ-FB76AF44** — Design batch, queue, scheduled, and long-running work for visibility, retry, cancellation, reconciliation, and recovery.
- [ ] **USEQ-4B6D8A92** — Provide clear procedures for data correction and exceptional manual processing.
- [ ] **USEQ-BFBBB220** — Ensure customer-facing status and support information can be updated during control-plane failure.
- [ ] **USEQ-E25419F0** — Verify operational tasks are accessible and usable by intended staff.
- [ ] **USEQ-1BD6FF0C** — Measure operational toil and eliminate recurring avoidable manual work.

### People, roles, competence, and staffing readiness

- [ ] **USEQ-0BFDE64A** — Confirm every required operational and support role is staffed for the intended service hours and risk.
- [ ] **USEQ-9730C439** — Provide primary, secondary, management, specialist, supplier, and executive escalation coverage as needed.
- [ ] **USEQ-20DF5585** — Verify personnel have required access, equipment, credentials, training, authority, and contact information.
- [ ] **USEQ-501984F9** — Assess demonstrated competence for high-impact tasks.
- [ ] **USEQ-905419B1** — Provide practical exercises for deployments, rollback, incidents, failover, restoration, data repair, and customer communication.
- [ ] **USEQ-8E55BA3A** — Avoid reliance on one person, device, location, time zone, or undocumented skill.
- [ ] **USEQ-2B75F633** — Define handover, leave, offboarding, role-change, and contractor-transition procedures.
- [ ] **USEQ-3B38E750** — Ensure working hours, on-call load, and alert volume are sustainable.
- [ ] **USEQ-E46B4EE7** — Provide accessible accommodations and inclusive tools for operational personnel.
- [ ] **USEQ-8814F7F9** — Confirm support teams understand new behavior, limitations, known issues, eligibility, pricing, and customer remedies.
- [ ] **USEQ-A58F8861** — Confirm legal, privacy, security, safety, finance, and communications contacts can be reached during an event.

### Knowledge transfer and documentation acceptance

- [ ] **USEQ-E04DCCED** — Transfer architecture, dependencies, data flows, trust boundaries, operating assumptions, and failure modes.
- [ ] **USEQ-4B01668D** — Transfer deployment, configuration, migration, rollback, failover, restoration, and decommissioning procedures.
- [ ] **USEQ-A7916E27** — Transfer service objectives, dashboards, alerts, synthetic checks, logs, traces, and diagnostic queries.
- [ ] **USEQ-F42A3E39** — Transfer supplier contracts, support routes, quotas, escalation terms, and planned changes.
- [ ] **USEQ-B4F1C28A** — Transfer open incidents, defects, risks, exceptions, workarounds, debt, and expiry dates.
- [ ] **USEQ-37F13993** — Transfer data classifications, retention, rights, legal holds, privacy obligations, and sensitive operations.
- [ ] **USEQ-A245042A** — Transfer access-control models, privileged processes, break-glass methods, keys, certificates, and secret-recovery procedures.
- [ ] **USEQ-4DE727BA** — Provide customer-support scripts, troubleshooting guides, user communications, and accessibility support.
- [ ] **USEQ-0CA85D67** — Identify authoritative documents and archive superseded material.
- [ ] **USEQ-B0CE4AC1** — Validate documentation through task execution by the receiving team.
- [ ] **USEQ-C35F42F2** — Record unanswered questions and close or accept them before transition completion.
- [ ] **USEQ-C7B53656** — Ensure knowledge remains available after project personnel or suppliers depart.

### Environment, tool, and access readiness

- [ ] **USEQ-CC3D9793** — Verify production accounts, networks, domains, certificates, storage, compute, databases, queues, registries, and observability are provisioned correctly.
- [ ] **USEQ-DF105206** — Verify production configuration and feature-flag control are versioned, reviewed, and recoverable.
- [ ] **USEQ-C59CC418** — Verify least-privilege access for operations, support, responders, automation, and suppliers.
- [ ] **USEQ-CCFEF50D** — Test break-glass access and its audit trail.
- [ ] **USEQ-FFCE14A7** — Verify monitoring and paging work independently enough to detect failure of the primary service.
- [ ] **USEQ-EB2AAD9A** — Verify service desk, case management, status communication, incident coordination, and evidence repositories are ready.
- [ ] **USEQ-3D557FAC** — Verify backups, restore tooling, alternate regions, replacement equipment, and required licenses are available.
- [ ] **USEQ-CFCD783F** — Confirm quotas, limits, support plans, spending controls, and provider contacts match launch demand.
- [ ] **USEQ-6560F301** — Detect non-production credentials, debug paths, test recipients, test payment routes, and synthetic data in production configuration.
- [ ] **USEQ-1A30E6F1** — Confirm administrative tools are protected, observable, and recoverable.
- [ ] **USEQ-5559DBAC** — Confirm every operational tool has an owner and support path.

### Cutover, migration, and early-life support

- [ ] **USEQ-BB017E08** — Define the exact cutover sequence, dependencies, checkpoints, approvals, communications, and rollback points.
- [ ] **USEQ-A9709FE9** — Rehearse cutover using representative data, volume, permissions, integrations, and timings.
- [ ] **USEQ-2E28DB48** — Establish a protected recovery point before irreversible change.
- [ ] **USEQ-FFD30078** — Prevent overlapping authoritative writers unless conflict behavior is explicitly designed.
- [ ] **USEQ-9D527478** — Reconcile data, configuration, permissions, jobs, queues, caches, indexes, integrations, and financial totals after cutover.
- [ ] **USEQ-D127155D** — Validate critical user and operator journeys in production.
- [ ] **USEQ-A82DE51D** — Use progressive adoption or cohort-based transition where it materially reduces risk.
- [ ] **USEQ-B6713B9D** — Maintain enhanced monitoring and staffing for a defined early-life period.
- [ ] **USEQ-D329E9E1** — Define how incidents are attributed and owned during overlap.
- [ ] **USEQ-3EB7D7A3** — Track support volume, user failure, performance, error, data integrity, accessibility, abuse, and business outcomes.
- [ ] **USEQ-2E22A580** — Pause expansion or roll back when predefined thresholds are crossed.
- [ ] **USEQ-5AB3637F** — Avoid declaring transition complete before the observation and reconciliation period ends.
- [ ] **USEQ-8F0B7D76** — Record the final acceptance decision and residual obligations.

### Support, service desk, and customer readiness

- [ ] **USEQ-CA4E11E2** — Define supported users, service hours, channels, languages, accessibility accommodations, and response expectations.
- [ ] **USEQ-7593DD19** — Provide a clear path for urgent security, privacy, safety, fraud, accessibility, and data-integrity issues.
- [ ] **USEQ-7839483B** — Train support staff to authenticate requesters without unsafe data collection or social-engineering bypass.
- [ ] **USEQ-80257BB3** — Give support staff only the data and privileges needed for their role.
- [ ] **USEQ-1B3A2CEF** — Make support actions attributable and auditable.
- [ ] **USEQ-99F82FB7** — Provide escalation paths that reach people able to diagnose and change the service.
- [ ] **USEQ-E27FB6F6** — Define incident linkage so related support cases contribute to impact assessment.
- [ ] **USEQ-744C959B** — Prepare communications for maintenance, degradation, data issues, security events, delays, and rollback.
- [ ] **USEQ-F189D5B3** — Ensure public status information is accurate and does not expose sensitive detail.
- [ ] **USEQ-CD98BD1C** — Provide user migration, onboarding, recovery, and rollback guidance.
- [ ] **USEQ-8BD86B01** — Preserve customer commitments, entitlements, preferences, accessibility needs, and case history during service transition.
- [ ] **USEQ-A86B1DF6** — Measure resolution quality, recurrence, customer effort, accessibility, and escalation effectiveness rather than closure count alone.

### Organizational adoption and change enablement

- [ ] **USEQ-F5A6CAF3** — Identify stakeholders whose responsibilities, incentives, workflows, skills, or authority will change.
- [ ] **USEQ-C8079013** — Assess readiness, resistance, constraints, and likely unintended consequences.
- [ ] **USEQ-C14B308B** — Explain the purpose, impact, timeline, support, and decision rights honestly.
- [ ] **USEQ-AD2234A0** — Involve affected people in design, pilots, and transition decisions.
- [ ] **USEQ-175BD173** — Avoid using training to compensate for a fundamentally unusable or unsafe design.
- [ ] **USEQ-148C33A9** — Update policies, procedures, performance measures, roles, and incentives to match the new operating model.
- [ ] **USEQ-18C79536** — Provide accessible training and support in the languages and formats needed.
- [ ] **USEQ-4619C514** — Stage adoption so feedback can change the rollout.
- [ ] **USEQ-38E54A2A** — Preserve a safe fallback for users who cannot transition immediately where obligations permit.
- [ ] **USEQ-A809C9C6** — Monitor workarounds, shadow systems, abandonment, error, burden, and inequitable impact.
- [ ] **USEQ-2DD0C134** — Verify promised benefits and identify displaced costs or harms.
- [ ] **USEQ-627B9907** — Stop or redesign adoption when evidence shows the change is not producing acceptable outcomes.

### Acceptance, closure, and post-transition governance

- [ ] **USEQ-7CE19A5C** — Conduct formal operational acceptance against agreed criteria and current evidence.
- [ ] **USEQ-02C84E3D** — Have the receiving owner explicitly accept or reject responsibility.
- [ ] **USEQ-6808A5B2** — Record conditional acceptance with compensating controls, monitoring, owner, remediation date, and expiry.
- [ ] **USEQ-00F86BC0** — Do not treat project completion or budget exhaustion as operational acceptance.
- [ ] **USEQ-354C8516** — Confirm unresolved work has transferred to active owners and backlogs.
- [ ] **USEQ-E8D0E95A** — Confirm temporary access, environments, flags, data copies, supplier accounts, and elevated support are removed or governed.
- [ ] **USEQ-A157B5DF** — Confirm old systems, routes, jobs, credentials, and communications are decommissioned or intentionally retained.
- [ ] **USEQ-B049F5B4** — Review benefits, incidents, support demand, reliability, costs, human workload, and customer outcomes after transition.
- [ ] **USEQ-3C44939B** — Conduct a transition retrospective and verify corrective actions.
- [ ] **USEQ-E1759D28** — Reopen acceptance if material assumptions prove false or required controls cease to operate.
- [ ] **USEQ-66B7CDAC** — Treat missing ownership, untested recovery, inadequate staffing, failed migration reconciliation, or unsupported operations as no-go conditions.

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
