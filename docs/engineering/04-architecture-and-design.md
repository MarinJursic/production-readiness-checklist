# Architecture and design

_Phase 4 of 16 in the [complete engineering review](00-overview.md)._

System understanding, modularity, state, distribution, reliability, scalability, interoperability, and sustainability.

Assess every control as **Pass**, **Fail**, **Blocked**, or **Not Applicable**. A Pass needs current evidence for the exact product, revision, artifact, configuration, environment, and data state under review.

## System Understanding and Threat Modeling

_Consolidated from `quality standards/04-architecture-design/01-system-understanding-and-threat-modeling.md`; 13 non-duplicative controls._

### Universal controls

- [ ] **USEQ-B7EF282B** — Maintain current system-context, component, deployment, and end-to-end data-flow diagrams.
- [ ] **USEQ-85AF022E** — Mark trust boundaries, external interfaces, administrative interfaces, support interfaces, authentication points, authorization points, and encryption boundaries.
- [ ] **USEQ-7E6F1FE2** — Mark where sensitive data is collected, transformed, stored, cached, logged, exported, backed up, transferred, and deleted.
- [ ] **USEQ-ABB67AFA** — Identify every stateful component, asynchronous boundary, queue, common dependency, shared failure domain, and control-plane dependency.
- [ ] **USEQ-7FEE35E5** — Identify single points of failure, circular dependencies, and hidden dependencies.
- [ ] **USEQ-6BDC92D9** — Record material architecture decisions, rejected alternatives, and tradeoffs.
- [ ] **USEQ-92FE8B5B** — Update diagrams after material changes and verify documentation against the deployed environment.

### Threat and abuse modeling

- [ ] **USEQ-71BDCCBE** — Include all actors, trust boundaries, entry points, data stores, external systems, and administrative paths.
- [ ] **USEQ-25F8A29F** — Model identity spoofing, unauthorized modification, information disclosure, denial of service, repudiation, and privilege escalation.
- [ ] **USEQ-47D7BDB8** — Model business-logic abuse, insider threats, compromised administrators, dependency compromise, and build-system compromise.
- [ ] **USEQ-A1FE2F54** — Model tenant escape, scraping, spam, fraud, automation, economic abuse, backup attacks, and support-channel social engineering.
- [ ] **USEQ-E37A2788** — Convert identified threats into requirements, controls, and tests.
- [ ] **USEQ-60FB538E** — Revisit the threat model after material scope or design changes.

## Performance and User-Experience Efficiency

_Consolidated from `quality standards/04-architecture-design/02-performance-and-user-experience-efficiency.md`; 17 non-duplicative controls._

### Universal controls

- [ ] **USEQ-266638A7** — Performance objectives are defined for every critical journey.
- [ ] **USEQ-50B6A12E** — Objectives use meaningful percentiles rather than averages alone.
- [ ] **USEQ-C8C2C6FB** — Frontend and backend performance budgets are documented.
- [ ] **USEQ-FE0DF97B** — Public pages meet approved loading, interaction, visual-stability, rendering, and responsiveness targets.
- [ ] **USEQ-5EFABA0A** — API latency targets are defined by endpoint or journey.
- [ ] **USEQ-29D46E0C** — Database-query latency, count, and volume are measured.
- [ ] **USEQ-2F6D8929** — Excessive queries, repeated fetches, blocking operations, and unbounded responses are eliminated.
- [ ] **USEQ-DED09D7D** — Payload sizes are bounded and compression is enabled appropriately.
- [ ] **USEQ-E555726F** — Images, fonts, scripts, styles, media, and critical resources are delivered efficiently.
- [ ] **USEQ-757AF87D** — Noncritical work is deferred without harming accessibility, correctness, or user trust.
- [ ] **USEQ-F2B2306F** — Cache and CDN behavior is validated.
- [ ] **USEQ-C79EF036** — Connection, thread, worker, and resource pools are sized and monitored.
- [ ] **USEQ-770BCF5A** — Memory, CPU, connection, file-descriptor, storage, and resource leaks are investigated.
- [ ] **USEQ-7ECBBCC3** — Long-running and resource-intensive requests have limits.
- [ ] **USEQ-F1EB0182** — Performance is tested with realistic authentication, authorization, data volume, cache state, devices, and networks.
- [ ] **USEQ-115511BC** — Client telemetry does not materially harm performance.
- [ ] **USEQ-7D6F12FE** — Performance objectives connect to user and business outcomes.

