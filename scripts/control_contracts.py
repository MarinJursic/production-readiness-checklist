#!/usr/bin/env python3
"""Generate the compact, machine-readable review contract for every control."""

from __future__ import annotations

import argparse
import hashlib
import importlib.util
import json
import sys
from collections import defaultdict
from pathlib import Path
from typing import Any

from jsonschema import Draft202012Validator


ROOT = Path(__file__).resolve().parents[1]
REGISTRY = ROOT / "catalog" / "control-id-registry.json"
OUTPUT = ROOT / "catalog" / "control-contracts.json"
SCHEMA = ROOT / "schemas" / "control-contracts.schema.json"


def load_acceptance_module() -> Any:
    path = ROOT / "scripts" / "generate_control_acceptance_review.py"
    spec = importlib.util.spec_from_file_location("prc_acceptance_review", path)
    if spec is None or spec.loader is None:
        raise ValueError("cannot load acceptance-review generator")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def canonical_json(value: Any) -> bytes:
    return json.dumps(value, ensure_ascii=False, separators=(",", ":"), sort_keys=True).encode("utf-8")


def normalized_statement(value: str) -> str:
    return " ".join(value.casefold().split()).rstrip(".")


def evaluation_class(flags: dict[str, Any]) -> str:
    if flags["human"] and (flags["runtime"] or flags["repository"]):
        return "mixed"
    if flags["human"]:
        return "human_external"
    if flags["runtime"] and flags["repository"]:
        return "mixed"
    if flags["runtime"]:
        return "environment"
    if flags["repository"]:
        return "repository"
    return "unclassified"


def automation_class(kind: str, flags: dict[str, Any]) -> str:
    if kind == "human_external":
        return "human_or_external_required"
    if kind == "environment":
        return "environment_evidence_required"
    if kind == "mixed":
        return "mixed_evidence_required"
    if kind == "repository" and not flags["compound"] and not flags["vague"]:
        return "deterministic_candidate"
    return "ai_advisory_candidate"


def evidence_authorities(kind: str) -> list[str]:
    values: dict[str, list[str]] = {
        "repository": ["repository", "artifact"],
        "environment": ["environment"],
        "human_external": ["human"],
        "mixed": ["repository", "artifact", "environment", "human"],
        "unclassified": ["declared", "human"],
    }
    return values[kind]


def build_document() -> dict[str, Any]:
    acceptance = load_acceptance_module()
    registry_bytes = REGISTRY.read_bytes()
    registry = json.loads(registry_bytes)
    entries = registry["entries"]
    contexts = acceptance.control_context(entries)
    grouped: dict[str, list[str]] = defaultdict(list)
    for entry in entries:
        grouped[normalized_statement(entry["statement"])].append(entry["id"])
    canonical = {control_id: min(ids) for ids in grouped.values() for control_id in ids}

    contracts = []
    for entry in entries:
        flags = acceptance.flags_for(entry, contexts.get(entry["id"], []))
        kind = evaluation_class(flags)
        contract: dict[str, Any] = {
            "control_id": entry["id"],
            "revision": entry["revision"],
            "contract_status": "generated_unreviewed" if entry["status"] == "active" else "retired",
            "canonical_control_id": canonical[entry["id"]],
            "evaluation_class": kind,
            "automation_class": automation_class(kind, flags),
            "applicability_class": "conditional" if flags["conditional"] else "scope_required",
            "atomicity": "compound_review_required" if flags["compound"] else "apparently_atomic",
            "complete_inventory_required": bool(flags["absolute"]),
            "negative_condition": bool(flags["negative"]),
            "project_thresholds_required": bool(flags["vague"]),
            "evidence_authorities": evidence_authorities(kind),
            "not_applicable_proof": "The trigger is affirmatively absent for the recorded scope and the reason is evidence-bound.",
            "reviewer_status": "generated_unreviewed",
        }
        contract["contract_sha256"] = hashlib.sha256(canonical_json(contract)).hexdigest()
        contracts.append(contract)

    return {
        "schema_version": "prc.control-contracts/v0.1",
        "registry_sha256": hashlib.sha256(registry_bytes).hexdigest(),
        "registry_version": registry["registry_version"],
        "contract_count": len(contracts),
        "contracts": contracts,
    }


def generated_text() -> str:
    return json.dumps(build_document(), ensure_ascii=False, separators=(",", ":")) + "\n"


def validate(document: dict[str, Any]) -> None:
    schema = json.loads(SCHEMA.read_text(encoding="utf-8"))
    errors = sorted(Draft202012Validator(schema).iter_errors(document), key=lambda item: list(item.path))
    if errors:
        raise ValueError("control contract schema failed: " + "; ".join(error.message for error in errors[:10]))
    ids = [item["control_id"] for item in document["contracts"]]
    if ids != sorted(ids) or len(ids) != len(set(ids)):
        raise ValueError("control contracts must be uniquely ordered by control ID")
    for item in document["contracts"]:
        digest = item["contract_sha256"]
        unsigned = {key: value for key, value in item.items() if key != "contract_sha256"}
        if hashlib.sha256(canonical_json(unsigned)).hexdigest() != digest:
            raise ValueError(f"contract digest mismatch for {item['control_id']}")


def generate() -> None:
    document = build_document()
    validate(document)
    OUTPUT.write_text(json.dumps(document, ensure_ascii=False, separators=(",", ":")) + "\n", encoding="utf-8")
    print(f"generated {OUTPUT.relative_to(ROOT)}")


def check() -> None:
    document = build_document()
    validate(document)
    expected = json.dumps(document, ensure_ascii=False, separators=(",", ":")) + "\n"
    if not OUTPUT.exists() or OUTPUT.read_text(encoding="utf-8") != expected:
        raise ValueError("catalog/control-contracts.json is stale; run python3 scripts/control_contracts.py generate")
    print(f"verified {OUTPUT.relative_to(ROOT)}")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("command", choices=("generate", "check"))
    args = parser.parse_args()
    try:
        generate() if args.command == "generate" else check()
    except (OSError, ValueError, json.JSONDecodeError) as error:
        print(f"control contract generation failed: {error}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
