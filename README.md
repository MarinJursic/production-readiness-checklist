<div align="center">

<img src="docs/assets/social-card.png" alt="Production Readiness Checklist — Ship with evidence, not optimism" width="100%">

# Production Readiness Checklist

**10,042 evidence-driven controls for engineering and shipping software with confidence.**

[![Controls](https://img.shields.io/badge/controls-10%2C042-2563eb)](docs/engineering/00-overview.md)
[![Validate](https://github.com/MarinJursic/production-readiness-checklist/actions/workflows/validate.yml/badge.svg)](https://github.com/MarinJursic/production-readiness-checklist/actions/workflows/validate.yml)
[![Docs](https://img.shields.io/badge/docs-GitHub%20Pages-0f766e)](https://marinjursic.github.io/production-readiness-checklist/)
[![GitHub stars](https://img.shields.io/github/stars/MarinJursic/production-readiness-checklist?logo=github&color=f59e0b)](https://github.com/MarinJursic/production-readiness-checklist)
[![License: MIT](https://img.shields.io/badge/license-MIT-f59e0b.svg)](LICENSE)

[Begin the complete review](docs/engineering/00-overview.md) · [Check a release quickly](docs/guides/getting-started.md) · [Use with an AI agent](docs/guides/ai-assisted-review.md) · [Contribute](CONTRIBUTING.md)

⭐ [Star this project](https://github.com/MarinJursic/production-readiness-checklist) · 🤝 [Help improve it](CONTRIBUTING.md) · [Share on LinkedIn](https://www.linkedin.com/sharing/share-offsite/?url=https%3A%2F%2Fmarinjursic.github.io%2Fproduction-readiness-checklist%2F) · [Share on X](https://twitter.com/intent/tweet?text=Production%20Readiness%20Checklist%3A%2010%2C042%20evidence-driven%20software%20engineering%20and%20release%20controls.&url=https%3A%2F%2Fmarinjursic.github.io%2Fproduction-readiness-checklist%2F)

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

The experimental scanner turns the checklist into a repeatable repository assessment. A normal scan is deliberately read-only: it inventories files, evaluates the 40-assertion `prc/core-repository@1.0` profile, prints a summary, and creates a private standalone HTML report. It does **not** fix files, run project code, install project dependencies, or turn unknown evidence into a pass.

### One-time setup from source

You need Git and Go 1.27 or a compatible later supported toolchain.

```bash
git clone https://github.com/MarinJursic/production-readiness-checklist.git
cd production-readiness-checklist
go mod verify
go build -trimpath -o prc ./cmd/prc
```

The `prc` binary is now ready. Keep it in this directory, because the scanner automatically finds the compatible bundled `catalog/` beside it. Published scanner release archives already contain the binary, catalog, adapter manifests, schemas, and scanner guides together; verify a downloaded archive as described in the [release guide](docs/scanner/releases.md), extract it, and run from anywhere without installing project dependencies.

To use `prc` without the `./` prefix, add this entire directory to `PATH`; do not move only the binary away from its compatible catalog. On macOS or Linux, for example: `export PATH="/absolute/path/to/production-readiness-checklist:$PATH"`.

### Scan a project

From the extracted or cloned scanner directory:

```bash
./prc scan /path/to/your/project
```

To scan the current directory, use:

```bash
./prc scan
```

The command accepts options before or after the project path. Its terminal output is intentionally short:

```text
Run: 91c2…
Profile: prc/core-repository@1.0
Target: example-api (4c0e…)
Terminal state: no_go
Assessment counts: fail=12, unknown=2, manual_review=1, pass=15, not_applicable=10
Verified findings: 12

[FAIL] PRC-A-CORE-001 (high/required): Required repository documentation is missing.
… 4 more findings are in the detailed report.
Scan mode: report only; no fixes were applied.
Detailed report: /Users/you/Library/Caches/prc/reports/example-api-91c2….html
```

Open the reported HTML file in a browser. It contains every verified finding with severity, gate, control IDs, exact file locations, evidence IDs, remediation class, and stable fingerprint. It also contains every assertion—including pass, blocked, unknown, manual-review, and not-applicable results—with required evidence and observed evidence shown separately. Reports are created with private file permissions and are stored outside the scanned project by default, so generating a report cannot change the inventory being assessed.

A result-bearing exit code is not a crash: `0` means the selected profile passed, `1` means an active gate failed, and `2` means evidence remains incomplete or blocked. Use `--exit-policy never` only when a script should always receive `0` after a completed report; the report still preserves the real terminal state.

Useful report options:

```bash
# Choose a new output path; an existing file is never overwritten.
./prc scan /path/to/project --report /safe/path/readiness.html

# Print JSON for another tool. Machine formats do not create an extra HTML file.
./prc scan /path/to/project --format json --exit-policy never > readiness.json

# Explicitly suppress the default HTML report.
./prc scan /path/to/project --no-report
```

For a Node project, npm's standard custom-script syntax is `npm run scan`, not `npm scan`. After placing `prc` on `PATH`, this optional `package.json` entry provides that shortcut without making the scanner Node-specific:

```json
{
  "scripts": {
    "scan": "prc scan ."
  }
}
```

Then run `npm run scan`. The native `prc scan` command remains the simplest choice for Go, Python, Java, Rust, infrastructure, and mixed repositories.

### What it checks today

The default profile checks repository governance, immutable source identity, dependency and runtime declarations, discoverable tests, GitHub Actions safety, private-key armor, OpenAPI contracts, Go HTTP timeout hazards, container definitions, Terraform locks, and Kubernetes workload policy. Focused API, Kubernetes, supply-chain, and infrastructure-as-code profiles are also available. Reviewed offline OCI adapters exist for Gitleaks, Syft, Grype, and Checkov, but external analyzers are never launched by the simple command; they require an explicit `verify-local` capability grant, an exact profile binding, and pinned local inputs. The [infrastructure policy guide](docs/scanner/infrastructure-policy.md) includes the exact Checkov command and its limits.

The scanner does not yet inspect every production concern automatically. Unsupported runtime, organizational, environment, and human evidence stays visibly blocked or manual. That is intentional: the report describes what was actually proven for one target and evidence set, not an unqualified claim that software has no defects.

Fixing is a separate workflow. `prc scan` never calls it. The bounded `prc fix` command works only in isolated candidate directories and supports a deliberately small set of independently verifiable changes; it never merges, deploys, or releases anything automatically.

Continue with the [complete scanner quick start](docs/scanner/getting-started.md), [CLI and exit codes](docs/scanner/cli-contract.md), [diagnostics](docs/scanner/doctor.md), [read-only agent integration](docs/scanner/mcp-agent-integration.md), [project configuration](docs/scanner/configuration.md), [state and history](docs/scanner/state-and-history.md), [supply-chain scanning](docs/scanner/supply-chain.md), and [isolated remediation](docs/scanner/remediation.md). Architecture details live in the [product contract](docs/architecture/product-contract.md), [trust model](docs/architecture/trust-model.md), [adapter protocol](docs/architecture/adapters.md), [evidence model](docs/architecture/evidence-and-results.md), and [remediation contract](docs/architecture/remediation-contract.md).

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
