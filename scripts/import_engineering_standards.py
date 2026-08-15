#!/usr/bin/env python3
"""Consolidate the supplied engineering standards archives into lifecycle manuals."""

from __future__ import annotations

import argparse
import hashlib
import re
import unicodedata
from dataclasses import dataclass, field
from pathlib import Path, PurePosixPath
from zipfile import ZipFile


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_OUTPUT = ROOT / "docs" / "engineering"
PRC_DIR = ROOT / "docs" / "checklists"

CHECKBOX = re.compile(r"^\s*- \[[ xX]\]\s+(.+)$")
EXISTING_ID = re.compile(
    r"^\*\*(?:PRC-\d{2}-\d{3}|USEQ-[A-F0-9]{8}|[A-Z]{3}-\d{5})\*\*"
    r"\s*(?:[—-]\s*)?"
)
HEADING = re.compile(r"^(#{1,6})\s+(.+?)\s*$")
LINK = re.compile(r"\[([^]]+)]\((https?://[^)]+)\)")

GAP_INCLUDED_SECTIONS = {
    "expanded gap-closure controls",
    "required evidence",
    "category no-go conditions",
}
PRC_CONTROL_COUNT = 1421

FINAL_META_SECTIONS = {
    "coverage audit, gap analysis, and consolidation map",
    "gap audit and consolidation map",
    "corpus maintenance and versioning",
    "volume completion rule",
}

# These controls are not exact text matches, but a manual review found that they
# add no obligation beyond an existing, more specific control.
SEMANTIC_EQUIVALENTS = {
    "use templates where they improve completeness without forcing irrelevant content",
    "define performance and resource budgets for critical operations before optimizing",
    "add regression tests for every material security defect",
    "prevent livelock, starvation, priority inversion, and unbounded contention",
}


@dataclass(frozen=True)
class Category:
    number: int
    slug: str
    title: str
    summary: str

    @property
    def filename(self) -> str:
        return f"{self.number:02d}-{self.slug}.md"


CATEGORIES = [
    Category(1, "governance-and-foundations", "Governance and foundations", "Applicability, evidence, quality attributes, ownership, risk, ethics, suppliers, people, and continual improvement."),
    Category(2, "product-and-requirements", "Product and requirements", "Strategy, discovery, requirements, prioritization, outcomes, experimentation, and product lifecycle."),
    Category(3, "user-experience-web-and-content", "User experience, web, and content", "Research, accessibility, localization, interaction design, frontend quality, content, SEO, and web-platform behavior."),
    Category(4, "architecture-and-design", "Architecture and design", "System understanding, modularity, state, distribution, reliability, scalability, interoperability, and sustainability."),
    Category(5, "code-quality-and-implementation", "Code quality and implementation", "Correctness, readability, contracts, errors, resources, concurrency, dependencies, review, testability, and maintainability."),
    Category(6, "application-services-and-apis", "Application services and APIs", "Backend services, APIs, identity flows, tenant isolation, content processing, jobs, queues, caching, and integrations."),
    Category(7, "data-and-information-lifecycle", "Data and information lifecycle", "Governance, contracts, semantics, modeling, quality, lineage, migrations, analytics, records, and preservation."),
    Category(8, "security-and-cryptography", "Security and cryptography", "Secure development, identity, application security, supply chain, vulnerability response, monitoring, cryptographic agility, and post-quantum planning."),
    Category(9, "privacy-and-data-protection", "Privacy and data protection", "Privacy governance, notices, consent, rights, minimization, retention, sharing, de-identification, and sensitive data."),
    Category(10, "verification-and-testing", "Verification and testing", "Test strategy, unit through acceptance testing, nonfunctional assurance, test data, static analysis, defects, and root-cause learning."),
    Category(11, "developer-experience-platform-and-delivery", "Developer experience, platform, and delivery", "Developer flow, platform engineering, engineering economics, source control, CI/CD, trusted builds, infrastructure as code, and progressive delivery."),
    Category(12, "operations-sre-and-support", "Operations, SRE, and support", "Infrastructure, SLOs, observability, alerting, on-call, incidents, continuity, capacity, FinOps, support, and retirement."),
    Category(13, "documentation-and-knowledge", "Documentation and knowledge", "Documentation governance, requirements and decisions, architecture, code, APIs, data, user help, runbooks, releases, and incidents."),
    Category(14, "trust-safety-and-ecosystems", "Trust, safety, and ecosystems", "Safety-by-design, abuse response, content integrity, authenticity, recommendations, notifications, extensions, plugins, and marketplaces."),
    Category(15, "ai-ml-and-ai-assisted-development", "AI, ML, and AI-assisted development", "AI governance, data and model supply chains, evaluation, MLOps, agents, human oversight, incidents, and AI-assisted engineering."),
    Category(16, "specialized-domains-and-release-assurance", "Specialized domains and release assurance", "Triggered modules for specialized products and the controls that connect lifecycle work to final production approval."),
]
CATEGORY_BY_NUMBER = {category.number: category for category in CATEGORIES}

