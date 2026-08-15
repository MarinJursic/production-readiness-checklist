# Getting started

You can get useful signal from this repository in 15 minutes, then deepen the review in proportion to the release risk.

## 1. Name the thing being approved

Record the source commit, immutable artifact digest, configuration version, database migrations, feature-flag state, target environment, and release window. If those change, affected evidence must be reviewed again.

Use the [release assessment template](../records/release-assessment.md) as the working record.

## 2. Run the no-go screen

Review all [immediate no-go conditions](../checklists/01-release-foundations.md#2-immediate-no-go-conditions). They are phrased as dangerous conditions: check one when the condition is true. Any checked no-go item stops the release until it is resolved or the release scope changes.

Do this before the full assessment. It prevents a long checklist from obscuring a single decisive failure.

## 3. Decide what applies

The ten core tracks cover most web applications. The conditional modules become mandatory whenever their trigger exists:

- payments, billing, subscriptions, or money movement;
- multi-tenant SaaS;
- user-generated content, communities, or marketplaces;
- email, SMS, push, or notifications;
- localization and internationalization;
- public content and SEO;
- PWA or offline operation;
- real-time, collaborative, or event-driven features;
- high-risk administrative tools;
- AI, machine learning, or LLM features;
- safety-critical or physically consequential behavior.

Record why an entire track or individual control is not applicable. Unexplained exclusions are gaps.

## 4. Assign owners by track

One coordinator should own the assessment, but evidence should come from the people accountable for each area: product, engineering, quality, security, reliability, data, privacy, accessibility, legal, support, and business risk.

Avoid assigning a team name alone. Each unresolved item needs a person who can make or escalate the decision.

## 5. Evaluate controls with evidence

Use these statuses consistently:

| Status | Meaning |
| --- | --- |
| Pass | Current evidence demonstrates the control for the exact release and environment |
| Fail | Evidence demonstrates that the control is not met |
| Blocked | Evidence is missing, inaccessible, stale, or cannot yet be collected |
| Not Applicable | The trigger is absent and the written rationale has been reviewed |

Use [evidence records](../records/evidence-record.md) for material controls. An implementation description is not proof that the deployed control works.

## 6. Handle exceptions explicitly

Do not relabel a failure as passed. Create a [risk exception](../records/risk-exception.md) with the business impact, accountable risk owner, compensating control, monitoring, remediation date, and automatic expiry.

Some controls are release blockers and should not be waived through an ordinary exception process.

## 7. Decide and retain

Assemble the final evidence package, collect the required sign-offs, and use the [go/no-go record](../records/go-no-go-decision.md). Store the assessment beside the release record so that it remains traceable after the next deployment.

## Suggested review depths

| Release profile | Suggested approach |
| --- | --- |
| Low-risk internal change | No-go screen, changed tracks, smoke tests, rollback evidence, named approver |
| Normal customer-facing release | All applicable core tracks, conditional modules, independent review, complete evidence package |
| High-impact or regulated launch | Full assessment, specialist review, rehearsals, external testing where warranted, formal sign-off |

Depth can vary; honesty about missing evidence cannot.
