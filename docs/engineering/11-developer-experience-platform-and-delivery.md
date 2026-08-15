# Developer experience, platform, and delivery

_Phase 11 of 16 in the [complete engineering review](00-overview.md)._

Developer flow, platform engineering, engineering economics, source control, CI/CD, trusted builds, infrastructure as code, and progressive delivery.

Assess every control as **Pass**, **Fail**, **Blocked**, or **Not Applicable**. A Pass needs current evidence for the exact product, revision, artifact, configuration, environment, and data state under review.

## Source Control and Change Management

_Consolidated from `quality standards/12-delivery-cicd/01-source-control-and-change-management.md`; 9 non-duplicative controls._

### Universal controls

- [ ] **USEQ-836D2261** — Repository access uses individual identities, strong authentication, and least privilege.
- [ ] **USEQ-3D3EF3D3** — Material changes require independent review; critical areas require designated reviewers.
- [ ] **USEQ-0BBDD49A** — Administrative review bypass cannot occur silently.
- [ ] **USEQ-4B76B0C6** — Merge requirements include all required automated checks.
- [ ] **USEQ-2163B1AC** — Security-sensitive, infrastructure, pipeline, authorization, cryptography, payment, and migration changes receive heightened review.
- [ ] **USEQ-32DAA868** — Generated code has a trusted source and reproducible regeneration process.
- [ ] **USEQ-2C510B3A** — Emergency changes use a documented break-glass process and receive retrospective review.
- [ ] **USEQ-AE4E848F** — Repository backup and recovery are tested.
- [ ] **USEQ-3560AE8E** — Repository webhooks, applications, bots, and deploy keys are inventoried and restricted.

## Build, CI/CD, and Artifact Integrity

_Consolidated from `quality standards/12-delivery-cicd/02-build-cicd-and-artifact-integrity.md`; 11 non-duplicative controls._

### Universal controls

- [ ] **USEQ-394F34E8** — Builds run automatically from reviewed source and are repeatable from declared inputs.
- [ ] **USEQ-D1076D92** — Dependency versions and build inputs are pinned sufficiently to prevent unexpected substitution.
- [ ] **USEQ-267DAD35** — The production artifact is built once and promoted rather than independently rebuilt per environment.
- [ ] **USEQ-3381F5D0** — Deployable artifacts are immutable and have cryptographic digests.
- [ ] **USEQ-2F6AA097** — Build provenance records source, builder, inputs, workflow, time, and output digest.
- [ ] **USEQ-4480E17F** — Build and signing identities use least privilege; keys are protected outside source and output.
- [ ] **USEQ-73488048** — Untrusted contribution workflows cannot access production secrets or overwrite trusted artifacts.
- [ ] **USEQ-A0338698** — Failed or canceled checks cannot be represented as passed.
- [ ] **USEQ-4C73564C** — Previous known-good artifacts remain available.
- [ ] **USEQ-7112CE32** — Build logs and caches do not expose secrets, sensitive source, or unverified cross-trust artifacts.
- [ ] **USEQ-D4ABAEAF** — Build-platform compromise has containment and recovery procedures.

## Environments, Configuration, Feature Flags, and Secrets

_Consolidated from `quality standards/12-delivery-cicd/03-environments-configuration-flags-and-secrets.md`; 24 non-duplicative controls._

### Environment separation

- [ ] **USEQ-9ED6000B** — Production credentials do not work outside production; non-production users cannot reach production control planes.
- [ ] **USEQ-23C9AA06** — Test accounts, bypasses, simulators, debug endpoints, and development consoles are absent or inaccessible in production.
- [ ] **USEQ-DAC810AD** — Non-production email, payment, notifications, and webhooks cannot affect real users or systems accidentally.
- [ ] **USEQ-AF8234B1** — Environment boundaries are documented, tested, logged, and periodically reviewed.

### Configuration

- [ ] **USEQ-7A52E01F** — Configuration is versioned or otherwise auditable and changes require review.
- [ ] **USEQ-447535A7** — Secure defaults and explicit environment-specific values are used.
- [ ] **USEQ-5E325DC6** — Production configuration can be reconstructed and rolled back.
- [ ] **USEQ-294721C4** — Debug mode, verbose exceptions, profiling, test routes, and development consoles are disabled or restricted.
- [ ] **USEQ-B66A3064** — Time zones, locales, encodings, units, time synchronization, resource limits, and timeouts are explicit.
- [ ] **USEQ-6A62F162** — Source maps and diagnostic artifacts do not expose secrets or inappropriate internal information.
- [ ] **USEQ-DE280DA4** — Configuration ownership, review dates, emergency changes, and approvals are recorded.

### Feature flags

- [ ] **USEQ-3E50B6BE** — Every flag has an owner, purpose, default state, and removal date where temporary.
- [ ] **USEQ-7BDEFE8A** — Flag-service failure produces safe behavior.
- [ ] **USEQ-ADCFB9BE** — High-impact flag changes are restricted and audited.
- [ ] **USEQ-002EFF2A** — High-risk new functionality has a kill switch where appropriate.
- [ ] **USEQ-98C06345** — Untrusted clients cannot manipulate flags to bypass controls.

### Secrets

- [ ] **USEQ-A088A7BA** — Secrets do not appear in source, artifacts, frontend bundles, logs, tickets, documentation, chat, or analytics.
- [ ] **USEQ-DDA7DCA9** — Least privilege and short-lived credentials are used where possible.
- [ ] **USEQ-5970F4E7** — Secret retrieval and use are audited where practical.
- [ ] **USEQ-F341514B** — Rotation and emergency revocation are tested.
- [ ] **USEQ-402D416B** — Expiry and rotation failures are monitored.
- [ ] **USEQ-2CE27D9F** — Backup and recovery preserve or securely reconstruct required secrets.
- [ ] **USEQ-BF6C9A29** — Break-glass credentials are protected, monitored, and tested.
- [ ] **USEQ-24DBEE2F** — Historical secret exposure is addressed by rotation or revocation, not merely deletion.

## Continuous Integration

_Consolidated from `quality standards/12-delivery-cicd/04-continuous-integration.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-69535E01** — Integrate small complete changes frequently rather than allowing long-lived divergence.
- [ ] **USEQ-782703AC** — Keep the authoritative branch releasable or restore it promptly when a gate fails.
- [ ] **USEQ-0B41367E** — Require review and required checks before protected-branch change.
- [ ] **USEQ-7B40E59A** — Use one trusted pipeline definition under the same change control as source.
- [ ] **USEQ-8CA318E4** — Run builds and tests from clean, declared inputs.
- [ ] **USEQ-8574E516** — Fail closed when mandatory checks cannot run or produce ambiguous status.
- [ ] **USEQ-D8198F66** — Separate trusted branch and release jobs from untrusted contribution contexts.
- [ ] **USEQ-E3F935AB** — Prevent untrusted jobs from accessing secrets, signing material, deployment rights, or mutable trusted caches.
- [ ] **USEQ-CB522CDF** — Run fast high-signal checks early and broader risk-based checks before acceptance.
- [ ] **USEQ-007A0FED** — Verify formatting, compilation, static analysis, unit tests, contracts, secrets, dependencies, policy, and packaging as applicable.
- [ ] **USEQ-BAE610E4** — Make check output diagnostic and link failures to evidence.
- [ ] **USEQ-3B38BEB9** — Track flaky tests, unstable infrastructure, and long queue time as integration defects.
- [ ] **USEQ-F97BCD7F** — Avoid rerunning failed jobs until they pass without understanding the cause.
- [ ] **USEQ-F4D6E3B3** — Pin toolchain, action, plugin, and dependency versions sufficiently for trustworthy results.
- [ ] **USEQ-E307C67C** — Use isolated workspaces and clean up credentials and artifacts after each run.
- [ ] **USEQ-37F01160** — Make generated files reproducible and detect uncommitted or stale generation.
- [ ] **USEQ-609408D4** — Prevent merge queues from accepting a combination never tested together.
- [ ] **USEQ-FF347E02** — Record source revision, pipeline version, inputs, environment, and artifacts for each result.
- [ ] **USEQ-B022185A** — Monitor lead time, failure rate, recovery time, queue time, and feedback latency without using them as individual performance scores.
- [ ] **USEQ-A0BC4022** — Exercise the pipeline's ability to reject deliberately broken, vulnerable, or unsigned changes.

## Continuous Delivery and Deployment

_Consolidated from `quality standards/12-delivery-cicd/05-continuous-delivery-and-deployment.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-8DA2AF16** — Build a release artifact once and promote the identical immutable artifact through environments.
- [ ] **USEQ-517E9DBD** — Separate artifact creation from deployment authorization.
- [ ] **USEQ-5FBAC758** — Verify provenance, signature, digest, policy, and environment eligibility before promotion.
- [ ] **USEQ-455A9A23** — Use least-privilege deployment identities scoped to intended targets.
- [ ] **USEQ-E2718EA6** — Store deployment definitions, configuration references, and policy as reviewed versioned artifacts.
- [ ] **USEQ-20FFD65B** — Require explicit approval where risk, law, contract, or segregation of duties demands it.
- [ ] **USEQ-F987D323** — Make automated approvals evidence-based and auditable.
- [ ] **USEQ-82AC6F44** — Use progressive exposure and health gates proportionate to blast radius.
- [ ] **USEQ-BC77A748** — Validate user-facing success, business integrity, latency, errors, saturation, security, and dependency health during rollout.
- [ ] **USEQ-D63DD56F** — Pause or roll back automatically or operationally when predefined thresholds fail.
- [ ] **USEQ-1510AAE1** — Keep previous known-good artifacts and configuration available.
- [ ] **USEQ-E09E9C8D** — Test rollback and safe roll-forward, including data and event compatibility.
- [ ] **USEQ-61A0F909** — Prevent deployment when required evidence is stale, missing, or tied to another artifact.
- [ ] **USEQ-79029B47** — Serialize or coordinate changes that cannot safely overlap.
- [ ] **USEQ-98255DC6** — Make deployment idempotent and resumable after interruption.
- [ ] **USEQ-C3E713ED** — Record who or what authorized, initiated, changed, verified, paused, and completed the deployment.
- [ ] **USEQ-913347A6** — Prevent direct production mutation outside controlled emergency procedures.
- [ ] **USEQ-BE458954** — Test emergency release paths without normalizing bypass.
- [ ] **USEQ-5A5B48C9** — Separate environment secrets and trust relationships.
- [ ] **USEQ-44042363** — Measure deployment outcomes and improve the system after failures and near misses.

## Reproducible, Hermetic, and Trusted Builds

