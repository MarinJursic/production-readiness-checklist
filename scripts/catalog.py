#!/usr/bin/env python3
"""Validate and generate the versioned scanner catalog pilot."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import sys
import unicodedata
from pathlib import Path
from typing import Any

import yaml
from jsonschema import Draft202012Validator


ROOT = Path(__file__).resolve().parents[1]
CATALOG = ROOT / "catalog"
SCHEMAS = ROOT / "schemas"
REGISTRY = CATALOG / "control-id-registry.json"
GENERATED_PROFILE = ROOT / "docs" / "scanner" / "core-repository.md"
CONTROL_PATTERNS = (
    re.compile(r"^- \[ \] \*\*(PRC-[0-9]{2}-[0-9]{3})\*\* — (.+)$"),
    re.compile(r"^- \[ \] \*\*(USEQ-[A-F0-9]{8})\*\* — (.+)$"),
)


class CatalogError(ValueError):
    """A deterministic catalog validation failure."""


def load_yaml(path: Path) -> dict[str, Any]:
    value = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(value, dict):
        raise CatalogError(f"{path.relative_to(ROOT)}: expected a mapping")
    return value


def load_json(path: Path) -> dict[str, Any]:
    value = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(value, dict):
        raise CatalogError(f"{path.relative_to(ROOT)}: expected an object")
    return value


def validate_schema(instance: dict[str, Any], schema_name: str, source: Path) -> None:
    schema = load_json(SCHEMAS / schema_name)
    errors = sorted(
        Draft202012Validator(schema).iter_errors(instance),
        key=lambda error: tuple(str(part) for part in error.absolute_path),
    )
    if not errors:
        return
    details = []
    for error in errors:
        location = ".".join(str(part) for part in error.absolute_path) or "<root>"
        details.append(f"{source.relative_to(ROOT)}:{location}: {error.message}")
    raise CatalogError("\n".join(details))


def validate_schema_documents() -> None:
    for path in sorted(SCHEMAS.glob("*.schema.json")):
        schema = load_json(path)
        try:
            Draft202012Validator.check_schema(schema)
        except Exception as error:
            raise CatalogError(f"{path.relative_to(ROOT)}: invalid JSON Schema: {error}") from error


def semantic_digest(statement: str) -> str:
    value = unicodedata.normalize("NFKC", statement)
    value = re.sub(r"[*_`]", "", value)
    value = re.sub(r"\s+", " ", value).strip().rstrip(".").casefold()
    return hashlib.sha256(value.encode("utf-8")).hexdigest()


def source_control_entries() -> tuple[list[dict[str, Any]], str]:
    paths = sorted((ROOT / "docs" / "checklists").glob("*.md"))
    paths.extend(sorted((ROOT / "docs" / "engineering").glob("[0-9][0-9]-*.md")))
    entries: list[dict[str, Any]] = []
    source_hasher = hashlib.sha256()
    seen: set[str] = set()
    for path in paths:
        relative = path.relative_to(ROOT).as_posix()
        content = path.read_text(encoding="utf-8")
        source_hasher.update(relative.encode("utf-8"))
        source_hasher.update(b"\0")
        source_hasher.update(content.encode("utf-8"))
        source_hasher.update(b"\0")
        for line_number, line in enumerate(content.splitlines(), start=1):
            for pattern in CONTROL_PATTERNS:
                match = pattern.match(line)
                if not match:
                    continue
                control_id, statement = match.groups()
                if control_id in seen:
                    raise CatalogError(f"duplicate control ID while generating registry: {control_id}")
                seen.add(control_id)
                entries.append(
                    {
                        "id": control_id,
                        "status": "active",
                        "revision": 1,
                        "statement": statement,
                        "semantic_sha256": semantic_digest(statement),
                        "source": {"path": relative, "line": line_number},
                    }
                )
                break
    entries.sort(key=lambda entry: entry["id"])
    return entries, source_hasher.hexdigest()


def previous_registry_entries() -> dict[str, dict[str, Any]]:
    if not REGISTRY.exists():
        return {}
    document = load_json(REGISTRY)
    return {entry["id"]: entry for entry in document.get("entries", [])}


def registry_document() -> dict[str, Any]:
    source_entries, source_digest = source_control_entries()
    previous = previous_registry_entries()
    entries = []
    active_ids = set()
    for entry in source_entries:
        control_id = entry["id"]
        active_ids.add(control_id)
        old = previous.get(control_id)
        if old:
            entry["revision"] = old["revision"] + (old["statement"] != entry["statement"])
        entries.append(entry)
    for control_id, old in previous.items():
        if control_id in active_ids:
            continue
        retired = dict(old)
        retired["status"] = "retired"
        entries.append(retired)
    entries.sort(key=lambda entry: entry["id"])
    return {
        "schema_version": "prc.control-id-registry/v0.1",
        "registry_version": "0.1.0",
        "source_sha256": source_digest,
        "entry_count": len(entries),
        "entries": entries,
    }


def yaml_documents(directory: str) -> list[tuple[Path, dict[str, Any]]]:
    return [(path, load_yaml(path)) for path in sorted((CATALOG / directory).glob("*.yaml"))]


def index_unique(records: list[dict[str, Any]], kind: str) -> dict[str, dict[str, Any]]:
    indexed: dict[str, dict[str, Any]] = {}
    for record in records:
        record_id = record["id"]
        if record_id in indexed:
            raise CatalogError(f"duplicate {kind} ID: {record_id}")
        indexed[record_id] = record
    return indexed


def validate_catalog() -> tuple[dict[str, dict[str, Any]], dict[str, dict[str, Any]], list[dict[str, Any]]]:
    validate_schema_documents()
    objective_records: list[dict[str, Any]] = []
    assertion_records: list[dict[str, Any]] = []
    profiles: list[dict[str, Any]] = []
    for path, document in yaml_documents("objectives"):
        validate_schema(document, "objective-catalog.schema.json", path)
        objective_records.extend(document["objectives"])
    for path, document in yaml_documents("assertions"):
        validate_schema(document, "assertion-catalog.schema.json", path)
        assertion_records.extend(document["assertions"])
    for path, document in yaml_documents("profiles"):
        validate_schema(document, "profile.schema.json", path)
        profiles.append(document)

    objectives = index_unique(objective_records, "objective")
    assertions = index_unique(assertion_records, "assertion")
    profile_index = index_unique(profiles, "profile")
    if len(profile_index) != len(profiles):
        raise CatalogError("profile indexing failed")

    generated_registry = registry_document()
    registry = {entry["id"]: entry for entry in generated_registry["entries"]}
    for objective in objectives.values():
        control_id = objective["id"]
        entry = registry.get(control_id)
        if entry is None:
            raise CatalogError(f"objective {control_id} is absent from the control registry")
        if entry["statement"] != objective["statement"]:
            raise CatalogError(f"objective {control_id} statement differs from its source control")
        if entry["source"] != objective["source"]:
            raise CatalogError(f"objective {control_id} source location is stale")
        if entry["revision"] != objective["revision"]:
            raise CatalogError(f"objective {control_id} revision differs from the registry")
        for assertion_id in objective["assertion_ids"]:
            assertion = assertions.get(assertion_id)
            if assertion is None:
                raise CatalogError(f"objective {control_id} references unknown assertion {assertion_id}")
            if control_id not in assertion["control_ids"]:
                raise CatalogError(f"assertion {assertion_id} does not map back to {control_id}")

    for assertion in assertions.values():
        for control_id in assertion["control_ids"]:
            if control_id not in registry:
                raise CatalogError(f"assertion {assertion['id']} references unknown control {control_id}")
            objective = objectives.get(control_id)
            if objective and assertion["id"] not in objective["assertion_ids"]:
                raise CatalogError(f"objective {control_id} does not map back to {assertion['id']}")
    for profile in profiles:
        for assertion_id in profile["assertion_ids"]:
            if assertion_id not in assertions:
                raise CatalogError(f"profile {profile['id']} references unknown assertion {assertion_id}")
    return objectives, assertions, profiles


def registry_text() -> str:
    return json.dumps(registry_document(), indent=2, ensure_ascii=False) + "\n"


def profile_markdown(
    objectives: dict[str, dict[str, Any]],
    assertions: dict[str, dict[str, Any]],
    profile: dict[str, Any],
) -> str:
    lines = [
        "<!-- Generated by scripts/catalog.py. Do not edit directly. -->",
        "",
        f"# {profile['title']}",
        "",
        f"Profile: `{profile['id']}@{profile['version']}`",
        "",
        profile["description"],
        "",
        "This pilot demonstrates the structured control-to-assertion model. A listed",
        "assertion is not a claim that the scanner implementation exists unless its",
        "implementation is present and verified by the scanner test suite.",
        "",
        "## Assertions",
        "",
        "| Assertion | Objective | Severity | Gate | Evidence | Remediation |",
        "| --- | --- | --- | --- | --- | --- |",
    ]
    for assertion_id in profile["assertion_ids"]:
        assertion = assertions[assertion_id]
        control_links = []
        for control_id in assertion["control_ids"]:
            objective = objectives.get(control_id)
            if objective:
                source = objective["source"]
                relative = "../" + source["path"].removeprefix("docs/")
                control_links.append(f"[{control_id}]({relative})")
            else:
                control_links.append(f"`{control_id}`")
        evidence = ", ".join(item["kind"] for item in assertion["evidence_required"])
        lines.append(
            f"| `{assertion_id}` — {assertion['title']} | {', '.join(control_links)} | "
            f"{assertion['severity']} | {assertion['gate']} | {evidence} | "
            f"{assertion['remediation_class']} |"
        )
    lines.extend(
        [
            "",
            "## Result interpretation",
            "",
            "The profile is satisfied only when every applicable required assertion has",
            "current passing evidence and no configured gate blocks the result. Manual or",
            "external evidence remains visible and cannot be inferred from repository files.",
        ]
    )
    return "\n".join(lines) + "\n"


def generated_files() -> dict[Path, str]:
    objectives, assertions, profiles = validate_catalog()
    profile = next((item for item in profiles if item["id"] == "prc/core-repository"), None)
    if profile is None:
        raise CatalogError("missing prc/core-repository profile")
    return {
        REGISTRY: registry_text(),
        GENERATED_PROFILE: profile_markdown(objectives, assertions, profile),
    }


def generate() -> None:
    for path, content in generated_files().items():
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(content, encoding="utf-8")
        print(f"generated {path.relative_to(ROOT)}")


def check() -> None:
    drifted = []
    for path, expected in generated_files().items():
        if not path.exists() or path.read_text(encoding="utf-8") != expected:
            drifted.append(path.relative_to(ROOT).as_posix())
    if drifted:
        raise CatalogError(
            "generated catalog files are stale: " + ", ".join(drifted) + "; run scripts/catalog.py generate"
        )
    print("catalog schemas, references, registry, and generated documentation are valid")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("command", choices=("validate", "generate", "check"))
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    try:
        if args.command == "generate":
            generate()
        elif args.command == "check":
            check()
        else:
            objectives, assertions, profiles = validate_catalog()
            print(
                f"validated {len(objectives)} objectives, {len(assertions)} assertions, "
                f"and {len(profiles)} profiles"
            )
    except (CatalogError, json.JSONDecodeError, yaml.YAMLError) as error:
        print(f"catalog validation failed: {error}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
