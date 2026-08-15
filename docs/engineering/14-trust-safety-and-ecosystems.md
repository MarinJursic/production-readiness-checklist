# Trust, safety, and ecosystems

_Phase 14 of 16 in the [complete engineering review](00-overview.md)._

Safety-by-design, abuse response, content integrity, authenticity, recommendations, notifications, extensions, plugins, and marketplaces.

Assess every control as **Pass**, **Fail**, **Blocked**, or **Not Applicable**. A Pass needs current evidence for the exact product, revision, artifact, configuration, environment, and data state under review.

## User-Generated Content, Social Features, and Marketplaces

_Consolidated from `quality standards/15-conditional-domains/02-user-generated-content-social-and-marketplaces.md`; 17 non-duplicative controls._

### Universal controls

- [ ] **USEQ-560779CA** — Acceptable-use, content, conduct, and enforcement policies are defined.
- [ ] **USEQ-A7D75987** — Users can report prohibited content, harmful conduct, impersonation, fraud, and safety threats.
- [ ] **USEQ-94462871** — Blocking, muting, audience, privacy, and safety controls work reliably.
- [ ] **USEQ-6F4DE347** — Moderation queues have owners, severity rules, response objectives, and escalation paths.
- [ ] **USEQ-9EFA79C0** — Moderator actions are authorized, attributable, reviewable, and audited.
- [ ] **USEQ-84EE02EE** — Appeals and correction procedures exist where appropriate.
- [ ] **USEQ-925A28B6** — Spam, bots, scraping, harassment, stalking, impersonation, brigading, and coordinated abuse are addressed.
- [ ] **USEQ-F0DC21E4** — Uploaded and linked content is scanned and isolated appropriately.
- [ ] **USEQ-6BE6FEE6** — Illegal-content, copyright, evidence-preservation, lawful-request, and emergency-disclosure procedures are defined.
- [ ] **USEQ-AF8BAAB3** — Public, private, group, follower, direct-message, and restricted audience controls are clear and enforced server-side.
- [ ] **USEQ-E0664353** — Deleted or private content does not remain publicly cached, searchable, recommended, or indexed.
- [ ] **USEQ-5AF9805B** — Search, ranking, notifications, recommendations, and previews do not expose restricted content.
- [ ] **USEQ-0A61FF1D** — Seller, buyer, provider, creator, or organizer verification is applied where risk requires it.
- [ ] **USEQ-320189B8** — Fraud, dispute, refund, reputation, review, and trust-and-safety workflows are integrated.
- [ ] **USEQ-E7BF7973** — Moderation personnel have suitable wellbeing safeguards and restricted data access.
- [ ] **USEQ-BC1EE838** — Threats to life or safety have an emergency escalation process.
- [ ] **USEQ-DABCEC8A** — Content-retention and account-deletion behavior is consistent with legal, safety, and privacy requirements.

## Email, SMS, Push, and Notifications

_Consolidated from `quality standards/15-conditional-domains/03-email-sms-push-and-notifications.md`; 14 non-duplicative controls._

### Universal controls

- [ ] **USEQ-C773D0E2** — Sender domains, identities, certificates, and provider accounts are controlled.
- [ ] **USEQ-BE8B8B0B** — Email authentication uses current domain-authentication practices.
- [ ] **USEQ-9B4FADDD** — Bounce, complaint, rejection, suppression, reputation, and delivery-delay signals are monitored.
- [ ] **USEQ-C259ED15** — Unsubscribe and preference changes work promptly and are not deceptive.
- [ ] **USEQ-E68ADD70** — Transactional, security, operational, and marketing messages are classified correctly.
- [ ] **USEQ-9041D400** — Sensitive information is not exposed unnecessarily in lock-screen content, email previews, SMS, subject lines, or push payloads.
- [ ] **USEQ-AF1D0C8C** — Reset, verification, invitation, and approval messages cannot be reused or redirected improperly.
- [ ] **USEQ-30B95EF0** — Retries do not create duplicates, message storms, harassment, or spam.
- [ ] **USEQ-D7B8731A** — Templates handle localization, escaping, long content, rendering differences, dark mode, plain text, and accessible structure.
- [ ] **USEQ-AFA6B800** — Provider outage, delay, duplicate delivery, and out-of-order delivery are handled.
- [ ] **USEQ-BF028AE2** — Frequency and quiet-hour controls prevent harassment and abuse where appropriate.
- [ ] **USEQ-9B51BA8F** — Phone-number and email-address reassignment risk is considered.
- [ ] **USEQ-ADF8072F** — Production cannot accidentally use test recipients, and testing cannot accidentally use production recipients.
- [ ] **USEQ-F56AD22C** — Security notifications are resistant to spoofing and tell users how to verify legitimacy safely.

