# Control contracts

The checklist contains broad objectives. A broad objective is not automatically
an executable test. `catalog/control-contracts.json` gives every registered
control a compact machine-readable triage contract without pretending that a
generated classification is an expert-approved acceptance test.

Every contract records:

- whether repository, environment, human, or mixed evidence is expected;
- whether deterministic checking or AI advice may be useful;
- whether the wording appears compound, conditional, negative, absolute, or
  dependent on a project-specific threshold;
- the evidence authority that a future assertion may require;
- the proof needed before Not Applicable may be proposed; and
- whether the contract is generated and unreviewed or expert reviewed.

Exact duplicate statements keep separate stable IDs but point to one canonical
control ID. This is an alias, not permission to silently merge controls with
different scope or meaning.

Generated contracts are routing and review data only. They cannot produce a
Pass. Before a contract becomes executable, a control owner must split compound
promises, define measurable applicability and evidence, approve technology-
neutral acceptance criteria, and add positive, negative, Not Applicable,
unusual-layout, and adversarial fixtures.

Every `prc.run/v0.12` complete scan binds the SHA-256 of this whole contract
document and the individual contract SHA-256 on each of its 10,042 control
results. A missing, stale, reordered, or edited contract file stops the scan.
The ordinary local profile still executes only its explicitly registered narrow
assertions. Optional AI review receives the contract as a routing hint and must
keep repository, environment, and human evidence authority separate.

The generated
[`research/CONTROL_ACCEPTANCE_CRITERIA_REVIEW.md`](https://github.com/MarinJursic/production-readiness-checklist/blob/main/research/CONTROL_ACCEPTANCE_CRITERIA_REVIEW.md)
lists proposed acceptance criteria, today's real scanner behavior, future check
method, and wording warnings for every control. It is a review queue, not a
claim that 10,042 expert-approved acceptance tests already exist.

Regenerate and verify the file with:

```bash
python3 scripts/control_contracts.py generate
python3 scripts/control_contracts.py check
```
