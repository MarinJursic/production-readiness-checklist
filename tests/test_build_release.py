from __future__ import annotations

import datetime as dt
import gzip
import json
import pathlib
import re
import tarfile
import tempfile
import unittest
import zipfile

from scripts import build_release


class ReleaseBuilderTests(unittest.TestCase):
    def test_release_python_lock_matches_dev_versions_and_hashes_every_wheel(self) -> None:
        development: dict[str, str] = {}
        for line in (build_release.ROOT / "requirements-dev.lock.txt").read_text(encoding="utf-8").splitlines():
            if line and not line.startswith("#"):
                name, version = line.split("==", 1)
                development[name.lower()] = version
        release_text = (build_release.ROOT / "requirements-release.lock.txt").read_text(encoding="utf-8")
        records = re.findall(
            r"(?m)^([A-Za-z0-9_.-]+)==([^ \\]+) \\\n+\s+--hash=sha256:([0-9a-f]{64})$",
            release_text,
        )
        self.assertEqual({name.lower(): version for name, version, _digest in records}, development)
        self.assertEqual(len({digest for _name, _version, digest in records}), len(development))

    def test_release_version_uses_strict_semver_without_build_metadata(self) -> None:
        for value in ("0.1.0", "1.2.3-rc.1", "10.0.0-alpha-beta"):
            self.assertIsNotNone(build_release.SEMVER.fullmatch(value), value)
        for value in ("v1.2.3", "01.2.3", "1.02.3", "1.2.3-01", "1.2.3-alpha.", "1.2.3+build"):
            self.assertIsNone(build_release.SEMVER.fullmatch(value), value)

    def test_timestamp_is_timezone_required_and_normalized(self) -> None:
        normalized, epoch = build_release.parse_timestamp("2026-08-23T12:34:56+02:00")
        self.assertEqual(normalized, "2026-08-23T10:34:56Z")
        self.assertEqual(epoch, int(dt.datetime(2026, 8, 23, 10, 34, 56, tzinfo=dt.timezone.utc).timestamp()))
        with self.assertRaisesRegex(ValueError, "timezone"):
            build_release.parse_timestamp("2026-08-23T10:34:56")

    def test_sbom_normalization_binds_release_identity_deterministically(self) -> None:
        old_reference = f"pkg:golang/{build_release.MODULE}@v0.0.0?type=module"
        source_document = {
            "bomFormat": "CycloneDX",
            "specVersion": "1.6",
            "version": 1,
            "metadata": {
                "tools": [{"name": "cyclonedx-gomod"}],
                "component": {
                    "bom-ref": old_reference,
                    "type": "application",
                    "name": build_release.MODULE,
                    "version": "v0.0.0",
                    "purl": old_reference,
                },
            },
            "dependencies": [{"ref": old_reference, "dependsOn": []}],
        }
        with tempfile.TemporaryDirectory() as temporary:
            source = pathlib.Path(temporary) / "bom.json"
            source.write_text(json.dumps(source_document), encoding="utf-8")
            first = build_release.normalized_sbom(source, "0.1.0", "a" * 40)
            second = build_release.normalized_sbom(source, "0.1.0", "a" * 40)
        self.assertEqual(first, second)
        normalized = json.loads(first)
        reference = f"pkg:golang/{build_release.MODULE}@0.1.0?type=module"
        self.assertEqual(normalized["metadata"]["component"]["bom-ref"], reference)
        self.assertEqual(normalized["dependencies"][0]["ref"], reference)
        self.assertEqual(
            normalized["serialNumber"],
            build_release.release_sbom_serial("0.1.0", "a" * 40),
        )
        self.assertRegex(
            normalized["serialNumber"],
            r"^urn:uuid:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$",
        )
        self.assertNotEqual(
            normalized["serialNumber"],
            build_release.release_sbom_serial("0.1.1", "a" * 40),
        )
        self.assertNotEqual(
            normalized["serialNumber"],
            build_release.release_sbom_serial("0.1.0", "b" * 40),
        )
        self.assertEqual(
            normalized["metadata"]["properties"],
            [
                {"name": "prc:build:version", "value": "0.1.0"},
                {"name": "prc:source:commit", "value": "a" * 40},
            ],
        )

    def test_sbom_rejects_nondeterministic_metadata(self) -> None:
        document = {
            "bomFormat": "CycloneDX",
            "specVersion": "1.6",
            "serialNumber": "urn:uuid:unsafe",
            "metadata": {"component": {"name": build_release.MODULE}},
        }
        with tempfile.TemporaryDirectory() as temporary:
            source = pathlib.Path(temporary) / "bom.json"
            source.write_text(json.dumps(document), encoding="utf-8")
            with self.assertRaisesRegex(ValueError, "builder owns"):
                build_release.normalized_sbom(source, "0.1.0", "b" * 40)

    def test_sbom_rejects_predeclared_scanner_identity(self) -> None:
        document = {
            "bomFormat": "CycloneDX",
            "specVersion": "1.6",
            "metadata": {
                "component": {"name": build_release.MODULE},
                "properties": [{"name": "prc:source:commit", "value": "spoofed"}],
            },
        }
        with tempfile.TemporaryDirectory() as temporary:
            source = pathlib.Path(temporary) / "bom.json"
            source.write_text(json.dumps(document), encoding="utf-8")
            with self.assertRaisesRegex(ValueError, "scanner-owned"):
                build_release.normalized_sbom(source, "0.1.0", "b" * 40)

    def test_archives_are_reproducible_and_rooted(self) -> None:
        entries = [("prc", b"binary", 0o755), ("LICENSE", b"license\n", 0o644)]
        timestamp = dt.datetime(2026, 8, 23, 10, 0, tzinfo=dt.timezone.utc)
        epoch = int(timestamp.timestamp())
        with tempfile.TemporaryDirectory() as temporary:
            root = pathlib.Path(temporary)
            first_tar = root / "first.tar.gz"
            second_tar = root / "second.tar.gz"
            build_release.create_tar_gz(first_tar, "prc_0.1.0_linux_amd64", entries, epoch)
            build_release.create_tar_gz(second_tar, "prc_0.1.0_linux_amd64", entries, epoch)
            self.assertEqual(first_tar.read_bytes(), second_tar.read_bytes())
            with gzip.open(first_tar, "rb") as compressed:
                with tarfile.open(fileobj=compressed, mode="r:") as archive:
                    self.assertEqual(
                        archive.getnames(),
                        ["prc_0.1.0_linux_amd64/LICENSE", "prc_0.1.0_linux_amd64/prc"],
                    )
                    self.assertEqual(archive.getmember("prc_0.1.0_linux_amd64/prc").mode, 0o755)

            first_zip = root / "first.zip"
            second_zip = root / "second.zip"
            build_release.create_zip(first_zip, "prc_0.1.0_windows_amd64", entries, timestamp)
            build_release.create_zip(second_zip, "prc_0.1.0_windows_amd64", entries, timestamp)
            self.assertEqual(first_zip.read_bytes(), second_zip.read_bytes())
            with zipfile.ZipFile(first_zip) as archive:
                self.assertEqual(
                    archive.namelist(),
                    ["prc_0.1.0_windows_amd64/LICENSE", "prc_0.1.0_windows_amd64/prc"],
                )

    def test_release_support_contains_runtime_catalog_and_user_guides(self) -> None:
        names = {name for name, _data, _mode in build_release.release_support_files()}
        for expected in (
            "DISCLOSURE",
            "THIRD_PARTY_NOTICES.md",
            "adapters/checkov-v3.3.8.yaml",
            "catalog/profiles/core-repository.yaml",
            "docs/scanner/getting-started.md",
        ):
            self.assertIn(expected, names)

    def test_npm_support_keeps_every_control_source_without_website_media(self) -> None:
        support = build_release.npm_support_files()
        names = {name for name, _data, _mode in support}
        for expected in (
            "THIRD_PARTY_NOTICES.md",
            "adapters/checkov-v3.3.8.yaml",
            "catalog/control-contracts.json.gz",
            "catalog/control-id-registry.json.gz",
            "catalog/profiles/core-repository.yaml",
            "docs/checklists/00-readiness-principle.md",
            "docs/engineering/16-specialized-domains-and-release-assurance.md",
            "fixtures/benchmarks/core-native/suite.yaml",
            "packs/core-native.yaml",
            "schemas/control-review-output.schema.json",
        ):
            self.assertIn(expected, names)
        self.assertFalse(any(name.startswith("docs/assets/") for name in names))
        self.assertFalse(any(name.startswith("docs/scanner/") for name in names))
        self.assertNotIn("catalog/control-contracts.json", names)
        self.assertNotIn("catalog/control-id-registry.json", names)
        compact = dict((name, data) for name, data, _mode in support)
        self.assertEqual(
            gzip.decompress(compact["catalog/control-contracts.json.gz"]),
            (build_release.ROOT / "catalog" / "control-contracts.json").read_bytes(),
        )
        self.assertEqual(
            gzip.decompress(compact["catalog/control-id-registry.json.gz"]),
            (build_release.ROOT / "catalog" / "control-id-registry.json").read_bytes(),
        )
        self.assertLessEqual(len(support), build_release.MAX_NPM_SUPPORT_FILES)
        self.assertLessEqual(sum(len(data) for _name, data, _mode in support), build_release.MAX_NPM_SUPPORT_BYTES)

    def test_current_release_manifest_schema_includes_exact_npm_packages(self) -> None:
        schema = json.loads((build_release.ROOT / "schemas" / "release-manifest-v0.3.schema.json").read_text())
        self.assertEqual(schema["properties"]["schema_version"]["const"], "prc.release-manifest/v0.3")
        archives = schema["properties"]["artifacts"]
        self.assertTrue(archives["uniqueItems"])
        self.assertEqual(len(archives["allOf"]), 6)
        self.assertEqual(len(archives["items"]["oneOf"]), 6)
        npm_packages = schema["properties"]["npm_packages"]
        self.assertEqual((npm_packages["minItems"], npm_packages["maxItems"]), (7, 7))
        self.assertTrue(npm_packages["uniqueItems"])
        self.assertEqual(len(npm_packages["allOf"]), 7)
        self.assertEqual(len(npm_packages["items"]["oneOf"]), 7)

    def test_checksums_are_sorted_and_reject_nonfiles(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            root = pathlib.Path(temporary)
            (root / "z.txt").write_bytes(b"z")
            (root / "a.txt").write_bytes(b"a")
            build_release.write_checksums(root)
            names = [line.split("  ", 1)[1] for line in (root / "SHA256SUMS").read_text().splitlines()]
            self.assertEqual(names, ["a.txt", "z.txt"])
            (root / "directory").mkdir()
            with self.assertRaisesRegex(ValueError, "not a regular file"):
                build_release.write_checksums(root)


if __name__ == "__main__":
    unittest.main()
