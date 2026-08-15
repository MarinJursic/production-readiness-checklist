# Evidence challenge

Use a second reviewer—or a fresh agent context—to challenge a completed assessment.

```text
Act as an independent reviewer of the attached production-readiness assessment. Follow CLAUDE.md, but do not repeat the original audit from scratch unless necessary.

For every claimed Pass or Not Applicable result:
- verify that the cited evidence exists;
- verify that it supports the exact PRC control rather than a nearby or weaker claim;
- verify that it covers the exact release, artifact, configuration, environment, migration state, and feature flags in scope;
- verify that it demonstrates behavior rather than intent;
- assess whether it is current enough;
- identify self-review, circular evidence, or untested assumptions;
- challenge Not Applicable rationales against product features, data, users, jurisdictions, dependencies, and risks.

Also identify:
- any immediate no-go condition hidden by wording, status, or exception;
- failures relabeled as accepted risk without an authorized owner and expiry;
- production or organizational controls incorrectly passed from repository evidence;
- material controls omitted from the assessment;
- changes made after evidence collection;
- contradictions between test results, configuration, documentation, and the assessment.

Return a table with Control ID, Original status, Challenged status, Reason, Evidence gap, and Required action. Do not calculate a score or make the final release decision.
```
