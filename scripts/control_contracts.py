#!/usr/bin/env python3
"""Build reviewed control contracts from the validated classification corpus."""

from __future__ import annotations

import argparse
import hashlib
import json
import sys
from collections import Counter
from pathlib import Path
from typing import Any

from jsonschema import Draft202012Validator


ROOT = Path(__file__).resolve().parents[1]
REGISTRY = ROOT / "catalog" / "control-id-registry.json"
FINAL_ROOT = ROOT / "research" / "control-classification" / "final"
SUMMARY = ROOT / "research" / "control-classification" / "summary.json"
BINDINGS = ROOT / "catalog" / "control-check-bindings.json"
BINDING_ARTIFACT_SCHEMA = ROOT / "schemas" / "control-check-bindings.schema.json"
OUTPUT = ROOT / "catalog" / "control-contracts.json"
SCHEMA = ROOT / "schemas" / "control-contracts.schema.json"

SCHEMA_VERSION = "prc.control-contracts/v0.2"
GENERATOR_ID = "prc.control-contracts@0.2"
FINAL_SCHEMA = "prc.control-classification-final/v0.1"
SUMMARY_SCHEMA = "prc.control-classification-summary/v0.1"
REGISTRY_SCHEMA = "prc.control-id-registry/v0.1"
BINDING_SCHEMA = "prc.control-check-bindings/v0.1"

DETERMINISTIC_ROUTES = {
    "local_static",
    "artifact_verification",
    "bounded_execution",
    "external_readonly_query",
    "structured_record_validation",
    "deterministic_composite",
}
NONDETERMINISTIC_ROUTES = {
    "contextual_judgment",
    "accountable_human_decision",
    "specialist_or_legal_judgment",
    "empirical_protocol_undefined",
    "contract_incomplete",
    "mixed",
    "unbounded_claim",
}
CHECKER_FAMILIES = {
    "inventory_fact",
    "structured_document",
    "package_metadata",
    "ci_policy",
    "container_iac",
    "source_ast",
    "artifact_integrity",
    "analysis_adapter",
    "execution_evidence",
    "environment_evidence",
    "structured_record",
}
EVIDENCE_AUTHORITIES = {
    "repository",
    "artifact",
    "executed",
    "environment",
    "external_registry",
    "structured_record",
}
NONDETERMINISTIC_REASONS = {
    "contextual_judgment",
    "human_accountability",
    "legal_or_specialist_authority",
    "undefined_protocol",
    "undefined_scope",
    "undefined_threshold",
    "mixed_with_nondeterministic_child",
    "unbounded_scope",
    "classification_not_approved",
}
FINAL_ENVELOPE_FIELDS = {
    "schema_version",
    "packet_id",
    "methodology_sha256",
    "registry_version",
    "registry_sha256",
    "source_path",
    "control_count",
    "controls",
}
FINAL_ROW_FIELDS = {
    "control_id",
    "revision",
    "semantic_sha256",
    "classification",
    "route",
    "reason",
    "deterministic_clauses",
    "nondeterministic_reason_codes",
    "decision_basis",
    "skeptical_verdict",
    "counterexample_analysis",
}
BINDING_ENVELOPE_FIELDS = {
    "schema_version",
    "generator_id",
    "registry_version",
    "registry_sha256",
    "methodology_sha256",
    "classification_corpus_sha256",
    "implementation_registry",
    "binding_count",
    "bindings",
}
BINDING_FIELDS = {
    "control_id",
    "revision",
    "semantic_sha256",
    "final_row_sha256",
    "route",
    "aggregation",
    "applicability_contract",
    "clauses",
}
CLAUSE_FIELDS = {"statement", "checker_family", "evidence_authority"}


class ContractError(ValueError):
    """Reviewed inputs cannot safely produce the contract catalog."""


def canonical_json(value: Any) -> bytes:
    return json.dumps(value, ensure_ascii=False, separators=(",", ":"), sort_keys=True).encode("utf-8")


def digest_value(value: Any) -> str:
    return hashlib.sha256(canonical_json(value)).hexdigest()