_Consolidated from `quality standards/12-delivery-cicd/06-reproducible-hermetic-and-trusted-builds.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-8555A412** — Declare source, dependencies, toolchain, environment, configuration, and build parameters as explicit inputs.
- [ ] **USEQ-E9C8C459** — Pin mutable inputs or verify them cryptographically.
- [ ] **USEQ-907EC8FF** — Restrict undeclared network access during trusted builds where practical.
- [ ] **USEQ-DED09248** — Use clean isolated builders and prevent cross-tenant or cross-trust contamination.
- [ ] **USEQ-1D16BCDA** — Make locale, time zone, timestamps, file order, random seeds, and environment influence deterministic or recorded.
- [ ] **USEQ-269DB62C** — Separate build-time secrets from output and prevent them from entering artifacts or logs.
- [ ] **USEQ-F73B5A06** — Use minimum build permissions and short-lived credentials.
- [ ] **USEQ-4B359B84** — Prevent source from modifying trusted builder policy or accessing unrelated release secrets.
- [ ] **USEQ-2A35153B** — Generate immutable output digests and machine-verifiable provenance.
- [ ] **USEQ-BE8EF1D3** — Rebuild independently and compare outputs or meaningful normalized content for high-assurance releases.
- [ ] **USEQ-B910BD14** — Detect generated content that differs from checked-in or declared source.
- [ ] **USEQ-D60E6950** — Verify final packages contain only intended files, permissions, metadata, and dependencies.
- [ ] **USEQ-37514564** — Scan and inspect the final artifact, not only its source inputs.
- [ ] **USEQ-48587D52** — Sign artifacts only after required verification succeeds.
- [ ] **USEQ-587FE097** — Protect signing keys and separate signing authority from untrusted build execution.
- [ ] **USEQ-C33F05AA** — Retain builders, inputs, provenance, artifacts, and logs for incident investigation.
- [ ] **USEQ-1A4ABD7D** — Maintain a clean-room rebuild procedure after suspected compromise.
- [ ] **USEQ-764FB544** — Rotate and revoke signing identity and redistribute corrected artifacts when necessary.
- [ ] **USEQ-6B385742** — Test reproducibility after toolchain and platform updates.
- [ ] **USEQ-FC20FB0F** — Document remaining nondeterminism and why it does not undermine trust.

## Infrastructure as Code and Environment Engineering

_Consolidated from `quality standards/12-delivery-cicd/07-infrastructure-as-code-and-environment-engineering.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-E4DC3BC6** — Represent material infrastructure, policy, networking, identity, observability, storage, and deployment configuration as versioned code or equivalent declarative records.
- [ ] **USEQ-F53D21B2** — Apply review, protected branches, traceability, testing, and release controls to infrastructure changes.
- [ ] **USEQ-9F1306DB** — Separate reusable modules from environment-specific composition.
- [ ] **USEQ-FC258A07** — Define secure defaults and prohibit dangerous options through policy where practical.
- [ ] **USEQ-59037C06** — Validate syntax, types, references, policy, cost, security, and destructive impact before apply.
- [ ] **USEQ-E0C505B4** — Generate a change plan and require human review for high-impact or destructive differences.
- [ ] **USEQ-AAC85EFD** — Use least-privilege automation identities and avoid shared administrator credentials.
- [ ] **USEQ-595C5316** — Keep secrets out of source and state files; protect state as sensitive production data.
- [ ] **USEQ-A4FA3E49** — Detect and reconcile unauthorized or manual drift.
- [ ] **USEQ-7EA84BA2** — Make creation, update, replacement, scaling, failover, and destruction idempotent where possible.
- [ ] **USEQ-73830565** — Protect critical resources against accidental deletion and uncontrolled replacement.
- [ ] **USEQ-5617B518** — Version and migrate state safely.
- [ ] **USEQ-506A79E0** — Test modules and complete environments in isolated accounts or namespaces.
- [ ] **USEQ-3AFFAE81** — Reconstruct an environment from source, artifacts, configuration, data recovery, and documented external prerequisites.
- [ ] **USEQ-47EDD512** — Record provider defaults and implicit resources that code does not create.
- [ ] **USEQ-CFDDFBBF** — Define import and ownership procedures for existing resources.
- [ ] **USEQ-FCC3E68C** — Coordinate changes across application, schema, network, identity, and observability dependencies.
- [ ] **USEQ-632F6BEB** — Plan for provider API deprecation and module lifecycle.
- [ ] **USEQ-6A024846** — Measure provisioning time, drift, failed changes, rollback, policy violations, and orphaned resources.
- [ ] **USEQ-17F99C8E** — Exercise disaster recovery without relying on inaccessible production control planes.

## Deployment Strategies and Progressive Delivery

_Consolidated from `quality standards/12-delivery-cicd/08-deployment-strategies-and-progressive-delivery.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-80674A43** — Choose rolling, canary, blue-green, cohort, shadow, parallel, regional, or other strategies according to risk and state compatibility.
- [ ] **USEQ-92842952** — Define the smallest safe initial exposure and an explicit expansion sequence.
- [ ] **USEQ-2A0FF5B4** — Ensure cohorts are representative enough to reveal relevant failures.
- [ ] **USEQ-FDF15AD0** — Keep control and treatment populations comparable when using comparative gates.
- [ ] **USEQ-4896CCAE** — Separate code deployment, feature exposure, data migration, and customer communication when decoupling reduces risk.
- [ ] **USEQ-DC837F09** — Define stop, hold, rollback, roll-forward, incident, and communication thresholds before launch.
- [ ] **USEQ-203A1CDD** — Verify old and new versions can coexist with schemas, messages, caches, sessions, and clients.
- [ ] **USEQ-AED1FAFC** — Prevent one-way data changes from making rollback appear available when it is not.
- [ ] **USEQ-906DDAA4** — Warm capacity, caches, connections, and dependencies before traffic shift where needed.
- [ ] **USEQ-5554C95B** — Drain and terminate old versions safely after in-flight work completes or transfers.
- [ ] **USEQ-C8FC60DB** — Observe user success, business integrity, latency, errors, saturation, dependency health, security, and support demand.
- [ ] **USEQ-AEFC8DCA** — Account for delayed effects, background work, billing cycles, data pipelines, and long-lived sessions before expansion.
- [ ] **USEQ-722D65DC** — Prevent automatic rollout from outrunning observability or human response.
- [ ] **USEQ-A66E17A6** — Use kill switches for high-risk behavior, but test them and keep defaults safe.
- [ ] **USEQ-FD5F6635** — Keep deployment and feature state visible in incident telemetry.
- [ ] **USEQ-33B2FE53** — Verify traffic routing, stickiness, failover, and health-check semantics.
- [ ] **USEQ-993B7855** — Test reverse traffic shift and cleanup before launch.
- [ ] **USEQ-2CFDB0CB** — Record every exposure change and decision rationale.
- [ ] **USEQ-FFB58250** — Define when the rollout is complete and temporary infrastructure can be removed.
- [ ] **USEQ-AF20E3D9** — Conduct a review when rollout gates fail or produce unexplained anomalies.

## Database and Data Release Engineering

_Consolidated from `quality standards/12-delivery-cicd/09-database-and-data-release-engineering.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-3F4E7A24** — Version every schema, migration, data transform, permission, index, and pipeline change.
- [ ] **USEQ-A97FB4A9** — Analyze application, worker, report, export, cache, search, integration, and rollback compatibility.
- [ ] **USEQ-29E45716** — Use staged compatible changes when versions coexist.
- [ ] **USEQ-11851D0A** — Run destructive or long migrations separately from application deployment when separation reduces risk.
- [ ] **USEQ-D7A998B9** — Estimate duration, locks, I/O, storage, log, replication, and downstream impact using representative data.
- [ ] **USEQ-1F838C1D** — Validate preconditions and stop on unexpected state.
- [ ] **USEQ-8C9DCBEF** — Take a verified recovery point before irreversible change.
- [ ] **USEQ-EEB2E713** — Make migrations idempotent, resumable, observable, and safe under retry.
- [ ] **USEQ-A69FE821** — Throttle and checkpoint long-running changes.
- [ ] **USEQ-ADDBE0D0** — Monitor locks, errors, lag, throughput, capacity, queue growth, and user impact.
- [ ] **USEQ-1612BF52** — Verify counts, checksums, constraints, totals, samples, and business invariants.
- [ ] **USEQ-48E8D418** — Prevent old code from writing data incompatible with the new representation.
- [ ] **USEQ-33DC8C50** — Keep old fields and indexes until consumers and rollback windows are complete.
- [ ] **USEQ-D692F341** — Define safe rollback or explicit roll-forward when data cannot be reversed.
- [ ] **USEQ-B18B8ADB** — Reconcile dual writes and backfills continuously.
- [ ] **USEQ-FC5FCD6B** — Review and audit manual corrections.
- [ ] **USEQ-64F2139F** — Protect retention, deletion, consent, residency, and legal-hold semantics during migration.
- [ ] **USEQ-22A52F1C** — Test restored backups against the new application and schema.
- [ ] **USEQ-5C9EF78B** — Remove temporary elevated permissions and migration infrastructure.
- [ ] **USEQ-89C3CCB0** — Record exact migration version, operator, start, end, outcome, and evidence.

## Rollback, Roll-Forward, and Kill Switches

