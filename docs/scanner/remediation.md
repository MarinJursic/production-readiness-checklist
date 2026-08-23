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
validated but rejected candidate is still printed and exits `1`. Invalid input,
an ineligible baseline, or an infrastructure failure exits `2`. Use
`--format json` for the versioned `prc.remediation-candidate/v0.2` record.

With `--config`, the exact canonical configuration digest and project identity
are recorded in the fix contract. The configured profile is mandatory,
`remediation.enabled` must be true, and command-line file or line limits cannot
raise the configured ceilings. The configured attempt ceiling is recorded even
though this command performs exactly one attempt. Scanner defaults, configured
protected paths, and the configuration file itself are unioned into the guard
set. Baseline, candidate, and final source-integrity scans all rebind the same
configuration.

## What acceptance verifies

The scanner creates a content-addressed fix contract with the baseline run and
inventory, exact allowed paths, protected paths, network denial, one-attempt
limit, and file and line budgets. After applying the fix in the copy, it:

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

## Apply one validated R2 proposal

`remediate-proposal` is the scanner-owned bridge from a validated Codex or
Claude Code `suggest` result to an isolated candidate. It never asks the provider
to edit files and never executes provider-authored commands. The source task must
still match the current workspace inventory, the assertion must be classified
R2, and the candidate destination must be new and outside the source tree.

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
