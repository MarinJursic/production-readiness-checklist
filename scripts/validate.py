#!/usr/bin/env python3
"""Validate checklist integrity and local Markdown links without dependencies."""

from __future__ import annotations

import re
import sys
import unicodedata
from collections import Counter, defaultdict
from pathlib import Path
from urllib.parse import unquote


ROOT = Path(__file__).resolve().parents[1]
CHECKLIST_DIR = ROOT / "docs" / "checklists"
ENGINEERING_DIR = ROOT / "docs" / "engineering"
EXPECTED_CONTROLS = 1421
EXPECTED_ENGINEERING_CONTROLS = 8621
EXPECTED_PAGES = {
    "00-readiness-principle.md",
    "01-release-foundations.md",
    "02-product-risk-architecture.md",
    "03-source-build-supply-chain.md",
    "04-environments-quality-experience.md",
    "05-application-security.md",
    "06-data-privacy-performance.md",
    "07-reliability-operations.md",
    "08-maintenance-vendors-compliance.md",
    "09-conditional-modules.md",
    "10-evidence-and-decision.md",
}
EXPECTED_ENGINEERING_PAGES = {
    "00-overview.md",
    "01-governance-and-foundations.md",
    "02-product-and-requirements.md",
    "03-user-experience-web-and-content.md",
    "04-architecture-and-design.md",
    "05-code-quality-and-implementation.md",
    "06-application-services-and-apis.md",
    "07-data-and-information-lifecycle.md",
    "08-security-and-cryptography.md",
    "09-privacy-and-data-protection.md",
    "10-verification-and-testing.md",
    "11-developer-experience-platform-and-delivery.md",
    "12-operations-sre-and-support.md",
    "13-documentation-and-knowledge.md",
    "14-trust-safety-and-ecosystems.md",
    "15-ai-ml-and-ai-assisted-development.md",
    "16-specialized-domains-and-release-assurance.md",
    "source-manifest.md",
}
CONTROL = re.compile(
    r"^- \[ \] \*\*(PRC-(\d{2})-(\d{3}))\*\* — (.+)$", re.MULTILINE
)
PLAIN_CHECKBOX = re.compile(r"^- \[ \] (?!\*\*PRC-)", re.MULTILINE)
ENGINEERING_CONTROL = re.compile(
    r"^- \[ \] \*\*(USEQ-[A-F0-9]{8})\*\* — (.+)$", re.MULTILINE
)
ENGINEERING_PLAIN_CHECKBOX = re.compile(r"^- \[ \] (?!\*\*USEQ-)", re.MULTILINE)
MARKDOWN_LINK = re.compile(r"!?\[[^\]]*\]\(([^)]+)\)")
MARKDOWN_HEADING = re.compile(r"^#{1,6}[ \t]+(.+?)[ \t]*#*[ \t]*$", re.MULTILINE)
IMPLEMENTATION_SPECIFIC = re.compile(
    r"\b(?:JavaScript|TypeScript|Node\.js|Vue\.js|Svelte|Next\.js|Nuxt|Django|"
    r"Flask|FastAPI|NestJS|Python|PHP|C\+\+|C#|Kotlin|Scala|Elixir|Haskell|"
    r"AWS|Azure|Google Cloud|GCP|GitHub|GitLab|Bitbucket|Docker|Kubernetes|"
    r"Terraform|Jenkins|CircleCI|Datadog|New Relic|Sentry|Stripe|Cloudflare|"
    r"Vercel|Netlify|PostgreSQL|MySQL|MongoDB|Redis|Kafka|RabbitMQ|OpenAI|"
    r"Anthropic|Claude)\b",
    re.IGNORECASE,
)


def markdown_files() -> list[Path]:
    files = list(ROOT.glob("*.md"))
    files.extend((ROOT / "docs").rglob("*.md"))
    files.extend((ROOT / ".github").rglob("*.md"))
    return sorted(set(files))


