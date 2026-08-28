from __future__ import annotations

import hashlib
import io
import pathlib
import tarfile
import tempfile
import unittest
from unittest import mock

from scripts import smoke_release


class SmokeReleaseTests(unittest.TestCase):
    def test_supported_host_mapping_is_exact(self) -> None:
        with mock.patch("platform.system", return_value="Windows"), mock.patch("platform.machine", return_value="ARM64"):
            self.assertEqual(smoke_release.host(), ("windows", "arm64", "win32", "arm64"))
        with mock.patch("platform.system", return_value="Plan9"), mock.patch("platform.machine", return_value="mips"):
            with self.assertRaisesRegex(RuntimeError, "unsupported"):
                smoke_release.host()

    def test_tar_extraction_rejects_links_and_parent_paths(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            root = pathlib.Path(temporary)
            for name, kind in (("../escape", tarfile.REGTYPE), ("root/link", tarfile.SYMTYPE)):
                archive = root / ("bad-" + kind.decode("ascii") + ".tar.gz")
                with tarfile.open(archive, "w:gz") as output:
                    member = tarfile.TarInfo(name)
                    member.type = kind
                    member.size = 1 if kind == tarfile.REGTYPE else 0
                    output.addfile(member, io.BytesIO(b"x") if member.size else None)
                with self.assertRaisesRegex(ValueError, "archive"):
                    smoke_release.extract_archive(archive, root / "out")

    def test_complete_report_contract_is_required(self) -> None:
        valid = {
            "schema_version": "prc.run/v0.12",
            "results": [{}],
            "control_results": [{}] * 10_042,
            "control_catalog": {"control_count": 10_042},
        }
        smoke_release.verify_result(valid, "test")
        valid["control_results"] = valid["control_results"][:-1]
        with self.assertRaisesRegex(RuntimeError, "incomplete"):
            smoke_release.verify_result(valid, "test")

    def test_json_process_output_is_always_strict_utf8(self) -> None:
        completed = mock.Mock(returncode=0, stdout='{"message":"ready — yes"}', stderr="")
        with mock.patch.object(smoke_release.subprocess, "run", return_value=completed) as run:
            self.assertEqual(smoke_release.run_json(["scanner"], "test"), {"message": "ready — yes"})
        self.assertEqual(run.call_args.kwargs["encoding"], "utf-8")
        self.assertEqual(run.call_args.kwargs["errors"], "strict")
        self.assertNotIn("text", run.call_args.kwargs)

    def test_global_command_uses_npm_exposed_shim(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            root = pathlib.Path(temporary)
            (root / "bin").mkdir()
            posix = root / "bin" / "prc"
            posix.write_text("#!/bin/sh\n", encoding="utf-8")
            self.assertEqual(
                smoke_release.global_prc_command(root, ["version"], windows=False),
                [str(posix), "version"],
            )
            windows = root / "prc.cmd"
            windows.write_text("@echo off\r\n", encoding="utf-8")
            command = smoke_release.global_prc_command(
                root, ["scan", "space path"], windows=True, command_processor="/Windows/System32/cmd.exe",
            )
            self.assertEqual(command[:4], ["/Windows/System32/cmd.exe", "/d", "/s", "/c"])
            self.assertIn("space path", command[4])

    def test_global_command_fails_when_npm_did_not_expose_it(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            with self.assertRaisesRegex(RuntimeError, "did not expose"):
                smoke_release.global_prc_command(pathlib.Path(temporary), [], windows=False)

    def test_manifest_cannot_authorize_an_oversized_artifact(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            root = pathlib.Path(temporary)
            artifact = root / "artifact.tgz"
            artifact.write_bytes(b"small")
            item = {
                "name": artifact.name,
                "size": smoke_release.MAXIMUM_ARTIFACT_BYTES + 1,
                "sha256": hashlib.sha256(b"small").hexdigest(),
            }
            with self.assertRaisesRegex(ValueError, "unsafe"):
                smoke_release.checked_artifact(root, item)


if __name__ == "__main__":
    unittest.main()
