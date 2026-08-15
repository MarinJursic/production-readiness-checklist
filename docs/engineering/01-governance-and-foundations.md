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

## Final Gap Closure — Governance, Acquisition, Automated Decisions, and Rights

_Consolidated from `final consolidated corpus/01-governance-product-requirements-risk-ethics-lifecycle.md#Final Gap Closure — Governance, Acquisition, Automated Decisions, and Rights`; 215 non-duplicative controls._

### Compliance and obligations management

- [ ] **USEQ-0A917506** — Maintain one authoritative register of applicable laws, regulations, contracts, standards, licenses, policies, certifications, customer commitments, and internal obligations.
- [ ] **USEQ-FC410330** — Record the jurisdiction, authority, scope, effective date, owner, interpretation, evidence, review date, and change source for every obligation.
- [ ] **USEQ-228A7EAC** — Translate each applicable obligation into testable product, process, operational, documentation, and evidence requirements.
- [ ] **USEQ-0A0A5D07** — Distinguish mandatory obligations from voluntary commitments, guidance, aspirations, and marketing claims.
- [ ] **USEQ-5826CC12** — Resolve conflicts among legal, contractual, ethical, security, privacy, accessibility, safety, and product requirements through documented authorized decisions.
- [ ] **USEQ-0A75AF7B** — Monitor official sources for amendments, enforcement guidance, court decisions, regulator interpretations, and transition deadlines that can change applicability.
- [ ] **USEQ-4BD7C569** — Reassess obligations before entering a new market, serving a new user population, processing a new data type, changing a supplier, or materially changing functionality.
- [ ] **USEQ-899D534D** — Ensure compliance owners have sufficient authority, independence, competence, resources, and access to decision-makers.
- [ ] **USEQ-E5464005** — Prevent commercial or delivery pressure from suppressing a material compliance concern.
- [ ] **USEQ-E49A534A** — Require qualified review when an obligation is ambiguous, high impact, or outside the team's competence.
- [ ] **USEQ-33C53EA1** — Link obligations to requirements, risks, controls, tests, operating procedures, records, owners, and release decisions.
- [ ] **USEQ-56A86581** — Verify that controls operate in the actual production context rather than only existing in policy documents.
- [ ] **USEQ-68C68A3D** — Maintain evidence of compliance decisions and the factual basis used at the time.
- [ ] **USEQ-61C0266E** — Track nonconformities, corrective actions, due dates, owners, effectiveness checks, and recurrence.
- [ ] **USEQ-61985D0A** — Escalate repeated or systemic nonconformity beyond the team that created it.
- [ ] **USEQ-4F97A2C6** — Define which nonconformities block release, require customer notification, or trigger incident response.
- [ ] **USEQ-2E9F63DE** — Review whether accepted exceptions remain legally and contractually permissible.
- [ ] **USEQ-E34BC22A** — Ensure temporary waivers expire automatically and cannot become permanent through inattention.
- [ ] **USEQ-A98DF5A6** — Verify that customer-specific commitments are reflected in tenant, region, configuration, support, and reporting behavior.
- [ ] **USEQ-39F64583** — Prevent unsupported certifications, seals, compliance labels, security claims, privacy claims, accessibility claims, and quality claims.
- [ ] **USEQ-8BAFBF7E** — State the exact product, service, location, version, period, and control scope covered by every assurance claim.
- [ ] **USEQ-DE933B4A** — Track certification and audit surveillance dates, evidence windows, exclusions, dependencies, and renewal actions.
- [ ] **USEQ-D673A4A5** — Ensure changes cannot move the product outside a certified or contracted scope without review.
- [ ] **USEQ-EDCF9DA1** — Maintain a controlled process for regulator, auditor, customer, and lawful-authority requests.
- [ ] **USEQ-214C3BC6** — Protect privilege, confidentiality, personal data, trade secrets, and investigation integrity when responding to requests.
- [ ] **USEQ-1ECFF549** — Record why an obligation was judged Not Applicable and have the judgment reviewed by an authorized person.
- [ ] **USEQ-EBBC0B9E** — Reopen applicability judgments when facts, interpretations, product scope, or law change.
- [ ] **USEQ-D2C209FC** — Measure control effectiveness and outcome quality rather than counting policies, training sessions, or completed forms alone.
- [ ] **USEQ-5AFB222A** — Test compliance processes through sampling, walkthroughs, incident simulations, and independent review.
- [ ] **USEQ-D4694365** — Provide a confidential path for reporting suspected noncompliance and protect reporters from retaliation.
- [ ] **USEQ-EB301D83** — Preserve required records for the applicable period and dispose of them defensibly when the obligation ends.
- [ ] **USEQ-F60CED4C** — Ensure compliance evidence remains understandable and retrievable after personnel, vendor, or technology changes.
- [ ] **USEQ-EF95F68C** — Include compliance impacts in architecture decisions, backlog prioritization, procurement, release planning, and retirement.
- [ ] **USEQ-517D13D7** — Treat supplier-operated obligations as shared responsibilities and verify the supplier's evidence rather than assuming transfer of accountability.
- [ ] **USEQ-648076DD** — Review the compliance management system after material incidents, enforcement actions, audit findings, or recurring customer complaints.

