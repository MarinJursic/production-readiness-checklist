# Specialized domains and release assurance

_Phase 16 of 16 in the [complete engineering review](00-overview.md)._

Triggered modules for specialized products and the controls that connect lifecycle work to final production approval.

Assess every control as **Pass**, **Fail**, **Blocked**, or **Not Applicable**. A Pass needs current evidence for the exact product, revision, artifact, configuration, environment, and data state under review.

## Payments, Billing, and Money Movement

_Consolidated from `quality standards/15-conditional-domains/01-payments-billing-and-money-movement.md`; 18 non-duplicative controls._

### Universal controls

- [ ] **USEQ-5FC4DFAA** — Payment credentials never enter logs, analytics, support tools, error reports, or ordinary databases unless explicitly required and protected.
- [ ] **USEQ-04ED44A7** — Amount, currency, tax, discount, exchange-rate, fee, commission, and rounding behavior is correct.
- [ ] **USEQ-333EF029** — Duplicate-charge and duplicate-payout protection is implemented.
- [ ] **USEQ-197BC003** — Payment, refund, payout, and callback requests are idempotent.
- [ ] **USEQ-6592AF43** — Payment webhooks are authenticated or signed and replay-protected.
- [ ] **USEQ-CA5B22D9** — Pending, authorized, captured, settled, failed, canceled, refunded, disputed, reversed, and chargeback states are modeled correctly.
- [ ] **USEQ-98E13A36** — The product does not claim financial success before authoritative confirmation.
- [ ] **USEQ-FB9EBF41** — Reconciliation compares application records with provider, processor, bank, or ledger records.
- [ ] **USEQ-4CDBDA9F** — Full refund, partial refund, reversal, chargeback, dispute, and failed-payment workflows are tested.
- [ ] **USEQ-4E901444** — Retry behavior cannot cause accidental duplicate charges or payouts.
- [ ] **USEQ-8F8CB05C** — Regional authentication, consumer-protection, and financial requirements are applied where relevant.
- [ ] **USEQ-513F01A1** — Fraud controls are monitored and reviewed for discriminatory, accessibility, and customer-harm impact.
- [ ] **USEQ-A47801E4** — Payout-destination and bank-detail changes receive enhanced verification.
- [ ] **USEQ-BE1302B9** — Financial exports and reports are access-controlled and audited.
- [ ] **USEQ-E056DB88** — Provider outage, delay, duplicate event, out-of-order event, and reconciliation failure are tested.
- [ ] **USEQ-F59B2D4E** — Invoices, receipts, statements, tax documents, and disclosures are correct.
- [ ] **USEQ-C12F67B1** — Subscription renewal, trial conversion, proration, upgrade, downgrade, pause, cancellation, grace period, and retry behavior are tested.
- [ ] **USEQ-6DF3B045** — Customer support cannot alter money movement without controlled authorization, separation of duties where required, and audit.

## Real-Time, Streaming, Collaboration, and Event-Driven Features

_Consolidated from `quality standards/15-conditional-domains/04-realtime-streaming-and-collaboration.md`; 12 non-duplicative controls._

### Universal controls

- [ ] **USEQ-40B1DC25** — Connection authentication and authorization remain valid for the lifetime of the connection.
- [ ] **USEQ-2DF2A48C** — Browser-origin and cross-site connection risks are addressed.
- [ ] **USEQ-9A67EA71** — Messages and events are schema-validated.
- [ ] **USEQ-65100565** — Message size, rate, fan-out, subscription count, and connection count are bounded.
- [ ] **USEQ-085B1CFB** — Ordering guarantees and conflict-resolution behavior are documented.
- [ ] **USEQ-AC384A6E** — Duplicate, delayed, reordered, and lost message behavior is handled.
- [ ] **USEQ-61B2D449** — Backpressure prevents memory, queue, and downstream-resource exhaustion.
- [ ] **USEQ-880F8C81** — Presence, typing, cursor, draft, and collaboration data do not leak across users or tenants.
- [ ] **USEQ-1F0175D9** — Permission changes, logout, suspension, and deletion terminate or constrain active streams.
- [ ] **USEQ-50CCE7A2** — Fan-out capacity and hot-room or hot-topic behavior are load-tested.
- [ ] **USEQ-58F7B452** — Event history and replay preserve authorization, retention, deletion, and tenant boundaries.
- [ ] **USEQ-44AC2BF0** — Collaborative edits handle concurrency, conflict, undo, attribution, and recovery correctly.

## Safety-Critical and Physically Consequential Systems

_Consolidated from `quality standards/15-conditional-domains/06-safety-critical-and-physically-consequential-systems.md`; 5 non-duplicative controls._

### Universal controls

- [ ] **USEQ-49E96333** — Establish fail-safe and fail-secure behavior.
- [ ] **USEQ-C010CAA8** — Safety alerts are perceptible, timely, actionable, and accessible.
- [ ] **USEQ-742FE91A** — Safety thresholds are conservative, documented, and validated.
- [ ] **USEQ-2A52C1E5** — Sensor, input, model, communication, timing, power, and actuator failure modes are tested as applicable.
- [ ] **USEQ-69A5E188** — Recall, disablement, rollback, and customer-notification procedures exist.

## Media, Voice, Video, and Live Communications

_Consolidated from `quality standards/15-conditional-domains/07-media-voice-video-and-live-communications.md`; 13 non-duplicative controls._

### Universal controls

- [ ] **USEQ-EBD1622C** — Camera, microphone, screen, and device permissions are requested only when needed.
- [ ] **USEQ-E6602317** — Permission and recording state are visible and revocable.
- [ ] **USEQ-A1A11D01** — Signaling, room membership, invitations, and media access are authenticated and authorized.
- [ ] **USEQ-2534B0B7** — Session and room identifiers are sufficiently unpredictable where required.
- [ ] **USEQ-7F4FC4B0** — Users cannot join, observe, record, or retrieve communications without authorization.
- [ ] **USEQ-FA268344** — Media relay, direct peer connectivity, network-address exposure, and metadata leakage are threat-modeled.
- [ ] **USEQ-033AF627** — Recording requires appropriate consent, indication, retention, access control, and deletion.
- [ ] **USEQ-B9628E94** — Transcription, translation, moderation, and analysis follow privacy and consent requirements.
- [ ] **USEQ-DF4B3471** — Background noise, echo, poor bandwidth, jitter, packet loss, device changes, and reconnection are tested.
- [ ] **USEQ-BD86BC9E** — Participant identity, speaking state, recording state, muting, removal, and access controls are clear and reliable.
- [ ] **USEQ-4C121133** — Abuse-reporting and moderation controls exist where needed.
- [ ] **USEQ-A3195B6F** — Encryption and key-management behavior match the communication risk.
- [ ] **USEQ-51374A62** — Live media cannot expose private network, location, device, or identity information beyond accepted risk.

