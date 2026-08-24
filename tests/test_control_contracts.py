from __future__ import annotations

import importlib.util
import json
import unittest
from collections import Counter
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SCRIPT = ROOT / "scripts" / "control_contracts.py"
SPEC = importlib.util.spec_from_file_location("control_contracts", SCRIPT)
assert SPEC is not None and SPEC.loader is not None
CONTRACTS = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(CONTRACTS)


class ControlContractTests(unittest.TestCase):
    def test_generated_contracts_are_current_complete_and_valid(self) -> None:
        document = CONTRACTS.build_document()
        CONTRACTS.validate(document)
        actual = json.loads(CONTRACTS.OUTPUT.read_text(encoding="utf-8"))
        self.assertEqual(actual, document)
        self.assertEqual(document["contract_count"], 10_042)
        self.assertEqual(len({item["control_id"] for item in document["contracts"]}), 10_042)
        self.assertTrue(all(item["reviewer_status"] == "generated_unreviewed" for item in document["contracts"] if item["contract_status"] != "retired"))

    def test_contracts_do_not_pretend_generated_triage_is_expert_review(self) -> None:
        document = CONTRACTS.build_document()
        statuses = Counter(item["contract_status"] for item in document["contracts"])
        self.assertEqual(statuses["reviewed"], 0)
        self.assertGreater(statuses["generated_unreviewed"], 10_000)
        aliases = [item for item in document["contracts"] if item["canonical_control_id"] != item["control_id"]]
        registry = json.loads(CONTRACTS.REGISTRY.read_text(encoding="utf-8"))
        statements = {
            item["id"]: CONTRACTS.normalized_statement(item["statement"])
            for item in registry["entries"]
        }
        self.assertTrue(all(statements[item["canonical_control_id"]] == statements[item["control_id"]] for item in aliases))


if __name__ == "__main__":
    unittest.main()