## Search, Ranking, Recommendations, and Personalization

_Consolidated from `quality standards/15-conditional-domains/09-search-ranking-recommendations-and-personalization.md`; 13 non-duplicative controls._

### Universal controls

- [ ] **USEQ-07F09D35** — Search and recommendation inputs, indexes, models, and ranking rules are inventoried.
- [ ] **USEQ-6D5C3152** — Access control is enforced before indexing, retrieval, ranking, preview, and recommendation.
- [ ] **USEQ-65B2307E** — Deleted, private, blocked, restricted, or expired content is removed or suppressed promptly.
- [ ] **USEQ-64BF1EA1** — Query expansion, autocomplete, spelling correction, snippets, and facets do not leak restricted information.
- [ ] **USEQ-34F34DA1** — Ranking cannot be manipulated trivially through spam, poisoning, coordinated behavior, or hidden metadata.
- [ ] **USEQ-B3D62C45** — Personalization uses only authorized and disclosed data.
- [ ] **USEQ-0ED5176E** — Users can understand or control personalization where required.
- [ ] **USEQ-897DFBBC** — Sensitive traits are not inferred or used without proper review and authorization.
- [ ] **USEQ-A51E620A** — Fairness, discrimination, safety, diversity, and filter-bubble risks are assessed where material.
- [ ] **USEQ-387B93A0** — Search and recommendation quality is evaluated across languages, user groups, cold starts, sparse data, and adversarial inputs.
- [ ] **USEQ-EF56A225** — Caching preserves user, tenant, audience, and permission boundaries.
- [ ] **USEQ-BD559A2D** — Index lag, model drift, and ranking changes are monitored.
- [ ] **USEQ-33AC50F2** — A kill switch or rollback exists for harmful ranking or recommendation changes.

## Trust, Safety, Content Integrity, and Extension Ecosystems Master Checklist

_Consolidated from `gap supplement/04-trust-safety-content-integrity-and-extension-ecosystems.md`; 191 non-duplicative controls._

### Expanded gap-closure controls

#### Safety-by-design governance

- [ ] **USEQ-BD9BA43E** — Apply this checklist whenever users can communicate, publish, transact, discover other users, distribute content, install extensions, delegate capabilities, or materially affect one another.
- [ ] **USEQ-AEDA607A** — Make user safety a product and governance responsibility rather than assigning it solely to moderation or support after harm occurs.
- [ ] **USEQ-AFD36B1D** — Assign executive, product, engineering, policy, operations, legal, privacy, security, accessibility, and regional ownership for material safety risks.
- [ ] **USEQ-596F4B00** — Define the rights, dignity, autonomy, privacy, expression, access, and safety outcomes the product is intended to preserve.
- [ ] **USEQ-9121E475** — Create a harm taxonomy covering physical, psychological, financial, reputational, privacy, discrimination, exploitation, harassment, fraud, coercion, self-harm, child safety, and societal risks as applicable.
- [ ] **USEQ-71D94266** — Identify higher-risk users and contexts, including children, survivors, marginalized groups, public figures, workers, sellers, creators, and users under coercion.
- [ ] **USEQ-E211C089** — Model how product incentives, growth loops, ranking, virality, defaults, anonymity, payments, and friction can increase or reduce harm.
- [ ] **USEQ-B12D033E** — Include misuse, abuse, coordinated behavior, insider action, compromised accounts, synthetic identities, automation, and cross-platform migration in threat and safety models.
- [ ] **USEQ-76490994** — Define risk appetite, release blockers, escalation thresholds, emergency actions, and authority to slow or stop growth.
- [ ] **USEQ-9FB97874** — Fund safety, moderation, investigation, appeals, law-enforcement response where lawful, and victim support proportionately to product reach and risk.
- [ ] **USEQ-4BE942F0** — Review safety impact before launches, market expansion, new messaging capabilities, recommendation changes, monetization changes, privacy changes, and reduced friction.
- [ ] **USEQ-C624557C** — Use staged rollout and heightened monitoring for features that change reach, contactability, discoverability, monetization, or user power.
- [ ] **USEQ-95CED227** — Document safety assumptions and test them with affected users, civil-society expertise, domain specialists, and adversarial research where impact warrants it.
- [ ] **USEQ-D5D07405** — Publish internal safety objectives and hold leaders accountable for outcome trends, not only content-removal volume.

#### Policies, notice, and user expectations

