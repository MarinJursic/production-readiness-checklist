# Documentation and knowledge

_Phase 13 of 16 in the [complete engineering review](00-overview.md)._

Documentation governance, requirements and decisions, architecture, code, APIs, data, user help, runbooks, releases, and incidents.

Assess every control as **Pass**, **Fail**, **Blocked**, or **Not Applicable**. A Pass needs current evidence for the exact product, revision, artifact, configuration, environment, and data state under review.

## Documentation Governance

_Consolidated from `quality standards/14-documentation/01-documentation-governance.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-482D2055** — Define required information products by audience, decision, task, risk, and lifecycle stage.
- [ ] **USEQ-3223AD64** — Assign an owner, reviewers, source of truth, status, effective version, and review date to each material document.
- [ ] **USEQ-02294CB7** — Store documentation where intended users can discover and access it during normal work and incidents.
- [ ] **USEQ-A84035BD** — Use version control or equivalent history for material technical and operational information.
- [ ] **USEQ-987B8A80** — Link documentation to the product, service, release, interface, data set, control, or process it describes.
- [ ] **USEQ-ED3A1441** — Separate normative requirements from guidance, examples, history, and opinion.
- [ ] **USEQ-0404C9C9** — Use consistent templates where they improve completeness without forcing irrelevant content.
- [ ] **USEQ-2DB92982** — Write in clear language appropriate to the audience and define unavoidable terminology.
- [ ] **USEQ-832F56DB** — Make documentation accessible, searchable, navigable, and usable on supported devices.
- [ ] **USEQ-0916418C** — Protect secrets, personal data, vulnerabilities, customer information, and privileged procedures.
- [ ] **USEQ-FCEFDDB5** — Review documentation as part of the same change that alters behavior.
- [ ] **USEQ-EA8D292A** — Detect and label stale, draft, experimental, deprecated, superseded, and archived information.
- [ ] **USEQ-1F28B445** — Provide feedback and correction channels.
- [ ] **USEQ-46DF4D23** — Use executable examples, generated references, schema extraction, or validation where they reduce drift.
- [ ] **USEQ-91A9954E** — Do not generate documentation automatically when the result cannot be trusted or reviewed.
- [ ] **USEQ-BFCB0797** — Retain decision and audit records according to legal and operational needs.
- [ ] **USEQ-A055B3FA** — Remove duplicate and contradictory sources of truth.
- [ ] **USEQ-D2AFE23D** — Measure findability, task success, freshness, support deflection, errors, and user feedback.
- [ ] **USEQ-3ECAC2EB** — Preserve essential documentation outside a single provider or identity failure domain.
- [ ] **USEQ-C069FDB5** — Retire information when the supported product or obligation ends.

## Product, Requirements, and Decision Documentation

_Consolidated from `quality standards/14-documentation/02-product-requirements-and-decision-documentation.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-EDE4A090** — Document product purpose, target users, affected parties, outcomes, non-goals, constraints, and assumptions.
- [ ] **USEQ-A9393FF5** — Preserve research methods, participant context, findings, limitations, and consent handling.
- [ ] **USEQ-FDD0B29E** — Record requirements with identifiers, source, priority, rationale, acceptance criteria, status, and traceability.
- [ ] **USEQ-A8E8FB02** — Separate user need, requirement, design choice, implementation task, and test evidence.
- [ ] **USEQ-8B2B337A** — Document conflicting requirements and the authority and rationale for resolution.
- [ ] **USEQ-F6971DAB** — Record roadmap confidence, dependencies, commitments, options, and change history.
- [ ] **USEQ-CCDB14AF** — Document experiment hypothesis, population, metrics, guardrails, duration, analysis, limitations, and decision.
- [ ] **USEQ-7AC80095** — Preserve product decisions with alternatives, evidence, trade-offs, consequences, owner, and review trigger.
- [ ] **USEQ-C3ABA651** — Record legal, privacy, accessibility, security, safety, and reliability constraints in product artifacts.
- [ ] **USEQ-D5861A2B** — Keep pricing, entitlement, cancellation, retention, and support behavior documented consistently.
- [ ] **USEQ-B0A0C709** — Link acceptance evidence to the exact release and requirement.
- [ ] **USEQ-19C10937** — Record known limitations, deferred scope, residual risk, and customer impact honestly.
- [ ] **USEQ-27C3E4F4** — Update product and support documentation when behavior changes.
- [ ] **USEQ-5303F9AF** — Preserve rejected alternatives when their rationale may prevent repeated analysis.
- [ ] **USEQ-6065823A** — Make sensitive research and commercial strategy access-controlled.
- [ ] **USEQ-38F7FBD9** — Use diagrams, examples, scenarios, and state models where text alone is ambiguous.
- [ ] **USEQ-75CA82FD** — Ensure documentation is understandable to engineering, design, quality, operations, support, legal, and leadership.
- [ ] **USEQ-53A776EA** — Do not treat backlog items as complete requirements documentation by default.
- [ ] **USEQ-80892063** — Archive superseded plans while retaining effective history.
- [ ] **USEQ-6017DF2F** — Review whether delivered outcomes match the documented strategy and assumptions.

