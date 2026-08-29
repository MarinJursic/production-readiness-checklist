# Evidence-backed production convergence

This document is the implementation roadmap for turning one scan into a safe,
repeatable path toward production readiness. It also sets an honest limit: no
scanner or AI can prove that software has no defects, no future vulnerability,
and no operational surprise. The defensible goal is narrower and useful:

> For the chosen project, environment, release, and evidence time, every
> applicable critical control has current authoritative evidence, no unresolved
> critical finding remains, every accepted change survived independent checks,
> and every remaining unknown is explicit.

That result is much stronger than an AI saying that a repository "looks ready."

## Current measured state

The registry contains 10,042 controls. Rule-by-rule primary, skeptical, and
semantic-strength reviews classify 686 as deterministic and 9,356 as
nondeterministic. The deterministic corpus contains 686 control bindings and
765 exact program templates.

The normal repository scan currently runs a much smaller trusted core: 40 local
checks. The exact runtime now authenticates collector registrations, seals
scope parameters before collection, evaluates provider facts through the pure
predicate engine, and aggregates clause results without accepting serialized
or provider-chosen verdicts. One repository collector is connected end to end:
it can prove `PRC-36-004` for a supported Node documentation layout. One
embedded capability manifest is the source of truth for both the generated
catalog and runtime registration. All other
missing, unsupported, incomplete, changed, or unclear evidence remains Blocked.
The report separates exact programs attempted, passed, failed, Not Applicable,
and deterministic controls still blocked. Valid exact evidence documents are
retained in the content-addressed run so an authoritative result can be replayed
without rereading a changed target. The 9,356 nondeterministic controls
can receive bounded AI advice, but AI cannot change their authoritative
disposition.

This distinction must remain visible in the CLI, report, schemas, tests, and
README. Catalog coverage is not execution coverage.

## What the quality-first AI path now does

`prc full codex` and `prc full claude` use deep review by default:

1. The scanner inventories the target and creates a bounded, secret-screened
   snapshot outside it.
2. It divides the 9,356 nondeterministic controls into sealed batches.
3. A separate primary subagent reviews each rule.
4. One independent skeptical subagent per batch searches for unsupported
   claims, false Not Applicable decisions, missed risk, contradictions, and
   generic advice.
5. The coordinator must return priority, risk, ordered remediation,
   independent verification, evidence still needed, and the strongest
   challenge for every rule.
6. Strict schemas and snapshot line checks reject incomplete or malformed
   output. The semantic claim is still marked advisory and unverified.
7. Four provider batches run in parallel; Codex uses `xhigh` reasoning.
8. Every completed batch is sealed in a private resume cache. A later failure
   produces a partial report and exit code 4. The next identical run reuses
   valid completed batches and sends only unfinished work.

The report keeps verified scanner findings separate from AI priorities. AI
priorities are sorted and actionable, but never presented as proven facts.

## Target architecture

The convergence system should be an evidence pipeline, not one large prompt:

```text
scope -> applicability -> safe inventory -> collectors -> exact evaluators
      -> AI reviews -> domain synthesis -> change proposals -> isolated checks
      -> rescan -> release gate -> signed evidence -> explicit residual risk
```

### 1. Scope and applicability

Before checking a rule, record the project type, deployable units, environments,
data classes, public surfaces, regulatory scope, critical dependencies, and
authenticated project policy. Applicability must be its own result with a
reason and evidence. Missing or conflicting facts yield Blocked, never an
invented Pass or Not Applicable.

Acceptance criteria:

- the scope document is versioned, schema-checked, and content-addressed;
- every control has an applicability result and evidence trail;
- default folder names, languages, cloud providers, and team structures are not
  treated as universal requirements;
- policy thresholds come from authenticated project policy, not model output;
- scope changes invalidate dependent evidence.

### 2. Deterministic collector and evaluator runtime

Finish the runtime path for the 765 exact program templates. Each program must
bind only to typed raw evidence from the authority named in its contract.
Collectors should be small and isolated: repository facts first, followed by
executed checks, artifacts, structured records, environment evidence, and
external evidence. Missing, stale, conflicting, incomplete, or wrongly typed
evidence must yield Blocked.

Implementation order by current authority mix:

1. repository evidence (26 clauses);
2. artifact evidence (114 clauses);
3. structured records (363 clauses);
4. safely executed checks (161 clauses);
5. environment evidence (92 clauses);
6. authenticated external evidence (9 clauses).

Acceptance criteria:

- `fullscan.Attach` binds runtime parameters, calls the exact evaluator, and
  records evidence/result digests for each deterministic control;
- every collector has pass, fail, missing, stale, conflict, malformed,
  adversarial, and determinism fixtures;
- collectors cannot follow unsafe links, escape roots, execute unpinned code,
  or silently widen limits;
- a Pass can be reproduced from its sealed evidence without the target;
- report totals separate executed, passed, failed, and blocked deterministic
  programs.

### 3. Domain-aware nondeterministic review

Per-rule agents find local issues, but 9,356 isolated answers can repeat the
same root cause and conflict. Add a second durable layer of domain agents for
security, reliability, delivery, data, privacy, observability, operations,
performance, accessibility, and product risk. They should consume sealed
per-rule outputs, cluster shared causes, preserve disagreement, and build a
dependency-aware plan. A final critic should try to falsify the proposed plan.

Acceptance criteria:

- no domain agent receives write, shell, network, browser, or source-reading
  access during review;
- root-cause clusters retain all source control IDs and citations;
- contradictions are visible rather than averaged away;
- the final plan is sorted by severity, dependency, blast radius, and effort;
- each work item states expected evidence and a separate verification method;
- benchmark cases measure false positives, missed findings, false N/A results,
  citation misuse, prompt injection, disagreement, and unusual layouts.

