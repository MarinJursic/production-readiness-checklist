# Production Readiness Checklist: research findings and improvement plan

**Research date:** 2026-08-24

**Repository revision reviewed:** `4f985c28d0252ccf64aeae2461959072ad2ffef1` plus the uncommitted implementation described below
**Scope:** command-line usability, npm distribution, Codex and Claude Code authentication, the 10,042-control corpus, public Reddit feedback, generic coverage, evidence quality, and the path toward safe automated remediation

## Executive decision

The project does not need another large, undifferentiated batch of controls.
Its 16 lifecycle areas already cover the major technology-neutral production
readiness domains. The largest gap is implementation depth: every control is in
the report, but only 43 narrow assertions currently map to 26 control
objectives. Most controls therefore remain honestly unresolved.

The next release should focus on three outcomes:

1. make installation, login, and scanning feel like a normal CLI;
2. prove AI claims rather than merely checking that they cite a real line; and
3. grow reviewed, tested evidence producers and useful profiles in measured
   stages.

The end goal should remain an evidence engine that can propose and verify safe
candidate fixes. It must not claim that a repository alone proves a production
environment, that an AI citation proves a fact, or that any scanner guarantees
zero defects or universal production readiness.

## Genericity rules to apply before reviewing any control

These rules answer the recurring concern that a control can accidentally force
one framework, folder layout, vendor, or engineering fashion onto every
project.

1. State the required **outcome**, not one preferred implementation.
2. Put technology-specific checks in a conditional profile or adapter, never in
   the universal wording.
3. Never require universal folder names such as `src`, `components`, or
   `services`. Verify that the real layout has understandable boundaries,
   ownership, navigation, and build behavior for that project.
4. Never assume a specific language, package manager, cloud, CI system,
   database, deployment model, team size, or release method.
5. Separate a configurable policy value from the universal rule. A timeout,
   retention period, coverage target, error budget, or approval count belongs
   in a profile or project policy unless an external obligation fixes it.
6. A missing file name is not proof that an outcome is missing. A present file
   name is not proof that the outcome exists.
7. A control that requires repository, built-artifact, deployed-environment,
   external, organizational, legal, or human evidence must say which authority
   is required. Lower-authority evidence cannot be upgraded.
8. Not Applicable requires positive evidence that the trigger is absent for the
   stated scope. Missing evidence is not Not Applicable.
9. Split compound controls into atomic assertions before automating them. One
   passing sub-check must not pass the whole control.
10. Add a new control only when it describes a distinct outcome not already
    covered. Prefer mappings, examples, assertions, or profiles over duplicates.

## What is implemented now

### Scanner and reports

- All 10,042 registered controls are included in a complete report.
- `prc quick` runs an 18-assertion, language-neutral high-signal profile;
  `prc scan` runs the 40-assertion core profile; and `prc full codex|claude`
  runs the core profile plus advisory review of every active control.
- The catalog contains 43
  assertions mapped to 26 objectives across core repository, infrastructure,
  and supply-chain groups.
- Terminal output already uses a symbol and a word as well as color: green
  Pass, red Fail, yellow incomplete/manual states, and plain output when color
  is unavailable.
- A normal `prc scan` is report-only. It does not run target project code,
  install target dependencies, or call remediation.
- The default HTML report is private, stored outside the target, never silently
  overwrites an existing file, and prints its exact location.
- JSON, Markdown, HTML, SARIF, and JUnit output paths exist for automation.
- AI results are advisory and cannot create an authoritative Pass or final Not
  Applicable result.

### Short CLI path added in this implementation

```text
prc quick
prc scan
prc login codex
prc full codex
prc login claude
prc full claude
prc auth
prc logout codex
```

All three scan commands target the current directory when no path is supplied.
The existing `--ai codex` and `--ai claude` spellings remain available and,
like `prc full`, explicitly permit screened source excerpts to be processed
remotely. The longer advanced flags remain available
for one-control tests, effort, concurrency, timeout, model, resume location, and
Claude cost limits.

The login commands call the providers' official CLIs. Codex officially supports
browser-based ChatGPT sign-in and API-key sign-in. Claude Code exposes
`claude auth login`, `logout`, and `status`. The scanner stores each login under a
private scanner-only configuration directory rather than loading the provider's
normal configuration, plugins, hooks, instructions, MCP servers, or sessions.
Scans still use a new temporary home and filtered environment. Supported API-key
environment variables remain an alternative.

