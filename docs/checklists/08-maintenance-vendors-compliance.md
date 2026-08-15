# Maintenance, vendors, and compliance

> Keep the service supportable and verify third-party, legal, and regulatory obligations.

Sections 36–38 of the master checklist. For each applicable item, capture status, owner, evidence, release, environment, reviewer, and evidence age.

## 36. Maintainability and long-term operability

- [ ] **PRC-36-001** — A current README or equivalent operating entry point exists.
- [ ] **PRC-36-002** — Architecture, data flow, dependencies, deployment, and recovery are documented.
- [ ] **PRC-36-003** — Local or isolated development setup is reproducible.
- [ ] **PRC-36-004** — Build and test commands are documented.
- [ ] **PRC-36-005** — Coding, review, testing, and release conventions are documented.
- [ ] **PRC-36-006** — Critical code has appropriate maintainers.
- [ ] **PRC-36-007** — Complexity and tightly coupled components have been reviewed.
- [ ] **PRC-36-008** — Dead code and obsolete paths are removed.
- [ ] **PRC-36-009** — Tests are understandable and maintainable.
- [ ] **PRC-36-010** — Critical behavior is not protected only by brittle end-to-end tests.
- [ ] **PRC-36-011** — Dependency-update ownership and cadence exist.
- [ ] **PRC-36-012** — Runtime and platform end-of-support dates are tracked.
- [ ] **PRC-36-013** — Technical debt that affects reliability or security is recorded.
- [ ] **PRC-36-014** — Temporary workarounds have owners and expiry dates.
- [ ] **PRC-36-015** — Feature flags and compatibility shims have removal plans.
- [ ] **PRC-36-016** — API and data deprecation policies exist.
- [ ] **PRC-36-017** — Operational toil is measured and reduced.
- [ ] **PRC-36-018** — Routine maintenance does not require broad production privileges.
- [ ] **PRC-36-019** — Support documentation and troubleshooting guidance are current.
- [ ] **PRC-36-020** — Support personnel are trained on the released behavior.
- [ ] **PRC-36-021** — Support access to customer data is restricted and audited.
- [ ] **PRC-36-022** — Known issues are communicated accurately.
- [ ] **PRC-36-023** — Release notes describe user-relevant changes.
- [ ] **PRC-36-024** — Maintenance windows and status communications are defined.
- [ ] **PRC-36-025** — Ownership has sufficient staffing and knowledge redundancy.
- [ ] **PRC-36-026** — Critical knowledge is not held by one person.
- [ ] **PRC-36-027** — Decommissioning, data export, retention, and secure deletion have a documented lifecycle plan.

---

## 37. Third-party and vendor readiness

- [ ] **PRC-37-001** — Maintain an inventory of all third parties involved in delivery or data processing.
- [ ] **PRC-37-002** — Record each provider’s purpose, owner, data access, permissions, regions, and criticality.
- [ ] **PRC-37-003** — Security and privacy due diligence is proportionate to risk.
- [ ] **PRC-37-004** — Contracts cover service levels, security, permitted data use, breach notification, deletion, portability, and termination.
- [ ] **PRC-37-005** — Subprocessors are known where required.
- [ ] **PRC-37-006** — Provider credentials and scopes follow least privilege.
- [ ] **PRC-37-007** — Sandbox and production accounts are separated.
- [ ] **PRC-37-008** — Integration secrets are rotatable.
- [ ] **PRC-37-009** — Provider quotas, rate limits, and concurrency limits are understood.
- [ ] **PRC-37-010** — Provider availability and support commitments align with the product’s objectives.
- [ ] **PRC-37-011** — Timeouts, retries, circuit breaking, fallback, and reconciliation are implemented.
- [ ] **PRC-37-012** — Provider responses are treated as untrusted input.
- [ ] **PRC-37-013** — Webhooks and callbacks are authenticated and replay-protected.
- [ ] **PRC-37-014** — Provider API versions and deprecation dates are monitored.
- [ ] **PRC-37-015** — Breaking-change notifications reach an accountable owner.
- [ ] **PRC-37-016** — Provider outages and malformed responses have been tested.
- [ ] **PRC-37-017** — Provider data-use changes are monitored.
- [ ] **PRC-37-018** — A kill switch can disable a compromised or harmful integration.
- [ ] **PRC-37-019** — A replacement, exit, or data-export plan exists for critical providers.
- [ ] **PRC-37-020** — The app can identify records affected by a provider incident.
- [ ] **PRC-37-021** — Provider status and escalation contacts are included in incident runbooks.
- [ ] **PRC-37-022** — Terminated providers lose credentials and data access.
- [ ] **PRC-37-023** — Vendor claims relied upon for compliance are evidenced rather than assumed.
- [ ] **PRC-37-024** — No provider can silently expand data collection beyond approved behavior.

