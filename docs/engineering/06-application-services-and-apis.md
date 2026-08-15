# Application services and APIs

_Phase 6 of 16 in the [complete engineering review](00-overview.md)._

Backend services, APIs, identity flows, tenant isolation, content processing, jobs, queues, caching, and integrations.

Assess every control as **Pass**, **Fail**, **Blocked**, or **Not Applicable**. A Pass needs current evidence for the exact product, revision, artifact, configuration, environment, and data state under review.

## APIs, Web Services, Webhooks, and Integrations

_Consolidated from `quality standards/07-backend-services/01-apis-webhooks-and-integrations.md`; 31 non-duplicative controls._

### Universal controls

- [ ] **USEQ-AD050959** — Maintain a complete API and integration inventory.
- [ ] **USEQ-A87B5F2B** — Every API has an owner, documented consumers, versioned schemas, and versioned behavior.
- [ ] **USEQ-4ECE79C1** — Unknown or prohibited fields are rejected or handled according to a documented safe contract.
- [ ] **USEQ-E87054DC** — Content type, character encoding, request size, URL size, header size, nesting, field count, query complexity, and batch size are bounded.
- [ ] **USEQ-F652D754** — Expensive queries, filtering, sorting, pagination, and aggregation have cost and authorization controls.
- [ ] **USEQ-4F1DF89C** — Request methods have correct and consistent semantics.
- [ ] **USEQ-C44A6ECD** — Idempotency and concurrency control exist where retries can repeat effects or updates can collide.
- [ ] **USEQ-F6ACEF15** — Every outbound call has an appropriate timeout.
- [ ] **USEQ-C57DC947** — Retries are bounded, delayed with backoff and jitter where appropriate, and limited to safe or idempotent operations.
- [ ] **USEQ-F8A92BB0** — Circuit breaking, isolation, or fallback exists for unstable dependencies where appropriate.
- [ ] **USEQ-554B725E** — Error responses do not expose stack traces, queries, secrets, internal topology, or sensitive data.
- [ ] **USEQ-80BFBA34** — Rate limits and quotas protect users, tenants, infrastructure, and dependencies and cannot be trivially bypassed.
- [ ] **USEQ-BBED00FD** — Cross-origin access uses explicit trusted origins; credentialed requests do not use wildcard origins.
- [ ] **USEQ-89E3635B** — Browser-authenticated state-changing requests have CSRF protection.
- [ ] **USEQ-DEAA8863** — Redirect and callback destinations are strictly validated.
- [ ] **USEQ-B6B803BC** — Server-side outbound destinations are constrained against server-side request forgery, DNS rebinding, dangerous schemes, redirects, and metadata services.
- [ ] **USEQ-E0FB8D9E** — URL, host, proxy, and forwarded-header parsing is consistent and explicitly trusted.
- [ ] **USEQ-CE328AB7** — Request smuggling, response splitting, cache poisoning, and parser ambiguity are assessed.
- [ ] **USEQ-31202793** — Webhooks are authenticated or signed and include replay protection.
- [ ] **USEQ-8C1E3CB4** — Webhook consumers are idempotent; retry, ordering, failure, and reconciliation behavior is defined.
- [ ] **USEQ-4667AD11** — API keys are scoped, attributable, rotatable, and revocable.
- [ ] **USEQ-F500AC1C** — Secrets are absent from responses, documentation examples, and sample requests.
- [ ] **USEQ-5184C8B1** — Backward compatibility, breaking-change, deprecation, and retirement policies are defined and tested.
- [ ] **USEQ-F1736D64** — Contract tests cover providers and consumers.
- [ ] **USEQ-96565CD6** — API documentation matches deployed behavior and contains no real credentials or personal data.

### Conditional interface checks

- [ ] **USEQ-EB908BF8** — Graph-oriented APIs bound query depth, breadth, aliases, recursion, batching, and computational cost.
- [ ] **USEQ-84E3AA0F** — WebSocket and streaming connections authenticate securely and verify browser origin where applicable.
- [ ] **USEQ-9115A638** — Authorization changes and logout invalidate or constrain long-lived connections.
- [ ] **USEQ-13717FE0** — Reconnection, ordering, deduplication, backpressure, heartbeat, and stale-connection behavior are tested.
- [ ] **USEQ-A2756DDC** — File-transfer APIs apply all file-security controls.
- [ ] **USEQ-79910F4D** — Public APIs have abuse controls, usage policies, and a security contact.