### npm path

The launcher, six native platform packages, release builder, integrity manifest,
and tests exist. The public package is still marked development/private and was
not published at the time of this review. Do not claim that registry installation
works until a release note confirms publication.

After publication, the normal path is:

```bash
npm install -D @marinjursic/prc
npx prc quick
```

Use `npx prc scan` for the core local profile.

For a project shortcut:

```json
{
  "scripts": {
    "scan": "prc scan"
  }
}
```

Then run `npm run scan`. npm does not allow a project to define an arbitrary
top-level `npm scan` command, so `npm run scan` is the shortest honest npm
script form. `npx prc scan` resolves the locally installed binary after the
dependency is installed.

For a pinned, security-sensitive installation:

```bash
npm install -D -E --ignore-scripts --no-audit --no-fund @marinjursic/prc@X.Y.Z
npm exec --offline --no -- prc scan
```

The normal command is intentionally short. The hardened command remains
documented for users who need exact versions, disabled dependency hooks, and a
run that fails instead of fetching a missing package.

### 10,042-control acceptance guide

The former 27,237,260-byte Markdown file was too large for common editors and
repository viewers. It is now generated as:

- `research/control-acceptance-criteria/README.md`: index and methodology;
- ten `part-NNN.md` files, split only between complete controls; and
- a small compatibility pointer at
  `research/CONTROL_ACCEPTANCE_CRITERIA_REVIEW.md`.

Every part is below the strict 3,000,000-byte limit. Generation checks the
limit, exact part set, ordering, and that all 10,042 control markers occur once.
No control contract was shortened or omitted.

## Expected user journey

### Short path

1. The user installs the published package with
   `npm install -D @marinjursic/prc`.
2. The user runs `npx prc quick` for the small screen or `npx prc scan` for the
   core local scan.
3. The terminal names the scanner, target, profile, scan-only mode, and number
   of deterministic checks.
4. Each narrow result shows a symbol, Pass/Fail/Blocked/Manual word, stable ID,
   and short reason. Color is an extra cue, not the only cue.
5. The summary separates the local profile result from the unresolved full
   catalog. It states that all 10,042 controls were included.
6. The final line prints the exact private HTML report path.
7. The user opens that report to see findings, observed evidence, required
   evidence, limits, exclusions, and every unresolved control.
8. Nothing in the project was fixed or executed.

### AI-assisted path

1. The user installs the official Codex or Claude Code CLI.
2. The user runs `prc login codex` or `prc login claude` and completes the
   provider's official sign-in flow.
3. `prc auth` shows whether the scanner can use the private login.
4. The user runs `prc full codex` or `prc full claude`.
5. Production Readiness Checklist inventories and hashes the target locally, builds a bounded snapshot,
   excludes sensitive names, screens token/key shapes, and marks repository
   text as untrusted.
6. The provider receives sealed batches and no target path, general shell,
   source-reading, write, web, browser, install, or ambient MCP authority.
7. Completed batches resume from private state. The run may take a long time
   and use many tokens.
8. AI advice appears separately in the final report and cannot turn itself into
   verified evidence or modify files.

## Threat review of each user-facing step

