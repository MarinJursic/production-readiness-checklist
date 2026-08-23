# npm distribution, terminal output, and safe-scan plan

**Status:** implemented and tested; public npm publication still requires maintainer registry setup

**Original reviewed repository revision:** `c95b152ac2f9dcc2558960ed92d7f64b6494da76`

**Research date:** 2026-08-23  
**Scope:** installation, launch, read-only scanning, terminal output, report creation, and the path from 10,042 written controls to honest machine checks

## Decision

Yes, the scanner can be made easy to install and run with npm.

The recommended user flow is:

```sh
npm install --save-dev --save-exact --ignore-scripts @marinjursic/prc@0.2.0
npx prc scan .
```

The package name was checked against the public npm registry on 2026-08-23. The launcher and six platform names returned `E404`, meaning no public package was visible. That does not prove the maintainer owns the `@marinjursic` npm scope. Public instructions therefore continue to say “after publication” until the maintainer configures the scope and trusted publisher.

The repository now builds the exact package design below as seven deterministic release tarballs. Local offline installation with `--ignore-scripts` was tested end to end against the packaged native scanner. No package was published as part of this implementation.

For a project that wants a short repeatable command:

```json
{
  "scripts": {
    "scan": "prc scan ."
  }
}
```

Then the normal command is:

```sh
npm run scan
```

`npm scan` cannot be added as a custom npm command. `npm run scan` is npm's standard custom-script form. A user can also run `npx prc scan .` after the package is installed.

The npm package must be an optional way to get the same native `prc` scanner. The direct release archive must remain available for Go, Python, Java, Rust, infrastructure, air-gapped, and non-Node users.

## Non-negotiable safety decision

Do not build an npm package that downloads a binary in `postinstall`.

An install-time downloader is simple to write but creates the wrong trust model:

- npm normally allows package lifecycle scripts unless the user disables or restricts them;
- an install script can read files, use the network, start child processes, and change the project;
- corporate and security-aware users often disable install scripts;
- a second download is not protected by the package lock's tarball integrity;
- proxy, redirect, release-account, DNS, and partial-download failures add more unsafe cases;
- the package would fail in offline and restricted environments.

The npm documentation confirms that `ignore-scripts` disables package scripts and that current npm versions support explicit install-script allow rules. The scanner package should work with `--ignore-scripts`, not ask users to allow it.

## Package design

### Package set

Use one small launcher package and one package for each supported platform:

```text
@marinjursic/prc
@marinjursic/prc-darwin-arm64
@marinjursic/prc-darwin-x64
@marinjursic/prc-linux-arm64
@marinjursic/prc-linux-x64
@marinjursic/prc-windows-arm64
@marinjursic/prc-windows-x64
```

The public registry lookup found no visible package under these names. Scope
ownership and trusted-publisher access still need to be confirmed by the
maintainer. Do not publish placeholder packages merely to hold names.

The launcher package should contain only:

- `package.json`;
- a small `bin/prc.js` launcher;
- `README.md`;
- `LICENSE`;
- exact-version optional dependencies on the six platform packages.

Each platform package should contain only:

- its exact native `prc` binary;
- the runtime catalog beside the binary;
- required packs, adapter manifests, schemas, and notices;
- a minimal `package.json` with exact `os` and `cpu` limits;
- a build manifest containing the scanner version, source commit, and SHA-256 hashes.

The current binary finds its catalog beside its executable. Therefore the platform package should use a layout such as:

```text
@marinjursic/prc-darwin-arm64/
└── bin/
    ├── prc
    ├── catalog/
    ├── adapters/
    ├── packs/
    └── schemas/
```

The launcher resolves the exact platform package, finds the binary by an absolute path relative to that installed package, validates the expected version and packaged hash manifest, and launches it with an argument array. It must not search `PATH`, build a shell string, use `eval`, download a fallback, or choose a package from target-project data.

### Package rules

The release tests must reject a package when any of these is true:

- `preinstall`, `install`, `postinstall`, `prepare`, or any other install-time hook exists;
- a runtime dependency exists in the launcher;
- an optional platform dependency uses a version range instead of the exact release version;
- an unknown file is included in a tarball;
- a symlink, device, socket, executable script, source map with local paths, secret, credential, or private key is included unexpectedly;
- the binary version, source commit, or hash differs from the release manifest;
- the launcher can fall back to a network download or a binary found on `PATH`;
- the package does not work when installed with `--ignore-scripts`;
- the package can load JavaScript or configuration from the project being scanned.

### Node support

