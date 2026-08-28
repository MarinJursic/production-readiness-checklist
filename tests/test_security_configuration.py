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
