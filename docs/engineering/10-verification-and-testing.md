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

## Advanced Testing, Conformance, and Continuous Assurance

_Consolidated from `final consolidated corpus/06-testing-verification-conformance-assurance.md#Advanced Testing, Conformance, and Continuous Assurance`; 316 non-duplicative controls._

### Assurance strategy, scope, and independence

- [ ] **USEQ-CE5D5E0A** — Define an assurance strategy that connects product risks, quality attributes, requirements, architecture, implementation, deployment, operations, and evidence rather than treating testing as an isolated phase.
- [ ] **USEQ-27F0F898** — Identify every system boundary, user group, role, tenant, device, interface, dependency, environment, data class, operational mode, and failure mode included in the assurance scope.
- [ ] **USEQ-D148EDEA** — Identify material exclusions, assumptions, test limitations, unavailable environments, inaccessible supplier components, and residual uncertainty.
- [ ] **USEQ-0511465F** — Rank features and components by potential harm, change frequency, complexity, defect history, exposure, reversibility, detectability, and recovery difficulty.
- [ ] **USEQ-D5BDAE4B** — Allocate deeper and more independent assurance to high-risk, safety-related, security-sensitive, privacy-sensitive, financially consequential, or widely reused components.
- [ ] **USEQ-C807D350** — Separate test design, implementation, execution, result interpretation, and release approval sufficiently to reduce confirmation bias for material controls.
- [ ] **USEQ-63C4A075** — Ensure the test organization can stop a release without commercial or schedule retaliation when release criteria are not met.
- [ ] **USEQ-3C8C8565** — Define the required evidence for each lifecycle stage and quality attribute before implementation begins.
- [ ] **USEQ-2FBAD908** — Specify which assurance activities are automated, manually performed, independently witnessed, formally verified, externally assessed, or continuously monitored.
- [ ] **USEQ-600C7E40** — Use multiple complementary verification methods for critical claims rather than relying on one tool or test layer.
- [ ] **USEQ-F76F8CFD** — Review the assurance strategy after material changes in architecture, risk, regulation, suppliers, users, deployment model, incident history, or threat landscape.
- [ ] **USEQ-7D51BFD7** — Maintain competence requirements for people designing, performing, and reviewing specialized testing.
- [ ] **USEQ-62725A3B** — Treat supplier test reports, certifications, and attestations as inputs requiring scope and freshness validation rather than unquestioned proof.
- [ ] **USEQ-A0245E21** — Align assurance depth with the required confidence and consequence of failure, not simply project size.
- [ ] **USEQ-70949753** — Document the conditions under which production evidence may supplement preproduction testing and the controls that make such testing safe.

### Test basis, requirements traceability, and executable acceptance

- [ ] **USEQ-BCD4F01B** — Identify the authoritative test basis for every behavior, including requirements, user needs, contracts, standards, policies, architecture decisions, threat models, data rules, and operational objectives.
- [ ] **USEQ-5A0ADC36** — Resolve ambiguous, contradictory, untestable, unverifiable, or missing requirements before relying on test execution.
- [ ] **USEQ-22497A56** — Express critical requirements with measurable acceptance criteria, observable outcomes, tolerances, preconditions, and failure behavior.
- [ ] **USEQ-506562A9** — Trace every critical requirement and risk to one or more verification activities and every verification result back to the requirement or risk it addresses.
- [ ] **USEQ-BD7F03F5** — Trace tests to the exact product version, configuration, feature flags, schema, data state, dependency versions, and environment evaluated.
- [ ] **USEQ-2460BD77** — Identify requirements that cannot be fully tested and define alternative analysis, review, proof, monitoring, or operational controls.
- [ ] **USEQ-87D05822** — Include negative requirements and prohibited states, not only desired behavior.
- [ ] **USEQ-A6BEC6D1** — Include invariants that must remain true across operations, concurrency, retries, recovery, migration, and partial failure.
- [ ] **USEQ-ADEFDCCF** — Include quantitative tolerances for accuracy, latency, capacity, availability, reliability, numerical error, and data quality where relevant.
- [ ] **USEQ-38D4FA5F** — Include accessibility, privacy, security, interoperability, supportability, and recoverability in acceptance criteria rather than treating them as optional audits.
- [ ] **USEQ-94DBFEBA** — Keep acceptance criteria independent of one implementation when the requirement is intended to remain technology-neutral.
- [ ] **USEQ-ED0E5F15** — Version and review the test basis so that requirement changes invalidate affected test evidence.
- [ ] **USEQ-FD79CE2B** — Ensure examples and mockups do not silently replace normative behavior unless explicitly approved.
- [ ] **USEQ-C15B6052** — Record unresolved questions and assumptions with owners and deadlines.
- [ ] **USEQ-3407C631** — Require business, engineering, operations, security, privacy, accessibility, and data stakeholders to approve material acceptance criteria within their remit.

### Test-oracle design and correctness of expected results

- [ ] **USEQ-9262DCB4** — Define how each test determines correctness rather than assuming that successful execution implies a correct result.
- [ ] **USEQ-43727413** — Use independent reference calculations, authoritative data, specifications, invariants, models, or alternative implementations for high-impact outcomes.
- [ ] **USEQ-0F70FDC0** — Avoid generating expected results with the same logic or defect-prone implementation being tested.
- [ ] **USEQ-44FB8B29** — Validate test oracles themselves through review, known examples, cross-checks, and controlled faults.
- [ ] **USEQ-E9D5AC04** — Distinguish exact expected values from acceptable ranges, statistical distributions, temporal bounds, eventual outcomes, and qualitative judgments.
- [ ] **USEQ-C33A820B** — Define oracles for intermediate state, side effects, messages, logs, external calls, audit records, notifications, and resource changes, not only final user-visible output.
- [ ] **USEQ-352EFEBC** — Verify that expected error type, message class, retryability, status, and recovery behavior are correct.
- [ ] **USEQ-A9D5B9D2** — Use metamorphic relations, differential comparison, invariants, reconciliation, or anomaly detection where a complete oracle is unavailable.
- [ ] **USEQ-6A433F76** — Control nondeterminism, randomness, time, concurrency, external data, and model variability when interpreting results.
- [ ] **USEQ-12F78258** — Define statistical confidence, sample size, uncertainty, and acceptable error for probabilistic systems.
- [ ] **USEQ-3A5FB09B** — Detect false positives and false negatives in automated oracles and treat frequent manual overrides as a test-system defect.
- [ ] **USEQ-3BE1C974** — Version expected outputs and golden files and review intentional changes rather than accepting broad snapshots automatically.
- [ ] **USEQ-5AA1B557** — Prevent stale snapshots from normalizing unintended behavior.
- [ ] **USEQ-D14E0A6C** — Ensure localization, formatting, time, and ordering differences are evaluated according to semantic correctness rather than brittle textual equality.
- [ ] **USEQ-13837854** — Document the limitations and blind spots of every material oracle.

### Systematic functional test-design techniques

