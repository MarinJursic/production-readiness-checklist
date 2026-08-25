#!/usr/bin/env python3
"""Generate a plain-language acceptance and scanner-coverage review for every control."""

from __future__ import annotations

import argparse
import json
import re
import sys
from collections import Counter, defaultdict
from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
REGISTRY = ROOT / "catalog" / "control-id-registry.json"
OUTPUT = ROOT / "research" / "CONTROL_ACCEPTANCE_CRITERIA_REVIEW.md"
OUTPUT_DIRECTORY = ROOT / "research" / "control-acceptance-criteria"
INDEX_OUTPUT = OUTPUT_DIRECTORY / "README.md"
MAX_PART_BYTES = 3_000_000
TARGET_PART_BYTES = 2_850_000

HEADING = re.compile(r"^(#{1,6})\s+(.+?)\s*$")

VAGUE_TERMS = (
    "appropriate",
    "appropriately",
    "adequate",
    "adequately",
    "sufficient",
    "sufficiently",
    "reasonable",
    "reasonably",
    "regular",
    "regularly",
    "timely",
    "material",
    "critical",
    "high-risk",
    "high risk",
    "proportionate",
    "relevant",
    "meaningful",
    "effective",
    "effectively",
    "as needed",
    "where feasible",
)

ABSOLUTE = re.compile(r"\b(?:all|every|always|never|none|only|complete|entire)\b", re.I)
CONDITIONAL = re.compile(
    r"\b(?:if|when|whenever|where applicable|as applicable|where required|"
    r"when required|unless|where appropriate|as appropriate)\b",
    re.I,
)
NEGATIVE = re.compile(
    r"^(?:do not|never|no\b|avoid\b|prevent\b|prohibit\b|reject\b|"
    r"a .+ (?:has not|cannot|does not|is not|remains)|an .+ (?:has not|cannot|does not|is not|remains))",
    re.I,
)
NUMBER_OR_VERSION = re.compile(
    r"\b(?:\d+(?:\.\d+){0,3}|RFC\s*\d+|ISO(?:/IEC)?\s*\d+|WCAG\s*\d+(?:\.\d+)?|"
    r"PCI DSS\s*v?\d+(?:\.\d+)*)\b",
    re.I,
)
IMPLEMENTATION_WORDS = re.compile(
    r"\b(?:folder|folders|directory|directories|filename|file name|exact path|"
    r"GitHub|GitLab|Bitbucket|Docker|Kubernetes|Terraform|Jenkins|CircleCI|"
    r"AWS|Azure|Google Cloud|GCP|Datadog|New Relic|Sentry|Stripe|Cloudflare|"
    r"Node\.js|npm|Python|Java|Rust|OpenAPI|PostgreSQL|MySQL|MongoDB|Redis|Kafka)\b",
    re.I,
)
HUMAN_OR_EXTERNAL = re.compile(
    r"\b(?:accountable|approval|approve|attestation|audit|contract|customer|executive|"
    r"governance|human|interview|legal|law|manager|owner|policy|regulation|regulatory|"
    r"reviewer|risk acceptance|stakeholder|training|trained|vendor|workforce)\b",
    re.I,
)
RUNTIME = re.compile(
    r"\b(?:alert|availability|backup|capacity|cluster|database|deployed|deployment|"
    r"disaster|environment|failover|incident|latency|load|monitor|network|on-call|"
    r"production|recovery|region|restore|runtime|service level|SLO|RPO|RTO|traffic)\b",
    re.I,
)
REPOSITORY = re.compile(
    r"\b(?:artifact|build|code|configuration|dependency|documentation|file|format|"
    r"interface|manifest|package|repository|schema|source|static analysis|test|workflow)\b",
    re.I,
)


def load_json(path: Path) -> dict[str, Any]:
    value = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(value, dict):
        raise ValueError(f"{path}: expected a JSON object")
    return value


def load_yaml(path: Path) -> dict[str, Any]:
    value = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(value, dict):
        raise ValueError(f"{path}: expected a YAML mapping")
    return value