## Capacity, Scalability, and Overload Control

_Consolidated from `quality standards/04-architecture-design/03-capacity-scalability-and-overload-control.md`; 14 non-duplicative controls._

### Universal controls

- [ ] **USEQ-7C17A636** — Forecast data, file, index, queue, log, telemetry, and backup growth.
- [ ] **USEQ-59F42DD9** — Identify the maximum safe capacity of every critical component.
- [ ] **USEQ-5AFFD98A** — Load tests use realistic request mixes, authorization, data distributions, and dependency behavior.
- [ ] **USEQ-C3395046** — Provider quotas, account limits, concurrency limits, and regional limits are reviewed.
- [ ] **USEQ-EDA216F8** — Database connections, locks, storage, replicas, and input/output have sufficient headroom.
- [ ] **USEQ-35154E0E** — Queue, cache, load-balancer, network, worker, and telemetry systems have sufficient headroom.
- [ ] **USEQ-1E625D17** — Scaling does not overload databases or dependencies.
- [ ] **USEQ-2AD26751** — Capacity remains sufficient during the failures the architecture promises to survive.
- [ ] **USEQ-46DF8678** — Rate limits and per-user or per-tenant quotas protect shared capacity fairly.
- [ ] **USEQ-18CF671E** — Retry storms and thundering-herd behavior are prevented.
- [ ] **USEQ-A072D92D** — DDoS and volumetric-abuse scenarios are reviewed.
- [ ] **USEQ-980A606F** — Capacity alerts provide enough time to act.
- [ ] **USEQ-FC5F896F** — Expected and peak operational cost are modeled.
- [ ] **USEQ-EEE57AB0** — Cost controls do not silently terminate essential operation.

## Reliability, Resilience, and Failure Engineering

_Consolidated from `quality standards/04-architecture-design/04-reliability-resilience-and-failure-engineering.md`; 17 non-duplicative controls._

### Universal controls

- [ ] **USEQ-7ED2ED65** — SLOs cover important availability, latency, correctness, durability, freshness, and quality outcomes.
- [ ] **USEQ-F64F6149** — Product and engineering stakeholders approve the SLOs.
- [ ] **USEQ-745E0AA2** — Critical dependency reliability is compatible with the product SLO or engineered around.
- [ ] **USEQ-D2574104** — Retries are bounded, delayed, jittered where appropriate, and safe or idempotent.
- [ ] **USEQ-2D2086B8** — Bulkheads or equivalent isolation prevent one workload from exhausting all capacity.
- [ ] **USEQ-4D1F0F17** — Optional functionality failure does not unnecessarily fail the whole product.
- [ ] **USEQ-37A811EC** — Instances stop receiving new work before shutdown.
- [ ] **USEQ-942757E2** — In-flight work is completed, transferred, canceled safely, or retried during shutdown.
- [ ] **USEQ-77C403FB** — Shared control planes, identity systems, networks, storage, DNS, providers, and dependencies are included in failure analysis.
- [ ] **USEQ-5DDDC794** — Network delay, packet loss, partition, and DNS failure have been tested.
- [ ] **USEQ-2BF79B16** — Dependency timeout, malformed response, throttling, partial failure, and outage have been tested.
- [ ] **USEQ-43D7583F** — Database failover, replication lag, and split-brain prevention have been tested.
- [ ] **USEQ-282499CC** — Disk-full, memory pressure, CPU exhaustion, file-descriptor exhaustion, connection exhaustion, and quota exhaustion have been tested.
- [ ] **USEQ-93CA6FC2** — Clock skew, time-service failure, certificate expiry, and key expiry have been considered.
- [ ] **USEQ-DD1DC2AE** — Scheduled jobs handle overlapping and missed execution.
- [ ] **USEQ-997091C4** — Controlled fault injection or equivalent resilience tests are completed for high-risk paths.
- [ ] **USEQ-90967BDC** — Recovery time from tested failures meets approved objectives.

