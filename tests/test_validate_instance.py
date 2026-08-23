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
    def test_catalog_manifest_and_bundle_conform(self) -> None:
        digest = "a" * 64
        manifest = {
            "schema_version": "prc.catalog-manifest/v0.1",
            "catalog_version": "0.1.0",
            "catalog_digest": digest,
            "objective_count": 1,
            "assertion_count": 1,
            "profile_count": 1,
        }
        bundle = {
            "schema_version": "prc.catalog-bundle/v0.1",
            "manifest": manifest,
            "objectives": [{
                "id": "USEQ-AAAAAAAA",
                "revision": 1,
                "title": "Test objective",
                "statement": "Test objective.",
                "source": {"path": "docs/source.md", "line": 1},
                "domains": ["repository"],
                "automation_class": "automated",
                "assertion_ids": ["PRC-A-TEST-001"],
            }],
            "assertions": [{
                "id": "PRC-A-TEST-001",
                "revision": 1,
                "control_ids": ["USEQ-AAAAAAAA"],
                "title": "Test assertion",
                "statement": "Test assertion.",
                "applicability": "true",
                "evidence_required": [{
                    "kind": "repository-file",
                    "minimum_authority": "repository",
                    "description": "Test evidence.",
                }],
                "implementation_id": "prc.native.test@0.1",
                "severity": "high",
                "gate": "required",
                "remediation_class": "R0",
            }],
            "profiles": [{
                "schema_version": "prc.profile/v0.1",
                "id": "prc/test",
                "version": "1.0",
                "title": "Test profile",
                "description": "Test profile.",
                "assertion_ids": ["PRC-A-TEST-001"],
                "terminal_policy": {
                    "block_on": ["high"],
                    "allow_manual_remaining": True,
                },
            }],
        }
        self.assertEqual(
            validate_instance.validation_errors(
                manifest, "catalog-manifest.schema.json"
            ),
            [],
        )
        self.assertEqual(
            validate_instance.validation_errors(
                bundle, "catalog-bundle.schema.json"
            ),
            [],
        )

    def test_canonical_finding_conforms(self) -> None:
        digest = "a" * 64
        finding = {
            "schema_version": "prc.finding/v0.1",
            "id": digest,
            "fingerprint": "b" * 64,
            "assertion_id": "PRC-A-CORE-008",
            "control_ids": ["USEQ-AAAAAAAA"],
            "title": "Action references are immutable",
            "summary": "One workflow action uses a mutable reference.",
            "severity": "critical",
            "gate": "no-go",
            "remediation_class": "R2",
            "subject": {
                "kind": "project",
                "id": "example",
                "inventory_digest": digest,
            },
            "locations": [{"path": ".github/workflows/ci.yml", "line": 8}],
            "evidence_ids": ["c" * 64],
        }
        self.assertEqual(
            validate_instance.validation_errors(finding, "finding.schema.json"),
            [],
        )

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
        capabilities = {
            "read_workspace": True,
            "write_scratch": False,
            "process": "deny",
            "network": "deny",
            "network_hosts": [],
            "secret_handles": [],
        }
        plan = {
            "schema_version": "prc.plan/v0.6",
            "digest": digest,
            "engine_version": "prc.engine/v0.1",
            "target_name": "example",
            "inventory_digest": digest,
            "profile_id": "prc/core-repository",
            "profile_version": "0.3",
            "profile_digest": digest,
            "catalog_digest": digest,
            "artifact_digests": [],
            "target_environments": [],
            "execution_mode": "inspect",
            "capability_policy": capabilities,
            "implementations": [{
                "id": "prc.native.file-present@0.1",
                "kind": "built-in",
                "assertion_ids": ["PRC-A-CORE-001"],
                "capabilities": capabilities,
                "status": "available",
            }],
            "adapters": [],
            "nodes": [
                {
                    "id": "inventory:" + digest,
                    "kind": "inventory",
                    "depends_on": [],
                    "capabilities": capabilities,
                    "status": "ready",
                },
                {
                    "id": "assertion:PRC-A-CORE-001",
                    "kind": "assertion",
                    "depends_on": ["inventory:" + digest],
                    "assertion_id": "PRC-A-CORE-001",
                    "implementation_id": "prc.native.file-present@0.1",
                    "capabilities": capabilities,
                    "status": "ready",
                },
                {
                    "id": "gate:core-repository",
                    "kind": "gate",
                    "depends_on": ["assertion:PRC-A-CORE-001"],
                    "capabilities": {
                        **capabilities,
                        "read_workspace": False,
                    },
                    "status": "ready",
                },
            ],
            "assertions": [{
                "assertion_id": "PRC-A-CORE-001",
                "assertion_revision": 2,
                "definition_digest": digest,
                "implementation_id": "prc.native.file-present@0.1",
                "applicability": "applicable",
                "applicability_evaluator": "cel-go/v0.30.0+prc-inventory/v0.3",
                "applicability_reason": "CEL expression evaluated to true.",
            }],
        }
        self.assertEqual(
            validate_instance.validation_errors(plan, "plan.schema.json"), []
        )
        self.assertEqual(
            validate_instance.validation_errors(plan, "plan-v0.6.schema.json"), []
        )
        del plan["assertions"][0]["applicability_reason"]
        self.assertTrue(
            validate_instance.validation_errors(plan, "plan.schema.json")
        )

    def test_v03_plan_contract_remains_validatable(self) -> None:
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
            "assertions": [],
        }
        self.assertEqual(
            validate_instance.validation_errors(plan, "plan-v0.3.schema.json"), []
        )
        inventory = {
            "schema_version": "prc.inventory/v0.3",
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
            "components": [],
            "relations": [],
            "facts": [],
            "files": [],
        }
        run = {
            "schema_version": "prc.run/v0.3",
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
            validate_instance.validation_errors(run, "run-result-v0.3.schema.json"), []
        )

    def test_v04_run_contract_remains_validatable(self) -> None:
        digest = "a" * 64
        plan = {
            "schema_version": "prc.plan/v0.4",
            "digest": digest,
            "engine_version": "prc.engine/v0.1",
            "target_name": "example",
            "inventory_digest": digest,
            "profile_id": "prc/core-repository",
            "profile_version": "0.3",
            "profile_digest": digest,
            "artifact_digests": [],
            "target_environments": [],
            "assertions": [],
        }
        inventory = {
            "schema_version": "prc.inventory/v0.3",
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
            "components": [],
            "relations": [],
            "facts": [],
            "files": [],
        }
        run = {
            "schema_version": "prc.run/v0.4",
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
            validate_instance.validation_errors(plan, "plan-v0.4.schema.json"), []
        )
        self.assertEqual(
            validate_instance.validation_errors(run, "run-result-v0.4.schema.json"), []
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
            "schema_version": "prc.fix-contract/v0.3",
            "task_id": digest,
            "baseline_run_id": digest,
            "baseline_inventory_digest": digest,
            "finding_id": digest,
            "finding_fingerprint": "b" * 64,
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
            "attempt": 1,
            "max_attempts": 1,
            "acceptance": ["The target assertion passes."],
        }
        candidate = {
            "schema_version": "prc.remediation-candidate/v0.3",
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
            "proposal_finding_id": "c" * 64,
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

        v02_contract = {
            **contract,
            "schema_version": "prc.fix-contract/v0.2",
        }
        del v02_contract["finding_id"]
        del v02_contract["finding_fingerprint"]
        v02_candidate = {
            **candidate,
            "schema_version": "prc.remediation-candidate/v0.2",
            "contract": v02_contract,
        }
        self.assertEqual(
            validate_instance.validation_errors(
                v02_candidate, "remediation-candidate-v0.2.schema.json"
            ),
            [],
        )

        legacy_contract = {
            **v02_contract,
            "schema_version": "prc.fix-contract/v0.1",
        }
        del legacy_contract["attempt"]
        legacy_candidate = {
            **candidate,
            "schema_version": "prc.remediation-candidate/v0.1",
            "contract": legacy_contract,
        }
        self.assertEqual(
            validate_instance.validation_errors(
                legacy_candidate, "remediation-candidate-v0.1.schema.json"
            ),
            [],
        )

    def test_minimal_remediation_run_conforms(self) -> None:
        digest = "a" * 64
        inventory = {
            "schema_version": "prc.inventory/v0.3",
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
            "schema_version": "prc.plan/v0.6",
            "digest": digest,
            "engine_version": "prc.engine/v0.1",
            "target_name": "example",
            "inventory_digest": digest,
            "profile_id": "prc/core-repository",
            "profile_version": "0.3",
            "profile_digest": digest,
            "catalog_digest": digest,
            "artifact_digests": [],
            "target_environments": [],
            "execution_mode": "inspect",
            "capability_policy": {
                "read_workspace": True,
                "write_scratch": False,
                "process": "deny",
                "network": "deny",
                "network_hosts": [],
                "secret_handles": [],
            },
            "implementations": [],
            "adapters": [],
            "nodes": [
                {
                    "id": "inventory:" + digest,
                    "kind": "inventory",
                    "depends_on": [],
                    "capabilities": {
                        "read_workspace": True,
                        "write_scratch": False,
                        "process": "deny",
                        "network": "deny",
                        "network_hosts": [],
                        "secret_handles": [],
                    },
                    "status": "ready",
                },
                {
                    "id": "gate:core-repository",
                    "kind": "gate",
                    "depends_on": [],
                    "capabilities": {
                        "read_workspace": False,
                        "write_scratch": False,
                        "process": "deny",
                        "network": "deny",
                        "network_hosts": [],
                        "secret_handles": [],
                    },
                    "status": "ready",
                },
            ],
            "assertions": [],
        }
        final_run = {
            "schema_version": "prc.run/v0.8",
            "run_id": digest,
            "started_at": "2026-08-23T12:00:00Z",
            "completed_at": "2026-08-23T12:00:01Z",
            "plan": plan,
            "inventory": inventory,
            "adapter_executions": [],
            "results": [],
            "findings": [],
            "terminal_state": "profile_satisfied",
        }
        remediation_run = {
            "schema_version": "prc.remediation-run/v0.4",
            "run_id": digest,
            "started_at": "2026-08-23T12:00:00Z",
            "completed_at": "2026-08-23T12:00:01Z",
            "profile_id": "prc/core-repository",
            "source_inventory_digest": digest,
            "candidate_root": "/tmp/candidates",
            "result_workspace": "/tmp/example",
            "final_inventory_digest": digest,
            "original_unchanged": True,
            "max_attempts": 3,
            "max_files": 20,
            "max_changed_lines": 200,
            "usage": {"attempts": 0, "changed_files": 0, "changed_lines": 0},
            "candidates": [],
            "provider_executions": [],
            "final_run": final_run,
            "gate_state": "profile_satisfied",
            "terminal_state": "profile_satisfied",
            "remaining": [],
            "stop_reasons": [],
        }
        self.assertEqual(
            validate_instance.validation_errors(
                remediation_run, "remediation-run.schema.json"
            ),
            [],
        )

        v07_run = {**final_run, "schema_version": "prc.run/v0.7"}
        self.assertEqual(
            validate_instance.validation_errors(
                v07_run, "run-result-v0.7.schema.json"
            ),
            [],
        )
        v03_remediation_run = {
            **remediation_run,
            "schema_version": "prc.remediation-run/v0.3",
            "final_run": v07_run,
        }
        self.assertEqual(
            validate_instance.validation_errors(
                v03_remediation_run, "remediation-run-v0.3.schema.json"
            ),
            [],
        )

        # Exercise the non-empty adapter branch of the current plan schema. The
        # adapter registry uses "authorized", while its executable DAG node uses
        # the distinct topological state "ready".
        verify_plan = json.loads(json.dumps(plan))
        verify_plan["execution_mode"] = "verify-local"
        verify_plan["capability_policy"]["write_scratch"] = True
        verify_plan["capability_policy"]["process"] = "oci"
        assertion_id = "PRC-A-CORE-013"
        implementation_id = "prc.native.analysis-evidence@0.1"
        adapter_id = "prc.adapter.static-analysis@0.1"
        adapter_capabilities = {
            "read_workspace": True,
            "write_scratch": True,
            "process": "oci",
            "network": "deny",
            "network_hosts": [],
            "secret_handles": [],
        }
        verify_plan["implementations"] = [{
            "id": implementation_id,
            "kind": "adapter-evidence",
            "assertion_ids": [assertion_id],
            "capabilities": {
                "read_workspace": False,
                "write_scratch": False,
                "process": "deny",
                "network": "deny",
                "network_hosts": [],
                "secret_handles": [],
            },
            "status": "available",
        }]
        verify_plan["adapters"] = [{
            "adapter_id": adapter_id,
            "manifest_sha256": digest,
            "observation_kinds": ["static-analysis"],
            "capabilities": adapter_capabilities,
            "status": "authorized",
        }]
        inventory_node = verify_plan["nodes"][0]
        adapter_node_id = "adapter:" + digest
        assertion_node_id = "assertion:" + assertion_id
        verify_plan["nodes"] = [
            inventory_node,
            {
                "id": adapter_node_id,
                "kind": "adapter",
                "depends_on": [inventory_node["id"]],
                "adapter_id": adapter_id,
                "manifest_sha256": digest,
                "capabilities": adapter_capabilities,
                "status": "ready",
            },
            {
                "id": assertion_node_id,
                "kind": "assertion",
                "depends_on": [inventory_node["id"], adapter_node_id],
                "assertion_id": assertion_id,
                "implementation_id": implementation_id,
                "capabilities": verify_plan["implementations"][0]["capabilities"],
                "status": "ready",
            },
            {
                "id": "gate:core-repository",
                "kind": "gate",
                "depends_on": [assertion_node_id],
                "capabilities": verify_plan["implementations"][0]["capabilities"],
                "status": "ready",
            },
        ]
        verify_plan["assertions"] = [{
            "assertion_id": assertion_id,
            "assertion_revision": 1,
            "definition_digest": digest,
            "implementation_id": implementation_id,
            "applicability": "applicable",
            "applicability_evaluator": "cel-go/v0.30.0+prc-inventory/v0.3",
            "applicability_reason": "CEL expression evaluated to true.",
        }]
        self.assertEqual(
            validate_instance.validation_errors(verify_plan, "plan.schema.json"),
            [],
        )

        v05_plan = {**plan, "schema_version": "prc.plan/v0.5"}
        for field in (
            "execution_mode",
            "capability_policy",
            "implementations",
            "adapters",
            "nodes",
        ):
            del v05_plan[field]
        v06_run = {
            **final_run,
            "schema_version": "prc.run/v0.6",
            "plan": v05_plan,
        }
        self.assertEqual(
            validate_instance.validation_errors(
                v06_run, "run-result-v0.6.schema.json"
            ),
            [],
        )

        legacy_plan = {**v05_plan, "schema_version": "prc.plan/v0.4"}
        del legacy_plan["catalog_digest"]
        legacy_final_run = {
            **final_run,
            "schema_version": "prc.run/v0.4",
            "plan": legacy_plan,
        }
        del legacy_final_run["findings"]
        legacy_remediation_run = {
            **remediation_run,
            "schema_version": "prc.remediation-run/v0.1",
            "final_run": legacy_final_run,
        }
        del legacy_remediation_run["provider_executions"]
        self.assertEqual(
            validate_instance.validation_errors(
                legacy_remediation_run, "remediation-run-v0.1.schema.json"
            ),
            [],
        )

        v05_run = {
            **final_run,
            "schema_version": "prc.run/v0.5",
            "plan": v05_plan,
        }
        del v05_run["findings"]

        v02_remediation_run = {
            **remediation_run,
            "schema_version": "prc.remediation-run/v0.2",
            "final_run": v05_run,
        }
        self.assertEqual(
            validate_instance.validation_errors(
                v02_remediation_run, "remediation-run-v0.2.schema.json"
            ),
            [],
        )

        self.assertEqual(
            validate_instance.validation_errors(
                v05_run, "run-result-v0.5.schema.json"
            ),
            [],
        )

    def test_doctor_report_conforms_and_warning_cannot_be_required(self) -> None:
        report = {
            "schema_version": "prc.doctor/v0.1",
            "generated_at": "2026-08-23T12:00:00Z",
            "platform": "linux",
            "architecture": "amd64",
            "target": "/tmp/project",
            "catalog_root": "/tmp/catalog",
            "ready": True,
            "summary": {"passed": 2, "warnings": 3, "failed": 0},
            "checks": [
                {
                    "id": "target.inventory",
                    "status": "pass",
                    "required": True,
                    "summary": "Target is valid.",
                    "details": [],
                },
                {
                    "id": "catalog.load",
                    "status": "pass",
                    "required": True,
                    "summary": "Catalog is valid.",
                    "details": [],
                },
                *[
                    {
                        "id": identifier,
                        "status": "warn",
                        "required": False,
                        "summary": "Capability was not requested.",
                        "details": [],
                    }
                    for identifier in (
                        "state-store",
                        "candidate-workspace",
                        "oci-runtime",
                    )
                ],
            ],
        }
        self.assertEqual(
            validate_instance.validation_errors(report, "doctor.schema.json"),
            [],
        )
        report["checks"][2]["required"] = True
        self.assertTrue(
            validate_instance.validation_errors(report, "doctor.schema.json")
        )

    def test_history_report_conforms(self) -> None:
        digest = "a" * 64
        report = {
            "schema_version": "prc.history/v0.1",
            "generated_at": "2026-08-23T12:00:00Z",
            "state_path": "/tmp/prc-state/state.sqlite",
            "runs": [
                {
                    "run_id": digest,
                    "started_at": "2026-08-23T11:59:59Z",
                    "completed_at": "2026-08-23T12:00:00Z",
                    "target_name": "example",
                    "profile_id": "prc/core-repository",
                    "profile_version": "0.3",
                    "inventory_digest": digest,
                    "terminal_state": "no_go",
                    "pass_count": 10,
                    "fail_count": 2,
                    "blocked_count": 18,
                }
            ],
        }
        self.assertEqual(
            validate_instance.validation_errors(report, "history.schema.json"),
            [],
        )

    def test_state_check_report_conforms(self) -> None:
        report = {
            "schema_version": "prc.state-check/v0.2",
            "checked_at": "2026-08-23T12:00:00Z",
            "state_path": "/tmp/prc-state/state.sqlite",
            "integrity": "ok",
            "counts": {
                "runs": 1,
                "results": 30,
                "findings": 2,
                "evidence": 20,
                "inventory_files": 100,
                "inventory_facts": 10,
                "audit_events": 1,
            },
        }
        self.assertEqual(
            validate_instance.validation_errors(
                report, "state-check.schema.json"
            ),
            [],
        )

    def test_checked_in_adapter_manifest_conforms(self) -> None:
        path = ROOT / "fixtures" / "adapters" / "fixture-adapter.yaml"
        instance = yaml.safe_load(path.read_text(encoding="utf-8"))
        self.assertEqual(
            validate_instance.validation_errors(instance, "adapter-manifest.schema.json"), []
        )

        legacy = {
            "schema_version": "prc.adapter-manifest/v0.1",
            "id": "prc.adapter.fixture@0.1",
            "title": "Legacy fixture",
            "trust": "first-party-sandboxed",
            "runner": "oci",
            "image": "registry.example/fixture@sha256:" + "a" * 64,
            "command": ["/adapter"],
            "capabilities": {
                "read_workspace": True,
                "write_scratch": True,
                "network": "deny",
                "network_hosts": [],
                "secret_handles": [],
                "child_processes": False,
            },
            "resources": {
                "timeout_seconds": 60,
                "memory_mb": 512,
                "cpus": 1,
                "pids": 64,
                "tmpfs_mb": 64,
                "max_line_bytes": 262144,
                "max_messages": 10000,
                "max_stdin_bytes": 1048576,
                "max_stdout_bytes": 4194304,
                "max_stderr_bytes": 65536,
            },
        }
        self.assertEqual(
            validate_instance.validation_errors(
                legacy, "adapter-manifest-v0.1.schema.json"
            ),
            [],
        )

    def test_checked_in_adapter_registry_conforms(self) -> None:
        path = ROOT / "fixtures" / "adapters" / "fixture-registry.yaml"
        instance = yaml.safe_load(path.read_text(encoding="utf-8"))
        self.assertEqual(
            validate_instance.validation_errors(instance, "adapter-registry.schema.json"),
            [],
        )
        report = {
            "schema_version": "prc.adapter-registry-report/v0.1",
            "registry": instance,
            "digest": "a" * 64,
        }
        self.assertEqual(
            validate_instance.validation_errors(
                report, "adapter-registry-report.schema.json"
            ),
            [],
        )

    def test_checked_in_benchmark_suite_and_report_conform(self) -> None:
        path = ROOT / "fixtures" / "benchmarks" / "core-native" / "suite.yaml"
        suite = yaml.safe_load(path.read_text(encoding="utf-8"))
        self.assertEqual(
            validate_instance.validation_errors(suite, "benchmark-suite.schema.json"),
            [],
        )
        digest = "a" * 64
        report = {
            "schema_version": "prc.benchmark-report/v0.1",
            "suite_id": suite["id"],
            "suite_digest": digest,
            "corpus_digest": "b" * 64,
            "catalog_digest": "c" * 64,
            "profile_id": suite["profile_id"],
            "evaluated_at": "2026-08-23T12:00:00Z",
            "quality_budget": suite["quality_budget"],
            "summary": {
                "cases": 1,
                "expectations": 1,
                "matched": 1,
                "mismatched": 0,
                "deterministic_cases": 1,
                "expected_outcomes": {
                    "pass": 1,
                    "fail": 0,
                    "not_applicable": 0,
                    "unknown": 0,
                    "manual_review": 0,
                    "stale": 0,
                    "conflicting": 0,
                },
            },
            "metrics": {
                "true_positive": 0,
                "false_positive": 0,
                "false_negative": 0,
                "true_negative": 1,
                "precision": 1,
                "recall": 1,
                "false_positive_rate": 0,
            },
            "cases": [{
                "id": "example",
                "target": "targets/example",
                "inventory_digest": "d" * 64,
                "run_id": "e" * 64,
                "deterministic": True,
                "passed": True,
                "expectations": [{
                    "assertion_id": "PRC-A-CORE-001",
                    "expected_assessment": "pass",
                    "expected_execution": "completed",
                    "actual_assessment": "pass",
                    "actual_execution": "completed",
                    "summary": "Observed README.md.",
                    "matched": True,
                }],
            }],
            "quality_failures": [],
            "passed": True,
        }
        self.assertEqual(
            validate_instance.validation_errors(report, "benchmark-report.schema.json"),
            [],
        )

    def test_checked_in_pack_and_report_conform(self) -> None:
        manifest = yaml.safe_load(
            (ROOT / "packs" / "core-foundation.yaml").read_text(encoding="utf-8")
        )
        self.assertEqual(
            validate_instance.validation_errors(manifest, "pack.schema.json"),
            [],
        )
        report = {
            "schema_version": "prc.pack-report/v0.1",
            "manifest": manifest,
            "digest": "a" * 64,
            "suite_digest": manifest["benchmark"]["suite_sha256"],
            "benchmark_corpus_digest": "c" * 64,
            "catalog_digest": "b" * 64,
        }
        self.assertEqual(
            validate_instance.validation_errors(report, "pack-report.schema.json"),
            [],
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
        legacy_task = {
            **task,
            "schema_version": "prc.agent-task/v0.1",
        }
        del legacy_task["finding_id"]
        del legacy_task["finding_fingerprint"]
        self.assertEqual(
            validate_instance.validation_errors(
                legacy_task, "agent-task-v0.1.schema.json"
            ),
            [],
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
            "schema_version": "prc.adapter-execution/v0.2",
            "execution_id": digest,
            "adapter_run_id": "b" * 64,
            "adapter_id": "prc.adapter.fixture@0.1",
            "manifest_sha256": "c" * 64,
            "image": "registry.example/prc/fixture@sha256:" + "d" * 64,
            "resolution": {
                "source": "registry",
                "publisher_id": "prc-project",
                "trust": "first-party-sandboxed",
                "registry_id": "prc.adapter-registry.fixtures@0.1",
                "registry_revision": 1,
                "registry_digest": "9" * 64,
            },
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
        legacy_execution = json.loads(json.dumps(execution))
        legacy_execution["schema_version"] = "prc.adapter-execution/v0.1"
        del legacy_execution["resolution"]
        self.assertEqual(
            validate_instance.validation_errors(
                legacy_execution, "adapter-execution-v0.1.schema.json"
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