## Authentication, Enrollment, and Account Recovery

_Consolidated from `quality standards/07-backend-services/02-authentication-enrollment-and-recovery.md`; 25 non-duplicative controls._

### Universal controls

- [ ] **USEQ-D6B72D8C** — Select identity and authentication assurance appropriate to account-compromise harm.
- [ ] **USEQ-45326065** — Multifactor enrollment, addition, removal, replacement, and recovery are protected against takeover.
- [ ] **USEQ-C6330D44** — Password handling follows current guidance and does not silently truncate or unexpectedly transform input.
- [ ] **USEQ-A7F41DCE** — Login defenses address credential stuffing and brute force without creating trivial denial of service.
- [ ] **USEQ-8DEBE078** — Login, registration, and recovery responses avoid unnecessary account-existence disclosure.
- [ ] **USEQ-339F2E8F** — Enrollment verifies the correct communication channel or identity evidence.
- [ ] **USEQ-5D7958B0** — Recovery is no weaker than normal authentication and cannot be redirected through an unverified channel.
- [ ] **USEQ-C8EED75F** — Recovery tokens are single-use, short-lived, unpredictable, and stored safely.
- [ ] **USEQ-5013A7F0** — Password, authenticator, or recovery changes invalidate appropriate sessions and obsolete recovery material.
- [ ] **USEQ-75ECA584** — Changes to email, phone, payout, identity, recovery, or security settings notify the user and receive enhanced verification where warranted.
- [ ] **USEQ-7AD37548** — Test, dormant, abandoned, and terminated accounts are handled according to policy.
- [ ] **USEQ-C3A74331** — Service accounts do not use inappropriate interactive recovery mechanisms.
- [ ] **USEQ-E27BE14A** — Authentication events are audited and suspicious authentication can trigger investigation or step-up controls.
- [ ] **USEQ-EBEBDB32** — Users can review and revoke relevant sessions and devices where appropriate.
- [ ] **USEQ-C6C9BCF9** — Recovery codes are protected, single-use, and regenerable securely.
- [ ] **USEQ-5A69500A** — Support procedures include social-engineering resistance.

### Federation and delegated authorization

- [ ] **USEQ-52AF78CA** — Issuer, audience, signature, algorithm, expiration, nonce, and state are validated as required.
- [ ] **USEQ-36E65D58** — Redirect identifiers are matched exactly according to protocol security requirements.
- [ ] **USEQ-9170A512** — Authorization-code flows use current proof-key protections where applicable.
- [ ] **USEQ-50604945** — Token substitution, replay, and mix-up attacks are addressed.
- [ ] **USEQ-281FCFA1** — Automatic account linking cannot merge an attacker's identity with a victim account.
- [ ] **USEQ-F493DD34** — Tokens are limited to minimum required permissions and audiences.
- [ ] **USEQ-642071D4** — Access and refresh tokens are stored safely.
- [ ] **USEQ-7BC31490** — Refresh-token rotation and reuse detection are used where appropriate.
- [ ] **USEQ-15527916** — Revocation, logout, key rotation, metadata refresh, and provider outage behavior are tested.

## Authorization, Privilege Management, and Tenant Isolation

_Consolidated from `quality standards/07-backend-services/03-authorization-and-tenant-isolation.md`; 12 non-duplicative controls._

### Universal controls

- [ ] **USEQ-871674ED** — Authorization is denied by default and enforced server-side at every trusted boundary.
- [ ] **USEQ-E4D0202F** — Access checks use the current user, tenant, resource, operation, and resource state.
- [ ] **USEQ-10AFC75B** — Object ownership is verified on every direct and indirect reference.
- [ ] **USEQ-24A910E7** — Horizontal, vertical, function-level, and property-level privilege escalation are tested.
- [ ] **USEQ-825DC852** — Bulk, export, search, report, and aggregation endpoints enforce the same controls as individual records.
- [ ] **USEQ-873CC668** — Role changes take effect promptly across sessions, caches, queues, services, indexes, and long-lived connections.
- [ ] **USEQ-002BF231** — High-impact actions require reauthentication or additional approval where appropriate.
- [ ] **USEQ-66605220** — Administrative, support, and impersonation actions require authorization, reason, visible indication where appropriate, and audit.
- [ ] **USEQ-5D7FE64E** — Impersonation does not expose secrets the support user should not receive.
- [ ] **USEQ-54A140C8** — Background jobs preserve correct authorization and tenant context.
- [ ] **USEQ-54DF7A38** — Caches, search indexes, files, exports, backups, logs, analytics, and queues preserve user and tenant boundaries.
- [ ] **USEQ-92DFA818** — Automated tests cover all material role-operation-resource-state combinations.