def control_context(entries: list[dict[str, Any]]) -> dict[str, list[str]]:
    """Return the Markdown heading trail at each control's source line."""

    by_path: dict[str, list[dict[str, Any]]] = defaultdict(list)
    for entry in entries:
        by_path[entry["source"]["path"]].append(entry)

    result: dict[str, list[str]] = {}
    for relative, source_entries in by_path.items():
        requested = {entry["source"]["line"]: entry["id"] for entry in source_entries}
        headings: dict[int, str] = {}
        lines = (ROOT / relative).read_text(encoding="utf-8").splitlines()
        for line_number, line in enumerate(lines, start=1):
            match = HEADING.match(line)
            if match:
                level = len(match.group(1))
                headings[level] = match.group(2)
                headings = {key: value for key, value in headings.items() if key <= level}
            control_id = requested.get(line_number)
            if control_id:
                result[control_id] = [headings[key] for key in sorted(headings)]
    return result


def assertion_index() -> dict[str, list[dict[str, Any]]]:
    index: dict[str, list[dict[str, Any]]] = defaultdict(list)
    for path in sorted((ROOT / "catalog" / "assertions").glob("*.yaml")):
        for assertion in load_yaml(path)["assertions"]:
            for control_id in assertion["control_ids"]:
                index[control_id].append(assertion)
    for assertions in index.values():
        assertions.sort(key=lambda item: item["id"])
    return index


def matched_terms(statement: str, terms: tuple[str, ...]) -> list[str]:
    lowered = statement.casefold()
    return sorted({term for term in terms if term.casefold() in lowered})


def looks_compound(statement: str) -> bool:
    conjunctions = len(re.findall(r"\b(?:and|or)\b", statement, re.I))
    list_separators = statement.count(",") + statement.count(";")
    return conjunctions >= 2 or conjunctions >= 1 and list_separators >= 2 or list_separators >= 4


def flags_for(entry: dict[str, Any], headings: list[str]) -> dict[str, Any]:
    statement = entry["statement"]
    vague = matched_terms(statement, VAGUE_TERMS)
    heading_text = " ".join(headings)
    return {
        "compound": looks_compound(statement),
        "vague": vague,
        "absolute": bool(ABSOLUTE.search(statement)),
        "conditional": bool(CONDITIONAL.search(statement)) or "conditional" in heading_text.casefold(),
        "negative": bool(NEGATIVE.search(statement)) or "no-go" in heading_text.casefold(),
        "number_or_version": NUMBER_OR_VERSION.findall(statement),
        "implementation": IMPLEMENTATION_WORDS.findall(statement),
        "human": bool(HUMAN_OR_EXTERNAL.search(statement)),
        "runtime": bool(RUNTIME.search(statement)),
        "repository": bool(REPOSITORY.search(statement)),
    }


def proposed_check_class(flags: dict[str, Any]) -> str:
    if flags["human"] and (flags["runtime"] or flags["repository"]):
        return "Mixed evidence"
    if flags["human"]:
        return "Human or external evidence"
    if flags["runtime"] and flags["repository"]:
        return "Mixed repository and environment evidence"
    if flags["runtime"]:
        return "Authorized environment check"
    if flags["repository"]:
        return "Repository or artifact check"
    return "Mixed evidence"


def applicability_text(headings: list[str], flags: dict[str, Any]) -> str:
    context = " > ".join(headings) if headings else "No source heading was found"
    if flags["conditional"]:
        return (
            f"Use this rule only when its stated trigger and the source section `{context}` apply. "
            "The scanner must record the trigger evidence. It may return Not Applicable only when "
            "the trigger is proven absent and a reason is saved."
        )
    return (
        f"Consider this rule inside the selected profile and the source section `{context}`. "
        "Before checking it, the scanner must name the exact product part, release, environment, "
        "data, and people it covers. A broad heading does not make the rule apply to every project."
    )


