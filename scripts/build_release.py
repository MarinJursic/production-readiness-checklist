#!/usr/bin/env python3
"""Build deterministic, self-describing Production Readiness Scanner releases."""

from __future__ import annotations

import argparse
import datetime as dt
import gzip
import hashlib
import io
import json
import os
import pathlib
import platform
import re
import shutil
import stat
import subprocess
import sys
import tarfile
import tempfile
import zipfile
from typing import Any

try:
    from scripts import npm_distribution
except ModuleNotFoundError:  # Direct `python scripts/build_release.py` execution.
    import npm_distribution


ROOT = pathlib.Path(__file__).resolve().parents[1]
MODULE = "github.com/MarinJursic/production-readiness-checklist"
SEMVER = re.compile(
    r"^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)"
    r"(?:-(?:(?:0|[1-9][0-9]*)|(?:[0-9A-Za-z-]*[A-Za-z-][0-9A-Za-z-]*))"
    r"(?:\.(?:(?:0|[1-9][0-9]*)|(?:[0-9A-Za-z-]*[A-Za-z-][0-9A-Za-z-]*)))*)?$"
)
COMMIT = re.compile(r"^(?:[0-9a-f]{40}|[0-9a-f]{64})$")
TARGETS = (
    ("linux", "amd64"),
    ("linux", "arm64"),
    ("darwin", "amd64"),
    ("darwin", "arm64"),
    ("windows", "amd64"),
    ("windows", "arm64"),
)
MAX_RELEASE_INPUT_FILES = 10_000
MAX_RELEASE_INPUT_BYTES = 256 * 1024 * 1024
MAX_SBOM_BYTES = 16 * 1024 * 1024


def parse_timestamp(value: str) -> tuple[str, int]:
    candidate = value[:-1] + "+00:00" if value.endswith("Z") else value
    try:
        parsed = dt.datetime.fromisoformat(candidate)
    except ValueError as error:
        raise ValueError("built-at must be an RFC 3339 timestamp") from error
    if parsed.tzinfo is None:
        raise ValueError("built-at must include a timezone")
    normalized = parsed.astimezone(dt.timezone.utc).replace(microsecond=0)
    return normalized.isoformat().replace("+00:00", "Z"), int(normalized.timestamp())