- [ ] **USEQ-95E87C5D** — Publish clear, accessible, localized rules for permitted and prohibited content, conduct, commerce, automation, impersonation, and extension behavior.
- [ ] **USEQ-2A90DDC8** — Make policies specific enough that users, moderators, developers, and reviewers can apply them consistently.
- [ ] **USEQ-82358F79** — Explain the scope of private, public, encrypted, ephemeral, moderated, and externally visible spaces without making false confidentiality promises.
- [ ] **USEQ-9ED88D70** — Explain what automated systems, human reviewers, trusted flaggers, merchants, developers, or partners may review and why.
- [ ] **USEQ-F6F7B3BA** — Provide age, geographic, and product-specific rules where risks and obligations differ.
- [ ] **USEQ-698F9E31** — Notify users of material policy changes before enforcement where feasible and preserve prior versions.
- [ ] **USEQ-59123933** — Separate safety rules from promotional language and avoid hiding consequential terms in inaccessible or excessive text.
- [ ] **USEQ-07E9A6B6** — Give developers and extension publishers explicit technical, privacy, permission, review, update, and enforcement requirements.
- [ ] **USEQ-1EEF3A16** — State consequences, enforcement ranges, appeal paths, evidence retention, and restoration possibilities.
- [ ] **USEQ-69D30DCF** — Avoid retroactive punishment for conduct that was permitted unless urgent safety or legal necessity justifies it and the rationale is documented.
- [ ] **USEQ-6CE2FA9B** — Provide policy examples without creating an exhaustive evasion manual.
- [ ] **USEQ-5DD0664C** — Ensure policy and product behavior agree; do not promise controls, review, confidentiality, or response times that the system cannot provide.

#### User empowerment, consent, and protective controls

- [ ] **USEQ-AD777CE0** — Provide understandable privacy, audience, contact, discovery, tagging, messaging, location, recommendation, and data-sharing controls.
- [ ] **USEQ-A780B514** — Use protective defaults proportionate to age, vulnerability, sensitivity, and expected context.
- [ ] **USEQ-22F0FF7A** — Allow users to block, mute, restrict, unfollow, leave, remove followers, limit replies, filter keywords, and control contact where relevant.
- [ ] **USEQ-9D270EFD** — Make safety controls available before, during, and after an interaction rather than only after harm.
- [ ] **USEQ-A7EF66C3** — Ensure blocking applies consistently across profiles, messages, search, recommendations, notifications, groups, shared content, and alternate identifiers as intended.
- [ ] **USEQ-486EFD38** — Prevent a blocked or restricted actor from learning unnecessary information through errors, presence, read receipts, recommendations, or indirect notifications.
- [ ] **USEQ-B65E9725** — Allow users to review and revoke connected applications, sessions, devices, delegates, payment mandates, and extension permissions.
- [ ] **USEQ-2417A7C7** — Request dangerous permissions at the time of meaningful use and explain the purpose in user language.
- [ ] **USEQ-D6445268** — Make optional permissions and data collection revocable without forcing account deletion when the core service can continue.
- [ ] **USEQ-77B28528** — Provide friction, previews, warnings, cooldowns, confirmation, and undo for actions with high abuse or irreversible impact.
- [ ] **USEQ-29E4B82E** — Allow users to preserve evidence safely before blocking, deleting, or leaving where needed.
- [ ] **USEQ-7C0C7DE6** — Provide discreet safety exits and account-protection options for users under surveillance or coercion.
- [ ] **USEQ-ABDF7D26** — Do not expose safety choices publicly when doing so could increase retaliation or stigma.
- [ ] **USEQ-0CB908CD** — Provide accessible controls and alternatives for users with disabilities, limited literacy, low bandwidth, or constrained devices.

#### Reporting, moderation, and case handling

