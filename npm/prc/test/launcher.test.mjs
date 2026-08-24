import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import { chmod, cp, mkdir, readFile, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { spawnSync } from "node:child_process";
import { mkdtemp } from "node:fs/promises";
import test from "node:test";

const SOURCE = join(dirname(fileURLToPath(import.meta.url)), "..", "bin", "everylast.js");
const VERSION = "1.2.3-test.1";
const COMMIT = "a".repeat(40);
const PLATFORM_NAMES = {
  "darwin-arm64": "@marinjursic/everylast-darwin-arm64",
  "darwin-x64": "@marinjursic/everylast-darwin-x64",
  "linux-arm64": "@marinjursic/everylast-linux-arm64",
  "linux-x64": "@marinjursic/everylast-linux-x64",
  "win32-arm64": "@marinjursic/everylast-windows-arm64",
  "win32-x64": "@marinjursic/everylast-windows-x64"
};

function digest(value) {
  return createHash("sha256").update(value).digest("hex");
}

async function fixture({ corruptBinary = false, omitPlatform = false } = {}) {
  const root = await mkdtemp(join(tmpdir(), "everylast-npm-test-"));
  const scope = join(root, "node_modules", "@marinjursic");
  const launcher = join(scope, "everylast");
  await mkdir(join(launcher, "bin"), { recursive: true });
  await cp(SOURCE, join(launcher, "bin", "everylast.js"));
  await chmod(join(launcher, "bin", "everylast.js"), 0o755);
  await writeFile(join(launcher, "package.json"), JSON.stringify({ name: "@marinjursic/everylast", version: VERSION, type: "module" }));

  const key = `${process.platform}-${process.arch}`;
  const packageName = PLATFORM_NAMES[key];
  assert.ok(packageName, `test host ${key} is supported`);
  const binaryName = process.platform === "win32" ? "everylast.exe" : "everylast";
  const platformRoot = join(scope, packageName.slice("@marinjursic/".length));
  const resultPath = join(root, "arguments.json");
  let binary = Buffer.from(`#!/usr/bin/env node\nrequire("node:fs").writeFileSync(process.env.PRC_TEST_RESULT, JSON.stringify(process.argv.slice(2))); process.exit(2);\n`);
  const manifest = {
    schema_version: "prc.npm-platform/v0.1", package_name: packageName, version: VERSION,
    source_commit: COMMIT, built_at: "2026-08-23T10:00:00Z", binary_path: `bin/${binaryName}`,
    binary_sha256: digest(binary)
  };
  const manifestBytes = Buffer.from(`${JSON.stringify(manifest)}\n`);
  const bindings = {};
  for (const [platform, name] of Object.entries(PLATFORM_NAMES)) {
    bindings[platform] = { package_name: name, manifest_sha256: platform === key ? digest(manifestBytes) : "0".repeat(64) };
  }
  await writeFile(join(launcher, "platforms.json"), JSON.stringify({ schema_version: "prc.npm-platforms/v0.1", version: VERSION, platforms: bindings }));
  if (!omitPlatform) {
    await mkdir(join(platformRoot, "bin"), { recursive: true });
    await writeFile(join(platformRoot, "package.json"), JSON.stringify({ name: packageName, version: VERSION }));
    await writeFile(join(platformRoot, "manifest.json"), manifestBytes);
    if (corruptBinary) binary = Buffer.concat([binary, Buffer.from("// changed\n")]);
    await writeFile(join(platformRoot, "bin", binaryName), binary);
    await chmod(join(platformRoot, "bin", binaryName), 0o755);
  }
  return { launcher: join(launcher, "bin", "everylast.js"), resultPath };
}

test("launcher verifies the package and forwards literal arguments and exit code", async () => {
  const item = await fixture();
  const argumentsToForward = ["scan", "space path", "; touch SHOULD_NOT_EXIST", "$(uname)", "line\nbreak"];
  const result = spawnSync(process.execPath, [item.launcher, ...argumentsToForward], {
    encoding: "utf8", env: { ...process.env, PRC_TEST_RESULT: item.resultPath }
  });
  assert.equal(result.status, 2, result.stderr);
  assert.deepEqual(JSON.parse(await readFile(item.resultPath, "utf8")), argumentsToForward);
});

test("launcher fails closed when the binary digest changes", async () => {
  const item = await fixture({ corruptBinary: true });
  const result = spawnSync(process.execPath, [item.launcher, "scan", "."], { encoding: "utf8" });
  assert.equal(result.status, 3);
  assert.match(result.stderr, /does not match its release manifest/);
});

test("launcher has no PATH or network fallback for a missing platform package", async () => {
  const item = await fixture({ omitPlatform: true });
  const result = spawnSync(process.execPath, [item.launcher, "scan", "."], { encoding: "utf8" });
  assert.equal(result.status, 3);
  assert.match(result.stderr, /required platform package/);
});