- [ ] **USEQ-9D9247ED** — Partition input and state spaces into meaningful equivalence classes and test valid and invalid representatives.
- [ ] **USEQ-931139A0** — Test lower, upper, just-below, just-above, empty, minimum, maximum, overflow, underflow, and unbounded boundary conditions.
- [ ] **USEQ-26215F72** — Use decision tables for behavior determined by combinations of rules, conditions, permissions, plans, or statuses.
- [ ] **USEQ-B1D3501A** — Use state-transition models for workflows, lifecycles, sessions, payments, approvals, jobs, devices, and other stateful behavior.
- [ ] **USEQ-85C74003** — Test every valid state transition, prohibited transition, repeated transition, out-of-order transition, and recovery transition that matters.
- [ ] **USEQ-51299A7B** — Use scenario and use-case testing for complete user and business journeys, including alternate, cancellation, timeout, and failure paths.
- [ ] **USEQ-7EB73115** — Use cause-effect analysis for complex relationships between inputs, conditions, actions, and outcomes.
- [ ] **USEQ-9C6928E4** — Use syntax-based testing for parsers, protocols, formats, commands, queries, expressions, and structured inputs.
- [ ] **USEQ-FC9B8485** — Use classification trees or equivalent models when many dimensions define input classes and configurations.
- [ ] **USEQ-DBCC55AA** — Use error guessing informed by defect history, architecture, implementation risks, incidents, and domain expertise.
- [ ] **USEQ-68B5113B** — Test repeated, duplicate, stale, delayed, missing, reordered, conflicting, and partially applied operations.
- [ ] **USEQ-34E35CBB** — Test create, read, update, delete, restore, archive, import, export, merge, split, and transfer operations where supported.
- [ ] **USEQ-E5DE8BF1** — Test roles, permissions, tenant relationships, ownership, object states, and data classifications systematically rather than through a few representative accounts.
- [ ] **USEQ-751D7D55** — Test normal, degraded, maintenance, recovery, migration, mixed-version, and shutdown modes.
- [ ] **USEQ-834CF191** — Record why each selected technique and coverage depth is sufficient for the identified risk.

### Structural coverage and change-impact adequacy

- [ ] **USEQ-29D7F4A9** — Measure structural coverage only where it provides useful evidence and never treat a percentage alone as proof of correctness.
- [ ] **USEQ-66961F40** — Select statement, branch, condition, decision, path, data-flow, call, exception, state, or requirements coverage according to component risk and language or runtime characteristics.
- [ ] **USEQ-E42BE90E** — Require coverage of security, authorization, error, recovery, cancellation, timeout, and boundary branches in critical code.
- [ ] **USEQ-35E8657D** — Investigate unexecuted critical code and either test it, remove it, justify it, or verify it by another method.
- [ ] **USEQ-7246E373** — Exclude generated, unreachable, defensive, or platform-specific code from metrics only with explicit rationale and review.
- [ ] **USEQ-48327490** — Measure changed-code coverage and impacted-behavior coverage for each release.
- [ ] **USEQ-ACC28191** — Use dependency, call-graph, data-flow, ownership, schema, and runtime evidence to identify change impact.
- [ ] **USEQ-D95A54B5** — Retest consumers, producers, migrations, operational tooling, documentation, and recovery procedures affected by a change.
- [ ] **USEQ-86276367** — Retain regression tests for every material defect and incident where a stable automated check is practical.
- [ ] **USEQ-018702B6** — Detect declining coverage, widening untested risk, and test deletion during review.
- [ ] **USEQ-0CF24972** — Avoid optimizing teams for raw coverage numbers that encourage trivial tests or conceal weak assertions.
- [ ] **USEQ-7FA69076** — Use coverage gaps to guide additional test design rather than lowering thresholds without analysis.
- [ ] **USEQ-931196A7** — Review whether mocks, stubs, fakes, or excluded processes prevent measured coverage from representing production behavior.
- [ ] **USEQ-49CDA8B6** — Confirm that tests execute the intended implementation rather than a test-only path or obsolete artifact.
- [ ] **USEQ-999DF568** — Combine structural coverage with mutation, requirements, risk, data, interface, state, and operational coverage for high-confidence decisions.

### Combinatorial interaction and configuration testing

- [ ] **USEQ-555AB134** — Identify parameters whose interactions can change behavior, including configuration, browser, device, locale, identity, role, tenant, data shape, dependency, network, deployment, and feature flags.
- [ ] **USEQ-9CC14B49** — Model valid and invalid values and constraints explicitly before generating combinations.
- [ ] **USEQ-A36B7B5C** — Use pairwise or higher-strength t-way coverage when exhaustive combinations are infeasible and interaction faults are plausible.
- [ ] **USEQ-8088BC71** — Select interaction strength based on risk, defect history, coupling, and consequence rather than applying pairwise testing mechanically.
- [ ] **USEQ-577B1DBC** — Include critical known combinations even when a generated covering array would otherwise omit them.
- [ ] **USEQ-6A1FC94F** — Verify the achieved combinatorial coverage after constraints, filtering, unavailable environments, and test failures are considered.
- [ ] **USEQ-4BA4B352** — Review constraints for accidental exclusion of valid or risky combinations.
- [ ] **USEQ-2755BC17** — Test invalid combinations and incompatible configurations for clear, safe rejection.
- [ ] **USEQ-A4B99182** — Use mixed-strength coverage so high-risk parameter groups receive deeper interaction coverage.
- [ ] **USEQ-EB2A1CBF** — Include version skew, rollout states, migration phases, and old/new client-server combinations.
- [ ] **USEQ-8D624C7F** — Include feature-flag interactions and default-state behavior when the flag service is unavailable.
- [ ] **USEQ-55D313EE** — Include combinations of accessibility settings, locale, viewport, input mode, consent, personalization, and authentication for user-facing products.
- [ ] **USEQ-E4DDFD53** — Include failure-injection parameters such as latency, timeout, partial response, quota, and retry state where practical.
- [ ] **USEQ-EB0C46E7** — Retain generated models, constraints, seeds, coverage reports, and failing combinations as evidence.
- [ ] **USEQ-E30A638C** — Update the interaction model after configuration, topology, feature, dependency, and supported-platform changes.

### Property-based, model-based, and invariant testing

- [ ] **USEQ-388A3A32** — Identify domain properties and invariants that should hold over broad input and state spaces.
- [ ] **USEQ-4297DBFF** — Generate valid, invalid, adversarial, boundary, and structurally diverse inputs from explicit models.
- [ ] **USEQ-470D5F97** — Ensure generators respect domain constraints while still exploring unusual combinations.
- [ ] **USEQ-06C48E9A** — Use shrinking or minimization to produce understandable counterexamples.
- [ ] **USEQ-C14B57D1** — Retain minimized failing seeds and examples as deterministic regression tests.
- [ ] **USEQ-0F46BD7A** — Model stateful systems with commands, preconditions, postconditions, transitions, and observable state.
- [ ] **USEQ-4CBA36F2** — Compare implementation behavior with an independently reviewed abstract model where feasible.
- [ ] **USEQ-9588447A** — Test algebraic properties such as identity, inverse, associativity, commutativity, monotonicity, idempotence, conservation, or round-trip behavior when applicable.
- [ ] **USEQ-1619E778** — Test serialization, parsing, migration, compression, encryption, export, and import for valid round-trip and canonicalization properties.
- [ ] **USEQ-6988A46D** — Test that authorization, tenant isolation, and privacy invariants hold for generated operations and states.
- [ ] **USEQ-27C90287** — Control generator bias so common values do not crowd out rare but important structures.
- [ ] **USEQ-98163D5F** — Bound generated input size and execution cost while retaining periodic deeper campaigns.
- [ ] **USEQ-BC09905F** — Validate model and generator correctness independently of the implementation.
- [ ] **USEQ-3C07904C** — Use coverage feedback to identify unvisited states and weak generators.
- [ ] **USEQ-C07FE87E** — Do not allow nondeterministic property tests to become an uninvestigated source of flaky failures.
- [ ] **USEQ-0560AF4B** — Use model-based testing standards and tooling only after confirming that the model represents the intended behavior rather than the current implementation defect.

