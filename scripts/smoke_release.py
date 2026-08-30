#!/usr/bin/env python3
"""Run one release's native archive and npm launcher on the current host."""

from __future__ import annotations

import argparse
import hashlib
import json
import os
import pathlib
import platform
import shutil
import subprocess
import sys
import tarfile
import tempfile
import zipfile
from typing import Any

try:
    from scripts import build_release, npm_distribution
except ModuleNotFoundError:  # Direct `python scripts/smoke_release.py` execution.
    import build_release
    import npm_distribution


HOSTS = {
    ("Linux", "x86_64"): ("linux", "amd64", "linux", "x64"),
    ("Linux", "aarch64"): ("linux", "arm64", "linux", "arm64"),
    ("Darwin", "x86_64"): ("darwin", "amd64", "darwin", "x64"),
    ("Darwin", "arm64"): ("darwin", "arm64", "darwin", "arm64"),
    ("Windows", "AMD64"): ("windows", "amd64", "win32", "x64"),
    ("Windows", "ARM64"): ("windows", "arm64", "win32", "arm64"),
}

MAXIMUM_ARTIFACT_BYTES = 512 * 1024 * 1024
MAXIMUM_EXPANDED_ARCHIVE_BYTES = 1024 * 1024 * 1024


def host() -> tuple[str, str, str, str]:
    result = HOSTS.get((platform.system(), platform.machine()))
    if result is None:
        raise RuntimeError(f"unsupported smoke-test host: {platform.system()}/{platform.machine()}")
    return result


def load_json(path: pathlib.Path, label: str) -> dict[str, Any]:
    if path.is_symlink() or not path.is_file() or path.stat().st_size > 32 * 1024 * 1024:
        raise ValueError(f"{label} must be a bounded regular file")
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
    except (OSError, UnicodeDecodeError, json.JSONDecodeError) as error:
        raise ValueError(f"{label} must contain valid UTF-8 JSON") from error
    if not isinstance(value, dict):
        raise ValueError(f"{label} must contain a JSON object")
    return value


def safe_destination(root: pathlib.Path, name: str) -> pathlib.Path:
    pure = pathlib.PurePosixPath(name.replace("\\", "/"))
    if pure.is_absolute() or not pure.parts or ".." in pure.parts:
        raise ValueError(f"unsafe release archive path: {name}")
    destination = root.joinpath(*pure.parts).resolve()
    if root.resolve() not in destination.parents:
        raise ValueError(f"release archive path escapes extraction root: {name}")
    return destination


def extract_archive(archive: pathlib.Path, destination: pathlib.Path) -> None:
    seen: set[str] = set()
    expanded = 0
    if archive.name.endswith(".tar.gz"):
        with tarfile.open(archive, mode="r:gz") as source:
            for member in source.getmembers():
                if not member.isfile() or member.issym() or member.islnk() or member.name in seen:
                    raise ValueError(f"release archive contains a non-file or duplicate entry: {member.name}")
                expanded += member.size
                if member.size < 0 or expanded > MAXIMUM_EXPANDED_ARCHIVE_BYTES:
                    raise ValueError("release archive exceeds the expanded-byte limit")
                seen.add(member.name)
                target = safe_destination(destination, member.name)
                target.parent.mkdir(parents=True, exist_ok=True)
                opened = source.extractfile(member)
                if opened is None:
                    raise ValueError(f"release archive file could not be read: {member.name}")
                with target.open("xb") as output:
                    shutil.copyfileobj(opened, output, length=1024 * 1024)
                if target.stat().st_size != member.size:
                    raise ValueError(f"release archive member changed while reading: {member.name}")
                target.chmod(member.mode & 0o777)
        return
    if archive.suffix == ".zip":
        with zipfile.ZipFile(archive) as source:
            for member in source.infolist():
                unix_mode = member.external_attr >> 16
                if member.is_dir() or (unix_mode & 0o170000) == 0o120000 or member.filename in seen:
                    raise ValueError(f"release archive contains a non-file or duplicate entry: {member.filename}")
                expanded += member.file_size
                if member.file_size < 0 or expanded > MAXIMUM_EXPANDED_ARCHIVE_BYTES:
                    raise ValueError("release archive exceeds the expanded-byte limit")
                seen.add(member.filename)
                target = safe_destination(destination, member.filename)
                target.parent.mkdir(parents=True, exist_ok=True)
                with source.open(member) as opened, target.open("xb") as output:
                    shutil.copyfileobj(opened, output, length=1024 * 1024)
                if target.stat().st_size != member.file_size:
                    raise ValueError(f"release archive member changed while reading: {member.filename}")
        return
    raise ValueError(f"unsupported release archive: {archive.name}")


