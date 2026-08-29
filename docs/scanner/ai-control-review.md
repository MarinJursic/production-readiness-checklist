# Safe AI review of nondeterministic controls

`prc scan` always puts all 10,042 controls in the report. By default it checks
only facts that the local scanner can prove safely. Broad questions stay
`needs_review` instead of being guessed.

An optional AI review can add advice for the 9,356 controls whose reviewed
classification is nondeterministic. The 686 deterministic controls are not
handed to AI for a verdict. You choose Codex or Claude Code. The AI result is
extra information; it cannot change the scanner's real pass, fail, blocked, or
review state, and it never fixes files.

## Before you start

Install one supported command-line tool:

- `codex` for Codex; or
- `claude` for Claude Code.

Sign in once through the scanner. This launches the provider's official login flow but
stores its credentials in a private scanner-only directory:

```bash
prc login codex
# or: prc login claude

prc auth
```

The scanner does not reuse the provider's normal user configuration, sessions,
instructions, plugins, hooks, or MCP servers. A scan uses the scanner-only login
with a fresh temporary home and a small runtime environment allowlist. Use
`prc logout codex` or `prc logout claude` to remove the saved scanner login.

Supported API-key environment variables remain an alternative: Codex accepts
`OPENAI_API_KEY` or `CODEX_API_KEY`; Claude accepts `ANTHROPIC_API_KEY`,
`ANTHROPIC_AUTH_TOKEN`, or `CLAUDE_CODE_OAUTH_TOKEN`. Only the selected
credential is forwarded. Remove a temporary variable from the shell after use.

AI review can send source excerpts to a remote model and can cost money. Read
the provider's data and billing rules first. The short `--ai` option is explicit
permission for screened source processing. The advanced form requires the
separate `--allow-remote-source-processing` switch.

## Try one control first

The shortest full-review command is:

```bash
prc full codex
```

Use `prc full claude` for Claude Code. `prc scan --ai codex|claude` uses the
same guarded review engine but keeps the advanced defaults of standard depth,
one worker, and high reasoning. To test only one control first, use the advanced
form:

```bash
prc scan /path/to/project \
  --review-provider codex \
  --review-control PRC-02-001 \
  --allow-remote-source-processing
```

For Claude Code, replace the provider name:

```bash
prc scan /path/to/project \
  --review-provider claude \
  --review-control PRC-02-001 \
  --allow-remote-source-processing
```

Open the `Detailed report:` path printed at the end. Find `PRC-02-001`. Its
normal scanner state is still present. The AI section separately shows:

- its suggested result;
- whether it thinks the rule applies;
- confidence;
- priority, reason, risk, and advice;
- ordered remediation and independent verification steps;
- evidence still needed and the strongest skeptical challenge;
- exact excerpt lines it used;
- separate citation-location and claim verification states; and
- what it could not prove.

`snapshot_location_validated` means the path and line existed in the exact
screened snapshot bound to the task. It does not mean the line supports the AI
sentence. The claim therefore remains `advisory_unverified` until an independent
typed verifier or a qualified person proves it.

## Review all 9,356 nondeterministic controls

The short commands review all active controls whose reviewed classification is
nondeterministic. All 10,042 controls still appear in the report, but the 686
deterministic controls never receive an AI verdict. They remain Blocked until
their exact program gets complete authoritative evidence.

```bash
prc full codex
prc full claude
```

The short `prc full` path is quality-first: it uses deep review, four parallel
provider workers, and Codex `xhigh` reasoning. The advanced form can choose
different settings, for example a one-control standard review:

```bash
prc scan /path/to/project \
  --review-provider codex \
  --review-control PRC-02-001 \
  --review-depth standard \
  --allow-remote-source-processing
```

This is deliberately slow and expensive. The scanner makes batches of at most
eight controls. For each batch, it tells the top AI process to create exactly
one separate primary subagent per control. Deep mode also creates one separate
skeptical subagent for the batch, runs the independent work concurrently, and
requires the coordinator to preserve the strongest objection or counterexample
for every result. With 9,356 nondeterministic controls and the default batch
size, expect about 1,170 provider calls, 9,356 primary rule reviews, and about
1,170 batch-skeptic reviews. The exact token and money cost depends on the
chosen provider and model.

Before the first provider call, the terminal shows the exact number of controls,
batches, cached batches, workers, and the private resume directory. Small
reviews print every completed batch. Large reviews print bounded progress at
percentage changes, including elapsed time and completed control counts, so
thousands of calls do not flood the terminal.

Codex reviews use its JSONL event stream. When the installed Codex CLI reports
usage on `turn.completed`, the scanner totals input, cached-input, output, and
reasoning tokens for new batches. Claude JSON output provides an estimated cost;
the scanner totals that estimate and labels it as an estimate rather than a
bill. A nonzero Claude `--review-max-cost-usd` remains an enforced limit for
each new batch, including subagent spend on supported Claude Code versions.
Codex does not expose an equivalent hard dollar limit through this path.