### Metamorphic, differential, and consistency testing

- [ ] **USEQ-46CC6FF5** — Use metamorphic testing when a direct expected result is expensive, unavailable, probabilistic, or incomplete.
- [ ] **USEQ-B33F46EE** — Define transformations for which the relationship between source and follow-up outputs is known.
- [ ] **USEQ-F56C5FD1** — Review metamorphic relations for domain validity and avoid relations that merely encode implementation assumptions.
- [ ] **USEQ-C18B4EE2** — Test invariance under irrelevant changes such as ordering, formatting, identifier renaming, partitioning, or equivalent representations when appropriate.
- [ ] **USEQ-7A1E3A3A** — Test predictable changes under scaling, translation, duplication, aggregation, filtering, rotation, or perturbation where the domain supports them.
- [ ] **USEQ-B6117366** — Use differential testing across independent implementations, versions, providers, algorithms, configurations, or execution paths.
- [ ] **USEQ-7D8BA55E** — Investigate disagreements rather than automatically treating the majority result as correct.
- [ ] **USEQ-57974680** — Normalize only irrelevant differences before comparison and preserve semantically meaningful distinctions.
- [ ] **USEQ-7F024C3A** — Use cross-version differential tests to detect unintended compatibility and migration changes.
- [ ] **USEQ-6893D078** — Use cross-environment differential tests to detect configuration drift and platform-dependent behavior.
- [ ] **USEQ-6598B38C** — Test replicas, caches, indexes, reports, exports, analytics, and source systems for defined consistency relationships.
- [ ] **USEQ-9F067CE7** — Test reconciliation identities such as totals, balances, counts, checksums, and conservation constraints.
- [ ] **USEQ-87C8A49A** — Apply statistical comparison and confidence intervals to probabilistic or noisy outputs.
- [ ] **USEQ-5773BAC5** — Retain relation definitions, transformations, seeds, normalization rules, and discrepancies as evidence.
- [ ] **USEQ-282E9928** — Revalidate relations after domain, specification, algorithm, or model changes.

### Mutation testing and assertion-strength assessment

- [ ] **USEQ-FB314DAD** — Use mutation testing selectively to evaluate whether tests detect meaningful defects rather than merely execute code.
- [ ] **USEQ-531A3AB1** — Prioritize security, business-rule, calculation, parsing, state-transition, error-handling, and shared-library code for mutation analysis.
- [ ] **USEQ-A02A84FC** — Choose mutation operators that represent plausible defect classes for the implementation and domain.
- [ ] **USEQ-9549E767** — Exclude equivalent, invalid, uncompilable, or irrelevant mutants only through documented rules and review.
- [ ] **USEQ-9EAF1998** — Investigate surviving mutants as potential weak assertions, missing scenarios, unreachable code, unclear requirements, or equivalent behavior.
- [ ] **USEQ-EC571C66** — Add tests that kill meaningful survivors without coupling tests to private implementation details unnecessarily.
- [ ] **USEQ-26292CF1** — Use mutation scores as diagnostic evidence rather than a universal release target.
- [ ] **USEQ-1E2966BB** — Measure changed-code or risk-focused mutation coverage when full-system mutation is impractical.
- [ ] **USEQ-FF127148** — Control computational cost through sampling, prioritization, incremental analysis, and parallel execution without hiding high-risk survivors.
- [ ] **USEQ-8331BF1D** — Ensure test timeouts and flaky tests do not create false mutant kills or survivals.
- [ ] **USEQ-5CC91231** — Review automatically generated assertions to ensure they encode desired behavior rather than current accidental output.
- [ ] **USEQ-C74ADD8A** — Track recurring survivor patterns to improve test design standards and code architecture.
- [ ] **USEQ-D4A286B1** — Retain operator set, excluded mutants, equivalent-mutant decisions, results, and follow-up actions.
- [ ] **USEQ-AED1AA2D** — Re-run relevant mutation analysis after major refactoring of critical code or test infrastructure.
- [ ] **USEQ-1FC19709** — Do not use mutation testing to compensate for missing requirements, unrealistic environments, or absent end-to-end assurance.

### Fuzzing, robustness, and adversarial input campaigns

- [ ] **USEQ-41808EC5** — Identify parsers, protocols, APIs, file formats, deserializers, native interfaces, compilers, interpreters, validators, and state machines suitable for fuzzing.
- [ ] **USEQ-135F858E** — Build harnesses that reach the real processing logic with minimal unrelated setup.
- [ ] **USEQ-065E46BE** — Use valid seed corpora representing supported formats, historical defects, boundary cases, and real structural diversity.
- [ ] **USEQ-55F7A391** — Use structure-aware or grammar-aware generation where random bytes cannot reach meaningful states.
- [ ] **USEQ-E1E75EE3** — Use coverage guidance or equivalent feedback to identify unexplored code and states.
- [ ] **USEQ-D84D31C1** — Instrument for memory errors, undefined behavior, leaks, hangs, excessive resource use, assertion failures, and invariant violations.
- [ ] **USEQ-7C972922** — Set resource limits so fuzz infrastructure cannot damage shared systems or conceal denial-of-service defects.
- [ ] **USEQ-D7CB7744** — Run continuous or recurring campaigns for high-risk components rather than one-time pre-release fuzzing.
- [ ] **USEQ-27B484EB** — Minimize and deduplicate failures while preserving the original crashing input and environment.
- [ ] **USEQ-D3975F34** — Triage every unique failure to root cause and affected product versions.
- [ ] **USEQ-6602CBDA** — Retain fixed fuzz cases in regression corpora.
- [ ] **USEQ-5ADACCB5** — Test stateful sequences, authentication states, concurrency, retries, and protocol transitions where relevant.
- [ ] **USEQ-E7DDC5C7** — Fuzz both accepted and rejected inputs and verify safe rejection without partial side effects.
- [ ] **USEQ-2347989F** — Isolate fuzzing from production credentials, personal data, external customers, and uncontrolled outbound traffic.
- [ ] **USEQ-B4D243DD** — Measure campaign duration, execution rate, coverage growth, unique paths, untriaged failures, and corpus health without treating any one metric as proof of completeness.

### Concurrency, asynchronous, and distributed-systems assurance

