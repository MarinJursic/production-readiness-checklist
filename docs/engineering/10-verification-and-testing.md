# Verification and testing

_Phase 10 of 16 in the [complete engineering review](00-overview.md)._

Test strategy, unit through acceptance testing, nonfunctional assurance, test data, static analysis, defects, and root-cause learning.

Assess every control as **Pass**, **Fail**, **Blocked**, or **Not Applicable**. A Pass needs current evidence for the exact product, revision, artifact, configuration, environment, and data state under review.

## Test Strategy and Evidence

_Consolidated from `quality standards/11-testing-quality-assurance/01-test-strategy-and-evidence.md`; 13 non-duplicative controls._

### Universal controls

- [ ] **USEQ-62B6FFE6** — A documented test strategy exists for the release's risk level.
- [ ] **USEQ-B2D342EF** — Unit, component, integration, contract, system, end-to-end, smoke, and regression testing are included where appropriate.
- [ ] **USEQ-20EF5A60** — Testing prioritizes critical and high-risk behavior rather than relying on aggregate code coverage.
- [ ] **USEQ-10FA9C47** — Test data represents realistic size, shape, encoding, and edge conditions.
- [ ] **USEQ-7B4C2AB9** — Personal and confidential test data is synthetic, masked, or properly authorized.
- [ ] **USEQ-80A6223A** — Release-gating tests are sufficiently deterministic.
- [ ] **USEQ-876F7B66** — Flaky tests are fixed or formally quarantined with a reviewed risk record.
- [ ] **USEQ-B818C295** — Migration, upgrade, mixed-version, downgrade where supported, rollback, and roll-forward scenarios are tested.
- [ ] **USEQ-8D7753F8** — Performance testing covers load, stress, spike, endurance, and capacity limits.
- [ ] **USEQ-5C6A77E1** — Security tests cover source, dependencies, deployed behavior, infrastructure, configuration, and business logic.
- [ ] **USEQ-55DBD644** — Critical calculations use independent verification or reconciliation.
- [ ] **USEQ-E01F5F5B** — Tests verify the exact artifact intended for production.
- [ ] **USEQ-BA5EE5B3** — Independent penetration testing is completed where risk, contract, or regulation warrants it.

## Unit and Component Testing