def duplicate_values(values: list[str]) -> list[str]:
    """Return repeated values once, in deterministic order."""

    return sorted(value for value, count in Counter(values).items() if count > 1)


def normalize_control(text: str) -> str:
    value = unicodedata.normalize("NFKC", text)
    value = value.replace("’", "'").replace("‘", "'")
    value = value.replace("“", '"').replace("”", '"')
    value = re.sub(r"[*_`]", "", value)
    value = re.sub(r"\s+", " ", value).strip().rstrip(".")
    return value.casefold()


def markdown_anchor(heading: str) -> str:
    """Approximate Python-Markdown's default heading anchor generation."""

    value = re.sub(r"<[^>]+>", "", heading)
    value = re.sub(r"[*_`~]", "", value)
    value = unicodedata.normalize("NFKD", value).casefold()
    value = re.sub(r"[^\w\- ]", "", value)
    return re.sub(r"[-\s]+", "-", value).strip("-")


def markdown_anchors(content: str) -> set[str]:
    anchors: set[str] = set()
    counts: Counter[str] = Counter()
    for match in MARKDOWN_HEADING.finditer(content):
        base = markdown_anchor(match.group(1))
        if not base:
            continue
        index = counts[base]
        anchors.add(base if index == 0 else f"{base}_{index}")
        counts[base] += 1
    return anchors


def validate_pages(
    errors: list[str],
) -> tuple[list[str], dict[int, list[int]], list[str]]:
    actual_pages = {path.name for path in CHECKLIST_DIR.glob("*.md")}
    if actual_pages != EXPECTED_PAGES:
        errors.append(
            "Checklist page set differs: "
            f"missing={sorted(EXPECTED_PAGES - actual_pages)}, "
            f"unexpected={sorted(actual_pages - EXPECTED_PAGES)}"
        )

    ids: list[str] = []
    normalized_controls: list[str] = []
    numbers: dict[int, list[int]] = defaultdict(list)
    for path in sorted(CHECKLIST_DIR.glob("*.md")):
        content = path.read_text(encoding="utf-8")
        for match in PLAIN_CHECKBOX.finditer(content):
            line = content.count("\n", 0, match.start()) + 1
            errors.append(f"{path.relative_to(ROOT)}:{line}: checklist item has no PRC ID")
        for match in CONTROL.finditer(content):
            control_id, section, item, text = match.groups()
            technology = IMPLEMENTATION_SPECIFIC.search(text)
            if technology:
                line = content.count("\n", 0, match.start()) + 1
                errors.append(
                    f"{path.relative_to(ROOT)}:{line}: {control_id} depends on "
                    f"a language, framework, or vendor ({technology.group(0)})"
                )
            ids.append(control_id)
            numbers[int(section)].append(int(item))
            normalized_controls.append(normalize_control(text))

    if len(ids) != EXPECTED_CONTROLS:
        errors.append(f"Expected {EXPECTED_CONTROLS} controls, found {len(ids)}")
    if len(set(ids)) != len(ids):
        duplicates = duplicate_values(ids)
        errors.append(f"Duplicate control IDs: {duplicates}")
    if len(set(normalized_controls)) != len(normalized_controls):
        errors.append("Duplicate normalized PRC control text found")

    for section, items in sorted(numbers.items()):
        expected = list(range(1, len(items) + 1))
        if items != expected:
            errors.append(f"PRC-{section:02d} control sequence is not contiguous: {items}")
    return ids, numbers, normalized_controls