_Consolidated from `quality standards/12-delivery-cicd/10-rollback-rollforward-and-kill-switches.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-D9F7262D** — Define which changes can be rolled back, which require roll-forward, and which require data restoration or reconciliation.
- [ ] **USEQ-2F23054A** — Keep previous trusted application, infrastructure, configuration, and policy artifacts available.
- [ ] **USEQ-52FBF8D0** — Test rollback using the same automation and permissions used in production.
- [ ] **USEQ-B540AC6B** — Ensure schema, messages, sessions, caches, files, and external side effects remain compatible with rollback.
- [ ] **USEQ-92B08EA4** — Avoid destructive data change before the rollback window closes unless risk is explicitly accepted.
- [ ] **USEQ-2B807409** — Define rollback triggers based on user impact and integrity, not only technical health.
- [ ] **USEQ-66FE3634** — Make rollback authority and initiation paths available during an incident.
- [ ] **USEQ-3264C430** — Ensure rollback does not repeat charges, notifications, jobs, exports, or other side effects.
- [ ] **USEQ-EB2A1C1A** — Reconcile data written by the failed version.
- [ ] **USEQ-36BF8E85** — Use kill switches for isolated risky behavior and define safe failure when the flag system is unavailable.
- [ ] **USEQ-24436BCA** — Restrict and audit kill-switch access.
- [ ] **USEQ-CE7FB317** — Test kill switches in production-like conditions and periodically in production-safe drills.
- [ ] **USEQ-F343F4EA** — Prevent stale clients, workers, and long-lived connections from continuing disabled behavior.
- [ ] **USEQ-C708B71B** — Make configuration and feature state visible to responders.
- [ ] **USEQ-6C422A58** — Define recovery from a bad configuration separately from bad code.
- [ ] **USEQ-2CA6BC3B** — Use roll-forward when reverting would violate data or protocol compatibility.
- [ ] **USEQ-059BE10D** — Verify user journeys and business invariants after containment.
- [ ] **USEQ-C98A9821** — Communicate rollback effects to support and customers where material.
- [ ] **USEQ-05EB12CC** — Remove temporary rollback compatibility and emergency flags after stability is established.
- [ ] **USEQ-1A767899** — Capture lessons from every rollback or emergency disablement.

## Developer Experience, Platform Engineering, and Engineering Economics Master Checklist

_Consolidated from `gap supplement/03-developer-experience-platform-engineering-and-economics.md`; 203 non-duplicative controls._

### Expanded gap-closure controls

#### Developer experience as an engineering-system outcome

- [ ] **USEQ-9DFE53B9** — Treat developer experience as the ability of people to understand, change, verify, deliver, operate, and retire software safely and effectively.
- [ ] **USEQ-E834F54B** — Assign accountable ownership for developer experience across team boundaries rather than treating it as an informal tooling concern.
- [ ] **USEQ-41F60C9D** — Measure developer experience through multiple dimensions including feedback loops, cognitive load, flow, satisfaction, collaboration, quality, and delivery outcomes.
- [ ] **USEQ-DE7FAAF5** — Segment findings by role, tenure, product area, location, accessibility needs, and workflow without exposing individuals or creating unsafe comparisons.
- [ ] **USEQ-8AC9F599** — Use qualitative research, workflow observation, support demand, incident evidence, and tool telemetry together rather than relying on one survey or metric.
- [ ] **USEQ-E3C0E730** — Prioritize improvements by repeated friction, risk, time lost, error probability, interruption cost, and impact on user outcomes.
- [ ] **USEQ-C8FCDEDB** — Publish the engineering-system roadmap, service ownership, supported workflows, known limitations, and ways to request improvement.
- [ ] **USEQ-496042E7** — Treat internal developers, operators, data practitioners, security teams, and automated agents as distinct platform users with different needs and permissions.
- [ ] **USEQ-351D0CF1** — Include accessibility and inclusive design in internal tools, documentation, terminals, dashboards, portals, development environments, and support processes.
- [ ] **USEQ-C9B06FE8** — Prevent local heroics and undocumented workarounds from being accepted as evidence that the engineering system works.
- [ ] **USEQ-3A8FC2D5** — Turn recurring onboarding questions, build failures, deployment tickets, access requests, and operational toil into product and platform improvements.
- [ ] **USEQ-EC6588B5** — Review developer-experience risks whenever tools, organizational structure, architecture, security policy, or delivery processes materially change.

#### Feedback loops and inner development cycle

- [ ] **USEQ-B09E273B** — Define target feedback times for editing, validation, formatting, compilation or analysis, focused tests, integration tests, review, build, deployment, and production verification.
- [ ] **USEQ-C262BE4B** — Measure end-to-end waiting time and queue time rather than only tool execution time.
- [ ] **USEQ-DAC53EAB** — Make fast local or isolated verification representative enough to catch common defects before expensive shared pipelines.
- [ ] **USEQ-3C5DA5E8** — Keep slow, broad, or expensive checks available without making every small change wait for unrelated coverage.
- [ ] **USEQ-C1D8AF64** — Use test selection, caching, parallelism, incremental analysis, and prebuilt environments only when correctness and reproducibility remain trustworthy.
- [ ] **USEQ-E29091F1** — Surface failures with the exact failing step, actionable cause, relevant logs, ownership, and remediation guidance.
- [ ] **USEQ-37226F1E** — Ensure developers can reproduce material pipeline failures locally or in an accessible diagnostic environment.
- [ ] **USEQ-CC581ACE** — Prevent flaky infrastructure, intermittent tests, and overloaded shared environments from becoming normalized background noise.
- [ ] **USEQ-BDB4D329** — Prioritize repair of broken mainline builds and shared development infrastructure according to their organization-wide cost.
- [ ] **USEQ-36CB4809** — Make code review assignment, status, required expertise, unresolved comments, and decision authority visible.
- [ ] **USEQ-AB3A0B9D** — Set review response expectations that preserve quality without creating long invisible queues.
- [ ] **USEQ-658ADA6B** — Provide early production-like feedback for configuration, migrations, permissions, performance, accessibility, security, and observability.
- [ ] **USEQ-D040219F** — Measure how often failures are detected at each lifecycle stage and move detection earlier when economically and technically justified.
- [ ] **USEQ-0F9A3D1B** — Ensure feedback remains available during provider, network, identity, or platform degradation through documented fallback paths where required.

#### Cognitive load, flow, and work design

- [ ] **USEQ-011FC4C6** — Limit the number of unrelated systems, interfaces, credentials, concepts, and manual decisions required for common workflows.
- [ ] **USEQ-E8B11D40** — Provide one discoverable authoritative entry point for services, ownership, documentation, environments, deployment, support, and operational status.
- [ ] **USEQ-1211B978** — Use consistent terminology, identity, navigation, commands, templates, and policy messages across internal capabilities.
- [ ] **USEQ-C7AADF90** — Make the safe path easier than bypassing controls or assembling an unsupported custom path.
- [ ] **USEQ-DBBC45BB** — Reduce context switching by limiting unnecessary meetings, interruptions, simultaneous priorities, handoffs, and approval queues.
- [ ] **USEQ-B2FE9B81** — Protect focused work while preserving explicit channels for incidents and genuinely urgent requests.
- [ ] **USEQ-5D046A65** — Limit work in progress at team and portfolio levels and finish high-value work before starting additional low-priority work.
- [ ] **USEQ-64A222AA** — Decompose work into small independently valuable batches without losing system-level design and integration thinking.
- [ ] **USEQ-8CC09A8B** — Keep ownership boundaries clear enough that teams know where to change behavior and whom to involve.
- [ ] **USEQ-23A2452C** — Avoid shifting platform, security, operations, or compliance complexity onto every application team through poorly designed self-service.
- [ ] **USEQ-A0FAB7FB** — Provide explanations and recovery paths when automated policy blocks work; do not return opaque denial messages.
- [ ] **USEQ-390C0893** — Design documentation and tools for recognition and guided choice rather than requiring memorization of volatile procedures.
- [ ] **USEQ-66FBAB61** — Track cognitive load introduced by legacy systems, fragmented platforms, excessive customization, incompatible conventions, and organizational coupling.
- [ ] **USEQ-D41144E6** — Retire redundant tools and paths after migration rather than preserving every historical option indefinitely.

#### Onboarding, environments, and daily usability

- [ ] **USEQ-C4A6187E** — Define the minimum supported path from account creation to a verified useful change and measure it with new or reset users.
- [ ] **USEQ-262F3070** — Automate account, repository, documentation, environment, test-data, secret, and deployment access according to role and least privilege.
- [ ] **USEQ-D031FE99** — Ensure onboarding does not depend on a particular person being available.
- [ ] **USEQ-72AF77C1** — Provide a current architectural orientation, service map, glossary, development workflow, quality gates, support route, and first-task guide.
- [ ] **USEQ-E3FF70A7** — Keep development setup reproducible from documented declarations and validate it regularly from a clean state.
- [ ] **USEQ-1EA7115D** — Minimize the use of production data and production access in ordinary development.
- [ ] **USEQ-AC9B4FCF** — Provide representative synthetic, generated, masked, or contract-based test data with known limitations.
- [ ] **USEQ-CD7A2F8E** — Support isolated and ephemeral environments when shared environments create collision, queueing, privacy, or reliability problems.
- [ ] **USEQ-4A4087C7** — Control environment drift and make differences from production explicit.
- [ ] **USEQ-63ACFD44** — Ensure local and remote development workflows meet security, accessibility, performance, and support requirements.
- [ ] **USEQ-2C0ADD3B** — Provide safe seed, reset, migration, and cleanup procedures for development and test state.
- [ ] **USEQ-7E8B9D53** — Make common debugging, tracing, profiling, and dependency simulation capabilities discoverable and permissioned.
- [ ] **USEQ-598AC30F** — Document offline, low-bandwidth, travel, and device-replacement recovery where work continuity requires it.
- [ ] **USEQ-DEFAB837** — Revalidate onboarding after major platform, identity, repository, architecture, or organizational changes.

#### Platform engineering as a product

- [ ] **USEQ-37A8A1E2** — Treat the internal platform as a product with named users, product management, service ownership, roadmap, support, research, adoption goals, and lifecycle funding.
- [ ] **USEQ-5BFBA215** — Build platform capabilities from observed user workflows and organization-specific constraints rather than copying a reference stack blindly.
- [ ] **USEQ-C5BF78A7** — Deliver small usable capabilities early and evolve through feedback instead of attempting a complete big-bang platform.
- [ ] **USEQ-0E7CE9E5** — Offer self-service for frequent, well-understood, automatable operations while retaining human review where consequence requires it.
- [ ] **USEQ-F2F34933** — Provide paved or golden paths that integrate security, privacy, observability, reliability, cost, documentation, and compliance by default.
- [ ] **USEQ-F2B89A92** — Allow justified escape paths and extensions without letting every team create an ungoverned parallel platform.
- [ ] **USEQ-52F5768B** — Define platform personas, supported use cases, service levels, prerequisites, quotas, costs, responsibilities, and exit paths.
- [ ] **USEQ-57EB708C** — Use stable, versioned platform interfaces and avoid forcing consumers to depend on internal implementation details.
- [ ] **USEQ-77692B0A** — Provide a service and software catalog with ownership, criticality, dependencies, runbooks, lifecycle, data classification, and operational health.
- [ ] **USEQ-8E00CB58** — Provide reusable templates only when their generated outputs remain understandable, owned, upgradeable, and removable.
- [ ] **USEQ-8AF5B8F4** — Make platform policies testable before deployment and explain violations in user terms.
- [ ] **USEQ-E5D8E500** — Minimize ticket-based provisioning for routine capabilities and measure ticket demand as evidence of missing self-service or unclear ownership.
- [ ] **USEQ-AC94BEB4** — Prevent platform adoption metrics from rewarding forced use while ignoring task success, satisfaction, reliability, or downstream outcomes.
- [ ] **USEQ-8B054594** — Provide migration support and coexistence periods when changing platform interfaces or golden paths.
- [ ] **USEQ-707DB5CA** — Retire deprecated paths only after consumer inventory, communication, migration evidence, and recovery planning are complete.

#### Platform security, compliance, and governance

- [ ] **USEQ-D23E6AE9** — Embed least privilege, strong identity, secret handling, provenance, policy checks, logging, and approved network patterns into default platform paths.
- [ ] **USEQ-DA324337** — Separate platform control-plane privileges from ordinary application-team privileges.
- [ ] **USEQ-79C11C18** — Make high-impact platform actions attributable, reviewable, time-bounded, and auditable.
- [ ] **USEQ-21D7885D** — Prevent templates and self-service workflows from creating public exposure, excessive privilege, unbounded cost, unsupported versions, or missing ownership by default.
- [ ] **USEQ-F9F2BEED** — Validate user-supplied parameters before provisioning and bound resource, region, data, network, and identity choices.
- [ ] **USEQ-2E333332** — Use policy as code where reliable, but retain an accountable exception and appeal process.
- [ ] **USEQ-85AD04E2** — Test policy behavior against intended allow, deny, exception, rollback, and degraded cases.
- [ ] **USEQ-6E35D785** — Ensure platform automation cannot access unrelated tenant, team, environment, repository, or production secrets.
- [ ] **USEQ-8BC1E850** — Inventory nonhuman users, automation agents, CI identities, service accounts, and AI agents as first-class platform consumers.
- [ ] **USEQ-94C7E270** — Provide scoped temporary credentials rather than distributing long-lived shared credentials.
- [ ] **USEQ-F5A77AE2** — Make compliance evidence an output of normal workflows rather than a separate manual reconstruction where practical.
- [ ] **USEQ-23AEDF92** — Keep platform-generated configurations, artifacts, attestations, and audit records traceable to request, policy, source, and approver.
- [ ] **USEQ-3712C085** — Threat-model the platform as a high-leverage dependency and test compromise, abuse, supply-chain, and insider scenarios.
- [ ] **USEQ-C5BBCBD4** — Maintain break-glass recovery that works when normal identity, documentation, or platform control planes are unavailable.

#### Platform reliability and operability

- [ ] **USEQ-104A029B** — Define user-facing SLIs and SLOs for critical platform workflows such as environment creation, build, test, artifact publication, deployment, rollback, and secret delivery.
- [ ] **USEQ-105A2D16** — Monitor task success from the consumer perspective, not only platform component uptime.
- [ ] **USEQ-1B99F055** — Provide status, incident communication, ownership, support, escalation, and postmortem processes for internal platform services.
- [ ] **USEQ-E2721BE2** — Design the platform so failure of an optional capability does not block unrelated delivery unnecessarily.
- [ ] **USEQ-859CA94E** — Test queue saturation, provider throttling, control-plane outage, identity failure, registry failure, network partition, and corrupt template scenarios.
- [ ] **USEQ-6FE68803** — Keep a supported manual or alternate path for genuinely critical operations when the primary platform is unavailable, and test it.
- [ ] **USEQ-99B27008** — Prevent a platform outage from removing access to the runbooks, credentials, artifacts, or communication needed for recovery.
- [ ] **USEQ-85373024** — Apply capacity planning, quotas, fairness, backpressure, and cost controls to shared capabilities.
- [ ] **USEQ-D8F0B0C5** — Publish maintenance windows, deprecation schedules, breaking changes, and reliability limitations.
- [ ] **USEQ-23DC86B8** — Keep rollback and version compatibility for platform interfaces, templates, policies, and generated assets.
- [ ] **USEQ-A96AE686** — Measure adoption alongside failure rate, support burden, lead time, stability, user satisfaction, and cost.
- [ ] **USEQ-C326D283** — Review whether platform centralization creates unacceptable correlated failure or organizational bottlenecks.

#### Balanced engineering measurement

- [ ] **USEQ-A84EEA72** — Use the current DORA delivery metrics as team or service-level outcome signals: change lead time, deployment frequency, failed deployment recovery time, change fail rate, and deployment rework rate.
- [ ] **USEQ-CAC06908** — Define events and denominators consistently before comparing delivery metrics over time.
- [ ] **USEQ-401CFDDA** — Measure metrics at a level where one flow of work and one operating context make interpretation meaningful.
- [ ] **USEQ-641B8692** — Do not use delivery metrics to rank individual engineers or teams with materially different contexts.
- [ ] **USEQ-EF22FA21** — Pair throughput metrics with stability, quality, security, privacy, accessibility, user outcome, well-being, and sustainability measures.
- [ ] **USEQ-6B03300D** — Use SPACE dimensions—satisfaction and well-being, performance, activity, communication and collaboration, and efficiency and flow—as a multidimensional lens rather than one composite score.
- [ ] **USEQ-D038E87B** — Use DevEx dimensions—feedback loops, cognitive load, and flow state—to identify causes rather than merely reporting sentiment.
- [ ] **USEQ-51FAFA2F** — Protect survey confidentiality and avoid collecting more employee data than needed.
- [ ] **USEQ-4BA2DDFB** — Publish metric definitions, data sources, exclusions, known biases, confidence, and intended decisions.
- [ ] **USEQ-23DD5632** — Audit telemetry for gaming, selection bias, missing work, bot activity, reclassification, and changed process semantics.
- [ ] **USEQ-AA2239E3** — Prefer trends and distributions to isolated averages and vanity totals.
- [ ] **USEQ-98265229** — Use qualitative inquiry to explain metric changes before acting on them.
- [ ] **USEQ-67683E93** — Retire metrics that no longer drive decisions or that create harmful incentives.
- [ ] **USEQ-078B9EDE** — Evaluate interventions with a baseline, expected mechanism, counter-metrics, observation period, and rollback criteria.

#### Engineering economics and business cases

- [ ] **USEQ-FD4E143E** — State the decision, alternatives, objectives, constraints, affected stakeholders, time horizon, and decision owner for every material investment.
- [ ] **USEQ-2A4B44EE** — Compare continuing, improving, replacing, buying, building, outsourcing, partnering, open sourcing, delaying, reducing scope, and stopping where plausible.
- [ ] **USEQ-01C9C8B4** — Estimate full lifecycle cost including discovery, acquisition, implementation, integration, migration, data conversion, security, privacy, accessibility, testing, operations, support, training, compliance, vendor management, downtime, exit, and retirement.
- [ ] **USEQ-FD955F18** — Include opportunity cost and the cost of delayed value, not only direct expenditure.
- [ ] **USEQ-01C3D941** — Identify benefits as measurable changes in user outcomes, revenue, risk, cost, capacity, quality, or strategic option value.
- [ ] **USEQ-606B2B76** — Avoid claiming time savings as cash savings unless the released capacity can be used or removed in a defined way.
- [ ] **USEQ-5906325E** — Model ranges and scenarios rather than one unjustifiably precise point estimate.
- [ ] **USEQ-B7F42705** — Discount future costs and benefits consistently when the decision horizon makes present-value analysis material.
- [ ] **USEQ-D09828A7** — Separate sunk costs from future decision-relevant costs.
- [ ] **USEQ-402B58A8** — Identify irreversible commitments, lock-in, switching cost, and the value of preserving options.
- [ ] **USEQ-8B570320** — Include expected loss from reliability, security, privacy, safety, legal, supplier, schedule, and adoption risks.
- [ ] **USEQ-BF1DBF5E** — Use sensitivity analysis to identify assumptions that dominate the decision.
- [ ] **USEQ-FB93CC82** — Define leading indicators and stop, pivot, expand, and review thresholds before committing the full investment.
- [ ] **USEQ-CEE325A1** — Update the business case with actual cost, benefit, adoption, quality, and risk evidence after delivery.
- [ ] **USEQ-8A77D825** — Record who accepted uncertainty and residual risk rather than hiding them in optimistic estimates.

#### Cost estimation and forecasting

- [ ] **USEQ-B9C50500** — Define the estimate's purpose, scope, decision date, confidence needs, currency, price basis, and included lifecycle stages.
- [ ] **USEQ-9AEFE830** — Create a technical and operational baseline describing intended capability, architecture, users, scale, data, integrations, quality targets, schedule, and constraints.
- [ ] **USEQ-73709F91** — Use a work breakdown that covers all material work and avoids double counting.
- [ ] **USEQ-A898AA51** — Document ground rules, assumptions, exclusions, dependencies, productivity factors, staffing model, inflation, exchange rates, and contingency treatment.
- [ ] **USEQ-F4AA620F** — Use relevant historical data and normalize it for size, complexity, technology, team maturity, quality obligations, and economic conditions.
- [ ] **USEQ-8596CFFB** — Select estimation methods appropriate to maturity and triangulate high-impact estimates using more than one method where possible.
- [ ] **USEQ-BFF44BBF** — Separate estimate uncertainty from management reserve and from intentionally deferred scope.
- [ ] **USEQ-A67CAF6F** — Model correlation among risks rather than summing independent optimistic contingencies blindly.
- [ ] **USEQ-8DBEC4BA** — Use probabilistic ranges for uncertain cost and schedule where consequence warrants it.
- [ ] **USEQ-7C70CDE8** — Perform sensitivity and risk analysis and identify the variables that most affect outcome.
- [ ] **USEQ-B74D3C1D** — Reconcile estimates with available capacity, hiring lead time, procurement time, and critical dependencies.
- [ ] **USEQ-5E3A5C8B** — Have material estimates independently reviewed for scope completeness, assumptions, data quality, method, arithmetic, bias, and credibility.
- [ ] **USEQ-DD3C0E57** — Document estimates so another qualified person can reproduce the logic and inputs.
- [ ] **USEQ-5E58EFDA** — Compare actuals with estimates at consistent milestones and explain variance by cause.
- [ ] **USEQ-0B9A4977** — Re-estimate after material scope, architecture, supplier, staffing, risk, quality, or schedule changes.
- [ ] **USEQ-E0A45D6A** — Do not convert an early exploratory estimate into a commitment without increasing evidence and explicitly changing its confidence classification.

#### Scheduling, dependencies, and delivery forecasting

- [ ] **USEQ-961773F5** — Represent all material work, milestones, external dependencies, approvals, migrations, training, operational readiness, and retirement activities in the delivery model.
- [ ] **USEQ-884A77F0** — Define task logic and dependencies rather than choosing dates independently.
- [ ] **USEQ-AB5BD18C** — Identify the critical and near-critical paths and monitor changes to them.
- [ ] **USEQ-3682C3C5** — Model resource availability, specialized skills, leave, support obligations, review capacity, environments, and supplier lead times.
- [ ] **USEQ-1B383C61** — Avoid planning people at full theoretical utilization; preserve capacity for coordination, incidents, learning, and uncertainty.
- [ ] **USEQ-4C01CA36** — Use ranges or probabilistic forecasts when task duration and scope are uncertain.
- [ ] **USEQ-82AC181E** — Track blocked time, queue time, rework, dependency age, and decision latency.
- [ ] **USEQ-163CD3CD** — Keep milestones tied to verifiable outcomes and evidence, not percentage-complete estimates.
- [ ] **USEQ-EEBE0AC0** — Reforecast using actual throughput and discovered work rather than preserving obsolete dates.
- [ ] **USEQ-65304AB1** — Make scope, time, cost, quality, and risk trade-offs explicit when constraints conflict.
- [ ] **USEQ-D1A09FF0** — Prevent schedule pressure from silently removing testing, accessibility, security, documentation, recovery, or maintenance work.
- [ ] **USEQ-57586986** — Define release, migration, rollback, stabilization, and support windows with adequate observation time.
- [ ] **USEQ-2E6ADA23** — Escalate dependencies before they reach the critical path rather than after they miss a date.
- [ ] **USEQ-E630B3A9** — Maintain a decision log explaining changes to committed scope and schedule.

#### Portfolio prioritization and investment governance

- [ ] **USEQ-0070EAA7** — Prioritize work using expected user and business value, strategic fit, urgency, risk reduction, learning value, dependency enablement, lifecycle cost, and capacity—not stakeholder volume alone.
- [ ] **USEQ-98ACEC61** — Reserve capacity for reliability, security, privacy, accessibility, platform health, maintenance, debt, and retirement.
- [ ] **USEQ-B54B7203** — Evaluate portfolio concentration by product, provider, technology, region, revenue source, data class, and key person.
- [ ] **USEQ-E49765AE** — Limit simultaneous initiatives to the organization's actual integration and change capacity.
- [ ] **USEQ-0D1C652B** — Fund discovery and risk retirement before making irreversible large commitments.
- [ ] **USEQ-3D020109** — Use staged investment with explicit evidence gates for uncertain initiatives.
- [ ] **USEQ-5C653F02** — Define success, failure, stop, pivot, and sunset criteria before launch.
- [ ] **USEQ-BBF30CFF** — Stop or reduce work when evidence no longer supports the expected outcome rather than protecting sunk cost.
- [ ] **USEQ-96ECBC49** — Review whether a local optimization shifts greater cost or risk to customers, support, operations, security, another team, or the future.
- [ ] **USEQ-0ACF35D3** — Include decommissioning and benefit-realization work in portfolio capacity.
- [ ] **USEQ-82D9539C** — Make priority decisions and trade-offs visible so teams are not forced to infer them from interruptions.
- [ ] **USEQ-BB0B4D99** — Review the portfolio after material market, regulatory, supplier, incident, capacity, or strategy changes.

#### Build, buy, open source, and supplier acquisition

- [ ] **USEQ-1B32D90B** — Define the capability and quality requirements before selecting a product, provider, framework, platform, or outsourcing arrangement.
- [ ] **USEQ-735862F8** — Evaluate functional fit, security, privacy, accessibility, reliability, performance, scalability, interoperability, data portability, operability, maintainability, support, and total cost.
- [ ] **USEQ-590D0692** — Assess supplier ownership, financial health, governance, development practices, vulnerability handling, roadmap, support model, and concentration risk.
- [ ] **USEQ-D6C76249** — Inventory sub-suppliers and critical external dependencies where risk warrants it.
- [ ] **USEQ-4C04425F** — Require evidence for security, accessibility, compliance, data location, deletion, backup, incident notification, and continuity claims.
- [ ] **USEQ-D6453C61** — Evaluate integration and migration cost using a representative proof, not only vendor demonstrations.
- [ ] **USEQ-56890F18** — Test contractual service levels against actual user objectives and define remedies, escalation, and evidence access.
- [ ] **USEQ-672B794C** — Retain rights to necessary data, configuration, logs, documentation, artifacts, and export formats.
- [ ] **USEQ-565684CB** — Define termination, migration assistance, credential revocation, data return, deletion, transition period, and exit cost before commitment.
- [ ] **USEQ-8754EEDE** — Avoid proprietary lock-in unless its value exceeds modeled switching risk and an accountable owner accepts it.
- [ ] **USEQ-34334A55** — Review open-source governance, license, maintainer health, release cadence, security response, transitive dependencies, and contribution strategy.
- [ ] **USEQ-2CBC3606** — Define ownership for forks, patches, vendored code, unsupported versions, and upstream contributions.
- [ ] **USEQ-41AD9DE5** — Use escrow, source access, alternate providers, or internal capability where loss of a supplier would create unacceptable continuity risk.
- [ ] **USEQ-FE5899D5** — Continuously monitor supplier changes that can alter price, terms, data use, functionality, security, support, or geographic availability.

#### People, skills, sustainability, and organizational resilience

- [ ] **USEQ-D91347DE** — Maintain a capability map for product, domain, architecture, code, data, security, privacy, accessibility, operations, and supplier knowledge.
- [ ] **USEQ-8394400E** — Ensure every critical service and recovery procedure has more than one qualified person and accessible documentation.
- [ ] **USEQ-8D1332CC** — Plan training before introducing technology or controls that teams are expected to operate safely.
- [ ] **USEQ-C1AF4BE4** — Provide mentoring, review, communities of practice, and protected learning time rather than relying only on courses.
- [ ] **USEQ-457EF747** — Design roles and on-call expectations that are sustainable and compatible with health, leave, and local obligations.
- [ ] **USEQ-44E36BD7** — Track burnout, interruption, after-hours load, incident frequency, support burden, and chronic understaffing as system risks.
- [ ] **USEQ-A48592E0** — Do not use heroics as a capacity plan or performance expectation.
- [ ] **USEQ-CCCA9A78** — Include contractors and suppliers in necessary knowledge transfer, documentation, incident, security, and exit processes.
- [ ] **USEQ-1479A1E6** — Ensure performance management does not reward risky volume, hidden work, gate bypass, knowledge hoarding, or incident concealment.
- [ ] **USEQ-51E72868** — Create psychologically safe routes to raise quality, security, ethics, accessibility, and schedule concerns.
- [ ] **USEQ-9EDBBD06** — Protect whistleblowing and escalation from retaliation.
- [ ] **USEQ-447508A8** — Plan succession, offboarding, access removal, artifact transfer, and ownership reassignment.
- [ ] **USEQ-59C8FAD1** — Review whether reorganizations create orphaned services, broken escalation paths, duplicated platforms, or unsupported commitments.

#### Engineering-system release blockers

- [ ] **USEQ-BF8D8060** — Do not declare the engineering system ready when a routine critical workflow depends on undocumented manual intervention or one person's availability.
- [ ] **USEQ-C67086D4** — Do not introduce a mandatory platform path that cannot meet required reliability, security, accessibility, and recovery objectives.
- [ ] **USEQ-8BBAAB6D** — Do not use individual productivity rankings, raw activity counts, or one delivery metric as a quality or performance verdict.
- [ ] **USEQ-4156513F** — Do not commit a material investment whose scope, lifecycle cost, dominant assumptions, alternatives, and exit path are unknown.
- [ ] **USEQ-3F0AD88C** — Do not commit to a date that omits critical dependencies, quality work, migration, stabilization, or support capacity.
- [ ] **USEQ-A70E005C** — Do not buy or build a critical capability without an accountable owner, operating model, support plan, and retirement path.
- [ ] **USEQ-8BA058A0** — Do not accept chronic broken builds, flaky tests, access delays, deployment tickets, or environment instability as normal developer experience.
- [ ] **USEQ-31618E9F** — Do not force migration away from a supported path without consumer inventory, compatibility evidence, communication, and recovery.
- [ ] **USEQ-92FFBA7B** — Do not use platform adoption as success evidence when user outcomes, stability, throughput, satisfaction, or total cost deteriorate.
- [ ] **USEQ-D933D63C** — Do not proceed when workload, on-call, knowledge concentration, or staffing makes safe operation unsustainable.

## Platform Engineering, Developer Experience, and Engineering Effectiveness

_Consolidated from `final consolidated corpus/07-delivery-configuration-platform-developer-experience-cicd.md#Platform Engineering, Developer Experience, and Engineering Effectiveness`; 238 non-duplicative controls._