def read_object(path: Path, label: str) -> tuple[dict[str, Any], bytes]:
    try:
        data = path.read_bytes()
        value = json.loads(data)
    except (OSError, json.JSONDecodeError) as error:
        raise ContractError(f"cannot read {label} {path}: {error}") from error
    if not isinstance(value, dict):
        raise ContractError(f"{label} {path} must be a JSON object")
    return value, data


def load_registry(path: Path) -> tuple[dict[str, dict[str, Any]], dict[str, Any], str]:
    document, data = read_object(path, "control registry")
    entries = document.get("entries")
    if (
        document.get("schema_version") != REGISTRY_SCHEMA
        or not isinstance(entries, list)
        or document.get("entry_count") != len(entries)
        or not entries
    ):
        raise ContractError("control registry envelope or entry count is invalid")
    active: dict[str, dict[str, Any]] = {}
    for entry in entries:
        if not isinstance(entry, dict) or not isinstance(entry.get("id"), str):
            raise ContractError("control registry contains an invalid entry")
        control_id = entry["id"]
        if entry.get("status") == "active":
            if control_id in active:
                raise ContractError(f"duplicate active registry control {control_id}")
            if not isinstance(entry.get("revision"), int) or entry["revision"] < 1:
                raise ContractError(f"{control_id}: invalid registry revision")
            if not isinstance(entry.get("semantic_sha256"), str):
                raise ContractError(f"{control_id}: invalid registry semantic digest")
            active[control_id] = entry
    if not active:
        raise ContractError("control registry contains no active controls")
    return active, document, hashlib.sha256(data).hexdigest()


def _validate_clause(clause: Any, control_id: str) -> None:
    if not isinstance(clause, dict) or set(clause) != CLAUSE_FIELDS:
        raise ContractError(f"{control_id}: deterministic clause fields are invalid")
    if not isinstance(clause["statement"], str) or not clause["statement"].strip():
        raise ContractError(f"{control_id}: deterministic clause statement is empty")
    if clause["checker_family"] not in CHECKER_FAMILIES:
        raise ContractError(f"{control_id}: unsupported checker family")
    if clause["evidence_authority"] not in EVIDENCE_AUTHORITIES:
        raise ContractError(f"{control_id}: unsupported evidence authority")


def _validate_final_row(row: Any, source: dict[str, Any]) -> None:
    control_id = source["id"]
    if not isinstance(row, dict) or set(row) != FINAL_ROW_FIELDS:
        raise ContractError(f"{control_id}: final classification fields are invalid")
    if (
        row.get("control_id") != control_id
        or row.get("revision") != source.get("revision")
        or row.get("semantic_sha256") != source.get("semantic_sha256")
    ):
        raise ContractError(f"{control_id}: final classification is stale")
    if not isinstance(row.get("reason"), str) or len(row["reason"].strip()) < 20:
        raise ContractError(f"{control_id}: final classification reason is too short")
    clauses = row.get("deterministic_clauses")
    reason_codes = row.get("nondeterministic_reason_codes")
    if not isinstance(clauses, list) or not isinstance(reason_codes, list):
        raise ContractError(f"{control_id}: final clauses or reason codes are invalid")
    for clause in clauses:
        _validate_clause(clause, control_id)
    if row.get("classification") == "deterministic":
        if (
            row.get("route") not in DETERMINISTIC_ROUTES
            or row.get("decision_basis") != "strength_audit_confirmed"
            or row.get("skeptical_verdict") != "confirmed_deterministic"
            or not clauses
            or reason_codes
            or not isinstance(row.get("counterexample_analysis"), str)
            or len(row["counterexample_analysis"].strip()) < 20
        ):
            raise ContractError(f"{control_id}: deterministic classification lacks strength-audit proof")
    elif row.get("classification") == "nondeterministic":
        if (
            row.get("route") not in NONDETERMINISTIC_ROUTES
            or not reason_codes
            or any(code not in NONDETERMINISTIC_REASONS for code in reason_codes)
        ):
            raise ContractError(f"{control_id}: nondeterministic classification is incomplete")
        if row.get("decision_basis") == "primary_nondeterministic":
            if row.get("skeptical_verdict") is not None or row.get("counterexample_analysis") is not None:
                raise ContractError(f"{control_id}: primary nondeterministic row has skeptical fields")
            if row.get("route") != "mixed" and clauses:
                raise ContractError(f"{control_id}: nondeterministic row retains unsafe clauses")
        elif row.get("decision_basis") in {"skeptically_rejected", "strength_audit_reclassified"}:
            if (
                (
                    row.get("decision_basis") == "skeptically_rejected"
                    and row.get("skeptical_verdict") != "rejected_nondeterministic"
                )
                or (
                    row.get("decision_basis") == "strength_audit_reclassified"
                    and row.get("skeptical_verdict") != "confirmed_deterministic"
                )
                or clauses
                or not isinstance(row.get("counterexample_analysis"), str)
                or len(row["counterexample_analysis"].strip()) < 20
            ):
                raise ContractError(f"{control_id}: reviewed nondeterministic decision is incomplete")
        else:
            raise ContractError(f"{control_id}: invalid nondeterministic decision basis")
    else:
        raise ContractError(f"{control_id}: unsupported final classification")


