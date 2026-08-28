#!/usr/bin/env python3
"""Build deterministic npm tarballs for the native PRC scanner release."""

from __future__ import annotations

import gzip
import hashlib
import io
import json
import pathlib
import tarfile
from typing import Any


ROOT = pathlib.Path(__file__).resolve().parents[1]
SCOPE = "@marinjursic"
LAUNCHER_NAME = f"{SCOPE}/prc"
PLATFORMS = {
    ("darwin", "arm64"): ("darwin-arm64", f"{SCOPE}/prc-darwin-arm64", "darwin", "arm64"),
    ("darwin", "amd64"): ("darwin-x64", f"{SCOPE}/prc-darwin-x64", "darwin", "x64"),
    ("linux", "arm64"): ("linux-arm64", f"{SCOPE}/prc-linux-arm64", "linux", "arm64"),
    ("linux", "amd64"): ("linux-x64", f"{SCOPE}/prc-linux-x64", "linux", "x64"),
    ("windows", "arm64"): ("win32-arm64", f"{SCOPE}/prc-windows-arm64", "win32", "arm64"),
    ("windows", "amd64"): ("win32-x64", f"{SCOPE}/prc-windows-x64", "win32", "x64"),
}
MAXIMUM_SUPPORT_FILES = 512
MAXIMUM_SUPPORT_BYTES = 24 * 1024 * 1024
MAXIMUM_SUPPORT_FILE_BYTES = 16 * 1024 * 1024
MAXIMUM_PLATFORM_TARBALL_BYTES = 24 * 1024 * 1024
MAXIMUM_LAUNCHER_TARBALL_BYTES = 128 * 1024


def allowed_support_name(name: str) -> bool:
    prefixes = (
        "adapters/",
        "catalog/",
        "docs/checklists/",
        "docs/engineering/",
        "fixtures/benchmarks/",
        "packs/",
        "schemas/",
    )
    return name == "THIRD_PARTY_NOTICES.md" or any(
        len(name) > len(prefix) and name.startswith(prefix) for prefix in prefixes
    )


def json_bytes(document: Any) -> bytes:
    return (json.dumps(document, indent=2, sort_keys=True, ensure_ascii=False) + "\n").encode("utf-8")