### Acquisition, supply, outsourcing, and exit governance

- [ ] **USEQ-C097977F** — Define the business capability, outcomes, constraints, risks, data, interfaces, service levels, and acceptance criteria before selecting a supplier or product.
- [ ] **USEQ-E8FAD9D8** — Compare build, buy, reuse, partner, and retire options using lifecycle value, risk, reversibility, opportunity cost, and organizational capability.
- [ ] **USEQ-AD81EAAC** — Include migration, integration, operation, support, security, privacy, accessibility, continuity, and exit costs in acquisition decisions.
- [ ] **USEQ-53949496** — Assess whether the supplier's financial, operational, security, quality, legal, and staffing capacity is proportionate to criticality.
- [ ] **USEQ-37920856** — Identify fourth parties and other subcontractors that can materially affect the service or data.
- [ ] **USEQ-32E592BA** — Evaluate concentration risk where multiple critical capabilities depend on one provider, region, control plane, identity system, or corporate group.
- [ ] **USEQ-34730C2F** — Define contractual ownership and permitted use of source code, binaries, configuration, data, metadata, telemetry, models, prompts, documentation, and derived outputs.
- [ ] **USEQ-76A225D3** — Define which party is responsible for defects, vulnerabilities, incidents, accessibility barriers, data errors, legal requests, and customer support.
- [ ] **USEQ-15770275** — Include measurable service levels and service-level remedies that reflect real user and business impact.
- [ ] **USEQ-330DBEE1** — Specify security, privacy, accessibility, quality, continuity, records, audit, and incident-notification obligations in enforceable terms.
- [ ] **USEQ-9DE72FA9** — Require timely notice of material changes to ownership, subprocessors, locations, technology, security posture, data use, support, or end-of-life plans.
- [ ] **USEQ-CBCD7DAA** — Preserve rights to obtain sufficient evidence, audit reports, test results, architecture information, and incident details.
- [ ] **USEQ-49502674** — Ensure confidentiality terms do not prevent necessary security disclosure, regulatory reporting, or customer notification.
- [ ] **USEQ-7A7A29CE** — Define vulnerability disclosure, remediation targets, emergency coordination, and customer advisory responsibilities.
- [ ] **USEQ-3C3DC5B8** — Define data return, export format, completeness, integrity, metadata, provenance, and deletion obligations at termination.
- [ ] **USEQ-033D8FA9** — Require deletion evidence for data, credentials, backups, replicas, caches, support copies, and subprocessors where applicable.
- [ ] **USEQ-7A559BF6** — Define portability for identities, permissions, configuration, workflows, schemas, integrations, keys, logs, and historical records.
- [ ] **USEQ-CA6C9422** — Avoid contractual or technical lock-in that makes safe exit economically or operationally infeasible without explicit acceptance.
- [ ] **USEQ-7BC25239** — Identify proprietary extensions and document the consequence of losing them during migration.
- [ ] **USEQ-B97EFA06** — Require open, documented, or independently implementable interfaces when interoperability and exitability are material.
- [ ] **USEQ-123B91B1** — Define source escrow, continuity licensing, documentation access, or equivalent safeguards when supplier failure would be unacceptable.
- [ ] **USEQ-6C97F8ED** — Verify that escrowed or transferred material is complete, current, buildable, deployable, maintainable, and legally usable.
- [ ] **USEQ-AEBAE6E8** — Define transition assistance, knowledge transfer, overlap periods, support continuity, and key-person availability.
- [ ] **USEQ-237A0006** — Test exports, backups, migration tooling, and alternate-provider integration before an emergency forces their use.
- [ ] **USEQ-83E0FB27** — Maintain an exit plan for each critical supplier from the start of the relationship.
- [ ] **USEQ-6AB53D5A** — Assign an exit owner and estimate time, cost, dependencies, data movement, customer impact, and residual obligations.
- [ ] **USEQ-E6A9CCF0** — Identify supplier-provided controls that would disappear during exit and replace them before cutover.
- [ ] **USEQ-C8B933AD** — Prevent test or evaluation access from becoming unreviewed production dependency.
- [ ] **USEQ-57E8F571** — Require formal acceptance against agreed criteria before treating acquired software or services as production ready.
- [ ] **USEQ-C11220A0** — Separate supplier demonstration from independent verification using representative workloads and failure conditions.
- [ ] **USEQ-B1DA6120** — Verify licensing quantities, use rights, geographic restrictions, support rights, and renewal conditions before deployment.
- [ ] **USEQ-02A1C0FF** — Track notice periods, price changes, minimum commitments, automatic renewals, and termination windows.
- [ ] **USEQ-B6B75132** — Prevent procurement incentives from rewarding low initial price while hiding long-term risk or switching cost.
- [ ] **USEQ-F0EB40C4** — Define governance for supplier-requested emergency changes and exceptions.
- [ ] **USEQ-265DA109** — Review supplier performance using outcomes, incidents, risk trends, roadmap delivery, support quality, and corrective-action effectiveness.
- [ ] **USEQ-14A323CB** — Escalate chronic supplier underperformance and maintain objective replacement triggers.
- [ ] **USEQ-CA5F0120** — Ensure supplier personnel access uses individual identities, least privilege, time limits, logging, and prompt revocation.
- [ ] **USEQ-C744DA1F** — Remove supplier access, integrations, credentials, certificates, routes, accounts, and data when the relationship ends.
- [ ] **USEQ-DD784EF2** — Preserve contractual and technical evidence needed for disputes, audits, investigations, and transition.
- [ ] **USEQ-B4C136B8** — Reassess critical suppliers after acquisition, merger, restructuring, leadership change, breach, sanctions exposure, or material service degradation.
- [ ] **USEQ-66028C4E** — Document the residual risks accepted when a supplier cannot satisfy a required control.
- [ ] **USEQ-DD621797** — Ensure the risk owner accepting a supplier exception is independent of the commercial incentive to proceed.

