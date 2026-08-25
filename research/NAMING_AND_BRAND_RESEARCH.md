# Vuk naming and brand research

Research date: 2026-08-25

## Decision

The public product name is **Vuk** (pronounced “vook,” the Croatian word for
wolf).

- CLI command: `vuk`
- npm package: `@marinjursic/vuk`
- Tagline: **Know what's left before you ship.**
- Visual idea: a compact wolf head that works as terminal-safe ASCII art
- Brand colors: midnight blue and cyan; status colors remain green, red,
  yellow, and blue

`Vuk` is three letters, one spoken beat, and easy to type. The wolf gives the
CLI an identity without describing one implementation detail or limiting the
project to a checklist. It also gives the project a distinctive connection to
its Croatian maintainer.

## What current successful names have in common

GitHub Trending was reviewed across its daily, weekly, and monthly views.
Names such as Codex, Plane, Hermes, Buzz, Needle, Pi, OpenClaw, Maka, and
OpenViking show the recurring pattern:

- one to three spoken beats;
- a short lowercase command;
- a person, creature, object, or crisp invented word instead of a sentence;
- enough visual character for an icon, mascot, or terminal mark; and
- room for the product to grow without renaming it after every feature change.

`Vuk` did not meet that bar. It read like a slogan fragment, took three
spoken beats, and had no natural character or object to draw. `Vuk` matches the
short creature-name pattern exemplified by names such as Raptor and OpenClaw.

Sources reviewed:

- [GitHub Trending today](https://github.com/trending)
- [GitHub Trending this week](https://github.com/trending?since=weekly)
- [GitHub Trending this month](https://github.com/trending?since=monthly)

## Collision screening

The search covered current GitHub repositories, developer-tool search results,
exact `NAME scan` phrases, npm package metadata, and npm executable names. A
candidate was rejected when it already named an active scanner, coding agent,
developer CLI, or closely related code-quality product.

| Name | Result |
| --- | --- |
| **Vuk** | Selected. Targeted searches found no active developer scanner, coding agent, or CLI using `vuk` or `vuk scan`. The old unscoped npm package is an unrelated Vue component library and exposes no `vuk` executable. `@marinjursic/vuk` was unclaimed on the research date. |
| Bura | Strong runner-up. Short and distinctive, but wind is harder to turn into a recognizable small terminal character than a wolf. The unscoped npm name is occupied by an unrelated React component package with no executable. |
| Zebu | Usable animal mascot and no direct scanner collision found, but the name has no meaningful connection to the maintainer or product. The unscoped npm name is occupied. |
| Norma | Clear standards meaning, but it overlaps existing developer and modeling tools and has a weaker terminal mascot. |
| Rova | Unscoped npm name was free, but it is too easy to confuse with Atlassian Rovo and the active `rove` developer CLIs. |
| Rex | Rejected. It is already used by active agent, framework, and workflow CLIs. |
| Moth | Rejected. It is already used by an active coding-agent debugging tool. |
| Tusk | Rejected. It is an active testing and AI code-review CLI. |
| Osprey | Rejected. It is an active agent framework and CLI. |
| Bubo | Rejected. It is used by active agent-orchestration and AI code-review projects. |
| Manul | Rejected. Active tools already expose the exact `manul scan` command, and an older security fuzzer uses the same name. |

This is a practical collision screen, not a legal trademark opinion. A formal
trademark clearance should precede a large commercial launch.

## Package and command decision

The package remains scoped:

```bash
npm install -D @marinjursic/vuk
npx vuk scan
```

The scope makes package ownership clear and avoids depending on an unrelated
unscoped package owner. The installed terminal command is still the short
`vuk`, so normal use does not carry the scope.

## Terminal identity

Human terminal output starts with a small static wolf. The art is embedded in
the Go binary; it does not need Figlet, a font, an image download, or another
runtime package.

```text
        /\       /\
       /  \_____/  \
      /             \
     |   o       o   |
      \      ^      /
       \   /___\   /
       `-._____.-'
           VUK
  Know what's left before you ship.
```

The banner follows these rules:

- show it only for human-facing help and interactive scan output;
- never add it to JSON, SARIF, JUnit, MCP, or redirected machine output;
- use color only on a TTY and honor `NO_COLOR` and `TERM=dumb`;
- keep the art readable without color;
- keep Pass, Fail, Blocked, and Manual as words as well as symbols; and
- sanitize project-controlled text before printing it beside trusted UI.

These rules follow the [Command Line Interface Guidelines](https://clig.dev/)
and the [`NO_COLOR` convention](https://no-color.org/). The existing terminal
layer already enforces the TTY, no-color, status-label, and untrusted-text
boundaries; the wolf replaces only the old brand banner.

## Compatibility boundary

The public name, command, npm package, cache directory, terminal output,
documentation, and new release filenames use Vuk. Stable machine contracts do
not need a breaking rename. These remain unchanged:

- repository URL and Go module path;
- internal `cmd/prc` source directory;
- versioned `prc.*` schema identifiers;
- `prc/` profile IDs;
- `PRC-*` control IDs;
- MCP tool names such as `prc_scan`; and
- signed adapter publisher fields and other digest-pinned registry identities.

Keeping these identifiers avoids breaking saved reports, profiles,
integrations, and old release manifests while giving all new user-facing paths
the Vuk name.
