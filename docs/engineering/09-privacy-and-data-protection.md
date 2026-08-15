# Privacy and data protection

_Phase 9 of 16 in the [complete engineering review](00-overview.md)._

Privacy governance, notices, consent, rights, minimization, retention, sharing, de-identification, and sensitive data.

Assess every control as **Pass**, **Fail**, **Blocked**, or **Not Applicable**. A Pass needs current evidence for the exact product, revision, artifact, configuration, environment, and data state under review.

## Privacy and Data Protection

_Consolidated from `quality standards/10-privacy/01-privacy-and-data-protection.md`; 16 non-duplicative controls._

### Universal controls

- [ ] **USEQ-362ADF8D** — Product behavior matches published privacy notices and user choices.
- [ ] **USEQ-375E775E** — New data uses undergo privacy review.
- [ ] **USEQ-46C19961** — Consent is used only where appropriate and is specific, informed, granular, recorded, and withdrawable.
- [ ] **USEQ-26C38E0D** — Cookie, tracker, analytics, advertising, attribution, session-replay, and experimentation inventories match production behavior.
- [ ] **USEQ-68E3287E** — Third-party tools receive only approved data.
- [ ] **USEQ-85394A39** — Sensitive fields are redacted from logs, monitoring, support tools, recordings, analytics, and error reports.
- [ ] **USEQ-1917583D** — Requester verification is proportionate and does not collect excessive additional data.
- [ ] **USEQ-06C60675** — Retention periods are explicit and enforced.
- [ ] **USEQ-51B2160E** — Deletion covers primary data, caches, indexes, derived data, exports, and downstream systems.
- [ ] **USEQ-5CCD452D** — Pseudonymization and anonymization claims are validated.
- [ ] **USEQ-D295EBB5** — Reidentification risk is assessed.
- [ ] **USEQ-19407D75** — Data export cannot expose other users' or tenants' information.
- [ ] **USEQ-EDE4FE5C** — Children's, biometric, health, precise-location, communications, financial, identity, and other sensitive data trigger enhanced review.
- [ ] **USEQ-FE438B01** — Processor and subprocessor contracts cover permitted use, security, deletion, assistance, audit, and incident notification.
- [ ] **USEQ-0CA8AE68** — A personal-data incident can be identified, scoped, contained, and reported within applicable deadlines.
- [ ] **USEQ-009AD2AD** — Interfaces avoid manipulative consent, cancellation, and privacy choices.

## Privacy Governance and Engineering

