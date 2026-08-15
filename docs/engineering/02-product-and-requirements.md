# Product and requirements

_Phase 2 of 16 in the [complete engineering review](00-overview.md)._

Strategy, discovery, requirements, prioritization, outcomes, experimentation, and product lifecycle.

Assess every control as **Pass**, **Fail**, **Blocked**, or **Not Applicable**. A Pass needs current evidence for the exact product, revision, artifact, configuration, environment, and data state under review.

## Product and Business Readiness

_Consolidated from `quality standards/02-product/01-product-and-business-readiness.md`; 14 non-duplicative controls._

### Product definition

- [ ] **USEQ-13EC8AE3** — The problem, intended users, release scope, non-goals, and deferred features are documented.
- [ ] **USEQ-429BBD1E** — Critical user journeys are identified and ranked by user and business impact.
- [ ] **USEQ-9DBA28BF** — Critical journeys cover onboarding, normal use, error recovery, account recovery, cancellation, and data/account deletion where applicable.
- [ ] **USEQ-AFDBA1B7** — Behavior is specified for unavailable dependencies, partial failures, empty data, malformed data, incomplete data, stale data, and conflicting data.
- [ ] **USEQ-6925AA24** — Pricing, entitlements, quotas, limits, billing, cancellation, refund, reversal, dispute, and undo behavior are unambiguous.
- [ ] **USEQ-CE594458** — Customer-support obligations, service hours, maintenance expectations, SLOs, SLAs, and contractual commitments are documented and aligned.
- [ ] **USEQ-F951491B** — Product and marketing claims are supported by implemented behavior and evidence.

### Launch objectives

- [ ] **USEQ-42AA91B8** — Define measurable launch success, failure, stop, and rollback criteria.
- [ ] **USEQ-6E15410E** — Define adoption, conversion, quality, reliability, support, safety, and abuse indicators as applicable.
- [ ] **USEQ-E84CA4EA** — Define who may pause, continue, expand, or roll back the launch.
- [ ] **USEQ-B636D78E** — Establish whether launch is internal, private, regional, percentage-based, invitation-only, or public.
- [ ] **USEQ-CD424A12** — Ensure product, marketing, sales, support, and operations use the same release date and scope.
- [ ] **USEQ-B319D366** — Ensure help content, screenshots, tours, pricing pages, documentation, and legal text match released behavior.
- [ ] **USEQ-2A991D71** — Notify customers of breaking changes, downtime, migrations, or changed data use where required.

## Product Strategy, Vision, and Outcomes

_Consolidated from `quality standards/02-product/02-product-strategy-vision-and-outcomes.md`; 14 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-3F2DC642** — State the user or societal problem independently of a preferred solution.
- [ ] **USEQ-232B216A** — Define target users, buyers, operators, administrators, support teams, and affected non-users.
- [ ] **USEQ-E1BD3483** — Define the intended outcome and observable evidence that the problem is reduced.
- [ ] **USEQ-66F98F70** — Define strategic non-goals and choices to prevent unbounded scope.
- [ ] **USEQ-4E6FBFF6** — Validate that the product aligns with organizational mission, capabilities, risk appetite, and operating model.
- [ ] **USEQ-F06D1CFF** — Describe the value exchange, including what users provide in money, data, attention, effort, or lock-in.
- [ ] **USEQ-F7E05701** — Identify alternatives, substitutes, and the cost of not building the product.
- [ ] **USEQ-BA54C71F** — Define quality attributes as part of the value proposition rather than downstream technical concerns.
- [ ] **USEQ-FB65AC3C** — Define ethical, legal, accessibility, privacy, security, reliability, and sustainability boundaries.
- [ ] **USEQ-1BF8A0A2** — Specify assumptions about demand, willingness, behavior, channels, support, and unit economics.
- [ ] **USEQ-02161588** — Test the riskiest assumptions before scaling irreversible investment.
- [ ] **USEQ-92DBE30F** — Ensure strategy can be understood and used to make prioritization trade-offs.
- [ ] **USEQ-0E97500F** — Review strategy when market evidence, regulation, technology, incidents, or user needs change.
- [ ] **USEQ-48034FED** — Do not confuse a roadmap, feature list, architecture, or revenue target with a product strategy.