## Sessions, Cookies, and Tokens

_Consolidated from `quality standards/07-backend-services/04-sessions-cookies-and-tokens.md`; 13 non-duplicative controls._

### Universal controls

- [ ] **USEQ-FC13E22C** — Session identifiers are generated from a cryptographically secure source and have sufficient unpredictability.
- [ ] **USEQ-B79A2729** — Session identifiers are accepted only through intended secure channels.
- [ ] **USEQ-F38F9324** — Browser cookies use secure transport, script-access, cross-site, domain, and path restrictions appropriate to their purpose.
- [ ] **USEQ-931DDFC8** — Session identifiers rotate after authentication and privilege changes.
- [ ] **USEQ-A784E8C1** — Logout revokes server-side access.
- [ ] **USEQ-B836719F** — Recovery and credential changes invalidate appropriate sessions.
- [ ] **USEQ-5FD43C2D** — Persistent login uses separate, appropriately constrained credentials.
- [ ] **USEQ-92AF7C8C** — State-changing browser operations have CSRF protection.
- [ ] **USEQ-7255F599** — Tokens validate issuer, audience, subject where needed, type, signature, algorithm, expiration, and not-before time.
- [ ] **USEQ-D912D5FB** — Token replay is addressed; refresh credentials are protected against reuse where appropriate.
- [ ] **USEQ-B061EEB0** — Tokens and session identifiers do not appear in URLs, referrers, analytics, logs, or error reports.
- [ ] **USEQ-956C315B** — Client-side token storage is minimized and threat-modeled.
- [ ] **USEQ-4589A6C9** — Logout, revocation, and compromise behavior is tested across all components.

## File and Content Processing

_Consolidated from `quality standards/07-backend-services/05-file-and-content-processing.md`; 13 non-duplicative controls._

### Universal controls

- [ ] **USEQ-685621E2** — File size, count, archive depth, expanded size, dimensions, page count, and processing time are bounded.
- [ ] **USEQ-2C48ADF4** — File type is determined from content and structure rather than trusting user metadata alone.
- [ ] **USEQ-43F66C3F** — Allowed file types are explicitly listed.
- [ ] **USEQ-8F094721** — User filenames are not trusted as filesystem paths.
- [ ] **USEQ-7326A422** — Stored names are unpredictable where guessing creates risk.
- [ ] **USEQ-10DBEE38** — Untrusted content uses an isolated origin where browser execution could be dangerous.
- [ ] **USEQ-ACD5C4B4** — Access control is checked during upload, listing, preview, transformation, and download.
- [ ] **USEQ-0EEFF2CC** — Download responses use safe content type and disposition behavior and prevent content sniffing.
- [ ] **USEQ-273EB093** — Image, document, media, archive, and metadata parsers are patched and isolated appropriately.
- [ ] **USEQ-F614699C** — Temporary, abandoned, failed, and orphaned files are protected and cleaned up.
- [ ] **USEQ-F69D7216** — Retention and deletion rules apply to originals and derived versions.
- [ ] **USEQ-9C66FEF7** — Processing failures do not expose internal paths or parser details.
- [ ] **USEQ-7C1A282F** — Content-moderation and legal controls apply where users can share files.

## Multi-Tenant Systems and Customer Isolation

_Consolidated from `quality standards/07-backend-services/06-multi-tenant-systems.md`; 8 non-duplicative controls._

### Universal controls

- [ ] **USEQ-1A813A2E** — Tenant context is mandatory at every data, processing, caching, indexing, queueing, file, logging, and administrative boundary.
- [ ] **USEQ-866A2A93** — Tenant isolation is tested in databases, caches, search, queues, files, logs, analytics, backups, exports, and support tools.
- [ ] **USEQ-72B786B4** — Tenant configuration, feature flags, custom code, integrations, and domains cannot affect other tenants.
- [ ] **USEQ-C89DD0B2** — Cross-tenant timing, error, count, identifier, and existence side channels are minimized.
- [ ] **USEQ-AEE95887** — Support access to a tenant requires authorization, purpose, visibility where appropriate, and audit.
- [ ] **USEQ-0111A6D5** — Tenant creation, cloning, migration, split, merge, and deletion preserve isolation and auditability.
- [ ] **USEQ-3E18FD8A** — Shared templates, catalogs, and global objects cannot leak private tenant data.
- [ ] **USEQ-80B4087C** — Per-tenant data residency, retention, encryption, and compliance commitments are enforced.