QUALITY_DIRECTORY_MAP = {
    "00-core": 1,
    "01-governance": 1,
    "02-product": 2,
    "03-ux-ui-accessibility": 3,
    "04-architecture-design": 4,
    "05-code-implementation": 5,
    "06-frontend": 3,
    "07-backend-services": 6,
    "08-data": 7,
    "09-security": 8,
    "10-privacy": 9,
    "11-testing-quality-assurance": 10,
    "12-delivery-cicd": 11,
    "13-operations-sre": 12,
    "14-documentation": 13,
    "16-release-production": 16,
}

TRUST_SAFETY_CONDITIONALS = {"02", "03", "09"}
AI_CONDITIONALS = {"05"}
GAP_CATEGORY_MAP = {1: 3, 2: 5, 3: 11, 4: 14, 5: 15, 6: 7, 7: 8}

QUALITY_SKIPS = {
    "00-core/10-category-index.md": "navigation-only source index",
    "16-release-production/01-universal-production-readiness-master-checklist.md": "covered by the repository's stable PRC production-readiness controls",
}


@dataclass
class SourceDocument:
    package: str
    path: str
    title: str
    category: int | None
    text: str
    mode: str
    skip_reason: str = ""
    total_controls: int = 0
    candidate_controls: int = 0
    imported_controls: int = 0
    covered_controls: int = 0
    references: list[tuple[str, str]] = field(default_factory=list)
    selected_controls: list[CandidateControl] = field(default_factory=list)


@dataclass(frozen=True)
class CandidateControl:
    text: str
    key: str
    headings: tuple[str, ...]


def archive_markdown(path: Path, package_name: str) -> dict[str, str]:
    documents: dict[str, str] = {}
    with ZipFile(path) as archive:
        for info in archive.infolist():
            member = PurePosixPath(info.filename)
            if member.is_absolute() or ".." in member.parts:
                raise ValueError(f"Unsafe archive member: {info.filename}")
            if info.is_dir() or member.suffix.lower() != ".md":
                continue
            parts = member.parts
            if not parts or parts[0] != package_name:
                raise ValueError(f"Unexpected archive root for {info.filename}")
            relative = PurePosixPath(*parts[1:]).as_posix()
            documents[relative] = archive.read(info).decode("utf-8")
    return documents


def title_from(text: str, fallback: str) -> str:
    for line in text.splitlines():
        match = HEADING.match(line)
        if match and len(match.group(1)) == 1:
            return match.group(2).strip()
    return fallback.replace("-", " ").removesuffix(".md").title()


def normalize_control(text: str) -> str:
    value = EXISTING_ID.sub("", text.strip())
    value = unicodedata.normalize("NFKC", value)
    value = value.replace("’", "'").replace("‘", "'")
    value = value.replace("“", '"').replace("”", '"')
    value = re.sub(r"[*_`]", "", value)
    value = re.sub(r"\s+", " ", value).strip().rstrip(".")
    return value.casefold()


def universalize_control(text: str) -> str:
    """Remove programming-language wording without weakening the obligation."""
    replacements = {
        "Test both initial server responses and rendered output for JavaScript-dependent pages.":
            "Test both initial server responses and rendered output for pages that depend on client-side execution.",
        "JavaScript-disabled or failed-script behavior degrades acceptably where required.":
            "Client-side execution being unavailable or failing degrades acceptably where required.",
    }
    return replacements.get(text, text)


def control_id(key: str) -> str:
    return f"USEQ-{hashlib.sha256(key.encode('utf-8')).hexdigest()[:8].upper()}"


