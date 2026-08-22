# Scanner product contract

The Production Readiness Scanner is an evidence and policy engine for a declared
target. It does not claim that a nontrivial system has no defects, certify an
organization, or make a release decision for an accountable owner.

## Verifiable promise

Given a declared target, release, readiness profile, evidence scope, policy, and
capability budget, the scanner will:

1. create an immutable plan of applicable and unresolved assertions;
2. collect observations through versioned, policy-authorized evidence producers;
3. evaluate each planned assertion without upgrading an inference into a fact;
4. preserve unknown, blocked, stale, conflicting, and manual-review results;
5. expose the evidence and reasoning behind every result;
6. remediate only findings whose fix contract and capability class permit it;
7. independently verify every candidate change; and
8. stop in an explicit terminal state.

The scanner may say that a profile is satisfied for a particular target and
evidence set. It must not reduce that scoped result to the unqualified statement
that a project is "production ready."

## Completion states

A scan or remediation run ends in exactly one of these states:

| State | Meaning |
| --- | --- |
| `profile_satisfied` | Every applicable assertion in the selected profile has current passing evidence and no gate blocks the result. |
| `machine_work_complete_manual_evidence_remaining` | Every policy-eligible automated fix is exhausted, but human or external evidence remains. |
| `assessment_incomplete` | Required evidence or applicability context is unavailable. |
| `no_go` | One or more configured blocking gates failed. |
| `policy_stopped` | A requested action exceeded the run's capability policy. |
| `budget_exhausted` | The attempt, time, cost, or change budget was reached. |
| `environment_blocked` | A required tool or evidence source could not be used safely. |

No aggregate percentage overrides these states. One no-go result can block a
release regardless of how many unrelated assertions pass.

## Truth ownership

The engine owns control selection, evidence verification, assertion evaluation,
gates, policy, and patch acceptance. Tools and connectors produce observations.
Coding agents produce untrusted patch candidates and explanations.

An agent response, model confidence, filename, README statement, configuration
file, or source-code comment is not proof of runtime or production behavior.
Evidence authority is determined by the assertion's proof obligation and the
environment and artifact to which the evidence applies.

## Non-goals

The scanner will not:

- prove the absence of every defect;
- infer production state from repository intent alone;
- treat provenance, an SBOM, a test file, or a policy file as proof that all
  associated behavior is correct;
- make legal, regulatory, accessibility, safety, risk-acceptance, or release
  decisions for accountable humans;
- let a target repository, adapter, or agent change scanner policy during a run;
- deploy from the general remediation loop; or
- weaken checks, thresholds, tests, baselines, or suppressions to manufacture a
  passing result.

## Compatibility promise

Catalog, profile, result, evidence, adapter, and provider formats are versioned
independently. Until a format is declared stable, consumers must pin a supported
version. Breaking changes require a migration path and explicit release notes.
