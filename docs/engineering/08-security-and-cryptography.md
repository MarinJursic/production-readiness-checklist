# Security and cryptography

_Phase 8 of 16 in the [complete engineering review](00-overview.md)._

Secure development, identity, application security, supply chain, vulnerability response, monitoring, cryptographic agility, and post-quantum planning.

Assess every control as **Pass**, **Fail**, **Blocked**, or **Not Applicable**. A Pass needs current evidence for the exact product, revision, artifact, configuration, environment, and data state under review.

## Dependencies, SBOM, and Licenses

_Consolidated from `quality standards/09-security/01-dependencies-sbom-and-licenses.md`; 15 non-duplicative controls._

### Universal controls

- [ ] **USEQ-CD388C32** — Generate an SBOM for every released artifact.
- [ ] **USEQ-F1E36337** — Include direct, transitive, runtime, operating-system, container, client, server, plugin, model, and embedded components as applicable.
- [ ] **USEQ-26843351** — Record exact versions, package identifiers, sources, hashes, relationships, licenses, and required notices.
- [ ] **USEQ-91621C41** — Link the SBOM to the released artifact digest and verify it against the final package.
- [ ] **USEQ-93155005** — Ensure no undeclared or unexpected components are present.
- [ ] **USEQ-321F26E3** — Remove unused, unsupported, and end-of-life dependencies.
- [ ] **USEQ-93206DB8** — Confirm dependencies come from intended publishers and registries.
- [ ] **USEQ-9CB28275** — Monitor upstream advisories and continuously scan released components.
- [ ] **USEQ-74B47B32** — Prioritize vulnerabilities using exploitation, reachability, exposure, criticality, data sensitivity, and compensating controls.
- [ ] **USEQ-F36AAE1B** — Define remediation targets by actual risk and verify fixes or mitigations after deployment.
- [ ] **USEQ-D2C8D527** — Maintain an emergency replacement process for compromised dependencies.
- [ ] **USEQ-3D9B945D** — Document ownership of vendored, forked, and internally patched components.
- [ ] **USEQ-1A4E996C** — Verify license compatibility and fulfill attribution, notice, source-offer, and redistribution obligations.
- [ ] **USEQ-5E5663FA** — Prevent unintended inclusion of proprietary or confidential code in distributable artifacts.
- [ ] **USEQ-E5774CFB** — Retain SBOMs for historical releases and investigations.

## Transport, Browser, DNS, and Network Security

_Consolidated from `quality standards/09-security/02-transport-browser-dns-and-network-security.md`; 15 non-duplicative controls._

### Universal controls

- [ ] **USEQ-7437D577** — Certificates cover correct names and trust chains.
- [ ] **USEQ-6A0D1D7E** — Certificate issuance, renewal, expiration, and revocation are monitored and tested.
- [ ] **USEQ-BD9DA6D5** — Strict transport policy is enabled only after confirming subdomain and recovery readiness.
- [ ] **USEQ-532FDBCB** — A Content Security Policy or equivalent browser execution policy is designed, deployed, tested, and monitored where applicable.
- [ ] **USEQ-8CFBF610** — Cross-origin opener, resource, and embedding isolation policies are assessed where relevant.
- [ ] **USEQ-F25290EA** — Third-party scripts and external resources are inventoried, justified, and controlled.
- [ ] **USEQ-C3FFF6E2** — External static resources use integrity protection or an equivalently controlled delivery mechanism where practical.
- [ ] **USEQ-C9411EEA** — Ingress exposes only required services and ports.
- [ ] **USEQ-433F17DA** — Proxy, load-balancer, host, and forwarded-header trust is configured explicitly.
- [ ] **USEQ-CCB23489** — The origin cannot be bypassed when a CDN, reverse proxy, gateway, or application firewall is relied upon.
- [ ] **USEQ-FA0BCE37** — Production networks do not implicitly trust development, testing, office, or personal networks.
- [ ] **USEQ-A812272C** — Registrar access uses strong authentication, restricted roles, transfer protection, and modification protection.
- [ ] **USEQ-34A4E09F** — DNS changes, TTLs, failover, and recovery are documented and tested where required.
- [ ] **USEQ-6671CDB9** — A standards-compliant security contact file or equivalent reporting path is published where appropriate.
- [ ] **USEQ-69F5FCA9** — External attack-surface monitoring detects newly exposed hosts, domains, ports, and services.

## Cryptography and Key Management

_Consolidated from `quality standards/09-security/03-cryptography-and-key-management.md`; 13 non-duplicative controls._

### Universal controls

- [ ] **USEQ-CB2472B5** — Cryptographic requirements derive from a documented threat model and data classification.
- [ ] **USEQ-EEE7D855** — Only approved, maintained algorithms, libraries, and protocol implementations are used.
- [ ] **USEQ-99B102B8** — Custom cryptographic algorithms and ad hoc protocols are absent.
- [ ] **USEQ-2D34262D** — Encryption provides confidentiality, integrity, and authenticity as required by the threat model.
- [ ] **USEQ-F3E22AF9** — Keys are stored separately from protected data where practical.
- [ ] **USEQ-C45F8665** — Keys do not appear in source, logs, artifacts, frontend code, tickets, or ordinary configuration.
- [ ] **USEQ-CF133A7A** — Rotation does not make required existing data unrecoverable.
- [ ] **USEQ-49F94E8A** — Required keys and recovery credentials are securely included in continuity and recovery planning.
- [ ] **USEQ-1C8FFFBF** — Key deletion and cryptographic-erasure behavior are defined.
- [ ] **USEQ-66A51D5D** — Certificates, signing keys, token keys, encryption keys, and webhook keys have lifecycle owners.
- [ ] **USEQ-7FAB83B4** — Signature verification rejects inappropriate algorithms, keys, malformed input, and algorithm confusion.
- [ ] **USEQ-A42AA818** — Random identifiers and security tokens use cryptographically secure generation where unpredictability is required.
- [ ] **USEQ-8BD5C6C7** — Sensitive plaintext is not copied unnecessarily into analytics, logs, caches, indexes, backups, or support tools.

## Security Validation and Vulnerability Management

_Consolidated from `quality standards/09-security/04-security-validation-and-vulnerability-management.md`; 14 non-duplicative controls._

### Universal controls

- [ ] **USEQ-A55D5542** — The release is assessed against the applicable OWASP Application Security Verification Standard requirements.
- [ ] **USEQ-A9032C59** — Dependency, SBOM, and license analysis has run.
- [ ] **USEQ-6EB61A35** — Secret scanning has run against current changes and relevant history.
- [ ] **USEQ-212AE882** — Infrastructure, configuration, and policy analysis has run.
- [ ] **USEQ-E111813C** — Deployed application and API security testing has run.
- [ ] **USEQ-A54FFC98** — Authentication, account recovery, session, authorization, and tenant-isolation testing has run.
- [ ] **USEQ-F470D5E2** — Injection, unsafe-input, browser, file-processing, business-logic, abuse, and fraud testing has run where applicable.
- [ ] **USEQ-A5C687CC** — Penetration testing covers authenticated users, multiple roles, multiple tenants, APIs, administrative tools, and business logic as applicable.
- [ ] **USEQ-65F8195D** — Findings have owners, risk ratings, remediation targets, and retest evidence.
- [ ] **USEQ-215CA1ED** — Relevant known-exploited vulnerabilities are blockers unless non-exposure and mitigation are convincingly demonstrated.
- [ ] **USEQ-D6455092** — Unsupported and unpatchable components have explicit treatment plans.
- [ ] **USEQ-376AA3EF** — A vulnerability-disclosure policy and security-reporting channel exist.
- [ ] **USEQ-B0C31F8B** — Security reports receive acknowledgment, triage, remediation tracking, coordinated disclosure handling, and appropriate credit.
- [ ] **USEQ-8C2B79FA** — Security monitoring covers account takeover, privilege changes, secret use, suspicious exports, abuse, fraud, and control tampering.

## Security Incident Response and Crisis Readiness

_Consolidated from `quality standards/09-security/05-incident-response-and-crisis-readiness.md`; 14 non-duplicative controls._

### Universal controls

- [ ] **USEQ-B5A25CFD** — Incident severity levels and declaration authority are defined.
- [ ] **USEQ-6A5A824B** — Roles cover incident command, technical response, communications, documentation, security, privacy/legal, and business ownership as needed.
- [ ] **USEQ-E91F44FF** — Contact details are current and remain available during identity, network, provider, email, chat, or documentation failure.
- [ ] **USEQ-2EDD7F91** — Procedures cover detection, triage, containment, eradication, recovery, verification, and closure.
- [ ] **USEQ-EFB31A9B** — Procedures distinguish outages, security incidents, privacy incidents, fraud events, safety events, and data-integrity events.
- [ ] **USEQ-1364ED3A** — Evidence-preservation and forensic procedures exist where needed.
- [ ] **USEQ-6E828BED** — Credentials, sessions, certificates, signing material, and encryption keys can be revoked or rotated rapidly.
- [ ] **USEQ-E4837A19** — Compromised builds and dependencies can be identified, contained, rebuilt, and replaced.
- [ ] **USEQ-C961296D** — Customer, regulator, partner, insurer, law-enforcement, and public communication templates and decision paths exist where applicable.
- [ ] **USEQ-DF6272E2** — Alternate operations are possible during critical provider failure where required.
- [ ] **USEQ-17C5B126** — Tabletop exercises cover realistic scenarios such as account takeover, data exfiltration, destructive access, ransomware, dependency compromise, regional outage, data corruption, and provider failure.
- [ ] **USEQ-9D9785D5** — Lessons are converted into tests, controls, runbooks, monitoring, training, and architecture changes.
- [ ] **USEQ-E968A368** — Incident metrics include detection, acknowledgment, containment, recovery, impact, customer harm, and recurrence.
- [ ] **USEQ-191DAB9C** — The response plan can operate when normal identity, cloud, chat, email, ticketing, or documentation systems are unavailable.

