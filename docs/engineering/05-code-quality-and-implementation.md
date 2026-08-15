# Code quality and implementation

_Phase 5 of 16 in the [complete engineering review](00-overview.md)._

Correctness, readability, contracts, errors, resources, concurrency, dependencies, review, testability, and maintainability.

Assess every control as **Pass**, **Fail**, **Blocked**, or **Not Applicable**. A Pass needs current evidence for the exact product, revision, artifact, configuration, environment, and data state under review.

## Functional Correctness and Business Logic

_Consolidated from `quality standards/05-code-implementation/01-functional-correctness-and-business-logic.md`; 18 non-duplicative controls._

### Universal controls

- [ ] **USEQ-1A47E1A7** — Every critical journey passes positive, negative, boundary, invalid-input, empty-state, maximum-size, and high-volume tests.
- [ ] **USEQ-6CEC2D29** — Duplicate submissions, refreshes, retries, browser navigation, and network interruption do not create unintended duplicate effects or corrupt state.
- [ ] **USEQ-CA554C6A** — Race conditions are tested around balances, inventory, permissions, limits, quotas, and state transitions.
- [ ] **USEQ-C2C592DA** — Invalid state-transition ordering is rejected.
- [ ] **USEQ-FBD441AF** — Partial failures cannot leave ambiguous or silently inconsistent results.
- [ ] **USEQ-4C913037** — Eventually consistent operations expose correct pending, completed, and failed states.
- [ ] **USEQ-F667636D** — Reconciliation detects missing, duplicated, and mismatched records.
- [ ] **USEQ-24D602A2** — Money uses correct decimal, currency, tax, exchange-rate, and rounding behavior.
- [ ] **USEQ-814D36E1** — Unicode, normalization, collation, case behavior, and confusable characters do not create duplicates or bypasses.
- [ ] **USEQ-55D4825F** — Limits and quotas are enforced at a trusted server-side boundary.
- [ ] **USEQ-5FF60A0C** — Approval and multi-person authorization workflows cannot be bypassed.
- [ ] **USEQ-4DE12FD4** — Suspension and termination affect all relevant sessions, tokens, workers, integrations, and cached permissions.
- [ ] **USEQ-D0B81370** — Long-running operations resume safely or fail clearly.
- [ ] **USEQ-2136C6A9** — Scheduled jobs handle missed runs, overlap, retry, and duplicate execution.
- [ ] **USEQ-E608D673** — Background jobs preserve correct user, authorization, and tenant context.
- [ ] **USEQ-04FA8EC6** — Dead-lettered work can be inspected and replayed safely.
- [ ] **USEQ-D138E3DB** — Support tools cannot create unsupported or impossible states.
- [ ] **USEQ-A0B1E367** — Error messages provide a safe, actionable next step without exposing sensitive details.

## Input Validation, Output Encoding, and Safe Processing

_Consolidated from `quality standards/05-code-implementation/02-input-validation-encoding-and-safe-processing.md`; 14 non-duplicative controls._

### Universal controls

- [ ] **USEQ-948F079B** — All untrusted input sources are identified.
- [ ] **USEQ-7C5F1359** — Input is validated for type, syntax, range, length, structure, and allowed values.
- [ ] **USEQ-4E8F864C** — Canonicalization and normalization occur before security-sensitive comparison.
- [ ] **USEQ-EC27C7B3** — Allowlists are used where a finite valid set exists.
- [ ] **USEQ-3A437981** — Operating-system commands are not constructed from untrusted input.
- [ ] **USEQ-4E590049** — Directory, query, search, template, expression, mail, logging, and interpreter contexts use safe APIs.
- [ ] **USEQ-D96C2840** — Rich text or markup is sanitized with an appropriate maintained sanitizer.
- [ ] **USEQ-1A1943A0** — Unsafe deserialization is avoided; expected types and structures are bounded.
- [ ] **USEQ-94B4DD17** — XML processing disables unnecessary external entities and resources.
- [ ] **USEQ-F5787145** — Path traversal, archive traversal, symbolic-link abuse, and filesystem races are prevented or controlled.
- [ ] **USEQ-7008ECFF** — Server-side request-forgery protections constrain protocol, destination, redirects, DNS resolution, and metadata endpoints.
- [ ] **USEQ-032834EB** — Header, response-splitting, log, CSV/formula, and email injection are addressed.
- [ ] **USEQ-B50A09FA** — Object and property injection through automatic binding or merging is addressed.
- [ ] **USEQ-8E9BA9B9** — Unicode confusables and normalization cannot bypass validation or identity checks.

## Maintainability and Long-Term Operability

_Consolidated from `quality standards/05-code-implementation/03-maintainability-and-long-term-operability.md`; 14 non-duplicative controls._

### Universal controls

- [ ] **USEQ-FDCA6C71** — A current operating entry point or README exists.
- [ ] **USEQ-F564C015** — Architecture, data flow, dependencies, deployment, recovery, and local setup are documented.
- [ ] **USEQ-AB553CE4** — Build, test, migration, deployment, and troubleshooting commands are documented.
- [ ] **USEQ-F1784541** — Coding, review, testing, security, release, and support conventions are documented.
- [ ] **USEQ-7CBEA80E** — Critical code and systems have appropriate maintainers.
- [ ] **USEQ-3FF4CC5B** — Complexity, coupling, and operational burden are reviewed.
- [ ] **USEQ-B6E04832** — Dead code, obsolete paths, stale flags, and unsupported compatibility layers are removed.
- [ ] **USEQ-6EF014F9** — Runtime, platform, protocol, certificate, and dependency end-of-support dates are tracked.
- [ ] **USEQ-2B9BF81E** — Technical debt affecting reliability, security, privacy, accessibility, or operability is recorded.
- [ ] **USEQ-BBB6533E** — API, data, feature, and client deprecation policies exist.
- [ ] **USEQ-FA1B43F8** — Support personnel are trained on released behavior.
- [ ] **USEQ-967C2427** — Known issues and limitations are communicated accurately.
- [ ] **USEQ-F534424C** — Maintenance windows and communications are defined.
- [ ] **USEQ-51159B9E** — Decommissioning, end-of-service export, retention, archival, and secure deletion have documented plans.

## Universal Code Quality

_Consolidated from `quality standards/05-code-implementation/04-universal-code-quality.md`; 22 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-8FA820FA** — Make code behavior satisfy explicit requirements and preserve documented invariants.
- [ ] **USEQ-6992FD80** — Choose representations, algorithms, and control flow that make correctness apparent.
- [ ] **USEQ-5D133AD0** — Keep units small enough to understand while preserving cohesive behavior.
- [ ] **USEQ-D34491E9** — Use clear names that reflect domain meaning, units, ownership, and side effects.
- [ ] **USEQ-C7003C4F** — Prefer explicit data flow, dependencies, state transitions, and error paths.
- [ ] **USEQ-3AB6A8DF** — Avoid surprising implicit conversion, mutation, global state, order dependence, and hidden I/O.
- [ ] **USEQ-AE80C3B2** — Validate untrusted input at trust boundaries and preserve validated types internally.
- [ ] **USEQ-1248940A** — Keep authorization, privacy, safety, and business rules in enforceable trusted locations.
- [ ] **USEQ-166192E3** — Handle every resource, transaction, lock, stream, handle, and subscription through a defined lifecycle.
- [ ] **USEQ-1D7ED71D** — Use abstractions only when they reduce total understanding and change cost.
- [ ] **USEQ-0CC65394** — Keep comments focused on intent, constraints, trade-offs, and non-obvious reasons rather than restating syntax.
- [ ] **USEQ-B6439E9E** — Delete unreachable, obsolete, duplicated, debug, and commented-out code.
- [ ] **USEQ-660F0A5E** — Treat warnings, static-analysis findings, and undefined behavior according to an explicit policy.
- [ ] **USEQ-2C162579** — Make exceptional and degraded behavior deliberate rather than accidental fall-through.
- [ ] **USEQ-5BA83821** — Use deterministic behavior where nondeterminism is not required.
- [ ] **USEQ-476447F8** — Avoid unbounded recursion, allocation, concurrency, retries, queues, and input size.
- [ ] **USEQ-B80088F6** — Make important behavior observable without logging secrets or sensitive data.
- [ ] **USEQ-509CC760** — Keep code reviewable through focused changes and stable formatting.
- [ ] **USEQ-F5EA5C82** — Refactor when structure obscures correctness, but verify behavior before and after change.
- [ ] **USEQ-22CB2FA4** — Do not optimize without evidence; do not ignore known high-impact inefficiency.

