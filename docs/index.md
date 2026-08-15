---
description: 1,421 evidence-driven checks for shipping secure, reliable, and supportable web applications.
---

# Production Readiness Checklist

Use this documentation to decide whether a web application is ready for production. It provides **1,421 technology-neutral controls**, stable evidence IDs, reusable decision records, and guidance for human or AI-assisted reviews.

The checklist does not produce a readiness score. A release is ready only when every applicable requirement has current evidence and no unresolved risk exceeds the organization's tolerance.

## Start here

| Your goal | Start with | What you will produce |
| --- | --- | --- |
| Review an upcoming release | [15-minute quick start](guides/getting-started.md) | A scoped assessment and an initial no-go decision |
| Ask Claude or another agent to review a repository | [AI-assisted review](guides/ai-assisted-review.md) | Evidence-backed findings with explicit unknowns |
| Introduce the process to a team | [Readiness principle](checklists/00-readiness-principle.md) | Agreed owners, evidence rules, and decision criteria |
| Inspect the complete control set | [Release foundations](checklists/01-release-foundations.md) | A track-by-track production-readiness review |

## Review process

1. Identify the exact artifact, configuration, migrations, flags, environment, and release scope.
2. Check the [immediate no-go conditions](checklists/01-release-foundations.md#2-immediate-no-go-conditions).
3. Select the applicable core and conditional tracks.
4. Collect current evidence for every applicable control.
5. Record failed, blocked, and not-applicable controls without converting them into a score.
6. Obtain the required independent reviews and accountable sign-offs.
7. Record the go/no-go decision, deploy progressively, and verify production.

## Checklist tracks

| Track | Covers |
| --- | --- |
| [Release foundations](checklists/01-release-foundations.md) | Release identity, scope, ownership, and immediate no-go conditions |
| [Product, risk, and architecture](checklists/02-product-risk-architecture.md) | Product intent, risk ownership, architecture, and threat modeling |
| [Source, build, and supply chain](checklists/03-source-build-supply-chain.md) | Source control, CI/CD, provenance, dependencies, SBOMs, and licenses |
| [Environments, quality, and experience](checklists/04-environments-quality-experience.md) | Configuration, testing, frontend behavior, compatibility, and accessibility |
| [Application security](checklists/05-application-security.md) | Identity, authorization, input handling, transport, secrets, and cryptography |
| [Data, privacy, and performance](checklists/06-data-privacy-performance.md) | Data integrity, migrations, privacy, capacity, caching, and overload behavior |
| [Reliability and operations](checklists/07-reliability-operations.md) | Resilience, observability, recovery, deployment, rollback, and incident response |
| [Maintenance, vendors, and compliance](checklists/08-maintenance-vendors-compliance.md) | Operability, third parties, legal obligations, and regulatory modules |
| [Conditional feature modules](checklists/09-conditional-modules.md) | Payments, SaaS, AI, real-time systems, user-generated content, and other optional capabilities |
| [Evidence, sign-off, and decision](checklists/10-evidence-and-decision.md) | Evidence packages, exceptions, sign-offs, deployment verification, and go/no-go records |

## Evidence states

Use one of four states for each applicable control:

| State | Meaning |
| --- | --- |
| **Pass** | Current evidence demonstrates that the control is satisfied. |
| **Fail** | Evidence demonstrates that the control is not satisfied. |
| **Blocked** | The control cannot yet be evaluated because required access, evidence, or work is missing. |
| **Not applicable** | The control does not apply, and the reason is recorded. |

> Every applicable requirement has current evidence; no known risk exceeds the organization's tolerance; critical user journeys meet defined reliability and security objectives; and the organization can detect, contain, roll back, restore, support, and communicate failures.

## Templates and project links

- Record the review with the [release assessment](records/release-assessment.md), [evidence record](records/evidence-record.md), [risk exception](records/risk-exception.md), and [go/no-go decision](records/go-no-go-decision.md) templates.
- Review the [references, standards snapshot, and limitations](references.md).
- [Star the repository on GitHub](https://github.com/MarinJursic/production-readiness-checklist) or share the documentation on [LinkedIn](https://www.linkedin.com/sharing/share-offsite/?url=https%3A%2F%2Fmarinjursic.github.io%2Fproduction-readiness-checklist%2F) or [X](https://twitter.com/intent/tweet?text=Production%20Readiness%20Checklist%3A%201%2C421%20evidence-driven%20checks%20for%20shipping%20web%20applications%20with%20confidence.&url=https%3A%2F%2Fmarinjursic.github.io%2Fproduction-readiness-checklist%2F).
