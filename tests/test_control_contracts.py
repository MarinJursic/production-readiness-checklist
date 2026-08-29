from __future__ import annotations

import copy
import hashlib
import importlib.util
import json
import tempfile
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
    def test_generated_catalog_is_current_complete_reviewed_and_valid(self) -> None:
        document = CONTRACTS.build_document()
        CONTRACTS.validate(document)
        actual_bytes = CONTRACTS.OUTPUT.read_bytes()
        self.assertEqual(json.loads(actual_bytes), document)
        self.assertEqual(actual_bytes.decode("utf-8"), CONTRACTS.generated_text(document))
        self.assertEqual(document["schema_version"], "prc.control-contracts/v0.2")
        self.assertEqual(document["contract_count"], 10_042)
        self.assertEqual(document["binding_count"], 686)
        self.assertEqual(
            Counter(item["classification"] for item in document["contracts"]),
            {"deterministic": 686, "nondeterministic": 9_356},
        )
        self.assertTrue(all(item["contract_status"] == "reviewed" for item in document["contracts"]))
        self.assertTrue(
            all(item["reviewer_status"] == "agent_reviewed" for item in document["contracts"])
        )
        self.assertEqual(
            document["control_check_bindings_sha256"],
            hashlib.sha256(CONTRACTS.BINDINGS.read_bytes()).hexdigest(),
        )

    def test_deterministic_contracts_map_one_to_one_to_exact_bindings(self) -> None:
        document = CONTRACTS.build_document()
        source = json.loads(CONTRACTS.BINDINGS.read_text(encoding="utf-8"))
        bindings = {item["control_id"]: item for item in source["bindings"]}
        deterministic = {
            item["control_id"]: item
            for item in document["contracts"]
            if item["classification"] == "deterministic"
        }
        self.assertEqual(set(deterministic), set(bindings))
        for control_id, contract in deterministic.items():
            binding = bindings[control_id]
            self.assertEqual(contract["deterministic_binding_id"], f"{control_id}@{contract['revision']}")
            self.assertEqual(contract["deterministic_binding_sha256"], CONTRACTS.digest_value(binding))
            self.assertEqual(contract["classification_row_sha256"], binding["final_row_sha256"])
        nondeterministic = [
            item for item in document["contracts"] if item["classification"] == "nondeterministic"
        ]
        self.assertTrue(
            all(
                item["deterministic_binding_id"] is None
                and item["deterministic_binding_sha256"] is None
                for item in nondeterministic
            )
        )

    def linked_final_corpus(self, root: Path) -> Path:
        final_root = root / "final"
        final_root.mkdir()
        for source in CONTRACTS.FINAL_ROOT.glob("*.json"):
            (final_root / source.name).symlink_to(source)
        return final_root

    def test_rejects_missing_final_packet(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            final_root = self.linked_final_corpus(Path(temporary))
            next(final_root.glob("*.json")).unlink()
            with self.assertRaisesRegex(CONTRACTS.ContractError, "differ from summary"):
                CONTRACTS.build_document(final_root=final_root)

    def test_rejects_duplicate_final_control(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            final_root = self.linked_final_corpus(Path(temporary))
            target = next(
                path
                for path in final_root.glob("*.json")
                if len(json.loads(path.read_text(encoding="utf-8"))["controls"]) > 1
            )
            document = json.loads(target.read_text(encoding="utf-8"))
            target.unlink()
            document["controls"][1] = copy.deepcopy(document["controls"][0])
            target.write_text(json.dumps(document), encoding="utf-8")
            with self.assertRaisesRegex(CONTRACTS.ContractError, "duplicate final classification"):
                CONTRACTS.build_document(final_root=final_root)

    def test_rejects_stale_binding_and_missing_binding(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            binding_path = Path(temporary) / "bindings.json"
            document = json.loads(CONTRACTS.BINDINGS.read_text(encoding="utf-8"))
            document["bindings"][0]["final_row_sha256"] = "0" * 64
            binding_path.write_text(json.dumps(document), encoding="utf-8")
            with self.assertRaisesRegex(CONTRACTS.ContractError, "binding is stale"):
                CONTRACTS.build_document(bindings_path=binding_path)

            document = json.loads(CONTRACTS.BINDINGS.read_text(encoding="utf-8"))
            document["bindings"].pop()
            document["binding_count"] -= 1
            binding_path.write_text(json.dumps(document), encoding="utf-8")
            with self.assertRaisesRegex(CONTRACTS.ContractError, "not one-to-one"):
                CONTRACTS.build_document(bindings_path=binding_path)

    def test_schema_forbids_binding_on_nondeterministic_contract(self) -> None:
        document = CONTRACTS.build_document()
        changed = copy.deepcopy(document)
        contract = next(
            item for item in changed["contracts"] if item["classification"] == "nondeterministic"
        )
        contract["deterministic_binding_id"] = f"{contract['control_id']}@{contract['revision']}"
        contract["deterministic_binding_sha256"] = "a" * 64
        contract["contract_sha256"] = CONTRACTS.digest_value(
            {key: value for key, value in contract.items() if key != "contract_sha256"}
        )
        with self.assertRaisesRegex(CONTRACTS.ContractError, "schema failed"):
            CONTRACTS.validate(changed)


if __name__ == "__main__":
    unittest.main()
