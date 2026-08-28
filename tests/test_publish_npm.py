from __future__ import annotations

import json
import os
import pathlib
import subprocess
import tempfile
import unittest
from unittest import mock

import yaml

from scripts import npm_distribution, publish_npm


def completed(command: list[str], code: int = 0, stdout: str = "", stderr: str = "") -> subprocess.CompletedProcess[str]:
    return subprocess.CompletedProcess(command, code, stdout, stderr)


class PublishNpmTests(unittest.TestCase):
    commit = "a" * 40

    def trusted_environment(self, version: str = "1.2.3") -> dict[str, str]:
        tag = f"scanner-v{version}"
        return {
            "GITHUB_ACTIONS": "true",
            "GITHUB_REPOSITORY": publish_npm.TRUSTED_REPOSITORY,
            "GITHUB_SHA": self.commit,
            "GITHUB_REF": f"refs/tags/{tag}",
            "GITHUB_REF_TYPE": "tag",
            "GITHUB_REF_NAME": tag,
            "GITHUB_WORKFLOW_REF": (
                f"{publish_npm.TRUSTED_REPOSITORY}/{publish_npm.TRUSTED_WORKFLOW}@refs/tags/{tag}"
            ),
            "ACTIONS_ID_TOKEN_REQUEST_URL": "https://oidc.invalid",
            "ACTIONS_ID_TOKEN_REQUEST_TOKEN": "short-lived",
        }

    def packages(self) -> list[publish_npm.Package]:
        root = pathlib.Path("/release")
        return [
            publish_npm.Package(
                "@marinjursic/prc-linux-x64", "1.2.3", root / "platform.tgz",
                "sha512-platform", "platform", self.commit,
            ),
            publish_npm.Package(
                "@marinjursic/prc", "1.2.3", root / "launcher.tgz",
                "sha512-launcher", "launcher", self.commit,
            ),
        ]

    def test_stage_is_platform_first_idempotent_and_binds_registry_bytes(self) -> None:
        calls: list[list[str]] = []
        lookups = {
            "@marinjursic/prc-linux-x64@1.2.3": [None],
            "@marinjursic/prc@1.2.3": ["sha512-launcher"],
        }

        def run(command: list[str]) -> subprocess.CompletedProcess[str]:
            calls.append(command)
            if command == ["node", "--version"]:
                return completed(command, stdout="v22.14.0\n")
            if command == ["npm", "--version"]:
                return completed(command, stdout="12.0.2\n")
            if command[:2] == ["npm", "view"]:
                value = lookups[command[2]].pop(0)
                if value is None:
                    return completed(command, 1, stderr=json.dumps({"error": {"code": "E404"}}))
                return completed(command, stdout=json.dumps(value))
            if command[:3] == ["npm", "stage", "publish"]:
                return completed(command, stdout=json.dumps({
                    "@marinjursic/prc-linux-x64": {
                        "id": "@marinjursic/prc-linux-x64@1.2.3",
                        "name": "@marinjursic/prc-linux-x64",
                        "version": "1.2.3",
                        "integrity": "sha512-platform",
                        "stageId": "stage-platform",
                    },
                }))
            raise AssertionError(command)

        publish_npm.stage(self.packages(), run, self.trusted_environment())
        stages = [command for command in calls if command[:3] == ["npm", "stage", "publish"]]
        self.assertEqual(len(stages), 1)
        self.assertEqual(stages[0][3], "/release/platform.tgz")
        self.assertIn("--json", stages[0])
        self.assertIn("--provenance", stages[0])
        self.assertIn("--ignore-scripts", stages[0])

    def test_default_runner_ignores_saved_npm_configuration(self) -> None:
        def run(command: list[str], **kwargs: object) -> subprocess.CompletedProcess[str]:
            self.assertTrue(pathlib.Path(str(kwargs["cwd"])).is_dir())
            environment = kwargs["env"]
            self.assertIsInstance(environment, dict)
            assert isinstance(environment, dict)
            user_config = pathlib.Path(environment["NPM_CONFIG_USERCONFIG"])
            global_config = pathlib.Path(environment["NPM_CONFIG_GLOBALCONFIG"])
            npm_cache = pathlib.Path(environment["NPM_CONFIG_CACHE"])
            self.assertNotEqual(user_config, global_config)
            self.assertTrue(user_config.is_file())
            self.assertTrue(global_config.is_file())
            self.assertEqual(npm_cache.parent, user_config.parent)
            self.assertNotIn("NPM_CONFIG__AUTH", environment)
            self.assertNotIn("GH_TOKEN", environment)
            self.assertNotIn("GITHUB_TOKEN", environment)
            self.assertNotIn("NODE_AUTH_TOKEN", environment)
            self.assertNotIn("NPM_TOKEN", environment)
            self.assertEqual(environment["ACTIONS_ID_TOKEN_REQUEST_URL"], "https://oidc.invalid")
            return completed(command, stdout="11.5.1\n")

        with mock.patch.dict(os.environ, {
            "NPM_CONFIG__AUTH": "saved-token",
            "GH_TOKEN": "github-cli-token",
            "GITHUB_TOKEN": "actions-token",
            "NODE_AUTH_TOKEN": "legacy-node-token",
            "NPM_TOKEN": "legacy-npm-token",
            "ACTIONS_ID_TOKEN_REQUEST_URL": "https://oidc.invalid",
        }), mock.patch("scripts.publish_npm.subprocess.run", side_effect=run):
            result = publish_npm.default_run(["npm", "--version"])
        self.assertEqual(result.returncode, 0)

    def test_release_workflow_does_not_enable_legacy_npm_token_auth(self) -> None:
        workflow = yaml.safe_load(
            (publish_npm.build_release.ROOT / ".github/workflows/release-scanner.yml").read_text(
                encoding="utf-8"
            )
        )
        setup = next(
            step
            for step in workflow["jobs"]["publish"]["steps"]
            if step.get("name") == "Set up Node.js for npm trusted publishing"
        )
        self.assertNotIn("registry-url", setup.get("with", {}))
        install = next(
            step
            for step in workflow["jobs"]["publish"]["steps"]
            if step.get("name") == "Install hash-bound npm trusted-publishing client"
        )
        self.assertEqual(install["env"]["PRC_NPM_VERSION"], "12.0.2")
        self.assertRegex(install["env"]["PRC_NPM_INTEGRITY"], r"^sha512-[A-Za-z0-9+/]+={0,2}$")
        self.assertIn('npm install --global --ignore-scripts', install["run"])
        workflow_text = (publish_npm.build_release.ROOT / ".github/workflows/release-scanner.yml").read_text(
            encoding="utf-8"
        )
        self.assertIn("Stage npm packages with trusted publishing", workflow_text)
        self.assertNotIn("npm publish", workflow_text)

    def test_existing_different_package_bytes_fail_closed(self) -> None:
        def run(command: list[str]) -> subprocess.CompletedProcess[str]:
            if command == ["node", "--version"]:
                return completed(command, stdout="v24.1.0")
            if command == ["npm", "--version"]:
                return completed(command, stdout="12.0.2")
            return completed(command, stdout=json.dumps("sha512-other"))

        with self.assertRaisesRegex(RuntimeError, "different bytes"):
            publish_npm.stage(self.packages(), run, self.trusted_environment())

    def test_token_or_old_toolchain_is_rejected_before_stage(self) -> None:
        token_environment = self.trusted_environment() | {"NPM_TOKEN": "secret"}
        with self.assertRaisesRegex(RuntimeError, "token-free"):
            publish_npm.stage(self.packages(), lambda command: completed(command), token_environment)
        def old_run(command: list[str]) -> subprocess.CompletedProcess[str]:
            return completed(command, stdout="v22.13.0" if command[0] == "node" else "12.0.2")
        with self.assertRaisesRegex(RuntimeError, "requires Node"):
            publish_npm.stage(self.packages(), old_run, self.trusted_environment())

        def old_npm_run(command: list[str]) -> subprocess.CompletedProcess[str]:
            return completed(command, stdout="v24.19.0" if command[0] == "node" else "11.17.0")
        with self.assertRaisesRegex(RuntimeError, "npm >=12.0.2"):
            publish_npm.stage(self.packages(), old_npm_run, self.trusted_environment())

    def test_public_verification_waits_for_publish_time_scanning(self) -> None:
        integrity_calls = 0

        def run(command: list[str]) -> subprocess.CompletedProcess[str]:
            nonlocal integrity_calls
            if command[:2] == ["npm", "view"]:
                self.assertIn("--prefer-online", command)
                if command[3] == "dist.integrity":
                    integrity_calls += 1
                    if integrity_calls < 3:
                        return completed(command, 1, stderr=json.dumps({"error": {"code": "E404"}}))
                    return completed(command, stdout=json.dumps("sha512-platform"))
                if command[3] == "dist.attestations":
                    return completed(command, stdout=json.dumps({
                        "url": (
                            "https://registry.npmjs.org/-/npm/v1/attestations/"
                            "%40marinjursic%2Fprc-linux-x64%401.2.3"
                        ),
                        "provenance": {"predicateType": "https://slsa.dev/provenance/v1"},
                    }))
            raise AssertionError(command)

        with mock.patch("scripts.publish_npm.time.sleep") as sleep:
            publish_npm.verify_public(self.packages()[:1], run)
        self.assertEqual(sleep.call_args_list, [mock.call(15), mock.call(30)])

    def test_public_verification_rejects_wrong_integrity(self) -> None:
        def run(command: list[str]) -> subprocess.CompletedProcess[str]:
            if command[:2] == ["npm", "view"]:
                return completed(command, stdout=json.dumps("sha512-wrong"))
            raise AssertionError(command)

        with self.assertRaisesRegex(RuntimeError, "different bytes"):
            publish_npm.verify_public(self.packages()[:1], run)

    def test_public_verification_rejects_unbound_provenance(self) -> None:
        def run(command: list[str]) -> subprocess.CompletedProcess[str]:
            if command[:2] == ["npm", "view"] and command[3] == "dist.integrity":
                return completed(command, stdout=json.dumps("sha512-platform"))
            if command[:2] == ["npm", "view"] and command[3] == "dist.attestations":
                return completed(command, stdout=json.dumps({
                    "url": "https://registry.npmjs.org/-/npm/v1/attestations/%40other%2Fpackage%401.2.3",
                    "provenance": {"predicateType": "https://slsa.dev/provenance/v1"},
                }))
            raise AssertionError(command)

        with self.assertRaisesRegex(RuntimeError, "unbound provenance"):
            publish_npm.verify_public(self.packages()[:1], run)

    def test_missing_provenance_is_not_mistaken_for_verified(self) -> None:
        package = self.packages()[0]
        for output in ("", "null\n", "{}\n"):
            with self.subTest(output=output):
                self.assertFalse(publish_npm.registry_has_provenance(
                    lambda command, value=output: completed(command, stdout=value), package
                ))

    def test_finalize_workflow_verifies_public_bytes_before_publishing_release(self) -> None:
        workflow_path = publish_npm.build_release.ROOT / ".github/workflows/finalize-scanner-release.yml"
        workflow = yaml.safe_load(workflow_path.read_text(encoding="utf-8"))
        job = workflow["jobs"]["finalize"]
        self.assertEqual(job["permissions"], {"contents": "write"})
        workflow_text = workflow_path.read_text(encoding="utf-8")
        self.assertIn("--verify-public", workflow_text)
        self.assertIn("sha256sum --check SHA256SUMS", workflow_text)
        self.assertIn("gh release edit", workflow_text)
        self.assertLess(workflow_text.index("--verify-public"), workflow_text.index("gh release edit"))

    def test_stage_rejects_local_or_wrong_workflow_identity(self) -> None:
        with self.assertRaisesRegex(RuntimeError, "exact trusted tagged release workflow"):
            publish_npm.stage(self.packages(), environment={})
        wrong = self.trusted_environment() | {"GITHUB_SHA": "b" * 40}
        with self.assertRaisesRegex(RuntimeError, "exact trusted tagged release workflow"):
            publish_npm.stage(self.packages(), environment=wrong)

    def test_stage_rejects_unbound_json_response(self) -> None:
        package = self.packages()[0]
        with self.assertRaisesRegex(RuntimeError, "exact package identity"):
            publish_npm.staged_package_id(json.dumps({package.name: {
                "id": f"{package.name}@{package.version}",
                "name": package.name,
                "version": package.version,
                "integrity": package.integrity,
            }}), package)

    def test_load_packages_binds_manifest_tarball_identity_and_launcher_last(self) -> None:
        version = "1.2.3-test.1"
        commit = "a" * 40
        with tempfile.TemporaryDirectory() as temporary:
            root = pathlib.Path(temporary)
            binaries: dict[tuple[str, str], pathlib.Path] = {}
            for target in npm_distribution.PLATFORMS:
                binary = root / f"{target[0]}-{target[1]}"
                binary.write_bytes(b"binary-" + repr(target).encode())
                binaries[target] = binary
            output = root / "release"
            artifacts = npm_distribution.build_packages(
                version=version,
                commit=commit,
                built_at="2026-08-24T10:00:00Z",
                epoch=1787565600,
                binaries=binaries,
                support=[("catalog/example.json", b"{}\n", 0o644)],
                output=output,
            )
            manifest = output / f"prc_{version}_release-manifest.json"
            manifest.write_text(json.dumps({
                "schema_version": "prc.release-manifest/v0.3",
                "product": "prc-scanner",
                "version": version,
                "source_commit": commit,
                "npm_packages": artifacts,
            }), encoding="utf-8")
            packages = publish_npm.load_packages(output, manifest, version)
            self.assertEqual(len(packages), 7)
            self.assertEqual(packages[-1].name, "@marinjursic/prc")
            self.assertTrue(all(package.integrity.startswith("sha512-") for package in packages))
            oversized = json.loads(manifest.read_text(encoding="utf-8"))
            oversized["npm_packages"][0]["size"] = publish_npm.MAXIMUM_NPM_PACKAGE_BYTES + 1
            manifest.write_text(json.dumps(oversized), encoding="utf-8")
            with self.assertRaisesRegex(ValueError, "inconsistent"):
                publish_npm.load_packages(output, manifest, version)
            oversized["npm_packages"][0]["size"] = artifacts[0]["size"]
            manifest.write_text(json.dumps(oversized), encoding="utf-8")
            tampered = packages[0].path
            tampered.write_bytes(tampered.read_bytes() + b"tampered")
            with self.assertRaisesRegex(ValueError, "does not match"):
                publish_npm.load_packages(output, manifest, version)


if __name__ == "__main__":
    unittest.main()
