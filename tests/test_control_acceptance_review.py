import importlib.util
import json
import re
import unittest
from collections import Counter
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SCRIPT = ROOT / "scripts" / "generate_control_acceptance_review.py"
SPEC = importlib.util.spec_from_file_location("generate_control_acceptance_review", SCRIPT)
assert SPEC is not None and SPEC.loader is not None
GENERATOR = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(GENERATOR)


class ControlAcceptanceReviewTests(unittest.TestCase):
    def test_generated_review_is_current_and_complete(self) -> None:
        expected = GENERATOR.generated_text()
        actual = GENERATOR.OUTPUT.read_text(encoding="utf-8")
        self.assertEqual(actual, expected)

        entries = json.loads(GENERATOR.REGISTRY.read_text(encoding="utf-8"))["entries"]
        self.assertEqual(len(entries), 10_042)
        expected_ids = Counter(entry["id"] for entry in entries)
        begin_ids = Counter(re.findall(r"<!-- BEGIN CONTROL ([A-Z0-9-]+) -->", actual))
        end_ids = Counter(re.findall(r"<!-- END CONTROL ([A-Z0-9-]+) -->", actual))
        self.assertEqual(begin_ids, expected_ids)
        self.assertEqual(end_ids, expected_ids)

    def test_current_coverage_numbers_are_not_overstated(self) -> None:
        text = GENERATOR.generated_text()
        self.assertIn("Controls connected to current scanner assertions: **26**", text)
        self.assertIn("Current assertions connected to those controls: **43**", text)
        self.assertIn("Not checked today.", text)


if __name__ == "__main__":
    unittest.main()