- [ ] **USEQ-6092D596** — Provide accessible reporting from the content, account, transaction, message, extension, or behavior being reported.
- [ ] **USEQ-41AFAEC3** — Allow reporting without requiring the reporter to maintain contact with the reported party.
- [ ] **USEQ-2EB40D3B** — Offer report categories that support triage while allowing free-text context and urgent-safety signals.
- [ ] **USEQ-BD179B81** — Acknowledge reports and provide a case reference without exposing protected investigation details.
- [ ] **USEQ-8A7710F9** — Prioritize cases using severity, imminence, vulnerability, scale, recurrence, virality, illegality, and evidence-loss risk.
- [ ] **USEQ-FF6C3B4B** — Define service objectives for acknowledgment, triage, emergency escalation, decision, appeal, and restoration by severity.
- [ ] **USEQ-9BE743D4** — Route credible imminent threats, child-safety risk, self-harm emergencies, financial theft, account compromise, and coordinated abuse to specialized procedures.
- [ ] **USEQ-36EBFCB9** — Preserve relevant content, metadata, decisions, and chain of custody according to law, policy, privacy, and retention requirements.
- [ ] **USEQ-5D0E8382** — Protect reporters, moderators, witnesses, and victims from retaliation, doxxing, and unnecessary disclosure.
- [ ] **USEQ-0750C577** — Use trained reviewers and provide specialized regional, linguistic, cultural, accessibility, legal, and subject-matter support.
- [ ] **USEQ-FC571D84** — Provide moderators with sufficient context while minimizing unnecessary exposure to personal or traumatic material.
- [ ] **USEQ-8D51224C** — Use rotation, wellness protections, clinical support, access controls, and workload limits for personnel exposed to harmful content.
- [ ] **USEQ-4A26124E** — Require consistent decision records including policy basis, evidence, confidence, action, duration, scope, reviewer, and escalation.
- [ ] **USEQ-1E2471A3** — Use quality assurance, calibration, double review, sampling, and disagreement analysis for consequential decisions.
- [ ] **USEQ-E9E9307A** — Prevent automated tools from making irreversible high-impact enforcement decisions without the required confidence, safeguards, and human review.
- [ ] **USEQ-C08DB232** — Provide proportionate actions such as warning, friction, reduced distribution, feature restriction, temporary suspension, removal, preservation, and permanent exclusion.
- [ ] **USEQ-623ABC6E** — Apply actions across linked accounts or devices only under a documented, privacy-reviewed, error-tolerant policy.
- [ ] **USEQ-C934C5F6** — Notify affected users with the rule, content or behavior, action, duration, consequences, and appeal route unless doing so would create material risk.
- [ ] **USEQ-75B5AC01** — Provide accessible, independent, timely appeals with review by someone not solely repeating the original automated decision.
- [ ] **USEQ-FC548DD3** — Restore content, access, reputation, reach, funds, and associated state when an enforcement error is reversed.
- [ ] **USEQ-C11B8851** — Track false positives, false negatives, inconsistent decisions, appeal reversals, recurrence, and unresolved backlog.

#### Recommendation, amplification, and virality controls

- [ ] **USEQ-613FA405** — Treat ranking, recommendation, trending, autocomplete, notifications, resharing, and search as distribution systems with safety responsibilities.
- [ ] **USEQ-8CB830F6** — Define which safety, quality, provenance, age, sensitivity, and policy signals affect eligibility, ranking, and reach.
- [ ] **USEQ-AE03B731** — Prevent prohibited or restricted content from re-entering recommendation through copies, edits, alternate encodings, or related accounts.
- [ ] **USEQ-34281A6E** — Test whether optimization objectives create incentives for outrage, harassment, misinformation, addiction, exploitation, fraud, or unsafe challenges.
- [ ] **USEQ-1077D020** — Use counter-metrics for harmful exposure, regret, complaints, blocks, hides, rapid exits, appeals, and vulnerable-user impact.
- [ ] **USEQ-DA6EBB44** — Limit rapid mass distribution for new, untrusted, sensitive, or unusually fast-growing content until risk signals are evaluated where appropriate.
- [ ] **USEQ-F6C33ECE** — Provide users meaningful control over recommendation topics, sources, personalization, sensitive content, and chronological alternatives where relevant.
- [ ] **USEQ-435FB22F** — Do not infer or exploit highly sensitive traits for distribution without a lawful, necessary, and user-respecting basis.
- [ ] **USEQ-F981278F** — Evaluate disparate exposure, suppression, monetization, and error rates across relevant groups and languages.
- [ ] **USEQ-4676585B** — Ensure downranking and reduced distribution have policy, governance, testing, monitoring, and appeal appropriate to their impact.
- [ ] **USEQ-075A14AD** — Prevent coordinated manipulation through engagement rings, fake accounts, purchased interactions, adversarial metadata, and automated activity.
- [ ] **USEQ-06F3E324** — Review recommendation changes through offline evaluation, adversarial testing, limited rollout, causal analysis, and post-launch monitoring.
- [ ] **USEQ-C776200D** — Provide researchers and auditors appropriate transparency and privacy-preserving access when required or beneficial to safety assurance.

#### Identity, impersonation, authenticity, and account integrity