## Architecture Governance and Decision Records

_Consolidated from `quality standards/04-architecture-design/05-architecture-governance-and-decision-records.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-1A22CF85** — Identify the system of interest, environment, stakeholders, concerns, boundaries, assumptions, and lifecycle stage.
- [ ] **USEQ-75264C89** — Use viewpoints and models suited to each stakeholder concern rather than one overloaded diagram.
- [ ] **USEQ-E4AB9272** — Describe structure, behavior, data, deployment, trust, operations, failure, and evolution where material.
- [ ] **USEQ-58B869F6** — Maintain consistency among architecture descriptions, source, infrastructure, configuration, data flows, and deployment reality.
- [ ] **USEQ-B5A29AFF** — Record material decisions with context, options, forces, evidence, trade-offs, consequences, owner, status, and revisit triggers.
- [ ] **USEQ-1FF9E601** — Separate architecture constraints from temporary implementation choices.
- [ ] **USEQ-1BC99AA6** — Define architectural invariants and verify them automatically where practical.
- [ ] **USEQ-8141B679** — Assess quality-attribute trade-offs explicitly instead of optimizing one dimension in isolation.
- [ ] **USEQ-F20D0829** — Use prototypes, models, benchmarks, and experiments to reduce high-impact uncertainty.
- [ ] **USEQ-2E0281F9** — Review architecture before irreversible commitment and after material change.
- [ ] **USEQ-5EB6F214** — Measure architecture conformance without freezing legitimate evolution.
- [ ] **USEQ-A0F4735A** — Document intentional violations and migration plans.
- [ ] **USEQ-82F9AE01** — Retire superseded decisions while preserving history and rationale.
- [ ] **USEQ-1573BEF1** — Make architecture understandable to implementers, operators, security reviewers, product owners, and maintainers.
- [ ] **USEQ-36885CEE** — Treat architecture as a continuous decision process, not a one-time phase or document.

## Modularity, Cohesion, Coupling, and Boundaries

_Consolidated from `quality standards/04-architecture-design/06-modularity-cohesion-coupling-and-boundaries.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-5DDCF0A7** — Partition the system around cohesive responsibilities and change drivers.
- [ ] **USEQ-23D2EB3B** — Give each module a clear purpose, owner, public contract, and prohibited dependencies.
- [ ] **USEQ-304665CA** — Keep internal representation and implementation details hidden behind stable boundaries.
- [ ] **USEQ-8381689B** — Minimize knowledge of remote internals, shared mutable state, temporal order, deployment topology, and incidental data shape.
- [ ] **USEQ-14246510** — Prefer explicit dependencies over ambient globals, service locators, hidden context, or action at a distance.
- [ ] **USEQ-0F67185A** — Direct dependencies toward more stable policies and abstractions where this reduces change propagation.
- [ ] **USEQ-6B114385** — Avoid cyclic dependencies or contain and document them when temporarily unavoidable.
- [ ] **USEQ-68249BC7** — Use dependency graphs and change-impact evidence to identify excessive coupling.
- [ ] **USEQ-2CD29CD6** — Ensure modules can be tested in isolation with realistic contracts.
- [ ] **USEQ-78091092** — Prevent shared utility modules from becoming unowned dumping grounds.
- [ ] **USEQ-0D9E38D0** — Keep boundary crossings observable and failure-aware.
- [ ] **USEQ-52DDF534** — Align consistency, authorization, privacy, and ownership boundaries intentionally.
- [ ] **USEQ-2A32526C** — Do not split components solely to imitate a fashionable architecture or maximize service count.
- [ ] **USEQ-88759DC9** — Do not combine unrelated responsibilities solely to avoid interface design.
- [ ] **USEQ-54ED85FF** — Refactor boundaries when recurring changes, incidents, or coordination costs show the current partition is wrong.

## SOLID Design Principles