## Product Discovery and Problem Validation

_Consolidated from `quality standards/02-product/03-product-discovery-and-problem-validation.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-43785162** — Use multiple evidence sources rather than relying on stakeholder opinion or a single research method.
- [ ] **USEQ-B77E43E7** — Recruit participants representative of relevant abilities, contexts, languages, devices, roles, and risk profiles.
- [ ] **USEQ-0ED90387** — Separate observed behavior from stated preference and researcher interpretation.
- [ ] **USEQ-1EC01EEF** — Document research questions, sampling, methods, consent, limitations, and analysis approach.
- [ ] **USEQ-D64D10C7** — Protect participant privacy, dignity, safety, compensation fairness, and withdrawal rights.
- [ ] **USEQ-372C9D37** — Validate problem frequency, severity, current workarounds, switching costs, and consequences.
- [ ] **USEQ-6D65C0B7** — Investigate why the problem exists before optimizing the visible symptom.
- [ ] **USEQ-638BEA5E** — Include edge users and people most likely to be harmed or excluded.
- [ ] **USEQ-FB9037E1** — Use prototypes and experiments appropriate to the uncertainty being reduced.
- [ ] **USEQ-446AFA18** — Avoid presenting speculative or unsafe prototypes as operational promises.
- [ ] **USEQ-AA7C1D23** — Triangulate qualitative findings with behavioral, operational, market, and support evidence.
- [ ] **USEQ-2E0A4FD3** — Record disconfirming evidence and avoid confirmation-biased research plans.
- [ ] **USEQ-EE58AEEE** — Validate desirability, usability, feasibility, viability, accessibility, security, privacy, and operability together.
- [ ] **USEQ-630E7D69** — Stop discovery activities when additional evidence no longer changes decisions materially.
- [ ] **USEQ-82E8E2F9** — Translate validated findings into traceable requirements, design constraints, and measurable outcomes.

## Requirements Engineering

_Consolidated from `quality standards/02-product/04-requirements-engineering.md`; 17 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-735B6BCE** — Identify all requirement sources and their authority, including users, operators, law, contracts, risk, standards, and dependencies.
- [ ] **USEQ-7B2F50E9** — Distinguish needs, goals, constraints, functional requirements, quality requirements, interface requirements, and design decisions.
- [ ] **USEQ-FD9E56C3** — Write requirements that are necessary, singular, unambiguous, feasible, verifiable, consistent, and traceable.
- [ ] **USEQ-639F791F** — State the subject, required behavior or property, conditions, thresholds, units, and tolerances.
- [ ] **USEQ-60E89506** — Avoid subjective terms such as fast, intuitive, secure, scalable, or robust without measurable meaning.
- [ ] **USEQ-00B1BC4F** — Define normal, alternate, exception, degraded, recovery, and prohibited behavior.
- [ ] **USEQ-3640A26A** — Specify data, state, timing, ordering, concurrency, retention, and authorization semantics.
- [ ] **USEQ-E6110D89** — Specify accessibility, privacy, security, reliability, support, observability, and lifecycle requirements explicitly.
- [ ] **USEQ-9C6943D9** — Resolve conflicts and record the rationale and authority for resolution.
- [ ] **USEQ-2AFFB86B** — Validate requirements with representative stakeholders and intended users.
- [ ] **USEQ-15B9169F** — Link every approved requirement to acceptance criteria and planned verification.
- [ ] **USEQ-F8007542** — Control requirement changes and analyze downstream impact before approval.
- [ ] **USEQ-52E6B01C** — Preserve history and effective dates; do not silently rewrite accepted requirements.
- [ ] **USEQ-CAAD1092** — Identify derived requirements and the assumptions from which they follow.
- [ ] **USEQ-E6A9E01F** — Remove requirements that no longer support an approved need rather than preserving them by inertia.

### Category no-go conditions

- [ ] **USEQ-72F577BC** — Critical quality attributes are implied but not specified.
- [ ] **USEQ-E26272FA** — Requirements cannot be verified objectively or traced to an authorized need.

## Prioritization, Roadmaps, and Scope Control

_Consolidated from `quality standards/02-product/05-prioritization-roadmaps-and-scope-control.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-90F20378** — Prioritize outcomes and risks before individual feature requests.
- [ ] **USEQ-92888739** — Use explicit criteria such as user impact, strategic fit, risk reduction, urgency, learning value, cost, reversibility, and dependency value.
- [ ] **USEQ-D0EAD736** — Include reliability, security, privacy, accessibility, maintenance, support, and debt work in the same prioritization system.
- [ ] **USEQ-1CA42AE8** — Make stakeholder influence and mandatory obligations visible rather than hiding them inside scores.
- [ ] **USEQ-137336BA** — Represent uncertainty and sensitivity instead of treating prioritization scores as objective truth.
- [ ] **USEQ-79915E9D** — Identify prerequisites, sequencing constraints, migration needs, and operational readiness work.
- [ ] **USEQ-60BEFCD2** — Define the smallest coherent scope that can safely produce evidence or value.
- [ ] **USEQ-80609BBD** — Protect critical quality attributes from being traded away implicitly to meet dates.
- [ ] **USEQ-DE6C27A5** — Use roadmaps to communicate intent, hypotheses, sequence, and confidence—not false certainty.
- [ ] **USEQ-5E98FA21** — Distinguish committed obligations from targets, options, and exploratory work.
- [ ] **USEQ-01BD0F61** — Define entry and exit criteria for roadmap stages.
- [ ] **USEQ-01BCFAAE** — Reassess priorities when assumptions, incidents, capacity, regulation, or user evidence changes.
- [ ] **USEQ-3B778F2C** — Limit work in progress and stop starting work that cannot be completed to the required quality.
- [ ] **USEQ-5B61643E** — Make descoping decisions explicit and analyze their impact on support, operations, safety, and future change.
- [ ] **USEQ-10689ED1** — Document who made material priority decisions and why.