As of this review, Node 24 is Active LTS and Node 22 is Maintenance LTS. The
launcher uses `engines.node` set to `>=22.14.0`, which includes npm trusted
publishing's current Node floor. Test both supported LTS lines. Recheck the
official Node release table at each scanner release; do not keep an end-of-life
Node version merely because the wrapper still happens to run.

The native binary does the scan. Node is only the launcher, so the JavaScript should use a very small set of stable built-in APIs and no third-party packages.

## Safe user commands

### Recommended project install

```sh
npm install --save-dev --save-exact --ignore-scripts @marinjursic/prc@0.2.0
npx prc scan .
```

Why each option exists:

- `--save-dev` keeps the scanner out of application runtime dependencies.
- `--save-exact` avoids silently moving to a new scanner version.
- `--ignore-scripts` prevents install-time lifecycle code from this and other newly installed packages from running.
- the exact `@0.2.0` makes the first install clear and repeatable.

The command still changes `package.json`, `package-lock.json`, and `node_modules`. The documentation must say that clearly.

### One-time run without adding a dependency

```sh
npm exec --yes --ignore-scripts --package=@marinjursic/prc@0.2.0 -- prc scan .
```

This is convenient but it downloads code into npm's cache and executes it. It is not safer than a reviewed, exact project dependency. Show it as the quick trial, not the strongest repeatable setup.

### Repeatable project command

After the exact package is installed, add:

```json
{
  "scripts": {
    "scan": "prc scan ."
  }
}
```

Then run:

```sh
npm run scan
```

### Stronger verification after install

```sh
npm audit signatures
npx prc version
npx prc scan .
```

`npm audit signatures` checks registry signatures and available provenance attestations. Provenance says where and how a package was built; it does not prove that the code is harmless. The README must say this plainly.

## Implemented terminal experience

### Main design

The terminal shows every selected assertion in a short line for the current
40-check core profile. When future profiles become large, `--summary` can show
only totals and items that need attention, while `--all-results` shows every
result.

Color is useful, but color must never be the only meaning:

| Result | Symbol | Label | TTY color | Meaning |
| --- | --- | --- | --- | --- |
| Pass | `✓` | `PASS` | green | The exact assertion has acceptable evidence. |
| Fail | `✗` | `FAIL` | red | The assertion was checked and failed. |
| Blocked or unknown | `!` | `BLOCKED` | yellow | Required proof or safe execution is missing. |
| Manual review | `?` | `MANUAL` | blue or yellow | A person or external authority must review it. |
| Not applicable | `–` | `N/A` | dim/gray | The trigger was proven absent and a reason was recorded. |
| Execution error | `×` | `ERROR` | red | The check could not complete; this is not a Pass. |

Rules for color:

- default to `--color auto`;
- use color only when standard output is a real terminal;
- support `--color always` and `--color never`;
- honor a nonempty `NO_COLOR` environment variable;
- disable color for `TERM=dumb`;
- never put ANSI bytes in JSON, Markdown, SARIF, JUnit, redirected human output, log files, or reports;
- sanitize all control characters from target names, paths, tool messages, and finding text before terminal rendering;
- do not emit terminal hyperlinks by default;
- keep the word label and symbol even when color is off.

The Go `x/term` package supplies terminal detection. The `NO_COLOR` convention supplies a simple user-wide opt-out.

### Example output

```text
Production Readiness Scanner 0.2.0

Target   invoice-api
Profile  prc/core-repository@1.0
Mode     scan only — no fixes and no project scripts

Checking 40 assertions...

  ✓ PASS     README is present
  ✓ PASS     License is present
  ✗ FAIL     Dependency resolution is reproducible
             package.json was found, but no matching lock file was found
             Control: USEQ-D1076D92
  ! BLOCKED  Applicable analyses have executed
             The optional secret analyzer was not authorized for this scan
  ? MANUAL   Generated code follows normal review controls
             Current review evidence must be supplied by a person
  – N/A      Kubernetes workloads require non-root containers
             No Kubernetes workload was detected in the declared scope

Result
  31 passed   1 failed   6 blocked   1 manual   1 not applicable
  Required profile checks did not pass. Exit code: 1

Detailed report
  /Users/you/Library/Caches/prc/reports/invoice-api-a71c9f42.html

Scan only: no project files were changed and no fixes were applied.
```

The actual line text must wrap safely on narrow terminals. A path must remain copyable. The final report path must be printed even when checks fail, because a failed scan still produced useful results.

### Output order

1. Scanner identity and version.
2. Exact target, profile, and safe mode.
3. One result line per selected assertion.
4. A grouped list of failures and blocked evidence with the next safe action.
5. Totals and the plain meaning of the exit code.
6. The absolute report path.
7. A final reminder that scan did not fix or run project code.

