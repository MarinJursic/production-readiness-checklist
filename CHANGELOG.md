# Changelog

All notable changes to this project will be documented here.

The project follows [Semantic Versioning](https://semver.org/) for published releases.

## [Unreleased]

### Added

- A public invitation to propose missing controls, corrections, documentation improvements, and tooling contributions.
- A documented long-term vision for a technology-neutral AI readiness scanner that evaluates available evidence against every applicable control and produces a complete gap report.

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
