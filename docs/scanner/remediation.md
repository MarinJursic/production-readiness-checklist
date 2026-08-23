# Bounded isolated remediation

The experimental `remediate` command supports two deterministic fixers:

- `PRC-A-CORE-014` appends one line-feed byte to each recognized source file
  that lacks one; and
- `PRC-A-CORE-022` clears group-write and other-write permission bits while
  preserving every file byte and all other permission bits.

Neither fixer executes target code, calls a model, uses the network, edits the
original workspace, or performs version-control operations.

The README-presence assertion is deliberately R2. Writing useful project
documentation requires project-specific judgment, so the scanner does not create
a placeholder README and claim the finding is fixed.

## Create an isolated candidate

Choose a new destination outside the target tree. Its parent must already exist,
and the destination itself must not exist.

```bash
./prc remediate \
  --catalog-root /path/to/production-readiness-checklist \
  --target /path/to/project \
  --config /path/to/project/production-readiness.yaml \
  --assertion PRC-A-CORE-014 \
  --candidate-dir /safe/path/prc-candidate \
  --max-files 20 \
  --max-changed-lines 20
```

The command exits `0` only when the candidate passes every acceptance check. A
validated but rejected candidate is still printed and exits `8`. Invalid input
exits `3`; a policy-denied remediation exits `5`. Use
`--format json` for the versioned `prc.remediation-candidate/v0.4` record.

With `--config`, the exact canonical configuration digest and project identity
are recorded in the fix contract. The configured profile is mandatory,
`remediation.enabled` must be true, and command-line file or line limits cannot
raise the configured ceilings. The configured attempt ceiling is recorded even
though this command performs exactly one attempt. Scanner defaults, configured
protected paths, and the configuration file itself are unioned into the guard
set. Baseline, candidate, and final source-integrity scans all rebind the same
configuration.

## What acceptance verifies

The scanner creates a content-addressed fix contract with the baseline run,
inventory, exact canonical finding ID and stable finding fingerprint, allowed
paths, protected paths, network denial, one-attempt limit, and file and line
budgets. After applying the fix in the copy, it:

1. inventories the candidate from fresh bytes;
2. walks the raw candidate tree so excluded directories cannot hide additions;
3. rejects additions, deletions, symlinks, non-regular entries,
   protected-path changes, and any change outside the allowlist;
4. verifies the fixer's exact byte-and-mode postcondition: one appended
   line-feed with the mode unchanged, or bytes unchanged with only group-write
   and other-write bits cleared;
5. rescans the candidate and requires the target assertion to pass; and
6. requires every assertion that passed in the baseline to remain passing; and
7. re-inventories the original target and requires it to remain byte-for-byte
   and mode-for-mode identical to the baseline.

The candidate directory is preserved for review. Acceptance is permission to
inspect or continue testing that isolated candidate; it is not authorization to
merge, deploy, release, accept risk, or claim that the full profile is satisfied.

## Run the bounded remediation loop

Without a provider, `prc fix` repeatedly applies only the registered R1 fixers.
Each accepted candidate becomes the source of a new sibling candidate, so fixes
compose without changing the original project. The candidate root must be a new
path outside the target tree. The scanner validates that destination first but
does not create it until an eligible attempt has passed task, source, policy,
and budget preflight.

```bash
./prc fix \
  --catalog-root /path/to/production-readiness-checklist \
  --target /path/to/project \
  --config /path/to/project/production-readiness.yaml \
  --candidate-root /safe/path/prc-remediation-run \
  --max-attempts 3 \
  --max-files 20 \
  --max-changed-lines 200 \
  --max-duration-seconds 1800 \
  --format json > remediation-run.json
```

