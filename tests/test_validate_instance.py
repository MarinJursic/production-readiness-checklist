from __future__ import annotations

import importlib.util
import unittest
from pathlib import Path


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


if __name__ == "__main__":
    unittest.main()