_Consolidated from `quality standards/04-architecture-design/07-solid-design-principles.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-E0114593** — Give each unit one coherent reason to change at its chosen level of abstraction.
- [ ] **USEQ-FE7895A8** — Separate policies that change for different stakeholders, rates, or reasons.
- [ ] **USEQ-1CDB0D32** — Prefer extension through stable composition or explicit variation points when recurring change is evidenced.
- [ ] **USEQ-9FB8CD95** — Do not introduce speculative extension mechanisms without a credible variation need.
- [ ] **USEQ-41BCD681** — Ensure replacements preserve all required preconditions, postconditions, invariants, side effects, error semantics, timing, and security properties.
- [ ] **USEQ-09BDE112** — Do not strengthen accepted inputs or weaken promised outputs in substitutes.
- [ ] **USEQ-D55AB3F0** — Keep interfaces focused on cohesive client needs rather than provider convenience.
- [ ] **USEQ-53F17C85** — Prevent clients from depending on operations or data they do not use.
- [ ] **USEQ-377A5E6F** — Depend on stable contracts at policy boundaries when doing so reduces volatility and improves testability.
- [ ] **USEQ-5B6C66DD** — Keep dependency inversion from degenerating into unnecessary wrappers over stable trivial behavior.
- [ ] **USEQ-589E6801** — Make ownership and lifetime of injected dependencies explicit.
- [ ] **USEQ-6F19C350** — Test contract conformance across implementations and substitutes.
- [ ] **USEQ-F873BE85** — Review SOLID choices together with simplicity, performance, locality, and cognitive load.
- [ ] **USEQ-7CEC4E2E** — Use composition over inheritance when inheritance does not represent a true substitutable relationship.
- [ ] **USEQ-1CAE0CA5** — Refactor violations when they cause change amplification, fragile tests, unsafe substitution, or duplicated policy.

## DRY, KISS, YAGNI, and Simplicity

_Consolidated from `quality standards/04-architecture-design/08-dry-kiss-yagni-and-simplicity.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-63328DA4** — Eliminate duplicated authoritative knowledge that can diverge, not merely visually similar code.
- [ ] **USEQ-56BA4427** — Keep one clear source of truth for rules, schemas, constants, permissions, and state transitions where practical.
- [ ] **USEQ-3299DB14** — Allow local duplication when premature unification would couple unrelated change drivers.
- [ ] **USEQ-14D7C72C** — Prefer the simplest design that satisfies current requirements, quality attributes, and credible near-term evolution.
- [ ] **USEQ-F89D2CC5** — Measure simplicity by ease of understanding, change, testing, operation, and failure recovery—not line count alone.
- [ ] **USEQ-EC772E63** — Avoid frameworks, layers, services, abstractions, configuration, and indirection without demonstrated benefit.
- [ ] **USEQ-030DFE3C** — Do not implement optionality, scale, integrations, or policies unsupported by an evidenced requirement.
- [ ] **USEQ-6CC29A15** — Preserve reversible seams around high-uncertainty decisions without building the unused future behind them.
- [ ] **USEQ-B7421394** — Remove dead code, stale flags, unused fields, obsolete interfaces, and abandoned compatibility paths.
- [ ] **USEQ-74681E5C** — Choose the least powerful language, data format, permission, and mechanism that safely solves the problem.
- [ ] **USEQ-DC1822ED** — Make ordinary paths obvious and exceptional paths explicit.
- [ ] **USEQ-DEDCD75F** — Prefer boring, understood technology for undifferentiated problems unless evidence justifies novelty.
- [ ] **USEQ-0B5AA3F8** — Reevaluate abstractions after the third real variation rather than applying a mechanical occurrence count.
- [ ] **USEQ-8C36A70B** — Document intentional duplication or complexity where trade-offs make it safer.
- [ ] **USEQ-4A3D04B2** — Treat cognitive load and operational burden as first-class costs in design review.

## Reusability and Software Product Lines

