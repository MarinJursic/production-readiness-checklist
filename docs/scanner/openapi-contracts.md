# OpenAPI contract analysis

`PRC-A-API-001` is a language-neutral, no-execution check for the bounded root
structure of detected OpenAPI YAML and JSON documents. Inventory recognizes the
conventional `openapi.yaml`, `openapi.yml`, and `openapi.json` names, plus
bounded YAML files with a top-level OpenAPI 3.x marker. Each detection is
recorded as an `api-description` component and a sourced fact with an explicit
limitation: a description does not prove that its API is implemented, reachable,
or deployed.

## What the rule proves

For published OpenAPI 3.0, 3.1, and 3.2 feature versions, the native check
verifies that:

- the file contains exactly one parseable YAML or JSON document with an object
  root;
- `openapi` is a nonempty supported semantic version string;
- `info` is an object with nonempty string `title` and `version` fields;
- OpenAPI 3.0 has a `paths` object; and
- OpenAPI 3.1 and 3.2 have at least one object-valued `paths`, `components`, or
  `webhooks` field.

These are requirements from the authoritative
[OpenAPI 3.0.4 specification](https://spec.openapis.org/oas/v3.0.4.html),
[OpenAPI 3.1.2 specification](https://spec.openapis.org/oas/v3.1.2.html), and
[OpenAPI 3.2.0 specification](https://spec.openapis.org/oas/v3.2.0.html).
OpenAPI 3.1 and 3.2 descriptions may legitimately describe only reusable
components or webhooks, while 3.0 requires `paths`; an empty Paths Object is
allowed by the specifications.

Duplicate mapping keys, non-string mapping keys, missing required
metadata, and incorrect root-field types produce a finding with bounded source
locations. Invalid syntax produces Error/Unknown. A syntactically valid but
unsupported feature version, such as a future 3.3 document, also produces
Error/Unknown instead of a false Pass or Fail; a new benchmarked implementation
version is required before the scanner claims support.

## Bounds and evidence

The implementation inspects at most 256 detected documents and 64 MiB in total.
The shared native reader caps each parsed file at 4 MiB and verifies its bytes
against the content-addressed inventory before parsing. YAML shape inspection is
bounded to 100,000 nodes and 128 levels. The scanner records file-hash evidence,
never document content, and reports at most 100 structural problems.

## Deliberate limitations

This is not a complete OpenAPI conformance validator. It does not resolve
references, retrieve remote documents, validate every nested object, compare the
description with application routes or deployed behavior, lint API design, test
requests, infer exposure, or decide compatibility. Full conformance, contract
testing, and runtime drift require separately versioned adapters with explicit
filesystem, process, network, authentication, target, and destructive-request
policies.