## Geolocation, Sensors, and Device Capabilities

_Consolidated from `quality standards/15-conditional-domains/08-geolocation-sensors-and-device-capabilities.md`; 12 non-duplicative controls._

### Universal controls

- [ ] **USEQ-90C7782A** — Device capabilities are requested only when necessary.
- [ ] **USEQ-9F507275** — Permission requests are contextual, understandable, revocable, and accessible.
- [ ] **USEQ-967F9F61** — Denial of permission has a safe and usable path.
- [ ] **USEQ-001D20FE** — Precise data is not collected when coarse data is sufficient.
- [ ] **USEQ-CD98743E** — Location and sensor-data retention is minimized.
- [ ] **USEQ-8DA60AC4** — Access is auditable where appropriate.
- [ ] **USEQ-49F1368C** — Shared-device and background-access risks are addressed.
- [ ] **USEQ-12A5A8D0** — Spoofed, stale, inaccurate, unavailable, or manipulated sensor data is handled safely.
- [ ] **USEQ-FD11B478** — Map, geocoding, routing, and location providers receive only approved data.
- [ ] **USEQ-374B2E5E** — Location sharing has clear audience, precision, duration, and expiry controls.
- [ ] **USEQ-49F3B2EC** — Safety implications of exposing home, work, travel, or real-time location are reviewed.
- [ ] **USEQ-948F1E8C** — Permission changes take effect promptly across caches, sessions, background jobs, and integrations.

## Analytics, Experimentation, Advertising, and Attribution

_Consolidated from `quality standards/15-conditional-domains/10-analytics-experimentation-advertising-and-attribution.md`; 15 non-duplicative controls._

### Universal controls

- [ ] **USEQ-DCB37BC2** — Every event, property, identifier, destination, and retention period is inventoried.
- [ ] **USEQ-03BBB392** — Event names and semantics are versioned and documented.
- [ ] **USEQ-E990B2BF** — Analytics never receives secrets, credentials, full payment data, prohibited health data, or unnecessary sensitive content.
- [ ] **USEQ-376714FD** — Consent and opt-out choices are enforced at collection and transmission time.
- [ ] **USEQ-1DB8EA36** — User and tenant isolation is preserved in analytics and reporting.
- [ ] **USEQ-E52D44A3** — Experiment assignment is deterministic and does not bypass authorization, pricing, safety, or legal controls.
- [ ] **USEQ-7E879111** — Experiment combinations and interactions are bounded and tested.
- [ ] **USEQ-9DB9E528** — Experiments have owners, hypotheses, stop criteria, monitoring, and expiry dates.
- [ ] **USEQ-838BB8B2** — Harmful or degraded variants can be disabled immediately.
- [ ] **USEQ-36303028** — Metric definitions resist double counting, bot traffic, missing events, retries, and clock differences.
- [ ] **USEQ-6F0BAC43** — Attribution and advertising code does not alter critical behavior or leak sensitive data.
- [ ] **USEQ-A77F523C** — Session replay and recording tools mask sensitive fields and are privacy-reviewed.
- [ ] **USEQ-5C3A4846** — Data quality, loss, duplication, delay, schema drift, and provider outage are monitored.
- [ ] **USEQ-4BF18A72** — Deletion and privacy requests propagate to analytics providers where required.
- [ ] **USEQ-7C1F74A2** — Reports do not expose small-group or individual data improperly.

## E-Commerce, Orders, Inventory, Fulfillment, and Returns

_Consolidated from `quality standards/15-conditional-domains/11-ecommerce-orders-inventory-fulfillment-and-returns.md`; 12 non-duplicative controls._

### Universal controls

- [ ] **USEQ-FAC20782** — Product, price, tax, discount, availability, inventory, delivery, and return information is accurate.
- [ ] **USEQ-D4A8D9FB** — Inventory reservation, decrement, release, oversell, backorder, and cancellation behavior is correct under concurrency.
- [ ] **USEQ-AA20A64E** — Cart and checkout totals cannot be manipulated by the client.
- [ ] **USEQ-100B6DA2** — Coupon, promotion, referral, loyalty, and gift-card rules cannot be abused or stacked improperly.
- [ ] **USEQ-2F559873** — Order creation is idempotent and reconciled with payment and inventory systems.
- [ ] **USEQ-F072E319** — Partial fulfillment, split shipment, substitution, backorder, cancellation, refund, and return workflows are tested.
- [ ] **USEQ-9818DD41** — Delivery addresses, contact details, and fulfillment instructions are protected.
- [ ] **USEQ-96A83B0F** — Shipping, customs, tax, and regional restrictions are applied correctly.
- [ ] **USEQ-2B9E1A68** — Order status reflects authoritative state and handles delayed or out-of-order provider events.
- [ ] **USEQ-8848D08B** — Marketplace seller and buyer permissions, payouts, disputes, and fees are isolated and audited.
- [ ] **USEQ-E1CC2FBF** — Customer support cannot modify prices, refunds, inventory, or fulfillment without controlled authorization.
- [ ] **USEQ-2DA04192** — Product recalls and safety notifications can identify and contact affected customers.

## Blockchain, Smart Contracts, Digital Assets, and Irreversible Ledgers

_Consolidated from `quality standards/15-conditional-domains/12-blockchain-smart-contracts-and-irreversible-ledgers.md`; 12 non-duplicative controls._

### Universal controls

