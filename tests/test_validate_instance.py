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
    def test_checked_in_project_configuration_conforms(self) -> None:
        path = ROOT / "fixtures" / "config" / "production-readiness.yaml"
        instance = yaml.safe_load(path.read_text(encoding="utf-8"))
        self.assertEqual(
            validate_instance.validation_errors(
                instance, "project-config.schema.json"
            ),
            [],
        )

    def test_minimal_inventory_conforms(self) -> None:
        instance = {
            "schema_version": "prc.inventory/v0.3",
            "target_name": "example",
            "digest": "a" * 64,
            "file_count": 0,
            "source_files": 0,
            "package_ecosystems": [],
            "manifests": [],
            "lock_files": [],
            "container_files": [],
            "symlinks": [],
            "ci": {"github_actions": False, "workflow_files": []},
            "infrastructure": {"terraform_files": [], "kubernetes_files": []},
            "components": [{"id": "repository:.", "kind": "repository", "path": "."}],
            "relations": [],
            "facts": [{
                "key": "repository.detected",
                "value": "true",
                "source": ".",
                "detector": "prc.inventory.repository",
                "detector_version": "0.3",
                "confidence": 1,
                "scope_path": ".",
                "limitations": [],
            }],
            "files": [],
        }
        self.assertEqual(validate_instance.validation_errors(instance, "inventory.schema.json"), [])

        legacy = {
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
        self.assertEqual(validate_instance.validation_errors(legacy, "inventory-v0.1.schema.json"), [])

        legacy_v2 = {**instance, "schema_version": "prc.inventory/v0.2"}
        self.assertEqual(
            validate_instance.validation_errors(
                legacy_v2, "inventory-v0.2.schema.json"
            ),
            [],
        )

    def test_invalid_inventory_is_rejected(self) -> None:
        instance = {"schema_version": "wrong"}
        errors = validate_instance.validation_errors(instance, "inventory.schema.json")
        self.assertTrue(errors)
        self.assertTrue(any("prc.inventory/v0.3" in error for error in errors))

    def test_legacy_run_contract_remains_validatable(self) -> None:
        digest = "a" * 64
        legacy = {
            "schema_version": "prc.run/v0.1",
            "run_id": digest,
            "started_at": "2026-08-23T12:00:00Z",
            "completed_at": "2026-08-23T12:00:01Z",
            "plan": {
                "schema_version": "prc.plan/v0.1",
                "digest": digest,
                "target_name": "example",
                "inventory_digest": digest,
                "profile_id": "prc/core-repository",
                "profile_version": "0.2",
                "assertions": [],
            },
            "inventory": {
                "schema_version": "prc.inventory/v0.1",
                "target_name": "example",
                "digest": digest,
                "file_count": 0,
                "source_files": 0,
                "package_ecosystems": [],
                "manifests": [],
                "lock_files": [],
                "ci": {"github_actions": False},
                "files": [],
            },
            "results": [],
            "terminal_state": "profile_satisfied",
        }
        self.assertEqual(
            validate_instance.validation_errors(
                legacy, "run-result-v0.1.schema.json"
            ),
            [],
        )

    def test_current_plan_records_applicability_reason(self) -> None:
        digest = "a" * 64
        plan = {
            "schema_version": "prc.plan/v0.3",
            "digest": digest,
            "target_name": "example",
            "inventory_digest": digest,
            "profile_id": "prc/core-repository",
            "profile_version": "0.3",
            "artifact_digests": [],
            "target_environments": [],
            "assertions": [{
                "assertion_id": "PRC-A-CORE-001",
                "implementation_id": "prc.native.file-present@0.1",
                "applicability": "applicable",
                "applicability_evaluator": "cel-go/v0.30.0+prc-inventory/v0.3",
                "applicability_reason": "CEL expression evaluated to true.",
            }],
        }
        self.assertEqual(
            validate_instance.validation_errors(plan, "plan.schema.json"), []
        )
        del plan["assertions"][0]["applicability_reason"]
        self.assertTrue(
            validate_instance.validation_errors(plan, "plan.schema.json")
        )

    def test_v02_run_contract_remains_validatable(self) -> None:
        digest = "a" * 64
        inventory = {
            "schema_version": "prc.inventory/v0.2",
            "target_name": "example",
            "digest": digest,
            "file_count": 0,
            "source_files": 0,
            "package_ecosystems": [],
            "manifests": [],
            "lock_files": [],
            "container_files": [],
            "symlinks": [],
            "ci": {"github_actions": False, "workflow_files": []},
            "infrastructure": {"terraform_files": [], "kubernetes_files": []},
            "components": [{"id": "repository:.", "kind": "repository", "path": "."}],
            "relations": [],
            "facts": [],
            "files": [],
        }
        plan = {
            "schema_version": "prc.plan/v0.2",
            "digest": digest,
            "target_name": "example",
            "inventory_digest": digest,
            "profile_id": "prc/core-repository",
            "profile_version": "0.3",
            "assertions": [{
                "assertion_id": "PRC-A-CORE-001",
                "implementation_id": "prc.native.file-present@0.1",
                "applicability": "applicable",
                "applicability_evaluator": "cel-go/v0.30.0+prc-inventory/v0.2",
                "applicability_reason": "CEL expression evaluated to true.",
            }],
        }
        run = {
            "schema_version": "prc.run/v0.2",
            "run_id": digest,
            "started_at": "2026-08-23T12:00:00Z",
            "completed_at": "2026-08-23T12:00:01Z",
            "plan": plan,
            "inventory": inventory,
            "adapter_executions": [],
            "results": [],
            "terminal_state": "profile_satisfied",
        }
        self.assertEqual(
            validate_instance.validation_errors(
                run, "run-result-v0.2.schema.json"
            ),
            [],
        )

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
        r2_contract = {
            **contract,
            "remediation_class": "R2",
            "fixer_id": "prc.provider-proposal@0.1",
            "provider": "codex",
            "proposal_task_id": digest,
            "proposal_sha256": "b" * 64,
        }
        self.assertEqual(
            validate_instance.validation_errors(
                r2_contract, "fix-contract.schema.json"
            ),
            [],
        )
        del r2_contract["proposal_sha256"]
        self.assertTrue(
            validate_instance.validation_errors(
                r2_contract, "fix-contract.schema.json"
            )
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

    def test_checked_in_agent_task_and_output_conform(self) -> None:
        task = json.loads(
            (ROOT / "fixtures" / "providers" / "suggest-task.json").read_text(
                encoding="utf-8"
            )
        )
        output = json.loads(
            (ROOT / "fixtures" / "providers" / "valid-output.json").read_text(
                encoding="utf-8"
            )
        )
        malicious = json.loads(
            (
                ROOT
                / "fixtures"
                / "providers"
                / "malicious-capability-output.json"
            ).read_text(encoding="utf-8")
        )
        self.assertEqual(
            validate_instance.validation_errors(task, "agent-task.schema.json"), []
        )
        self.assertEqual(
            validate_instance.validation_errors(output, "agent-output.schema.json"), []
        )
        self.assertTrue(
            validate_instance.validation_errors(
                malicious, "agent-output.schema.json"
            )
        )

        digest = "a" * 64
        execution = {
            "schema_version": "prc.agent-execution/v0.1",
            "execution_id": digest,
            "provider": "codex",
            "task_id": task["task_id"],
            "executable_sha256": digest,
            "output_schema_sha256": digest,
            "started_at": "2026-08-23T10:00:00Z",
            "completed_at": "2026-08-23T10:00:01Z",
            "duration_ms": 1000,
            "stdout_path": "/tmp/stdout.log",
            "stdout_sha256": digest,
            "stdout_bytes": 1,
            "stderr_path": "/tmp/stderr.log",
            "stderr_sha256": digest,
            "stderr_bytes": 0,
            "output": output,
        }
        self.assertEqual(
            validate_instance.validation_errors(
                execution, "agent-execution.schema.json"
            ),
            [],
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

    def test_bound_adapter_execution_conforms(self) -> None:
        digest = "a" * 64
        execution = {
            "schema_version": "prc.adapter-execution/v0.1",
            "execution_id": digest,
            "adapter_run_id": "b" * 64,
            "adapter_id": "prc.adapter.fixture@0.1",
            "manifest_sha256": "c" * 64,
            "image": "registry.example/prc/fixture@sha256:" + "d" * 64,
            "subject": {
                "target_name": "example",
                "inventory_digest": "e" * 64,
            },
            "started_at": "2026-08-23T12:00:00Z",
            "completed_at": "2026-08-23T12:00:01Z",
            "duration_ms": 1000,
            "diagnostics_sha256": "f" * 64,
            "diagnostics_bytes": 0,
            "transcript": {
                "logs": [],
                "observations": [{
                    "id": "OBS-1",
                    "kind": "fixture-result",
                    "outcome": "not_found",
                    "summary": "No fixture match.",
                    "locations": [],
                }],
                "artifacts": [],
                "summary": {
                    "status": "completed",
                    "counts": {"observations": 1},
                },
            },
        }
        self.assertEqual(
            validate_instance.validation_errors(
                execution, "adapter-execution.schema.json"
            ),
            [],
        )
        execution["transcript"]["observations"][0]["assessment"] = "pass"
        self.assertTrue(
            validate_instance.validation_errors(
                execution, "adapter-execution.schema.json"
            )
        )


if __name__ == "__main__":
    unittest.main()