### Platform and developer-experience strategy

- [ ] **USEQ-2837A62E** — Treat internal engineering capabilities as products with identifiable users, outcomes, owners, roadmaps, support, and lifecycle decisions.
- [ ] **USEQ-2087E7B1** — Define which developer, operator, analyst, data, security, support, and partner personas the platform serves.
- [ ] **USEQ-7872D116** — Research user workflows and friction directly rather than inferring needs only from leadership or tool owners.
- [ ] **USEQ-2ADEDFF1** — Define the outcomes the platform should improve, such as safe delivery, feedback speed, reliability, compliance, onboarding, or cognitive load.
- [ ] **USEQ-9624E443** — Do not create a platform merely to centralize control or standardize tools without demonstrated user value.
- [ ] **USEQ-C98C9BC6** — Distinguish mandatory organization controls from optional accelerators and clearly explain why each mandatory control exists.
- [ ] **USEQ-F8CE1AED** — Define what the platform owns, what product teams own, and what suppliers own.
- [ ] **USEQ-84F487EA** — Publish a service boundary, supported use cases, non-goals, dependencies, and escalation model.
- [ ] **USEQ-D88D34D6** — Prioritize capabilities using user impact, delivery risk, security, reliability, cost, and total organizational effort.
- [ ] **USEQ-9D173BAA** — Fund platform maintenance, support, security, migration, and deprecation—not only initial construction.
- [ ] **USEQ-2BC4E867** — Use an incremental “thinnest viable platform” or equivalent approach instead of attempting to build every capability before learning from use.
- [ ] **USEQ-624E74B9** — Keep platform adoption voluntary where risk permits and earn adoption through value; where controls are mandatory, make the compliant path the easiest path.
- [ ] **USEQ-F61BB327** — Review whether the platform increases autonomy or merely moves tickets and bottlenecks to a different team.
- [ ] **USEQ-BD39D438** — Measure and manage platform-induced lock-in, concentration risk, and loss of local expertise.