## Architecture and ADR Documentation

_Consolidated from `quality standards/14-documentation/03-architecture-and-adr-documentation.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-9F4859B3** — Identify the system of interest, environment, stakeholders, concerns, and scope of each architecture description.
- [ ] **USEQ-C331A72F** — Maintain context, container or component, deployment, data-flow, trust-boundary, and operational views as applicable.
- [ ] **USEQ-1294BBC5** — Document interfaces, dependencies, protocols, ownership, service objectives, and failure assumptions.
- [ ] **USEQ-8CB4B293** — Mark authoritative data, stateful components, consistency boundaries, and recovery paths.
- [ ] **USEQ-25A500F5** — Show identity, authorization, encryption, secret, privacy, and tenant boundaries.
- [ ] **USEQ-58746DC8** — Show queues, asynchronous workflows, scheduled work, caches, search, analytics, and third-party transfers.
- [ ] **USEQ-38FEC5D8** — Document quality-attribute scenarios and architectural tactics.
- [ ] **USEQ-6C132BDE** — Record architecture decisions with context, options, forces, evidence, decision, consequences, owner, status, and revisit trigger.
- [ ] **USEQ-F7ACB1C5** — Link decisions to requirements, threats, experiments, benchmarks, migrations, and incidents.
- [ ] **USEQ-51CB6303** — Keep diagrams and descriptions consistent with deployed reality.
- [ ] **USEQ-CD43693B** — Use stable names and legends and avoid diagrams that require undocumented tribal knowledge.
- [ ] **USEQ-78AF3880** — Separate current, target, transitional, and historical architecture.
- [ ] **USEQ-906B5632** — Document intentional constraint violations and remediation plans.
- [ ] **USEQ-A8C8399B** — Include operational and recovery dependencies that are often omitted from functional diagrams.
- [ ] **USEQ-407E94C2** — Make architecture discoverable to implementers and responders.
- [ ] **USEQ-0EA4D9A7** — Protect sensitive topology and security details while keeping necessary knowledge available.
- [ ] **USEQ-16B42B20** — Automate conformance or drift checks for critical invariants where practical.
- [ ] **USEQ-3520DA55** — Retire superseded decisions without deleting rationale.
- [ ] **USEQ-606151BB** — Review after material change and at planned intervals.
- [ ] **USEQ-84233857** — Test documentation usefulness during onboarding, design review, and incident response.

## Code, API, and Data Documentation

_Consolidated from `quality standards/14-documentation/04-code-api-and-data-documentation.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-D4D4DA52** — Provide a reliable entry point explaining purpose, ownership, prerequisites, setup, build, test, run, and support.
- [ ] **USEQ-A69B684D** — Document public modules, interfaces, functions, commands, events, and configuration according to consumer need.
- [ ] **USEQ-6D750AA3** — State inputs, outputs, units, formats, errors, side effects, authorization, idempotency, ordering, and resource behavior.
- [ ] **USEQ-2CFFD683** — Document API authentication, schemas, pagination, filtering, rate, quota, versioning, deprecation, and examples.
- [ ] **USEQ-CA239954** — Use generated API and schema references only when enriched with semantic and failure guidance.
- [ ] **USEQ-9EA30976** — Keep examples executable or automatically validated.
- [ ] **USEQ-06617A33** — Never include real secrets, credentials, personal data, or unsafe production endpoints in examples.
- [ ] **USEQ-DD2CC11E** — Document data entities, fields, identifiers, null semantics, units, quality, source, lineage, retention, and ownership.
- [ ] **USEQ-F4E837FE** — Document migrations, backward compatibility, mixed-version behavior, and consumer actions.
- [ ] **USEQ-DF76D73B** — Explain non-obvious algorithms, invariants, concurrency, performance, and security constraints.
- [ ] **USEQ-4A0C3F91** — Keep comments near code for local reasoning and external documentation for stable consumer contracts.
- [ ] **USEQ-21E0B53A** — Link runbooks and troubleshooting for operationally significant components.
- [ ] **USEQ-B2599C93** — Document extension and contribution rules for reusable components.
- [ ] **USEQ-71737F7A** — Label experimental, internal, unstable, deprecated, and unsupported interfaces clearly.
- [ ] **USEQ-C42FEE7F** — Provide change logs for consumer-relevant behavior.
- [ ] **USEQ-0BA694D6** — Avoid duplicating the same normative contract in multiple manually maintained locations.
- [ ] **USEQ-319F5870** — Test documentation from a clean consumer environment.
- [ ] **USEQ-4989E2A5** — Record supported toolchain and platform versions.
- [ ] **USEQ-D6EFE595** — Provide a security contact and vulnerability-reporting path for public packages or APIs.
- [ ] **USEQ-9FC49F3D** — Review documentation during code and release review.