def count_checkboxes(text: str) -> int:
    return sum(1 for line in text.splitlines() if CHECKBOX.match(line))


def duplicate_occurrences(documents: dict[str, str]) -> int:
    keys: list[str] = []
    for path, text in documents.items():
        if not re.match(r"^\d{2}-", PurePosixPath(path).name):
            continue
        for line in text.splitlines():
            checkbox = CHECKBOX.match(line)
            if checkbox:
                control = universalize_control(
                    EXISTING_ID.sub("", checkbox.group(1).strip())
                )
                keys.append(normalize_control(control))
    return len(keys) - len(set(keys))


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as archive:
        for chunk in iter(lambda: archive.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def references_from(text: str) -> list[tuple[str, str]]:
    references: list[tuple[str, str]] = []
    seen: set[str] = set()
    in_references = False
    for line in text.splitlines():
        heading = HEADING.match(line)
        if heading and len(heading.group(1)) == 2:
            in_references = heading.group(2).strip().lower() in {
                "standards and authoritative references",
                "standards and source references",
            }
            continue
        if not in_references:
            continue
        for label, url in LINK.findall(line):
            if url not in seen:
                references.append((label.strip(), url.strip()))
                seen.add(url)
    return references


def candidate_controls(document: SourceDocument) -> list[CandidateControl]:
    candidates: list[CandidateControl] = []
    heading_stack: dict[int, str] = {}
    current_h2 = ""
    for line in document.text.splitlines():
        heading = HEADING.match(line)
        if heading:
            level = len(heading.group(1))
            text = heading.group(2).strip()
            if level == 2:
                current_h2 = text.casefold()
            if level >= 2:
                heading_stack[level] = text
                for deeper in [item for item in heading_stack if item > level]:
                    del heading_stack[deeper]
            continue

        checkbox = CHECKBOX.match(line)
        if not checkbox:
            continue
        if document.mode == "gap" and current_h2 not in GAP_INCLUDED_SECTIONS:
            continue
        text = universalize_control(EXISTING_ID.sub("", checkbox.group(1).strip()))
        key = normalize_control(text)
        headings = tuple(heading_stack[level] for level in sorted(heading_stack))
        candidates.append(CandidateControl(text=text, key=key, headings=headings))
    return candidates


def quality_category(path: str) -> int | None:
    parts = PurePosixPath(path).parts
    if len(parts) < 2:
        return None
    directory, filename = parts[0], parts[-1]
    if directory == "15-conditional-domains":
        prefix = filename.split("-", 1)[0]
        if prefix in TRUST_SAFETY_CONDITIONALS:
            return 14
        if prefix in AI_CONDITIONALS:
            return 15
        return 16
    return QUALITY_DIRECTORY_MAP.get(directory)


def final_category(volume: int, heading: str) -> int:
    """Route final-corpus sections into the existing 16-phase lifecycle."""
    normalized = heading.casefold()
    if volume == 1:
        if any(
            term in normalized
            for term in (
                "product and business",
                "product strategy",
                "product discovery",
                "requirements engineering",
                "prioritization",
                "product metrics",
                "product lifecycle",
            )
        ):
            return 2
        return 1
    if volume == 2:
        return 3
    if volume == 3:
        if any(term in normalized for term in ("frontend", "browser", "offline", "forms")):
            return 3
        if any(
            term in normalized
            for term in (
                "api",
                "backend",
                "background jobs",
                "caching",
                "integration, transition",
            )
        ):
            return 6
        if any(
            term in normalized
            for term in (
                "code quality",
                "functional correctness",
                "input validation",
                "maintainability",
                "readability",
                "abstractions",
                "error handling",
                "resource lifecycle",
                "concurrent",
                "configuration",
                "dependency",
                "refactoring",
                "code review",
                "testability",
                "reusability in code",
            )
        ):
            return 5
        return 4
    if volume == 4:
        return 9 if "privacy" in normalized or "sensitive" in normalized or "cross-border" in normalized else 7
    if volume == 5:
        if any(term in normalized for term in ("fraud", "abuse resistance", "trust, safety")):
            return 14
        return 8
    if volume == 6:
        return 10
    if volume == 7:
        return 11
    if volume == 8:
        return 12
    if volume == 9:
        return 13
    if volume == 10:
        return 15
    if volume == 11:
        if any(
            term in normalized
            for term in (
                "user-generated",
                "notifications",
                "search, ranking",
                "advertising",
                "marketplaces",
            )
        ):
            return 14
        return 16
    return 16


def split_h2_sections(text: str) -> list[tuple[str, str]]:
    sections: list[tuple[str, str]] = []
    heading = ""
    lines: list[str] = []
    for line in text.splitlines():
        match = HEADING.match(line)
        if match and len(match.group(1)) == 2:
            if heading:
                sections.append((heading, "\n".join(lines)))
            heading = match.group(2).strip()
            lines = [line]
        elif heading:
            lines.append(line)
    if heading:
        sections.append((heading, "\n".join(lines)))
    return sections


def source_documents(
    quality_docs: dict[str, str],
    gap_docs: dict[str, str],
    final_docs: dict[str, str],
) -> list[SourceDocument]:
    documents: list[SourceDocument] = []
    for path, text in sorted(quality_docs.items()):
        skip_reason = QUALITY_SKIPS.get(path, "")
        category = None if skip_reason else quality_category(path)
        if category is None and not skip_reason:
            skip_reason = "package documentation or manifest, not a product control checklist"
        documents.append(
            SourceDocument(
                package="quality standards",
                path=path,
                title=title_from(text, path),
                category=category,
                text=text,
                mode="quality",
                skip_reason=skip_reason,
                total_controls=count_checkboxes(text),
                references=references_from(text),
            )
        )

    for path, text in sorted(gap_docs.items()):
        match = re.match(r"^(0[1-7])-", PurePosixPath(path).name)
        category = GAP_CATEGORY_MAP.get(int(match.group(1))) if match else None
        skip_reason = "" if category else "package documentation, manifest, or adoption audit rather than product controls"
        documents.append(
            SourceDocument(
                package="gap supplement",
                path=path,
                title=title_from(text, path),
                category=category,
                text=text,
                mode="gap",
                skip_reason=skip_reason,
                total_controls=count_checkboxes(text),
                references=references_from(text),
            )
        )

    for path, text in sorted(final_docs.items()):
        match = re.match(r"^(\d{2})-", PurePosixPath(path).name)
        if not match:
            documents.append(
                SourceDocument(
                    package="final consolidated corpus",
                    path=path,
                    title=title_from(text, path),
                    category=None,
                    text=text,
                    mode="final",
                    skip_reason="package documentation, manifest, or audit report rather than product controls",
                    total_controls=count_checkboxes(text),
                    references=references_from(text),
                )
            )
            continue

        volume = int(match.group(1))
        for heading, section_text in split_h2_sections(text):
            heading_key = heading.casefold()
            if volume == 12:
                skip_reason = "covered by the repository's stable PRC production-readiness controls"
                category = None
            elif heading_key in FINAL_META_SECTIONS:
                skip_reason = "corpus maintenance, audit, or completion instruction rather than a product control section"
                category = None
            else:
                skip_reason = ""
                category = final_category(volume, heading)
            documents.append(
                SourceDocument(
                    package="final consolidated corpus",
                    path=f"{path}#{heading}",
                    title=heading,
                    category=category,
                    text=section_text,
                    mode="final",
                    skip_reason=skip_reason,
                    total_controls=count_checkboxes(section_text),
                    references=references_from(section_text),
                )
            )
    return documents


def existing_prc_keys() -> set[str]:
    keys: set[str] = set()
    for path in sorted(PRC_DIR.glob("*.md")):
        for line in path.read_text(encoding="utf-8").splitlines():
            match = CHECKBOX.match(line)
            if match:
                keys.add(normalize_control(match.group(1)))
    return keys


def select_controls(
    documents: list[SourceDocument], seen: set[str]
) -> dict[int, int]:
    """Select controls with earlier source packages taking precedence over addenda."""
    known_ids: dict[str, str] = {}
    category_counts = {category.number: 0 for category in CATEGORIES}
    for modes in ({"quality", "gap"}, {"final"}):
        for category in CATEGORIES:
            for document in documents:
                if (
                    document.mode not in modes
                    or document.category != category.number
                    or document.skip_reason
                ):
                    continue
                candidates = candidate_controls(document)
                document.candidate_controls = len(candidates)
                for candidate in candidates:
                    if candidate.key in seen or candidate.key in SEMANTIC_EQUIVALENTS:
                        document.covered_controls += 1
                        continue
                    identifier = control_id(candidate.key)
                    prior_key = known_ids.get(identifier)
                    if prior_key and prior_key != candidate.key:
                        raise ValueError(f"Control ID collision for {identifier}")
                    known_ids[identifier] = candidate.key
                    seen.add(candidate.key)
                    document.selected_controls.append(candidate)
                    document.imported_controls += 1
                    category_counts[category.number] += 1
    return category_counts


def render_category(
    category: Category,
    documents: list[SourceDocument],
) -> tuple[str, int, list[tuple[str, str]]]:
    body: list[str] = [
        f"# {category.title}",
        "",
        f"_Phase {category.number} of {len(CATEGORIES)} in the [complete engineering review](00-overview.md)._",
        "",
        category.summary,
        "",
        "Assess every control as **Pass**, **Fail**, **Blocked**, or **Not Applicable**. A Pass needs current evidence for the exact product, revision, artifact, configuration, environment, and data state under review.",
        "",
    ]
    imported = 0
    references: list[tuple[str, str]] = []
    reference_urls: set[str] = set()

    for document in documents:
        if document.category != category.number or document.skip_reason:
            continue
        included = document.selected_controls
        imported += len(included)
        if not included:
            continue

        body.extend(
            [
                f"## {document.title}",
                "",
                f"_Consolidated from `{document.package}/{document.path}`; {len(included):,} non-duplicative controls._",
                "",
            ]
        )
        previous_headings: tuple[str, ...] = ()
        for candidate in included:
            display_headings = candidate.headings
            if display_headings and display_headings[0] == document.title:
                display_headings = display_headings[1:]
            common = 0
            for left, right in zip(previous_headings, display_headings):
                if left != right:
                    break
                common += 1
            for index, heading in enumerate(display_headings[common:], start=common):
                if body[-1] != "":
                    body.append("")
                level = min(3 + index, 6)
                body.extend([f"{'#' * level} {heading}", ""])
            body.append(f"- [ ] **{control_id(candidate.key)}** — {candidate.text}")
            previous_headings = display_headings
        body.append("")

        for label, url in document.references:
            if url not in reference_urls:
                references.append((label, url))
                reference_urls.add(url)

    if references:
        body.extend(["## Standards and source references", ""])
        for label, url in references:
            body.append(f"- [{label}]({url})")
        body.append("")

    previous_link = "00-overview.md" if category.number == 1 else CATEGORY_BY_NUMBER[category.number - 1].filename
    if category.number == len(CATEGORIES):
        next_link = "../checklists/01-release-foundations.md"
        next_label = "Production readiness: release foundations"
    else:
        next_category = CATEGORY_BY_NUMBER[category.number + 1]
        next_link = next_category.filename
        next_label = f"Phase {next_category.number}: {next_category.title}"
    body.extend(
        [
            "---",
            "",
            f"[Previous phase]({previous_link}) · [Next: {next_label}]({next_link})",
            "",
        ]
    )
    return "\n".join(body), imported, references


def render_overview(
    category_counts: dict[int, int],
    source_documents_count: int,
    source_total: int,
    imported: int,
) -> str:
    lines = [
        "# Complete engineering review",
        "",
        "This is the start-to-finish path for evaluating a software product across its full lifecycle. Complete phases 1–16 in order, then run the existing production-readiness review before release approval.",
        "",
        f"The integrated corpus contains **{imported + PRC_CONTROL_COUNT:,} unique controls**: **{imported:,} lifecycle and quality controls** with `USEQ-` IDs plus **{PRC_CONTROL_COUNT:,} production-readiness controls** with `PRC-` IDs. The import reviewed {source_total:,} checkbox lines from {source_documents_count:,} source documents and removed repeated boilerplate, explicit consolidated copies, mirrored production controls, and exact overlaps.",
        "",
        "## How to use the sequence",
        "",
        "1. Define the product, lifecycle, organization, release, and evidence boundaries.",
        "2. Work through every phase in order; use Not Applicable only with a reviewed rationale.",
        "3. Apply every specialized module whose trigger exists.",
        "4. Preserve evidence against the stable `USEQ-` or `PRC-` identifier.",
        "5. Finish with the production no-go screen, release evidence, sign-offs, and deployment verification.",
        "",
        "The phases are ordered for navigation, not as a waterfall mandate. Iterative teams can revisit them continuously, but final approval still requires every applicable control to have a disposition.",
        "",
        "## Why this structure",
        "",
        "The sequence follows the whole-lifecycle scope of [ISO/IEC/IEEE 12207:2026](https://www.iso.org/standard/90219.html) and groups implementation concerns using the [SWEBOK Guide V4.0 knowledge areas](https://www.computer.org/education/bodies-of-knowledge/software-engineering/topics). Quality is treated as a set of explicit product attributes using [ISO/IEC 25010:2023](https://www.iso.org/standard/78176.html), while secure development remains outcome- and risk-based in line with the [NIST Secure Software Development Framework](https://csrc.nist.gov/Projects/ssdf).",
        "",
        "Specialized controls remain in the same journey rather than becoming separate products. Accessibility uses [WCAG 2.2](https://www.w3.org/TR/WCAG22/), application security maps to [OWASP ASVS](https://owasp.org/www-project-application-security-verification-standard/), and AI risk work is informed by the [NIST AI Risk Management Framework](https://www.nist.gov/itl/ai-risk-management-framework). The [references page](../references.md) explains scope and interpretation limits.",
        "",
        "## Lifecycle phases",
        "",
        "| Phase | Category | Unique controls |",
        "| ---: | --- | ---: |",
    ]
    for category in CATEGORIES:
        lines.append(
            f"| {category.number} | [{category.title}]({category.filename}) | {category_counts[category.number]:,} |"
        )
    lines.extend(
        [
            f"| Final release review | [Production readiness](../checklists/01-release-foundations.md) | {PRC_CONTROL_COUNT:,} |",
            "",
            "## Decision rule",
            "",
            "> Completion is not a score. One material failure can block approval regardless of how many unrelated controls pass.",
            "",
            "Use the [source consolidation manifest](source-manifest.md) to trace each imported source document to its destination and understand what was removed as duplicate or already covered.",
            "",
            f"[Begin phase 1: {CATEGORIES[0].title}]({CATEGORIES[0].filename})",
            "",
        ]
    )
    return "\n".join(lines)


def escape_cell(value: str) -> str:
    return value.replace("|", "\\|").replace("\n", " ")


def render_manifest(
    documents: list[SourceDocument],
    source_documents_count: int,
    source_total: int,
    imported: int,
    final_imported: int,
    archive_digests: list[tuple[str, str]],
    final_duplicate_occurrences: int,
) -> str:
    lines = [
        "# Source consolidation manifest",
        "",
        "This manifest records how the three supplied archives were incorporated. Archive content was treated as source material; operational instructions inside it were not executed.",
        "",
        "## Consolidation result",
        "",
        f"- Source documents reviewed: **{source_documents_count:,}**",
        f"- Source checkbox lines reviewed: **{source_total:,}**",
        f"- Non-duplicative lifecycle controls: **{imported:,}**",
        f"- Additional controls imported from the final corpus: **{final_imported:,}**",
        f"- Existing production-readiness controls retained: **{PRC_CONTROL_COUNT:,}**",
        f"- Final unique control set: **{imported + PRC_CONTROL_COUNT:,}**",
        "- Published lifecycle manuals: **16**, followed by the existing production-readiness review",
        "",
        "## Source archives",
        "",
    ]
    for filename, digest in archive_digests:
        lines.append(f"- `{filename}` — SHA-256 `{digest}`")
    lines.extend(
        [
            "",
            "## Consolidation rules",
            "",
            "- The quality archive's production master was not copied because this repository already maintains that review with stable `PRC-` identifiers.",
            "- The gap supplement's `Consolidated controls from the prior corpus` sections were not copied a second time. Only expanded gap-closure controls, category evidence, and category no-go controls were candidates for import.",
            f"- The final corpus was deduplicated globally across all volumes. Its {final_duplicate_occurrences:,} repeated cross-volume occurrences were not reissued as separate controls.",
            "- The final corpus's release volume was not copied because it mirrors the repository's stable `PRC-` production gate.",
            "- Corpus audits, consolidation maps, maintenance rules, and volume-completion instructions were treated as checklist-authoring material rather than product controls.",
            "- Repeated applicability boilerplate and normalized exact duplicates were retained once at their earliest lifecycle phase.",
            "- Exact matches to existing `PRC-` controls were not reissued with new identifiers.",
            "- Four manually reviewed near-matches were excluded because an existing control expresses the same or a stronger obligation.",
            "- Programming-language wording was generalized to equivalent client-side execution terminology before import.",
            "- Imported controls use deterministic `USEQ-` identifiers derived from normalized control text so evidence references remain stable across regeneration.",
            "- Similar-looking controls were retained when their wording expressed a materially different scope, condition, or evidence obligation.",
            "",
            "## Source inventory",
            "",
            "| Package | Source document | Source checks | Candidates | Imported | Covered or excluded | Destination or treatment |",
            "| --- | --- | ---: | ---: | ---: | ---: | --- |",
        ]
    )
    for document in documents:
        destination = (
            f"Phase {document.category}: {CATEGORY_BY_NUMBER[document.category].title}"
            if document.category
            else document.skip_reason
        )
        excluded = document.total_controls - document.imported_controls
        lines.append(
            "| "
            + " | ".join(
                [
                    escape_cell(document.package),
                    f"`{escape_cell(document.path)}`",
                    f"{document.total_controls:,}",
                    f"{document.candidate_controls:,}",
                    f"{document.imported_controls:,}",
                    f"{excluded:,}",
                    escape_cell(destination),
                ]
            )
            + " |"
        )
    lines.extend(
        [
            "",
            "## Interpretation boundary",
            "",
            "This project is an independent implementation-oriented synthesis. It does not reproduce or replace authoritative standards, laws, contracts, regulatory rules, or certification criteria. Applicable authoritative requirements prevail.",
            "",
        ]
    )
    return "\n".join(lines)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("--quality-archive", type=Path, required=True)
    parser.add_argument("--gap-archive", type=Path, required=True)
    parser.add_argument("--final-archive", type=Path, required=True)
    parser.add_argument("--output", type=Path, default=DEFAULT_OUTPUT)
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    quality_package = "universal-software-engineering-quality-standards"
    gap_package = "universal-software-engineering-gap-supplement"
    final_package = "universal-software-engineering-final-checklists"
    quality_docs = archive_markdown(args.quality_archive, quality_package)
    gap_docs = archive_markdown(args.gap_archive, gap_package)
    final_docs = archive_markdown(args.final_archive, final_package)
    documents = source_documents(quality_docs, gap_docs, final_docs)
    source_total = sum(document.total_controls for document in documents)
    source_documents_count = len(quality_docs) + len(gap_docs) + len(final_docs)

    args.output.mkdir(parents=True, exist_ok=True)
    seen = existing_prc_keys()
    category_counts = select_controls(documents, seen)
    total_imported = sum(category_counts.values())
    final_imported = sum(
        document.imported_controls
        for document in documents
        if document.mode == "final"
    )
    for category in CATEGORIES:
        rendered, imported, _ = render_category(category, documents)
        (args.output / category.filename).write_text(rendered, encoding="utf-8")
        if imported != category_counts[category.number]:
            raise ValueError(f"Rendered count mismatch for phase {category.number}")

    (args.output / "00-overview.md").write_text(
        render_overview(
            category_counts, source_documents_count, source_total, total_imported
        ),
        encoding="utf-8",
    )
    (args.output / "source-manifest.md").write_text(
        render_manifest(
            documents,
            source_documents_count,
            source_total,
            total_imported,
            final_imported,
            [
                (args.quality_archive.name, sha256_file(args.quality_archive)),
                (args.gap_archive.name, sha256_file(args.gap_archive)),
                (args.final_archive.name, sha256_file(args.final_archive)),
            ],
            duplicate_occurrences(final_docs),
        ),
        encoding="utf-8",
    )

    print(
        f"Imported {total_imported:,} non-duplicative controls into "
        f"{len(CATEGORIES)} lifecycle manuals from {source_total:,} source checkbox lines."
    )
    for category in CATEGORIES:
        print(f"{category.number:02d} {category.title}: {category_counts[category.number]:,}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
