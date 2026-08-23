from __future__ import annotations

import importlib.util
import json
import unittest
from pathlib import Path

import yaml


ROOT = Path(__file__).resolve().parents[1]
SPEC = importlib.util.spec_from_file_location(
    "prc_validate_instance", ROOT / "scripts" / "validate_instance.py"
)
assert SPEC and SPEC.loader
validate_instance = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(validate_instance)


class ScannerOutputSchemaTests(unittest.TestCase):
    def test_minimal_inventory_conforms(self) -> None:
        instance = {
            "schema_version": "prc.inventory/v0.1",
            "target_name": "example",
            "digest": "a" * 64,
            "file_count": 0,
            "source_files": 0,
            "package_ecosystems": [],
            "manifests": [],
            "lock_files": [],
            "ci": {"github_actions": False},
            "files": [],
        }
        self.assertEqual(validate_instance.validation_errors(instance, "inventory.schema.json"), [])

    def test_invalid_inventory_is_rejected(self) -> None:
        instance = {"schema_version": "wrong"}
        errors = validate_instance.validation_errors(instance, "inventory.schema.json")
        self.assertTrue(errors)
        self.assertTrue(any("prc.inventory/v0.1" in error for error in errors))

    def test_minimal_remediation_candidate_conforms(self) -> None:
        digest = "a" * 64
        contract = {
            "schema_version": "prc.fix-contract/v0.1",
            "task_id": digest,
            "baseline_run_id": digest,
            "baseline_inventory_digest": digest,
            "assertion_id": "PRC-A-CORE-014",
            "control_ids": ["USEQ-DAF77C8F"],
            "goal": "Append one final line-feed byte.",
            "fixer_id": "prc.fixer.final-newline@0.1",
            "remediation_class": "R1",
            "allowed_paths": ["app.py"],
            "protected_paths": [".git/"],
            "network": "deny",
            "max_changed_lines": 20,
            "max_files": 20,
            "max_attempts": 1,
            "acceptance": ["The target assertion passes."],
        }
        candidate = {
            "schema_version": "prc.remediation-candidate/v0.1",
            "candidate_id": digest,
            "candidate_path": "/tmp/candidate",
            "contract": contract,
            "candidate_inventory_digest": digest,
            "candidate_run_id": digest,
            "changes": [{
                "path": "app.py",
                "kind": "modified",
                "before_sha256": digest,
                "after_sha256": "b" * 64,
                "before_mode": 420,
                "after_mode": 420,
                "added_lines": 1,
                "removed_lines": 0,
            }],
            "before_assessment": "fail",
            "after_assessment": "pass",
            "accepted": True,
            "reasons": [],
        }
        self.assertEqual(
            validate_instance.validation_errors(contract, "fix-contract.schema.json"), []
        )
        self.assertEqual(
            validate_instance.validation_errors(
                candidate, "remediation-candidate.schema.json"
            ),
            [],
        )

    def test_checked_in_adapter_manifest_conforms(self) -> None:
        path = ROOT / "fixtures" / "adapters" / "fixture-adapter.yaml"
        instance = yaml.safe_load(path.read_text(encoding="utf-8"))
        self.assertEqual(
            validate_instance.validation_errors(instance, "adapter-manifest.schema.json"), []
        )

    def test_adapter_jsonl_messages_conform_and_authority_attack_fails(self) -> None:
        valid_path = ROOT / "fixtures" / "adapters" / "valid-output.jsonl"
        for line in valid_path.read_text(encoding="utf-8").splitlines():
            instance = json.loads(line)
            self.assertEqual(
                validate_instance.validation_errors(instance, "adapter-message.schema.json"), []
            )
        malicious_path = ROOT / "fixtures" / "adapters" / "malicious-authority-output.jsonl"
        first_line = malicious_path.read_text(encoding="utf-8").splitlines()[0]
        self.assertTrue(
            validate_instance.validation_errors(
                json.loads(first_line), "adapter-message.schema.json"
            )
        )


if __name__ == "__main__":
    unittest.main()
