from __future__ import annotations

import importlib.util
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SCRIPT = ROOT / "scripts" / "control_classification_review.py"
SPEC = importlib.util.spec_from_file_location("control_classification_review", SCRIPT)
assert SPEC is not None and SPEC.loader is not None
REVIEW = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(REVIEW)


class ControlClassificationReviewTests(unittest.TestCase):
    def setUp(self) -> None:
        self.expected = {"control_id": "PRC-01-001", "revision": 2, "semantic_sha256": "a" * 64}

    def test_accepts_complete_deterministic_review(self) -> None:
        row = {
            "control_id": "PRC-01-001",
            "revision": 2,
            "semantic_sha256": "a" * 64,
            "classification": "deterministic",
            "route": "local_static",
            "reason": "A bounded repository file defines the complete binary result.",
            "deterministic_clauses": [
                {
                    "statement": "The required repository file exists.",
                    "checker_family": "inventory_fact",
                    "evidence_authority": "repository",
                }
            ],
            "nondeterministic_reason_codes": [],
            "review_status": "primary_agent_reviewed",
        }
        REVIEW.validate_review(row, self.expected)

    def test_accepts_nondeterministic_review_without_partial_checker(self) -> None:
        row = {
            "control_id": "PRC-01-001",
            "revision": 2,
            "semantic_sha256": "a" * 64,
            "classification": "nondeterministic",
            "route": "contextual_judgment",
            "reason": "Two informed reviewers can disagree about whether the design is clear.",
            "deterministic_clauses": [],
            "nondeterministic_reason_codes": ["contextual_judgment"],
            "review_status": "primary_agent_reviewed",
        }
        REVIEW.validate_review(row, self.expected)

    def test_rejects_partial_checker_as_deterministic(self) -> None:
        row = {
            "control_id": "PRC-01-001",
            "revision": 2,
            "semantic_sha256": "a" * 64,
            "classification": "deterministic",
            "route": "local_static",
            "reason": "A filename hint alone does not prove the complete control requirement.",
            "deterministic_clauses": [],
            "nondeterministic_reason_codes": [],
            "review_status": "primary_agent_reviewed",
        }
        with self.assertRaisesRegex(ValueError, "deterministic review is incomplete"):
            REVIEW.validate_review(row, self.expected)

    def test_rejects_ai_as_evidence_authority(self) -> None:
        row = {
            "control_id": "PRC-01-001",
            "revision": 2,
            "semantic_sha256": "a" * 64,
            "classification": "deterministic",
            "route": "local_static",
            "reason": "An AI answer cannot serve as authoritative proof for a pass result.",
            "deterministic_clauses": [
                {
                    "statement": "An AI model says the control passes.",
                    "checker_family": "analysis_adapter",
                    "evidence_authority": "ai",
                }
            ],
            "nondeterministic_reason_codes": [],
            "review_status": "primary_agent_reviewed",
        }
        with self.assertRaisesRegex(ValueError, "unsupported evidence authority"):
            REVIEW.validate_review(row, self.expected)

    def test_accepts_skeptical_confirmation_without_rejection_fields(self) -> None:
        row = {
            "control_id": "PRC-01-001",
            "revision": 2,
            "semantic_sha256": "a" * 64,
            "verdict": "confirmed_deterministic",
            "reason": "The bounded checker proves both pass and fail for the complete statement.",
            "counterexample_analysis": "No hidden violation remains when the complete authoritative inventory is sealed.",
            "nondeterministic_route": None,
            "nondeterministic_reason_codes": [],
            "review_status": "skeptical_agent_reviewed",
        }
        REVIEW.validate_skeptic_review(row, self.expected)

    def test_skeptical_rejection_requires_nondeterministic_reason(self) -> None:
        row = {
            "control_id": "PRC-01-001",
            "revision": 2,
            "semantic_sha256": "a" * 64,
            "verdict": "rejected_nondeterministic",
            "reason": "The proposed parser checks only the presence of a document.",
            "counterexample_analysis": "The document can exist while its guidance is wrong for the actual project.",
            "nondeterministic_route": None,
            "nondeterministic_reason_codes": [],
            "review_status": "skeptical_agent_reviewed",
        }
        with self.assertRaisesRegex(ValueError, "lacks a nondeterministic classification"):
            REVIEW.validate_skeptic_review(row, self.expected)

    def test_strengthened_clause_rejects_acceptance_contract_wrapper(self) -> None:
        clause = {
            "statement": "For the sealed assessed scope, verify the versioned acceptance contract for: Test rollback in production.",
            "checker_family": "structured_record",
            "evidence_authority": "structured_record",
        }
        with self.assertRaisesRegex(ValueError, "forbidden generic weakening"):
            REVIEW.validate_strengthened_clause(clause, "PRC-01-001")

    def test_strengthened_clause_rejects_known_generic_duplicate(self) -> None:
        clause = {
            "statement": "Authenticated bounded execution covers every scenario in the approved test matrix for the exact revision and records the required expected outcome for each case.",
            "checker_family": "execution_evidence",
            "evidence_authority": "executed",
        }
        with self.assertRaisesRegex(ValueError, "forbidden generic weakening"):
            REVIEW.validate_strengthened_clause(clause, "PRC-01-001")

    def test_strengthened_clause_rejects_provider_conclusion_wrapper(self) -> None:
        clause = {
            "statement": (
                "Read-only effective configuration, state, and event evidence for every subject "
                "in the complete assessed inventory directly demonstrates the operating behavior "
                "required by the original control; a policy alone is insufficient."
            ),
            "checker_family": "environment_evidence",
            "evidence_authority": "environment",
        }
        with self.assertRaisesRegex(ValueError, "forbidden generic weakening"):
            REVIEW.validate_strengthened_clause(clause, "PRC-01-001")

    def test_strengthened_clause_rejects_generic_record_wrapper(self) -> None:
        clause = {
            "statement": (
                "For every subject in the complete assessed inventory, an authenticated versioned "
                "record directly contains every identity, category, field, relation, threshold, "
                "date, scope item, and status explicitly required by the original control: Define it."
            ),
            "checker_family": "structured_record",
            "evidence_authority": "structured_record",
        }
        with self.assertRaisesRegex(ValueError, "forbidden generic weakening"):
            REVIEW.validate_strengthened_clause(clause, "PRC-01-001")

    def test_strengthened_clause_accepts_direct_full_strength_promise(self) -> None:
        clause = {
            "statement": "Authenticated rollback execution uses the exact production automation and production permission identities.",
            "checker_family": "execution_evidence",
            "evidence_authority": "executed",
        }
        REVIEW.validate_strengthened_clause(clause, "PRC-01-001")


if __name__ == "__main__":
    unittest.main()
