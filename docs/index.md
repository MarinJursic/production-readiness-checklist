# Production Readiness Checklist

Ship with evidence, not optimism.

This project provides 1,421 technology-neutral controls for assessing whether a web application is ready for production. It is designed for real release work: focused review tracks, stable control IDs, reusable evidence templates, and honest support for AI-assisted audits.

## Choose your path

### I am preparing a release

Start with the [15-minute quick start](guides/getting-started.md), copy the [release assessment](records/release-assessment.md), then inspect the [immediate no-go conditions](checklists/01-release-foundations.md#2-immediate-no-go-conditions).

### I want an AI agent to review my repository

Read [AI-assisted review](guides/ai-assisted-review.md), keep `CLAUDE.md` at the repository root, and use one of the [review prompts](prompts/full-readiness-review.md).

### I am adopting this for a team

Begin with release foundations, assign track owners, define evidence freshness, and customize applicability before making any control a required gate. Do not use a percentage score as the go/no-go decision.

## Review flow

1. Identify the exact release and scope.
2. Stop on immediate no-go conditions.
3. Select applicable core and conditional tracks.
4. Collect current evidence for each applicable control.
5. Record failures, blocked items, and time-bounded exceptions.
6. Obtain independent review and required sign-offs.
7. Make and retain the go/no-go decision.
8. Verify the production deployment and continue observation.

## Core tracks

- [The readiness principle](checklists/00-readiness-principle.md)
- [Release foundations](checklists/01-release-foundations.md)
- [Product, risk, and architecture](checklists/02-product-risk-architecture.md)
- [Source, build, and supply chain](checklists/03-source-build-supply-chain.md)
- [Environments, quality, and experience](checklists/04-environments-quality-experience.md)
- [Application security](checklists/05-application-security.md)
- [Data, privacy, and performance](checklists/06-data-privacy-performance.md)
- [Reliability and operations](checklists/07-reliability-operations.md)
- [Maintenance, vendors, and compliance](checklists/08-maintenance-vendors-compliance.md)
- [Conditional feature modules](checklists/09-conditional-modules.md)
- [Evidence, sign-off, and decision](checklists/10-evidence-and-decision.md)

## Working principle

> Every applicable requirement has current evidence; no known risk exceeds the organization’s tolerance; critical user journeys meet defined reliability and security objectives; and the organization can detect, contain, roll back, restore, support, and communicate failures.

See [references and scope](references.md) for the standards snapshot and limitations.