### Category no-go conditions

- [ ] **USEQ-E95604AB** — Undefined or implementation-dependent behavior can affect critical correctness.
- [ ] **USEQ-E66B032D** — Critical logic cannot be explained, reviewed, or tested by another qualified engineer.

## Readability, Naming, and Style

_Consolidated from `quality standards/05-code-implementation/05-readability-naming-and-style.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-DAF77C8F** — Use a consistent machine-enforced style so review focuses on behavior and design.
- [ ] **USEQ-9AF5256A** — Name concepts using the ubiquitous domain vocabulary agreed with stakeholders.
- [ ] **USEQ-51F9CDC3** — Name booleans and predicates so truth meaning is unambiguous.
- [ ] **USEQ-73138957** — Include units, time basis, currency, encoding, and coordinate system where confusion is plausible.
- [ ] **USEQ-63B6EFE7** — Distinguish identifiers, display labels, secrets, tokens, user input, and trusted values through naming and types.
- [ ] **USEQ-0A590DC8** — Avoid abbreviations unless they are standard and unambiguous in the relevant domain.
- [ ] **USEQ-AA3301B8** — Avoid misleading legacy names when behavior has changed.
- [ ] **USEQ-DC209C03** — Keep related logic physically and conceptually close.
- [ ] **USEQ-A69D553C** — Order code to reveal the normal flow before exceptional detail where practical.
- [ ] **USEQ-403F240B** — Use early returns or structured decomposition to limit unnecessary nesting.
- [ ] **USEQ-1B3ECB70** — Avoid dense expressions that hide evaluation order, side effects, or error handling.
- [ ] **USEQ-4697DF3E** — Replace unexplained literals with named domain concepts when meaning is not obvious.
- [ ] **USEQ-E621FB99** — Use comments to explain why constraints exist, including external quirks and security assumptions.
- [ ] **USEQ-1729029B** — Keep documentation synchronized with behavior or remove misleading documentation.
- [ ] **USEQ-084F8B56** — Write error messages that identify context and action without leaking sensitive information.
- [ ] **USEQ-79C6A5E7** — Use consistent terminology across code, schemas, APIs, UI, logs, tests, and documentation.
- [ ] **USEQ-9914FF54** — Make public contracts readable without requiring knowledge of internal implementation.
- [ ] **USEQ-F8A9898C** — Use examples for complex contracts and edge semantics.
- [ ] **USEQ-CD3455A4** — Prefer clarity over cleverness, novelty, terseness, or stylistic performance.
- [ ] **USEQ-E67B4E95** — Review readability from the perspective of a competent maintainer unfamiliar with the change.

## Abstractions, Interfaces, and Contracts

_Consolidated from `quality standards/05-code-implementation/06-abstractions-interfaces-and-contracts.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-F39FA898** — Abstract stable domain meaning rather than superficial syntax similarity.
- [ ] **USEQ-C484D3B0** — Define each public contract's purpose, inputs, outputs, errors, side effects, timing, ordering, ownership, and security expectations.
- [ ] **USEQ-4EDDA552** — State preconditions, postconditions, invariants, and exceptional behavior.
- [ ] **USEQ-84D5B294** — Keep contracts smaller than implementations and free from unnecessary implementation detail.
- [ ] **USEQ-85DDD1B8** — Hide representation so it can evolve without changing consumers.
- [ ] **USEQ-6F3283A1** — Use types or schemas that make invalid states difficult or impossible to represent.
- [ ] **USEQ-51A3C767** — Avoid boolean flags and parameter combinations whose semantics are unclear.
- [ ] **USEQ-82F98FD7** — Prefer cohesive parameter objects or explicit variants when arguments must change together.
- [ ] **USEQ-04B47B17** — Prevent null, missing, sentinel, and optional values from carrying multiple ambiguous meanings.
- [ ] **USEQ-66DD52CA** — Make mutability and ownership explicit.
- [ ] **USEQ-781C734B** — Do not expose writable internal collections or mutable references unintentionally.
- [ ] **USEQ-48E6F214** — Keep interface segregation aligned to consumer needs.
- [ ] **USEQ-BF0669D4** — Provide stable error categories rather than forcing consumers to parse text.
- [ ] **USEQ-1C9AF441** — Define compatibility and versioning before publishing a contract outside one change boundary.
- [ ] **USEQ-12574552** — Test contracts independently from implementations and across every substitute.
- [ ] **USEQ-47C5760F** — Avoid wrapping every dependency when no policy, isolation, testability, or evolution benefit exists.
- [ ] **USEQ-D0944F6C** — Avoid generic abstractions that erase meaningful constraints.
- [ ] **USEQ-73588D8A** — Document performance and resource characteristics when consumers can depend on them.
- [ ] **USEQ-38607061** — Define cancellation, timeout, retry, and idempotency semantics for remote or long-running contracts.
- [ ] **USEQ-A1CE1F63** — Retire obsolete interfaces through monitored migration rather than indefinite parallel support.

## Error Handling and Defensive Programming

_Consolidated from `quality standards/05-code-implementation/07-error-handling-and-defensive-programming.md`; 22 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-0903254E** — Define an error taxonomy that distinguishes validation, conflict, authorization, not-found, dependency, resource, timeout, cancellation, integrity, and internal failures.
- [ ] **USEQ-DA42055F** — Detect errors as close as practical to their source while handling them at the layer that has recovery context.
- [ ] **USEQ-CD710602** — Do not ignore, overwrite, convert to success, or silently log-and-continue after an error that invalidates the operation.
- [ ] **USEQ-50C6FF56** — Preserve causal context while redacting secrets and sensitive data.
- [ ] **USEQ-208AD931** — Use structured error values rather than parsing human-readable strings.
- [ ] **USEQ-8F5AF188** — Translate errors at boundaries without losing retryability, severity, and causal identity.
- [ ] **USEQ-50874EC8** — Fail closed for authorization, identity, privacy, safety, and integrity controls.
- [ ] **USEQ-F69CB29C** — Avoid broad catch-all handling that masks programming defects or corrupt state.
- [ ] **USEQ-60CEF31B** — Ensure cleanup, rollback, lock release, and cancellation occur on every exit path.
- [ ] **USEQ-13405F0D** — Make retries explicit, bounded, delayed, and safe for the operation.
- [ ] **USEQ-CF15A8BE** — Separate user-facing guidance from diagnostic detail.
- [ ] **USEQ-0E8AFB14** — Avoid exposing internal paths, queries, stack traces, keys, topology, or tenant information.
- [ ] **USEQ-F07D674E** — Treat assertion and invariant failures as defects requiring investigation, not routine user errors.
- [ ] **USEQ-3B20FCE0** — Validate external assumptions and provider responses before use.
- [ ] **USEQ-0C9AFBF3** — Use defensive copies or immutable values when callers could corrupt shared state.
- [ ] **USEQ-BACFAF4F** — Bound input, recursion, allocation, loops, retries, and concurrency against resource exhaustion.
- [ ] **USEQ-AA460D26** — Define safe fallback and degraded behavior before failures occur.
- [ ] **USEQ-F2F9031A** — Make partial success explicit and reconcilable.
- [ ] **USEQ-A45DF449** — Test error injection at every material dependency and state transition.
- [ ] **USEQ-718A1C3D** — Monitor unexpected error classes and repeated recoveries as signals of systemic defects.

### Category no-go conditions

- [ ] **USEQ-F5CCDCCC** — A critical operation can report success after losing, duplicating, corrupting, or failing to authorize its effect.
- [ ] **USEQ-A1DBD310** — Errors can leave security or business invariants silently violated.

## Resource Lifecycle, Memory, and Cleanup