def run_json(command: list[str], label: str, cwd: pathlib.Path | None = None) -> dict[str, Any]:
    completed = subprocess.run(
        command,
        cwd=cwd,
        check=False,
        stdin=subprocess.DEVNULL,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        encoding="utf-8",
        errors="strict",
    )
    if completed.returncode != 0:
        detail = (completed.stderr.strip() or completed.stdout.strip() or "no diagnostic output")[:4_000]
        raise RuntimeError(f"{label} failed with exit code {completed.returncode}: {detail}")
    try:
        value = json.loads(completed.stdout)
    except json.JSONDecodeError as error:
        raise RuntimeError(f"{label} did not return JSON") from error
    if not isinstance(value, dict):
        raise RuntimeError(f"{label} did not return a JSON object")
    return value


def global_prc_command(
    prefix: pathlib.Path, arguments: list[str], *, windows: bool,
    command_processor: str | None = None,
) -> list[str]:
    """Return the real command npm exposes for one isolated global prefix."""
    shim = prefix / ("prc.cmd" if windows else "bin/prc")
    if shim.is_symlink() or not shim.is_file():
        raise RuntimeError("global npm install did not expose the prc command")
    if not windows:
        return [str(shim), *arguments]
    if not command_processor or not pathlib.Path(command_processor).is_absolute():
        raise RuntimeError("Windows command processor is unavailable for the global prc smoke test")
    command_line = subprocess.list2cmdline([str(shim), *arguments])
    return [command_processor, "/d", "/s", "/c", command_line]


def record(records: Any, **wanted: str) -> dict[str, Any]:
    matches = [item for item in records if isinstance(item, dict) and all(item.get(key) == value for key, value in wanted.items())]
    if len(matches) != 1:
        raise ValueError(f"release manifest does not contain exactly one matching record: {wanted}")
    return matches[0]


def checked_artifact(release: pathlib.Path, item: dict[str, Any]) -> pathlib.Path:
    name, digest, size = item.get("name"), item.get("sha256"), item.get("size")
    path = (release.resolve() / str(name)).resolve()
    if (
        path.parent != release.resolve()
        or path.is_symlink()
        or not path.is_file()
        or not isinstance(size, int)
        or isinstance(size, bool)
        or size < 1
        or size > MAXIMUM_ARTIFACT_BYTES
    ):
        raise ValueError(f"release artifact is missing or unsafe: {name}")
    before = path.stat()
    hasher = hashlib.sha256()
    total = 0
    with path.open("rb") as source:
        opened = os.fstat(source.fileno())
        if not os.path.samestat(before, opened):
            raise ValueError(f"release artifact changed while opening: {name}")
        while chunk := source.read(1024 * 1024):
            total += len(chunk)
            if total > MAXIMUM_ARTIFACT_BYTES:
                raise ValueError(f"release artifact exceeds its byte limit: {name}")
            hasher.update(chunk)
        final_open = os.fstat(source.fileno())
    final_path = path.stat()
    if (
        total != size
        or hasher.hexdigest() != digest
        or not os.path.samestat(before, final_open)
        or not os.path.samestat(before, final_path)
        or final_open.st_size != before.st_size
        or final_open.st_mtime_ns != before.st_mtime_ns
    ):
        raise ValueError(f"release artifact does not match its manifest: {name}")
    return path


def verify_result(value: dict[str, Any], label: str) -> None:
    if (
        value.get("schema_version") != "prc.run/v0.13"
        or not isinstance(value.get("results"), list)
        or not isinstance(value.get("control_results"), list)
        or len(value["control_results"]) != 10_042
        or value.get("control_catalog", {}).get("control_count") != 10_042
    ):
        raise RuntimeError(f"{label} returned an incomplete or incorrectly bound scan")