## User Help, Training, and Support Content

_Consolidated from `quality standards/14-documentation/05-user-help-training-and-support-content.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-51331060** — Base content on user tasks, context, knowledge, risk, language, and support needs.
- [ ] **USEQ-D3ACA67F** — Provide concise in-context guidance before forcing users to leave the task.
- [ ] **USEQ-280A2D3E** — Explain prerequisites, permissions, cost, privacy, consequences, and expected outcomes.
- [ ] **USEQ-02E9D4BC** — Use step sequences that match current product labels and behavior.
- [ ] **USEQ-20D1BA9B** — Cover normal use, errors, recovery, cancellation, deletion, accessibility, security, and support escalation.
- [ ] **USEQ-A17A56AB** — Use screenshots and media only when maintained, accessible, localized, and genuinely helpful.
- [ ] **USEQ-D5C35624** — Provide text alternatives, captions, transcripts, headings, landmarks, and accessible document formats.
- [ ] **USEQ-93775EA6** — Avoid making critical instructions available only through video, image, color, or hover.
- [ ] **USEQ-E10C9176** — Write troubleshooting from symptoms to safe diagnosis and recovery.
- [ ] **USEQ-BFA3670F** — Distinguish user-actionable issues from incidents that require support or operator intervention.
- [ ] **USEQ-7A1E53B9** — Warn before destructive, irreversible, or security-sensitive actions.
- [ ] **USEQ-94EB3B4A** — Avoid exposing internal diagnostics, secrets, or unsafe workarounds.
- [ ] **USEQ-6056DDE7** — Define content version, product version, locale, owner, and review date.
- [ ] **USEQ-F5953572** — Test instructions with representative users in a clean environment.
- [ ] **USEQ-EC281E9C** — Track failed searches, abandonment, support contacts, and feedback to improve content.
- [ ] **USEQ-EC821EEE** — Keep legal, pricing, privacy, and security statements aligned with authoritative policy.
- [ ] **USEQ-F2B8BABF** — Provide offline or alternate access for information needed during outages where appropriate.
- [ ] **USEQ-11BCC1E3** — Archive obsolete content and remove it from search recommendations.
- [ ] **USEQ-10AFC3D2** — Coordinate release notes, in-product guidance, support scripts, and knowledge articles.
- [ ] **USEQ-41A50E2B** — Use training assessment to verify competence for high-impact tasks.

## Operations, Runbooks, Release, and Incident Documentation

_Consolidated from `quality standards/14-documentation/06-operations-runbooks-release-and-incident-documentation.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-69C45485** — Provide runbooks for every paging alert, critical dependency, common failure, recovery operation, and high-risk manual task.
- [ ] **USEQ-02EBD336** — State purpose, trigger, prerequisites, access, safety warnings, expected evidence, and escalation.
- [ ] **USEQ-30B49FB6** — Use commands and queries that are reviewed, scoped, idempotent where possible, and clearly mark destructive steps.
- [ ] **USEQ-BA28ABA4** — Include diagnosis, containment, mitigation, rollback, roll-forward, recovery, validation, and communication.
- [ ] **USEQ-6AF3EC1F** — Define stop conditions and when to declare or escalate an incident.
- [ ] **USEQ-95A45226** — Keep contact, supplier, status, and alternate communication information current.
- [ ] **USEQ-6CDA3BBC** — Document deployment artifacts, configuration, migrations, flags, health gates, rollback, and post-deployment verification.
- [ ] **USEQ-642AD6F8** — Record releases with source, digest, approvals, changes, known issues, risk, and support impact.
- [ ] **USEQ-E1D5C77F** — Maintain incident timelines, hypotheses, decisions, actions, impact, evidence, communications, and recovery verification.
- [ ] **USEQ-2BDDA007** — Separate real-time incident notes from the later reviewed postmortem.
- [ ] **USEQ-461D6C18** — Preserve evidence and sensitive information with appropriate access and retention.
- [ ] **USEQ-29CFDA1F** — Test runbooks through drills and actual use and capture corrections.
- [ ] **USEQ-9444CDD3** — Avoid requiring unavailable systems, identities, or networks during disaster recovery.
- [ ] **USEQ-4EA55BD1** — Make time-sensitive procedures usable under stress through concise structure and verified automation.
- [ ] **USEQ-C13E3C06** — Provide data repair and reconciliation procedures with safeguards.
- [ ] **USEQ-D97EB15E** — Document backup restoration, failover, failback, and clean-environment recovery.
- [ ] **USEQ-381295F2** — Keep runbooks aligned with current architecture and permissions.
- [ ] **USEQ-44B2EFFA** — Record temporary incident changes and ensure cleanup.
- [ ] **USEQ-BCD83655** — Track corrective actions to verified completion.
- [ ] **USEQ-02E3B3DC** — Retire obsolete runbooks and redirect users to the effective procedure.

