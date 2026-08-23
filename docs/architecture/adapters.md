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
available. Run result v0.3 can embed these records. A scan may execute one
adapter only when an applicable assertion binds the exact adapter ID, manifest
SHA-256 digest, and observation kind:

```bash
prc scan \
  --target /path/to/project \
  --catalog-root /path/to/trusted/catalog \
  --adapter-manifest /path/to/pinned-adapter.yaml \
  --adapter-runtime docker
```

Authorization is checked before the OCI runtime is invoked. The adapter cannot
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