The loop evaluates findings in profile order and stops predictably when it has
closed all eligible findings, a cumulative budget is exhausted, a provider
cannot return a candidate, or an independently checked candidate is rejected.
File and changed-line usage accumulates across accepted candidates. One
wall-clock budget spans planning, provider execution, candidate generation, and
verification. It is propagated as a hard deadline to provider and verifier
child processes and checked at scanner phase boundaries; a synchronous
candidate operation that returns after the deadline is preserved but rejected,
never accepted. Command-line values cannot raise limits declared in project
configuration. Go contexts propagate deadline cancellation, and commands
started with `CommandContext` are interrupted when that context completes
([context package](https://pkg.go.dev/context),
[os/exec package](https://pkg.go.dev/os/exec)).

The `prc.remediation-run/v0.8` report records every actual attempt, including
proposals rejected before candidate creation. Each attempt binds its sequence,
mode, exact finding and fingerprint, scanner-owned task, before and after
inventory digests, provider execution or failure and candidate when present, timestamps,
outcome, machine-readable reason code, and exact scanner rejection reason. The
scanner verifies this linkage before computing the run content ID. The report
also preserves every candidate, provider transcript digest, cumulative budget
usage, final fresh assessment, final isolated workspace, and a reason code for
every unresolved result. Its embedded v0.9 scan result preserves adapter
resolution provenance; frozen v0.3 through v0.7 remediation schemas retain their
version-pinned dependency graphs. Every unresolved failure includes its
canonical finding ID and stable fingerprint. Its terminal states are:

- `profile_satisfied`: every required result in the selected profile passed;
- `machine_work_complete`: no registered deterministic R1 failure remains, but
  manual evidence, blocked checks, or higher-risk work can still remain;
- `stopped_by_policy_or_budget`: an eligible fix could not run within policy;
- `candidate_rejected`: independent acceptance rejected an attempted fix;
- `provider_stopped`: the provider returned `unable` or `needs_escalation`
  without a patch; or
- `provider_failed`: scanner preflight, transcript persistence, process,
  timeout, output-bound, postflight-integrity, or output-protocol validation
  failed and the loop did not retry.

Exit status `0` is reserved for `profile_satisfied`. A no-go gate exits `1`,
incomplete or blocked assessment work and `provider_stopped` exit `2`, a policy
or budget stop exits `5`, provider failure exits `4`, and candidate rejection
exits `8`. Caller cancellation retains exit `7` while preserving its failure
record.
`machine_work_complete` is not a production-readiness claim. The default loop
does not invoke an agent. No loop mode runs project commands, deploys, merges,
or performs version-control operations.

## Opt in to one scanner-planned R2 task

`prc fix --provider` connects the bounded loop to the read-only Codex or Claude
Code provider adapter. This is not general repository autonomy. The current
task planner supports only a failing `PRC-A-CORE-010` test-discovery assertion:
it selects one bounded source file, derives a small allowlist of new test paths,
binds the exact triggering finding into the sealed task, and asks for exactly
one non-vacuous test file. Every other R2 finding remains
`no_safe_agent_task` until it has a dedicated planner and sufficiently strong
verification.

```bash
./prc fix \
  --catalog-root /path/to/production-readiness-checklist \
  --target /path/to/project \
  --config /path/to/project/production-readiness.yaml \
  --candidate-root /safe/path/prc-remediation-run \
  --provider codex \
  --allow-remote-source-processing \
  --verifier-runtime docker \
  --verifier-image registry.example/prc/python-verifier@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef \
  --max-attempts 3 \
  --max-duration-seconds 1800 \
  --format json > remediation-run.json
```

The remote-processing flag is mandatory because the sealed prompt contains the
selected source file. High-confidence secret-like material in that input causes
a policy stop before the candidate root or provider output directory is created;
the error identifies only the source path and detector family. The provider
receives no source-workspace access, shell, network tool, MCP server, secret,
edit tool, or candidate workspace. Its output
directory is private and separate from both source and candidate. Each provider
attempt preserves `agent-task.json`, bounded stdout/stderr transcripts, their
digests, the executable digest, and the schema digest. A proposal is parsed and
applied by the scanner exactly once in a fresh candidate, then passes the same
structural, anti-gaming, target-assertion, regression, and source-integrity
checks as `remediate-proposal`.

`--verifier-image` is mandatory with `--provider`. The reference must include an
explicit registry and immutable SHA-256 digest, and the image must already be
present because the runner passes `--pull=never`. The provider and repository
cannot choose or alter the image, command, or limits. The scanner infers one
supported command solely from the sealed source path: `go test ./...` for Go,
`python -m pytest -q` for Python, or `node --test` for plain JavaScript. The
scanner never invokes package scripts. TypeScript proposal planning remains
fail-closed because no scanner-owned verifier is registered for that ecosystem.
Because verification is network-denied and starts with empty scratch caches, the
selected image must already contain the required toolchain and test runner, and
the candidate must vendor or otherwise carry every dependency it needs. The
project does not yet publish a canonical verifier image; operators are
responsible for reviewing and pinning one for each supported ecosystem.

The candidate is mounted read-only into a read-only container with no network,
all Linux capabilities dropped, `no-new-privileges`, the caller's non-root
numeric user, bounded CPU, memory and swap, processes, file descriptors,
scratch space, time, stdout, and stderr. These flags implement documented
[Docker runtime isolation and resource controls](https://docs.docker.com/engine/containers/run/).
The output record contains hashes and byte counts rather than raw test output,
binds the configured candidate identity separately from its raw workspace-byte
inventory, and verifies that the candidate bytes did not change during the
run. A test failure, timeout, output limit, unavailable runtime or image, or
integrity change rejects the candidate; none is converted to a pass.

If invocation fails before a valid provider output exists, the scanner writes a
content-addressed `prc.agent-failure/v0.1` record. It uses a scanner-authored
safe reason, links the sealed task and provider identities, distinguishes the
failure stage and reason code, and records whichever bounded transcripts were
successfully persisted. The failed attempt consumes one attempt, is terminal,
and is never retried automatically.

Current acceptance reconstructs the proposed test before candidate creation,
requires a conventionally collectable declaration and a recognized behavioral
failure check, and establishes that the test was added without weakening
existing tests or regressing prior scanner passes. It then runs the supported
scanner-owned test command in the isolated verifier and requires exit status
zero. Go documents `go test ./...` as package-list mode over packages below the
current directory, and a failing Go test returns a nonzero status
([Go command reference](https://go.dev/cmd/go/)). Passing the suite proves only
that the configured command passed in that image; it does not prove complete
behavioral coverage. Broader R2 autonomy stays disabled until each task has an
assertion-specific behavioral contract and verifier.

## Apply one validated R2 proposal

`remediate-proposal` is the scanner-owned bridge from a validated Codex or
Claude Code `suggest` result to an isolated candidate. It never asks the provider
to edit files and never executes provider-authored commands. The source task must
still match the current workspace inventory, the assertion must be classified
R2, its exact finding ID and stable fingerprint must still match a freshly
reproduced failure, and the candidate destination must be new and outside the
source tree. The v0.3 fix contract records both the task's triggering finding
ID and the freshly verified baseline finding ID so the provider handoff remains
auditable end to end.

```bash
./prc remediate-proposal \
  --catalog-root /path/to/production-readiness-checklist \
  --target /path/to/project \
  --config /path/to/project/production-readiness.yaml \
  --provider codex \
  --task /safe/path/task.json \
  --output /safe/path/validated-provider-output.json \
  --candidate-dir /safe/path/prc-r2-candidate \
  --verifier-runtime docker \
  --verifier-image registry.example/prc/python-verifier@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef \
  --max-files 20 \
  --max-changed-lines 200
```

The strict internal unified-diff parser supports exact-context modification and
mode-`0644` text-file addition. It rejects deletions, renames, copies, binary or
mode changes, CRLF patch encoding, non-newline markers, malformed coordinates,
context mismatches, unlisted or protected paths, and budget overruns before
acceptance. A conservative anti-gaming pass also rejects changes to existing test
or specification files, newly introduced suppression or skip directives,
constant assertions, empty Go tests, test-shaped files without a collectable
declaration, and invocation-only tests without a recognized behavioral failure
check. Language-aware payload rules also reject generated tests that request
process or shell execution, network clients, environment or secret access,
filesystem mutation or absolute-path reads, dynamic evaluation, deserialization
execution, or long encoded payloads before creating a candidate. New focused
local assertions remain permitted. The payload audit is deliberately
conservative and syntactic; the independent sandbox remains the containment
boundary for obfuscation or parser gaps. The raw-tree, byte, mode, target assertion, baseline
regression, and source-integrity audits then run from fresh inventories.

Current R2 acceptance proves the declared deterministic scanner postcondition,
structural anti-gaming checks, scanner non-regression envelope, and successful
execution of the registered sandbox command. The anti-gaming checks catch known
unsafe patch shapes, and the test execution proves the suite passed; neither
alone proves that a new test is behaviorally sufficient. Broader code-changing
autonomy still requires task-specific verification and stronger behavioral
assertions before it can be enabled by policy.
