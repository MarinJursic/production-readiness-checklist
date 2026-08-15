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