## Public Forms, Anonymous Access, and Abuse Resistance

_Consolidated from `quality standards/09-security/06-public-forms-anonymous-access-and-abuse-resistance.md`; 10 non-duplicative controls._

### Universal controls

- [ ] **USEQ-F77F6419** — Anonymous and unauthenticated endpoints are inventoried.
- [ ] **USEQ-D4E90196** — Submission size, rate, cost, frequency, and resource use are bounded.
- [ ] **USEQ-DA8F4EC3** — Spam, scraping, credential stuffing, enumeration, denial of service, and automated abuse are addressed.
- [ ] **USEQ-068808BE** — Anti-automation controls are accessible and do not create disproportionate user harm.
- [ ] **USEQ-B8EC811C** — Challenges do not become the sole security control for authorization or fraud prevention.
- [ ] **USEQ-4CD39055** — Rate limits account for shared networks, privacy relays, attackers rotating identifiers, and legitimate high-volume users.
- [ ] **USEQ-C16A799E** — Contact, signup, invite, reset, search, and feedback forms cannot be used to harass or message-bomb third parties.
- [ ] **USEQ-C8A6D225** — File and URL submissions apply all content-processing and server-side request protections.
- [ ] **USEQ-97A2A4F4** — Abuse signals, false positives, appeals, and manual-review paths are monitored.
- [ ] **USEQ-50A7DF0C** — Emergency controls can restrict abusive traffic without disabling essential legitimate access unnecessarily.

## Security Governance and Risk Management

_Consolidated from `quality standards/09-security/07-security-governance-and-risk-management.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-6C65E57D** — Define security objectives from business context, assets, threats, legal duties, contracts, and risk appetite.
- [ ] **USEQ-6DE0A64D** — Assign accountable owners for information, services, identities, systems, controls, vulnerabilities, and incidents.
- [ ] **USEQ-4101D65F** — Maintain a current inventory of assets, data, interfaces, dependencies, identities, secrets, and externally exposed surfaces.
- [ ] **USEQ-543BE90B** — Use a consistent risk method that records threat, vulnerability, likelihood, impact, existing controls, residual risk, and owner.
- [ ] **USEQ-73C4644A** — Prioritize based on realistic exploitability, exposure, reachability, user harm, business impact, and control strength.
- [ ] **USEQ-3AC61448** — Select controls proportionate to risk and verify they operate in the actual environment.
- [ ] **USEQ-D1B6C513** — Separate policy requirements from implementation guidance and document justified tailoring.
- [ ] **USEQ-1B6BE375** — Define minimum baselines for products, environments, suppliers, identities, endpoints, and data classes.
- [ ] **USEQ-24FE911C** — Require independent review for material security decisions and exceptions.
- [ ] **USEQ-514A892D** — Make security exceptions time-bounded with compensating controls, monitoring, and automatic expiry.
- [ ] **USEQ-16C353A9** — Integrate security objectives into funding, planning, requirements, architecture, procurement, delivery, and operations.
- [ ] **USEQ-FBE36F84** — Protect security staff independence while keeping product teams accountable for owned risk.
- [ ] **USEQ-61067AF7** — Measure control effectiveness, detection, remediation, recurrence, and risk reduction rather than finding counts alone.
- [ ] **USEQ-53911D5E** — Review inherited provider controls and identify customer responsibilities explicitly.
- [ ] **USEQ-5EE950A3** — Maintain vulnerability disclosure, researcher communication, and coordinated remediation processes.
- [ ] **USEQ-4F0EE479** — Use incidents, threat intelligence, exercises, audits, and near misses to update controls.
- [ ] **USEQ-1A267F87** — Prevent security work from being deferred indefinitely as nonfunctional or post-launch.
- [ ] **USEQ-2C25E805** — Review risk after material scope, data, identity, supplier, exposure, architecture, or threat changes.
- [ ] **USEQ-AD2024D1** — Maintain executive visibility of risk that exceeds approved tolerance.
- [ ] **USEQ-728DA235** — Do not claim security based solely on compliance, certification, penetration tests, or absence of known findings.

## Secure Software Development Lifecycle

_Consolidated from `quality standards/09-security/08-secure-software-development-lifecycle.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-4892B305** — Define role-specific secure-development responsibilities and competency requirements.
- [ ] **USEQ-8B4975B0** — Protect source, build, test, artifact, deployment, signing, and production environments according to their trust.
- [ ] **USEQ-A4C4F453** — Translate threat models and security objectives into traceable requirements and verification.
- [ ] **USEQ-966F7118** — Use approved patterns for identity, authorization, input handling, secrets, cryptography, logging, and error responses.
- [ ] **USEQ-D0514E08** — Review security-sensitive architecture and code with qualified independent reviewers.
- [ ] **USEQ-9BE73656** — Scan source, dependencies, secrets, infrastructure, configuration, and artifacts using maintained rules.
- [ ] **USEQ-63F37D23** — Treat automated findings as evidence inputs requiring triage, not automatic truth or automatic dismissal.
- [ ] **USEQ-7AF842BD** — Perform manual testing for business logic, authorization, identity recovery, tenant isolation, and abuse cases.
- [ ] **USEQ-8C8E061F** — Verify the exact release artifact and configuration rather than only source branches.
- [ ] **USEQ-7C42AB74** — Generate and retain provenance, dependency inventory, SBOM, signatures, test evidence, and approvals.
- [ ] **USEQ-12D7F021** — Define release gates based on risk, exploit activity, exposure, and remediation capability.
- [ ] **USEQ-7BE3DC3E** — Keep a secure emergency patch and release path that preserves traceability and review.
- [ ] **USEQ-A8038EF0** — Monitor released components for newly disclosed vulnerabilities and compromised suppliers.
- [ ] **USEQ-A4E69FA6** — Maintain response procedures for leaked secrets, malicious packages, build compromise, and signing-key compromise.
- [ ] **USEQ-000D26B5** — Convert findings and incidents into requirements, tests, patterns, training, and architecture improvements.
- [ ] **USEQ-B634CA7E** — Review security debt and unsupported components as planned lifecycle obligations.
- [ ] **USEQ-D495A353** — Apply equivalent controls to scripts, migrations, infrastructure, pipelines, internal tools, and low-code assets.
- [ ] **USEQ-4C935BDD** — Ensure suppliers can provide sufficient security evidence and notification.
- [ ] **USEQ-C62DB062** — Track security requirements through deprecation and retirement, including data and credential cleanup.
- [ ] **USEQ-96C8AFAD** — Measure security outcomes without rewarding concealment or excessive low-value findings.

## Threat Modeling and Abuse-Case Analysis

