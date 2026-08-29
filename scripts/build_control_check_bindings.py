#!/usr/bin/env python3
"""Build fail-closed checker bindings from the reviewed control classifications."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import sys
from pathlib import Path
from typing import Any

from jsonschema import Draft202012Validator


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_REGISTRY = ROOT / "catalog" / "control-id-registry.json"
DEFAULT_FINAL = ROOT / "research" / "control-classification" / "final"
DEFAULT_OUTPUT = ROOT / "catalog" / "control-check-bindings.json"
DEFAULT_SCHEMA = ROOT / "schemas" / "control-check-bindings.schema.json"

SCHEMA_VERSION = "prc.control-check-bindings/v0.1"
GENERATOR_ID = "prc.build-control-check-bindings@0.1"
FINAL_SCHEMA = "prc.control-classification-final/v0.1"
REGISTRY_SCHEMA = "prc.control-id-registry/v0.1"
CONTROL_ID = re.compile(r"^(?:PRC-[0-9]{2}-[0-9]{3}|USEQ-[A-F0-9]{8})$")
DIGEST = re.compile(r"^[0-9a-f]{64}$")

IMPLEMENTATIONS = {
    "inventory_fact": "prc.check.inventory-fact@0.1",
    "structured_document": "prc.check.structured-document@0.1",
    "package_metadata": "prc.check.package-metadata@0.1",
    "ci_policy": "prc.check.ci-policy@0.1",
    "container_iac": "prc.check.container-iac@0.1",
    "source_ast": "prc.check.source-ast@0.1",
    "artifact_integrity": "prc.check.artifact-integrity@0.1",
    "analysis_adapter": "prc.check.analysis-adapter@0.1",
    "execution_evidence": "prc.check.execution-evidence@0.1",
    "environment_evidence": "prc.check.environment-evidence@0.1",
    "structured_record": "prc.check.structured-record@0.1",
}

IMPLEMENTATION_AUTHORITIES = {
    "inventory_fact": ["repository"],
    "structured_document": ["repository"],
    "package_metadata": ["artifact", "external_registry", "repository"],
    "ci_policy": ["environment", "repository"],
    "container_iac": ["environment", "repository"],
    "source_ast": ["repository"],
    "artifact_integrity": ["artifact"],
    "analysis_adapter": ["executed"],
    "execution_evidence": ["executed"],
    "environment_evidence": ["environment", "external_registry"],
    "structured_record": ["structured_record"],
}

IMPLEMENTATION_CAPABILITIES = {
    "inventory_fact": "read_workspace",
    "structured_document": "read_workspace",
    "package_metadata": "read_workspace_or_external_readonly",
    "ci_policy": "read_workspace_or_environment_readonly",
    "container_iac": "read_workspace_or_environment_readonly",
    "source_ast": "read_workspace",
    "artifact_integrity": "read_artifact",
    "analysis_adapter": "executed_evidence",
    "execution_evidence": "executed_evidence",
    "environment_evidence": "external_readonly",
    "structured_record": "structured_record_evidence",
}

SUPPORTED_AUTHORITIES = {
    "repository",
    "artifact",
    "executed",
    "environment",
    "external_registry",
    "structured_record",
}

DETERMINISTIC_ROUTES = {
    "local_static",
    "artifact_verification",
    "bounded_execution",
    "external_readonly_query",
    "structured_record_validation",
    "deterministic_composite",
}

FINAL_ENVELOPE_FIELDS = {
    "schema_version", "packet_id", "methodology_sha256", "registry_version",
    "registry_sha256", "source_path", "control_count", "controls",
}
FINAL_ROW_FIELDS = {
    "control_id", "revision", "semantic_sha256", "classification", "route",
    "reason", "deterministic_clauses", "nondeterministic_reason_codes",
    "decision_basis", "skeptical_verdict", "counterexample_analysis",
}
CLAUSE_FIELDS = {"statement", "checker_family", "evidence_authority"}


class BindingError(ValueError):
    """The reviewed inputs cannot safely produce executable bindings."""


def canonical_bytes(value: Any) -> bytes:
    return json.dumps(value, ensure_ascii=False, sort_keys=True, separators=(",", ":")).encode("utf-8")


def digest_value(value: Any) -> str:
    return hashlib.sha256(canonical_bytes(value)).hexdigest()


def read_object(path: Path, label: str) -> tuple[dict[str, Any], bytes]:
    try:
        data = path.read_bytes()
        value = json.loads(data)
    except (OSError, json.JSONDecodeError) as error:
        raise BindingError(f"cannot read {label} {path}: {error}") from error
    if not isinstance(value, dict):
        raise BindingError(f"{label} {path} must be a JSON object")
    return value, data


def load_registry(path: Path) -> tuple[dict[str, dict[str, Any]], dict[str, Any], str]:
    document, data = read_object(path, "registry")
    if document.get("schema_version") != REGISTRY_SCHEMA:
        raise BindingError("unsupported control registry schema")
    entries = document.get("entries")
    if not isinstance(entries, list) or document.get("entry_count") != len(entries) or not entries:
        raise BindingError("registry entry count is invalid")
    indexed: dict[str, dict[str, Any]] = {}
    for entry in entries:
        if not isinstance(entry, dict):
            raise BindingError("registry entry must be an object")
        control_id = entry.get("id")
        if not isinstance(control_id, str) or not CONTROL_ID.fullmatch(control_id) or control_id in indexed:
            raise BindingError(f"invalid or duplicate registry control {control_id!r}")
        if not isinstance(entry.get("revision"), int) or entry["revision"] < 1:
            raise BindingError(f"{control_id}: invalid registry revision")
        if not isinstance(entry.get("semantic_sha256"), str) or not DIGEST.fullmatch(entry["semantic_sha256"]):
            raise BindingError(f"{control_id}: invalid registry semantic digest")
        indexed[control_id] = entry
    return indexed, document, hashlib.sha256(data).hexdigest()


def load_final_rows(final_dir: Path, registry: dict[str, dict[str, Any]], registry_version: str,
                    registry_sha256: str) -> tuple[list[dict[str, Any]], str, str]:
    paths = sorted(final_dir.glob("*.json"))
    if not paths:
        raise BindingError(f"no final classification packets in {final_dir}")
    corpus = hashlib.sha256()
    rows: list[dict[str, Any]] = []
    seen: set[str] = set()
    methodology_sha256 = ""
    for path in paths:
        document, data = read_object(path, "final classification")
        corpus.update(path.name.encode("utf-8"))
        corpus.update(b"\0")
        corpus.update(data)
        corpus.update(b"\0")
        if set(document) != FINAL_ENVELOPE_FIELDS or document.get("schema_version") != FINAL_SCHEMA:
            raise BindingError(f"{path}: final classification envelope is invalid")
        if document.get("registry_version") != registry_version or document.get("registry_sha256") != registry_sha256:
            raise BindingError(f"{path}: final classification is bound to a different registry")
        packet_methodology = document.get("methodology_sha256")
        if not isinstance(packet_methodology, str) or not DIGEST.fullmatch(packet_methodology):
            raise BindingError(f"{path}: invalid methodology digest")
        if methodology_sha256 and methodology_sha256 != packet_methodology:
            raise BindingError("final packets use different methodology digests")
        methodology_sha256 = packet_methodology
        packet_rows = document.get("controls")
        if not isinstance(packet_rows, list) or document.get("control_count") != len(packet_rows):
            raise BindingError(f"{path}: final control count is invalid")
        for row in packet_rows:
            if not isinstance(row, dict) or set(row) != FINAL_ROW_FIELDS:
                raise BindingError(f"{path}: final row fields are invalid")
            control_id = row.get("control_id")
            if control_id in seen:
                raise BindingError(f"duplicate final classification for {control_id}")
            source = registry.get(control_id)
            if source is None:
                raise BindingError(f"final classification contains unknown control {control_id}")
            if row.get("revision") != source["revision"] or row.get("semantic_sha256") != source["semantic_sha256"]:
                raise BindingError(f"{control_id}: final classification is stale")
            if row.get("classification") not in {"deterministic", "nondeterministic"}:
                raise BindingError(f"{control_id}: unsupported classification")
            clauses = row.get("deterministic_clauses")
            if not isinstance(clauses, list):
                raise BindingError(f"{control_id}: clauses must be a list")
            if row["classification"] == "deterministic":
                if (row.get("decision_basis") not in {"two_stage_confirmed", "strength_audit_confirmed"} or
                        row.get("skeptical_verdict") != "confirmed_deterministic" or
                        row.get("route") not in DETERMINISTIC_ROUTES or not clauses or
                        row.get("nondeterministic_reason_codes") != []):
                    raise BindingError(f"{control_id}: deterministic row lacks reviewed strength proof")
                for clause in clauses:
                    validate_clause(clause, control_id)
            # A primary nondeterministic `mixed` row may retain objectively
            # separable clauses for explanation. It still receives no runtime
            # binding because the complete control is nondeterministic.
            seen.add(control_id)
            rows.append(row)
    missing = sorted(set(registry) - seen)
    if missing:
        raise BindingError(f"final classification is missing {len(missing)} controls; first is {missing[0]}")
    return rows, methodology_sha256, corpus.hexdigest()


def validate_clause(clause: Any, control_id: str) -> None:
    if not isinstance(clause, dict) or set(clause) != CLAUSE_FIELDS:
        raise BindingError(f"{control_id}: invalid deterministic clause")
    if not isinstance(clause["statement"], str) or not clause["statement"].strip():
        raise BindingError(f"{control_id}: empty deterministic clause")
    if clause["checker_family"] not in IMPLEMENTATIONS:
        raise BindingError(f"{control_id}: unsupported checker family {clause['checker_family']!r}")
    if clause["evidence_authority"] not in SUPPORTED_AUTHORITIES:
        raise BindingError(f"{control_id}: unsupported evidence authority {clause['evidence_authority']!r}")


def implementation_registry() -> list[dict[str, Any]]:
    records = []
    for family, implementation_id in sorted(IMPLEMENTATIONS.items()):
        contract = {
            "checker_family": family,
            "implementation_id": implementation_id,
            "supported_evidence_authorities": IMPLEMENTATION_AUTHORITIES[family],
            "capability_class": IMPLEMENTATION_CAPABILITIES[family],
            "external_provider_registration": "required_when_authority_is_environment_or_external_registry",
        }
        records.append({
            **contract,
            "implementation_contract_sha256": digest_value(contract),
            "registration_state": "runtime_registration_required",
            "on_unregistered": "blocked",
            "external_provider_claimed": False,
        })
    return records


def result_contract() -> dict[str, Any]:
    return {
        "pass": {
            "result": "pass",
            "requires_all": [
                "applicability_resolved_applicable",
                "exact_implementation_registered",
                "provider_contract_registered",
                "required_evidence_authority_available",
                "complete_bounded_evidence",
                "full_clause_proven",
            ],
        },
        "fail": {
            "result": "fail",
            "requires_all": [
                "applicability_resolved_applicable",
                "exact_implementation_registered",
                "provider_contract_registered",
                "required_evidence_authority_available",
                "complete_bounded_evidence",
                "bounded_violation_observed",
            ],
        },
        "blocked": {
            "result": "blocked",
            "when_any": [
                "applicability_missing_or_ambiguous",
                "implementation_unregistered",
                "provider_contract_missing",
                "required_evidence_authority_unavailable",
                "unsupported_target",
                "complete_inventory_unavailable",
                "evidence_missing_stale_partial_or_conflicting",
                "rate_limited_or_ambiguous_external_result",
            ],
        },
    }


def binding_for(row: dict[str, Any], implementations: dict[str, dict[str, Any]]) -> dict[str, Any]:
    clauses = []
    for ordinal, source in enumerate(row["deterministic_clauses"], start=1):
        clause_digest = digest_value(source)
        clauses.append(
            {
                "clause_id": clause_digest,
                "ordinal": ordinal,
                "statement": source["statement"],
                "checker_family": source["checker_family"],
                "evidence_authority": source["evidence_authority"],
                "implementation_id": IMPLEMENTATIONS[source["checker_family"]],
                "implementation_contract_sha256": implementations[source["checker_family"]]["implementation_contract_sha256"],
                "implementation_registration": "required_before_execution",
                "provider_contract": "required_before_execution",
                "external_provider_claimed": False,
                "result_contract": result_contract(),
            }
        )
    return {
        "control_id": row["control_id"],
        "revision": row["revision"],
        "semantic_sha256": row["semantic_sha256"],
        "final_row_sha256": digest_value(row),
        "route": row["route"],
        "aggregation": "all_clauses_pass",
        "applicability_contract": {
            "applicable_result": "applicable",
            "not_applicable_requires": "authoritative_bounded_absence_proof",
            "unresolved_result": "blocked",
        },
        "clauses": clauses,
    }


def build(registry_path: Path, final_dir: Path) -> dict[str, Any]:
    registry, registry_document, registry_sha256 = load_registry(registry_path)
    rows, methodology_sha256, corpus_sha256 = load_final_rows(
        final_dir, registry, registry_document["registry_version"], registry_sha256
    )
    registry_records = implementation_registry()
    implementations = {item["checker_family"]: item for item in registry_records}
    bindings = [binding_for(row, implementations) for row in rows if row["classification"] == "deterministic"]
    bindings.sort(key=lambda item: item["control_id"])
    if len({item["control_id"] for item in bindings}) != len(bindings):
        raise BindingError("generated bindings contain duplicate controls")
    return {
        "schema_version": SCHEMA_VERSION,
        "generator_id": GENERATOR_ID,
        "registry_version": registry_document["registry_version"],
        "registry_sha256": registry_sha256,
        "methodology_sha256": methodology_sha256,
        "classification_corpus_sha256": corpus_sha256,
        "implementation_registry": registry_records,
        "binding_count": len(bindings),
        "bindings": bindings,
    }


def validate_artifact(document: dict[str, Any], schema_path: Path) -> None:
    schema, _ = read_object(schema_path, "binding schema")
    errors = sorted(Draft202012Validator(schema).iter_errors(document), key=lambda e: tuple(map(str, e.absolute_path)))
    if errors:
        first = errors[0]
        location = ".".join(map(str, first.absolute_path)) or "<root>"
        raise BindingError(f"binding schema error at {location}: {first.message}")
    ids = [binding["control_id"] for binding in document["bindings"]]
    if ids != sorted(ids) or len(ids) != len(set(ids)) or document["binding_count"] != len(ids):
        raise BindingError("bindings must be unique, ordered, and count-bound")
    implementations = {item["checker_family"]: item for item in document["implementation_registry"]}
    if set(implementations) != set(IMPLEMENTATIONS):
        raise BindingError("implementation registry does not cover exactly the reviewed checker families")
    for binding in document["bindings"]:
        if [clause["ordinal"] for clause in binding["clauses"]] != list(range(1, len(binding["clauses"]) + 1)):
            raise BindingError(f"{binding['control_id']}: clause ordinals are not contiguous")
        if len({clause["clause_id"] for clause in binding["clauses"]}) != len(binding["clauses"]):
            raise BindingError(f"{binding['control_id']}: duplicate clause IDs")
        for clause in binding["clauses"]:
            if implementations[clause["checker_family"]]["implementation_id"] != clause["implementation_id"]:
                raise BindingError(f"{binding['control_id']}: implementation mapping is inconsistent")
            if implementations[clause["checker_family"]]["implementation_contract_sha256"] != clause["implementation_contract_sha256"]:
                raise BindingError(f"{binding['control_id']}: implementation contract digest is inconsistent")
            if clause["evidence_authority"] not in implementations[clause["checker_family"]]["supported_evidence_authorities"]:
                raise BindingError(f"{binding['control_id']}: checker family cannot provide the reviewed evidence authority")
            if clause["external_provider_claimed"]:
                raise BindingError(f"{binding['control_id']}: declarative binding claims an external provider")


def serialize(document: dict[str, Any]) -> str:
    return json.dumps(document, indent=2, ensure_ascii=False) + "\n"


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("command", choices=("generate", "validate"))
    parser.add_argument("--registry", type=Path, default=DEFAULT_REGISTRY)
    parser.add_argument("--final-dir", type=Path, default=DEFAULT_FINAL)
    parser.add_argument("--output", type=Path, default=DEFAULT_OUTPUT)
    parser.add_argument("--schema", type=Path, default=DEFAULT_SCHEMA)
    parser.add_argument("--stdout", action="store_true")
    args = parser.parse_args()
    try:
        expected = build(args.registry, args.final_dir)
        validate_artifact(expected, args.schema)
        if args.command == "generate":
            text = serialize(expected)
            if args.stdout:
                sys.stdout.write(text)
            else:
                args.output.parent.mkdir(parents=True, exist_ok=True)
                args.output.write_text(text, encoding="utf-8")
                print(f"generated {expected['binding_count']} deterministic control bindings at {args.output}")
        else:
            actual, _ = read_object(args.output, "binding artifact")
            validate_artifact(actual, args.schema)
            if canonical_bytes(actual) != canonical_bytes(expected):
                raise BindingError("binding artifact is stale or differs from deterministic generation")
            print(f"validated {actual['binding_count']} deterministic control bindings")
        return 0
    except BindingError as error:
        print(f"control check binding build failed: {error}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