def smoke(release: pathlib.Path, version: str, commit: str) -> None:
    if not build_release.is_release_version(version) or not build_release.COMMIT.fullmatch(commit):
        raise ValueError("version or commit is invalid")
    goos, goarch, npm_os, npm_cpu = host()
    manifest_path = release / f"prc_{version}_release-manifest.json"
    manifest = load_json(manifest_path, "release manifest")
    if manifest.get("schema_version") != "prc.release-manifest/v0.3" or manifest.get("version") != version or manifest.get("source_commit") != commit:
        raise ValueError("release manifest identity does not match the requested smoke test")

    archive_item = record(manifest.get("artifacts", []), goos=goos, goarch=goarch)
    archive = checked_artifact(release, archive_item)
    platform_package = checked_artifact(release, record(manifest.get("npm_packages", []), kind="platform", os=npm_os, cpu=npm_cpu))
    launcher_package = checked_artifact(release, record(manifest.get("npm_packages", []), kind="launcher"))

    with tempfile.TemporaryDirectory(prefix="prc-release-smoke-") as temporary:
        root = pathlib.Path(temporary)
        extracted = root / "native"
        extracted.mkdir()
        extract_archive(archive, extracted)
        binary_name = "prc.exe" if goos == "windows" else "prc"
        binary = extracted / f"prc_{version}_{goos}_{goarch}" / binary_name
        version_info = run_json([str(binary), "version", "--format", "json"], "native version check")
        if version_info.get("version") != version or version_info.get("revision") != commit:
            raise RuntimeError("native archive reports the wrong release identity")

        project = root / "project"
        project.mkdir()
        (project / "README.md").write_text("# Release smoke fixture\n", encoding="utf-8")
        native_result = run_json(
            [str(binary), "scan", str(project), "--format", "json", "--no-report", "--exit-policy", "never"],
            "native scan",
            cwd=project,
        )
        verify_result(native_result, "native archive")

        npm_project = root / "npm"
        npm_project.mkdir()
        npm_executable = shutil.which("npm.cmd" if os.name == "nt" else "npm")
        node_executable = shutil.which("node")
        if not npm_executable or not node_executable:
            raise RuntimeError("Node.js and npm are required for the npm launcher smoke test")
        npm_environment = {
            name: value for name, value in os.environ.items()
            if not name.upper().startswith("NPM_CONFIG_")
        }
        npm_cache = root / "npm-cache"
        npm_cache.mkdir()
        npm_environment["NPM_CONFIG_CACHE"] = str(npm_cache)
        for filename, variable in (
            ("user.npmrc", "NPM_CONFIG_USERCONFIG"),
            ("global.npmrc", "NPM_CONFIG_GLOBALCONFIG"),
        ):
            config = root / filename
            config.touch(mode=0o600, exist_ok=False)
            npm_environment[variable] = str(config)
        install = subprocess.run(
            [
                npm_executable, "install", "--ignore-scripts", "--offline", "--no-audit", "--no-fund",
                "--package-lock=false", str(platform_package), str(launcher_package),
            ],
            cwd=npm_project,
            env=npm_environment,
            check=False,
            stdin=subprocess.DEVNULL,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            encoding="utf-8",
            errors="strict",
        )
        if install.returncode != 0:
            raise RuntimeError(f"offline npm tarball install failed: {(install.stderr or install.stdout)[:4_000]}")
        launcher = npm_project / "node_modules" / "@marinjursic" / "prc" / "bin" / "prc.js"
        npm_version = run_json([node_executable, str(launcher), "version", "--format", "json"], "npm launcher version check")
        if npm_version.get("version") != version or npm_version.get("revision") != commit:
            raise RuntimeError("npm launcher reports the wrong release identity")
        npm_result = run_json(
            [node_executable, str(launcher), "scan", str(project), "--format", "json", "--no-report", "--exit-policy", "never"],
            "npm launcher scan",
            cwd=project,
        )
        verify_result(npm_result, "npm launcher")

        global_prefix = root / "global"
        global_install = subprocess.run(
            [
                npm_executable, "install", "--global", "--prefix", str(global_prefix),
                "--ignore-scripts", "--offline", "--no-audit", "--no-fund",
                "--package-lock=false", str(platform_package), str(launcher_package),
            ],
            cwd=root,
            env=npm_environment,
            check=False,
            stdin=subprocess.DEVNULL,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            encoding="utf-8",
            errors="strict",
        )
        if global_install.returncode != 0:
            detail = (global_install.stderr or global_install.stdout)[:4_000]
            raise RuntimeError(f"offline global npm install failed: {detail}")
        command_processor = shutil.which("cmd.exe") if os.name == "nt" else None
        global_version = run_json(
            global_prc_command(
                global_prefix, ["version", "--format", "json"],
                windows=os.name == "nt", command_processor=command_processor,
            ),
            "global prc version check",
        )
        if global_version.get("version") != version or global_version.get("revision") != commit:
            raise RuntimeError("global prc command reports the wrong release identity")
        global_result = run_json(
            global_prc_command(
                global_prefix,
                ["scan", str(project), "--format", "json", "--no-report", "--exit-policy", "never"],
                windows=os.name == "nt", command_processor=command_processor,
            ),
            "global prc scan",
            cwd=project,
        )
        verify_result(global_result, "global prc command")
        if (project / "node_modules").exists():
            raise RuntimeError("global prc install or scan wrote node_modules into the target project")
    print(f"release smoke test passed for {goos}/{goarch}")


def parser() -> argparse.ArgumentParser:
    result = argparse.ArgumentParser(description=__doc__)
    result.add_argument("--release", required=True, type=pathlib.Path)
    result.add_argument("--version", required=True)
    result.add_argument("--commit", required=True)
    return result


def main() -> int:
    try:
        arguments = parser().parse_args()
        smoke(arguments.release.resolve(), arguments.version, arguments.commit)
    except (OSError, ValueError, RuntimeError, subprocess.SubprocessError, tarfile.TarError, zipfile.BadZipFile) as error:
        print(f"release smoke test failed: {error}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
