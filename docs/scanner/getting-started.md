# Scanner quick start

The scanner CLI inventories a repository, creates an immutable plan, evaluates
native assertions, records evidence, and reports explicit unresolved states.
It also binds all 10,042 registered controls into every complete scan report.
The narrow assertions do not overclaim that they proved a whole broad control.

After a global install, the normal start is just:

```bash
cd /path/to/project
prc setup
prc
```

`prc setup` is an optional one-time local preflight; it does not run project
code, containers, or AI. Bare `prc` runs the 40 core local checks. Use `prc /path/to/project` to scan a
different directory, `prc quick` for an 18-check screen, or `prc full codex`
for the core scan plus advisory AI review of all 9,356 reviewed nondeterministic
controls. `prc full claude` selects Claude Code instead. The commands do not fix
files or execute project code. They print a summary and create a private,
standalone HTML report outside the target. Click the exact path printed as
`Detailed report:` in a supported terminal, or open the same plain path normally,
to review verified findings and every incomplete or manual result. `quick` still includes every
control in that report; it only runs fewer local assertions.

The profile evaluator can consume a live, sandboxed adapter execution only when
the selected profile binds its exact ID, manifest digest, and observation kind.
The generated-code review remains manual. The default analysis-evidence
assertion pins the reviewed Gitleaks 8.30.0 current-tree adapter, but remains
blocked in ordinary inspect mode. After the exact immutable image has been
reviewed and pulled once, `prc verify [project]` is the short explicit command.
It selects `verify-local` and the bundled manifest, passes `--pull=never`, turns
off container networking, and scans a sealed read-only snapshot. The
[adapter contract](../architecture/adapters.md#inspect-and-validate) documents
the exact image command and important coverage limitations.
One deterministic R1 remediation is available for recognized source files that
lack a final line-feed byte, and another restricts broadly writable file modes.
The bounded `fix` loop can compose those repairs in isolated sibling candidates
without editing the target.

## One-time setup

### npm release package

For a one-time global installation and a short command in every project:

```bash
npm install -g @marinjursic/prc
cd /path/to/project
prc setup
prc
```

Use `prc /path/to/project` to scan another folder. npm stores this tool
under its global prefix and exposes `prc` on `PATH`; it does not add a
`node_modules` directory or package entry to the scanned project. Do not use
`sudo` to work around a global-install permission error; install Node with a
version manager instead.

The package is one user-facing install. npm quietly chooses the one native
binary that matches this computer instead of downloading all six builds. There
is no `npx` prefix after installation. Check, update, or remove it with:

```bash
prc version
prc update
npm install -g @marinjursic/prc@latest
npm uninstall -g @marinjursic/prc
```

The startup screen prints the exact project path. If the inventory reaches its
8 GiB safety limit, check that path first and run inside the intended project
root, or use `prc /exact/project/path`. The error identifies the next file
that would cross the limit. Clear generated cache data when appropriate; do not
delete real project data or raise the guard just to force a result.

For a large directory that is genuinely non-source data, a root `.prcignore`
can name an exact relative directory and a reviewed reason:

```text
recordings | Local generated demo recordings are not project source.
```

The scanner refuses traversal, symlinks, missing paths, wildcards, and any
excluded directory containing recognized source, configuration, deployment,
CI, documentation, environment, or security-policy files. Every accepted
omission is visible in the inventory identity and report. See
[project configuration](configuration.md#large-local-non-source-directories).

The npm launcher has no install hooks or third-party JavaScript dependencies.
It uses one exact native platform package and never downloads a fallback or
updates itself in the background. Every launch verifies the binary and bundled
runtime files against the release-bound manifest. The package retains the
complete 10,042-control catalog and scanner evidence, while excluding website
media and contributor-only documentation. The global tool is convenient for
one developer machine, but it is not pinned by a project lock file.

For a pinned install with dependency hooks disabled, use
`npm install -D -E --ignore-scripts --no-audit --no-fund @marinjursic/prc@X.Y.Z`,
then `npm exec --offline --no -- prc scan`. To give a Node project one short
repeatable command, add `"scan": "prc"` to `package.json`, then run
`npm run scan`. Use `npm run --ignore-scripts scan` to skip local `prescan` and
`postscan` hooks. npm does not let a project add a custom top-level `npm scan`
command.

### Build from source

Go 1.27 or a compatible later supported toolchain is required.

```bash
go mod verify
go build -trimpath -o prc ./cmd/prc
./prc version
```

Run these commands from the Production Readiness Checklist repository. The
binary automatically discovers the compatible `catalog/` in the current
directory or beside the executable. Keep release-archive files together when
moving the binary. If `go` is not found, install the supported Go toolchain and
open a new shell so its binary directory is on `PATH`.

Add the entire extracted or cloned scanner directory to `PATH` when you want to
run `prc` without `./`; moving only the executable separates it from its
compatible catalog. For a current macOS or Linux shell:

```bash
export PATH="/absolute/path/to/production-readiness-checklist:$PATH"
```

`prc version --format json` reports the scanner version, exact source revision,
reproducible source timestamp, and Go toolchain embedded in a release build.
Downloadable scanner releases use a separate `scanner-vX.Y.Z` tag namespace and
bundle the compatible catalog, adapter manifests, packs, schemas, and scanner
guides with every binary. See
[releases and verification](releases.md) before trusting a downloaded artifact.

The native CLI stays language-neutral and remains available without Node.

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

The scanner normally discovers its bundled catalog. Pass `--catalog-root` only
when intentionally testing a different local catalog:

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

## Scan and read the report

```bash
./prc scan /path/to/project
```

Options may appear before or after the project path. The default report is
created with mode `0600` under the operating system's user cache directory and
is never written inside the target. Its name includes the target and run ID.
It contains:

- all 10,042 registered controls and an honest disposition for each one;
- a plain-language split between verified problems, narrow passes, unresolved
  local checks, manual decisions, controls needing evidence, and AI advice;
- every verified finding with severity, gate, controls, locations, evidence,
  remediation class, finding ID, and stable fingerprint; and
- every assertion result with applicability, execution state, required
  evidence, observed evidence, and locations.

When optional features are used, the same report also retains dual-signature
verification for imported authoritative evidence and a compact, scanner-owned
AI improvement plan grouped only by exact cause matches.

The scanner creates report files with exclusive creation and will not overwrite
an existing path. Use `--report /safe/path/readiness.html` to choose a new path,
or `--no-report` to explicitly suppress the default HTML report.

The large circle is the pass rate for applicable **local checks**, not a claim
that the full project is production-ready. Categories with only one or two
applicable checks say **Limited evidence** even when those checks pass. The full
catalog is kept as compact inert data and the browser renders only 25 matching
controls at a time. Open the newest report and inspect scanner-owned disk use
from any directory with:

```bash
prc report
prc report list
prc cache status
```

Cache deletion always needs an explicit class, for example
`prc cache clean --reports --older-than 720h`.

`prc scan` has no code path to the remediation commands. A missing final newline
or broad file mode may appear as a finding, but the target bytes and modes remain
unchanged. `prc fix` is a separate, explicitly invoked candidate workflow.

## Optional signed deterministic evidence

Trusted collector teams can attach one previously created evidence bundle
without loading collector code into `prc`:

```bash
prc scan /path/to/project \
  --evidence-bundle evidence.json \
  --evidence-trust-store trust-store.json \
  --evidence-policy-signature policy-signature.json \
  --evidence-signature authority-evidence-signature.json
```

All four files are required. Two different trusted Ed25519 keys sign ordered
subjects: one policy key signs the program and runtime-input digest before
collection, and one key limited to the evidence authority signs the completed
bundle's canonical digest after observation.
The bundle must match the current catalog and the exact inventory built by this
scan. It contains typed facts, not executable code or a provider verdict. See
[signed authoritative evidence bundles](authoritative-evidence-bundles.md) for
the producer contract and every fail-closed check.

When evidence comes from several authorities, place the trust store, bundles,
and signatures in one private directory and use the one-file form:

```bash
prc scan /path/to/project --evidence-set /private/evidence/evidence-set.json
```

Use `prc coverage` to see the exact difference between rules that have reviewed
predicates, clauses the scanner can collect by itself, and clauses supported by
the signed external import route.

To inspect what a missing collector must produce, filter the authenticated
producer contract instead of reading the generated catalog by hand:

```bash
prc evidence requirements
prc evidence requirements --authority repository --collector-status missing
prc evidence requirements --control PRC-36-004 --format json
```

Human output is intended for review. JSON output conforms to
`schemas/evidence-requirements.schema.json` and includes typed raw facts,
pre-sealed inputs, source, inventory, normalization, completeness, and
freshness requirements. It never counts the contract itself as observed
evidence.

Before attaching a multi-authority set to a full scan, verify it against the
same project inventory and catalog:

```bash
prc evidence verify-set \
  --set /private/evidence/evidence-set.json \
  /path/to/project
```

Use `--format json` for a schema-checked verification record. A successful
preflight confirms the signatures and internal bindings. It does not certify
the truth of producer observations or production readiness, and any failing or
blocked predicate stays failing or blocked.

## Optional Codex or Claude review

A normal scan does not contact an AI provider. Use the provider's official CLI
sign-in flow once, then start the review with one short option:

```bash
prc login codex
prc full codex --plan
prc full codex
```

Replace `codex` with `claude` for Claude Code. `prc auth` shows whether each
private scanner login is ready, and `prc logout codex` or `prc logout claude`
removes it. A supported API-key environment variable remains an alternative.

The `--plan` run performs source screening and exact batching without resolving
or starting the provider and without creating resume data. It shows the source
bytes, omissions, controls, batches, workers, timeout, 1,500-batch default
ceiling, and 24-hour default total deadline.

The short full command reviews all 9,356 reviewed nondeterministic controls.
The 686 reviewed deterministic controls are never decided by AI. A supported
exact program can produce a verified result from sealed authoritative evidence;
all remaining controls stay Blocked until their required collector and evidence
are available. The first shipped exact collector recognizes a root Node
`package.json` with usable `build` and `test` scripts and proves that both public
commands appear in inventoried Markdown code. It does not run either command.
The sealed task asks the coordinator to assign every nondeterministic control to
a separate primary subagent inside a sealed batch. Deep mode also asks for one
independent skeptical subagent for the batch, then keeps the strongest objection
in the structured result instead of hiding disagreement. Current provider
output cannot independently prove that each internal subagent actually ran, so
the scanner never counts that orchestration as evidence. Four batches run in parallel, and Codex
uses `xhigh` reasoning. Completed batches resume from private state outside the
target. This full run can take a long time and use many tokens. The provider
receives bounded, secret-screened excerpts but no target
workspace path or source-reading, shell, write, install, web, browser, or MCP
tools. Its result stays advisory and never turns a control into a verified
Pass. A cited path and line are checked against the exact screened snapshot,
but the claim remains `advisory_unverified` because a real line does not prove
the AI understood it. `prc full` and `--ai` are also explicit permission to send
those screened excerpts to the remote provider. Read [safe AI control review](ai-control-review.md) for
one-control testing, cost warnings, result meanings, and stop conditions.

The terminal shows the review plan before the first remote call, then bounded
batch and control progress with elapsed time. New Codex batches show provider-
reported token totals when available. New Claude batches show the provider's
client-side cost estimate, plus any enforced per-batch limit you selected.
Cached batches are clearly counted but their old usage is not guessed.
If a later batch fails, the scanner still writes a clearly marked partial
report, returns exit code `4`, and keeps the valid completed batches. Repeating
the same command checks and reuses those batches before continuing.

## Preserve evidence and history

```bash
mkdir -m 0700 /safe/path/prc-state

./prc scan \
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
./prc scan . --format json --exit-policy never
```

Explicit machine formats write to stdout and do not also create the default HTML
file. Add `--report /safe/path/readiness.html` when both forms are wanted.

The short local-only CI preset writes SARIF 2.1.0 and creates no HTML report:

```bash
prc ci > prc-results.sarif
```

It keeps the normal readiness exit code so a failed or incomplete gate is not
silently converted to success.

To let Codex, Claude Code, or another local agent request the same read-only
plans, scans, and assertion explanations, use the
[path-locked MCP server](mcp-agent-integration.md). It deliberately exposes no
write, remediation, provider, adapter, process, network, or state tool.

The same scan can emit a scoped Markdown report, self-contained accessible HTML,
SARIF 2.1.0 failed findings, or JUnit XML with failures, execution errors, and
manual/not-applicable skips kept distinct:

```bash
./prc scan . --format markdown --exit-policy never
./prc scan . --format html --exit-policy never
./prc scan . --format sarif --exit-policy never
./prc scan . --format junit --exit-policy never
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
v0.1 schemas remain available for archived consumers. Adapter-manifest v0.1
through v0.3 are likewise frozen while v0.4 is the current executable contract.
Their local dependency graphs use versioned filenames, so later unversioned
aliases cannot change an archived contract.

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

The profile has 40 atomic assertions. Go, container, Terraform, Kubernetes, and
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

## Optional offline infrastructure policy scan

The focused `prc/iac@0.1` profile authorizes the exact checked-in Checkov 3.3.8
manifest when Terraform, Kubernetes, or container definitions are inventoried.
This is not part of the simple native scan because it launches a pinned OCI
analyzer. Read the [infrastructure policy profile](infrastructure-policy.md),
review the manifest, pre-pull its exact digest, and grant local verification
explicitly:

```bash
docker pull docker.io/bridgecrew/checkov@sha256:c64ffb6d6fc8087c896341a2c697770a04a1cf558db04fa7b8129d8ca6bce336

./prc scan /path/to/project \
  --profile prc/iac \
  --mode verify-local \
  --adapter-manifest adapters/checkov-v3.3.8.yaml
```

The scanner mounts a sealed inventory read-only, denies network access, uses
scanner-owned offline policy, and prevents target `.checkov.yml` settings from
hiding results. A policy violation becomes a detailed finding. Parsing errors,
inline suppressions, unsupported output, or unsafe metadata fail closed. The
command still produces a report only and never applies an infrastructure
change.

Read [bounded isolated remediation](remediation.md) to run the deterministic R1
loop, create and verify one R1 candidate, or apply one validated R2 provider
proposal outside the source workspace.

Read [read-only agent providers](../architecture/agent-providers.md) to inspect
or explicitly run the experimental Codex and Claude Code `suggest` adapters.
Provider proposals are never applied or treated as evidence of a passing check.
