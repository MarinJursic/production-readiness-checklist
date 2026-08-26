# Production Readiness Checklist for npm

This is the small npm launcher for the native `prc` scanner. A normal scan reads
the project and writes a report; it does not fix files, install project
dependencies, or run project scripts.

The launcher has no third-party JavaScript dependencies and no install scripts.
It selects one exact platform package, checks its release-bound manifest and
binary SHA-256, and starts the native scanner without a shell or network
fallback.

For a one-time global installation and a short command in every project:

```sh
npm install -g --ignore-scripts @marinjursic/prc
prc quick
```

Use `prc scan` for the 40-check core local scan. After `prc login codex` or
`prc login claude`, use `prc full codex` or `prc full claude` for advisory AI
review of every active control. Every mode keeps all 10,042 controls visible in
the report; AI advice cannot create a verified pass.

For a pinned, security-sensitive install, use
`npm install -D -E --ignore-scripts --no-audit --no-fund @marinjursic/prc@VERSION`,
then `npm exec --offline --no -- prc scan`.

Add `"scan": "prc scan"` to the project's `scripts` object when you want to
run `npm run scan`. Use `npm run --ignore-scripts scan` to skip any local
`prescan` and `postscan` hooks. npm does not support a custom top-level
`npm scan` command. The report path is printed at the end and is clickable in
supported terminals. See the main
[project documentation](https://marinjursic.github.io/production-readiness-checklist/)
for the full result meanings, AI review opt-in, and release verification.
