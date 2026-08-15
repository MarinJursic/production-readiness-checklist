# Governance and foundations

_Phase 1 of 16 in the [complete engineering review](00-overview.md)._

Applicability, evidence, quality attributes, ownership, risk, ethics, suppliers, people, and continual improvement.

Assess every control as **Pass**, **Fail**, **Blocked**, or **Not Applicable**. A Pass needs current evidence for the exact product, revision, artifact, configuration, environment, and data state under review.

## Applicability, Evidence, Exceptions, and Audit Rules

_Consolidated from `quality standards/00-core/01-applicability-evidence-and-exceptions.md`; 61 non-duplicative controls._

### Applicability and evidence

- [ ] **USEQ-B14FDA8E** — Every control is assessed as **Pass**, **Fail**, **Blocked**, or **Not Applicable**.
- [ ] **USEQ-97B150FF** — Every Not Applicable result has a written, reviewed rationale tied to system context.
- [ ] **USEQ-64B6A3A8** — Every applicable control has one accountable owner and current, reproducible evidence.
- [ ] **USEQ-72088469** — Evidence identifies the exact product, source revision, artifact, configuration, environment, data state, and date assessed.
- [ ] **USEQ-CD2F9D5E** — Material changes invalidate affected evidence and trigger proportionate reassessment.
- [ ] **USEQ-6E42CC84** — Exceptions are time-bounded, risk-owned, monitored, documented, and automatically expire.
- [ ] **USEQ-2798897B** — A control operated by a supplier or platform remains in scope; responsibility may be shared but accountability is not outsourced.
- [ ] **USEQ-B9A1C5BC** — A serious unresolved failure in this category can block release regardless of aggregate checklist completion.

### Universal applicability rules

- [ ] **USEQ-E8F15536** — Establish the product, service, system, organization, release, and lifecycle boundaries being assessed.
- [ ] **USEQ-5657EBD9** — List every component, repository, artifact, environment, data store, interface, user class, role, supplier, location, and operational process in scope.
- [ ] **USEQ-C8C131DB** — Identify which conditional-domain modules are triggered by capabilities, users, data, jurisdiction, contract, or risk.
- [ ] **USEQ-75497572** — Assess every checklist item; omission is not a valid status.
- [ ] **USEQ-3A8A5F50** — Use only Pass, Fail, Blocked, or Not Applicable as assessment statuses.
- [ ] **USEQ-683B6299** — Treat Blocked as a failure for final approval unless an authorized decision explicitly states otherwise.
- [ ] **USEQ-8137C823** — Require a specific rationale for every Not Applicable decision.
- [ ] **USEQ-D7F6AF3B** — Reject Not Applicable rationales based only on inconvenience, cost, ownership elsewhere, tool limitations, or lack of implementation.
- [ ] **USEQ-F0272B02** — Assign one accountable owner to every applicable control and every exception.
- [ ] **USEQ-7375F93B** — Identify the qualified reviewer and preserve reviewer independence for material controls.
- [ ] **USEQ-5FB50410** — Tie evidence to the exact source revision, artifact digest, configuration, schema, flags, environment, dependency versions, and data state.
- [ ] **USEQ-4004F77F** — Reject evidence from a different release or environment unless equivalence is demonstrated.
- [ ] **USEQ-59870003** — Define evidence validity periods according to change rate, threat change, supplier change, and control volatility.
- [ ] **USEQ-569D2912** — Invalidate evidence after any material change that can affect the conclusion.
- [ ] **USEQ-DB6244DF** — Use sampling only with a documented population, method, coverage, confidence, limitations, and risk rationale.
- [ ] **USEQ-22F43D4B** — Prefer direct evidence of operation over policy statements, plans, intentions, or screenshots without context.
- [ ] **USEQ-56521362** — Require evidence that compensating controls work under the conditions in which the primary control is absent.
- [ ] **USEQ-41204E99** — Record residual risk after controls, not only inherent risk before controls.
- [ ] **USEQ-97E5A6C1** — Prohibit averaging or scoring that allows a critical failure to be offset by unrelated passes.
- [ ] **USEQ-08DC1879** — Define release-blocking and nonblocking severity before assessment results are known.
- [ ] **USEQ-01B6090D** — Escalate conflicts between quality attributes to an authorized multidisciplinary decision owner.
- [ ] **USEQ-08130ED1** — Preserve evidence in a durable, access-controlled, searchable, and immutable or tamper-evident form proportionate to risk.
- [ ] **USEQ-3AC76487** — Protect assessment evidence from exposing secrets, personal data, vulnerabilities, or customer information.
- [ ] **USEQ-0DAF2834** — Review provider and inherited controls while preserving customer-side responsibilities.
- [ ] **USEQ-38E46EF2** — Repeat assessment after incidents, escaped defects, major findings, jurisdiction changes, supplier changes, and significant scale change.
- [ ] **USEQ-35B75024** — Maintain historical assessments with historical releases and do not overwrite past conclusions.
- [ ] **USEQ-24057956** — Use the corpus as a baseline and add domain, organization, product, and jurisdiction controls without silently deleting universal controls.

### Minimum evidence record

- [ ] **USEQ-D135E510** — Control identifier and complete control text.
- [ ] **USEQ-ECC769B0** — Assessment status.
- [ ] **USEQ-4A3D6129** — Applicability rationale.
- [ ] **USEQ-6D3C345E** — Accountable owner.
- [ ] **USEQ-8D90ACE7** — Implementing owner or team.
- [ ] **USEQ-1015BBD0** — Qualified reviewer and approval.
- [ ] **USEQ-9FE92874** — Evidence location and access classification.
- [ ] **USEQ-5E2E1223** — Evidence date and expiry.
- [ ] **USEQ-F1DF60DF** — System, component, environment, source revision, artifact digest, configuration, schema, flags, and dependency context.
- [ ] **USEQ-F5B88725** — Test method, data, tools, versions, assumptions, and limitations.
- [ ] **USEQ-C8D18281** — Observed result and threshold.
- [ ] **USEQ-2989132F** — Linked defect, risk, exception, incident, or corrective action.
- [ ] **USEQ-89722999** — Reassessment trigger.

### Exception and risk-acceptance rules

- [ ] **USEQ-9EEBC544** — Identify the exact unmet control and why remediation is not complete.
- [ ] **USEQ-62B384A0** — Describe affected users, systems, data, operations, and business outcomes.
- [ ] **USEQ-4F2BE7CE** — Assess likelihood, impact, exploitability, reach, detectability, recoverability, and duration.
- [ ] **USEQ-D947B87C** — Document existing and proposed compensating controls.
- [ ] **USEQ-EB9477C4** — Demonstrate compensating controls with current evidence.
- [ ] **USEQ-0D775A08** — Define monitoring and alert thresholds specific to the accepted risk.
- [ ] **USEQ-C208D665** — Assign an authorized risk owner with authority over the consequences.
- [ ] **USEQ-27A4BEB6** — Set a remediation owner, funded plan, target date, and automatic expiry.
- [ ] **USEQ-4E4FE3D7** — Define conditions that immediately revoke the exception.
- [ ] **USEQ-238A650D** — Prevent exceptions from silently renewing or becoming permanent architecture.
- [ ] **USEQ-EB30F057** — Review cumulative risk from multiple individually accepted exceptions.
- [ ] **USEQ-841F7705** — Communicate the exception to operators, support, security, privacy, product, and customers where appropriate.
- [ ] **USEQ-DAE761C7** — Close the exception only after remediation and independent verification.

## Universal Software and Product Quality Attribute Model

