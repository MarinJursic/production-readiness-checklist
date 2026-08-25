# Read-only agent providers

The experimental provider layer connects the scanner to installed Codex and
Claude Code CLIs without giving either provider authority over scanner truth.
The local CLI binary and remote provider are still trusted process and data
processing dependencies. The current mode is deliberately `suggest`: an agent receives only content-addressed
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

- a content-addressed `prc.agent-task/v0.2` task bound to one canonical finding
  ID and stable fingerprint;
- explicit acknowledgement that relevant source may be processed by a remote
  provider;
- an exact workspace inventory plus bounded, content-addressed task inputs;
- a fail-closed preflight that prevents obvious private keys, distinctive
  provider-token shapes, and credential-bearing URLs from entering a remote
  task without logging the matched value;
- a private execution/output directory disjoint from the source workspace;
- denied agent-tool network access and no task secrets;
- schema-constrained, non-interactive output;
- scanner-enforced time and output limits;
- executable and output-schema digests rechecked immediately before execution;
  and
- filtered process environments that exclude unrelated credentials such as
  cloud, repository, and deployment tokens.

Vuk can reuse only a login created by `vuk login codex` or `vuk login claude`.
Those commands call the provider's official authentication flow with a private
Vuk credential directory. They do not reuse the provider's normal user
configuration, sessions, instructions, plugins, hooks, or MCP servers. Supported
API-key environment variables remain an alternative. Each scan still gets a
new private home; only the selected credential, basic runtime variables, and
scanner-owned overrides reach the process.

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

These command flags cannot contain a malicious local executable. A replaced or
compromised `codex` or `claude` program runs as the current operating-system
user and could ignore every argument before the scanner detects the changed
digest. Install the CLI from its official source, keep it updated, inspect the
resolved path and digest with `vuk doctor`, and use a separate OS account or
strong external sandbox when the host contains secrets the CLI must never see.
The scanner also stops when it can see local Claude managed settings that may
force hooks, plugins, or MCP configuration, but it cannot inspect every policy
delivered by a provider server.

Repository text is always untrusted data, including comments that resemble
instructions or the scanner's task delimiter. The scanner JSON-encodes the
entire content-addressed task, so delimiter characters inside a source file are
escaped and cannot terminate the authoritative envelope. Checked-in adversarial
tests reconstruct that envelope and require the hostile fixture to round-trip
only as an input value. This is defense in depth: output schema validation,
capability denial, patch-path validation, and independent candidate checks remain
mandatory even when the prompt boundary is intact. These controls implement the
OWASP guidance to [identify external content, constrain privileges, validate
outputs, and test adversarially](https://genai.owasp.org/llmrisk/llm01-prompt-injection/).

These controls are not an R2 write sandbox. Neither provider may mutate a
candidate workspace. An agent proposal remains untrusted data. The
scanner-owned [`remediate-proposal`](../scanner/remediation.md#apply-one-validated-r2-proposal)
path can explicitly parse one validated proposal into a fresh isolated
candidate and run deterministic acceptance checks; the provider does not apply
or approve that result. The bounded `vuk fix --provider` path composes those
same two stages only for a scanner-planned missing-test task.

## Create and seal a task

Start from the checked-in
[`fixtures/providers/suggest-task.json`](https://github.com/MarinJursic/production-readiness-checklist/blob/main/fixtures/providers/suggest-task.json)
example.
For a draft, leave `task_id` empty, keep arrays sorted, tailor the goal and path
allowlists, and make the remote-processing decision explicit. Seal it without
modifying the draft:

```bash
./vuk provider seal-task \
  --file /path/to/draft-task.json \
  --workspace /path/to/project \
  --config /path/to/project/production-readiness.yaml \
  > /safe/path/task.json
```

Sealing reads only the sorted `relevant_paths`, rejects non-regular, binary, or
larger-than-256-KiB inputs, and embeds their text and SHA-256 digests. Total input
text is limited to 768 KiB. Before remote processing, sealing also rejects
high-confidence secret-like input without including the matched material in the
error. This conservative guard is not a comprehensive repository secret scan;
it covers selected high-precision generic and provider credential families
recognized by [GitHub's supported secret-scanning pattern
catalog](https://docs.github.com/en/code-security/reference/secret-security/supported-secret-scanning-patterns),
and projects should still run a dedicated secret scanner. Sealing also binds the
current workspace inventory digest into the task. The task's `finding_id` names the exact scan finding that caused
the task to be created; both it and `finding_fingerprint` are revalidated
against a fresh baseline before a proposal can be applied. For a manual draft,
copy both values from the same canonical finding in a current JSON scan; the
bounded loop does this automatically. When configured, that inventory includes
the declared-scope digest, and the scanner merges its default guards, configured
protected paths, and the in-target configuration path into the sealed task.
Changing a task field, finding binding, configuration declaration, or any
inventoried workspace file invalidates the execution plan.

## Inspect a launch plan

Create a new output directory outside the workspace, then inspect the exact
executable identity, arguments, environment names, schema identity, and
capabilities before any provider call:

```bash
install -d -m 700 /safe/path/provider-output
export OPENAI_API_KEY='your-provider-key'

./vuk provider plan \
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
./vuk provider run \
  --provider codex \
  --task /safe/path/task.json \
  --workspace /path/to/project \
  --output-dir /safe/path/provider-output
```

The scanner writes raw standard output and diagnostics as mode-`0600` transcript
files in the output directory, records their byte counts and SHA-256 digests,
and emits a `prc.agent-execution/v0.1` record after valid output. A failed
invocation instead returns a content-addressed `prc.agent-failure/v0.1` record
with scanner-authored stage and reason codes plus complete or partial transcript
metadata. Treat transcripts as sensitive source-derived evidence.

## Validate without running

Golden and adversarial provider outputs can be checked independently:

```bash
./vuk provider validate-output \
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
workflow when the task's assertion is R2. `provider run` itself never applies a
proposal. The optional [bounded R2 loop](../scanner/remediation.md#opt-in-to-one-scanner-planned-r2-task)
constructs and seals its own narrow task, invokes the same read-only provider
protocol, and hands valid output to the isolated scanner-owned application
path. A proposal is never applied to the source workspace.