def load_final_rows(
    final_root: Path,
    summary_path: Path,
    registry: dict[str, dict[str, Any]],
    registry_version: str,
    registry_sha256: str,
) -> tuple[dict[str, dict[str, Any]], str, str, str]:
    summary, summary_bytes = read_object(summary_path, "classification summary")
    packets = summary.get("packets")
    if (
        summary.get("schema_version") != SUMMARY_SCHEMA
        or summary.get("registry_version") != registry_version
        or summary.get("registry_sha256") != registry_sha256
        or not isinstance(packets, list)
        or not packets
    ):
        raise ContractError("classification summary is stale or invalid")
    expected_names = []
    packet_summaries: dict[str, dict[str, Any]] = {}
    for packet in packets:
        if not isinstance(packet, dict) or not isinstance(packet.get("packet_id"), str):
            raise ContractError("classification summary has an invalid packet")
        name = f"{packet['packet_id']}.json"
        expected_names.append(name)
        packet_summaries[name] = packet
    if len(expected_names) != len(set(expected_names)):
        raise ContractError("classification summary contains duplicate packets")
    actual_names = sorted(path.name for path in final_root.glob("*.json"))
    if sorted(expected_names) != actual_names:
        missing = sorted(set(expected_names) - set(actual_names))
        extra = sorted(set(actual_names) - set(expected_names))
        raise ContractError(f"final classification packets differ from summary; missing={missing}, extra={extra}")

    rows: dict[str, dict[str, Any]] = {}
    counts: Counter[str] = Counter()
    corpus = hashlib.sha256()
    methodology_sha256 = summary.get("methodology_sha256")
    for name in actual_names:
        path = final_root / name
        document, data = read_object(path, "final classification")
        corpus.update(name.encode("utf-8"))
        corpus.update(b"\0")
        corpus.update(data)
        corpus.update(b"\0")
        if set(document) != FINAL_ENVELOPE_FIELDS or document.get("schema_version") != FINAL_SCHEMA:
            raise ContractError(f"{path}: final classification envelope is invalid")
        if (
            document.get("registry_version") != registry_version
            or document.get("registry_sha256") != registry_sha256
            or document.get("methodology_sha256") != methodology_sha256
            or f"{document.get('packet_id')}.json" != name
            or document.get("source_path") != packet_summaries[name].get("source_path")
        ):
            raise ContractError(f"{path}: final classification envelope is stale")
        packet_rows = document.get("controls")
        if not isinstance(packet_rows, list) or document.get("control_count") != len(packet_rows):
            raise ContractError(f"{path}: final classification count is invalid")
        for row in packet_rows:
            if not isinstance(row, dict):
                raise ContractError(f"{path}: final classification row is invalid")
            control_id = row.get("control_id")
            if control_id in rows:
                raise ContractError(f"duplicate final classification for {control_id}")
            source = registry.get(control_id)
            if source is None:
                raise ContractError(f"unknown final classification control {control_id}")
            _validate_final_row(row, source)
            rows[control_id] = row
            counts[row["classification"]] += 1
        packet_counts = Counter(row["classification"] for row in packet_rows)
        packet_summary = packet_summaries[name]
        if (
            packet_summary.get("control_count") != len(packet_rows)
            or packet_summary.get("deterministic") != packet_counts["deterministic"]
            or packet_summary.get("nondeterministic") != packet_counts["nondeterministic"]
        ):
            raise ContractError(f"{name}: classification packet summary is stale")
    if set(rows) != set(registry):
        missing = sorted(set(registry) - set(rows))
        extra = sorted(set(rows) - set(registry))
        raise ContractError(f"final corpus does not match active registry; missing={missing[:5]}, extra={extra[:5]}")
    if (
        summary.get("control_count") != len(rows)
        or summary.get("deterministic") != counts["deterministic"]
        or summary.get("nondeterministic") != counts["nondeterministic"]
    ):
        raise ContractError("classification summary counts are stale")
    return rows, methodology_sha256, corpus.hexdigest(), hashlib.sha256(summary_bytes).hexdigest()