## Final Gap Closure — Knowledge Management, Competence, Learning, and Continuity

_Consolidated from `final consolidated corpus/09-documentation-knowledge-training-support-service-management.md#Final Gap Closure — Knowledge Management, Competence, Learning, and Continuity`; 121 non-duplicative controls._

### Knowledge-management governance

- [ ] **USEQ-37E3F8D7** — Define which knowledge is critical to product value, engineering, operation, security, privacy, compliance, support, and continuity.
- [ ] **USEQ-70E7772F** — Assign accountable owners for knowledge domains and authoritative sources.
- [ ] **USEQ-C14D3998** — Establish objectives for knowledge creation, capture, validation, sharing, reuse, protection, retention, and retirement.
- [ ] **USEQ-41D9395F** — Align knowledge management with organizational strategy, risk, lifecycle, and user needs.
- [ ] **USEQ-145DDFF2** — Include employees, contractors, suppliers, partners, customers, and communities where they create or rely on critical knowledge.
- [ ] **USEQ-FC217284** — Distinguish public, internal, confidential, restricted, personal, customer-owned, legally privileged, and export-controlled knowledge.
- [ ] **USEQ-18C257FA** — Ensure knowledge practices respect intellectual property, privacy, security, competition, and employment obligations.
- [ ] **USEQ-415D2615** — Provide resources, roles, tools, and incentives for maintaining knowledge, not merely producing new artifacts.
- [ ] **USEQ-633F8CDB** — Measure whether people can find, understand, trust, and apply needed knowledge.
- [ ] **USEQ-D40AFC25** — Review the knowledge-management system after incidents, reorganizations, acquisitions, supplier changes, and key-person departures.