def acceptance_items(statement: str, flags: dict[str, Any]) -> list[str]:
    items = [
        "The exact assessment target and the complete applicable scope are recorded.",
    ]
    if flags["negative"]:
        items.append(
            "The unsafe condition named by the rule is shown to be absent. For a no-go rule, "
            "missing proof is Blocked, not Pass."
        )
    else:
        items.append(f"Current trustworthy evidence shows that this statement is true: “{statement}”")
    if flags["compound"]:
        items.append(
            "Each separate promise in the sentence is turned into its own atomic assertion, and "
            "every required assertion passes. One passing part cannot hide a failing part."
        )
    if flags["absolute"]:
        items.append(
            "Words such as all, every, never, or only require a complete inventory of the relevant "
            "items. A sample is not enough unless the rule itself defines sampling."
        )
    if flags["vague"]:
        terms = ", ".join(f"`{term}`" for term in flags["vague"])
        items.append(
            f"Before evaluation, the project or selected profile replaces {terms} with measurable "
            "limits, named owners, and evidence freshness. The scanner must not invent one global value."
        )
    if flags["number_or_version"]:
        items.append(
            "Any named number, date, or standard version is confirmed to be the selected current "
            "requirement for this project; a newer or stricter binding rule wins."
        )
    items.extend(
        [
            "The evidence belongs to the exact source revision, artifact, configuration, environment, "
            "and time being assessed.",
            "No equally strong or stronger evidence contradicts the result. If scope or proof is "
            "missing, stale, unsafe to collect, or conflicting, the result stays Blocked or Unknown.",
        ]
    )
    return items


def proposed_method(check_class: str, flags: dict[str, Any]) -> list[str]:
    common = [
        "First split the rule into one assertion per observable promise and bind each assertion to a "
        "versioned checker and an allowed evidence type.",
    ]
    if check_class == "Repository or artifact check":
        common.append(
            "Build an inventory from file contents, manifests, parsers, build graphs, and explicit "
            "project configuration. Common file or folder names may be discovery hints, but they must "
            "not be the only way to pass."
        )
        common.append(
            "Run a bounded, read-only parser or analyzer. Save the exact files, artifact digests, "
            "checker version, limits, and observed result. Do not run the target project's scripts."
        )
    elif check_class == "Authorized environment check":
        common.append(
            "Require an explicit, least-privilege connector for the named non-production or production "
            "environment. A normal repository scan cannot pass this rule."
        )
        common.append(
            "Collect a bounded query, observation, or non-destructive test result and bind it to the "
            "exact environment, release, time, and project-defined threshold."
        )
    elif check_class == "Human or external evidence":
        common.append(
            "Ask for a dated evidence record from the accountable person or authority. Validate its "
            "required fields, signature or trusted source, scope, freshness, and expiry."
        )
        common.append(
            "The scanner may check that evidence exists and is well formed. It must not replace legal, "
            "risk, accessibility, safety, management, or release decisions with an AI guess."
        )
    elif check_class in {"Mixed evidence", "Mixed repository and environment evidence"}:
        common.append(
            "Use separate repository, artifact, environment, and human assertions as needed. Keep the "
            "authority of each evidence source visible and require every applicable part."
        )
        common.append(
            "Repository intent cannot prove deployed behavior, and a runtime snapshot cannot prove the "
            "source or review process by itself."
        )
    else:
        common.append(
            "A control author must define the trigger, atomic assertions, acceptable evidence, and "
            "failure meaning before this rule can enter a machine-run profile. Until then it remains a "
            "human checklist item."
        )
    if flags["implementation"]:
        common.append(
            "Keep any named technology, standard, path, or directory inside a conditional adapter or "
            "profile. Provide an equivalent evidence route for other project designs when the objective "
            "is broader than that implementation."
        )
    return common