_Consolidated from `quality standards/11-testing-quality-assurance/02-unit-and-component-testing.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-92EEEDDD** — Define the unit or component boundary around cohesive behavior rather than arbitrary files or methods.
- [ ] **USEQ-5410EA72** — Test externally meaningful behavior and invariants rather than incidental implementation sequence.
- [ ] **USEQ-2FCE8AAC** — Cover representative equivalence classes, boundaries, invalid inputs, state transitions, and error paths.
- [ ] **USEQ-FB830B78** — Test pure logic across all material branches and rules.
- [ ] **USEQ-C964A6AD** — Test serialization, precision, time, locale, normalization, and identifier edge cases where relevant.
- [ ] **USEQ-F843BEDE** — Control time, randomness, scheduling, environment, and external effects for determinism.
- [ ] **USEQ-554CBC58** — Use real domain values and builders that create valid data by default.
- [ ] **USEQ-C56D105D** — Keep each test independent and safe to run in any order and in parallel.
- [ ] **USEQ-F0E6FA63** — Make test names describe condition, behavior, and expected outcome.
- [ ] **USEQ-BCD4E731** — Ensure assertions prove the intended behavior and fail when that behavior is broken.
- [ ] **USEQ-BEA54B0C** — Avoid excessive mocking that merely repeats implementation or permits impossible behavior.
- [ ] **USEQ-48C6DC0C** — Use contract-constrained fakes for external boundaries where appropriate.
- [ ] **USEQ-2241BCAB** — Test resource cleanup, cancellation, retry, and idempotency locally where owned.
- [ ] **USEQ-AE1A5D9B** — Use property-based, generative, mutation, or model-based testing for high-risk combinatorial logic.
- [ ] **USEQ-5D8DD5E5** — Keep tests fast enough for frequent feedback without sacrificing meaningful coverage.
- [ ] **USEQ-A567F4BA** — Measure critical behavior and mutation sensitivity rather than treating line coverage as the goal.
- [ ] **USEQ-913335AC** — Refactor brittle tests when harmless internal change causes widespread failure.
- [ ] **USEQ-3DD8385B** — Do not expose production internals solely to satisfy tests without a design benefit.
- [ ] **USEQ-263390BE** — Review skipped, quarantined, or flaky tests with owner and expiry.
- [ ] **USEQ-C4149EA3** — Convert escaped local defects into focused regression tests.

## Integration and Contract Testing

_Consolidated from `quality standards/11-testing-quality-assurance/03-integration-and-contract-testing.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-28F1E368** — Identify every material internal and external boundary and its contract owner.
- [ ] **USEQ-14D0AB33** — Test request, response, event, file, database, authentication, authorization, and error schemas.
- [ ] **USEQ-EF768482** — Verify serialization, encoding, precision, ordering, pagination, null, optional, and unknown-field behavior.
- [ ] **USEQ-B0707D0F** — Test provider and consumer expectations independently to support compatible evolution.
- [ ] **USEQ-D044F47B** — Run contract tests against faithful provider implementations or approved simulators.
- [ ] **USEQ-DF85751F** — Test persistence constraints, transactions, isolation, migrations, and query behavior using representative engines and settings.
- [ ] **USEQ-A20F1BCD** — Test authentication, authorization, tenant, quota, and rate behavior across boundaries.
- [ ] **USEQ-09FC670C** — Inject timeout, throttling, malformed, partial, duplicate, reordered, stale, and unavailable dependency responses.
- [ ] **USEQ-122F73E2** — Verify retry, idempotency, circuit breaking, fallback, reconciliation, and observability.
- [ ] **USEQ-1FEBE72F** — Test webhook signature, replay, ordering, retry, and duplicate delivery.
- [ ] **USEQ-597ABA9E** — Test protocol and schema version negotiation and mixed-version operation.
- [ ] **USEQ-24D28788** — Prevent shared mutable test environments from causing cross-test contamination.
- [ ] **USEQ-786C6BF4** — Use isolated data and cleanup with evidence that cleanup cannot delete another test's resources.
- [ ] **USEQ-C013247E** — Do not hide provider incompatibility behind mocks that are more permissive than reality.
- [ ] **USEQ-09078504** — Run a controlled subset against real third-party sandboxes or conformance endpoints where available.
- [ ] **USEQ-82DB4AE3** — Record provider version and configuration with results.
- [ ] **USEQ-0029BEB2** — Monitor contract failures as release blockers for affected consumers.
- [ ] **USEQ-5A2F54E4** — Coordinate breaking changes through published migration and deprecation.
- [ ] **USEQ-D2C62906** — Test data residency, privacy, and logging at integration boundaries.
- [ ] **USEQ-AD3AD837** — Convert production integration failures into reproducible contract cases.

## System, End-to-End, and Acceptance Testing