def load_bindings(
    path: Path,
    final_rows: dict[str, dict[str, Any]],
    registry_version: str,
    registry_sha256: str,
    methodology_sha256: str,
    classification_corpus_sha256: str,
) -> tuple[dict[str, tuple[dict[str, Any], str]], str, str]:
    document, data = read_object(path, "control check bindings")
    binding_schema, _ = read_object(BINDING_ARTIFACT_SCHEMA, "control check binding schema")
    binding_errors = sorted(
        Draft202012Validator(binding_schema).iter_errors(document),
        key=lambda item: tuple(map(str, item.absolute_path)),
    )
    if binding_errors:
        first = binding_errors[0]
        location = ".".join(map(str, first.absolute_path)) or "<root>"
        raise ContractError(f"control check binding schema failed at {location}: {first.message}")
    if set(document) != BINDING_ENVELOPE_FIELDS or document.get("schema_version") != BINDING_SCHEMA:
        raise ContractError("control check binding envelope is invalid")
    if (
        document.get("registry_version") != registry_version
        or document.get("registry_sha256") != registry_sha256
        or document.get("methodology_sha256") != methodology_sha256
        or document.get("classification_corpus_sha256") != classification_corpus_sha256
    ):
        raise ContractError("control check bindings are stale against reviewed inputs")
    source_bindings = document.get("bindings")
    if not isinstance(source_bindings, list) or document.get("binding_count") != len(source_bindings):
        raise ContractError("control check binding count is invalid")
    bindings: dict[str, tuple[dict[str, Any], str]] = {}
    ordered_ids: list[str] = []
    for binding in source_bindings:
        if not isinstance(binding, dict) or set(binding) != BINDING_FIELDS:
            raise ContractError("control check binding fields are invalid")
        control_id = binding.get("control_id")
        if not isinstance(control_id, str) or control_id in bindings:
            raise ContractError(f"duplicate or invalid deterministic binding {control_id}")
        row = final_rows.get(control_id)
        if row is None or row["classification"] != "deterministic":
            raise ContractError(f"{control_id}: binding does not map to a deterministic control")
        if (
            binding.get("revision") != row["revision"]
            or binding.get("semantic_sha256") != row["semantic_sha256"]
            or binding.get("final_row_sha256") != digest_value(row)
            or binding.get("route") != row["route"]
            or not isinstance(binding.get("clauses"), list)
            or not binding["clauses"]
            or len(binding["clauses"]) != len(row["deterministic_clauses"])
        ):
            raise ContractError(f"{control_id}: deterministic binding is stale")
        for ordinal, (bound_clause, reviewed_clause) in enumerate(
            zip(binding["clauses"], row["deterministic_clauses"]), start=1
        ):
            if (
                bound_clause.get("ordinal") != ordinal
                or bound_clause.get("clause_id") != digest_value(reviewed_clause)
                or bound_clause.get("statement") != reviewed_clause["statement"]
                or bound_clause.get("checker_family") != reviewed_clause["checker_family"]
                or bound_clause.get("evidence_authority") != reviewed_clause["evidence_authority"]
            ):
                raise ContractError(f"{control_id}: deterministic binding clause is stale")
        ordered_ids.append(control_id)
        bindings[control_id] = (binding, digest_value(binding))
    if ordered_ids != sorted(ordered_ids):
        raise ContractError("deterministic bindings are not ordered by control ID")
    deterministic_ids = {
        control_id for control_id, row in final_rows.items() if row["classification"] == "deterministic"
    }
    if set(bindings) != deterministic_ids:
        missing = sorted(deterministic_ids - set(bindings))
        extra = sorted(set(bindings) - deterministic_ids)
        raise ContractError(
            "deterministic controls and check bindings are not one-to-one; "
            f"missing={missing[:5]}, extra={extra[:5]}"
        )
    return bindings, hashlib.sha256(data).hexdigest(), document["schema_version"]


