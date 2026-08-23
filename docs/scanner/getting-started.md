# Scanner quick start

The scanner CLI is an experimental deterministic vertical slice. It inventories a
repository, creates an immutable plan for `prc/core-repository@0.8`, evaluates
native assertions, records evidence, and reports explicit unresolved states.

The profile evaluator can consume a live, sandboxed adapter execution only when
the selected profile binds its exact ID, manifest digest, and observation kind.
The generated-code review remains manual, and the default analysis-evidence
assertion remains blocked because no production adapter is authorized yet. A
separate experimental OCI adapter protocol and runner are available for protocol
and sandbox development.
One deterministic R1 remediation is available for recognized source files that
lack a final line-feed byte, and another restricts broadly writable file modes.
The bounded `fix` loop can compose those repairs in isolated sibling candidates
without editing the target.

## Build

Go 1.27 or a compatible later supported toolchain is required.

```bash
go mod verify
go build -trimpath -o prc ./cmd/prc
```

`prc version --format json` reports the scanner version, exact source revision,
reproducible source timestamp, and Go toolchain embedded in a release build.
Downloadable scanner releases use a separate `scanner-vX.Y.Z` tag namespace and
bundle the compatible catalog, packs, and schemas with every binary. See
[releases and verification](releases.md) before trusting a downloaded artifact.

## Inventory without executing target code

```bash
./prc inventory --target /path/to/project
```

The inventory walks regular files without following symlinks or executing project
commands. Inventory v0.3 records content digests plus a deterministic component
graph and sourced, confidence-bearing facts for package manifests, CI workflows,
container build definitions, Terraform, Kubernetes, OpenAPI descriptions, and
symlinks. Detection is
not proof that a component is built or deployed; each fact states that limitation.
When a validated project configuration is supplied, its digest, declared scope,
and explicitly limited facts are bound into the inventory identity. Versioned
v0.1, v0.2, and v0.3 schema files remain available for consumers of pinned
output contracts; only the unversioned alias advances with a later contract.

Before enabling persistent evidence, OCI adapters, or isolated remediation, run
[`prc doctor`](doctor.md) with the corresponding paths and executables. It probes
only explicitly requested host capabilities and never launches target code,
containers, or agent providers.

## Create a deterministic plan

Optionally validate or export the exact catalog first:

```bash
./prc catalog validate --catalog-root /path/to/production-readiness-checklist
./prc catalog bundle --catalog-root /path/to/production-readiness-checklist > catalog-bundle.json
```

The manifest and bundle are deterministic and contain no local path or timestamp.
They identify validated definitions; release signing is a separate trust step.

Run this command from the Production Readiness Checklist repository, or pass its
location through `--catalog-root`:

```bash
./prc plan \
  --catalog-root /path/to/production-readiness-checklist \
  --target /path/to/project \
  --config /path/to/project/production-readiness.yaml \
  --profile prc/core-repository
```

The same target inventory, profile, catalog, and execution mode produce the same
plan digest. Plan v0.6 binds the full governing catalog, exact profile and
assertion definitions, bounded CEL evaluator, an applicability reason for every
assertion, configuration digest, declared project ID, artifact digests, target
environments, implementation registry, capability envelope, and dependency
DAG. `inspect` is the default and denies external processes, writable scratch,
network, and secrets. Use `--mode verify-local` only when reviewing a plan that
contains an authorized no-network OCI adapter. Invalid, unavailable,
non-Boolean, or resource-exhausting expressions become `undetermined`; they
never silently become Not Applicable. See
[bounded applicability evaluation](../architecture/applicability.md) for the
available inventory fields and limits.

## Scan and preserve evidence

```bash
mkdir -m 0700 /safe/path/prc-state

./prc scan \
  --catalog-root /path/to/production-readiness-checklist \
  --target /path/to/project \
  --profile prc/core-repository \
  --state-dir /safe/path/prc-state
```

`--state-dir` is optional. When supplied, evidence envelopes and complete runs
are stored by digest, then indexed transactionally in an embedded SQLite
database. The state root must be private and local; do not put it inside the
target unless target-local scanner state is intentional. See
[durable state and history](state-and-history.md) for the storage contract and
`history` commands.

After a later change, use [`prc diff`](diff-and-invalidation.md) with a canonical
base run to see exactly which rule inputs were invalidated before rescanning.

JSON is available for automation:

```bash
./prc scan --target . --catalog-root . --format json --exit-policy never
```

To let Codex, Claude Code, or another local agent request the same read-only
plans, scans, and assertion explanations, use the
[path-locked MCP server](mcp-agent-integration.md). It deliberately exposes no
write, remediation, provider, adapter, process, network, or state tool.

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
canonical, fingerprinted findings. Unknown, blocked, manual, and not-applicable
results remain in the canonical JSON, Markdown, HTML, and JUnit reports instead
of being mislabeled as source findings. Run v0.9 adds bounded source locations
to assertion results so native analyses can feed exact file, line, and column
data into canonical findings. It keeps a stable finding fingerprint separate
from the content-addressed finding ID and embeds the reviewable v0.6 execution
plan. Adapter executions bind their local or registry authorization provenance.
Versioned run v0.1 through v0.8, plan v0.1 through v0.6, inventory v0.1
through v0.3, adapter-execution v0.1 through v0.2, evidence v0.1, and finding
v0.1 schemas remain available for archived consumers. Their local dependency
graphs use versioned filenames, so later unversioned aliases cannot change an
archived contract.

`--exit-policy profile` is the default and uses the [stable CLI exit-code
contract](cli-contract.md): a failed active gate is `1`, while incomplete,
blocked, or manual evidence is `2`. `no-go` never hides incomplete execution.
`never` is an explicit report-generation override that returns `0` after a
completed scan but does not change the terminal state in the report; do not use
it as a release gate.

## Fix eligible deterministic findings

```bash
./prc fix \
  --target /path/to/project \
  --catalog-root /path/to/production-readiness-checklist \
  --candidate-root /safe/path/prc-remediation-run \
  --format json
```

The command creates a new candidate per accepted R1 fix and rescans from fresh
evidence after every attempt. It returns success only when the selected profile
is satisfied. Read [bounded isolated remediation](remediation.md) for budgets,
configuration policy, terminal states, and the machine-work-complete report.

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
- discoverable tests, final line-feed bytes, merge-conflict markers, broad
  file-write permissions, and committed private-key armor;
- direct Go `net/http` package helpers backed by mutable global
  `http.DefaultClient` state;
- immutable container base identities and final-stage non-root users;
- Terraform provider locks; and
- Kubernetes non-root and container resource policies.

The profile has 34 atomic assertions. Go, container, Terraform, Kubernetes, and
OpenAPI assertions are planned as not applicable when their corresponding
sourced inventory facts are absent.

The engine rejects repository path and symlink escapes and detects a target file
that changes between inventory and evidence collection. Missing evidence,
invalid or unevaluable applicability expressions, unavailable adapters, and
manual review are never converted into Pass.

Read the [sandboxed adapter protocol](../architecture/adapters.md) to validate an
adapter transcript, inspect an OCI execution plan, run an already-present
digest-pinned adapter image, or understand live profile-authorized evidence
consumption.

Read [bounded isolated remediation](remediation.md) to run the deterministic R1
loop, create and verify one R1 candidate, or apply one validated R2 provider
proposal outside the source workspace.

Read [read-only agent providers](../architecture/agent-providers.md) to inspect
or explicitly run the experimental Codex and Claude Code `suggest` adapters.
Provider proposals are never applied or treated as evidence of a passing check.