### Critical-knowledge inventory and risk

- [ ] **USEQ-8EEF087F** — Identify knowledge whose loss, inaccuracy, or inaccessibility could cause material failure or delay.
- [ ] **USEQ-2A3A7103** — Include architecture rationale, domain rules, data semantics, security assumptions, operational procedures, supplier knowledge, and historical constraints.
- [ ] **USEQ-1EDE2DA3** — Identify tacit knowledge held by one person or a small group.
- [ ] **USEQ-27857C7E** — Identify knowledge dependent on one vendor, proprietary tool, unsupported format, inaccessible repository, or expiring credential.
- [ ] **USEQ-C13A3455** — Assess knowledge concentration, obsolescence, misinformation, confidentiality, and availability risks.
- [ ] **USEQ-D2F9BDB2** — Record the owner, audience, source, validation method, location, sensitivity, dependencies, and review cycle.
- [ ] **USEQ-9EB594E8** — Prioritize continuity actions according to impact and replacement difficulty.
- [ ] **USEQ-EF4FC8D9** — Track key-person and key-supplier risk explicitly.
- [ ] **USEQ-126AA834** — Reassess critical knowledge when systems, markets, obligations, or organizational structures change.

### Creation, capture, and validation

- [ ] **USEQ-43FCA7C7** — Capture decisions, rationale, alternatives, assumptions, and consequences at the time they are made.
- [ ] **USEQ-6D7905EF** — Capture lessons from design, implementation, testing, incidents, support, migrations, audits, and retirement.
- [ ] **USEQ-43428C3C** — Convert tacit knowledge into usable artifacts, demonstrations, mentoring, automation, or shared practice where appropriate.
- [ ] **USEQ-237035A5** — Avoid documenting only ideal procedures while omitting real constraints and known workarounds.
- [ ] **USEQ-21FFD735** — Record the evidence and authority supporting material statements.
- [ ] **USEQ-2A391824** — Distinguish fact, decision, opinion, hypothesis, draft, obsolete guidance, and unresolved question.
- [ ] **USEQ-9E8B3204** — Validate critical knowledge through peer review, execution, testing, or independent confirmation.
- [ ] **USEQ-D49EE764** — Include authorship, ownership, version, date, scope, status, and review date.
- [ ] **USEQ-916C3F57** — Preserve links to source requirements, systems, incidents, changes, and evidence.
- [ ] **USEQ-D67C717E** — Avoid creating documentation solely to satisfy a gate when nobody can use or maintain it.

### Findability, access, and usability

