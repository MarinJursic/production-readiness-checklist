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
`--format json` for the versioned `prc.remediation-candidate/v0.3` record.

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
path outside the target tree.

```bash
./prc fix \
  --catalog-root /path/to/production-readiness-checklist \
  --target /path/to/project \
  --config /path/to/project/production-readiness.yaml \
  --candidate-root /safe/path/prc-remediation-run \
  --max-attempts 3 \
  --max-files 20 \
  --max-changed-lines 200 \
  --format json > remediation-run.json
```

The loop evaluates findings in profile order and stops predictably when it has
closed all eligible findings, a cumulative budget is exhausted, a provider
cannot return a candidate, or an independently checked candidate is rejected.
File and changed-line usage accumulates across accepted candidates;
command-line values cannot raise limits declared in project configuration.

The `prc.remediation-run/v0.3` report records every candidate, provider
execution and transcript digest, cumulative budget usage, the final fresh
assessment, the final isolated workspace, and a reason code for every
unresolved result. Every unresolved failure includes its canonical finding ID
and stable fingerprint. Its terminal states are:

- `profile_satisfied`: every required result in the selected profile passed;
- `machine_work_complete`: no registered deterministic R1 failure remains, but
  manual evidence, blocked checks, or higher-risk work can still remain;
- `stopped_by_policy_or_budget`: an eligible fix could not run within policy;
- `candidate_rejected`: independent acceptance rejected an attempted fix; or
- `provider_stopped`: the provider returned `unable` or `needs_escalation`
  without a patch.

Exit status `0` is reserved for `profile_satisfied`. A no-go gate exits `1`,
incomplete or blocked assessment work and `provider_stopped` exit `2`, a policy
or budget stop exits `5`, and candidate rejection exits `8`.
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
  --max-attempts 3 \
  --format json > remediation-run.json
```

The remote-processing flag is mandatory because the sealed prompt contains the
selected source file. The provider receives no source-workspace access, shell,
network tool, MCP server, secret, edit tool, or candidate workspace. Its output
directory is private and separate from both source and candidate. Each provider
attempt preserves `agent-task.json`, bounded stdout/stderr transcripts, their
digests, the executable digest, and the schema digest. A proposal is parsed and
applied by the scanner exactly once in a fresh candidate, then passes the same
structural, anti-gaming, target-assertion, regression, and source-integrity
checks as `remediate-proposal`.

Current acceptance establishes that a discoverable, structurally non-vacuous
test was added without weakening existing tests or regressing prior scanner
passes. It does not execute project tests or prove behavioral coverage. Broader
R2 autonomy stays disabled until sandboxed, scanner-owned verification commands
and assertion-specific behavioral contracts are implemented.

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
  --max-files 20 \
  --max-changed-lines 200
```

The strict internal unified-diff parser supports exact-context modification and
mode-`0644` text-file addition. It rejects deletions, renames, copies, binary or
mode changes, CRLF patch encoding, non-newline markers, malformed coordinates,
context mismatches, unlisted or protected paths, and budget overruns before
acceptance. A conservative anti-gaming pass also rejects changes to existing test
or specification files, newly introduced suppression or skip directives,
constant assertions, and empty Go tests before creating a candidate. New focused
tests remain permitted. The raw-tree, byte, mode, target assertion, baseline
regression, and source-integrity audits then run from fresh inventories.

Current R2 acceptance proves only the declared deterministic scanner
postcondition, structural anti-gaming checks, and non-regression envelope. The
anti-gaming checks catch known unsafe patch shapes; they do not prove that a new
test is behaviorally sufficient. Broader code-changing autonomy must add
sandboxed project-specific verification commands and stronger behavioral
assertions before it can be enabled by policy.
