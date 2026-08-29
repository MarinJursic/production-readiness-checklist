#!/usr/bin/env python3
"""Build the authenticated catalog of exact deterministic clause programs.

The reviewed classification answers whether a clause can be decided without
judgment. The program definitions separately provide typed raw facts,
independently sealed parameters, a closed predicate, and adversarial fixtures.
Predicates are never inferred from prose and providers never supply verdicts.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import sys
from collections import Counter
from pathlib import Path
from typing import Any, Iterable

from jsonschema import Draft202012Validator
from jsonschema.exceptions import SchemaError


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_BINDINGS = ROOT / "catalog" / "control-check-bindings.json"
DEFAULT_BINDING_SCHEMA = ROOT / "schemas" / "control-check-bindings.schema.json"
DEFAULT_PROGRAM_SCHEMA = ROOT / "schemas" / "control-check-program.schema.json"
DEFAULT_CATALOG_SCHEMA = ROOT / "schemas" / "control-check-program-catalog.schema.json"
DEFAULT_DEFINITION_SCHEMA = ROOT / "schemas" / "control-check-program-definitions.schema.json"
DEFAULT_DEFINITION_DIR = ROOT / "research" / "control-classification" / "program-specs"
DEFAULT_OUTPUT = ROOT / "catalog" / "control-check-programs.json"

SCHEMA_VERSION = "prc.control-check-program-catalog/v0.4"
GENERATOR_ID = "prc.build-control-check-programs@0.4"
PROGRAM_SCHEMA_VERSION = "prc.control-check-program/v0.1"
BINDING_SCHEMA_VERSION = "prc.control-check-bindings/v0.1"
EXPECTED_CONTROL_COUNT = 686
EXPECTED_CLAUSE_COUNT = 765

OPAQUE_VERDICT_NAMES = {
    "clause_satisfied", "is_compliant", "compliant", "control_passed",
    "provider_verdict", "verified", "valid", "passed", "result",
}

FORBIDDEN_GENERIC_DEFINITION_TEXT = (
    '"/semantic/',
    "schema binding for every promise",
    "tokenized statement words",
    "the entire clause reduces to closed typed comparisons over independently authoritative raw values and sealed expectations",
    "a tempting partial signal is present, but the exact predicate still rejects the broken required value",
    "is bounded by direct raw facts",
    "closed operations",
    "jointly decide every named part of",
    "can look correct while",
    "is missing or mismatched for the exact promise",
    "the program reads",
    "the passing relation is specific",
    "false-pass scenario for",
    "evidence that may still appear healthy",
    "decisive broken observation",
    "is bounded to the sealed",
    "its pinned schema exposes these concrete members rather than statement tokens",
    "fixture rejects a complete-looking record",
)

RUNTIME_REQUIREMENTS = {
    "subject_id": "inject_at_runtime",
    "subjects": "inject_complete_bounded_inventory_at_runtime",
    "inventory_sha256": "inject_at_runtime",
    "maximum_evidence_age_seconds": "inject_approved_freshness_at_runtime",
    "allow_not_applicable": "inject_reviewed_applicability_at_runtime",
    "applicability_proof_contract_sha256": "inject_at_runtime",
    "sealed_parameters": "compile_from_declared_trusted_origin_before_requesting_evidence",
    "sealed_parameters_bound_by": "program_sha256",
    "evidence_provider_may_supply_parameters": False,
    "provider_registration": "required_before_evidence",
    "provider_claimed": False,
    "domain_evidence_collector": "not_shipped_or_registered",
    "missing_capability_result": "blocked",
}


class ProgramBuildError(ValueError):
    pass


def canonical_bytes(value: Any) -> bytes:
    return json.dumps(value, ensure_ascii=False, sort_keys=True, separators=(",", ":")).encode("utf-8")


def digest_value(value: Any) -> str:
    return hashlib.sha256(canonical_bytes(value)).hexdigest()


def statement_sha256(statement: str) -> str:
    return hashlib.sha256(statement.encode("utf-8")).hexdigest()


def _reject_duplicate_keys(pairs: list[tuple[str, Any]]) -> dict[str, Any]:
    result: dict[str, Any] = {}
    for key, value in pairs:
        if key in result:
            raise ProgramBuildError(f"duplicate JSON object key {key!r}")
        result[key] = value
    return result


def read_object(path: Path, label: str) -> tuple[dict[str, Any], bytes]:
    data = path.read_bytes()
    try:
        value = json.loads(data, object_pairs_hook=_reject_duplicate_keys)
    except UnicodeDecodeError as error:
        raise ProgramBuildError(f"{label} {path} is not valid UTF-8") from error
    if not isinstance(value, dict):
        raise ProgramBuildError(f"{label} {path} must be a JSON object")
    return value, data


def validate_schema(document: dict[str, Any], schema: dict[str, Any], label: str) -> None:
    errors = sorted(
        Draft202012Validator(schema).iter_errors(document),
        key=lambda item: tuple(map(str, item.absolute_path)),
    )
    if errors:
        first = errors[0]
        location = ".".join(map(str, first.absolute_path)) or "<root>"
        raise ProgramBuildError(f"{label} schema failed at {location}: {first.message}")


def default_definition_paths() -> tuple[Path, ...]:
    paths = tuple(sorted(DEFAULT_DEFINITION_DIR.glob("*.json")))
    if not paths:
        raise ProgramBuildError("no deterministic program definition files exist")
    return paths


def load_bindings(bindings_path: Path, schema_path: Path) -> tuple[dict[str, Any], bytes]:
    document, data = read_object(bindings_path, "control check bindings")
    schema, _ = read_object(schema_path, "control check binding schema")
    Draft202012Validator.check_schema(schema)
    validate_schema(document, schema, "control check bindings")
    if document["schema_version"] != BINDING_SCHEMA_VERSION:
        raise ProgramBuildError("control check binding schema version is unsupported")
    bindings = document["bindings"]
    if document["binding_count"] != EXPECTED_CONTROL_COUNT or len(bindings) != EXPECTED_CONTROL_COUNT:
        raise ProgramBuildError(f"deterministic binding control count is not exactly {EXPECTED_CONTROL_COUNT}")
    control_ids = [binding["control_id"] for binding in bindings]
    if control_ids != sorted(control_ids) or len(control_ids) != len(set(control_ids)):
        raise ProgramBuildError("deterministic bindings must be ordered and unique")
    implementations = {item["checker_family"]: item for item in document["implementation_registry"]}
    if len(implementations) != len(document["implementation_registry"]):
        raise ProgramBuildError("implementation registry contains duplicate checker families")
    identities: set[tuple[str, int, str]] = set()
    count = 0
    for binding in bindings:
        ordinals = [clause["ordinal"] for clause in binding["clauses"]]
        if ordinals != list(range(1, len(ordinals) + 1)):
            raise ProgramBuildError(f"{binding['control_id']}: clause ordinals are not contiguous")
        for clause in binding["clauses"]:
            count += 1
            identity = (binding["control_id"], clause["ordinal"], clause["clause_id"])
            if identity in identities:
                raise ProgramBuildError(f"duplicate deterministic clause identity {identity}")
            identities.add(identity)
            reviewed = {
                "statement": clause["statement"],
                "checker_family": clause["checker_family"],
                "evidence_authority": clause["evidence_authority"],
            }
            if clause["clause_id"] != digest_value(reviewed):
                raise ProgramBuildError(f"{binding['control_id']} clause identity is stale")
            implementation = implementations.get(clause["checker_family"])
            if implementation is None or (
                clause["implementation_id"] != implementation["implementation_id"]
                or clause["implementation_contract_sha256"] != implementation["implementation_contract_sha256"]
                or clause["evidence_authority"] not in implementation["supported_evidence_authorities"]
            ):
                raise ProgramBuildError(f"{binding['control_id']} implementation binding is stale")
    if count != EXPECTED_CLAUSE_COUNT:
        raise ProgramBuildError(f"deterministic clause count is not exactly {EXPECTED_CLAUSE_COUNT}")
    return document, data


def expression_operations(expression: dict[str, Any]) -> list[str]:
    operation = expression.get("op")
    if not isinstance(operation, str):
        raise ProgramBuildError("predicate expression is missing a string operation")
    result = [operation]
    child = expression.get("arg")
    if child is not None:
        if not isinstance(child, dict):
            raise ProgramBuildError("predicate arg must be an object")
        result.extend(expression_operations(child))
    children = expression.get("args", [])
    if not isinstance(children, list):
        raise ProgramBuildError("predicate args must be an array")
    for item in children:
        if not isinstance(item, dict):
            raise ProgramBuildError("predicate child must be an object")
        result.extend(expression_operations(item))
    return result


def expression_references(expression: dict[str, Any]) -> tuple[set[str], set[str]]:
    facts: set[str] = set()
    parameters: set[str] = set()
    for key in ("fact", "other_fact", "third_fact"):
        value = expression.get(key)
        if value is not None:
            if not isinstance(value, str):
                raise ProgramBuildError(f"predicate {key} reference must be a string")
            facts.add(value)
    value = expression.get("parameter")
    if value is not None:
        if not isinstance(value, str):
            raise ProgramBuildError("predicate parameter reference must be a string")
        parameters.add(value)
    child = expression.get("arg")
    if isinstance(child, dict):
        child_facts, child_parameters = expression_references(child)
        facts.update(child_facts)
        parameters.update(child_parameters)
    for child in expression.get("args", []):
        child_facts, child_parameters = expression_references(child)
        facts.update(child_facts)
        parameters.update(child_parameters)
    return facts, parameters


def validate_fixture_contracts(definition: dict[str, Any]) -> None:
    fact_types = {item["fact_id"]: item["fact_type"] for item in definition["raw_fact_contracts"]}
    parameter_types = {item["parameter_id"]: item["parameter_type"] for item in definition["sealed_parameter_contracts"]}
    if len(fact_types) != len(definition["raw_fact_contracts"]):
        raise ProgramBuildError(f"{definition['control_id']} has duplicate raw fact contracts")
    if len(parameter_types) != len(definition["sealed_parameter_contracts"]):
        raise ProgramBuildError(f"{definition['control_id']} has duplicate sealed parameter contracts")
    for name, fixture in definition["fixtures"].items():
        fixture_facts = fixture["facts"]
        fixture_parameters = fixture["parameters"]
        if set(fixture_parameters) != set(parameter_types):
            raise ProgramBuildError(f"{definition['control_id']} {name} fixture does not bind every sealed parameter exactly once")
        if name != "blocked" and set(fixture_facts) != set(fact_types):
            raise ProgramBuildError(f"{definition['control_id']} {name} fixture does not provide every raw fact exactly once")
        if name == "blocked" and not set(fixture_facts).issubset(fact_types):
            raise ProgramBuildError(f"{definition['control_id']} blocked fixture contains an undeclared raw fact")
        for key, value in fixture_parameters.items():
            if not isinstance(value, dict) or value.get("type") != parameter_types[key]:
                raise ProgramBuildError(f"{definition['control_id']} {name} fixture has wrong parameter type for {key}")
        for key, value in fixture_facts.items():
            if not isinstance(value, dict) or value.get("type") != fact_types[key]:
                raise ProgramBuildError(f"{definition['control_id']} {name} fixture has wrong fact type for {key}")
        if name == "blocked" and set(fixture_facts) == set(fact_types) and all(
            value.get("complete") is True for value in fixture_facts.values()
        ):
            raise ProgramBuildError(f"{definition['control_id']} blocked fixture does not omit or mark any fact incomplete")


def validate_definition(definition: dict[str, Any], binding: dict[str, Any], clause: dict[str, Any]) -> None:
    identity = (definition["control_id"], definition["clause_ordinal"], definition["clause_id"])
    if definition["classification_status"] != "exact_predicate" or definition["classification_error_reason"] is not None:
        raise ProgramBuildError(f"deterministic definition is not an exact predicate: {identity}")
    if (
        definition["control_revision"] != binding["revision"]
        or definition["control_semantic_sha256"] != binding["semantic_sha256"]
        or definition["clause_statement"] != clause["statement"]
        or definition["corrected_checker_family"] != clause["checker_family"]
        or definition["corrected_required_authority"] != clause["evidence_authority"]
    ):
        raise ProgramBuildError(f"stale program definition for {identity}")
    serialized_definition = json.dumps(definition, ensure_ascii=False, sort_keys=True).lower()
    for fragment in FORBIDDEN_GENERIC_DEFINITION_TEXT:
        if fragment in serialized_definition:
            raise ProgramBuildError(f"forbidden generic predicate delegation in {identity}: {fragment}")
    for contract in definition["raw_fact_contracts"]:
        leaf = contract["fact_id"].rsplit(".", 1)[-1].lower()
        if leaf in OPAQUE_VERDICT_NAMES or leaf in {"record-paths", "required-paths", "record_paths", "required_paths"}:
            raise ProgramBuildError(f"opaque delegated verdict fact in {identity}: {contract['fact_id']}")
        if contract["authority"] != clause["evidence_authority"]:
            raise ProgramBuildError(f"raw fact authority differs from corrected binding for {identity}")
    predicate = definition["predicate"]
    operations = expression_operations(predicate)
    if sorted(set(operations)) != definition["required_runtime_ops"]:
        raise ProgramBuildError(f"required runtime operations are stale for {identity}")
    facts, parameters = expression_references(predicate)
    declared_facts = {item["fact_id"] for item in definition["raw_fact_contracts"]}
    declared_parameters = {item["parameter_id"] for item in definition["sealed_parameter_contracts"]}
    if any(key.rsplit(".", 1)[-1].lower() in {"record-paths", "required-paths", "record_paths", "required_paths"} for key in declared_parameters):
        raise ProgramBuildError(f"forbidden generic schema-path parameter in {identity}")
    if facts != declared_facts:
        raise ProgramBuildError(f"predicate raw fact references do not exactly match contracts for {identity}")
    if parameters != declared_parameters:
        raise ProgramBuildError(f"predicate parameter references do not exactly match contracts for {identity}")
    validate_fixture_contracts(definition)


def load_definitions(
    paths: Iterable[Path], schema_path: Path, bindings: dict[str, Any], binding_sha256: str,
) -> tuple[dict[tuple[str, int, str], dict[str, Any]], str, str]:
    schema, schema_bytes = read_object(schema_path, "control check program definition schema")
    Draft202012Validator.check_schema(schema)
    expected = {
        (binding["control_id"], clause["ordinal"], clause["clause_id"]): (binding, clause)
        for binding in bindings["bindings"] for clause in binding["clauses"]
    }
    merged: dict[tuple[str, int, str], dict[str, Any]] = {}
    corpus = hashlib.sha256()
    scopes: set[str] = set()
    path_list = tuple(sorted(paths))
    if not path_list:
        raise ProgramBuildError("no deterministic program definitions were supplied")
    for path in path_list:
        document, data = read_object(path, "control check program definitions")
        validate_schema(document, schema, f"control check program definitions {path.name}")
        if document["source_binding_catalog_sha256"] != binding_sha256:
            raise ProgramBuildError(f"{path}: source binding catalog digest is stale")
        if document["definition_count"] != len(document["definitions"]):
            raise ProgramBuildError(f"{path}: definition count is stale")
        scopes.add(document["scope"])
        corpus.update(path.name.encode("utf-8")); corpus.update(b"\0"); corpus.update(data); corpus.update(b"\0")
        previous: tuple[str, int, str] | None = None
        for definition in document["definitions"]:
            identity = (definition["control_id"], definition["clause_ordinal"], definition["clause_id"])
            if previous is not None and identity <= previous:
                raise ProgramBuildError(f"{path}: definitions are not in strict binding order")
            previous = identity
            if identity in merged:
                raise ProgramBuildError(f"duplicate program definition {identity}")
            source = expected.get(identity)
            if source is None:
                raise ProgramBuildError(f"program definition has no current deterministic binding {identity}")
            binding, clause = source
            expected_scope = "structured_record" if clause["checker_family"] == "structured_record" else "non_structured"
            if document["scope"] != expected_scope:
                raise ProgramBuildError(f"{path}: {identity} belongs to {expected_scope}")
            validate_definition(definition, binding, clause)
            merged[identity] = definition
    if scopes != {"structured_record", "non_structured"}:
        raise ProgramBuildError("structured_record and non_structured definition scopes are both required")
    mismatch = sorted(set(expected) ^ set(merged))
    if mismatch or len(merged) != EXPECTED_CLAUSE_COUNT:
        first = mismatch[0] if mismatch else None
        raise ProgramBuildError(f"program definitions do not cover exactly {EXPECTED_CLAUSE_COUNT} current clauses; first mismatch is {first}")
    for field in ("review_reason", "counterexample_analysis"):
        counts = Counter(definition[field].strip() for definition in merged.values())
        duplicate = next((text for text, count in counts.items() if count > 1), None)
        if duplicate is not None:
            raise ProgramBuildError(f"program definitions reuse a generic {field}: {duplicate[:120]}")
    return merged, hashlib.sha256(schema_bytes).hexdigest(), corpus.hexdigest()


def template_for(
    binding: dict[str, Any], clause: dict[str, Any], definition: dict[str, Any], program_schema_sha256: str,
) -> dict[str, Any]:
    control_id = binding["control_id"]
    template = {
        "template_id": digest_value({"control_id": control_id, "control_revision": binding["revision"], "clause_id": clause["clause_id"]}),
        "program_schema_version": PROGRAM_SCHEMA_VERSION,
        "program_schema_sha256": program_schema_sha256,
        "control_id": control_id,
        "control_revision": binding["revision"],
        "control_semantic_sha256": binding["semantic_sha256"],
        "clause_ordinal": clause["ordinal"],
        "clause_id": clause["clause_id"],
        "clause_statement": clause["statement"],
        "clause_statement_sha256": statement_sha256(clause["statement"]),
        "checker_family": clause["checker_family"],
        "required_authority": clause["evidence_authority"],
        "implementation_id": clause["implementation_id"],
        "implementation_contract_sha256": clause["implementation_contract_sha256"],
        "review_status": "predicate_defined_provider_unregistered",
        "predicate_defined": True,
        "end_to_end_runnable": False,
        "predicate_shape": definition["predicate"]["op"],
        "review_reason": definition["review_reason"],
        "counterexample_analysis": definition["counterexample_analysis"],
        "raw_fact_contracts": definition["raw_fact_contracts"],
        "sealed_parameter_contracts": definition["sealed_parameter_contracts"],
        "predicate": definition["predicate"],
        "required_runtime_ops": definition["required_runtime_ops"],
        "collector_contract": definition["collector_contract"],
        "provider_capability_status": "unregistered",
        "runtime_requirements": dict(RUNTIME_REQUIREMENTS),
    }
    template["template_sha256"] = digest_value(template)
    return template


def build_document(
    *, bindings_path: Path = DEFAULT_BINDINGS,
    binding_schema_path: Path = DEFAULT_BINDING_SCHEMA,
    program_schema_path: Path = DEFAULT_PROGRAM_SCHEMA,
    definition_schema_path: Path = DEFAULT_DEFINITION_SCHEMA,
    definition_paths: Iterable[Path] | None = None,
) -> dict[str, Any]:
    bindings, binding_bytes = load_bindings(bindings_path, binding_schema_path)
    binding_sha256 = hashlib.sha256(binding_bytes).hexdigest()
    definitions, definition_schema_sha256, definition_corpus_sha256 = load_definitions(
        default_definition_paths() if definition_paths is None else definition_paths,
        definition_schema_path, bindings, binding_sha256,
    )
    program_schema, program_schema_bytes = read_object(program_schema_path, "control check program schema")
    Draft202012Validator.check_schema(program_schema)
    if program_schema.get("properties", {}).get("schema_version", {}).get("const") != PROGRAM_SCHEMA_VERSION:
        raise ProgramBuildError("control check program schema version is unsupported")
    program_schema_sha256 = hashlib.sha256(program_schema_bytes).hexdigest()
    templates = []
    for binding in bindings["bindings"]:
        for clause in binding["clauses"]:
            identity = (binding["control_id"], clause["ordinal"], clause["clause_id"])
            templates.append(template_for(binding, clause, definitions[identity], program_schema_sha256))
    identities = [(item["control_id"], item["clause_ordinal"], item["clause_id"]) for item in templates]
    if identities != sorted(identities) or len({item["template_id"] for item in templates}) != len(templates):
        raise ProgramBuildError("program templates must be ordered with unique clause identities")
    controls = {item["control_id"] for item in templates}
    shapes = Counter(item["predicate_shape"] for item in templates)
    unsigned = {
        "schema_version": SCHEMA_VERSION,
        "generator_id": GENERATOR_ID,
        "program_schema_version": PROGRAM_SCHEMA_VERSION,
        "program_schema_sha256": program_schema_sha256,
        "binding_schema_version": bindings["schema_version"],
        "binding_catalog_sha256": binding_sha256,
        "definition_schema_sha256": definition_schema_sha256,
        "definition_corpus_sha256": definition_corpus_sha256,
        "registry_version": bindings["registry_version"],
        "registry_sha256": bindings["registry_sha256"],
        "methodology_sha256": bindings["methodology_sha256"],
        "classification_corpus_sha256": bindings["classification_corpus_sha256"],
        "control_count": len(controls),
        "template_count": len(templates),
        "predicate_defined_count": len(templates),
        "implementation_missing_count": 0,
        "provider_capability_missing_count": len(templates),
        "end_to_end_runnable_template_count": 0,
        "end_to_end_runnable_control_count": 0,
        "blocked_control_count": len(controls),
        "classification_error_count": 0,
        "predicate_shape_counts": dict(sorted(shapes.items())),
        "templates": templates,
    }
    return {**unsigned, "catalog_sha256": digest_value(unsigned)}


def validate_catalog(document: dict[str, Any], schema_path: Path = DEFAULT_CATALOG_SCHEMA) -> None:
    schema, _ = read_object(schema_path, "program catalog schema")
    Draft202012Validator.check_schema(schema)
    validate_schema(document, schema, "program catalog")
    templates = document["templates"]
    identities = [(item["control_id"], item["clause_ordinal"], item["clause_id"]) for item in templates]
    if identities != sorted(identities) or len({item["template_id"] for item in templates}) != len(templates):
        raise ProgramBuildError("program catalog templates are not ordered and unique")
    if (
        document["control_count"] != EXPECTED_CONTROL_COUNT
        or document["template_count"] != EXPECTED_CLAUSE_COUNT
        or document["predicate_defined_count"] != EXPECTED_CLAUSE_COUNT
        or document["implementation_missing_count"] != 0
        or document["provider_capability_missing_count"] != EXPECTED_CLAUSE_COUNT
        or document["end_to_end_runnable_template_count"] != 0
        or document["end_to_end_runnable_control_count"] != 0
        or document["blocked_control_count"] != EXPECTED_CONTROL_COUNT
        or document["classification_error_count"] != 0
    ):
        raise ProgramBuildError("program capability counts are stale")
    if dict(sorted(Counter(item["predicate_shape"] for item in templates).items())) != document["predicate_shape_counts"]:
        raise ProgramBuildError("predicate shape counts are stale")
    for template in templates:
        expected_id = digest_value({"control_id": template["control_id"], "control_revision": template["control_revision"], "clause_id": template["clause_id"]})
        if template["template_id"] != expected_id:
            raise ProgramBuildError(f"template identity mismatch for {template['control_id']}")
        if template["program_schema_sha256"] != document["program_schema_sha256"]:
            raise ProgramBuildError(f"program schema binding mismatch for {template['template_id']}")
        if template["clause_statement_sha256"] != statement_sha256(template["clause_statement"]):
            raise ProgramBuildError(f"clause statement digest mismatch for {template['template_id']}")
        if not template["predicate_defined"] or template["review_status"] != "predicate_defined_provider_unregistered" or template["end_to_end_runnable"]:
            raise ProgramBuildError(f"predicate/provider status is inconsistent for {template['template_id']}")
        if template["provider_capability_status"] != "unregistered" or template["collector_contract"]["provider_status"] != "unregistered":
            raise ProgramBuildError(f"unregistered collector was incorrectly claimed for {template['template_id']}")
        facts, parameters = expression_references(template["predicate"])
        if facts != {item["fact_id"] for item in template["raw_fact_contracts"]}:
            raise ProgramBuildError(f"catalog predicate fact contracts are stale for {template['template_id']}")
        if parameters != {item["parameter_id"] for item in template["sealed_parameter_contracts"]}:
            raise ProgramBuildError(f"catalog predicate parameter contracts are stale for {template['template_id']}")
        if sorted(set(expression_operations(template["predicate"]))) != template["required_runtime_ops"]:
            raise ProgramBuildError(f"catalog predicate operation list is stale for {template['template_id']}")
        unsigned = {key: value for key, value in template.items() if key != "template_sha256"}
        if digest_value(unsigned) != template["template_sha256"]:
            raise ProgramBuildError(f"template digest mismatch for {template['template_id']}")
    unsigned = {key: value for key, value in document.items() if key != "catalog_sha256"}
    if digest_value(unsigned) != document["catalog_sha256"]:
        raise ProgramBuildError("catalog digest mismatch")


def serialize(document: dict[str, Any]) -> str:
    return json.dumps(document, ensure_ascii=False, indent=2) + "\n"


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("command", choices=("generate", "check"))
    parser.add_argument("--bindings", type=Path, default=DEFAULT_BINDINGS)
    parser.add_argument("--binding-schema", type=Path, default=DEFAULT_BINDING_SCHEMA)
    parser.add_argument("--program-schema", type=Path, default=DEFAULT_PROGRAM_SCHEMA)
    parser.add_argument("--definition-schema", type=Path, default=DEFAULT_DEFINITION_SCHEMA)
    parser.add_argument("--catalog-schema", type=Path, default=DEFAULT_CATALOG_SCHEMA)
    parser.add_argument("--definitions", type=Path, nargs="*")
    parser.add_argument("--output", type=Path, default=DEFAULT_OUTPUT)
    args = parser.parse_args()
    try:
        document = build_document(
            bindings_path=args.bindings, binding_schema_path=args.binding_schema,
            program_schema_path=args.program_schema, definition_schema_path=args.definition_schema,
            definition_paths=args.definitions,
        )
        validate_catalog(document, args.catalog_schema)
        expected = serialize(document)
        if args.command == "generate":
            args.output.write_text(expected, encoding="utf-8")
            print(f"generated {document['template_count']} exact clause programs; collectors remain fail-closed until registered")
        else:
            actual, _ = read_object(args.output, "generated program catalog")
            validate_catalog(actual, args.catalog_schema)
            if canonical_bytes(actual) != canonical_bytes(document) or args.output.read_text(encoding="utf-8") != expected:
                raise ProgramBuildError("program catalog is stale or nondeterministic")
            print(f"verified {document['template_count']} exact clause programs")
        return 0
    except (OSError, ProgramBuildError, SchemaError, json.JSONDecodeError) as error:
        print(f"control check program build failed: {error}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