- [ ] **USEQ-25471447** — Provide one discoverable entry point or catalog for authoritative engineering and operational knowledge.
- [ ] **USEQ-5EB04CC7** — Use consistent titles, terminology, taxonomy, tags, identifiers, and ownership metadata.
- [ ] **USEQ-F92EB50C** — Make search return authoritative and current material before duplicates and obsolete copies.
- [ ] **USEQ-7786CD26** — Clearly mark superseded, draft, experimental, and archived content.
- [ ] **USEQ-82735880** — Redirect or link from obsolete locations to current guidance where practical.
- [ ] **USEQ-1D316628** — Keep access no broader than necessary but no narrower than legitimate work requires.
- [ ] **USEQ-CEA15543** — Ensure critical procedures remain available during network, identity, collaboration-tool, or provider failure.
- [ ] **USEQ-2DD8ED47** — Provide offline or alternate access to emergency information where appropriate.
- [ ] **USEQ-BEE02B64** — Make knowledge accessible to people with disabilities.
- [ ] **USEQ-2D68989A** — Use language, examples, and structure appropriate to the intended audience’s competence.
- [ ] **USEQ-0CE9689B** — Provide multilingual material where operational or customer obligations require it.
- [ ] **USEQ-C362DCCB** — Avoid hidden knowledge in private messages, personal notes, local devices, or inaccessible meeting recordings.
- [ ] **USEQ-DD8E42D2** — Ensure machine-readable knowledge does not replace necessary human-understandable explanation.

### Maintenance, freshness, and retirement

- [ ] **USEQ-331D2821** — Assign review intervals based on criticality and rate of change.
- [ ] **USEQ-76C53F97** — Trigger review after material product, architecture, supplier, policy, incident, or regulatory change.
- [ ] **USEQ-90E787BB** — Detect broken links, missing owners, expired review dates, inaccessible repositories, and conflicting guidance.
- [ ] **USEQ-2210FD4E** — Make documentation change part of the same workflow as the product or process change.
- [ ] **USEQ-0FA18749** — Test critical runbooks and instructions rather than relying on editorial review alone.
- [ ] **USEQ-74890670** — Remove or archive obsolete knowledge promptly enough to prevent unsafe use.
- [ ] **USEQ-A2D587D9** — Preserve historical versions needed for support, audit, legal, research, and incident reconstruction.
- [ ] **USEQ-435CCB76** — Record why material guidance was changed or retired.
- [ ] **USEQ-D0AA1A12** — Ensure redirects, references, training, and automation are updated when authoritative knowledge moves.
- [ ] **USEQ-2C8DAAB6** — Prevent generated or imported knowledge from bypassing review and ownership.

### Competence models and role readiness

- [ ] **USEQ-7BE1905F** — Define required knowledge, skills, judgment, authorization, and experience for each material role.
- [ ] **USEQ-BD044858** — Distinguish awareness, working competence, specialist competence, and authority to approve or act.
- [ ] **USEQ-A1E84D31** — Assess current competence against role requirements.
- [ ] **USEQ-C5984409** — Create development plans for material gaps.
- [ ] **USEQ-E80FC077** — Verify competence through demonstrated performance, not attendance or self-report alone.
- [ ] **USEQ-3097646C** — Use realistic tasks, simulations, reviews, supervised practice, or certification where appropriate.
- [ ] **USEQ-2E8AFDED** — Reassess competence after long absence, role change, major system change, or serious error.
- [ ] **USEQ-C0F7014C** — Restrict high-impact actions to people with current competence and authority.
- [ ] **USEQ-FEA30171** — Avoid creating artificial credential barriers unrelated to actual role capability.
- [ ] **USEQ-952AFCF1** — Recognize accessibility accommodations and different learning pathways without reducing outcome requirements.
- [ ] **USEQ-4BC23A67** — Track competence expiry where skills degrade or external credentials lapse.
- [ ] **USEQ-27B4D8D1** — Include suppliers and contractors in competence requirements.

### Onboarding, role changes, and offboarding