def validate_engineering_pages(errors: list[str]) -> tuple[list[str], list[str]]:
    actual_pages = {path.name for path in ENGINEERING_DIR.glob("*.md")}
    if actual_pages != EXPECTED_ENGINEERING_PAGES:
        errors.append(
            "Engineering page set differs: "
            f"missing={sorted(EXPECTED_ENGINEERING_PAGES - actual_pages)}, "
            f"unexpected={sorted(actual_pages - EXPECTED_ENGINEERING_PAGES)}"
        )

    ids: list[str] = []
    normalized_controls: list[str] = []
    for path in sorted(ENGINEERING_DIR.glob("[0-9][0-9]-*.md")):
        if path.name == "00-overview.md":
            continue
        content = path.read_text(encoding="utf-8")
        for match in ENGINEERING_PLAIN_CHECKBOX.finditer(content):
            line = content.count("\n", 0, match.start()) + 1
            errors.append(
                f"{path.relative_to(ROOT)}:{line}: checklist item has no USEQ ID"
            )
        for match in ENGINEERING_CONTROL.finditer(content):
            control_id, text = match.groups()
            technology = IMPLEMENTATION_SPECIFIC.search(text)
            if technology:
                line = content.count("\n", 0, match.start()) + 1
                errors.append(
                    f"{path.relative_to(ROOT)}:{line}: {control_id} depends on "
                    f"a language, framework, or vendor ({technology.group(0)})"
                )
            normalized = normalize_control(text)
            ids.append(control_id)
            normalized_controls.append(normalized)

    if len(ids) != EXPECTED_ENGINEERING_CONTROLS:
        errors.append(
            f"Expected {EXPECTED_ENGINEERING_CONTROLS} USEQ controls, found {len(ids)}"
        )
    if len(set(ids)) != len(ids):
        errors.append("Duplicate USEQ control IDs found")
    if len(set(normalized_controls)) != len(normalized_controls):
        errors.append("Duplicate normalized USEQ control text found")
    return ids, normalized_controls


def validate_links(errors: list[str]) -> int:
    checked = 0
    repository_root = ROOT.resolve()
    for path in markdown_files():
        content = path.read_text(encoding="utf-8")
        for match in MARKDOWN_LINK.finditer(content):
            target = match.group(1).strip()
            if target.startswith(("http://", "https://", "mailto:")):
                continue
            path_part, _, fragment = target.partition("#")
            target_path = unquote(path_part.split("?", 1)[0])
            checked += 1
            resolved = path.resolve() if not target_path else (path.parent / target_path).resolve()
            if not resolved.is_relative_to(repository_root):
                line = content.count("\n", 0, match.start()) + 1
                errors.append(
                    f"{path.relative_to(ROOT)}:{line}: local link escapes repository {target!r}"
                )
                continue
            if not resolved.exists():
                line = content.count("\n", 0, match.start()) + 1
                errors.append(
                    f"{path.relative_to(ROOT)}:{line}: broken local link {target!r}"
                )
                continue
            if fragment and resolved.suffix.lower() == ".md":
                anchors = markdown_anchors(resolved.read_text(encoding="utf-8"))
                decoded_fragment = unquote(fragment).casefold()
                if decoded_fragment not in anchors:
                    line = content.count("\n", 0, match.start()) + 1
                    errors.append(
                        f"{path.relative_to(ROOT)}:{line}: missing Markdown anchor "
                        f"{fragment!r} in {resolved.relative_to(repository_root)}"
                    )
    return checked


def main() -> int:
    errors: list[str] = []
    ids, numbers, production_controls = validate_pages(errors)
    engineering_ids, engineering_controls = validate_engineering_pages(errors)
    cross_layer_duplicates = set(production_controls) & set(engineering_controls)
    if cross_layer_duplicates:
        errors.append(
            f"Lifecycle controls duplicate {len(cross_layer_duplicates)} production controls"
        )
    link_count = validate_links(errors)

    if errors:
        print("Validation failed:")
        for error in errors:
            print(f"- {error}")
        return 1

    sections = len(numbers)
    print(
        f"Validated {len(ids) + len(engineering_ids)} unique controls "
        f"({len(engineering_ids)} lifecycle and {len(ids)} production-readiness) "
        f"across {sections} numbered production sections and {link_count} local Markdown links."
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
