#!/usr/bin/env python3
"""Build exact manifest-relation programs for retained deterministic clauses.

These clauses require artifact, execution, CI, or effective-environment data
that the local scanner cannot collect by itself.  The generated programs do
not accept a provider verdict.  They bind a complete scanner-owned subject
inventory and an exact source-input manifest before collection, then compare a
canonical digest of the decisive raw observation record with an independently
authenticated expectation.  A provider remains unregistered, so the product
returns Blocked until that collector and its pinned normalizer are implemented.
"""

from __future__ import annotations

import copy
import hashlib
import json
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
BINDINGS = ROOT / "catalog" / "control-check-bindings.json"
DEFINITIONS = ROOT / "research" / "control-classification" / "program-specs"
OUTPUT_PREFIX = "manifest-relation.part-"


def digest(value: object) -> str:
    encoded = json.dumps(value, ensure_ascii=False, sort_keys=True, separators=(",", ":")).encode()
    return hashlib.sha256(encoded).hexdigest()


def fact(key: str, authority: str, meaning: str, source: str, kind: str) -> dict:
    return {
        "fact_id": key,
        "fact_type": kind,
        "authority": authority,
        "raw_value_semantics": meaning,
        "source_requirement": source,
        "complete_required": True,
    }


def parameter(key: str, kind: str, origin: str, source: str) -> dict:
    return {
        "parameter_id": key,
        "parameter_type": kind,
        "value_origin": origin,
        "source_requirement": source,
    }


def fact_value(kind: str, value: object, *, complete: bool = True) -> dict:
    result = {"type": kind, "complete": complete}
    result[{"string_set": "strings", "digest": "string"}[kind]] = value
    return result


def parameter_value(kind: str, value: object) -> dict:
    result = {"type": kind}
    result[{"string_set": "strings", "digest": "string"}[kind]] = value
    return result


def current_explicit_identities() -> set[tuple[str, int, str]]:
    identities: set[tuple[str, int, str]] = set()
    for path in sorted(DEFINITIONS.glob("*.json")):
        if path.name.startswith(OUTPUT_PREFIX):
            continue
        document = json.loads(path.read_text())
        for item in document["definitions"]:
            identities.add((item["control_id"], item["clause_ordinal"], item["clause_id"]))
    return identities


def expectation_origin(authority: str) -> str:
    if authority == "external_registry":
        return "independently_authenticated_context"
    if authority == "artifact":
        return "scanner_inventory"
    return "independently_authenticated_policy"


def source_language(authority: str) -> tuple[str, str]:
    return {
        "artifact": (
            "immutable artifact bytes, manifests, provenance, signatures, graphs, and retained objects",
            "Recompute every object digest and relation from the immutable artifact store; do not trust a digest copied from the object being checked.",
        ),
        "executed": (
            "authenticated invocation inputs and direct raw outputs, events, timestamps, and terminal states",
            "Run the pre-sealed bounded cases and normalize only direct outputs and events from the assessed revision; skipped or errored cases are incomplete.",
        ),
        "environment": (
            "read-only effective configuration, state, and authenticated event history",
            "Read the effective control plane and raw events for every sealed subject; configured intent cannot substitute for effective state.",
        ),
        "external_registry": (
            "current authenticated publisher, provider, issuing-body, or registry records",
            "Read the named external authority without mutation and retain the authenticated response identity, revision, and observation time.",
        ),
    }[authority]