_Consolidated from `quality standards/10-privacy/02-privacy-governance-and-engineering.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-E171ADC5** — Assign accountable privacy owners and data stewards for every product and processing purpose.
- [ ] **USEQ-2B4BFFCE** — Maintain an inventory of personal data, sources, purposes, lawful basis or authorization, recipients, locations, retention, and deletion.
- [ ] **USEQ-178B5045** — Map personal-data flows across clients, services, logs, analytics, support, suppliers, backups, models, and exports.
- [ ] **USEQ-BDE2681A** — Apply privacy-protective defaults and collect only what is necessary for the stated purpose.
- [ ] **USEQ-44C28E2D** — Separate mandatory processing from optional processing and communicate the difference.
- [ ] **USEQ-EBEB9674** — Assess privacy risk before new collection, inference, combination, sharing, monitoring, or automated decision use.
- [ ] **USEQ-C82F78B9** — Translate privacy principles into technical requirements, data models, access, retention, UI, telemetry, and tests.
- [ ] **USEQ-79A3CEBF** — Use purpose limitation and access controls to prevent unauthorized secondary use.
- [ ] **USEQ-2CDD7FC7** — Minimize identifiers and linkability across contexts.
- [ ] **USEQ-D84D7490** — Use pseudonymization, aggregation, isolation, or on-device processing where they meaningfully reduce risk.
- [ ] **USEQ-2B0B18B4** — Include privacy threats such as reidentification, inference, surveillance, coercion, exclusion, and loss of autonomy.
- [ ] **USEQ-C8E8FC88** — Review suppliers for permitted use, location, retention, deletion, security, subprocessing, and incident duties.
- [ ] **USEQ-C5D655EE** — Keep privacy records and notices consistent with deployed behavior.
- [ ] **USEQ-A637D569** — Monitor data flows and trackers for unauthorized change.
- [ ] **USEQ-8B9E2BD8** — Provide privacy review for high-risk processing and material product changes.
- [ ] **USEQ-0978812C** — Treat privacy incidents and user complaints as design feedback.
- [ ] **USEQ-E3F0E7D9** — Measure rights completion, deletion effectiveness, access minimization, unexpected collection, and risk reduction.
- [ ] **USEQ-FF1E023B** — Train role-appropriate personnel before granting personal-data access.
- [ ] **USEQ-58593369** — Define privacy behavior through deprecation, migration, backup, and retirement.
- [ ] **USEQ-9DE0E53C** — Do not claim anonymization, consent, or compliance without evidence supporting the exact processing.

## Privacy Notices, Consent, and Preferences

_Consolidated from `quality standards/10-privacy/03-notices-consent-and-preferences.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-87EEAE5A** — Describe what data is collected, why, from whom, how used, shared, retained, protected, and controlled.
- [ ] **USEQ-B15DD85E** — Use concise contextual notice at the point of decision with access to complete detail.
- [ ] **USEQ-2144F9CE** — Keep notice language accurate for actual production data flows and suppliers.
- [ ] **USEQ-CFA62DF6** — Distinguish required processing from optional analytics, personalization, advertising, sharing, or communication.
- [ ] **USEQ-EDB93290** — Use consent only when it is an appropriate legal and ethical basis.
- [ ] **USEQ-5F4C0EE4** — Make consent specific, informed, unambiguous, granular, freely given, and recorded where required.
- [ ] **USEQ-BB6D6EB2** — Do not bundle unrelated purposes or condition service on unnecessary consent.
- [ ] **USEQ-19FC62D9** — Prevent nonessential processing before valid permission where required.
- [ ] **USEQ-3AB0815F** — Make refusal and withdrawal as easy and understandable as acceptance.
- [ ] **USEQ-215CE910** — Propagate preference changes promptly to all relevant systems, suppliers, caches, and future processing.
- [ ] **USEQ-539BFF83** — Define behavior for previously collected data after withdrawal.
- [ ] **USEQ-3B78A945** — Record notice version, purpose, choice, time, context, and proof without collecting excessive evidence.
- [ ] **USEQ-BCD2092F** — Allow users to review and change preferences through an accessible interface.
- [ ] **USEQ-5BF23D4C** — Avoid manipulative hierarchy, preselection, repeated nagging, misleading toggles, or false necessity.
- [ ] **USEQ-3E3BB1A3** — Test comprehension and operation with representative users and assistive technologies.
- [ ] **USEQ-379AB97F** — Re-consent only when purpose or processing materially changes and the basis requires it.
- [ ] **USEQ-379D10F1** — Avoid using inactivity as agreement.
- [ ] **USEQ-4DA9C46D** — Handle conflicting preferences across devices, accounts, organizations, and channels explicitly.
- [ ] **USEQ-B2092183** — Monitor preference enforcement and investigate processing that contradicts user choice.
- [ ] **USEQ-80938132** — Retire obsolete consent records according to retention and audit needs.

## Data Rights, Retention, Deletion, and Portability

