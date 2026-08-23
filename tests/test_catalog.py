from __future__ import annotations

import copy
import importlib.util
import unittest
from pathlib import Path
from unittest.mock import patch


ROOT = Path(__file__).resolve().parents[1]
SPEC = importlib.util.spec_from_file_location("prc_catalog", ROOT / "scripts" / "catalog.py")
assert SPEC and SPEC.loader
catalog = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(catalog)


class CatalogValidationTests(unittest.TestCase):
    def test_catalog_references_and_schemas_validate(self) -> None:
        objectives, assertions, profiles = catalog.validate_catalog()
        self.assertEqual(len(objectives), 24)
        self.assertEqual(len(assertions), 34)
        self.assertEqual([profile["id"] for profile in profiles], ["prc/core-repository"])

    def test_registry_freezes_every_current_control(self) -> None:
        registry = catalog.registry_document()
        entries = registry["entries"]
        self.assertEqual(registry["entry_count"], 10_042)
        self.assertEqual(len({entry["id"] for entry in entries}), 10_042)
        self.assertTrue(all(entry["revision"] >= 1 for entry in entries))
        self.assertTrue(all(len(entry["semantic_sha256"]) == 64 for entry in entries))

    def test_checked_in_generated_files_have_no_drift(self) -> None:
        for path, expected in catalog.generated_files().items():
            with self.subTest(path=path):
                self.assertTrue(path.exists())
                self.assertEqual(path.read_text(encoding="utf-8"), expected)

    def test_duplicate_ids_are_rejected(self) -> None:
        record = {"id": "PRC-A-CORE-001"}
        with self.assertRaisesRegex(catalog.CatalogError, "duplicate assertion ID"):
            catalog.index_unique([record, copy.deepcopy(record)], "assertion")

    def test_semantic_digest_ignores_formatting_but_not_meaning(self) -> None:
        self.assertEqual(
            catalog.semantic_digest("Use **current** evidence."),
            catalog.semantic_digest("use current evidence"),
        )
        self.assertNotEqual(
            catalog.semantic_digest("Use current evidence."),
            catalog.semantic_digest("Use optional evidence."),
        )

    def test_normative_change_increments_revision_without_changing_id(self) -> None:
        old = {
            "id": "USEQ-00000000",
            "status": "active",
            "revision": 4,
            "statement": "Old statement.",
            "semantic_sha256": catalog.semantic_digest("Old statement."),
            "source": {"path": "docs/example.md", "line": 1},
        }
        new = copy.deepcopy(old)
        new["statement"] = "New statement."
        new["semantic_sha256"] = catalog.semantic_digest(new["statement"])
        with (
            patch.object(catalog, "source_control_entries", return_value=([new], "digest")),
            patch.object(catalog, "previous_registry_entries", return_value={old["id"]: old}),
        ):
            result = catalog.registry_document()["entries"][0]

        self.assertEqual(result["id"], old["id"])
        self.assertEqual(result["revision"], 5)


if __name__ == "__main__":
    unittest.main()