_Consolidated from `quality standards/04-architecture-design/09-reusability-and-software-product-lines.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-D46D1FB1** — Reuse stable knowledge and capabilities only when consumers share genuinely compatible semantics and quality needs.
- [ ] **USEQ-0B6B90C2** — Define the reusable asset's scope, supported contexts, assumptions, constraints, ownership, and service level.
- [ ] **USEQ-AF935A66** — Separate common core behavior from explicit, bounded variation points.
- [ ] **USEQ-4A789DEC** — Avoid configuration that permits invalid or untested combinations.
- [ ] **USEQ-31292BBF** — Design contracts for consumers rather than exposing provider internals.
- [ ] **USEQ-18804590** — Version reusable assets and publish compatibility, deprecation, and migration policy.
- [ ] **USEQ-B88AF508** — Test every supported variant and representative combinations.
- [ ] **USEQ-07414832** — Provide reference examples, integration tests, and failure semantics.
- [ ] **USEQ-5CFEC70E** — Track consumers so security fixes and breaking changes can be coordinated.
- [ ] **USEQ-7F342BF9** — Avoid copy-and-paste forks that lose provenance and maintenance unless isolation is an intentional trade-off.
- [ ] **USEQ-76C0FA3F** — Avoid central shared platforms that become mandatory bottlenecks without consumer governance.
- [ ] **USEQ-ED0786E2** — Measure reuse by reduced total lifecycle cost and consistent quality, not raw adoption count.
- [ ] **USEQ-DEC303C0** — Allow teams to leave a shared asset through documented interfaces and data portability.
- [ ] **USEQ-D45B9C85** — Fund maintenance, support, security, documentation, and evolution of shared assets.
- [ ] **USEQ-B3DBC424** — Retire variants and consumers that no longer justify their complexity.

## State Management

_Consolidated from `quality standards/04-architecture-design/10-state-management.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-26BE905D** — Identify every stateful element and assign an authoritative owner.
- [ ] **USEQ-B98CE1BB** — Distinguish domain state, session state, presentation state, configuration, cache, derived state, and transient workflow state.
- [ ] **USEQ-75AA9FEC** — Define state schema, invariants, legal transitions, initiating actors, guards, side effects, and terminal states.
- [ ] **USEQ-AE50A4B1** — Maintain one authoritative source for each fact and make replicas or projections explicitly derived.
- [ ] **USEQ-C64C66B8** — Prevent impossible, ambiguous, partially committed, or silently divergent state.
- [ ] **USEQ-A8E226BC** — Define consistency, durability, freshness, ordering, and conflict-resolution requirements per state.
- [ ] **USEQ-3044BFD8** — Make transitions atomic where partial effects are unacceptable.
- [ ] **USEQ-1E57B9D5** — Use idempotency, versioning, compare-and-set, or equivalent controls for repeated and concurrent updates.
- [ ] **USEQ-C58A7755** — Define initialization, hydration, synchronization, invalidation, expiry, persistence, migration, and reset behavior.
- [ ] **USEQ-0CC0C3BA** — Preserve user work across retries, navigation, process restart, failover, and offline recovery where required.
- [ ] **USEQ-48B58C2F** — Avoid hidden mutable global state and uncontrolled two-way synchronization.
- [ ] **USEQ-734F3AB7** — Ensure authorization and tenant context are part of security-sensitive state decisions.
- [ ] **USEQ-1AA2772B** — Make state changes auditable and diagnosable without exposing sensitive values.
- [ ] **USEQ-229E45DB** — Test transition coverage, illegal transitions, races, stale reads, duplicate events, and recovery.
- [ ] **USEQ-6EEAAF05** — Provide repair and reconciliation for state divergence that cannot be eliminated.

## Concurrency, Consistency, and Transactions