def review_notes(flags: dict[str, Any]) -> list[str]:
    notes = []
    if flags["compound"]:
        notes.append("May bundle several promises; split it before machine evaluation.")
    if flags["vague"]:
        notes.append(
            "Contains project-dependent words: " + ", ".join(f"`{term}`" for term in flags["vague"])
            + ". Define them in the selected profile or project policy."
        )
    if flags["absolute"]:
        notes.append("Uses an absolute word; passing requires proof that the relevant inventory is complete.")
    if flags["conditional"]:
        notes.append("Has a condition or sits in a conditional section; make the trigger explicit and testable.")
    if flags["negative"]:
        notes.append("Is negative or no-go wording; do not turn missing evidence into proof of absence.")
    if flags["number_or_version"]:
        notes.append(
            "Names a number or version: " + ", ".join(f"`{value}`" for value in flags["number_or_version"])
            + ". Add an owner and review date so it cannot silently become stale."
        )
    if flags["implementation"]:
        notes.append(
            "Names a possible implementation detail: "
            + ", ".join(f"`{value}`" for value in flags["implementation"])
            + ". Keep it conditional; do not force the same design on every project."
        )
    if flags["human"]:
        notes.append("Likely needs accountable human or external evidence; it should not be auto-passed from source files.")
    if not notes:
        notes.append(
            "No automatic wording warning was found. A control author still needs to approve its "
            "applicability, evidence, and atomic assertions before automation."
        )
    return notes


def current_scanner_text(assertions: list[dict[str, Any]]) -> list[str]:
    if not assertions:
        return [
            "Included in every complete `prc scan` report as `needs_review`, but not deterministically "
            "checked today. No missing implementation is turned into a Pass.",
            "With an explicitly enabled Codex or Claude review, the sealed task requires the coordinator to assign this rule to one separate subagent "
            "plus bounded, secret-screened repository excerpts. Its candidate, evidence, advice, and "
            "limitations are advisory only; they cannot create a verified Pass or final Not Applicable result."
        ]
    lines = [
        f"Included in the complete report and partly represented by {len(assertions)} catalog assertion(s). "
        "A passing narrow assertion produces `partially_verified`, not a complete Pass for this broad "
        "control. A failing narrow assertion can produce `confirmed_failure`:"
    ]
    for assertion in assertions:
        evidence = "; ".join(
            f"{item['kind']} ({item['minimum_authority']}): {item['description']}"
            for item in assertion["evidence_required"]
        )
        lines.append(
            f"`{assertion['id']}` — {assertion['statement']} Applies when "
            f"`{assertion['applicability']}`. Implementation: `{assertion['implementation_id']}`. "
            f"Required evidence: {evidence}"
        )
    return lines


def summary_counts(all_flags: dict[str, dict[str, Any]]) -> Counter[str]:
    counts: Counter[str] = Counter()
    for flags in all_flags.values():
        for name in ("compound", "absolute", "conditional", "negative", "human"):
            counts[name] += bool(flags[name])
        counts["vague"] += bool(flags["vague"])
        counts["number_or_version"] += bool(flags["number_or_version"])
        counts["implementation"] += bool(flags["implementation"])
    return counts


