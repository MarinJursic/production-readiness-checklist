from __future__ import annotations

import gzip
import hashlib
import pathlib
import re
import unittest

import yaml


ROOT = pathlib.Path(__file__).resolve().parents[1]
HISTORICAL_GITLEAKS_CONFIG = "scanner/adapter/gitleaks-v8.30.0.toml"


class SecurityConfigurationTests(unittest.TestCase):
    def test_python_workflow_installs_require_hashes_and_binary_wheels(self) -> None:
        workflow_paths = (
            ROOT / ".github" / "workflows" / "pages.yml",
            ROOT / ".github" / "workflows" / "validate.yml",
        )
        install_count = 0
        for path in workflow_paths:
            source = path.read_text(encoding="utf-8")
            for match in re.finditer(r"(?:python -m )?pip install", source):
                install_count += 1
                command = source[match.start() : match.start() + 240]
                self.assertIn("--require-hashes", command, path)
                self.assertIn("--only-binary=:all:", command, path)
        self.assertEqual(install_count, 6)

    def test_python_lock_files_pin_every_artifact_by_hash(self) -> None:
        lock_names = (
            "requirements-dev.lock.txt",
            "requirements-docs.lock.txt",
            "requirements-audit.lock.txt",
        )
        record_pattern = re.compile(r"(?m)^[A-Za-z0-9_.-]+==")
        for name in lock_names:
            source = (ROOT / name).read_text(encoding="utf-8")
            starts = [match.start() for match in record_pattern.finditer(source)]
            self.assertTrue(starts, name)
            for index, start in enumerate(starts):
                end = starts[index + 1] if index + 1 < len(starts) else len(source)
                self.assertIn("--hash=sha256:", source[start:end], f"{name}:{index + 1}")

        audit_input = (ROOT / "requirements-audit.txt").read_text(encoding="utf-8")
        self.assertRegex(audit_input, r"(?m)^pip-audit==[0-9]+(?:\.[0-9]+)+$")

    def test_vulnerability_reporting_and_scorecard_context_are_explicit(self) -> None:
        security_policy = (ROOT / "SECURITY.md").read_text(encoding="utf-8")
        self.assertIn(
            "https://github.com/MarinJursic/production-readiness-checklist/security/advisories/new",
            security_policy,
        )

        scorecard = yaml.safe_load((ROOT / ".scorecard.yml").read_text(encoding="utf-8"))
        reasons = {
            check: annotation["reasons"][0]["reason"]
            for annotation in scorecard["annotations"]
            for check in annotation["checks"]
        }
        self.assertEqual(
            reasons,
            {
                "pinned-dependencies": "remediated",
                "branch-protection": "remediated",
                "code-review": "not-applicable",
            },
        )

    def test_secret_scanning_excludes_only_hash_bound_historical_upstream_rules(self) -> None:
        configuration = yaml.safe_load(
            (ROOT / ".github" / "secret_scanning.yml").read_text(encoding="utf-8")
        )
        self.assertEqual(configuration, {"paths-ignore": [HISTORICAL_GITLEAKS_CONFIG]})
        self.assertFalse((ROOT / HISTORICAL_GITLEAKS_CONFIG).exists())

        archive = ROOT / f"{HISTORICAL_GITLEAKS_CONFIG}.gz"
        archive_data = archive.read_bytes()
        source = (ROOT / "scanner" / "adapter" / "gitleaks.go").read_text(encoding="utf-8")
        expanded_digest = self._constant(source, "GitleaksConfigSHA256")
        archive_digest = self._constant(source, "GitleaksConfigArchiveSHA256")

        self.assertEqual(hashlib.sha256(archive_data).hexdigest(), archive_digest)
        self.assertEqual(hashlib.sha256(gzip.decompress(archive_data)).hexdigest(), expanded_digest)

    def test_repository_security_gates_keep_their_required_depth(self) -> None:
        codeql = (ROOT / ".github" / "workflows" / "codeql.yml").read_text(
            encoding="utf-8"
        )
        self.assertIn("queries: security-extended", codeql)
        for language in ("go", "python", "actions"):
            self.assertIn(f"language: {language}", codeql)

        secret_scan = (ROOT / ".github" / "workflows" / "secret-scan.yml").read_text(
            encoding="utf-8"
        )
        self.assertIn("fetch-depth: 0", secret_scan)
        self.assertIn("git --no-banner --redact .", secret_scan)

        dependabot = yaml.safe_load(
            (ROOT / ".github" / "dependabot.yml").read_text(encoding="utf-8")
        )
        ecosystems = {update["package-ecosystem"] for update in dependabot["updates"]}
        self.assertEqual(ecosystems, {"github-actions", "pip", "gomod"})

    @staticmethod
    def _constant(source: str, name: str) -> str:
        match = re.search(rf'(?m)^\s*{name}\s*=\s*"([0-9a-f]{{64}})"$', source)
        if match is None:
            raise AssertionError(f"missing {name}")
        return match.group(1)


if __name__ == "__main__":
    unittest.main()
