# Scanner quick start

The scanner CLI is an experimental deterministic vertical slice. It inventories a
repository, creates an immutable plan for `prc/core-repository@0.1`, evaluates
native assertions, records evidence, and reports explicit unresolved states.

It does not yet run external analysis adapters or remediate findings. The
generated-code review remains manual, and the analysis-evidence assertion remains
blocked until an adapter supplies current executed evidence.

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
- discoverable source-level tests.

The engine rejects repository path and symlink escapes and detects a target file
that changes between inventory and evidence collection. Missing evidence,
unsupported applicability expressions, unavailable adapters, and manual review
are never converted into Pass.
