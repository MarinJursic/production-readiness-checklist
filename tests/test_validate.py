from __future__ import annotations

import importlib.util
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch


ROOT = Path(__file__).resolve().parents[1]
SPEC = importlib.util.spec_from_file_location("prc_validate", ROOT / "scripts" / "validate.py")
assert SPEC and SPEC.loader
validate = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(validate)


class NormalizeControlTests(unittest.TestCase):
    def test_normalizes_case_quotes_markup_and_terminal_period(self) -> None:
        left = validate.normalize_control("Use **Current** ‘Evidence’.")
        right = validate.normalize_control("use current 'evidence'")
        self.assertEqual(left, right)

    def test_duplicate_values_are_unique_and_sorted(self) -> None:
        self.assertEqual(validate.duplicate_values(["b", "a", "b", "b", "c", "a"]), ["a", "b"])


class TechnologyNeutralityTests(unittest.TestCase):
    def test_vendor_detection_is_case_insensitive(self) -> None:
        for value in ("AWS", "aws", "AwS", "kubernetes", "PYTHON"):
            with self.subTest(value=value):
                self.assertIsNotNone(validate.IMPLEMENTATION_SPECIFIC.search(value))

    def test_generic_terms_do_not_match(self) -> None:
        self.assertIsNone(validate.IMPLEMENTATION_SPECIFIC.search("managed object storage"))


class MarkdownAnchorTests(unittest.TestCase):
    def test_generates_expected_anchors_and_duplicate_suffixes(self) -> None:
        content = "# Release gate\n\n## 2. Immediate no-go conditions\n\n# Release gate\n"
        self.assertEqual(
            validate.markdown_anchors(content),
            {"release-gate", "release-gate_1", "2-immediate-no-go-conditions"},
        )

    def test_validator_reports_missing_file_anchor_and_escape(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            docs = root / "docs"
            docs.mkdir()
            (docs / "target.md").write_text("# Existing heading\n", encoding="utf-8")
            (docs / "source.md").write_text(
                "[good](target.md#existing-heading)\n"
                "[bad anchor](target.md#missing)\n"
                "[bad file](absent.md)\n"
                "[escape](../../outside.md)\n",
                encoding="utf-8",
            )
            errors: list[str] = []
            with (
                patch.object(validate, "ROOT", root),
                patch.object(validate, "CHECKLIST_DIR", docs / "checklists"),
                patch.object(validate, "ENGINEERING_DIR", docs / "engineering"),
            ):
                checked = validate.validate_links(errors)

            self.assertEqual(checked, 4)
            self.assertEqual(len(errors), 3)
            self.assertTrue(any("missing Markdown anchor" in error for error in errors))
            self.assertTrue(any("broken local link" in error for error in errors))
            self.assertTrue(any("escapes repository" in error for error in errors))

    def test_same_page_anchor_is_checked(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            docs = root / "docs"
            docs.mkdir()
            (docs / "page.md").write_text(
                "# Present\n\n[valid](#present)\n[invalid](#absent)\n",
                encoding="utf-8",
            )
            errors: list[str] = []
            with patch.object(validate, "ROOT", root):
                checked = validate.validate_links(errors)

            self.assertEqual(checked, 2)
            self.assertEqual(len(errors), 1)
            self.assertIn("missing Markdown anchor", errors[0])


class ClassificationArtifactTests(unittest.TestCase):
    def test_current_classification_and_binding_documents_are_consistent(self) -> None:
        errors: list[str] = []
        count = validate.validate_classification_documents(errors)
        self.assertEqual(errors, [])
        self.assertEqual(count, 686)

    def test_generated_validator_failure_is_reported_fail_closed(self) -> None:
        errors: list[str] = []
        completed = validate.subprocess.CompletedProcess(
            args=["validator"], returncode=1, stdout="stale binding artifact\n"
        )
        with patch.object(validate.subprocess, "run", return_value=completed):
            validate.run_generated_validator(["validator"], "Binding validation", errors)
        self.assertEqual(errors, ["Binding validation failed: stale binding artifact"])


class BenchmarkDocumentationTests(unittest.TestCase):
    def test_current_benchmark_counts_are_measured_and_documented(self) -> None:
        errors: list[str] = []
        self.assertEqual(validate.validate_benchmark_documentation(errors), (36, 144))
        self.assertEqual(errors, [])

    def test_stale_benchmark_count_is_rejected(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            suite = root / "fixtures" / "benchmarks" / "core-native"
            docs = root / "docs" / "scanner"
            suite.mkdir(parents=True)
            docs.mkdir(parents=True)
            (suite / "suite-comprehensive.yaml").write_text(
                "cases:\n  - id: one\n    expectations:\n      - {assertion_id: A}\n",
                encoding="utf-8",
            )
            (docs / "benchmarks.md").write_text(
                "The suite has a 2-case, 2-expectation fixture corpus.\n",
                encoding="utf-8",
            )
            errors: list[str] = []
            with patch.object(validate, "ROOT", root):
                self.assertEqual(validate.validate_benchmark_documentation(errors), (1, 1))
            self.assertEqual(len(errors), 1)
            self.assertIn("stale", errors[0])


if __name__ == "__main__":
    unittest.main()
