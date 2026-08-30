# Signed authoritative evidence bundles

The normal `prc` scan is self-contained and needs no evidence bundle. This
advanced interface lets a separate trusted collector supply typed facts for one
or more reviewed deterministic programs without loading a plug-in or running
collector code inside the scanner.

It is a verification interface, not an evidence generator. `prc` does not hold
private keys, sign bundles, contact a registry, or decide that a collector is
trustworthy.

## Run a scan with one bundle

```bash
prc scan /path/to/project \
  --evidence-bundle evidence.json \
  --evidence-trust-store trust-store.json \
  --evidence-policy-signature policy-signature.json \
  --evidence-signature authority-evidence-signature.json
```

All four options are required together. The single-bundle form accepts at most one bundle. The
bundle may contain 1–765 entries from exactly one authority and may be no larger
than 32 MiB. It is strict JSON; duplicate keys, unknown fields, trailing data,
unsafe links, and files changed during reading are rejected. Signature identity
uses the scanner's canonical typed encoding, so harmless presentation whitespace
does not create a different signed subject and a saved run can reconstruct it.

The accepted authorities are `repository`, `artifact`, `executed`,
`environment`, `external_registry`, and `structured_record`.

## Run one scan with every authority

A production assessment normally needs more than one authority. Put one trust
store, up to six authority bundles, and their signatures in one private
directory. Then add a small `evidence-set.json` beside them:

```json
{
  "schema_version": "prc.authoritative-evidence-set/v0.1",
  "trust_store_file": "trust-store.json",
  "bundles": [
    {
      "authority": "artifact",
      "bundle_file": "artifact.json",
      "policy_signature_file": "artifact-policy.json",
      "evidence_signature_file": "artifact-evidence.json"
    },
    {
      "authority": "repository",
      "bundle_file": "repository.json",
      "policy_signature_file": "repository-policy.json",
      "evidence_signature_file": "repository-evidence.json"
    }
  ]
}
```

Entries must use unique authorities in alphabetical order. Every referenced
value is a sibling file name; paths, links, reused files, duplicate authorities,
and more than six bundles are rejected. Run the complete set with one option:

```bash
prc scan /path/to/project --evidence-set /private/evidence/evidence-set.json
```

`--evidence-set` cannot be combined with the four single-bundle options. The
scanner verifies the entire set before attaching any imported result. A bad
signature or malformed later bundle therefore stops the import instead of
leaving a partly trusted assessment. The set can carry evidence for all 765
exact clauses, but it does not create that evidence: the named authorities must
still observe and sign complete facts.

## Why two signatures are required

Two different Ed25519 keys sign two related, immutable subjects:

1. Before collection, a policy key with the `control-policy-bundle` scope signs
   the scanner-defined SHA-256 digest of the exact reviewed programs, collector
   IDs, scope, freshness, applicability contract, and typed parameter values.
   Observations are deliberately excluded because they do not exist yet.
2. After collection, a different key, limited to the bundle authority, signs
   the scanner-defined canonical SHA-256 of the completed typed bundle and
   attests the observed facts.
   For example, repository evidence requires the
   `control-evidence-repository` scope.

The trust store rejects duplicate public keys, so different key IDs cannot hide
the same key material. Policy must be signed no later than evidence collection;
the evidence signature must be issued no earlier than the observations it
attests. Both keys must be active and valid at scan time.

This separation follows the same basic safety idea as SLSA artifact
verification: trust policy is selected independently, the exact subject digest
is checked, and signer identity and expected inputs are verified rather than
accepted from the artifact itself. See the official
[SLSA artifact verification](https://slsa.dev/spec/v1.2/verifying-artifacts) and
[provenance](https://slsa.dev/spec/v1.2/provenance) specifications.

## What is bound

Before evaluating one entry, the scanner confirms all of the following:

- the policy digest matches the policy signature and the canonical bundle digest
  matches the evidence signature;
- the bundle's raw program-catalog digest matches the scanner's current
  `catalog/control-check-programs.json` bytes;
- the bundle inventory digest matches the exact project inventory created for
  this scan;
- the template ID, control revision, clause, implementation digest, authority,
  predicate, parameter names, and parameter types match the reviewed catalog;
- program and evidence subjects are bound to the same inventory;
- evidence is structurally complete enough for the pure evaluator; and
- the same template was not already evaluated by a built-in collector or an
  earlier bundle entry.

The scanner never accepts a provider verdict. Bundle entries contain a reviewed
program and normalized typed evidence. Only the scanner's closed predicate
evaluator can produce Pass, Fail, Blocked, or Not Applicable.

The retained `prc.run/v0.13` result records both signature-verification records,
the trust-store digest, policy digest, canonical bundle digest, catalog digest,
inventory digest, authority, entry count, every materialized program, exact
evidence, and resulting clause evaluation. A history reload reconstructs the
canonical bundle and replays every predicate before accepting the run. The HTML
report keeps this under
**Scan metadata → Signed authoritative evidence**.

That shape also follows the evidence-traceability model in NIST OSCAL assessment
results, where subjects, observations, evidence, findings, risks, and timestamps
remain linked instead of collapsing into one unsupported status. See the
official [OSCAL assessment-results model](https://pages.nist.gov/OSCAL/learn/concepts/layer/assessment/assessment-results/).

## Producer contract

A producer must generate documents conforming to:

- `schemas/authoritative-evidence-bundle.schema.json`;
- `schemas/authoritative-evidence-set.schema.json` when several authorities are
  supplied together;
- `schemas/trust-store.schema.json`; and
- two `schemas/signature.schema.json` envelopes.

The policy signature uses artifact kind `control-policy-bundle`, the bundle ID
as `artifact_id`, and `scanner/evidencebundle.PolicySHA256` as `sha256`. That
function canonicalizes only the policy projection and can be called before
observations are populated. The evidence signature uses the exact authority
scope listed above, the same bundle ID, and
`scanner/evidencebundle.BundleSHA256` as `sha256` after observations are
populated.

The signature payload is domain-separated by the scanner. Collector authors
should use the reviewed `scanner/trust.SigningPayload` contract rather than
inventing a different byte format. The scanner intentionally provides no
private-key command. Key creation, custody, approval, rotation, revocation, and
signing belong in an independently reviewed producer system.

## Safe failure behavior

The import stops before attachment on a missing file, symlink, oversized file,
unknown field, changed byte, wrong digest, wrong authority, untrusted or revoked
key, reused key, invalid time, stale catalog, stale inventory, altered predicate,
wrong parameter type, malformed evidence, or duplicate template. It does not
fall back to unsigned evidence and does not turn import failure into a guessed
Fail or Pass.

The ordinary local scan still works without this interface. An exact program
with no built-in collector and no accepted bundle stays Blocked.