_Consolidated from `quality standards/04-architecture-design/11-concurrency-consistency-and-transactions.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-F0C80A36** — Identify shared resources, concurrent actors, ordering dependencies, contention points, and critical invariants.
- [ ] **USEQ-7AA2FAF0** — State the required consistency and isolation semantics rather than relying on platform defaults.
- [ ] **USEQ-B8D912A4** — Keep transaction boundaries aligned with business invariants.
- [ ] **USEQ-0F565CB8** — Minimize transaction duration and avoid interactive or external work inside locks where possible.
- [ ] **USEQ-38CF7C87** — Prevent lost updates, dirty reads, nonrepeatable decisions, duplicate effects, overselling, double spending, and stale authorization.
- [ ] **USEQ-F35EF8BD** — Use optimistic or pessimistic control deliberately based on conflict rate and failure cost.
- [ ] **USEQ-EEE5F2CD** — Make retried operations idempotent or detect duplicate execution.
- [ ] **USEQ-C5529D7D** — Define ordering and conflict-resolution behavior for asynchronous and replicated updates.
- [ ] **USEQ-5EE1E09A** — Avoid distributed transactions unless their guarantees and failure modes are required and understood.
- [ ] **USEQ-C19A4201** — Use sagas or compensating actions only when compensation semantics are valid and observable.
- [ ] **USEQ-CF270649** — Prevent deadlock, livelock, starvation, priority inversion, and unbounded contention.
- [ ] **USEQ-EB0270A6** — Bound queues, locks, semaphores, pools, and parallelism.
- [ ] **USEQ-55723B56** — Propagate cancellation, deadlines, and backpressure across concurrent work.
- [ ] **USEQ-619FC931** — Test with deterministic scheduling, stress, race detection, fault injection, and invariant checking as appropriate.
- [ ] **USEQ-515CE6F8** — Reconcile and repair after partial failure without concealing lost or duplicated work.

## Distributed and Event-Driven Systems

_Consolidated from `quality standards/04-architecture-design/12-distributed-and-event-driven-systems.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-AF6746E0** — Treat every remote call as fallible, delayed, duplicated, reordered, throttled, or unavailable.
- [ ] **USEQ-FFA324F3** — Define timeouts, deadlines, retry eligibility, backoff, jitter, circuit breaking, and fallback per dependency.
- [ ] **USEQ-334E0DEF** — Prevent retry amplification across layers.
- [ ] **USEQ-3D388A4E** — Define message delivery, ordering, partitioning, retention, replay, and deduplication semantics.
- [ ] **USEQ-69C4BB2C** — Make consumers idempotent when duplicate delivery is possible.
- [ ] **USEQ-A06A1E8F** — Use immutable event facts and version schemas compatibly.
- [ ] **USEQ-903E625E** — Do not expose internal database change as a public event contract without an intentional compatibility layer.
- [ ] **USEQ-1459B13E** — Use outbox, inbox, transactional messaging, or reconciliation where state and event consistency matters.
- [ ] **USEQ-8B7F37FF** — Bound queues and define overload, dead-letter, poison-message, and replay behavior.
- [ ] **USEQ-5D391FD0** — Preserve tenant, identity, trace, privacy, and authorization context across asynchronous boundaries.
- [ ] **USEQ-4CE2BAE8** — Make causality and correlation observable without assuming a single global clock.
- [ ] **USEQ-AC6ABDAD** — Design for network partitions and independent service evolution.
- [ ] **USEQ-4D924CD6** — Avoid synchronous dependency chains that create fragile availability coupling.
- [ ] **USEQ-AADFF352** — Test partial failure, stale replicas, clock skew, split-brain, failover, rebalancing, and consumer lag.
- [ ] **USEQ-85E0AA6D** — Provide reconciliation capable of detecting missing, duplicated, delayed, or contradictory effects.

## API and Integration Architecture