### Organizational change, adoption, and benefits realization

- [ ] **USEQ-7D97FA29** — Identify every group whose work, rights, incentives, responsibilities, workload, or outcomes will change.
- [ ] **USEQ-F4BFF6A0** — Involve affected users, operators, support teams, administrators, and downstream partners early enough to influence design.
- [ ] **USEQ-7A6C75E2** — Document the current process, target process, transition states, retained exceptions, and retirement conditions.
- [ ] **USEQ-BD17C161** — Distinguish technical deployment from successful adoption and verified benefit realization.
- [ ] **USEQ-F637AD14** — Define adoption outcomes, behavioral indicators, user outcomes, business outcomes, safety indicators, and failure criteria before launch.
- [ ] **USEQ-B8EE5089** — Establish baseline measurements before changing the system or process.
- [ ] **USEQ-C4685A7E** — Validate that the change addresses the actual problem rather than merely introducing new technology.
- [ ] **USEQ-183A85DE** — Identify incentives that can cause gaming, unsafe shortcuts, shadow systems, duplicate work, or data-quality degradation.
- [ ] **USEQ-9C2558C8** — Assess workload transfer to users, support, operations, security, compliance, and external partners.
- [ ] **USEQ-4D532A00** — Ensure automation does not silently shift difficult work or risk to less powerful users.
- [ ] **USEQ-AECB2ABF** — Provide role-specific training, practice, support, and reference material before affected people are expected to perform new tasks.
- [ ] **USEQ-F60A847F** — Verify competence through observed performance or assessment rather than attendance alone.
- [ ] **USEQ-7908A553** — Provide accessible training and alternatives for different languages, abilities, schedules, locations, and levels of experience.
- [ ] **USEQ-871E097A** — Pilot changes with representative users, data, volume, constraints, and exception cases.
- [ ] **USEQ-6CD7A9F9** — Define safe coexistence between old and new processes during transition.
- [ ] **USEQ-05F6D1E7** — Prevent duplicate records, conflicting authorities, and inconsistent decisions while both processes operate.
- [ ] **USEQ-88FC79BA** — Define which system is authoritative at every transition stage.
- [ ] **USEQ-22C353B1** — Provide a supported path for correcting migration, adoption, and training errors.
- [ ] **USEQ-7A86F837** — Monitor adoption by user segment to detect exclusion, differential failure, or hidden burden.
- [ ] **USEQ-7AE11356** — Collect qualitative feedback as well as usage metrics.
- [ ] **USEQ-4FC73955** — Protect users from retaliation for reporting that a change is unsafe, unusable, inaccessible, or ineffective.
- [ ] **USEQ-F2A1587B** — Maintain a visible decision path for pausing, modifying, narrowing, or reversing the change.
- [ ] **USEQ-7B44BE46** — Avoid interpreting forced use as evidence of satisfaction or success.
- [ ] **USEQ-2184CE2C** — Confirm that users understand changed responsibilities, limits, escalation routes, and consequences.
- [ ] **USEQ-36B73F2A** — Ensure support capacity matches expected questions, incidents, and transition failures.
- [ ] **USEQ-DED02C8E** — Update policies, contracts, job aids, controls, dashboards, reports, and audit procedures with the change.
- [ ] **USEQ-C4FE996D** — Remove obsolete instructions and clearly label temporary transition guidance.
- [ ] **USEQ-6CD51BC3** — Measure intended benefits after stabilization and compare them with cost, harm, risk, and displaced work.
- [ ] **USEQ-5ECB509F** — Investigate benefits that appear only through metric definition changes, excluded populations, or unmeasured downstream costs.
- [ ] **USEQ-0C2A39BB** — Retire old systems and processes only after legal, operational, data, support, and recovery obligations are satisfied.
- [ ] **USEQ-13D789B8** — Preserve access to historical records needed to understand decisions made under the old process.
- [ ] **USEQ-5DEBD48E** — Capture lessons and update future change assumptions, estimates, training, and rollout methods.
- [ ] **USEQ-F54385C8** — Assign an accountable owner for benefits realization beyond the deployment date.
- [ ] **USEQ-6A646372** — Close the change only after adoption, outcome, support, control, and retirement criteria are met.