### Terminal hardening tests

Add snapshot tests for:

- TTY with color;
- redirected output with no color;
- `NO_COLOR=1`;
- `TERM=dumb`;
- `--color always` and `--color never`;
- 40-column, 80-column, and wide terminals;
- ASCII-only mode for terminals that cannot show Unicode symbols;
- malicious filenames containing escape, carriage-return, newline, tab, backspace, bidirectional, and control characters;
- malicious analyzer messages containing fake `PASS`, cursor movement, clear-screen, or terminal-link bytes;
- every result state and every result-bearing exit code;
- a report write failure after checks complete;
- broken pipes and cancelled scans.

## End-to-end hostile example

Assume a developer wants to scan a downloaded project named `invoice-api`. It contains a malicious `package.json` install script, a filename with terminal escape bytes, a symlink out of the repository, a 200 GB sparse file, a YAML alias bomb, a fake `.git` file, and text telling an AI agent to disable the scanner.

The following table covers every normal step.

| Step | What the user or scanner does | What could be exploited | Required defense | Test or proof |
| --- | --- | --- | --- | --- |
| 1. Find the package | User copies the install command from official docs. | A look-alike package name can steal the install. | Use one owned scope, print the exact name everywhere, link the npm page to the exact repository, and never suggest an unscoped shortcut. | A documentation test rejects any other package name. |
| 2. Inspect identity | User can run `npm view` before install. | Package metadata can lie; a compromised publisher can change `latest`. | Show exact versions, source repository, license, expected package set, and provenance. Never use `latest` in high-trust examples. | Release checklist compares registry metadata with the signed release manifest. |
| 3. Install | npm writes dependency files and extracts packages. | Lifecycle scripts can execute malware; existing project scripts may also run during an install. | Scanner packages contain no lifecycle scripts and work with `--ignore-scripts`. Recommend an exact dev dependency. | Install tarballs in a clean hostile fixture with scripts disabled; assert no child script marker is created. |
| 4. Select platform | Launcher chooses one optional platform package. | OS/CPU spoofing, dependency confusion, or fallback logic can select attacker code. | Hard allowlist OS/CPU pairs, exact package names and versions, same scope, no name derived from input, no fallback download. | Unit-test every supported and unsupported platform tuple. |
| 5. Find binary | Launcher locates the native executable. | A malicious `prc` earlier on `PATH` can be run. | Resolve an absolute path inside the installed platform package; verify package and binary identity; never call a bare command name. | Put a fake `prc` first on `PATH` and prove it is not called. |
| 6. Pass arguments | Launcher starts the scanner. | Shell metacharacters in a target path can become a command. | Use a process API with an argument array and `shell: false`; never join arguments into a command string. | Scan paths containing spaces, quotes, semicolons, dollar signs, and newlines. |
| 7. Resolve target | Scanner opens `invoice-api`. | Path traversal or a changing symlink can move scope outside the project. | Resolve and lock the root; reject unsafe root changes; use root-relative file handles; never follow target symlinks silently. | Race and symlink fixtures must end Blocked without reading outside the root. |
| 8. Inventory files | Scanner walks and hashes regular files. | Huge files, sparse files, deep trees, hard-link repetition, permission traps, and file races can cause denial of service or inconsistent evidence. | Keep entry limits, add per-file and total-byte budgets, elapsed-time and open-file limits, cancellation, and root-scoped open-after-stat checks. Report budget exhaustion, not Pass. | Adversarial tests cover huge, sparse, deep, changing, unreadable, hard-linked, and special files. |
| 9. Parse structured data | Scanner reads manifests, YAML, OpenAPI, and source syntax. | Parser bombs, huge nesting, aliases, duplicate keys, invalid encodings, and quadratic input can exhaust resources or change meaning. | Enforce byte, document, node, depth, alias, duplicate-key, and problem-count limits before semantic checks. Fail closed. | Parser fuzzing plus fixed alias, nesting, duplicate-key, and multi-document fixtures. |
| 10. Decide applicability | Scanner selects assertions. | A project can hide technology or claim Not Applicable in its own text. | Use bounded inventory facts and separately validated project declarations. A declaration is context, not proof. Undetermined applicability stays incomplete. | Tests try false declarations, excluded files, missing fields, and over-budget CEL. |
| 11. Run native checks | Scanner reads the project but does not run it. | A malicious `package.json`, Makefile, test, editor config, or prompt can ask the scanner or an agent to execute code or weaken rules. | Simple scan never runs package managers, tests, shells, hooks, plugins, project binaries, or agent instructions. Treat all repository text as data. | A fixture places marker-writing commands in every common script location; no marker may appear. |
| 12. Run optional analyzers | Only an explicit `verify-local` scan may run pinned OCI analyzers. | A compromised image, config injection, socket mount, network access, or output forgery can escape or lie. | Exact image digest and manifest, no network, read-only sealed snapshot, no host socket, bounded CPU/memory/PIDs/time/output, scanner-owned config, typed output, and fail-closed protocol checks. | Integration tests inspect the exact container plan and send malicious output. |
| 13. Render terminal | Scanner prints results. | A filename or analyzer message can clear the screen, forge a green PASS, alter clipboard data, or hide text. | Escape or replace control bytes, bidi controls, ANSI, OSC, CR, and untrusted newlines. Add color only after sanitizing data. | Golden tests contain known terminal-control attacks. |
| 14. Write report | Scanner creates a private HTML file outside the target. | HTML injection, symlink overwrite, path collision, permission leak, disk exhaustion, or placing the report inside the scanned tree can corrupt evidence. | Keep HTML auto-escaping, private `0700` directory and `0600` new file, no overwrite, no target directory, safe path components, bounded report size, and cleanup on partial failure. Check every parent against symlink races. | Tests cover HTML payloads, existing files, symlinked parents, full disk, permission denial, and target-inside-report paths. |
| 15. Exit | Scanner returns a result-bearing code. | Wrappers can treat fail/incomplete as a crash or erase the code. | Launcher forwards signals and the exact child exit code. Human output explains that `1` and `2` are assessment results. | Process tests cover all documented codes and signals. |
| 16. Open report | User opens the printed path. | A report could contain active remote content or leak data. | Standalone HTML contains no scripts, remote fonts, images, trackers, or network links created from untrusted content. Keep evidence summaries minimal and redact secrets. | Browser test confirms no network requests and safe escaping. |

