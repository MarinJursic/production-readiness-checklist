# Bounded remediation contract

Remediation is a separate, policy-controlled phase. A failed assertion does not
automatically grant an agent permission to change a repository.

## Remediation classes

| Class | Capability | Default policy |
| --- | --- | --- |
| R0 | Explain or suggest only | Allowed |
| R1 | Deterministic, behavior-preserving edit | Allowed when the implementation is trusted |
| R2 | Agent-authored repository change | Experimental and isolated |
| R3 | Dependency or build behavior change | Explicit opt-in and stronger verification |
| R4 | Infrastructure or deployment-definition change | Human authorization required |
| R5 | External staging-system mutation | Separate connector policy and approval |
| R6 | Production, destructive, legal, financial, or risk decision | Prohibited in the general loop |

## Fix contract

Every attempt is bound to a signed or hashed contract containing:

- exact finding content ID, stable finding fingerprint, assertion, target, and
  baseline result identifiers;
- desired postcondition;
- permitted and protected paths;
- permitted commands, tools, network, secrets, and external systems;
- maximum changed files, lines, bytes, attempts, time, and cost;
- required new or updated regression tests;
- invariant and regression checks;
- policy and scanner inputs that must remain unchanged; and
- the deterministic acceptance sequence.

The agent receives the minimum context needed for this contract. Unrelated source,
evidence, credentials, and previous agent memory are excluded.

## Candidate lifecycle

1. Capture an immutable baseline and reproduce the finding.
2. Create a fresh isolated workspace or worktree.
3. Apply one deterministic fix or run one provider task.
4. Parse the result as untrusted structured output.
5. Reject path, capability, policy, threshold, baseline, or suppression violations.
6. Run contract-specific tests and the original finding implementation.
7. Run impacted regression checks and anti-gaming validation.
8. Rescan affected assertions from fresh evidence.
9. Accept only when the original finding closes and no protected result regresses.
10. Preserve the candidate, transcript, evidence, and rejection or acceptance reason.

The current v0.8 remediation-run record adds a scanner-validated attempt chain,
binds any launched R2 test verification to its exact candidate inventory, and
records the scanner-owned wall-clock ceiling for the complete loop.
An accepted attempt must link to exactly one content-addressed candidate; an
agent attempt must link to exactly one successful provider execution or failed
provider invocation; a pre-candidate
policy rejection records its exact safe scanner reason without inventing a
candidate; and the accepted digest chain must end at the reported final
inventory. Missing, duplicate, reordered, or cross-linked records invalidate
the run before its content ID is computed.

An agent does not validate its own patch. A second model can challenge a candidate,
but only the independent deterministic acceptance sequence can approve it.

## Termination

The loop stops when:

- the selected profile is satisfied;
- only manual, external, prohibited, or inapplicable work remains;
- a no-go finding cannot be remediated within policy;
- an attempt makes no measurable improvement;
- the same root cause repeats without progress;
- any capability, change, attempt, or duration budget is exhausted; or
- the environment cannot execute the verifier safely.

The general loop never deploys and never accepts residual risk.

## Implemented R1 pilot

The current CLI implements this lifecycle for final-newline and restrictive-mode
R1 assertions. Its trusted fixers make exact byte or mode transformations in a
new external workspace. The acceptance audit verifies the raw tree, content
hashes, permission modes, protected paths, budgets, target result, and baseline
passing results. The bounded `fix` loop can repeat these registered operations in
sibling candidates, rescan after every accepted attempt, enforce cumulative
attempt, file, and line limits, and report why every unresolved item remains.
See the [R1 remediation guide](../scanner/remediation.md) for the exact commands
and limitations.

## Implemented R2 proposal pilot

The [read-only provider layer](agent-providers.md) can produce bounded Codex or
Claude Code proposals. The separate `remediate-proposal` command validates and
parses one proposal, applies it to a fresh external candidate without invoking
the provider, and runs raw-tree, exact-byte, mode, budget, target-result,
regression, and source-integrity audits. The provider never approves its own
work. Both one-shot remediation paths bind the exact verified baseline finding,
project configuration, protected paths, and file, line, and attempt ceilings
into a v0.3 fix contract. The default loop does not invoke providers. When
explicitly enabled,
`fix --provider` composes the read-only provider and scanner-owned proposal
paths only for a missing-test assertion with one bounded source input and new
test-file allowlist. The R2 contract also preserves the exact triggering
finding ID from the sealed provider task and requires both its content ID and
stable fingerprint to match a freshly reproduced failure. An R2 candidate cannot
be accepted without a separate digest-pinned OCI verifier. The scanner, not the
provider or repository, selects the Go, Python, or plain-JavaScript test argv;
the read-only
candidate runs without network, privileges, secrets, or host-write access and
with bounded resources and output. General R2 repair planning, TypeScript
verification, merges, deployments, and releases remain unimplemented.
