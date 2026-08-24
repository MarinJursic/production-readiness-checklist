#!/usr/bin/env python3
"""Verify and publish one Everylast npm release with npm trusted publishing."""

from __future__ import annotations

import argparse
import base64
import gzip
import hashlib
import io
import json
import os
import pathlib
import re
import subprocess
import sys
import tarfile
import tempfile
from dataclasses import dataclass
from typing import Any, Callable

try:
    from scripts import build_release, npm_distribution
except ModuleNotFoundError:  # Direct `python scripts/publish_npm.py` execution.
    import build_release
    import npm_distribution


REGISTRY = "https://registry.npmjs.org/"
MINIMUM_NODE = (22, 14, 0)
MINIMUM_NPM = (11, 5, 1)
MAXIMUM_NPM_PACKAGE_BYTES = 256 * 1024 * 1024
MAXIMUM_EXPANDED_PACKAGE_BYTES = 256 * 1024 * 1024
VERSION_TEXT = re.compile(r"^v?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:[-+].*)?$")
Run = Callable[[list[str]], subprocess.CompletedProcess[str]]


@dataclass(frozen=True)
class Package:
    name: str
    version: str
    path: pathlib.Path
    integrity: str
    kind: str


def default_run(command: list[str]) -> subprocess.CompletedProcess[str]:
    environment = {
        name: value
        for name, value in os.environ.items()
        if not name.upper().startswith("NPM_CONFIG_")
    }
    # Do not let a saved user/global npm token or repository .npmrc influence
    # trusted publishing. GitHub's OIDC request variables remain available.
    with tempfile.TemporaryDirectory(prefix="everylast-npm-publish-") as directory:
        for filename, variable in (("user.npmrc", "NPM_CONFIG_USERCONFIG"), ("global.npmrc", "NPM_CONFIG_GLOBALCONFIG")):
            path = pathlib.Path(directory, filename)
            path.touch(mode=0o600, exist_ok=False)
            environment[variable] = str(path)
        return subprocess.run(
            command,
            cwd=directory,
            env=environment,
            check=False,
            stdin=subprocess.DEVNULL,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
        )


def clean_text(value: str, limit: int = 2_000) -> str:
    return "".join(character if character.isprintable() else "?" for character in value)[:limit]


def parse_tool_version(value: str, label: str) -> tuple[int, int, int]:
    match = VERSION_TEXT.fullmatch(value.strip())
    if not match:
        raise ValueError(f"{label} returned an unrecognized version")
    return tuple(int(part) for part in match.groups())


def command_output(run: Run, command: list[str], label: str) -> str:
    completed = run(command)
    if completed.returncode != 0:
        detail = clean_text(completed.stderr.strip() or completed.stdout.strip() or "no diagnostic output")
        raise RuntimeError(f"{label} failed with exit code {completed.returncode}: {detail}")
    return completed.stdout.strip()


def npm_error_code(completed: subprocess.CompletedProcess[str]) -> str:
    for candidate in (completed.stdout, completed.stderr):
        try:
            document = json.loads(candidate)
        except (TypeError, json.JSONDecodeError):
            continue
        if isinstance(document, dict):
            error = document.get("error")
            if isinstance(error, dict) and isinstance(error.get("code"), str):
                return error["code"]
            if isinstance(document.get("code"), str):
                return document["code"]
    return ""


def read_tgz_package(path: pathlib.Path) -> tuple[dict[str, Any], set[str]]:
    if path.is_symlink() or not path.is_file() or path.stat().st_size > MAXIMUM_NPM_PACKAGE_BYTES:
        raise ValueError(f"npm artifact must be a bounded regular file: {path.name}")
    try:
        with gzip.open(path, "rb") as compressed:
            raw = compressed.read(MAXIMUM_EXPANDED_PACKAGE_BYTES + 1)
        if len(raw) > MAXIMUM_EXPANDED_PACKAGE_BYTES:
            raise ValueError(f"expanded npm artifact is too large: {path.name}")
        with tarfile.open(fileobj=io.BytesIO(raw), mode="r:") as archive:
            names: set[str] = set()
            metadata: bytes | None = None
            for member in archive.getmembers():
                pure = pathlib.PurePosixPath(member.name)
                if (
                    not member.isfile()
                    or member.issym()
                    or member.islnk()
                    or pure.is_absolute()
                    or ".." in pure.parts
                    or member.name in names
                ):
                    raise ValueError(f"unsafe npm archive member in {path.name}: {member.name}")
                names.add(member.name)
                if member.name == "package/package.json":
                    source = archive.extractfile(member)
                    metadata = source.read(1024 * 1024 + 1) if source else None
            if metadata is None or len(metadata) > 1024 * 1024:
                raise ValueError(f"npm package metadata is missing or too large: {path.name}")
    except (OSError, EOFError, tarfile.TarError) as error:
        raise ValueError(f"invalid npm tarball: {path.name}") from error
    try:
        document = json.loads(metadata)
    except (UnicodeDecodeError, json.JSONDecodeError) as error:
        raise ValueError(f"invalid npm package metadata: {path.name}") from error
    if not isinstance(document, dict):
        raise ValueError(f"npm package metadata is not an object: {path.name}")
    return document, names


