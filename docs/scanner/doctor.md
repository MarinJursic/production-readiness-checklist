# Diagnose scanner capabilities

`prc doctor` checks whether the local host can perform the scanner operations
you intend to use. It does not execute project code, containers, Codex, or Claude
Code. Requested executable capabilities are resolved, identity-checked, and
hashed without being launched.

## Check read-only scanning

```bash
prc setup /path/to/project
# or, for the full diagnostic record:
prc doctor --target /path/to/project
```

Published builds discover their bundled catalog automatically. Pass
`--catalog-root` only when deliberately testing a source checkout or another
catalog build. Target inventory and catalog validation are required. State storage, candidate
workspaces, an OCI runtime, and agent providers are reported as optional warnings
unless their corresponding flags are supplied.

## Probe state and remediation filesystems

Create the state directory with private permissions before probing it:

```bash
mkdir -m 0700 /safe/path/prc-state

prc doctor \
  --target /path/to/project \
  --state-dir /safe/path/prc-state \
  --candidate-parent /safe/path/candidates
```

The state probe verifies private-directory enforcement and the hard-link
publication primitive used for immutable evidence records. The candidate probe
creates a private temporary sibling directory, verifies that it is disjoint from
the target, and tests file creation and mode enforcement. Probe directories are
removed before the command returns.

Do not place `--candidate-parent` inside the target. An ancestor such as a shared
temporary root is allowed only when the actual private probe directory is a
disjoint sibling of the target.

## Inspect optional executables

```bash
prc doctor \
  --target /path/to/project \
  --oci-runtime docker \
  --provider codex \
  --provider claude
```

The command you request must be named `docker` or `podman` for OCI, or `codex`
or `claude` for a provider. Normal package-manager command links are supported:
Doctor checks the requested command name, follows the link, and hashes the real
versioned executable. The report records both paths when they differ. A matching
filename is not treated as trust; execution planning can bind the exact digest
later. Doctor does not test authentication, make network requests, pull images,
or prove that a daemon is running.

After the exact Gitleaks image in the
[adapter guide](../architecture/adapters.md#inspect-and-validate) is present,
`prc verify /path/to/project` runs the core scan plus that one bundled adapter.
It never downloads the image and fails closed if the digest is unavailable.

## Automation output

```bash
prc doctor --target . --format json > doctor.json
python scripts/validate_instance.py doctor.schema.json doctor.json
```

JSON conforms to `prc.doctor/v0.1`. Each check is `pass`, `warn`, or `fail` and
states whether it is required. The command exits `0` when all requested required
capabilities pass and `2` when at least one required capability fails. Warnings
never silently claim that an untested optional capability is available.
