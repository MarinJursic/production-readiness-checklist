# CLI and exit-code contract

The scanner uses one stable process exit-code table across commands. Automation
must inspect both the exit code and the versioned report; it must never interpret
“a process ran” as “the profile passed.”

| Code | Meaning |
| ---: | --- |
| `0` | The command completed and the selected gate passed. |
| `1` | An active assessment gate failed (`no_go`). |
| `2` | The assessment or requested environment check is incomplete or blocked. |
| `3` | CLI input, project configuration, catalog, target, or state-location configuration is invalid. |
| `4` | An adapter, provider, or execution protocol failed. |
| `5` | Policy denied the requested operation or stopped it at a capability/change budget. |
| `6` | An internal invariant, canonical-state operation, serialization, or output operation failed. |
| `7` | The operation was cancelled. |
| `8` | A candidate patch was produced but failed independent acceptance. |

Codes are never borrowed directly from a child tool. An adapter maps its own
process behavior into the scanner protocol; the scanner then decides whether
the tool completed, found something, was partial, was unsupported, or failed.
Requesting a live adapter without the explicit `--mode verify-local` capability
grant returns `5` before the adapter runtime is resolved or invoked.
Registry revocation, lifecycle, and trust-policy denials also return `5` before
the runtime is selected. A malformed registry is a configuration failure.

Errors are written to stderr with their stable class, for example:

```text
error [PRC-EXIT-5]: adapter ... is not authorized by an applicable assertion
```

Machine-readable reports remain on stdout. A completed assessment that does not
pass is still a valid report and normally has no stderr error.

`prc version` emits a human-readable build identity. `prc version --format json`
emits `prc.version/v0.1`, including the semantic scanner version, exact source
revision, reproducible source timestamp, and Go toolchain. Development builds
report `unknown` for identity fields that were not injected by the release
builder; consumers must not treat that as release provenance.

`adapter fixture-validate` also uses `1` when its versioned fixture report is
valid but one or more recorded expectations or determinism checks fail. Invalid
suite structure, unsafe paths, manifest digest drift, or an attempted limit
increase use `3`; the command never executes the adapter.

## Scan policy

`prc scan` inspects the current directory; `prc scan /path/to/project` inspects
one explicit project. Scan options may appear before or after that path. The
equivalent advanced `--target` flag remains available, but it cannot be combined
with the positional project path.

`prc quick` selects the 18-assertion `prc/quick` local profile. `prc scan`
selects the 40-assertion core local profile. `prc verify` selects that same
core profile plus the exact bundled Gitleaks manifest in `verify-local` mode.
It owns the bundled catalog, profile, execution mode, and adapter selection; it rejects AI and
custom-adapter flags. The OCI plan always uses `--pull=never`, no network, a
read-only root, removed capabilities, resource limits, and a sealed read-only
snapshot, so the reviewed digest must already exist in Docker or Podman.
`prc full codex` and
`prc full claude` select the core profile plus advisory AI review of all 9,356
reviewed nondeterministic controls. The 686 reviewed deterministic controls are
never handed to AI for a verdict; supported exact programs use authoritative
collectors and all other programs remain honestly Blocked. `full` uses deep review, four parallel provider
batches, and requests one primary subagent per rule plus one independent skeptical
subagent per batch. Provider output cannot independently attest those internal
invocations, so orchestration is never treated as evidence. Codex full review
also selects `xhigh` reasoning. `quick`
rejects AI-provider flags, and both aliases reject a profile override so their
meaning cannot be silently changed. Every mode still includes all 10,042
controls in the complete report.

`prc setup [project]` is the short first-run preflight. It validates the target
and bundled catalog and reports optional provider availability without running
project code, a provider, or a container. `prc report` opens the newest private
scanner report; `prc report path` prints its path and `prc report list` lists
bounded recent results. `prc update` is the only ordinary command that checks
the npm registry; it never installs anything or runs in the background.

