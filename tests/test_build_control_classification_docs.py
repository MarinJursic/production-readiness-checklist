from __future__ import annotations

import importlib.util
import hashlib
import json
import sys
import tempfile
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SCRIPT = ROOT / "scripts" / "build_control_classification_docs.py"
SPEC = importlib.util.spec_from_file_location("build_control_classification_docs", SCRIPT)
assert SPEC is not None and SPEC.loader is not None
DOCS = importlib.util.module_from_spec(SPEC)
sys.modules[SPEC.name] = DOCS
SPEC.loader.exec_module(DOCS)


class BuildControlClassificationDocsTests(unittest.TestCase):
    def fixture(self, root: Path) -> tuple[Path, Path, Path, Path]:
        final_root = root / "final"
        packet_root = root / "packets"
        final_root.mkdir()
        packet_root.mkdir()
        method = "m" * 64
        registry = "r" * 64
        packet_id = "docs-example-md-part-001"
        source_path = "docs/example.md"
        sources = [
            {
                "control_id": "PRC-01-001",
                "revision": 1,
                "semantic_sha256": "a" * 64,
                "statement": "The release manifest identifies the exact artifact.",
                "source": {"path": source_path, "line": 1},
                "heading_trail": ["Example"],
            },
            {
                "control_id": "PRC-01-002",
                "revision": 1,
                "semantic_sha256": "b" * 64,
                "statement": "The product is easy to use for every intended user.",
                "source": {"path": source_path, "line": 2},
                "heading_trail": ["Example"],
            },
        ]
        packet = {
            "schema_version": "prc.control-classification-packet/v0.1",
            "packet_id": packet_id,
            "methodology_sha256": method,
            "registry_version": "0.1.0",
            "registry_sha256": registry,
            "source_path": source_path,
            "control_count": 2,
            "controls": sources,
        }
        deterministic = {
            "control_id": "PRC-01-001",
            "revision": 1,
            "semantic_sha256": "a" * 64,
            "classification": "deterministic",
            "route": "artifact_verification",
            "reason": "The artifact digest gives the complete bounded identity fact.",
            "deterministic_clauses": [
                {
                    "statement": "The release manifest digest equals the assessed artifact digest.",
                    "checker_family": "artifact_integrity",
                    "evidence_authority": "artifact",
                }
            ],
            "nondeterministic_reason_codes": [],
            "decision_basis": "strength_audit_confirmed",
            "skeptical_verdict": "confirmed_deterministic",
            "counterexample_analysis": "Changing the artifact changes its digest and fails the equality check.",
        }
        rejected = {
            "control_id": "PRC-01-002",
            "revision": 1,
            "semantic_sha256": "b" * 64,
            "classification": "nondeterministic",
            "route": "contextual_judgment",
            "reason": "Ease of use across all intended users requires contextual user research.",
            "deterministic_clauses": [],
            "nondeterministic_reason_codes": ["contextual_judgment"],
            "decision_basis": "skeptically_rejected",
            "skeptical_verdict": "rejected_nondeterministic",
            "counterexample_analysis": "A scripted happy path can pass while a keyboard-only user cannot finish it.",
        }
        final = {
            "schema_version": "prc.control-classification-final/v0.1",
            "packet_id": packet_id,
            "methodology_sha256": method,
            "registry_version": "0.1.0",
            "registry_sha256": registry,
            "source_path": source_path,
            "control_count": 2,
            "controls": [deterministic, rejected],
        }
        registry = {
            "schema_version": "prc.control-id-registry/v0.1",
            "registry_version": "0.1.0",
            "source_sha256": {},
            "entry_count": 2,
            "entries": [
                {
                    "id": source["control_id"],
                    "status": "active",
                    "revision": source["revision"],
                    "statement": source["statement"],
                    "semantic_sha256": source["semantic_sha256"],
                    "source": source["source"],
                }
                for source in sources
            ],
        }
        registry_path = root / "control-id-registry.json"
        registry_path.write_text(json.dumps(registry), encoding="utf-8")
        summary = {
            "schema_version": "prc.control-classification-summary/v0.1",
            "methodology_sha256": method,
            "registry_version": "0.1.0",
            "registry_sha256": hashlib.sha256(registry_path.read_bytes()).hexdigest(),
            "control_count": 2,
            "deterministic": 1,
            "nondeterministic": 1,
            "decision_counts": {"skeptically_rejected": 1, "strength_audit_confirmed": 1},
            "route_counts": {"artifact_verification": 1, "contextual_judgment": 1},
            "packets": [
                {
                    "packet_id": packet_id,
                    "source_path": source_path,
                    "control_count": 2,
                    "deterministic": 1,
                    "nondeterministic": 1,
                }
            ],
        }
        (packet_root / f"{packet_id}.json").write_text(json.dumps(packet), encoding="utf-8")
        (final_root / f"{packet_id}.json").write_text(json.dumps(final), encoding="utf-8")
        summary_path = root / "summary.json"
        summary_path.write_text(json.dumps(summary), encoding="utf-8")
        for document in (packet, final):
            document["registry_sha256"] = summary["registry_sha256"]
        (packet_root / f"{packet_id}.json").write_text(json.dumps(packet), encoding="utf-8")
        (final_root / f"{packet_id}.json").write_text(json.dumps(final), encoding="utf-8")
        return final_root, packet_root, summary_path, registry_path

    def test_renders_every_required_human_readable_field(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            final_root, packet_root, summary, registry = self.fixture(Path(temporary))
            corpus = DOCS.load_corpus(
                final_root=final_root, packet_root=packet_root, summary_path=summary,
                registry_path=registry,
            )
            files = DOCS.render_files(corpus, target_part_bytes=2_000, max_part_bytes=3_000)
            rendered = "\n".join(files.values())
            self.assertEqual(rendered.count("## `PRC-01-001`"), 1)
            self.assertEqual(rendered.count("## `PRC-01-002`"), 1)
            self.assertIn("The release manifest identifies the exact artifact.", rendered)
            self.assertIn("**Classification:** **Deterministic**", rendered)
            self.assertIn("**Route:** `artifact_verification`", rendered)
            self.assertIn("**Checker family:** `artifact_integrity`", rendered)
            self.assertIn("**Evidence authority:** `artifact`", rendered)
            self.assertIn("Skeptical rejection counterexample", rendered)
            self.assertIn("keyboard-only user", rendered)
            self.assertIn("AI is advisory only for nondeterministic rules", files[Path("README.md")])

    def test_rejects_missing_final_packet(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            final_root, packet_root, summary, registry = self.fixture(Path(temporary))
            next(final_root.glob("*.json")).unlink()
            with self.assertRaisesRegex(ValueError, "missing final packets"):
                DOCS.load_corpus(
                    final_root=final_root, packet_root=packet_root, summary_path=summary,
                    registry_path=registry,
                )

    def test_rejects_duplicate_control(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            final_root, packet_root, summary, registry = self.fixture(Path(temporary))
            packet_path = next(packet_root.glob("*.json"))
            final_path = next(final_root.glob("*.json"))
            packet = json.loads(packet_path.read_text())
            final = json.loads(final_path.read_text())
            packet["controls"][1].update(
                control_id="PRC-01-001", semantic_sha256="a" * 64
            )
            final["controls"][1].update(
                control_id="PRC-01-001", semantic_sha256="a" * 64
            )
            packet_path.write_text(json.dumps(packet))
            final_path.write_text(json.dumps(final))
            with self.assertRaisesRegex(ValueError, "duplicate final classification"):
                DOCS.load_corpus(
                    final_root=final_root, packet_root=packet_root, summary_path=summary,
                    registry_path=registry,
                )

    def test_rejects_missing_control_even_when_local_counts_are_rewritten(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            final_root, packet_root, summary_path, registry = self.fixture(Path(temporary))
            packet_path = next(packet_root.glob("*.json"))
            final_path = next(final_root.glob("*.json"))
            packet = json.loads(packet_path.read_text())
            final = json.loads(final_path.read_text())
            summary = json.loads(summary_path.read_text())
            packet["controls"].pop()
            final["controls"].pop()
            packet["control_count"] = final["control_count"] = 1
            summary.update(control_count=1, deterministic=1, nondeterministic=0)
            summary["packets"][0].update(
                control_count=1, deterministic=1, nondeterministic=0
            )
            packet_path.write_text(json.dumps(packet))
            final_path.write_text(json.dumps(final))
            summary_path.write_text(json.dumps(summary))
            with self.assertRaisesRegex(ValueError, "does not match the active registry"):
                DOCS.load_corpus(
                    final_root=final_root,
                    packet_root=packet_root,
                    summary_path=summary_path,
                    registry_path=registry,
                )

    def test_rejects_stale_rule_binding(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            final_root, packet_root, summary, registry = self.fixture(Path(temporary))
            final_path = next(final_root.glob("*.json"))
            final = json.loads(final_path.read_text())
            final["controls"][0]["semantic_sha256"] = "z" * 64
            final_path.write_text(json.dumps(final))
            with self.assertRaisesRegex(ValueError, "stale or bound to another rule"):
                DOCS.load_corpus(
                    final_root=final_root, packet_root=packet_root, summary_path=summary,
                    registry_path=registry,
                )

    def test_split_files_stay_bounded_and_stale_output_fails(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            final_root, packet_root, summary, registry = self.fixture(root)
            corpus = DOCS.load_corpus(
                final_root=final_root, packet_root=packet_root, summary_path=summary,
                registry_path=registry,
            )
            files = DOCS.render_files(corpus, target_part_bytes=800, max_part_bytes=1_600)
            parts = {path: text for path, text in files.items() if path.name != "README.md"}
            self.assertGreaterEqual(len(parts), 2)
            self.assertTrue(all(len(text.encode("utf-8")) < 1_600 for text in parts.values()))
            output = root / "output"
            DOCS.write_files(files, output)
            DOCS.check_files(files, output)
            changed = next(output.rglob("*-part-*.md"))
            changed.write_text(changed.read_text() + "stale\n")
            with self.assertRaisesRegex(ValueError, "is stale"):
                DOCS.check_files(files, output)


if __name__ == "__main__":
    unittest.main()
