<div align="center">

<img src="docs/assets/social-card.png" alt="Vuk — Know what's left before you ship" width="100%">

# Vuk

**Know what's left before you ship.**

10,042 evidence-driven controls plus a read-only scanner for understanding what a project has proved, what failed, and what still needs review.

[![Controls](https://img.shields.io/badge/controls-10%2C042-2563eb)](docs/engineering/00-overview.md)
[![Validate](https://github.com/MarinJursic/production-readiness-checklist/actions/workflows/validate.yml/badge.svg)](https://github.com/MarinJursic/production-readiness-checklist/actions/workflows/validate.yml)
[![Docs](https://img.shields.io/badge/docs-GitHub%20Pages-0f766e)](https://marinjursic.github.io/production-readiness-checklist/)
[![GitHub stars](https://img.shields.io/github/stars/MarinJursic/production-readiness-checklist?logo=github&color=f59e0b)](https://github.com/MarinJursic/production-readiness-checklist)
[![License: MIT](https://img.shields.io/badge/license-MIT-f59e0b.svg)](LICENSE)

[Begin the complete review](docs/engineering/00-overview.md) · [Check a release quickly](docs/guides/getting-started.md) · [Use with an AI agent](docs/guides/ai-assisted-review.md) · [Contribute](CONTRIBUTING.md)

⭐ [Star this project](https://github.com/MarinJursic/production-readiness-checklist) · 🤝 [Help improve it](CONTRIBUTING.md) · [Share on LinkedIn](https://www.linkedin.com/sharing/share-offsite/?url=https%3A%2F%2Fmarinjursic.github.io%2Fproduction-readiness-checklist%2F) · [Share on X](https://twitter.com/intent/tweet?text=Vuk%3A%20know%20what%27s%20left%20before%20you%20ship.&url=https%3A%2F%2Fmarinjursic.github.io%2Fproduction-readiness-checklist%2F)

</div>

---

Production readiness is not “all tests are green.” It is an evidence-backed decision about the exact product and artifact you intend to operate.

This repository provides one start-to-finish review path:

1. **Engineering lifecycle:** 8,621 unique `USEQ-*` controls across 16 phases, from governance and requirements through implementation, delivery, operations, AI, and specialized domains.
2. **Production decision:** 1,421 `PRC-*` controls across ten release tracks, ending with evidence, sign-off, deployment, and verification.

The source material contained 25,359 checkbox lines across 215 documents. It was consolidated into 16 navigable lifecycle manuals, with exact duplicates, repeated boilerplate, mirrored production controls, and 1,969 repeated cross-volume occurrences in the final corpus removed. The [source consolidation manifest](docs/engineering/source-manifest.md) records the treatment of every source document and section.

> [!IMPORTANT]
> No checklist or AI review can prove that a nontrivial application has zero defects. A credible decision requires current evidence, explicit risk ownership, and tested recovery paths. One material failure can block approval regardless of how many unrelated controls pass.

## Start here

For a product-wide assessment, begin with the [complete engineering review](docs/engineering/00-overview.md) and follow phases 1–16 in order. Then complete the production-release tracks.

For an imminent release with an established engineering baseline:

1. Copy the [release assessment template](docs/records/release-assessment.md).
2. Run the [immediate no-go screen](docs/checklists/01-release-foundations.md#2-immediate-no-go-conditions).
3. Review affected lifecycle controls and every applicable production track.
4. Record Pass, Fail, Blocked, or Not Applicable against each stable control ID.
5. Complete the [evidence and decision track](docs/checklists/10-evidence-and-decision.md).

## Complete engineering review

| Phases | Coverage |
| --- | --- |
| [1. Governance](docs/engineering/01-governance-and-foundations.md) · [2. Product](docs/engineering/02-product-and-requirements.md) · [3. UX](docs/engineering/03-user-experience-web-and-content.md) · [4. Architecture](docs/engineering/04-architecture-and-design.md) | Scope, ownership, risk, requirements, user experience, and system design |
| [5. Code](docs/engineering/05-code-quality-and-implementation.md) · [6. Services and APIs](docs/engineering/06-application-services-and-apis.md) · [7. Data](docs/engineering/07-data-and-information-lifecycle.md) · [8. Security](docs/engineering/08-security-and-cryptography.md) · [9. Privacy](docs/engineering/09-privacy-and-data-protection.md) · [10. Testing](docs/engineering/10-verification-and-testing.md) | Product construction and verification |
| [11. Platform and delivery](docs/engineering/11-developer-experience-platform-and-delivery.md) · [12. Operations and SRE](docs/engineering/12-operations-sre-and-support.md) · [13. Documentation](docs/engineering/13-documentation-and-knowledge.md) | Build, deploy, observe, recover, support, and maintain |
| [14. Trust and safety](docs/engineering/14-trust-safety-and-ecosystems.md) · [15. AI and ML](docs/engineering/15-ai-ml-and-ai-assisted-development.md) · [16. Specialized domains](docs/engineering/16-specialized-domains-and-release-assurance.md) | Triggered and emerging product risks, connected to final release assurance |

The order is for dependable navigation, not a waterfall mandate. Iterative teams can revisit phases continuously, but every applicable control still needs a disposition before final approval.

## Final production review

The ten production tracks cover release identity and no-go conditions; product risk and architecture; source, build, and supply chain; environments, quality, and experience; application security; data, privacy, and performance; reliability and operations; maintenance, vendors, and compliance; conditional features; and the final evidence-backed decision.

[Start the production review](docs/checklists/01-release-foundations.md) or read the [readiness principle](docs/checklists/00-readiness-principle.md).

## Use it with Claude or another coding agent

The root [`CLAUDE.md`](CLAUDE.md) and copy-ready prompts teach an agent to distinguish verified evidence from assumptions, cite repository evidence, and leave production-only or organizational questions Blocked.

```text
Perform a read-only, full-lifecycle review using CLAUDE.md.
Review docs/engineering phases 1–16, then docs/checklists for the final release gate.
Cite evidence for every Pass and list every unknown separately.
Do not modify code and do not make the final release decision.
```

- [Full repository review prompt](docs/prompts/full-readiness-review.md)
- [Release-diff review prompt](docs/prompts/release-diff-review.md)
- [Evidence challenge prompt](docs/prompts/evidence-challenge.md)
- [AI-assisted review guide](docs/guides/ai-assisted-review.md)

## Scanner: quickest path

The experimental scanner turns the checklist into a repeatable repository assessment. Every scan puts all **10,042 controls** in the detailed report. The quick profile runs 18 high-signal local checks, while the default `prc/core-repository@1.0` profile runs 40 narrow deterministic checks. A narrow check can prove one exact fact, but it cannot silently mark a broader control as fully passed. Everything that still needs proof stays visible.

A normal scan is read-only. It does **not** fix files, run project code, install project dependencies, or turn missing evidence into a pass.

### Easiest setup with npm

The npm package source, launcher, six native platform packages, and release builder are implemented. The public package names were still unpublished when checked on August 25, 2026, so the commands below start working after the first release is listed in this project's release notes.

```bash
npm install -D @marinjursic/vuk
npx vuk quick
```

Use `npx vuk scan` for the larger 40-check local scan. That is the normal short path. The package has no install scripts and no third-party JavaScript dependencies. It does not download a binary after installation. npm selects one platform package; the launcher checks its release manifest and native binary hash, then starts it without a shell.

For a security-sensitive or repeatable install, pin a version from the scanner release notes and disable dependency install hooks:

```bash
npm install -D -E --ignore-scripts --no-audit --no-fund @marinjursic/vuk@X.Y.Z
npm exec --offline --no -- vuk scan
```

The install places reviewed package files in `node_modules` and updates `package.json` and `package-lock.json`; review those changes before committing them. The offline run fails instead of downloading a missing package.

For a short command that the whole project can reuse, add this to `package.json`:

```json
{
  "scripts": {
    "scan": "vuk scan"
  }
}
```

Run it with `npm run scan`. npm reserves its top-level command names, so a project cannot add a custom `npm scan` command. Use `npm run --ignore-scripts scan` when you also want npm to skip `prescan` and `postscan` hooks.

### One-time setup from source

You need Git and Go 1.27 or a compatible later supported toolchain.

```bash
git clone https://github.com/MarinJursic/production-readiness-checklist.git
cd production-readiness-checklist
go mod verify
go build -trimpath -o vuk ./cmd/prc
./vuk doctor
```

The `vuk` binary is now ready. `doctor` checks that the target and bundled catalog
can be read and explains any missing optional tool; it does not run the target.
Keep the binary in this directory, because the scanner automatically finds the compatible bundled `catalog/` beside it. Published scanner release archives already contain the binary, catalog, adapter manifests, schemas, and scanner guides together; verify a downloaded archive as described in the [release guide](docs/scanner/releases.md), extract it, and run from anywhere without installing project dependencies.

To use `vuk` without the `./` prefix, add this entire directory to `PATH`; do not move only the binary away from its compatible catalog. On macOS or Linux, for example: `export PATH="/absolute/path/to/production-readiness-checklist:$PATH"`.

### Scan a project

Choose one of three clear levels:

```bash
# Fast local screen: 18 high-signal checks, no AI.
./vuk quick /path/to/your/project

# Core local scan: 40 checks, no AI.
./vuk scan /path/to/your/project

# Full catalog AI advice, after provider login.
./vuk full codex /path/to/your/project
# or: ./vuk full claude /path/to/your/project
```

Every level still lists all 10,042 controls in the report. `quick` means fewer
local checks and less terminal noise; it does not mean the other controls
passed. `full` runs the core local scan and asks the selected AI provider for
advice on every active control. AI advice remains unverified.

From the extracted or cloned scanner directory:

```bash
./vuk scan /path/to/your/project
```

To scan the current directory, use:

```bash
./vuk scan
```

The command accepts options before or after the project path. The terminal shows every deterministic result with a word and symbol. On a real terminal, Pass is green and Fail is red; redirected and machine output has no color:

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

Vuk 0.1.0-dev

Run: 91c2…
Profile: prc/core-repository@1.0
Target: example-api (4c0e…)
Mode: scan only — no fixes and no project scripts

Checking 40 deterministic assertions...

  ✓ PASS     PRC-A-CORE-001  Observed README.md.
  ✗ FAIL     PRC-A-CORE-007  No supported lock file was found for node.
  ! BLOCKED  PRC-A-CORE-013  The optional analysis was not authorized.
  ? MANUAL   PRC-A-CORE-012  An accountable reviewer must supply evidence.

Result
Local profile result: no_go
Full catalog result: needs_review
Assessment counts: fail=1, unknown=1, manual_review=1, pass=37
Verified findings: 1
Narrow checks passed: 37
Local checks unresolved: 1
Manual decisions: 1
Complete control catalog: 10042/10042 controls included
Controls still needing review or evidence: 10020
Advisory AI reviews: 0

Scan mode: report only; no fixes were applied.
Detailed report: /Users/you/Library/Caches/vuk/reports/example-api-91c2….html
```

Open the reported HTML file in a browser. It contains all 10,042 controls, every verified finding, every narrow assertion result, exact evidence, and any advisory AI review. `needs_review` means the scanner has not proved the broad rule. `partially_verified` means linked narrow checks passed; it is still not a complete Pass. Reports are private and stored outside the scanned project by default, so creating one does not change the project being scanned.

The report separates verified problems, narrow checks that passed, unresolved
local checks, manual decisions, broad controls still needing evidence, and AI
advice. An AI citation can be marked `snapshot_location_validated` only when it
points to a real line in the exact screened snapshot. Its claim is still marked
`advisory_unverified`: a real line can be irrelevant or misunderstood.

A result-bearing exit code is not a crash: `0` means the selected profile passed, `1` means an active gate failed, and `2` means evidence remains incomplete or blocked. Use `--exit-policy never` only when a script should always receive `0` after a completed report; the report still preserves the real terminal state.

Useful report options:

```bash
# Choose a new output path; an existing file is never overwritten.
./vuk scan /path/to/project --report /safe/path/readiness.html

# Print JSON for another tool. Machine formats do not create an extra HTML file.
./vuk scan /path/to/project --format json --exit-policy never > readiness.json

# Explicitly suppress the default HTML report.
./vuk scan /path/to/project --no-report
```

The native `vuk scan` command remains available for Go, Python, Java, Rust, infrastructure, air-gapped, and mixed repositories that do not use npm.

### Optional deep review with Codex or Claude Code

The ordinary scan is local and deterministic. AI review is a separate option for broad or subjective rules such as project-appropriate folder structure. First sign in through an installed provider CLI, then scan:

```bash
vuk login codex
vuk full codex
```

For Claude Code, use:

```bash
vuk login claude
vuk full claude
```

`vuk auth` shows login status and `vuk logout codex` or `vuk logout claude` clears Vuk's saved login. These commands use each provider's official sign-in flow. Vuk stores that login in a private Vuk-only directory, separate from the provider's normal configuration, plugins, instructions, and sessions. Existing supported API-key environment variables remain an alternative.

`vuk full` uses the same guarded path as `vuk scan --ai`. Selecting the provider is also your explicit permission to send bounded, secret-screened source excerpts to that remote provider. The provider receives no target workspace path and gets no shell, file-reading, write, web, MCP, or install tools. Do not enable AI review for source that its provider is not allowed to process.

The full run is intentionally large: controls are sent in sealed batches of at most eight, and the coordinator must create one separate subagent for every control. Completed batches are saved outside the target and reused when the same scan resumes. This can take a long time and use many tokens. AI results are always labeled advisory; they cannot create a verified Pass, a final Not Applicable decision, or modify the project.

Advanced options such as reviewing one control, changing effort, or setting a Claude cost limit remain available in the [safe AI control review](docs/scanner/ai-control-review.md).

### What it checks today

The default profile checks repository governance, immutable source identity, dependency and runtime declarations, discoverable tests, GitHub Actions safety, private-key armor, OpenAPI contracts, Go HTTP timeout hazards, container definitions, Terraform locks, and Kubernetes workload policy. Focused API, Kubernetes, supply-chain, and infrastructure-as-code profiles are also available. Reviewed offline OCI adapters exist for Gitleaks, Syft, Grype, and Checkov, but external analyzers are never launched by the simple command; they require an explicit `verify-local` capability grant, an exact profile binding, and pinned local inputs. The [infrastructure policy guide](docs/scanner/infrastructure-policy.md) includes the exact Checkov command and its limits.

The scanner includes every production concern in its report, but it does not pretend every concern can be proved from source code. Unsupported runtime, organizational, environment, legal, and human evidence stays visibly blocked or in review. That is intentional: the report describes what was actually proven for one target and evidence set, not an unqualified claim that software has no defects.

Fixing is a separate workflow. `vuk scan` never calls it. The bounded `vuk fix` command works only in isolated candidate directories and supports a deliberately small set of independently verifiable changes; it never merges, deploys, or releases anything automatically.

Continue with the [complete scanner quick start](docs/scanner/getting-started.md), [safe start-to-finish walkthrough](docs/scanner/security-walkthrough.md), [CLI and exit codes](docs/scanner/cli-contract.md), [diagnostics](docs/scanner/doctor.md), [read-only agent integration](docs/scanner/mcp-agent-integration.md), [project configuration](docs/scanner/configuration.md), [state and history](docs/scanner/state-and-history.md), [supply-chain scanning](docs/scanner/supply-chain.md), and [isolated remediation](docs/scanner/remediation.md). The [research findings and improvement plan](research/PROJECT_RESEARCH_AND_IMPROVEMENT_PLAN.md) records the Reddit audit, standards review, coverage gaps, and prioritized next work. Architecture details live in the [product contract](docs/architecture/product-contract.md), [trust model](docs/architecture/trust-model.md), [adapter protocol](docs/architecture/adapters.md), [evidence model](docs/architecture/evidence-and-results.md), and [remediation contract](docs/architecture/remediation-contract.md).

## Evidence, not checkbox theater

| Field | Example |
| --- | --- |
| Control | `USEQ-1A2B3C4D` or `PRC-34-017` |
| Status | Pass / Fail / Blocked / Not Applicable |
| Owner | One accountable person |
| Evidence | Test run, configuration export, dashboard, decision record, or drill report |
| Scope | Product, commit, artifact, configuration, data, and environment |
| Reviewed | Reviewer, date, freshness, and limitations |
| Exception | Risk owner, compensating control, remediation date, and expiry |

Use the [evidence record](docs/records/evidence-record.md) and [risk exception](docs/records/risk-exception.md) templates to keep reviews auditable.

## Repository structure

```text
.
├── adapters/                 # Reviewed, immutable external-analyzer manifests
├── CLAUDE.md                  # Guardrails for AI-assisted reviews
├── catalog/                   # Versioned objectives, assertions, and profiles
├── cmd/prc/                   # Experimental scanner CLI
├── docs/
│   ├── engineering/           # 16-phase lifecycle review + source manifest
│   ├── checklists/            # 10 final production tracks
│   ├── guides/                # Adoption and AI-review guidance
│   ├── prompts/               # Copy-ready agent prompts
│   ├── records/               # Evidence, risk, release, and decision templates
│   └── scanner/               # Scanner usage and generated profile documentation
├── fixtures/                  # Positive, negative, and adversarial scanner cases
├── scanner/                   # Inventory, evaluation, adapters, and bounded remediation
├── schemas/                   # Versioned catalog and scanner JSON schemas
├── scripts/                   # Generation, validation, and integrity tooling
└── .github/                   # Issue forms, PR template, and CI workflows
```

## Standards and scope

The structure was researched against whole-lifecycle and engineering-quality sources including ISO/IEC/IEEE 12207:2026, SWEBOK Guide V4.0, ISO/IEC 25010:2023, NIST SSDF 1.1, OWASP ASVS, WCAG 2.2, NIST AI RMF, and DORA guidance. See [references and scope](docs/references.md) for primary links and limitations.

This project is engineering guidance, not legal advice, a certification, or a substitute for qualified security, privacy, accessibility, safety, or compliance review.

## Contributing

This project should improve through real-world experience. If a control is missing, duplicated, unclear, outdated, incorrectly categorized, or difficult to verify, please help fix it. Contributions can expand the checklist, improve existing wording and evidence guidance, strengthen the documentation and tooling, or advance the future AI scanner.

- [Propose a missing or improved control](https://github.com/MarinJursic/production-readiness-checklist/issues/new?template=control-proposal.yml)
- [Report a documentation or navigation problem](https://github.com/MarinJursic/production-readiness-checklist/issues/new?template=documentation.yml)
- Read the [contribution guide](CONTRIBUTING.md) and open a focused pull request

First-time contributors are welcome. Please report vulnerabilities privately as described in [SECURITY.md](SECURITY.md).

If this project helps your team, [star the repository](https://github.com/MarinJursic/production-readiness-checklist), share the [documentation site](https://marinjursic.github.io/production-readiness-checklist/), or use GitHub’s **Cite this repository** action via [CITATION.cff](CITATION.cff).

Released under the [MIT License](LICENSE).