_Consolidated from `quality standards/11-testing-quality-assurance/04-system-end-to-end-and-acceptance-testing.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-804A9C95** — Prioritize critical user and operator journeys by impact rather than attempting exhaustive UI-path enumeration.
- [ ] **USEQ-59137891** — Trace each acceptance test to approved requirements and business outcomes.
- [ ] **USEQ-738ED752** — Use production-like topology, configuration, identity, data shape, integrations, and artifact builds.
- [ ] **USEQ-08D8D136** — Test onboarding, normal use, permission variation, recovery, cancellation, deletion, support, and administration as applicable.
- [ ] **USEQ-7F711C38** — Test alternate, error, degraded, offline, timeout, retry, and partial-success journeys.
- [ ] **USEQ-C50D762D** — Verify business invariants, notifications, background work, analytics, audit, and downstream effects.
- [ ] **USEQ-618DC5B8** — Use multiple roles, tenants, accounts, locales, devices, browsers, and accessibility contexts.
- [ ] **USEQ-ED737106** — Validate actual external boundaries where safe and use contract-tested substitutes otherwise.
- [ ] **USEQ-3368AA1A** — Keep test setup and teardown deterministic, isolated, and observable.
- [ ] **USEQ-ADC0AFF5** — Avoid making every detailed rule dependent on slow end-to-end tests.
- [ ] **USEQ-3F958D15** — Use stable semantic selectors and user-observable synchronization rather than implementation timing.
- [ ] **USEQ-110F3E80** — Record artifacts sufficient to diagnose failure across services.
- [ ] **USEQ-3A691B37** — Separate product acceptance from mere technical execution success.
- [ ] **USEQ-1A04F649** — Include support and operational acceptance for runbooks, monitoring, capacity, and recovery.
- [ ] **USEQ-A8C0F294** — Verify data created, updated, exported, corrected, and deleted across all relevant stores.
- [ ] **USEQ-05687BCB** — Test deployment, mixed-version, migration, rollback, and post-deployment smoke behavior.
- [ ] **USEQ-A587AD05** — Manage flaky tests as defects with owner, cause, and expiry.
- [ ] **USEQ-F7B75AC7** — Run production-safe synthetic journeys for critical paths after launch.
- [ ] **USEQ-B4A1B852** — Obtain qualified stakeholder acceptance for high-impact requirements.
- [ ] **USEQ-9665D391** — Convert escaped journey defects into the lowest reliable regression layer plus an end-to-end sentinel when valuable.

## Regression, Change-Impact, and Compatibility Testing

_Consolidated from `quality standards/11-testing-quality-assurance/05-regression-change-impact-and-compatibility-testing.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-A9DFA5B1** — Analyze change impact across requirements, modules, data, interfaces, consumers, dependencies, environments, and operations.
- [ ] **USEQ-02AF9A1D** — Select regression tests from risk and dependency evidence, not only changed files.
- [ ] **USEQ-D3436F27** — Maintain a stable critical regression suite that represents business, security, privacy, accessibility, and reliability invariants.
- [ ] **USEQ-5258BB43** — Test backward, forward, mixed-version, schema, protocol, browser, device, and data compatibility as applicable.
- [ ] **USEQ-36EFB529** — Test changed defaults, configuration, flags, permissions, quotas, and feature combinations.
- [ ] **USEQ-4A9B0A0E** — Include known defect regressions and historically fragile areas.
- [ ] **USEQ-C2E84396** — Use differential or golden-master comparison carefully where behavior is large but not fully specified.
- [ ] **USEQ-B7BBB198** — Review intended behavior changes so they do not appear as unexplained regression failures.
- [ ] **USEQ-5D5B4D27** — Run broad regression before irreversible or high-blast-radius changes.
- [ ] **USEQ-821187B8** — Use production-representative data distributions and older valid records.
- [ ] **USEQ-1DA0B1A4** — Test upgrades from every supported version and supported downgrade or roll-forward path.
- [ ] **USEQ-48F4A75B** — Verify client and server caching does not preserve incompatible state.
- [ ] **USEQ-1B98C843** — Measure regression effectiveness through escaped defects and mutation or fault sensitivity.
- [ ] **USEQ-F1DECC35** — Retire obsolete tests only after confirming the protected behavior or risk no longer exists.
- [ ] **USEQ-20F430F8** — Avoid allowing test-suite duration to grow without architecture and selection strategy.
- [ ] **USEQ-42745DA4** — Use staged test selection with mandatory high-risk gates and broader asynchronous evidence before final approval.
- [ ] **USEQ-54517178** — Verify the final artifact after packaging and deployment transformation.
- [ ] **USEQ-1EBA61FD** — Include security control and observability regression.
- [ ] **USEQ-EE8433B8** — Preserve a known-good baseline and explain material differences.
- [ ] **USEQ-4826FE5C** — Convert incidents caused by interaction effects into cross-component regression coverage.

## Test Automation Quality

