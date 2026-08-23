# Declared project configuration

Repository discovery cannot safely infer release scope, risk, deployed features,
data classes, or execution authority. `production-readiness.yaml` is the
versioned declaration for those facts and for the scanner's local capability
budget.

The current `prc.config/v0.1` contract is deliberately restrictive. It supports
configuration validation and content identity; scan-plan consumption is not yet
enabled. Declaring a feature or environment is context for future applicability,
not proof that the feature is deployed or behaves correctly.

## Validate a configuration

Start from the checked-in
[`fixtures/config/production-readiness.yaml`](https://github.com/MarinJursic/production-readiness-checklist/blob/main/fixtures/config/production-readiness.yaml)
example, then run:

```bash
prc config validate --file production-readiness.yaml
prc config validate --file production-readiness.yaml --format json
```

The JSON output includes a SHA-256 digest over the canonical decoded document.
Mapping-key order does not change the digest. Security-relevant lists must be
sorted and duplicate-free so equivalent inputs have one representation.

## Scope declarations

The document records:

- stable project identity and risk profile;
- selected scanner profile, source reference, artifact digests, and target
  environments;
- included component roots and reviewed exclusions with rationales;
- feature flags that discovery cannot establish safely;
- data classifications, regulatory-context tags, and evidence classes that must
  never be retained; and
- execution and remediation budgets.

Paths are repository-relative, normalized slash paths. Absolute paths,
backslashes, traversal, duplicate paths, and unsorted declarations are rejected.
Excluding a component does not delete it or cause the inventory walker to ignore
it in v0.1.

## Capability boundary

Configuration v0.1 accepts only the capabilities the current engine can enforce:

- `network: deny`;
- an empty `allow_commands` list;
- `production_connected: false`;
- bounded parallelism and duration;
- bounded remediation attempts, files, and changed lines; and
- a nonempty protected-path list.

A request for network access, target commands, or production connectivity fails
validation instead of being accepted and ignored. Future configuration versions
may add capabilities only together with an enforceable runner and threat-model
update.

## Parser and file safety

The loader accepts one regular YAML file no larger than 1 MiB. It rejects
symlinks, unknown fields, duplicate keys, multiple YAML documents, unsafe paths,
unsupported enum values, capability expansion, and out-of-range budgets. The
JSON Schema and Go semantic validator are both exercised in CI.