## Current scanner safety review

The current code already has useful protections:

- `prc scan` is separate from `prc fix` and does not call the fix path;
- it does not run project dependency installation or target-project scripts;
- symlinks are inventoried but not followed;
- regular files are opened and compared with their earlier file identity before hashing;
- inventory entries, structured input, OpenAPI input, Go analysis input, sensitive-material input, adapter snapshots, adapter data, processes, output, and time have several explicit bounds;
- adapter execution needs the explicit `verify-local` mode and an exact profile binding;
- OCI adapters use exact image digests, no network, read-only snapshots, limited resources, scanner-owned configuration, and typed result binding;
- HTML uses Go's escaping template engine;
- default reports are outside the project, created privately, and never overwrite an existing file;
- missing or unsafe evidence remains blocked or unknown.

Remaining gaps before public npm beta:

1. **Total inventory bytes:** the inventory limits file count but hashes every regular-file byte. Add per-file, total-byte, elapsed-time, and open-file budgets to stop huge-file denial of service.
2. **Whole-walk path races:** some inventory reads still use path-based `WalkDir` plus later opens. Move all target reads to root-scoped handles or record a blocked result when identity changes.
3. **Report parent path:** the final report directory is checked, but every existing parent should be opened without following unsafe links before file creation.
4. **Registry publication:** configure ownership and npm trusted publishing, then publish all platform packages before the launcher. Never use a long-lived token when OIDC is available.
5. **Cross-platform install matrix:** the release contract builds and binds all six packages and runs the launcher tests on Linux. Add live clean-install jobs on macOS ARM64/x64 and Windows ARM64/x64 as runner availability permits.

## Review of every current assertion

This table checks whether the current pass rule is useful and whether it is too tied to one project shape.

