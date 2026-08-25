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
selects the 40-assertion core local profile. `prc full codex` and
`prc full claude` select the core profile plus advisory AI review of every
active control. `quick` rejects AI-provider flags, and both aliases reject a
profile override so their meaning cannot be silently changed. Every mode still
includes all 10,042 controls in the complete report.

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
