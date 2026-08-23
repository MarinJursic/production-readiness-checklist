# Software supply-chain profile

The focused `prc/supply-chain@0.1` profile combines repository-native checks
for dependency and build-input hygiene with `PRC-A-SUPPLY-001`, which generates
a normalized CycloneDX 1.7 software bill of materials from the sealed current
repository inventory. SBOM generation runs only in the explicit
`verify-local` capability mode through the digest-pinned Syft 1.51.0 adapter.

Pull the reviewed image by its immutable digest before scanning; the scanner
never pulls an image during a run:

```bash
docker pull \
  ghcr.io/anchore/syft@sha256:d2dc3ec86cb2b4e7ddb226ba0305c4523b7c0694c45d9f576b42b4c2f5ce7aa8

mkdir -m 0700 /safe/path/prc-state

prc scan \
  --target PATH \
  --catalog-root PATH_TO_RELEASE \
  --profile prc/supply-chain \
  --mode verify-local \
  --adapter-manifest PATH_TO_RELEASE/adapters/syft-v1.51.0.yaml \
  --state-dir /safe/path/prc-state \
  --format json
```

The run's adapter transcript contains a `sha256:` artifact descriptor. Its
payload is stored at:

```text
/safe/path/prc-state/artifacts/sha256/<first-two-hex>/<sha256-hex>
```

The payload uses the registered
`application/vnd.cyclonedx+json;version=1.7` media type. CycloneDX identifies
JSON Schema as its reference format and documents `bom.json` and `*.cdx.json`
as conventional names. See the
[CycloneDX specification](https://github.com/CycloneDX/specification) and the
[Syft v1.51.0 release](https://github.com/anchore/syft/releases/tag/v1.51.0).

## What the generated-SBOM assertion proves

The assertion passes only when the exact catalog-pinned manifest produces a
completed `sbom-generation` observation with the configured `value` outcome.
The engine, not Syft, maps that outcome to Pass. The execution remains bound to
the inventory digest, manifest digest, image digest, scanner-owned command,
snapshot digest, normalized artifact digest, and explicit local or registry
authorization provenance.

For identical sealed input, the normalizer removes Syft's optional generation
timestamp, document serial number, and source component `bom-ref`, then sorts
components, properties, dependencies, and dependency targets before encoding.
No package, file, license, hash, package URL, or dependency content is removed.

## Profile coverage

In addition to the generated SBOM, the profile checks:

- repository license presence;
- dependency lock or checksum coverage and nonempty dependency inputs;
- configured dependency-update automation;
- immutable GitHub Actions references;
- declared runtime versions;
- immutable container base references; and
- locked Terraform provider selections.

## Deliberate limitations

This is a source-repository SBOM, not an SBOM for a compiled binary, container
image, installer, deployment bundle, or exact production release. Discovery is
limited by Syft's cataloger coverage and the scanner's regular-file inventory.
Generation does not prove component reachability, dependency support,
vulnerability status, exploitability, license compatibility, provenance,
signature validity, or deployed equivalence. Those are separate assertions and
must never be inferred from `PRC-A-SUPPLY-001` passing.
