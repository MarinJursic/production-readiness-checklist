# Safe start-to-finish scanner walkthrough

This page explains what happens from installation through a finished report,
what can go wrong at each step, and what the scanner does about it. Safe
behavior should not depend on knowing security jargon.

## 1. Choose one exact release

Find a `scanner-vX.Y.Z` release in this repository and use that exact version.
Do not copy a similar package name from a search result.

For a Node project:

```bash
npm install --save-dev --save-exact --ignore-scripts --no-audit --no-fund @marinjursic/everylast@X.Y.Z
```

What this changes:

- npm adds the exact scanner version to `devDependencies`;
- npm installs the small launcher and one matching native package;
- npm writes or updates `package-lock.json`; and
- npm does not run package install hooks because `--ignore-scripts` is present.

What could be attacked:

- A misspelled name could install somebody else's package. Use the exact scoped
  name `@marinjursic/everylast` from this repository.
- A floating version such as `latest` can change between machines. Use the exact
  release number and commit the lock file.
- Any package manager or registry can be compromised. Check the release's npm
  provenance and compare the version with the GitHub release before trusting
  it. The first npm publication also needs a one-time human bootstrap, which is
  clearly documented rather than hidden in automation.
- Install hooks are a common place for unwanted code. The scanner packages have
  no hooks, no third-party JavaScript runtime dependencies, and the command
  disables hooks anyway.

The npm launcher never downloads a binary. It reads the exact platform manifest,
checks the packaged native binary's SHA-256, and starts that file directly
without a shell. It does not fall back to a program named `prc` on `PATH`.

## 2. Add the easy command

Add this one line to the existing `scripts` object in `package.json`:

```json
{
  "scripts": {
    "scan": "everylast scan"
  }
}
```

The shortest repeatable command is:

```bash
npm run scan
```

For the stricter run, use `npm run --ignore-scripts scan`; the named `scan`
script runs, but local `prescan` and `postscan` hooks do not. If you do not want
to edit `package.json`, use this after the exact package is already installed:

```bash
npm exec --offline --no -- everylast scan
```

`--offline --no` makes a missing local command fail instead of asking to fetch
something from the registry. `npm scan` is not valid npm syntax for a custom
command.

## 3. Build a safe inventory

The scanner reads the target tree and records file paths, sizes, modes, hashes,
recognized project facts, exclusions, and Git identity. It does not run the
project's package manager, build, tests, scripts, hooks, containers, or code.

A hostile repository might contain links that point outside the project, huge
files, changing files, strange Git configuration, invalid text, or generated
trees designed to waste time. The inventory therefore:

- never follows repository symlinks;
- uses entry-count, per-file, and total-byte limits;
- checks a file again after reading it;
- runs Git with a small scanner-owned configuration and scopes status to the
  target worktree;
- excludes only clear caches, generated output, scanner state, and Git internals;
- records every exclusion and limit instead of silently calling the scan
  complete; and
- stops on unsafe or changing input instead of guessing.

## 4. Run the local checks

The normal profile currently has 40 narrow checks. They cover exact repository
facts such as a usable lock file, workflow parsing, pinned action references,
private-key armor, selected Kubernetes fields, OpenAPI structure, and direct Go
HTTP helper calls.

Parsers can have bugs, and one file pattern does not prove a broad engineering
promise. Each parser has input limits and returns an error or incomplete result
when it cannot safely decide. A narrow passing check produces only its narrow
evidence. It does not turn a large linked control into a complete Pass.

For example, finding `README.md` proves only that a nonempty root README exists.
It does not prove the instructions are correct. Detecting
`pull_request_target` asks for security review; trigger presence alone is not
reported as proof of a vulnerability. The six-hour workflow timeout is a named
value in the selected core profile, not a claim that every project must use the
same limit.

## 5. Include all 10,042 controls honestly

Every complete report contains every registered control and its exact control
contract. The generated contracts say what kind of evidence appears necessary,
whether complete inventory or a project threshold is needed, whether the rule
looks compound or negative, and what must be shown for Not Applicable.

All generated contracts are marked `generated_unreviewed`. They route future
work; they are not 10,042 approved automatic tests. Broad controls remain
`needs_review`, `blocked`, or `partially_verified` until the right evidence
exists. This is how the scanner uses every rule without inventing thousands of
false passes.

The full proposed acceptance review for every control is split into bounded
parts under `research/control-acceptance-criteria/`. Its `README.md` is the
index. A control owner must still approve and test a contract before it can
become a trusted automatic check.