_Consolidated from `quality standards/05-code-implementation/08-resource-lifecycle-memory-and-cleanup.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-A3BE3A81** — Identify every finite resource used by the code, including memory, files, sockets, connections, threads, tasks, locks, handles, subscriptions, and temporary storage.
- [ ] **USEQ-0836A8FC** — Make resource ownership and transfer explicit.
- [ ] **USEQ-00041688** — Acquire resources as late as practical and release them as early as safe.
- [ ] **USEQ-FEF94073** — Use structured lifetime mechanisms so cleanup occurs on success, error, cancellation, and timeout.
- [ ] **USEQ-1AF33210** — Prevent double release, use-after-release, dangling references, and ownership ambiguity.
- [ ] **USEQ-12D8A914** — Bound pools, queues, buffers, caches, batches, and concurrency.
- [ ] **USEQ-7233DE8D** — Set timeouts and cancellation for operations that can wait indefinitely.
- [ ] **USEQ-743D7427** — Avoid holding scarce resources across user input, network calls, or long computation without justification.
- [ ] **USEQ-4204983E** — Size pools from measured demand and downstream capacity rather than arbitrary defaults.
- [ ] **USEQ-C16F962C** — Prevent pool exhaustion from causing uncontrolled retry or cascading failure.
- [ ] **USEQ-353D14C8** — Clean temporary, orphaned, expired, and abandoned resources reliably.
- [ ] **USEQ-21DBD1D6** — Handle partial acquisition by releasing already-acquired resources.
- [ ] **USEQ-4A482940** — Avoid unbounded object retention through listeners, closures, caches, registries, or diagnostic context.
- [ ] **USEQ-3A5B2F4F** — Measure memory, allocation, handle, connection, and thread behavior under sustained load.
- [ ] **USEQ-3E450704** — Test resource exhaustion and recovery.
- [ ] **USEQ-EB0A5E6E** — Make spill, eviction, rejection, and load-shedding behavior explicit.
- [ ] **USEQ-E9F27341** — Protect shared resources from noisy neighbors.
- [ ] **USEQ-24A3BDC9** — Avoid performing cleanup that can delete another operation's or tenant's resources.
- [ ] **USEQ-13140E4F** — Make resource limits configurable, validated, documented, and observable.
- [ ] **USEQ-1A94D768** — Verify that shutdown drains, transfers, persists, or safely abandons in-flight work.

## Concurrent, Asynchronous, and Parallel Code

_Consolidated from `quality standards/05-code-implementation/09-concurrent-asynchronous-and-parallel-code.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-D1B25825** — Use concurrency only where it provides required responsiveness, throughput, isolation, or latency benefit.
- [ ] **USEQ-D464A4C5** — Define which state is shared, immutable, thread-confined, message-owned, or synchronized.
- [ ] **USEQ-0C811996** — Prefer ownership transfer and immutability over shared mutation.
- [ ] **USEQ-8D0CC00B** — Protect every shared invariant with a coherent synchronization strategy.
- [ ] **USEQ-EDE64E87** — Avoid blocking operations on event loops, cooperative schedulers, or latency-critical threads.
- [ ] **USEQ-6944D6B8** — Bound task creation, fan-out, parallelism, queues, and in-flight work.
- [ ] **USEQ-6BCF0D32** — Propagate cancellation, deadlines, tenant context, identity, trace context, and request scope correctly.
- [ ] **USEQ-9136B800** — Do not detach background work without ownership, error handling, lifecycle, and shutdown behavior.
- [ ] **USEQ-2854EA57** — Define ordering and completion semantics explicitly.
- [ ] **USEQ-38212874** — Prevent races between timeout, cancellation, completion, retry, and cleanup.
- [ ] **USEQ-1AC46AF9** — Avoid holding locks while invoking untrusted, remote, or reentrant code.
- [ ] **USEQ-B3FAEE44** — Use lock ordering or lock-free designs deliberately to prevent deadlock.
- [ ] **USEQ-A5FFFB7A** — Protect against starvation, livelock, priority inversion, and unfair scheduling.
- [ ] **USEQ-8EF6EF5A** — Make async errors observable and attributable to the initiating work.
- [ ] **USEQ-41B0B0C5** — Preserve idempotency when work can be retried or duplicated.
- [ ] **USEQ-B4C195BA** — Test under randomized scheduling, high contention, slow dependencies, cancellation, and shutdown.
- [ ] **USEQ-4C15F6BF** — Use race detectors, deterministic schedulers, model checking, or stress tools where available and appropriate.
- [ ] **USEQ-4D867D6C** — Document memory visibility and consistency assumptions where the execution model requires it.
- [ ] **USEQ-272FFC6D** — Treat nondeterministic test failures as defects rather than normal noise.
- [ ] **USEQ-95839798** — Verify that parallel execution does not violate rate, quota, ordering, or downstream capacity limits.

## Configuration and Feature-Flag Code Quality

_Consolidated from `quality standards/05-code-implementation/10-configuration-and-feature-flag-code-quality.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-2C7B2D29** — Define a schema, type, allowed range, default, sensitivity, owner, and description for every configuration value.
- [ ] **USEQ-88C70738** — Validate all configuration before serving work and fail safely on invalid critical values.
- [ ] **USEQ-3AE6779A** — Make precedence across defaults, files, environment, remote sources, and command-line input explicit.
- [ ] **USEQ-B0283141** — Avoid configuration values whose meaning changes by undocumented convention.
- [ ] **USEQ-EEC9584F** — Do not store secrets in ordinary configuration or expose them through diagnostics.
- [ ] **USEQ-26B1A7AC** — Separate environment identity from behavior and avoid environment-name conditionals scattered through code.
- [ ] **USEQ-C0A2F9FB** — Version configuration with the release or maintain equivalent auditable traceability.
- [ ] **USEQ-98F39A05** — Make dynamic configuration updates atomic, observable, authorized, and reversible.
- [ ] **USEQ-81C2964A** — Define stale-cache, unavailable-source, and partial-update behavior.
- [ ] **USEQ-F3EE8B38** — Treat flags as temporary lifecycle-managed code with owner, purpose, default, scope, creation date, and removal date.
- [ ] **USEQ-340326EC** — Test both flag states and material combinations before exposure.
- [ ] **USEQ-CC9C98B9** — Use safe defaults when flag evaluation fails.
- [ ] **USEQ-985ADD85** — Restrict and audit high-impact flag changes.
- [ ] **USEQ-27A8638E** — Avoid nested flags and flag interactions that create untestable state spaces.
- [ ] **USEQ-6D9F904A** — Remove completed rollout flags and dead branches promptly.
- [ ] **USEQ-7FF8A889** — Do not use flags as permanent authorization, pricing, compliance, or data-isolation mechanisms without appropriate governance.
- [ ] **USEQ-BB0F7387** — Ensure configuration changes cannot bypass validation, deployment controls, or segregation of duties.
- [ ] **USEQ-C5D04742** — Redact sensitive values from logs, UIs, metrics, and support tools.
- [ ] **USEQ-6A2615C9** — Provide configuration drift detection and reconstruction.
- [ ] **USEQ-555D7538** — Include configuration in incident, rollback, and reproducibility evidence.

## Dependency Selection and Hygiene

_Consolidated from `quality standards/05-code-implementation/11-dependency-selection-and-hygiene.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-D0D3A457** — Add a dependency only when its total lifecycle benefit exceeds implementing or avoiding the capability.
- [ ] **USEQ-AA278192** — Evaluate maintainer identity, release history, support, security practice, community health, documentation, and abandonment risk.
- [ ] **USEQ-184D1160** — Use trusted sources and verify package identity, namespace, signature, hash, or provenance as available.
- [ ] **USEQ-71972C15** — Pin versions or constraints tightly enough to prevent unreviewed substitution.
- [ ] **USEQ-2FA33CEA** — Record direct and transitive dependencies in a machine-readable inventory.
- [ ] **USEQ-D4E7F5C7** — Separate build, development, optional, and runtime dependencies.
- [ ] **USEQ-D49F7FE4** — Minimize privileged, networked, native, dynamically loaded, or executable dependencies.
- [ ] **USEQ-DD00A1EA** — Review license, notice, patent, redistribution, and source-disclosure obligations before use.
- [ ] **USEQ-3EC4C6D1** — Avoid multiple libraries providing the same capability without a justified migration plan.
- [ ] **USEQ-B894A653** — Wrap volatile or high-risk dependencies behind a policy boundary when replacement or containment is valuable.
- [ ] **USEQ-CB6E238D** — Do not rely on undocumented behavior or private interfaces.
- [ ] **USEQ-EA32C05E** — Test supported versions and upgrade paths.
- [ ] **USEQ-B09FE69D** — Automate updates while preserving review, test, provenance, and rollback gates.
- [ ] **USEQ-9A64DA9A** — Monitor newly disclosed vulnerabilities and compromised releases after deployment.
- [ ] **USEQ-6A7D79B1** — Prioritize remediation by exploitability, reachability, exposure, impact, and controls.
- [ ] **USEQ-02FECB7F** — Remove unused, abandoned, end-of-life, and superseded dependencies.
- [ ] **USEQ-55D2754E** — Maintain an emergency process to revoke, replace, rebuild, and redeploy a compromised dependency.
- [ ] **USEQ-F0B1DD9F** — Track local patches and forks with ownership and upstream strategy.
- [ ] **USEQ-4559EED5** — Prevent dependency confusion, typosquatting, namespace collision, and untrusted install scripts.
- [ ] **USEQ-696A7B67** — Verify the final artifact contains only intended dependency content.

