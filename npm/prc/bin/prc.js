#!/usr/bin/env node

import { createHash } from "node:crypto";
import { closeSync, createReadStream, fstatSync, lstatSync, openSync, readSync, realpathSync } from "node:fs";
import { dirname, isAbsolute, join, relative, sep } from "node:path";
import { fileURLToPath } from "node:url";
import { spawn } from "node:child_process";

const MAX_JSON_BYTES = 1024 * 1024;
const MAX_BINARY_BYTES = 64 * 1024 * 1024;
const MAX_SUPPORT_FILES = 512;
const MAX_SUPPORT_BYTES = 24 * 1024 * 1024;
const MAX_SUPPORT_FILE_BYTES = 16 * 1024 * 1024;
const DIGEST = /^[0-9a-f]{64}$/;
const VERSION = /^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-[0-9A-Za-z.-]+)?$/;
const COMMIT = /^(?:[0-9a-f]{40}|[0-9a-f]{64})$/;
const EXPECTED_PACKAGES = Object.freeze({
  "darwin-arm64": "@marinjursic/prc-darwin-arm64",
  "darwin-x64": "@marinjursic/prc-darwin-x64",
  "linux-arm64": "@marinjursic/prc-linux-arm64",
  "linux-x64": "@marinjursic/prc-linux-x64",
  "win32-arm64": "@marinjursic/prc-windows-arm64",
  "win32-x64": "@marinjursic/prc-windows-x64"
});

function safeText(value) {
  let result = "";
  for (const character of String(value)) {
    const code = character.codePointAt(0);
    const unsafe = code <= 0x1f || (code >= 0x7f && code <= 0x9f) ||
      code === 0x061c || code === 0x200e || code === 0x200f ||
      (code >= 0x202a && code <= 0x202e) || (code >= 0x2066 && code <= 0x2069);
    result += unsafe ? `\\u${code.toString(16).toUpperCase().padStart(code <= 0xffff ? 4 : 8, "0")}` : character;
  }
  return result;
}

function fail(message) {
  process.stderr.write(`prc npm launcher: ${safeText(message)}\n`);
  process.exitCode = 3;
}

function readBoundedJSON(path, label) {
  return JSON.parse(readBoundedFile(path, label).toString("utf8"));
}

function readBoundedFile(path, label) {
  const before = lstatSync(path);
  if (!before.isFile() || before.isSymbolicLink() || before.size > MAX_JSON_BYTES) {
    throw new Error(`${label} is not a bounded regular file`);
  }
  const descriptor = openSync(path, "r");
  try {
    const opened = fstatSync(descriptor);
    if (!sameFile(before, opened)) throw new Error(`${label} changed while it was opened`);
    const data = Buffer.alloc(before.size + 1);
    let total = 0;
    while (total < data.length) {
      const count = readSync(descriptor, data, total, data.length - total, null);
      if (count === 0) break;
      total += count;
    }
    const afterOpen = fstatSync(descriptor);
    const afterPath = lstatSync(path);
    if (total !== before.size || !sameFile(before, afterOpen) || !sameFile(before, afterPath)) {
      throw new Error(`${label} changed while it was read`);
    }
    return data.subarray(0, total);
  } finally {
    closeSync(descriptor);
  }
}

function sameFile(before, after) {
  return after.isFile() && !after.isSymbolicLink() && before.dev === after.dev && before.ino === after.ino &&
    before.size === after.size && before.mtimeMs === after.mtimeMs;
}

function exactKeys(value, expected, label) {
  if (!value || typeof value !== "object" || Array.isArray(value)) {
    throw new Error(`${label} is not an object`);
  }
  const actual = Object.keys(value).sort();
  const wanted = [...expected].sort();
  if (actual.length !== wanted.length || actual.some((key, index) => key !== wanted[index])) {
    throw new Error(`${label} has unexpected fields`);
  }
}