def sha256(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def file_identity(path: pathlib.Path) -> dict[str, Any]:
    data = path.read_bytes()
    return {"sha256": sha256(data), "size": len(data)}


def write_checksums(root: pathlib.Path) -> None:
    paths = sorted(path for path in root.iterdir() if path.name != "SHA256SUMS")
    if not paths:
        raise ValueError("release directory contains no checksum subjects")
    for path in paths:
        if path.is_symlink() or not path.is_file():
            raise ValueError(f"release checksum subject is not a regular file: {path.name}")
    text = "".join(f"{file_identity(path)['sha256']}  {path.name}\n" for path in paths)
    (root / "SHA256SUMS").write_text(text, encoding="ascii", newline="\n")


def run_json(command: list[str], *, environment: dict[str, str] | None = None) -> dict[str, Any]:
    completed = subprocess.run(
        command,
        cwd=ROOT,
        env=environment,
        check=False,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )
    if completed.returncode != 0:
        detail = completed.stderr.strip() or completed.stdout.strip() or "no diagnostic output"
        raise RuntimeError(f"command failed ({completed.returncode}): {' '.join(command)}: {detail}")
    try:
        result = json.loads(completed.stdout)
    except json.JSONDecodeError as error:
        raise RuntimeError(f"command did not emit JSON: {' '.join(command)}") from error
    if not isinstance(result, dict):
        raise RuntimeError(f"command did not emit a JSON object: {' '.join(command)}")
    return result


def normalized_sbom(source: pathlib.Path, version: str, commit: str) -> bytes:
    if source.is_symlink() or not source.is_file() or source.stat().st_size > MAX_SBOM_BYTES:
        raise ValueError("SBOM must be a regular JSON file no larger than 16 MiB")
    try:
        document = json.loads(source.read_text(encoding="utf-8"))
    except (OSError, UnicodeDecodeError, json.JSONDecodeError) as error:
        raise ValueError("SBOM must contain valid UTF-8 JSON") from error
    if not isinstance(document, dict) or document.get("bomFormat") != "CycloneDX" or document.get("specVersion") != "1.6":
        raise ValueError("SBOM must be a CycloneDX 1.6 JSON document")
    if "serialNumber" in document:
        raise ValueError("release SBOM must omit a nondeterministic serial number")
    metadata = document.get("metadata")
    component = metadata.get("component") if isinstance(metadata, dict) else None
    if not isinstance(component, dict) or component.get("name") != MODULE:
        raise ValueError("SBOM metadata must identify the scanner Go module")
    if "timestamp" in metadata:
        raise ValueError("release SBOM must omit a nondeterministic timestamp")

    old_reference = component.get("bom-ref")
    old_purl = component.get("purl")
    reference = f"pkg:golang/{MODULE}@{version}?type=module"
    component["bom-ref"] = reference
    component["purl"] = reference
    component["version"] = version

    replacements = {
        value: reference for value in (old_reference, old_purl) if isinstance(value, str) and value
    }

    def replace_references(value: Any) -> Any:
        if isinstance(value, str):
            return replacements.get(value, value)
        if isinstance(value, list):
            return [replace_references(item) for item in value]
        if isinstance(value, dict):
            return {key: replace_references(item) for key, item in value.items()}
        return value

    document = replace_references(document)
    metadata = document["metadata"]
    properties = metadata.setdefault("properties", [])
    if not isinstance(properties, list):
        raise ValueError("SBOM metadata properties must be an array")
    for item in properties:
        if not isinstance(item, dict) or not isinstance(item.get("name"), str):
            raise ValueError("SBOM metadata properties must contain named objects")
        if item["name"] in {"prc:build:version", "prc:source:commit"}:
            raise ValueError("SBOM input cannot predeclare scanner-owned release identity")
    properties.extend(
        [
            {"name": "prc:build:version", "value": version},
            {"name": "prc:source:commit", "value": commit},
        ]
    )
    properties.sort(key=lambda item: (str(item.get("name", "")), str(item.get("value", ""))))
    return (json.dumps(document, indent=2, sort_keys=True, ensure_ascii=False) + "\n").encode("utf-8")


def release_support_files() -> list[tuple[str, bytes, int]]:
    files: list[pathlib.Path] = [
        ROOT / "LICENSE",
        ROOT / "README.md",
        ROOT / "THIRD_PARTY_NOTICES.md",
    ]
    for directory in (
        "adapters",
        "catalog",
        "docs",
        "fixtures/benchmarks",
        "packs",
        "schemas",
    ):
        files.extend(sorted((ROOT / directory).rglob("*")))
    result: list[tuple[str, bytes, int]] = []
    total_bytes = 0
    for path in files:
        if path.is_symlink():
            raise ValueError(f"release support files cannot contain symlinks: {path.relative_to(ROOT)}")
        if path.is_dir():
            continue
        if not path.is_file():
            raise ValueError(f"release support entry is not a regular file: {path.relative_to(ROOT)}")
        data = path.read_bytes()
        total_bytes += len(data)
        if len(result) >= MAX_RELEASE_INPUT_FILES or total_bytes > MAX_RELEASE_INPUT_BYTES:
            raise ValueError("release support files exceed the bounded packaging limits")
        result.append((path.relative_to(ROOT).as_posix(), data, 0o644))
    return result


def create_tar_gz(path: pathlib.Path, root_name: str, entries: list[tuple[str, bytes, int]], epoch: int) -> None:
    with path.open("wb") as output:
        with gzip.GzipFile(filename="", mode="wb", fileobj=output, compresslevel=9, mtime=epoch) as compressed:
            with tarfile.open(fileobj=compressed, mode="w", format=tarfile.USTAR_FORMAT) as archive:
                for name, data, mode in sorted(entries):
                    information = tarfile.TarInfo(f"{root_name}/{name}")
                    information.size = len(data)
                    information.mtime = epoch
                    information.mode = mode
                    information.uid = 0
                    information.gid = 0
                    information.uname = ""
                    information.gname = ""
                    archive.addfile(information, io.BytesIO(data))


def create_zip(path: pathlib.Path, root_name: str, entries: list[tuple[str, bytes, int]], timestamp: dt.datetime) -> None:
    earliest = dt.datetime(1980, 1, 1, tzinfo=dt.timezone.utc)
    archive_time = max(timestamp, earliest)
    date_time = (
        archive_time.year,
        archive_time.month,
        archive_time.day,
        archive_time.hour,
        archive_time.minute,
        archive_time.second - archive_time.second % 2,
    )
    with zipfile.ZipFile(path, mode="w", compression=zipfile.ZIP_DEFLATED, compresslevel=9) as archive:
        for name, data, mode in sorted(entries):
            information = zipfile.ZipInfo(f"{root_name}/{name}", date_time=date_time)
            information.compress_type = zipfile.ZIP_DEFLATED
            information.create_system = 3
            information.external_attr = (stat.S_IFREG | mode) << 16
            archive.writestr(information, data)


def host_target() -> tuple[str, str]:
    operating_system = {"Linux": "linux", "Darwin": "darwin", "Windows": "windows"}.get(platform.system())
    architecture = {"x86_64": "amd64", "AMD64": "amd64", "aarch64": "arm64", "arm64": "arm64"}.get(platform.machine())
    if not operating_system or not architecture:
        raise RuntimeError("release verification requires a supported host OS and architecture")
    return operating_system, architecture


def build_release(args: argparse.Namespace) -> None:
    if len(args.version) > 128 or not SEMVER.fullmatch(args.version):
        raise ValueError("version must be semantic version X.Y.Z with an optional prerelease")
    if not COMMIT.fullmatch(args.commit):
        raise ValueError("commit must be a lowercase 40- or 64-character hexadecimal revision")
    built_at, epoch = parse_timestamp(args.built_at)
    timestamp = dt.datetime.fromtimestamp(epoch, tz=dt.timezone.utc)
    output = args.output.resolve()
    if output.exists():
        raise ValueError("output directory already exists")
    output.parent.mkdir(parents=True, exist_ok=True)
    if args.sbom.resolve() == output or output in args.sbom.resolve().parents:
        raise ValueError("SBOM input cannot be inside the output directory")

    staging = pathlib.Path(tempfile.mkdtemp(prefix=".prc-release-", dir=output.parent))
    try:
        binaries = staging / "binaries"
        distribution = staging / "distribution"
        binaries.mkdir(mode=0o700)
        distribution.mkdir(mode=0o700)
        ldflags = " ".join(
            (
                "-s",
                "-w",
                f"-X=main.version={args.version}",
                f"-X=main.revision={args.commit}",
                f"-X=main.builtAt={built_at}",
            )
        )
        base_environment = os.environ.copy()
        base_environment.update({"CGO_ENABLED": "0", "SOURCE_DATE_EPOCH": str(epoch)})
        built: dict[tuple[str, str], pathlib.Path] = {}
        for goos, goarch in TARGETS:
            environment = base_environment | {"GOOS": goos, "GOARCH": goarch}
            suffix = ".exe" if goos == "windows" else ""
            binary = binaries / f"prc_{goos}_{goarch}{suffix}"
            subprocess.run(
                [
                    args.go,
                    "build",
                    "-trimpath",
                    "-buildvcs=false",
                    "-ldflags",
                    ldflags,
                    "-o",
                    str(binary),
                    "./cmd/prc",
                ],
                cwd=ROOT,
                env=environment,
                check=True,
            )
            built[(goos, goarch)] = binary

        host_binary = built[host_target()]
        version_information = run_json([str(host_binary), "version", "--format", "json"])
        expected_version = {
            "schema_version": "prc.version/v0.1",
            "version": args.version,
            "revision": args.commit,
            "built_at": built_at,
        }
        if any(version_information.get(key) != value for key, value in expected_version.items()):
            raise RuntimeError("built scanner did not report the requested release identity")
        go_version = version_information.get("go_version")
        if not isinstance(go_version, str) or not go_version.startswith("go1."):
            raise RuntimeError("built scanner did not report a Go toolchain version")

        catalog = run_json([str(host_binary), "catalog", "validate", "--catalog-root", str(ROOT), "--format", "json"])
        packs: list[dict[str, Any]] = []
        for pack_path in sorted((ROOT / "packs").glob("*.yaml")):
            report = run_json(
                [
                    str(host_binary),
                    "pack",
                    "validate",
                    "--catalog-root",
                    str(ROOT),
                    "--file",
                    str(pack_path),
                    "--format",
                    "json",
                ]
            )
            manifest = report.get("manifest")
            if not isinstance(manifest, dict) or not isinstance(manifest.get("id"), str):
                raise RuntimeError(f"pack validation returned no manifest identity: {pack_path.name}")
            packs.append(
                {
                    "id": manifest["id"],
                    "path": pack_path.relative_to(ROOT).as_posix(),
                    "digest": report["digest"],
                    "suite_digest": report["suite_digest"],
                    "benchmark_corpus_digest": report["benchmark_corpus_digest"],
                    "catalog_digest": report["catalog_digest"],
                }
            )

        sbom_name = f"prc_{args.version}.cdx.json"
        sbom_path = distribution / sbom_name
        sbom_path.write_bytes(normalized_sbom(args.sbom, args.version, args.commit))
        support = release_support_files()
        artifacts: list[dict[str, Any]] = []
        for goos, goarch in TARGETS:
            suffix = ".exe" if goos == "windows" else ""
            entries = [(f"prc{suffix}", built[(goos, goarch)].read_bytes(), 0o755), *support]
            stem = f"prc_{args.version}_{goos}_{goarch}"
            root_name = stem
            if goos == "windows":
                archive_path = distribution / f"{stem}.zip"
                create_zip(archive_path, root_name, entries, timestamp)
                archive_format = "zip"
            else:
                archive_path = distribution / f"{stem}.tar.gz"
                create_tar_gz(archive_path, root_name, entries, epoch)
                archive_format = "tar.gz"
            identity = file_identity(archive_path)
            artifacts.append(
                {
                    "name": archive_path.name,
                    "goos": goos,
                    "goarch": goarch,
                    "archive_format": archive_format,
                    **identity,
                }
            )

        npm_staging = staging / "npm"
        npm_packages = npm_distribution.build_packages(
            version=args.version,
            commit=args.commit,
            built_at=built_at,
            epoch=epoch,
            binaries=built,
            support=support,
            output=npm_staging,
        )
        for package in npm_packages:
            os.replace(npm_staging / package["name"], distribution / package["name"])

        manifest = {
            "schema_version": "prc.release-manifest/v0.2",
            "product": "prc-scanner",
            "version": args.version,
            "source_commit": args.commit,
            "built_at": built_at,
            "go_version": go_version,
            "catalog": catalog,
            "packs": packs,
            "artifacts": artifacts,
            "npm_packages": npm_packages,
            "sbom": {
                "name": sbom_name,
                "format": "CycloneDX",
                "spec_version": "1.6",
                **file_identity(sbom_path),
            },
        }
        manifest_path = distribution / f"prc_{args.version}_release-manifest.json"
        manifest_path.write_text(json.dumps(manifest, indent=2, sort_keys=True) + "\n", encoding="utf-8")

        write_checksums(distribution)
        os.replace(distribution, output)
    except BaseException:
        shutil.rmtree(staging, ignore_errors=True)
        raise
    shutil.rmtree(staging, ignore_errors=True)


def parser() -> argparse.ArgumentParser:
    result = argparse.ArgumentParser(description=__doc__)
    result.add_argument("--version", required=True, help="scanner semantic version without a v prefix")
    result.add_argument("--commit", required=True, help="exact lowercase source revision")
    result.add_argument("--built-at", required=True, help="reproducible RFC 3339 source timestamp")
    result.add_argument("--sbom", required=True, type=pathlib.Path, help="timestamp-free CycloneDX 1.6 module SBOM")
    result.add_argument("--output", required=True, type=pathlib.Path, help="new output directory")
    result.add_argument("--go", default="go", help="Go toolchain executable")
    return result


def main() -> int:
    try:
        build_release(parser().parse_args())
    except (OSError, ValueError, RuntimeError, subprocess.CalledProcessError) as error:
        print(f"release build failed: {error}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
