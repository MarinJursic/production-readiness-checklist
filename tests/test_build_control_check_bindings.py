from __future__ import annotations

import importlib.util
import json
import tempfile
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SCRIPT = ROOT / "scripts" / "build_control_check_bindings.py"
SPEC = importlib.util.spec_from_file_location("build_control_check_bindings", SCRIPT)
assert SPEC and SPEC.loader
bindings = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(bindings)


class BindingFixture:
    def __init__(self, root: Path) -> None:
        self.root = root
        self.registry = root / "control-id-registry.json"
        self.final = root / "final"
        self.final.mkdir()
        self.schema = ROOT / "schemas" / "control-check-bindings.schema.json"
        self.rows = [self._row("USEQ-AAAAAAAA", True), self._row("USEQ-BBBBBBBB", False)]
        self._write_registry()
        self.write_final()

    @staticmethod
    def _row(control_id: str, deterministic: bool) -> dict:
        common = {
            "control_id": control_id,
            "revision": 1,
            "semantic_sha256": ("a" if control_id.endswith("AAAAAAAA") else "b") * 64,
            "reason": "A sufficiently detailed reviewed reason for this classification.",
        }
        if deterministic:
            return {
                **common,
                "classification": "deterministic",
                "route": "external_readonly_query",
                "deterministic_clauses": [{
                    "statement": "The exact registry identity matches the requested identity.",
                    "checker_family": "environment_evidence",
                    "evidence_authority": "external_registry",
                }],
                "nondeterministic_reason_codes": [],
                "decision_basis": "two_stage_confirmed",
                "skeptical_verdict": "confirmed_deterministic",
                "counterexample_analysis": "Wrong identity, missing identity, or stale registry evidence is blocked.",
            }
        return {
            **common,
            "classification": "nondeterministic",
            "route": "contextual_judgment",
            "deterministic_clauses": [],
            "nondeterministic_reason_codes": ["contextual_judgment"],
            "decision_basis": "primary_nondeterministic",
            "skeptical_verdict": None,
            "counterexample_analysis": None,
        }

    def _write_registry(self) -> None:
        document = {
            "schema_version": "prc.control-id-registry/v0.1",
            "registry_version": "0.1.0",
            "source_sha256": "c" * 64,
            "entry_count": 2,
            "entries": [
                {"id": "USEQ-AAAAAAAA", "status": "active", "revision": 1, "statement": "A", "semantic_sha256": "a" * 64, "source": {"path": "docs/a.md", "line": 1}},
                {"id": "USEQ-BBBBBBBB", "status": "active", "revision": 1, "statement": "B", "semantic_sha256": "b" * 64, "source": {"path": "docs/b.md", "line": 1}},
            ],
        }
        self.registry.write_text(json.dumps(document, indent=2) + "\n", encoding="utf-8")

    def write_final(self) -> None:
        registry_sha = bindings.hashlib.sha256(self.registry.read_bytes()).hexdigest()
        document = {
            "schema_version": "prc.control-classification-final/v0.1",
            "packet_id": "fixture-part-001",
            "methodology_sha256": "d" * 64,
            "registry_version": "0.1.0",
            "registry_sha256": registry_sha,
            "source_path": "docs/fixture.md",
            "control_count": len(self.rows),
            "controls": self.rows,
        }
        (self.final / "fixture-part-001.json").write_text(json.dumps(document, indent=2) + "\n", encoding="utf-8")


class BuildControlCheckBindingsTest(unittest.TestCase):
    def setUp(self) -> None:
        self.temporary = tempfile.TemporaryDirectory()
        self.addCleanup(self.temporary.cleanup)
        self.fixture = BindingFixture(Path(self.temporary.name))

    def build(self) -> dict:
        return bindings.build(self.fixture.registry, self.fixture.final)

    def test_emits_exactly_one_binding_for_confirmed_deterministic_control(self) -> None:
        artifact = self.build()
        bindings.validate_artifact(artifact, self.fixture.schema)
        self.assertEqual(artifact["binding_count"], 1)
        self.assertEqual([item["control_id"] for item in artifact["bindings"]], ["USEQ-AAAAAAAA"])
        clause = artifact["bindings"][0]["clauses"][0]
        self.assertEqual(clause["implementation_id"], "prc.check.environment-evidence@0.1")
        self.assertRegex(clause["implementation_contract_sha256"], r"^[0-9a-f]{64}$")
        self.assertEqual(clause["evidence_authority"], "external_registry")
        self.assertFalse(clause["external_provider_claimed"])
        self.assertEqual(clause["result_contract"]["blocked"]["result"], "blocked")
        self.assertIn("provider_contract_missing", clause["result_contract"]["blocked"]["when_any"])
        descriptor = next(item for item in artifact["implementation_registry"] if item["checker_family"] == "environment_evidence")
        self.assertEqual(descriptor["registration_state"], "runtime_registration_required")
        self.assertEqual(descriptor["on_unregistered"], "blocked")
        self.assertFalse(descriptor["external_provider_claimed"])
        self.assertEqual(descriptor["implementation_contract_sha256"], clause["implementation_contract_sha256"])

    def test_output_is_deterministic_and_clause_ids_are_stable(self) -> None:
        first = self.build()
        second = self.build()
        self.assertEqual(bindings.serialize(first), bindings.serialize(second))
        self.assertEqual(first["bindings"][0]["clauses"][0]["clause_id"], second["bindings"][0]["clauses"][0]["clause_id"])

    def test_rejects_stale_rule_revision_or_semantic_hash(self) -> None:
        self.fixture.rows[0]["semantic_sha256"] = "e" * 64
        self.fixture.write_final()
        with self.assertRaisesRegex(bindings.BindingError, "stale"):
            self.build()

    def test_rejects_missing_final_control(self) -> None:
        self.fixture.rows.pop()
        self.fixture.write_final()
        with self.assertRaisesRegex(bindings.BindingError, "missing 1 controls"):
            self.build()

    def test_rejects_duplicate_final_control(self) -> None:
        self.fixture.rows.append(dict(self.fixture.rows[0]))
        self.fixture.write_final()
        with self.assertRaisesRegex(bindings.BindingError, "duplicate final classification"):
            self.build()

    def test_never_binds_a_nondeterministic_row_even_if_it_retains_partial_clauses(self) -> None:
        self.fixture.rows[1]["deterministic_clauses"] = [{
            "statement": "A partial check must not survive nondeterministic classification.",
            "checker_family": "inventory_fact",
            "evidence_authority": "repository",
        }]
        self.fixture.write_final()
        artifact = self.build()
        self.assertEqual(artifact["binding_count"], 1)
        self.assertNotIn("USEQ-BBBBBBBB", {item["control_id"] for item in artifact["bindings"]})

    def test_project_corpus_generates_all_confirmed_bindings(self) -> None:
        artifact = bindings.build(
            ROOT / "catalog" / "control-id-registry.json",
            ROOT / "research" / "control-classification" / "final",
        )
        bindings.validate_artifact(artifact, ROOT / "schemas" / "control-check-bindings.schema.json")
        self.assertEqual(artifact["binding_count"], 686)
        self.assertEqual(len(artifact["implementation_registry"]), 11)


if __name__ == "__main__":
    unittest.main()