- [ ] **USEQ-0C858E3A** — Define when real identity, pseudonymity, anonymity, organization verification, age assurance, or role verification is necessary and proportionate.
- [ ] **USEQ-7F8BF05A** — Avoid collecting stronger identity evidence than the risk requires.
- [ ] **USEQ-2AEAE9F2** — Distinguish verified control of an account or organization from endorsement, expertise, factual accuracy, or general trustworthiness.
- [ ] **USEQ-2BBFCDF2** — Prevent deceptive impersonation of people, organizations, support agents, merchants, public authorities, and automated systems.
- [ ] **USEQ-33DC094C** — Provide clear disclosure for bots, synthetic personas, delegated operators, and materially automated accounts where users could be misled.
- [ ] **USEQ-F6D77BAF** — Protect usernames, organization names, custom domains, badges, and profile identifiers from confusing or abusive collision.
- [ ] **USEQ-F009B337** — Detect account takeover, credential stuffing, session theft, recovery abuse, SIM or email reassignment, and support-channel social engineering.
- [ ] **USEQ-D805F610** — Use step-up verification and cooling periods for changes to identity, recovery, payout, advertising, administrative, or high-reach settings.
- [ ] **USEQ-BC52C708** — Notify users of material identity and security changes through independent channels where appropriate.
- [ ] **USEQ-31D2242A** — Provide recovery paths for compromised accounts without giving attackers an easier route than normal authentication.
- [ ] **USEQ-C1AD085C** — Preserve evidence and limit further harm during account-takeover investigation.
- [ ] **USEQ-361A211A** — Allow users to report impersonation even when they do not have an account.
- [ ] **USEQ-BD3A717B** — Define memorialization, succession, organization transfer, employee departure, and abandoned-account behavior.

#### Fraud, spam, bots, and coordinated abuse

- [ ] **USEQ-8DE587F8** — Model economic incentives, attacker return, victim loss, platform subsidy, and externalized cost for each abuse class.
- [ ] **USEQ-41719BE8** — Use layered controls across enrollment, authentication, device, network, behavior, content, transaction, reputation, payout, and recovery.
- [ ] **USEQ-F2E3DC32** — Bound actions per account, identity, tenant, device, network, payment instrument, destination, and time window using risk-aware limits.
- [ ] **USEQ-3E29C386** — Detect burst, fan-out, graph, velocity, replay, template, shared-infrastructure, and collusion patterns.
- [ ] **USEQ-5DEEFFE1** — Prevent rate-limit evasion through alternate endpoints, encodings, identities, tenants, or protocol layers.
- [ ] **USEQ-EDC9A641** — Use challenges and friction that are accessible, privacy-preserving, and proportionate rather than relying on one CAPTCHA.
- [ ] **USEQ-19CE9D99** — Separate suspicious activity from proven abuse and avoid irreversible action based on one weak signal.
- [ ] **USEQ-0D9379C8** — Keep fraud models and rules observable, versioned, tested, monitored for drift, and protected from unauthorized manipulation.
- [ ] **USEQ-0CC5A940** — Protect high-risk state changes, payouts, refunds, transfers, promotions, invitations, reviews, and messaging from automation abuse.
- [ ] **USEQ-15AABFBA** — Delay or reserve irreversible value transfer until authoritative risk and settlement signals are available where appropriate.
- [ ] **USEQ-649D5697** — Reconcile platform records with payment, fulfillment, identity, and provider records to detect hidden losses.
- [ ] **USEQ-B1FBD274** — Share abuse intelligence internally and with trusted partners under privacy, legal, accuracy, and purpose controls.
- [ ] **USEQ-6D72C114** — Design containment so blocking one attack does not create widespread denial of service for legitimate users.
- [ ] **USEQ-1BA3C45A** — Measure attacker adaptation and rotate controls without relying on secrecy as the primary defense.

#### Children, vulnerable users, and high-risk interactions

- [ ] **USEQ-53975D3E** — Determine whether children are likely users even when the service is not explicitly directed to them.
- [ ] **USEQ-0BBF963F** — Use age-appropriate defaults for discoverability, contact, location, profiling, advertising, spending, sharing, and public posting.
- [ ] **USEQ-95346626** — Avoid dark patterns that pressure minors or vulnerable users to disclose data, spend money, maintain streaks, or engage continuously.
- [ ] **USEQ-BC356E9B** — Limit adult-to-minor contact, grooming patterns, coercive requests, sexualization, and migration to less protected channels where applicable.
- [ ] **USEQ-FB8D5668** — Provide guardian controls only when lawful and designed not to endanger children in abusive households.
- [ ] **USEQ-3FC789B0** — Use age assurance proportionate to risk and evaluate privacy, accessibility, bias, circumvention, and exclusion.
- [ ] **USEQ-4194C938** — Provide child-understandable explanations, reporting, blocking, and help.
- [ ] **USEQ-FCB00C8A** — Route child-safety incidents to trained teams with strict access, evidence, reporting, and wellness controls.
- [ ] **USEQ-B0F21CF7** — Protect users seeking help for self-harm, abuse, addiction, financial distress, or violence from exploitation and harmful recommendation loops.
- [ ] **USEQ-7B563568** — Test safety controls with representative youth and vulnerable-user expertise.
- [ ] **USEQ-315C17BD** — Avoid public indicators that reveal a user's protected status or safety settings.