### Service catalog, ownership, and discoverability

- [ ] **USEQ-75B23248** — Maintain a searchable catalog of services, applications, libraries, data products, pipelines, environments, infrastructure, owners, dependencies, criticality, and lifecycle status.
- [ ] **USEQ-0537910D** — Assign an accountable owner and operational contact to every production component.
- [ ] **USEQ-D80B97CE** — Record source repository, documentation, runbook, dashboards, alerts, SLOs, data classifications, APIs, dependencies, and deployment location for each catalog entry.
- [ ] **USEQ-010C8741** — Identify orphaned, duplicate, abandoned, unsupported, and unowned components.
- [ ] **USEQ-312DD67A** — Keep catalog data synchronized automatically where reliable and review manually supplied metadata for staleness.
- [ ] **USEQ-ACF53040** — Distinguish authoritative metadata from discovered or inferred metadata.
- [ ] **USEQ-DEDBCA46** — Expose dependency and ownership relationships without revealing sensitive topology to unauthorized users.
- [ ] **USEQ-86095E9D** — Integrate catalog ownership with incident routing, vulnerability notifications, deprecation notices, and cost allocation.
- [ ] **USEQ-C2E31CA3** — Define registration and retirement workflows so new and deleted components do not bypass inventory.
- [ ] **USEQ-8972276D** — Provide stable identifiers for catalog entities across renames, moves, and platform migrations.
- [ ] **USEQ-489B6ACB** — Allow teams to discover approved components, examples, policies, and reusable capabilities without tribal knowledge.
- [ ] **USEQ-6385F8F8** — Monitor coverage and freshness of the catalog as an operational control.
- [ ] **USEQ-5FDA8D4E** — Do not treat catalog presence as evidence that documentation, support, or ownership is actually effective.
- [ ] **USEQ-FB32D303** — Keep emergency contact information available when the normal catalog or identity platform is unavailable.

### Golden paths, paved roads, templates, and reusable capabilities

- [ ] **USEQ-07FDA5F4** — Provide recommended paths for common workflows using secure, observable, supportable, and maintained defaults.
- [ ] **USEQ-4B12B7BF** — Design paths around complete user outcomes such as creating, testing, deploying, operating, and retiring a service—not isolated tool installation.
- [ ] **USEQ-2603D264** — Keep recommended paths versioned, documented, testable, and reproducible.
- [ ] **USEQ-2B0F618E** — Generate the minimum necessary starter material and avoid creating large opaque templates teams cannot understand.
- [ ] **USEQ-1FA6A3A8** — Ensure generated projects have clear ownership, dependency provenance, update paths, tests, and removal instructions.
- [ ] **USEQ-372F71E0** — Provide examples that follow current standards and are continuously tested.
- [ ] **USEQ-69A27DFC** — Use shared libraries and components only when their support, compatibility, and change model are clear.
- [ ] **USEQ-2BCB3A9A** — Avoid forcing unrelated products into one architecture when risk and requirements differ.
- [ ] **USEQ-710C661D** — Provide reviewed escape hatches for legitimate exceptions, with explicit ownership and risk treatment.
- [ ] **USEQ-8CF50CBD** — Ensure escaping a recommended path does not silently remove security, privacy, accessibility, reliability, or audit controls.
- [ ] **USEQ-ADE1D409** — Track adoption, abandonment, support burden, defect rates, and user satisfaction for each path.
- [ ] **USEQ-B3750A1C** — Deprecate obsolete paths with migration tooling, communication, and sufficient support windows.
- [ ] **USEQ-A33D1C6F** — Allow teams to inspect and understand generated configuration and infrastructure.
- [ ] **USEQ-9BCE1201** — Prevent templates from embedding secrets, obsolete dependencies, insecure defaults, or organization-specific identifiers.

### Self-service and workflow automation

- [ ] **USEQ-2E72ABDB** — Enable authorized users to perform frequent low-risk actions without waiting for manual tickets.
- [ ] **USEQ-2B4C917C** — Automate policy checks and approvals that are deterministic; reserve human approval for judgment, exceptions, and high-impact risk.
- [ ] **USEQ-4B11585B** — Use clear inputs, previews, validation, cost implications, and consequences before self-service actions execute.
- [ ] **USEQ-E31ED114** — Make destructive and high-impact operations reversible, confirmable, or subject to additional authorization as appropriate.
- [ ] **USEQ-8CBDC7EA** — Keep self-service workflows idempotent and safe to retry.
- [ ] **USEQ-C2F428A8** — Provide progress, failure, rollback, and recovery information for long-running operations.
- [ ] **USEQ-3B63475F** — Do not hide partial failure or leave ambiguous resource ownership.
- [ ] **USEQ-5B6344BA** — Use scoped identities and least privilege for automation.
- [ ] **USEQ-B737F580** — Record actor, request, policy decision, artifact, configuration, outcome, and affected resources for auditable actions.
- [ ] **USEQ-6325634F** — Apply quotas, rate limits, budgets, and guardrails to prevent accidental or malicious resource exhaustion.
- [ ] **USEQ-4B82601E** — Prevent users from provisioning unsupported versions or bypassing required lifecycle controls.
- [ ] **USEQ-DBFE6DD8** — Provide an escalation path when automation cannot satisfy a valid request.
- [ ] **USEQ-CB367A5D** — Measure whether self-service reduces total wait time and toil rather than only ticket count.
- [ ] **USEQ-A3EF0CB5** — Test automation against duplicate requests, cancellation, dependency failure, and platform outages.

### Environment creation, onboarding, and reproducibility

- [ ] **USEQ-587143E3** — Provide a documented, testable path from a new authorized workstation or workspace to a functioning development environment.
- [ ] **USEQ-DBC5BFD6** — Minimize required manual installation, local machine mutation, and hidden prerequisites.
- [ ] **USEQ-0B335E6B** — Pin or verify toolchain, runtime, dependency, and generator versions sufficiently for reproducibility.
- [ ] **USEQ-AB6CABAF** — Keep setup scripts idempotent and safe to rerun.
- [ ] **USEQ-1629715E** — Validate prerequisites early and produce actionable errors.
- [ ] **USEQ-1AAC867A** — Provide sanitized representative test data and avoid copying production data by default.
- [ ] **USEQ-53AE4D3E** — Use separate identities, credentials, data, endpoints, and side-effect controls for non-production environments.
- [ ] **USEQ-D5E12F9F** — Make environment reset and cleanup reliable.
- [ ] **USEQ-8C2638B8** — Ensure local, remote, containerized, virtual, or hosted development options meet security and accessibility requirements.
- [ ] **USEQ-78829D95** — Provide an alternative when a developer cannot use a specific hardware platform, biometric method, visual tool, or network location.
- [ ] **USEQ-1B1561EC** — Document meaningful differences between local, test, staging, and production behavior.
- [ ] **USEQ-C59AAC06** — Measure onboarding task success and time without turning individual speed into a performance score.
- [ ] **USEQ-547DB4A5** — Test onboarding instructions with people unfamiliar with the system.
- [ ] **USEQ-697FCFA9** — Remove obsolete setup paths and prevent conflicting documentation from remaining discoverable.

### Feedback loops, builds, tests, and review latency

- [ ] **USEQ-B739ED54** — Keep the fastest relevant validation close to the developer and defer expensive checks without losing release assurance.
- [ ] **USEQ-9450BC5E** — Define target feedback times for formatting, compilation, static analysis, unit tests, integration tests, review, builds, and deployments.
- [ ] **USEQ-03172925** — Monitor feedback latency percentiles and queue time rather than averages alone.
- [ ] **USEQ-65A42809** — Prioritize deterministic, actionable failures over large volumes of noisy output.
- [ ] **USEQ-62AAA93F** — Fail early on configuration, dependency, schema, and policy errors that would invalidate later work.
- [ ] **USEQ-F9B57D3D** — Parallelize independent checks without introducing nondeterminism or excessive resource cost.
- [ ] **USEQ-FEDEFC91** — Use caching only when cache keys, invalidation, isolation, provenance, and poisoning risks are controlled.
- [ ] **USEQ-0F71EDF0** — Make it easy to reproduce a CI failure locally or in an isolated diagnostic environment.
- [ ] **USEQ-31021B4A** — Provide logs and artifacts sufficient to diagnose failures without broad privileged access.
- [ ] **USEQ-6E3F668E** — Track flaky tests, intermittent infrastructure failures, canceled jobs, and reruns separately from product defects.
- [ ] **USEQ-725337F7** — Prevent teams from normalizing repeated reruns as a valid delivery process.
- [ ] **USEQ-E454B7C4** — Keep code-review queues visible and balance reviewer load and expertise.
- [ ] **USEQ-66533E37** — Encourage small batches that can be reviewed, tested, deployed, and reversed safely.
- [ ] **USEQ-5A81540A** — Measure end-to-end feedback loops, including waiting for environments, data, approvals, reviews, and external teams.

### Internal APIs, contracts, and integration experience