- [ ] **USEQ-FB281AC4** — On-chain and off-chain components, networks, contracts, bridges, oracles, custodians, signers, and keys are inventoried.
- [ ] **USEQ-4EE7054B** — Smart-contract code and upgrades receive independent review and testing.
- [ ] **USEQ-5A5F7869** — Irreversible actions require clear confirmation, address validation, amount validation, and network validation.
- [ ] **USEQ-367B41AA** — Replay, reorganization, finality, nonce, fee, front-running, sandwiching, and transaction-replacement behavior is addressed.
- [ ] **USEQ-102C77C7** — Private keys and signing devices follow strong key-management controls.
- [ ] **USEQ-416C0FA5** — Multi-party approval is used for high-value or administrative actions where appropriate.
- [ ] **USEQ-76F26E01** — Contract ownership, upgrade authority, pause authority, and emergency controls are documented and restricted.
- [ ] **USEQ-20A50596** — Oracle, bridge, exchange, custodian, and wallet-provider failures are threat-modeled.
- [ ] **USEQ-C38E0EC6** — Chain congestion, network fork, delayed finality, and provider outage are handled.
- [ ] **USEQ-14B6379E** — User balances reconcile across internal records and authoritative ledgers.
- [ ] **USEQ-F6401644** — Token, asset, sanctions, tax, consumer-protection, and financial-regulatory applicability is reviewed.
- [ ] **USEQ-FA3950BD** — Recovery limitations and irreversible-risk warnings are clear and accurate.

## IoT, Device Control, Industrial, and Cyber-Physical Integrations

_Consolidated from `quality standards/15-conditional-domains/13-iot-device-control-industrial-and-cyber-physical.md`; 12 non-duplicative controls._

### Universal controls

- [ ] **USEQ-427D419B** — Devices, firmware, gateways, protocols, commands, sensors, actuators, and ownership states are inventoried.
- [ ] **USEQ-D9C174A7** — Device identity, enrollment, ownership transfer, reset, and decommissioning are secure.
- [ ] **USEQ-0E3D9A51** — Commands are authenticated, authorized, scoped, replay-protected, and auditable.
- [ ] **USEQ-3A5CABA3** — Unsafe or conflicting commands are rejected.
- [ ] **USEQ-922D07FC** — Offline, delayed, duplicated, reordered, and partially executed commands are handled safely.
- [ ] **USEQ-01750D1D** — Device firmware and configuration have trusted update, rollback, and integrity verification.
- [ ] **USEQ-7C316720** — Lost, stolen, compromised, or resold devices can be revoked and reset.
- [ ] **USEQ-72F94CB0** — Local physical access and network exposure are threat-modeled.
- [ ] **USEQ-A69FFAC6** — Safety interlocks remain effective if the web application, network, cloud, or identity provider fails.
- [ ] **USEQ-BC5663A0** — Telemetry accuracy, spoofing, calibration, drift, and missing-data behavior are monitored.
- [ ] **USEQ-958BBF9F** — Device and user privacy, location, audio, video, and household data are minimized and protected.
- [ ] **USEQ-2D51DCAA** — Emergency disablement, field support, recall, and end-of-life procedures exist.

## Open-Source Projects and Public Packages

_Consolidated from `quality standards/15-conditional-domains/14-open-source-projects-and-public-packages.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-6497B434** — Publish a clear license and ensure all included code and assets are legally compatible.
- [ ] **USEQ-D73106BF** — Define governance, maintainer roles, decision process, code of conduct, and conflict resolution.
- [ ] **USEQ-89F5E97A** — Provide contribution, development, testing, review, security, and release guidance.
- [ ] **USEQ-6C1118DF** — Require contributor provenance or attestation appropriate to project risk.
- [ ] **USEQ-94A60F9D** — Protect maintainer accounts, repositories, packages, domains, registries, and signing keys.
- [ ] **USEQ-23087C10** — Use protected branches, independent review, and trusted release automation.
- [ ] **USEQ-888978BC** — Publish release artifacts, checksums, provenance, SBOM, compatibility, and change notes.
- [ ] **USEQ-B88FA04F** — Define supported versions, security maintenance, deprecation, and end-of-life.
- [ ] **USEQ-BC106078** — Provide a private vulnerability-reporting path and coordinated disclosure policy.
- [ ] **USEQ-DBBB57B9** — Triage reports promptly and publish advisories with affected versions and mitigation.
- [ ] **USEQ-A1F8E702** — Avoid exposing reporter identity or exploit detail before users can protect themselves.
- [ ] **USEQ-8BB2D4E1** — Review dependencies, install scripts, generated content, examples, and package metadata.
- [ ] **USEQ-01F8B782** — Prevent package-name takeover, abandoned namespace, and unauthorized release.
- [ ] **USEQ-795582EF** — Track downstream consumers and ecosystem impact for breaking or security changes.
- [ ] **USEQ-C714CD1A** — Use semantic or otherwise documented versioning consistently.
- [ ] **USEQ-93F18228** — Provide reproducible setup and tests for contributors.
- [ ] **USEQ-266D92FB** — Make public documentation accessible and free of secrets or private data.
- [ ] **USEQ-F4CC9891** — Plan maintainer succession and avoid single-person release authority.
- [ ] **USEQ-8A380A72** — Disclose commercial sponsorship, conflicts of interest, and governance changes.
- [ ] **USEQ-ADC4E444** — Archive or transfer abandoned projects responsibly rather than leaving vulnerable packages silently active.

## Mobile, Desktop, and Installed Clients

_Consolidated from `quality standards/15-conditional-domains/15-mobile-desktop-and-installed-clients.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-FBF0A7B0** — Define supported operating systems, versions, device classes, architectures, input modes, and assistive technologies.
- [ ] **USEQ-C56B91A7** — Sign distributable applications and protect signing and store accounts.
- [ ] **USEQ-B7D00E9D** — Use trusted distribution channels and verify update authenticity and integrity.
- [ ] **USEQ-222045D3** — Make updates compatible, resumable, rollback-aware, and safe under interruption.
- [ ] **USEQ-328B6D42** — Define minimum and recommended versions and communicate end-of-support.
- [ ] **USEQ-C7371C71** — Request device permissions only at the moment of need with clear purpose.
- [ ] **USEQ-ADE1D9DE** — Provide functional alternatives when optional permissions are denied.
- [ ] **USEQ-01CFD666** — Protect local databases, files, caches, logs, notifications, clipboard, backups, and screenshots according to sensitivity.
- [ ] **USEQ-8955879F** — Clear or protect user data during logout, account switch, uninstall, device transfer, and remote revocation.
- [ ] **USEQ-B9555C7E** — Handle offline changes, conflict, replay, and synchronization safely.
- [ ] **USEQ-0EC9D8F4** — Prevent stale clients from bypassing current authorization, safety, or data rules.
- [ ] **USEQ-AA7C2DC9** — Validate deep links, inter-application messages, extensions, custom schemes, and shared files.
- [ ] **USEQ-51B9E38A** — Use platform security storage for credentials and keys where appropriate.
- [ ] **USEQ-B5BDC412** — Avoid embedding long-lived secrets in distributed binaries.
- [ ] **USEQ-0EA47DFB** — Detect rooted, jailbroken, tampered, or compromised environments only where signals are reliable and user impact is justified.
- [ ] **USEQ-0E6064FC** — Test battery, memory, storage, network, background, interruption, suspension, and resume behavior.
- [ ] **USEQ-5053A39C** — Provide accessible keyboard, screen-reader, scaling, contrast, and reduced-motion behavior.
- [ ] **USEQ-4C38FC1C** — Respect enterprise management and privacy boundaries where supported.
- [ ] **USEQ-389EAD2A** — Make crash and usage telemetry consented, minimized, redacted, and diagnosable.
- [ ] **USEQ-F43D141E** — Test clean install, upgrade from supported versions, migration, repair, uninstall, and reinstallation.