#### Crisis, emergency, and severe-harm response

- [ ] **USEQ-197A27D2** — Define criteria and authority for emergency feature restriction, account freeze, content preservation, distribution reduction, payment hold, or service shutdown.
- [ ] **USEQ-DC481166** — Maintain current contacts and procedures for emergency services, child-protection channels, payment providers, hosting providers, legal counsel, and trusted experts where applicable.
- [ ] **USEQ-7545CE12** — Verify jurisdiction, legal basis, identity, scope, urgency, and preservation requirements for emergency requests.
- [ ] **USEQ-28ED6E96** — Minimize disclosure and record every emergency access or data release.
- [ ] **USEQ-2C2015DC** — Provide an out-of-band response path when the normal support, identity, or platform systems are compromised.
- [ ] **USEQ-251F513E** — Exercise scenarios involving credible threats, coordinated violence, mass fraud, viral harmful content, child exploitation, doxxing, extortion, and insider abuse as relevant.
- [ ] **USEQ-FEB8AB8C** — Coordinate public communication without amplifying harmful material, exposing victims, or prejudicing investigation.
- [ ] **USEQ-324E1782** — Define when and how affected users are notified after emergency actions.
- [ ] **USEQ-B38FDE82** — Review emergency actions promptly and reverse overbroad restrictions when safe.
- [ ] **USEQ-0AE9EFF0** — Conduct post-incident analysis that addresses product incentives, detection gaps, staffing, process, tooling, and governance.

#### Content provenance, synthetic media, and authenticity signals

- [ ] **USEQ-C6361B6A** — Define which content types and workflows require provenance, origin, edit history, capture context, signer identity, or synthetic-content disclosure.
- [ ] **USEQ-A3D6E17A** — Use cryptographically bound provenance standards such as C2PA where they meaningfully improve authenticity and ecosystem interoperability.
- [ ] **USEQ-12397117** — Treat provenance as evidence about claims and history, not proof that content is true, safe, unbiased, or authorized.
- [ ] **USEQ-87AC1895** — Verify signatures, certificate chains, trust lists, timestamps, asset binding, manifest integrity, revocation, and supported algorithms.
- [ ] **USEQ-E9A1411C** — Preserve provenance through transcoding, resizing, editing, composition, export, syndication, and archival where technically and legally feasible.
- [ ] **USEQ-279B0168** — Make provenance loss, invalidity, redaction, and unsupported claims visible rather than presenting a false trusted state.
- [ ] **USEQ-466DE378** — Protect signing keys, capture credentials, author identities, and provenance services as high-value security assets.
- [ ] **USEQ-B70CAE6C** — Allow privacy-preserving redaction and selective disclosure without silently rewriting history.
- [ ] **USEQ-E8CBABAC** — Provide accessible human-readable explanations of provenance and avoid binary trusted/untrusted labels that overstate confidence.
- [ ] **USEQ-A48CB2BE** — Record automated generation, material editing, model or tool identity, source assets, and responsible publisher when required by policy or risk.
- [ ] **USEQ-175ACA22** — Use visible labels, watermarks, metadata, detection, and user education as complementary controls; do not rely on one mechanism as perfect detection.
- [ ] **USEQ-59C5F9CA** — Evaluate provenance spoofing, stripping, laundering, replay, key compromise, trust-list abuse, and misleading authentic-but-harmful content.
- [ ] **USEQ-9C34ECD6** — Monitor ecosystem support and maintain graceful behavior when receiving unknown or newer provenance formats.

#### Extension, plugin, app, and integration ecosystems