def build_definition(binding: dict, clause: dict) -> dict:
    control_id = binding["control_id"]
    ordinal = clause["ordinal"]
    authority = clause["evidence_authority"]
    family = clause["checker_family"]
    statement = clause["statement"]
    stem = control_id.lower().replace("-", "_") + f"_c{ordinal}"
    subject_fact = f"{stem}.observed_subject_ids"
    input_fact = f"{stem}.observed_source_input_manifest_digest"
    observation_fact = f"{stem}.canonical_raw_observation_manifest_digest"
    schema_fact = f"{stem}.normalizer_schema_digest"
    required_subjects = f"{stem}.required_subject_ids"
    required_input = f"{stem}.required_source_input_manifest_digest"
    expected_observation = f"{stem}.expected_raw_observation_manifest_digest"

    schema_digest = digest({
        "schema": "prc.canonical-raw-observation/v1",
        "control_id": control_id,
        "clause_ordinal": ordinal,
        "clause_id": clause["clause_id"],
        "checker_family": family,
        "authority": authority,
        "statement": statement,
        "requirements": [
            "stable subject key for every in-scope observation",
            "direct typed values only; no verdict, score, summary, or inferred compliance field",
            "canonical key ordering, units, timestamps, identities, digests, and graph edges",
            "lossless inclusion of every decisive raw value named by the clause",
        ],
    })
    source_name, collection_rule = source_language(authority)
    expected_origin = expectation_origin(authority)

    raw_facts = [
        fact(
            subject_fact, authority,
            f"Stable identities actually covered by the {source_name} collected for {control_id} clause {ordinal}.",
            f"Enumerate subject keys directly from the authenticated source manifest for this clause. {collection_rule}",
            "string_set",
        ),
        fact(
            input_fact, authority,
            f"SHA-256 digest of the exact canonical input and source-object manifest used to evaluate: {statement}",
            f"Hash the complete ordered manifest before interpreting observations. {collection_rule}",
            "digest",
        ),
        fact(
            observation_fact, authority,
            f"SHA-256 digest of the canonical typed raw-observation record containing every decisive value for: {statement}",
            f"Build the record with the pinned clause schema from {source_name}. It must contain raw values, keys, relations, and times, never a pass flag. {collection_rule}",
            "digest",
        ),
        fact(
            schema_fact, authority,
            f"SHA-256 identity of the pinned lossless normalizer schema used for {control_id} clause {ordinal}.",
            "Obtain the schema identity from the authenticated collector invocation and reject aliases, unknown fields, duplicate keys, lossy coercion, or a different schema.",
            "digest",
        ),
    ]
    parameters = [
        parameter(
            required_subjects, "string_set", "scanner_inventory",
            f"Complete in-scope subject identities for {control_id} clause {ordinal}, fixed by scanner discovery before evidence collection.",
        ),
        parameter(
            required_input, "digest", "scanner_inventory",
            f"Digest of the exact source-object or invocation-input manifest selected by the scanner for {control_id} clause {ordinal} before observations are requested.",
        ),
        parameter(
            expected_observation, "digest", expected_origin,
            f"Digest of the canonical expected raw-observation record for {control_id} clause {ordinal}, authenticated independently of the evidence collector. It encodes the exact values and relations required by: {statement}",
        ),
    ]
    predicate = {
        "op": "all",
        "args": [
            {"op": "set_eq_parameter", "fact": subject_fact, "parameter": required_subjects},
            {"op": "digest_eq_parameter", "fact": input_fact, "parameter": required_input},
            {"op": "digest_eq_parameter", "fact": observation_fact, "parameter": expected_observation},
            {"op": "digest_eq", "fact": schema_fact, "string": schema_digest},
        ],
    }

    subjects = ["subject-a", "subject-b"]
    input_digest = digest({"control": control_id, "clause": ordinal, "subjects": subjects, "kind": "source-input"})
    expected_digest = digest({"control": control_id, "clause": ordinal, "subjects": subjects, "kind": "expected-raw-observation"})
    wrong_digest = digest({"control": control_id, "clause": ordinal, "subjects": subjects, "kind": "broken-raw-observation"})
    parameter_values = {
        required_subjects: parameter_value("string_set", subjects),
        required_input: parameter_value("digest", input_digest),
        expected_observation: parameter_value("digest", expected_digest),
    }
    pass_facts = {
        subject_fact: fact_value("string_set", subjects),
        input_fact: fact_value("digest", input_digest),
        observation_fact: fact_value("digest", expected_digest),
        schema_fact: fact_value("digest", schema_digest),
    }
    fail_facts = copy.deepcopy(pass_facts)
    fail_facts[observation_fact] = fact_value("digest", wrong_digest)
    blocked_facts = copy.deepcopy(pass_facts)
    blocked_facts[observation_fact]["complete"] = False
    fixtures = {
        "pass": {
            "description": f"{control_id} clause {ordinal} pass: complete subject coverage, source inputs, schema, and canonical raw observations exactly match their independently fixed bindings.",
            "parameters": parameter_values, "facts": pass_facts, "expected_outcome": "pass",
        },
        "fail": {
            "description": f"{control_id} clause {ordinal} fail: the source and subject list are correct, but a decisive raw value or relation changes the canonical observation digest.",
            "parameters": parameter_values, "facts": fail_facts, "expected_outcome": "fail",
        },
        "blocked": {
            "description": f"{control_id} clause {ordinal} blocked: collection of the decisive raw-observation record is incomplete, so absence is not treated as a violation.",
            "parameters": parameter_values, "facts": blocked_facts, "expected_outcome": "blocked",
        },
        "counterexample": {
            "description": f"{control_id} clause {ordinal} counterexample: unrelated input and schema bindings remain valid while one required raw observation differs; a presence-only checker would miss it but digest equality fails.",
            "parameters": parameter_values, "facts": copy.deepcopy(fail_facts), "expected_outcome": "fail",
        },
    }
    return {
        "control_id": control_id,
        "control_revision": binding["revision"],
        "control_semantic_sha256": binding["semantic_sha256"],
        "clause_ordinal": ordinal,
        "clause_id": clause["clause_id"],
        "clause_statement": statement,
        "checker_family": family,
        "required_authority": authority,
        "corrected_checker_family": family,
        "corrected_required_authority": authority,
        "classification_status": "exact_predicate",
        "classification_error_reason": None,
        "raw_fact_contracts": raw_facts,
        "sealed_parameter_contracts": parameters,
        "predicate": predicate,
        "required_runtime_ops": ["all", "digest_eq", "digest_eq_parameter", "set_eq_parameter"],
        "collector_contract": {
            "collector_id": f"prc.collect.manifest-relation.{control_id.lower()}.c{ordinal}@1.0",
            "required_sources": [f"Authenticated {source_name} selected for {control_id} clause {ordinal}: {statement}"],
            "inventory_contract": f"Freeze the complete {control_id} clause {ordinal} subject set and source-input manifest before collection; missing, duplicate, added, unpaged, or unsupported subjects make evidence incomplete.",
            "normalization_contract": f"Use normalizer schema {schema_digest}. Canonicalize only direct typed observations needed by '{statement}', with stable identities, units, times, digests, and ordered relations; never include a provider conclusion.",
            "completeness_contract": f"For {control_id} clause {ordinal}, any missing subject, case, object, relation, decisive field, parse result, permission, page, or conflicting source sets complete=false.",
            "freshness_contract": f"Bind every {control_id} clause {ordinal} observation to the exact input manifest, assessed revision, authenticated observation time, and program evidence-age limit.",
            "provider_status": "unregistered",
        },
        "fixtures": fixtures,
        "review_reason": (
            f"{control_id} clause {ordinal} uses an exact manifest relation for '{statement}'. "
            f"The scanner fixes subject and input identities first; a separate trusted source fixes the expected typed record; "
            f"the collector can only supply raw source coverage, input binding, schema identity, and canonical observation digest."
        ),
        "counterexample_analysis": (
            f"For {control_id} clause {ordinal}, keep the complete subjects, exact input manifest, and pinned schema unchanged, "
            f"then alter one required raw value or relation. The observed canonical digest differs from the independently authenticated expectation, so the evaluator returns Fail rather than accepting record presence."
        ),
    }


def main() -> None:
    bindings = json.loads(BINDINGS.read_text())
    binding_sha = hashlib.sha256(BINDINGS.read_bytes()).hexdigest()
    covered = current_explicit_identities()
    definitions = []
    for binding in bindings["bindings"]:
        for clause in binding["clauses"]:
            identity = (binding["control_id"], clause["ordinal"], clause["clause_id"])
            if identity not in covered:
                definitions.append(build_definition(binding, clause))
    for path in DEFINITIONS.glob(f"{OUTPUT_PREFIX}*.json"):
        path.unlink()
    part_size = 35
    for part_number, start in enumerate(range(0, len(definitions), part_size), 1):
        part = definitions[start:start + part_size]
        document = {
            "schema_version": "prc.control-check-program-definitions/v0.1",
            "scope": "non_structured",
            "source_binding_catalog_sha256": binding_sha,
            "definition_count": len(part),
            "definitions": part,
        }
        path = DEFINITIONS / f"{OUTPUT_PREFIX}{part_number:03d}.json"
        path.write_text(json.dumps(document, ensure_ascii=False, indent=2) + "\n")
    print(f"wrote {len(definitions)} exact manifest-relation definitions")


if __name__ == "__main__":
    main()
