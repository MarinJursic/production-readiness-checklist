#!/usr/bin/env python3
"""Split the master checklist into the topic pages used by the documentation site."""

from __future__ import annotations

import argparse
import re
from dataclasses import dataclass
from pathlib import Path


@dataclass(frozen=True)
class Track:
    filename: str
    title: str
    description: str
    first_section: int
    last_section: int


TRACKS = (
    Track(
        "01-release-foundations.md",
        "Release foundations",
        "Define what readiness means, stop unsafe launches early, and identify the exact release under review.",
        1,
        3,
    ),
    Track(
        "02-product-risk-architecture.md",
        "Product, risk, and architecture",
        "Align business intent, accountable ownership, risk tolerance, and system understanding.",
        4,
        6,
    ),
    Track(
        "03-source-build-supply-chain.md",
        "Source, build, and supply chain",
        "Protect source changes, build integrity, dependencies, provenance, and licensing.",
        7,
        9,
    ),
    Track(
        "04-environments-quality-experience.md",
        "Environments, quality, and experience",
        "Verify configuration, functional correctness, testing, frontend behavior, and accessibility.",
        10,
        14,
    ),
    Track(
        "05-application-security.md",
        "Application security",
        "Review interfaces, identity, authorization, sessions, untrusted input, files, transport, and cryptography.",
        15,
        22,
    ),
    Track(
        "06-data-privacy-performance.md",
        "Data, privacy, and performance",
        "Protect data integrity and privacy while proving user-experience and capacity targets.",
        23,
        26,
    ),
    Track(
        "07-reliability-operations.md",
        "Reliability and operations",
        "Prepare infrastructure, observability, response, recovery, deployment, and post-launch verification.",
        27,
        35,
    ),
    Track(
        "08-maintenance-vendors-compliance.md",
        "Maintenance, vendors, and compliance",
        "Keep the service supportable and verify third-party, legal, and regulatory obligations.",
        36,
        38,
    ),
    Track(
        "09-conditional-modules.md",
        "Conditional feature modules",
        "Apply deeper controls for payments, multi-tenancy, AI, real-time systems, and other triggered features.",
        39,
        39,
    ),
    Track(
        "10-evidence-and-decision.md",
        "Evidence, sign-off, and decision",
        "Assemble the evidence package, obtain accountable sign-off, and make the final go/no-go decision.",
        40,
        43,
    ),
)

TOP_LEVEL_SECTION = re.compile(r"^# (\d+)\. ", re.MULTILINE)


def demote_headings(markdown: str) -> str:
    """Make original top-level headings fit beneath each track's page title."""
    lines = []
    for line in markdown.splitlines():
        if line.startswith("#"):
            line = f"#{line}"
        lines.append(line)
    return "\n".join(lines).strip()


def add_control_ids(markdown: str) -> str:
    """Attach stable, human-readable IDs to imported checklist items."""
    current_section: int | None = None
    control_number = 0
    rendered: list[str] = []

    for line in markdown.splitlines():
        section_match = re.match(r"^# (\d+)\. ", line)
        if section_match:
            current_section = int(section_match.group(1))
            control_number = 0
        if line.startswith("- [ ] "):
            if current_section is None:
                raise ValueError("Checklist item found before a numbered section")
            control_number += 1
            control_id = f"PRC-{current_section:02d}-{control_number:03d}"
            line = line.replace("- [ ] ", f"- [ ] **{control_id}** — ", 1)
        rendered.append(line)

    return "\n".join(rendered)


def section_offsets(source: str) -> dict[int, int]:
    offsets: dict[int, int] = {}
    for match in TOP_LEVEL_SECTION.finditer(source):
        offsets[int(match.group(1))] = match.start()
    expected = set(range(1, 44))
    if set(offsets) != expected:
        missing = sorted(expected - set(offsets))
        extra = sorted(set(offsets) - expected)
        raise ValueError(f"Expected sections 1-43; missing={missing}, extra={extra}")
    return offsets


def render_track(track: Track, source: str, offsets: dict[int, int]) -> str:
    start = offsets[track.first_section]
    end = offsets.get(track.last_section + 1, len(source))
    content = demote_headings(add_control_ids(source[start:end]))
    return (
        f"# {track.title}\n\n"
        f"> {track.description}\n\n"
        f"Sections {track.first_section}–{track.last_section} of the master checklist. "
        "For each applicable item, capture status, owner, evidence, release, environment, reviewer, and evidence age.\n\n"
        f"{content}\n"
    )


def render_intro(source: str, offsets: dict[int, int]) -> str:
    intro = source[: offsets[1]].strip()
    intro = re.sub(
        r"^# Universal Web Application Production-Readiness Master Checklist\n+",
        "",
        intro,
    )
    return intro.strip()


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("source", type=Path, help="Path to the master Markdown checklist")
    parser.add_argument(
        "--output",
        type=Path,
        default=Path("docs/checklists"),
        help="Directory for generated topic pages",
    )
    args = parser.parse_args()

    source = args.source.read_text(encoding="utf-8")
    offsets = section_offsets(source)
    args.output.mkdir(parents=True, exist_ok=True)

    for track in TRACKS:
        rendered = render_track(track, source, offsets)
        (args.output / track.filename).write_text(rendered, encoding="utf-8")

    intro = render_intro(source, offsets)
    (args.output / "00-readiness-principle.md").write_text(
        "# The readiness principle\n\n"
        "> Production readiness is an evidence-backed risk decision, not a promise of perfection.\n\n"
        f"{intro}\n",
        encoding="utf-8",
    )

    count = source.count("- [ ]")
    print(f"Generated {len(TRACKS) + 1} pages containing {count} checklist items.")


if __name__ == "__main__":
    main()