def artifact_hashes(path: pathlib.Path) -> tuple[int, str, bytes]:
    if path.is_symlink() or not path.is_file():
        raise ValueError(f"npm artifact is missing or unsafe: {path.name}")
    before = path.stat()
    if before.st_size < 1 or before.st_size > MAXIMUM_NPM_PACKAGE_BYTES:
        raise ValueError(f"npm artifact exceeds the package size limit: {path.name}")
    sha256 = hashlib.sha256()
    sha512 = hashlib.sha512()
    with path.open("rb") as source:
        opened = os.fstat(source.fileno())
        if not os.path.samestat(before, opened):
            raise ValueError(f"npm artifact changed while opening: {path.name}")
        total = 0
        while chunk := source.read(1024 * 1024):
            total += len(chunk)
            if total > MAXIMUM_NPM_PACKAGE_BYTES:
                raise ValueError(f"npm artifact exceeds the package size limit: {path.name}")
            sha256.update(chunk)
            sha512.update(chunk)
        final_open = os.fstat(source.fileno())
    final_path = path.stat()
    if (
        total != before.st_size
        or not os.path.samestat(before, final_open)
        or not os.path.samestat(before, final_path)
        or final_open.st_size != before.st_size
        or final_open.st_mtime_ns != before.st_mtime_ns
    ):
        raise ValueError(f"npm artifact changed while hashing: {path.name}")
    return total, sha256.hexdigest(), sha512.digest()


def load_packages(release: pathlib.Path, manifest_path: pathlib.Path, version: str) -> list[Package]:
    if not build_release.SEMVER.fullmatch(version):
        raise ValueError("version must be an exact semantic version")
    if manifest_path.is_symlink() or not manifest_path.is_file() or manifest_path.stat().st_size > 16 * 1024 * 1024:
        raise ValueError("release manifest must be a bounded regular file")
    try:
        manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    except (OSError, UnicodeDecodeError, json.JSONDecodeError) as error:
        raise ValueError("release manifest must be valid UTF-8 JSON") from error
    if (
        not isinstance(manifest, dict)
        or manifest.get("schema_version") != "prc.release-manifest/v0.3"
        or manifest.get("product") != "everylast"
        or manifest.get("version") != version
    ):
        raise ValueError("release manifest identity does not match the requested npm release")
    records = manifest.get("npm_packages")
    if not isinstance(records, list) or len(records) != 7:
        raise ValueError("release manifest must bind exactly seven npm packages")

    expected_names = {npm_distribution.LAUNCHER_NAME, *(binding[1] for binding in npm_distribution.PLATFORMS.values())}
    packages: list[Package] = []
    seen: set[str] = set()
    release_root = release.resolve()
    for record in records:
        if not isinstance(record, dict):
            raise ValueError("release manifest contains an invalid npm package record")
        name, filename, digest, size, kind = (
            record.get("package_name"),
            record.get("name"),
            record.get("sha256"),
            record.get("size"),
            record.get("kind"),
        )
        if (
            name not in expected_names
            or name in seen
            or filename != npm_distribution.package_filename(name, version)
            or not isinstance(digest, str)
            or not re.fullmatch(r"[0-9a-f]{64}", digest)
            or not isinstance(size, int)
            or isinstance(size, bool)
            or size < 1
            or size > MAXIMUM_NPM_PACKAGE_BYTES
            or kind not in {"launcher", "platform"}
            or (name == npm_distribution.LAUNCHER_NAME) != (kind == "launcher")
        ):
            raise ValueError("release manifest contains an inconsistent npm package record")
        path = (release_root / filename).resolve()
        if path.parent != release_root or path.is_symlink() or not path.is_file():
            raise ValueError(f"npm artifact is missing or unsafe: {filename}")
        actual_size, actual_sha256, actual_sha512 = artifact_hashes(path)
        if actual_size != size or actual_sha256 != digest:
            raise ValueError(f"npm artifact does not match the release manifest: {filename}")
        metadata, members = read_tgz_package(path)
        if (
            metadata.get("name") != name
            or metadata.get("version") != version
            or metadata.get("license") != "MIT"
            or metadata.get("scripts") is not None
            or metadata.get("publishConfig") != {"access": "public"}
            or "package/package.json" not in members
        ):
            raise ValueError(f"npm tarball identity or script-free contract is invalid: {filename}")
        integrity = "sha512-" + base64.b64encode(actual_sha512).decode("ascii")
        packages.append(Package(name, version, path, integrity, kind))
        seen.add(name)
    if seen != expected_names:
        raise ValueError("release manifest does not contain the exact npm package set")
    return sorted(packages, key=lambda item: (item.kind == "launcher", item.name))