- [ ] **USEQ-C8EE5F64** — Test concurrent reads, writes, updates, deletion, retries, cancellation, timeout, failover, and recovery against stated consistency and ordering guarantees.
- [ ] **USEQ-0F2BF3BE** — Exercise race windows by controlling scheduling, synchronization, delays, barriers, and fault timing rather than relying only on random load.
- [ ] **USEQ-A7961220** — Test lost updates, duplicate effects, stale reads, write skew, deadlock, livelock, starvation, priority inversion, and resource contention.
- [ ] **USEQ-2068317C** — Verify atomicity, isolation, durability, ordering, idempotency, deduplication, and exactly-once claims at the actual system boundary.
- [ ] **USEQ-C1796402** — Test message loss, duplication, reordering, delay, corruption, redelivery, poison messages, and consumer restarts.
- [ ] **USEQ-33BB177D** — Test clock skew, leap behavior, timer drift, delayed scheduling, and inconsistent time sources.
- [ ] **USEQ-C48DFA3E** — Test network partitions, split brain, leader changes, replica lag, quorum loss, and reconciliation.
- [ ] **USEQ-9559DD2C** — Test cancellation and shutdown while operations are in progress.
- [ ] **USEQ-4B2D2D23** — Verify that retries do not repeat non-idempotent side effects or create retry storms.
- [ ] **USEQ-A347CC98** — Use deterministic schedulers, model checkers, simulation, history analysis, or linearizability checking for high-risk concurrency where feasible.
- [ ] **USEQ-78CFBE7F** — Capture enough event, trace, version, and correlation data to reconstruct failing interleavings.
- [ ] **USEQ-B80F01F3** — Retain reproducible schedules, seeds, histories, and reduced scenarios for regression.
- [ ] **USEQ-38085CC5** — Test multi-region, mixed-version, and rolling-deployment interleavings.
- [ ] **USEQ-062D18B9** — Verify authorization and tenant context across asynchronous boundaries.
- [ ] **USEQ-D0B99C3C** — Confirm that recovery and reconciliation preserve business invariants and do not silently discard conflicting data.

### Fault injection, resilience, recovery, and chaos assurance

- [ ] **USEQ-0424F6D4** — Derive fault-injection scenarios from architecture, dependency maps, incident history, threat models, recovery objectives, and common failure modes.
- [ ] **USEQ-0A29EB40** — Test dependency timeout, refusal, throttling, malformed response, partial response, slow response, stale data, and authentication failure.
- [ ] **USEQ-84B2AB3C** — Test process crash, node loss, zone loss, region loss, storage failure, disk full, memory pressure, CPU exhaustion, connection exhaustion, and quota exhaustion as applicable.
- [ ] **USEQ-A60A22BB** — Test certificate expiry, secret rotation, key revocation, DNS failure, clock error, and control-plane unavailability.
- [ ] **USEQ-7797AE7B** — Verify timeouts, retries, circuit breakers, bulkheads, backpressure, load shedding, fallbacks, and graceful degradation.
- [ ] **USEQ-B3A4C8AE** — Verify that resilience controls preserve correctness, security, privacy, and tenant isolation rather than only availability.
- [ ] **USEQ-160599FA** — Test backup restoration, point-in-time recovery, failover, failback, rebuild, and clean-environment recovery.
- [ ] **USEQ-A936D8F1** — Measure actual recovery point, recovery time, data loss, duplicate effects, and service degradation.
- [ ] **USEQ-279CBD42** — Run experiments in progressively more realistic environments with explicit blast-radius, stop, and rollback controls.
- [ ] **USEQ-E29C86AA** — Prevent fault injection from contacting real users, moving real money, corrupting uncontrolled data, or violating contractual commitments.
- [ ] **USEQ-6A752B54** — Observe whether alerts, runbooks, on-call response, communications, and escalation work during the experiment.
- [ ] **USEQ-C46CE499** — Inject compound and correlated failures where shared dependencies make them plausible.
- [ ] **USEQ-1CEFFE2F** — Test recovery under peak or degraded capacity rather than only in idle environments.
- [ ] **USEQ-C05A5D50** — Retain experiment hypotheses, scope, safeguards, telemetry, outcomes, surprises, and corrective actions.
- [ ] **USEQ-04DF2D4E** — Repeat material experiments after topology, dependency, capacity, recovery, and control changes.

### Performance, capacity, endurance, and resource assurance

- [ ] **USEQ-F6D37FEA** — Define performance and capacity objectives per critical journey, interface, workload, tenant, and operational mode using meaningful percentiles and error tolerances.
- [ ] **USEQ-ABD5E7E6** — Use representative data volume, distribution, cache state, authentication, authorization, dependency behavior, network conditions, and client capabilities.
- [ ] **USEQ-B48C8A09** — Test normal, peak, burst, sustained, seasonal, launch, failure, recovery, and abusive workload profiles.
- [ ] **USEQ-968D920D** — Measure latency, throughput, errors, saturation, queueing, resource consumption, contention, and user-perceived responsiveness together.
- [ ] **USEQ-B12F69C6** — Test cold start, warm start, scale up, scale down, deployment, failover, cache miss, and degraded dependency conditions.
- [ ] **USEQ-8FCC3759** — Identify maximum safe capacity and how the system fails beyond it.
- [ ] **USEQ-9C7D9575** — Test load shedding, rate limiting, quotas, backpressure, prioritization, and noisy-neighbor isolation.
- [ ] **USEQ-2F1099C0** — Run endurance tests long enough to reveal leaks, fragmentation, drift, compaction, backlog, cache decay, and gradual degradation.
- [ ] **USEQ-BF863ED1** — Test the largest expected object, tenant, account, batch, file, report, query, and historical data range.
- [ ] **USEQ-3DD76FA3** — Test client performance on representative low-power devices, accessibility settings, high latency, packet loss, and limited bandwidth where relevant.
- [ ] **USEQ-0D187CDB** — Measure third-party scripts, services, models, analytics, and tags separately from first-party behavior.
- [ ] **USEQ-D108E212** — Set regression budgets and compare against a stable baseline using statistically meaningful methods.
- [ ] **USEQ-AED5CE75** — Confirm that test generators and monitoring do not become the bottleneck or distort results.
- [ ] **USEQ-8AE79A6F** — Correlate performance limits with cost, SLOs, provider quotas, autoscaling limits, and recovery capacity.
- [ ] **USEQ-1A70007C** — Retain workload model, generators, data, topology, configuration, versions, raw measurements, analysis, and confidence limits.

### Conformance, interoperability, compatibility, and protocol assurance