[OpenAI documents the Codex JSONL usage event](https://learn.chatgpt.com/docs/non-interactive-mode#make-output-machine-readable).
[Claude Code documents JSON cost estimates and their limits](https://code.claude.com/docs/en/headless#pipe-data-through-claude).
Cached review files contain sealed advice, not old billing records, so a run
that reuses cached batches does not invent or re-count their past tokens or
cost. The final terminal summary says how many new batches supplied accounting.

The scanner can require this orchestration in the sealed task and verify that
one final result returns for every control. Current provider output does not
offer trustworthy proof of each internal subagent call, so the scanner does not
treat claimed orchestration as evidence. A provider that skips a requested
subagent can at most produce untrusted advisory text, never a verified Pass.

Completed batches are stored privately outside the target project. If a later
batch fails or the run is stopped, the scanner writes a partial report with all
completed, schema-checked reviews, returns an execution error, and leaves the
target unchanged. Run the same command again: matching finished batches are
checked and reused, and only unfinished work is sent again. A partial report is
never labeled complete and no partial AI result can change a control's real
disposition.

## What is sent to the provider

The scanner first creates a temporary, private copy of selected text files. It:

- uses the exact files and hashes from the local inventory;
- does not follow symlinks;
- skips `.env`, `.npmrc`, credentials, private-key files, AI instruction files,
  and similar sensitive names;
- rejects known high-confidence token and private-key shapes before any provider
  starts;
- skips binary and invalid text;
- limits each copied file, the total copy, path list, excerpt size, file count,
  and provider output size; and
- removes the temporary copy at the end.

The provider gets a sealed JSON task containing control text, deterministic
check context, repository paths, and selected excerpts. Repository text is
marked as untrusted data. A file that says “ignore the scanner” is not treated
as an instruction.

The scanner never gives the provider the source workspace path. For Codex it
disables shell, web search, apps, browser, computer use, MCP servers, and project
rules, uses a read-only sandbox, and allows only subagent coordination. For
Claude Code it allows only the `Agent` tool and blocks Bash, Read, Glob, Grep,
Write, Edit, web tools, and questions. Provider output must match a strict JSON
schema, the sealed task ID, every requested control in order, safe evidence
paths and real line numbers, and fixed size limits.

This screen is a safety layer, not a complete secret scanner. Do not opt in with
source that you are not allowed to send to the selected provider.

The scanner stops if it can see local Claude managed settings because those can
force hooks, plugins, or MCP policy that an ordinary setting cannot override.
It cannot inspect every policy delivered by a provider's server. The provider
binary and remote service therefore remain trusted dependencies even though
they receive no target path or general-purpose workspace tools.

## Result meanings

AI review uses candidate words on purpose:

| AI value | Simple meaning |
| --- | --- |
| `advisory_pass_candidate` | The shown excerpts look consistent with the rule, but this is not a verified Pass. |
| `advisory_fail_candidate` | The shown excerpts suggest a problem worth fixing or checking. |
| `needs_evidence` | Source excerpts cannot prove the rule. Runtime, production, legal, company, or human proof may be needed. |
| `not_applicable_candidate` | The AI thinks the trigger is absent. A trusted check or person must still make the final decision. |

Folder structure is a good example. The AI may explain that a surprising layout
makes ownership or navigation unclear, but it must not require universal names
such as `src`, `components`, or `services`. Advice must fit the project's real
language, build system, size, boundaries, and conventions.

## Useful controls

| Option | Default | Meaning |
| --- | ---: | --- |
| `--ai` | off | Short form: choose `codex` or `claude` and explicitly allow screened remote processing. |
| `--review-provider` | `none` | Choose `codex` or `claude`. No AI starts by default. |
| `--review-control ID` | all active | Review one named control; repeat for a small test set. |
| `--review-batch-size` | `8` | Provider calls contain 1–8 controls; every control still gets its own subagent. |
| `--review-depth` | `standard` | `standard` uses one primary subagent per rule; `deep` also adds an independent skeptical subagent per batch. `prc full` selects `deep`. |
| `--review-workers` | `1` | Run 1–4 provider calls at once. One is safer for cost and rate limits. |
| `--review-timeout` | `30m` | Limit for each resumable batch. |
| `--review-effort` | `high` | Codex supports `high` or `xhigh`; Claude uses `high`. |
| `--review-model` | provider default | Pin a provider model when needed. |
| `--review-state-dir` | user cache | Choose a private resume directory outside the target. |
| `--review-max-cost-usd` | no limit | Claude-only enforced cost limit for each new batch. The terminal separately labels Claude's reported total as a client estimate. Codex does not expose the same hard CLI limit. |

More workers can finish sooner but increase simultaneous cost and the chance of
provider rate limits. More controls per call reduce top-level calls but make the
coordinator's job larger. Advanced `prc scan` defaults to one worker and
standard depth. The simple `prc full` command deliberately selects four workers
and deep review because it is the quality-first path.

## What stops the run

The run stops with an execution error if:

- a secret-like value is found before remote review;
- the provider is missing, changes during the run, times out, or exits badly;
- the schema changes during the run;
- output is too large, malformed, incomplete, reordered, duplicated, or cites
  a path or line outside the screened copy; or
- a saved batch does not match its sealed task.

Valid completed batches remain saved and appear in a clearly marked partial
report. No target file is changed, and no partial AI result is allowed to make
a scanner rule pass.
