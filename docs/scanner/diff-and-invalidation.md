# Diff-aware evidence invalidation

`prc diff` compares a canonical prior run with a freshly inventoried target and
answers a narrower question than `scan`: which prior conclusions had relevant
inputs changed?

```bash
./prc diff \
  --state-dir /safe/local/path/prc-state \
  --base-run <run-id> \
  --target /path/to/project \
  --catalog-root /path/to/production-readiness-checklist
```

Use `--format json` for `prc.invalidation/v0.1`. The report includes changed
files and inventory dimensions, one decision for each prior or current
assertion, stable reason codes, and aggregate counts.

## What the decisions mean

| Conclusion | Meaning | Can the old result be returned as current? |
| --- | --- | --- |
| `invalidated` | A rule, applicability outcome, relevant repository input, declared scope, or evidence-freshness condition changed. | No. |
| `unchanged_inputs` | The implementation-specific inputs did not change. | Only when `reuse_allowed` is also true. |
| `new` | The current profile newly includes the assertion. | No. |
| `removed` | The current profile no longer includes the assertion. | No current evaluation is requested. |

The scanner deliberately distinguishes input equivalence from evidence reuse.
Repository evidence is bound to the complete inventory digest. If an unrelated
file changes, an assertion may be `unchanged_inputs` while `reuse_allowed` is
false and `fresh_evaluation_required` is true: its old evidence still names the
old target. `prc diff` does not manufacture a new evidence envelope or relabel a
prior Pass.

Reuse is permitted only when the current inventory and plan identities exactly
match the base run and no rule requires fresh non-repository evidence. Manual,
executed, artifact, and environment evidence is not carried forward without a
separate validity-window policy.

## Bound rule identity

Plans produced as `prc.plan/v0.6` bind:

- the scanner engine contract;
- the exact governing catalog, including every parsed objective, assertion, and profile;
- the exact profile definition;
- every assertion revision and complete definition, including parameters;
- the implementation identifier and applicability evaluator;
- project configuration, artifacts, and target environments;
- the selected execution mode, capability policy, implementation registry, and
  dependency DAG; and
- the complete inventory identity.

Each definition is represented by a SHA-256 digest over its canonical JSON
representation. A changed objective statement or assertion parameter therefore
invalidates the old conclusion even if its implementation name is unchanged.

Runs recorded with an older plan schema remain readable, but they lack one or
more rule bindings. Their assertions receive
`base_plan_lacks_rule_binding` and require a fresh evaluation. This is a
conservative migration rule, not a claim that the earlier result was wrong.

## Native dependency map

The versioned invalidation rules map each native implementation to the inputs it
actually reads. Examples include candidate document paths for `file-present`,
all workflow definitions for workflow parsers, source files for final-newline,
manifest and lock sets for dependency checks, container definitions for
container checks, and Terraform/Kubernetes resources for infrastructure checks.

Unknown implementations default to all repository files. That may perform more
fresh work than necessary, but it cannot preserve a stale Pass due to an
unrecognized dependency. Adapter and manual-evidence assertions also require
freshness proof instead of being inferred from a filename diff.

`prc diff` is an analysis and planning command. `prc scan` still evaluates the
current target afresh; automatic cache reuse is intentionally deferred until a
future evidence-rebinding protocol can preserve subject identity and validity
without weakening the trust model.
