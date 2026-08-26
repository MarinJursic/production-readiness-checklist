from __future__ import annotations

import gzip
import hashlib
import io
import json
import pathlib
import tarfile
import tempfile
import unittest

from scripts import npm_distribution


def archive_files(path: pathlib.Path) -> dict[str, bytes]:
    with gzip.open(path, "rb") as compressed:
        with tarfile.open(fileobj=io.BytesIO(compressed.read()), mode="r:") as archive:
            for member in archive.getmembers():
                if not member.isfile() or member.issym() or member.islnk():
                    raise AssertionError(f"unsafe tar member {member.name}")
            return {member.name: archive.extractfile(member).read() for member in archive.getmembers()}


class NpmDistributionTests(unittest.TestCase):
    def test_packages_are_reproducible_script_free_and_exactly_bound(self) -> None:
        version = "0.2.0-test.1"
        commit = "a" * 40
        with tempfile.TemporaryDirectory() as temporary:
            root = pathlib.Path(temporary)
            binaries: dict[tuple[str, str], pathlib.Path] = {}
            for target in npm_distribution.PLATFORMS:
                suffix = ".exe" if target[0] == "windows" else ""
                path = root / f"prc-{target[0]}-{target[1]}{suffix}"
                path.write_bytes(f"binary:{target[0]}:{target[1]}".encode())
                binaries[target] = path
            support = [("catalog/example.json", b"{}\n", 0o644), ("docs/checklists/example.md", b"# Example\n", 0o644)]
            first, second = root / "first", root / "second"
            first_artifacts = npm_distribution.build_packages(
                version=version, commit=commit, built_at="2026-08-23T10:00:00Z", epoch=1787479200,
                binaries=binaries, support=support, output=first,
            )
            second_artifacts = npm_distribution.build_packages(
                version=version, commit=commit, built_at="2026-08-23T10:00:00Z", epoch=1787479200,
                binaries=binaries, support=support, output=second,
            )
            self.assertEqual(first_artifacts, second_artifacts)
            self.assertEqual(len(first_artifacts), 7)
            for artifact in first_artifacts:
                self.assertEqual((first / artifact["name"]).read_bytes(), (second / artifact["name"]).read_bytes())
                self.assertEqual(hashlib.sha256((first / artifact["name"]).read_bytes()).hexdigest(), artifact["sha256"])

            launcher_name = npm_distribution.package_filename(npm_distribution.LAUNCHER_NAME, version)
            launcher_files = archive_files(first / launcher_name)
            self.assertEqual(
                set(launcher_files),
                {"package/LICENSE", "package/README.md", "package/bin/prc.js", "package/package.json", "package/platforms.json"},
            )
            launcher_metadata = json.loads(launcher_files["package/package.json"])
            launcher = launcher_files["package/bin/prc.js"]
            self.assertNotIn("scripts", launcher_metadata)
            self.assertNotIn("dependencies", launcher_metadata)
            self.assertNotIn("private", launcher_metadata)
            self.assertEqual(launcher_metadata["publishConfig"], {"access": "public", "provenance": True})
            for forbidden in (b"shell: true", b"node:http", b"node:https", b"node:net", b"node:tls", b"fetch("):
                self.assertNotIn(forbidden, launcher)
            self.assertEqual(set(launcher_metadata["optionalDependencies"].values()), {version})
            self.assertEqual(launcher_metadata["engines"], {"node": ">=22.14.0"})
            bindings = json.loads(launcher_files["package/platforms.json"])

            platform_name = "@marinjursic/prc-linux-x64"
            platform_files = archive_files(first / npm_distribution.package_filename(platform_name, version))
            metadata = json.loads(platform_files["package/package.json"])
            manifest_bytes = platform_files["package/manifest.json"]
            manifest = json.loads(manifest_bytes)
            self.assertNotIn("scripts", metadata)
            self.assertNotIn("dependencies", metadata)
            self.assertEqual(metadata["os"], ["linux"])
            self.assertEqual(metadata["cpu"], ["x64"])
            self.assertEqual(metadata["publishConfig"], {"access": "public", "provenance": True})
            self.assertEqual(bindings["platforms"]["linux-x64"]["manifest_sha256"], hashlib.sha256(manifest_bytes).hexdigest())
            self.assertEqual(manifest["binary_sha256"], hashlib.sha256(platform_files["package/bin/prc"]).hexdigest())
            self.assertEqual(manifest["schema_version"], "prc.npm-platform/v0.2")
            self.assertEqual(manifest["support_file_count"], 2)
            self.assertEqual(manifest["support_bytes"], len(b"{}\n") + len(b"# Example\n"))
            self.assertEqual(
                manifest["support_files"],
                [
                    {"path": "bin/catalog/example.json", "size": 3, "sha256": hashlib.sha256(b"{}\n").hexdigest()},
                    {"path": "bin/docs/checklists/example.md", "size": 10, "sha256": hashlib.sha256(b"# Example\n").hexdigest()},
                ],
            )
            self.assertIn("package/bin/catalog/example.json", platform_files)
            self.assertIn("package/bin/docs/checklists/example.md", platform_files)

    def test_package_filename_never_uses_scope_syntax(self) -> None:
        self.assertEqual(
            npm_distribution.package_filename("@marinjursic/prc-linux-x64", "1.2.3"),
            "marinjursic-prc-linux-x64-1.2.3.tgz",
        )


if __name__ == "__main__":
    unittest.main()