def _compatibility_fields(row: dict[str, Any]) -> dict[str, Any]:
    route = row["route"]
    reason_codes = set(row["nondeterministic_reason_codes"])
    if row["classification"] == "deterministic":
        authorities = sorted({clause["evidence_authority"] for clause in row["deterministic_clauses"]})
        repository = bool(set(authorities) & {"repository", "artifact"})
        external = bool(set(authorities) & {"executed", "environment", "external_registry", "structured_record"})
        if repository and external:
            evaluation = "mixed"
        elif repository:
            evaluation = "repository"
        else:
            evaluation = "environment"
        automation = {
            "repository": "deterministic_candidate",
            "environment": "environment_evidence_required",
            "mixed": "mixed_evidence_required",
        }[evaluation]
    else:
        route_authorities = {
            "contextual_judgment": ["human"],
            "accountable_human_decision": ["human"],
            "specialist_or_legal_judgment": ["human"],
            "empirical_protocol_undefined": ["environment", "human"],
            "contract_incomplete": ["declared", "human"],
            "mixed": ["repository", "environment", "human"],
            "unbounded_claim": ["human"],
        }
        authorities = route_authorities[route]
        evaluation = {
            "contextual_judgment": "human_external",
            "accountable_human_decision": "human_external",
            "specialist_or_legal_judgment": "human_external",
            "empirical_protocol_undefined": "mixed",
            "contract_incomplete": "unclassified",
            "mixed": "mixed",
            "unbounded_claim": "unclassified",
        }[route]
        automation = {
            "human_external": "human_or_external_required",
            "mixed": "mixed_evidence_required",
            "unclassified": "ai_advisory_candidate",
        }[evaluation]
    return {
        "evaluation_class": evaluation,
        "automation_class": automation,
        "applicability_class": "scope_required",
        "atomicity": (
            "compound_review_required"
            if route in {"mixed", "deterministic_composite"}
            else "apparently_atomic"
        ),
        "complete_inventory_required": bool(
            reason_codes & {"undefined_scope", "unbounded_scope"}
            or any(
                clause["checker_family"] == "inventory_fact"
                for clause in row["deterministic_clauses"]
            )
        ),
        "negative_condition": "unbounded_scope" in reason_codes,
        "project_thresholds_required": "undefined_threshold" in reason_codes,
        "evidence_authorities": authorities,
    }


