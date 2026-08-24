# Software supply-chain profile

The focused `prc/supply-chain@0.2` profile combines repository-native checks
with two executed analysis gates:

- `PRC-A-SUPPLY-001` generates a normalized CycloneDX 1.7 software bill of
  materials from the sealed repository inventory with Syft 1.51.0; and
- `PRC-A-SUPPLY-002` checks packages discoverable in that same inventory for
  known vulnerabilities with Grype 0.116.1 and a fresh, identity-bound offline
  database.

Both run only in the explicit `verify-local` capability mode. The scanner never
pulls images or downloads vulnerability data during a scan.

## Prepare the immutable tools and database

Pull the two reviewed images by immutable digest:

```bash
docker pull \
  ghcr.io/anchore/syft@sha256:d2dc3ec86cb2b4e7ddb226ba0305c4523b7c0694c45d9f576b42b4c2f5ce7aa8
docker pull \
  ghcr.io/anchore/grype@sha256:1e71065c0a4cff3e6bd3b8add525ffac4343eb4971694eb90a31cf6d4d3e85db
```

Database acquisition is a separate, explicitly networked operator action. The
following command asks the exact reviewed Grype image to download and validate
the current database into an operator-controlled directory:

```bash
mkdir -m 0700 /safe/path/grype-db

docker run --rm --pull=never \
  --user "$(id -u):$(id -g)" \
  --mount type=bind,src=/safe/path/grype-db,dst=/grype-db \
  --env GRYPE_CHECK_FOR_APP_UPDATE=false \
  --env GRYPE_DB_CACHE_DIR=/grype-db \
  ghcr.io/anchore/grype@sha256:1e71065c0a4cff3e6bd3b8add525ffac4343eb4971694eb90a31cf6d4d3e85db \
  db update
```

Review this acquisition in environments with egress or software-source policy.
The scan itself mounts the resulting directory read-only, runs with
`--network=none`, verifies the database's official archive checksum, and rejects
a database built more than 120 hours earlier. Refresh it before that boundary.

## Run both gates

```bash
mkdir -m 0700 /safe/path/prc-state

everylast scan \
  --target PATH \
  --catalog-root PATH_TO_RELEASE \
  --profile prc/supply-chain \
  --mode verify-local \
  --adapter-manifest PATH_TO_RELEASE/adapters/syft-v1.51.0.yaml \
  --adapter-manifest PATH_TO_RELEASE/adapters/grype-v0.116.1.yaml \
  --adapter-data 'prc.adapter.grype@0.116/grype-db=/safe/path/grype-db' \
  --state-dir /safe/path/prc-state \
  --format json
```

The complete adapter set is authorized before either container starts. It runs
in deterministic adapter-ID and manifest-digest order under one configured
deadline. The database directory is hashed before planning and before and after
execution; its digest, file count, byte count, and reserved container
destination are recorded without persisting the host path.

Each transcript contains a `sha256:` artifact descriptor. With `--state-dir`,
the corresponding immutable payload is stored at:

```text
/safe/path/prc-state/artifacts/sha256/<first-two-hex>/<sha256-hex>
```

The SBOM uses
`application/vnd.cyclonedx+json;version=1.7`. The vulnerability report uses
`application/vnd.prc.grype.vulnerability-report+json;version=1` and the public
`schemas/grype-vulnerability-report.schema.json` contract.

## What the gates prove

The generated-SBOM assertion passes only when the exact catalog-pinned Syft
manifest produces a completed `sbom-generation` observation with the configured
`value` outcome. For identical sealed input, the normalizer removes Syft's
generation timestamp, document serial number, and source component `bom-ref`,
then canonically orders the retained component, file, license, hash, package
URL, and dependency content.

The known-vulnerability assertion passes only when the exact catalog-pinned
Grype manifest reports `not_found`. Any `found` observation fails the critical
gate. The normalized result retains vulnerability and alias identity, affected
package and locations, severity, best available CVSS and EPSS values, CISA KEV
and known-ransomware flags, Grype risk, fix state and versions, and database
build and provider provenance. Target `.grype.yaml`, VEX, and ignore policy
cannot suppress the scanner-owned analysis; any ignored match or unsupported
package alert makes the adapter fail closed.

The engine, not either tool, maps these factual outcomes to Pass or Fail. Both
executions remain bound to the inventory, snapshot, manifest, image, command,
artifact, and external-data digests.

## Profile coverage

In addition to the two executed analysis gates, the profile checks:

- repository license presence;
- dependency lock or checksum coverage and nonempty dependency inputs;
- configured dependency-update automation;
- immutable GitHub Actions references;
- declared runtime versions;
- immutable container base references; and
- locked Terraform provider selections.

## Deliberate limitations

This is source-repository analysis, not analysis of a compiled binary, container
image, installer, deployment bundle, or exact production release. Discovery is
limited by Syft and Grype cataloger coverage and the scanner's regular-file
inventory. A clean result is time-bound to the recorded database and does not
prove that no vulnerability exists, that every dependency is reachable, that
the built artifact matches the source tree, or that severity and exploitability
metadata are complete. Vulnerability feeds, package matching, aliases, fixes,
CVSS, EPSS, KEV, and risk can contain delays, false positives, false negatives,
or incomplete data. License compatibility, provenance, signatures, runtime
configuration, and deployed equivalence remain separate assertions.

See the [CycloneDX specification](https://github.com/CycloneDX/specification),
[Syft v1.51.0 release](https://github.com/anchore/syft/releases/tag/v1.51.0),
[Grype v0.116.1 release](https://github.com/anchore/grype/releases/tag/v0.116.1),
and [Grype database documentation](https://oss.anchore.com/docs/guides/vulnerability/database/).