## Refactoring and Legacy Code

_Consolidated from `quality standards/05-code-implementation/12-refactoring-and-legacy-code.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-9B50CC48** — Define the problem, desired quality improvement, boundaries, risks, and measurable completion criteria.
- [ ] **USEQ-E52095B1** — Establish characterization tests or other evidence for current required behavior before structural change.
- [ ] **USEQ-9D7728A6** — Distinguish required behavior from defects, accidental quirks, and obsolete compatibility.
- [ ] **USEQ-7B68CEDF** — Use small, reviewable, reversible steps with continuous verification.
- [ ] **USEQ-B9589547** — Avoid mixing broad behavior change, dependency upgrade, migration, and refactoring without a clear reason.
- [ ] **USEQ-48D1C437** — Preserve public contracts or provide a controlled migration.
- [ ] **USEQ-7E07C9BB** — Add seams around tightly coupled or untestable areas before changing internals.
- [ ] **USEQ-3AC09F62** — Use parallel run, shadow comparison, dual read, or reconciliation where behavior equivalence is difficult to prove.
- [ ] **USEQ-E9BADF4B** — Measure defect rate, change lead time, complexity, performance, and operational burden before and after.
- [ ] **USEQ-6CD08E3B** — Remove dead paths, compatibility shims, stale flags, and old data only after consumers have migrated.
- [ ] **USEQ-BED93E1D** — Keep data migrations resumable and auditable.
- [ ] **USEQ-57C309A4** — Do not translate code mechanically into a new language or framework without reconsidering domain boundaries and failure semantics.
- [ ] **USEQ-B3207962** — Prioritize high-interest debt that creates recurring incidents, vulnerability, delay, or inability to change.
- [ ] **USEQ-F763C13B** — Stop modernization that expands scope without reducing the target risk or cost.
- [ ] **USEQ-89E1A165** — Preserve historical decision context and document new constraints.
- [ ] **USEQ-2C1D6F3B** — Verify rollback or safe roll-forward at each migration stage.
- [ ] **USEQ-DA219C26** — Avoid permanent dual systems without an owner and exit date.
- [ ] **USEQ-DEE55F32** — Include operators, support, security, data, and users in modernization impact review.
- [ ] **USEQ-07A0BBCB** — Treat unexplained legacy behavior as an investigation target, not permission to break it.
- [ ] **USEQ-354B29AD** — Retire obsolete infrastructure, credentials, data, and suppliers after migration completes.

## Code Review and Work-Product Review

_Consolidated from `quality standards/05-code-implementation/13-code-review-and-work-product-review.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-04CD3D8D** — Define review objectives, required reviewers, independence, evidence, and approval rules according to risk.
- [ ] **USEQ-8B82476F** — Keep changes small and focused enough for reliable human review.
- [ ] **USEQ-39D96811** — Provide context, requirement links, design rationale, risk, test evidence, migration, rollout, and rollback information.
- [ ] **USEQ-8C7532C7** — Review behavior and consequences rather than formatting already enforced by tools.
- [ ] **USEQ-9E82B4E8** — Verify requirements, invariants, authorization, validation, error handling, concurrency, data, and operational behavior.
- [ ] **USEQ-561E2A9E** — Examine deleted and generated content as well as added lines.
- [ ] **USEQ-6FF0CECF** — Check whether tests can fail for the defect they claim to detect.
- [ ] **USEQ-3FEE7961** — Challenge assumptions, edge cases, misuse, partial failure, and recovery.
- [ ] **USEQ-5143F148** — Use specialists for security, privacy, accessibility, data, cryptography, performance, or safety-sensitive changes.
- [ ] **USEQ-02E7213F** — Prevent authors from approving their own material changes.
- [ ] **USEQ-235B0220** — Resolve review comments explicitly; do not dismiss findings through status changes without rationale.
- [ ] **USEQ-D61F1909** — Distinguish blocking defects, required follow-up, questions, and optional suggestions.
- [ ] **USEQ-2BA69F57** — Re-review materially changed code after approval.
- [ ] **USEQ-7DB8FB57** — Record exceptional bypasses and perform prompt retrospective review.
- [ ] **USEQ-42214EC2** — Use checklists as memory aids, not substitutes for understanding.
- [ ] **USEQ-683FB10E** — Avoid review overload, rubber-stamping, and excessive queues through ownership and change sizing.
- [ ] **USEQ-A87A04F3** — Measure escapes and review effectiveness rather than comment count or review speed alone.
- [ ] **USEQ-58E7D957** — Create automated rules from recurring objective review findings.
- [ ] **USEQ-7A2E33C3** — Maintain respectful, specific, evidence-based review communication.
- [ ] **USEQ-CA466120** — Verify the merged artifact and configuration still match the reviewed change.

## Testability

_Consolidated from `quality standards/05-code-implementation/14-testability.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-839B5A77** — Separate pure decision logic from I/O, time, randomness, process state, and external dependencies where practical.
- [ ] **USEQ-28C12031** — Inject or otherwise control clocks, randomness, identifiers, schedulers, and external boundaries.
- [ ] **USEQ-792699A0** — Expose behavior through stable public contracts rather than test-only access to internals.
- [ ] **USEQ-0880766A** — Make important state transitions and outcomes observable without exposing sensitive data.
- [ ] **USEQ-4FD13567** — Support deterministic setup, execution, teardown, and repeatable seeds.
- [ ] **USEQ-E384EC6B** — Avoid hidden globals, ambient singletons, static caches, and order-dependent initialization.
- [ ] **USEQ-0D13F05F** — Design dependency boundaries that can be replaced by faithful fakes, simulators, or contract-tested test doubles.
- [ ] **USEQ-1117639C** — Keep test doubles behaviorally constrained so they do not create false confidence.
- [ ] **USEQ-BAE66E14** — Provide safe test hooks only when production security and integrity cannot be weakened.
- [ ] **USEQ-2C8EA20F** — Make data creation concise, valid by default, and capable of expressing edge states.
- [ ] **USEQ-B3264D1A** — Allow failure, timeout, cancellation, retry, and partial-response injection at material boundaries.
- [ ] **USEQ-7E74EFCD** — Make asynchronous completion, queues, and eventual consistency waitable by state rather than arbitrary delay.
- [ ] **USEQ-7F30BB3A** — Provide stable identifiers and correlation for cross-component verification.
- [ ] **USEQ-10B21BF1** — Keep units cohesive enough that failures localize the cause.
- [ ] **USEQ-2EF78BB9** — Avoid nondeterministic concurrency and network dependence in lower-level tests.
- [ ] **USEQ-430CE5AC** — Design migrations and background jobs to support dry run, checkpoint, and reconciliation.
- [ ] **USEQ-12F0CBBD** — Expose health and diagnostic signals suitable for production-safe verification.
- [ ] **USEQ-AC43E131** — Ensure tests can distinguish absence of work from successful completion.
- [ ] **USEQ-7330C584** — Measure test setup complexity and time as design feedback.
- [ ] **USEQ-36AB1DA6** — Refactor production design when testing requires excessive private knowledge or brittle global setup.

## Reusability in Code

_Consolidated from `quality standards/05-code-implementation/15-reusability-in-code.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-B338A49B** — Extract reuse only when multiple consumers share stable semantics and change drivers.
- [ ] **USEQ-A31D6133** — Keep reusable units cohesive and independent of consumer-specific workflow.
- [ ] **USEQ-BD0179D8** — Define explicit inputs, outputs, errors, side effects, resource ownership, and extension points.
- [ ] **USEQ-DEA324DE** — Use configuration for genuine variation, not to combine incompatible responsibilities.
- [ ] **USEQ-AD85128A** — Provide safe defaults and prevent invalid configuration combinations.
- [ ] **USEQ-288CCD01** — Avoid leaking one consumer's terminology, permissions, tenant context, or data into the shared contract.
- [ ] **USEQ-457E0384** — Keep dependencies minimal and visible.
- [ ] **USEQ-25CA518A** — Version and document compatibility commitments.
- [ ] **USEQ-4B011C6E** — Provide examples for ordinary use and guidance for unsupported use.
- [ ] **USEQ-0845C21B** — Test each supported context and representative integration.
- [ ] **USEQ-BEC9299D** — Track consumers and communicate security fixes, deprecations, and breaking changes.
- [ ] **USEQ-B704594F** — Avoid global utility modules without ownership or domain boundaries.
- [ ] **USEQ-F32FF3B0** — Prefer local duplication over a shared abstraction that causes coordinated change across unrelated products.
- [ ] **USEQ-1E288863** — Measure whether reuse reduces total defects, effort, and inconsistency after maintenance cost.
- [ ] **USEQ-D15F66C8** — Provide a clear path to fork or replace when shared governance no longer serves consumers.
- [ ] **USEQ-41D3FA73** — Retire unused exports and variants.
- [ ] **USEQ-459360CC** — Keep reusable security and validation functions context-aware rather than universally permissive.
- [ ] **USEQ-EFD9E266** — Do not expose internal implementation merely to make reuse convenient.
- [ ] **USEQ-5E19F6A1** — Avoid inheritance or mixins that create hidden coupling and fragile override behavior.
- [ ] **USEQ-E2E29325** — Treat reuse as a product with ownership, support, documentation, and lifecycle.