_Consolidated from `quality standards/10-privacy/04-data-rights-retention-deletion-and-portability.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-849F8B32** — Determine which rights and response timelines apply to each user, data category, and jurisdiction.
- [ ] **USEQ-09FCED67** — Provide accessible request channels and avoid unnecessary barriers.
- [ ] **USEQ-7F3E7E81** — Verify requester identity proportionately without collecting excessive new personal data.
- [ ] **USEQ-7D31E28F** — Search all relevant primary, derived, archived, support, analytics, model, supplier, and export systems.
- [ ] **USEQ-5525702D** — Prevent disclosure of another person's or tenant's data in access and portability responses.
- [ ] **USEQ-E19DD9C1** — Provide data in understandable and machine-usable form where required.
- [ ] **USEQ-6C6EDCD8** — Correct inaccurate data and propagate material corrections downstream.
- [ ] **USEQ-94E79DB4** — Apply restriction and objection to future processing, not merely visible UI.
- [ ] **USEQ-966E836A** — Define retention from purpose, law, contract, dispute, security, and operational need.
- [ ] **USEQ-E1AD559F** — Delete or irreversibly deidentify data when retention expires unless a documented hold applies.
- [ ] **USEQ-86610CCA** — Propagate deletion to caches, indexes, replicas, derived data, exports, test copies, and processors.
- [ ] **USEQ-0B8FCCA8** — Document backup deletion behavior and prevent restored backups from silently resurrecting deleted use.
- [ ] **USEQ-B61915BC** — Record legal holds with scope, authority, owner, review, and release.
- [ ] **USEQ-0A64E949** — Track requests, deadlines, exceptions, decisions, and completion evidence.
- [ ] **USEQ-F10C6116** — Notify relevant recipients of correction or deletion where required.
- [ ] **USEQ-2D7549DE** — Test rights workflows end to end using representative accounts and complex data relationships.
- [ ] **USEQ-855E3643** — Prevent account deletion from becoming a deceptive cancellation substitute or vice versa.
- [ ] **USEQ-BAFDDC4B** — Preserve minimum evidence needed to demonstrate completion without recreating deleted profiles.
- [ ] **USEQ-16CE0274** — Measure timeliness, completeness, error, rework, and downstream propagation.
- [ ] **USEQ-7C6F29CF** — Review edge cases involving minors, joint accounts, employers, organizations, fraud, and litigation.

## De-Identification, Pseudonymization, and Anonymization

_Consolidated from `quality standards/10-privacy/05-deidentification-pseudonymization-and-anonymization.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-45B3FEAE** — Define the intended threat model, recipients, data context, linkage sources, and acceptable reidentification risk.
- [ ] **USEQ-56568A68** — Distinguish masking, tokenization, pseudonymization, aggregation, synthetic data, and anonymization accurately.
- [ ] **USEQ-D923770C** — Do not treat removal of direct identifiers as sufficient anonymization.
- [ ] **USEQ-69EB450F** — Assess uniqueness, sparsity, rare combinations, free text, location, time, device, and behavioral linkage risk.
- [ ] **USEQ-52C69AFA** — Minimize attributes and precision before applying complex transformation.
- [ ] **USEQ-BAB3D94B** — Keep reidentification keys separate, access-controlled, audited, and purpose-limited.
- [ ] **USEQ-15335623** — Use different pseudonyms across contexts when cross-context linkage is unnecessary.
- [ ] **USEQ-D7788ABA** — Protect mappings and quasi-identifiers with controls appropriate to the original sensitivity.
- [ ] **USEQ-F26B6FC0** — Evaluate attacks using realistic auxiliary information and motivated recipients.
- [ ] **USEQ-6F3EE651** — Measure utility and privacy together and document the trade-off.
- [ ] **USEQ-A8E1D86B** — Apply contractual and policy prohibitions on reidentification in addition to technical controls.
- [ ] **USEQ-A75CE960** — Review outputs, small groups, outliers, queries, and repeated releases for differencing risk.
- [ ] **USEQ-51D461F5** — Use privacy-preserving query controls, noise, thresholds, or aggregation where appropriate.
- [ ] **USEQ-A496904D** — Validate synthetic data for memorization, leakage, bias, and unrealistic artifacts.
- [ ] **USEQ-5B9EFCA5** — Reassess risk when new data sources, recipients, methods, or public information appear.
- [ ] **USEQ-0B0862CD** — Avoid irreversible public release unless risk is demonstrably acceptable.
- [ ] **USEQ-FBD875B1** — Label transformed data and permitted uses clearly.
- [ ] **USEQ-7759F06C** — Monitor misuse and downstream redistribution.
- [ ] **USEQ-09D0C73F** — Delete or rotate mappings according to purpose and retention.
- [ ] **USEQ-2CDB3A37** — Make every anonymization claim bounded to context and supported by current evidence.

## Cross-Border Transfers, Processors, and Data Sharing