def generated_text() -> str:
    registry = load_json(REGISTRY)
    entries = registry["entries"]
    contexts = control_context(entries)
    assertions = assertion_index()
    all_flags = {
        entry["id"]: flags_for(entry, contexts.get(entry["id"], [])) for entry in entries
    }
    counts = summary_counts(all_flags)
    mapped_controls = sum(1 for entry in entries if entry["id"] in assertions)
    assertion_count = sum(len(values) for values in assertions.values())

    lines = [
        "<!-- Generated by scripts/generate_control_acceptance_review.py. Do not edit by hand. -->",
        "",
        "# Review of all production-readiness control acceptance criteria",
        "",
        f"- Registry: `{registry['schema_version']}` / `{registry['registry_version']}`",
        f"- Source digest: `{registry['source_sha256']}`",
        f"- Controls reviewed: **{len(entries):,}**",
        f"- Controls connected to current scanner assertions: **{mapped_controls:,}**",
        f"- Current assertions connected to those controls: **{assertion_count:,}**",
        "",
        "## Read this first: the important problems found",
        "",
        "This file reviews every control, but it does not claim that every control is already machine-checkable. "
        "Every `prc scan` report now contains all 10,042 controls. The deterministic profile proves only its "
        "narrow assertions; the remaining controls stay visibly `needs_review`. An optional AI review can add "
        "a one-subagent-per-control review task, but current provider output cannot prove every internal subagent call happened and no advice can turn subjective judgment into a verified Pass. Turning "
        "broad objectives directly into hard-coded rules would create false passes and one-size-fits-all designs.",
        "",
        "### The rule design that should be used everywhere",
        "",
        "1. A broad control is an objective, not a test.",
        "2. Each objective must be split into small assertions that each prove one thing.",
        "3. Applicability must come from detected facts plus reviewed project configuration. It must never be guessed.",
        "4. Common names such as `src`, `components`, `tests`, `README.md`, or `Dockerfile` may help discovery, but "
        "they cannot be the only accepted project structure.",
        "5. The scanner should discover components from file contents, manifests, parsers, build graphs, and "
        "explicit project configuration. A user may point to an unusual layout without weakening the rule.",
        "6. Words such as “appropriate”, “sufficient”, and “regularly” need a project or profile value before a "
        "machine can judge them. The scanner must not invent one value for all projects.",
        "7. Repository files prove repository facts. They do not prove production behavior, legal compliance, "
        "backup recovery, staff readiness, or a human approval.",
        "8. Missing, stale, conflicting, unsafe-to-collect, or incomplete evidence is Blocked or Unknown, never Pass.",
        "9. A Not Applicable result needs a recorded trigger check and a reason.",
        "10. Every Pass must name the exact revision, artifact, configuration, environment, checker version, time, "
        "and evidence used.",
        "",
        "### Automatic wording review",
        "",
        "These counts are review warnings, not proof that the rules are wrong. One rule can have several warnings.",
        "",
        "| Warning | Controls | Why it needs attention |",
        "| --- | ---: | --- |",
        f"| Several promises may be bundled together | {counts['compound']:,} | Split them so one passing part cannot hide a failure. |",
        f"| Vague or project-dependent words | {counts['vague']:,} | Define a measurable project or profile value. |",
        f"| Absolute words such as all, every, or never | {counts['absolute']:,} | A Pass needs a complete inventory, not a sample. |",
        f"| Conditional wording or conditional section | {counts['conditional']:,} | Make the trigger explicit and prove Not Applicable. |",
        f"| Negative or no-go wording | {counts['negative']:,} | Missing evidence must not be mistaken for absence. |",
        f"| Human, legal, policy, owner, or external evidence signal | {counts['human']:,} | Keep accountable decisions outside automatic source checks. |",
        f"| Number, date, RFC, or standard version | {counts['number_or_version']:,} | Add a review date and update route. |",
        f"| Possible technology, vendor, filename, path, or directory detail | {counts['implementation']:,} | Keep it conditional and allow equivalent evidence. |",
        "",
        "### Main problems in the current executable layer",
        "",
        "The current 43 catalog assertions are much safer than pretending all 10,042 rules are deterministic, but several still "
        "need refinement before a broad release:",
        "",
        "- Presence-only checks for README, license, contribution, security, conduct, ownership, manifests, and "
        "lock files can pass empty-quality or misleading content. Add content/schema checks without forcing one name.",
        "- Test discovery, CI discovery, dependency updates, runtime declarations, and ownership rely on a limited "
        "set of known files or conventions. Add detector adapters and an explicit reviewed configuration route.",
        "- The six-hour workflow timeout is a universal fixed value. Make the maximum a named profile or project policy.",
        "- The blanket `pull_request_target` ban is safe for a hardened profile but rejects secure designs. A later "
        "data-flow check can distinguish untrusted checkout or execution from safe metadata-only use.",
        "- Final-newline enforcement is style, not a universal production gate. Keep it advisory or in a style profile.",
        "- Local Unix group/other write bits do not reliably describe repository policy across Windows and different "
        "checkout tools. Treat that result as local-workspace evidence only.",
        "- Terraform locks should apply to runnable root modules, not every directory containing a `.tf` file.",
        "- Kubernetes resource limits, capability allowlists, seccomp choices, and explicit security context fields "
        "are policy-specific. Put strict defaults in named hardened profiles and allow reviewed equivalent controls.",
        "- Static container and Kubernetes files cannot prove the effective deployed image, admission policy, or "
        "runtime identity. Add optional environment evidence without weakening the repository result.",
        "- The Go HTTP helper checks are conservative patterns and can flag code that configures global defaults "
        "safely. Keep the finding honest and add data-flow evidence rather than silent suppression.",
        "- “Applicable analyses have executed” is broader than its current secret-scanner binding. Split each "
        "analysis class into its own assertion so the title cannot overclaim.",
        "",
        "### How to read each control below",
        "",
        "Each entry keeps the exact source rule, gives a proposed pass contract, states what `prc scan` really "
        "checks today, and proposes a future method. The proposal is a review starting point. It must be approved "
        "by a control owner and tested against positive, negative, unusual-layout, and adversarial fixtures before "
        "it becomes executable.",
        "",
        "## All controls",
        "",
    ]

    previous_source = ""
    for entry in sorted(entries, key=lambda item: (item["source"]["path"], item["source"]["line"], item["id"])):
        source = entry["source"]
        if source["path"] != previous_source:
            lines.extend([f"## `{source['path']}`", ""])
            previous_source = source["path"]
        context = contexts.get(entry["id"], [])
        flags = all_flags[entry["id"]]
        check_class = proposed_check_class(flags)
        lines.extend(
            [
                f"<!-- BEGIN CONTROL {entry['id']} -->",
                f"### {entry['id']}",
                "",
                f"- Source: `{source['path']}:{source['line']}`",
                f"- Section: `{' > '.join(context) if context else 'unknown'}`",
                f"- Status/revision: `{entry['status']}` / `{entry['revision']}`",
                "",
                f"> {entry['statement']}",
                "",
                "**When it applies**",
                "",
                applicability_text(context, flags),
                "",
                "**Proposed acceptance criteria**",
                "",
            ]
        )
        lines.extend(f"- {item}" for item in acceptance_items(entry["statement"], flags))
        lines.extend(["", "**What the program checks today**", ""])
        lines.extend(f"- {item}" for item in current_scanner_text(assertions.get(entry["id"], [])))
        lines.extend(["", f"**Proposed future check: {check_class}**", ""])
        lines.extend(f"- {item}" for item in proposed_method(check_class, flags))
        lines.extend(["", "**Review notes**", ""])
        lines.extend(f"- {item}" for item in review_notes(flags))
        lines.extend(["", f"<!-- END CONTROL {entry['id']} -->", ""])

    return "\n".join(lines).rstrip() + "\n"


