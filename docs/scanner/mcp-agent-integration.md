# Read-only MCP agent integration

The scanner can expose its deterministic inspect surface to Codex, Claude Code,
or another local Model Context Protocol client. The MCP process is a read-only
observer for one target selected by the operator at startup. It does not give a
model a new write path, execute project commands, launch adapters or providers,
store state, apply remediation, accept risk, or make a release decision.

This separation is intentional: the agent may edit a workspace through its own
reviewed tools, while the scanner independently rebuilds evidence and decides
whether the exact target now satisfies the selected profile.

## Start a path-locked server

Build the scanner, then give the server absolute paths for both the trusted
catalog and the repository being assessed:

```bash
go build -trimpath -o /absolute/path/prc ./cmd/prc

/absolute/path/prc mcp serve \
  --catalog-root /absolute/path/production-readiness-checklist \
  --target /absolute/path/project \
  --profile prc/core-repository
```

An optional validated configuration can be locked into the assessment:

```bash
/absolute/path/prc mcp serve \
  --catalog-root /absolute/path/production-readiness-checklist \
  --target /absolute/path/project \
  --config /absolute/path/project/production-readiness.yaml
```

The server resolves catalog, target, and configuration symlinks once. The
catalog is validated and loaded once, and the selected profile cannot change
during that process. The target and configuration contents are read again for
every plan or scan so edits made outside the server invalidate their prior
content identities.

The stdio process writes only newline-delimited JSON-RPC messages to stdout.
Diagnostics use stderr. It supports the MCP `2025-11-25` and `2025-06-18`
handshake revisions, advertises a fixed tool list, and emits both structured
content and the same serialized JSON as text for older clients. See the official
[MCP lifecycle](https://modelcontextprotocol.io/specification/2025-11-25/basic/lifecycle),
[stdio transport](https://modelcontextprotocol.io/specification/2025-11-25/basic/transports),
and [tool result](https://modelcontextprotocol.io/specification/2025-11-25/server/tools)
contracts.

## Available tools

| Tool | Input | Result | Capability boundary |
| --- | --- | --- | --- |
| `prc_plan` | Empty object | Versioned inspect plan and execution DAG | Reads inventory; runs nothing |
| `prc_scan` | Empty object | Versioned run identity, gate, bounded inventory summary, assertion results, evidence, and findings | Native inspect implementations only |
| `prc_explain` | One `assertion_id` | Exact catalog assertion and sorted linked objectives | Catalog read only |

All three definitions set MCP `readOnlyHint: true`, `destructiveHint: false`,
and `openWorldHint: false`. Those annotations help a client present the tools,
but the implementation also enforces the boundary: tool arguments cannot supply
a target, catalog, configuration, command, adapter, provider, network host,
secret, output path, or remediation mode.

`prc_scan` omits the complete file inventory from the agent-facing payload to
bound context size. It retains the inventory digest, plan digest, canonical run
ID, evidence-linked result set, and canonical findings. Scanner errors are
returned as visible MCP tool errors; they are never rewritten as passing results.

## Connect Codex

Codex local clients support command-launched stdio MCP servers. Register the
scanner with the same absolute, operator-reviewed arguments:

```bash
codex mcp add production-readiness -- \
  /absolute/path/prc mcp serve \
  --catalog-root /absolute/path/production-readiness-checklist \
  --target /absolute/path/project \
  --profile prc/core-repository
```

Codex stores server configuration in user or trusted-project configuration;
review that scope before enabling it. See the official
[Codex MCP documentation](https://developers.openai.com/codex/mcp/).

## Connect Claude Code

Use a local scope unless the exact command should deliberately be shared with a
team:

```bash
claude mcp add --transport stdio --scope local production-readiness -- \
  /absolute/path/prc mcp serve \
  --catalog-root /absolute/path/production-readiness-checklist \
  --target /absolute/path/project \
  --profile prc/core-repository
```

Confirm the connection with `claude mcp get production-readiness` or `/mcp`.
Claude Code documents local stdio commands and the security implications of MCP
configuration in its official
[MCP guide](https://code.claude.com/docs/en/mcp).

## Evidence-driven agent loop

An agent using this server should follow one narrow loop:

1. Call `prc_plan` and preserve every blocked or undetermined node.
2. Call `prc_scan` and use its `run_id`, `plan_digest`, finding IDs, and evidence
   as the baseline.
3. Use `prc_explain` before changing code for an unfamiliar assertion.
4. Propose or make one bounded change through the agent host's normal reviewed
   filesystem tools, not through this MCP server.
5. Call `prc_scan` again. A changed workspace must produce a new inventory,
   plan, and run identity.
6. Continue only while changes remain within the user's authority and the
   scanner reports independently verifiable findings. Stop on unknown, manual,
   policy, production-only, or organizational evidence that code cannot prove.

The current `prc/core-repository@0.9` profile contains 40 repository assertions.
A satisfied profile means only that this versioned profile passed for that exact
inventory and evidence set. It does not mean that all 10,042 checklist controls,
runtime behavior, production infrastructure, organizational processes, or
unknown defects were autonomously verified.

## Fail-closed limits

The server rejects duplicate JSON keys, unknown request and tool fields,
non-integer numeric request IDs, path-like tool arguments, messages larger than
1 MiB, and successful tool payloads larger than 8 MiB. Calls are accepted only
after the MCP initialization lifecycle completes. A clean stdin close ends the
process normally; startup configuration failures use CLI exit `3`, and a broken
stdio transport uses exit `4`.

There is intentionally no MCP remediation tool. Use the scanner's
[bounded isolated remediation](remediation.md) commands separately when their
explicit candidate and policy controls fit the task.