## Cloud-Native and Managed Platforms

_Consolidated from `quality standards/15-conditional-domains/16-cloud-native-and-managed-platforms.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-B7FB0BC0** — Document the shared-responsibility boundary for every managed service and control.
- [ ] **USEQ-40B56A39** — Separate production, nonproduction, security, logging, and shared-platform accounts or projects according to risk.
- [ ] **USEQ-45C59E66** — Use workload and human identities with least privilege and strong authentication.
- [ ] **USEQ-0A87718E** — Inventory regions, zones, endpoints, public exposure, service policies, keys, and quotas.
- [ ] **USEQ-242E52E4** — Encrypt and control data according to sensitivity and residency requirements.
- [ ] **USEQ-4842963A** — Protect provider control planes, organization roots, billing, DNS, and identity administration.
- [ ] **USEQ-542A6C62** — Use policy and automation to prevent prohibited configurations.
- [ ] **USEQ-8155C7A6** — Monitor control-plane, quota, billing, identity, network, region, and managed-service failures.
- [ ] **USEQ-0DCA89D1** — Design resilience across actual failure domains rather than nominal instance count.
- [ ] **USEQ-B68F5F66** — Test provider-service outage, throttling, regional loss, identity failure, and quota exhaustion.
- [ ] **USEQ-28D7D9E6** — Maintain backups and recovery that are isolated from ordinary account compromise.
- [ ] **USEQ-FE43BE53** — Avoid provider defaults that create public access, weak retention, or broad identity trust.
- [ ] **USEQ-7ACBEBF3** — Version infrastructure and detect drift.
- [ ] **USEQ-43BE2B24** — Review managed-service upgrade, deprecation, and maintenance behavior.
- [ ] **USEQ-A842153A** — Model data egress, inter-region transfer, support access, and supplier subprocessor implications.
- [ ] **USEQ-C2368A23** — Maintain provider support and escalation appropriate to impact.
- [ ] **USEQ-FFE8903C** — Forecast cost and prevent runaway resource creation without disabling critical service.
- [ ] **USEQ-29CD22E1** — Define exit, export, migration, key, DNS, and continuity plans for critical services.
- [ ] **USEQ-A4D95CE3** — Avoid claiming multi-cloud portability that has not been exercised.
- [ ] **USEQ-C0562936** — Include provider incidents and control changes in internal risk and postmortem processes.

## Containers and Orchestration

_Consolidated from `quality standards/15-conditional-domains/17-containers-and-orchestration.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-7D4DA192** — Use minimal trusted base images pinned by immutable identity.
- [ ] **USEQ-A2E62F38** — Build images from reviewed source with provenance and SBOM.
- [ ] **USEQ-726D8BC7** — Scan final images and monitor them after release.
- [ ] **USEQ-A3559C64** — Run workloads as non-root and without unnecessary privilege, capabilities, host access, or writable filesystem.
- [ ] **USEQ-13CBD6C6** — Define resource requests, limits, priorities, disruption, and eviction behavior.
- [ ] **USEQ-B1114E9E** — Use workload identities rather than shared node or cluster credentials.
- [ ] **USEQ-43AEBF8E** — Keep secrets out of images and ordinary environment output.
- [ ] **USEQ-84D35810** — Restrict network communication to declared flows.
- [ ] **USEQ-071F5865** — Protect admission, policy, registry, control-plane, and cluster administration.
- [ ] **USEQ-D72B3DB5** — Verify image signatures or attestations before admission where risk warrants.
- [ ] **USEQ-2C2AF8DA** — Separate tenants and trust levels across namespaces, nodes, clusters, or accounts according to isolation needs.
- [ ] **USEQ-901CB8A9** — Define readiness, liveness, startup, graceful termination, and drain behavior.
- [ ] **USEQ-5154706D** — Test rescheduling, rolling update, node loss, zone loss, storage failure, and control-plane degradation.
- [ ] **USEQ-CD9A34E1** — Protect persistent volumes, snapshots, backups, and restoration.
- [ ] **USEQ-7E6911D8** — Keep cluster, node, runtime, image, policy, and add-on versions supported and patched.
- [ ] **USEQ-4E9A0CB6** — Prevent untrusted workloads from reaching metadata, host sockets, control APIs, or other tenants.
- [ ] **USEQ-41A84AD1** — Audit privileged operations, exec, port-forward, secret access, and policy changes.
- [ ] **USEQ-8EAA66DC** — Bound autoscaling and prevent scale storms against dependencies.
- [ ] **USEQ-884ADA51** — Detect configuration drift and unmanaged workloads.
- [ ] **USEQ-472604C8** — Practice cluster rebuild and workload recovery from trusted source.

## Serverless and Managed Execution

