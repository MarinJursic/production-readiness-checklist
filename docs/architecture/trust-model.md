# Scanner trust and threat model

The scanner processes repositories, build artifacts, tool output, model output,
and external evidence that may be malformed, compromised, or intentionally
hostile. Safe operation therefore starts from a default-deny capability model.

## Trust zones

| Zone | Contents | Default trust |
| --- | --- | --- |
| Control plane | Catalog, profiles, policy, schemas, tool registry, verifier | Trusted and immutable for the run |
| Execution plane | Adapters, tools, builds, browsers, coding agents | Isolated and untrusted |
| Target plane | Repository, dependencies, artifacts, configuration, logs | Untrusted evidence source |
| External plane | CI, registry, cloud, cluster, telemetry, ticketing | Untrusted transport with explicitly configured authority |

The execution and target planes must not be able to mutate the control plane.
The verifier must consume immutable inputs captured before candidate remediation.
Inventory-file and benchmark-fixture reads bind an opened file back to the
root-scoped entry that was inspected. Git identity is read only through an
`os.Root` rooted at the target: metadata paths and symlinks cannot escape it,
and external `gitdir` indirection is rejected rather than followed.

## Threats in scope

- prompt injection in instructions, comments, documentation, fixtures, logs, test
  names, generated files, issue text, dependency metadata, and tool output;
- command injection through filenames, configuration, shell fragments, or adapter
  fields;
- path traversal, symlink escape, archive bombs, oversized output, and resource
  exhaustion;
- malicious build scripts, compilers, package-manager hooks, test runners, browser
  targets, and language servers;
- evidence forgery, replay, staleness, scope confusion, artifact substitution, and
  contradictory sources;
- an agent weakening tests, policies, thresholds, assertions, baselines, or
  suppressions instead of fixing a defect;
- poisoned tools, mutable tags, compromised dependencies, and unverified plugins;
- credential theft, unintended network access, and cross-run state leakage; and
- unsafe production actions caused by a repository scan or remediation request.

## Capability policy

Every executable unit receives an explicit capability manifest. Omitted
capabilities are denied.

The manifest covers at least:

- readable and writable paths;
- executable identities and arguments;
- environment variables and secret handles;
- network destinations and methods;
- CPU, memory, process, file, output, and wall-clock budgets;
- external systems and allowed operations;
- remediation class and protected paths; and
- whether a human authorization is required.

Scanner policy is snapshotted and hashed before execution. A candidate that
changes policy inputs is rejected and cannot validate itself.

## Isolation tiers

| Tier | Intended workload | Required boundary |
| --- | --- | --- |
| 0 | File parsing and metadata inspection | No target code execution |
| 1 | Curated native read-only checks | OS sandbox, read-only target, no network |
| 2 | Curated third-party tools | Rootless pinned container, read-only mounts, no socket, no network by default |
| 3 | Untrusted builds, parsers, tests, and agents | Disposable gVisor, Kata, microVM, or VM-equivalent boundary |
| 4 | Staging or production evidence | Separate read-only connector identity and explicit organization policy |

A container is packaging and an isolation layer, not a universal hostile-code
security boundary. Implementations must support a stronger Tier 3 runner rather
than treating ordinary OCI isolation as sufficient.

## Agent boundary

Agent integrations receive a scanner-created context bundle and an isolated
candidate workspace. Repository `AGENTS.md`, `CLAUDE.md`, hooks, MCP settings,
plugins, and similar files are target data, not scanner instructions.

Agents may propose changes only inside the fix contract. They cannot determine
Pass, approve exceptions, alter policy, access unrelated files, or deploy.

## Immediate stops

Execution stops safely when the scanner observes:

- a capability request outside policy;
- a sandbox or path-boundary violation;
- unexpected credential, network, or external-system access;
- ambiguous mutation of a protected file;
- evidence identity or digest mismatch;
- a tool or provider version outside its supported contract; or
- an unbounded or repeatedly non-improving remediation loop.

The stop is recorded as evidence. It is never silently converted into a pass.
