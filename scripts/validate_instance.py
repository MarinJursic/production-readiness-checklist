#!/usr/bin/env python3
"""Validate scanner JSON output against a checked-in schema."""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path
from typing import Any

from jsonschema import Draft202012Validator, FormatChecker
from referencing import Registry, Resource


ROOT = Path(__file__).resolve().parents[1]
SCHEMAS = ROOT / "schemas"


def load_json(path: Path) -> Any:
    return json.loads(path.read_text(encoding="utf-8"))


def schema_registry() -> tuple[Registry, dict[str, dict[str, Any]]]:
    registry = Registry()
    schemas: dict[str, dict[str, Any]] = {}
    for path in sorted(SCHEMAS.glob("*.schema.json")):
        schema = load_json(path)
        schema_id = schema.get("$id")
        if not isinstance(schema_id, str):
            raise ValueError(f"{path.relative_to(ROOT)} has no $id")
        registry = registry.with_resource(schema_id, Resource.from_contents(schema))
        schemas[path.name] = schema
    return registry, schemas


def validation_errors(instance: Any, schema_name: str) -> list[str]:
    registry, schemas = schema_registry()
    schema = schemas.get(schema_name)
    if schema is None:
        return [f"unknown schema {schema_name!r}"]
    validator = Draft202012Validator(schema, registry=registry, format_checker=FormatChecker())
    errors = sorted(
        validator.iter_errors(instance),
        key=lambda error: tuple(str(part) for part in error.absolute_path),
    )
    details = []
    for error in errors:
        location = ".".join(str(part) for part in error.absolute_path) or "<root>"
        details.append(f"{location}: {error.message}")
    return details


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("schema", help="schema filename under schemas/")
    parser.add_argument("instance", type=Path, help="JSON instance to validate")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    try:
        instance = load_json(args.instance)
        errors = validation_errors(instance, args.schema)
    except (OSError, ValueError, json.JSONDecodeError) as error:
        print(f"output validation failed: {error}", file=sys.stderr)
        return 1
    if errors:
        print("output validation failed:", file=sys.stderr)
        for error in errors:
            print(f"- {error}", file=sys.stderr)
        return 1
    print(f"{args.instance} conforms to {args.schema}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