`prc ci [project]` is a fixed, local-only alias for SARIF stdout with no HTML
report. It rejects AI and output-format overrides. Its assessment exit code is
unchanged; callers may use ordinary scan options such as a configuration or
profile-independent evidence input, but must not weaken the release gate with
`--exit-policy never`.

Human output creates one detailed standalone HTML report by default. The file is
created privately outside the target, its absolute path is printed, and an
existing file is never overwritten. `--report PATH` chooses a new file;
`--no-report` disables the automatic report. JSON, Markdown, HTML, SARIF, and
JUnit stdout formats do not create an extra report unless `--report` is also
given. Report creation is output generation, not remediation: `scan` never
invokes `fix`, `remediate`, or a write-capable target process.

Human output uses `--color auto` by default. `--color always` and
`--color never` are explicit alternatives. A nonempty `NO_COLOR` environment
variable and `TERM=dumb` disable color. Color is never the only meaning: every
line keeps a symbol and a Pass, Fail, Blocked, Manual, N/A, or Error word.
Machine formats reject the human-only `--color` option and never contain ANSI
bytes. Untrusted paths and messages are escaped to one printable line before
terminal rendering.

Every complete scan includes all 10,042 registered controls. Broad controls
without complete proof remain `needs_review`; narrow passing assertions produce
only `partially_verified`. `prc full codex|claude` and the equivalent
`--ai codex|claude` form acknowledge screened remote source processing. The advanced form uses
`--review-provider` plus `--allow-remote-source-processing`. Those providers add
strict-schema advisory candidates and cannot change the authoritative
disposition. Cited paths and lines are snapshot-location validated, while the
AI claims remain explicitly `advisory_unverified`. Provider launch,
timeout, secret-screening, or protocol failures return `4`. See
[safe AI control review](ai-control-review.md).

`prc full codex --plan` and `prc full claude --plan` perform the same source
screening, control selection, and batching without resolving or starting the
provider and without creating resume state. Full review defaults to a maximum
of 1,500 batches and a 24-hour whole-run deadline. `--review-max-batches` and
`--review-max-duration` may change those bounds within their validated ranges.
They do not promise a token or whole-run money ceiling.

Human full-review output includes a live completion bar, completed batch and
rule counts, exact active top-level provider jobs, requested inner-review
slots, elapsed time, and a rough ETA. Codex token usage and Claude's reported
cost update only after completed batches. Provider events do not expose a
trustworthy live state for every inner subagent, so the CLI labels those counts
as requested rather than confirmed running agents. The default terminal keeps
only the six highest-priority local problems and eight highest-priority AI
reviews. `--details` prints every local check and completed AI review; advanced
`--review-details` expands only AI reviews. Every result remains in the HTML
report either way.

The optional root `.prcreviewignore` file narrows remote AI context only. Each
line is an exact `relative/file | reviewed reason`; directories, globs,
traversal, missing files, symlinks, and files already omitted by remote-review
safety rules are rejected. The file and every target remain in the local
inventory. Accepted targets are rehashed before omission and appear as explicit
sealed task limitations. This mechanism cannot change authoritative control
results.

The root `.prcignore` file is a narrow inventory escape hatch for reviewed
non-source directories that would otherwise cross the size guard. It uses exact
`relative/directory | reason` entries, never globs, and refuses unsafe paths,
symlinks, missing directories, and directories containing recognized project
source, configuration, deployment, CI, documentation, environment, or security
files. Accepted omissions are explicit inventory facts and part of the scan
identity; they are not evidence that omitted data is harmless.

An advanced scan may import one offline deterministic evidence bundle with
`--evidence-bundle`, `--evidence-trust-store`,
`--evidence-policy-signature`, and `--evidence-signature`. The four options are
atomic: supplying only some is configuration error `3`. A cryptographic trust,
scope, digest, time, catalog, inventory, program, authority, or evidence failure
is policy denial `5`. The scanner never executes bundle code and rejects a
template already evaluated by a built-in collector. See
[signed authoritative evidence bundles](authoritative-evidence-bundles.md).

