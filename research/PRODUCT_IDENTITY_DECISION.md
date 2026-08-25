# Product identity decision

Decision date: 2026-08-25

## Decision

The project has no separate product or mascot name. Its public name is the
descriptive repository name: **Production Readiness Checklist**.

- CLI command: `prc`
- npm package: `@marinjursic/prc`
- Repository: `production-readiness-checklist`
- Terminal identity: a neutral checklist panel, not a mascot
- Status colors: green for Pass, red for Fail, yellow for Blocked, and blue
  for Manual Review

`prc` is only the short command users type. It is an abbreviation of the
project name, not a second brand.

## Why this is clearer

The repository already explains its purpose in its name. Keeping that name:

- tells a new user what the project does before they open it;
- avoids maintaining a separate brand, pronunciation, mascot, and package
  story;
- keeps the website, scanner, package, reports, and release files aligned;
- restores the original `prc` command and package contract; and
- leaves room for the scanner to grow without pretending the checklist and
  scanner are different products.

## Command and package contract

Normal use stays short:

```bash
npm install -D @marinjursic/prc
npx prc quick
npx prc scan
```

The package stays scoped so ownership is clear. The installed executable is
the short `prc` command.

## Terminal identity

Human-facing help and scan output use a compact checklist panel:

```text
  ╭────────────────────────────────────────────────────╮
  │  ✓  PRODUCTION READINESS CHECKLIST                 │
  │     Know what's ready and what still needs work.   │
  ╰────────────────────────────────────────────────────╯
```

The panel is built into the Go binary and needs no font, image, download, or
runtime package. It follows these rules:

- show it only in human-facing help and interactive scan output;
- never add it to JSON, SARIF, JUnit, MCP, or redirected machine output;
- use color only on a TTY and honor `NO_COLOR` and `TERM=dumb`;
- keep the text readable without color;
- show Pass, Fail, Blocked, and Manual Review as words as well as symbols; and
- sanitize project-controlled text before displaying it.

## Stable machine contracts

The descriptive identity also matches the stable machine-facing contract:

- repository URL and Go module path;
- internal `cmd/prc` source directory;
- versioned `prc.*` schema identifiers;
- `prc/` profile IDs;
- `PRC-*` control IDs;
- MCP tool names such as `prc_scan`; and
- the `prc` cache, release, report, and package paths.
