# Bounded R1 remediation

The experimental `remediate` command supports one deterministic fixer:
`PRC-A-CORE-014`, which appends one line-feed byte to each recognized source file
that lacks one. It does not execute target code, call a model, use the network,
edit the original workspace, or perform version-control operations.

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
  --assertion PRC-A-CORE-014 \
  --candidate-dir /safe/path/prc-candidate \
  --max-files 20 \
  --max-changed-lines 20
```

The command exits `0` only when the candidate passes every acceptance check. A
validated but rejected candidate is still printed and exits `1`. Invalid input,
an ineligible baseline, or an infrastructure failure exits `2`. Use
`--format json` for the versioned `prc.remediation-candidate/v0.1` record.

## What acceptance verifies

The scanner creates a content-addressed fix contract with the baseline run and
inventory, exact allowed paths, protected paths, network denial, one-attempt
limit, and file and line budgets. After applying the fix in the copy, it:

1. inventories the candidate from fresh bytes;
2. walks the raw candidate tree so excluded directories cannot hide additions;
3. rejects additions, deletions, symlinks, non-regular entries, mode changes,
   protected-path changes, and any change outside the allowlist;
4. verifies that each allowed file differs by exactly one appended line-feed
   byte;
5. rescans the candidate and requires the target assertion to pass; and
6. requires every assertion that passed in the baseline to remain passing.

The candidate directory is preserved for review. Acceptance is permission to
inspect or continue testing that isolated candidate; it is not authorization to
merge, deploy, release, accept risk, or claim that the full profile is satisfied.