- [ ] **USEQ-8B67E945** — Provide role-specific onboarding that covers purpose, users, architecture, data, security, privacy, quality, operations, and support.
- [ ] **USEQ-6DD88B91** — Give new personnel a clear map of systems, owners, authoritative knowledge, decision rights, and escalation paths.
- [ ] **USEQ-FEAB4A25** — Provide a safe environment and bounded tasks before granting high-impact production authority.
- [ ] **USEQ-C0FE3EC3** — Assess onboarding effectiveness by time to safe independent contribution and error patterns.
- [ ] **USEQ-C8406B27** — Update access, training, responsibilities, and knowledge when roles change.
- [ ] **USEQ-DE657E35** — Conduct structured knowledge transfer before planned departures.
- [ ] **USEQ-7508DA74** — Capture unresolved work, decisions, risks, contacts, credentials custody, and supplier context.
- [ ] **USEQ-4EDCE1DB** — Recover organizational assets and revoke access during offboarding.
- [ ] **USEQ-4277B789** — Preserve necessary knowledge without retaining personal data or private material unnecessarily.
- [ ] **USEQ-8BCF9F5F** — Review whether a departure creates unacceptable single-person dependency or support gap.

### Mentoring, communities, and shared practice

- [ ] **USEQ-4BC7BCDC** — Establish communities of practice for domains where consistency and learning materially improve outcomes.
- [ ] **USEQ-6077092C** — Give communities clear purpose, sponsorship, decision boundaries, and time to operate.
- [ ] **USEQ-8C107E0D** — Encourage mentoring, pairing, review, rotation, and shadowing for critical skills.
- [ ] **USEQ-4D3A335F** — Avoid using informal communities as a substitute for accountable ownership and maintained standards.
- [ ] **USEQ-07FB1FF5** — Capture reusable outcomes from discussions without recording sensitive conversation unnecessarily.
- [ ] **USEQ-9328A714** — Share incident and defect learning across teams that face similar risks.
- [ ] **USEQ-CE78EF77** — Encourage challenge, questions, and dissent without status-based suppression.
- [ ] **USEQ-47E26039** — Recognize teaching, documentation, review, and community maintenance as valuable work.
- [ ] **USEQ-F062AD31** — Evaluate whether communities improve consistency, reuse, onboarding, and problem resolution.
- [ ] **USEQ-B30208C6** — Retire or redesign communities that no longer produce useful outcomes.

### Learning from work, incidents, and evidence

- [ ] **USEQ-25BCCC4A** — Create feedback loops from production, support, testing, audits, user research, and business outcomes into engineering knowledge.
- [ ] **USEQ-0EAA2DD8** — Analyze near misses, recovered errors, rejected changes, and successful prevention as well as incidents.
- [ ] **USEQ-2B75385B** — Distinguish local mistake from systemic incentives, tooling, design, workload, or knowledge failure.
- [ ] **USEQ-9B52148F** — Convert lessons into changed requirements, code, tests, controls, procedures, training, and ownership.
- [ ] **USEQ-E9DAF617** — Verify that corrective actions changed outcomes rather than merely creating documents.
- [ ] **USEQ-F86DE0A9** — Share lessons proportionately without exposing personal, customer, privileged, or security-sensitive information.
- [ ] **USEQ-F50CD9BA** — Track recurring lessons that have not produced effective change.
- [ ] **USEQ-D4161A24** — Preserve incident chronology, decisions, evidence, and uncertainty for future learning.
- [ ] **USEQ-0BCCB40A** — Revisit old decisions when new evidence invalidates assumptions.
- [ ] **USEQ-918085CF** — Measure whether known problems recur after declared remediation.

### Succession, resilience, and bus-factor reduction

- [ ] **USEQ-2C89D675** — Identify roles and knowledge domains with insufficient qualified backup.
- [ ] **USEQ-0057EA20** — Maintain at least one practical continuity path for every critical capability.
- [ ] **USEQ-89EC9701** — Cross-train personnel through real supervised work, not documentation alone.
- [ ] **USEQ-D8FB8E99** — Rotate operational and review responsibility enough to prevent hidden dependence.
- [ ] **USEQ-A988FCED** — Ensure emergency access and recovery do not depend on one person’s device, memory, or availability.
- [ ] **USEQ-5BA62E06** — Maintain supplier and external-expert alternatives where internal capability is intentionally limited.
- [ ] **USEQ-8ED5CE14** — Preserve build, deployment, recovery, and signing knowledge independent of one tool administrator.
- [ ] **USEQ-623582AC** — Test continuity during leave, simulated unavailability, and organizational change.
- [ ] **USEQ-4E17EC00** — Avoid concentration of authority that prevents independent review or creates coercion risk.
- [ ] **USEQ-56D4A68C** — Treat an irreplaceable single knowledge holder for a critical production capability as a material risk.

