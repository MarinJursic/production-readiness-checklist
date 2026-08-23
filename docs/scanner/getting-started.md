# Scanner quick start

The scanner CLI is an experimental deterministic vertical slice. It inventories a
repository, creates an immutable plan for `prc/core-repository@0.3`, evaluates
native assertions, records evidence, and reports explicit unresolved states.

The profile evaluator can consume a live, sandboxed adapter execution only when
the selected profile binds its exact ID, manifest digest, and observation kind.
The generated-code review remains manual, and the default analysis-evidence
assertion remains blocked because no production adapter is authorized yet. A
separate experimental OCI adapter protocol and runner are available for protocol
and sandbox development.
One deterministic R1 remediation is available for recognized source files that
lack a final line-feed byte; it produces an isolated candidate instead of editing
the target.

## Build

Go 1.27 or a compatible later supported toolchain is required.

```bash
go mod verify
go build -trimpath -o prc ./cmd/prc
```

## Inventory without executing target code

```bash
./prc inventory --target /path/to/project
```

The inventory walks regular files without following symlinks or executing project
commands. Inventory v0.3 records content digests plus a deterministic component
graph and sourced, confidence-bearing facts for package manifests, CI workflows,
container build definitions, Terraform, Kubernetes, and symlinks. Detection is
not proof that a component is built or deployed; each fact states that limitation.
When a validated project configuration is supplied, its digest, declared scope,
and explicitly limited facts are bound into the inventory identity. Frozen v0.1
and v0.2 schemas remain available for consumers migrating pinned output contracts.

## Create a deterministic plan

Run this command from the Production Readiness Checklist repository, or pass its
location through `--catalog-root`:

```bash
./prc plan \
  --catalog-root /path/to/production-readiness-checklist \
  --target /path/to/project \
  --config /path/to/project/production-readiness.yaml \
  --profile prc/core-repository
```

The same target inventory, profile, and catalog produce the same plan digest.
Plan v0.3 records the exact bounded CEL evaluator, an applicability reason for
every assertion, the configuration digest, declared project ID, artifact
digests, and target environments. Invalid, unavailable, non-Boolean, or
resource-exhausting expressions become `undetermined`; they never silently
become Not Applicable. See
[bounded applicability evaluation](../architecture/applicability.md) for the
available inventory fields and limits.

## Scan and preserve evidence

```bash
./prc scan \
  --catalog-root /path/to/production-readiness-checklist \
  --target /path/to/project \
  --profile prc/core-repository \
  --state-dir /safe/path/prc-state
```

`--state-dir` is optional. When supplied, evidence envelopes are stored by digest
and the complete run is written atomically. Do not put the state directory inside
the target unless target-local scanner state is intentional.

JSON is available for automation:

```bash
./prc scan --target . --catalog-root . --format json --exit-policy never
```

The same scan can emit a scoped Markdown report, self-contained accessible HTML,
SARIF 2.1.0 failed findings, or JUnit XML with failures, execution errors, and
manual/not-applicable skips kept distinct:

```bash
./prc scan --target . --catalog-root . --format markdown --exit-policy never
./prc scan --target . --catalog-root . --format html --exit-policy never
./prc scan --target . --catalog-root . --format sarif --exit-policy never
./prc scan --target . --catalog-root . --format junit --exit-policy never
```

SARIF intentionally contains only failed assertions that can be represented as
findings. Unknown, blocked, manual, and not-applicable results remain in the
canonical JSON, Markdown, HTML, and JUnit reports instead of being mislabeled as
source findings.

`--exit-policy profile` is the default and exits nonzero unless the profile is
fully satisfied. `no-go` exits nonzero only for a no-go terminal state. `never`
always returns success after a completed scan but does not change the result in
the report.

## Explain an assertion

```bash
./prc explain --catalog-root . PRC-A-CORE-008
```

## Current native checks

- repository governance files and an immutable Git source identity;
- detected dependency manifests, nonempty lock inputs, update configuration,
  and runtime-version declarations;
- GitHub Actions syntax, jobs, bounded runtimes, permissions, immutable action
  references, and unsafe untrusted triggers;
- discoverable tests, final line-feed bytes, merge-conflict markers, and broad
  file-write permissions;
- immutable container base identities and final-stage non-root users;
- Terraform provider locks; and
- Kubernetes non-root and container resource policies.

The profile has 30 atomic assertions. Container, Terraform, and Kubernetes
assertions are planned as not applicable when their corresponding sourced
inventory facts are absent.

The engine rejects repository path and symlink escapes and detects a target file
that changes between inventory and evidence collection. Missing evidence,
invalid or unevaluable applicability expressions, unavailable adapters, and
manual review are never converted into Pass.

Read the [sandboxed adapter protocol](../architecture/adapters.md) to validate an
adapter transcript, inspect an OCI execution plan, run an already-present
digest-pinned adapter image, or understand live profile-authorized evidence
consumption.

Read [bounded isolated remediation](remediation.md) to create and verify an R1
candidate or to apply one validated R2 provider proposal outside the source
workspace.

Read [read-only agent providers](../architecture/agent-providers.md) to inspect
or explicitly run the experimental Codex and Claude Code `suggest` adapters.
Provider proposals are never applied or treated as evidence of a passing check.