| Assertion | Current rule | Review and planned change |
| --- | --- | --- |
| `PRC-A-CORE-001` | Nonempty root README with one of three names | Useful convention but name- and location-based, and any byte passes. Discover an operating guide by content and manifest links; keep root README as the default hint and allow a reviewed path. Check minimum required topics separately. |
| `PRC-A-CORE-002` | Nonempty root license with one of four names | Presence is too weak. Detect common names plus package metadata, identify license text/SPDX where possible, and keep compatibility review separate and manual. |
| `PRC-A-CORE-003` | Nonempty root contribution guide | Public collaborative projects need this; private projects may use another process. Detect links and configured paths, then check development/test/review/security guidance as separate assertions. |
| `PRC-A-CORE-004` | Nonempty security policy at two paths | Do not pass arbitrary text. Check reporting route, supported versions, response expectations, and safe contact information. Allow a reviewed external policy URL only with captured evidence. |
| `PRC-A-CORE-005` | Nonempty code of conduct file | Make it applicable to public or multi-party collaboration rather than every possible repository. Check an accepted content route, not only a filename. |
| `PRC-A-CORE-006` | Nonempty CODEOWNERS in three locations | Too tied to one hosting convention and does not prove durable ownership. Accept other ownership maps and signed team evidence; keep repository ownership and organizational resilience as separate assertions. |
| `PRC-A-CORE-007` | Every detected ecosystem has a lock/checksum file | Good goal but ecosystem presence is not enough to know build intent. Evaluate each runnable component; allow an evidence-backed reason for libraries or ecosystems with another reproducible mechanism. |
| `PRC-A-CORE-008` | Every external GitHub Action pinned to full commit SHA | Good strict GitHub profile rule. Keep it conditional. Also record the human-readable release behind the SHA and update ownership. |
| `PRC-A-CORE-009` | A supported CI file exists | Can miss custom CI and cannot prove CI ran. Add detector adapters and explicit configuration, then add a separate exact-commit CI-run evidence assertion. |
| `PRC-A-CORE-010` | Recognized test path or test command | This is the folder-layout problem the user warned about. Discover test frameworks and commands from parsers/build graphs; allow explicit reviewed commands; prove the command ran separately. Do not require `tests/`, `test/`, or another fixed layout. |
| `PRC-A-CORE-011` | GitHub workflow permissions are explicit | Useful but must evaluate effective workflow and job permissions by event. Presence alone is not least privilege. Keep GitHub-specific applicability. |
| `PRC-A-CORE-012` | Manual evidence for generated code review | Correctly manual, but not every project contains generated or AI-suggested code. Add an explicit trigger and retain manual evidence when it applies. |
| `PRC-A-CORE-013` | “Applicable analyses” with current evidence | The title overclaims because the current binding is only the pinned secret analyzer. Split secret, dependency, license, source, IaC, and other analysis classes into distinct assertions. |
| `PRC-A-CORE-014` | Every recognized source file ends in LF | This is a style check, not a universal production gate. Move it to an advisory style profile and respect generated/binary classifications and project policy. |
| `PRC-A-CORE-015` | Target resolves to a Git commit | Strong evidence for Git repositories, but archives and artifacts need another immutable subject digest. Generalize to exact subject identity with Git as one implementation. |
| `PRC-A-CORE-016` | Dependabot or Renovate file exists | Too vendor- and filename-specific, and presence does not prove operation. Accept other update systems and a documented manual path; verify recent update evidence separately. |
| `PRC-A-CORE-017` | Workflow is YAML with a mapping root | Good syntax signal but not full workflow validity. Rename it to the exact syntax claim and add provider schema validation separately. |
| `PRC-A-CORE-018` | Every workflow has a jobs map | Good structural GitHub check. Handle reusable and generated workflow forms according to the provider schema. |
| `PRC-A-CORE-019` | Every direct job has a timeout of at most 360 minutes | The hard six-hour value is one-size-fits-all. Keep a safe upper bound in a named profile and allow a lower project value; report reusable jobs clearly. |
| `PRC-A-CORE-020` | No `pull_request_target` | A safe hardened default, but some metadata-only designs can be secure. Keep the ban in the strict profile; later add a data-flow check for untrusted checkout, secrets, writes, and artifact trust. |
| `PRC-A-CORE-021` | No complete merge-conflict markers | Useful signal. Avoid false positives in fixtures and documentation by using file type/context and explicit reviewed test-fixture classification, not a broad ignored folder. |
| `PRC-A-CORE-022` | No group/other-writable local files | Local filesystem evidence is not portable repository policy and is weak on Windows. Make it local-workspace advisory evidence and assess packaged artifact permissions separately. |
| `PRC-A-CORE-023` | Manifests are nonempty | Too weak. A broken or comment-only manifest passes. Parse each supported manifest and tie it to a component. |
| `PRC-A-CORE-024` | Lock/checksum files are nonempty | Too weak. Parse them, check consistency with their manifest, and record unsupported formats as Blocked. |
| `PRC-A-CORE-025` | Each ecosystem has a known runtime declaration | Can force irrelevant files and miss unusual layouts. Evaluate runnable components, accept supported manifest/toolchain/config evidence, and allow reviewed external build-image identity. |
| `PRC-A-CORE-026` | Every container stage is `scratch` or digest-pinned | Strong supply-chain rule for a hardened profile. Keep dynamic or unresolved base identity Blocked. Other profiles may accept a separately resolved, signed base policy, but must not call a tag immutable. |
| `PRC-A-CORE-027` | Final stage declares a static non-root `USER` | Safe static rule but can flag a base image whose effective user is already non-root. Keep it for repository-only hardened scans; an image inspection assertion can prove effective user separately. |
| `PRC-A-CORE-028` | Every directory with Terraform has `.terraform.lock.hcl` | Incorrect for reusable child modules. Discover runnable root modules from Terraform structure and reviewed configuration, then require locks only where initialization selects providers. |
| `PRC-A-CORE-029` | Kubernetes workload explicitly avoids user zero | Useful hardened repository rule. A cluster admission policy may supply stronger environment evidence, but source intent and deployed enforcement must remain separate results. |
| `PRC-A-CORE-030` | Every Kubernetes container has requests and limits | Too universal. Requests and CPU/memory limits depend on workload and policy; use named profiles and per-resource requirements. Do not force identical resource fields on every project. |
| `PRC-A-K8S-001` | No host namespaces, HostProcess, hostPath, or privileged container | Strong restricted-workload profile, but system workloads can have reviewed needs. Keep explicit applicability and exception evidence; never silently ignore namespaces. |
| `PRC-A-K8S-002` | Linux containers set `allowPrivilegeEscalation: false` | Good restricted baseline. Keep Windows and other workload applicability explicit and show whether an admission policy supplies enforcement. |
| `PRC-A-K8S-003` | Drop all capabilities, add at most `NET_BIND_SERVICE` | This is one strict policy, not every valid workload. Use a profile capability allowlist and require a reason for additions. |
| `PRC-A-K8S-004` | Use `RuntimeDefault` or `Localhost` seccomp | Good Linux restricted baseline. Keep it conditional and validate `Localhost` details or external enforcement when available. |
| `PRC-A-CORE-031` | No recognized private-key armor in bounded files | Good no-go signal, but it proves only the named patterns and scanned bytes. Keep the limit visible and use the optional secret analyzer for broader evidence. |
| `PRC-A-GO-001` | Ban `net/http` client helpers | Conservative and useful, but can flag a globally configured default client. Improve with data-flow evidence; until then state the pattern exactly and provide a safe code change, not a silent waiver. |
| `PRC-A-GO-002` | Ban `net/http` server helpers | Conservative and useful because those helpers do not expose timeout fields. Add wrapper/data-flow detection if a project supplies safe bounds elsewhere. |
| `PRC-A-API-001` | OpenAPI root and version structure | Good conditional format check. Keep supported OpenAPI versions explicit and treat unsupported/newer versions as Blocked rather than invalid by guess. |
| `PRC-A-API-002` | Direct operations have valid responses | Good format-level rule. Resolve references only with strict local/network policy and bounds; do not confuse a valid contract with implemented behavior. |
| `PRC-A-API-003` | Declared operation IDs are nonempty and unique | Good because it does not require every operation to have an ID. A separate profile may require IDs for client generation. |
| `PRC-A-IAC-001` | Pinned offline Checkov reports no selected violation | Good scoped external evidence when the exact policy set, parser result, image, and database inputs are recorded. Keep “no selected finding” distinct from “infrastructure is secure.” |
| `PRC-A-SUPPLY-001` | Generate CycloneDX from repository inventory | Good repository SBOM claim. Keep its limitation visible: it is not automatically the final built artifact, container, or runtime SBOM. |
| `PRC-A-SUPPLY-002` | No known dependency vulnerability in offline repository scan | Good only when the database is fresh, identity-bound, complete enough for the package types, and parse errors fail closed. Severity and release gates must come from policy; do not make every advisory the same no-go. |

