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

`adapter fixture-validate` also uses `1` when its versioned fixture report is
valid but one or more recorded expectations or determinism checks fail. Invalid
suite structure, unsafe paths, manifest digest drift, or an attempted limit
increase use `3`; the command never executes the adapter.

## Scan policy

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
provider launch, timeout, or protocol failure returns `4`; a valid provider
`unable` or `needs_escalation` response is recorded as `provider_stopped` and
returns `2`. Otherwise its embedded assessment gate maps through the scan table,
so machine work ending with manual or unavailable evidence is `2`, not a
production-readiness claim.

`doctor` returns `2` when a requested required capability is unavailable. Its
JSON report still distinguishes each passing, warning, and failing probe.