## High-Risk Administrative and Support Tools

_Consolidated from `quality standards/07-backend-services/07-high-risk-admin-and-support-tools.md`; 15 non-duplicative controls._

### Universal controls

- [ ] **USEQ-08FE86F1** — Administrative and support interfaces are separately inventoried and threat-modeled.
- [ ] **USEQ-0ADF382B** — Strong, preferably phishing-resistant authentication is mandatory for privileged access where risk warrants it.
- [ ] **USEQ-70CC4AE3** — Privileges are granular and follow least privilege.
- [ ] **USEQ-2829A59B** — Access is just-in-time or time-bounded where practical.
- [ ] **USEQ-1EB641B7** — High-impact actions require recent authentication or additional approval where appropriate.
- [ ] **USEQ-BB90D77D** — Bulk exports and changes require confirmation, scope preview, authorization, and audit.
- [ ] **USEQ-D933D2D3** — Destructive commands have safe limits, dry-run or preview behavior where practical, and clear recovery instructions.
- [ ] **USEQ-312E1130** — Impersonation is authorized, justified, visibly indicated, time-bounded, and logged.
- [ ] **USEQ-B5AA5454** — Sensitive values are masked unless explicitly needed.
- [ ] **USEQ-17D96426** — Emergency and break-glass access triggers review.
- [ ] **USEQ-4BDEADF5** — Administrative APIs receive the same or stronger testing as user-facing APIs.
- [ ] **USEQ-AC203899** — Support tooling cannot bypass tenant, data-residency, retention, or legal boundaries.
- [ ] **USEQ-7A2D2764** — Sessions terminate promptly after privilege removal.
- [ ] **USEQ-8AFCE597** — Every action is attributable to an individual identity.
- [ ] **USEQ-C2F34517** — Privileged activity is monitored for unusual access, bulk operations, exports, and control changes.

## Enterprise Identity, Provisioning, and Organization Lifecycle

_Consolidated from `quality standards/07-backend-services/08-enterprise-identity-provisioning-and-organization-lifecycle.md`; 12 non-duplicative controls._

### Universal controls

- [ ] **USEQ-447370B0** — Enterprise federation, domain verification, directory synchronization, provisioning, and deprovisioning are inventoried.
- [ ] **USEQ-52027B79** — Organization and domain claiming cannot be hijacked.
- [ ] **USEQ-FFAB6C6F** — Federation setup requires verified administrative control.
- [ ] **USEQ-EFFC5073** — Provisioning creates users with minimum access.
- [ ] **USEQ-ECFAF330** — Deprovisioning promptly disables sessions, tokens, API keys, devices, and background access.
- [ ] **USEQ-FA59F542** — Group and role mapping is deterministic, reviewed, and resistant to privilege escalation.
- [ ] **USEQ-0FC2ABAB** — Directory synchronization handles rename, merge, duplicate, suspension, reactivation, and deletion safely.
- [ ] **USEQ-347028DE** — Just-in-time provisioning cannot bypass organization policy.
- [ ] **USEQ-9DF78FDA** — Multiple identity providers and migration between providers do not create orphaned or duplicate identities.
- [ ] **USEQ-3BB3C453** — Break-glass access exists for identity-provider outage and is monitored.
- [ ] **USEQ-9BD8A2CC** — Enterprise audit logs, exports, retention, data residency, and administrative controls meet contractual commitments.
- [ ] **USEQ-9069FC2B** — Organization transfer, merger, split, and termination preserve data and authorization boundaries.

## Backend Service Design

