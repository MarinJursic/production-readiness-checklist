# Bounded private-key armor inspection

`PRC-A-CORE-031` is a deterministic, no-network guard against committing
recognized armored private-key blocks. It inspects content-addressed files from
the current inventory and never copies a matched marker or following key bytes
into evidence, findings, logs, or reports. A failure identifies only the file
path, inventory and content digests, and a redacted summary.

The native rules cover private and encrypted private-key labels standardized by
[RFC 7468](https://www.rfc-editor.org/info/rfc7468/), legacy RSA, EC, and DSA
private-key labels, the
[OpenSSH private-key format](https://github.com/openssh/openssh-portable/blob/master/PROTOCOL.key),
and armored PGP private-key blocks. Public-key and certificate armor is not a
failure.

## Fail-closed bounds

The check scans recognized source, configuration, documentation,
infrastructure, environment, and key-file paths. Repository dotfiles are also
included, except platform metadata that is already excluded by inventory. It
stops as Blocked/Unknown instead of returning Pass when any of these bounds
would prevent complete inspection:

- 10,000 candidate files;
- 8 MiB per candidate file; or
- 256 MiB across all candidate files.

Every file is reopened without following a symlink and checked against its
inventoried size, mode, and SHA-256 before its bytes are inspected. A target
mutation is an execution error, never a passing result. Reports include at most
100 violating locations while retaining the exact total count.

## Deliberate limitations

This check does not inspect Git history, ignored dependency/cache directories,
binary or unrecognized file types, remote artifacts, container layers, secret
stores, deployed environments, or arbitrary vendor API-token formats. It does
not prove that a recognized key is live, but a detected private-key block is
still a critical no-go because verification must not expose or exercise it.

Broader secret detection belongs to a maintained, digest-pinned analysis
adapter with a tool-specific benchmark corpus. Until such an adapter is
published and authorized, `PRC-A-CORE-013` remains independently Blocked; this
narrow native check must not be presented as complete secret-scanning coverage.