- [ ] **USEQ-2C95C98C** — Identify every normative external standard, protocol, schema, format, API contract, accessibility rule, and certification claim applicable to the product.
- [ ] **USEQ-326F8926** — Separate mandatory requirements, optional features, profiles, extensions, implementation-defined behavior, and prohibited behavior.
- [ ] **USEQ-14C96036** — Use official conformance suites where available and validate their version, scope, limitations, and expected results.
- [ ] **USEQ-E6518281** — Test both producer and consumer behavior at each interface.
- [ ] **USEQ-9647D489** — Test unknown fields, optional fields, extension points, ordering differences, whitespace, encoding, canonicalization, size limits, and future-version data.
- [ ] **USEQ-F3290C67** — Test backward, forward, mixed-version, downgrade, upgrade, and deprecation behavior.
- [ ] **USEQ-94A0A73B** — Test multiple independent implementations or providers to detect shared assumptions and vendor lock-in.
- [ ] **USEQ-DDCB79E8** — Test error semantics, retries, idempotency, timeouts, status codes, headers, content negotiation, and rate-limit behavior.
- [ ] **USEQ-DDC867E4** — Test malformed, ambiguous, duplicate, conflicting, and noncanonical representations.
- [ ] **USEQ-21F2CD0E** — Test locale, Unicode, time, number, identifier, and serialization interoperability.
- [ ] **USEQ-062F69DE** — Preserve unknown data where the contract requires forward compatibility.
- [ ] **USEQ-A0AEBC24** — Reject unsupported or unsafe versions clearly without insecure fallback.
- [ ] **USEQ-5132E054** — Record deviations, implementation choices, extensions, and customer-visible limitations.
- [ ] **USEQ-B85CC7B0** — Do not claim standards conformance outside the exact version, profile, scope, environment, and evidence tested.
- [ ] **USEQ-3F433ECC** — Repeat conformance testing after standard revisions, parser changes, dependency changes, and protocol upgrades.

### Static analysis, formal methods, and specification assurance

- [ ] **USEQ-D3B10800** — Select static analysis techniques according to defect classes, languages, runtimes, architecture, and risk.
- [ ] **USEQ-70E9D5AA** — Run type, data-flow, control-flow, taint, dependency, configuration, infrastructure, secret, license, and policy analysis as applicable.
- [ ] **USEQ-49154A5B** — Configure rules and severity thresholds explicitly rather than relying on unreviewed defaults.
- [ ] **USEQ-D330E2E8** — Treat suppressions as reviewed, attributable, time-bounded decisions with rationale.
- [ ] **USEQ-6A3C60D1** — Investigate recurring false positives as opportunities to improve rules, code patterns, or tool configuration.
- [ ] **USEQ-AE94E396** — Verify that generated code, vendored code, build scripts, migrations, infrastructure, queries, and configuration are included where they can affect quality.
- [ ] **USEQ-3C921989** — Use formal specifications for critical state machines, protocols, authorization rules, distributed algorithms, safety properties, or calculations where the consequence justifies it.
- [ ] **USEQ-F13CA22A** — State assumptions, invariants, preconditions, postconditions, liveness, safety, and environmental constraints precisely.
- [ ] **USEQ-729EFDDA** — Validate that the formal model corresponds to the implemented and deployed system.
- [ ] **USEQ-8D18549C** — Use model checking, theorem proving, symbolic execution, abstract interpretation, or runtime verification when they materially increase assurance.
- [ ] **USEQ-099AEFD2** — Independently review proofs, models, abstractions, solver assumptions, and unverified components.
- [ ] **USEQ-C679FD8C** — Test the boundaries between formally verified and unverified components.
- [ ] **USEQ-9EB4BECC** — Retain tool versions, rules, models, proof artifacts, assumptions, suppressions, and unresolved findings.
- [ ] **USEQ-80C8E969** — Re-run affected analyses after source, compiler, dependency, configuration, model, or platform changes.
- [ ] **USEQ-B8AE7F5F** — Avoid representing tool success as proof of properties outside the tool model or analyzed scope.

### Test automation architecture and reliability

- [ ] **USEQ-C2DD853F** — Treat test code, harnesses, fixtures, generators, simulators, mocks, and pipelines as production-quality engineering assets.
- [ ] **USEQ-8987771C** — Give test infrastructure owners, review requirements, version control, documentation, observability, support, and lifecycle plans.
- [ ] **USEQ-96CDA61D** — Keep tests independent enough that one failure does not create misleading cascades.
- [ ] **USEQ-BFA06C49** — Control time, randomness, network, identity, external services, shared state, and parallel execution to improve reproducibility.
- [ ] **USEQ-CB1343BF** — Use deterministic seeds and retain them for generated and randomized failures.
- [ ] **USEQ-DF29F97A** — Avoid fixed sleeps; synchronize on observable conditions with bounded timeouts.
- [ ] **USEQ-123DA8C7** — Make failures explain which requirement, step, input, state, and observation differed.
- [ ] **USEQ-A6EE8FB6** — Capture diagnostics without exposing secrets or unnecessary personal data.
- [ ] **USEQ-4A334C27** — Prevent tests from passing because assertions were skipped, exceptions swallowed, fixtures failed, or dependencies were unavailable.
- [ ] **USEQ-9B3A9AA9** — Fail closed when required test setup, data, instrumentation, or result publication is incomplete.
- [ ] **USEQ-95988B28** — Verify mocks and simulators against real contracts and periodically against real dependencies.
- [ ] **USEQ-65CF6581** — Minimize shared mutable test data and clean resources reliably after success, failure, cancellation, and timeout.
- [ ] **USEQ-05E81260** — Run tests in parallel only when isolation and ordering assumptions are explicit and verified.
- [ ] **USEQ-F2F23289** — Monitor test duration, queue time, retry rate, flakiness, infrastructure failures, cost, and diagnostic quality.
- [ ] **USEQ-9AFF0C25** — Remove obsolete, duplicated, low-value, or misleading tests while preserving necessary coverage.
- [ ] **USEQ-E0D501EC** — Secure test systems so untrusted code and test data cannot access production secrets or damage production resources.

### Test data, environments, and production representativeness

- [ ] **USEQ-B5965E6B** — Define which production characteristics each test environment must reproduce and which differences remain material limitations.
- [ ] **USEQ-FEB0BCAC** — Version infrastructure, configuration, schema, dependencies, feature flags, identities, and test data with the tested artifact.
- [ ] **USEQ-DDE86544** — Use synthetic, generated, masked, or specifically authorized data rather than uncontrolled production copies.
- [ ] **USEQ-8FA110C8** — Preserve realistic distributions, relationships, cardinality, sparsity, skew, encoding, history, and edge cases without retaining unnecessary personal data.
- [ ] **USEQ-C5C286DB** — Include malformed, incomplete, duplicated, conflicting, stale, future-dated, historical, deleted, archived, and maximum-size records.
- [ ] **USEQ-366D4955** — Include multiple tenants, roles, plans, locales, currencies, time zones, consent states, lifecycle states, and accessibility settings as applicable.
- [ ] **USEQ-C58ED0EB** — Protect test data according to its actual sensitivity and delete it according to retention rules.
- [ ] **USEQ-AE04D918** — Prevent nonproduction notifications, payments, webhooks, integrations, and automated actions from reaching real users or counterparties.
- [ ] **USEQ-A34619FA** — Reset environments reproducibly and detect residual state or cross-test contamination.
- [ ] **USEQ-A6D5B399** — Use production-like topology and security boundaries for migration, performance, resilience, and deployment tests.
- [ ] **USEQ-0691B9EF** — Validate provider sandboxes and emulators against production behavior and document differences.
- [ ] **USEQ-60A4CD09** — Reserve controlled production testing for cases that cannot be represented elsewhere and apply blast-radius, monitoring, and rollback safeguards.
- [ ] **USEQ-3BCF09B5** — Detect environment drift and invalidate evidence when material drift occurs.
- [ ] **USEQ-2B7192B0** — Retain environment manifests, data-generation rules, masking methods, seeds, and known deviations.
- [ ] **USEQ-A8FC5F26** — Ensure test environments remain available, supported, patched, and affordable enough for required assurance.