def sha256(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def package_filename(name: str, version: str) -> str:
    return f"{name.removeprefix('@').replace('/', '-')}-{version}.tgz"


def create_tgz(path: pathlib.Path, entries: list[tuple[str, bytes, int]], epoch: int) -> None:
    seen: set[str] = set()
    with path.open("wb") as output:
        with gzip.GzipFile(filename="", mode="wb", fileobj=output, compresslevel=9, mtime=epoch) as compressed:
            with tarfile.open(fileobj=compressed, mode="w", format=tarfile.USTAR_FORMAT) as archive:
                for name, data, mode in sorted(entries):
                    if not name or name.startswith("/") or "\\" in name or ".." in pathlib.PurePosixPath(name).parts or name in seen:
                        raise ValueError(f"unsafe or duplicate npm package path: {name}")
                    seen.add(name)
                    information = tarfile.TarInfo(f"package/{name}")
                    information.size = len(data)
                    information.mtime = epoch
                    information.mode = mode
                    information.uid = information.gid = 0
                    information.uname = information.gname = ""
                    archive.addfile(information, io.BytesIO(data))


def common_metadata(name: str, version: str, description: str) -> dict[str, Any]:
    return {
        "name": name,
        "version": version,
        "description": description,
        "license": "MIT",
        "contentPolicy": {"class": "dual-use"},
        "author": "Marin Jursic",
        "homepage": "https://marinjursic.github.io/production-readiness-checklist/",
        "repository": {
            "type": "git",
            "url": "git+https://github.com/MarinJursic/production-readiness-checklist.git",
        },
        "bugs": "https://github.com/MarinJursic/production-readiness-checklist/issues",
        "publishConfig": {"access": "public", "provenance": True},
    }


def validate_support(support: list[tuple[str, bytes, int]]) -> list[tuple[str, bytes, int]]:
    if not support or len(support) > MAXIMUM_SUPPORT_FILES:
        raise ValueError("npm runtime support file count is outside its package limit")
    normalized = sorted(support)
    total = 0
    previous = ""
    for name, data, mode in normalized:
        pure = pathlib.PurePosixPath(name)
        if (
            not name
            or name <= previous
            or pure.is_absolute()
            or pure.as_posix() != name
            or "\\" in name
            or ".." in pure.parts
            or not allowed_support_name(name)
            or any(ord(character) < 0x20 or 0x7F <= ord(character) <= 0x9F for character in name)
            or mode != 0o644
            or len(data) > MAXIMUM_SUPPORT_FILE_BYTES
        ):
            raise ValueError(f"unsafe npm runtime support entry: {name}")
        previous = name
        total += len(data)
        if total > MAXIMUM_SUPPORT_BYTES:
            raise ValueError("npm runtime support files exceed their byte limit")
    return normalized


def package_artifact(path: pathlib.Path, package_name: str, kind: str, os_name: str = "", cpu: str = "") -> dict[str, Any]:
    data = path.read_bytes()
    result: dict[str, Any] = {
        "name": path.name,
        "package_name": package_name,
        "kind": kind,
        "sha256": sha256(data),
        "size": len(data),
    }
    if os_name:
        result["os"] = os_name
        result["cpu"] = cpu
    return result


def build_packages(
    *,
    version: str,
    commit: str,
    built_at: str,
    epoch: int,
    binaries: dict[tuple[str, str], pathlib.Path],
    support: list[tuple[str, bytes, int]],
    output: pathlib.Path,
) -> list[dict[str, Any]]:
    output.mkdir(mode=0o700)
    support = validate_support(support)
    launcher_path = ROOT / "npm" / "prc" / "bin" / "prc.js"
    launcher_readme_path = ROOT / "npm" / "prc" / "README.md"
    license_path = ROOT / "LICENSE"
    disclosure_path = ROOT / "DISCLOSURE"
    npm_disclosure_path = ROOT / "npm" / "prc" / "DISCLOSURE"
    for source in (launcher_path, launcher_readme_path, license_path, disclosure_path, npm_disclosure_path):
        if source.is_symlink() or not source.is_file():
            raise ValueError(f"npm source must be a regular file: {source.relative_to(ROOT)}")
    launcher = launcher_path.read_bytes()
    forbidden_launcher_fragments = (
        b"postinstall",
        b"child_process.exec",
        b"shell: true",
        b"node:http",
        b"node:https",
        b"node:http2",
        b"node:net",
        b"node:tls",
        b"fetch(",
    )
    if not launcher.startswith(b"#!/usr/bin/env node\n") or any(
        fragment in launcher for fragment in forbidden_launcher_fragments
    ):
        raise ValueError("npm launcher violates the no-shell, no-install-hook contract")
    license_data = license_path.read_bytes()
    disclosure_data = disclosure_path.read_bytes()
    if disclosure_data != npm_disclosure_path.read_bytes() or not disclosure_data.isascii():
        raise ValueError("npm dual-use disclosure copies must be identical plain text")
    platform_readme = (
        "# Production Readiness Checklist native package\n\n"
        "This is one operating-system and CPU build of the native `prc` scanner. It is an "
        "internal package installed automatically by `@marinjursic/prc`; install and invoke "
        "the launcher instead of depending on this package directly.\n\n"
        "The package has no npm dependencies and no install, update, or lifecycle scripts. "
        "It does not download a binary after installation. The launcher selects the exact "
        "package for the current machine, verifies the release-bound manifest, checks the "
        "SHA-256 and size of the native executable and every runtime support file, and fails "
        "closed if anything is missing or changed.\n\n"
        "A normal scan reads a user-selected project and writes a private report outside that "
        "project. It does not fix project files, install project dependencies, or run project "
        "scripts. Security analyzers and Codex or Claude review require explicit opt-in. See "
        "the launcher package and project documentation for installation, supported commands, "
        "the dual-use disclosure, checksums, SBOM, and provenance information.\n"
    ).encode("utf-8")

    platform_bindings: dict[str, dict[str, str]] = {}
    artifacts: list[dict[str, Any]] = []
    for target, (key, package_name, npm_os, npm_cpu) in sorted(PLATFORMS.items()):
        binary_path = binaries.get(target)
        if binary_path is None or binary_path.is_symlink() or not binary_path.is_file():
            raise ValueError(f"missing regular release binary for {target[0]}/{target[1]}")
        binary = binary_path.read_bytes()
        binary_name = "prc.exe" if target[0] == "windows" else "prc"
        support_records = [
            {"path": f"bin/{name}", "size": len(data), "sha256": sha256(data)}
            for name, data, _mode in support
        ]
        manifest = {
            "schema_version": "prc.npm-platform/v0.2",
            "package_name": package_name,
            "version": version,
            "source_commit": commit,
            "built_at": built_at,
            "binary_path": f"bin/{binary_name}",
            "binary_sha256": sha256(binary),
            "support_file_count": len(support_records),
            "support_bytes": sum(record["size"] for record in support_records),
            "support_files": support_records,
        }
        manifest_data = json_bytes(manifest)
        platform_bindings[key] = {"package_name": package_name, "manifest_sha256": sha256(manifest_data)}
        metadata = common_metadata(package_name, version, f"Native {npm_os}/{npm_cpu} binary for Production Readiness Checklist")
        metadata.update(
            {
                "os": [npm_os],
                "cpu": [npm_cpu],
                "files": ["bin", "manifest.json", "README.md", "LICENSE", "DISCLOSURE"],
            }
        )
        entries = [
            ("LICENSE", license_data, 0o644),
            ("DISCLOSURE", disclosure_data, 0o644),
            ("README.md", platform_readme, 0o644),
            ("manifest.json", manifest_data, 0o644),
            ("package.json", json_bytes(metadata), 0o644),
            (f"bin/{binary_name}", binary, 0o755),
        ]
        entries.extend((f"bin/{name}", data, mode) for name, data, mode in support)
        destination = output / package_filename(package_name, version)
        create_tgz(destination, entries, epoch)
        if destination.stat().st_size > MAXIMUM_PLATFORM_TARBALL_BYTES:
            raise ValueError(f"npm platform package exceeds its compressed size budget: {package_name}")
        artifacts.append(package_artifact(destination, package_name, "platform", npm_os, npm_cpu))

    optional_dependencies = {binding[1]: version for binding in PLATFORMS.values()}
    launcher_metadata = common_metadata(LAUNCHER_NAME, version, "Safe launcher for Production Readiness Checklist")
    launcher_metadata.update(
        {
            "type": "module",
            "bin": {"prc": "bin/prc.js"},
            "files": ["bin/prc.js", "platforms.json", "README.md", "LICENSE", "DISCLOSURE"],
            "engines": {"node": ">=22.14.0"},
            "keywords": ["production-readiness", "scanner", "security", "checklist", "devops"],
            "optionalDependencies": dict(sorted(optional_dependencies.items())),
        }
    )
    platforms = {
        "schema_version": "prc.npm-platforms/v0.2",
        "version": version,
        "platforms": dict(sorted(platform_bindings.items())),
    }
    readme = launcher_readme_path.read_text(encoding="utf-8").replace("@marinjursic/prc@VERSION", f"@marinjursic/prc@{version}")
    launcher_entries = [
        ("LICENSE", license_data, 0o644),
        ("DISCLOSURE", disclosure_data, 0o644),
        ("README.md", readme.encode("utf-8"), 0o644),
        ("bin/prc.js", launcher, 0o755),
        ("package.json", json_bytes(launcher_metadata), 0o644),
        ("platforms.json", json_bytes(platforms), 0o644),
    ]
    destination = output / package_filename(LAUNCHER_NAME, version)
    create_tgz(destination, launcher_entries, epoch)
    if destination.stat().st_size > MAXIMUM_LAUNCHER_TARBALL_BYTES:
        raise ValueError("npm launcher package exceeds its compressed size budget")
    artifacts.append(package_artifact(destination, LAUNCHER_NAME, "launcher"))
    return sorted(artifacts, key=lambda item: item["package_name"])
