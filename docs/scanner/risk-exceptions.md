# Signed risk exceptions

A risk exception is an accountable human decision about one exact failed
finding. It is not a scanner result, cannot turn a failure into a pass, and
does not change the run's terminal state. The scanner has no command that
creates, approves, or signs an exception.

The `prc.risk-exception/v0.1` record binds:

- the immutable run, inventory, profile, project, commit, artifact, and
  environment scope;
- the exact finding ID, stable fingerprint, assertion, and controls;
- distinct requester, risk owner, and independent reviewer identities;
- likelihood, impact, rationale, and worst credible outcome;
- compensating controls with content-addressed evidence references;
- monitoring, containment response, remediation owner, plan, and due date; and
- UTC approval and automatic expiry no more than 366 days apart.

Verification rejects a missing or non-failing finding, scope drift,
self-approval, duplicate reviewers, future approval, expired decision,
unbounded duration, unsigned edits, revoked keys, or a key without the
`risk-exception` scope. The signature must have been issued after approval and
before verification.

## Verify against immutable history

The bound run must already exist in the private scanner state store. Loading it
through the history store revalidates its nested plan, inventory, findings,
evidence, adapter executions, and content identities before exception
verification:

```bash
prc exception verify \
  --file /path/to/PRC-EXC-001.yaml \
  --state-dir /path/to/private-state \
  --trust-store /path/to/risk-owner-keys.yaml \
  --signature /path/to/PRC-EXC-001.signature.yaml \
  --verified-at 2026-08-23T13:00:00Z \
  --format json
```

The output disposition is `accepted_risk_exception`, but its `gate_effect`
always states that the finding remains failed and the scanner terminal state
is unchanged. A release process may present that signed decision to accountable
people; general agent remediation and scanner policy cannot consume it as a
Pass or silently suppress the finding.

The machine contracts are `risk-exception.schema.json` and
`risk-exception-verification.schema.json`. Compensating evidence references
are digests, not proof by themselves; the approving workflow remains
responsible for retaining and independently authenticating the referenced
evidence.