def registry_integrity(run: Run, package: Package) -> str | None:
    command = [
        "npm", "view", f"{package.name}@{package.version}", "dist.integrity",
        "--json", "--registry", REGISTRY,
    ]
    completed = run(command)
    if completed.returncode != 0:
        if npm_error_code(completed) == "E404":
            return None
        detail = clean_text(completed.stderr.strip() or completed.stdout.strip() or "no diagnostic output")
        raise RuntimeError(f"npm registry lookup failed for {package.name}: {detail}")
    try:
        integrity = json.loads(completed.stdout)
    except json.JSONDecodeError as error:
        raise RuntimeError(f"npm registry returned invalid JSON for {package.name}") from error
    if not isinstance(integrity, str) or not integrity.startswith("sha512-"):
        raise RuntimeError(f"npm registry returned no SHA-512 integrity for {package.name}")
    return integrity


def publish(packages: list[Package], run: Run = default_run) -> None:
    for secret_name in ("NODE_AUTH_TOKEN", "NPM_TOKEN"):
        if os.environ.get(secret_name):
            raise RuntimeError(f"{secret_name} is set; this release requires token-free npm trusted publishing")
    node = parse_tool_version(command_output(run, ["node", "--version"], "Node.js version check"), "Node.js")
    npm = parse_tool_version(command_output(run, ["npm", "--version"], "npm version check"), "npm")
    if node < MINIMUM_NODE or npm < MINIMUM_NPM:
        raise RuntimeError("npm trusted publishing requires Node.js >=22.14.0 and npm >=11.5.1")

    for package in packages:
        existing = registry_integrity(run, package)
        if existing is not None:
            if existing != package.integrity:
                raise RuntimeError(f"npm already has different bytes for immutable package {package.name}@{package.version}")
            print(f"verified existing {package.name}@{package.version}")
            continue
        tag = "next" if "-" in package.version else "latest"
        command = [
            "npm", "publish", str(package.path), "--access", "public", "--tag", tag,
            "--provenance", "--ignore-scripts", "--registry", REGISTRY,
        ]
        command_output(run, command, f"npm publish {package.name}@{package.version}")
        published = registry_integrity(run, package)
        if published != package.integrity:
            raise RuntimeError(f"npm did not return the expected integrity for {package.name}@{package.version}")
        print(f"published and verified {package.name}@{package.version} ({tag})")


def parser() -> argparse.ArgumentParser:
    result = argparse.ArgumentParser(description=__doc__)
    result.add_argument("--release", required=True, type=pathlib.Path, help="verified release directory")
    result.add_argument("--manifest", required=True, type=pathlib.Path, help="release manifest path")
    result.add_argument("--version", required=True, help="exact scanner/npm semantic version")
    return result


def main() -> int:
    try:
        arguments = parser().parse_args()
        packages = load_packages(arguments.release, arguments.manifest.resolve(), arguments.version)
        publish(packages)
    except (OSError, ValueError, RuntimeError) as error:
        print(f"npm release failed: {error}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