| Step | Plausible attack or failure | Current protection | Remaining work |
| --- | --- | --- | --- |
| Find package | Typosquatted or similarly named package | Scoped canonical name in docs | Publish verified links from release notes only |
| Install | Lifecycle script runs with user authority | Production Readiness Checklist packages have no install scripts or third-party JS runtime dependencies; hardened command uses `--ignore-scripts` | Publish package provenance, signatures/checksums, and reproducible package comparison |
| Resolve native binary | Launcher downloads or executes a PATH replacement | Exact platform package, release manifest, SHA-256 check, no shell, no download fallback | Release and revocation tests on every platform |
| `npx` execution | Missing local package causes a registry fetch or prompt | Normal flow installs locally first | Document offline `npm exec --offline --no` for sensitive use and test missing-local failure |
| Inventory hostile repository | Symlink escape, huge files, changing files, strange names, repository code execution | No target scripts, no symlink following, budgets, identity/hashes, rechecks, recorded exclusions | Expand archive/parser/Unicode/path and resource-exhaustion corpus |
| Deterministic check | Parser bug or framework-specific assumption | Narrow assertion wording, bounds, explicit profile, incomplete state | Add benchmark repositories, precision/recall budgets, and more unusual-layout fixtures |
| Provider login | Normal user hooks/plugins/settings alter the scan | Dedicated private scanner credential root and official provider auth command | Test OS keyring/file variations and enterprise-managed-policy failure behavior |
| Provider launch | Fake or compromised `codex`/`claude` binary ignores flags | Resolved executable name/digest, strict launch plan, filtered environment, temporary home | Optional external OS sandbox; signed provider verification where available |
| Remote review | Secret leaves the host | Sensitive-name exclusion, token/key screening, byte/file/output limits, explicit `--ai` consent | Make clear the screen is not a complete secret scanner; add pluggable preflight and organization policy |
| Prompt handling | Repository text instructs the model to escape rules | Sealed JSON, untrusted-data label, no general tools, strict result schema | Larger adversarial corpus and provider-version regression suite |
| Citation | AI cites a real line that does not prove its sentence | Exact snapshot/path/line/range/digest validation | Add semantic claim-to-observation verification described below |
| Result | Users read advisory text as proof | Advisory-only state and separate authoritative result | Stronger report labels and “why this is not verified” explanations |
| Remediation | Scan silently changes code or weakens gates | Scan and fix are separate; fixes use isolated candidates and verification | Expand only with protected paths, anti-gaming, budgets, human gates, and no automatic deployment |

## Public Reddit feedback audit

This audit covers the project-related submissions and public comments visible
from the maintainer's submitted-post history on 2026-08-24. Removed posts and
comments that Reddit did not expose cannot be reconstructed reliably. Praise,
generic promotion advice, and author replies are recorded as context but do not
become product requirements.

| Feedback | Does it make sense? | Project status | Decision |
| --- | --- | --- | --- |
| Reviewing 10,000 controls without context is impractical; the validator looked nested and string-heavy; regex case handling could be brittle. | Yes on review scale and maintainability. The regex concern is a check request, not proof of a bug. | `scripts/validate.py` is now split into focused functions and uses case-insensitive matching/case folding. The 10K review is split into bounded files. Semantic duplicates remain difficult. | Keep validator refactoring and add an offline semantic-overlap candidate report with human adjudication; never auto-delete controls. |
| Stable IDs are useful, but 10K items are unwieldy and invite scope creep. | Yes. A complete catalog should not force a beginner to process a 10K list as one user task. | Quick, Core, and Full commands now disclose 18 local checks, 40 local checks, or core plus all-control AI advice. Every report still preserves the full catalog. | Validate the commands with first-time users and tune the quick assertion set only through measured benchmark results. |
| Ask another model to challenge a claim and produce a counterexample; use machine/kernel truth for actions that actually happened. | Partly. A challenger can find weak reasoning, but a second model is not independent proof. Host/runtime observations can be higher-authority evidence when collected safely. | AI cannot produce verified Pass; evidence authorities are separated. | Add a deterministic citation/claim verifier first. Offer a second-model challenge only as advisory. Add signed runtime evidence through explicit connectors later. |
| Separate evidence found from evidence required and record time. | Yes. | `AssertionResult` separates required and observed evidence. HTML now shows authority, stable observation ID, observation time, and a warning that collection time is not proof of freshness. | Add project-specific freshness policies and stale-state tests where authoritative evidence producers are introduced. |
| Do not trust a persuasive filename or README claim; parse content and hash what was inspected. | Yes. | Inventory hashes source; checks are narrow; AI input is untrusted; README existence proves only README existence. | Add adversarial “good-sounding filename/stub/README claim” fixtures to every relevant assertion family. |
| Deterministic extraction should establish facts; the LLM should explain rather than invent them. | Yes, with one nuance: AI can propose a candidate concern that still needs proof. | Deterministic and advisory states are separate; AI cannot create verified Pass/NA. | Keep AI candidates, but label whether each sentence is copied observation, deterministic inference, or unverified advice. |
| The work needs stages, prompt-injection quarantine, world-truth checks, and human gates. | Yes. | Sealed prompts, excluded instruction files, capability denial, resume batches, evidence authority, and human/manual states exist. | Follow phased releases and publish adversarial evaluation results before widening authority. |
| A valid path and line citation is not verification of the claim at that line. | Yes; this is the most important unresolved feedback. | Paths, line ranges, snapshot digests, task IDs, ordering, and schemas are validated. Reports now record `snapshot_location_validated` separately from `advisory_unverified`, including an adversarial irrelevant-citation test. Semantic support is not yet independently re-parsed. | Continue with typed claim-to-observation verifiers; never label a location-valid citation as semantic proof. |
| One FOSS response said the project looked worth reviewing. | Positive but not actionable. | No change needed. | No product requirement. |
| A beta-testing bot asked posts to explain the product, link, and requested feedback clearly. | Useful for future outreach, not scanner design. | Project posts already describe scope and requested feedback. | Keep as communication guidance only. |

