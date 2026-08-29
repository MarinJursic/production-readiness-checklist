#!/usr/bin/env python3
"""Create and validate disjoint rule-by-rule control classification work packets."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import sys
from collections import defaultdict
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
REGISTRY = ROOT / "catalog" / "control-id-registry.json"
SOURCE_ROOT = ROOT
PACKET_ROOT = ROOT / "research" / "control-classification" / "packets"
PRIMARY_ROOT = ROOT / "research" / "control-classification" / "primary"
SKEPTIC_PACKET_ROOT = ROOT / "research" / "control-classification" / "skeptic-packets"
SKEPTIC_ROOT = ROOT / "research" / "control-classification" / "skeptic"
FINAL_ROOT = ROOT / "research" / "control-classification" / "final"
STRENGTH_ROOT = ROOT / "research" / "control-classification" / "strength"
SUMMARY = ROOT / "research" / "control-classification" / "summary.json"
METHODOLOGY = ROOT / "docs" / "architecture" / "control-classification.md"
MAX_PACKET_CONTROLS = 250
MAX_SKEPTIC_CONTROLS = 100
HEADING = re.compile(r"^(#{1,6})\s+(.+?)\s*$")

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

GENERIC_WEAKENING_PREFIXES = (
    "For the sealed assessed scope, verify the versioned acceptance contract for:",
    "For the sealed assessed scope, verify every objective predicate explicitly required by:",
    "For the sealed assessed scope, verify the complete versioned acceptance contract for:",
    "For the sealed assessed scope, the configured checker verifies the versioned acceptance contract for:",
    "For every subject in the complete assessed inventory, an authenticated versioned record directly contains every identity, category, field, relation, threshold, date, scope item, and status explicitly required by the original control:",
    "Read-only effective configuration, state, and event evidence for every subject in the complete assessed inventory directly demonstrates the operating behavior required by the original control;",
    "Authenticated raw execution evidence covers every case in the pre-sealed scenario and subject inventory, is bound to the exact tested revisions and inputs, and records the objective observations required by the original control:",
    "Every item required by the original control is present in the complete retention inventory, remains retrievable by immutable identity, and preserves all named bindings and fields:",
    "Exact bytes, locally recomputed digests, immutable manifests, and authenticated provenance or signature records for the complete artifact inventory directly satisfy every artifact relation required by the original control:",
    "Read-only effective gate configuration implements the exact blocking, ordering, retry, failure, approval, or revocation behavior required by the original control for every in-scope gate:",
)
GENERIC_WEAKENING_STATEMENTS = {
    "Run the registered bounded protocol and verify authenticated evidence records the exact inputs, versions, identities, expected transitions, negative cases, outputs, and successful result.",
    "The authoritative record contains every explicit step type, scenario, method, role, metric, deviation, risk link, owner, expiry, or assurance binding named by this control.",
    "Authenticated bounded execution covers every scenario in the approved test matrix for the exact revision and records the required expected outcome for each case.",
    "The authoritative record contains complete in-scope entries with every explicit standard, matrix dimension, owner, state, field, budget, version, or lifecycle attribute required by this control.",
    "Bound authenticated evidence satisfies every exact isolation, timing, revision, participant, reassessment, or independent-review condition.",
    "Complete authoritative inventories enumerate every in-scope role, change, training requirement, practice environment, interface, or high-impact evaluation subject.",
    "The generated or final artifact digest, source digest, generator identity, schema, template, toolchain, provenance, and reviewed artifact binding satisfy the approved immutable manifest.",
    "Bound source, artifact, browser, or operational evidence satisfies every exact state, value, timing, isolation, or threshold condition in the approved contract.",
    "Complete authoritative inventories enumerate every in-scope page, journey, component, locale, asset, route, or runtime subject required by this control.",
    "The immutable artifact, dependency manifest, lock constraints, and provenance evidence match the approved identities, versions, digests, and intended-content set exactly.",
    "A bounded negative fixture exercises each configured trigger and produces the required non-passing, paused, rejected, revoked, or otherwise safe outcome; an unavailable, ambiguous, skipped, or errored fixture cannot pass.",
}


def load_registry() -> dict[str, Any]:
    return json.loads(REGISTRY.read_text(encoding="utf-8"))


def canonical_bytes(value: Any) -> bytes:
    return json.dumps(value, ensure_ascii=False, sort_keys=True, separators=(",", ":")).encode("utf-8")


def digest_value(value: Any) -> str:
    return hashlib.sha256(canonical_bytes(value)).hexdigest()


def slug(path: str) -> str:
    return re.sub(r"[^a-z0-9]+", "-", path.casefold()).strip("-")


def heading_context(path: str, requested_lines: set[int]) -> dict[int, list[str]]:
    headings: dict[int, str] = {}
    result: dict[int, list[str]] = {}
    lines = (SOURCE_ROOT / path).read_text(encoding="utf-8").splitlines()
    for number, line in enumerate(lines, start=1):
        match = HEADING.match(line)
        if match:
            level = len(match.group(1))
            headings[level] = match.group(2)
            headings = {key: value for key, value in headings.items() if key <= level}
        if number in requested_lines:
            result[number] = [headings[key] for key in sorted(headings)]
    return result


def generate_packets() -> None:
    registry = load_registry()
    grouped: dict[str, list[dict[str, Any]]] = defaultdict(list)
    for entry in registry["entries"]:
        grouped[entry["source"]["path"]].append(entry)
    PACKET_ROOT.mkdir(parents=True, exist_ok=True)
    expected = set()
    for source_path, entries in sorted(grouped.items()):
        contexts = heading_context(source_path, {entry["source"]["line"] for entry in entries})
        controls = []
        for entry in sorted(entries, key=lambda item: item["source"]["line"]):
            controls.append(
                {
                    "control_id": entry["id"],
                    "revision": entry["revision"],
                    "semantic_sha256": entry["semantic_sha256"],
                    "statement": entry["statement"],
                    "source": entry["source"],
                    "heading_trail": contexts.get(entry["source"]["line"], []),
                }
            )
        for chunk_index in range(0, len(controls), MAX_PACKET_CONTROLS):
            chunk = controls[chunk_index : chunk_index + MAX_PACKET_CONTROLS]
            part_number = chunk_index // MAX_PACKET_CONTROLS + 1
            packet_id = f"{slug(source_path)}-part-{part_number:03d}"
            name = packet_id + ".json"
            expected.add(name)
            document = {
                "schema_version": "prc.control-classification-packet/v0.1",
                "packet_id": packet_id,
                "methodology_sha256": hashlib.sha256(METHODOLOGY.read_bytes()).hexdigest(),
                "registry_version": registry["registry_version"],
                "registry_sha256": hashlib.sha256(REGISTRY.read_bytes()).hexdigest(),
                "source_path": source_path,
                "control_count": len(chunk),
                "controls": chunk,
            }
            (PACKET_ROOT / name).write_text(json.dumps(document, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")
    for path in PACKET_ROOT.glob("*.json"):
        if path.name not in expected:
            path.unlink()
    print(f"generated {len(expected)} review packets containing {registry['entry_count']} controls")


def validate_clause(clause: Any, control_id: str) -> None:
    if not isinstance(clause, dict) or set(clause) != {"statement", "checker_family", "evidence_authority"}:
        raise ValueError(f"{control_id}: invalid deterministic clause fields")
    if not isinstance(clause["statement"], str) or not clause["statement"].strip():
        raise ValueError(f"{control_id}: deterministic clause needs a statement")
    if clause["checker_family"] not in CHECKER_FAMILIES:
        raise ValueError(f"{control_id}: unsupported checker family {clause['checker_family']!r}")
    if clause["evidence_authority"] not in EVIDENCE_AUTHORITIES:
        raise ValueError(f"{control_id}: unsupported evidence authority {clause['evidence_authority']!r}")


def validate_strengthened_clause(clause: Any, control_id: str) -> None:
    validate_clause(clause, control_id)
    statement = clause["statement"].strip()
    if statement in GENERIC_WEAKENING_STATEMENTS or statement.startswith(GENERIC_WEAKENING_PREFIXES):
        raise ValueError(f"{control_id}: strengthened clause uses a forbidden generic weakening template")


def validate_review(review: Any, expected: dict[str, Any]) -> None:
    fields = {
        "control_id",
        "revision",
        "semantic_sha256",
        "classification",
        "route",
        "reason",
        "deterministic_clauses",
        "nondeterministic_reason_codes",
        "review_status",
    }
    if not isinstance(review, dict) or set(review) != fields:
        raise ValueError(f"{expected['control_id']}: review fields do not match the contract")
    control_id = expected["control_id"]
    if review["control_id"] != control_id:
        raise ValueError(f"expected {control_id}, received {review['control_id']!r}")
    if review["revision"] != expected["revision"] or review["semantic_sha256"] != expected["semantic_sha256"]:
        raise ValueError(f"{control_id}: review is not bound to the current rule revision")
    if not isinstance(review["reason"], str) or len(review["reason"].strip()) < 20:
        raise ValueError(f"{control_id}: classification reason is too short")
    if review["review_status"] != "primary_agent_reviewed":
        raise ValueError(f"{control_id}: primary review has an invalid status")
    clauses = review["deterministic_clauses"]
    reasons = review["nondeterministic_reason_codes"]
    if not isinstance(clauses, list) or not isinstance(reasons, list) or len(reasons) != len(set(reasons)):
        raise ValueError(f"{control_id}: clauses or reason codes are invalid")
    if review["classification"] == "deterministic":
        if review["route"] not in DETERMINISTIC_ROUTES or not clauses or reasons:
            raise ValueError(f"{control_id}: deterministic review is incomplete")
        for clause in clauses:
            validate_clause(clause, control_id)
    elif review["classification"] == "nondeterministic":
        if review["route"] not in NONDETERMINISTIC_ROUTES or not reasons:
            raise ValueError(f"{control_id}: nondeterministic review is incomplete")
        if any(reason not in NONDETERMINISTIC_REASONS for reason in reasons):
            raise ValueError(f"{control_id}: unsupported nondeterministic reason code")
        if review["route"] != "mixed" and clauses:
            raise ValueError(f"{control_id}: only mixed nondeterministic controls may retain deterministic partial clauses")
        for clause in clauses:
            validate_clause(clause, control_id)
    else:
        raise ValueError(f"{control_id}: unsupported binary classification")


def validate_packet_metadata(packet: dict[str, Any], packet_path: Path) -> None:
    registry = load_registry()
    if packet.get("methodology_sha256") != hashlib.sha256(METHODOLOGY.read_bytes()).hexdigest():
        raise ValueError(f"{packet_path}: packet uses a stale classification methodology")
    if (
        packet.get("registry_version") != registry["registry_version"]
        or packet.get("registry_sha256") != hashlib.sha256(REGISTRY.read_bytes()).hexdigest()
    ):
        raise ValueError(f"{packet_path}: packet uses a stale control registry")
    controls = packet.get("controls")
    if not isinstance(controls, list) or packet.get("control_count") != len(controls):
        raise ValueError(f"{packet_path}: packet control count is invalid")


def validate_primary() -> None:
    packet_paths = sorted(PACKET_ROOT.glob("*.json"))
    if not packet_paths:
        raise ValueError("classification packets are missing; run generate-packets")
    seen: set[str] = set()
    expected_total = 0
    for packet_path in packet_paths:
        packet = json.loads(packet_path.read_text(encoding="utf-8"))
        validate_packet_metadata(packet, packet_path)
        expected = packet["controls"]
        expected_total += len(expected)
        review_path = PRIMARY_ROOT / packet_path.name
        if not review_path.exists():
            raise ValueError(f"missing primary review: {review_path}")
        review = json.loads(review_path.read_text(encoding="utf-8"))
        envelope_fields = {
            "schema_version",
            "packet_id",
            "methodology_sha256",
            "registry_version",
            "registry_sha256",
            "source_path",
            "controls",
        }
        if not isinstance(review, dict) or set(review) != envelope_fields:
            raise ValueError(f"{review_path}: review envelope fields do not match the contract")
        if (
            review.get("schema_version") != "prc.control-classification-review/v0.1"
            or review.get("packet_id") != packet["packet_id"]
            or review.get("methodology_sha256") != packet["methodology_sha256"]
            or review.get("registry_version") != packet["registry_version"]
            or review.get("registry_sha256") != packet["registry_sha256"]
            or review.get("source_path") != packet["source_path"]
        ):
            raise ValueError(f"{review_path}: review envelope does not match packet")
        rows = review.get("controls")
        if not isinstance(rows, list) or len(rows) != len(expected):
            raise ValueError(f"{review_path}: expected {len(expected)} reviews")
        for row, source in zip(rows, expected):
            validate_review(row, source)
            if row["control_id"] in seen:
                raise ValueError(f"duplicate review for {row['control_id']}")
            seen.add(row["control_id"])
    registry = load_registry()
    if expected_total != registry["entry_count"] or len(seen) != registry["entry_count"]:
        raise ValueError(f"review count {len(seen)} does not match registry {registry['entry_count']}")
    print(f"validated one primary classification for all {len(seen)} controls")


def validate_packet(packet_name: str, *, verbose: bool = True) -> None:
    packet_path = PACKET_ROOT / packet_name
    if packet_path.suffix != ".json" or packet_path.parent != PACKET_ROOT:
        raise ValueError("packet must be a JSON filename from the packet directory")
    if not packet_path.exists():
        raise ValueError(f"classification packet does not exist: {packet_path}")
    packet = json.loads(packet_path.read_text(encoding="utf-8"))
    validate_packet_metadata(packet, packet_path)
    review_path = PRIMARY_ROOT / packet_path.name
    if not review_path.exists():
        raise ValueError(f"missing primary review: {review_path}")
    review = json.loads(review_path.read_text(encoding="utf-8"))
    envelope_fields = {
        "schema_version",
        "packet_id",
        "methodology_sha256",
        "registry_version",
        "registry_sha256",
        "source_path",
        "controls",
    }
    if not isinstance(review, dict) or set(review) != envelope_fields:
        raise ValueError(f"{review_path}: review envelope fields do not match the contract")
    if (
        review.get("schema_version") != "prc.control-classification-review/v0.1"
        or review.get("packet_id") != packet["packet_id"]
        or review.get("methodology_sha256") != packet["methodology_sha256"]
        or review.get("registry_version") != packet["registry_version"]
        or review.get("registry_sha256") != packet["registry_sha256"]
        or review.get("source_path") != packet["source_path"]
    ):
        raise ValueError(f"{review_path}: review envelope does not match packet")
    expected = packet["controls"]
    rows = review.get("controls")
    if not isinstance(rows, list) or len(rows) != len(expected):
        raise ValueError(f"{review_path}: expected {len(expected)} reviews")
    seen: set[str] = set()
    for row, source in zip(rows, expected):
        validate_review(row, source)
        if row["control_id"] in seen:
            raise ValueError(f"duplicate review for {row['control_id']}")
        seen.add(row["control_id"])
    if verbose:
        print(f"validated {len(seen)} primary classifications in {packet_path.name}")


def show_progress() -> None:
    packet_paths = sorted(PACKET_ROOT.glob("*.json"))
    reviewed_packets = 0
    reviewed_controls = 0
    deterministic = 0
    nondeterministic = 0
    invalid: list[str] = []
    for packet_path in packet_paths:
        review_path = PRIMARY_ROOT / packet_path.name
        if not review_path.exists():
            continue
        try:
            validate_packet(packet_path.name, verbose=False)
        except (OSError, ValueError, json.JSONDecodeError) as error:
            invalid.append(f"{packet_path.name}: {error}")
            continue
        review = json.loads(review_path.read_text(encoding="utf-8"))
        reviewed_packets += 1
        reviewed_controls += len(review["controls"])
        deterministic += sum(row["classification"] == "deterministic" for row in review["controls"])
        nondeterministic += sum(row["classification"] == "nondeterministic" for row in review["controls"])
    total_controls = sum(
        json.loads(path.read_text(encoding="utf-8"))["control_count"] for path in packet_paths
    )
    print(
        f"progress: {reviewed_packets}/{len(packet_paths)} packets; "
        f"{reviewed_controls}/{total_controls} controls; "
        f"{deterministic} deterministic; {nondeterministic} nondeterministic"
    )
    for error in invalid:
        print(f"invalid: {error}", file=sys.stderr)
    if invalid:
        raise ValueError(f"{len(invalid)} completed packet reviews are invalid")


def generate_skeptic_packets() -> None:
    validate_primary()
    proposals_by_packet: dict[str, list[dict[str, Any]]] = {}
    for packet_path in sorted(PACKET_ROOT.glob("*.json")):
        packet = json.loads(packet_path.read_text(encoding="utf-8"))
        review = json.loads((PRIMARY_ROOT / packet_path.name).read_text(encoding="utf-8"))
        proposals: list[dict[str, Any]] = []
        for source, row in zip(packet["controls"], review["controls"]):
            if row["classification"] != "deterministic":
                continue
            proposals.append(
                {
                    "control_id": source["control_id"],
                    "revision": source["revision"],
                    "semantic_sha256": source["semantic_sha256"],
                    "statement": source["statement"],
                    "source": source["source"],
                    "heading_trail": source["heading_trail"],
                    "primary_packet_id": packet["packet_id"],
                    "primary_route": row["route"],
                    "primary_reason": row["reason"],
                    "primary_deterministic_clauses": row["deterministic_clauses"],
                }
            )
        if proposals:
            proposals_by_packet[packet["packet_id"]] = proposals
    SKEPTIC_PACKET_ROOT.mkdir(parents=True, exist_ok=True)
    expected_names: set[str] = set()
    registry = load_registry()
    proposal_total = 0
    for primary_packet_id, proposals in proposals_by_packet.items():
        proposal_total += len(proposals)
        for chunk_index in range(0, len(proposals), MAX_SKEPTIC_CONTROLS):
            chunk = proposals[chunk_index : chunk_index + MAX_SKEPTIC_CONTROLS]
            part_number = chunk_index // MAX_SKEPTIC_CONTROLS + 1
            packet_id = f"skeptic-{primary_packet_id}-part-{part_number:03d}"
            name = packet_id + ".json"
            expected_names.add(name)
            document = {
                "schema_version": "prc.control-classification-skeptic-packet/v0.1",
                "packet_id": packet_id,
                "primary_packet_id": primary_packet_id,
                "methodology_sha256": hashlib.sha256(METHODOLOGY.read_bytes()).hexdigest(),
                "registry_version": registry["registry_version"],
                "registry_sha256": hashlib.sha256(REGISTRY.read_bytes()).hexdigest(),
                "control_count": len(chunk),
                "controls": chunk,
            }
            (SKEPTIC_PACKET_ROOT / name).write_text(
                json.dumps(document, indent=2, ensure_ascii=False) + "\n", encoding="utf-8"
            )
    for path in SKEPTIC_PACKET_ROOT.glob("*.json"):
        if path.name not in expected_names:
            path.unlink()
    print(f"generated {len(expected_names)} skeptical-review packets for {proposal_total} deterministic proposals")


def validate_skeptic_review(review: Any, expected: dict[str, Any]) -> None:
    fields = {
        "control_id",
        "revision",
        "semantic_sha256",
        "verdict",
        "reason",
        "counterexample_analysis",
        "nondeterministic_route",
        "nondeterministic_reason_codes",
        "review_status",
    }
    control_id = expected["control_id"]
    if not isinstance(review, dict) or set(review) != fields:
        raise ValueError(f"{control_id}: skeptical review fields do not match the contract")
    if (
        review["control_id"] != control_id
        or review["revision"] != expected["revision"]
        or review["semantic_sha256"] != expected["semantic_sha256"]
    ):
        raise ValueError(f"{control_id}: skeptical review is not bound to the current rule revision")
    if not isinstance(review["reason"], str) or len(review["reason"].strip()) < 20:
        raise ValueError(f"{control_id}: skeptical-review reason is too short")
    if not isinstance(review["counterexample_analysis"], str) or len(review["counterexample_analysis"].strip()) < 20:
        raise ValueError(f"{control_id}: counterexample analysis is too short")
    if review["review_status"] != "skeptical_agent_reviewed":
        raise ValueError(f"{control_id}: skeptical review has an invalid status")
    reasons = review["nondeterministic_reason_codes"]
    if not isinstance(reasons, list) or len(reasons) != len(set(reasons)):
        raise ValueError(f"{control_id}: skeptical reason codes are invalid")
    if review["verdict"] == "confirmed_deterministic":
        if review["nondeterministic_route"] is not None or reasons:
            raise ValueError(f"{control_id}: confirmed proposal contains rejection fields")
    elif review["verdict"] == "rejected_nondeterministic":
        if review["nondeterministic_route"] not in NONDETERMINISTIC_ROUTES or not reasons:
            raise ValueError(f"{control_id}: rejected proposal lacks a nondeterministic classification")
        if any(reason not in NONDETERMINISTIC_REASONS for reason in reasons):
            raise ValueError(f"{control_id}: unsupported skeptical reason code")
    else:
        raise ValueError(f"{control_id}: unsupported skeptical verdict")


def validate_skeptic_packet(packet_name: str, *, verbose: bool = True) -> int:
    packet_path = SKEPTIC_PACKET_ROOT / packet_name
    if packet_path.suffix != ".json" or packet_path.parent != SKEPTIC_PACKET_ROOT:
        raise ValueError("skeptical packet must be a JSON filename from the skeptical packet directory")
    if not packet_path.exists():
        raise ValueError(f"skeptical classification packet does not exist: {packet_path}")
    packet = json.loads(packet_path.read_text(encoding="utf-8"))
    validate_packet_metadata(packet, packet_path)
    review_path = SKEPTIC_ROOT / packet_path.name
    if not review_path.exists():
        raise ValueError(f"missing skeptical review: {review_path}")
    review = json.loads(review_path.read_text(encoding="utf-8"))
    envelope_fields = {
        "schema_version",
        "packet_id",
        "primary_packet_id",
        "methodology_sha256",
        "registry_version",
        "registry_sha256",
        "controls",
    }
    if not isinstance(review, dict) or set(review) != envelope_fields:
        raise ValueError(f"{review_path}: skeptical review envelope fields do not match the contract")
    if (
        review.get("schema_version") != "prc.control-classification-skeptic-review/v0.1"
        or review.get("packet_id") != packet["packet_id"]
        or review.get("primary_packet_id") != packet["primary_packet_id"]
        or review.get("methodology_sha256") != packet["methodology_sha256"]
        or review.get("registry_version") != packet["registry_version"]
        or review.get("registry_sha256") != packet["registry_sha256"]
    ):
        raise ValueError(f"{review_path}: skeptical review envelope does not match packet")
    expected = packet["controls"]
    rows = review.get("controls")
    if not isinstance(rows, list) or len(rows) != len(expected):
        raise ValueError(f"{review_path}: expected {len(expected)} skeptical reviews")
    seen: set[str] = set()
    for row, source in zip(rows, expected):
        validate_skeptic_review(row, source)
        if row["control_id"] in seen:
            raise ValueError(f"duplicate skeptical review for {row['control_id']}")
        seen.add(row["control_id"])
    if verbose:
        print(f"validated {len(seen)} skeptical classifications in {packet_path.name}")
    return len(seen)


def validate_skeptic() -> None:
    packet_paths = sorted(SKEPTIC_PACKET_ROOT.glob("*.json"))
    if not packet_paths:
        raise ValueError("skeptical classification packets are missing; run generate-skeptic-packets")
    total = 0
    seen: set[str] = set()
    for packet_path in packet_paths:
        validate_skeptic_packet(packet_path.name, verbose=False)
        review = json.loads((SKEPTIC_ROOT / packet_path.name).read_text(encoding="utf-8"))
        for row in review["controls"]:
            if row["control_id"] in seen:
                raise ValueError(f"duplicate skeptical review for {row['control_id']}")
            seen.add(row["control_id"])
        total += len(review["controls"])
    proposal_total = sum(
        json.loads(path.read_text(encoding="utf-8"))["control_count"] for path in packet_paths
    )
    if total != proposal_total or len(seen) != proposal_total:
        raise ValueError("skeptical review coverage does not match deterministic proposals")
    print(f"validated one skeptical classification for all {total} deterministic proposals")


def show_skeptic_progress() -> None:
    packet_paths = sorted(SKEPTIC_PACKET_ROOT.glob("*.json"))
    reviewed_packets = 0
    reviewed_controls = 0
    confirmed = 0
    rejected = 0
    invalid: list[str] = []
    for packet_path in packet_paths:
        review_path = SKEPTIC_ROOT / packet_path.name
        if not review_path.exists():
            continue
        try:
            validate_skeptic_packet(packet_path.name, verbose=False)
        except (OSError, ValueError, json.JSONDecodeError) as error:
            invalid.append(f"{packet_path.name}: {error}")
            continue
        review = json.loads(review_path.read_text(encoding="utf-8"))
        reviewed_packets += 1
        reviewed_controls += len(review["controls"])
        confirmed += sum(row["verdict"] == "confirmed_deterministic" for row in review["controls"])
        rejected += sum(row["verdict"] == "rejected_nondeterministic" for row in review["controls"])
    total_controls = sum(
        json.loads(path.read_text(encoding="utf-8"))["control_count"] for path in packet_paths
    )
    print(
        f"skeptical progress: {reviewed_packets}/{len(packet_paths)} packets; "
        f"{reviewed_controls}/{total_controls} proposals; "
        f"{confirmed} confirmed; {rejected} rejected"
    )
    for error in invalid:
        print(f"invalid: {error}", file=sys.stderr)
    if invalid:
        raise ValueError(f"{len(invalid)} completed skeptical reviews are invalid")


def load_strength_reviews(
    expected_order: list[str],
    expected_by_id: dict[str, dict[str, Any]],
) -> dict[str, dict[str, Any]]:
    paths = sorted(STRENGTH_ROOT.glob("part-*.json"))
    if not paths:
        raise ValueError("semantic-strength reviews are missing")
    expected_envelope = {
        "schema_version", "review_id", "methodology_sha256", "registry_version",
        "registry_sha256", "control_count", "controls",
    }
    expected_fields = {
        "control_id", "revision", "semantic_sha256", "original_statement",
        "old_deterministic_clauses", "verdict", "route", "reason",
        "deterministic_clauses", "nondeterministic_reason_codes",
        "counterexample_analysis", "review_status", "row_sha256",
    }
    registry = load_registry()
    methodology_sha256 = hashlib.sha256(METHODOLOGY.read_bytes()).hexdigest()
    registry_sha256 = hashlib.sha256(REGISTRY.read_bytes()).hexdigest()
    rows: list[dict[str, Any]] = []
    seen: set[str] = set()
    for part, path in enumerate(paths, start=1):
        document = json.loads(path.read_text(encoding="utf-8"))
        if not isinstance(document, dict) or set(document) != expected_envelope:
            raise ValueError(f"{path}: strength-review envelope fields do not match the contract")
        if (
            document["schema_version"] != "prc.control-classification-strength-review/v0.1"
            or document["review_id"] != f"deterministic-strength-part-{part:03d}"
            or document["methodology_sha256"] != methodology_sha256
            or document["registry_version"] != registry["registry_version"]
            or document["registry_sha256"] != registry_sha256
            or document["control_count"] != len(document["controls"])
        ):
            raise ValueError(f"{path}: strength-review envelope is stale or invalid")
        for row in document["controls"]:
            if not isinstance(row, dict) or set(row) != expected_fields:
                raise ValueError(f"{path}: strength-review row fields do not match the contract")
            control_id = row["control_id"]
            expected = expected_by_id.get(control_id)
            if expected is None or control_id in seen:
                raise ValueError(f"{path}: unknown or duplicate strength-review control {control_id!r}")
            source = expected["source"]
            if (
                row["revision"] != source["revision"]
                or row["semantic_sha256"] != source["semantic_sha256"]
                or row["original_statement"] != source["statement"]
                or row["old_deterministic_clauses"] != expected["clauses"]
                or row["review_status"] != "strength_audited"
            ):
                raise ValueError(f"{control_id}: strength review is not bound to the current baseline")
            unsigned = {key: value for key, value in row.items() if key != "row_sha256"}
            if row["row_sha256"] != digest_value(unsigned):
                raise ValueError(f"{control_id}: strength-review row digest is stale")
            if not isinstance(row["reason"], str) or len(row["reason"].strip()) < 20:
                raise ValueError(f"{control_id}: strength-review reason is too short")
            if not isinstance(row["counterexample_analysis"], str) or len(row["counterexample_analysis"].strip()) < 20:
                raise ValueError(f"{control_id}: strength-review counterexample is too short")
            clauses = row["deterministic_clauses"]
            reason_codes = row["nondeterministic_reason_codes"]
            if not isinstance(clauses, list) or not isinstance(reason_codes, list) or len(reason_codes) != len(set(reason_codes)):
                raise ValueError(f"{control_id}: strength-review clauses or reason codes are invalid")
            if row["verdict"] == "strength_audit_confirmed":
                if row["route"] not in DETERMINISTIC_ROUTES or not clauses or reason_codes:
                    raise ValueError(f"{control_id}: confirmed strength review is incomplete")
                for clause in clauses:
                    validate_strengthened_clause(clause, control_id)
            elif row["verdict"] == "strength_audit_reclassified":
                if row["route"] not in NONDETERMINISTIC_ROUTES or clauses or not reason_codes:
                    raise ValueError(f"{control_id}: reclassified strength review is incomplete")
                if any(code not in NONDETERMINISTIC_REASONS for code in reason_codes):
                    raise ValueError(f"{control_id}: strength review has unsupported reason code")
            else:
                raise ValueError(f"{control_id}: unsupported strength-review verdict")
            seen.add(control_id)
            rows.append(row)
    actual_order = [row["control_id"] for row in rows]
    if actual_order != expected_order:
        missing = [control_id for control_id in expected_order if control_id not in seen]
        first = missing[0] if missing else actual_order[0] if actual_order else None
        raise ValueError(f"strength reviews do not cover deterministic controls in exact corpus order; first mismatch {first}")
    return {row["control_id"]: row for row in rows}


def generate_final() -> None:
    validate_primary()
    validate_skeptic()
    skeptical_by_id: dict[str, dict[str, Any]] = {}
    for path in sorted(SKEPTIC_ROOT.glob("*.json")):
        document = json.loads(path.read_text(encoding="utf-8"))
        for row in document["controls"]:
            skeptical_by_id[row["control_id"]] = row
    baseline_by_id: dict[str, dict[str, Any]] = {}
    strength_expected_order: list[str] = []
    strength_expected_by_id: dict[str, dict[str, Any]] = {}
    for packet_path in sorted(PACKET_ROOT.glob("*.json")):
        packet = json.loads(packet_path.read_text(encoding="utf-8"))
        primary = json.loads((PRIMARY_ROOT / packet_path.name).read_text(encoding="utf-8"))
        for source, row in zip(packet["controls"], primary["controls"]):
            if row["classification"] == "deterministic":
                skeptical = skeptical_by_id.get(row["control_id"])
                if skeptical is None:
                    raise ValueError(f"missing skeptical decision for {row['control_id']}")
                if skeptical["verdict"] == "confirmed_deterministic":
                    baseline = {
                        "classification": "deterministic", "route": row["route"],
                        "reason": skeptical["reason"], "clauses": row["deterministic_clauses"],
                        "reason_codes": [], "decision_basis": "two_stage_confirmed",
                        "skeptical_verdict": skeptical["verdict"],
                        "counterexample_analysis": skeptical["counterexample_analysis"],
                    }
                    strength_expected_order.append(row["control_id"])
                    strength_expected_by_id[row["control_id"]] = {
                        "source": source, "clauses": row["deterministic_clauses"]
                    }
                else:
                    baseline = {
                        "classification": "nondeterministic", "route": skeptical["nondeterministic_route"],
                        "reason": skeptical["reason"], "clauses": [],
                        "reason_codes": skeptical["nondeterministic_reason_codes"],
                        "decision_basis": "skeptically_rejected", "skeptical_verdict": skeptical["verdict"],
                        "counterexample_analysis": skeptical["counterexample_analysis"],
                    }
            else:
                baseline = {
                    "classification": "nondeterministic", "route": row["route"],
                    "reason": row["reason"], "clauses": row["deterministic_clauses"],
                    "reason_codes": row["nondeterministic_reason_codes"],
                    "decision_basis": "primary_nondeterministic", "skeptical_verdict": None,
                    "counterexample_analysis": None,
                }
            baseline_by_id[row["control_id"]] = baseline
    strength_by_id = load_strength_reviews(strength_expected_order, strength_expected_by_id)
    FINAL_ROOT.mkdir(parents=True, exist_ok=True)
    expected_names: set[str] = set()
    counts = {"deterministic": 0, "nondeterministic": 0}
    decision_counts: dict[str, int] = defaultdict(int)
    route_counts: dict[str, int] = defaultdict(int)
    packet_summaries: list[dict[str, Any]] = []
    for packet_path in sorted(PACKET_ROOT.glob("*.json")):
        packet = json.loads(packet_path.read_text(encoding="utf-8"))
        primary = json.loads((PRIMARY_ROOT / packet_path.name).read_text(encoding="utf-8"))
        final_rows: list[dict[str, Any]] = []
        packet_counts = {"deterministic": 0, "nondeterministic": 0}
        for source, row in zip(packet["controls"], primary["controls"]):
            baseline = baseline_by_id[row["control_id"]]
            strength = strength_by_id.get(row["control_id"])
            if strength is not None:
                classification = "deterministic" if strength["verdict"] == "strength_audit_confirmed" else "nondeterministic"
                route = strength["route"]
                reason = strength["reason"]
                clauses = strength["deterministic_clauses"]
                reason_codes = strength["nondeterministic_reason_codes"]
                decision_basis = strength["verdict"]
                skeptical_verdict = baseline["skeptical_verdict"]
                counterexample_analysis = strength["counterexample_analysis"]
            else:
                classification = baseline["classification"]
                route = baseline["route"]
                reason = baseline["reason"]
                clauses = baseline["clauses"]
                reason_codes = baseline["reason_codes"]
                decision_basis = baseline["decision_basis"]
                skeptical_verdict = baseline["skeptical_verdict"]
                counterexample_analysis = baseline["counterexample_analysis"]
            final_rows.append(
                {
                    "control_id": source["control_id"],
                    "revision": source["revision"],
                    "semantic_sha256": source["semantic_sha256"],
                    "classification": classification,
                    "route": route,
                    "reason": reason,
                    "deterministic_clauses": clauses,
                    "nondeterministic_reason_codes": reason_codes,
                    "decision_basis": decision_basis,
                    "skeptical_verdict": skeptical_verdict,
                    "counterexample_analysis": counterexample_analysis,
                }
            )
            counts[classification] += 1
            packet_counts[classification] += 1
            decision_counts[decision_basis] += 1
            route_counts[route] += 1
        expected_names.add(packet_path.name)
        final_document = {
            "schema_version": "prc.control-classification-final/v0.1",
            "packet_id": packet["packet_id"],
            "methodology_sha256": packet["methodology_sha256"],
            "registry_version": packet["registry_version"],
            "registry_sha256": packet["registry_sha256"],
            "source_path": packet["source_path"],
            "control_count": len(final_rows),
            "controls": final_rows,
        }
        (FINAL_ROOT / packet_path.name).write_text(
            json.dumps(final_document, indent=2, ensure_ascii=False) + "\n", encoding="utf-8"
        )
        packet_summaries.append(
            {
                "packet_id": packet["packet_id"],
                "source_path": packet["source_path"],
                "control_count": len(final_rows),
                "deterministic": packet_counts["deterministic"],
                "nondeterministic": packet_counts["nondeterministic"],
            }
        )
    for path in FINAL_ROOT.glob("*.json"):
        if path.name not in expected_names:
            path.unlink()
    registry = load_registry()
    summary = {
        "schema_version": "prc.control-classification-summary/v0.1",
        "methodology_sha256": hashlib.sha256(METHODOLOGY.read_bytes()).hexdigest(),
        "registry_version": registry["registry_version"],
        "registry_sha256": hashlib.sha256(REGISTRY.read_bytes()).hexdigest(),
        "control_count": counts["deterministic"] + counts["nondeterministic"],
        "deterministic": counts["deterministic"],
        "nondeterministic": counts["nondeterministic"],
        "decision_counts": dict(sorted(decision_counts.items())),
        "route_counts": dict(sorted(route_counts.items())),
        "packets": packet_summaries,
    }
    SUMMARY.write_text(json.dumps(summary, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")
    print(
        f"generated final classifications for {summary['control_count']} controls: "
        f"{counts['deterministic']} deterministic and {counts['nondeterministic']} nondeterministic"
    )


def validate_final() -> None:
    packet_paths = sorted(PACKET_ROOT.glob("*.json"))
    seen: set[str] = set()
    counts = {"deterministic": 0, "nondeterministic": 0}
    fields = {
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
    for packet_path in packet_paths:
        packet = json.loads(packet_path.read_text(encoding="utf-8"))
        final_path = FINAL_ROOT / packet_path.name
        if not final_path.exists():
            raise ValueError(f"missing final classification: {final_path}")
        document = json.loads(final_path.read_text(encoding="utf-8"))
        envelope_fields = {
            "schema_version",
            "packet_id",
            "methodology_sha256",
            "registry_version",
            "registry_sha256",
            "source_path",
            "control_count",
            "controls",
        }
        if not isinstance(document, dict) or set(document) != envelope_fields:
            raise ValueError(f"{final_path}: final envelope fields do not match the contract")
        if (
            document["schema_version"] != "prc.control-classification-final/v0.1"
            or document["packet_id"] != packet["packet_id"]
            or document["methodology_sha256"] != packet["methodology_sha256"]
            or document["registry_version"] != packet["registry_version"]
            or document["registry_sha256"] != packet["registry_sha256"]
            or document["source_path"] != packet["source_path"]
        ):
            raise ValueError(f"{final_path}: final envelope does not match packet")
        rows = document["controls"]
        if not isinstance(rows, list) or document["control_count"] != len(rows) or len(rows) != len(packet["controls"]):
            raise ValueError(f"{final_path}: final control count does not match packet")
        for row, source in zip(rows, packet["controls"]):
            control_id = source["control_id"]
            if not isinstance(row, dict) or set(row) != fields:
                raise ValueError(f"{control_id}: final classification fields do not match the contract")
            if (
                row["control_id"] != control_id
                or row["revision"] != source["revision"]
                or row["semantic_sha256"] != source["semantic_sha256"]
            ):
                raise ValueError(f"{control_id}: final classification uses a stale rule revision")
            if control_id in seen:
                raise ValueError(f"duplicate final classification for {control_id}")
            seen.add(control_id)
            if not isinstance(row["reason"], str) or len(row["reason"].strip()) < 20:
                raise ValueError(f"{control_id}: final reason is too short")
            clauses = row["deterministic_clauses"]
            reasons = row["nondeterministic_reason_codes"]
            if not isinstance(clauses, list) or not isinstance(reasons, list):
                raise ValueError(f"{control_id}: final clauses or reason codes are invalid")
            for clause in clauses:
                validate_clause(clause, control_id)
            if row["classification"] == "deterministic":
                for clause in clauses:
                    validate_strengthened_clause(clause, control_id)
                if (
                    row["route"] not in DETERMINISTIC_ROUTES
                    or not clauses
                    or reasons
                    or row["decision_basis"] not in {"two_stage_confirmed", "strength_audit_confirmed"}
                    or row["skeptical_verdict"] != "confirmed_deterministic"
                    or not isinstance(row["counterexample_analysis"], str)
                    or len(row["counterexample_analysis"].strip()) < 20
                ):
                    raise ValueError(f"{control_id}: deterministic final classification lacks two-stage proof")
            elif row["classification"] == "nondeterministic":
                if row["route"] not in NONDETERMINISTIC_ROUTES or not reasons:
                    raise ValueError(f"{control_id}: nondeterministic final classification is incomplete")
                if any(reason not in NONDETERMINISTIC_REASONS for reason in reasons):
                    raise ValueError(f"{control_id}: final nondeterministic reason code is unsupported")
                if row["decision_basis"] == "primary_nondeterministic":
                    if row["skeptical_verdict"] is not None or row["counterexample_analysis"] is not None:
                        raise ValueError(f"{control_id}: primary nondeterministic result contains skeptical fields")
                    if row["route"] != "mixed" and clauses:
                        raise ValueError(f"{control_id}: final nondeterministic result retains unsafe clauses")
                elif row["decision_basis"] == "skeptically_rejected":
                    if (
                        row["skeptical_verdict"] != "rejected_nondeterministic"
                        or not isinstance(row["counterexample_analysis"], str)
                        or len(row["counterexample_analysis"].strip()) < 20
                        or clauses
                    ):
                        raise ValueError(f"{control_id}: rejected final classification is incomplete")
                elif row["decision_basis"] == "strength_audit_reclassified":
                    if (
                        row["skeptical_verdict"] != "confirmed_deterministic"
                        or not isinstance(row["counterexample_analysis"], str)
                        or len(row["counterexample_analysis"].strip()) < 20
                        or clauses
                    ):
                        raise ValueError(f"{control_id}: strength-reclassified final classification is incomplete")
                else:
                    raise ValueError(f"{control_id}: nondeterministic final decision basis is invalid")
            else:
                raise ValueError(f"{control_id}: unsupported final classification")
            counts[row["classification"]] += 1
    registry = load_registry()
    if len(seen) != registry["entry_count"]:
        raise ValueError(f"final count {len(seen)} does not match registry {registry['entry_count']}")
    summary = json.loads(SUMMARY.read_text(encoding="utf-8"))
    if (
        summary.get("control_count") != len(seen)
        or summary.get("deterministic") != counts["deterministic"]
        or summary.get("nondeterministic") != counts["nondeterministic"]
    ):
        raise ValueError("classification summary counts do not match final packets")
    print(
        f"validated final classifications for {len(seen)} controls: "
        f"{counts['deterministic']} deterministic and {counts['nondeterministic']} nondeterministic"
    )


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "command",
        choices=(
            "generate-packets",
            "generate-final",
            "generate-skeptic-packets",
            "progress",
            "skeptic-progress",
            "validate-packet",
            "validate-final",
            "validate-primary",
            "validate-skeptic-packet",
            "validate-skeptic",
        ),
    )
    parser.add_argument("packet", nargs="?")
    args = parser.parse_args()
    try:
        if args.command == "generate-packets":
            if args.packet:
                raise ValueError("generate-packets does not accept a packet name")
            generate_packets()
        elif args.command == "progress":
            if args.packet:
                raise ValueError("progress does not accept a packet name")
            show_progress()
        elif args.command == "skeptic-progress":
            if args.packet:
                raise ValueError("skeptic-progress does not accept a packet name")
            show_skeptic_progress()
        elif args.command == "generate-skeptic-packets":
            if args.packet:
                raise ValueError("generate-skeptic-packets does not accept a packet name")
            generate_skeptic_packets()
        elif args.command == "generate-final":
            if args.packet:
                raise ValueError("generate-final does not accept a packet name")
            generate_final()
        elif args.command == "validate-packet":
            if not args.packet:
                raise ValueError("validate-packet requires a packet filename")
            validate_packet(args.packet)
        elif args.command == "validate-skeptic-packet":
            if not args.packet:
                raise ValueError("validate-skeptic-packet requires a packet filename")
            validate_skeptic_packet(args.packet)
        elif args.command == "validate-skeptic":
            if args.packet:
                raise ValueError("validate-skeptic does not accept a packet name")
            validate_skeptic()
        elif args.command == "validate-final":
            if args.packet:
                raise ValueError("validate-final does not accept a packet name")
            validate_final()
        else:
            if args.packet:
                raise ValueError("validate-primary does not accept a packet name")
            validate_primary()
    except (OSError, ValueError, json.JSONDecodeError) as error:
        print(f"control classification review failed: {error}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
