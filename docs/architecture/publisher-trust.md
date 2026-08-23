# Offline publisher trust

Content digests prove identity, not authorship. The scanner therefore supports
detached Ed25519 verification for pack and adapter-registry artifacts through
an operator-selected trust store. It never discovers keys, downloads trust
metadata, or handles publisher private keys.

The trust store is an explicit local root of trust. Each key has an ID, exact
public key, artifact-kind scopes, validity interval, and active or revoked
status. Revocation is fail-closed and applies even to a signature created
before the key was revoked. A key authorized for a pack cannot sign an adapter
registry unless that second scope is also present.

The detached signature covers a domain-separated canonical payload containing:

- signature schema and algorithm;
- artifact kind and versioned artifact ID;
- the artifact's canonical SHA-256 identity;
- publisher key ID; and
- UTC issuance time.

Verification requires an explicit UTC time so historical automation remains
reproducible. The issuance and verification times must both fall inside the
key validity interval, and a future-issued signature is rejected. The output
record binds the signature digest and canonical trust-store digest.

## Verify a pack

Pack verification first validates the catalog binding and reruns the pinned
benchmark. Only then does it verify the detached signature against the
resulting canonical pack digest:

```bash
prc pack verify \
  --catalog-root . \
  --file /path/to/pack.yaml \
  --trust-store /path/to/trust-store.yaml \
  --signature /path/to/pack.signature.yaml \
  --verified-at 2026-08-23T13:00:00Z \
  --format json
```

## Verify an adapter registry

Registry verification first resolves every non-revoked manifest and checks its
ID, publisher, lifecycle, engine compatibility, and digest pins:

```bash
prc adapter registry-verify \
  --file /path/to/adapter-registry.yaml \
  --trust-store /path/to/trust-store.yaml \
  --signature /path/to/registry.signature.yaml \
  --verified-at 2026-08-23T13:00:00Z \
  --format json
```

The schemas are `trust-store.schema.json`, `signature.schema.json`, and
`signature-verification.schema.json`. The current repository does not publish
an official release trust store or signatures yet. Verification support is not
itself a key ceremony, secure private-key service, transparency log,
reproducible release, or revocation-distribution channel; those remain release
engineering requirements.
