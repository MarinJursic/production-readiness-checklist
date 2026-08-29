# Control classification strength gate

This document records the third review gate added after the frozen classification
methodology was used for the primary and skeptical reviews. The frozen methodology
file remains unchanged so all historical packet and review digests stay verifiable.

Every deterministic proposal that survived primary and skeptical review receives
one additional rule-by-rule strength review. The reviewer compares the proposed
atomic clauses with the complete original statement and constructs a concrete
pass-while-broken counterexample. A clause that only checks for an acceptance
contract, repeats the original statement, or accepts a provider conclusion does
not preserve the control. The strength review must replace it with exact evidence
conditions or reclassify the complete control as nondeterministic.

The review is stored as a separate digest-bound overlay under
`research/control-classification/strength/`. This preserves the historical primary
and skeptical decisions while making the strength decision authoritative in the
generated final corpus.

## Compile-feasibility rule

A strengthened deterministic clause must be expressible as a closed predicate over
supported typed raw facts and independently sealed scanner or policy inputs. A
provider cannot supply both expected and observed sides of a comparison and cannot
supply an opaque conclusion such as `clause_satisfied`.

The predicate language uses a small set of reusable audited operators. Missing,
stale, conflicting, incomplete, wrongly typed, or unsupported evidence produces
`Blocked`, never `Pass` or an invented `Fail`. If the full original promise cannot
be represented honestly with reusable operators and an authoritative evidence
contract, the strength overlay reclassifies the control as nondeterministic.

## Rebuild and validation

Run the complete sequence documented in
[`research/control-classification/README.md`](../../research/control-classification/README.md).
`generate-final` validates the strength overlay before writing the final corpus;
all downstream bindings, contracts, human-readable docs, and program catalogs must
then be regenerated and checked for staleness.