_Consolidated from `quality standards/00-core/02-universal-quality-attribute-model.md`; 62 non-duplicative controls._

### Quality planning rules

- [ ] **USEQ-DBA8DAE2** — Select and rank quality attributes according to users, context, criticality, law, contract, threat, failure cost, and product strategy.
- [ ] **USEQ-706F2C82** — Define measurable scenarios and acceptance thresholds for every material quality attribute.
- [ ] **USEQ-8CB19F4F** — State trade-offs explicitly and prevent one quality attribute from being optimized by silently harming another.
- [ ] **USEQ-8D9F2BD3** — Use both product-quality evidence and real-world quality-in-use evidence.
- [ ] **USEQ-03A9E0A8** — Measure distributions, edge populations, and degraded conditions rather than averages alone.
- [ ] **USEQ-6DD9714D** — Define quality budgets and regression limits before implementation and launch.
- [ ] **USEQ-F30FF972** — Assign ownership for each quality objective and each cross-attribute conflict.
- [ ] **USEQ-A5FC207A** — Maintain traceability from quality objectives to architecture tactics, code, tests, monitoring, and operational decisions.
- [ ] **USEQ-4D15E15B** — Review quality objectives after material change, incident, scale increase, new user population, or new regulation.
- [ ] **USEQ-0EDC049D** — Do not declare overall quality from a single metric, tool score, certification, test level, or absence of reported defects.

### Product-quality characteristics

#### Functional suitability

- [ ] **USEQ-086B1E7E** — All required functions are present for the intended tasks and contexts.
- [ ] **USEQ-96E895D9** — Results, calculations, state transitions, permissions, and side effects are correct.
- [ ] **USEQ-4386DEDD** — Functions enable users to complete intended tasks without unnecessary or inappropriate behavior.
- [ ] **USEQ-AF45F27F** — Exceptional, degraded, recovery, cancellation, and prohibited behavior is defined and verified.

#### Performance efficiency

- [ ] **USEQ-83E0B216** — Latency, throughput, responsiveness, and completion time meet defined objectives under specified conditions.
- [ ] **USEQ-F3FC8AD5** — Compute, memory, storage, network, energy, and external-service use are proportionate to outcomes.
- [ ] **USEQ-31D05A91** — Capacity limits and safe overload behavior are known and tested.
- [ ] **USEQ-769066B0** — Performance remains acceptable across representative devices, networks, data volumes, tenants, and failure conditions.

#### Compatibility

- [ ] **USEQ-5104BA7F** — The product coexists safely with other software, workloads, users, and tenants in shared environments.
- [ ] **USEQ-09B7C249** — Information exchange preserves syntax, semantics, precision, identity, ordering, privacy, and authorization.
- [ ] **USEQ-2DF76BC4** — Supported versions, formats, clients, protocols, and environments interoperate as promised.
- [ ] **USEQ-BB1B1C78** — Incompatible changes use controlled negotiation, migration, deprecation, and retirement.

#### Interaction capability

- [ ] **USEQ-311C0D5C** — Users can recognize whether the product is appropriate to their goal.
- [ ] **USEQ-2AAD63D5** — Users can learn and operate the product effectively with reasonable effort.
- [ ] **USEQ-DB083136** — The product prevents, identifies, explains, and supports recovery from user error.
- [ ] **USEQ-E539A311** — The experience is inclusive across relevant abilities, language, literacy, device, input, and context.
- [ ] **USEQ-1865C898** — Help, feedback, status, and self-description make the system understandable.
- [ ] **USEQ-4BD5F349** — Engagement patterns preserve autonomy and do not manipulate users.

#### Reliability

- [ ] **USEQ-C3A94BE9** — The product performs required functions without unacceptable faults under specified conditions.
- [ ] **USEQ-427C1469** — The service is available according to user-centric objectives.
- [ ] **USEQ-02155A34** — Faults are contained and do not cascade unnecessarily.
- [ ] **USEQ-1156B06C** — State and service can recover within approved time, point, integrity, and continuity objectives.
- [ ] **USEQ-6DAFFC83** — Reliability remains valid during deployment, dependency failure, overload, maintenance, and disaster scenarios.

#### Security

- [ ] **USEQ-69A431A2** — Information and operations are confidential to unauthorized parties.
- [ ] **USEQ-ED5F3FB0** — Data, code, configuration, identities, and actions retain integrity.
- [ ] **USEQ-971B3218** — Material actions are attributable and appropriately auditable.
- [ ] **USEQ-88D972E4** — Identities, services, devices, artifacts, and messages are authenticated at required assurance.
- [ ] **USEQ-E3BBAE93** — The product resists attack, abuse, control bypass, compromise, and unauthorized recovery.
- [ ] **USEQ-CC0B63FE** — Security remains effective across development, supply chain, deployment, operation, and retirement.

#### Maintainability

- [ ] **USEQ-72137EBC** — The system is modular enough to localize understanding, change, testing, and failure.
- [ ] **USEQ-72AA3B4E** — Useful assets can be reused without unsafe coupling or semantic mismatch.
- [ ] **USEQ-636F1705** — Defects, behavior, dependencies, state, and performance can be analyzed efficiently.
- [ ] **USEQ-D5503076** — Changes can be implemented, reviewed, verified, deployed, and reversed safely.
- [ ] **USEQ-2A27E944** — Ownership, documentation, tests, tooling, and environments support future maintainers.

#### Flexibility

- [ ] **USEQ-0BDD9349** — The product can adapt to required environments, users, workloads, and configuration without uncontrolled redesign.
- [ ] **USEQ-B0202719** — The product scales to required demand while preserving quality and safe overload behavior.
- [ ] **USEQ-2329A8C0** — Installation, deployment, configuration, upgrade, and migration are reliable and supportable.
- [ ] **USEQ-3F04BFB9** — Critical components and providers can be replaced where continuity or strategy requires it.
- [ ] **USEQ-98F2CE0F** — The product is testable through controllable inputs, observable outputs, isolation, and reproducibility.

#### Safety

- [ ] **USEQ-C5CB4E01** — Operational constraints prevent the system from entering unacceptable hazardous states.
- [ ] **USEQ-543A1786** — Foreseeable hazards and affected people, property, business, software, and environment are identified.
- [ ] **USEQ-0737D264** — The product fails safe or limits harm when required assumptions or dependencies fail.
- [ ] **USEQ-1042DCD3** — Users and operators receive timely, perceptible, actionable warnings.
- [ ] **USEQ-8F9C2013** — Integration with other systems does not invalidate safety controls.
- [ ] **USEQ-709154B4** — Residual safety risk is independently reviewed and accepted only by authorized leadership.

### Quality in use

- [ ] **USEQ-B878B774** — Representative users achieve intended goals accurately and completely.
- [ ] **USEQ-0E7B4A58** — Users achieve goals with acceptable time, effort, cognitive load, cost, and resource use.
- [ ] **USEQ-6997B2CE** — Users experience acceptable trust, comfort, usefulness, and satisfaction without manipulation.
- [ ] **USEQ-15BC2C00** — Use avoids unacceptable economic, health, safety, privacy, security, social, and environmental risk.
- [ ] **USEQ-09CD9184** — Quality remains acceptable across the full range of intended and foreseeable contexts.
- [ ] **USEQ-3980950F** — The product supports users when context changes through interruption, disability, stress, low connectivity, device limitation, or emergency.
- [ ] **USEQ-3FCB5EB3** — Real-world outcomes are monitored and used to revise product-quality assumptions.

## Standards and Frameworks Crosswalk