## Product Metrics, Experimentation, and Learning

_Consolidated from `quality standards/02-product/06-product-metrics-experimentation-and-learning.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-0A03EEA5** — Define a small set of outcome metrics linked to user value and product purpose.
- [ ] **USEQ-590E9EFF** — Pair outcome metrics with guardrails for quality, safety, privacy, accessibility, reliability, abuse, and cost.
- [ ] **USEQ-D470C55C** — Define event semantics and data-quality requirements before instrumenting them.
- [ ] **USEQ-8180A0E5** — Verify that telemetry reflects actual user behavior and not duplicate, missing, synthetic, or bot events.
- [ ] **USEQ-9C795830** — Pre-register experiment hypothesis, population, unit, allocation, duration, primary measure, guardrails, and stopping rules.
- [ ] **USEQ-F5E0C7E3** — Use randomization or a credible causal method when making causal claims.
- [ ] **USEQ-FA653F6A** — Prevent sample-ratio mismatch, cross-group contamination, novelty effects, peeking, and repeated-testing bias.
- [ ] **USEQ-C9444031** — Segment results to detect harmed cohorts without data dredging or privacy-invasive profiling.
- [ ] **USEQ-B58791C5** — Do not ship a statistically significant result that is practically harmful or strategically irrelevant.
- [ ] **USEQ-5391DDEE** — Treat inconclusive and negative results as learning rather than pressure to manipulate analysis.
- [ ] **USEQ-11E13BB4** — Preserve experiment configuration, code, data lineage, analysis, and decision records.
- [ ] **USEQ-86D51237** — Ensure users are not exposed to materially unsafe or deceptive variants.
- [ ] **USEQ-883443BC** — Define feature-flag cleanup, long-term holdout, and post-launch validation plans.
- [ ] **USEQ-37C2F56F** — Reconcile product analytics with operational, financial, support, and qualitative evidence.
- [ ] **USEQ-8EA9F41B** — Review whether optimizing a metric is degrading the underlying user outcome.

## Product Lifecycle, Deprecation, and Retirement

_Consolidated from `quality standards/02-product/07-product-lifecycle-deprecation-and-retirement.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-0B5ECDD2** — Define lifecycle states and entry, support, maintenance, deprecation, and retirement criteria.
- [ ] **USEQ-E363C389** — Publish support windows and compatibility commitments that the organization can sustain.
- [ ] **USEQ-E028DE12** — Track adoption, usage, critical customers, integrations, dependencies, and data before deprecation.
- [ ] **USEQ-3ACFDDBF** — Provide proportionate notice, migration guidance, export, tooling, and support.
- [ ] **USEQ-43188F99** — Avoid breaking security, accessibility, legal, or data-access obligations during reduced support.
- [ ] **USEQ-33C4E64B** — Define how defects and vulnerabilities are handled during maintenance and extended-support periods.
- [ ] **USEQ-B5593A92** — Prevent new dependencies on components already marked for retirement.
- [ ] **USEQ-B7C71F97** — Measure migration progress and contact customers at risk of disruption.
- [ ] **USEQ-65B667B8** — Keep old and new versions interoperable for the documented transition where required.
- [ ] **USEQ-BD2D6F7B** — Define final data export, retention, deletion, legal hold, and backup behavior.
- [ ] **USEQ-CD275A7D** — Revoke credentials, domains, certificates, endpoints, jobs, queues, infrastructure, and supplier access during retirement.
- [ ] **USEQ-79D1353B** — Preserve required records and evidence after operational systems are removed.
- [ ] **USEQ-0151D776** — Verify that retired assets are no longer reachable, billed, monitored as active, or processing data.
- [ ] **USEQ-8429659A** — Provide rollback or continuity options for failed migrations where practical.
- [ ] **USEQ-94B6F552** — Conduct a retirement review and capture lessons for future lifecycle planning.

