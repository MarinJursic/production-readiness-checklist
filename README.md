<div align="center">

<img src="docs/assets/social-card.png" alt="Production Readiness Checklist — Ship with evidence, not optimism" width="100%">

# Production Readiness Checklist

**1,421 evidence-driven checks for shipping web applications with confidence.**

[![Checklist](https://img.shields.io/badge/checks-1%2C421-2563eb)](docs/checklists/01-release-foundations.md)
[![Validate](https://github.com/MarinJursic/production-readiness-checklist/actions/workflows/validate.yml/badge.svg)](https://github.com/MarinJursic/production-readiness-checklist/actions/workflows/validate.yml)
[![Docs](https://img.shields.io/badge/docs-GitHub%20Pages-0f766e)](https://marinjursic.github.io/production-readiness-checklist/)
[![GitHub stars](https://img.shields.io/github/stars/MarinJursic/production-readiness-checklist?logo=github&color=f59e0b)](https://github.com/MarinJursic/production-readiness-checklist)
[![License: MIT](https://img.shields.io/badge/license-MIT-f59e0b.svg)](LICENSE)

[Start the checklist](docs/checklists/01-release-foundations.md) · [Use with Claude or another AI](docs/guides/ai-assisted-review.md) · [Copy a release template](docs/records/release-assessment.md)

⭐ [Star this project](https://github.com/MarinJursic/production-readiness-checklist) · [Share on LinkedIn](https://www.linkedin.com/sharing/share-offsite/?url=https%3A%2F%2Fmarinjursic.github.io%2Fproduction-readiness-checklist%2F) · [Share on X](https://twitter.com/intent/tweet?text=Production%20Readiness%20Checklist%3A%201%2C421%20evidence-driven%20checks%20for%20shipping%20web%20applications%20with%20confidence.&url=https%3A%2F%2Fmarinjursic.github.io%2Fproduction-readiness-checklist%2F)

</div>

---

Production readiness is not “all tests are green.” It is a documented decision that the exact artifact you intend to ship is reliable, secure, supportable, recoverable, and within your organization’s risk tolerance.

This repository turns a comprehensive master checklist into a practical review system. The controls are split into focused tracks, each control has a stable ID, and ready-to-copy templates help teams attach owners, evidence, exceptions, and sign-offs.

> [!IMPORTANT]
> No checklist or AI review can prove that a nontrivial application has zero defects. A credible readiness claim requires current evidence, explicit risk ownership, and tested recovery paths.

## Start here

1. Copy [the release assessment template](docs/records/release-assessment.md) into your project or release workspace.
2. Review the [immediate no-go conditions](docs/checklists/01-release-foundations.md#2-immediate-no-go-conditions) before spending time on the full assessment.
3. Mark each track as required, not applicable, or deferred. A “not applicable” decision needs a written reason.
4. Evaluate every applicable control as **Pass**, **Fail**, **Blocked**, or **Not Applicable**, with an owner and current evidence.
5. Complete the [evidence and decision track](docs/checklists/10-evidence-and-decision.md) and record the final decision.

New to the process? Read the [15-minute quick start](docs/guides/getting-started.md).

## Checklist tracks

| Track | What it covers |
| --- | --- |
| [Readiness principle](docs/checklists/00-readiness-principle.md) | What a defensible “production-ready” claim means |
| [1. Release foundations](docs/checklists/01-release-foundations.md) | Operating rules, no-go conditions, release identity, and scope |
| [2. Product, risk, and architecture](docs/checklists/02-product-risk-architecture.md) | Product intent, ownership, risk assessment, architecture, and threat modeling |
| [3. Source, build, and supply chain](docs/checklists/03-source-build-supply-chain.md) | Source control, CI/CD, provenance, dependencies, SBOMs, and licenses |
| [4. Environments, quality, and experience](docs/checklists/04-environments-quality-experience.md) | Configuration, secrets, correctness, testing, frontend behavior, and accessibility |
| [5. Application security](docs/checklists/05-application-security.md) | APIs, identity, authorization, sessions, input, files, transport, and cryptography |
| [6. Data, privacy, and performance](docs/checklists/06-data-privacy-performance.md) | Data integrity, migrations, privacy, performance, capacity, and overload |
| [7. Reliability and operations](docs/checklists/07-reliability-operations.md) | Resilience, infrastructure, observability, incidents, recovery, and deployment |
| [8. Maintenance, vendors, and compliance](docs/checklists/08-maintenance-vendors-compliance.md) | Long-term operation, third parties, legal applicability, and regulatory modules |
| [9. Conditional feature modules](docs/checklists/09-conditional-modules.md) | Payments, SaaS, UGC, notifications, i18n, PWA, real-time, admin, AI, and safety-critical systems |
| [10. Evidence, sign-off, and decision](docs/checklists/10-evidence-and-decision.md) | Evidence package, required sign-offs, go/no-go rules, and declaration |

## Use it with Claude or another coding agent

The repository includes a root [`CLAUDE.md`](CLAUDE.md) and purpose-built prompts. They teach an agent to distinguish verified evidence from assumptions, cite repository evidence, and return gaps instead of silently marking controls as passed.

Example request:

```text
Review this repository for production readiness using the instructions in CLAUDE.md.
Start with the no-go controls and the tracks relevant to this application.
Do not modify code. Return an evidence-linked report and list every unknown separately.
```

Choose a workflow:

- [Full repository readiness review](docs/prompts/full-readiness-review.md)
- [Release-diff review](docs/prompts/release-diff-review.md)
- [Evidence challenge](docs/prompts/evidence-challenge.md)
- [AI-assisted review guide](docs/guides/ai-assisted-review.md)

AI can inspect code, tests, infrastructure definitions, and documentation. It usually cannot prove production configuration, organizational readiness, legal applicability, restore success, or real-world operational behavior without external evidence. Those controls must remain **Blocked** or **Unknown** until a person supplies evidence.

## Evidence, not checkbox theater

A checked box without evidence is only a claim. For each control, record:

| Field | Example |
| --- | --- |
| Control | `PRC-34-017` |
| Status | Pass / Fail / Blocked / Not Applicable |
| Owner | One accountable person |
| Evidence | Test run, configuration export, dashboard, decision record, or drill report |
| Release | Commit, artifact digest, migrations, configuration, and flag state |
| Environment | The exact environment tested |
| Reviewed | Reviewer, date, and evidence expiry |
| Exception | Risk owner, compensating control, remediation date, and expiry |

Use the [evidence record](docs/records/evidence-record.md) and [risk exception](docs/records/risk-exception.md) templates to keep reviews auditable.

## What is included

```text
.
├── CLAUDE.md                  # Guardrails for AI-assisted audits
├── docs/
│   ├── checklists/            # 10 focused tracks + readiness principle
│   ├── guides/                # Adoption and AI-review guidance
│   ├── prompts/               # Copy-ready agent review prompts
│   └── records/               # Release, evidence, risk, and decision templates
├── scripts/                   # Repository integrity validation
└── .github/                   # Issue forms, PR template, and CI workflows
```

## Standards snapshot

The original checklist was compiled against a standards snapshot dated **August 14, 2026**, including NIST SSDF 1.1, OWASP ASVS 5.0.0, WCAG 2.2, SLSA 1.2, NIST SP 800-63 Revision 4, and NIST SP 800-61 Revision 3. Standards and laws evolve; verify current applicability before using the checklist as a release gate. See [references and scope](docs/references.md).

This project is engineering guidance, not legal advice, a certification, or a substitute for qualified security, privacy, accessibility, safety, or compliance review.

## Contributing

Corrections, clearer controls, stronger evidence examples, and additional standards mappings are welcome. Read [CONTRIBUTING.md](CONTRIBUTING.md) before opening a pull request. Please report vulnerabilities privately as described in [SECURITY.md](SECURITY.md).

## Help it reach more teams

If this checklist prevents one avoidable production incident, consider [starring the repository](https://github.com/MarinJursic/production-readiness-checklist), sharing the [documentation site](https://marinjursic.github.io/production-readiness-checklist/), or sending it to a team preparing a launch. GitHub also exposes a **Cite this repository** action using [CITATION.cff](CITATION.cff).

Released under the [MIT License](LICENSE).
