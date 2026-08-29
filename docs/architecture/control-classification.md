# Control classification contract

This contract defines how every active production-readiness control is classified
as `deterministic` or `nondeterministic`. The safe default is
`nondeterministic`. There is no target percentage for either category.

The classification is about the complete meaning of the broad control. A script
that checks only one useful part is partial evidence; it does not make the parent
control deterministic.

## Deterministic

A control is deterministic only when the same sealed evidence and the same
versioned policy inputs always produce the same applicability and result without
AI or human judgment. Every condition below is required:

1. The exact subject and assessment scope are known.
2. Applicability has an objective `true`, `false`, or `unknown` predicate.
3. The control is atomic, or every child and the complete aggregation are defined.
4. Pass and Fail are measurable in both directions.
5. Every threshold comes from a named standard, profile, or approved project
   policy. The checker does not invent a default.
6. Each evidence source has authority for the fact it proves.
7. `all`, `every`, `none`, `never`, and similar claims have a complete inventory.
8. Negative claims require evidence of absence. Missing evidence is not absence.
9. Scope binding, freshness, contradictions, and invalidation are defined.
10. Evaluation needs no AI interpretation, interview, legal decision, risk
    acceptance, user judgment, or expert choice.
11. The evaluator is bounded, read-only, versioned, and may return Blocked or
    Unknown instead of guessing.
12. Passing proves the whole control and failing proves it is not met. A
    one-directional warning rule is not complete verification.

A deterministic control may need a read-only environment connector, a signed test
result, or an approved numeric policy. Missing access makes the result Blocked; it
does not make the evidence up.

## Nondeterministic

A control is nondeterministic when any required part depends on contextual
interpretation, human authority, legal or ethical judgment, user research,
undefined words or thresholds, an unbounded universe, subjective applicability,
or a mixture containing a nondeterministic child.

AI may find relevant evidence, explain uncertainty, or draft review questions. AI
is never an evidence authority and cannot create a verified Pass or a final Not
Applicable result.

## Routes

Deterministic routes are:

- `local_static`: bounded parsing of sealed repository content.
- `artifact_verification`: hashes, signatures, SBOMs, provenance, or schemas.
- `bounded_execution`: an isolated test or analyzer with fixed inputs and outcomes.
- `external_readonly_query`: bounded CI, cloud, registry, monitoring, identity, or
  deployment evidence.
- `structured_record_validation`: the exact fact is the presence or state of an
  authenticated record.
- `deterministic_composite`: every child is deterministic and aggregation is
  complete.

Nondeterministic routes are:

- `contextual_judgment`: design, clarity, maintainability, proportionality, or
  similar meaning-based review.
- `accountable_human_decision`: approval, risk acceptance, or release authority.
- `specialist_or_legal_judgment`: legal, privacy, safety, accessibility, or other
  specialist interpretation.
- `empirical_protocol_undefined`: a study, drill, performance test, or sample has
  no fixed protocol and outcome.
- `contract_incomplete`: scope, applicability, threshold, term, or expected result
  is not sufficiently defined.
- `mixed`: at least one required child is nondeterministic.
- `unbounded_claim`: complete scope cannot be safely enumerated or observed.

## Rule-by-rule review

For each control, the reviewer must read the exact statement and its full heading
trail, preserve its strength, split independent promises, define applicability,
identify evidence authority, and perform two tests:

1. **Counterexample test:** Can the proposed checker Pass while a real violation of
   the control remains? If yes, the checker is partial and the broad control is not
   deterministic.
2. **Reviewer-variance test:** Can two informed reviewers reach different answers
   from identical evidence without either violating the contract? If yes, the
   control is nondeterministic.

Every control receives one primary review. Every proposed deterministic control
then receives a separate skeptical review. A disagreement stays nondeterministic
until it is adjudicated. Generated keyword labels are input hints only and cannot
count as either review.

## Deterministic implementation gate

Each confirmed deterministic control must have a complete atomic decomposition and
versioned checker bindings. Implementations should reuse a small set of reviewed
checker families. They must not copy one ad-hoc script per control when the same
safe parser or verifier can serve several controls.

Every binding requires Pass, Fail, Blocked, and, when allowed, Not Applicable
fixtures. It also requires unusual-layout, malformed, oversized, symlink,
path-escape, stale, conflicting, hidden-inventory, repeated-run, and read-only
tests relevant to its evidence family. Missing, malformed, stale, conflicting, or
truncated evidence never becomes Pass.