## Standards and source references

- [ISO/IEC 25010:2023 — Product quality model](https://www.iso.org/standard/78176.html)
- [ISO/IEC 25019:2023 — Quality-in-use model](https://www.iso.org/standard/78177.html)
- [ISO/IEC/IEEE 29148:2018 — Requirements engineering](https://www.iso.org/standard/72089.html)
- [ISO/IEC 38500:2024 — Governance of IT](https://www.iso.org/standard/81684.html)
- [ISO 9241-210:2019 — Human-centred design](https://www.iso.org/standard/77520.html)
- [ISO/IEC/IEEE 12207:2026 — Software life cycle processes](https://www.iso.org/standard/90219.html)
- [ISO/IEC/IEEE 16326:2019 — Project management](https://www.iso.org/standard/74397.html)
- [ISO 31000:2018 — Risk management guidelines](https://www.iso.org/standard/65694.html)
- [ISO/IEC/IEEE 15939:2017 — Measurement process](https://www.iso.org/standard/71197.html)
- [ISO/IEC 42001:2023 — AI management systems](https://www.iso.org/standard/81230.html)
- [ISO/IEC 20000-1:2018 — Service management system requirements](https://www.iso.org/standard/70636.html)
- [ISO 22301:2019 — Business continuity management systems](https://www.iso.org/standard/75106.html)

---

[Previous phase](01-governance-and-foundations.md) · [Next: Phase 3: User experience, web, and content](03-user-experience-web-and-content.md)
