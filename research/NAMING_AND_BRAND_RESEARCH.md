# Everylast naming and brand research

Research date: 2026-08-25

## Decision

The public product name is **Everylast**.

- CLI command: `everylast`
- npm package: `@marinjursic/everylast`
- Tagline: **Know what's left before you ship.**
- Visual idea: a row of unfinished scan points ending in a clear check mark
- Colors: midnight blue, violet, and mint

The name is short, easy to say, easy to type, and does not lock the project to
one language or to a checklist-only product. It supports the wider goal: keep
going through the evidence until every last unresolved item is visible.

## Names checked

The shortlist was checked against current web search results, developer-tool
results, GitHub results, and npm package names.

| Name | Result |
| --- | --- |
| Everylast | Selected. No exact competing developer product was found in the searches performed, and both the unscoped and `@marinjursic` npm package names returned as unclaimed on the research date. |
| Shipshape | Rejected. The name is already used by active code-quality and software-security products. |
| Plimsoll | Rejected. It is already used by an active Go load-line linter. |
| Checkride | Rejected. It is already used by a developer verification CLI. |
| Plumbline | Rejected. It is used by existing verification and agent projects. |
| Seaworthy | Rejected. It is used by an active software-security scanner. |

Examples used for collision checking include
[Shipshape](https://github.com/dmytri/shipshape),
[ShipShape Labs](https://www.shipshapelabs.com/),
[Plimsoll](https://pkg.go.dev/github.com/sourcehaven-bv/plimsoll@v0.3.0/cmd/plimsoll),
[Checkride](https://github.com/robmclarty/checkride/blob/main/README.md),
[Plumbline](https://github.com/EmilioCarrion/plumbline), and
[Seaworthy](https://seaworthycode.com/en/docs).

## Compatibility boundary

The public name, command, package, terminal output, documentation, and new
release filenames use Everylast. Existing machine contracts do not need a
breaking rename. These remain stable:

- the repository URL and Go module path;
- the internal `cmd/prc` source directory;
- versioned `prc.*` schema identifiers;
- `prc/` profile IDs;
- `PRC-*` control IDs; and
- existing MCP tool names such as `prc_scan`.
- signed adapter publisher fields and other digest-pinned registry identities.

Keeping those identifiers avoids breaking saved reports, profile references,
integrations, and old release manifests. New public examples use `everylast`.

## Important limit

This is practical product-name research, not a legal trademark clearance. A
formal clearance search and registration decision should be completed before a
large public launch. Package and domain availability can also change, so they
must be checked again immediately before reservation or publication.