- [ ] **USEQ-BF82CD25** — Apply a formal intake, review, publication, update, monitoring, enforcement, and retirement lifecycle to third-party extensions and integrations.
- [ ] **USEQ-0B382DF6** — Require an accurate description of functionality, data access, remote services, permissions, monetization, ownership, support, and security practices.
- [ ] **USEQ-ABE691EB** — Require only the minimum permissions and allow optional permissions to be requested at meaningful use time where supported.
- [ ] **USEQ-2900DB5A** — Make permission prompts specific, understandable, attributable, and revisitable.
- [ ] **USEQ-E6B9618B** — Prevent extensions from loading unreviewed remote executable code or silently changing behavior outside the approved update channel.
- [ ] **USEQ-9EF29D9E** — Require reviewable source or reproducible correspondence between submitted source and distributed artifact for high-risk ecosystems.
- [ ] **USEQ-EC3D9D3D** — Scan dependencies, secrets, malware, obfuscation, suspicious behavior, policy violations, and known vulnerabilities.
- [ ] **USEQ-5CC601B2** — Restrict dynamic evaluation, native execution, privileged APIs, cross-origin access, filesystem access, identity access, and payment access according to risk.
- [ ] **USEQ-D2594826** — Sandbox extension execution and isolate storage, network, UI, processes, credentials, and tenant context where practical.
- [ ] **USEQ-9BBB4EEA** — Prevent an extension from weakening host security headers, transport security, authentication, authorization, privacy, accessibility, or audit controls.
- [ ] **USEQ-632E8E55** — Require encrypted transport and authenticated endpoints for extension data exchange.
- [ ] **USEQ-5EDDE5B2** — Require transparent data collection, purpose, retention, sharing, sale, consent, and deletion behavior.
- [ ] **USEQ-48EB755B** — Test extension performance, memory, battery, startup, stability, offline behavior, accessibility, localization, and failure containment.
- [ ] **USEQ-3452812F** — Prevent one faulty or malicious extension from crashing, blocking, corrupting, or exfiltrating the host or other extensions.
- [ ] **USEQ-7A7CA22D** — Use signed packages and a secure update framework resistant to rollback, freeze, mix-and-match, repository compromise, and signing-key compromise.
- [ ] **USEQ-245DA35C** — Support staged rollout, compatibility checks, rollback or disablement, revocation, and emergency blocklisting.
- [ ] **USEQ-96481324** — Make extension identity, publisher changes, ownership transfers, and signing-key changes reviewable and visible.
- [ ] **USEQ-9DB20BBB** — Require updates that add permissions, data use, remote services, or materially different behavior to receive renewed review and user notice.
- [ ] **USEQ-724B59B6** — Monitor abuse reports, crash rates, permission use, network destinations, update anomalies, dormant takeover, and publisher compromise.
- [ ] **USEQ-55B8A0A8** — Provide users a clear inventory, enable/disable control, permission review, update history, and removal path.
- [ ] **USEQ-374B203B** — Remove credentials, data access, webhooks, tokens, stored data, background jobs, and delegated rights when an integration is removed.
- [ ] **USEQ-575DA746** — Define compatibility and deprecation policy for host APIs without forcing insecure legacy support indefinitely.
- [ ] **USEQ-130FF33D** — Provide an appeal process for developers while prioritizing user safety during credible compromise.

#### Ecosystem and marketplace governance

- [ ] **USEQ-A9C693F2** — Verify merchants, developers, advertisers, creators, and high-reach accounts proportionately to their ability to cause harm.
- [ ] **USEQ-3F7F7E38** — Separate platform endorsement, verification, popularity, ranking, and paid placement visually and semantically.
- [ ] **USEQ-5242F580** — Disclose commissions, fees, sponsorship, ranking influence, and platform conflicts that materially affect user choice.
- [ ] **USEQ-9F589E1F** — Prevent review manipulation, fake testimonials, retaliatory reviews, review gating, and undisclosed incentives.
- [ ] **USEQ-7750BA91** — Protect sellers, creators, and developers from arbitrary enforcement through published rules, notice, evidence, appeal, and restoration.
- [ ] **USEQ-1AA1AED9** — Prevent marketplace design from pushing risk, support, refunds, tax, safety, or legal obligations onto users without clarity.
- [ ] **USEQ-941BFF73** — Define responsibility for product safety, prohibited goods, intellectual property, counterfeit claims, fulfillment, refunds, and disputes.
- [ ] **USEQ-6CD3ED9A** — Monitor dormant account or package takeover, ownership transfers, sudden behavior changes, and supply-chain compromise.
- [ ] **USEQ-DF8621BD** — Apply progressive trust and limits rather than granting full reach, payouts, or permissions immediately to unproven participants.
- [ ] **USEQ-A847268E** — Provide transparent complaint and dispute resolution with evidence and deadlines.
- [ ] **USEQ-7ADC4406** — Measure ecosystem health using user harm, fraud loss, quality, diversity, concentration, developer sustainability, false enforcement, and support—not gross volume alone.

#### Transparency, audit, and continuous improvement