_Consolidated from `quality standards/00-core/03-standards-and-frameworks-crosswalk.md`; 8 non-duplicative controls._

### How to use this crosswalk

- [ ] **USEQ-E2576411** — Treat this as a broad baseline catalog, not a claim that every listed reference applies to every product.
- [ ] **USEQ-429A53EC** — Create an applicability register based on system type, users, data, jurisdictions, sector, contracts, safety impact, suppliers, and deployment model.
- [ ] **USEQ-42342CFA** — Verify the current published edition directly with the issuing body before audit, procurement, certification, or compliance claims.
- [ ] **USEQ-16C2ACD8** — Distinguish published standards from drafts, committee drafts, proposed revisions, guidance, and awareness lists.
- [ ] **USEQ-AC5D537E** — Obtain and use the authoritative text when conformance is required; this corpus does not reproduce copyrighted requirements.
- [ ] **USEQ-1F1280E6** — Add sector-specific standards for medical, automotive, aviation, rail, industrial, financial, government, telecommunications, energy, education, and other regulated domains.
- [ ] **USEQ-46BFC86D** — Record which clauses or outcomes map to local controls, owners, evidence, and exceptions.
- [ ] **USEQ-9E85DA5A** — Review the register at least annually and after new jurisdictions, products, data, incidents, contracts, or standards revisions.

## Universal Definition of Done

_Consolidated from `quality standards/00-core/04-universal-definition-of-done.md`; 35 non-duplicative controls._

### A work item is not done until every applicable condition passes

#### Product and requirements

- [ ] **USEQ-5BE34669** — The user, business, regulatory, risk, or platform outcome is explicit.
- [ ] **USEQ-CD293CE7** — Scope, non-scope, assumptions, constraints, dependencies, and acceptance criteria are approved.
- [ ] **USEQ-C7ED66D6** — Functional and quality requirements are testable and traceable.
- [ ] **USEQ-FAB421AD** — Relevant user research and stakeholder validation are complete.
- [ ] **USEQ-4F358A61** — Privacy, security, accessibility, reliability, data, support, and lifecycle requirements are included.

#### Design and architecture

- [ ] **USEQ-8016290C** — The design satisfies requirements with justified trade-offs and no unnecessary complexity.
- [ ] **USEQ-F40624A5** — System, data, identity, trust, state, failure, operational, and migration implications are understood.
- [ ] **USEQ-FE5FCB37** — Material decisions and alternatives are recorded.
- [ ] **USEQ-41DCA413** — Threat, abuse, privacy, accessibility, and failure analysis is complete.
- [ ] **USEQ-4E5214F8** — Compatibility, rollout, rollback or roll-forward, and retirement are designed.

#### Implementation

- [ ] **USEQ-46552445** — Code is correct, readable, cohesive, appropriately simple, and follows approved standards.
- [ ] **USEQ-9FE5701E** — SOLID, DRY, KISS, YAGNI, reuse, and abstraction principles are applied contextually rather than mechanically.
- [ ] **USEQ-F83FD40A** — Input, output, identity, authorization, errors, concurrency, resources, configuration, and secrets are handled safely.
- [ ] **USEQ-2282B35E** — Dependencies are necessary, approved, supported, licensed, pinned, inventoried, and monitored.
- [ ] **USEQ-00E007D5** — Dead code, debug paths, temporary privileges, stale flags, and obsolete compatibility are removed or explicitly lifecycle-managed.

#### Verification

- [ ] **USEQ-4F6B6D16** — Unit, component, integration, contract, system, acceptance, regression, and exploratory testing are complete as applicable.
- [ ] **USEQ-7173A5C6** — Positive, negative, boundary, race, retry, failure, recovery, migration, and misuse cases are covered.
- [ ] **USEQ-C2EE6972** — Security, privacy, accessibility, performance, reliability, compatibility, and data quality objectives have evidence.
- [ ] **USEQ-F266538C** — Tests can fail for the defects they claim to detect and are not silently flaky or skipped.
- [ ] **USEQ-21B6514F** — The exact releasable artifact and production-representative configuration are verified.

#### Delivery and operations

- [ ] **USEQ-AB3A2C88** — Build provenance, artifact digest, SBOM, signature or attestation, and approvals are available as required.
- [ ] **USEQ-5FF05D93** — Configuration, migrations, feature flags, infrastructure, dependencies, and deployment order are versioned and reviewed.
- [ ] **USEQ-A8AA064C** — Monitoring, logs, metrics, traces, alerts, dashboards, SLOs, runbooks, capacity, and support are ready.
- [ ] **USEQ-1D2662B7** — Rollback, roll-forward, restoration, reconciliation, and incident response are tested.
- [ ] **USEQ-78BFB1B9** — Ownership, on-call, escalation, supplier support, and customer communication are ready.

#### Documentation and lifecycle

- [ ] **USEQ-75AAF559** — User, support, developer, API, data, architecture, test, security, privacy, and operational documentation is current as applicable.
- [ ] **USEQ-F62D0973** — Known limitations, residual risks, exceptions, and deferred work are explicit.
- [ ] **USEQ-12852D61** — Release notes, migration, deprecation, and support information are prepared.
- [ ] **USEQ-395CFEA4** — Temporary controls have owners and expiry dates.
- [ ] **USEQ-1ADFF5EF** — Evidence and decisions are retained and linked to the exact change.

#### Final acceptance

- [ ] **USEQ-3B007012** — All release blockers are resolved.
- [ ] **USEQ-8864F7C9** — No known residual risk exceeds authorized tolerance.
- [ ] **USEQ-AE200A41** — Required independent reviews and sign-offs are complete.
- [ ] **USEQ-6F9C0838** — The change can be detected, contained, disabled, recovered, and communicated if it fails.
- [ ] **USEQ-FB98DCA5** — Post-deployment validation and observation criteria are defined.

## Master No-Go Gates

_Consolidated from `quality standards/00-core/05-master-no-go-gates.md`; 25 non-duplicative controls._

### Absolute release and approval blockers