Sources: [r/codereview discussion](https://www.reddit.com/r/codereview/comments/1vpq6h3/code_review_wanted_validating_10000_checklist/),
[r/sideprojects post](https://www.reddit.com/r/sideprojects/comments/1vpausn/im_building_an_opensource_production_readiness/),
[r/ClaudeCoding discussion](https://www.reddit.com/r/ClaudeCoding/comments/1vpq908/how_do_you_stop_claude_from_turning_repository/),
and [r/LLMDevs discussion](https://www.reddit.com/r/LLMDevs/comments/1vperza/designing_an_opensource_llm_scanner_that_must/).

## Generic coverage research

The corpus was compared with current primary sources including NIST SSDF 1.1,
NIST CSF 2.0, OWASP ASVS 5.0.0, SLSA 1.2, WCAG 2.2, NIST AI RMF, the FinOps
Framework, and the Software Carbon Intensity specification. These sources
support outcome-oriented, versioned mappings. They do not justify copying every
framework item into a universal checklist.

A lexical sample across `docs/checklists` and `docs/engineering` found broad
existing coverage: post-quantum/crypto agility 22 matches; sustainability,
carbon, energy, or environmental impact 110; retirement/decommission/end of
life/sunset 132; non-human/workload identities or service accounts 12;
residency/sovereignty/localization 47; SBOM 28; VEX 5; prompt injection 7;
accessibility/WCAG 498; FinOps/unit economics/allocation/anomaly 7; safety
case/hazard/functional safety 10; memory safety 4; model/system cards 3. These
counts only locate review candidates; they do not prove quality or completeness.

### No missing top-level area was established

The corpus already includes governance, product, UX and accessibility,
architecture, code quality, services/APIs, data, security and cryptography,
privacy, testing, platform and delivery, operations/SRE/support, documentation,
trust and safety, AI/ML, and specialized domains. It also contains substantial
coverage of supply chain, cost, sustainability, post-quantum migration,
decommissioning, data residency, mobile, desktop, IoT, cloud native,
blockchain, safety-critical systems, and public/open-source packages.

Absence of the term `OSPO` is not a universal gap. Requiring an Open Source
Program Office would be inappropriate for many projects; the underlying open
source governance, licensing, contribution, and maintenance outcomes already
belong in conditional organization profiles.

### Depth candidates worth expert review

These are candidates for crosswalk and wording review, not instructions to add
bulk controls immediately:

1. **Non-human identity lifecycle:** issuance, owner, least privilege,
   rotation, workload binding, monitoring, and retirement are present but less
   explicit than human identity coverage.
2. **Safety-case and hazard evidence:** strengthen the link from hazard,
   mitigation, residual risk, verification evidence, and accountable approval
   for projects whose failure can harm people or physical systems.
3. **AI model/system transparency and evaluation traceability:** clarify
   versioned model/system cards, evaluation-set provenance, intended-use limits,
   known failure modes, and change-triggered re-evaluation without making these
   universal for non-AI systems.
4. **Current standards crosswalks:** map, do not duplicate, NIST SSDF 1.1, NIST
   CSF 2.0, OWASP ASVS 5.0.0, SLSA 1.2 including its Source Track, WCAG 2.2,
   NIST AI RMF/GenAI profile, FinOps, and SCI. Store source version and review
   date with every mapping.
5. **VEX decision evidence:** the corpus mentions VEX, but future artifact
   workflows should preserve why a vulnerability is affected/not affected,
   the exact product/version, evidence, author, time, and expiry.

## Prioritized implementation plan

### P0 — finish and release the easy, trustworthy CLI

- Complete docs and tests for `prc login`, `auth`, `logout`, and `--ai`.
- Test browser login, status, logout, API-key fallback, missing provider,
  oversized output, provider failure, unsafe directory permissions, symlinks,
  and enterprise-managed settings on supported operating systems.
- Publish the seven npm packages only through the existing reviewed release
  workflow; remove development/private markers only in release artifacts.
- Add npm provenance, release checksums/signatures, SBOM, build provenance, and
  a documented package/release revocation process.
- Keep `npm install -D @marinjursic/prc`, `npx prc quick`, and
  `npx prc scan` as the normal docs;
  keep the exact offline form next to the security explanation.

**Done when:** a new user can install a published version, run one local scan,
find the report, login to either supported provider, run an AI scan, inspect
status, and logout without reading an advanced flag page; all paths have
cross-platform integration tests and no hidden download or install script.

### P1 — make 10,042 controls usable without hiding scope

**Implemented foundation:** the CLI now exposes `quick`, core `scan`, and
provider-named `full` commands; all retain the complete catalog. The report's
first page now separates verified problems, narrow passes, unresolved local
checks, manual decisions, controls needing evidence, and AI advice. The
remaining work below is progress, prioritization, cancellation, and usability
measurement.

- Add named intents such as `quick`, `core`, and `full`, or an interactive
  first-run selector. The command must show selected controls/assertions before
  spending tokens.
- Keep all 10,042 controls in a complete report, but prioritize applicable,
  high-impact, changed, and currently provable items in the terminal and report.
- Show elapsed batches, completed/remaining control counts, resume location,
  provider/model, concurrency, timeouts, and any enforceable cost boundary.
- Add clean cancellation and resume summaries. Never promise a precise Codex
  cost limit when its CLI cannot enforce one.
- Improve the HTML report's first page: “proven problems,” “proven passes,”
  “needs evidence,” “manual decisions,” and “AI advice” must be visibly
  separate.

**Done when:** usability tests show a first-time user can explain what was
checked, what remains unknown, where the report is, whether remote processing
occurred, and why the result is not a universal readiness certificate.

### P1 — verify AI citations and claims

**Implemented foundation:** every new AI result records citation-location state
separately from claim state. Real snapshot locations are
`snapshot_location_validated`; no citation is `not_cited`; all current AI claim
text is `advisory_unverified`. This closes the reporting ambiguity but does not
perform semantic proof. The typed verifier work below remains required.

- Extend AI output so each claim has a stable claim ID, exact quoted or parsed
  observation, expected evidence authority, source digest, range, and
  verification state.
- Re-read the exact snapshot bytes after response validation.
- For simple claims, use typed deterministic verifiers: exact text/value,
  parsed manifest field, workflow permission, route/middleware relationship,
  configuration value, or dependency edge.
- Record `citation_valid`, `claim_matches_observation`,
  `authority_sufficient`, and `human_review_required` separately.
- If no verifier exists, retain `advisory_unverified`; never turn it into Pass.
- Optionally run a second model as a challenger that produces counterexamples,
  but keep that result advisory.

**Done when:** fixtures with real-but-irrelevant citations, stale lines, swapped
files, correct numbers with wrong meaning, README promises, and fabricated
runtime outcomes are rejected or remain explicitly unverified.

### P1 — build a public adversarial and quality benchmark

- Add small language-diverse and unusual-layout repositories with known truth.
- Cover misleading filenames, empty stubs, copied policy text, fake badges,
  prompt injection, secret bait, Unicode/path tricks, symlinks, large files,
  archives, changing files, ambiguous monorepos, vendored/generated code, and
  contradictory evidence.
- Measure precision, recall, unsupported rate, blocked rate, flakiness,
  reproducibility, scan time, memory, and provider-version regressions per
  assertion/profile.
- Require published quality budgets before an assertion can become a default
  gate.

**Done when:** every default-gating assertion has positive, negative,
Not Applicable, unusual-layout, limit, and adversarial cases plus a published
quality result.

### P2 — turn generated contracts into reviewed evidence producers

- Triage by impact and applicability, not control ID order.
- Review each generated contract with a domain owner before calling it trusted.
- Split compound objectives into atomic assertions and map equivalent controls
  instead of writing duplicate checks.
- Add native parsers first, then immutable offline adapters, then explicit
  environment/external connectors. Preserve evidence authority and freshness.
- Add versioned crosswalks to current primary standards. A mapping update must
  not silently change old run conclusions.
- Improve semantic duplicate detection as an offline review report. It should
  show candidate pairs, similarity reasons, source sections, and possible
  conflict; a person decides whether to merge or retain them.

**Done when:** each newly automated assertion has an approved contract,
applicability rule, implementation version, evidence authority, limits,
fixtures, benchmark result, and no broader-control overclaim.

### P2 — controlled runtime and organizational evidence

- Define read-only connector contracts for deployment state, monitoring,
  incident systems, artifact registries, identity systems, and policy records.
- Require explicit capability, subject/environment identity, least-privilege
  credentials, bounded collection, freshness, provenance, redaction, and
  revocation.
- Keep legal conclusions, risk acceptance, final Not Applicable decisions, and
  release approval human-owned.

**Done when:** repository configuration can no longer be confused with deployed
state, and every external observation states subject, environment, collector,
time, freshness, scope, and limitations.

### P3 — expand remediation only after verification is strong

- Keep `scan` and `fix` separate forever.
- Expand fixes from deterministic R1 changes to bounded R2 agent proposals only
  when the related assertion has independent verification.
- Use isolated candidates, allowlisted paths, protected policy/test/evidence
  files, no network/secrets by default, token/time/change budgets, baseline and
  regression tests, anti-gaming checks, and predictable termination.
- Never automatically merge, deploy, release, accept risk, weaken gates, delete
  tests, or claim that a source fix proves production behavior.

**Done when:** a proposed fix can be reproduced and independently closes its
atomic assertion without weakening unrelated safeguards; all remaining unknowns
stay visible and deployment/release approval remains outside the loop.

## Reconciliation with the earlier scanner roadmap

The supplied `PRODUCTION_READINESS_SCANNER_ROADMAP.md` was treated as research
input, not as executable instructions. Much of its foundation is now present:
machine-readable objectives/assertions/profiles, immutable inventory and plans,
evidence/result schemas, fail-closed capability handling, OCI adapter contracts,
HTML/JSON/Markdown/SARIF/JUnit outputs, state/history, diff support, MCP read-only
inspection, provider adapters, isolated candidate remediation, and explicit
unknown/manual states.

Its most important remaining recommendations match this review: approved
control-contract migration, much broader assertion coverage, public benchmarks,
independent claim verification, external evidence connectors, release supply
chain hardening, and cautious staged remediation. The roadmap's core correction
still holds: agents should generate bounded candidates; the scanner engine owns
truth; people own risk and release decisions.

## Primary research sources

- [OpenAI Codex authentication](https://learn.chatgpt.com/docs/auth)
- [Claude Code CLI reference](https://code.claude.com/docs/en/cli-reference)
- [Claude Code getting started](https://code.claude.com/docs/en/getting-started)
- [npm install](https://docs.npmjs.com/cli/v11/commands/npm-install/)
- [npm exec and npx](https://docs.npmjs.com/cli/v11/commands/npm-exec/)
- [NIST Secure Software Development Framework 1.1](https://csrc.nist.gov/pubs/sp/800/218/final)
- [NIST Cybersecurity Framework 2.0](https://csrc.nist.gov/pubs/cswp/29/the-nist-cybersecurity-framework-csf-20/final)
- [OWASP ASVS 5.0.0](https://owasp.org/www-project-application-security-verification-standard/)
- [SLSA 1.2](https://slsa.dev/spec/v1.2/)
- [WCAG 2.2](https://www.w3.org/TR/WCAG22/)
- [NIST AI Risk Management Framework](https://www.nist.gov/itl/ai-risk-management-framework)
- [FinOps Framework](https://www.finops.org/framework/)
- [Software Carbon Intensity specification](https://sci.greensoftware.foundation/)