## 6. Write the report

By default the scanner writes one private HTML file outside the target project
and prints its exact path. It creates a new file and never overwrites an old
one. The target bytes and file modes are checked before and after the scan in
the test suite.

The terminal uses green for Pass, red for Fail, yellow for blocked or manual
work, and plain text when output is redirected. The report separates:

- the result of the 40-check local profile;
- the state of the full 10,042-control catalog;
- verified findings and exact evidence;
- controls that still need proof; and
- inventory limits and exclusions.

Reports can contain project names, paths, evidence summaries, and optional AI
review text. Treat them as project data. Store them in a private location and do
not upload them automatically.

## 7. Optional AI review

A normal scan never contacts Codex, Claude, or another remote model. Sign in
through PRC once, then use the short AI option:

```bash
everylast login codex
everylast full codex
```

Use `claude` in both commands for Claude Code. The login is kept in a private
Everylast-only directory; normal provider settings, sessions, instructions, plugins,
hooks, and MCP servers are not loaded. `--ai` is also explicit permission to
send bounded, screened source excerpts to the chosen provider. Supported
temporary API-key environment variables and the longer
`--review-provider`/`--allow-remote-source-processing` form remain available.
`everylast full codex` is the short spelling of `everylast scan --ai codex`.

Before the provider starts, the scanner makes a private snapshot of bounded text
excerpts. It skips sensitive names and known key or token shapes, gives the
provider no target path, and marks repository text as untrusted data. Output
must cite a path and line that was actually shown. A valid location is recorded
as `snapshot_location_validated`, while the claim remains
`advisory_unverified`; a real line can still be irrelevant or misunderstood.
Missing, omitted, runtime, legal, company, or human evidence stays a limitation.

Prompt injection is still possible whenever untrusted text reaches a model.
The provider gets no general shell, source-reading, write, browser, web, or MCP
tool. Codex gets only subagent coordination for full review; Claude gets only
its Agent tool. The coordinator is required to create one separate subagent per
rule. The scanner checks the final schema, task ID, order, paths, lines, sizes,
and completeness, but current provider output does not give it trustworthy
proof that the provider really made every requested internal subagent call. A
provider that ignores the task can therefore produce only rejected or untrusted
advice, never verified evidence.

The local provider program itself is a bigger trust boundary. A malicious
`codex` or `claude` executable runs as your OS user and could ignore command
flags. The scanner hashes it and isolates configuration and environment, but it
cannot make a malicious executable safe. Install the official CLI, inspect its
path with `everylast doctor`, and use a separate OS account or stronger outer sandbox
on a sensitive computer. The remote provider also receives the approved
excerpts and may charge money. Start with one rule.

AI output is advice only. It cannot change the local result, create a verified
Pass, make the final Not Applicable decision, or fix a file.

## 8. Optional fixes

`everylast scan` has no path to the fix code. Fixing starts only with the separate
`everylast fix` command.

The fix system copies the target into a new private sibling candidate. It never
edits the original project. Deterministic fixes are limited to scanner-owned
changes with exact checks. The current AI path supports only one planned
missing-test task; the provider returns an untrusted patch proposal and never
edits a candidate itself.

A proposed fix might try to escape its path list, weaken tests, add a skip or
suppress a finding. The scanner rejects protected paths, unapproved files,
change-size limits, known test weakening, constant or empty tests, command
requests, policy changes, and capability expansion. It applies an accepted
proposal once to a fresh candidate, rescans fresh bytes, and runs a pinned,
network-denied verifier container for the supported test command.

The container runtime and verifier image are also trusted dependencies. The
image must already exist, use an immutable digest, and contain its dependencies;
the scanner never pulls it during the fix. Even a successful candidate is not
merged, committed, deployed, or called production ready. A person reviews the
candidate and decides what to keep.

## 9. Read the final state correctly

`profile_satisfied` means the selected executable profile has current acceptable
evidence. It does not mean all 10,042 broad controls passed or that no bug exists.
`needs_review`, `environment_blocked`, `no_go`, `machine_work_complete`, and
similar states keep the missing decision or evidence visible.

The safe end of a normal run is therefore a report, not an automatic rewrite:

```text
Scan mode: report only; no fixes were applied.
Detailed report: /private/path/example-api-91c2....html
```

Open that file, start with red failures and no-go findings, then review yellow
blocked/manual items. Green means only the exact displayed check passed.