_Consolidated from `quality standards/09-security/09-threat-modeling-and-abuse-case-analysis.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-84B924A4** — Define the system, release, feature, or workflow in scope and its assumptions.
- [ ] **USEQ-E7AC5DFE** — Inventory assets, sensitive operations, users, administrators, service identities, suppliers, and attackers.
- [ ] **USEQ-5B4990BE** — Map entry points, data flows, trust boundaries, stores, queues, callbacks, support paths, and recovery mechanisms.
- [ ] **USEQ-BDEAA860** — Include accidental misuse, malicious outsiders, compromised users, insiders, compromised suppliers, automation, and physical access where relevant.
- [ ] **USEQ-92F77DDE** — Model identity spoofing, tampering, disclosure, denial, privilege escalation, repudiation, fraud, and business-logic abuse.
- [ ] **USEQ-CC3D4986** — Model cross-user, cross-tenant, cross-environment, and cross-region isolation failure.
- [ ] **USEQ-72ED476D** — Model account recovery, support intervention, impersonation, and administrative override.
- [ ] **USEQ-1B4D56C1** — Model dependency, package, build, pipeline, artifact, secret, and control-plane compromise.
- [ ] **USEQ-06DD37A7** — Model privacy harm, surveillance, reidentification, inference, coercion, and harmful secondary use.
- [ ] **USEQ-BF3435E9** — Model resource exhaustion, economic denial, quota abuse, scraping, spam, and automation.
- [ ] **USEQ-71F002F7** — Model data corruption, replay, duplication, ordering, race, rollback, and recovery attacks.
- [ ] **USEQ-1962D4C8** — Model observability, alerting, backup, and incident-response blind spots.
- [ ] **USEQ-9F85335D** — Rank scenarios using evidence and local context rather than generic labels alone.
- [ ] **USEQ-E2160FE7** — Convert every material scenario into prevention, detection, response, test, and residual-risk decisions.
- [ ] **USEQ-5568FB20** — Validate assumptions with architecture, source, configuration, deployment, and operator evidence.
- [ ] **USEQ-E965683E** — Use adversarial walkthroughs and abuse stories with product, security, engineering, operations, privacy, and support participants.
- [ ] **USEQ-B3D6D462** — Update the model after material design, integration, identity, data, deployment, or incident changes.
- [ ] **USEQ-242A5671** — Track unresolved threats with owners and deadlines.
- [ ] **USEQ-3F6B5754** — Test mitigations against bypass and alternate paths.
- [ ] **USEQ-2907CE90** — Preserve the threat model as a living decision artifact, not a one-time diagram.

## Application Security Engineering

_Consolidated from `quality standards/09-security/10-application-security-engineering.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-2861D3C4** — Use an explicit application-security verification level appropriate to impact and exposure.
- [ ] **USEQ-F0B58ED0** — Enforce authentication and authorization at every trusted operation and resource boundary.
- [ ] **USEQ-C873B2C6** — Deny by default and prevent object, function, property, and tenant authorization bypass.
- [ ] **USEQ-24970613** — Validate and normalize untrusted input at trusted boundaries.
- [ ] **USEQ-2EDA261C** — Use parameterized or safe structured APIs for interpreters, queries, commands, templates, paths, and parsers.
- [ ] **USEQ-4437D261** — Encode output for the exact destination context and sanitize intentionally allowed active content.
- [ ] **USEQ-69BEF304** — Protect state-changing browser operations against cross-site request forgery.
- [ ] **USEQ-5AAE97D2** — Prevent server-side request forgery through destination, scheme, redirect, DNS, and metadata controls.
- [ ] **USEQ-0D643188** — Prevent insecure deserialization, unsafe dynamic evaluation, parser ambiguity, path traversal, and archive traversal.
- [ ] **USEQ-39C0F69C** — Use safe error handling that does not reveal internals or convert security failure into success.
- [ ] **USEQ-65C105BC** — Protect session identifiers, tokens, recovery material, and sensitive state throughout their lifecycle.
- [ ] **USEQ-6BB86F34** — Apply rate, quota, anti-automation, and resource limits at meaningful identities and operations.
- [ ] **USEQ-B603E870** — Protect uploads, downloads, generated files, and derived content through validation, isolation, access control, and cleanup.
- [ ] **USEQ-CF640F4D** — Apply browser isolation and transport controls appropriate to content and trust.
- [ ] **USEQ-8874907B** — Keep sensitive data out of URLs, logs, analytics, caches, client storage, and support tools unless explicitly justified.
- [ ] **USEQ-2BD39E9A** — Test business invariants, race conditions, duplicate requests, workflow skipping, and partial failure.
- [ ] **USEQ-D7E6AD97** — Protect administrative and support actions with stronger controls and audit.
- [ ] **USEQ-AD54BF40** — Review third-party scripts, components, integrations, and redirects as part of the application attack surface.
- [ ] **USEQ-18CA6B82** — Use production-safe monitoring for authentication anomalies, access denials, control bypass, abuse, and data exfiltration.
- [ ] **USEQ-15FAFA1A** — Retest security after material fixes, configuration changes, migrations, and platform updates.

## Identity, Access, and Privileged Security

_Consolidated from `quality standards/09-security/11-identity-access-and-privileged-security.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-1E4C01BE** — Inventory human, service, workload, device, emergency, supplier, and anonymous identities.
- [ ] **USEQ-C3393B74** — Assign each identity a unique accountable owner and purpose.
- [ ] **USEQ-F14BDA5C** — Use assurance proportional to account-compromise impact and transaction risk.
- [ ] **USEQ-84835ED7** — Require phishing-resistant authentication for privileged and high-impact access where feasible.
- [ ] **USEQ-C4E8C5CD** — Protect enrollment, recovery, authenticator replacement, and identity linking at least as strongly as normal authentication.
- [ ] **USEQ-DB24C217** — Use least privilege, separation of duties, and time-bounded elevation.
- [ ] **USEQ-64636061** — Avoid shared accounts; where unavoidable, control, attribute, rotate, and review their use.
- [ ] **USEQ-32615AD6** — Use short-lived, scoped workload credentials rather than embedded long-lived secrets.
- [ ] **USEQ-C3D28E99** — Verify issuer, audience, subject, signature, algorithm, time, nonce, state, and token type.
- [ ] **USEQ-9EC03586** — Keep delegated authorization scopes minimal and consent meaningful.
- [ ] **USEQ-47AB3CEB** — Revoke sessions, tokens, keys, connections, and cached permissions promptly after access changes.
- [ ] **USEQ-F93D51D5** — Review access periodically according to impact and event-triggered changes.
- [ ] **USEQ-AC84BB6E** — Automate joiner, mover, leaver, supplier, tenant, and service-account lifecycle where reliable.
- [ ] **USEQ-182E040C** — Protect break-glass access, test it, alert on use, and require retrospective review.
- [ ] **USEQ-6BFA62B6** — Log authentication, recovery, privilege, access-policy, impersonation, and sensitive-action events.
- [ ] **USEQ-4CCBD31C** — Detect credential stuffing, brute force, token replay, impossible use, dormant accounts, and anomalous privilege.
- [ ] **USEQ-68CA391B** — Prevent identity existence disclosure and social-engineering bypass through support channels.
- [ ] **USEQ-9D9A377E** — Use independent approval for high-impact access grants and destination changes.
- [ ] **USEQ-8A3AEC6A** — Ensure authorization survives background processing, caching, search, export, and data replication.
- [ ] **USEQ-6F7E9AA4** — Remove orphaned identities, keys, roles, groups, sessions, and entitlements.

## Secrets Management

_Consolidated from `quality standards/09-security/12-secrets-management.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-B68869B4** — Inventory passwords, API keys, tokens, certificates, signing keys, encryption keys, recovery codes, and connection secrets.
- [ ] **USEQ-ABB424A2** — Assign every secret an owner, purpose, scope, environment, consumer, creation date, rotation policy, and revocation path.
- [ ] **USEQ-4FECA33F** — Store secrets in an approved system designed for access control, audit, and lifecycle management.
- [ ] **USEQ-0092B3DF** — Use workload identity and short-lived credentials instead of static secrets where supported.
- [ ] **USEQ-E85957E3** — Keep secrets out of source, history, artifacts, images, frontend bundles, tickets, chat, documentation, tests, analytics, and logs.
- [ ] **USEQ-A8D2A9F6** — Use separate secrets across environments, applications, tenants, and purposes where compromise isolation requires it.
- [ ] **USEQ-B7E2A7B6** — Grant secret access to the smallest identity and operation scope.
- [ ] **USEQ-088D659E** — Authenticate secret consumers and protect delivery in transit and at rest.
- [ ] **USEQ-7BB03237** — Avoid passing secrets through command lines, URLs, process listings, or broadly inherited environment state where exposed.
- [ ] **USEQ-DB8A0FB0** — Rotate secrets automatically where safe and test rotation before relying on it.
- [ ] **USEQ-901C7657** — Support emergency revocation and replacement without an uncontrolled full outage.
- [ ] **USEQ-B886B94F** — Monitor retrieval, failure, unusual use, stale versions, expiry, and rotation drift.
- [ ] **USEQ-2409DA14** — Invalidate old versions after migration and verify no consumers remain.
- [ ] **USEQ-0726DF71** — Scan current content and history for accidental disclosure.
- [ ] **USEQ-5819670C** — Treat deletion from source as insufficient after exposure; rotate or revoke the secret.
- [ ] **USEQ-B48719C8** — Protect backups and recovery material containing secrets.
- [ ] **USEQ-797135EE** — Control and review break-glass secrets separately.
- [ ] **USEQ-B740DAFC** — Avoid sharing secrets between human and machine use.
- [ ] **USEQ-4F6A4888** — Redact secrets consistently in errors, traces, crash dumps, support bundles, and diagnostics.
- [ ] **USEQ-27196AF2** — Exercise compromise scenarios for production credentials and signing keys.

## Software Supply-Chain Security