- [ ] **USEQ-90D7EA9F** — Treat platform interfaces as supported products with explicit contracts, versioning, ownership, and compatibility policies.
- [ ] **USEQ-A87EE218** — Publish schemas, examples, error semantics, rate limits, quotas, authentication, authorization, idempotency, and lifecycle states.
- [ ] **USEQ-C38502CA** — Provide machine-readable interfaces where automation is expected and human-readable interfaces where judgment is required.
- [ ] **USEQ-CDE257DC** — Keep command-line, user-interface, API, and configuration behavior semantically consistent.
- [ ] **USEQ-15FDD113** — Use stable identifiers and avoid exposing volatile internal implementation details.
- [ ] **USEQ-B0E23E13** — Validate inputs before creating partial resources.
- [ ] **USEQ-CF3EEC00** — Provide preview, plan, or dry-run behavior for consequential changes where practical.
- [ ] **USEQ-1D76E3B6** — Define asynchronous operation status, cancellation, retry, timeout, and reconciliation behavior.
- [ ] **USEQ-2702E53D** — Provide contract tests or conformance suites for clients and providers.
- [ ] **USEQ-74089988** — Use deprecation notices that reach actual consumers and include migration guidance and deadlines.
- [ ] **USEQ-6DDF8C21** — Monitor use of deprecated interfaces and identify unknown consumers before removal.
- [ ] **USEQ-B74367BE** — Keep error messages actionable and linked to current documentation.
- [ ] **USEQ-6707A6FE** — Do not require consumers to scrape logs or user interfaces when a stable interface is appropriate.
- [ ] **USEQ-4F9EA2A9** — Test interface accessibility and localization when humans interact with it.

### Guardrails, policy as code, and secure defaults

- [ ] **USEQ-48758704** — Encode repeatable mandatory controls in versioned, reviewed, tested policy where practical.
- [ ] **USEQ-A3C4E315** — Keep policy source, rationale, owner, scope, severity, remediation guidance, and exception process visible.
- [ ] **USEQ-91D01B47** — Test policy rules against allowed and prohibited examples before enforcement.
- [ ] **USEQ-CFDE8E86** — Roll out high-impact policy changes progressively and measure false positives and operational impact.
- [ ] **USEQ-DA1E0181** — Fail closed for critical security and authorization controls unless a documented safety analysis requires another behavior.
- [ ] **USEQ-C663C195** — Use warning or advisory modes only when the residual risk is explicitly accepted and monitored.
- [ ] **USEQ-A04912A7** — Provide safe defaults for identity, network access, encryption, logging, backups, resource limits, and exposure.
- [ ] **USEQ-DF36238E** — Do not make essential security capabilities a premium or optional path when users reasonably expect safe defaults.
- [ ] **USEQ-7B592261** — Allow scoped time-bounded exceptions with owner, evidence, compensating controls, and automatic expiry.
- [ ] **USEQ-C6CEA1DC** — Prevent policy bypass through alternate tools, direct provider access, unmanaged accounts, or emergency paths.
- [ ] **USEQ-4F850A3C** — Log and review break-glass and exception use.
- [ ] **USEQ-B8027624** — Keep policy evaluation available and performant enough not to create pressure for bypass.
- [ ] **USEQ-BBB80D01** — Provide automated remediation where it is safe and transparent.
- [ ] **USEQ-49857976** — Review whether guardrails continue to address real risk rather than obsolete architecture.

### Platform reliability, SLOs, observability, and support

- [ ] **USEQ-BA808893** — Define user-centered SLIs and SLOs for critical platform journeys such as build, test, provision, deploy, retrieve secrets, and diagnose incidents.
- [ ] **USEQ-791F51DA** — Include correctness, latency, availability, freshness, capacity, and support response where relevant.
- [ ] **USEQ-8121C780** — Ensure platform objectives are compatible with the reliability objectives of dependent products.
- [ ] **USEQ-B52AE8A4** — Monitor platform dependencies, quotas, certificate expiry, identity, DNS, artifact storage, and control planes.
- [ ] **USEQ-87EFFA25** — Provide synthetic checks for critical self-service journeys.
- [ ] **USEQ-263E0518** — Expose service status and known incidents to platform users.
- [ ] **USEQ-2ADCF482** — Use error budgets or an equivalent policy to balance feature work with reliability work.
- [ ] **USEQ-40DBA950** — Provide runbooks and on-call ownership for platform components whose failure blocks delivery or recovery.
- [ ] **USEQ-B6288BBC** — Ensure support personnel can diagnose problems without inappropriate access to product secrets or customer data.
- [ ] **USEQ-62FDFE99** — Define support channels, hours, severity, escalation, and expected response.
- [ ] **USEQ-D2E7BAC3** — Capture repeated support requests as product or documentation defects.
- [ ] **USEQ-98754BDB** — Test platform disaster recovery and restore, including source, configuration, registries, metadata, identity, and signing material.
- [ ] **USEQ-0B9A6478** — Ensure product teams can continue critical incident response if the normal platform is unavailable.
- [ ] **USEQ-F7346E05** — Publish maintenance, degradation, and emergency change procedures.

### Cognitive load, flow, accessibility, and inclusive tooling

- [ ] **USEQ-58A2631E** — Measure whether engineers can understand and complete common tasks without memorizing undocumented system details.
- [ ] **USEQ-C802111B** — Reduce unnecessary tool switching, duplicate data entry, repeated authentication, and manual translation between systems.
- [ ] **USEQ-6924F2DC** — Provide clear defaults while revealing advanced options progressively.
- [ ] **USEQ-C31E7F6F** — Use consistent terminology, navigation, status, and error patterns across engineering tools.
- [ ] **USEQ-09A914D6** — Keep documentation searchable from the point of need.
- [ ] **USEQ-A3E724DF** — Protect focus time and minimize unnecessary notifications, approvals, meetings, and interrupt-driven work.
- [ ] **USEQ-6A7F9660** — Design tooling so users can recover from errors without recreating environments or waiting for privileged intervention.
- [ ] **USEQ-14E7881C** — Make internal tools keyboard accessible and compatible with relevant assistive technologies.
- [ ] **USEQ-91534D82** — Provide sufficient contrast, scalable text, non-color cues, captions, and accessible documents in engineering systems.
- [ ] **USEQ-7C591DBF** — Ensure remote, distributed, part-time, and differently located team members have equivalent access to workflows and information.
- [ ] **USEQ-CE9D204D** — Avoid relying on informal spoken knowledge, office proximity, or one time zone for essential decisions.
- [ ] **USEQ-DC4B7B0A** — Provide psychologically safe ways to report friction, errors, and harmful metrics.
- [ ] **USEQ-DAD6FF8C** — Do not interpret struggle with inaccessible or poorly designed tools as individual underperformance.
- [ ] **USEQ-C7C144EB** — Include cognitive load, feedback loops, and flow state in improvement research.

### Documentation, learning, and knowledge continuity

- [ ] **USEQ-4275296E** — Provide task-oriented documentation for common workflows and reference documentation for interfaces and policies.
- [ ] **USEQ-EB41F750** — Keep one discoverable authoritative source for each instruction or contract.
- [ ] **USEQ-01E6CBBE** — Generate reference material from source where reliable while retaining reviewed explanations and examples.
- [ ] **USEQ-3B0A5ED4** — Test documentation examples and commands continuously where practical.
- [ ] **USEQ-891BB07B** — Label versions, prerequisites, scope, owner, last review date, and known limitations.
- [ ] **USEQ-726AC92D** — Archive or redirect obsolete documentation so search does not lead users to unsupported paths.
- [ ] **USEQ-D916210E** — Provide onboarding journeys, architecture overviews, troubleshooting guides, and conceptual explanations.
- [ ] **USEQ-63D39C0B** — Record material decisions and trade-offs in durable decision records.
- [ ] **USEQ-33FAD702** — Ensure runbooks identify safe and destructive actions clearly.
- [ ] **USEQ-24B646C7** — Provide learning environments that cannot affect production or real users.
- [ ] **USEQ-D82CEC22** — Use communities of practice, office hours, and support channels to supplement—not replace—durable documentation.
- [ ] **USEQ-2B54E2F5** — Measure documentation task success, search failure, repeated questions, and stale-page reports.
- [ ] **USEQ-983B5F26** — Plan knowledge transfer before owner changes, organizational moves, or supplier exit.
- [ ] **USEQ-4BDCFFCF** — Keep emergency and recovery documentation available outside the systems it is needed to restore.

### Delivery performance and developer-experience measurement

- [ ] **USEQ-CB8F61D4** — Use multiple outcome and experience measures; no single metric represents developer productivity.
- [ ] **USEQ-7AA62376** — Use DORA delivery metrics at an appropriate service or team level to understand throughput and instability.
- [ ] **USEQ-B0B6D105** — As of the 2026 DORA model, consider change lead time, deployment frequency, failed deployment recovery time, change fail rate, and deployment rework rate together.
- [ ] **USEQ-83FDD06F** — Use SPACE dimensions to balance satisfaction and well-being, performance, activity, communication and collaboration, and efficiency and flow.
- [ ] **USEQ-4586EC60** — Use DevEx dimensions such as feedback loops, cognitive load, and flow state to diagnose friction.
- [ ] **USEQ-F58BC741** — Combine quantitative telemetry with confidential qualitative research.
- [ ] **USEQ-20223EB5** — Measure systems and team conditions rather than ranking individuals by commits, lines, tickets, reviews, hours, or tool activity.
- [ ] **USEQ-95ED4C79** — Do not use activity metrics as direct proxies for value, quality, difficulty, or individual contribution.
- [ ] **USEQ-D208E392** — Expect measures to change behavior and assess gaming, fear, inequity, and surveillance risk before collection.
- [ ] **USEQ-2184C875** — Collect only data necessary for the stated improvement purpose and apply privacy, access, retention, and transparency controls.
- [ ] **USEQ-F3060494** — Segment metrics by workflow and context without exposing or penalizing small groups.
- [ ] **USEQ-5355662B** — Look for trade-offs: faster throughput can coexist with more rework, instability, burnout, or support burden.
- [ ] **USEQ-5A698497** — Use trends and experiments to evaluate interventions rather than claiming causality from correlation alone.
- [ ] **USEQ-CF1326C5** — Publish metric definitions, limitations, owners, and decision rules.

### Continuous improvement and value-stream management

- [ ] **USEQ-2F8B10C6** — Map the end-to-end path from idea or defect to validated user outcome, including queues, handoffs, approvals, and rework.
- [ ] **USEQ-B8A5DC21** — Measure touch time, wait time, blocked time, rework, abandonment, and failure demand.
- [ ] **USEQ-FF953623** — Improve the largest system constraint rather than optimizing local activity that increases downstream queues.
- [ ] **USEQ-6A5A19D4** — Work in small batches to shorten feedback and reduce integration and rollback risk.
- [ ] **USEQ-A13C5E1E** — Use experiments with baseline, hypothesis, expected effect, guardrails, and review date.
- [ ] **USEQ-B517A398** — Include product quality, user value, reliability, security, accessibility, privacy, and team well-being as guardrails.
- [ ] **USEQ-9C9CBB20** — Review whether standardization removes friction or prevents teams from solving distinct problems.
- [ ] **USEQ-E325A12A** — Automate repetitive deterministic toil and remove the source of recurring toil where possible.
- [ ] **USEQ-7ABD2C5B** — Keep manual judgment where ambiguity, ethics, safety, or contextual risk requires it.
- [ ] **USEQ-DC44A0E3** — Prioritize improvements from observed friction, incident lessons, support demand, and user research.
- [ ] **USEQ-BD769EB9** — Share successful patterns without declaring them universal before context is understood.
- [ ] **USEQ-F0AE0BFF** — Retire metrics and processes that no longer lead to useful decisions.
- [ ] **USEQ-188FE891** — Review improvement effects after enough time to observe adoption and secondary consequences.
- [ ] **USEQ-441F16D2** — Capture lessons in platform capabilities, policies, documentation, and organization design.

### Cost, capacity, and sustainability of engineering systems

