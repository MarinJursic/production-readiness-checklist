# Rule-by-rule classification audit

This directory records the explicit agent review of every active control as
`deterministic` or `nondeterministic`. Generated labels in
`catalog/control-contracts.json` are triage hints and are not accepted as reviews.

## Layout

- `packets/` freezes each rule's wording, revision, source heading, and semantic
  digest in bounded review batches.
- `primary/` contains one primary agent review for every packet.
- `skeptic/` contains an independent counterexample review for every deterministic
  proposal.
- `strength/` contains a third rule-by-rule review of every proposal that survived
  the first two gates. It must prove that the final clauses preserve the complete
  original statement and can compile to reusable typed predicates over raw facts.
- Final classifications and checker bindings are generated only after all three
  review gates pass. Historical primary and skeptical files remain unchanged for
  provenance; the strength overlay is authoritative when it confirms or reclassifies
  a previously deterministic control.

The review contract is in
[`docs/architecture/control-classification.md`](../../docs/architecture/control-classification.md).
The machine schema is
[`schemas/control-classification-review.schema.json`](../../schemas/control-classification-review.schema.json).

## Commands

Generate packets from the current registry:

```console
python3 scripts/control_classification_review.py generate-packets
```

Validate one completed packet:

```console
python3 scripts/control_classification_review.py validate-packet PACKET.json
```

Show audited progress:

```console
python3 scripts/control_classification_review.py progress
```

Validate exact one-to-one coverage after all primary reviews are complete:

```console
python3 scripts/control_classification_review.py validate-primary
```

Generate smaller skeptical-review packets containing only proposed deterministic
controls, then validate one or all completed skeptical reviews:

```console
python3 scripts/control_classification_review.py generate-skeptic-packets
python3 scripts/control_classification_review.py validate-skeptic-packet PACKET.json
python3 scripts/control_classification_review.py validate-skeptic
```

After the strength overlay contains exactly one digest-bound row for every
skeptically confirmed proposal, generate and validate the authoritative final
corpus:

```console
python3 scripts/control_classification_review.py generate-final
python3 scripts/control_classification_review.py validate-final
```

`generate-final` validates the complete strength overlay before it writes output.
It rejects missing, duplicate, stale, reordered, or digest-mismatched strength
rows and clauses that use a generic acceptance-contract or provider-conclusion
wrapper instead of naming the exact evidence conditions.

Regenerate and verify every downstream view after a final decision or clause
changes:

```console
python3 scripts/build_control_check_bindings.py generate
python3 scripts/build_control_check_bindings.py validate
python3 scripts/control_contracts.py generate
python3 scripts/control_contracts.py check
python3 scripts/build_control_classification_docs.py generate
python3 scripts/build_control_classification_docs.py check
python3 scripts/build_control_check_programs.py generate
python3 scripts/build_control_check_programs.py check
```

The program catalog reports predicate definitions separately from registered
collectors and end-to-end runnable controls. Missing implementation or authority
must stay `Blocked`; generated bindings never imply that an evidence provider
exists.

Validation fails closed when a rule is missing, duplicated, reordered, changed,
or reviewed under a stale registry or methodology. A partial deterministic check
does not make its broader parent rule deterministic.