- [ ] **USEQ-6E292740** — A critical user, operator, administrator, recovery, cancellation, or deletion journey lacks successful end-to-end evidence.
- [ ] **USEQ-B366F417** — A material requirement or quality objective is undefined, unverified, or contradicted by observed behavior.
- [ ] **USEQ-770D5CDF** — Authentication, authorization, privileged access, account recovery, or tenant isolation lacks independent verification.
- [ ] **USEQ-AD254104** — Known critical or actively exploited exposure remains without demonstrated mitigation and authorized exception.
- [ ] **USEQ-C24545ED** — A secret, signing key, credential, token, personal record, or sensitive artifact was exposed and not revoked, rotated, contained, and investigated.
- [ ] **USEQ-0EB787DF** — A destructive, irreversible, or high-volume data change was not rehearsed and cannot be verified or recovered safely.
- [ ] **USEQ-B85F89C0** — Backup restoration, disaster recovery, rollback, or safe roll-forward cannot meet approved objectives.
- [ ] **USEQ-39949FBC** — Expected peak, burst, largest-tenant, abuse, or dependency-failure behavior is unknown.
- [ ] **USEQ-8124B601** — Critical user journeys cannot meet approved availability, latency, correctness, durability, freshness, accessibility, security, privacy, or safety objectives.
- [ ] **USEQ-26BF004F** — Monitoring cannot detect critical journey failure, material data corruption, control bypass, or loss of telemetry.
- [ ] **USEQ-AE082DBE** — No qualified and authorized responder is available for launch and early operation.
- [ ] **USEQ-5958B1BD** — Incident, escalation, customer communication, supplier escalation, or breach procedures required by risk do not exist.
- [ ] **USEQ-1F972B41** — A mandatory legal, regulatory, contractual, accessibility, privacy, payment, safety, or procurement obligation is unmet.
- [ ] **USEQ-9D2FAC3F** — The production artifact cannot be traced to reviewed source, declared dependencies, trusted build, tests, provenance, and approval.
- [ ] **USEQ-9B701F53** — An unsupported or end-of-life critical component lacks a tested compensating control and funded replacement plan.
- [ ] **USEQ-502289C2** — Known cross-user or cross-tenant disclosure, data corruption, duplicate charging, irreversible financial effect, or unsafe act remains possible without a tested control.
- [ ] **USEQ-C4B996D3** — A single person, credential, device, region, provider, undocumented procedure, or unavailable control plane is the only recovery path contrary to approved objectives.
- [ ] **USEQ-7AA7DE0D** — Evidence predates material untested change or refers to another artifact or configuration.
- [ ] **USEQ-08296C2C** — A blocker was relabeled, suppressed, waived, or marked Not Applicable without authorized evidence-based risk acceptance.
- [ ] **USEQ-5BCE4800** — The team cannot state concrete hold, stop, rollback, roll-forward, incident, and communication thresholds before deployment.
- [ ] **USEQ-6D298FCF** — A serious anomaly, test failure, data mismatch, or security signal remains unexplained.
- [ ] **USEQ-2094042A** — The only rationale for proceeding is deadline, revenue, sunk cost, executive pressure, or confidence without evidence.
- [ ] **USEQ-A2B2E114** — The product uses manipulative, discriminatory, inaccessible, unsafe, or unlawful behavior that has not been removed or explicitly and lawfully constrained.
- [ ] **USEQ-E72E7A13** — The system cannot preserve essential records, data integrity, or user rights through failure, migration, recovery, and retirement.
- [ ] **USEQ-4401E253** — The organization cannot identify the accountable owner of the product, service, data, security, privacy, operations, and residual risk.

## Final Evidence Package

_Consolidated from `quality standards/00-core/06-final-evidence-package.md`; 21 non-duplicative controls._

### Universal controls

- [ ] **USEQ-D4A97289** — Release identifier, source commits, tags, artifact digests, configuration versions, migration versions, and feature-flag states.
- [ ] **USEQ-0015A567** — Final scope, asset inventory, architecture diagrams, deployment diagrams, and data-flow diagrams.
- [ ] **USEQ-E15E9EEE** — Data inventory, classification, retention schedule, and privacy records.
- [ ] **USEQ-73265165** — Threat model, abuse-case assessment, risk register, and accepted exceptions.
- [ ] **USEQ-0820243E** — Applicable legal, privacy, accessibility, security, regulatory, and contractual requirements.
- [ ] **USEQ-1777943C** — Requirements-to-test and control-to-evidence traceability.
- [ ] **USEQ-153FAFE7** — Functional, integration, contract, system, regression, and user-acceptance test results.
- [ ] **USEQ-BB6A3F92** — Browser, device, localization, usability, and accessibility results.
- [ ] **USEQ-E01943D2** — Performance, load, stress, spike, endurance, and capacity results.
- [ ] **USEQ-5ECCFB7E** — Resilience, dependency-failure, failover, rollback, roll-forward, and recovery results.
- [ ] **USEQ-AE5734B3** — Security scan, code review, threat-validation, and penetration-test results.
- [ ] **USEQ-B39689E7** — SBOM, license report, build provenance, artifact signatures, and attestation evidence.
- [ ] **USEQ-0C482ED0** — Infrastructure, network, DNS, certificate, secrets, and production-configuration review.
- [ ] **USEQ-C3A340EE** — Backup, restoration, disaster-recovery, and continuity evidence.
- [ ] **USEQ-4A35B528** — Incident-response, on-call, alerting, runbook, and communications drill evidence.
- [ ] **USEQ-F475A3F9** — Monitoring, dashboard, synthetic-check, logging, audit, and telemetry review.
- [ ] **USEQ-0DBF3D1F** — Third-party risk, contract, subprocessor, and continuity review.
- [ ] **USEQ-41AB3A61** — Open defect, risk, and exception registers.
- [ ] **USEQ-8C3A67DF** — Every exception's owner, impact, compensating control, monitoring, remediation plan, and expiry.
- [ ] **USEQ-68373154** — Deployment, rollback or roll-forward, and post-deployment validation plans.
- [ ] **USEQ-291D1CF9** — Required approvals, sign-offs, and final decision rationale.

## Corpus Maintenance and Versioning

_Consolidated from `quality standards/00-core/08-corpus-maintenance-and-versioning.md`; 20 non-duplicative controls._

### Maintenance controls

- [ ] **USEQ-96C0CF6A** — Assign a corpus owner and qualified maintainers across product, engineering, UX, security, privacy, data, quality, delivery, operations, and legal domains.
- [ ] **USEQ-8A990F83** — Review the corpus at least annually and after major standards revisions.
- [ ] **USEQ-8B838859** — Review it after every material incident, audit finding, penetration-test finding, accessibility failure, privacy harm, safety event, or near miss.
- [ ] **USEQ-E5AB0C75** — Review it when entering a new jurisdiction, sector, user population, deployment model, or contractual regime.
- [ ] **USEQ-56E10E2D** — Review it after adding new sensitive data, identity, payment, AI, device, physical-control, or high-impact functionality.
- [ ] **USEQ-692CA8B7** — Track authoritative standards and mark published, amended, withdrawn, superseded, and draft status accurately.
- [ ] **USEQ-0A48A6ED** — Do not replace current published baselines with drafts unless a documented transition decision requires early adoption.
- [ ] **USEQ-1CD7806F** — Preserve historical corpus versions and the exact version used for each release or audit.
- [ ] **USEQ-6074CE61** — Record every added, changed, moved, or retired control with rationale and owner.
- [ ] **USEQ-9F4E9139** — Do not silently weaken controls during consolidation or reorganization.
- [ ] **USEQ-87104C15** — Use incidents, defects, user research, support demand, threat intelligence, supplier change, and operational evidence to add or refine controls.
- [ ] **USEQ-96ACBB49** — Review duplicate and overlapping controls while preserving standalone usability of category files.
- [ ] **USEQ-3726C504** — Validate all Markdown links, headings, checklist syntax, front matter, encoding, and relative paths automatically.
- [ ] **USEQ-24066ACE** — Validate that every category file contains purpose, applicability, universal controls, evidence expectations, and interpretation boundaries.
- [ ] **USEQ-ACB1E213** — Search for accidental secrets, personal data, internal URLs, and proprietary information before publication.
- [ ] **USEQ-D3AEAB87** — Review framework, vendor, language, and platform references and keep universal files technology-neutral.
- [ ] **USEQ-61E97483** — Place genuinely technology-specific controls only in clearly conditional modules.
- [ ] **USEQ-C6F40249** — Measure corpus use, control ambiguity, Not Applicable patterns, exception recurrence, and escaped defects.
- [ ] **USEQ-8CB08DFC** — Retire obsolete controls only after confirming the risk or requirement no longer exists.
- [ ] **USEQ-27EF115F** — Publish release notes and a migration guide for material corpus changes.

## Governance, Ownership, and Risk

_Consolidated from `quality standards/01-governance/01-governance-ownership-and-risk.md`; 19 non-duplicative controls._

### Ownership