## Plan for all 10,042 controls

The companion file [`CONTROL_ACCEPTANCE_CRITERIA_REVIEW.md`](CONTROL_ACCEPTANCE_CRITERIA_REVIEW.md) contains every control. For each one it records:

- its exact source and source section;
- when it may apply;
- a proposed plain-language Pass contract;
- the exact current scanner coverage, if any;
- a proposed future evidence/check class;
- wording and genericity warnings.

It is generated by:

```sh
python3 scripts/generate_control_acceptance_review.py generate
python3 scripts/generate_control_acceptance_review.py check
```

The generated review is intentionally honest: 26 controls currently map to 43 assertions. The other controls say “not checked today.” A normal scan must never show 10,042 as the number of checks it ran.

### How a control becomes executable

1. Select one useful control for a named profile.
2. State its trigger in machine-readable form.
3. Split every `and`/`or` list into one observable assertion per promise.
4. Define the exact acceptable evidence and minimum authority.
5. Decide whether the check is repository, artifact, environment, or human evidence.
6. Define Pass, Fail, Blocked, Unknown, Stale, Conflicting, and Not Applicable behavior.
7. Add an implementation that is read-only and bounded.
8. Add normal, failing, unusual-layout, incomplete, hostile, and race fixtures.
9. Measure false passes and false failures.
10. Add it to an opt-in profile first.
11. Promote it to a default profile only after real project testing.