## Universal Code Quality, Design, and Implementation Master Checklist

_Consolidated from `gap supplement/02-universal-code-quality-design-and-implementation.md`; 205 non-duplicative controls._

### Expanded gap-closure controls

#### How to apply engineering principles without dogma

- [ ] **USEQ-B3C67404** — Treat correctness, security, privacy, safety, accessibility, and recoverability as constraints; optimize elegance, reuse, and brevity only within those constraints.
- [ ] **USEQ-53A25470** — Use SOLID, DRY, KISS, YAGNI, separation of concerns, information hiding, composition, immutability, and related principles as context-sensitive heuristics rather than unconditional laws.
- [ ] **USEQ-60CC7E5A** — Document the quality attribute or change scenario that justifies a nontrivial abstraction, pattern, layer, cache, concurrency mechanism, framework, or generated solution.
- [ ] **USEQ-21B58F87** — Prefer the simplest design that satisfies current evidenced requirements and known near-term constraints without closing necessary evolution paths.
- [ ] **USEQ-E17453CC** — Do not invoke YAGNI to omit a known requirement, required control, migration path, operability need, or foreseeable high-impact failure treatment.
- [ ] **USEQ-41F7BE92** — Do not invoke DRY to merge coincidentally similar concepts that have different reasons to change, ownership, lifecycle, authorization, or semantics.
- [ ] **USEQ-0124C8F5** — Do not preserve duplication when repeated defects or synchronized changes show that one authoritative abstraction is warranted.
- [ ] **USEQ-B89A1447** — Do not invoke SOLID to create an excessive number of indirections, interfaces, containers, files, services, or extension points that increase total cognitive load.
- [ ] **USEQ-08D5F342** — Do not optimize for minimum line count; optimize for clear intent, correct behavior, safe change, and total lifecycle cost.
- [ ] **USEQ-9887E01F** — Choose consistency with the codebase when alternatives are equivalent, but change an established convention when evidence shows it causes defects or unacceptable cost.
- [ ] **USEQ-E0F101C2** — Make trade-offs explicit when a principle conflicts with latency, memory, safety, compatibility, auditability, or simplicity.
- [ ] **USEQ-D18F9478** — Review principles at system boundaries and data flows, not only within individual functions or classes.

#### Specifications, contracts, invariants, and correctness

- [ ] **USEQ-A191CD19** — Express externally observable behavior as testable preconditions, postconditions, invariants, state transitions, error outcomes, timing expectations, and side-effect guarantees.
- [ ] **USEQ-9072DA14** — Identify which invariants are enforced by types, schemas, constructors, validation, storage constraints, transactions, authorization, tests, monitoring, or manual process.
- [ ] **USEQ-76DF8B5A** — Represent invalid states as unconstructable or immediately rejected where practical.
- [ ] **USEQ-6DAA216C** — Keep domain rules independent from presentation, transport, persistence, and vendor details unless coupling is inherent to the requirement.
- [ ] **USEQ-B01D50C5** — Make units, coordinate systems, currencies, time zones, precision, encoding, normalization, and ordering explicit in names and contracts.
- [ ] **USEQ-076FF26C** — Specify behavior for zero, empty, missing, duplicate, stale, maximum, minimum, overflow, underflow, malformed, unauthorized, concurrent, canceled, and partially completed inputs.
- [ ] **USEQ-C6B6695E** — Define equality, identity, uniqueness, ordering, hashing, comparison, and canonicalization consistently for domain values.
- [ ] **USEQ-B39802DF** — Use one authoritative rule for each business invariant and ensure every entry point reaches it.
- [ ] **USEQ-F46D430E** — Ensure a successful result means the complete promised effect occurred or a clearly modeled pending state was created.
- [ ] **USEQ-3305DCE8** — Ensure failure cannot be mistaken for success through default values, swallowed errors, partial writes, stale caches, or ambiguous status.
- [ ] **USEQ-C3DFBAB2** — Verify critical calculations with an independent oracle, reconciliation, property, model, or second implementation where consequence warrants it.
- [ ] **USEQ-1E0653D0** — Make assumptions about external systems executable through contract tests, runtime validation, or monitored assertions.
- [ ] **USEQ-75B210AD** — Use assertions for internal impossible states and diagnostics, not as the sole validation of untrusted input or security policy.
- [ ] **USEQ-92DFE757** — Define deterministic conflict resolution when multiple valid updates, versions, or sources can disagree.
- [ ] **USEQ-44BCDEC5** — Preserve invariants across retries, replay, rollback, restoration, migration, mixed versions, and disaster recovery.

#### Boundaries, modularity, coupling, and ownership

- [ ] **USEQ-611FF9D1** — Align module boundaries with cohesive capabilities, data ownership, trust boundaries, scaling needs, and independent reasons to change.
- [ ] **USEQ-2500EBD9** — Keep dependency direction consistent with architectural policy and make forbidden dependencies mechanically detectable where practical.
- [ ] **USEQ-3E378FD4** — Prevent internal representation, storage schema, framework types, transport objects, and vendor models from leaking across boundaries without an intentional contract.
- [ ] **USEQ-D6F735CA** — Minimize the number of modules that can mutate a given state or enforce a given invariant.
- [ ] **USEQ-5F91F410** — Give every public module, service, package, component, data set, and interface an owner and lifecycle policy.
- [ ] **USEQ-FD78B783** — Keep boundary contracts smaller and more stable than internal implementations.
- [ ] **USEQ-EE97AF92** — Avoid cyclic dependencies; when a cycle is unavoidable, document the shared invariant and why extraction would be worse.
- [ ] **USEQ-65E57265** — Prefer explicit dependency injection or parameterization over hidden service location, ambient globals, import-time side effects, or implicit thread context.
- [ ] **USEQ-8265FFBB** — Separate policy from mechanism and pure decision logic from effectful execution where that improves verification and reuse.
- [ ] **USEQ-AC6193E0** — Co-locate behavior with the data and invariants it owns unless doing so creates a stronger boundary violation.
- [ ] **USEQ-197BE325** — Prevent shared utility areas from becoming unowned collections of unrelated behavior.
- [ ] **USEQ-F826255C** — Track fan-in, fan-out, change coupling, dependency depth, and cross-boundary change frequency for high-risk components.
- [ ] **USEQ-5F8E12FC** — Test boundaries independently and through representative integration paths.
- [ ] **USEQ-5B7AE895** — Keep failure containment boundaries aligned with ownership and operational recovery.

#### Abstraction, reuse, and duplication