_Consolidated from `quality standards/11-testing-quality-assurance/06-test-automation-quality.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-CECAF01C** — Automate checks that are repeatable, valuable, and sufficiently deterministic.
- [ ] **USEQ-FA8213FA** — Keep automation code to production-quality standards for review, versioning, security, and maintenance.
- [ ] **USEQ-43C3E4B8** — Make failures identify the violated expectation, actual result, context, and evidence.
- [ ] **USEQ-CC168172** — Keep tests independent and safe for parallel or repeated execution.
- [ ] **USEQ-F21A03D2** — Control time, randomness, external dependencies, and data where practical.
- [ ] **USEQ-74328706** — Use retries only to diagnose known infrastructure instability, never to convert flaky failures into passes.
- [ ] **USEQ-7247BFD2** — Track flakiness, quarantine duration, ownership, and business risk.
- [ ] **USEQ-9A51EF4B** — Prevent test credentials, personal data, and secrets from leaking through fixtures and reports.
- [ ] **USEQ-256B9A15** — Version test data, simulators, schemas, baselines, and expected results with compatible production behavior.
- [ ] **USEQ-57E6BCF0** — Validate test tools and custom matchers against known passing and failing examples.
- [ ] **USEQ-0CDDDBAD** — Ensure negative tests fail for the intended reason rather than an earlier setup problem.
- [ ] **USEQ-B339ABB5** — Use layered tests to balance speed, fidelity, and diagnosis.
- [ ] **USEQ-33111C22** — Avoid duplicated tests that add cost without new risk coverage.
- [ ] **USEQ-CC741332** — Make environment and prerequisite failures distinct from product failures.
- [ ] **USEQ-3F834ADF** — Preserve logs, traces, screenshots, videos, diffs, and seed values as appropriate.
- [ ] **USEQ-949E24F8** — Review disabled, skipped, expected-failure, and muted tests before release.
- [ ] **USEQ-3073031B** — Measure suite duration, reliability, fault detection, maintenance cost, and escape rate.
- [ ] **USEQ-D53E745C** — Refactor automation when product implementation details make tests brittle.
- [ ] **USEQ-34A3D75D** — Protect CI workers and test environments from untrusted code and cross-run contamination.
- [ ] **USEQ-9993C108** — Periodically test that release gates actually block a deliberately failing change.

## Test Data and Test Environments

_Consolidated from `quality standards/11-testing-quality-assurance/07-test-data-and-test-environments.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-8E3B0906** — Define the production characteristics that each test environment must reproduce and those intentionally simulated.
- [ ] **USEQ-E0D9E9AC** — Version infrastructure, configuration, schemas, dependencies, feature flags, and seed data.
- [ ] **USEQ-173F222D** — Use synthetic data by default and authorize, minimize, mask, isolate, and expire production-derived data.
- [ ] **USEQ-D350CA89** — Preserve referential, statistical, temporal, locale, encoding, and edge-case characteristics needed for valid tests.
- [ ] **USEQ-AD7590EE** — Include large, sparse, skewed, duplicate, malformed, legacy, deleted, and sensitive-shape records without exposing real identities.
- [ ] **USEQ-CFC1A943** — Assign unique test identities, tenants, namespaces, and resources to prevent cross-run interference.
- [ ] **USEQ-337C90B2** — Prevent nonproduction notifications, payments, webhooks, and integrations from affecting real users or systems.
- [ ] **USEQ-60E233B7** — Use deterministic setup and teardown and detect orphaned resources.
- [ ] **USEQ-1A1BC247** — Reset environments from reproducible automation rather than manual undocumented repair.
- [ ] **USEQ-5A7D89CD** — Detect drift between intended and actual environment configuration.
- [ ] **USEQ-DDA5AB90** — Protect test environments according to the sensitivity of code, credentials, and data they contain.
- [ ] **USEQ-D0CF0FA3** — Keep production secrets and trust relationships out of lower environments.
- [ ] **USEQ-32F22082** — Simulate dependency latency, error, throttling, and outage realistically.
- [ ] **USEQ-5C661F9B** — Measure whether environment limitations invalidate test conclusions.
- [ ] **USEQ-8CF5D4A8** — Maintain capacity sufficient for representative performance and concurrency tests.
- [ ] **USEQ-03A1613E** — Record environment version and health with every material result.
- [ ] **USEQ-F5919ABA** — Prevent shared long-lived environments from hiding order dependence and contamination.
- [ ] **USEQ-2B8B156B** — Test migration from realistic older schemas and data.
- [ ] **USEQ-56252FA6** — Review access and remove dormant test accounts, keys, data, and infrastructure.
- [ ] **USEQ-D49F972C** — Destroy or sanitize environments and data when their purpose ends.