### Deterministic automated decisions, rules engines, scoring, and workflow logic

- [ ] **USEQ-2A4D03EE** — Inventory every automated rule, score, eligibility decision, prioritization, routing decision, threshold, and enforcement action that materially affects people or organizations.
- [ ] **USEQ-DE0C89CA** — Record the decision purpose, owner, legal basis, ethical rationale, affected populations, inputs, outputs, consequences, and appeal path.
- [ ] **USEQ-B989239A** — Classify decisions by potential harm, reversibility, scale, sensitivity, and degree of human reliance.
- [ ] **USEQ-BA4A6745** — Require stronger review, evidence, monitoring, and human oversight as decision impact increases.
- [ ] **USEQ-FA2A720B** — Define the authoritative source and version for every rule and decision table.
- [ ] **USEQ-0518839E** — Prevent undocumented business rules from existing only in code, spreadsheets, individual memory, or vendor configuration.
- [ ] **USEQ-A9E124F3** — Translate policy language into unambiguous executable logic with reviewed examples and counterexamples.
- [ ] **USEQ-7EA42112** — Record ambiguities, discretionary areas, conflicts, exceptions, and precedence among rules.
- [ ] **USEQ-CC7A27FE** — Make rule effective dates, expiry dates, jurisdiction, product scope, and population scope explicit.
- [ ] **USEQ-5BF321CF** — Ensure historical decisions can be reconstructed using the rules, data, configuration, and time context in force at the time.
- [ ] **USEQ-D0AF7A75** — Validate every input for accuracy, completeness, timeliness, provenance, semantic meaning, and permitted use.
- [ ] **USEQ-CE3B9271** — Define behavior for missing, disputed, stale, inconsistent, or low-confidence input.
- [ ] **USEQ-7881AC46** — Prevent proxies from defeating legal, ethical, fairness, or purpose limitations.
- [ ] **USEQ-A1EF5F3D** — Test rules across representative demographic, geographic, linguistic, accessibility, and socioeconomic groups where outcomes can differ.
- [ ] **USEQ-15661B6F** — Test boundary values and combinations near thresholds for disproportionate or unstable effects.
- [ ] **USEQ-69F2B791** — Detect feedback loops in which earlier automated outcomes shape future inputs and reinforce error or disadvantage.
- [ ] **USEQ-05AA60BA** — Evaluate whether a threshold creates cliff effects that should be replaced by review, graduated handling, or uncertainty disclosure.
- [ ] **USEQ-85E43B93** — Provide understandable reasons for consequential outcomes at a level appropriate to affected users and reviewers.
- [ ] **USEQ-3D5B49D0** — Distinguish an explanation of the actual decision from generic policy text.
- [ ] **USEQ-E94AE5FB** — Avoid exposing security-sensitive details while still providing meaningful reason and redress.
- [ ] **USEQ-ECA2F700** — Provide a timely, accessible, and effective path to challenge, correct, and appeal consequential decisions.
- [ ] **USEQ-E6116080** — Ensure reviewers can change an outcome and are not evaluated solely on agreement with automation.
- [ ] **USEQ-AAE859CB** — Prevent automation bias by presenting uncertainty, conflicting evidence, limitations, and alternatives.
- [ ] **USEQ-EE4F8C75** — Define when human review is mandatory and what competence, authority, time, and information the reviewer needs.
- [ ] **USEQ-FDC0790E** — Record overrides, reasons, outcomes, and patterns without discouraging justified intervention.
- [ ] **USEQ-5F638352** — Monitor override rates, appeal success, reversal, complaints, subgroup outcomes, drift, and downstream harm.
- [ ] **USEQ-6E4D7B2E** — Investigate high disagreement between automated outcomes and expert review.
- [ ] **USEQ-16C7AB3F** — Separate policy approval from technical implementation approval.
- [ ] **USEQ-9C43C510** — Require independent review for rules that affect rights, safety, money, access, employment, education, healthcare, housing, identity, or essential services.
- [ ] **USEQ-A3E0DD69** — Version rules, test cases, explanatory text, data mappings, and monitoring together.
- [ ] **USEQ-F9A87984** — Use dual control for emergency rule changes with high impact.
- [ ] **USEQ-507353D7** — Simulate the population and operational effect of material rule changes before activation.
- [ ] **USEQ-DA024B8E** — Use phased rollout, shadow evaluation, or retrospective analysis where direct rollout creates avoidable risk.
- [ ] **USEQ-C059B588** — Define rollback and remediation for incorrect decisions already made.
- [ ] **USEQ-72675725** — Identify every affected record and user when a rule defect is discovered.
- [ ] **USEQ-BAD34E15** — Correct downstream systems, notifications, balances, permissions, reports, and derived decisions after remediation.
- [ ] **USEQ-A3D3BEAB** — Notify affected parties when required or when doing so materially enables correction and trust.
- [ ] **USEQ-EF58B259** — Prevent users from manipulating decision inputs beyond intended discretion while avoiding unfair barriers to legitimate correction.
- [ ] **USEQ-E0356C82** — Test adversarial gaming, collusion, automation, and economic exploitation of rules.
- [ ] **USEQ-4F9E9ECD** — Ensure rate limits and fraud controls do not silently create discriminatory or inaccessible outcomes.
- [ ] **USEQ-1D490309** — Revalidate rules when source policy, law, population, data collection, product behavior, or operational context changes.
- [ ] **USEQ-4892D554** — Retire rules and related data when their purpose or authority ends.
- [ ] **USEQ-411D10D5** — Keep a complete audit trail of rule creation, approval, testing, deployment, execution, override, and retirement.
- [ ] **USEQ-CC58936A** — Publish appropriate transparency about the existence and purpose of material automated decision-making.
- [ ] **USEQ-2D7EA97F** — Prohibit claims of neutrality or objectivity when the rule embeds policy choices, value judgments, or incomplete data.

