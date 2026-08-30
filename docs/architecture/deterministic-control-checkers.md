# Deterministic control checkers

This document defines how a control that passed the primary review, independent
skeptical review, and semantic-strength audit becomes executable scanner logic.
Classification alone is not implementation, and a useful partial check is not
complete verification of the broader control.

## Exact program definitions

Every retained deterministic clause has a reviewed program definition under
`research/control-classification/program-specs/`. A definition names the raw
typed facts, independently sealed parameters, closed predicate, collector
contract, and four executable fixtures: pass, fail, blocked, and a
pass-while-broken counterexample that must fail.

The current validated corpus contains 686 deterministic controls and 765
clauses. All 765 clauses have executable predicates and all four definition
fixtures pass through the Go evaluator. This does not mean all 765 clauses have
collectors. The runtime ships the first reviewed repository collector for
`PRC-36-004`. It also accepts an offline, authority-scoped bundle for any exact
template when two different trusted keys approve the two ordered subjects: one
policy key signs programs and runtime inputs before collection, and one
evidence-authority key signs the complete bundle after collection. This does
not create evidence or replace built-in collectors; every unavailable or
unauthenticated source still remains Blocked.

The runtime uses an immutable exact collector-ID registry. It rejects missing,
duplicate, nil, and wrong-authority providers. The scanner seals scope, policy,
inventory, applicability, and freshness into the program before calling a
provider. A provider returns typed facts, never a verdict. Only the pure
predicate evaluator can produce Pass, Fail, or Not Applicable. Externally
constructed or serialized execution results cannot be attached as authoritative
because attachment accepts only runtime-authenticated in-process results.

Valid evidence is retained with the content-addressed run and bound back to its
program, control, clause, inventory, subject, authority, and observation time.
Run-state validation rejects missing, duplicate, unreferenced, or mismatched
replay evidence. The HTML report keeps the normalized facts in a collapsed
technical section.

The first collector is intentionally narrow. It reopens the exact hashed
`package.json` and Markdown files without following links, rejects changed or
oversized input and duplicate JSON keys, and recognizes public build and test
invocations only inside Markdown code. It never runs either script. A supported
positive case can Pass; missing, prose-only, malformed, oversized, changed, or
unsupported evidence stays Blocked instead of becoming a guessed Fail.

The generated `catalog/control-check-programs.json` contains the runtime form of
those definitions. The loader authenticates the catalog against the reviewed
binding catalog, control registry, program schema, and definition corpus. It
rejects duplicate JSON keys, stale digests, unknown fields, duplicate or missing
clauses, unsupported operations, undeclared facts or parameters, opaque verdict
fields, generic schema/profile delegation, and malformed or oversized input.

The predicate language has no shell commands, file access, network access,
regular expressions, dynamic functions, plug-ins, or code execution. It only
supports bounded boolean composition and reviewed operations over typed scalar,
set, map, timestamp, and directed-graph values.

Expected values use three separate trust lanes. Scanner-discovered subject and
object inventories use `scanner_inventory`; project thresholds and approved
values use `independently_authenticated_policy`; current publisher or issuing
body facts use `independently_authenticated_context`. A program compiler must
fill the declared lane before evidence is requested. Evidence providers cannot
supply or replace any sealed parameter.

For clauses whose decisive data is a large artifact graph, bounded execution,
or effective-environment snapshot, the predicate uses an exact manifest
relation. It binds the complete subject list, source/input-manifest digest,
pinned normalizer schema digest, and canonical typed raw-observation digest.
The expected observation digest comes from a different trust lane. The
normalizer may contain direct values and relations only—never a provider verdict,
score, or compliance summary. A missing field, case, object, page, relation, or
conflict makes the evidence incomplete and therefore Blocked.

## Binding invariant

Every confirmed deterministic control has exactly one full-control binding. The
binding is tied to the control revision and semantic digest, the reviewed
classification row and methodology, and one or more versioned clause checkers.
Every clause must pass for the control to pass.

A binding may reuse a reviewed checker family, but it may not reuse an outcome
from a different control, revision, subject, inventory, artifact, environment,
or evidence authority. Existing catalog assertions remain partial until they are
explicitly referenced by a complete binding whose clauses cover the whole rule.

## Outcomes

Clause checkers return one of three outcomes:

- `pass`: complete, current, authoritative evidence proves the exact clause;
- `fail`: complete, current, authoritative evidence proves the exact negative
  condition; or
- `blocked`: the checker cannot decide safely.

Not Applicable is decided separately by a versioned applicability predicate. An
unknown trigger is Blocked, not Not Applicable. A full control passes only when
applicability is known and every applicable clause passes. Any blocked clause
blocks the full control.

The following conditions always block instead of guessing: missing or wrong
authority, missing provider, denied capability, unsupported target, incomplete
inventory, stale or conflicting evidence, schema or digest mismatch, rate limit,
incomplete pagination, ambiguous result, and an unavailable implementation.

## Checker families

The implementation uses a small registry of reviewed, data-driven families:

- inventory facts;
- structured documents;
- package metadata;
- CI policy;
- container and infrastructure configuration;
- source syntax trees;
- artifact integrity;
- pinned analysis adapters;
- bounded execution evidence;
- environment evidence; and
- authenticated structured records.

Each family has a versioned implementation and immutable digest. A rule-owned
binding supplies only reviewed parameters. It cannot broaden the evidence source
or weaken the family contract.

## External evidence

Repository inspection is read-only and content-bound. External checks require an
explicit read-only provider with an allowlisted host, named credential handle,
bounded pagination and timeout, stable subject and tenant scope, and captured
query identity. The default local scan has no network or secret authority, so an
external clause is Blocked until such a provider is configured.

An offline producer can instead supply one strict
`prc.authoritative-evidence-bundle/v0.1`. The bundle is limited to one authority,
32 MiB, and 765 ordered unique entries. Its raw catalog and inventory digests
must match the current scan. The policy signature must predate observation; the
separate authority signature must postdate it. The runtime revalidates the full
reviewed program identity, exact predicate, parameter names and types, evidence
authority, subject, completeness, freshness, and applicability before calling
the pure evaluator. The run retains both verification records and every exact
program and evidence document. History validation reconstructs the canonical
bundle and replays each retained evaluation. See
[signed authoritative evidence bundles](../scanner/authoritative-evidence-bundles.md).

AI review is advisory. It may explain a blocked or nondeterministic rule and
point to possible evidence, but it is never an accepted evidence authority for a
deterministic Pass, Fail, or final Not Applicable result.

## Required tests

Every checker family must pass one shared conformance suite and its own fixtures.
The shared suite covers Pass, Fail, missing evidence, unsupported targets,
incomplete inventories, wrong authority, stale and conflicting evidence, stale
control and implementation digests, malformed and oversized inputs, symlinks and
path escape, denied capabilities, repeatability, read-only behavior, and a
pass-while-broken counterexample.

The Go fixture harness executes all four fixtures for every retained clause and
requires exact one-to-one coverage. Bindings are generated only from the
validated final classification corpus.
Validation fails when a confirmed deterministic control is missing, duplicated,
stale, mapped to a nondeterministic result, or bound to a checker that cannot
prove both directions of the whole clause.