### Generic structure rule

Never define a broad control like this:

> Pass only when source code is in `src/`, components are in `components/`, and tests are in `tests/`.

Use this model instead:

- detect components from manifests, imports, build files, source parsers, and explicit configuration;
- detect tests from framework registration, build graph tasks, package commands, and source syntax;
- let a project declare unusual paths in a validated configuration;
- confirm declared paths exist and belong to the inventory;
- record which detector found each component and its confidence/limits;
- keep undetected or conflicting structure Blocked;
- check real outcomes such as ownership, test execution, dependency direction, and deployable boundaries rather than preferred folder names.

## Release and supply-chain plan

### Build

1. Start from an immutable scanner tag and exact commit.
2. Run the existing full Go, Python, documentation, catalog, benchmark, adapter, and release tests.
3. Build all native release binaries deterministically with the existing release builder.
4. Produce the normal release archives first; npm must reuse those exact tested binary bytes.
5. Generate the SBOM, checksums, binary attestations, and release manifest.
6. Build the six platform npm tarballs and one launcher tarball from an allowlist.
7. Run a secret scan against the exact unpacked and packed contents.
8. Use `npm pack --json` and independently inspect the file list, sizes, modes, scripts, dependencies, and hashes.
9. Install each tarball into a clean empty project with `--ignore-scripts` and no network after the tarball is present.
10. Run `prc version`, a positive fixture scan, a failing fixture scan, and a hostile fixture scan from each package.

### Publish

1. Use npm trusted publishing with OIDC instead of a long-lived npm write token.
2. Restrict the publisher to the exact public repository and exact release workflow.
3. Use a protected release environment and two-person review for the first stable releases.
4. Publish platform packages first, then the launcher with the exact same version.
5. Publish provenance automatically and retain registry signatures.
6. Never reuse a version. If any package is wrong, publish a new patch version and deprecate the bad version.
7. After publishing, install the exact registry version in clean environments, run `npm audit signatures`, compare hashes, and repeat the fixture scans.
8. Keep the direct archive checksums and attestations so npm is not the only distribution route.

Provenance is useful proof of build origin, but npm itself warns that provenance does not prove a package has no malicious code. The release process still needs source review, package allowlists, clean tests, secret scans, and incident response.

## Implementation phases to add to the project goal

### Phase A — Terminal safety before color

- Add a single terminal renderer package.
- Sanitize every untrusted string before rendering.
- Add TTY, `NO_COLOR`, `TERM=dumb`, color-mode, width, and ASCII-mode handling.
- Render every assertion state with a symbol, word, and optional color.
- Keep JSON and other machine formats byte-stable and color-free.
- Add the hostile terminal fixture suite.

**Done when:** no target-controlled byte can move the cursor, clear the screen, create a terminal link, forge a result line, or enter a machine format as ANSI output.

### Phase B — Inventory resource and race hardening

- Add per-file, total-byte, elapsed-time, and open-file budgets to inventory.
- Move all target reads to root-scoped handles or equivalent safe opens.
- Add cancellation checks during walking and hashing.
- Harden report parent creation against symlink races.
- Add huge, sparse, hard-link, unreadable, changing-file, and deep-tree fixtures.

**Done when:** hostile projects stop with a clear Blocked/budget result and never read or write outside authorized roots.

### Phase C — npm package prototype

- Create the launcher and six platform package definitions.
- Add exact optional dependency selection.
- Add version/hash verification and signal/exit-code forwarding.
- Add package-content allowlists and no-script/no-dependency checks.
- Test local tarballs on macOS, Linux, and Windows for both supported architectures where runners exist.

**Done when:** an install with `--ignore-scripts` works, an unsupported platform fails clearly, no network is used after install, a fake `PATH` binary is ignored, and direct archive and npm scans return equivalent results.

### Phase D — Attested beta publishing

- Verify the final npm scope and package names.
- Configure trusted publishing and provenance.
- Publish prerelease versions first.
- Run post-publish signature, content, identity, and end-to-end tests.
- Document exact install, one-time run, repeatable script, upgrade, uninstall, and offline paths.

**Done when:** a new user can follow the README from an empty folder to a private HTML report, verify package signatures, understand every result state, and see that scan made no fixes.

