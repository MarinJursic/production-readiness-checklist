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

- finding, assertion, target, and baseline result identifiers;
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

An agent does not validate its own patch. A second model can challenge a candidate,
but only the independent deterministic acceptance sequence can approve it.

## Termination

The loop stops when:

- the selected profile is satisfied;
- only manual, external, prohibited, or inapplicable work remains;
- a no-go finding cannot be remediated within policy;
- an attempt makes no measurable improvement;
- the same root cause repeats without progress;
- any capability or change budget is exhausted; or
- the environment cannot execute the verifier safely.

The general loop never deploys and never accepts residual risk.

## Implemented R1 pilot

The current CLI implements this lifecycle only for `PRC-A-CORE-014`. Its trusted
fixer appends one final line-feed byte to exact allowlisted source paths in a new
external workspace. The acceptance audit verifies the raw tree, content hashes,
permission modes, protected paths, budgets, target result, and baseline passing
results. See the [R1 remediation guide](../scanner/remediation.md) for the exact
command and limitations. Agent-authored R2 changes and broader autonomous loops
remain unimplemented.
