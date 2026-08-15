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