def split_review(text: str) -> dict[Path, str]:
    """Return a small compatibility page, one index, and bounded control parts."""

    divider = "## All controls\n\n"
    preamble, found, controls_text = text.partition(divider)
    if not found:
        raise ValueError("generated review has no all-controls divider")
    block_pattern = re.compile(
        r"<!-- BEGIN CONTROL ([A-Z0-9-]+) -->\n.*?<!-- END CONTROL \1 -->\n",
        re.DOTALL,
    )
    source_pattern = re.compile(r"^- Source: `([^`]+):\d+`$", re.MULTILINE)
    blocks: list[tuple[str, str, str]] = []
    for match in block_pattern.finditer(controls_text):
        block = match.group(0)
        source = source_pattern.search(block)
        if source is None:
            raise ValueError(f"control {match.group(1)} has no source path")
        blocks.append((match.group(1), source.group(1), block))
    if not blocks:
        raise ValueError("generated review has no control blocks")

    groups: list[list[tuple[str, str, str]]] = []
    current: list[tuple[str, str, str]] = []
    current_bytes = 2048
    previous_source = ""
    for control_id, source, block in blocks:
        addition = len(block.encode("utf-8"))
        if source != previous_source:
            addition += len(f"## `{source}`\n\n".encode("utf-8"))
        if current and current_bytes + addition > TARGET_PART_BYTES:
            groups.append(current)
            current = []
            current_bytes = 2048
            previous_source = ""
            addition = len(block.encode("utf-8")) + len(f"## `{source}`\n\n".encode("utf-8"))
        current.append((control_id, source, block))
        current_bytes += addition
        previous_source = source
    if current:
        groups.append(current)

    files: dict[Path, str] = {}
    part_rows: list[str] = []
    total_parts = len(groups)
    for index, group in enumerate(groups, start=1):
        relative = Path(f"part-{index:03d}.md")
        part_lines = [
            "<!-- Generated by scripts/generate_control_acceptance_review.py. Do not edit by hand. -->",
            "",
            f"# Control acceptance criteria — part {index} of {total_parts}",
            "",
            "[Back to the acceptance-criteria index](README.md).",
            "",
            f"Controls in this file: **{len(group):,}** (`{group[0][0]}` through `{group[-1][0]}` in source order).",
            "",
        ]
        previous_source = ""
        for _control_id, source, block in group:
            if source != previous_source:
                part_lines.extend([f"## `{source}`", ""])
                previous_source = source
            part_lines.extend([block.rstrip(), ""])
        part_text = "\n".join(part_lines).rstrip() + "\n"
        part_bytes = len(part_text.encode("utf-8"))
        if part_bytes >= MAX_PART_BYTES:
            raise ValueError(f"{relative} is {part_bytes} bytes; limit is below {MAX_PART_BYTES}")
        files[relative] = part_text
        part_rows.append(
            f"| [{relative.name}]({relative.name}) | {len(group):,} | `{group[0][0]}` | "
            f"`{group[-1][0]}` | {part_bytes:,} |"
        )

    index_lines = [
        preamble.rstrip(),
        "",
        "## Split control files",
        "",
        "The per-control review is split at control boundaries so it stays easy to open, diff, and review. "
        f"Every part is below **{MAX_PART_BYTES:,} bytes**. The order is the same as the registry source order.",
        "",
        "| File | Controls | First control | Last control | Bytes |",
        "| --- | ---: | --- | --- | ---: |",
        *part_rows,
        "",
    ]
    files[Path("README.md")] = "\n".join(index_lines).rstrip() + "\n"
    files[Path("../CONTROL_ACCEPTANCE_CRITERIA_REVIEW.md")] = (
        "<!-- Generated by scripts/generate_control_acceptance_review.py. Do not edit by hand. -->\n\n"
        "# Review of all production-readiness control acceptance criteria\n\n"
        "The complete 10,042-control review is split into files smaller than 3 MB. "
        "Open the [acceptance-criteria index](control-acceptance-criteria/README.md) for the method, "
        "findings, and links to every part.\n"
    )
    return files


