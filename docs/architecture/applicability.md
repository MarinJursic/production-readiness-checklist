# Bounded applicability evaluation

Applicability decides whether an assertion belongs in a specific assessment. It
does not decide whether the assertion passes. The scanner evaluates catalog
applicability with Common Expression Language (CEL) and records the outcome,
evaluator identity, and reason in the immutable scan plan.

## Fail-closed contract

Every expression must compile, evaluate within the configured resource budget,
and return a Boolean. The result mapping is deliberately small:

| CEL outcome | Planned applicability |
| --- | --- |
| `true` | `applicable` |
| `false` | `not_applicable` |
| Compile error, missing field, non-Boolean result, or evaluation error | `undetermined` |
| Parser, recursion, expression-size, or runtime-cost limit reached | `undetermined` |

`undetermined` never becomes Not Applicable or Pass. A required assertion with
undetermined applicability keeps the assessment incomplete.

Plan v0.3 records `applicability_reason` for every planned assertion and
identifies the evaluator as `cel-go/v0.30.0+prc-inventory/v0.3`. Frozen v0.1 and
v0.2 plan schemas remain available for validating archived plans.

## Available inventory view

Expressions receive only a deterministic projection named `inventory`. Target
file bytes, the target root path, environment variables, credentials, and
process or network capabilities are not exposed.

The projection contains:

- `file_count` and `source_files`;
- `package_ecosystems`, `manifests`, `lock_files`, `container_files`, and
  `symlinks`;
- `ci.github_actions` and `ci.workflow_files`;
- `infrastructure.terraform_files` and
  `infrastructure.kubernetes_files`;
- `components`, limited to component ID, kind, path, and ecosystem; and
- `fact_values`, mapping sourced inventory fact keys to their string values; and
- `declared`, containing the bound project ID, risk profile, profile, release
  scope, features, and data-context tags, or `configured: false` when absent.

Examples:

```text
inventory.package_ecosystems.size() > 0
inventory.ci.github_actions == true
inventory.components.exists(c, c.kind == "container-build")
inventory.components.exists(c, c.kind == "api-description" && c.ecosystem == "openapi")
```

Inventory detection remains evidence with stated limitations. An expression can
select assertions using detected facts; it cannot turn those facts into proof of
runtime behavior or deployment state.

## Resource and capability boundary

CEL is used as a non-Turing-complete, side-effect-free expression evaluator. The
scanner additionally limits an expression to 4,096 source characters, parser
recursion depth 64, and 10,000 runtime cost units, with periodic interruption
checks. Compiled programs are cached by exact expression for deterministic,
concurrent reuse.

No custom functions perform file, process, network, clock, random, or environment
access. Catalog authors cannot use applicability to execute a target project or
an adapter. Applicability expressions should be short predicates over the
documented projection; assertion implementations remain responsible for
collecting and evaluating evidence.

## Design references

- The [CEL project](https://github.com/cel-expr/cel-go) describes CEL as a
  non-Turing-complete, type-checkable, side-effect-free expression language and
  documents program reuse.
- The [cel-go API](https://pkg.go.dev/github.com/google/cel-go/cel) documents
  parser recursion and expression-size limits, runtime cost limits, and
  interruption checks used by this implementation.
