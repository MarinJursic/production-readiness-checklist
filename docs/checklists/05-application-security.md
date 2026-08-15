# Application security

> Review interfaces, identity, authorization, sessions, untrusted input, files, transport, and cryptography.

Sections 15–22 of the master checklist. For each applicable item, capture status, owner, evidence, release, environment, reviewer, and evidence age.

## 15. APIs, web services, and integrations

- [ ] **PRC-15-001** — Maintain a complete API inventory.
- [ ] **PRC-15-002** — Every API has an owner and documented consumers.
- [ ] **PRC-15-003** — Interface schemas and behavior are versioned.
- [ ] **PRC-15-004** — Authentication is enforced on every protected operation.
- [ ] **PRC-15-005** — Authorization is evaluated for every object, property, operation, and tenant.
- [ ] **PRC-15-006** — Input is validated against an explicit schema.
- [ ] **PRC-15-007** — Unknown or prohibited fields are rejected or handled safely.
- [ ] **PRC-15-008** — Content types and character encodings are constrained.
- [ ] **PRC-15-009** — Request body, URL, header, nesting, field, query, and batch sizes are bounded.
- [ ] **PRC-15-010** — Expensive queries and filters have cost controls.
- [ ] **PRC-15-011** — Pagination has safe maximums.
- [ ] **PRC-15-012** — Sorting and filtering cannot bypass access control.
- [ ] **PRC-15-013** — Mass assignment and unexpected property binding are prevented.
- [ ] **PRC-15-014** — HTTP methods have correct and consistent semantics.
- [ ] **PRC-15-015** — Idempotency is implemented where retries can repeat effects.
- [ ] **PRC-15-016** — Concurrency control prevents lost updates where required.
- [ ] **PRC-15-017** — Timeouts exist for all outbound calls.
- [ ] **PRC-15-018** — Retries are bounded and use backoff and jitter.
- [ ] **PRC-15-019** — Retries occur only for safe or idempotent operations.
- [ ] **PRC-15-020** — Circuit breaking or equivalent isolation exists for unstable dependencies.
- [ ] **PRC-15-021** — Errors use consistent safe structures and appropriate status codes.
- [ ] **PRC-15-022** — Error responses do not expose stack traces, queries, secrets, or internal topology.
- [ ] **PRC-15-023** — Rate limits and quotas protect users, tenants, infrastructure, and dependencies.
- [ ] **PRC-15-024** — Rate limits cannot be trivially bypassed through alternate identifiers or routes.
- [ ] **PRC-15-025** — CORS uses explicit trusted origins where cross-origin access is needed.
- [ ] **PRC-15-026** — Credentialed cross-origin requests do not use wildcard origins.
- [ ] **PRC-15-027** — Browser-authenticated state-changing operations have CSRF protection.
- [ ] **PRC-15-028** — Sensitive credentials and tokens never appear in URLs.
- [ ] **PRC-15-029** — Redirects and callback destinations are strictly validated.
- [ ] **PRC-15-030** — Server-side outbound destinations are constrained to prevent SSRF.
- [ ] **PRC-15-031** — URL parsing is consistent across proxies, application code, and security controls.
- [ ] **PRC-15-032** — Host-header and forwarded-header trust is explicitly configured.
- [ ] **PRC-15-033** — Request smuggling, response splitting, cache poisoning, and parser ambiguity have been assessed.
- [ ] **PRC-15-034** — Webhooks are authenticated or signed.
- [ ] **PRC-15-035** — Webhooks include replay protection.
- [ ] **PRC-15-036** — Webhook consumers are idempotent.
- [ ] **PRC-15-037** — Webhook retries, ordering, failure, and reconciliation are defined.
- [ ] **PRC-15-038** — API keys are scoped, rotated, revocable, and attributable.
- [ ] **PRC-15-039** — Secrets are not returned in responses or documentation examples.
- [ ] **PRC-15-040** — Backward compatibility is tested.
- [ ] **PRC-15-041** — Breaking-change and deprecation policies exist.
- [ ] **PRC-15-042** — Old versions have retirement dates and monitoring.
- [ ] **PRC-15-043** — Contract tests cover external providers and consumers.
- [ ] **PRC-15-044** — API documentation matches deployed behavior.
- [ ] **PRC-15-045** — Production examples contain no real credentials or personal data.

