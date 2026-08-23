# Changelog

All notable changes to this project will be documented here.

The project follows [Semantic Versioning](https://semver.org/) for published releases.

## [Unreleased]

### Added

- Immutable CodeQL gates for Go, Python, and GitHub Actions, pull-request
  dependency review across all dependency scopes, and complete-history Gitleaks
  scanning with redacted output. Scanner releases repeat secret and pinned
  dependency audits before publication.
- A deterministic cross-platform scanner release builder and tag-only workflow
  with embedded build identity, bundled catalog/packs/schemas, byte-for-byte
  rebuild comparison, catalog objective sources, pack benchmark corpus,
  checksums, CycloneDX 1.6 SBOM, release compatibility manifest, and signed SLSA
  and SBOM attestations.
- A scanner release verification, supported-version, failure, compromise, and
  revocation contract that never equates a valid signature with a secure or
  production-ready artifact.
- A checksum- and provenance-bound self-scan for each scanner release, produced
  by the packaged binary and bundled catalog while preserving every manual,
  blocked, unknown, and failed result.
- A public invitation to propose missing controls, corrections, documentation improvements, and tooling contributions.
- A documented long-term vision for a technology-neutral AI readiness scanner that evaluates available evidence against every applicable control and produces a complete gap report.
- A deterministic inventory v0.2 graph with sourced, confidence-bearing package, CI, container, Terraform, Kubernetes, and symlink facts.
- A 34-assertion core repository profile with additional source-integrity,
  workflow-safety, dependency, runtime, private-key armor, container, Terraform,
  and Kubernetes checks.
- A syntax-aware, no-execution Go check for direct `net/http` package helpers
  whose request bound depends on mutable global `http.DefaultClient` state,
  including aliased and dot imports, fail-closed parser limits, and measured
  Pass, Fail, and execution-error fixtures.
- A syntax-aware Go HTTP server check that detects package-level serving helpers
  whose internally constructed `http.Server` cannot configure request timeout
  fields, with alias, dot-import, shadowing, source-location, and benchmark
  coverage; both Go timeout checks exclude `_test.go` files from production
  gates.
- Language-neutral OpenAPI discovery with provenance-bearing inventory
  components and a bounded structural check for the required 3.0, 3.1, and 3.2
  root metadata, including source locations and fail-closed future-version,
  malformed-input, and resource-limit behavior.
- Run result v0.9 source locations, allowing native analyses to bind validated
  file, line, and column coordinates directly into canonical findings and SARIF
  while retaining the frozen v0.8 schema for archived consumers.
- Profile-authorized live OCI adapter evidence, content-addressed execution records, exact inventory binding, and deterministic observation-to-assessment evaluation.
- Bounded CEL applicability evaluation over a deterministic inventory projection,
  with fail-closed resource limits and reason-bearing plan v0.2 records.
- A byte-preserving R1 remediation for broadly writable repository files, with
  exact mode postconditions and original-workspace integrity verification.
- A conservative R2 anti-gaming audit that blocks existing-test rewrites,
  suppression and skip directives, constant assertions, and empty Go tests.
- A strict, content-addressed `prc.config/v0.1` project configuration contract
  with default-deny capability semantics and adversarial parser validation.
- Configuration-bound inventory v0.3, plan v0.3, and run v0.3 identities with
  sourced declared facts, scoped CEL fields, and frozen v0.2 schemas.
- Configuration-governed remediation v0.2 contracts that bind project identity,
  protected paths, and non-expandable file, line, and attempt ceilings into
  deterministic and provider-proposal candidates.
- A bounded deterministic R1 remediation loop with isolated sequential
  candidates, cumulative budgets, fresh rescans, fail-closed acceptance, and a
  versioned machine-work-complete report that classifies every remaining result.
- A non-executing `prc doctor` command and versioned diagnostics report for
  target/catalog validity, private evidence storage, candidate filesystem
  isolation, OCI runtime identity, and optional agent-provider identity.
- A CGO-free embedded SQLite state index with strict migrations, transactional
  run/result/evidence/inventory indexing, canonical-record identity checks,
  audit events, integrity checks, and `history list`/`history show` queries.
- Diff-aware evidence invalidation with implementation-specific dependency
  rules, conservative non-repository freshness handling, and plan/run v0.4
  identities that bind the engine, profile, assertion revisions, parameters,
  and implementation definitions while preserving v0.3 schema access.
- A stable eight-class CLI failure contract that separates gate failures,
  incomplete assessments, configuration, execution, policy, internal,
  cancellation, and rejected-candidate outcomes.
- A path-locked, read-only MCP stdio server for Codex, Claude Code, and other
  compatible clients, with fixed plan, scan, and assertion-explanation tools,
  strict lifecycle and message limits, structured output schemas, and no
  scanner-owned mutation, process, network, adapter, or provider capability.

### Fixed

- Prevented hostile Git metadata and filesystem races from escaping scanner
  roots during inventory and benchmark materialization, with adversarial tests
  and pinned static-analysis gates for traversal and walk-race regressions.
- Upgraded the pinned documentation renderer to `pymdown-extensions` 11.0.1,
  closing its path-traversal file disclosure and exponential-backtracking
  denial-of-service advisories.
- Added CI vulnerability audits for pinned Python dependencies and reachable Go
  call paths, plus monthly Go module update monitoring.
- Pinned every transitive local dependency of archived scan and remediation
  schemas to a version-specific contract, preventing future mutable schema
  aliases from changing validation of historical records.
- Aligned the header star icon and replaced the oversized footer star button with a compact project-action group.
- Materialized intentionally unsafe container, Terraform, and Kubernetes
  benchmark states only in temporary copies, preventing test data from creating
  false production findings when the scanner evaluates its own repository.

## [2.1.0] - 2026-08-15

### Added

- 3,704 non-duplicative lifecycle controls from the final consolidated engineering corpus, bringing the project to 10,042 controls.
- Coverage for acquisition and exit, automated decisions and rights, operator experience, transition and interoperability, eDiscovery, physical security, assurance cases, toolchain trust, service acceptance, knowledge continuity, scientific reproducibility, and digital trust services.
- Visible GitHub star actions in the documentation header and footer.
- Validation that rejects programming-language-, framework-, or vendor-dependent control wording.

### Changed

- The Get started tab now opens the 15-minute quick start, while Engineering lifecycle continues to open the full review overview.
- Reworded client-execution controls to remain technology-neutral.
- Audited the final corpus globally, removing 1,969 repeated cross-volume occurrences, exact overlaps, four reviewed semantic equivalents, authoring instructions, and the mirrored production-release master.

## [2.0.0] - 2026-08-15

### Added

- 4,917 non-duplicative `USEQ-*` controls organized into 16 ordered engineering-lifecycle manuals.
- Complete start-to-finish review map that connects lifecycle assurance to the existing 1,421-control production gate.
- Source consolidation manifest with archive hashes and document-level import accounting for 14,132 source checkbox lines across 197 documents.
- Reproducible archive consolidation script and validation for deterministic lifecycle control IDs, counts, duplicates, pages, and links.
- Primary-source rationale based on ISO/IEC/IEEE 12207, SWEBOK Guide V4.0, ISO/IEC 25010, NIST SSDF, WCAG, OWASP ASVS, NIST AI RMF, and DORA guidance.

### Changed

- Expanded the project from a release-only checklist to a 6,338-control engineering and production-readiness review system.
- Reorganized documentation navigation, onboarding, AI prompts, citation metadata, and contributor guidance around the complete review sequence.

## [1.1.1] - 2026-08-15

### Changed

- Replaced the promotional homepage with a conventional documentation overview.
- Added a scannable start table, review process, track reference, and evidence-state reference.
- Added navigation tracking, breadcrumbs, and search suggestions.
- Reduced the custom stylesheet to restrained typography, table, focus, and color refinements.
- Simplified the repository header action to a compact GitHub icon.

### Removed

- Marketing hero, showcase cards, statistics, animated panels, and promotional page footers.
- Custom sharing JavaScript; sharing remains available through unobtrusive project links.

## [1.1.0] - 2026-08-15

### Added

- Responsive landing-page hero, review-path cards, release flow, and track explorer.
- Accessible star and share actions on the homepage and every checklist page.
- Native Web Share support with a clipboard fallback and live status messaging.
- Open Graph, X card, and `WebSite` structured metadata with a 1280×640 social image.
- GitHub citation metadata through `CITATION.cff`.

### Changed

- Theme icons now show a sun in light mode and a moon in dark mode.
- Typography, surfaces, tables, cards, focus states, and mobile spacing received a full visual refresh.
- Page-level edit/view actions and previous/next navigation make contributing and browsing easier.

## [1.0.0] - 2026-08-15

### Added

- 1,421 controls organized into ten production-readiness tracks.
- Stable `PRC-*` identifiers for evidence and issue tracking.
- Release assessment, evidence, risk exception, and go/no-go templates.
- Claude and tool-agnostic AI-assisted review guidance.
- MkDocs documentation site and GitHub Pages workflow.
- Integrity validation for control IDs, counts, sequences, and local links.
- Open-source contribution, conduct, security, and issue-reporting guidance.
