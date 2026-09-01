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
checked-in release suite includes completed, unsupported, timeout, malformed,
resource-limit, undeclared-output, and explicit evaluator-authority attack
transcripts.

Manifest v0.4 also defines two closed native-output protocols. The first,
`prc-adapter-gitleaks-json-v1`, accepts only Gitleaks 8.30.0's reviewed
official image digest, exact scanner-owned current-tree command, and JSON
report contract. The scanner supplies the SHA-256-pinned upstream default
ruleset on standard input, forces full redaction, ignores target-owned Gitleaks
configuration, ignore files, and `gitleaks:allow` comments, and disables archive
and recursive-decoding expansion. The normalizer rejects unknown or duplicate
fields, unredacted findings, history or symlink metadata, paths outside the
snapshot, inconsistent fingerprints, invalid coordinates, and excessive
findings. It retains only rule identity and normalized location metadata.

The exact redacted native report is represented by a media type, byte count,
and SHA-256 artifact descriptor in the transcript; raw report content is not
copied into the durable run record. This preserves content provenance without
persisting matched source context. An empty report becomes an explicit
`not_found` observation. One or more findings become `found` observations, but
neither the tool nor normalizer can declare an assertion pass or failure.

The second closed protocol, `prc-adapter-syft-cyclonedx-json-v1`, accepts only
Syft 1.51.0's reviewed non-root image digest and the scanner-owned directory
command. The scanner injects an empty configuration at a reserved `.prc` path,
so target `.syft.yaml` or environment defaults cannot disable catalogers or
redirect output. The normalizer requires Syft's exact tool identity and
CycloneDX 1.7 document envelope, rejects duplicate keys, malformed component,
package-URL, dependency, or metadata structures, and caps the component and
output counts.

Syft generates a fresh timestamp and serial number for identical input. The
normalizer removes those optional values and the source component's
non-semantic `bom-ref`, then canonically orders components, properties, and
dependencies. The resulting `application/vnd.cyclonedx+json;version=1.7`
artifact is content-addressed and byte-for-byte deterministic for the same
sealed inventory. When `--state-dir` is supplied, the bytes are written before
the run that references them. A successful `value` observation proves only
that this repository-inventory SBOM was generated; it does not claim that a
built artifact is complete, vulnerability-free, acceptably licensed, or
production ready.

The third closed protocol, `prc-adapter-grype-json-v1`, accepts only Grype
0.116.1's reviewed image digest, exact scanner-owned command, and enriched JSON
report. It injects scanner-owned configuration that disables application and
database updates, external sources, target ignores, and VEX policy. The
container runs without network access against one explicitly supplied,
read-only `grype-db` directory. The runner hashes and seals that directory while
the normalizer independently requires a valid schema-v6 database, an official
archive URL with a SHA-256 checksum, a build age no greater than 120 hours, and
valid EPSS, KEV, and NVD provider provenance.

Grype output is decoded with duplicate and unknown-field rejection and bounded
nested records. Ignored matches, package alerts, escaping paths, unsupported
tool or database identities, incomplete matching evidence, or unsafe numeric
values fail the execution. The deterministic normalized artifact retains only
the findings and provenance needed for assessment; report timestamps and bulk
raw metadata do not affect its digest. A `not_found` observation is time-bound
evidence for the catalog gate, while any `found` observation remains a factual
tool result for the engine to assess.

An assertion binding owns the meaning of an observation outcome. Bindings that
omit an outcome policy retain the conservative analysis default:
`not_found` passes and `found` fails. A binding may instead declare a nonempty
`pass_outcomes` list and a disjoint, possibly empty `fail_outcomes` list from
the closed `found`, `not_found`, `value`, and `unsupported` vocabulary. Any
unconfigured outcome remains Unknown; a configured pass mixed with another
unconfigured outcome is Conflicting. The SBOM assertion therefore treats only
`value` as passing and has no failing tool outcome: malformed, missing,
unsupported, partial, or contradictory output never becomes Pass.

## Capability manifest