_Consolidated from `quality standards/07-backend-services/09-backend-service-design.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-1CC90997** — Define each service's cohesive capability, owner, consumers, data authority, and operational objectives.
- [ ] **USEQ-AE48DCE5** — Split services only when independent ownership, scaling, security, change, or failure isolation justifies the boundary.
- [ ] **USEQ-3307A379** — Avoid shared databases or mutable internals that bypass service contracts without explicit governance.
- [ ] **USEQ-CB9942F4** — Keep business rules in trusted domain logic rather than transport handlers or client code.
- [ ] **USEQ-ED684D68** — Use explicit request, command, query, event, and error contracts.
- [ ] **USEQ-10891705** — Validate identity, authorization, tenant, input, quotas, and invariants at every entry point.
- [ ] **USEQ-1CE1D773** — Make write operations idempotent where clients or infrastructure can retry.
- [ ] **USEQ-FA593D00** — Keep transaction boundaries aligned to invariants and avoid unbounded remote transactions.
- [ ] **USEQ-41F351E9** — Define timeouts, retries, cancellation, backpressure, and overload behavior.
- [ ] **USEQ-078CA238** — Make dependency failure degrade safely and observably.
- [ ] **USEQ-271D3DFF** — Use structured logging, metrics, traces, health, and audit events tied to user and business outcomes.
- [ ] **USEQ-4265879C** — Separate readiness from liveness and stop accepting work before shutdown.
- [ ] **USEQ-9BD0F34D** — Drain, transfer, checkpoint, or safely retry in-flight work during deployment and failure.
- [ ] **USEQ-2D916DB0** — Keep configuration, secrets, schema, and deployment versions traceable.
- [ ] **USEQ-DAE49FB5** — Support backward-compatible rolling deployment where required.
- [ ] **USEQ-8DAB072E** — Test service contracts, data behavior, failure, recovery, load, security, and migration independently.
- [ ] **USEQ-C947839C** — Prevent one tenant, operation, or dependency from exhausting all resources.
- [ ] **USEQ-A88B60B7** — Document support, runbook, SLO, capacity, backup, recovery, and retirement obligations.
- [ ] **USEQ-D2AB670F** — Avoid exposing internal implementation or database structure as a permanent external contract.
- [ ] **USEQ-D0769D00** — Measure whether the service boundary reduces total coordination and failure cost over time.

## Background Jobs, Scheduling, and Queues

_Consolidated from `quality standards/07-backend-services/10-background-jobs-scheduling-and-queues.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-16F68399** — Assign every job type, queue, schedule, and dead-letter stream an owner.
- [ ] **USEQ-5918958F** — Define trigger, input schema, tenant and identity context, priority, deadline, retry, ordering, and completion semantics.
- [ ] **USEQ-1CFCCEFB** — Assume at-least-once execution unless stronger guarantees are proven.
- [ ] **USEQ-DA5D8874** — Make side-effecting jobs idempotent or protect effects with deduplication keys.
- [ ] **USEQ-881FD155** — Persist sufficient state to resume after worker, process, node, or region failure.
- [ ] **USEQ-4B1B01AF** — Bound attempts, age, execution time, concurrency, queue depth, payload, and resource use.
- [ ] **USEQ-56F9433C** — Use backoff and jitter and prevent retry synchronization.
- [ ] **USEQ-CDCA58F2** — Distinguish retryable dependency failure from permanent validation, authorization, or business rejection.
- [ ] **USEQ-AAB22D62** — Prevent poison work from blocking unrelated work.
- [ ] **USEQ-056E5454** — Provide dead-letter inspection, privacy protection, alerting, remediation, and safe replay.
- [ ] **USEQ-13ECE301** — Maintain ordering only where the business requires it and account for partitioning and rebalancing.
- [ ] **USEQ-CC43CE3C** — Protect against overlapping scheduled runs and missed runs.
- [ ] **USEQ-E679A7AF** — Define clock, time-zone, daylight-saving, calendar, and catch-up behavior.
- [ ] **USEQ-66D9DF06** — Propagate trace, tenant, identity, data classification, and causal context.
- [ ] **USEQ-7BF28954** — Make cancellation and supersession observable and safe.
- [ ] **USEQ-42ACA9BE** — Avoid embedding mutable large objects or secrets in queue payloads.
- [ ] **USEQ-BAE577C3** — Revalidate authorization and current state when delayed work executes.
- [ ] **USEQ-84A401CA** — Monitor queue age, throughput, failures, retries, lag, saturation, and business completion.
- [ ] **USEQ-D14A5097** — Test duplicate, delayed, reordered, lost, corrupt, expired, and partially processed work.
- [ ] **USEQ-6A7C0140** — Reconcile authoritative outcomes against job records to detect silent loss or duplication.

## Caching, Search, and Derived Data

