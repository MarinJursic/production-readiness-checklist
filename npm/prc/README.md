# Production Readiness Scanner for npm

This is the small npm launcher for the native `prc` scanner. A normal scan reads
the project and writes a report; it does not fix files, install project
dependencies, or run project scripts.

The launcher has no third-party JavaScript dependencies and no install scripts.
It selects one exact platform package, checks its release-bound manifest and
binary SHA-256, and starts the native scanner without a shell or network
fallback.

After a published version is available, install an exact version safely:

```sh
npm install --save-dev --save-exact --ignore-scripts @marinjursic/prc@VERSION
npx prc scan .
```

Add `"scan": "prc scan ."` to the project's `scripts` object when you want to
run `npm run scan`. The report path is printed at the end. See the main
[project documentation](https://marinjursic.github.io/production-readiness-checklist/)
for the full result meanings, AI review opt-in, and release verification.