- [ ] **USEQ-DCB9C605** — Abstract after understanding the stable common concept, not merely after observing matching syntax.
- [ ] **USEQ-9A07F40A** — Require shared abstractions to have a clear semantic contract, owner, versioning policy, compatibility policy, and consumer tests.
- [ ] **USEQ-5E68F12C** — Prefer composition of small capabilities over inheritance or extension hierarchies that expose fragile implementation details.
- [ ] **USEQ-2D91D43D** — Keep extension points narrow, explicit, permission-aware, and bounded by real use cases.
- [ ] **USEQ-4A1FF0C9** — Avoid speculative genericity, configuration surfaces, plugin mechanisms, and type parameters without demonstrated consumers.
- [ ] **USEQ-69AF721B** — Duplicate small, volatile code temporarily when sharing would couple independent domains; record the decision if material.
- [ ] **USEQ-3A790373** — Centralize security-sensitive validation, authorization, cryptographic use, privacy rules, and dangerous operations where divergence would create risk.
- [ ] **USEQ-6E130389** — Do not hide fundamentally different latency, failure, transaction, ownership, or consistency semantics behind one misleading interface.
- [ ] **USEQ-69283C1D** — Make reusable components safe by default and difficult to use incorrectly.
- [ ] **USEQ-4BBEF611** — Document thread safety, reentrancy, ownership, lifecycle, mutability, blocking behavior, side effects, and error semantics for reusable units.
- [ ] **USEQ-309520B0** — Test reusable components against multiple representative consumers rather than only their original use case.
- [ ] **USEQ-ACE422A4** — Version shared components conservatively and provide migration support for breaking changes.
- [ ] **USEQ-FB8D2947** — Measure reuse by reduced total change and defect cost, not by the number of consumers or abstraction layers.
- [ ] **USEQ-22929EAB** — Delete unused generality and retire abstractions whose maintenance cost exceeds their value.

#### Naming, readability, and cognitive load

- [ ] **USEQ-A9C943F1** — Use domain language consistently and maintain a glossary for terms whose meaning is not obvious or is contested.
- [ ] **USEQ-747F50E2** — Name values and operations by meaning and effect rather than storage type, implementation mechanism, or historical accident.
- [ ] **USEQ-182734ED** — Include units, scope, state, direction, or security sensitivity in names when omission can cause misuse.
- [ ] **USEQ-DEF452E6** — Avoid negated, ambiguous, overloaded, misleading, joke, temporary, and context-dependent names.
- [ ] **USEQ-18FF61DA** — Keep one level of abstraction within a unit of code when practical and extract detail that obscures the primary intent.
- [ ] **USEQ-2C7F27CA** — Use early validation or guard structure when it reduces nesting without hiding cleanup or transactional behavior.
- [ ] **USEQ-6F203597** — Keep control flow explicit enough that all exits, retries, cancellations, side effects, and error paths can be reviewed.
- [ ] **USEQ-14DC38B4** — Prefer readable intermediate values over compressed expressions when they expose assumptions or aid diagnostics.
- [ ] **USEQ-3EDFCA38** — Keep comments synchronized with behavior and remove comments that merely repeat syntax or preserve obsolete history.
- [ ] **USEQ-D4930961** — Explain why a surprising constraint, workaround, algorithm, threshold, ordering, security control, or compatibility behavior exists.
- [ ] **USEQ-4446A52B** — Reference decision records or issue identifiers for temporary workarounds and include removal conditions.
- [ ] **USEQ-0E8F2EB1** — Apply stable automated formatting so reviews focus on behavior rather than style churn.
- [ ] **USEQ-AF70F539** — Keep changes focused; separate mechanical refactoring, generated output, dependency updates, and behavior changes when doing so improves review confidence.
- [ ] **USEQ-60313694** — Ensure another qualified engineer can explain critical code, failure behavior, and invariants without consulting the original author.

#### State, effects, time, and determinism

- [ ] **USEQ-86ACCDA1** — Minimize mutable shared state and make the owner, lifetime, synchronization, and persistence of every material state explicit.
- [ ] **USEQ-2D7758AD** — Distinguish source-of-truth state, derived state, cached state, transient workflow state, and presentation state.
- [ ] **USEQ-D200CC6D** — Derive values instead of storing duplicates when derivation is reliable and affordable.
- [ ] **USEQ-D30F8D53** — When duplicated state is necessary, define synchronization, reconciliation, invalidation, conflict, and repair behavior.
- [ ] **USEQ-6121CE83** — Make side effects explicit at boundaries and separate them from pure transformations where practical.
- [ ] **USEQ-0DA007A6** — Pass time, randomness, identity, locale, configuration, and external effects through controllable interfaces when deterministic testing matters.
- [ ] **USEQ-85AD8D5E** — Use a consistent time model and distinguish wall-clock time, monotonic elapsed time, business date, event time, processing time, and effective time.
- [ ] **USEQ-27E5FAAF** — Never infer elapsed duration from a wall clock that may jump.
- [ ] **USEQ-A9FAC9E7** — Define clock-skew tolerance and authoritative time sources for distributed decisions, tokens, ordering, billing, and audit.
- [ ] **USEQ-9B8A7D56** — Use explicit random seeds for reproducible tests and experiments while using cryptographically secure randomness for security-sensitive values.
- [ ] **USEQ-F0B0FCA9** — Make serialization and deserialization preserve semantic meaning, version, precision, and unknown-field behavior.
- [ ] **USEQ-DD5BE9BC** — Ensure restarting, retrying, replaying, or resuming cannot silently repeat irreversible effects.
- [ ] **USEQ-1FC568C7** — Define state-machine transitions explicitly for complex workflows and reject illegal transitions.
- [ ] **USEQ-E343856B** — Use deterministic ordering or tie-breaking where nondeterministic output would harm tests, users, signatures, caches, or reconciliation.

#### Error handling, cancellation, and recovery

- [ ] **USEQ-61D9C5BA** — Define an error taxonomy that distinguishes invalid input, authorization denial, conflict, absence, timeout, cancellation, dependency failure, resource exhaustion, invariant violation, and unexpected defect.
- [ ] **USEQ-2490A721** — Preserve the original cause, operation context, correlation information, and retryability without exposing sensitive internals to untrusted callers.
- [ ] **USEQ-51DDEC07** — Handle errors at the layer that has enough context to recover, translate, compensate, or communicate; otherwise propagate them intact.
- [ ] **USEQ-FD48C8C6** — Do not catch broad failures merely to continue in an unknown or corrupted state.
- [ ] **USEQ-9EE161FC** — Do not silently discard errors from cleanup, background work, callbacks, asynchronous tasks, streams, or destructors when they affect correctness.
- [ ] **USEQ-13655EB8** — Make retry policies explicit, bounded, observable, and limited to operations that are safe to repeat.
- [ ] **USEQ-2ED61C97** — Use backoff, jitter, deadlines, and budgets so retries do not amplify overload or exceed user expectations.
- [ ] **USEQ-1230BB14** — Propagate cancellation and deadlines across internal and external calls where work should stop.
- [ ] **USEQ-7DD3C0A2** — Define what happens to partial work when a request is canceled, a worker terminates, or a dependency times out.
- [ ] **USEQ-97D83064** — Use compensation only when true rollback is impossible and verify compensation itself is retryable and reconcilable.
- [ ] **USEQ-703563EC** — Return user-actionable errors while retaining richer diagnostic context in protected telemetry.
- [ ] **USEQ-BCFA2E39** — Avoid using exceptional mechanisms for expected high-volume control flow when that harms clarity or performance.
- [ ] **USEQ-9E0F9B73** — Test error translation so one layer does not convert distinct failures into misleading success, absence, or generic retry.
- [ ] **USEQ-9B3E4081** — Verify recovery after process restart, duplicate delivery, network partition, storage failure, and partial commit.

#### Resources, concurrency, asynchronous work, and transactions

