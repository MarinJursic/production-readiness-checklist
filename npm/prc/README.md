# Production Readiness Checklist for npm

This is the small npm launcher for the native `prc` scanner. A normal scan reads
the project and writes a report; it does not fix files, install project
dependencies, or run project scripts.

The launcher has no third-party JavaScript dependencies and no install scripts.
It selects one exact platform package, checks its release-bound manifest,
native binary, and every bundled runtime file, and starts the scanner without a
shell, network fallback, post-install download, or background updater.

The package includes a plain-text `DISCLOSURE` because some opt-in scanner
features are security-related dual-use capabilities. The normal scan remains
local and read-only. Published versions are staged by an exact GitHub OIDC
identity, require human npm 2FA approval, and are verified against their npm
provenance and release-bound bytes before the matching GitHub release is made
public.

Install it once, open a project, and run it:

```sh
npm install -g @marinjursic/prc
cd /path/to/project
prc setup
prc
```

There is no `npx` prefix after this one-time install. `prc setup` is an optional
local preflight that does not run project code or AI. Run `prc version` to see
the installed version, `prc update` to check explicitly for a newer release,
`npm install -g @marinjursic/prc@latest` to update, or
`npm uninstall -g @marinjursic/prc` to remove the global command. npm installs
this one user-facing package and only the native package for the current
operating system and CPU.

The global install is outside the target project and does not add project
`node_modules` or edit its package files. Website media and contributor-only
documentation are excluded from the native npm package; all 10,042 controls,
their source text, schemas, and scanner runtime evidence remain included.

Bare `prc` runs the 40-check core scan in the current directory and writes a
private report outside the project. Use `prc /path/to/project` to scan another
folder or `prc quick` for an 18-check screen. After `prc login codex` or
`prc login claude`, use `prc full codex` or `prc full claude` for advisory AI
review of all 9,356 reviewed nondeterministic controls. Run
`prc full codex --plan` or `prc full claude --plan` first to screen the same
source and see the exact work and limits without starting a provider. Deep
review requests a primary subagent for every rule and an independent skeptical
subagent per batch; current provider output cannot prove every internal call.
The 686 reviewed deterministic controls remain blocked until their exact
programs receive authoritative evidence. Every mode keeps all 10,042 controls
visible in the report; AI advice cannot create a verified pass.

`prc verify [project]` is the short explicit local command for the core scan
plus the bundled, digest-pinned Gitleaks check. Its reviewed container image
must already exist: the scanner uses `--pull=never`, disables networking, and
mounts a sealed read-only snapshot. It does not install or download anything.

For CI or a team that deliberately wants a project lock-file entry, use
`npm install -D -E --ignore-scripts --no-audit --no-fund @marinjursic/prc@VERSION`,
then `npm exec --offline --no -- prc scan`.

With the global install, `prc ci > prc-results.sarif` is the short local-only
CI form. `prc report` reopens the newest private HTML report and
`prc cache status` shows scanner-owned disk use without deleting anything.

Add `"scan": "prc"` to the project's `scripts` object when you want to
run `npm run scan`. Use `npm run --ignore-scripts scan` to skip any local
`prescan` and `postscan` hooks. npm does not support a custom top-level
`npm scan` command. The report path is printed at the end and is clickable in
supported terminals. See the main
[project documentation](https://marinjursic.github.io/production-readiness-checklist/)
for the full result meanings, AI review opt-in, and release verification.