_Consolidated from `quality standards/15-conditional-domains/18-serverless-and-managed-execution.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-2BE63A79** — Inventory every function, trigger, event source, destination, identity, secret, region, and failure target.
- [ ] **USEQ-1A3AAD5A** — Use a separate least-privilege identity per workload or cohesive group.
- [ ] **USEQ-7FAE65A8** — Assume event retries and duplicate delivery unless guarantees prove otherwise.
- [ ] **USEQ-21DF7AD2** — Make handlers idempotent and safe under concurrent execution.
- [ ] **USEQ-AAC6692F** — Define timeout, memory, payload, temporary storage, connection, and concurrency limits.
- [ ] **USEQ-4919EBF1** — Protect downstream systems from burst scaling and retry amplification.
- [ ] **USEQ-2916316A** — Configure dead-letter, failure destination, replay, and reconciliation.
- [ ] **USEQ-C6F8CEBE** — Keep deployment packages minimal, signed or attested, and dependency-scanned.
- [ ] **USEQ-28FDB25A** — Do not place secrets in code, packages, logs, or event payloads.
- [ ] **USEQ-6843A15C** — Control public invocation, cross-account triggers, and event source permissions.
- [ ] **USEQ-4B7E5463** — Measure cold-start and scale behavior against user objectives.
- [ ] **USEQ-B43F9D53** — Reuse connections safely without assuming process lifetime.
- [ ] **USEQ-C47082A7** — Handle process reuse so one invocation cannot leak state to another user or tenant.
- [ ] **USEQ-AFE31124** — Propagate trace, identity, tenant, deadline, and correlation context.
- [ ] **USEQ-A5E51F84** — Monitor throttling, duration, errors, retries, age, concurrency, cost, and downstream saturation.
- [ ] **USEQ-C27ADF5F** — Test provider outage, delayed events, partial region failure, and control-plane failure.
- [ ] **USEQ-E28C525F** — Define versioning, aliases, traffic shift, rollback, and event compatibility.
- [ ] **USEQ-2BFBFB1A** — Avoid unbounded invocation chains and recursive triggers.
- [ ] **USEQ-BACD0649** — Forecast cost amplification from abuse, loops, and large events.
- [ ] **USEQ-6D182068** — Maintain a migration and continuity plan for critical provider-specific execution.

## Healthcare and Regulated Sensitive Systems

_Consolidated from `quality standards/15-conditional-domains/19-healthcare-and-regulated-sensitive-systems.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-2E4F8C0C** — Determine every applicable healthcare, medical-device, research, records, privacy, security, and professional obligation.
- [ ] **USEQ-7EEBE134** — Classify clinical, patient, subject, caregiver, provider, identity, insurance, and research data.
- [ ] **USEQ-EAE30A78** — Apply minimum-necessary access and purpose-based disclosure.
- [ ] **USEQ-CC64FD64** — Verify patient, professional, organization, device, and service identity proportionately.
- [ ] **USEQ-5136CD3A** — Protect emergency access while logging, limiting, and reviewing its use.
- [ ] **USEQ-053C255E** — Maintain complete tamper-evident audit of access, correction, disclosure, order, result, and administrative action where required.
- [ ] **USEQ-D2E3E2DF** — Preserve data provenance, units, reference ranges, time, author, status, correction, and version.
- [ ] **USEQ-8A2B11D6** — Prevent stale, incomplete, duplicate, misattributed, or unit-confused information from appearing authoritative.
- [ ] **USEQ-5C75B50F** — Use controlled vocabularies and interoperability standards where applicable.
- [ ] **USEQ-2855968A** — Validate clinical rules, calculations, alerts, and decision support with qualified domain experts.
- [ ] **USEQ-5F39178C** — Prevent alert fatigue and unsafe suppression.
- [ ] **USEQ-781C2D61** — Separate advisory output from authoritative diagnosis or order when appropriate.
- [ ] **USEQ-7ABEA012** — Provide human review, correction, escalation, and downtime procedures.
- [ ] **USEQ-AC848C87** — Test continuity during network, provider, device, identity, and regional failure.
- [ ] **USEQ-92E715A2** — Protect research consent, protocol, cohort, deidentification, and data-use restrictions.
- [ ] **USEQ-B931B65F** — Validate data migration and reconciliation with clinical or regulated invariants.
- [ ] **USEQ-5A73044C** — Train users for high-impact workflows and assess competence where required.
- [ ] **USEQ-A07EC898** — Maintain incident, breach, safety, and regulatory reporting procedures.
- [ ] **USEQ-99FF0E3C** — Use independent validation and documented change control for high-impact functionality.
- [ ] **USEQ-7AFC88F2** — Do not deploy when patient or subject harm cannot be bounded, detected, and mitigated.

## Public-Sector and High-Assurance Systems