### Flaky tests, quarantine, and trust in results

- [ ] **USEQ-B74D5F78** — Define a flaky test as an engineering defect requiring ownership, diagnosis, and remediation rather than a normal condition.
- [ ] **USEQ-DFD85174** — Measure flakiness by test, suite, environment, branch, dependency, and failure signature.
- [ ] **USEQ-7679C932** — Differentiate product defects, test defects, environment failures, dependency failures, and infrastructure failures.
- [ ] **USEQ-9CE7EDE5** — Do not automatically retry failures without preserving and reporting the first failure.
- [ ] **USEQ-01EEF7ED** — Limit retries and ensure a pass-after-retry remains visible.
- [ ] **USEQ-8825082A** — Quarantine only when the test cannot block unrelated delivery and the missing coverage is understood.
- [ ] **USEQ-676775D7** — Require a quarantine owner, reason, date, risk assessment, compensating coverage, and expiry.
- [ ] **USEQ-48996A09** — Prevent quarantined critical tests from silently disappearing from release evidence.
- [ ] **USEQ-58D50432** — Fix root causes such as timing assumptions, shared state, uncontrolled data, order dependence, resource leakage, and nondeterministic external systems.
- [ ] **USEQ-A3A1C911** — Use statistical analysis carefully to distinguish real intermittent product failures from test noise.
- [ ] **USEQ-CF484D87** — Retain logs, traces, screenshots, dumps, seeds, schedules, environment state, and previous outcomes for intermittent failures.
- [ ] **USEQ-8282F21D** — Escalate repeated intermittent failures in security, integrity, financial, safety, and authorization paths as release blockers.
- [ ] **USEQ-C0B34752** — Review whether test parallelism, virtualization, clocks, resource contention, and provider quotas create artificial instability.
- [ ] **USEQ-EA5C03BB** — Track quarantine age and remove or repair expired tests.
- [ ] **USEQ-A480FEDA** — Restore trust by demonstrating stable repeated execution after remediation rather than simply resetting history.

### Continuous assurance and production verification

- [ ] **USEQ-0541BA5A** — Define which critical properties must be verified continuously after deployment because preproduction testing cannot guarantee them.
- [ ] **USEQ-624D5AC5** — Use production-safe smoke tests and synthetic journeys for authentication, authorization, critical transactions, dependencies, and user-visible outcomes.
- [ ] **USEQ-8F1C7642** — Verify the deployed artifact digest, configuration, feature flags, schema, migrations, certificates, routes, and security controls after deployment.
- [ ] **USEQ-9FC91F6A** — Compare canary and control cohorts using predefined correctness, reliability, performance, business, security, and support indicators.
- [ ] **USEQ-3066B3F2** — Use shadow, replay, mirrored, or dark traffic only with privacy, consent, isolation, and side-effect controls.
- [ ] **USEQ-53C69573** — Validate invariants, reconciliation totals, data freshness, queue lag, duplicate effects, and silent failure indicators continuously.
- [ ] **USEQ-59BD56C3** — Use runtime assertions and policy enforcement for properties that must never be violated, with safe failure behavior.
- [ ] **USEQ-0A3F86F0** — Monitor test and verification telemetry for gaps, sampling bias, disabled checks, stale baselines, and pipeline failure.
- [ ] **USEQ-09C1C843** — Define automatic or human stop, rollback, kill-switch, and incident-declaration thresholds before rollout.
- [ ] **USEQ-DD71888C** — Prevent experimentation and monitoring from exposing sensitive data or materially degrading the service.
- [ ] **USEQ-A6B8EFC8** — Sample real outcomes for human or independent review where automation cannot assess semantic quality.
- [ ] **USEQ-50075879** — Correlate escaped defects and incidents with missing, ineffective, or bypassed pre-release assurance.
- [ ] **USEQ-619B8141** — Reassess supported configurations and user journeys using field data without allowing popularity to exclude minority or accessibility needs.
- [ ] **USEQ-4382A92C** — Retain production verification evidence with the exact release and observation period.
- [ ] **USEQ-23D4E8D7** — Treat unexplained anomalies as unresolved risk rather than assuming they are harmless noise.

### Defect management, root cause, and assurance effectiveness

- [ ] **USEQ-4050C11D** — Record each defect with affected requirement, users, versions, environments, data, severity, detectability, workaround, and owner.
- [ ] **USEQ-BC455E22** — Classify impact using user harm, security, privacy, integrity, financial, accessibility, safety, legal, operational, and reputational dimensions.
- [ ] **USEQ-F0B8342C** — Prioritize by actual risk and exposure rather than age or raw count alone.
- [ ] **USEQ-CD7C7286** — Reproduce defects against the exact released or candidate artifact when possible.
- [ ] **USEQ-1B6B2359** — Identify the first point in the lifecycle where the defect could reasonably have been prevented or detected.
- [ ] **USEQ-AC327620** — Perform root-cause analysis for serious, recurring, systemic, or escaped defects.
- [ ] **USEQ-ECAA299C** — Distinguish triggering condition, proximate technical cause, contributing conditions, organizational cause, and control failure.
- [ ] **USEQ-279BAC2D** — Search for variants across shared code, products, versions, configurations, data, and processes.
- [ ] **USEQ-59366EFE** — Verify fixes through targeted tests, broader regression, and production monitoring where appropriate.
- [ ] **USEQ-9B2BE606** — Measure escaped defects, recurrence, reopen rate, detection stage, time to detect, time to repair, customer exposure, and corrective-action completion.
- [ ] **USEQ-F2870EFC** — Use metrics to improve the system rather than rank individuals or encourage defect concealment.
- [ ] **USEQ-28C4F691** — Review whether tests failed to exist, failed to run, failed to assert, failed to represent production, or were ignored.
- [ ] **USEQ-003E4BB9** — Convert lessons into requirements, design rules, coding standards, test techniques, tooling, runbooks, and training.
- [ ] **USEQ-2BC9A7C5** — Close corrective actions only after evidence demonstrates effectiveness.
- [ ] **USEQ-7945FABA** — Periodically assess whether the assurance portfolio finds the defect classes and user harms that matter most.

### Assurance release blockers and evidence package

