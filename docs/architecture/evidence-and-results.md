# Evidence and result model

Production readiness is a set of proof obligations, not a count of checked boxes.
The scanner keeps normative intent, executable assertions, observations, evidence,
and policy decisions as separate versioned records.

## Core records

| Record | Purpose |
| --- | --- |
| Control objective | Human-readable outcome and rationale |
| Assertion | Atomic proof obligation with explicit applicability and evidence requirements |
| Implementation | Versioned producer and evaluator for an assertion |
| Observation | Fact reported by a producer before assertion evaluation |
| Evidence | Immutable bytes or external reference supporting an observation |
| Assertion result | Applicability, execution, and assessment for one assertion |
| Finding | Prioritized consequence and remediation eligibility |
| Gate | Policy decision derived from results |
| Exception | Scoped, owned, expiring risk acceptance |
| Fix contract | Authorized change scope and independent acceptance criteria |

One control objective can map to multiple assertions. Multiple implementations
can support an assertion for different ecosystems. An implementation version is
part of result identity.

## Proof obligation

An assertion declares before execution:

- the property to prove;
- applicability conditions;
- acceptable evidence types and minimum authority;
- target, artifact, release, and environment scope;
- freshness and invalidation rules;
- evaluation logic;
- failure severity and gate behavior;
- whether human or external evidence is required; and
- remediation class, if any.

`evidence_required` and `evidence_observed` are distinct fields. Missing evidence
does not prove that a property is absent, and a located file does not prove that
the file's intended behavior occurred.

## Independent result axes

Applicability:

- `applicable`
- `not_applicable`
- `undetermined`

Execution:

- `not_run`
- `completed`
- `blocked`
- `error`

Assessment:

- `pass`
- `fail`
- `unknown`
- `manual_review`
- `stale`
- `conflicting`

This prevents an adapter error, missing evidence, manual decision, and verified
failure from collapsing into the same label.

## Canonical findings and correlation

Run result v0.6 emits a `prc.finding/v0.1` record for every failed assertion.
It does not create findings for missing access, execution errors, manual review,
unknown applicability, or Not Applicable results; those remain visible on their
independent assertion-result axes.

Each finding records:

- a content-addressed ID for the exact inventory, summary, evidence, and policy
  classification;
- a stable fingerprint over the assertion, logical subject, and normalized
  source locations;
- the mapped assertion and control IDs;
- the exact repository or declared-project subject and inventory digest;
- severity, gate, and remediation class;
- only validated repository-relative locations; and
- the complete set of evidence IDs from its failed assertion result.

The fingerprint survives new evidence and explanatory wording changes so a
failure can be correlated across scans. The content ID changes when the exact
finding envelope changes. A declared project ID supplies stable cross-checkout
subject identity; without project configuration, correlation is scoped to the
scanner's repository target name.

The state database indexes both identities, normalized locations, and evidence
links while immutable run JSON remains authoritative. SARIF is derived from
these canonical findings. It no longer guesses findings or file locations from
arbitrary evidence source strings.

Frozen `prc.run/v0.5` and older schemas remain available for archived results.

## Evidence envelope

Every evidence item records:

- schema version and evidence identifier;
- target, commit, artifact digest, release, and environment as applicable;
- producer, implementation, tool, and policy versions;
- collection time, validity interval, and invalidation keys;
- source authority and collection method;
- media type, content digest, location, and optional redacted preview;
- sensitivity, retention, and access classification;
- scope, limitations, and assumptions;
- verification method and outcome; and
- relationships to supporting or contradictory observations.

Evidence content is stored by digest. Reports refer to evidence identifiers and
redacted previews rather than copying sensitive content indiscriminately.

The durable state layer keeps immutable JSON as the canonical record and adds a
transactional SQLite query index for runs, results, evidence metadata, inventory
files, sourced facts, relationships, and audit events. It verifies canonical
identities before indexing, uses strict tables and foreign keys, and repairs
derived rows by re-indexing an immutable run. Indexed record paths are relative
and must resolve inside the state root. See [durable state and run
history](../scanner/state-and-history.md).

The [diff-aware invalidation model](../scanner/diff-and-invalidation.md) binds
new plans to exact profile and assertion definitions, maps native
implementations to their repository inputs, and refuses to carry evidence to a
different target identity without revalidation.

## Authority is assertion-specific

Authority is not a universal numeric confidence score. A repository configuration
may be authoritative for declared configuration syntax while remaining weak
evidence of what is deployed. A signed production observation may prove deployed
state without proving recovery or user-visible behavior.

The evaluator compares the evidence's authority, scope, identity, and freshness
to the assertion's declared proof obligation. Model confidence cannot replace
this comparison.

## Exceptions and gates

An exception contains an owner, rationale, affected assertion and target scope,
compensating controls, approval identity, creation time, expiry, and invalidation
conditions. Expired or out-of-scope exceptions do not suppress findings.

Gates consume explicit assertion results and exceptions. They do not consume a
readiness percentage.