_Consolidated from `quality standards/15-conditional-domains/20-public-sector-and-high-assurance-systems.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-0A7E51B2** — Determine applicable public-records, procurement, accessibility, identity, security, privacy, localization, and continuity obligations.
- [ ] **USEQ-02E0F51C** — Design for equitable access across disability, language, device, bandwidth, literacy, geography, and lack of commercial identity.
- [ ] **USEQ-80BCEB64** — Provide non-digital or assisted alternatives where exclusion would deny essential service.
- [ ] **USEQ-3358D1EF** — Use transparent eligibility, decision, appeal, correction, and complaint processes.
- [ ] **USEQ-9D2A742A** — Avoid requiring unnecessary private accounts, tracking, or third-party platforms.
- [ ] **USEQ-08693ABE** — Protect high-value identity, benefit, tax, justice, licensing, voting, or civic records according to impact.
- [ ] **USEQ-D39A0138** — Provide complete auditability and records retention without exposing sensitive information.
- [ ] **USEQ-E3C13E8C** — Use strong supplier transparency, data ownership, exit, and continuity terms.
- [ ] **USEQ-0E684597** — Avoid proprietary lock-in that prevents public accountability or service continuity without explicit justification.
- [ ] **USEQ-0701D270** — Publish service status, maintenance, accessibility, privacy, security contact, and support information.
- [ ] **USEQ-A0B4A6E4** — Test peak civic events, deadlines, emergencies, and coordinated abuse.
- [ ] **USEQ-E1AB8D23** — Design for operation during disasters, regional failure, identity outage, and supplier failure.
- [ ] **USEQ-39F51DEC** — Protect public interfaces against denial, defacement, disinformation, impersonation, and data manipulation.
- [ ] **USEQ-24FDB6CD** — Use independent assurance for high-impact releases and security claims.
- [ ] **USEQ-9806BBFD** — Make algorithmic or rules-based decisions explainable and contestable where required.
- [ ] **USEQ-59716E09** — Maintain multilingual and accessible notices and forms.
- [ ] **USEQ-797FD404** — Preserve evidence for audit, oversight, legal challenge, and historical record.
- [ ] **USEQ-7C31FB6D** — Separate political or administrative preference from authorized policy and law.
- [ ] **USEQ-979BA5B0** — Review disparate impact and barriers to essential access.
- [ ] **USEQ-459E7E15** — Do not launch an essential service without demonstrated continuity and support capacity.

## Libraries, SDKs, CLIs, and Developer Tools

_Consolidated from `quality standards/15-conditional-domains/21-libraries-sdks-clis-and-developer-tools.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-5DB27264** — Define supported languages, runtimes, platforms, package managers, shells, architectures, and versions.
- [ ] **USEQ-D73379F6** — Use stable idiomatic contracts appropriate to each target ecosystem.
- [ ] **USEQ-36FDC9A4** — Keep default behavior safe, predictable, non-destructive, and suitable for automation.
- [ ] **USEQ-6F66FA4B** — Separate machine-readable output from human-readable diagnostics.
- [ ] **USEQ-4C2BBB0D** — Use meaningful exit status, structured errors, and actionable messages.
- [ ] **USEQ-09D01674** — Support noninteractive, quiet, dry-run, confirmation, timeout, retry, and cancellation behavior as appropriate.
- [ ] **USEQ-AAC6303F** — Do not expose secrets in command history, process listings, logs, examples, or error output.
- [ ] **USEQ-D027BA27** — Validate paths, URLs, configuration, environment, and untrusted project content.
- [ ] **USEQ-939DF129** — Avoid executing project or remote code implicitly.
- [ ] **USEQ-CC6BECF8** — Protect update, plugin, extension, and package-loading mechanisms.
- [ ] **USEQ-BB61ED88** — Version contracts and publish compatibility and deprecation policy.
- [ ] **USEQ-F8B8E814** — Provide migration guidance and test against representative downstream consumers.
- [ ] **USEQ-C6981D19** — Generate signed or attested packages with provenance, SBOM, and checksums.
- [ ] **USEQ-09B1BA8F** — Keep examples executable and documentation synchronized.
- [ ] **USEQ-D91045B1** — Provide deterministic behavior and output where automation depends on it.
- [ ] **USEQ-9AF6A86A** — Respect platform conventions without silently changing semantics across platforms.
- [ ] **USEQ-D2CD8EC4** — Make telemetry optional, transparent, minimized, and safe.
- [ ] **USEQ-933B0296** — Test fresh install, upgrade, downgrade where supported, offline, proxy, permission, and restricted-network behavior.
- [ ] **USEQ-AC6AC548** — Provide security reporting and supported-version information.
- [ ] **USEQ-65DF370E** — Avoid making a convenience tool a hidden production single point of failure.

## Games, Immersive, and High-Engagement Products