- [ ] **USEQ-8C84BDA1** — Every production service has named business and engineering owners.
- [ ] **USEQ-C8627A8D** — Every third-party integration, dashboard, alert, runbook, scheduled job, and background process has an owner.
- [ ] **USEQ-2F823834** — Ownership remains valid during leave, holidays, and personnel changes.
- [ ] **USEQ-8272E813** — Escalation paths include primary, secondary, management, security, privacy/legal, and provider contacts.
- [ ] **USEQ-170D2EA2** — Offboarding and role changes remove or adjust access promptly.
- [ ] **USEQ-17D2A81C** — Ownership records are periodically reviewed.

### Risk assessment

- [ ] **USEQ-6B39A8A1** — Classify confidentiality, integrity, availability, privacy, financial, fraud, safety, regulatory, contractual, reputation, and customer-trust impact.
- [ ] **USEQ-21A06C7A** — Identify external attackers, abusive users, insiders, compromised accounts, bots, automated agents, social engineering, and accidental failure sources.
- [ ] **USEQ-ACE196FE** — Identify foreseeable misuse and abuse, not only intended use.
- [ ] **USEQ-4B9CB799** — Give every risk an owner, likelihood, impact, treatment, evidence, and review date.
- [ ] **USEQ-0B8D7893** — Record accepted residual risks and the authorized people who accept them.
- [ ] **USEQ-786705B4** — Verify that launch scope remains within the organization's approved risk appetite.
- [ ] **USEQ-9843CD0A** — Ensure exceptions expire and compensating controls are tested.

### Control framework

- [ ] **USEQ-AFB0927D** — Map security controls to a recognized verification standard.
- [ ] **USEQ-02B7631F** — Map privacy controls to legal and organizational requirements.
- [ ] **USEQ-0C8F7443** — Map accessibility controls to the selected conformance target.
- [ ] **USEQ-D29A0C99** — Identify and resolve conflicts between product, legal, security, privacy, accessibility, and technical requirements.
- [ ] **USEQ-2C632814** — Document resolutions, tradeoffs, and approvals.
- [ ] **USEQ-DDBA1140** — Periodically confirm the selected control level still matches application risk.

## Third-Party and Supplier Readiness

_Consolidated from `quality standards/01-governance/02-third-party-and-supplier-readiness.md`; 15 non-duplicative controls._

### Universal controls

- [ ] **USEQ-4246F390** — Record each provider's purpose, owner, data access, permissions, regions, subprocessors, and criticality.
- [ ] **USEQ-AAED2E02** — Security, privacy, resilience, accessibility, and compliance due diligence is proportionate to risk.
- [ ] **USEQ-9C004476** — Contracts cover service levels, security, permitted data use, incident notification, deletion, portability, audit, and termination.
- [ ] **USEQ-0D759FDC** — Provider credentials and permissions follow least privilege.
- [ ] **USEQ-FC267359** — Integration secrets are rotatable and revocable.
- [ ] **USEQ-D669D379** — Provider quotas, rate limits, concurrency limits, regional limits, and support commitments are understood.
- [ ] **USEQ-DE1C021A** — Provider availability and recovery commitments align with product objectives.
- [ ] **USEQ-5DA6E2FF** — Provider webhooks and callbacks are authenticated and replay-protected.
- [ ] **USEQ-FBBA350E** — Breaking-change and incident notifications reach an accountable owner.
- [ ] **USEQ-8FC152BB** — Provider outages, delays, malformed responses, throttling, and partial failures are tested.
- [ ] **USEQ-FC8F0D82** — Provider data-use and contract changes are monitored.
- [ ] **USEQ-DB18DC98** — A replacement, exit, continuity, and data-export plan exists for critical providers.
- [ ] **USEQ-01D417C3** — The application can identify records and users affected by a provider incident.
- [ ] **USEQ-615B275E** — Provider status and escalation contacts are included in runbooks.
- [ ] **USEQ-B009F0F4** — No provider can silently expand data collection or processing beyond approved behavior.

## Legal, Regulatory, and Contractual Applicability

_Consolidated from `quality standards/01-governance/03-legal-regulatory-and-contractual-applicability.md`; 27 non-duplicative controls._

### Universal controls

- [ ] **USEQ-8211B7E2** — Identify applicable privacy, data-protection, cybersecurity, breach-notification, and data-localization laws.
- [ ] **USEQ-029FDBE7** — Identify accessibility laws, public-sector rules, procurement requirements, and sector-specific accessibility duties.
- [ ] **USEQ-6B2941CB** — Identify consumer-protection, unfair-practice, pricing, subscription, renewal, cancellation, refund, and complaint requirements.
- [ ] **USEQ-532689B4** — Identify marketing, advertising, email, SMS, telemarketing, tracking, and consent requirements.
- [ ] **USEQ-1364FDC7** — Identify payment, banking, financial-service, insurance, securities, tax, invoicing, currency, and marketplace requirements.
- [ ] **USEQ-846F99B4** — Identify healthcare, education, employment, housing, government, critical-infrastructure, and other sector requirements.
- [ ] **USEQ-3C9DE828** — Identify children's privacy, age-assurance, age-appropriate-design, and child-safety requirements.
- [ ] **USEQ-2FF57224** — Identify identity-verification, anti-fraud, anti-money-laundering, sanctions, export-control, and restricted-country obligations.
- [ ] **USEQ-F29992F2** — Identify record-retention, legal-hold, evidence-preservation, and lawful-request obligations.
- [ ] **USEQ-F499A92F** — Identify international-transfer and data-residency requirements.
- [ ] **USEQ-BC302091** — Identify content, intellectual-property, copyright, moderation, intermediary, and platform obligations.
- [ ] **USEQ-BD48566C** — Identify algorithmic-decision, automated-decision, AI, transparency, explanation, human-review, and appeal requirements.
- [ ] **USEQ-45790745** — Identify contractual and customer-specific security, privacy, availability, audit, and notification requirements.
- [ ] **USEQ-A6E04178** — Terms of service match actual behavior.
- [ ] **USEQ-685C66D7** — Privacy and cookie notices match actual data and tracking behavior.
- [ ] **USEQ-9FBCC7DE** — Acceptable-use, content, moderation, and safety policies are enforceable and communicated.
- [ ] **USEQ-EC4D34E4** — Trademarks, media, data sets, fonts, models, and other content have appropriate rights.
- [ ] **USEQ-CAA051A6** — Security, compliance, privacy, sustainability, accessibility, and performance claims are accurate and supportable.
- [ ] **USEQ-4C0378A1** — Required contact, company, pricing, cancellation, complaint, and regulatory information is displayed.
- [ ] **USEQ-96B1F594** — Renewal, review, filing, audit, and reporting deadlines are tracked.

### Frequently triggered regulatory modules

- [ ] **USEQ-903DD079** — **Payments:** determine actual payment-security scope and apply the current applicable payment-card or financial-security standard.
- [ ] **USEQ-7690E19B** — **Health information:** determine whether health-information laws or contracts apply and implement required administrative, physical, and technical safeguards.
- [ ] **USEQ-4CF7E38B** — **Children:** determine whether the product is directed to children or knows it has child users and apply parental-consent, minimization, advertising, design, safety, and deletion requirements.
- [ ] **USEQ-63BBB657** — **Accessibility:** determine whether regional accessibility laws, public-sector rules, procurement requirements, or sector requirements apply.
- [ ] **USEQ-64DF55C0** — **Financial services:** determine licensing, disclosures, suitability, anti-fraud, recordkeeping, customer-protection, and transaction-monitoring obligations.
- [ ] **USEQ-C3784FD4** — **Employment, housing, education, credit, insurance, and essential services:** determine whether nondiscrimination, fairness, explanation, appeal, recordkeeping, and human-review obligations apply.
- [ ] **USEQ-224DF94F** — **Government and critical infrastructure:** determine whether additional security, hosting, accessibility, records, resilience, and supply-chain obligations apply.