_Consolidated from `quality standards/07-backend-services/11-caching-search-and-derived-data.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-FA365A87** — Classify each cache or index as disposable, reconstructable, authoritative, or temporarily authoritative.
- [ ] **USEQ-27109AA1** — Define source of truth, key dimensions, freshness, expiry, invalidation, consistency, and reconstruction behavior.
- [ ] **USEQ-06FF47B4** — Include identity, tenant, authorization, locale, version, content negotiation, and privacy dimensions in keys where required.
- [ ] **USEQ-24FF04B1** — Prevent cache poisoning, key collision, unbounded cardinality, stampede, and eviction of critical entries.
- [ ] **USEQ-E104A243** — Never serve private or personalized content from a shared cache without correct isolation.
- [ ] **USEQ-CA1CF7A0** — Use negative caching only with safe expiry and authorization semantics.
- [ ] **USEQ-3E156CB0** — Make stale serving deliberate and communicate staleness when it affects decisions.
- [ ] **USEQ-64AE632C** — Bound memory, disk, entry size, result size, query complexity, and retention.
- [ ] **USEQ-2C6E1A16** — Protect search filters, facets, snippets, ranking, and aggregation with the same authorization as source data.
- [ ] **USEQ-719D7363** — Propagate deletion, correction, legal hold, consent, and tenant changes to derived stores.
- [ ] **USEQ-569D6922** — Monitor lag, hit rate, miss cost, eviction, index failures, partial updates, and reconstruction duration.
- [ ] **USEQ-31CA65E1** — Provide cache bypass and index rebuild procedures for incidents.
- [ ] **USEQ-700AB4BA** — Avoid requiring a cache to preserve correctness unless its durability and recovery are engineered accordingly.
- [ ] **USEQ-0386371D** — Prevent search indexing from exposing secrets, hidden fields, deleted data, or cross-tenant terms.
- [ ] **USEQ-86A0C9B1** — Test cold start, empty cache, stale entries, concurrent invalidation, failover, and partial indexing.
- [ ] **USEQ-305C6668** — Reconcile derived counts, documents, and permissions with authoritative data.
- [ ] **USEQ-DD834B76** — Version schemas and ranking behavior for compatible rollout.
- [ ] **USEQ-909674E0** — Do not use cache as an undocumented communication or locking mechanism.
- [ ] **USEQ-17F09290** — Make query and ranking behavior explainable enough to diagnose user-visible errors.
- [ ] **USEQ-5C56E75A** — Retire obsolete indexes and cached representations with verified cleanup.

## Standards and source references

- [RFC 9110 — HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110)
- [RFC 9457 — Problem Details for HTTP APIs](https://www.rfc-editor.org/rfc/rfc9457)
- [OWASP API Security Top 10 — 2023](https://owasp.org/API-Security/)
- [OWASP Application Security Verification Standard 5.0.0](https://owasp.org/www-project-application-security-verification-standard/)
- [NIST SP 800-63-4 — Digital Identity Guidelines](https://pages.nist.gov/800-63-4/)
- [RFC 9700 — OAuth 2.0 Security Best Current Practice](https://www.rfc-editor.org/rfc/rfc9700)
- [NIST Cybersecurity Framework 2.0](https://www.nist.gov/cyberframework)
- [RFC 8725 — JSON Web Token Best Current Practices](https://www.rfc-editor.org/rfc/rfc8725)
- [OWASP Web Security Testing Guide 4.2](https://owasp.org/www-project-web-security-testing-guide/)
- [ISO/IEC 27001:2022 — Information security management systems](https://www.iso.org/standard/27001)
- [ISO/IEC 25010:2023 — Product quality model](https://www.iso.org/standard/78176.html)
- [ISO/IEC/IEEE 42010:2022 — Architecture description](https://www.iso.org/standard/74393.html)
- [ISO/IEC/IEEE 12207:2026 — Software life cycle processes](https://www.iso.org/standard/90219.html)
- [ISO/IEC/IEEE 29119-4:2021 — Test techniques](https://www.iso.org/standard/79430.html)
- [Google SRE Workbook — Implementing SLOs](https://sre.google/workbook/implementing-slos/)
- [ISO/IEC 25012:2008 — Data quality model](https://www.iso.org/standard/35736.html)

---

[Previous phase](05-code-quality-and-implementation.md) · [Next: Phase 7: Data and information lifecycle](07-data-and-information-lifecycle.md)
