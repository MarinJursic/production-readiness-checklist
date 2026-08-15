# References and scope

The checklist is technology-neutral and synthesizes lifecycle governance, product quality, software engineering, secure development, accessibility, software supply chain, identity, incident response, operations, and reliability practices.

## Structure rationale

The 16-phase sequence follows the whole-lifecycle scope of ISO/IEC/IEEE 12207 and the breadth of the SWEBOK Guide while remaining usable as one navigable review. ISO/IEC 25010 provides a product-quality lens, and NIST SSDF keeps secure development outcome- and risk-oriented. WCAG, OWASP ASVS, NIST AI RMF, and other specialist sources inform triggered domains. The production tracks remain the final assessment of the exact artifact and operating environment.

The sequence is not a prescribed development process. Teams may work iteratively or concurrently and revisit controls when assumptions, architecture, dependencies, data, or release state change.

## Standards snapshot

The integrated checklist was reviewed on **August 15, 2026**. Verify current versions and applicability before making the controls a formal engineering or release gate.

- [ISO/IEC/IEEE 12207:2026 — software lifecycle processes](https://www.iso.org/standard/90219.html)
- [IEEE Computer Society — SWEBOK Guide V4.0 knowledge areas](https://www.computer.org/education/bodies-of-knowledge/software-engineering/topics)
- [ISO/IEC 25010:2023 — product quality model](https://www.iso.org/standard/78176.html)
- [NIST SP 800-218 — Secure Software Development Framework 1.1](https://csrc.nist.gov/pubs/sp/800/218/final)
- [OWASP Application Security Verification Standard](https://owasp.org/www-project-application-security-verification-standard/)
- [W3C Web Content Accessibility Guidelines 2.2](https://www.w3.org/TR/WCAG22/)
- [SLSA specification 1.2](https://slsa.dev/spec/v1.2/)
- [NIST SP 800-63 Revision 4 — Digital Identity Guidelines](https://pages.nist.gov/800-63-4/)
- [NIST SP 800-61 Revision 3 — Incident Response Recommendations](https://csrc.nist.gov/pubs/sp/800/61/r3/final)
- [FIRST CVSS v4.0 implementation guidance](https://www.first.org/cvss/v4.0/implementation-guide)
- [CISA software bill of materials resources](https://www.cisa.gov/topics/information-communications-technology-supply-chain-security/sbom)
- [Google SRE Workbook — implementing SLOs](https://sre.google/workbook/implementing-slos/)
- [DORA — software delivery performance metrics](https://dora.dev/guides/dora-metrics/)
- [web.dev — Core Web Vitals](https://web.dev/articles/vitals)
- [NIST AI Risk Management Framework](https://www.nist.gov/itl/ai-risk-management-framework)
- [OWASP Artificial Intelligence Security Verification Standard](https://owasp.org/www-project-artificial-intelligence-security-verification-standard-aisvs-docs/)

## Frequently triggered legal and regulatory references

- [EU General Data Protection Regulation](https://eur-lex.europa.eu/eli/reg/2016/679/oj/eng)
- [PCI Security Standards Council document library](https://www.pcisecuritystandards.org/document_library/)
- [US HHS HIPAA Security Rule resources](https://www.hhs.gov/hipaa/for-professionals/security/index.html)
- [US FTC Children’s Online Privacy Protection Rule](https://www.ftc.gov/legal-library/browse/rules/childrens-online-privacy-protection-rule-coppa)
- [European Accessibility Act](https://eur-lex.europa.eu/eli/dir/2019/882/oj/eng)

These links are starting points, not a complete applicability analysis.

## Limitations

- This project is not legal advice, certification, or a guarantee of security or reliability.
- Controls must be tailored to the application’s users, data, architecture, markets, contracts, and risk.
- A source-code review cannot substitute for production evidence, operating drills, qualified human review, or accountable risk acceptance.
- Standards, laws, threat techniques, provider behavior, and dependencies change over time.
- Some high-risk systems need controls beyond this checklist.

Use the [source consolidation manifest](engineering/source-manifest.md) to trace the supplied archive documents, import counts, deduplication rules, and archive hashes.

When requirements conflict, document the conflict, obtain qualified review, and record who approved the resolution.