## Nonfunctional and Quality-Attribute Testing

_Consolidated from `quality standards/11-testing-quality-assurance/08-nonfunctional-quality-attribute-testing.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-31B70C01** — Translate every material quality attribute into measurable scenarios with stimulus, environment, response, and threshold.
- [ ] **USEQ-B4739B5B** — Use representative workloads, data, users, topology, dependencies, devices, and failure conditions.
- [ ] **USEQ-B72AA355** — Test latency distributions, throughput, saturation, resource use, and efficiency together.
- [ ] **USEQ-4CCCEA59** — Test availability and correctness under dependency, node, zone, network, storage, and control-plane failure.
- [ ] **USEQ-BC77A6BC** — Test recovery time, recovery point, data integrity, failover, failback, and restored operation.
- [ ] **USEQ-A3A36C12** — Test security and privacy controls through adversarial and misuse scenarios.
- [ ] **USEQ-5BA95FC4** — Test accessibility through automation, keyboard, assistive technology, zoom, reflow, and human evaluation.
- [ ] **USEQ-0D1CBEE3** — Test usability with representative tasks and users.
- [ ] **USEQ-BF03BCAF** — Test compatibility across supported clients, versions, formats, locales, and environments.
- [ ] **USEQ-0CC6633E** — Test maintainability through change exercises, diagnostics, upgrade, and operator tasks where impact warrants it.
- [ ] **USEQ-B5C10C88** — Test scalability and overload beyond expected peak to identify safe failure limits.
- [ ] **USEQ-C97F0266** — Test long-duration operation for leaks, drift, backlog, and degradation.
- [ ] **USEQ-98EB3ECC** — Test sustainability and cost under representative workload where material.
- [ ] **USEQ-CC592135** — Verify monitoring and alerts during each quality test.
- [ ] **USEQ-E5A00B61** — Record limitations, confidence, and conditions under which results do not apply.
- [ ] **USEQ-DE808721** — Use independent review for high-impact thresholds and results.
- [ ] **USEQ-C2C6C7E0** — Do not extrapolate beyond measured scale without an explicit model and uncertainty.
- [ ] **USEQ-4861ED7B** — Repeat tests after material architecture, dependency, data, configuration, or infrastructure change.
- [ ] **USEQ-E12B5B11** — Compare results to baselines and approved budgets.
- [ ] **USEQ-2C0ECEB2** — Block release when a critical quality objective lacks valid evidence.

## Exploratory, Risk-Based, and Usability Testing

_Consolidated from `quality standards/11-testing-quality-assurance/09-exploratory-risk-based-and-usability-testing.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-7879005A** — Prioritize charters by user harm, business impact, novelty, complexity, change, defect history, and uncertainty.
- [ ] **USEQ-519085DD** — Define a mission and scope while preserving freedom to follow evidence.
- [ ] **USEQ-399DA6B9** — Use skilled testers with domain, user, technical, accessibility, security, and operational context as appropriate.
- [ ] **USEQ-AF91B7B7** — Vary data, sequence, timing, concurrency, roles, devices, networks, locales, and interruption.
- [ ] **USEQ-BBF7C191** — Investigate inconsistencies, ambiguous feedback, unexpected state, and recovery difficulty.
- [ ] **USEQ-341F2D0F** — Explore complete workflows and cross-channel transitions rather than isolated screens.
- [ ] **USEQ-D0D46EA6** — Use heuristics, models, tours, personas, abuse stories, and past incidents to generate tests.
- [ ] **USEQ-FE154839** — Capture concise notes, evidence, coverage, questions, risks, and follow-up.
- [ ] **USEQ-014C3738** — Distinguish observations, hypotheses, defects, design concerns, and user research findings.
- [ ] **USEQ-35B2FB35** — Reproduce and reduce defects to the smallest useful case.
- [ ] **USEQ-04D07FE1** — Pair with developers, designers, support, operators, or users for high-context investigation.
- [ ] **USEQ-F9EE8DA8** — Include users with disabilities and representative assistive technologies.
- [ ] **USEQ-7BA8AA27** — Avoid treating scripted checks as a replacement for exploration.
- [ ] **USEQ-BDF47008** — Time-box sessions and review whether additional exploration changes decisions.
- [ ] **USEQ-DE9F6640** — Feed discoveries into requirements, automated regression, design, support, and risk records.
- [ ] **USEQ-03907FCF** — Track untested areas and session limitations.
- [ ] **USEQ-EE9B03CE** — Use production-like or production-safe environments without exposing real users or data.
- [ ] **USEQ-F427D041** — Review quality across competing attributes rather than optimizing one path.
- [ ] **USEQ-D103850E** — Measure escaped defect themes and charter effectiveness qualitatively.
- [ ] **USEQ-A346C21D** — Preserve psychological safety so testers can challenge accepted behavior.

