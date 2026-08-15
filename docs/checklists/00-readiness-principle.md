# The readiness principle

> Production readiness is an evidence-backed risk decision, not a promise of perfection.

## First: what “production-ready without any issues” can honestly mean

No finite checklist, automated test suite, penetration test, or review can prove that a nontrivial web application contains **zero defects** or will experience **zero incidents**. Dependencies change, traffic behaves unexpectedly, external services fail, threat techniques evolve, and some defects are reachable only under conditions that testing did not reproduce.

A defensible production-readiness claim is:

> **Every applicable requirement has current evidence; no known risk exceeds the organization’s tolerance; critical user journeys meet defined reliability and security objectives; and the organization can detect, contain, roll back, restore, support, and communicate failures.**

That interpretation matches modern secure-development and reliability practice: NIST’s SSDF aims to reduce vulnerabilities and mitigate the impact of vulnerabilities that remain undetected, while SRE practice uses measurable service-level objectives and error budgets rather than an unsupported promise of perfect uptime. ([csrc.nist.gov](https://csrc.nist.gov/pubs/sp/800/218/final))

This checklist is technology-neutral. The complete review begins with the [16-phase engineering lifecycle](../engineering/00-overview.md) and ends with this production gate. The standards snapshot was reviewed on **August 15, 2026** and includes ISO/IEC/IEEE 12207:2026, SWEBOK Guide V4.0, ISO/IEC 25010:2023, NIST SSDF 1.1, OWASP ASVS, WCAG 2.2, SLSA 1.2, NIST SP 800-63 Revision 4, and NIST SP 800-61 Revision 3. See [references and scope](../references.md).

---
