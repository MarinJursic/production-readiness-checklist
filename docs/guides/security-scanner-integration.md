# Security scanner integration

Several controls require that an automated analysis *has run*: source-code security analysis, dependency and SBOM analysis, secret scanning, and infrastructure and configuration analysis. A scanner wired into automation carelessly can satisfy the wording of those controls while producing evidence that does not support a Pass. This page describes the properties an automated scanning integration must have before its output counts as evidence, without recommending any particular scanner.

The checklist is technology-neutral, and so is this guidance. Any analyzer that can demonstrate the properties below is acceptable; a well-known one that cannot demonstrate them is not.

## The evidence problem

A clean scan result has two possible causes: the analyzer ran and found nothing, or the analyzer did not meaningfully run. Job logs usually render both as a green check. Treating the absence of a finding as proof of security is already an anti-pattern in [AI-assisted review](ai-assisted-review.md), and it applies with equal force to automated tooling.

An integration is only useful as evidence if a reviewer can tell those two cases apart months later, without rerunning it.

## Required properties

**Fail loudly rather than cleanly.** An install failure, crash, timeout, authentication error, or unreadable target must fail the job. A scanner that reports success when it never executed is worse than no scanner, because it converts an unknown into a false Pass.

**Prove coverage.** Record how many files were analyzed and which were skipped or unreadable. A clean result over a partial, empty, or wrongly-rooted file set is not coverage. Coverage counts are also what let a reviewer judge whether the run actually addressed the release under assessment.

**Publish to a durable, reviewable destination.** Findings should land in a standard machine-readable format—SARIF is one open example—in a store that outlives job-log retention and that reviewers already consult. Findings that exist only in expired console output cannot satisfy the requirement that security evidence is retained with the release.

**Report first, enforce later.** A newly introduced analyzer should not block merges until its output has been tuned against the codebase. Enforcing an untuned scanner teaches contributors to bypass it, and a bypass habit is harder to reverse than a noisy backlog.

**Give findings owners and disposition.** Every finding needs an owner and a remediation target, and every dismissed finding needs a documented rationale. An unreviewed backlog is not evidence of anything, and suppressing an entire rule class to clear one is a decision that must be recorded rather than absorbed into configuration.

**Measure the analyzer before adopting it.** Run a candidate against the actual repository and count true and false positives before wiring it into shared automation. An analyzer that produces no true positives on the codebase adds maintenance surface and reviewer fatigue without adding assurance, and one that produces mostly false positives will be ignored within a release cycle.

## The pipeline is attack surface

Scanning automation runs with repository privileges against code that a contributor proposed. It deserves the same scrutiny as any other privileged component.

- Analysis of an outside contribution must not run in a trigger context that exposes repository credentials or secrets to contributor-controlled code.
- The job should hold the narrowest privilege that lets it publish results, and no write access it does not use. Adding scopes so that automation can post commentary expands the blast radius for a convenience.
- Untrusted values—branch names, change titles, author fields, file paths—must never be interpolated into a shell command or script body.
- The analyzer and any actions it depends on should be pinned to immutable versions or digests, verified for integrity before installation, and installed without executing package lifecycle scripts.
- The resulting automation identity belongs in the inventory of bots, apps, and credentials that the repository already maintains.

A scanner introduced to reduce risk should not be the most privileged and least reviewed job in the pipeline.

## Choosing the scan target

Scanning only the change under review is fast and keeps findings actionable, but it will not see a weakness introduced before the integration existed. Scanning the whole tree establishes a baseline and finds inherited problems, at the cost of an initial backlog that must be triaged rather than ignored.

Secret detection is the case where history matters: a credential removed in a later commit remains valid and remains reachable in earlier objects, so current-state scanning alone leaves it exposed. Rotation, not deletion, is what resolves an exposed secret.

## Controls this supports

| Control | What the integration must show |
| --- | --- |
| `PRC-31-003` | Source-code analysis executed against the assessed release, with coverage recorded. |
| `PRC-31-004` | Dependency and SBOM analysis executed, with the component set it covered. |
| `PRC-31-005`, `PRC-07-014` | Secret detection executed across current changes and repository history. |
| `PRC-31-006` | Infrastructure and configuration analysis executed against the definitions in scope. |
| `PRC-31-019` | Dismissed findings carry a recorded rationale rather than a silent suppression. |
| `PRC-31-020`, `PRC-31-021` | Findings carry owners and remediation targets, and fixes were retested. |
| `PRC-31-030` | Results are retained with the release rather than expiring with the job log. |
| `PRC-07-009` | Required automated checks are actually required to merge, once tuned. |
| `PRC-07-020` | The scanning automation identity is inventoried and restricted. |
| `PRC-09-016` | Released components continue to be scanned as new advisories appear. |

Verify each mapping against the current text in [source, build, and supply chain](../checklists/03-source-build-supply-chain.md) and [reliability and operations](../checklists/07-reliability-operations.md); control identifiers are stable, but wording may be refined.

## Anti-patterns

- Treating a green check as proof that the analysis ran.
- Recording a Pass because a scanner reported no findings, without confirming what it covered.
- Adopting an analyzer without first measuring its true and false positive rate on the repository.
- Adding a scanner whose findings no one is accountable for triaging.
- Granting the scan job write privileges or secret access it does not need to publish results.
- Pinning an analyzer to a mutable tag, or installing it with lifecycle scripts enabled.
- Blocking merges on an untuned analyzer, then normalizing bypass.
- Silencing a rule class to clear a backlog instead of documenting each dismissal.
- Citing one analyzer's clean result as evidence for controls it never evaluated.
- Reporting a count of findings or a percentage in place of blocker analysis.

## Scope note

Automated analysis supports the verification controls; it does not replace manual review of the highest-risk components, independent testing where warranted, or the accountable human decision described in [evidence, sign-off, and decision](../checklists/10-evidence-and-decision.md).