- [ ] **USEQ-F2A1D89A** — Attribute major platform costs to capabilities and outcomes without creating incentives to hide necessary use.
- [ ] **USEQ-D68D0A9D** — Forecast build, test, storage, network, runner, artifact, logging, and environment demand.
- [ ] **USEQ-1A4066FD** — Use quotas and budgets that protect shared resources while allowing legitimate bursts and incident response.
- [ ] **USEQ-49E28DF3** — Monitor queue saturation, wait time, cache effectiveness, idle environments, abandoned resources, and data growth.
- [ ] **USEQ-2C2CAF70** — Automatically expire or suspend unused ephemeral resources after a visible grace period.
- [ ] **USEQ-990ED243** — Keep retention aligned with legal, audit, diagnostic, and reproducibility needs.
- [ ] **USEQ-344C24AA** — Minimize duplicate artifact, dependency, test-data, and log storage while preserving integrity and provenance.
- [ ] **USEQ-28FE7C91** — Choose efficient architectures and workload scheduling without shifting hidden cost to developers or product reliability.
- [ ] **USEQ-3B939CA0** — Measure environmental and resource impact where material and avoid unsupported sustainability claims.
- [ ] **USEQ-69D55612** — Plan capacity for organizational growth, launch spikes, regional failure, and supplier quotas.
- [ ] **USEQ-39710937** — Ensure cost controls cannot silently delete evidence, disable security checks, or block emergency recovery.
- [ ] **USEQ-514E0CD0** — Review supplier pricing, egress, licensing, and exit costs before deep integration.
- [ ] **USEQ-07373DFE** — Provide teams with actionable cost visibility and safe optimization guidance.
- [ ] **USEQ-BBCAACB3** — Validate cost improvements against delivery, quality, reliability, and developer-experience guardrails.

### Portability, interoperability, and escape paths

- [ ] **USEQ-0E1655B8** — Document dependencies on proprietary APIs, formats, identity, networking, build services, registries, and control planes.
- [ ] **USEQ-8A0CF5CF** — Use open or well-documented interfaces where they materially reduce lock-in and improve interoperability.
- [ ] **USEQ-4018176E** — Keep source, configuration, metadata, artifacts, logs, and data exportable in usable forms.
- [ ] **USEQ-91A133AE** — Test restore or reconstruction outside the primary platform when continuity objectives require it.
- [ ] **USEQ-97037A86** — Define which platform abstractions intentionally hide provider details and which provider capabilities remain visible.
- [ ] **USEQ-49BA0A1C** — Avoid lowest-common-denominator abstractions that remove valuable capability without a real portability requirement.
- [ ] **USEQ-54067417** — Provide escape hatches that preserve essential security, audit, and ownership controls.
- [ ] **USEQ-AC2E75BF** — Document the operational burden transferred to teams that leave the recommended path.
- [ ] **USEQ-F300483C** — Keep external interfaces versioned independently of internal implementation where practical.
- [ ] **USEQ-47566BF2** — Maintain replacement and exit plans for critical suppliers and unsupported internal tools.
- [ ] **USEQ-30CCD67D** — Assess whether platform failure can prevent access to source, builds, deployment, observability, or recovery simultaneously.
- [ ] **USEQ-9934131E** — Prevent one identity or control-plane dependency from becoming the only emergency path.
- [ ] **USEQ-6EAB87B3** — Rehearse migration of at least the most critical artifacts and metadata.
- [ ] **USEQ-746C0E26** — Review portability requirements based on business continuity and bargaining risk, not ideology.

### Lifecycle, deprecation, and platform retirement

- [ ] **USEQ-D19BE27B** — Publish lifecycle states such as experimental, supported, deprecated, maintenance-only, and retired.
- [ ] **USEQ-9372CB8C** — Define support and compatibility expectations for each state.
- [ ] **USEQ-A7ABBFA6** — Identify actual consumers before deprecation and communicate through channels they use.
- [ ] **USEQ-8D189C10** — Provide migration tooling, documentation, examples, validation, and support.
- [ ] **USEQ-669B7C08** — Monitor remaining use and unknown consumers.
- [ ] **USEQ-6EE88A72** — Avoid indefinite dual operation without explicit cost and risk acceptance.
- [ ] **USEQ-68A53B3C** — Preserve required artifacts, evidence, audit history, and data during retirement.
- [ ] **USEQ-5AD5C402** — Revoke credentials, identities, network access, secrets, signing keys, webhooks, and automation associated with retired capabilities.
- [ ] **USEQ-45B505EF** — Delete or archive obsolete resources according to retention and privacy requirements.
- [ ] **USEQ-FF09ED1C** — Remove old paths from templates, documentation, policy, onboarding, and catalogs.
- [ ] **USEQ-7CBE6D66** — Confirm that retirement does not break incident recovery or historical rebuilds.
- [ ] **USEQ-36237DAD** — Review lessons from adoption, support, and retirement before replacing the capability.
- [ ] **USEQ-69BACE4D** — Keep a rollback window until migration validation is complete.
- [ ] **USEQ-C42204DA** — Reassign ownership explicitly when platform responsibilities move between teams or suppliers.

### Platform and developer-experience release blockers

- [ ] **USEQ-CB04A4EC** — Block rollout when the platform cannot identify owners, affected users, supported scope, or recovery responsibility.
- [ ] **USEQ-0CD5999D** — Block rollout when a mandatory path is less secure or less reliable than the path it replaces.
- [ ] **USEQ-A882ABAF** — Block rollout when self-service can create irreversible, privileged, publicly exposed, or unowned resources without adequate controls.
- [ ] **USEQ-24631885** — Block rollout when critical workflows cannot be diagnosed or recovered without the unavailable platform itself.
- [ ] **USEQ-742B58F7** — Block rollout when generated templates contain unsupported, vulnerable, unlicensed, or unowned components.
- [ ] **USEQ-7A5FFA5E** — Block rollout when policy enforcement produces material false positives without a working exception and support process.
- [ ] **USEQ-FF503F0F** — Block rollout when adoption is measured only by mandate or activity and user outcomes are unknown.
- [ ] **USEQ-DA19EDF9** — Block rollout when accessibility barriers prevent authorized engineers from completing essential tasks.
- [ ] **USEQ-9863D445** — Block rollout when telemetry constitutes undisclosed individual surveillance or creates harmful performance incentives.
- [ ] **USEQ-DE6F930B** — Require evidence from representative teams, including first-time users and high-risk workflows.
- [ ] **USEQ-34527B32** — Require documented SLOs, on-call ownership, rollback, support, migration, and deprecation plans for critical platform services.
- [ ] **USEQ-5EBA2B57** — Confirm that product teams retain enough knowledge and access to operate and recover their systems.
- [ ] **USEQ-7EA98BC7** — Confirm that cost, capacity, identity, supplier, and concentration risks are understood and accepted.
- [ ] **USEQ-08C12DD5** — Confirm that the platform improves the total value stream rather than only the platform team’s local metrics.

## Final Gap Closure — IT Asset Management, Configuration Management, and Toolchain Trust

_Consolidated from `final consolidated corpus/07-delivery-configuration-platform-developer-experience-cicd.md#Final Gap Closure — IT Asset Management, Configuration Management, and Toolchain Trust`; 127 non-duplicative controls._

### IT and software asset governance

- [ ] **USEQ-47D22600** — Define the assets governed across hardware, software, cloud, data, identities, certificates, domains, repositories, build systems, services, licenses, documentation, and suppliers.
- [ ] **USEQ-5EFC0327** — Assign an accountable owner, custodian, business purpose, criticality, lifecycle state, location, and support status to each material asset.
- [ ] **USEQ-6327CC41** — Maintain one authoritative or reconciled asset inventory with unique identifiers.
- [ ] **USEQ-F339AB78** — Discover assets continuously enough to detect unknown, abandoned, duplicated, shadow, and unauthorized resources.
- [ ] **USEQ-4408FB18** — Reconcile technical discovery with financial, procurement, identity, network, configuration, and service records.
- [ ] **USEQ-60452413** — Record relationships among assets, services, data, owners, dependencies, contracts, and customers.
- [ ] **USEQ-2FC61C58** — Classify assets by confidentiality, integrity, availability, privacy, safety, financial, and legal impact.
- [ ] **USEQ-BBF0269D** — Define lifecycle states and required controls from request and acquisition through deployment, operation, transfer, retirement, and disposal.
- [ ] **USEQ-E3B57237** — Prevent an asset from entering production without ownership, support, monitoring, security, backup, and retirement expectations.
- [ ] **USEQ-6DD907B9** — Review inventory completeness and accuracy periodically and after major organizational or platform change.
- [ ] **USEQ-D7D10F08** — Protect the asset inventory because it contains sensitive architecture and ownership information.
- [ ] **USEQ-A91E76D3** — Preserve historical asset and ownership state for incident, audit, cost, and lifecycle analysis.

### Acquisition, entitlement, and licensing

- [ ] **USEQ-C2BEC642** — Link acquired assets to approved requirements, funding, supplier, contract, entitlement, and accountable owner.
- [ ] **USEQ-BDA3E93B** — Verify license, subscription, support, geographic, user, device, capacity, environment, and redistribution rights before use.
- [ ] **USEQ-9F40C6C0** — Track renewal, notice, termination, price-adjustment, true-up, support, and end-of-life dates.
- [ ] **USEQ-90D325A9** — Prevent unauthorized software, services, plugins, packages, and devices from becoming production dependencies.
- [ ] **USEQ-C68F693C** — Identify free, trial, community, and evaluation services that lack production rights or support.
- [ ] **USEQ-A87D5AB1** — Avoid collecting more licenses or cloud capacity than can be governed responsibly.
- [ ] **USEQ-F929714B** — Detect under-licensing, over-licensing, duplicate subscriptions, dormant seats, and unused reserved capacity.
- [ ] **USEQ-C3E997C5** — Record supplier and product identifiers sufficiently to distinguish similarly named assets.
- [ ] **USEQ-20E81856** — Retain proof of entitlement and contractual terms for the required period.
- [ ] **USEQ-0B64BF62** — Ensure acquisition terms cover continuity, security updates, data return, audit, and exit where material.

### Asset assignment, custody, and use

- [ ] **USEQ-D738D8F6** — Assign individual custody for portable assets and privileged credentials where appropriate.
- [ ] **USEQ-40C843EB** — Record transfers, loans, repairs, loss, return, reassignment, and disposal.
- [ ] **USEQ-7D4EF0D8** — Require users to acknowledge handling responsibilities proportionate to sensitivity.
- [ ] **USEQ-11330086** — Restrict asset use to authorized purposes and environments.
- [ ] **USEQ-9D66D98F** — Prevent production assets from being repurposed for testing without controlled sanitization and approval.
- [ ] **USEQ-61F13688** — Ensure shared assets have clear administration, accountability, and access records.
- [ ] **USEQ-95BBBD9D** — Reconcile assets during personnel, supplier, office, and organizational transitions.
- [ ] **USEQ-1102BC52** — Investigate missing, stale, unresponsive, or ownerless assets.
- [ ] **USEQ-F82AEC89** — Revoke access and recover assets promptly during offboarding.

### Software, service, and cloud asset lifecycle

- [ ] **USEQ-BCC8A948** — Inventory deployed applications, services, libraries, images, functions, jobs, APIs, endpoints, agents, and management planes.
- [ ] **USEQ-F0FA28E9** — Record exact versions, editions, regions, accounts, projects, environments, and deployment identifiers.
- [ ] **USEQ-2C9BCEDC** — Identify unsupported, end-of-life, unmaintained, and unpatchable components.
- [ ] **USEQ-404CE25B** — Track public exposure and externally reachable management interfaces.
- [ ] **USEQ-DBD981E5** — Identify software installed outside approved deployment or package mechanisms.
- [ ] **USEQ-4521F0FD** — Detect abandoned cloud resources, storage, addresses, domains, certificates, secrets, snapshots, and service accounts.
- [ ] **USEQ-53BFB44D** — Link every running workload to reviewed source, build provenance, artifact, configuration, owner, and service.
- [ ] **USEQ-7E25ABA0** — Maintain approved and prohibited technology lists with documented exception processes.
- [ ] **USEQ-59EF6776** — Define patch, upgrade, replacement, and retirement plans before support expires.
- [ ] **USEQ-DE68E9CE** — Verify that decommissioning removes traffic, access, data, credentials, billing, monitoring, DNS, and dependencies.
- [ ] **USEQ-B7D56A01** — Retain required historical artifacts and documentation for support and investigation.

### Configuration-management planning