For several authorities, `--evidence-set FILE` replaces those four options.
The set contains at most one alphabetically ordered bundle per authority and
references only non-symlink sibling files. The scanner verifies every bundle
before attaching any result. The two forms are mutually exclusive.

`prc coverage` prints five separate measures: reviewed rule routing, exact
predicate coverage, advisory AI review support, built-in collector coverage,
and signed-import support.
`prc coverage --format json` emits the same data as
`prc.automatic-coverage/v0.1`. AI support must never be reported as a verified
result. Signed-import support must never be reported as built-in collection or
as evidence actually observed for a particular scan.

`prc evidence requirements` authenticates the exact-program catalog and the
embedded collector-capability manifest, then exports the producer contract for
all selected clauses. `--authority`, `--control`, and `--collector-status`
filter it without weakening a clause. JSON output uses
`prc.evidence-requirements/v0.1`. The command only describes required evidence;
it does not collect or evaluate it.

`prc evidence verify-set --set FILE [PROJECT]` inventories the project and
verifies every referenced bundle and signature as one atomic set. An invalid
set returns policy denial `5` and emits no success document. JSON success output
uses `prc.evidence-set-verification/v0.1` and keeps cryptographic verification
separate from Pass, Fail, Not Applicable, and Blocked counts. It is a preflight,
not a readiness verdict.

If a later provider batch fails, every earlier schema-checked batch remains in
the private resume cache and is included in a report marked `partial`. The
report is written before code `4` is returned. Repeating the same command
revalidates and reuses matching completed batches, then sends only unfinished
work. A partial advisory review cannot make a control Pass.

The default `--exit-policy profile` maps `profile_satisfied` to `0`, `no_go` to
`1`, and incomplete, blocked, or manual-evidence states to `2`. A policy or
budget terminal state maps to `5`.

`--exit-policy no-go` does not convert incomplete evidence, blocked execution,
or scanner errors into success. It is retained for gate-policy compatibility as
profiles gain more non-no-go completed-failure states.

`--exit-policy never` is an explicit report-generation override: after a scan
has completed and emitted its canonical result it returns `0` regardless of the
terminal assessment state. It does not alter that state, rewrite a finding, or
hide a pre-report configuration, adapter, policy, cancellation, or internal
error. Do not use `never` as a release gate.

## Remediation policy

One-shot `remediate` and `remediate-proposal` commands return `0` only for an
accepted candidate and `8` for a fully reported but rejected candidate. A
request that exceeds an allowed remediation class, protected-path rule, or
configured capability budget returns `5`.

The `fix` loop returns `0` only for `profile_satisfied`, `8` when candidate
acceptance stops the loop, and `5` when policy or budget stops it. A configured
provider launch, timeout, or protocol failure is recorded as `provider_failed`
and returns `4`; a valid provider
`unable` or `needs_escalation` response is recorded as `provider_stopped` and
returns `2`. Explicit caller cancellation remains code `7` even though its
failure evidence is preserved. Otherwise the embedded assessment gate maps through the scan table,
so machine work ending with manual or unavailable evidence is `2`, not a
production-readiness claim.

`doctor` returns `2` when a requested required capability is unavailable. Its
JSON report still distinguishes each passing, warning, and failing probe.

## MCP stdio process

`mcp serve` returns `0` after a clean stdin close. Invalid locked paths,
configuration, catalog, or profile fail startup with `3`; an unreadable or
oversized stdio transport message terminates with `4`. Individual valid
JSON-RPC protocol and tool errors are returned on stdout using MCP error
objects, so the process can continue serving later messages without converting
an error into a readiness result. See the
[read-only MCP agent integration](mcp-agent-integration.md).
