# Release-diff readiness review

Use this prompt when the application has a known-good baseline and you need to understand what a release changes or invalidates.

```text
Review the proposed release diff for engineering-quality and production-readiness impact. Follow CLAUDE.md, map affected lifecycle controls under docs/engineering/, and use docs/checklists/ for the final release gate.

Baseline release: [commit/tag/artifact]
Candidate release: [commit/tag/artifact]
Target environment: [environment]
Change summary: [summary or issue/PR links]
Critical journeys touched: [list]

Tasks:
1. Inspect the complete diff, including dependencies, lockfiles, migrations, infrastructure, workflows, configuration schemas, feature flags, tests, documentation, and generated artifacts.
2. Map each material change to affected USEQ and PRC control IDs.
3. Identify immediate no-go conditions introduced or made uncertain by the change.
4. Identify earlier evidence that is invalidated and must be collected again.
5. Identify compatibility risks during mixed-version deployment, rollback, roll-forward, queued work, cached data, and database migration.
6. Recommend the smallest sufficient regression, security, accessibility, performance, resilience, and post-deployment verification set.
7. Use Pass, Fail, Blocked, or Not Applicable. Cite exact evidence and state assumptions.
8. Do not modify files, run destructive commands, or make the final release decision.

Return:
- scope and assumptions;
- blockers;
- affected-control table with ID, impact, status, evidence, and required verification;
- invalidated-evidence list;
- recommended pre-deploy and post-deploy checks;
- human or external decisions needed.
```
