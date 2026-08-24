# Safe AI review of all controls

`everylast scan` always puts all 10,042 controls in the report. By default it checks
only facts that the local scanner can prove safely. Broad questions stay
`needs_review` instead of being guessed.

An optional AI review can add advice for those broad questions. You choose
Codex or Claude Code. The AI result is extra information; it cannot change the
scanner's real pass, fail, blocked, or review state, and it never fixes files.

## Before you start

Install one supported command-line tool:

- `codex` for Codex; or
- `claude` for Claude Code.

Sign in once through PRC. This launches the provider's official login flow but
stores its credentials in a private Everylast-only directory:

```bash
everylast login codex
# or: everylast login claude

everylast auth
```

PRC does not reuse the provider's normal user configuration, sessions,
instructions, plugins, hooks, or MCP servers. A scan uses the Everylast-only login
with a fresh temporary home and a small runtime environment allowlist. Use
`everylast logout codex` or `everylast logout claude` to remove the saved Everylast login.

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
everylast full codex
```

Use `everylast full claude` for Claude Code. `everylast scan --ai codex|claude` remains an
equivalent advanced-friendly spelling. To test only one control first, use the
advanced form:

```bash
everylast scan /path/to/project \
  --review-provider codex \
  --review-control PRC-02-001 \
  --allow-remote-source-processing
```

For Claude Code, replace the provider name:

```bash
everylast scan /path/to/project \
  --review-provider claude \
  --review-control PRC-02-001 \
  --allow-remote-source-processing
```

Open the `Detailed report:` path printed at the end. Find `PRC-02-001`. Its
normal scanner state is still present. The AI section separately shows:

- its suggested result;
- whether it thinks the rule applies;
- confidence;
- reason and advice;
- exact excerpt lines it used;
- separate citation-location and claim verification states; and
- what it could not prove.

`snapshot_location_validated` means the path and line existed in the exact
screened snapshot bound to the task. It does not mean the line supports the AI
sentence. The claim therefore remains `advisory_unverified` until an independent
typed verifier or a qualified person proves it.

## Review all 10,042 controls

The short commands review all active controls:

```bash
everylast full codex
everylast full claude
```

Use the advanced form only when changing defaults, for example Codex `xhigh`
effort:

```bash
everylast scan /path/to/project \
  --review-provider codex \
  --review-effort xhigh \
  --allow-remote-source-processing
```

This is deliberately slow and expensive. The scanner makes batches of at most
eight controls. For each batch, it tells the top AI process to create exactly
one separate subagent per control, wait for all of them, and return exactly one
result per control. With 10,042 active controls and the default batch size,
expect about 1,256 provider calls and 10,042 subagent reviews. The exact token
and money cost depends on the chosen provider and model.

The scanner can require this orchestration in the sealed task and verify that
one final result returns for every control. Current provider output does not
offer trustworthy proof of each internal subagent call, so the scanner does not
treat claimed orchestration as evidence. A provider that skips a requested
subagent can at most produce untrusted advisory text, never a verified Pass.

Completed batches are stored privately outside the target project. If a later
batch fails or the run is stopped, run the same command again. Matching finished
batches are checked and reused. The report is written only after every requested
batch has a valid result.

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
| `--review-workers` | `1` | Run 1–4 provider calls at once. One is safer for cost and rate limits. |
| `--review-timeout` | `30m` | Limit for each resumable batch. |
| `--review-effort` | `high` | Codex supports `high` or `xhigh`; Claude uses `high`. |
| `--review-model` | provider default | Pin a provider model when needed. |
| `--review-state-dir` | user cache | Choose a private resume directory outside the target. |
| `--review-max-cost-usd` | no limit | Claude-only cost limit for each batch. Codex does not expose the same hard CLI limit. |

More workers can finish sooner but increase simultaneous cost and the chance of
provider rate limits. More controls per call reduce top-level calls but make the
coordinator's job larger. The defaults favor correctness and easy resuming.

## What stops the run

The run stops with an execution error if:

- a secret-like value is found before remote review;
- the provider is missing, changes during the run, times out, or exits badly;
- the schema changes during the run;
- output is too large, malformed, incomplete, reordered, duplicated, or cites
  a path or line outside the screened copy; or
- a saved batch does not match its sealed task.

Valid completed batches remain saved. No target file is changed, and no partial
AI result is allowed to make a scanner rule pass.