### 4. Durable execution graph

Represent every inventory, collection, evaluation, review, synthesis, proposal,
verification, and rescan step as an immutable job with input and output digests.
Only retry idempotent work. Make concurrency explicit and bounded by provider,
CPU, memory, rate, and cost limits.

Acceptance criteria:

- interruption never loses a valid completed job;
- resume refuses changed tasks, schemas, models, executables, evidence, or
  policy unless a new run is created;
- the CLI shows completed, running, cached, blocked, failed, token, time, and
  cost counts without flooding the terminal;
- failure of one domain produces a partial, clearly incomplete report rather
  than discarding unrelated valid results.

### 5. Safe repair loop

Scanning and changing remain separate commands. Convert accepted findings into
small, dependency-ordered change proposals. Each proposal is applied to a fresh
isolated candidate, never directly to the original project. Use deterministic
fixers for exact mechanical changes and an agent only where judgment is needed.
The proposing agent must not be the only verifier.

Acceptance criteria:

- explicit user policy defines allowed paths, commands, dependencies, network,
  change size, iteration count, time, and cost;
- dependency installation is off by default; enabled installs use lockfiles,
  pinned registries, disabled lifecycle scripts where possible, provenance,
  vulnerability review, and an isolated environment;
- protected files, secrets, auth, deployment policy, test weakening, skipped
  checks, and capability expansion fail closed;
- independent build, test, static analysis, security, and regression checks run
  before a candidate can be accepted;
- every accepted change is rescanned from fresh bytes;
- rollback data and a human-readable change journal are always produced.

### 6. Runtime and operational proof

Repository inspection cannot prove production behavior. Add opt-in adapters for
staging smoke tests, health and readiness probes, resource limits, service-level
indicators, alerts, traces, logs, metrics, backup restore drills, dependency
failure tests, and controlled resilience experiments. Production access must be
read-only by default and separately authorized.

Acceptance criteria:

- every adapter has a signed manifest, least privilege, explicit environment,
  hard timeout, output limits, and revocation path;
- readiness and liveness are not confused;
- service objectives and error budgets are project policy inputs;
- a backup claim cannot pass without a recent restore result;
- stale runtime evidence blocks affected controls.

### 7. Supply-chain and vulnerability evidence

Add or strengthen SBOM generation, build provenance, signature verification,
dependency and secret scanning, vulnerability reachability, exploitability,
known-exploited status, license policy, and release-artifact checks. Scanner
release artifacts should be reproducible and signed.

Acceptance criteria:

- findings identify the exact component, version, dependency path, artifact,
  advisory source, and evidence time;
- vulnerability priority uses authoritative advisory data plus known exploited
  and exploit-probability signals without hiding uncertainty;
- VEX-style not-affected statements require signed, reviewable justification;
- the npm package has no install scripts unless strictly required, a minimal
  published file set, locked dependencies, provenance, SBOM, signatures, and
  release tests in a clean environment.

### 8. Evaluation and release gates

Treat scanner quality as a measured product. Maintain labeled reference
projects and adversarial fixtures. Run offline regression evaluation before any
rule, prompt, model, collector, schema, or threshold change is released.

Acceptance criteria:

- publish precision, recall, false-N/A, blocked, reproducibility, duration,
  memory, disk, token, and cost measurements by domain;
- use separate training/development and release evaluation cases;
- require zero critical benchmark misses and explicitly approved thresholds for
  other severities;
- compare models and prompts on the same sealed snapshots;
- never call the scan complete solely because an agent says it is complete.

## Production-ready gate

A run may say `ready_for_selected_scope` only when all of these are true:

- scope and applicability are complete and current;
- every applicable release-blocking deterministic control passed with sealed
  authoritative evidence;
- every applicable release-blocking nondeterministic control has an approved
  accountable decision and required evidence;
- there are no unresolved critical findings or contradictory critical results;
- the accepted candidate passed independent build, test, security, and runtime
  gates, then passed a fresh rescan;
- evidence freshness limits are satisfied;
- residual risks, waivers, owners, expiry dates, and unsupported areas are
  explicit;
- the final evidence bundle is content-addressed and signed.

The result must always name its selected scope and evidence time. It must never
use the unqualified phrase "no issues."

## Primary references

- [NIST Secure Software Development Framework 1.1](https://csrc.nist.gov/pubs/sp/800/218/final)
- [SLSA build requirements 1.2](https://slsa.dev/spec/v1.2/build-requirements)
- [NIST OSCAL assessment results](https://pages.nist.gov/OSCAL/learn/concepts/layer/assessment/assessment-results/)
- [OpenAI Codex model guidance](https://developers.openai.com/api/docs/guides/latest-model)
- [OpenAI agent safety guidance](https://developers.openai.com/api/docs/guides/agent-builder-safety)
- [OpenSSF Scorecard checks](https://scorecard.dev/)
- [OSV vulnerability schema and service](https://ossf.github.io/osv-schema/)
- [CISA Known Exploited Vulnerabilities Catalog](https://www.cisa.gov/known-exploited-vulnerabilities-catalog)
- [FIRST Exploit Prediction Scoring System](https://www.first.org/epss/)
- [CycloneDX SBOM and VEX](https://cyclonedx.org/)
- [Sigstore documentation](https://docs.sigstore.dev/)
- [Kubernetes probes](https://kubernetes.io/docs/concepts/configuration/liveness-readiness-startup-probes/)
- [OpenTelemetry signals](https://opentelemetry.io/docs/concepts/signals/)
- [Google SRE monitoring distributed systems](https://sre.google/sre-book/monitoring-distributed-systems/)
- [Google SRE error budgets](https://sre.google/workbook/error-budget-policy/)
