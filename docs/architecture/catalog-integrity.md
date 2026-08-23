# Catalog integrity

The scanner treats the catalog as executable policy even though catalog files do
not run code. A malformed objective, assertion, or profile can otherwise change
what is assessed, what counts as evidence, or which result blocks a release.
Catalog loading therefore fails closed before inventory planning begins.

## Runtime loading contract

The Go scanner accepts only one bounded YAML document per regular catalog file.
Files larger than 4 MiB, symbolic links, unknown fields, trailing YAML documents,
unsupported schema versions, and mixed catalog versions are rejected. The
loader also checks the file identity and size while reading so a concurrent
replacement cannot silently become policy.

Every loaded definition is checked for:

- valid and unique objective, assertion, implementation, and profile IDs;
- positive revisions and supported catalog and profile versions;
- nonempty bounded UTF-8 text without NUL bytes;
- allowed automation classes, evidence authorities, severities, gates, and
  remediation classes;
- unique domains, evidence requirements, control links, profile members, and
  terminal severities;
- safe repository-relative Markdown source paths and exact source line text;
- existing, bidirectional objective-to-assertion mappings; and
- profile references to assertions that are present in the same catalog.

Any violation is a CLI input error. The scanner does not partially load a
catalog, ignore an invalid record, or continue with a weakened profile.

## Repository-wide generation checks

Runtime validation protects scanner consumers. The repository's Python catalog
check additionally rebuilds and compares the complete 10,042-control registry,
verifies source digests and revisions, validates every JSON Schema, and checks
the generated core-profile documentation. Run both the scanner tests and the
repository check after changing catalog policy:

```bash
go test ./scanner/catalog ./scanner/engine ./cmd/prc
python3 scripts/catalog.py check
python3 scripts/validate.py
```

The two layers are intentionally complementary. A released scanner must reject
an unsafe catalog even when repository CI was skipped, while repository CI must
detect drift across source controls and generated artifacts that a small runtime
profile does not load.

## Reproducible inspection and distribution

Validate a catalog and print its path-independent identity before planning a
scan:

```bash
prc catalog validate --catalog-root /path/to/catalog
prc catalog validate --catalog-root /path/to/catalog --format json
```

`prc catalog bundle` emits `prc.catalog-bundle/v0.1` JSON containing the same
manifest plus every validated objective, assertion, and profile in stable ID
order:

```bash
prc catalog bundle --catalog-root /path/to/catalog > catalog-bundle.json
```

The bundle contains no build timestamp or absolute path, so two copies of the
same definitions produce identical bytes. Its manifest records the semantic
catalog version, exact catalog digest, and definition counts. The checked-in
JSON Schemas validate both the manifest and bundle. Signing and publisher trust
remain release-pipeline responsibilities; an unsigned bundle is not evidence of
who published it.

## Identity binding

A successful load produces a deterministic digest over every parsed objective,
assertion, profile, and catalog version. The plan records that digest together
with the exact profile and assertion definitions. A catalog policy change can
therefore invalidate earlier evidence even when the target repository did not
change.

The filesystem root is excluded from the digest. Identical trusted definitions
have the same identity when installed at different absolute paths.
