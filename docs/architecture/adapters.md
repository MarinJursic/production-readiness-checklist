# Sandboxed adapter protocol

External tools are untrusted evidence producers. They do not evaluate controls,
change policy, or mark assertions as passing. The scanner accepts only
schema-valid observations and leaves the final assertion assessment to the
deterministic engine.

## Protocol

`prc-adapter-jsonl-v1` is a line-delimited JSON protocol over standard input and
standard output. The scanner sends exactly three messages:

```json
{"type":"hello","protocol":"prc-adapter-jsonl-v1","run_id":"<64 lowercase hex characters>"}
{"type":"input","subject":{"target_name":"example","inventory_digest":"<64 lowercase hex characters>"},"facts":{},"config":{}}
{"type":"execute"}
```

An adapter may emit bounded `log`, `observation`, and `artifact` messages and
must end with exactly one `summary`. Standard output is protocol-only. Standard
error is bounded diagnostic text. A summary is an implementation execution
status, not a control result.

Raw standard error is not echoed in CLI JSON. The runner reports only its byte
count and SHA-256 digest. Protocol messages may still describe sensitive target
facts, so operators must treat the transcript as assessment evidence and apply
appropriate access, redaction, and retention controls.

The decoder rejects:

- unknown or duplicate JSON fields;
- blank, oversized, excessive, or trailing messages;
- absolute, non-normalized, or escaping artifact and source paths;
- malformed digests, statuses, locations, and counters;
- missing summaries or messages after a summary; and
- fields through which an adapter tries to declare `pass`, suppress a finding,
  or otherwise take evaluator authority.

The machine contract is
[`adapter-message.schema.json`](https://github.com/MarinJursic/production-readiness-checklist/blob/main/schemas/adapter-message.schema.json). The
checked-in fixtures include an accepted transcript and an explicit attempt to
inject a passing assessment.

## Capability manifest

Every external adapter has a strict
[`adapter-manifest.schema.json`](https://github.com/MarinJursic/production-readiness-checklist/blob/main/schemas/adapter-manifest.schema.json).
Manifest v0.2 binds the exact JSONL protocol and output schema, compatible
engine APIs, publisher and owner identities, immutable tool version and
supported format versions, declared observation kinds, maintenance state, and
known limitations. Output validation rejects an observation kind that is not
declared by that exact manifest. The archived v0.1 schema remains available for
record interpretation but is not accepted for new execution.

Trust is deliberately not self-declared by the adapter. A registry or explicit
local operator grant must bind trust to the exact manifest digest and publisher;
the manifest alone cannot promote itself to first-party or verified status.
The current experimental runner deliberately supports only a narrow subset:

- OCI execution through Docker or Podman;
- an image reference pinned by a `sha256` digest and explicit registry host;
- a read-only target workspace;
- no image pull during a scan;
- no network;
- no secret handles;
- no declared child processes;
- an optional bounded scratch `tmpfs`; and
- explicit wall-time, memory, CPU, process, line, message, stdin, stdout, and
  stderr limits.

The generated OCI command also drops all Linux capabilities, enables
`no-new-privileges`, runs as UID/GID 65532, makes the image root filesystem
read-only, and removes the named container. Timeout cleanup targets only the
scanner-generated container name. The plan records the OCI client binary digest
and execution stops if either the plan, manifest, or runtime binary changes
after capability evaluation.

An OCI runtime is a security control, not proof of perfect isolation. Operators
must keep the runtime and host patched, review adapter images, and treat runtime
or resource-control failure as an execution error. Podman documents that some
rootless resource limits depend on the host cgroup configuration.

## Registry lockfile and revocation

[`adapter-registry.schema.json`](https://github.com/MarinJursic/production-readiness-checklist/blob/main/schemas/adapter-registry.schema.json)
defines the local registry trust root. Each entry pins the adapter ID, manifest
SHA-256, publisher ID, registry-assigned trust, lifecycle status, and normalized
relative manifest path. The loader hashes and validates the manifest, rejects
publisher or lifecycle drift, rejects symlinked or escaping paths, and verifies
current engine compatibility.

`revoked` entries remain effective even when the compromised manifest has been
removed. Default resolution permits only `first-party-sandboxed` and
`verified-community` entries; it denies deprecated, unverified-community, and
local entries. A manifest cannot alter these registry decisions. The current
lockfile is a local trust root, not yet a signed distribution system; signed
registry releases and publisher-key verification remain required before a
public adapter ecosystem can be considered complete.

Validate and inspect a lockfile without executing anything:

```bash
prc adapter registry-validate \
  --file /path/to/adapter-registry.yaml \
  --format json
```

## Inspect and validate

Validate a transcript without executing an adapter:

```bash
prc adapter validate-output \
  --manifest fixtures/adapters/fixture-adapter.yaml \
  --file fixtures/adapters/valid-output.jsonl
```

Render the exact OCI command without starting a container:

```bash
prc adapter plan-oci \
  --manifest /path/to/pinned-adapter.yaml \
  --target /path/to/project \
  --runtime docker
```

Run an already-present pinned image in the bounded OCI environment and emit a
content-addressed execution record bound to the exact manifest and inventory:

```bash
prc adapter run-oci \
  --manifest /path/to/pinned-adapter.yaml \
  --target /path/to/project \
  --runtime docker
```

`--pull=never` means the command fails if the exact image is not already
available. Current run results embed these records. A scan may execute one
adapter only when an applicable assertion binds the exact adapter ID, manifest
SHA-256 digest, and observation kind:

```bash
prc scan \
  --target /path/to/project \
  --catalog-root /path/to/trusted/catalog \
  --mode verify-local \
  --adapter-manifest /path/to/pinned-adapter.yaml \
  --adapter-runtime docker
```

For registry-assigned trust and revocation, resolve the same catalog-pinned
adapter through a lockfile instead:

```bash
prc scan \
  --target /path/to/project \
  --catalog-root /path/to/trusted/catalog \
  --mode verify-local \
  --adapter-registry /path/to/adapter-registry.yaml \
  --adapter-id prc.adapter.example@1.0 \
  --adapter-runtime docker
```

`--adapter-manifest` is the explicit local-operator path and is mutually
exclusive with `--adapter-registry`. Both paths still require an exact manifest
digest binding in an applicable catalog assertion; registry approval cannot
authorize a catalog-unbound adapter.

The explicit mode grants only the reviewed no-network OCI capability envelope;
authorization is checked before the OCI runtime is invoked. The adapter cannot
declare an assertion assessment: the engine maps `found`, `not_found`,
`unsupported`, incomplete, and conflicting observations through the catalog's
immutable binding. The execution and each resulting evidence envelope are bound
to the exact inventory digest. Offline execution-record import is deliberately
unsupported because a content digest alone does not prove that a tool ran.

The default `PRC-A-CORE-013` binding list remains empty until a production
analysis adapter image is published, reviewed, and digest-pinned. It therefore
remains Blocked in an ordinary core scan. The checked-in fixture exercises the
plumbing in tests but is not authorized by the production profile.

## Design references

- [OCI content descriptors](https://specs.opencontainers.org/image-spec/descriptor/?v=v1.1.0)
  define digest-addressed references and require retrieved content to be checked
  against its digest.
- [Docker container run](https://docs.docker.com/reference/cli/docker/container/run/)
  documents read-only roots, capability removal, resource limits, and
  `no-new-privileges`.
- [Podman run](https://docs.podman.io/en/stable/markdown/podman-run.1.html)
  documents the corresponding rootless container and resource controls.
