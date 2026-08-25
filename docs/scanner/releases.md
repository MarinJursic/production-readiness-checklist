# Scanner releases and verification

Scanner binaries use `scanner-vX.Y.Z` tags so their experimental product version
does not collide with versions of the checklist corpus. A release is built only
from the exact tagged commit. The workflow reruns the source, catalog, schema,
security, benchmark, and pack gates before publishing anything.

## Release contents

Each scanner release contains:

- Linux, macOS, and Windows archives for AMD64 and ARM64;
- one dependency-free `@marinjursic/vuk` npm launcher tarball and six exact
  native npm platform tarballs for the same systems;
- a binary with the scanner version, source revision, source timestamp, and Go
  toolchain embedded in `vuk version --format json`;
- the exact compatible `catalog/`, `packs/`, and `schemas/` trees in every
  archive, together with the catalog's human-readable objective sources and the
  packs' benchmark fixtures;
- a versioned release manifest binding artifact digests to the catalog and pack
  validation digests;
- a timestamp-free [CycloneDX 1.6](https://cyclonedx.org/specification/overview/)
  module SBOM;
- a canonical self-scan of the exact tagged source, executed by the packaged
  binary against its bundled catalog without hiding blocked or manual results;
- `SHA256SUMS`; and
- GitHub-hosted Sigstore attestations for SLSA build provenance and the
  CycloneDX SBOM predicate.

The SBOM describes the scanner's Go module and resolved Go dependencies. It does
not claim that the bundled catalog data, schemas, host operating system, or
future adapter images are Go components; those files are integrity-bound by the
archive checksum and release manifest instead. GitHub documents what its
[artifact attestations](https://docs.github.com/en/actions/how-tos/secure-your-work/use-artifact-attestations/use-artifact-attestations)
prove and how consumers should verify them.

The builder fixes archive order, ownership, modes, paths, timestamps, binary
metadata, and source identity. CI builds the complete release twice and requires
a byte-for-byte directory comparison. This is a reproducibility check under the
same declared CI toolchain; it is not a claim that provenance or reproducibility
proves the artifact is secure.

## Verify before running

Download the archive, `SHA256SUMS`, release manifest, and SBOM from the same
GitHub release. Verify the checksum from the directory containing those files:

```bash
sha256sum --check SHA256SUMS
```

On macOS, verify one downloaded archive against the matching line instead:

```bash
shasum -a 256 vuk_0.1.0_darwin_arm64.tar.gz
grep 'vuk_0.1.0_darwin_arm64.tar.gz' SHA256SUMS
```

Then verify that GitHub's signed provenance binds the archive to this repository
and release workflow:

```bash
gh attestation verify vuk_0.1.0_linux_amd64.tar.gz \
  --repo MarinJursic/production-readiness-checklist
```

Verify the separate CycloneDX SBOM attestation using its recognized predicate:

```bash
gh attestation verify vuk_0.1.0_linux_amd64.tar.gz \
  --repo MarinJursic/production-readiness-checklist \
  --predicate-type https://cyclonedx.org/bom
```

Finally, inspect `vuk_X.Y.Z_release-manifest.json` and compare its
`source_commit`, catalog digest, pack digests, and artifact digest with the
assessment scope you intend to use. Inspect `vuk_X.Y.Z_self-scan.json` as a
normal `prc.run/v0.12` report: a valid signed self-assessment may still be
`environment_blocked` because organizational, production, or adapter evidence
is deliberately unavailable in the release job. After extraction:

```bash
./vuk_X.Y.Z_linux_amd64/vuk version --format json
```

The release manifest also binds every npm tarball. Before the packages are
published to npm, they can be tested directly from one release directory. On
Linux x64, for example:

```bash
mkdir npm-smoke && cd npm-smoke
npm install --ignore-scripts --offline --no-audit --no-fund --package-lock=false \
  ../marinjursic-vuk-linux-x64-X.Y.Z.tgz \
  ../marinjursic-vuk-X.Y.Z.tgz
./node_modules/.bin/vuk version --format json
./node_modules/.bin/vuk scan /path/to/project
```

The platform package contains the native binary and its exact catalog. The
launcher checks the platform manifest and binary SHA-256 and never downloads a
fallback or starts a binary found on `PATH`. Public npm publishing uses npm
trusted publishing from the pinned release workflow, so the workflow refuses
`NPM_TOKEN` and `NODE_AUTH_TOKEN`. It verifies all seven tarballs against the
release manifest, publishes the six native packages before the launcher,
verifies the registry SHA-512 for every package, and safely skips only an
already-published version with exactly matching bytes. Because npm versions are
immutable, any byte mismatch stops the release. The publisher also runs npm in
a new empty working directory with user, global, and environment npm
configuration removed, so a saved `.npmrc` token cannot silently replace OIDC.

The release job builds once, uploads those exact bytes, then runs the matching
native archive and npm launcher on Linux x64, Linux ARM64, macOS x64, macOS
ARM64, Windows x64, and Windows ARM64. Publication starts only after every host
has completed a real 10,042-control smoke scan.

### One-time npm owner setup

npm requires a package to exist before a trusted publisher can be configured.
The owner must therefore bootstrap each of the seven package names once from
the exact verified release tarballs, with npm's required human authentication.
Then configure the same trusted-publisher identity on every package:

- repository owner: `MarinJursic`;
- repository: `production-readiness-checklist`;
- workflow filename: `release-scanner.yml`; and
- no GitHub environment unless the workflow is later changed to use one.

After that one-time setup, `scanner-vX.Y.Z` tags use OIDC and no long-lived npm
publishing secret. The release workflow requires Node.js 22.14 or newer and npm
11.5.1 or newer for trusted publishing. npm provenance links a package to its
build source; it does not prove the package has no unsafe code.

Do not substitute a successful signature check for vulnerability review or a
production-readiness decision. An attestation proves the signed claim's origin
and integrity, not the absence of defects.

## Release failure and revocation

Publishing is fail closed: an invalid tag, test failure, benchmark regression,
schema failure, vulnerability finding, non-reproducible output, checksum error,
or attestation failure prevents release publication. A failed workflow is not a
release and its temporary artifacts are not supported.

If a signing identity, workflow dependency, tool, release asset, catalog, pack,
or scanner version is compromised or materially incorrect, maintainers will:

1. publish a private-to-public security advisory as coordination permits;
2. mark the affected release and version as revoked, with the exact artifact
   digests and reason;
3. remove or revoke affected attestations where the platform supports it;
4. never reuse the affected version or tag and never replace an asset under the
   same name;
5. publish a new patched version from a reviewed commit with new checksums,
   SBOM, provenance, and compatibility manifest; and
6. update bundled trust stores or registry revocations when a pack, adapter, or
   publisher identity is affected.

Consumers should stop using a revoked artifact even if its historical signature
still verifies: verification establishes who produced it, not whether it remains
approved.