---

## 38. Legal and compliance applicability

This section is a control framework, not legal advice. Applicability must be determined using qualified legal and compliance review for the app’s actual users, markets, data, transactions, and sector.

- [ ] **PRC-38-001** — List every jurisdiction in which the service is offered, operated, marketed, or monitors users.
- [ ] **PRC-38-002** — Identify applicable privacy and data-protection laws.
- [ ] **PRC-38-003** — Identify applicable accessibility laws and procurement requirements.
- [ ] **PRC-38-004** — Identify consumer-protection and unfair-practice requirements.
- [ ] **PRC-38-005** — Identify subscription, renewal, cancellation, refund, and pricing requirements.
- [ ] **PRC-38-006** — Identify marketing, advertising, email, SMS, telemarketing, and consent requirements.
- [ ] **PRC-38-007** — Identify payment, banking, financial-service, insurance, or securities requirements.
- [ ] **PRC-38-008** — Identify healthcare, education, employment, housing, government, or other sector requirements.
- [ ] **PRC-38-009** — Identify children’s and age-appropriate-design requirements.
- [ ] **PRC-38-010** — Identify identity-verification, fraud, and anti-money-laundering obligations where applicable.
- [ ] **PRC-38-011** — Identify record-retention and legal-hold obligations.
- [ ] **PRC-38-012** — Identify tax, invoice, currency, and marketplace obligations.
- [ ] **PRC-38-013** — Identify export-control, sanctions, and restricted-country requirements.
- [ ] **PRC-38-014** — Identify data-localization and international-transfer requirements.
- [ ] **PRC-38-015** — Identify content, intellectual-property, copyright, moderation, and intermediary obligations.
- [ ] **PRC-38-016** — Identify cybersecurity, incident-reporting, and breach-notification obligations.
- [ ] **PRC-38-017** — Identify algorithmic-decision and AI-specific obligations.
- [ ] **PRC-38-018** — Identify contractual and customer-specific security requirements.
- [ ] **PRC-38-019** — Terms of service match actual product behavior.
- [ ] **PRC-38-020** — Privacy and cookie notices match actual data behavior.
- [ ] **PRC-38-021** — Acceptable-use and content policies are enforceable and communicated.
- [ ] **PRC-38-022** — Data-processing agreements and subprocessor disclosures are current.
- [ ] **PRC-38-023** — Open-source notices and license obligations are fulfilled.
- [ ] **PRC-38-024** — Trademarks, media, data sets, fonts, and other content have appropriate rights.
- [ ] **PRC-38-025** — Security and compliance claims are accurate and supportable.
- [ ] **PRC-38-026** — Certifications are not represented beyond their actual scope.
- [ ] **PRC-38-027** — Accessibility claims are supported by evidence.
- [ ] **PRC-38-028** — Required contact, company, pricing, cancellation, and complaint information is displayed.
- [ ] **PRC-38-029** — Customer contracts and SLAs do not promise performance beyond demonstrated capacity.
- [ ] **PRC-38-030** — Insurance requirements and notification procedures are understood.
- [ ] **PRC-38-031** — Legal review has approved material launch changes.
- [ ] **PRC-38-032** — Renewal, review, and filing deadlines are tracked after launch.

### 38.1 Frequently triggered regulatory modules

- [ ] **PRC-38-033** — **Payments:** apply current PCI DSS requirements to the actual cardholder-data environment and minimize scope through tokenization or hosted collection where suitable. PCI DSS v4.0.1 is the current published revision in the PCI SSC document library. ([pcisecuritystandards.org](https://www.pcisecuritystandards.org/document_library/))
- [ ] **PRC-38-034** — **Health information:** determine whether HIPAA or another health-information regime applies and implement required administrative, physical, and technical safeguards. ([hhs.gov](https://www.hhs.gov/hipaa/for-professionals/security/index.html))
- [ ] **PRC-38-035** — **Children:** determine whether the product is directed to children or has actual knowledge of child users and apply parental-consent, minimization, advertising, safety, and deletion requirements. ([ftc.gov](https://www.ftc.gov/legal-library/browse/rules/childrens-online-privacy-protection-rule-coppa))
- [ ] **PRC-38-036** — **European accessibility:** determine whether the European Accessibility Act or related national rules apply to the service. ([eur-lex.europa.eu](https://eur-lex.europa.eu/eli/dir/2019/882/oj/eng))

---