_Consolidated from `quality standards/04-architecture-design/13-api-and-integration-architecture.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-C0D91FC2** — Design interfaces from consumer tasks and domain contracts rather than internal storage models.
- [ ] **USEQ-8A67B399** — Define ownership, intended consumers, trust level, availability, latency, consistency, and support expectations.
- [ ] **USEQ-7414620B** — Use explicit schemas for requests, responses, events, errors, and compatibility rules.
- [ ] **USEQ-EF1356E3** — Choose synchronous, asynchronous, batch, stream, or file integration according to semantics and failure tolerance.
- [ ] **USEQ-90712C5E** — Keep transport concerns separate from domain behavior where practical.
- [ ] **USEQ-1D682327** — Make authentication, authorization, tenant, privacy, and data-classification requirements part of the contract.
- [ ] **USEQ-31636598** — Use stable identifiers and avoid leaking implementation-specific keys without need.
- [ ] **USEQ-3A484AC5** — Define pagination, filtering, sorting, concurrency, idempotency, rate, quota, and bulk semantics.
- [ ] **USEQ-FF5EBA4B** — Use consistent error taxonomy and distinguish retryable, terminal, conflict, validation, authorization, and dependency failures.
- [ ] **USEQ-2D10E138** — Avoid chatty interfaces and cyclic service dependencies.
- [ ] **USEQ-41F02F93** — Support independent deployment through consumer-driven contracts and compatible evolution.
- [ ] **USEQ-D5680902** — Define deprecation, discovery, change notification, migration, and retirement.
- [ ] **USEQ-1B7657B0** — Make external calls observable and attributable to user and business journeys.
- [ ] **USEQ-CCDE0B28** — Provide reconciliation for integrations that can lose or duplicate effects.
- [ ] **USEQ-FA82C923** — Evaluate whether an interface should exist before creating another permanent compatibility obligation.

## Evolvability, Technical Debt, and Changeability

_Consolidated from `quality standards/04-architecture-design/14-evolvability-technical-debt-and-changeability.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-02B1B35B** — Identify expected change dimensions and isolate the most volatile policies from stable mechanisms.
- [ ] **USEQ-E2513FAC** — Prefer reversible decisions where uncertainty is high.
- [ ] **USEQ-AFAA0716** — Make schemas, protocols, configuration, data, and deployment capable of compatible staged evolution.
- [ ] **USEQ-A720D7E2** — Design migrations before introducing states that cannot be changed safely.
- [ ] **USEQ-44AED881** — Track technical debt as a specific obligation with cause, consequence, interest, owner, and remediation trigger.
- [ ] **USEQ-F8D2AC51** — Distinguish deliberate short-term debt from accidental complexity, defects, missing controls, and obsolete technology.
- [ ] **USEQ-1EBB28D7** — Do not use technical debt as a label to hide security, privacy, accessibility, or reliability nonconformance.
- [ ] **USEQ-12D1C014** — Measure change lead time, regression rate, dependency impact, coordination cost, and operational burden.
- [ ] **USEQ-F8A23259** — Use automated architecture and compatibility checks for critical invariants.
- [ ] **USEQ-DDA66628** — Keep replacement seams around volatile external providers and end-of-life components.
- [ ] **USEQ-F8700B34** — Remove compatibility layers after consumers migrate.
- [ ] **USEQ-096571F7** — Reserve capacity for modernization and debt reduction based on evidence of risk and cost.
- [ ] **USEQ-3D11DC31** — Validate that refactoring changes structure without silently changing behavior.
- [ ] **USEQ-49148D13** — Use strangler, parallel run, shadowing, dual read/write, or other migration patterns only with reconciliation and exit criteria.
- [ ] **USEQ-06593DDA** — Review whether the architecture still fits scale, team topology, risk, and product strategy.

## Interoperability, Portability, and Compatibility