## Quality Management and Continual Improvement

_Consolidated from `quality standards/01-governance/04-quality-management-and-continual-improvement.md`; 55 non-duplicative controls._

### Governance and lifecycle controls

- [ ] **USEQ-0897EF0F** — Define the category scope, intended outcomes, boundaries, stakeholders, users, and affected systems.
- [ ] **USEQ-F828BDE3** — Assign one accountable owner and identify all contributing and reviewing roles.
- [ ] **USEQ-7FE46690** — Inventory the artifacts, interfaces, data, dependencies, suppliers, environments, and lifecycle stages in scope.
- [ ] **USEQ-716AACBA** — Translate quality goals into explicit, testable requirements and acceptance criteria.
- [ ] **USEQ-F79AAA15** — Identify legal, contractual, security, privacy, accessibility, safety, reliability, and business constraints.
- [ ] **USEQ-6365A464** — Identify foreseeable misuse, failure modes, edge conditions, and conflicting stakeholder needs.
- [ ] **USEQ-770D6E54** — Record material assumptions and validate them with evidence before relying on them.
- [ ] **USEQ-0D80749F** — Record material design decisions, alternatives, trade-offs, consequences, and reversal conditions.
- [ ] **USEQ-06344248** — Prefer the least complex solution that satisfies current evidenced requirements and preserves necessary evolution paths.
- [ ] **USEQ-C7901DA7** — Maintain traceability from goals and risks through decisions, implementation, verification, deployment, and operational evidence.
- [ ] **USEQ-E075330B** — Apply independent review proportionate to impact and prevent authors from self-approving material controls.
- [ ] **USEQ-C7D1EC67** — Automate repeatable checks where automation is reliable, and retain human evaluation where judgment or user context is required.
- [ ] **USEQ-3B2B2AC9** — Test normal, negative, boundary, concurrency, failure, recovery, upgrade, rollback, and misuse behavior as applicable.
- [ ] **USEQ-9F30EF3C** — Use production-representative data shape, scale, topology, configuration, and dependency behavior for material verification.
- [ ] **USEQ-6D6113E1** — Measure outcomes and quality attributes rather than relying only on activity, output, or vanity metrics.
- [ ] **USEQ-C37654F3** — Make failures observable, diagnosable, containable, recoverable, and learnable.
- [ ] **USEQ-EC2F36AC** — Document operating procedures, limitations, ownership, support paths, and lifecycle expectations.
- [ ] **USEQ-C60E9C67** — Define deprecation, migration, retention, archival, and retirement behavior before those paths become urgent.
- [ ] **USEQ-69216AAA** — Review the controls after material change, incident, audit finding, user harm, supplier change, or new evidence.
- [ ] **USEQ-08BC8F27** — Convert recurring defects, support demand, security findings, and operational toil into permanent systemic improvements.

### Category-specific universal rules

- [ ] **USEQ-23D52699** — Establish a quality policy aligned with customer outcomes, safety, security, privacy, accessibility, reliability, and business objectives.
- [ ] **USEQ-740272E6** — Define measurable quality objectives at organization, portfolio, product, service, and release levels.
- [ ] **USEQ-DA6F3EA5** — Ensure quality objectives do not conflict with incentives that reward speed or output at the expense of outcomes.
- [ ] **USEQ-B73EE86B** — Define process owners, process inputs, expected outcomes, controls, records, and improvement mechanisms.
- [ ] **USEQ-715E3556** — Tailor processes explicitly to context without silently removing required outcomes.
- [ ] **USEQ-E77901AD** — Use prevention and early feedback before relying on downstream inspection.
- [ ] **USEQ-9A8916C2** — Integrate quality planning into discovery, requirements, architecture, implementation, verification, release, operation, and retirement.
- [ ] **USEQ-73C6FD52** — Use independent assurance for high-impact claims and controls.
- [ ] **USEQ-496084CD** — Analyze escapes, rework, incidents, complaints, support demand, and near misses as quality-system signals.
- [ ] **USEQ-5BC79E55** — Distinguish common-cause systemic problems from isolated special-cause events.
- [ ] **USEQ-B5A8BFB6** — Prioritize corrective action according to user harm, risk, recurrence, and systemic reach.
- [ ] **USEQ-CF27C8CE** — Verify corrective actions for effectiveness rather than closing them when work is merely completed.
- [ ] **USEQ-6B1C466A** — Prevent local optimizations from degrading end-to-end product or service quality.
- [ ] **USEQ-587925A5** — Maintain control of nonconforming work products and prevent accidental release or reuse.
- [ ] **USEQ-4246E9DC** — Use calibrated tools, test assets, data sets, and environments where measurement validity depends on them.
- [ ] **USEQ-2AFCA02A** — Review supplier quality and inherited controls as part of the same quality system.
- [ ] **USEQ-277CEB44** — Conduct periodic management review using evidence, trends, unresolved risk, customer feedback, and improvement results.
- [ ] **USEQ-225D52E9** — Preserve organizational learning so recurring failures do not reset with personnel or team changes.

### Required evidence

- [ ] **USEQ-928743A4** — Retain an approved scope and applicability record.
- [ ] **USEQ-2CD5DD6D** — Retain requirements and acceptance criteria linked to risks and intended outcomes.
- [ ] **USEQ-0BC431C9** — Retain decision records and review approvals for material choices.
- [ ] **USEQ-45C1DDB5** — Retain automated and manual verification results tied to the exact assessed revision.
- [ ] **USEQ-48BBA3DE** — Retain representative examples of successful, failed, boundary, and recovery behavior.
- [ ] **USEQ-E1B67CDC** — Retain current metrics, thresholds, trends, and the rationale for each threshold.
- [ ] **USEQ-321D7565** — Retain an open-defect and residual-risk register with owners and target dates.
- [ ] **USEQ-3358CD08** — Retain evidence that temporary exceptions and compensating controls are monitored and expire.
- [ ] **USEQ-397DD8D0** — Retain operational evidence showing the control continues to work after deployment.
- [ ] **USEQ-FB14A009** — Retain lessons and corrective actions from incidents, escapes, regressions, and near misses.

### Category no-go conditions

- [ ] **USEQ-A1AD3791** — Do not approve the category when a critical requirement lacks objective evidence.
- [ ] **USEQ-8A113D59** — Do not use checklist completion percentages to offset a material unresolved risk.
- [ ] **USEQ-9CC0C0B9** — Do not accept a control merely because a tool, supplier, framework, or platform claims to provide it.
- [ ] **USEQ-7D0D9C8E** — Do not mark a control Not Applicable solely because it is difficult, expensive, owned elsewhere, or not yet implemented.
- [ ] **USEQ-6993212D** — Do not proceed when the team cannot detect, contain, recover from, or communicate a foreseeable high-impact failure.
- [ ] **USEQ-961EF318** — Quality objectives are undefined, contradictory, or impossible to measure.
- [ ] **USEQ-F33B64FF** — Known systemic defects recur without an effective corrective-action program.

## Engineering Governance and Decision Rights