_Consolidated from `quality standards/10-privacy/06-cross-border-transfers-processors-and-sharing.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-EFA1156D** — Inventory every destination country, provider, subprocessor, recipient, support location, backup region, and remote-access path.
- [ ] **USEQ-E4AB060F** — Determine legal, contractual, residency, localization, and government-access constraints before transfer.
- [ ] **USEQ-BF01DCD0** — Document the purpose, data categories, frequency, security, retention, and onward-transfer rights.
- [ ] **USEQ-DFD1BBE8** — Transfer only the minimum data required for the approved purpose.
- [ ] **USEQ-E8F7E7D4** — Use contracts that restrict use, require security, control subprocessors, support rights, require deletion, and define incident notification.
- [ ] **USEQ-D90C9EFB** — Assess provider technical, organizational, legal, and jurisdictional risk proportionate to sensitivity.
- [ ] **USEQ-CD9A1808** — Encrypt transfers and stored data where required by risk.
- [ ] **USEQ-3919B6AF** — Restrict provider and support access through least privilege and audited workflows.
- [ ] **USEQ-FF7405F1** — Monitor provider, subprocessor, region, terms, and data-use changes.
- [ ] **USEQ-0BA5E8E2** — Prevent suppliers from using customer data for unrelated analytics, advertising, or model training without authorization.
- [ ] **USEQ-B806C61C** — Provide data-subject and customer transparency where required.
- [ ] **USEQ-5E302DDE** — Ensure deletion, correction, restriction, and consent changes propagate to processors.
- [ ] **USEQ-234514B3** — Maintain export and return capability for provider transition.
- [ ] **USEQ-5402DD1F** — Define retention and verified deletion at contract end.
- [ ] **USEQ-65334F09** — Revoke credentials, integrations, network access, and keys on termination.
- [ ] **USEQ-3C8EBF98** — Test provider outage, breach, data recovery, and exit procedures.
- [ ] **USEQ-8A899654** — Review bulk exports and human-readable support access as transfers.
- [ ] **USEQ-27E1C5E3** — Maintain evidence of transfer assessments, contracts, and supplementary controls.
- [ ] **USEQ-F9E745E7** — Reassess when law, provider ownership, infrastructure, processing purpose, or data sensitivity changes.
- [ ] **USEQ-0135FDA3** — Stop transfer when required safeguards cannot be maintained.

## Children’s, Sensitive, and High-Risk Data

_Consolidated from `quality standards/10-privacy/07-childrens-sensitive-and-high-risk-data.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-71E5442B** — Identify sensitive categories by jurisdiction, context, user expectation, harm, and inference—not labels alone.
- [ ] **USEQ-7A4F495F** — Avoid collecting sensitive data unless it is necessary for a documented high-value purpose.
- [ ] **USEQ-6231CEB2** — Use age-appropriate design and language when children may use or be affected by the service.
- [ ] **USEQ-4B4EEF3F** — Determine age assurance, parental authorization, and child autonomy requirements proportionately.
- [ ] **USEQ-092E7F6B** — Do not use dark patterns to encourage disclosure or extended engagement by children.
- [ ] **USEQ-7EAF6F92** — Disable unnecessary profiling, advertising, location, contact discovery, and public visibility by default for minors.
- [ ] **USEQ-B3F64417** — Protect biometric templates, health data, precise location, communications, identity documents, and financial data with enhanced access and encryption.
- [ ] **USEQ-E1B3DE1D** — Avoid deriving sensitive traits from ordinary behavior unless explicitly justified and governed.
- [ ] **USEQ-A3672E68** — Limit retention and secondary use more strictly than ordinary data.
- [ ] **USEQ-C1A114B0** — Prevent support, analytics, logs, and test systems from receiving unnecessary sensitive values.
- [ ] **USEQ-19E31877** — Require independent review for new sensitive-data processing or sharing.
- [ ] **USEQ-15A125C3** — Provide meaningful alternatives where refusal is possible.
- [ ] **USEQ-6D30DBF0** — Design for family, guardian, school, employer, and joint-account conflicts without assuming one actor always represents the user.
- [ ] **USEQ-C241B622** — Test disclosure, inference, enumeration, coercion, stalking, and account-recovery risks.
- [ ] **USEQ-DD1AEB92** — Provide rapid revocation, correction, deletion, and safety escalation.
- [ ] **USEQ-FE82D430** — Monitor for misuse, scraping, trafficking, discrimination, and harmful targeting.
- [ ] **USEQ-F8676A18** — Avoid public or permanent identifiers for minors and vulnerable users.
- [ ] **USEQ-474C667C** — Use vendors with appropriate contractual and technical safeguards.
- [ ] **USEQ-01769764** — Review safety and privacy as users age or their legal status changes.
- [ ] **USEQ-A9BDCB7D** — Stop processing when harm cannot be reduced below approved tolerance.

## Standards and source references

- [ISO/IEC 27701:2025 — Privacy information management systems](https://www.iso.org/standard/85819.html)
- [ISO/IEC 29100:2024 — Privacy framework](https://www.iso.org/standard/85938.html)
- [NIST Privacy Framework 1.0](https://www.nist.gov/privacy-framework)
- [ISO/IEC 25012:2008 — Data quality model](https://www.iso.org/standard/35736.html)
- [ISO/IEC 27001:2022 — Information security management systems](https://www.iso.org/standard/27001)
- [ISO 31000:2018 — Risk management guidelines](https://www.iso.org/standard/65694.html)

---

[Previous phase](08-security-and-cryptography.md) · [Next: Phase 10: Verification and testing](10-verification-and-testing.md)