- [ ] **USEQ-70E4676E** — Define which products, components, documents, schemas, models, infrastructure, policies, and tools are configuration items.
- [ ] **USEQ-9E197CB5** — Assign unique identifiers and version conventions to configuration items.
- [ ] **USEQ-16AFEC0D** — Define authoritative repositories and custodians for each configuration-item class.
- [ ] **USEQ-A9FA1137** — Establish baselines at meaningful lifecycle and release points.
- [ ] **USEQ-ADC79C0E** — Define who may propose, review, approve, implement, verify, and audit changes.
- [ ] **USEQ-A2BA66FC** — Scale control rigor according to criticality and consequence.
- [ ] **USEQ-C21F7C87** — Include supplier-provided and customer-controlled configuration where it affects outcomes.
- [ ] **USEQ-205CA0F6** — Define emergency and temporary-change procedures with retrospective review and expiry.
- [ ] **USEQ-AA3E4353** — Document configuration relationships and compatibility constraints.
- [ ] **USEQ-BB179760** — Preserve the configuration-management plan and evidence with the product lifecycle record.

### Configuration identification and baselines

- [ ] **USEQ-33E78AD8** — Identify every item required to reproduce, operate, support, verify, recover, and retire a release.
- [ ] **USEQ-989F2CFE** — Include source, dependencies, build definitions, tool versions, infrastructure, schema, data migrations, runtime configuration, flags, secrets references, certificates, models, prompts, documentation, and tests as applicable.
- [ ] **USEQ-4BC4CCAE** — Record immutable digests or equivalent identifiers for deployable and evidentiary artifacts.
- [ ] **USEQ-6A253B4C** — Define which parts of a baseline are immutable and which may change operationally.
- [ ] **USEQ-DCF98B7D** — Include environment-specific configuration and external-service versions in the operational baseline.
- [ ] **USEQ-899A62A9** — Prevent ambiguous labels such as “latest,” “current,” or “production” from being the sole identifier.
- [ ] **USEQ-5823C219** — Maintain links among requirements, changes, configuration items, builds, tests, approvals, and deployments.
- [ ] **USEQ-D50778E8** — Preserve prior known-good baselines needed for rollback and investigation.
- [ ] **USEQ-0889583B** — Detect missing, duplicate, conflicting, and orphaned configuration records.
- [ ] **USEQ-709323D1** — Verify the documented baseline against the actual deployed state.

### Change control and configuration status accounting

- [ ] **USEQ-FC3E698D** — Require every material configuration change to state purpose, scope, risk, dependencies, testing, rollout, rollback, owner, and approval.
- [ ] **USEQ-43388938** — Assess impacts across components, data, users, suppliers, obligations, documentation, and operations.
- [ ] **USEQ-79CA7209** — Prevent unauthorized direct change to controlled production configuration.
- [ ] **USEQ-504287F4** — Record who changed what, when, why, from which previous state, and with what result.
- [ ] **USEQ-A2D17FE8** — Track proposed, approved, implemented, verified, deployed, rolled-back, superseded, and retired states.
- [ ] **USEQ-E33E127B** — Ensure rejected and canceled changes cannot be deployed accidentally.
- [ ] **USEQ-11CF89FF** — Detect and reconcile configuration drift.
- [ ] **USEQ-A049228C** — Make emergency changes visible and subject to prompt retrospective review.
- [ ] **USEQ-755987E5** — Ensure temporary settings, elevated access, diagnostics, and feature flags expire or are removed.
- [ ] **USEQ-71E0C724** — Report baseline status, open deviations, pending changes, and unresolved discrepancies to decision-makers.
- [ ] **USEQ-D0653183** — Prevent change records from being closed before implementation and verification evidence exists.

### Configuration verification and audit

- [ ] **USEQ-3117AB90** — Verify that each baseline is complete, internally consistent, approved, retrievable, and reproducible.
- [ ] **USEQ-FEC62FEC** — Compare deployed, running, documented, and approved configuration.
- [ ] **USEQ-F752B649** — Audit whether change-control procedures were followed and effective.
- [ ] **USEQ-2F4DE3A8** — Verify that generated artifacts can be traced to declared inputs and tools.
- [ ] **USEQ-8CD13776** — Detect unapproved components, hidden defaults, local patches, manual edits, and undocumented overrides.
- [ ] **USEQ-0C0D4B9F** — Verify that rollback artifacts, migrations, keys, and documentation remain usable.
- [ ] **USEQ-9886EAF3** — Sample high-risk settings for secure and correct values.
- [ ] **USEQ-3082692A** — Reconcile configuration after incidents, failover, disaster recovery, and provider restoration.
- [ ] **USEQ-94277832** — Record discrepancies, owners, corrective actions, and effectiveness checks.
- [ ] **USEQ-868E9E2E** — Treat inability to identify the production baseline as a release blocker.

### Engineering-tool inventory and qualification

- [ ] **USEQ-C2FBDFF0** — Inventory tools that can create, modify, test, approve, build, sign, deploy, monitor, diagnose, migrate, or destroy production assets.
- [ ] **USEQ-35C25338** — Classify tools by the consequence of incorrect output, compromise, unavailability, or misuse.
- [ ] **USEQ-5A38E631** — Record tool owner, supplier, version, configuration, privileges, data access, support status, and integrations.
- [ ] **USEQ-92D287A4** — Define expected behavior and acceptance criteria for tools whose output supports critical decisions.
- [ ] **USEQ-190FA346** — Validate or qualify tools proportionate to their role and risk.
- [ ] **USEQ-53E9BA34** — Independently verify tool output when an undetected tool error could create unacceptable harm.
- [ ] **USEQ-5BF8AE95** — Use known-answer tests, reference cases, cross-tool comparison, or manual review where appropriate.
- [ ] **USEQ-31FA56CF** — Record tool limitations, false-positive and false-negative behavior, unsupported inputs, and operating constraints.
- [ ] **USEQ-6B714CFD** — Requalify tools after material version, configuration, plugin, model, rule-set, or environment changes.
- [ ] **USEQ-3BCEA6C5** — Prevent unreviewed extensions and plugins from gaining toolchain privileges.
- [ ] **USEQ-64FCAB05** — Verify AI-assisted tools separately for the tasks and context in which their output is relied upon.

### Toolchain access and integrity

- [ ] **USEQ-5B569CEA** — Use individual identities and strong authentication for engineering and management tools.
- [ ] **USEQ-BCD3314F** — Apply least privilege to repositories, pipelines, registries, signing, deployment, cloud, data, and monitoring systems.
- [ ] **USEQ-46326F4E** — Separate untrusted contribution contexts from secrets and production authority.
- [ ] **USEQ-DB33526C** — Protect tool configuration, policies, workflows, rules, and plugins with review and version control.
- [ ] **USEQ-1D612C9C** — Record administrative and high-impact tool actions.
- [ ] **USEQ-2CD634D1** — Monitor unusual access, token use, configuration change, artifact publication, and approval behavior.
- [ ] **USEQ-32C3102B** — Rotate and revoke tool credentials, tokens, certificates, deploy keys, and integration secrets.
- [ ] **USEQ-D268F72B** — Prevent one compromised tool from silently changing source, tests, artifacts, and evidence without independent detection.
- [ ] **USEQ-5AC58356** — Use artifact signing, provenance, immutable logs, separation of duties, or equivalent controls to reduce toolchain trust concentration.
- [ ] **USEQ-2E9EE695** — Back up and test recovery of repositories, configurations, pipelines, registries, evidence, and documentation.
- [ ] **USEQ-D8872F98** — Maintain alternate access and recovery methods for identity-provider or control-plane failure.
- [ ] **USEQ-9EF2F036** — Include toolchain compromise in incident-response and continuity exercises.

### Toolchain availability, performance, and developer safety

- [ ] **USEQ-714190FE** — Define service objectives for critical engineering tools and workflows.
- [ ] **USEQ-6879B287** — Monitor availability, latency, queueing, capacity, error rate, and data integrity.
- [ ] **USEQ-00E9100D** — Avoid making one tool or administrator the sole path to build, release, restore, or investigate.
- [ ] **USEQ-7213309B** — Provide safe degraded workflows for urgent security or recovery work.
- [ ] **USEQ-0D23F8EE** — Prevent tool outages from encouraging uncontrolled local workarounds or secret sharing.
- [ ] **USEQ-44ACC81E** — Protect developers from misleading partial success, stale results, and silently skipped checks.
- [ ] **USEQ-CC01C62E** — Make required gates, failures, retries, and evidence status clear.
- [ ] **USEQ-D478B533** — Preserve work and allow safe resumption after interruption.
- [ ] **USEQ-E66C02D0** — Test toolchain recovery and backlog processing after outage.
- [ ] **USEQ-AACF05BE** — Capacity-plan repositories, runners, registries, logs, artifact stores, and test environments.

### Retirement, disposal, and audit

- [ ] **USEQ-7B9E1A6F** — Define retirement criteria for assets, configurations, tools, services, licenses, and platforms.
- [ ] **USEQ-392818A7** — Confirm dependencies and consumers before retirement.
- [ ] **USEQ-C4ACDFF7** — Preserve required data, records, evidence, source, artifacts, keys, and documentation.
- [ ] **USEQ-B87013D8** — Remove access, integrations, routes, accounts, credentials, certificates, data, backups, and billing according to the retirement plan.
- [ ] **USEQ-72EC7BEC** — Sanitize or destroy physical and logical storage appropriately.
- [ ] **USEQ-A40FB25C** — Update inventories, ownership, diagrams, support records, and financial systems.
- [ ] **USEQ-D5B4C02F** — Verify no production traffic or authoritative updates still reach the retired asset.
- [ ] **USEQ-2581CD16** — Monitor for attempted use after retirement.
- [ ] **USEQ-06B4C16D** — Record residual obligations and retained components.
- [ ] **USEQ-AB436474** — Audit asset and configuration processes for completeness, accuracy, timeliness, and control effectiveness.
- [ ] **USEQ-662CE68B** — Treat unknown critical assets, unsupported production technology, untraceable configuration, or unqualified high-impact tools as no-go conditions.

## Standards and source references

- [ISO/IEC/IEEE 12207:2026 — Software life cycle processes](https://www.iso.org/standard/90219.html)
- [ISO/IEC 20246:2017 — Work product reviews](https://www.iso.org/standard/67597.html)
- [NIST SP 800-218 v1.1 — Secure Software Development Framework](https://csrc.nist.gov/pubs/sp/800/218/final)
- [SLSA Specification 1.2](https://slsa.dev/spec/v1.2/)
- [OpenSSF Security Baseline and Best Practices](https://baseline.openssf.org/)
- [ISO/IEC 27001:2022 — Information security management systems](https://www.iso.org/standard/27001)
- [DORA software delivery performance research](https://dora.dev/guides/dora-metrics/)
- [ISO/IEC 20000-1:2018 — Service management system requirements](https://www.iso.org/standard/70636.html)
- [NIST Cybersecurity Framework 2.0](https://www.nist.gov/cyberframework)
- [Google SRE Workbook — Implementing SLOs](https://sre.google/workbook/implementing-slos/)
- [ISO/IEC 25012:2008 — Data quality model](https://www.iso.org/standard/35736.html)
- [ISO/IEC/IEEE 29119-2:2021 — Test processes](https://www.iso.org/standard/79428.html)
- [ISO 22301:2019 — Business continuity management systems](https://www.iso.org/standard/75106.html)
- [DORA platform engineering capability](https://dora.dev/capabilities/platform-engineering/)
- [DORA research](https://dora.dev/research/)
- [CNCF Platforms White Paper](https://tag-app-delivery.cncf.io/whitepapers/platforms/)
- [CNCF Platform Engineering Maturity Model](https://tag-app-delivery.cncf.io/whitepapers/platform-eng-maturity-model/)
- [ACM Queue — DevEx: What Actually Drives Productivity](https://queue.acm.org/detail.cfm?id=3595878)
- [IEEE Computer Society SWEBOK v4](https://www.computer.org/education/bodies-of-knowledge/software-engineering)
- [GAO Cost Estimating and Assessment Guide](https://www.gao.gov/products/gao-20-195g)
- [GAO Schedule Assessment Guide](https://www.gao.gov/products/gao-16-89g)
- [NIST SP 800-161 Rev. 1 — Cybersecurity Supply Chain Risk Management](https://csrc.nist.gov/pubs/sp/800/161/r1/final)

---

[Previous phase](10-verification-and-testing.md) · [Next: Phase 12: Operations, SRE, and support](12-operations-sre-and-support.md)