_Consolidated from `quality standards/09-security/13-software-supply-chain-security.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-7D8D0C8D** — Inventory source repositories, build systems, runners, package sources, registries, artifact stores, signers, deployment systems, and suppliers.
- [ ] **USEQ-847A8B6A** — Protect repository, pipeline, artifact, and signing administration with strong identity and least privilege.
- [ ] **USEQ-D0B3EF17** — Require reviewed source and protected pipeline definitions for trusted releases.
- [ ] **USEQ-92BFE01F** — Separate untrusted contribution workflows from release secrets and trusted builders.
- [ ] **USEQ-647E3B49** — Pin and verify build inputs sufficiently to prevent unexpected substitution.
- [ ] **USEQ-B97B1238** — Use isolated, ephemeral, or clean build environments appropriate to risk.
- [ ] **USEQ-C7DF41EB** — Generate immutable artifact digests and provenance tied to exact source and build inputs.
- [ ] **USEQ-2F8BE0DD** — Sign or attest release artifacts and verify before promotion and deployment.
- [ ] **USEQ-C55F21C3** — Build once and promote the same artifact through environments.
- [ ] **USEQ-AEA71A77** — Generate an SBOM for the final artifact, including transitive and embedded components.
- [ ] **USEQ-CE0ED010** — Verify artifact content against declared inputs and detect unexpected files or dependencies.
- [ ] **USEQ-E19478CB** — Protect caches from cross-trust poisoning and stale unverified substitution.
- [ ] **USEQ-164F96FC** — Restrict build and deploy network access and credentials to minimum need.
- [ ] **USEQ-B3EBD693** — Prevent pull-request or untrusted jobs from modifying trusted release state.
- [ ] **USEQ-7F151E96** — Monitor dependencies and build components for compromise, vulnerability, and ownership change.
- [ ] **USEQ-CB1C22C1** — Retain previous trusted artifacts, provenance, and reproducible rebuild capability.
- [ ] **USEQ-C42CB7AE** — Define revocation, rebuild, resigning, redistributing, and customer-notification procedures.
- [ ] **USEQ-52B43AE6** — Test supply-chain incident response, including compromised package, runner, registry, and signing key.
- [ ] **USEQ-50870162** — Assess supplier security evidence and customer-responsibility boundaries.
- [ ] **USEQ-609A80DF** — Preserve chain-of-custody evidence from source to deployed artifact.

## Security Testing and Assurance

_Consolidated from `quality standards/09-security/14-security-testing-and-assurance.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-E1893821** — Define test scope from the threat model, architecture, data, roles, exposure, and business consequences.
- [ ] **USEQ-5F926D02** — Map applicable security requirements to test procedures and evidence.
- [ ] **USEQ-CF82BB56** — Use complementary source, composition, secret, infrastructure, configuration, artifact, and deployed-behavior analysis.
- [ ] **USEQ-00DFEB02** — Test unauthenticated, authenticated, privileged, support, service, and cross-tenant perspectives as applicable.
- [ ] **USEQ-036EA3AA** — Test normal, alternate, error, recovery, rate, race, replay, and workflow-bypass paths.
- [ ] **USEQ-A957585C** — Test identity enrollment, login, federation, recovery, authenticator change, session, and logout.
- [ ] **USEQ-3FC14ADD** — Test object, property, function, bulk, search, export, background, and administrative authorization.
- [ ] **USEQ-B032A272** — Test injection, unsafe parsing, request forgery, file handling, browser boundaries, and transport configuration.
- [ ] **USEQ-0BC88E26** — Test business abuse, fraud, automation, resource exhaustion, and economic attacks.
- [ ] **USEQ-99477FCA** — Verify findings manually when automated confidence is insufficient.
- [ ] **USEQ-FD9D76B6** — Record exact environment, artifact, configuration, accounts, data, methods, limitations, and evidence.
- [ ] **USEQ-D40A9147** — Protect testing from damaging production data, availability, privacy, or third parties.
- [ ] **USEQ-07ECAEBA** — Require independent penetration testing for high-impact or externally required systems.
- [ ] **USEQ-24C3D367** — Retest remediations and alternate bypass paths.
- [ ] **USEQ-436DFF0D** — Convert confirmed findings into durable regression tests where practical.
- [ ] **USEQ-5B1F49AA** — Prioritize findings using local exploitability, reachability, exposure, impact, and controls.
- [ ] **USEQ-49403BDE** — Track false-positive rationales and accepted risk with evidence and expiry.
- [ ] **USEQ-819E9D01** — Review testing blind spots and inaccessible components explicitly.
- [ ] **USEQ-0424B2D7** — Measure escape and recurrence rather than only scan volume.
- [ ] **USEQ-8D48E2E1** — Verify security controls continuously after deployment where safe.

## Security Monitoring, Detection, and Response

_Consolidated from `quality standards/09-security/15-security-monitoring-detection-and-response.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-C2ED6EAC** — Define high-value assets, attack paths, control failures, and abuse behaviors that require detection.
- [ ] **USEQ-95DF0EAA** — Collect sufficient authentication, authorization, privilege, configuration, secret, data-access, export, and administrative events.
- [ ] **USEQ-907F3DE0** — Preserve actor, target, tenant, source, time, outcome, reason, correlation, and release context.
- [ ] **USEQ-A0B84E24** — Protect telemetry integrity, availability, access, retention, and time synchronization.
- [ ] **USEQ-4321CA03** — Minimize sensitive content while retaining evidence needed for investigation.
- [ ] **USEQ-D1CD522E** — Develop detections from threat scenarios and observed attacker behavior rather than generic log availability.
- [ ] **USEQ-69FBE0AB** — Define expected false-positive, false-negative, latency, and coverage characteristics.
- [ ] **USEQ-32C9B568** — Test detections through simulation, purple-team exercises, and known benign edge cases.
- [ ] **USEQ-9CB52172** — Route alerts to accountable responders with severity, context, and actionable runbooks.
- [ ] **USEQ-9B860463** — Detect disabled controls, missing telemetry, logging gaps, queue backlog, and clock failure.
- [ ] **USEQ-451C60CD** — Correlate identity, endpoint, service, network, cloud, application, data, and supplier evidence where relevant.
- [ ] **USEQ-58BFE5F8** — Support rapid session, token, credential, key, account, integration, and artifact revocation.
- [ ] **USEQ-6EF12932** — Preserve forensic evidence and chain of custody according to need.
- [ ] **USEQ-A7F39100** — Separate containment from permanent remediation and verify safe recovery.
- [ ] **USEQ-D32BBF4E** — Monitor for recurrence after recovery.
- [ ] **USEQ-29DCFDFF** — Define customer, partner, regulator, insurer, and public communication paths.
- [ ] **USEQ-D731ED42** — Measure detection, acknowledgment, containment, recovery, impact, and recurrence.
- [ ] **USEQ-87A291FB** — Review undetected incidents and near misses to improve coverage.
- [ ] **USEQ-330B87B1** — Prevent alert fatigue by retiring low-value detections and tuning with evidence.
- [ ] **USEQ-67F368C0** — Ensure response remains possible when ordinary identity, communication, or cloud control planes fail.

## Zero Trust and Privileged Access

_Consolidated from `quality standards/09-security/16-zero-trust-and-privileged-access.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-47B7B401** — Treat network location and possession of an internal credential as insufficient proof of trust.
- [ ] **USEQ-1BA5440C** — Authenticate and authorize each material access using identity, resource, action, context, and current policy.
- [ ] **USEQ-5C1FB5C9** — Use least privilege for people, services, devices, pipelines, and suppliers.
- [ ] **USEQ-8B647FD9** — Separate ordinary and privileged identities and sessions.
- [ ] **USEQ-59A554AA** — Use time-bounded, request-based elevation for privileged work where practical.
- [ ] **USEQ-D6FDDCAD** — Require stronger authentication and recent verification for high-impact operations.
- [ ] **USEQ-0E951C55** — Segment resources to limit lateral movement and blast radius.
- [ ] **USEQ-F91F2234** — Constrain service-to-service communication to declared paths.
- [ ] **USEQ-A0F37A0E** — Validate device or workload posture only when signals are reliable and privacy-proportionate.
- [ ] **USEQ-629F9DEE** — Continuously revoke or reevaluate access after risk, role, session, device, or policy change.
- [ ] **USEQ-867BEC2B** — Prevent standing broad access through inherited groups and unmanaged service accounts.
- [ ] **USEQ-0B4C0AA8** — Record purpose, approver, scope, start, end, and actions for privileged access.
- [ ] **USEQ-6E8335B2** — Monitor unusual privilege grants, failed elevation, mass access, new destinations, and policy bypass.
- [ ] **USEQ-37F9FAA8** — Protect privileged workstations, administrative interfaces, and emergency channels.
- [ ] **USEQ-8B3A1AA3** — Require dual control for destructive or highly sensitive actions where appropriate.
- [ ] **USEQ-50D97CA4** — Test access decisions during identity, network, policy, and control-plane failure.
- [ ] **USEQ-08E1FA75** — Provide break-glass paths that are isolated, monitored, tested, and reviewed.
- [ ] **USEQ-89AE6890** — Review privileges based on actual use and remove dormant access.
- [ ] **USEQ-620C44CB** — Prevent support impersonation from bypassing user and tenant protections.
- [ ] **USEQ-98013A7E** — Make access architecture understandable and auditable end to end.

## Fraud, Abuse, and Automation Resistance