- [ ] **USEQ-758F3740** — Define acquisition, ownership, transfer, sharing, release, timeout, cancellation, and cleanup for every finite resource.
- [ ] **USEQ-929CC93B** — Use structured lifetime management so cleanup occurs on success, failure, cancellation, and early return.
- [ ] **USEQ-E112B91C** — Bound memory, recursion, stack, threads, tasks, processes, connections, files, sockets, handles, subscriptions, queues, buffers, batches, and temporary storage.
- [ ] **USEQ-07EE43CB** — Avoid holding scarce resources across slow or unbounded external calls unless necessary and justified.
- [ ] **USEQ-FEDFBF14** — Document thread safety, process safety, reentrancy, ordering, atomicity, isolation, and visibility guarantees.
- [ ] **USEQ-7AC74572** — Protect shared invariants with one coherent synchronization or transactional strategy.
- [ ] **USEQ-B938A5BB** — Use the narrowest lock scope compatible with correctness and define lock ordering to prevent deadlock.
- [ ] **USEQ-850EB655** — Avoid blocking operations on execution contexts that must remain responsive.
- [ ] **USEQ-414A8CDE** — Apply backpressure between producers and consumers and define overload behavior before queues become unbounded.
- [ ] **USEQ-C4C5F2C3** — Treat timeouts as uncertain outcomes when the remote side may have completed the operation.
- [ ] **USEQ-F6C7947D** — Use idempotency keys, deduplication, sequence numbers, or reconciliation where messages and requests may be repeated.
- [ ] **USEQ-2AD963FD** — Define delivery semantics and never assume exactly-once execution without an end-to-end proof.
- [ ] **USEQ-995147C5** — Keep transaction boundaries aligned with business invariants and avoid long transactions that span unreliable dependencies.
- [ ] **USEQ-9E11F421** — Use optimistic or pessimistic concurrency deliberately and surface conflicts clearly.
- [ ] **USEQ-1CB07EFE** — Test race conditions, starvation, deadlock, livelock, reordering, duplicate execution, partial visibility, and shutdown under load.

#### Secure and privacy-preserving implementation

- [ ] **USEQ-FB0E62E5** — Treat every external, cross-process, cross-tenant, stored, queued, cached, imported, and deserialized value as untrusted until validated for its use.
- [ ] **USEQ-74C08945** — Perform authorization at every trusted operation using current actor, tenant, resource, action, context, and policy.
- [ ] **USEQ-46A8E662** — Keep authentication, authorization, validation, privacy, and audit decisions out of client-only or caller-controlled code.
- [ ] **USEQ-191810F1** — Use safe parameterization and context-specific encoding for queries, commands, markup, templates, headers, logs, spreadsheets, and other interpreters.
- [ ] **USEQ-92C2BC98** — Constrain outbound destinations, protocols, redirects, DNS resolution, file paths, and parsers that process attacker-influenced input.
- [ ] **USEQ-FEEEAD00** — Use approved cryptography through maintained libraries and prevent callers from selecting insecure algorithms or modes.
- [ ] **USEQ-B00F5F93** — Avoid secret-dependent branching, comparison, logging, caching, or error behavior where side channels are material.
- [ ] **USEQ-887FDDCF** — Minimize collection, retention, copying, logging, tracing, caching, indexing, and exposure of personal or confidential data.
- [ ] **USEQ-0064D510** — Redact or tokenize sensitive values before they enter generic diagnostics, analytics, exception systems, or support tools.
- [ ] **USEQ-C6D97724** — Prevent one user, tenant, request, job, cache key, log context, or background task from inheriting another's security context.
- [ ] **USEQ-D502B36F** — Use least-privilege identities and capabilities for code, jobs, tests, migrations, and administrative utilities.
- [ ] **USEQ-0DED5F53** — Make dangerous APIs hard to call accidentally through explicit types, capability objects, validation, and review gates.
- [ ] **USEQ-DB12D4C5** — Fail closed for authorization and confidentiality while designing availability fallbacks deliberately.
- [ ] **USEQ-96F9AAB8** — Add regression tests for every material security and privacy defect.

#### Performance, scalability, and resource efficiency in code

- [ ] **USEQ-15F1FD62** — Define performance requirements and resource budgets for critical operations before optimizing.
- [ ] **USEQ-B3518752** — Use representative profiling and measurement to identify dominant costs rather than guessing from local code appearance.
- [ ] **USEQ-ED6D076A** — Select algorithms and data structures whose worst-case behavior is acceptable for adversarial and maximum supported input.
- [ ] **USEQ-A577B30D** — Bound user-controlled computational complexity, allocation, fan-out, query cardinality, decompression, parsing, and serialization.
- [ ] **USEQ-729C23DC** — Avoid repeated remote calls, repeated parsing, accidental nested scans, redundant serialization, and hidden per-item I/O in collection operations.
- [ ] **USEQ-00C6A928** — Batch work only when it preserves latency, fairness, memory, failure isolation, and transactional semantics.
- [ ] **USEQ-CFB9CE23** — Cache only with a defined key, authorization scope, freshness model, invalidation path, memory bound, and failure behavior.
- [ ] **USEQ-9D9D041C** — Prevent cache stampedes, hot keys, unbounded cardinality, and cross-user data reuse.
- [ ] **USEQ-FF566A37** — Avoid premature micro-optimization that obscures correctness; preserve benchmark and profiling evidence for non-obvious optimized code.
- [ ] **USEQ-43DA35AB** — Test performance under cold starts, warm state, realistic data skew, maximum object size, concurrency, and dependency degradation.
- [ ] **USEQ-3532A80B** — Track algorithmic and allocation regressions in automated tests where stable measurement is possible.
- [ ] **USEQ-B98C679D** — Prefer streaming or incremental processing when full materialization creates unacceptable memory or latency.
- [ ] **USEQ-7A4A1550** — Release resources promptly and avoid retaining graphs, closures, listeners, or caches longer than their useful lifetime.
- [ ] **USEQ-4A549DE1** — Consider energy, bandwidth, storage, and client-device cost when alternatives provide equivalent product outcomes.

#### Testability, verification, and proof obligations

- [ ] **USEQ-88C87EB4** — Design critical logic so inputs, outputs, dependencies, time, randomness, side effects, and failure modes can be controlled and observed.
- [ ] **USEQ-8140ABCD** — Keep unit boundaries aligned with meaningful behavior rather than testing private implementation trivia.
- [ ] **USEQ-DA8D2229** — Use examples for known cases and properties or invariants for broad input spaces.
- [ ] **USEQ-2C2DD5CB** — Use generative, fuzz, mutation, model-based, differential, concurrency, and fault-injection techniques where ordinary examples leave material risk.
- [ ] **USEQ-2959B58C** — Verify parser, protocol, state-machine, numerical, authorization, financial, cryptographic, and migration logic against independent models or oracles where warranted.
- [ ] **USEQ-F41CC7BD** — Require tests to fail for the defect they are intended to prevent and use mutation or deliberate fault seeding selectively to validate test effectiveness.
- [ ] **USEQ-0B7747C8** — Avoid mocks that reproduce implementation details while failing to represent real dependency contracts.
- [ ] **USEQ-B85907BD** — Use contract tests and representative emulators, sandboxes, or test instances for external systems.
- [ ] **USEQ-61D6A6CC** — Make test fixtures explicit, minimal, valid, isolated, deterministic, and representative of production edge conditions.
- [ ] **USEQ-B9ABC567** — Ensure tests clean up resources and cannot pass because of execution order, shared state, time zone, locale, machine speed, or network access unless intended.
- [ ] **USEQ-FC90E237** — Classify flaky tests as product defects in the engineering system and fix root causes rather than normalizing retries.
- [ ] **USEQ-04007F77** — Map critical requirements and risks to verification evidence and identify untested assumptions.
- [ ] **USEQ-3397D184** — Apply static analysis, type checking, conformance checking, symbolic execution, formal specification, model checking, proof, or runtime verification when consequence and tractability justify them.
- [ ] **USEQ-53209324** — Independently review proofs, models, generators, test oracles, and safety claims because defects can exist in the verification system itself.
- [ ] **USEQ-E7843EF8** — Test the built and configured artifact, not only source-level units.

#### Observability and diagnosability in implementation