### Intellectual property, confidential information, and rights provenance

- [ ] **USEQ-83CB306E** — Identify the owner and permitted uses of every material codebase, library, data set, model, design, document, media asset, brand asset, and generated artifact.
- [ ] **USEQ-4161FF42** — Obtain written intellectual-property assignments or licenses from employees, contractors, partners, and contributors where required.
- [ ] **USEQ-89CB5A91** — Verify that contributors have authority to provide the material they contribute.
- [ ] **USEQ-F8B4E93A** — Maintain provenance records for externally sourced code, data, content, models, fonts, images, audio, video, and documentation.
- [ ] **USEQ-3F1EA654** — Review license compatibility before combining, distributing, hosting, modifying, or embedding components.
- [ ] **USEQ-E687E904** — Fulfil attribution, notice, source-offer, reciprocity, patent, trademark, and redistribution obligations.
- [ ] **USEQ-18BB4ACD** — Prevent confidential, proprietary, customer, employer, or restricted material from entering repositories or products without authorization.
- [ ] **USEQ-1E23B466** — Define rules for employee side projects, prior inventions, open-source contributions, and use of employer resources.
- [ ] **USEQ-46E6D74A** — Define contribution terms for public and private repositories.
- [ ] **USEQ-E4A212F5** — Review contributor license agreements or developer certificates for appropriateness and enforceability.
- [ ] **USEQ-B34707EE** — Track patents, defensive publications, trademarks, domain names, and other rights material to the product.
- [ ] **USEQ-50A7E640** — Avoid asserting patent or trademark rights without qualified review.
- [ ] **USEQ-64E5AE64** — Assess freedom-to-operate or infringement risk where the product, market, or contract warrants it.
- [ ] **USEQ-53DAE09A** — Use clean-room or other controlled development methods when independently reimplementing restricted interfaces or behavior.
- [ ] **USEQ-20AD6393** — Preserve evidence that clean-room boundaries and information controls were followed.
- [ ] **USEQ-42D1678E** — Ensure reverse engineering, interoperability work, scraping, and data acquisition comply with applicable rights and restrictions.
- [ ] **USEQ-E49C70F6** — Verify rights to use personal data, licensed data, public data, synthetic data, and derived data for the intended purpose.
- [ ] **USEQ-D70EDF93** — Do not assume public availability means unrestricted reuse.
- [ ] **USEQ-97C2833E** — Verify that training, evaluation, retrieval, analytics, and generated-output uses are covered by appropriate rights.
- [ ] **USEQ-29EE3F45** — Define rights and obligations for user-generated content, feedback, support submissions, and uploaded files.
- [ ] **USEQ-CC00B307** — Avoid terms that claim broader user-content rights than needed for the service.
- [ ] **USEQ-2482AD39** — Provide processes for copyright, trademark, privacy, and other rights complaints.
- [ ] **USEQ-FF4BC88F** — Prevent takedown processes from being abused without removing effective remedies for legitimate claimants.
- [ ] **USEQ-864C229F** — Track geographic, field-of-use, time, user-count, device-count, and environment restrictions in licenses.
- [ ] **USEQ-D38BB2A8** — Monitor license expiration, renewal, support, and termination dates.
- [ ] **USEQ-A6A85301** — Ensure backups, disaster recovery, testing, and archive copies remain within license rights.
- [ ] **USEQ-EA6F3090** — Include intellectual-property and data-rights obligations in supplier and acquisition reviews.
- [ ] **USEQ-FD996784** — Prevent AI-assisted tools from receiving confidential or restricted material contrary to policy or contract.
- [ ] **USEQ-38A675BA** — Independently review generated output for unauthorized copying, confidential disclosure, incompatible licensing, and provenance uncertainty.
- [ ] **USEQ-9A5356E6** — Preserve trade secrets through access control, confidentiality terms, labeling, training, and reasonable handling practices.
- [ ] **USEQ-ACCAC2B8** — Revoke access and recover confidential material during offboarding and supplier termination.
- [ ] **USEQ-9A5C1988** — Separate public, internal, confidential, restricted, and customer-owned materials in repositories and documentation systems.
- [ ] **USEQ-3B7D73D6** — Review product names, domains, icons, and branding for conflicts before public launch.
- [ ] **USEQ-4414E1D1** — Define rights needed for decommissioning, customer export, long-term support, security patching, and legal retention.
- [ ] **USEQ-EE05022B** — Maintain a response plan for infringement claims, license violations, ownership disputes, and accidental disclosure.
- [ ] **USEQ-2967E6ED** — Correct notices, distributions, source offers, and customer communications promptly when a rights defect is found.