- [ ] **USEQ-ECC22163** — Block release when a critical requirement, user journey, safety property, authorization boundary, data invariant, recovery objective, or mandatory standard lacks current verification evidence.
- [ ] **USEQ-06A9124D** — Block release when tests ran against a different artifact, configuration, schema, feature-flag state, environment, or dependency set than the intended production release.
- [ ] **USEQ-4D20534E** — Block release when required tests were skipped, suppressed, quarantined, or passed only after unexplained retries.
- [ ] **USEQ-9BC293BB** — Block release when a serious intermittent, nondeterministic, or environment-dependent failure remains unexplained.
- [ ] **USEQ-F84F4311** — Block release when expected results, test data, tools, or environments are known to be invalid or unrepresentative.
- [ ] **USEQ-9D68BDCB** — Block release when performance, capacity, resilience, migration, rollback, or restoration behavior is unknown for expected operating conditions.
- [ ] **USEQ-941A1A9C** — Block release when critical test coverage relies exclusively on mocks, snapshots, superficial assertions, or aggregate percentages.
- [ ] **USEQ-C6187362** — Block release when material conformance, interoperability, accessibility, security, privacy, or regulatory claims exceed the tested scope.
- [ ] **USEQ-D601267F** — Retain the test strategy, scope, requirements traceability, models, data, environments, tools, versions, results, defects, exclusions, limitations, exceptions, and approvals.
- [ ] **USEQ-F55894B6** — Retain raw or reproducible evidence for high-impact calculations, performance measurements, resilience experiments, and independent assessments.
- [ ] **USEQ-DC9036BC** — Document which defect classes and conditions remain insufficiently tested.
- [ ] **USEQ-62ABC3B0** — Require an authorized owner to accept residual assurance uncertainty explicitly and time-bound any compensating monitoring.
- [ ] **USEQ-3D261C6C** — Verify that post-release monitoring and incident response can detect and contain risks that could not be fully tested.
- [ ] **USEQ-2AE6772C** — Make the final assurance statement specific to the tested release and never describe testing as proof that no defects exist.

## Final Gap Closure — Assurance Cases, Independent Assurance, and Conformance Claims

_Consolidated from `final consolidated corpus/06-testing-verification-conformance-assurance.md#Final Gap Closure — Assurance Cases, Independent Assurance, and Conformance Claims`; 107 non-duplicative controls._

### Assurance strategy and scope

- [ ] **USEQ-67C9E2FD** — Define the qualities, harms, obligations, claims, systems, versions, environments, users, and operating conditions that require assurance.
- [ ] **USEQ-FB59BCD1** — Scale assurance rigor according to consequence, uncertainty, novelty, exposure, complexity, and reversibility.
- [ ] **USEQ-0CF01C55** — Identify the stakeholders who rely on each assurance conclusion and the decisions it supports.
- [ ] **USEQ-D99B9218** — Define the required confidence level and what residual uncertainty is acceptable.
- [ ] **USEQ-0EB55A7C** — Distinguish evidence of implementation, evidence of control operation, evidence of outcome, and evidence of absence of known failure.
- [ ] **USEQ-1283600A** — Use multiple complementary assurance methods when one method cannot support the required confidence.
- [ ] **USEQ-19F7A4D0** — Define independence requirements for design review, verification, validation, audit, certification, and risk acceptance.
- [ ] **USEQ-BFC07AAB** — Ensure assurance planning begins with requirements and architecture rather than after implementation.
- [ ] **USEQ-511BD294** — Budget time, tools, environments, representative data, specialists, and remediation capacity for assurance.
- [ ] **USEQ-734F6033** — Record limitations, exclusions, assumptions, dependencies, and evidence freshness.
- [ ] **USEQ-2B085AD7** — Treat assurance as a lifecycle activity that continues through change, operation, incident, and retirement.

### Assurance-case structure

- [ ] **USEQ-A3008EA4** — Express each material assurance conclusion as a clear, bounded, testable claim.
- [ ] **USEQ-6B2A8983** — State the exact system, release, configuration, environment, period, and operational context covered by the claim.
- [ ] **USEQ-01943631** — Decompose broad claims into subclaims that can be supported independently.
- [ ] **USEQ-67DF1B66** — Provide an explicit argument explaining why the evidence supports each claim.
- [ ] **USEQ-D6EC7E3B** — Identify assumptions, context, definitions, inference rules, and dependencies used in the argument.
- [ ] **USEQ-ACEEAD05** — Link every subclaim to relevant evidence or record it as unsupported.
- [ ] **USEQ-0A2301F0** — Avoid circular reasoning in which a claim is used as evidence for itself.
- [ ] **USEQ-5DDCFB45** — Avoid replacing an outcome claim with evidence that a process or document exists.
- [ ] **USEQ-6CC47215** — Include counterclaims, defeaters, alternative explanations, and known contrary evidence.
- [ ] **USEQ-100AD69F** — Explain why identified defeaters are eliminated, controlled, monitored, or accepted.
- [ ] **USEQ-5C92AC6C** — Show how claims interact across security, privacy, safety, accessibility, reliability, correctness, and compliance.
- [ ] **USEQ-B9AA39F0** — Distinguish claims about design intent from claims about deployed behavior.
- [ ] **USEQ-14D36863** — Distinguish expected behavior from demonstrated behavior.
- [ ] **USEQ-E917726C** — Keep the assurance case readable and reviewable by the intended decision-makers.
- [ ] **USEQ-3C79C82C** — Version the assurance case and tie it to immutable release and evidence identifiers.

### Evidence quality and sufficiency

- [ ] **USEQ-993D84C7** — Define evidence acceptance criteria before relying on results.
- [ ] **USEQ-F473BC2E** — Prefer evidence produced from the exact production artifact and representative configuration.
- [ ] **USEQ-7D363841** — Record evidence source, method, tool, version, operator, environment, data, time, scope, and result.
- [ ] **USEQ-E6AD6762** — Establish provenance from requirement and source through build, deployment, execution, and observation.
- [ ] **USEQ-FF35E31B** — Protect evidence against unauthorized alteration and deletion.
- [ ] **USEQ-98FB3D02** — Verify tool calibration, validation, configuration, and known limitations where tool output is material.
- [ ] **USEQ-043971A5** — Confirm that test or audit samples represent the relevant population and risk.
- [ ] **USEQ-ADA14CE9** — Quantify uncertainty, confidence, coverage, error rates, and detection limits where meaningful.
- [ ] **USEQ-F16DCC81** — Distinguish absence of observed failure from evidence that failure is sufficiently unlikely.
- [ ] **USEQ-2963F4D9** — Do not treat test count, code coverage, policy count, or scan count alone as proof of quality.
- [ ] **USEQ-453FDEAE** — Combine static, dynamic, analytical, operational, and human evidence as appropriate.
- [ ] **USEQ-E0A35E0B** — Use independent reproduction for high-consequence calculations, models, migrations, or tests where practical.
- [ ] **USEQ-A6DE513F** — Resolve contradictory evidence rather than choosing the favorable result silently.
- [ ] **USEQ-2BADD314** — Preserve failed, inconclusive, excluded, and superseded evidence with its disposition.
- [ ] **USEQ-EB04BD9A** — Define evidence validity periods and events that make evidence stale.
- [ ] **USEQ-C278CEB3** — Recollect evidence after material changes to source, dependencies, configuration, data, environment, threat, or usage.
- [ ] **USEQ-1BA3206E** — Ensure evidence remains retrievable and interpretable for the required retention period.

### Assurance integrity levels and proportional rigor