Every external adapter has a strict
[`adapter-manifest.schema.json`](https://github.com/MarinJursic/production-readiness-checklist/blob/main/schemas/adapter-manifest.schema.json).
Manifest v0.4 binds the exact protocol and output schema, compatible
engine APIs, publisher and owner identities, immutable tool version and
supported format versions, declared observation kinds, maintenance state, and
known limitations. Output validation rejects an observation kind that is not
declared by that exact manifest. The archived v0.1 through v0.3 schemas remain
available for record interpretation but are not accepted for new execution.

Trust is deliberately not self-declared by the adapter. A registry or explicit
local operator grant must bind trust to the exact manifest digest and publisher;
the manifest alone cannot promote itself to first-party or verified status.
The current experimental runner deliberately supports only a narrow subset:

- OCI execution through Docker or Podman;
- an image reference pinned by a `sha256` digest and explicit registry host;
- a private, read-only snapshot containing the regular files in the sealed
  inventory plus any exact protocol-owned policy input described below;
- explicitly declared, content-hashed, read-only external data directories for
  reviewed native adapters whose immutable datasets are too large to copy into
  each project snapshot;
- no image pull during a scan;
- no network;
- no secret handles;
- no child processes for generic JSONL adapters, with narrowly reviewed,
  PID-bounded OS-task allowances for the pinned Gitleaks, Syft, and Grype
  binaries;
- an optional bounded scratch `tmpfs`; and
- explicit wall-time, memory, CPU, process, line, message, stdin, stdout, and
  stderr limits.

Before execution, the scanner reopens and hashes every inventoried regular file
while copying it to a private temporary snapshot. Scanner-excluded paths and
symlinks are not copied. Protocol-owned path remapping may prevent a target file
from changing analyzer policy without dropping its content: the Gitleaks
protocol relocates the root `.gitleaksignore`, scans its original bytes, and
maps any resulting location back to `.gitleaksignore`. The Syft and Grype
protocols add only their digest-bound scanner configurations under the
inventory-excluded `.prc` namespace. A changed type, size, or
digest stops execution. The snapshot has a 4 GiB safety ceiling, is mounted
read-only, and is removed after the run. Its own deterministic digest is sealed
into the OCI plan and checked again immediately before and after the container
runs, so observations cannot silently refer to different bytes than the
inventory.

External data is a separate trust boundary. A manifest can declare only a
bounded named directory mounted at its reserved `/prc-inputs/NAME` path. The
operator supplies the host directory explicitly. The scanner rejects empty,
oversized, over-populated, symlink-containing, non-regular, missing, extra, or
duplicate mounts; hashes every file before planning; seals the digest, file
count, byte count, and destination into the OCI plan; and rehashes the directory
immediately before and after execution. The durable execution records this
content identity but never the host path. A data directory changing during a
run is an execution error, not evidence. This mechanism does not download or
refresh data and does not let an adapter expand its declared mounts.

The generated OCI command also drops all Linux capabilities, enables
`no-new-privileges`, uses the invoking non-root UID/GID so private snapshot
permissions need not be broadened, disables swap beyond the memory ceiling,
caps open files, makes the image root filesystem read-only, and removes the
named container. The runner refuses root and Windows hosts. Timeout cleanup
targets only the scanner-generated container name. The plan records the OCI
client binary digest and execution stops if the plan, manifest, runtime binary,
or bound snapshot changes after capability evaluation.

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
local entries. A manifest cannot alter these registry decisions. Detached,
scoped Ed25519 verification is available through an explicitly selected
[publisher trust store](publisher-trust.md). The current repository does not
yet publish an official release trust store or signed registry, so a public
adapter distribution channel is not yet complete.

Validate and inspect a lockfile without executing anything:

```bash
prc adapter registry-validate \
  --file /path/to/adapter-registry.yaml \
  --format json
```

## Inspect and validate

Validate the full release fixture suite without executing an adapter:

```bash
prc adapter fixture-validate \
  --suite fixtures/adapters/fixture-suite.yaml \
  --format json
```

The strict
[`adapter-fixture-suite.schema.json`](https://github.com/MarinJursic/production-readiness-checklist/blob/main/schemas/adapter-fixture-suite.schema.json)
binds every transcript to the canonical digest of one exact manifest. A case
may reduce the manifest's line, message, or stdout ceiling to exercise a
resource failure, but it cannot increase adapter authority or resource limits.
The runner rejects symlinks and paths outside the suite, hashes the suite and
transcript corpus, and evaluates every case twice. It compares exact
disposition, summary status, error class, and observation ID, kind, and
outcome. CI fails when any expectation drifts or either evaluation differs.

The resulting
[`adapter-fixture-report.schema.json`](https://github.com/MarinJursic/production-readiness-checklist/blob/main/schemas/adapter-fixture-report.schema.json)
is release evidence for protocol compatibility only. It does not prove the
upstream analyzer's detection accuracy, image provenance, runtime isolation,
or suitability for the production profile. A real adapter still needs its own
tool-specific corpus, reviewed pinned image, registry entry, publisher trust,
and exact catalog binding.

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
available. For a manifest that declares external data, repeat
`--data NAME=/path/to/directory` for both `plan-oci` and `run-oci`. Current v0.3
execution records include required `resolution` and `data_inputs` identities.
An explicit manifest records the publisher and `local-explicit`
operator grant; a registry resolution additionally binds the registry ID,
revision, content digest, and registry-assigned trust. Changing any provenance
field changes the execution ID. Version-specific v0.1 and v0.2 records remain
valid for archived runs, but v0.1 records cannot be supplied as evidence to a
new scan. Current run results embed these records. A scan may execute up to 16
adapters when applicable assertions bind every exact adapter ID, manifest
SHA-256 digest, and observation kind. Repeat `--adapter-manifest` for an
explicit local set:

```bash
prc scan \
  --target /path/to/project \
  --catalog-root /path/to/trusted/catalog \
  --mode verify-local \
  --adapter-manifest /path/to/first-pinned-adapter.yaml \
  --adapter-manifest /path/to/second-pinned-adapter.yaml \
  --adapter-runtime docker
```

Scan-level data bindings use an adapter-qualified name so multi-adapter runs
cannot consume one another's inputs:

```bash
prc scan \
  --target /path/to/project \
  --catalog-root /path/to/trusted/catalog \
  --profile prc/supply-chain \
  --mode verify-local \
  --adapter-manifest /path/to/pinned-adapter.yaml \
  --adapter-data 'prc.adapter.example@1.0/example-db=/path/to/database-cache'
```

For registry-assigned trust and revocation, resolve the same catalog-pinned
adapter through a lockfile instead:

```bash
prc scan \
  --target /path/to/project \
  --catalog-root /path/to/trusted/catalog \
  --mode verify-local \
  --adapter-registry /path/to/adapter-registry.yaml \
  --adapter-id prc.adapter.first@1.0 \
  --adapter-id prc.adapter.second@1.0 \
  --adapter-runtime docker
```

`--adapter-manifest` is the explicit local-operator path and is mutually
exclusive with `--adapter-registry`. Adapter paths and IDs cannot repeat. Both
paths still require an exact manifest digest binding in an applicable catalog
assertion; registry approval cannot authorize a catalog-unbound adapter. The
scanner resolves, validates, and authorizes the complete set before invoking the
runtime, then executes the set in deterministic adapter-ID and manifest-digest
order. When project configuration supplies a maximum duration, one shared
deadline bounds the complete set rather than restarting the budget per adapter.

The explicit mode grants only the reviewed no-network OCI capability envelope;
authorization is checked before the OCI runtime is invoked. The adapter cannot
declare an assertion assessment: the engine maps `found`, `not_found`,
`unsupported`, incomplete, and conflicting observations through the catalog's
immutable binding. The execution and each resulting evidence envelope are bound
to the exact inventory digest. Offline execution-record import is deliberately
unsupported because a content digest alone does not prove that a tool ran.

The default `PRC-A-CORE-013` binding pins
`prc.adapter.gitleaks@8.30` and the canonical manifest digest for
`adapters/gitleaks-v8.30.0.yaml`. It remains Blocked in an ordinary inspect-mode
core scan because inspect mode never launches containers. To produce executed
evidence, an operator must first pull the exact image digest and then explicitly
run `prc verify`, which selects `--mode verify-local` with the bundled manifest.
The alias rejects custom adapter, profile, and AI overrides. The checked-in
protocol fixture still exercises generic JSONL plumbing in tests but is not
authorized by the production profile.

```bash
docker pull \
  ghcr.io/gitleaks/gitleaks@sha256:691af3c7c5a48b16f187ce3446d5f194838f91238f27270ed36eef6359a574d9

prc verify /path/to/project
```

The command scans only the sealed current-tree snapshot. It does not scan Git
history, files above 10 MiB, symlink targets, archives, or recursively decoded
content. Gitleaks rules are heuristic; both false positives and false negatives
remain possible, so this binding is one selected analysis class rather than a
claim of complete secret or static-analysis coverage.

Normal validation and scanner-release workflows pull that exact digest and run
the clean and suppression-resistant finding cases through the production OCI
runner. A nonempty test-image override that differs from the reviewed digest is
an error, not a skipped test.

## Design references

- [OCI content descriptors](https://specs.opencontainers.org/image-spec/descriptor/?v=v1.1.0)
  define digest-addressed references and require retrieved content to be checked
  against its digest.
- [Docker container run](https://docs.docker.com/reference/cli/docker/container/run/)
  documents read-only roots, capability removal, resource limits, and
  `no-new-privileges`.
- [Podman run](https://docs.podman.io/en/stable/markdown/podman-run.1.html)
  documents the corresponding rootless container and resource controls.
- [Gitleaks v8.30.0](https://github.com/gitleaks/gitleaks/releases/tag/v8.30.0)
  is the pinned upstream release; its official CLI documents directory scans,
  JSON reports, full redaction, timeouts, and bounded target-file size.
- [Gitleaks default configuration](https://github.com/gitleaks/gitleaks/blob/v8.30.0/config/gitleaks.toml)
  is vendored as a deterministic gzip archive, attributed in
  `THIRD_PARTY_NOTICES.md`, and verified in both compressed and exact
  decompressed forms before every execution.