async function sha256File(path, maximumBytes) {
  const before = lstatSync(path);
  if (!before.isFile() || before.isSymbolicLink() || before.size > maximumBytes) {
    throw new Error("a packaged scanner file is not a bounded regular file");
  }
  const hash = createHash("sha256");
  let total = 0;
  for await (const chunk of createReadStream(path)) {
    total += chunk.length;
    if (total > maximumBytes) throw new Error("a packaged scanner file exceeded its byte limit while it was verified");
    hash.update(chunk);
  }
  const after = lstatSync(path);
  if (total !== before.size || !sameFile(before, after)) {
    throw new Error("a packaged scanner file changed while it was verified");
  }
  return { digest: hash.digest("hex"), size: after.size };
}

function inside(root, path) {
  const value = relative(root, path);
  return value !== "" && value !== ".." && !value.startsWith(`..${sep}`) && !isAbsolute(value);
}

function safeSupportPath(value) {
  if (typeof value !== "string" || value.length < 1 || value.length > 512 || value.includes("\\")) return false;
  const segments = value.split("/");
  if (segments.some((segment) => segment === "" || segment === "." || segment === "..")) return false;
  const allowed = [
    "bin/adapters/", "bin/catalog/", "bin/docs/checklists/",
    "bin/docs/engineering/", "bin/fixtures/benchmarks/", "bin/packs/", "bin/schemas/"
  ];
  return value === "bin/THIRD_PARTY_NOTICES.md" || allowed.some(
    (prefix) => value.length > prefix.length && value.startsWith(prefix)
  );
}

