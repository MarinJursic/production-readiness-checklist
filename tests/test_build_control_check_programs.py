from __future__ import annotations

import copy
import hashlib
import importlib.util
import json
import tempfile
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SCRIPT = ROOT / "scripts" / "build_control_check_programs.py"
SPEC = importlib.util.spec_from_file_location("build_control_check_programs", SCRIPT)
assert SPEC is not None and SPEC.loader is not None
PROGRAMS = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(PROGRAMS)


class ControlCheckProgramTests(unittest.TestCase):
    def test_generated_catalog_is_current_complete_and_deterministic(self) -> None:
        document = PROGRAMS.build_document()
        PROGRAMS.validate_catalog(document)
        actual = PROGRAMS.DEFAULT_OUTPUT.read_text(encoding="utf-8")
        self.assertEqual(actual, PROGRAMS.serialize(document))
        self.assertEqual(json.loads(actual), document)
        self.assertEqual(document["control_count"], 686)
        self.assertEqual(document["template_count"], 765)
        self.assertEqual(document["predicate_defined_count"], 765)
        self.assertEqual(document["implementation_missing_count"], 0)
        capability_count = len(PROGRAMS.load_capabilities(
            PROGRAMS.DEFAULT_CAPABILITIES, PROGRAMS.DEFAULT_CAPABILITY_SCHEMA,
        ))
        self.assertEqual(document["provider_capability_missing_count"], 765 - capability_count)
        self.assertEqual(document["end_to_end_runnable_template_count"], capability_count)
        self.assertEqual(document["end_to_end_runnable_control_count"], capability_count)
        self.assertEqual(document["blocked_control_count"], 686 - capability_count)
        self.assertEqual(document["classification_error_count"], 0)
        self.assertEqual(
            document["binding_catalog_sha256"],
            hashlib.sha256(PROGRAMS.DEFAULT_BINDINGS.read_bytes()).hexdigest(),
        )
        self.assertEqual(
            document["program_schema_sha256"],
            hashlib.sha256(PROGRAMS.DEFAULT_PROGRAM_SCHEMA.read_bytes()).hexdigest(),
        )

    def test_every_template_maps_to_one_binding_and_one_exact_definition(self) -> None:
        catalog = PROGRAMS.build_document()
        bindings = json.loads(PROGRAMS.DEFAULT_BINDINGS.read_text(encoding="utf-8"))
        expected = {
            (binding["control_id"], clause["ordinal"], clause["clause_id"]): (binding, clause)
            for binding in bindings["bindings"] for clause in binding["clauses"]
        }
        self.assertEqual(len(expected), 765)
        identities = []
        for template in catalog["templates"]:
            identity = (template["control_id"], template["clause_ordinal"], template["clause_id"])
            identities.append(identity)
            binding, clause = expected.pop(identity)
            self.assertEqual(template["control_revision"], binding["revision"])
            self.assertEqual(template["control_semantic_sha256"], binding["semantic_sha256"])
            self.assertEqual(template["clause_statement"], clause["statement"])
            self.assertEqual(template["checker_family"], clause["checker_family"])
            self.assertEqual(template["required_authority"], clause["evidence_authority"])
            self.assertTrue(template["predicate_defined"])
            self.assertIsInstance(template["predicate"], dict)
            self.assertTrue(template["raw_fact_contracts"])
            self.assertTrue(template["required_runtime_ops"])
        self.assertEqual(identities, sorted(identities))
        self.assertFalse(expected)

    def test_only_deterministic_controls_have_programs_and_ai_owns_the_rest(self) -> None:
        catalog = PROGRAMS.build_document()
        contracts = json.loads((ROOT / "catalog" / "control-contracts.json").read_text(encoding="utf-8"))
        deterministic = {item["control_id"] for item in contracts["contracts"] if item["classification"] == "deterministic"}
        nondeterministic = {item["control_id"] for item in contracts["contracts"] if item["classification"] == "nondeterministic"}
        program_controls = {item["control_id"] for item in catalog["templates"]}
        self.assertEqual(program_controls, deterministic)
        self.assertEqual(len(deterministic), 686)
        self.assertEqual(len(nondeterministic), 9356)
        self.assertTrue(program_controls.isdisjoint(nondeterministic))

    def test_provider_manifest_is_the_exact_source_of_runnable_claims(self) -> None:
        capabilities = PROGRAMS.load_capabilities(
            PROGRAMS.DEFAULT_CAPABILITIES, PROGRAMS.DEFAULT_CAPABILITY_SCHEMA,
        )
        runnable = set()
        for template in PROGRAMS.build_document()["templates"]:
            self.assertTrue(template["predicate_defined"])
            collector_id = template["collector_contract"]["collector_id"]
            registered = collector_id in capabilities
            if registered:
                runnable.add(collector_id)
            status = "registered" if registered else "unregistered"
            self.assertEqual(template["review_status"], f"predicate_defined_provider_{status}")
            self.assertEqual(template["end_to_end_runnable"], registered)
            self.assertEqual(template["provider_capability_status"], status)
            self.assertEqual(template["collector_contract"]["provider_status"], status)
            runtime = template["runtime_requirements"]
            self.assertEqual(runtime["provider_claimed"], registered)
            self.assertEqual(
                runtime["domain_evidence_collector"],
                "shipped_and_registered" if registered else "not_shipped_or_registered",
            )
            self.assertFalse(runtime["evidence_provider_may_supply_parameters"])
            self.assertEqual(runtime["missing_capability_result"], "blocked")
        self.assertEqual(runnable, set(capabilities))

    def test_rejects_capability_that_does_not_match_the_reviewed_template(self) -> None:
        source = json.loads(PROGRAMS.DEFAULT_CAPABILITIES.read_text(encoding="utf-8"))
        source["capabilities"][0]["control_id"] = "PRC-36-001"
        with tempfile.TemporaryDirectory() as temporary:
            path = Path(temporary) / "providers.json"
            path.write_text(json.dumps(source), encoding="utf-8")
            with self.assertRaisesRegex(PROGRAMS.ProgramBuildError, "does not match reviewed clause"):
                PROGRAMS.build_document(capabilities_path=path)

    def test_definition_parts_are_reviewable_and_no_generic_oracle_remains(self) -> None:
        paths = PROGRAMS.default_definition_paths()
        self.assertGreaterEqual(len(paths), 2)
        for path in paths:
            self.assertLessEqual(path.stat().st_size, 3_000_000)
        encoded = json.dumps(PROGRAMS.build_document()).lower()
        for forbidden in (
            "clause_satisfied", "provider_verdict", "record-paths", "required-paths",
            "record_paths", "required_paths", '"/semantic/', "schema binding for every promise",
            "observed_typed_tuples", "required_typed_tuples",
            "the entire clause reduces to closed typed comparisons",
            "a tempting partial signal is present, but the exact predicate still rejects",
        ):
            self.assertNotIn(forbidden, encoded)

    def replacement_paths(self, changed_path: Path, original: Path) -> tuple[Path, ...]:
        return tuple(changed_path if path == original else path for path in PROGRAMS.default_definition_paths())

    def test_rejects_missing_duplicate_and_generic_program_definitions(self) -> None:
        original = PROGRAMS.default_definition_paths()[0]
        source = json.loads(original.read_text(encoding="utf-8"))
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)

            missing = copy.deepcopy(source)
            missing["definitions"].pop()
            missing["definition_count"] -= 1
            missing_path = root / original.name
            missing_path.write_text(json.dumps(missing), encoding="utf-8")
            with self.assertRaisesRegex(PROGRAMS.ProgramBuildError, "do not cover exactly"):
                PROGRAMS.build_document(definition_paths=self.replacement_paths(missing_path, original))

            duplicate = copy.deepcopy(source)
            duplicate["definitions"][-1] = copy.deepcopy(duplicate["definitions"][0])
            duplicate_path = root / ("duplicate-" + original.name)
            duplicate_path.write_text(json.dumps(duplicate), encoding="utf-8")
            with self.assertRaisesRegex(PROGRAMS.ProgramBuildError, "strict binding order|duplicate"):
                PROGRAMS.build_document(definition_paths=self.replacement_paths(duplicate_path, original))

            generic = copy.deepcopy(source)
            generic["definitions"][0]["review_reason"] += " schema binding for every promise"
            generic_path = root / ("generic-" + original.name)
            generic_path.write_text(json.dumps(generic), encoding="utf-8")
            with self.assertRaisesRegex(PROGRAMS.ProgramBuildError, "generic predicate delegation"):
                PROGRAMS.build_document(definition_paths=self.replacement_paths(generic_path, original))

    def test_catalog_and_template_digests_reject_corruption(self) -> None:
        document = PROGRAMS.build_document()
        changed = copy.deepcopy(document)
        changed["templates"][0]["clause_statement"] += " changed"
        with self.assertRaisesRegex(PROGRAMS.ProgramBuildError, "statement digest mismatch"):
            PROGRAMS.validate_catalog(changed)

        changed = copy.deepcopy(document)
        changed["catalog_sha256"] = "0" * 64
        with self.assertRaisesRegex(PROGRAMS.ProgramBuildError, "catalog digest mismatch"):
            PROGRAMS.validate_catalog(changed)


if __name__ == "__main__":
    unittest.main()