- [ ] **USEQ-0164140D** — Publish meaningful aggregate information about policies, reports, enforcement, appeals, automation, severe harms, and response performance where appropriate.
- [ ] **USEQ-8C2EE771** — Explain material recommendation, verification, provenance, and enforcement systems at a level users and oversight bodies can understand.
- [ ] **USEQ-99F86F23** — Disclose the limits and known error modes of automated moderation, age assurance, identity, fraud, and provenance systems.
- [ ] **USEQ-DF46D365** — Protect transparency reports against privacy leakage, reidentification, security evasion, and misleading aggregation.
- [ ] **USEQ-EB1EF776** — Provide independent audit, researcher access, regulator access, or trusted oversight where risk or obligation warrants it.
- [ ] **USEQ-AA903343** — Retain model, rule, policy, reviewer, evidence, and version information needed to reconstruct consequential decisions.
- [ ] **USEQ-2140FACB** — Track harm prevalence and exposure, not only the number of items removed.
- [ ] **USEQ-FB9F9E37** — Measure time to detection, interruption, victim support, recovery, appeal, and recurrence prevention.
- [ ] **USEQ-CCCAA8A8** — Review false-positive and false-negative impacts across languages, regions, disabilities, ages, and affected groups.
- [ ] **USEQ-DF0CB7C2** — Use red-team, tabletop, abuse simulation, user research, and incident learning to discover controls that ordinary testing misses.
- [ ] **USEQ-E3F9DBC4** — Convert recurring abuse into safer defaults, reduced capability, stronger boundaries, better education, and redesigned incentives.
- [ ] **USEQ-4C29738F** — Review this checklist after major incidents, market expansion, policy changes, new content formats, new extension capabilities, or attacker adaptation.

#### Trust-and-safety release blockers

- [ ] **USEQ-9B250F50** — Do not launch a high-risk interaction surface without reporting, blocking, escalation, evidence, enforcement, and appeal paths.
- [ ] **USEQ-69245CE0** — Do not launch when severe-harm reports cannot reach a qualified responder within the required objective.
- [ ] **USEQ-97EA05F5** — Do not launch a recommendation or virality change when harmful-exposure and abuse effects are unmeasured and uncontrolled.
- [ ] **USEQ-FE34F2BA** — Do not grant broad extension permissions, remote code execution, or privileged host access without review, containment, transparency, and revocation.
- [ ] **USEQ-3FF9D040** — Do not distribute updates whose source, signer, version, and authorized update path cannot be verified.
- [ ] **USEQ-0D13B7DC** — Do not present provenance, verification, or identity badges in a way that falsely implies truth, endorsement, expertise, or safety.
- [ ] **USEQ-307DFDE0** — Do not collect identity or age evidence beyond demonstrated necessity and protection capability.
- [ ] **USEQ-04AA676C** — Do not make irreversible enforcement or financial decisions from one opaque automated signal without required review and recovery.
- [ ] **USEQ-01689474** — Do not suppress known severe-harm trends to protect growth, engagement, reputation, or revenue metrics.
- [ ] **USEQ-9F34570B** — Do not accept a safety exception without an explicit harm analysis, protective alternative, monitoring, owner, remediation date, and expiry.

## Standards and source references

- [ISO/IEC 25010:2023 — Product quality model](https://www.iso.org/standard/78176.html)
- [NIST Cybersecurity Framework 2.0](https://www.nist.gov/cyberframework)
- [ISO/IEC 27701:2025 — Privacy information management systems](https://www.iso.org/standard/85819.html)
- [ISO/IEC 42001:2023 — AI management systems](https://www.iso.org/standard/81230.html)
- [ISO/IEC 23894:2023 — AI risk management](https://www.iso.org/standard/77304.html)
- [eSafety Safety by Design](https://www.esafety.gov.au/industry/safety-by-design)
- [C2PA Specifications 2.4](https://spec.c2pa.org/specifications/specifications/2.4/index.html)
- [C2PA Security Considerations](https://spec.c2pa.org/specifications/specifications/2.4/security/Security_Considerations.html)
- [Mozilla Add-on Policies](https://extensionworkshop.com/documentation/publish/add-on-policies/)
- [Mozilla permission guidance](https://extensionworkshop.com/documentation/develop/request-the-right-permissions/)
- [The Update Framework](https://theupdateframework.io/)
- [NIST AI 100-4 — Reducing Risks Posed by Synthetic Content](https://nvlpubs.nist.gov/nistpubs/ai/NIST.AI.100-4.pdf)

---

[Previous phase](13-documentation-and-knowledge.md) · [Next: Phase 15: AI, ML, and AI-assisted development](15-ai-ml-and-ai-assisted-development.md)