def build_document(
    *,
    registry_path: Path = REGISTRY,
    final_root: Path = FINAL_ROOT,
    summary_path: Path = SUMMARY,
    bindings_path: Path = BINDINGS,
) -> dict[str, Any]:
    registry, registry_document, registry_sha256 = load_registry(registry_path)
    final_rows, methodology_sha256, corpus_sha256, summary_sha256 = load_final_rows(
        final_root,
        summary_path,
        registry,
        registry_document["registry_version"],
        registry_sha256,
    )
    bindings, bindings_sha256, bindings_schema = load_bindings(
        bindings_path,
        final_rows,
        registry_document["registry_version"],
        registry_sha256,
        methodology_sha256,
        corpus_sha256,
    )
    contracts = []
    for control_id in sorted(registry):
        entry = registry[control_id]
        row = final_rows[control_id]
        binding = bindings.get(control_id)
        contract: dict[str, Any] = {
            "control_id": control_id,
            "revision": entry["revision"],
            "contract_status": "reviewed",
            "reviewer_status": "agent_reviewed",
            "classification": row["classification"],
            "classification_route": row["route"],
            "classification_decision_basis": row["decision_basis"],
            "classification_row_sha256": digest_value(row),
            "deterministic_binding_id": (
                f"{control_id}@{entry['revision']}" if binding is not None else None
            ),
            "deterministic_binding_sha256": binding[1] if binding is not None else None,
            "canonical_control_id": control_id,
            **_compatibility_fields(row),
            "not_applicable_proof": (
                "The trigger is affirmatively absent for the recorded scope and the reason is evidence-bound."
            ),
        }
        contract["contract_sha256"] = digest_value(contract)
        contracts.append(contract)
    return {
        "schema_version": SCHEMA_VERSION,
        "generator_id": GENERATOR_ID,
        "registry_sha256": registry_sha256,
        "registry_version": registry_document["registry_version"],
        "classification_methodology_sha256": methodology_sha256,
        "classification_summary_sha256": summary_sha256,
        "classification_corpus_sha256": corpus_sha256,
        "control_check_bindings_schema_version": bindings_schema,
        "control_check_bindings_sha256": bindings_sha256,
        "binding_count": len(bindings),
        "contract_count": len(contracts),
        "contracts": contracts,
    }


def generated_text(document: dict[str, Any] | None = None) -> str:
    return json.dumps(document or build_document(), ensure_ascii=False, separators=(",", ":")) + "\n"


def validate(document: dict[str, Any], schema_path: Path = SCHEMA) -> None:
    schema, _ = read_object(schema_path, "control contract schema")
    errors = sorted(
        Draft202012Validator(schema).iter_errors(document),
        key=lambda item: tuple(map(str, item.absolute_path)),
    )
    if errors:
        first = errors[0]
        location = ".".join(map(str, first.absolute_path)) or "<root>"
        raise ContractError(f"control contract schema failed at {location}: {first.message}")
    contracts = document["contracts"]
    ids = [item["control_id"] for item in contracts]
    if ids != sorted(ids) or len(ids) != len(set(ids)) or document["contract_count"] != len(ids):
        raise ContractError("control contracts must be uniquely ordered and count-bound")
    deterministic = 0
    binding_ids: set[str] = set()
    for item in contracts:
        digest = item["contract_sha256"]
        unsigned = {key: value for key, value in item.items() if key != "contract_sha256"}
        if digest_value(unsigned) != digest:
            raise ContractError(f"contract digest mismatch for {item['control_id']}")
        if item["classification"] == "deterministic":
            deterministic += 1
            binding_id = item["deterministic_binding_id"]
            expected_binding_id = f"{item['control_id']}@{item['revision']}"
            if binding_id != expected_binding_id:
                raise ContractError(
                    f"{item['control_id']}: deterministic binding identity does not match the contract"
                )
            if binding_id in binding_ids:
                raise ContractError(f"duplicate deterministic binding identity {binding_id}")
            binding_ids.add(binding_id)
        elif item["deterministic_binding_id"] is not None or item["deterministic_binding_sha256"] is not None:
            raise ContractError(f"{item['control_id']}: nondeterministic contract has a binding")
    if deterministic != document["binding_count"] or len(binding_ids) != document["binding_count"]:
        raise ContractError("deterministic contracts and binding count are not one-to-one")


def generate() -> None:
    document = build_document()
    validate(document)
    OUTPUT.write_text(generated_text(document), encoding="utf-8")
    print(f"generated {OUTPUT.relative_to(ROOT)} with {document['binding_count']} deterministic bindings")


def check() -> None:
    document = build_document()
    validate(document)
    expected = generated_text(document)
    if not OUTPUT.exists() or OUTPUT.read_text(encoding="utf-8") != expected:
        raise ContractError(
            "catalog/control-contracts.json is stale; run python3 scripts/control_contracts.py generate"
        )
    print(f"verified {OUTPUT.relative_to(ROOT)}")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("command", choices=("generate", "check"))
    args = parser.parse_args()
    try:
        generate() if args.command == "generate" else check()
    except (OSError, ContractError, json.JSONDecodeError) as error:
        print(f"control contract generation failed: {error}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