### Knowledge security and trustworthy use

- [ ] **USEQ-40F8A533** — Apply access control and handling rules according to knowledge sensitivity.
- [ ] **USEQ-CF05C656** — Prevent confidential or restricted knowledge from entering public tools, prompts, repositories, recordings, or support channels.
- [ ] **USEQ-26BF83D3** — Protect credentials, secrets, personal data, exploit details, and customer information in documentation.
- [ ] **USEQ-268DD942** — Provide safe examples and synthetic data rather than real restricted values.
- [ ] **USEQ-6213693F** — Monitor unauthorized bulk export or unusual access to high-value knowledge repositories.
- [ ] **USEQ-1C374A1F** — Verify external and AI-generated knowledge before relying on it for material decisions.
- [ ] **USEQ-4544D84B** — Preserve citations and provenance for imported guidance.
- [ ] **USEQ-D5DAA624** — Detect stale, fabricated, manipulated, or conflicting content.
- [ ] **USEQ-C12331DF** — Provide a correction and withdrawal process for dangerous guidance.
- [ ] **USEQ-A275C64A** — Back up and test recovery of critical knowledge repositories.
- [ ] **USEQ-1BA3A22D** — Maintain alternate access when the normal knowledge platform is unavailable.

### Knowledge-management assurance

- [ ] **USEQ-0DFA37A5** — Define indicators for findability, freshness, coverage, use, comprehension, competence, concentration risk, and recovery.
- [ ] **USEQ-B2E6E3E9** — Sample whether intended users can locate and correctly execute critical guidance.
- [ ] **USEQ-A28F8DFF** — Audit authoritative-source ownership, review status, access, and retirement.
- [ ] **USEQ-D7B09D2D** — Exercise knowledge continuity during incidents, staff absence, supplier loss, and platform outage.
- [ ] **USEQ-001E1D5A** — Include knowledge and competence gaps in release, transition, continuity, and risk decisions.
- [ ] **USEQ-AD8EEB84** — Treat unavailable critical procedures, unknown ownership, untested competence, or single-person recovery dependence as no-go conditions when impact is material.

## Standards and source references

- [ISO/IEC/IEEE 15289:2019 — Life-cycle information items](https://www.iso.org/standard/74909.html)
- [ISO/IEC/IEEE 26514:2022 — Information for users](https://www.iso.org/standard/77451.html)
- [ISO/IEC/IEEE 12207:2026 — Software life cycle processes](https://www.iso.org/standard/90219.html)
- [ISO/IEC/IEEE 29148:2018 — Requirements engineering](https://www.iso.org/standard/72089.html)
- [ISO/IEC/IEEE 16326:2019 — Project management](https://www.iso.org/standard/74397.html)
- [ISO/IEC/IEEE 42010:2022 — Architecture description](https://www.iso.org/standard/74393.html)
- [RFC 9110 — HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110)
- [RFC 9457 — Problem Details for HTTP APIs](https://www.rfc-editor.org/rfc/rfc9457)
- [ISO 9241-210:2019 — Human-centred design](https://www.iso.org/standard/77520.html)
- [W3C Web Content Accessibility Guidelines 2.2](https://www.w3.org/TR/WCAG22/)
- [ISO/IEC 20000-1:2018 — Service management system requirements](https://www.iso.org/standard/70636.html)
- [NIST SP 800-61 Rev. 3 — Incident Response](https://csrc.nist.gov/pubs/sp/800/61/r3/final)

---

[Previous phase](12-operations-sre-and-support.md) · [Next: Phase 14: Trust, safety, and ecosystems](14-trust-safety-and-ecosystems.md)