## Static Analysis, Formal Methods, and Verification

_Consolidated from `quality standards/11-testing-quality-assurance/10-static-analysis-formal-methods-and-verification.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-187751E3** — Select techniques according to consequence, defect class, language semantics, and system complexity.
- [ ] **USEQ-180F046B** — Enable compiler, type, lint, security, dependency, secret, and configuration checks appropriate to the codebase.
- [ ] **USEQ-1AE136B6** — Treat warnings according to a documented baseline and prevent uncontrolled growth.
- [ ] **USEQ-2BEC6E57** — Use sound or high-confidence analyses for critical properties where available.
- [ ] **USEQ-2C055DE5** — Encode domain invariants in types, schemas, contracts, assertions, and model constraints.
- [ ] **USEQ-D20DCA2A** — Use model checking or state-machine verification for high-risk concurrency, protocol, and workflow behavior where practical.
- [ ] **USEQ-1E45D431** — Use property-based and symbolic methods for parsers, numerical logic, authorization, and combinatorial inputs.
- [ ] **USEQ-5A042057** — Review false-positive suppressions with scope, rationale, owner, and expiry.
- [ ] **USEQ-E6197F5B** — Do not suppress findings globally when a narrow exception is possible.
- [ ] **USEQ-857E50CE** — Validate custom rules against known defects and safe examples.
- [ ] **USEQ-566068B3** — Protect analyzers and plugins as supply-chain dependencies.
- [ ] **USEQ-0C558F2C** — Run analysis on generated code, infrastructure, migrations, and configuration where relevant.
- [ ] **USEQ-C5555A64** — Use reproducible tool versions and rule sets.
- [ ] **USEQ-3A2F1122** — Track introduced findings separately from inherited debt.
- [ ] **USEQ-CB6715F0** — Require higher assurance for cryptography, safety, identity, authorization, and irreversible financial logic.
- [ ] **USEQ-E903DE38** — Use formal claims only within explicitly stated assumptions and model boundaries.
- [ ] **USEQ-B0A6D650** — Complement static proof with runtime, integration, and operational evidence.
- [ ] **USEQ-F267C2FC** — Review whether code transformations and optimizations preserve verified properties.
- [ ] **USEQ-468A609C** — Archive models, tool versions, assumptions, results, and review decisions.
- [ ] **USEQ-27C19FEA** — Convert recurring escaped defects into machine-enforced rules where reliable.

## Defect Management and Root-Cause Improvement

_Consolidated from `quality standards/11-testing-quality-assurance/11-defect-management-and-root-cause.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-57E0618C** — Provide clear channels for users, support, monitoring, security, testing, and suppliers to report defects.
- [ ] **USEQ-4F3E2904** — Capture affected user, journey, environment, version, data, time, symptoms, impact, and evidence.
- [ ] **USEQ-1981A1CF** — Classify severity by harm, reach, integrity, security, privacy, financial, legal, and recoverability impact.
- [ ] **USEQ-C3398A1B** — Separate severity from scheduling priority and make both explicit.
- [ ] **USEQ-5C76863C** — Triage duplicates without losing affected-user evidence or recurrence count.
- [ ] **USEQ-B9CD2703** — Contain high-impact defects before waiting for permanent correction.
- [ ] **USEQ-AFC1B6A5** — Identify the earliest point at which the defect could have been prevented or detected.
- [ ] **USEQ-12B340B5** — Analyze contributing technical, process, organizational, incentive, interface, and knowledge factors.
- [ ] **USEQ-E51510AC** — Avoid stopping at individual error or a single proximate cause.
- [ ] **USEQ-D480EB7C** — Correct corrupted or inconsistent data and reconcile downstream effects.
- [ ] **USEQ-AD42951B** — Verify fixes against the original failure, adjacent cases, alternate paths, and regression risk.
- [ ] **USEQ-520E6677** — Monitor after deployment for recurrence and unintended consequences.
- [ ] **USEQ-8C0B0798** — Communicate known impact, workarounds, and resolution accurately to affected stakeholders.
- [ ] **USEQ-23336C87** — Use time-bounded risk acceptance only with owner and compensating controls.
- [ ] **USEQ-5827E80F** — Track aging, recurrence, escape phase, detection source, and corrective-action effectiveness.
- [ ] **USEQ-23DEB246** — Prioritize systemic fixes that eliminate classes of defects.
- [ ] **USEQ-1E7A9936** — Update requirements, design, tests, monitoring, documentation, training, and review guidance.
- [ ] **USEQ-2B757A11** — Close defects only when completion criteria and evidence are satisfied.
- [ ] **USEQ-AE77FE17** — Preserve incident and defect history for trend analysis.
- [ ] **USEQ-92037A06** — Review recurring low-severity defects whose aggregate burden is material.