_Consolidated from `quality standards/09-security/17-fraud-abuse-and-automation-resistance.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-BF9990A3** — Define abuse cases from incentives, value flows, scarce resources, trust signals, and attacker economics.
- [ ] **USEQ-630C616C** — Identify who benefits, who pays, how abuse scales, and which legitimate users resemble attackers.
- [ ] **USEQ-8B9EAEFF** — Protect account creation, login, recovery, promotion, referral, messaging, search, checkout, payout, and support paths as applicable.
- [ ] **USEQ-8904C2F7** — Use layered controls across identity, device, network, behavior, rate, reputation, transaction, and human review.
- [ ] **USEQ-3522FC38** — Apply limits at identities and resources that attackers cannot trivially rotate.
- [ ] **USEQ-BEBDDC75** — Prevent controls from disproportionately blocking accessibility tools, shared networks, travelers, or low-resource users.
- [ ] **USEQ-4D26C1D7** — Avoid relying on opaque risk scores without validation, monitoring, and appeal.
- [ ] **USEQ-0141B321** — Protect anti-abuse thresholds, models, rules, and internal signals from unnecessary disclosure.
- [ ] **USEQ-68B110A3** — Detect distributed, low-and-slow, collusive, and cross-account behavior.
- [ ] **USEQ-7CF8B159** — Make high-value actions idempotent, reversible, delayed, or independently confirmed where practical.
- [ ] **USEQ-D1B44693** — Protect destination changes, payouts, refunds, credits, entitlements, inventory, and promotional value.
- [ ] **USEQ-2E0A4EBF** — Use progressive friction proportionate to risk rather than universal hostile challenges.
- [ ] **USEQ-FF2CEAE4** — Provide safe reporting, blocking, recovery, and support escalation for affected users.
- [ ] **USEQ-EA296705** — Monitor false positives, false negatives, user harm, attacker adaptation, and economic impact.
- [ ] **USEQ-DF44098D** — Test abuse controls adversarially and through realistic load.
- [ ] **USEQ-21CC19A4** — Prevent attackers from using error responses, timing, enumeration, or limits as an oracle.
- [ ] **USEQ-687982AA** — Coordinate fraud, security, trust and safety, support, legal, privacy, and product response.
- [ ] **USEQ-A5455E57** — Preserve evidence according to legal and privacy constraints.
- [ ] **USEQ-2876F29D** — Review whether product incentives or design create the abuse opportunity.
- [ ] **USEQ-987D4170** — Retire ineffective controls and harden successful controls against bypass.

## Cryptographic Agility, Obsolescence, and Long-Term Resilience Master Checklist

_Consolidated from `gap supplement/07-cryptographic-agility-obsolescence-and-long-term-resilience.md`; 202 non-duplicative controls._

### Expanded gap-closure controls

#### Cryptographic governance and inventory

- [ ] **USEQ-6F35FD47** — Assign accountable owners for cryptographic policy, architecture, key management, public-key infrastructure, code signing, trust stores, post-quantum migration, and long-term validation.
- [ ] **USEQ-F2939F3B** — Maintain a cryptographic inventory or bill of materials covering algorithms, modes, protocols, parameters, libraries, providers, hardware, certificates, keys, trust anchors, signatures, hashes, random generators, secret-sharing schemes, formats, and dependent systems.
- [ ] **USEQ-A9BC0B49** — Include cryptography embedded in operating systems, browsers, runtimes, databases, backups, appliances, network devices, identity providers, SaaS services, mobile apps, firmware, build systems, third-party SDKs, documents, archives, and partner integrations.
- [ ] **USEQ-19FCB4C3** — Record purpose, data protected, owner, implementation, location, version, algorithm, key size, protocol, certificate profile, key source, rotation, expiry, dependency, exposure, and replacement path.
- [ ] **USEQ-0AAFC440** — Discover hard-coded, undocumented, shadow, obsolete, proprietary, and supplier-managed cryptography continuously.
- [ ] **USEQ-E2EDFAA3** — Classify cryptographic uses by confidentiality lifetime, integrity lifetime, signature validity period, authentication consequence, exposure, scale, and recovery difficulty.
- [ ] **USEQ-FE5C3FA7** — Map applicable government, sector, contractual, export, certification, and organizational cryptographic requirements.
- [ ] **USEQ-CF4E7B53** — Define approved, transitional, deprecated, prohibited, and exception-only algorithms, protocols, modes, parameters, and providers.
- [ ] **USEQ-18C1D59D** — Review the policy at least when standards, threats, provider support, regulations, or system lifetimes change.
- [ ] **USEQ-CF1B0C98** — Require exceptions to identify affected data, threat horizon, compensating controls, migration date, owner, and expiry.
- [ ] **USEQ-730B2348** — Do not equate use of encryption with complete security; verify authentication, authorization, key protection, metadata leakage, endpoint security, and operational behavior.
- [ ] **USEQ-49F52D91** — Maintain architecture and data-flow diagrams showing where plaintext exists, where cryptographic boundaries begin and end, and which identities and keys are used.
- [ ] **USEQ-6EC323FE** — Ensure cryptographic governance includes availability and recoverability so protection does not make required information permanently inaccessible.

#### Protection lifetime and threat-horizon analysis

- [ ] **USEQ-D02DD651** — Define how long each category of information must remain confidential, authentic, integral, nonrepudiable, verifiable, and recoverable.
- [ ] **USEQ-3F9F8D06** — Consider collection-to-deletion lifetime, backup retention, archive retention, legal hold, intellectual-property value, national-security relevance, health and biometric permanence, and future aggregation.
- [ ] **USEQ-4365268E** — Assess harvest-now-decrypt-later risk when encrypted information may retain value after cryptographically relevant quantum computers become available.
- [ ] **USEQ-98778914** — Prioritize long-lived sensitive data and high-impact signing or identity systems for earlier migration.
- [ ] **USEQ-481596B9** — Distinguish confidentiality transition risk from signature, certificate, software-update, identity, and archival-validation transition risk.
- [ ] **USEQ-1419679C** — Identify data already captured by adversaries or broadly replicated where future decryption would still cause harm.
- [ ] **USEQ-1C1F4FFF** — Account for the time required to inventory, standardize, procure, test, certify, deploy, coordinate, re-encrypt, re-sign, and retire old mechanisms.
- [ ] **USEQ-46FA0922** — Compare required protection lifetime plus migration time with the plausible lifetime of the current cryptography.
- [ ] **USEQ-A18ADDBC** — Define cryptoperiods for keys according to algorithm, usage volume, exposure, compromise impact, operational environment, and applicable guidance.
- [ ] **USEQ-E6C238AF** — Prevent one default rotation period from being applied blindly to every key type and purpose.
- [ ] **USEQ-9E1CDEAF** — Review whether encrypted archives can be migrated before keys, algorithms, software, hardware, or organizational knowledge disappear.
- [ ] **USEQ-FFEC983D** — Document residual risk when immediate migration is impossible and reduce new accumulation under vulnerable mechanisms.

#### Algorithm and protocol lifecycle policy

- [ ] **USEQ-8E388226** — Use algorithms, modes, parameters, protocols, and implementations approved for the actual security purpose and assurance level.
- [ ] **USEQ-E0ABD6F9** — Prefer well-reviewed standard constructions and maintained implementations over custom cryptography.
- [ ] **USEQ-5AC12920** — Prohibit obsolete hashes, ciphers, modes, key sizes, padding schemes, random generators, and protocol versions according to current policy.
- [ ] **USEQ-9E10A498** — Separate encryption, message authentication, digital signature, key agreement, password hashing, key derivation, and random-generation requirements.
- [ ] **USEQ-38360471** — Use authenticated encryption where confidentiality and integrity are both required.
- [ ] **USEQ-EC647484** — Bind relevant context, version, identity, purpose, sequence, and metadata as authenticated data when substitution would be harmful.
- [ ] **USEQ-0A6799B2** — Prevent nonce, initialization-vector, salt, counter, sequence, or ephemeral-key reuse where the construction requires uniqueness.
- [ ] **USEQ-E2E01711** — Use approved entropy sources and handle startup, virtualization, cloning, constrained devices, and entropy failure.
- [ ] **USEQ-C6326840** — Use memory-hard or otherwise current password-hashing approaches with parameters reviewed against current hardware and threat levels.
- [ ] **USEQ-BDD9D750** — Validate certificates, signatures, revocation, names, purposes, algorithms, chains, constraints, timestamps, and trust anchors completely.
- [ ] **USEQ-F4DCC26A** — Disable insecure protocol fallback and prevent downgrade to weak algorithms, parameters, versions, or unauthenticated modes.
- [ ] **USEQ-3B0C74AE** — Use explicit protocol and message versioning and reject ambiguous or unknown critical fields safely.
- [ ] **USEQ-751D6B3F** — Design algorithm negotiation so an active attacker cannot remove stronger choices or select an unintended suite.
- [ ] **USEQ-0FD4DE36** — Do not expose unauthenticated error details that create padding, format, timing, or oracle attacks.
- [ ] **USEQ-5C6B799A** — Minimize distinguishable behavior for cryptographic failure while retaining operational diagnosability through protected telemetry.
- [ ] **USEQ-6B3C1DFB** — Use constant-time or otherwise side-channel-resistant operations where secrets and attacker observation make it necessary.
- [ ] **USEQ-293F39CE** — Review compression, deduplication, length leakage, traffic analysis, access patterns, and metadata exposure around encrypted content.
- [ ] **USEQ-EB0FC491** — Treat format, encoding, parsing, canonicalization, and serialization as part of signature and verification security.

#### Cryptographic agility architecture

- [ ] **USEQ-D27B79E2** — Separate business logic and data formats from one hard-coded algorithm, provider, certificate profile, key size, or cryptographic library.
- [ ] **USEQ-2FCEE2CF** — Use controlled abstraction boundaries that allow safe replacement without hiding security-critical parameters and failures.
- [ ] **USEQ-DBA0CD05** — Store algorithm, version, key identifier, parameters, and format metadata with ciphertext, signatures, tokens, or records when future interpretation requires it.
- [ ] **USEQ-7BEE2035** — Support multiple approved versions during migration without automatically enabling obsolete ones for new protection.
- [ ] **USEQ-C0921D64** — Make encryption and signature formats self-describing enough for controlled evolution but reject algorithm confusion and attacker-selected unsafe parameters.
- [ ] **USEQ-1EA1B2BE** — Avoid algorithm identifiers controlled solely by untrusted input.
- [ ] **USEQ-EC5E35C3** — Centralize policy while avoiding a single runtime or service failure that can disable every critical operation.
- [ ] **USEQ-18D118F0** — Design key and certificate references so rotation does not require rewriting unrelated business identifiers.
- [ ] **USEQ-7D9C45FA** — Separate key material from encrypted data and from ordinary application deployment.
- [ ] **USEQ-D6E1A3D6** — Provide compatibility adapters and explicit migration tools rather than permanent implicit fallback.
- [ ] **USEQ-74D66238** — Design protocols, storage, APIs, devices, and message formats with enough size and negotiation headroom for larger future keys, signatures, ciphertexts, certificates, and handshakes.
- [ ] **USEQ-2400F7A8** — Avoid fixed field sizes, database columns, packet assumptions, QR or barcode limits, and UI constraints that cannot accommodate replacement algorithms.
- [ ] **USEQ-B99BDC20** — Identify latency, bandwidth, storage, memory, CPU, hardware, battery, and handshake impacts before selecting migration mechanisms.
- [ ] **USEQ-F49BC55A** — Maintain test vectors and conformance tests for every supported algorithm and format.
- [ ] **USEQ-11E0B313** — Ensure cryptographic provider replacement does not change canonicalization, error behavior, key interpretation, randomness, or access control unexpectedly.
- [ ] **USEQ-24F0BF3A** — Use configuration and policy to disable a compromised mechanism rapidly without redeploying every dependent component where feasible.
- [ ] **USEQ-BB7B3DA7** — Test policy rollback carefully so emergency recovery cannot silently restore prohibited cryptography.

#### Key, secret, and trust-anchor lifecycle

- [ ] **USEQ-FAA06A19** — Generate keys in an approved environment using sufficient entropy and parameters.
- [ ] **USEQ-C75EB0DA** — Assign every key a unique identity, owner, purpose, algorithm, allowed operations, scope, environment, creation time, activation time, expiry, rotation, archival, recovery, and destruction policy.
- [ ] **USEQ-BA21D2ED** — Prevent one key from being reused across unrelated purposes, tenants, environments, regions, protocols, or trust domains unless explicitly justified.
- [ ] **USEQ-D2104A6B** — Use hardware-backed or isolated key protection when consequence and threat warrant it.
- [ ] **USEQ-8EF69F9B** — Restrict key export and distinguish exportable, non-exportable, wrapped, escrowed, recoverable, ephemeral, and derived keys.
- [ ] **USEQ-74803778** — Apply least privilege and separation of duties to key creation, activation, backup, use, rotation, recovery, revocation, and destruction.
- [ ] **USEQ-4BD8537A** — Use dual control or threshold authorization for high-impact root, signing, recovery, or master keys where appropriate.
- [ ] **USEQ-7C8016A4** — Audit administrative and cryptographic operations without logging key material or reusable secrets.
- [ ] **USEQ-A4A4C02E** — Rotate keys on schedule, before exposure limits, after relevant personnel or provider changes, and immediately upon credible compromise.
- [ ] **USEQ-9EAE18CD** — Support overlapping activation windows so rotation does not cause avoidable outage, while removing old keys after the authorized verification or decryption period.
- [ ] **USEQ-F90D3813** — Verify that every producer and consumer has adopted a new key before retiring the old one.
- [ ] **USEQ-DDAC56E9** — Handle key identifiers, caches, offline clients, queued messages, backups, replicas, and disaster-recovery systems during rotation.
- [ ] **USEQ-20C41B49** — Back up only keys that must be recoverable and protect backup copies at least as strongly as active keys.
- [ ] **USEQ-7BF71BAC** — Test key recovery and ensure recovery cannot bypass authorization, audit, or quorum requirements.
- [ ] **USEQ-A7A6C265** — Destroy keys using a verified method appropriate to media, hardware, replicas, backups, and cryptographic erasure claims.
- [ ] **USEQ-7DDE7DE6** — Define response for lost, unavailable, corrupt, duplicated, weak, misissued, leaked, or unauthorized keys.
- [ ] **USEQ-5CC4E744** — Revoke dependent credentials, certificates, tokens, signatures, and trust when a key is compromised.
- [ ] **USEQ-38DD5B78** — Protect secret zero, bootstrap credentials, root trust, recovery codes, signing keys, and automation identities as special high-risk assets.
- [ ] **USEQ-78D4173D** — Prevent development, test, sample, vendor-default, and demonstration keys from entering production.
- [ ] **USEQ-714052D8** — Ensure personnel departure, supplier termination, account transfer, and organizational change trigger required key and trust review.

#### Certificates, PKI, identity, and trust services

- [ ] **USEQ-C88DE1A5** — Maintain an inventory of private and public certificate authorities, trust anchors, issuing hierarchies, profiles, certificate types, domains, workloads, devices, and relying parties.
- [ ] **USEQ-77D5B2F5** — Define certificate policy, issuance authority, proofing, approval, names, extensions, key usage, validity, renewal, revocation, and audit requirements.
- [ ] **USEQ-AA881D4C** — Automate issuance and renewal where it reduces expiry risk without granting uncontrolled issuance authority.
- [ ] **USEQ-D89D115E** — Monitor certificate and trust-anchor expiration with enough lead time for coordinated replacement.
- [ ] **USEQ-190C6FE9** — Validate hostname or identity, chain, purpose, constraints, policy, name constraints, critical extensions, revocation as required, and algorithm strength.
- [ ] **USEQ-EB35BB62** — Pin trust only where lifecycle, backup, rotation, and failure consequences are fully managed.
- [ ] **USEQ-92905AD7** — Prevent unauthorized certificate issuance through restricted account access, domain controls, certificate transparency monitoring where applicable, and issuance policy.
- [ ] **USEQ-3BF8C390** — Use separate trust domains for production, testing, internal identity, public web, code signing, device identity, and document signing where appropriate.
- [ ] **USEQ-8F57BFF3** — Protect root and intermediate authorities with stronger controls, limited use, offline or highly isolated operation, and ceremony evidence where warranted.
- [ ] **USEQ-F757D85C** — Test CA compromise, misissuance, revocation, trust-store change, expired chain, clock failure, unavailable status service, and emergency replacement.
- [ ] **USEQ-AA2F29A5** — Define how offline, constrained, intermittently connected, and long-lived devices receive trust updates.
- [ ] **USEQ-9F1BA40E** — Prevent an expired or revoked identity from being restored by backup, cached trust, factory reset, or old client software.
- [ ] **USEQ-A882F312** — Record ownership transfers, domain changes, supplier changes, and organizational mergers in trust lifecycle planning.

#### Digital signatures, timestamping, and long-term validation

- [ ] **USEQ-CD2F3826** — Define whether each signature proves integrity, origin, approval, authorization, publication, software provenance, transaction intent, or legal commitment.
- [ ] **USEQ-5237024A** — Bind the signature to the exact canonical content, context, identity, purpose, version, and relevant metadata.
- [ ] **USEQ-7729F205** — Prevent signature wrapping, substitution, partial-content ambiguity, detached-content mismatch, and canonicalization attacks.
- [ ] **USEQ-13CC6B9B** — Validate signer authority at signing time and distinguish identity proof from authorization to approve the specific action.
- [ ] **USEQ-1BDF9889** — Use trusted time evidence when future validation depends on knowing when a signature existed.
- [ ] **USEQ-EE20ACBC** — Preserve certificate chains, revocation evidence, policy identifiers, timestamps, validation data, and signed attributes needed for later verification.
- [ ] **USEQ-A173C360** — Plan re-timestamping, evidence-record renewal, re-signing, or archival protection before algorithms, certificates, or trust services expire.
- [ ] **USEQ-EFB39081** — Do not overwrite an original signature when applying a new preservation signature; maintain the evidence chain.
- [ ] **USEQ-37431A85** — Define how redaction, transformation, format migration, aggregation, and metadata updates affect signature validity.
- [ ] **USEQ-1223A7F4** — Use reproducible or independently verifiable signing pipelines for software, models, packages, documents, and release metadata.
- [ ] **USEQ-B6D084C8** — Protect code-signing and update-signing keys from build workers and untrusted contributions.
- [ ] **USEQ-76246FEA** — Require threshold, review, or release evidence for high-impact signing operations.
- [ ] **USEQ-397EBD90** — Support emergency key revocation, signer replacement, artifact recall, trust-list update, and reissue.
- [ ] **USEQ-D33B37C6** — Make verification failure visible and distinguish unknown, unsupported, expired, revoked, malformed, untrusted, and modified states.
- [ ] **USEQ-E7707180** — Do not label content authentic merely because a cryptographic signature is valid; verify signer, authority, context, and policy.

#### Post-quantum readiness and migration

- [ ] **USEQ-C38475E5** — Maintain a prioritized quantum-readiness roadmap based on cryptographic inventory, protection lifetime, exposure, interoperability, and replacement complexity.
- [ ] **USEQ-052CC191** — Track official standards, implementation guidance, validation programs, protocol adoption, supplier roadmaps, and security analysis for post-quantum mechanisms.
- [ ] **USEQ-1C0AA009** — Plan adoption of standardized post-quantum key-establishment and signature mechanisms appropriate to each use rather than inventing proprietary replacements.
- [ ] **USEQ-B3A0471E** — Treat current NIST post-quantum standards as building blocks whose secure protocol integration, parameter selection, implementation, and operational use still require validation.
- [ ] **USEQ-EA1CA297** — Identify systems that cannot accommodate larger keys, ciphertexts, signatures, certificates, packets, handshakes, firmware, or storage.
- [ ] **USEQ-B05B3A8B** — Prioritize new long-lived sensitive data, high-value identity, software-update, archival-signature, and difficult-to-upgrade systems.
- [ ] **USEQ-75A7AB2C** — Engage hardware, operating-system, browser, network, identity, cloud, SaaS, device, certificate, and partner suppliers on supported algorithms and timelines.
- [ ] **USEQ-7F196A53** — Require procurement and renewal contracts to disclose cryptographic dependencies and provide credible migration support.
- [ ] **USEQ-4B7B4908** — Establish test environments for interoperability among legacy, hybrid, and post-quantum participants.
- [ ] **USEQ-7E780AB8** — Use hybrid mechanisms only with a documented security goal, sound composition, supported standards, bounded transition period, and complete downgrade protection.
- [ ] **USEQ-2D0784C6** — Do not assume hybrid automatically provides the stronger of two components; validate composition and failure behavior.
- [ ] **USEQ-34C19231** — Measure handshake, signing, verification, bandwidth, storage, memory, CPU, energy, queue, timeout, and denial-of-service impact.
- [ ] **USEQ-8E08DDE8** — Test fragmentation, maximum transmission units, certificate-chain size, message limits, proxy behavior, hardware acceleration, and constrained-device operation.
- [ ] **USEQ-36DC6205** — Validate randomness, key generation, side-channel resistance, implementation maturity, error handling, and safe failure.
- [ ] **USEQ-4731F951** — Use published test vectors, interoperability events, independent implementations, and validation programs where available.
- [ ] **USEQ-96AC5FBE** — Plan how encrypted stored data will be re-encrypted or otherwise protected and how old ciphertext and keys will be retired.
- [ ] **USEQ-A5E567FD** — Plan how signatures, certificates, software updates, device roots, archived records, and trust stores will transition without invalidating required evidence.
- [ ] **USEQ-9E674F8C** — Maintain a rollback or containment strategy for implementation defects without returning to permanently weak protection for new data.
- [ ] **USEQ-B15D431D** — Phase migration with inventory completion, lab validation, limited deployment, interoperability evidence, performance validation, broader rollout, and legacy retirement.
- [ ] **USEQ-B2FB844A** — Set dates for stopping new use of vulnerable mechanisms and for removing legacy verification or decryption according to official guidance and local risk.
- [ ] **USEQ-37595B5C** — Do not wait for a cryptographically relevant quantum computer to begin inventory, architecture changes, supplier coordination, and migration testing.
- [ ] **USEQ-4A7A983D** — Do not deploy immature or nonstandard mechanisms to production solely from urgency; balance migration risk against implementation and interoperability risk.

#### Migration engineering and compatibility

- [ ] **USEQ-B6B4F06C** — Create a dependency graph from each cryptographic mechanism to protocols, data, schemas, devices, clients, partners, certificates, keys, libraries, hardware, and recovery systems.
- [ ] **USEQ-67BFEE6A** — Identify long-tail clients, offline devices, embedded systems, exported data, immutable media, third-party archives, and contractual interfaces that cannot change quickly.
- [ ] **USEQ-FBE1B27D** — Define source and target states, dual-operation period, negotiation rules, data conversion, key transition, validation evidence, stop conditions, and retirement criteria.
- [ ] **USEQ-8E5D50F7** — Test new and old versions in every permitted combination and reject unintended downgrade.
- [ ] **USEQ-CACC86AD** — Use staged rollout and telemetry to identify incompatible clients, algorithm selection, failures, latency, and capacity impact.
- [ ] **USEQ-AA101A4B** — Ensure error handling distinguishes compatibility failure from attack, corruption, policy rejection, expired trust, and implementation defect.
- [ ] **USEQ-C7ACDDE3** — Keep migration tooling deterministic, reviewed, resumable, rate-limited, observable, and capable of verifying every converted object.
- [ ] **USEQ-7AE294DA** — Use manifests, counts, checksums, sampling, and independent validation for bulk re-encryption or re-signing.
- [ ] **USEQ-892CCDFE** — Prevent migration from losing original timestamps, metadata, provenance, signatures, access controls, retention, or legal holds.
- [ ] **USEQ-3B54FE0E** — Protect plaintext and old keys during migration and minimize the time both old and new secrets are simultaneously exposed.
- [ ] **USEQ-64274E48** — Test interruption, partial completion, replay, duplicate processing, rollback, provider outage, and disaster recovery during migration.
- [ ] **USEQ-3813F94F** — Coordinate certificate, key, software, protocol, policy, monitoring, documentation, and customer communication changes.
- [ ] **USEQ-06DE13DE** — Retire obsolete code paths, algorithms, keys, trust anchors, certificates, libraries, feature flags, exceptions, and monitoring after successful transition.
- [ ] **USEQ-DC746BF8** — Continue detecting attempted or accidental use of retired cryptography.

#### Supplier, hardware, and platform dependencies

- [ ] **USEQ-4321C4D8** — Require suppliers to identify cryptographic algorithms, libraries, hardware roots, certificates, protocols, key custody, validation status, and update mechanisms used by their products.
- [ ] **USEQ-935A57C7** — Assess whether provider-managed encryption permits customer control, rotation, revocation, export, deletion, recovery, regional separation, and audit as required.
- [ ] **USEQ-3402F9BF** — Verify claims such as hardware-backed, end-to-end encrypted, customer-managed, zero knowledge, quantum safe, validated, or certified against scope and implementation.
- [ ] **USEQ-52940ABD** — Track end-of-support dates for cryptographic modules, hardware security devices, operating systems, firmware, browsers, runtimes, certificate services, and network equipment.
- [ ] **USEQ-CC5DB2BB** — Ensure firmware and trust-store update paths remain available for the required device or system lifetime.
- [ ] **USEQ-DBC1725E** — Prevent unsupported hardware from becoming the only holder of critical non-exportable keys without redundancy and succession planning.
- [ ] **USEQ-4C5507AF** — Review multi-tenant key isolation, provider administrator access, backup, disaster recovery, lawful access, and geographic processing.
- [ ] **USEQ-27797557** — Define provider-exit procedures for keys, certificates, trust, encrypted data, audit records, and verification capability.
- [ ] **USEQ-AEFC3740** — Test whether data remains decryptable and signatures verifiable after provider outage or contract termination where continuity requires it.
- [ ] **USEQ-9D47D1FB** — Use independent escrow or recovery only where justified and govern it as an additional high-value attack surface.
- [ ] **USEQ-851CC8B6** — Monitor supplier security advisories and cryptographic roadmap changes.
- [ ] **USEQ-49331877** — Require material cryptographic provider changes to trigger reassessment and interoperability testing.

#### Cryptographic implementation assurance

- [ ] **USEQ-9E3E0AFD** — Use maintained, approved cryptographic libraries and platform APIs rather than duplicating primitives.
- [ ] **USEQ-DCE13F42** — Pin and verify cryptographic dependencies and monitor vulnerabilities, algorithm changes, validation status, and end of support.
- [ ] **USEQ-B0B4CA24** — Test known-answer vectors, invalid inputs, boundary values, malformed encodings, wrong keys, expired keys, revoked certificates, truncated data, and tampering.
- [ ] **USEQ-20672872** — Test interoperability with independent implementations where standards permit variation.
- [ ] **USEQ-7BF37D43** — Use static analysis, dynamic analysis, fuzzing, side-channel testing, protocol testing, and specialist review proportionate to impact.
- [ ] **USEQ-871E29A7** — Review random generation, nonce management, key derivation, parameter validation, error handling, serialization, canonicalization, and state reuse.
- [ ] **USEQ-F9A6A733** — Prevent debug modes, test hooks, deterministic randomness, fallback keys, insecure examples, and verbose cryptographic errors from entering production.
- [ ] **USEQ-A5548F38** — Review memory, process, crash dump, swap, core dump, log, telemetry, temporary file, and diagnostic exposure of keys and plaintext.
- [ ] **USEQ-AB55133A** — Clear sensitive material when feasible while recognizing compiler, runtime, garbage collection, copy, and hardware limitations.
- [ ] **USEQ-58024100** — Use privilege separation and process isolation for high-value signing and decryption services.
- [ ] **USEQ-A293539B** — Rate-limit expensive verification, handshake, password hashing, key generation, and malformed-input paths to resist denial of service.
- [ ] **USEQ-B28612F3** — Monitor cryptographic failures, negotiation, algorithm use, certificate expiry, key rotation, revocation, signing volume, anomaly, and policy violations.
- [ ] **USEQ-0317726D** — Ensure telemetry does not disclose keys, secrets, plaintext, sensitive identifiers, or reusable tokens.
- [ ] **USEQ-F7D94A9D** — Commission independent cryptographic design and implementation review for novel protocols, high-impact systems, and custom integration.

#### Long-term encrypted data and technology obsolescence

- [ ] **USEQ-92363FC5** — Maintain a technology-obsolescence register covering cryptography, formats, hardware, operating systems, storage, protocols, libraries, providers, and required expertise.
- [ ] **USEQ-EC985241** — Identify protected information that may outlive current software, hardware, identities, certificate authorities, organizations, or legal entities.
- [ ] **USEQ-3218D380** — Preserve algorithm identifiers, parameters, key references, certificates, validation evidence, software, format specifications, and representation information needed for future access.
- [ ] **USEQ-53CC648D** — Avoid encryption schemes whose keys, formats, or provider APIs cannot be migrated or recovered for required long-term records.
- [ ] **USEQ-DDCBBEB7** — Plan periodic re-encryption, re-signing, timestamp renewal, media refresh, format migration, and verification.
- [ ] **USEQ-60A7706B** — Test future access from a clean environment without hidden dependence on a retired production service.
- [ ] **USEQ-01F5E70A** — Retain old decryption or verification capability only in a controlled isolated environment for the authorized period.
- [ ] **USEQ-0DBA5354** — Prevent legacy compatibility from exposing obsolete mechanisms to ordinary network negotiation or new data protection.
- [ ] **USEQ-B2765DBC** — Document when cryptographic erasure is the intended deletion mechanism and verify all key copies, derivations, escrow, backups, and replicas are addressed.
- [ ] **USEQ-661859BF** — Coordinate preservation requirements with privacy deletion, legal holds, records schedules, provider exit, and disaster recovery.
- [ ] **USEQ-E091BF22** — Fund migration before the technology becomes unsupported or unavailable.

#### Compromise, incident response, and recovery

- [ ] **USEQ-4D94F1DC** — Define incident playbooks for key exposure, weak randomness, nonce reuse, certificate misissuance, CA compromise, signing compromise, algorithm break, library vulnerability, trust-store poisoning, and failed rotation.
- [ ] **USEQ-2C522C55** — Maintain protected emergency contacts and authority for revocation, suspension, reissue, algorithm disablement, and customer communication.
- [ ] **USEQ-3824B918** — Be able to identify every asset, ciphertext, signature, certificate, token, package, device, and relying party affected by a compromised key or mechanism.
- [ ] **USEQ-94B5FFD5** — Preserve evidence without continuing unsafe operation.
- [ ] **USEQ-35AB7FB4** — Revoke or distrust compromised material through every relevant online, offline, cached, embedded, partner, backup, and disaster-recovery path.
- [ ] **USEQ-0A440CD8** — Rotate dependent secrets and investigate whether compromise enabled impersonation, decryption, forgery, code execution, or historical exposure.
- [ ] **USEQ-BE1C0A2A** — Re-encrypt or re-sign affected information when doing so reduces residual risk and preserves required evidence.
- [ ] **USEQ-773E6FF9** — Publish advisories and update trust information promptly where external relying parties are affected.
- [ ] **USEQ-D6AA53D3** — Test emergency rotation and recovery regularly, including loss of the primary key-management or identity service.
- [ ] **USEQ-E02F6368** — Prevent old backup restoration from reintroducing revoked keys, certificates, trust anchors, policy, or vulnerable artifacts.
- [ ] **USEQ-90908C45** — Conduct root-cause review and update inventory, architecture, policy, tests, monitoring, and migration plans.

#### Cryptographic release blockers

- [ ] **USEQ-F2957839** — Do not release a system when material cryptographic use, keys, certificates, trust anchors, algorithms, or dependent data cannot be inventoried and owned.
- [ ] **USEQ-A08D470E** — Do not use custom cryptography or a proprietary protocol without compelling necessity and independent specialist review.
- [ ] **USEQ-BD8F255A** — Do not protect long-lived sensitive data with a mechanism whose remaining expected life is shorter than the data's protection lifetime plus migration time without explicit treatment.
- [ ] **USEQ-24D458E7** — Do not proceed when keys are embedded in source, public artifacts, client code, logs, tickets, shared documents, or unapproved storage.
- [ ] **USEQ-1781A243** — Do not deploy when nonce reuse, weak randomness, algorithm confusion, downgrade, incomplete certificate validation, or unauthenticated encryption is known or reasonably suspected.
- [ ] **USEQ-9FD289E2** — Do not rotate or revoke a critical key without verifying recovery, consumer adoption, dependent artifacts, and rollback implications.
- [ ] **USEQ-08D148E1** — Do not claim quantum safety, end-to-end encryption, hardware protection, validation, or certification beyond the exact assessed configuration and scope.
- [ ] **USEQ-D01E15AC** — Do not deploy a post-quantum or hybrid mechanism without standard alignment, interoperability, implementation assurance, performance testing, and downgrade protection.
- [ ] **USEQ-95379EA9** — Do not retire legacy verification or decryption while authorized records, partners, devices, archives, or recovery paths still require it.
- [ ] **USEQ-2E5E1B57** — Do not retain obsolete cryptography indefinitely for convenience; isolate, constrain, monitor, and remove it on a risk-owned schedule.

## Standards and source references

- [SPDX 3.0.1 / ISO/IEC 5962:2021](https://spdx.github.io/spdx-spec/v3.0.1/)
- [CycloneDX 1.7](https://cyclonedx.org/specification/overview/)
- [SLSA Specification 1.2](https://slsa.dev/spec/v1.2/)
- [OpenSSF Security Baseline and Best Practices](https://baseline.openssf.org/)
- [ISO/IEC 27001:2022 — Information security management systems](https://www.iso.org/standard/27001)
- [OWASP Application Security Verification Standard 5.0.0](https://owasp.org/www-project-application-security-verification-standard/)
- [RFC 9110 — HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110)
- [OWASP Web Security Testing Guide 4.2](https://owasp.org/www-project-web-security-testing-guide/)
- [NIST SP 800-218 v1.1 — Secure Software Development Framework](https://csrc.nist.gov/pubs/sp/800/218/final)
- [NIST SP 800-61 Rev. 3 — Incident Response](https://csrc.nist.gov/pubs/sp/800/61/r3/final)
- [NIST Cybersecurity Framework 2.0](https://www.nist.gov/cyberframework)
- [OWASP Top 10 — 2025](https://owasp.org/www-project-top-ten/)
- [ISO 31000:2018 — Risk management guidelines](https://www.iso.org/standard/65694.html)
- [ISO/IEC/IEEE 12207:2026 — Software life cycle processes](https://www.iso.org/standard/90219.html)
- [ISO/IEC/IEEE 42010:2022 — Architecture description](https://www.iso.org/standard/74393.html)
- [ISO/IEC/IEEE 16085:2021 — Life-cycle risk management](https://www.iso.org/standard/74385.html)
- [NIST SP 800-63-4 — Digital Identity Guidelines](https://pages.nist.gov/800-63-4/)
- [RFC 9700 — OAuth 2.0 Security Best Current Practice](https://www.rfc-editor.org/rfc/rfc9700)
- [RFC 8725 — JSON Web Token Best Current Practices](https://www.rfc-editor.org/rfc/rfc8725)
- [ISO/IEC/IEEE 29119-2:2021 — Test processes](https://www.iso.org/standard/79428.html)
- [OpenTelemetry Specifications](https://opentelemetry.io/docs/specs/)
- [NIST Post-Quantum Cryptography Project](https://csrc.nist.gov/projects/post-quantum-cryptography)
- [NIST FIPS 203, 204, and 205 announcement](https://www.nist.gov/news-events/news/2024/08/nist-releases-first-3-finalized-post-quantum-encryption-standards)
- [NIST CSWP 39upd1 — Considerations for Achieving Crypto Agility: Strategies and Practices](https://csrc.nist.gov/pubs/cswp/39/upd1/considerations-for-achieving-crypto-agility/final)
- [CISA Post-Quantum Cryptography Initiative](https://www.cisa.gov/quantum)
- [NIST SP 800-131A Rev. 2 — Transitioning the Use of Cryptographic Algorithms and Key Lengths](https://csrc.nist.gov/pubs/sp/800/131/a/r2/final)
- [NIST SP 800-57 Part 1 Rev. 5 — Key Management](https://csrc.nist.gov/pubs/sp/800/57/pt1/r5/final)
- [NIST Cryptographic Standards and Guidelines](https://csrc.nist.gov/projects/cryptographic-standards-and-guidelines)
- [ISO 14721:2025 — Open archival information system reference model](https://www.iso.org/standard/87471.html)

---

[Previous phase](07-data-and-information-lifecycle.md) · [Next: Phase 9: Privacy and data protection](09-privacy-and-data-protection.md)
