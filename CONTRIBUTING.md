# Contributing

Thanks for helping make production-readiness reviews clearer and more useful.

## Good contributions

- correct an inaccurate or outdated control;
- make a control more testable without weakening it;
- add an evidence example or reviewer prompt;
- improve accessibility, navigation, or adoption guidance;
- add a missing conditional module or standards mapping;
- fix broken links, automation, or documentation rendering.

Open an issue before a large structural change so contributors can agree on scope.

## Control-writing guidelines

A control should:

- state one verifiable outcome;
- remain technology-neutral unless it belongs to a specific module;
- describe the required property, not one vendor’s implementation;
- avoid claiming that any control guarantees security, reliability, or compliance;
- make conditional triggers explicit;
- use normative language consistently;
- include a primary-source reference when it depends on a standard or law.

Do not renumber existing `PRC-*` IDs. New controls receive the next unused number in their numbered section. If a control is retired, preserve its ID in the changelog rather than assigning it to a different requirement.

## Development workflow

1. Fork the repository and create a focused branch.
2. Make the smallest coherent change.
3. Run the integrity validator:

   ```bash
   python3 scripts/validate.py
   ```

4. Build the documentation site when navigation or rendering changes:

   ```bash
   python3 -m venv .venv
   . .venv/bin/activate
   pip install -r requirements-docs.txt
   mkdocs build --strict
   ```

5. Open a pull request using the repository template.

The `scripts/import_source.py` script records the one-time transformation from the original master document. It is not the normal editing workflow and overwrites generated checklist pages; contributors should edit the structured pages directly.

## Pull requests

Explain the problem, the control IDs affected, why the wording is correct, and how you validated it. Include primary references for standards-dependent changes. Keep unrelated formatting out of substantive control changes so reviewers can evaluate meaning precisely.

By contributing, you agree that your contribution is licensed under the MIT License.
