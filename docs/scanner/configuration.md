# Declared project configuration

Repository discovery cannot safely infer release scope, risk, deployed features,
data classes, or execution authority. `production-readiness.yaml` is the
versioned declaration for those facts and for the scanner's local capability
budget.

The current `prc.config/v0.1` contract is deliberately restrictive. Inventory,
plan, and scan commands accept it through `--config`. The scanner binds its
canonical digest and declarations into the inventory and plan identity.
Declaring a feature or environment is applicability context, not proof that the
feature is deployed or behaves correctly.

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
- selected scanner profile, exact source revision, artifact digests, and target
  environments;
- included component roots and reviewed exclusions with rationales;
- feature flags that discovery cannot establish safely;
- data classifications, regulatory-context tags, and evidence classes that must
  never be retained; and
- execution and remediation budgets.

Paths are repository-relative, normalized slash paths. Absolute paths,
backslashes, traversal, duplicate paths, and unsorted declarations are rejected.
Excluding a component does not delete it or cause the inventory walker to ignore
it in v0.1. A configuration inside the target must match the exact regular file
captured by inventory. An external configuration is also reopened and compared
by both raw-source and canonical digests. A change after validation fails the
run in either case.

```bash
prc inventory --target . --config production-readiness.yaml --format json
prc plan --target . --catalog-root . --config production-readiness.yaml
prc scan --target . --catalog-root . --config production-readiness.yaml
```

When `--profile` is omitted, plan and scan use the configured profile. An
explicitly selected profile must match it. Declared components are additive;
they do not erase discovered components or facts.

### Large local non-source directories

`components.exclude` is assessment scope only; it never hides files from the
content-hashed inventory. If the 8 GiB inventory guard is reached because the
project root contains a genuinely non-source directory such as generated demo
recordings, use the separate root `.prcignore` contract:

```text
# Exact relative directory, then a reviewed reason.
recordings | Local generated demo recordings are not project source.
```

This is deliberately not a normal glob ignore file. It accepts at most 100
exact canonical relative directories and no negation, wildcard, absolute path,
backslash, traversal, symlink, directory containing a symlink, or missing path. Every reason must be 10–300
characters. Before skipping anything, the scanner walks the directory names and
refuses the exclusion if it contains recognized source, manifests, lock files,
CI workflow YAML, Terraform, deployment YAML, documentation, security policy,
environment files, or other project-defining files. Accepted omissions appear
as `repository.user_exclusion` facts, are bound into the inventory digest, and
are called out in the terminal and HTML report. This guard reduces accidental
omissions; it is not proof that every skipped binary or data file is harmless.

### Remote AI context exclusions

`.prcreviewignore` is separate from `.prcignore`. It accepts exact
`relative/file | reviewed reason` entries only for source files that should
remain in every local check but must not be sent to an optional remote AI
review. This is useful for secret-scanner fixtures containing fake
credential-shaped strings. It cannot exclude directories or use globs, and
every accepted omission is named in the AI preview and sealed task
limitations. See [safe AI control review](ai-control-review.md#try-one-control-first).

`source_ref` must be empty or a lowercase 40–64 character hexadecimal revision.
When supplied, it must equal the Git revision inventoried from the target; an
unresolved variable, branch name, missing Git identity, or mismatch fails closed.
Git metadata is read only when its directory and references remain inside the
target root. A symlinked `.git` entry or a linked-worktree `gitdir` outside the
target is not followed, so that target has no inferred Git identity and cannot
satisfy a nonempty `source_ref`.

## Capability boundary

Configuration v0.1 accepts only deny-by-default capabilities and bounded resource
declarations:

- `network: deny`;
- an empty `allow_commands` list;
- `production_connected: false`;
- bounded parallelism and duration declarations;
- bounded remediation attempts, files, and changed lines; and
- a nonempty protected-path list.

A request for network access, target commands, or production connectivity fails
validation instead of being accepted and ignored. Future configuration versions
may add capabilities only together with an enforceable runner and threat-model
update.

The native scanner currently evaluates serially, so it remains below
`max_parallel`. `max_duration_seconds` is enforced for configured live adapter
execution and as the total `fix` loop deadline, including provider and verifier
child processes. Native file discovery uses fixed internal size and traversal
limits but is not yet interrupted by this setting. `remediate`,
`remediate-proposal`, and `provider seal-task` also accept `--config`.
Remediation must be enabled; configured duration, file, line, and attempt
ceilings cannot be raised by command-line flags; configured protected paths and
an in-target configuration source are always added to the immutable guard set.

## Parser and file safety

The loader accepts one regular YAML file no larger than 1 MiB. It rejects
symlinks, missing or unknown fields, duplicate keys, null values, YAML aliases,
multiple YAML documents, unsafe paths, unsupported enum values, capability
expansion, and out-of-range budgets. The JSON Schema and Go semantic validator
are both exercised in CI.