_Consolidated from `quality standards/01-governance/05-engineering-governance-and-decision-rights.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-41CEC06B** — Define which decisions are centralized, delegated, advisory, or prohibited.
- [ ] **USEQ-35C12A25** — Assign decision rights to roles with sufficient competence, authority, independence, and context.
- [ ] **USEQ-A56E3EB0** — Separate decision authority from personal status, seniority, or tool ownership.
- [ ] **USEQ-2048A133** — Define mandatory engineering guardrails and the risk basis for each guardrail.
- [ ] **USEQ-EA86F0CC** — Publish escalation paths for conflicts between product, engineering, security, privacy, reliability, accessibility, legal, and commercial goals.
- [ ] **USEQ-573E2FC1** — Require material decisions to state the problem, context, options, evidence, trade-offs, decision, owner, date, and review trigger.
- [ ] **USEQ-4D7CD9DF** — Make reversible decisions quickly while applying more rigor to irreversible or high-cost decisions.
- [ ] **USEQ-41A69859** — Require affected stakeholders to be consulted before decisions that transfer risk or operational burden to them.
- [ ] **USEQ-31ADB2B3** — Prevent architecture boards and review bodies from becoming unaccountable bottlenecks.
- [ ] **USEQ-1C45E26D** — Ensure exceptions identify the exact control waived, duration, compensating controls, owner, and exit plan.
- [ ] **USEQ-4E007EF6** — Review decision outcomes and retire decisions whose assumptions no longer hold.
- [ ] **USEQ-4384C947** — Make governance artifacts discoverable and understandable to the people expected to follow them.
- [ ] **USEQ-5657213E** — Ensure emergency decision paths preserve auditability and receive retrospective review.
- [ ] **USEQ-FED54377** — Define how teams challenge unsafe decisions without retaliation.
- [ ] **USEQ-77A6DDD9** — Measure governance by decision quality, risk outcomes, and flow—not meeting volume or document count.

## Project, Program, and Portfolio Management

_Consolidated from `quality standards/01-governance/06-project-program-and-portfolio-management.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-0B4B4DFD** — Connect every initiative to an explicit user, business, compliance, risk-reduction, or platform outcome.
- [ ] **USEQ-AA4F038D** — Define scope, non-scope, assumptions, constraints, dependencies, milestones, acceptance criteria, and completion conditions.
- [ ] **USEQ-9F8B6548** — Use outcome-based funding and review rather than treating initial estimates as immutable commitments.
- [ ] **USEQ-DF8B2B5D** — Maintain an integrated dependency map across products, teams, suppliers, data, infrastructure, and regulatory work.
- [ ] **USEQ-0BE4C2D3** — Identify the critical path and critical chain without hiding uncertainty behind single-point estimates.
- [ ] **USEQ-F2229F87** — Represent uncertainty with ranges, scenarios, confidence, and explicit assumptions.
- [ ] **USEQ-9BAB8E9B** — Reserve capacity for defects, security, reliability, maintenance, learning, and unplanned work.
- [ ] **USEQ-511069A8** — Prevent portfolio overcommitment from making every initiative simultaneously urgent.
- [ ] **USEQ-049681CB** — Define benefit owners and measure whether expected benefits are realized after delivery.
- [ ] **USEQ-EA0E51DC** — Stop or reshape initiatives when evidence invalidates the original case.
- [ ] **USEQ-B7629D08** — Coordinate cross-team interface, migration, release, and operational readiness obligations.
- [ ] **USEQ-BFDEE280** — Use transparent change control for material scope, schedule, cost, quality, or risk changes.
- [ ] **USEQ-2E9C12E9** — Treat unresolved staffing or skill gaps as delivery risks, not individual performance problems.
- [ ] **USEQ-9BF0B2DD** — Ensure project closure includes handover, operational acceptance, documentation, residual risk, and lessons learned.
- [ ] **USEQ-55341909** — Avoid rewarding deadline compliance when it produces hidden quality, security, or sustainability debt.

## Engineering Measurement and Metrics

_Consolidated from `quality standards/01-governance/07-engineering-measurement-and-metrics.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-DB947847** — Begin with a decision or information need before selecting a metric.
- [ ] **USEQ-DF84913E** — Define each metric's construct, population, unit, method, source, frequency, owner, limitations, and intended decisions.
- [ ] **USEQ-F1BEE0D8** — Validate that a metric measures the intended property rather than an easy proxy.
- [ ] **USEQ-BFBA2F5B** — Use balanced measures spanning outcomes, flow, quality, reliability, security, sustainability, and human factors.
- [ ] **USEQ-1163E0B6** — Use distributions and percentiles where averages hide important variation.
- [ ] **USEQ-2662F4B7** — Segment data when aggregate results can conceal harmed users, devices, regions, tenants, or cohorts.
- [ ] **USEQ-163BE262** — Define confidence, uncertainty, sampling, missing-data, and measurement-error treatment.
- [ ] **USEQ-7598EA63** — Protect measurement systems against manipulation, gaming, selection bias, and survivorship bias.
- [ ] **USEQ-5AAC4687** — Avoid using team-level system metrics as simplistic individual performance scores.
- [ ] **USEQ-13A8CCE3** — Review metric behavior after process or tooling changes that alter collection semantics.
- [ ] **USEQ-22853BEF** — Establish baselines and distinguish natural variation from meaningful change.
- [ ] **USEQ-9282AB20** — Predefine success, guardrail, stop, and rollback thresholds for experiments and launches.
- [ ] **USEQ-82522805** — Preserve raw evidence and transformation lineage sufficient to reproduce reported values.
- [ ] **USEQ-1C341AE4** — Retire metrics that no longer inform decisions or create harmful incentives.
- [ ] **USEQ-3DAEC856** — Pair quantitative evidence with qualitative user and operator evidence where context matters.

## People, Competence, Culture, and Sustainable Work

_Consolidated from `quality standards/01-governance/08-people-competence-culture-and-sustainable-work.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-CC59B784** — Define competency expectations for each material engineering, product, design, security, privacy, data, quality, and operations role.
- [ ] **USEQ-1BE6DBB2** — Assess competency using demonstrated outcomes and supervised practice rather than attendance alone.
- [ ] **USEQ-77B89A84** — Provide training before assigning high-impact responsibilities.
- [ ] **USEQ-7114A125** — Ensure reviewers have relevant expertise and sufficient time to review.
- [ ] **USEQ-805EACEA** — Use pairing, mentoring, rotation, and documentation to reduce single-person dependencies.
- [ ] **USEQ-6062519D** — Maintain succession and backup coverage for critical systems and responsibilities.
- [ ] **USEQ-D28F8B74** — Create psychologically safe channels for raising defects, risk, uncertainty, and ethical concerns.
- [ ] **USEQ-C1D6C65C** — Prohibit retaliation for stopping unsafe work or reporting material risk in good faith.
- [ ] **USEQ-AE14AF0D** — Avoid chronic overtime, interrupted leave, and unsustainable on-call as routine capacity strategies.
- [ ] **USEQ-866D6154** — Treat excessive toil and cognitive load as system defects to reduce.
- [ ] **USEQ-4DB71087** — Design team boundaries and ownership to minimize coordination cost while preserving end-to-end accountability.
- [ ] **USEQ-5F523DA1** — Ensure incentives reward quality, collaboration, learning, and customer outcomes rather than heroics and local output.
- [ ] **USEQ-456D1187** — Include accessibility, privacy, security, reliability, and incident skills in role-appropriate learning paths.
- [ ] **USEQ-7584EFCD** — Review access promptly when people join, change roles, take leave, or depart.
- [ ] **USEQ-FBBC8C19** — Capture critical tacit knowledge before reorganizations, outsourcing, or retirement.

## Ethics and Responsible Technology