### Phase E — Fix the current assertion genericity issues

- First fix test discovery, Terraform root-module detection, fixed workflow timeout, presence-only files, dependency-update vendors, runtime declarations, and policy-specific Kubernetes rules.
- Add project configuration routes that name unusual paths and policy values without letting the target redefine scanner truth.
- Move style and hardened checks into clearly named profiles where appropriate.
- Split the broad “applicable analyses” assertion.

**Done when:** unusual but valid repository layouts pass the same objective with equivalent evidence, and strict security policies remain available as explicit profiles.

### Phase F — Grow executable coverage carefully

- Use the 10,042-control review to select high-value controls.
- Prefer checks with strong evidence and low false-pass risk.
- Publish coverage by objective, assertion, profile, evidence authority, and checker quality—not a single control percentage.
- Never make 10,042 separate scanners just to claim full coverage.

**Done when:** each new assertion has an owner, trigger, proof contract, implementation, limitations, adversarial fixtures, benchmark data, and profile decision.

## Suggested commit sequence

1. `Document npm distribution, terminal UX, and hostile-install model`
2. `Generate acceptance review for all catalog controls`
3. `Sanitize human terminal output before adding color`
4. `Render accessible per-assertion terminal results`
5. `Bound inventory bytes and harden root-scoped reads`
6. `Add script-free npm launcher and platform packages`
7. `Verify npm tarballs and cross-platform install fixtures`
8. `Publish attested npm prerelease documentation`
9. `Make core assertions layout-aware and profile-specific`

Each commit should be independently tested. Publishing must remain a separate, deliberate release step after the local package prototype passes.

## Final acceptance checklist

- [ ] The official package name and scope are owned and linked to the exact repository.
- [x] Every npm package uses the same scanner version and exact optional dependency versions.
- [x] No package has install scripts or unexpected executable JavaScript.
- [x] The launcher has zero third-party runtime dependencies.
- [x] Installation works with `--ignore-scripts`.
- [x] The launcher never downloads a binary and never searches `PATH` for one.
- [x] Native binary, platform manifest, package, and release hashes are bound together.
- [ ] npm provenance and registry signature checks pass.
- [ ] Package tarballs and release archives contain no secret or unexpected file.
- [x] A basic scan never runs project scripts, package managers, hooks, plugins, agents, or network tools.
- [x] Terminal text is safe against control-character injection.
- [x] Every state has a word/symbol; color is optional and accessible.
- [x] `NO_COLOR`, redirected output, and all machine formats contain no ANSI bytes.
- [x] The final absolute report path is always clear when a report was written.
- [x] The report remains private, standalone, escaped, and outside the scanned project by default.
- [x] Result exit codes pass unchanged through the npm launcher.
- [ ] Huge or changing repositories stop safely under documented budgets.
- [x] The README distinguishes 40 deterministic checks from 10,042 included controls.
- [x] Every one of the 10,042 controls appears once in the generated acceptance review.
- [x] Every mapped control lists its exact current assertions; every unmapped control says it is not deterministically checked today.
- [x] Layout conventions are discovery hints, not universal acceptance rules.
- [x] A normal scan still only reports; fixes remain a separate explicit workflow.

## Primary research sources

- [npm `package.json` fields, including `files`, `bin`, `os`, and `cpu`](https://docs.npmjs.com/files/package.json/)
- [npm install behavior and install-script controls](https://docs.npmjs.com/cli/install/)
- [npm one-time package execution](https://docs.npmjs.com/cli/npm-exec/)
- [npm clean, lock-based installs](https://docs.npmjs.com/cli/commands/npm-ci/)
- [npm package scopes](https://docs.npmjs.com/using-npm/scope.html/)
- [npm package name guidance](https://docs.npmjs.com/package-name-guidelines/)
- [npm trusted publishing](https://docs.npmjs.com/trusted-publishers/)
- [npm provenance and its limits](https://docs.npmjs.com/generating-provenance-statements/)
- [npm audit and signature verification](https://docs.npmjs.com/cli/audit/)
- [npm public package content and testing guidance](https://docs.npmjs.com/creating-and-publishing-scoped-public-packages/)
- [GitHub artifact attestations for release binaries](https://docs.github.com/en/actions/how-tos/secure-your-work/use-artifact-attestations/use-artifact-attestations)
- [Official Node release status and LTS guidance](https://nodejs.org/en/about/previous-releases)
- [Go terminal detection package](https://pkg.go.dev/golang.org/x/term)
- [`NO_COLOR` convention](https://no-color.org/)