## Standards and source references

- [ISO/IEC/IEEE 29119-1:2022 — Software testing concepts](https://www.iso.org/standard/81291.html)
- [ISO/IEC/IEEE 29119-2:2021 — Test processes](https://www.iso.org/standard/79428.html)
- [ISO/IEC/IEEE 29119-3:2021 — Test documentation](https://www.iso.org/standard/79429.html)
- [ISO/IEC/IEEE 29119-4:2021 — Test techniques](https://www.iso.org/standard/79430.html)
- [RFC 9110 — HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110)
- [RFC 9457 — Problem Details for HTTP APIs](https://www.rfc-editor.org/rfc/rfc9457)
- [ISO/IEC/IEEE 29148:2018 — Requirements engineering](https://www.iso.org/standard/72089.html)
- [ISO/IEC/IEEE 12207:2026 — Software life cycle processes](https://www.iso.org/standard/90219.html)
- [ISO/IEC 25010:2023 — Product quality model](https://www.iso.org/standard/78176.html)
- [ISO/IEC 27701:2025 — Privacy information management systems](https://www.iso.org/standard/85819.html)
- [ISO/IEC 25019:2023 — Quality-in-use model](https://www.iso.org/standard/78177.html)
- [ISO 9241-210:2019 — Human-centred design](https://www.iso.org/standard/77520.html)
- [ISO/IEC 5055:2021 — Automated source code quality measures](https://www.iso.org/standard/80623.html)
- [NIST SP 800-218 v1.1 — Secure Software Development Framework](https://csrc.nist.gov/pubs/sp/800/218/final)
- [ISO 9001:2015 — Quality management systems](https://www.iso.org/standard/62085.html)
- [ISO/IEC/IEEE 16085:2021 — Life-cycle risk management](https://www.iso.org/standard/74385.html)

---

[Previous phase](09-privacy-and-data-protection.md) · [Next: Phase 11: Developer experience, platform, and delivery](11-developer-experience-platform-and-delivery.md)