_Consolidated from `quality standards/01-governance/09-ethics-and-responsible-technology.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-5A014C1A** — Identify people and communities who can be affected, including non-users and vulnerable groups.
- [ ] **USEQ-338D038D** — Assess intended benefits, foreseeable harms, misuse, power imbalances, and unequal distribution of risk.
- [ ] **USEQ-1216601A** — Avoid collecting, inferring, or exposing sensitive information without a necessary and justified purpose.
- [ ] **USEQ-00419A24** — Avoid manipulative, coercive, deceptive, addictive, or exploitative interaction patterns.
- [ ] **USEQ-63731FF6** — Provide meaningful transparency about consequential automated or policy-driven outcomes.
- [ ] **USEQ-8D2D4788** — Provide human review, appeal, correction, and redress where decisions can materially affect people.
- [ ] **USEQ-C9E9ADAC** — Test for systematically different quality or harm across relevant groups and contexts.
- [ ] **USEQ-219A351E** — Do not use protected or sensitive attributes as proxies without lawful, ethical, and evidenced justification.
- [ ] **USEQ-3C567287** — Limit surveillance, monitoring, scoring, and behavioral influence to proportionate purposes.
- [ ] **USEQ-0F61741E** — Protect user autonomy through reversible choices, clear defaults, and genuine refusal paths.
- [ ] **USEQ-4A92C813** — Define unacceptable uses and enforce them technically and contractually where practical.
- [ ] **USEQ-55927854** — Establish independent escalation for conflicts between commercial objectives and human welfare.
- [ ] **USEQ-7858A9FB** — Include environmental and labor impacts in major technology and supplier decisions.
- [ ] **USEQ-A7CFBEC2** — Reassess impact as usage, scale, context, data, models, and downstream integrations change.
- [ ] **USEQ-46335AA9** — Stop or constrain a capability when harms cannot be reduced below approved tolerance.

## Standards and source references

- [ISO/IEC/IEEE 12207:2026 — Software life cycle processes](https://www.iso.org/standard/90219.html)
- [ISO/IEC/IEEE 16085:2021 — Life-cycle risk management](https://www.iso.org/standard/74385.html)
- [ISO/IEC/IEEE 15289:2019 — Life-cycle information items](https://www.iso.org/standard/74909.html)
- [ISO 9001:2015 — Quality management systems](https://www.iso.org/standard/62085.html)
- [ISO/IEC 25010:2023 — Product quality model](https://www.iso.org/standard/78176.html)
- [ISO/IEC 25019:2023 — Quality-in-use model](https://www.iso.org/standard/78177.html)
- [ISO/IEC 25002:2024 — Quality model overview and usage](https://www.iso.org/standard/86715.html)
- [ISO/IEC/IEEE 15939:2017 — Measurement process](https://www.iso.org/standard/71197.html)
- [ISO/IEC/IEEE 15288:2023 — System life cycle processes](https://www.iso.org/standard/81702.html)
- [ISO/IEC/IEEE 29148:2018 — Requirements engineering](https://www.iso.org/standard/72089.html)
- [ISO/IEC/IEEE 42010:2022 — Architecture description](https://www.iso.org/standard/74393.html)
- [ISO/IEC/IEEE 16326:2019 — Project management](https://www.iso.org/standard/74397.html)
- [ISO/IEC 20246:2017 — Work product reviews](https://www.iso.org/standard/67597.html)
- [ISO/IEC/IEEE 29119-1:2022 — Software testing concepts](https://www.iso.org/standard/81291.html)
- [ISO/IEC/IEEE 29119-2:2021 — Test processes](https://www.iso.org/standard/79428.html)
- [ISO/IEC/IEEE 29119-3:2021 — Test documentation](https://www.iso.org/standard/79429.html)
- [ISO/IEC/IEEE 29119-4:2021 — Test techniques](https://www.iso.org/standard/79430.html)
- [ISO/IEC 27001:2022 — Information security management systems](https://www.iso.org/standard/27001)
- [ISO/IEC 27701:2025 — Privacy information management systems](https://www.iso.org/standard/85819.html)
- [ISO/IEC 29100:2024 — Privacy framework](https://www.iso.org/standard/85938.html)
- [ISO/IEC 25012:2008 — Data quality model](https://www.iso.org/standard/35736.html)
- [ISO/IEC 25024:2015 — Measurement of data quality](https://www.iso.org/standard/35749.html)
- [ISO/IEC 38500:2024 — Governance of IT](https://www.iso.org/standard/81684.html)
- [ISO/IEC 20000-1:2018 — Service management system requirements](https://www.iso.org/standard/70636.html)
- [ISO 22301:2019 — Business continuity management systems](https://www.iso.org/standard/75106.html)
- [ISO 31000:2018 — Risk management guidelines](https://www.iso.org/standard/65694.html)
- [ISO 9241-210:2019 — Human-centred design](https://www.iso.org/standard/77520.html)
- [ISO/IEC/IEEE 26514:2022 — Information for users](https://www.iso.org/standard/77451.html)
- [ISO/IEC 5055:2021 — Automated source code quality measures](https://www.iso.org/standard/80623.html)
- [ISO/IEC 21031:2024 — Software Carbon Intensity](https://www.iso.org/standard/86612.html)
- [ISO/IEC 42001:2023 — AI management systems](https://www.iso.org/standard/81230.html)
- [ISO/IEC 23894:2023 — AI risk management](https://www.iso.org/standard/77304.html)
- [NIST SP 800-218 v1.1 — Secure Software Development Framework](https://csrc.nist.gov/pubs/sp/800/218/final)
- [NIST SP 800-63-4 — Digital Identity Guidelines](https://pages.nist.gov/800-63-4/)
- [NIST SP 800-61 Rev. 3 — Incident Response](https://csrc.nist.gov/pubs/sp/800/61/r3/final)
- [NIST Cybersecurity Framework 2.0](https://www.nist.gov/cyberframework)
- [NIST Privacy Framework 1.0](https://www.nist.gov/privacy-framework)
- [NIST AI Risk Management Framework 1.0](https://www.nist.gov/itl/ai-risk-management-framework)
- [OWASP Application Security Verification Standard 5.0.0](https://owasp.org/www-project-application-security-verification-standard/)
- [OWASP Web Security Testing Guide 4.2](https://owasp.org/www-project-web-security-testing-guide/)
- [OWASP API Security Top 10 — 2023](https://owasp.org/API-Security/)
- [OWASP Top 10 — 2025](https://owasp.org/www-project-top-ten/)
- [W3C Web Content Accessibility Guidelines 2.2](https://www.w3.org/TR/WCAG22/)
- [SLSA Specification 1.2](https://slsa.dev/spec/v1.2/)
- [SPDX 3.0.1 / ISO/IEC 5962:2021](https://spdx.github.io/spdx-spec/v3.0.1/)
- [CycloneDX 1.7](https://cyclonedx.org/specification/overview/)
- [RFC 9700 — OAuth 2.0 Security Best Current Practice](https://www.rfc-editor.org/rfc/rfc9700)
- [RFC 8725 — JSON Web Token Best Current Practices](https://www.rfc-editor.org/rfc/rfc8725)
- [RFC 9110 — HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110)
- [RFC 9457 — Problem Details for HTTP APIs](https://www.rfc-editor.org/rfc/rfc9457)
- [OpenTelemetry Specifications](https://opentelemetry.io/docs/specs/)
- [Google SRE Workbook — Implementing SLOs](https://sre.google/workbook/implementing-slos/)
- [DORA software delivery performance research](https://dora.dev/guides/dora-metrics/)
- [OpenSSF Security Baseline and Best Practices](https://baseline.openssf.org/)

---

[Previous phase](00-overview.md) · [Next: Phase 2: Product and requirements](02-product-and-requirements.md)