_Consolidated from `quality standards/15-conditional-domains/22-games-immersive-and-high-engagement-products.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-77183203** — Define supported devices, controls, displays, performance modes, network conditions, and accessibility options.
- [ ] **USEQ-CC1754CB** — Maintain fair authoritative state for competitive or economic outcomes.
- [ ] **USEQ-4AEFAFB8** — Protect against cheating, automation, item duplication, economy manipulation, and account theft.
- [ ] **USEQ-96AE9F3B** — Validate matchmaking, ranking, rewards, randomization, and progression for integrity and unintended bias.
- [ ] **USEQ-925DD3FB** — Disclose probabilities, purchases, recurring costs, virtual currency, expiration, and irreversible choices clearly.
- [ ] **USEQ-FFAB4E8F** — Protect children and vulnerable users from exploitative monetization and social interaction.
- [ ] **USEQ-3DEFD074** — Provide controls for harassment, blocking, reporting, moderation, voice, chat, and user-generated content.
- [ ] **USEQ-7473B28D** — Provide motion, flashing, color, audio, input, text, difficulty, and cognitive accessibility options.
- [ ] **USEQ-F7140184** — Avoid forcing sensory effects or physical movement without alternatives and safety warnings.
- [ ] **USEQ-F84463E8** — Preserve progress across crash, disconnect, patch, device change, and account recovery.
- [ ] **USEQ-7083EEB1** — Test latency, synchronization, prediction, rollback, reconnect, and regional capacity.
- [ ] **USEQ-7A4816F3** — Protect local save, cloud save, downloadable content, mods, and plugins.
- [ ] **USEQ-5002400A** — Use anti-cheat controls proportionate to privacy, device control, and false-positive harm.
- [ ] **USEQ-7FABC8C3** — Provide appeal and recovery for enforcement actions.
- [ ] **USEQ-2A00C28A** — Monitor toxic behavior, fraud, economy inflation, crash, performance, and player safety.
- [ ] **USEQ-C873F83A** — Avoid deceptive engagement loops and make time or spending controls available where appropriate.
- [ ] **USEQ-C6F35C3F** — Respect streaming, recording, spectator, and privacy preferences.
- [ ] **USEQ-F2EE2423** — Test long-session heat, battery, memory, network, and motion comfort.
- [ ] **USEQ-D2EEA3DA** — Plan seasonal content, live operations, event rollback, and end-of-service continuity.
- [ ] **USEQ-EA1879A2** — Do not launch high-stakes economies or social features without abuse and support readiness.

## Release Control Record

_Consolidated from `quality standards/16-release-production/02-release-control-record.md`; 18 non-duplicative controls._

### Release identity

- [ ] **USEQ-7EFB51F6** — Record exact source commits and release tags.
- [ ] **USEQ-8C855A06** — Record immutable digests for every deployable service, worker, job, frontend bundle, and infrastructure artifact.
- [ ] **USEQ-18DBA736** — Record production configuration versions and feature-flag states.
- [ ] **USEQ-E03B5E82** — Record target regions, accounts, projects, environments, clusters, and domains.
- [ ] **USEQ-EDFA793D** — Record the deployment date and change window.
- [ ] **USEQ-A68072E7** — Record the previous known-good artifact and configuration.
- [ ] **USEQ-B1A0C473** — Record the rollback artifact or tested roll-forward plan.
- [ ] **USEQ-D1D57EA1** — Confirm source, artifacts, configuration, migrations, tests, and approvals all refer to the same release.

### Scope inventory

- [ ] **USEQ-8B5934D8** — List every API, webhook, callback, streaming connection, and machine-to-machine interface.
- [ ] **USEQ-7F798459** — List every worker, queue, scheduled task, batch process, administrative tool, and support tool.
- [ ] **USEQ-9EE38C0B** — List every database, object store, search index, cache, message broker, analytics store, and backup location.
- [ ] **USEQ-D8184B36** — List every identity, payment, communications, analytics, advertising, experimentation, support, logging, monitoring, CDN, DNS, hosting, and infrastructure provider.
- [ ] **USEQ-305EF1D6** — List all user types, administrative roles, support roles, service accounts, external actors, and partner systems.
- [ ] **USEQ-C7F3DF66** — List supported countries, languages, locales, currencies, time zones, browsers, devices, and viewport classes.
- [ ] **USEQ-CEECE86D** — List accessibility conformance targets.
- [ ] **USEQ-D81B1472** — List all data classifications and regulated or contractually restricted data.
- [ ] **USEQ-282C5B2F** — List excluded components and explain why they are outside the assessment.
- [ ] **USEQ-3D969CA9** — Identify shared platforms and dependencies that can affect the application even when another team operates them.

## Deployment and Release Engineering

_Consolidated from `quality standards/16-release-production/03-deployment-and-release-engineering.md`; 18 non-duplicative controls._

### Universal controls

- [ ] **USEQ-8CF2C791** — The deployment plan identifies exact artifacts, target environments, migrations, configuration changes, infrastructure changes, feature flags, and external dependencies.
- [ ] **USEQ-34987B94** — The deployer cannot silently substitute an unapproved artifact.
- [ ] **USEQ-2222C9D0** — Predeployment checks validate environment health, data state, provider health, quotas, and available capacity.
- [ ] **USEQ-CF32EB78** — Migration steps and ordering are explicit.
- [ ] **USEQ-EB431FD0** — Application, worker, scheduled-job, event, and schema compatibility is maintained during rollout.
- [ ] **USEQ-321B8525** — Progressive, canary, blue-green, cohort, regional, or equivalent controlled rollout is used where risk warrants it.
- [ ] **USEQ-DC6E1BDA** — Health gates examine user-facing success, latency, errors, saturation, data integrity, security, abuse, and business outcomes.
- [ ] **USEQ-58678F75** — Rollout pauses automatically or operationally on predefined signals.
- [ ] **USEQ-BBD42F35** — Stop, rollback, roll-forward, and incident-declaration criteria are documented before launch.
- [ ] **USEQ-1409AA23** — Rollback can be initiated quickly and has been rehearsed.
- [ ] **USEQ-9EEEEBB8** — Rollback is compatible with database, event, cache, and interface changes.
- [ ] **USEQ-99B06FF3** — Feature flags or kill switches can independently disable risky behavior where appropriate.
- [ ] **USEQ-3E2AD22B** — Traffic shifting, load-balancer behavior, cache invalidation, CDN invalidation, DNS changes, and certificate changes are planned and tested.
- [ ] **USEQ-96297CBB** — Rate limits, quotas, and capacities are adjusted where needed.
- [ ] **USEQ-3AC480C5** — No critical step exists only in someone's memory.
- [ ] **USEQ-4C6B2EA9** — Manual steps are minimized, scripted where practical, and peer-checked where material.
- [ ] **USEQ-7BC2F497** — The full procedure has been rehearsed in a production-like environment.
- [ ] **USEQ-161C6591** — Engineering, operations, security, data, support, communications, and provider contacts are available as needed.

## Post-Deployment Verification

_Consolidated from `quality standards/16-release-production/04-post-deployment-verification.md`; 14 non-duplicative controls._

### Universal controls

- [ ] **USEQ-CA947752** — Verify the deployed artifact digest, configuration, feature-flag state, schema version, and migration completion.
- [ ] **USEQ-5E936430** — Test login, logout, account recovery, and privileged authentication.
- [ ] **USEQ-C2B79849** — Test authorization with multiple roles and tenants where applicable.
- [ ] **USEQ-A5477858** — Verify public, authenticated, administrative, API, worker, scheduled-job, and batch paths.
- [ ] **USEQ-F30D75DE** — Verify payment, webhook, email, SMS, push, identity, storage, search, analytics, and other critical integrations.
- [ ] **USEQ-CA491689** — Verify transport security, certificates, redirects, cross-origin behavior, cache controls, and security headers.
- [ ] **USEQ-CCE6C371** — Verify CDN and application-cache behavior.
- [ ] **USEQ-3904EEC5** — Review errors, latency, traffic, saturation, queues, database health, replication, and dependency metrics.
- [ ] **USEQ-2A66A878** — Review logs and traces for unexpected exceptions, data leakage, authorization failures, and abnormal behavior.
- [ ] **USEQ-8780BABB** — Verify alerts, synthetic monitoring, dashboards, and paging.
- [ ] **USEQ-6FAC2881** — Confirm support receives expected telemetry and current documentation.
- [ ] **USEQ-6DFB4E7F** — Confirm no unexpected increase in abuse, fraud, failed payments, duplicate operations, data-integrity failures, or support volume.
- [ ] **USEQ-46A643BE** — Roll back, roll forward, or declare an incident rather than rationalizing unexplained anomalies.
- [ ] **USEQ-23CEBF4F** — Record the final go, hold, rollback, roll-forward, or incident decision.

## Required Sign-Offs

_Consolidated from `quality standards/16-release-production/05-required-signoffs.md`; 13 non-duplicative controls._

### Universal controls

- [ ] **USEQ-603F2296** — Product owner signs off on scope, behavior, acceptance criteria, launch objectives, and customer impact.
- [ ] **USEQ-DFE5DFDB** — Engineering owner signs off on architecture, implementation, tests, migrations, maintainability, and technical risk.
- [ ] **USEQ-9A317D5C** — Quality owner signs off on test completeness, evidence, defects, and regression status.
- [ ] **USEQ-4A5E2EB1** — Security owner signs off on threat modeling, application-security verification, findings, supply chain, and residual security risk.
- [ ] **USEQ-3077E62C** — Reliability or operations owner signs off on SLOs, capacity, observability, on-call, rollout, rollback, and recovery.
- [ ] **USEQ-F81F6CE4** — Data owner signs off on integrity, migration, retention, reconciliation, and restoration.
- [ ] **USEQ-17F410B4** — Privacy owner signs off on data inventory, purpose, rights, notices, transfers, and privacy risk.
- [ ] **USEQ-1C6662C7** — Accessibility owner or qualified reviewer signs off on target conformance and unresolved barriers.
- [ ] **USEQ-92006B53** — Legal or compliance owner signs off on applicable obligations, terms, contracts, regulatory risk, and required notices.
- [ ] **USEQ-615F9CB7** — Support owner signs off on support readiness, escalation, documentation, and customer communications.
- [ ] **USEQ-2A0C5DBF** — Business risk owner explicitly accepts remaining material risk.
- [ ] **USEQ-3437198D** — Executive or designated release authority makes the final decision for high-impact launches.
- [ ] **USEQ-9A466E24** — No person signs an area for which evidence is absent or for which they are not qualified or authorized.

## Go, Conditional-Go, and No-Go Decision Rule

_Consolidated from `quality standards/16-release-production/06-go-conditional-go-and-no-go.md`; 12 non-duplicative controls._

### GO

- [ ] **USEQ-A4723C18** — Every critical requirement has current evidence tied to the exact production artifact and configuration.
- [ ] **USEQ-98349D6B** — Critical journeys meet approved correctness, performance, security, privacy, accessibility, and reliability objectives.
- [ ] **USEQ-4CECD69F** — No known residual risk exceeds the organization's approved risk appetite.
- [ ] **USEQ-A8D674CC** — Every accepted risk has an accountable owner, compensating control, monitoring, remediation plan, and expiry.
- [ ] **USEQ-3FE0D592** — Rollout, rollback or roll-forward, restoration, incident response, and customer communication are demonstrably workable.
- [ ] **USEQ-4084EC49** — Required operational, support, security, legal, privacy, and business personnel are ready.

### CONDITIONAL GO

- [ ] **USEQ-A4B3AB62** — User, security, privacy, integrity, financial, legal, accessibility, reliability, and safety impact is bounded and understood.
- [ ] **USEQ-87C1B515** — The affected feature can be disabled, isolated, rolled back, or rolled forward safely.

### NO-GO

- [ ] **USEQ-88CFEFEC** — Evidence is missing, stale, contradictory, or tied to a different artifact.
- [ ] **USEQ-F96A1B86** — Rollback, safe roll-forward, restoration, or incident response is unavailable.
- [ ] **USEQ-F38CE10F** — Required legal, regulatory, contractual, privacy, accessibility, safety, or security approval is absent.
- [ ] **USEQ-6E882462** — The team relies on "we will monitor it" without a concrete, tested detection and mitigation path.

## Standards and source references

- [ISO/IEC 27001:2022 — Information security management systems](https://www.iso.org/standard/27001)
- [OWASP Application Security Verification Standard 5.0.0](https://owasp.org/www-project-application-security-verification-standard/)
- [ISO/IEC 25010:2023 — Product quality model](https://www.iso.org/standard/78176.html)
- [ISO/IEC/IEEE 15288:2023 — System life cycle processes](https://www.iso.org/standard/81702.html)
- [ISO 31000:2018 — Risk management guidelines](https://www.iso.org/standard/65694.html)
- [ISO/IEC 27701:2025 — Privacy information management systems](https://www.iso.org/standard/85819.html)
- [ISO/IEC 29100:2024 — Privacy framework](https://www.iso.org/standard/85938.html)
- [NIST Privacy Framework 1.0](https://www.nist.gov/privacy-framework)
- [ISO/IEC/IEEE 29148:2018 — Requirements engineering](https://www.iso.org/standard/72089.html)
- [ISO 22301:2019 — Business continuity management systems](https://www.iso.org/standard/75106.html)
- [OpenSSF Security Baseline and Best Practices](https://baseline.openssf.org/)
- [SLSA Specification 1.2](https://slsa.dev/spec/v1.2/)
- [SPDX 3.0.1 / ISO/IEC 5962:2021](https://spdx.github.io/spdx-spec/v3.0.1/)
- [NIST SP 800-218 v1.1 — Secure Software Development Framework](https://csrc.nist.gov/pubs/sp/800/218/final)
- [NIST Cybersecurity Framework 2.0](https://www.nist.gov/cyberframework)
- [Google SRE Workbook — Implementing SLOs](https://sre.google/workbook/implementing-slos/)
- [ISO/IEC 38500:2024 — Governance of IT](https://www.iso.org/standard/81684.html)
- [W3C Web Content Accessibility Guidelines 2.2](https://www.w3.org/TR/WCAG22/)
- [ISO/IEC/IEEE 26514:2022 — Information for users](https://www.iso.org/standard/77451.html)
- [ISO 9241-210:2019 — Human-centred design](https://www.iso.org/standard/77520.html)
- [ISO/IEC/IEEE 12207:2026 — Software life cycle processes](https://www.iso.org/standard/90219.html)
- [ISO/IEC/IEEE 15289:2019 — Life-cycle information items](https://www.iso.org/standard/74909.html)
- [DORA software delivery performance research](https://dora.dev/guides/dora-metrics/)
- [ISO/IEC/IEEE 29119-2:2021 — Test processes](https://www.iso.org/standard/79428.html)
- [ISO 9001:2015 — Quality management systems](https://www.iso.org/standard/62085.html)
- [ISO/IEC/IEEE 16085:2021 — Life-cycle risk management](https://www.iso.org/standard/74385.html)

---

[Previous phase](15-ai-ml-and-ai-assisted-development.md) · [Next: Production readiness: release foundations](../checklists/01-release-foundations.md)
