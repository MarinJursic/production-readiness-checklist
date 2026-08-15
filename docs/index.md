---
description: 10,042 evidence-driven controls for engineering and shipping secure, reliable, and supportable software.
---

# Production Readiness Checklist

Use this documentation to evaluate a software product from initial governance and requirements through design, implementation, operations, and final production approval. It provides **10,042 technology-neutral controls**, stable evidence IDs, reusable decision records, and guidance for human or AI-assisted reviews.

The control set has two connected layers: **8,621 lifecycle and quality controls** across 16 engineering phases, followed by **1,421 production-readiness controls** for the release decision. It does not produce a readiness score. Every applicable requirement needs current evidence, and one material failure can block approval.

## Start here

| Your goal | Start with | What you will produce |
| --- | --- | --- |
| Evaluate a product from start to finish | [Complete engineering review](engineering/00-overview.md) | A disposition for every applicable lifecycle and release control |
| Check an upcoming release quickly | [15-minute quick start](guides/getting-started.md) | A scoped assessment and an initial no-go decision |
| Ask Claude or another agent to inspect a repository | [AI-assisted review](guides/ai-assisted-review.md) | Evidence-backed findings with explicit unknowns |
| Understand where the added material came from | [Source consolidation manifest](engineering/source-manifest.md) | A document-by-document import and deduplication record |

## Complete review sequence

1. Define the product, lifecycle, organization, release, and evidence boundaries.
2. Work through [engineering phases 1–16](engineering/00-overview.md), recording Pass, Fail, Blocked, or Not Applicable for every control.
3. Apply all specialized modules whose trigger exists, including accessibility, privacy, AI, safety, marketplaces, and regulated domains.
4. Run the [production no-go screen](checklists/01-release-foundations.md#2-immediate-no-go-conditions) and the ten release tracks.
5. Collect current evidence, resolve or accept risk through the proper authority, obtain sign-offs, and verify the deployed system.

The order gives reviewers a dependable path through the material; it does not require waterfall development. Teams can assess phases iteratively and revisit controls whenever a design, dependency, environment, or release changes.

## Help expand the project

This is an open-source, community-improvable control set. If you find a missing, duplicated, unclear, outdated, or incorrectly categorized control, you can [propose a checklist improvement](https://github.com/MarinJursic/production-readiness-checklist/issues/new?template=control-proposal.yml). Contributions to evidence guidance, documentation, navigation, validation, and tooling are also welcome; see the [contribution guide](https://github.com/MarinJursic/production-readiness-checklist/blob/main/CONTRIBUTING.md).

## Long-term vision: an AI readiness scanner

The project’s longer-term goal is an open-source AI-assisted scanner that can inspect a repository and its available engineering evidence against this complete control set. It should map every result to a stable control ID, cite the evidence behind each conclusion, distinguish verified facts from unknowns, and generate a prioritized list of failures, blocked checks, manual verification work, and next actions.

The scanner does not exist as a complete product today. It must remain technology-neutral, treat unavailable production and organizational evidence as unknown rather than passed, and leave risk acceptance and release approval to accountable humans. Contributors are welcome to help design and build it.

## Coverage map

| Review layer | Categories | Purpose |
| --- | --- | --- |
| Foundations | Governance, product, requirements, UX, architecture | Establish intent, scope, ownership, risk, and system design |
| Construction | Code, services, APIs, data, security, privacy, testing | Verify the product and its implementation properties |
| Delivery and operation | Developer experience, platform, delivery, SRE, support, documentation | Verify that teams can build, deploy, observe, recover, and maintain it |
| Specialized assurance | Trust and safety, ecosystems, AI/ML, specialized domains | Apply triggered controls without breaking the main review sequence |
| Production decision | Ten production-readiness tracks | Approve the exact artifact, configuration, environment, and release |

## Evidence states

| State | Meaning |
| --- | --- |
| **Pass** | Current evidence demonstrates that the control is satisfied. |
| **Fail** | Evidence demonstrates that the control is not satisfied. |
| **Blocked** | Required access, evidence, or work is missing. |
| **Not Applicable** | The trigger is absent and a reviewed reason is recorded. |

> Every applicable requirement has current evidence; no known risk exceeds the organization's tolerance; critical user journeys meet defined reliability and security objectives; and the organization can detect, contain, roll back, restore, support, and communicate failures.

## Templates and project links

- Record the review with the [release assessment](records/release-assessment.md), [evidence record](records/evidence-record.md), [risk exception](records/risk-exception.md), and [go/no-go decision](records/go-no-go-decision.md) templates.
- Review the [references, standards snapshot, and limitations](references.md).
- [Contribute a control, correction, documentation improvement, or scanner capability](https://github.com/MarinJursic/production-readiness-checklist/blob/main/CONTRIBUTING.md).
- [Star the repository on GitHub](https://github.com/MarinJursic/production-readiness-checklist) or share the documentation on [LinkedIn](https://www.linkedin.com/sharing/share-offsite/?url=https%3A%2F%2Fmarinjursic.github.io%2Fproduction-readiness-checklist%2F) or [X](https://twitter.com/intent/tweet?text=Production%20Readiness%20Checklist%3A%2010%2C042%20evidence-driven%20software%20engineering%20and%20release%20controls.&url=https%3A%2F%2Fmarinjursic.github.io%2Fproduction-readiness-checklist%2F).