- [ ] **USEQ-A309E8E6** — Classify components, functions, data, and decisions by the consequence of failure or compromise.
- [ ] **USEQ-E19EECAD** — Define assurance rigor, independence, coverage, formality, and evidence retention for each class.
- [ ] **USEQ-A72314F7** — Apply the highest relevant assurance level to shared components whose failure can affect higher-criticality functions.
- [ ] **USEQ-B3F7905C** — Prevent low-assurance dependencies from silently determining high-assurance outcomes.
- [ ] **USEQ-4D9F7786** — Use stronger specification, traceability, review, analysis, and verification for higher-criticality behavior.
- [ ] **USEQ-96C2BFD1** — Use diverse or independent methods where common-mode verification failure would be unacceptable.
- [ ] **USEQ-4FF2DB8D** — Require qualified specialists for domains such as safety, cryptography, accessibility, privacy, finance, or formal verification when warranted.
- [ ] **USEQ-C6D28252** — Define when formal methods, exhaustive analysis, proof, model checking, or mathematically justified bounds are required.
- [ ] **USEQ-B7FB7430** — Document why a selected assurance level is sufficient.
- [ ] **USEQ-07963687** — Review integrity classification when architecture, exposure, users, or consequence changes.

### Independence, competence, and conflict controls

- [ ] **USEQ-3ABADEBD** — Define which assurance activities may be self-reviewed and which require organizational or technical independence.
- [ ] **USEQ-7E738595** — Prevent the person solely responsible for delivery outcomes from being the only approver of material residual risk.
- [ ] **USEQ-ADBD2C5E** — Ensure independent reviewers can access complete evidence, source material, environments, and relevant personnel.
- [ ] **USEQ-24426703** — Disclose reviewer conflicts of interest, incentives, prior involvement, and supplier relationships.
- [ ] **USEQ-9F1C648A** — Recuse, rotate, or supplement reviewers when independence is impaired.
- [ ] **USEQ-25DE0302** — Verify reviewer competence in the system, assurance method, and relevant domain.
- [ ] **USEQ-15D97639** — Protect reviewers from retaliation or pressure to weaken findings.
- [ ] **USEQ-ED2D15CF** — Preserve dissenting opinions and unresolved findings in the decision record.
- [ ] **USEQ-AFC0B80B** — Require management to respond explicitly to material adverse assurance conclusions.
- [ ] **USEQ-1D117647** — Separate advisory work from certification or audit decisions where the same party’s objectivity would be compromised.
- [ ] **USEQ-0C690882** — Periodically evaluate assurance-provider quality, consistency, and missed-defect patterns.

### Conformance, certification, and claims

- [ ] **USEQ-03CE6FB0** — Identify the exact normative requirements, version, profile, scope, and options used for a conformance claim.
- [ ] **USEQ-273415BD** — Distinguish full conformance, partial conformance, compatibility, alignment, certification, attestation, and self-assessment.
- [ ] **USEQ-EF9F966C** — Do not imply certification by citing a standard that was only used as guidance.
- [ ] **USEQ-F7D4B6BC** — Map every claimed requirement to evidence and disposition.
- [ ] **USEQ-FAD45C35** — Record permitted deviations, extensions, optional features, and inherited controls.
- [ ] **USEQ-70F1E6CA** — Test interoperability when conformance is intended to enable independent implementations to work together.
- [ ] **USEQ-0AAEEBDE** — Verify that declared exclusions do not invalidate the intended claim.
- [ ] **USEQ-93D22446** — State the time period and product versions covered by an audit or certificate.
- [ ] **USEQ-263D9220** — Ensure marketing, procurement, customer, and regulatory statements match the assurance scope exactly.
- [ ] **USEQ-2BC823F8** — Track certificate expiry, surveillance, conditions, exceptions, and renewal obligations.
- [ ] **USEQ-1281EF2B** — Reassess certification impact before material architecture, supplier, location, process, or product change.
- [ ] **USEQ-30397AF4** — Withdraw or correct claims promptly when evidence no longer supports them.
- [ ] **USEQ-2A1264FC** — Maintain procedures for regulator, customer, auditor, and public verification of appropriate claims.

### Assurance across suppliers and inherited controls

- [ ] **USEQ-01F30AE1** — Identify every claim that depends on a supplier, shared platform, managed service, open-source component, or customer-operated control.
- [ ] **USEQ-229D6F8B** — Define the boundary between inherited, shared, and directly operated controls.
- [ ] **USEQ-31EDBB67** — Obtain evidence sufficient for the actual dependency and risk rather than relying on a generic certificate alone.
- [ ] **USEQ-6EBAC76D** — Verify that supplier assurance scope includes the service, region, data, configuration, and time period used.
- [ ] **USEQ-2B60BF65** — Assess gaps between supplier controls and product requirements.
- [ ] **USEQ-4B7E3419** — Test integration assumptions and customer-configurable controls independently.
- [ ] **USEQ-7D7FB635** — Monitor supplier changes, incidents, expired attestations, and scope reductions.
- [ ] **USEQ-A892799B** — Include supplier evidence availability and transition rights in contracts.
- [ ] **USEQ-44E74665** — Maintain compensating controls for evidence that cannot be obtained.
- [ ] **USEQ-6647217A** — Do not transfer accountability merely because an activity is outsourced.

### Continuous and operational assurance

- [ ] **USEQ-7FBB7982** — Define production indicators that continue to support or challenge assurance claims.
- [ ] **USEQ-79BEB4A9** — Monitor correctness, safety, security, privacy, accessibility, reliability, performance, data quality, and abuse outcomes relevant to the case.
- [ ] **USEQ-272E383B** — Link incidents, near misses, defects, complaints, overrides, exceptions, and drift to affected claims.
- [ ] **USEQ-B354E1E9** — Reopen assurance conclusions when operational evidence contradicts assumptions.
- [ ] **USEQ-8F326500** — Detect configuration and deployment changes that invalidate prior evidence.
- [ ] **USEQ-F2A92C95** — Maintain continuous or scheduled control testing where risk warrants it.
- [ ] **USEQ-69FFB927** — Use production verification safely and within privacy and reliability constraints.
- [ ] **USEQ-DC2E122D** — Preserve enough historical telemetry to evaluate trends and reconstruct claim validity.
- [ ] **USEQ-83870806** — Review assurance cases after incidents, major changes, new threats, or changed obligations.
- [ ] **USEQ-73001BBE** — Retire claims and evidence when the corresponding system or obligation ends, subject to retention requirements.

### Assurance review and release gates

- [ ] **USEQ-BACE71B2** — Conduct a structured challenge review of each material assurance case.
- [ ] **USEQ-E4B0B27B** — Verify claims are precise, evidence is relevant, arguments are valid, and limitations are visible.
- [ ] **USEQ-CA8CE5F9** — Include qualified representatives from product, engineering, operations, security, privacy, accessibility, data, legal, and safety as applicable.
- [ ] **USEQ-6C6059E7** — Confirm that open defects and exceptions are reflected in the claims rather than hidden outside the case.
- [ ] **USEQ-5797F11A** — Treat unsupported critical claims, stale evidence, unresolved defeaters, or inadequate independence as no-go conditions.
- [ ] **USEQ-54B2A382** — Require explicit risk acceptance for residual uncertainty that remains within policy.
- [ ] **USEQ-498728F8** — Record who reviewed, challenged, approved, rejected, or conditionally accepted the case.
- [ ] **USEQ-9513B3E9** — Preserve the exact assurance case with the released artifact and production configuration.
- [ ] **USEQ-775F8E10** — Define post-release observations required before a provisional claim becomes final.
- [ ] **USEQ-A4A4897A** — Revoke approval when required conditions, monitoring, or compensating controls fail.

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