_Consolidated from `quality standards/04-architecture-design/15-interoperability-portability-and-compatibility.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-EA5D38F6** — Define supported environments, clients, protocols, formats, versions, locales, and integration partners.
- [ ] **USEQ-02149480** — Use open, stable, documented standards when they satisfy requirements.
- [ ] **USEQ-F91ACE3D** — Document intentional proprietary dependencies and the value that justifies them.
- [ ] **USEQ-D16FB36D** — Separate domain data and policy from provider-specific control planes where practical.
- [ ] **USEQ-8300060F** — Maintain export and import paths for data, configuration, identities, and audit history where continuity requires them.
- [ ] **USEQ-F46E5E03** — Define backward, forward, source, binary, schema, behavioral, and operational compatibility as applicable.
- [ ] **USEQ-F1600F03** — Test mixed-version operation and rolling upgrades where versions coexist.
- [ ] **USEQ-236C9C7F** — Use capability negotiation rather than unreliable version inference when appropriate.
- [ ] **USEQ-74E9981A** — Avoid undocumented behavior and implementation quirks as required contracts.
- [ ] **USEQ-B7411A36** — Preserve semantic meaning across serialization, precision, time zones, encodings, and locales.
- [ ] **USEQ-796DD1B7** — Make migration tooling resumable, observable, reversible or safely roll-forward, and reconcilable.
- [ ] **USEQ-7E1C1315** — Validate backups and recovery in the intended alternate environment where portability is a continuity control.
- [ ] **USEQ-4207A34E** — Track end-of-support dates for platforms and protocols.
- [ ] **USEQ-B923EBAE** — Do not promise portability that has never been tested.
- [ ] **USEQ-32F43E0B** — Review lock-in at major procurement and architecture decisions, including exit cost and time.

## Sustainable Software and Resource Efficiency

_Consolidated from `quality standards/04-architecture-design/16-sustainable-software-and-resource-efficiency.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-551A99E3** — Define sustainability objectives alongside performance, reliability, cost, and user outcomes.
- [ ] **USEQ-25AE787D** — Measure energy, carbon, compute, memory, storage, network transfer, and hardware utilization where material.
- [ ] **USEQ-CAB0969B** — Use a transparent measurement boundary, time period, allocation method, data source, and uncertainty statement.
- [ ] **USEQ-BB0D2BE5** — Reduce unnecessary computation, polling, transfer, storage, duplication, and retention.
- [ ] **USEQ-E75C31D5** — Prefer efficient algorithms and data structures where end-to-end impact is material.
- [ ] **USEQ-2C908E41** — Scale resources to actual demand without degrading required resilience or latency.
- [ ] **USEQ-5961E742** — Use batching, caching, compression, scheduling, and locality only when correctness and freshness remain satisfied.
- [ ] **USEQ-7BF9F781** — Extend useful system and device lifetime by avoiding needless client requirements and upgrade pressure.
- [ ] **USEQ-A0321258** — Consider embodied impact and hardware turnover, not only runtime energy.
- [ ] **USEQ-96E7974E** — Choose regions, schedules, and providers with environmental impact in mind where business and resilience constraints permit.
- [ ] **USEQ-AE0A9E6E** — Avoid shifting emissions or cost to users, suppliers, or other lifecycle stages without accounting for it.
- [ ] **USEQ-E9DA33B5** — Make measurement and optimization privacy-preserving.
- [ ] **USEQ-C30BCEBE** — Prevent efficiency targets from reducing accessibility, security, reliability, or user control.
- [ ] **USEQ-D0C9688D** — Include environmental requirements in supplier and architecture decisions.
- [ ] **USEQ-2FBFB7E2** — Review regressions and publish internally auditable progress against baselines.

## Standards and source references

- [ISO/IEC/IEEE 42010:2022 — Architecture description](https://www.iso.org/standard/74393.html)
- [NIST SP 800-218 v1.1 — Secure Software Development Framework](https://csrc.nist.gov/pubs/sp/800/218/final)
- [OWASP Application Security Verification Standard 5.0.0](https://owasp.org/www-project-application-security-verification-standard/)
- [ISO/IEC 25010:2023 — Product quality model](https://www.iso.org/standard/78176.html)
- [ISO/IEC 25019:2023 — Quality-in-use model](https://www.iso.org/standard/78177.html)
- [ISO/IEC/IEEE 15939:2017 — Measurement process](https://www.iso.org/standard/71197.html)
- [Google SRE Workbook — Implementing SLOs](https://sre.google/workbook/implementing-slos/)
- [ISO 22301:2019 — Business continuity management systems](https://www.iso.org/standard/75106.html)
- [ISO/IEC/IEEE 12207:2026 — Software life cycle processes](https://www.iso.org/standard/90219.html)
- [ISO/IEC 5055:2021 — Automated source code quality measures](https://www.iso.org/standard/80623.html)
- [ISO/IEC 20246:2017 — Work product reviews](https://www.iso.org/standard/67597.html)
- [ISO/IEC/IEEE 29148:2018 — Requirements engineering](https://www.iso.org/standard/72089.html)
- [ISO/IEC/IEEE 29119-4:2021 — Test techniques](https://www.iso.org/standard/79430.html)
- [RFC 9110 — HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110)
- [OWASP API Security Top 10 — 2023](https://owasp.org/API-Security/)
- [RFC 9457 — Problem Details for HTTP APIs](https://www.rfc-editor.org/rfc/rfc9457)
- [ISO/IEC 21031:2024 — Software Carbon Intensity](https://www.iso.org/standard/86612.html)

---

[Previous phase](03-user-experience-web-and-content.md) · [Next: Phase 5: Code quality and implementation](05-code-quality-and-implementation.md)
