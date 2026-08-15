# Complete engineering review

This is the start-to-finish path for evaluating a software product across its full lifecycle. Complete phases 1–16 in order, then run the existing production-readiness review before release approval.

The integrated corpus contains **10,042 unique controls**: **8,621 lifecycle and quality controls** with `USEQ-` IDs plus **1,421 production-readiness controls** with `PRC-` IDs. The import reviewed 25,359 checkbox lines from 215 source documents and removed repeated boilerplate, explicit consolidated copies, mirrored production controls, and exact overlaps.

## How to use the sequence

1. Define the product, lifecycle, organization, release, and evidence boundaries.
2. Work through every phase in order; use Not Applicable only with a reviewed rationale.
3. Apply every specialized module whose trigger exists.
4. Preserve evidence against the stable `USEQ-` or `PRC-` identifier.
5. Finish with the production no-go screen, release evidence, sign-offs, and deployment verification.

The phases are ordered for navigation, not as a waterfall mandate. Iterative teams can revisit them continuously, but final approval still requires every applicable control to have a disposition.

## Why this structure

The sequence follows the whole-lifecycle scope of [ISO/IEC/IEEE 12207:2026](https://www.iso.org/standard/90219.html) and groups implementation concerns using the [SWEBOK Guide V4.0 knowledge areas](https://www.computer.org/education/bodies-of-knowledge/software-engineering/topics). Quality is treated as a set of explicit product attributes using [ISO/IEC 25010:2023](https://www.iso.org/standard/78176.html), while secure development remains outcome- and risk-based in line with the [NIST Secure Software Development Framework](https://csrc.nist.gov/Projects/ssdf).

Specialized controls remain in the same journey rather than becoming separate products. Accessibility uses [WCAG 2.2](https://www.w3.org/TR/WCAG22/), application security maps to [OWASP ASVS](https://owasp.org/www-project-application-security-verification-standard/), and AI risk work is informed by the [NIST AI Risk Management Framework](https://www.nist.gov/itl/ai-risk-management-framework). The [references page](../references.md) explains scope and interpretation limits.

## Lifecycle phases

| Phase | Category | Unique controls |
| ---: | --- | ---: |
| 1 | [Governance and foundations](01-governance-and-foundations.md) | 638 |
| 2 | [Product and requirements](02-product-and-requirements.md) | 105 |
| 3 | [User experience, web, and content](03-user-experience-web-and-content.md) | 1,037 |
| 4 | [Architecture and design](04-architecture-and-design.md) | 241 |
| 5 | [Code quality and implementation](05-code-quality-and-implementation.md) | 839 |
| 6 | [Application services and APIs](06-application-services-and-apis.md) | 332 |
| 7 | [Data and information lifecycle](07-data-and-information-lifecycle.md) | 778 |
| 8 | [Security and cryptography](08-security-and-cryptography.md) | 916 |
| 9 | [Privacy and data protection](09-privacy-and-data-protection.md) | 136 |
| 10 | [Verification and testing](10-verification-and-testing.md) | 636 |
| 11 | [Developer experience, platform, and delivery](11-developer-experience-platform-and-delivery.md) | 752 |
| 12 | [Operations, SRE, and support](12-operations-sre-and-support.md) | 581 |
| 13 | [Documentation and knowledge](13-documentation-and-knowledge.md) | 241 |
| 14 | [Trust, safety, and ecosystems](14-trust-safety-and-ecosystems.md) | 235 |
| 15 | [AI, ML, and AI-assisted development](15-ai-ml-and-ai-assisted-development.md) | 632 |
| 16 | [Specialized domains and release assurance](16-specialized-domains-and-release-assurance.md) | 522 |
| Final release review | [Production readiness](../checklists/01-release-foundations.md) | 1,421 |

## Decision rule

> Completion is not a score. One material failure can block approval regardless of how many unrelated controls pass.

Use the [source consolidation manifest](source-manifest.md) to trace each imported source document to its destination and understand what was removed as duplicate or already covered.

[Begin phase 1: Governance and foundations](01-governance-and-foundations.md)
