#!/usr/bin/env python3
"""Validate checklist integrity and local Markdown links without dependencies."""

from __future__ import annotations

import re
import sys
from collections import defaultdict
from pathlib import Path
from urllib.parse import unquote


ROOT = Path(__file__).resolve().parents[1]
CHECKLIST_DIR = ROOT / "docs" / "checklists"
EXPECTED_CONTROLS = 1421
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
CONTROL = re.compile(r"^- \[ \] \*\*(PRC-(\d{2})-(\d{3}))\*\* — .+", re.MULTILINE)
PLAIN_CHECKBOX = re.compile(r"^- \[ \] (?!\*\*PRC-)", re.MULTILINE)
MARKDOWN_LINK = re.compile(r"!?\[[^\]]*\]\(([^)]+)\)")


def markdown_files() -> list[Path]:
    return sorted(
        path
        for path in ROOT.rglob("*.md")
        if ".git" not in path.parts and "site" not in path.parts
    )


def validate_pages(errors: list[str]) -> tuple[list[str], dict[int, list[int]]]:
    actual_pages = {path.name for path in CHECKLIST_DIR.glob("*.md")}
    if actual_pages != EXPECTED_PAGES:
        errors.append(
            "Checklist page set differs: "
            f"missing={sorted(EXPECTED_PAGES - actual_pages)}, "
            f"unexpected={sorted(actual_pages - EXPECTED_PAGES)}"
        )

    ids: list[str] = []
    numbers: dict[int, list[int]] = defaultdict(list)
    for path in sorted(CHECKLIST_DIR.glob("*.md")):
        content = path.read_text(encoding="utf-8")
        for match in PLAIN_CHECKBOX.finditer(content):
            line = content.count("\n", 0, match.start()) + 1
            errors.append(f"{path.relative_to(ROOT)}:{line}: checklist item has no PRC ID")
        for match in CONTROL.finditer(content):
            control_id, section, item = match.groups()
            ids.append(control_id)
            numbers[int(section)].append(int(item))

    if len(ids) != EXPECTED_CONTROLS:
        errors.append(f"Expected {EXPECTED_CONTROLS} controls, found {len(ids)}")
    if len(set(ids)) != len(ids):
        duplicates = sorted({control_id for control_id in ids if ids.count(control_id) > 1})
        errors.append(f"Duplicate control IDs: {duplicates}")

    for section, items in sorted(numbers.items()):
        expected = list(range(1, len(items) + 1))
        if items != expected:
            errors.append(f"PRC-{section:02d} control sequence is not contiguous: {items}")
    return ids, numbers


def validate_links(errors: list[str]) -> int:
    checked = 0
    for path in markdown_files():
        content = path.read_text(encoding="utf-8")
        for match in MARKDOWN_LINK.finditer(content):
            target = match.group(1).strip()
            if target.startswith(("http://", "https://", "mailto:", "#")):
                continue
            target_path = unquote(target.split("#", 1)[0].split("?", 1)[0])
            if not target_path:
                continue
            checked += 1
            resolved = (path.parent / target_path).resolve()
            if not resolved.exists():
                line = content.count("\n", 0, match.start()) + 1
                errors.append(
                    f"{path.relative_to(ROOT)}:{line}: broken local link {target!r}"
                )
    return checked


def main() -> int:
    errors: list[str] = []
    ids, numbers = validate_pages(errors)
    link_count = validate_links(errors)

    if errors:
        print("Validation failed:")
        for error in errors:
            print(f"- {error}")
        return 1

    sections = len(numbers)
    print(
        f"Validated {len(ids)} unique controls across {sections} numbered sections "
        f"and {link_count} local Markdown links."
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
