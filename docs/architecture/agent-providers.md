# Read-only agent providers

The experimental provider layer connects the scanner to installed Codex and
Claude Code CLIs without making either provider part of the trust base. Its
current mode is deliberately `suggest`: an agent receives only content-addressed
copies of declared relevant text files and returns a schema-constrained patch
proposal. It cannot inspect the source workspace, edit files, run shell commands,
use web or MCP tools, change capabilities, or mark an assertion as passed.

The implementation follows the current official documentation for [Codex
non-interactive mode](https://developers.openai.com/codex/noninteractive/),
[Codex sandboxing](https://developers.openai.com/codex/sandboxing/), [Claude Code
headless mode](https://code.claude.com/docs/en/headless), and [Claude Code
permissions](https://code.claude.com/docs/en/permissions). Provider flags are
also covered by golden launch-plan tests so an accidental dangerous flag fails
review.

## Capability boundary

Both launch plans require:

- a content-addressed `prc.agent-task/v0.1` task;
- explicit acknowledgement that relevant source may be processed by a remote
  provider;
- an exact workspace inventory plus bounded, content-addressed task inputs;
- a private execution/output directory disjoint from the source workspace;
- denied agent-tool network access and no task secrets;
- schema-constrained, non-interactive output;
- scanner-enforced time and output limits;
- executable and output-schema digests rechecked immediately before execution;
  and
- filtered process environments that exclude unrelated credentials such as
  cloud, repository, and deployment tokens.

Codex runs from the private output directory with ignored user configuration,
strict configuration, ephemeral sessions, the read-only sandbox, approval policy
`never`, no inherited shell environment, and the default shell tool disabled.
Hosted web search, MCP, multi-agent, goal, remote-plugin, app, browser,
computer-use, and image-generation features are also disabled. Claude Code runs
from the same isolated directory with an empty tool allowlist; Bash, reading,
editing, writing, notebook editing, web tools, slash commands, ambient MCP
servers, project setting sources, and session persistence are disabled. The
source text needed for either provider is inside the sealed task prompt. Claude's
optional provider-side cost limit is passed through; the current Codex CLI
adapter rejects a nonzero cost limit because it cannot enforce one.

These controls are not an R2 write sandbox. Neither provider may mutate a
candidate workspace. An agent proposal remains untrusted data. The separate
scanner-owned [`remediate-proposal`](../scanner/remediation.md#apply-one-validated-r2-proposal)
path can explicitly parse one validated proposal into a fresh isolated candidate
and run deterministic acceptance checks; the provider does not apply or approve
that result.

## Create and seal a task

Start from the checked-in
[`fixtures/providers/suggest-task.json`](https://github.com/MarinJursic/production-readiness-checklist/blob/main/fixtures/providers/suggest-task.json)
example.
For a draft, leave `task_id` empty, keep arrays sorted, tailor the goal and path
allowlists, and make the remote-processing decision explicit. Seal it without
modifying the draft:

```bash
./prc provider seal-task \
  --file /path/to/draft-task.json \
  --workspace /path/to/project \
  --config /path/to/project/production-readiness.yaml \
  > /safe/path/task.json
```

Sealing reads only the sorted `relevant_paths`, rejects non-regular, binary, or
larger-than-256-KiB inputs, and embeds their text and SHA-256 digests. Total input
text is limited to 768 KiB. It also binds the current workspace inventory digest
into the task. When configured, that inventory includes the declared-scope
digest, and the scanner merges its default guards, configured protected paths,
and the in-target configuration path into the sealed task. Changing a task
field, configuration declaration, or any inventoried workspace file invalidates
the execution plan.

## Inspect a launch plan

Create a new output directory outside the workspace, then inspect the exact
executable identity, arguments, environment names, schema identity, and
capabilities before any provider call:

```bash
install -d -m 700 /safe/path/provider-output

./prc provider plan \
  --provider codex \
  --task /safe/path/task.json \
  --workspace /path/to/project \
  --output-dir /safe/path/provider-output
```

Use `--provider claude` for Claude Code. `provider capabilities` reports the
static envelope for either adapter.

## Run explicitly

`provider run` invokes the installed provider and can consume provider quota or
incur provider charges. It should be called only after the operator reviews the
task's remote-source acknowledgement and launch plan.

```bash
./prc provider run \
  --provider codex \
  --task /safe/path/task.json \
  --workspace /path/to/project \
  --output-dir /safe/path/provider-output
```

The scanner writes raw standard output and diagnostics as mode-`0600` transcript
files in the output directory, records their byte counts and SHA-256 digests,
and emits a `prc.agent-execution/v0.1` record. Treat transcripts as sensitive
source-derived evidence.

## Validate without running

Golden and adversarial provider outputs can be checked independently:

```bash
./prc provider validate-output \
  --provider codex \
  --task fixtures/providers/suggest-task.json \
  --file fixtures/providers/valid-output.json
```

Validation rejects prose-only responses, duplicate JSON keys, task mismatches,
unsorted or duplicate files, changes outside the allowlist, protected paths,
reported command execution, requested capability expansion, trailing JSON, and
oversized output. Claude's outer JSON envelope must contain a non-error
`structured_output` value that passes the same validation.

After validation, follow the [isolated R2 proposal](../scanner/remediation.md#apply-one-validated-r2-proposal)
workflow when the task's assertion is R2. There is no automatic handoff from
`provider run`, and a proposal is never applied to the source workspace.
