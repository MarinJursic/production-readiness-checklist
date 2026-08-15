# Full production-readiness review

Copy the prompt below into Claude, Codex, or another coding agent with access to the application repository and this checklist.

```text
Perform a read-only, full-lifecycle engineering and production-readiness review of this repository.

Follow CLAUDE.md. Use docs/engineering/00-overview.md as the master sequence: review engineering phases 1–16 in order, then use docs/checklists/ for the final production gate. Also run the immediate no-go screen early so a decisive release blocker is not hidden by the larger review.

Release context:
- Release identifier: [fill in]
- Source commit or tag: [fill in]
- Artifact digest: [fill in or unknown]
- Target environment: [fill in]
- Configuration and feature-flag state: [fill in or reference]
- Critical user journeys: [fill in]
- Expected traffic and risk profile: [fill in]
- Conditional modules that apply: [fill in or ask the agent to assess]
- Lifecycle phases already covered by accepted evidence: [fill in or none]

Requirements:
- Inspect source, tests, dependency files, CI/CD, infrastructure definitions, migrations, runbooks, and documentation.
- Use both USEQ and PRC control IDs, and do not duplicate a finding when the same evidence satisfies related controls.
- Run only safe, read-only or local verification commands.
- Do not modify code or configuration.
- Use Pass, Fail, Blocked, or Not Applicable exactly as defined in CLAUDE.md.
- Cite precise evidence for every Pass.
- Never infer production or organizational evidence from source code.
- Put missing human or external evidence in a separate section.
- Do not produce a readiness percentage.
- Prioritize immediate no-go findings and controls that could cause security, privacy, data-integrity, financial, legal, safety, or unrecoverable operational harm.

Return the output structure required by CLAUDE.md and end with the smallest practical set of next actions needed to reach a defensible go/no-go decision.
```
