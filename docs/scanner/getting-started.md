# Scanner quick start

The scanner CLI is an experimental deterministic vertical slice. It inventories a
repository, creates an immutable plan for `prc/core-repository@0.2`, evaluates
native assertions, records evidence, and reports explicit unresolved states.

The profile evaluator does not yet consume external adapter observations. The
generated-code review remains manual, and the
analysis-evidence assertion remains blocked until an adapter supplies current
executed evidence that the engine can evaluate. A separate experimental OCI
adapter protocol and runner are available for protocol and sandbox development.
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
commands. File content digests, detected manifests, package ecosystems, lock
files, source files, and CI providers contribute to the inventory.

## Create a deterministic plan

Run this command from the Production Readiness Checklist repository, or pass its
location through `--catalog-root`:

```bash
./prc plan \
  --catalog-root /path/to/production-readiness-checklist \
  --target /path/to/project \
  --profile prc/core-repository
```

The same target inventory, profile, and catalog produce the same plan digest.

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

- required repository governance and documentation files;
- detected dependency ecosystems and corresponding lock or checksum files;
- GitHub Actions presence and full-commit action pins;
- explicit GitHub Actions permissions; and
- discoverable source-level tests;
- final line-feed bytes on recognized source files.

The engine rejects repository path and symlink escapes and detects a target file
that changes between inventory and evidence collection. Missing evidence,
unsupported applicability expressions, unavailable adapters, and manual review
are never converted into Pass.

Read the [sandboxed adapter protocol](../architecture/adapters.md) to validate an
adapter transcript, inspect an OCI execution plan, or run an already-present
digest-pinned adapter image outside profile evaluation.

Read [bounded isolated remediation](remediation.md) to create and verify an R1
candidate or to apply one validated R2 provider proposal outside the source
workspace.

Read [read-only agent providers](../architecture/agent-providers.md) to inspect
or explicitly run the experimental Codex and Claude Code `suggest` adapters.
Provider proposals are never applied or treated as evidence of a passing check.