### Speak-up culture, independence, conflicts, and decision integrity

- [ ] **USEQ-11CC414E** — Provide confidential and accessible channels for reporting quality, safety, security, privacy, accessibility, ethical, legal, and financial concerns.
- [ ] **USEQ-FD8476F7** — Prohibit retaliation against people who raise concerns in good faith.
- [ ] **USEQ-89CCF5C5** — Allow concerns to bypass the normal management chain when that chain has a conflict or is implicated.
- [ ] **USEQ-7C6CF5BB** — Define independent escalation to legal, compliance, security, safety, audit, risk, or governing bodies as appropriate.
- [ ] **USEQ-22A55748** — Record conflicts of interest for reviewers, approvers, auditors, procurement participants, and risk acceptors.
- [ ] **USEQ-7D5BF185** — Recuse or add independent review when a conflict can affect objectivity.
- [ ] **USEQ-F42C1888** — Prevent a person from accepting material risk solely to meet a target for which that person is rewarded.
- [ ] **USEQ-0081139A** — Separate implementation, verification, approval, and audit responsibilities in proportion to risk.
- [ ] **USEQ-5D1E5D62** — Ensure independent reviewers receive complete evidence and are not limited to management-selected summaries.
- [ ] **USEQ-FB0DCDDA** — Protect dissenting technical opinions and record unresolved disagreements in decision records.
- [ ] **USEQ-D06298BE** — Require decision-makers to address material contrary evidence explicitly.
- [ ] **USEQ-866CC7D2** — Do not use consensus pressure to convert unresolved risk into apparent agreement.
- [ ] **USEQ-C30D4EE0** — Define stop-work and release-blocking authority for credible high-impact concerns.
- [ ] **USEQ-CEECAF90** — Protect the ability to pause unsafe automation, deployment, data use, or customer communication.
- [ ] **USEQ-7EE0FA0B** — Investigate reports promptly, fairly, confidentially, and with appropriate evidence preservation.
- [ ] **USEQ-584CD846** — Communicate outcomes to reporters to the extent legally and operationally possible.
- [ ] **USEQ-22119AF5** — Track patterns of ignored warnings, repeated exceptions, near misses, and suppressed defects.
- [ ] **USEQ-F8C584BF** — Evaluate leaders on truthful risk communication and corrective action, not only delivery speed.
- [ ] **USEQ-7BB8745A** — Reward prevention, simplification, documentation, mentoring, and reliability work that may not create visible features.
- [ ] **USEQ-FA9114A6** — Prohibit manipulation of metrics, test scope, severity, evidence, or sampling to create a false appearance of readiness.
- [ ] **USEQ-FE265A51** — Require retrospective review of emergency approvals and break-glass decisions.
- [ ] **USEQ-6BA781DF** — Preserve auditability of who decided, what evidence was available, what alternatives existed, and why the decision was made.
- [ ] **USEQ-622B3666** — Review the effectiveness and trustworthiness of speak-up channels at least periodically and after serious incidents.

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