### 15.1 Conditional interface checks

- [ ] **PRC-15-046** — Graph-based APIs limit depth, breadth, recursion, aliases, batching, and computational cost.
- [ ] **PRC-15-047** — Graph schema introspection exposure is intentionally decided.
- [ ] **PRC-15-048** — WebSocket or streaming connections authenticate securely.
- [ ] **PRC-15-049** — Streaming connections verify origin where browser initiated.
- [ ] **PRC-15-050** — Authorization changes and logout invalidate long-lived connections.
- [ ] **PRC-15-051** — Reconnect, ordering, deduplication, backpressure, and heartbeat behavior are tested.
- [ ] **PRC-15-052** — File-transfer APIs apply the file-security controls in this checklist.
- [ ] **PRC-15-053** — Partner and machine-to-machine APIs use appropriately scoped identities.
- [ ] **PRC-15-054** — Public APIs have abuse controls, usage policies, and security contact information.

---

## 16. Identity, authentication, and account recovery

NIST SP 800-63 Revision 4, finalized in July 2025, provides current risk-based guidance for identity proofing, authentication, authenticators, and federation. OAuth deployments should also follow the current OAuth 2.0 Security Best Current Practice. ([pages.nist.gov](https://pages.nist.gov/800-63-4/))

- [ ] **PRC-16-001** — Inventory every login, enrollment, recovery, federation, API, support, and administrative authentication path.
- [ ] **PRC-16-002** — Select identity and authentication assurance appropriate to the harm of account compromise.
- [ ] **PRC-16-003** — Require stronger authentication for privileged and high-impact accounts.
- [ ] **PRC-16-004** — Support phishing-resistant authentication where risk warrants it.
- [ ] **PRC-16-005** — Multifactor enrollment is protected from account takeover.
- [ ] **PRC-16-006** — New authenticators require appropriate reauthentication.
- [ ] **PRC-16-007** — Removing or replacing an authenticator requires appropriate verification.
- [ ] **PRC-16-008** — Password handling follows current guidance.
- [ ] **PRC-16-009** — Passwords are not silently truncated or transformed unexpectedly.
- [ ] **PRC-16-010** — Common and compromised passwords are rejected where password authentication is used.
- [ ] **PRC-16-011** — Password managers and paste are supported.
- [ ] **PRC-16-012** — Login defenses address credential stuffing and brute force.
- [ ] **PRC-16-013** — Defenses avoid creating trivial denial-of-service against legitimate users.
- [ ] **PRC-16-014** — Login, registration, and recovery responses do not unnecessarily disclose whether an account exists.
- [ ] **PRC-16-015** — Account enrollment verifies the correct communication channel or identity evidence.
- [ ] **PRC-16-016** — Account recovery is no weaker than normal authentication.
- [ ] **PRC-16-017** — Recovery cannot be redirected through an unverified channel.
- [ ] **PRC-16-018** — Recovery tokens are single-use, short-lived, random, and stored safely.
- [ ] **PRC-16-019** — Password or authenticator reset invalidates appropriate sessions and recovery material.
- [ ] **PRC-16-020** — High-risk operations require recent authentication.
- [ ] **PRC-16-021** — Email, phone, payout, identity, and security-setting changes notify the user appropriately.
- [ ] **PRC-16-022** — Changes to recovery information use delayed or additional verification where warranted.
- [ ] **PRC-16-023** — Default credentials are absent.
- [ ] **PRC-16-024** — Test, dormant, abandoned, and terminated accounts are disabled.
- [ ] **PRC-16-025** — Service accounts cannot use interactive recovery mechanisms inappropriately.
- [ ] **PRC-16-026** — Authentication events are audited.
- [ ] **PRC-16-027** — Suspicious authentication can trigger step-up controls or investigation.
- [ ] **PRC-16-028** — Users can review and revoke relevant sessions and devices.
- [ ] **PRC-16-029** — Clock skew and token timing behavior are tested.
- [ ] **PRC-16-030** — Recovery codes are protected, single-use, and regenerated securely.
- [ ] **PRC-16-031** — Support personnel cannot bypass identity checks without a controlled, audited exception process.
- [ ] **PRC-16-032** — Social-engineering resistance is included in support procedures.

### 16.1 Federation and delegated authorization

- [ ] **PRC-16-033** — Issuer, audience, signature, algorithm, expiration, nonce, and state are validated.
- [ ] **PRC-16-034** — Redirect URIs are matched exactly according to the protocol’s security requirements.
- [ ] **PRC-16-035** — Authorization-code flows use current protections such as proof-key mechanisms where applicable.
- [ ] **PRC-16-036** — Deprecated or insecure delegated-authorization flows are not used.
- [ ] **PRC-16-037** — Token substitution and mix-up attacks are addressed.
- [ ] **PRC-16-038** — Identity-provider accounts are linked only after secure verification.
- [ ] **PRC-16-039** — Automatic account linking cannot merge an attacker’s identity with a victim’s account.
- [ ] **PRC-16-040** — Tokens are scoped to the minimum required audience and permissions.
- [ ] **PRC-16-041** — Access and refresh tokens are stored and transmitted safely.
- [ ] **PRC-16-042** — Refresh-token rotation and reuse detection are applied where appropriate.
- [ ] **PRC-16-043** — Revocation, logout, key rotation, and identity-provider outage behavior are tested.
- [ ] **PRC-16-044** — Federation metadata and signing keys are refreshed safely.
- [ ] **PRC-16-045** — Identity claims are not trusted beyond their documented assurance.
- [ ] **PRC-16-046** — Tenant or organization selection cannot be manipulated to cross boundaries.

---

## 17. Authorization, privileges, and tenant isolation

- [ ] **PRC-17-001** — Authorization is denied by default.
- [ ] **PRC-17-002** — Authorization is enforced server-side at every trusted boundary.
- [ ] **PRC-17-003** — A current role-permission-resource matrix exists.
- [ ] **PRC-17-004** — Every endpoint and action maps to a required permission.
- [ ] **PRC-17-005** — Access checks use the current user, tenant, resource, and operation.
- [ ] **PRC-17-006** — Object ownership is verified on every direct and indirect object reference.
- [ ] **PRC-17-007** — Horizontal privilege escalation has been tested.
- [ ] **PRC-17-008** — Vertical privilege escalation has been tested.
- [ ] **PRC-17-009** — Function-level and property-level authorization have been tested.
- [ ] **PRC-17-010** — Bulk, export, search, reporting, and aggregation endpoints enforce the same controls as individual records.
- [ ] **PRC-17-011** — Hidden fields and client-side controls cannot grant access.
- [ ] **PRC-17-012** — Role changes take effect promptly across sessions, caches, queues, and services.
- [ ] **PRC-17-013** — Removed users lose access to shared resources.
- [ ] **PRC-17-014** — Nested groups and inherited permissions are tested.
- [ ] **PRC-17-015** — Temporary permissions expire automatically.
- [ ] **PRC-17-016** — Privileged access follows least privilege and just-in-time principles where practical.
- [ ] **PRC-17-017** — High-impact actions require additional approval where appropriate.
- [ ] **PRC-17-018** — Administrative, support, and impersonation actions require reason, authorization, visible indication, and audit.
- [ ] **PRC-17-019** — Impersonation never exposes secrets the support user should not receive.
- [ ] **PRC-17-020** — Emergency access is time-bounded and reviewed.
- [ ] **PRC-17-021** — Authorization decisions are logged at an appropriate level.
- [ ] **PRC-17-022** — Background jobs retain correct authorization and tenant context.
- [ ] **PRC-17-023** — Caches do not return data authorized for another principal.
- [ ] **PRC-17-024** — Search indexes preserve authorization boundaries.
- [ ] **PRC-17-025** — Files, exports, backups, logs, and analytics preserve tenant boundaries.
- [ ] **PRC-17-026** — Queue messages cannot be reassigned across tenants through manipulated metadata.
- [ ] **PRC-17-027** — Tenant-specific encryption, keys, domains, and configuration are isolated as designed.
- [ ] **PRC-17-028** — Resource identifiers alone are never treated as authorization.
- [ ] **PRC-17-029** — Side channels do not disclose another tenant’s existence, activity, identifiers, or resource usage.
- [ ] **PRC-17-030** — Automated tests cover every role × operation × resource-state combination considered material.

---

## 18. Sessions, cookies, and tokens

- [ ] **PRC-18-001** — Session identifiers have sufficient unpredictability.
- [ ] **PRC-18-002** — Session identifiers are generated using a cryptographically secure source.
- [ ] **PRC-18-003** — Session identifiers are never accepted from insecure or unintended channels.
- [ ] **PRC-18-004** — Browser cookies use `Secure` where transmitted over HTTPS.
- [ ] **PRC-18-005** — Sensitive browser cookies use `HttpOnly` unless script access is essential and justified.
- [ ] **PRC-18-006** — `SameSite` behavior is intentionally selected and tested.
- [ ] **PRC-18-007** — Cookie domain and path scope are as narrow as practical.
- [ ] **PRC-18-008** — Cross-subdomain cookie trust is explicitly assessed.
- [ ] **PRC-18-009** — Session identifiers rotate after login and privilege changes.
- [ ] **PRC-18-010** — Session fixation is prevented.
- [ ] **PRC-18-011** — Idle and absolute timeouts match risk and user needs.
- [ ] **PRC-18-012** — Logout invalidates the server-side session or otherwise revokes access.
- [ ] **PRC-18-013** — Account recovery and credential change invalidate appropriate sessions.
- [ ] **PRC-18-014** — Concurrent-session policy is defined.
- [ ] **PRC-18-015** — Users can identify and revoke sessions where appropriate.
- [ ] **PRC-18-016** — “Remember me” behavior has separate, appropriately constrained credentials.
- [ ] **PRC-18-017** — Sensitive state-changing actions have CSRF protection.
- [ ] **PRC-18-018** — Tokens validate issuer, audience, subject, type, signature, algorithm, expiry, and not-before values.
- [ ] **PRC-18-019** — A token intended for one purpose cannot be accepted for another.
- [ ] **PRC-18-020** — Token replay risk is addressed.
- [ ] **PRC-18-021** — Refresh credentials are rotated or otherwise protected against replay where appropriate.
- [ ] **PRC-18-022** — Tokens and session identifiers are absent from URLs, referrers, analytics, logs, and error reports.
- [ ] **PRC-18-023** — Client-side storage of tokens is minimized and threat-modeled.
- [ ] **PRC-18-024** — Logout, revocation, and compromise behavior is tested across all application components.
- [ ] **PRC-18-025** — Session-related logs do not expose reusable credentials.

---

## 19. Input validation, encoding, injection, and unsafe processing

- [ ] **PRC-19-001** — All untrusted input is identified.
- [ ] **PRC-19-002** — Validation occurs at the trusted server-side boundary.
- [ ] **PRC-19-003** — Inputs are validated for type, syntax, range, length, structure, and allowed values.
- [ ] **PRC-19-004** — Canonicalization and normalization occur before security-sensitive comparisons.
- [ ] **PRC-19-005** — Validation uses allowlists where a finite valid set exists.
- [ ] **PRC-19-006** — Database operations use safe parameter binding.
- [ ] **PRC-19-007** — Operating-system operations avoid command construction from untrusted input.
- [ ] **PRC-19-008** — Directory, query, search, template, expression, mail, logging, and interpreter contexts use appropriate safe APIs.
- [ ] **PRC-19-009** — Output is encoded for its exact destination context.
- [ ] **PRC-19-010** — HTML or rich-text content is sanitized with an appropriate maintained sanitizer.
- [ ] **PRC-19-011** — Untrusted input cannot escape into executable script, style, markup, query, command, or template contexts.
- [ ] **PRC-19-012** — Dynamic evaluation of untrusted input is absent.
- [ ] **PRC-19-013** — Unsafe deserialization is avoided.
- [ ] **PRC-19-014** — Deserialization permits only expected types and bounded structures.
- [ ] **PRC-19-015** — XML processing disables unnecessary external entity and external resource behavior.
- [ ] **PRC-19-016** — Path traversal and archive traversal are prevented.
- [ ] **PRC-19-017** — Symbolic-link and filesystem race behavior is controlled.
- [ ] **PRC-19-018** — Server-side request forgery protections constrain protocol, destination, redirects, DNS resolution, and metadata endpoints.
- [ ] **PRC-19-019** — Regular expressions cannot be used for practical resource-exhaustion attacks.
- [ ] **PRC-19-020** — Header, response-splitting, log, CSV, formula, and email injection are addressed.
- [ ] **PRC-19-021** — Open redirects are prevented.
- [ ] **PRC-19-022** — Object and property injection through automatic binding or merge behavior is addressed.
- [ ] **PRC-19-023** — Unicode confusables and normalization do not bypass validation or identity checks.
- [ ] **PRC-19-024** — Error handling does not reflect unsafe input into executable contexts.
- [ ] **PRC-19-025** — Security tests cover injection across every supported parser and protocol.

---

## 20. File upload, download, and content processing

- [ ] **PRC-20-001** — File upload is enabled only where necessary.
- [ ] **PRC-20-002** — Maximum file size is enforced before excessive processing.
- [ ] **PRC-20-003** — File count, archive depth, expanded size, dimensions, page count, and processing time are bounded.
- [ ] **PRC-20-004** — File type is checked using content, structure, and expected extension rather than user-supplied metadata alone.
- [ ] **PRC-20-005** — Allowed types are explicitly listed.
- [ ] **PRC-20-006** — Active or executable content is rejected or neutralized where not required.
- [ ] **PRC-20-007** — Uploaded content is malware-scanned where risk warrants it.
- [ ] **PRC-20-008** — Files remain quarantined until required validation completes.
- [ ] **PRC-20-009** — User-provided filenames are not used as trusted filesystem paths.
- [ ] **PRC-20-010** — Stored filenames are non-predictable where guessing would create risk.
- [ ] **PRC-20-011** — Upload storage is non-executable.
- [ ] **PRC-20-012** — Untrusted files use an isolated origin where browser execution could be dangerous.
- [ ] **PRC-20-013** — Access control is checked when files are uploaded, listed, previewed, transformed, and downloaded.
- [ ] **PRC-20-014** — Download responses use safe content type and disposition behavior.
- [ ] **PRC-20-015** — Browser content sniffing is prevented.
- [ ] **PRC-20-016** — Archive extraction prevents traversal, symbolic-link abuse, and resource exhaustion.
- [ ] **PRC-20-017** — Image, document, media, and metadata parsers are patched and isolated appropriately.
- [ ] **PRC-20-018** — Re-encoding or content disarm is used where warranted.
- [ ] **PRC-20-019** — Generated exports escape spreadsheet formulas and other active content.
- [ ] **PRC-20-020** — Temporary files have restricted permissions and reliable deletion.
- [ ] **PRC-20-021** — Abandoned, failed, and orphaned uploads are cleaned up.
- [ ] **PRC-20-022** — Retention and deletion rules apply to stored files and derived versions.
- [ ] **PRC-20-023** — File-processing failure cannot expose internal paths or parser details.
- [ ] **PRC-20-024** — Large or malicious uploads cannot starve shared resources.
- [ ] **PRC-20-025** — Download links expire or are access-controlled appropriately.
- [ ] **PRC-20-026** — Content moderation and legal controls are applied where users can share files with others.

---

## 21. Transport, browser security, DNS, and network exposure

- [ ] **PRC-21-001** — All user and administrative traffic uses an approved encrypted transport configuration.
- [ ] **PRC-21-002** — Plaintext access is redirected or rejected safely.
- [ ] **PRC-21-003** — Certificates cover the correct names and trust chains.
- [ ] **PRC-21-004** — Certificate issuance, renewal, expiration, and revocation are monitored.
- [ ] **PRC-21-005** — Internal service traffic is encrypted where the threat model requires it.
- [ ] **PRC-21-006** — Mutual authentication is used for high-risk service connections where warranted.
- [ ] **PRC-21-007** — Mixed active content is absent.
- [ ] **PRC-21-008** — Strict transport policy is enabled after confirming subdomain and recovery readiness.
- [ ] **PRC-21-009** — Sensitive responses have appropriate cache controls.
- [ ] **PRC-21-010** — Shared caches cannot store or serve personalized data incorrectly.
- [ ] **PRC-21-011** — Content Security Policy is designed, deployed, and monitored where applicable.
- [ ] **PRC-21-012** — Framing is restricted to intended parents.
- [ ] **PRC-21-013** — MIME sniffing is disabled.
- [ ] **PRC-21-014** — Referrer leakage is controlled.
- [ ] **PRC-21-015** — Browser capability permissions are restricted.
- [ ] **PRC-21-016** — Cross-origin opener, resource, and embedding policies are assessed where isolation matters.
- [ ] **PRC-21-017** — CORS configuration is explicit and tested.
- [ ] **PRC-21-018** — Third-party scripts and resources are inventoried and justified.
- [ ] **PRC-21-019** — External static resources use integrity protection or an equivalently controlled delivery method where practical.
- [ ] **PRC-21-020** — Security headers are present on normal, error, redirect, authentication, and static responses as applicable.
- [ ] **PRC-21-021** — Ingress rules expose only required services and ports.
- [ ] **PRC-21-022** — Administrative services are not publicly exposed unnecessarily.
- [ ] **PRC-21-023** — Egress is restricted where practical, especially for high-value workloads.
- [ ] **PRC-21-024** — Proxy and load-balancer header trust is configured explicitly.
- [ ] **PRC-21-025** — The origin cannot be bypassed unintentionally when a CDN, WAF, or reverse proxy is relied upon.
- [ ] **PRC-21-026** — DDoS and volumetric-abuse protections match expected risk.
- [ ] **PRC-21-027** — Network segmentation limits lateral movement.
- [ ] **PRC-21-028** — Production networks do not trust development or office networks implicitly.
- [ ] **PRC-21-029** — Domain ownership is documented.
- [ ] **PRC-21-030** — Registrar access uses strong authentication and restricted roles.
- [ ] **PRC-21-031** — Domain transfer and modification protections are enabled.
- [ ] **PRC-21-032** — DNS records are inventoried.
- [ ] **PRC-21-033** — Dangling DNS records and abandoned hosted resources are removed.
- [ ] **PRC-21-034** — Subdomain takeover has been tested.
- [ ] **PRC-21-035** — DNS changes, TTLs, failover, and recovery are documented.
- [ ] **PRC-21-036** — DNS-security extensions are evaluated according to threat and operational capability.
- [ ] **PRC-21-037** — A standards-compliant `security.txt` or equivalent vulnerability-reporting path is published where appropriate.
- [ ] **PRC-21-038** — External attack-surface monitoring detects newly exposed hosts and services.

---

## 22. Cryptography and key management

- [ ] **PRC-22-001** — Cryptographic requirements derive from a documented threat and data classification.
- [ ] **PRC-22-002** — Only approved, maintained algorithms and protocol implementations are used.
- [ ] **PRC-22-003** — Custom cryptographic algorithms and ad hoc protocol designs are absent.
- [ ] **PRC-22-004** — Encryption provides integrity and authentication as required, not merely confidentiality.
- [ ] **PRC-22-005** — Passwords use a current adaptive one-way password-hashing method.
- [ ] **PRC-22-006** — Encryption keys are stored separately from protected data where practical.
- [ ] **PRC-22-007** — Keys are absent from source code, logs, artifacts, frontend code, and ordinary configuration files.
- [ ] **PRC-22-008** — Key access follows least privilege.
- [ ] **PRC-22-009** — Production and non-production use separate keys.
- [ ] **PRC-22-010** — Different purposes use separate keys where appropriate.
- [ ] **PRC-22-011** — Key generation uses approved secure randomness.
- [ ] **PRC-22-012** — Nonces and initialization values meet uniqueness and unpredictability requirements.
- [ ] **PRC-22-013** — Key rotation and versioning are supported.
- [ ] **PRC-22-014** — Rotation does not make existing required data unrecoverable.
- [ ] **PRC-22-015** — Compromised keys can be revoked promptly.
- [ ] **PRC-22-016** — Key compromise has a tested response playbook.
- [ ] **PRC-22-017** — Required keys are included securely in continuity and recovery planning.
- [ ] **PRC-22-018** — Key deletion and cryptographic erasure behavior are defined.
- [ ] **PRC-22-019** — Certificates, signing keys, token keys, encryption keys, and webhook keys each have lifecycle owners.
- [ ] **PRC-22-020** — Signature verification rejects inappropriate algorithms, keys, or malformed input.
- [ ] **PRC-22-021** — Random identifiers and security tokens use cryptographically secure generation.
- [ ] **PRC-22-022** — Cryptographic dependencies are inventoried for future algorithm migration.
- [ ] **PRC-22-023** — Sensitive plaintext is not copied into analytics, logs, caches, search indexes, or backups unnecessarily.

---