- [ ] **USEQ-B9CB9777** — Define which decisions, state transitions, external calls, retries, security events, and business outcomes require telemetry.
- [ ] **USEQ-BE07FB46** — Use stable event names, field definitions, units, identifiers, severity, and versioning so telemetry remains interpretable across releases.
- [ ] **USEQ-D1487D42** — Propagate correlation and causation context across asynchronous and distributed work without trusting caller-supplied identifiers blindly.
- [ ] **USEQ-86A37C88** — Make important failures distinguishable by cause, affected operation, dependency, tenant, release, and retryability.
- [ ] **USEQ-5866714D** — Avoid high-cardinality, secret-bearing, personal, attacker-controlled, or unbounded telemetry fields.
- [ ] **USEQ-E6F2C57A** — Ensure diagnostic logging cannot change business behavior, exhaust critical resources, or introduce deadlocks.
- [ ] **USEQ-21B0B78A** — Use sampling that preserves rare critical events and supports unbiased interpretation.
- [ ] **USEQ-6AB83E4E** — Make audit records tamper-resistant and semantically distinct from ordinary debug logs where accountability is required.
- [ ] **USEQ-0488BD90** — Provide enough state and decision context to reconstruct consequential outcomes without logging prohibited data.
- [ ] **USEQ-8693CC5C** — Test telemetry assertions, alert predicates, dashboards, and trace propagation as part of behavior.
- [ ] **USEQ-0D0A9B71** — Make telemetry pipeline failure observable and ensure the application has deliberate degradation behavior.
- [ ] **USEQ-0FB8250E** — Remove temporary debug logging, probes, and sensitive instrumentation before production unless explicitly approved.

#### Configuration, flags, dependencies, and generated code

- [ ] **USEQ-7BECCA66** — Define a schema, type, allowed range, default, owner, sensitivity, environment scope, and reload behavior for every material configuration value.
- [ ] **USEQ-154F3E57** — Fail startup or activation safely when required configuration is absent, malformed, contradictory, stale, or insecure.
- [ ] **USEQ-821892E1** — Do not let a missing flag service, configuration store, or secret silently enable risky behavior.
- [ ] **USEQ-4C6EE6DE** — Test all material feature-flag combinations, transitions, targeting rules, default states, and rollback paths.
- [ ] **USEQ-E6BA3845** — Remove stale flags, compatibility shims, dead branches, experimental paths, and obsolete configuration after their exit criteria are met.
- [ ] **USEQ-0D903230** — Pin and verify build and runtime dependencies sufficiently to prevent unexpected substitution.
- [ ] **USEQ-A7D84FB4** — Assess dependency API surface, transitive graph, maintenance health, license, provenance, security history, performance, and exit cost before adoption.
- [ ] **USEQ-49AA2475** — Wrap vendor or unstable dependencies only where the wrapper provides a real boundary, not ceremonial indirection.
- [ ] **USEQ-FA7B7D7F** — Treat generated code, models, schemas, clients, migrations, and assets as reproducible outputs with a trusted generator, version, review path, and diff policy.
- [ ] **USEQ-C8FBA56E** — Do not edit generated output manually unless the process explicitly preserves and verifies such changes.
- [ ] **USEQ-2C0D2973** — Ensure dependency upgrades, generator upgrades, and configuration changes receive behavior, compatibility, security, and performance testing.
- [ ] **USEQ-C0C75F62** — Remove unused dependencies and capabilities to reduce attack surface and maintenance cost.
- [ ] **USEQ-497725A5** — Provide a replacement or migration plan for critical unsupported, proprietary, or single-maintainer dependencies.

#### Review, change safety, technical debt, and maintainability

- [ ] **USEQ-96222956** — Require every material change to explain intent, scope, risk, alternatives, testing, deployment, rollback, data impact, and observability impact.
- [ ] **USEQ-E6E8E65D** — Use reviewers with relevant domain, security, privacy, accessibility, data, or operational expertise for high-impact changes.
- [ ] **USEQ-B42D60FB** — Review changed behavior in context of callers, consumers, data flows, concurrency, failure paths, and deployment order, not only the diff.
- [ ] **USEQ-8B195A70** — Use automated checks to remove mechanical review burden while never substituting automation for judgment on semantics and risk.
- [ ] **USEQ-3F1DB8A6** — Keep review size and latency within limits that preserve attention; split work by coherent behavior, not arbitrary line count.
- [ ] **USEQ-954F522D** — Record dissent and unresolved design questions for consequential changes.
- [ ] **USEQ-EE7867F6** — Treat technical debt as a specific future cost or risk with evidence, owner, affected outcomes, and trigger for repayment.
- [ ] **USEQ-FD32E7E9** — Do not label missing requirements, security defects, or production failures as harmless technical debt.
- [ ] **USEQ-8DC4A289** — Refactor behind characterization tests or stronger specifications when current behavior is poorly understood.
- [ ] **USEQ-BC601120** — Separate behavior-preserving refactoring from behavior changes where it improves review and rollback.
- [ ] **USEQ-49D3FE94** — Delete dead code and obsolete compatibility paths after confirming no runtime, data, customer, or external dependency remains.
- [ ] **USEQ-57E2CEAB** — Measure maintainability using defect patterns, change lead time, review difficulty, dependency structure, cognitive load, recovery cost, and ownership—not one complexity metric alone.
- [ ] **USEQ-6F163940** — Turn escaped defects into regression tests and systemic prevention when practical.
- [ ] **USEQ-473642A2** — Schedule maintenance before support deadlines, degradation, or concentration of knowledge creates an emergency.
- [ ] **USEQ-CA602113** — Ensure code ownership and documentation survive staff changes and do not depend on one person.

#### Code-specific release blockers

- [ ] **USEQ-97163345** — Do not release critical behavior whose invariants, authorization, failure modes, or data effects cannot be explained and tested.
- [ ] **USEQ-CD7D768A** — Do not release with known undefined, implementation-dependent, race-prone, overflow-prone, or precision-losing behavior in a critical path.
- [ ] **USEQ-F78B1390** — Do not release when an error, timeout, retry, cancellation, or restart can cause silent corruption, duplicate irreversible effects, or false success.
- [ ] **USEQ-205749EA** — Do not release unbounded work or allocation reachable from untrusted input when it can exhaust shared resources.
- [ ] **USEQ-DD8B007F** — Do not release code that logs reusable secrets, session material, payment credentials, or prohibited personal data.
- [ ] **USEQ-DF40914E** — Do not release a critical dependency, generated artifact, or build output whose source and version cannot be traced and reproduced.
- [ ] **USEQ-09B573EE** — Do not release high-impact code that was self-approved without the required independent review.
- [ ] **USEQ-C1797663** — Do not release a material refactor without evidence that required behavior and compatibility were preserved.
- [ ] **USEQ-2B15F92E** — Do not waive failing tests, static findings, warnings, or code-review concerns without a documented technical disposition and owner.
- [ ] **USEQ-8B433E08** — Do not claim high code quality solely from coverage percentage, linter cleanliness, low complexity, style conformance, or absence of known defects.

## Standards and source references

- [ISO/IEC 25010:2023 — Product quality model](https://www.iso.org/standard/78176.html)
- [ISO/IEC/IEEE 29148:2018 — Requirements engineering](https://www.iso.org/standard/72089.html)
- [ISO/IEC/IEEE 29119-4:2021 — Test techniques](https://www.iso.org/standard/79430.html)
- [OWASP Application Security Verification Standard 5.0.0](https://owasp.org/www-project-application-security-verification-standard/)
- [OWASP Top 10 — 2025](https://owasp.org/www-project-top-ten/)
- [OWASP Web Security Testing Guide 4.2](https://owasp.org/www-project-web-security-testing-guide/)
- [ISO/IEC/IEEE 12207:2026 — Software life cycle processes](https://www.iso.org/standard/90219.html)
- [ISO/IEC 5055:2021 — Automated source code quality measures](https://www.iso.org/standard/80623.html)
- [ISO/IEC 20246:2017 — Work product reviews](https://www.iso.org/standard/67597.html)
- [ISO/IEC/IEEE 42010:2022 — Architecture description](https://www.iso.org/standard/74393.html)
- [ISO/IEC 27001:2022 — Information security management systems](https://www.iso.org/standard/27001)
- [NIST SP 800-218 v1.1 — Secure Software Development Framework](https://csrc.nist.gov/pubs/sp/800/218/final)
- [SLSA Specification 1.2](https://slsa.dev/spec/v1.2/)
- [OpenSSF Security Baseline and Best Practices](https://baseline.openssf.org/)
- [SPDX 3.0.1 / ISO/IEC 5962:2021](https://spdx.github.io/spdx-spec/v3.0.1/)
- [ISO/IEC/IEEE 29119-2:2021 — Test processes](https://www.iso.org/standard/79428.html)
- [IEEE Computer Society SWEBOK v4](https://www.computer.org/education/bodies-of-knowledge/software-engineering)
- [SEI secure development resources](https://www.sei.cmu.edu/secure-development/)

---

[Previous phase](04-architecture-and-design.md) · [Next: Phase 6: Application services and APIs](06-application-services-and-apis.md)