async function main() {
  const launcherRoot = realpathSync(join(dirname(fileURLToPath(import.meta.url)), ".."));
  const launcherPackage = readBoundedJSON(join(launcherRoot, "package.json"), "launcher package metadata");
  if (launcherPackage.name !== "@marinjursic/prc" || !VERSION.test(launcherPackage.version)) {
    throw new Error("launcher package identity is invalid");
  }
  const platformDocument = readBoundedJSON(join(launcherRoot, "platforms.json"), "platform binding");
  exactKeys(platformDocument, ["schema_version", "version", "platforms"], "platform binding");
  if (platformDocument.schema_version !== "prc.npm-platforms/v0.2" || platformDocument.version !== launcherPackage.version) {
    throw new Error("platform binding does not match the launcher version");
  }
  exactKeys(platformDocument.platforms, Object.keys(EXPECTED_PACKAGES), "platform map");
  for (const [key, packageName] of Object.entries(EXPECTED_PACKAGES)) {
    const binding = platformDocument.platforms[key];
    exactKeys(binding, ["package_name", "manifest_sha256"], `platform binding ${key}`);
    if (binding.package_name !== packageName || !DIGEST.test(binding.manifest_sha256)) {
      throw new Error(`platform binding ${key} is invalid`);
    }
  }

  const key = `${process.platform}-${process.arch}`;
  const packageName = EXPECTED_PACKAGES[key];
  if (!packageName) {
    throw new Error(`unsupported platform ${key}; supported platforms are ${Object.keys(EXPECTED_PACKAGES).join(", ")}`);
  }
  const packageDirectoryName = packageName.slice("@marinjursic/".length);
  const scopeDirectory = dirname(launcherRoot);
  const candidates = [
    join(scopeDirectory, packageDirectoryName),
    join(launcherRoot, "node_modules", "@marinjursic", packageDirectoryName)
  ];
  let platformRoot;
  for (const candidate of candidates) {
    try {
      platformRoot = realpathSync(candidate);
      break;
    } catch (error) {
      if (error?.code !== "ENOENT") throw error;
    }
  }
  if (!platformRoot) {
    throw new Error(`required platform package ${packageName}@${launcherPackage.version} is not installed; reinstall the exact launcher version with optional dependencies enabled`);
  }
  const packageMetadata = readBoundedJSON(join(platformRoot, "package.json"), "platform package metadata");
  if (packageMetadata.name !== packageName || packageMetadata.version !== launcherPackage.version) {
    throw new Error("platform package identity does not match the launcher");
  }
  const manifestPath = join(platformRoot, "manifest.json");
  const manifestBytes = readBoundedFile(manifestPath, "platform package manifest");
  if (manifestBytes.length > MAX_JSON_BYTES || createHash("sha256").update(manifestBytes).digest("hex") !== platformDocument.platforms[key].manifest_sha256) {
    throw new Error("platform package manifest does not match the launcher binding");
  }
  const manifest = JSON.parse(manifestBytes.toString("utf8"));
  exactKeys(manifest, [
    "schema_version", "package_name", "version", "source_commit", "built_at", "binary_path", "binary_sha256",
    "support_file_count", "support_bytes", "support_files"
  ], "platform package manifest");
  const expectedBinary = process.platform === "win32" ? "bin/prc.exe" : "bin/prc";
  if (manifest.schema_version !== "prc.npm-platform/v0.2" || manifest.package_name !== packageName ||
      manifest.version !== launcherPackage.version || !COMMIT.test(manifest.source_commit) ||
      typeof manifest.built_at !== "string" || manifest.binary_path !== expectedBinary || !DIGEST.test(manifest.binary_sha256)) {
    throw new Error("platform package manifest identity is invalid");
  }
  const binaryPath = realpathSync(join(platformRoot, ...manifest.binary_path.split("/")));
  if (!inside(platformRoot, binaryPath)) throw new Error("packaged scanner binary escapes its package");
  const binaryIdentity = await sha256File(binaryPath, MAX_BINARY_BYTES);
  if (binaryIdentity.digest !== manifest.binary_sha256) {
    throw new Error("packaged scanner binary does not match its release manifest");
  }
  if (!Array.isArray(manifest.support_files) || !Number.isSafeInteger(manifest.support_file_count) ||
      manifest.support_file_count < 1 || manifest.support_file_count > MAX_SUPPORT_FILES ||
      manifest.support_files.length !== manifest.support_file_count || !Number.isSafeInteger(manifest.support_bytes) ||
      manifest.support_bytes < 1 || manifest.support_bytes > MAX_SUPPORT_BYTES) {
    throw new Error("platform package support-file inventory is invalid");
  }
  let supportBytes = 0;
  let previousSupportPath = "";
  for (const [index, support] of manifest.support_files.entries()) {
    exactKeys(support, ["path", "size", "sha256"], `platform support file ${index}`);
    if (!safeSupportPath(support.path) || support.path <= previousSupportPath || !DIGEST.test(support.sha256) ||
        !Number.isSafeInteger(support.size) || support.size < 0 || support.size > MAX_SUPPORT_FILE_BYTES) {
      throw new Error(`platform support file ${index} is invalid`);
    }
    previousSupportPath = support.path;
    supportBytes += support.size;
    if (supportBytes > MAX_SUPPORT_BYTES) throw new Error("platform package support files exceed their byte limit");
    const supportPath = realpathSync(join(platformRoot, ...support.path.split("/")));
    if (!inside(platformRoot, supportPath)) throw new Error("platform support file escapes its package");
    const identity = await sha256File(supportPath, MAX_SUPPORT_FILE_BYTES);
    if (identity.size !== support.size || identity.digest !== support.sha256) {
      throw new Error(`packaged support file ${safeText(support.path)} does not match its release manifest`);
    }
  }
  if (supportBytes !== manifest.support_bytes) {
    throw new Error("platform package support byte count does not match its release manifest");
  }

  const child = spawn(binaryPath, process.argv.slice(2), { stdio: "inherit", shell: false, windowsHide: true });
  const forwarded = new Set();
  for (const signal of ["SIGINT", "SIGTERM", "SIGHUP"]) {
    process.once(signal, () => {
      forwarded.add(signal);
      if (!child.killed) child.kill(signal);
    });
  }
  child.once("error", (error) => fail(`could not start the packaged scanner: ${error.message}`));
  child.once("exit", (code, signal) => {
    if (signal) {
      if (!forwarded.has(signal)) process.kill(process.pid, signal);
      else process.exitCode = 128 + ({ SIGHUP: 1, SIGINT: 2, SIGTERM: 15 }[signal] ?? 1);
    } else {
      process.exitCode = code ?? 4;
    }
  });
}

main().catch((error) => fail(error instanceof Error ? error.message : error));
