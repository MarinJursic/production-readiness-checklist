#!/usr/bin/env python3
"""Build the human-readable reference for the final control classifications."""

from __future__ import annotations

import argparse
import hashlib
import json
import os
import sys
from collections import Counter, defaultdict
from dataclasses import dataclass
from pathlib import Path
from typing import Any, Iterable


ROOT = Path(__file__).resolve().parents[1]
FINAL_ROOT = ROOT / "research" / "control-classification" / "final"
PACKET_ROOT = ROOT / "research" / "control-classification" / "packets"
SUMMARY = ROOT / "research" / "control-classification" / "summary.json"
REGISTRY = ROOT / "catalog" / "control-id-registry.json"
OUTPUT_ROOT = ROOT / "docs" / "control-classification"

# Parts are deliberately kept well below the user's 3 MB ceiling.
TARGET_PART_BYTES = 1_500_000
MAX_PART_BYTES = 2_500_000

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
DECISION_BASIS = {
    "primary_nondeterministic": (
        "The primary review found that the complete rule needs judgment or evidence "
        "that cannot produce a closed automatic pass/fail result."
    ),
    "skeptically_rejected": (
        "A primary deterministic proposal failed the independent counterexample review, "
        "so the final rule is nondeterministic."
    ),
    "two_stage_confirmed": (
        "The primary deterministic proposal passed an independent skeptical review of "
        "scope, authority, pass/fail directions, and counterexamples."
    ),
    "strength_audit_confirmed": (
        "A third rule-by-rule strength audit confirmed that the exact clauses preserve "
        "the original promise and remain mechanically decidable."
    ),
    "strength_audit_reclassified": (
        "A third rule-by-rule strength audit found that the earlier deterministic clause "
        "weakened the original promise, so the complete rule is nondeterministic."
    ),
}
FINAL_CONTROL_FIELDS = {
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
CLAUSE_FIELDS = {"statement", "checker_family", "evidence_authority"}


@dataclass(frozen=True)
class Control:
    """A final classification joined to its frozen exact source statement."""

    packet_id: str
    source_path: str
    statement: str
    row: dict[str, Any]


@dataclass(frozen=True)
class Corpus:
    controls: tuple[Control, ...]
    source_order: tuple[str, ...]
    summary: dict[str, Any]


@dataclass(frozen=True)
class Part:
    path: Path
    source_path: str
    controls: tuple[Control, ...]
    content: str


def load_json(path: Path) -> dict[str, Any]:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
    except FileNotFoundError as error:
        raise ValueError(f"missing classification input: {path}") from error
    if not isinstance(value, dict):
        raise ValueError(f"{path}: expected a JSON object")
    return value


def _require_text(value: Any, label: str) -> str:
    if not isinstance(value, str) or not value.strip():
        raise ValueError(f"{label} must be nonempty text")
    return value


def _validate_clause(clause: Any, control_id: str) -> None:
    if not isinstance(clause, dict) or set(clause) != CLAUSE_FIELDS:
        raise ValueError(f"{control_id}: deterministic clause fields do not match the contract")
    for field in sorted(CLAUSE_FIELDS):
        _require_text(clause[field], f"{control_id} clause {field}")
    if clause["checker_family"] not in CHECKER_FAMILIES:
        raise ValueError(f"{control_id}: unsupported checker family")
    if clause["evidence_authority"] not in EVIDENCE_AUTHORITIES:
        raise ValueError(f"{control_id}: unsupported evidence authority")


def _validate_final_row(row: Any, source: dict[str, Any]) -> None:
    control_id = source.get("control_id", "unknown control")
    if not isinstance(row, dict) or set(row) != FINAL_CONTROL_FIELDS:
        raise ValueError(f"{control_id}: final classification fields do not match the contract")
    if (
        row["control_id"] != source.get("control_id")
        or row["revision"] != source.get("revision")
        or row["semantic_sha256"] != source.get("semantic_sha256")
    ):
        raise ValueError(f"{control_id}: final classification is stale or bound to another rule")
    reason = _require_text(row["reason"], f"{control_id} reason")
    if len(reason.strip()) < 20:
        raise ValueError(f"{control_id}: final reason is too short")
    clauses = row["deterministic_clauses"]
    reason_codes = row["nondeterministic_reason_codes"]
    if not isinstance(clauses, list) or not isinstance(reason_codes, list):
        raise ValueError(f"{control_id}: clauses and reason codes must be lists")
    for clause in clauses:
        _validate_clause(clause, control_id)
    classification = row["classification"]
    basis = row["decision_basis"]
    if basis not in DECISION_BASIS:
        raise ValueError(f"{control_id}: unsupported decision basis {basis!r}")
    if classification == "deterministic":
        if (
            row["route"] not in DETERMINISTIC_ROUTES
            or not clauses
            or reason_codes
            or basis != "strength_audit_confirmed"
            or row["skeptical_verdict"] != "confirmed_deterministic"
        ):
            raise ValueError(f"{control_id}: deterministic final result is incomplete")
    elif classification == "nondeterministic":
        if row["route"] not in NONDETERMINISTIC_ROUTES or not reason_codes:
            raise ValueError(f"{control_id}: nondeterministic final result is incomplete")
        if any(code not in NONDETERMINISTIC_REASONS for code in reason_codes):
            raise ValueError(f"{control_id}: unsupported nondeterministic reason code")
        if basis in {"skeptically_rejected", "strength_audit_reclassified"}:
            expected_verdict = (
                "rejected_nondeterministic"
                if basis == "skeptically_rejected"
                else "confirmed_deterministic"
            )
            if row["skeptical_verdict"] != expected_verdict or clauses:
                raise ValueError(f"{control_id}: reviewed rejection verdict is missing")
            counterexample = _require_text(
                row["counterexample_analysis"], f"{control_id} counterexample analysis"
            )
            if len(counterexample.strip()) < 20:
                raise ValueError(f"{control_id}: counterexample analysis is too short")
        elif basis == "primary_nondeterministic":
            if row["skeptical_verdict"] is not None or row["counterexample_analysis"] is not None:
                raise ValueError(f"{control_id}: primary decision contains unexpected skeptical fields")
            if row["route"] != "mixed" and clauses:
                raise ValueError(f"{control_id}: unsafe clauses retained for a nondeterministic rule")
        else:
            raise ValueError(f"{control_id}: nondeterministic decision basis is invalid")
    else:
        raise ValueError(f"{control_id}: unsupported classification {classification!r}")


def load_corpus(
    *,
    final_root: Path = FINAL_ROOT,
    packet_root: Path = PACKET_ROOT,
    summary_path: Path = SUMMARY,
    registry_path: Path = REGISTRY,
) -> Corpus:
    """Load and fail-closed validate the final corpus and exact source statements."""

    summary = load_json(summary_path)
    if summary.get("schema_version") != "prc.control-classification-summary/v0.1":
        raise ValueError(f"{summary_path}: unsupported classification summary schema")
    packet_summaries = summary.get("packets")
    if not isinstance(packet_summaries, list) or not packet_summaries:
        raise ValueError(f"{summary_path}: packet summary is missing")

    registry = load_json(registry_path)
    if (
        registry.get("schema_version") != "prc.control-id-registry/v0.1"
        or registry.get("registry_version") != summary.get("registry_version")
        or registry.get("entry_count") != len(registry.get("entries", []))
        or hashlib.sha256(registry_path.read_bytes()).hexdigest() != summary.get("registry_sha256")
    ):
        raise ValueError("classification registry or its summary binding is stale")
    registry_rows = registry.get("entries")
    if not isinstance(registry_rows, list):
        raise ValueError("classification registry entries must be a list")
    registry_by_id: dict[str, dict[str, Any]] = {}
    for entry in registry_rows:
        if not isinstance(entry, dict) or not isinstance(entry.get("id"), str):
            raise ValueError("classification registry contains an invalid entry")
        if entry["id"] in registry_by_id:
            raise ValueError(f"classification registry contains duplicate ID {entry['id']}")
        registry_by_id[entry["id"]] = entry

    packet_ids: list[str] = []
    source_order: list[str] = []
    for packet in packet_summaries:
        if not isinstance(packet, dict):
            raise ValueError(f"{summary_path}: invalid packet summary")
        packet_id = _require_text(packet.get("packet_id"), "packet_id")
        source_path = _require_text(packet.get("source_path"), f"{packet_id} source_path")
        packet_ids.append(packet_id)
        if source_path not in source_order:
            source_order.append(source_path)
    if len(packet_ids) != len(set(packet_ids)):
        raise ValueError("classification summary contains duplicate packet IDs")

    expected_files = {f"{packet_id}.json" for packet_id in packet_ids}
    actual_final = {path.name for path in final_root.glob("*.json")}
    actual_packets = {path.name for path in packet_root.glob("*.json")}
    missing_final = sorted(expected_files - actual_final)
    unexpected_final = sorted(actual_final - expected_files)
    missing_packets = sorted(expected_files - actual_packets)
    if missing_final or unexpected_final or missing_packets:
        details = []
        if missing_final:
            details.append("missing final packets: " + ", ".join(missing_final))
        if unexpected_final:
            details.append("unexpected final packets: " + ", ".join(unexpected_final))
        if missing_packets:
            details.append("missing source packets: " + ", ".join(missing_packets))
        raise ValueError("; ".join(details))

    controls: list[Control] = []
    seen_controls: set[str] = set()
    classification_counts: Counter[str] = Counter()
    summary_by_packet = {packet["packet_id"]: packet for packet in packet_summaries}
    for packet_id in packet_ids:
        name = f"{packet_id}.json"
        frozen = load_json(packet_root / name)
        final = load_json(final_root / name)
        packet_summary = summary_by_packet[packet_id]
        envelope_pairs = (
            ("packet_id", packet_id),
            ("source_path", packet_summary["source_path"]),
            ("methodology_sha256", summary.get("methodology_sha256")),
            ("registry_version", summary.get("registry_version")),
            ("registry_sha256", summary.get("registry_sha256")),
        )
        if frozen.get("schema_version") != "prc.control-classification-packet/v0.1":
            raise ValueError(f"{packet_root / name}: unsupported packet schema")
        if final.get("schema_version") != "prc.control-classification-final/v0.1":
            raise ValueError(f"{final_root / name}: unsupported final schema")
        for field, expected in envelope_pairs:
            if frozen.get(field) != expected or final.get(field) != expected:
                raise ValueError(f"{name}: stale or mismatched {field}")
        source_rows = frozen.get("controls")
        final_rows = final.get("controls")
        if not isinstance(source_rows, list) or not isinstance(final_rows, list):
            raise ValueError(f"{name}: controls must be lists")
        declared = packet_summary.get("control_count")
        if (
            frozen.get("control_count") != len(source_rows)
            or final.get("control_count") != len(final_rows)
            or len(source_rows) != len(final_rows)
            or declared != len(final_rows)
        ):
            raise ValueError(f"{name}: missing or extra control rows")

        packet_counts: Counter[str] = Counter()
        for source, row in zip(source_rows, final_rows):
            if not isinstance(source, dict):
                raise ValueError(f"{name}: invalid frozen source row")
            control_id = _require_text(source.get("control_id"), f"{name} control_id")
            statement = _require_text(source.get("statement"), f"{control_id} statement")
            if control_id in seen_controls:
                raise ValueError(f"duplicate final classification for {control_id}")
            registry_entry = registry_by_id.get(control_id)
            if not registry_entry or (
                registry_entry.get("status") != "active"
                or registry_entry.get("revision") != source.get("revision")
                or registry_entry.get("semantic_sha256") != source.get("semantic_sha256")
                or registry_entry.get("statement") != statement
                or registry_entry.get("source", {}).get("path") != final["source_path"]
            ):
                raise ValueError(f"{control_id}: source packet is stale against the registry")
            _validate_final_row(row, source)
            seen_controls.add(control_id)
            classification_counts[row["classification"]] += 1
            packet_counts[row["classification"]] += 1
            controls.append(Control(packet_id, final["source_path"], statement, row))
        if (
            packet_summary.get("deterministic") != packet_counts["deterministic"]
            or packet_summary.get("nondeterministic") != packet_counts["nondeterministic"]
        ):
            raise ValueError(f"{name}: packet summary classification counts are stale")

    if summary.get("control_count") != len(controls):
        raise ValueError("classification summary control count is stale")
    if seen_controls != set(registry_by_id):
        missing = sorted(set(registry_by_id) - seen_controls)
        unexpected = sorted(seen_controls - set(registry_by_id))
        raise ValueError(
            "final corpus does not match the active registry; "
            f"missing={missing[:5]}, unexpected={unexpected[:5]}"
        )
    if summary.get("deterministic") != classification_counts["deterministic"]:
        raise ValueError("classification summary deterministic count is stale")
    if summary.get("nondeterministic") != classification_counts["nondeterministic"]:
        raise ValueError("classification summary nondeterministic count is stale")
    return Corpus(tuple(controls), tuple(source_order), summary)


def _blockquote(text: str) -> str:
    return "\n".join("> " + line for line in text.splitlines())


def render_control(control: Control) -> str:
    row = control.row
    label = "Deterministic" if row["classification"] == "deterministic" else "Nondeterministic"
    lines = [
        f"## `{row['control_id']}` — {label}",
        "",
        "**Exact statement**",
        "",
        _blockquote(control.statement),
        "",
        f"- **Classification:** **{label}**",
        f"- **Route:** `{row['route']}`",
        f"- **Decision basis:** `{row['decision_basis']}` — {DECISION_BASIS[row['decision_basis']]}",
        f"- **Reason:** {row['reason']}",
    ]
    reason_codes = row["nondeterministic_reason_codes"]
    if reason_codes:
        lines.append("- **Nondeterministic reason codes:** " + ", ".join(f"`{code}`" for code in reason_codes))

    clauses = row["deterministic_clauses"]
    if row["classification"] == "deterministic":
        lines.extend(["", "### Exact deterministic check"])
    elif clauses:
        lines.extend(
            [
                "",
                "### Partial deterministic clauses",
                "",
                "These clauses can be checked, but they do **not** prove the complete nondeterministic rule.",
            ]
        )
    if clauses:
        for index, clause in enumerate(clauses, start=1):
            lines.extend(
                [
                    "",
                    f"{index}. {clause['statement']}",
                    f"   - **Checker family:** `{clause['checker_family']}`",
                    f"   - **Evidence authority:** `{clause['evidence_authority']}`",
                ]
            )

    counterexample = row["counterexample_analysis"]
    if counterexample:
        heading = {
            "skeptically_rejected": "Skeptical rejection counterexample",
            "strength_audit_reclassified": "Strength-audit rejection counterexample",
        }.get(row["decision_basis"], "Skeptical counterexample test")
        lines.extend(["", f"### {heading}", "", counterexample])
    return "\n".join(lines).rstrip() + "\n"


def _part_relative_path(source_path: str, number: int) -> Path:
    source = Path(source_path)
    if source.parts and source.parts[0] == "docs":
        source = Path(*source.parts[1:])
    stem = source.with_suffix("")
    return stem.parent / f"{stem.name}-part-{number:03d}.md"


def _back_to_index(path: Path) -> str:
    depth = len(path.parent.parts)
    return "../" * depth + "README.md"


def _render_part(
    *, source_path: str, controls: tuple[Control, ...], part_number: int, part_count: int, path: Path
) -> str:
    counts = Counter(control.row["classification"] for control in controls)
    lines = [
        "<!-- Generated by scripts/build_control_classification_docs.py. Do not edit by hand. -->",
        "",
        f"# `{source_path}` — classification reference, part {part_number} of {part_count}",
        "",
        f"[Back to the classification index]({_back_to_index(path)}).",
        "",
        f"Controls: **{len(controls):,}** · Deterministic: **{counts['deterministic']:,}** · "
        f"Nondeterministic: **{counts['nondeterministic']:,}**",
        "",
    ]
    for control in controls:
        lines.extend([render_control(control).rstrip(), ""])
    return "\n".join(lines).rstrip() + "\n"


def _split_source(
    source_path: str,
    controls: list[Control],
    *,
    target_part_bytes: int,
    max_part_bytes: int,
) -> list[Part]:
    if target_part_bytes <= 0 or max_part_bytes <= target_part_bytes:
        raise ValueError("part byte limits are invalid")
    groups: list[list[Control]] = []
    current: list[Control] = []
    current_bytes = 0
    for control in controls:
        block_bytes = len(render_control(control).encode("utf-8"))
        if current and current_bytes + block_bytes > target_part_bytes:
            groups.append(current)
            current = []
            current_bytes = 0
        current.append(control)
        current_bytes += block_bytes
    if current:
        groups.append(current)

    parts: list[Part] = []
    for index, group in enumerate(groups, start=1):
        path = _part_relative_path(source_path, index)
        group_tuple = tuple(group)
        content = _render_part(
            source_path=source_path,
            controls=group_tuple,
            part_number=index,
            part_count=len(groups),
            path=path,
        )
        size = len(content.encode("utf-8"))
        if size >= max_part_bytes:
            raise ValueError(f"{path} is {size:,} bytes; limit is below {max_part_bytes:,}")
        parts.append(Part(path, source_path, group_tuple, content))
    return parts


def render_files(
    corpus: Corpus,
    *,
    target_part_bytes: int = TARGET_PART_BYTES,
    max_part_bytes: int = MAX_PART_BYTES,
) -> dict[Path, str]:
    by_source: dict[str, list[Control]] = defaultdict(list)
    for control in corpus.controls:
        by_source[control.source_path].append(control)

    parts: list[Part] = []
    for source_path in corpus.source_order:
        source_controls = by_source.pop(source_path, [])
        if not source_controls:
            raise ValueError(f"summary source has no controls: {source_path}")
        parts.extend(
            _split_source(
                source_path,
                source_controls,
                target_part_bytes=target_part_bytes,
                max_part_bytes=max_part_bytes,
            )
        )
    if by_source:
        raise ValueError("final corpus contains sources missing from the summary")

    files = {part.path: part.content for part in parts}
    if len(files) != len(parts):
        raise ValueError("generated documentation paths collide")
    emitted_ids = [control.row["control_id"] for part in parts for control in part.controls]
    expected_ids = [control.row["control_id"] for control in corpus.controls]
    if emitted_ids != expected_ids or len(emitted_ids) != len(set(emitted_ids)):
        raise ValueError("generated documentation would omit, duplicate, or reorder controls")

    summary = corpus.summary
    table_rows = []
    for part in parts:
        counts = Counter(control.row["classification"] for control in part.controls)
        size = len(part.content.encode("utf-8"))
        table_rows.append(
            f"| [{part.path.as_posix()}]({part.path.as_posix()}) | `{part.source_path}` | "
            f"{len(part.controls):,} | {counts['deterministic']:,} | "
            f"{counts['nondeterministic']:,} | {size:,} |"
        )
    index_lines = [
        "<!-- Generated by scripts/build_control_classification_docs.py. Do not edit by hand. -->",
        "",
        "# Rule-by-rule control classification reference",
        "",
        "This is the human-readable view of the validated final classification corpus. Each of the "
        f"**{summary['control_count']:,}** controls appears exactly once with its exact statement, final "
        "classification, route, reason, and decision basis.",
        "",
        f"- **Deterministic:** {summary['deterministic']:,}",
        f"- **Nondeterministic:** {summary['nondeterministic']:,}",
        f"- **Registry version:** `{summary['registry_version']}`",
        "",
        "## How to use this reference",
        "",
        "A **deterministic** rule has a bounded automatic check whose named evidence authority can prove "
        "both Pass and Fail for the complete rule. Its exact clauses, checker families, and evidence "
        "authorities are shown below the decision.",
        "",
        "A **nondeterministic** rule requires context, accountable human judgment, specialist or legal "
        "authority, a protocol that is not yet fully defined, or evidence whose complete scope cannot be "
        "sealed. **AI is advisory only for nondeterministic rules.** It may collect or summarize evidence "
        "and suggest concerns, but it must not issue the authoritative Pass/Fail decision.",
        "",
        "Decision basis values show whether the rule was nondeterministic at primary review, rejected by "
        "the independent skeptical review, confirmed there, or then confirmed/reclassified by the third "
        "semantic-strength audit.",
        "",
        "The [research basis](../architecture/control-classification-research.md) explains how NIST OSCAL, "
        "NIST assessment methods, policy-as-code, SLSA verification, and OpenSSF check design informed "
        "the boundary between repeatable proof and advice.",
        "",
        "## Split reference files",
        "",
        f"Files are split only at control boundaries and kept below **{max_part_bytes:,} bytes**, safely "
        "under the 3 MB review limit.",
        "",
        "| File | Source | Controls | Deterministic | Nondeterministic | Bytes |",
        "| --- | --- | ---: | ---: | ---: | ---: |",
        *table_rows,
        "",
    ]
    files[Path("README.md")] = "\n".join(index_lines)
    return files


def generated_files(
    *,
    final_root: Path = FINAL_ROOT,
    packet_root: Path = PACKET_ROOT,
    summary_path: Path = SUMMARY,
    registry_path: Path = REGISTRY,
    target_part_bytes: int = TARGET_PART_BYTES,
    max_part_bytes: int = MAX_PART_BYTES,
) -> dict[Path, str]:
    corpus = load_corpus(
        final_root=final_root,
        packet_root=packet_root,
        summary_path=summary_path,
        registry_path=registry_path,
    )
    return render_files(
        corpus, target_part_bytes=target_part_bytes, max_part_bytes=max_part_bytes
    )


def write_files(files: dict[Path, str], output_root: Path = OUTPUT_ROOT) -> None:
    output_root.mkdir(parents=True, exist_ok=True)
    expected = set(files)
    for existing in sorted(output_root.rglob("*.md")):
        relative = existing.relative_to(output_root)
        if relative not in expected:
            existing.unlink()
    for relative, content in sorted(files.items(), key=lambda item: item[0].as_posix()):
        destination = output_root / relative
        destination.parent.mkdir(parents=True, exist_ok=True)
        temporary = destination.with_name(destination.name + ".tmp")
        temporary.write_text(content, encoding="utf-8")
        os.replace(temporary, destination)
    for directory in sorted(
        (path for path in output_root.rglob("*") if path.is_dir()),
        key=lambda path: len(path.parts),
        reverse=True,
    ):
        if not any(directory.iterdir()):
            directory.rmdir()


def check_files(files: dict[Path, str], output_root: Path = OUTPUT_ROOT) -> None:
    expected = set(files)
    actual = {path.relative_to(output_root) for path in output_root.rglob("*.md")} if output_root.exists() else set()
    if actual != expected:
        missing = sorted(path.as_posix() for path in expected - actual)
        extra = sorted(path.as_posix() for path in actual - expected)
        raise ValueError(f"generated classification docs are missing or stale; missing={missing}, extra={extra}")
    for relative, content in files.items():
        destination = output_root / relative
        if destination.read_text(encoding="utf-8") != content:
            raise ValueError(
                f"{destination.relative_to(ROOT) if destination.is_relative_to(ROOT) else destination} "
                "is stale; run python3 scripts/build_control_classification_docs.py generate"
            )


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("command", choices=("generate", "check"))
    args = parser.parse_args()
    try:
        files = generated_files()
        if args.command == "generate":
            write_files(files)
            print(
                f"generated {len(files) - 1} classification reference parts for "
                f"{load_corpus().summary['control_count']} controls"
            )
        else:
            check_files(files)
            print(f"verified {len(files) - 1} classification reference parts")
    except (OSError, ValueError, json.JSONDecodeError) as error:
        print(f"control classification documentation failed: {error}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