def generated_files() -> dict[Path, str]:
    return split_review(generated_text())


def generate() -> None:
    OUTPUT_DIRECTORY.mkdir(parents=True, exist_ok=True)
    files = generated_files()
    expected_parts = {path.name for path in files if path.name.startswith("part-")}
    for stale in OUTPUT_DIRECTORY.glob("part-*.md"):
        if stale.name not in expected_parts:
            stale.unlink()
    for relative, content in files.items():
        destination = OUTPUT_DIRECTORY / relative
        destination.parent.mkdir(parents=True, exist_ok=True)
        destination.write_text(content, encoding="utf-8")
    print(f"generated {len(files) - 2} bounded parts and {INDEX_OUTPUT.relative_to(ROOT)}")


def check() -> None:
    files = generated_files()
    for relative, expected in files.items():
        destination = OUTPUT_DIRECTORY / relative
        if not destination.exists() or destination.read_text(encoding="utf-8") != expected:
            raise ValueError(
                f"{destination.relative_to(ROOT)} is stale; run "
                "python3 scripts/generate_control_acceptance_review.py generate"
            )
    expected_parts = {path.name for path in files if path.name.startswith("part-")}
    actual_parts = {path.name for path in OUTPUT_DIRECTORY.glob("part-*.md")}
    if actual_parts != expected_parts:
        raise ValueError("control acceptance review has unexpected or missing part files")
    print(f"verified {len(expected_parts)} bounded control-review parts")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("command", choices=("generate", "check"))
    args = parser.parse_args()
    try:
        generate() if args.command == "generate" else check()
    except (OSError, ValueError, json.JSONDecodeError, yaml.YAMLError) as error:
        print(f"control review generation failed: {error}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
