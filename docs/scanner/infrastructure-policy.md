# Infrastructure policy profile

`prc/iac@0.1` combines bounded native checks with a reviewed Checkov 3.3.8
analysis of inventoried Terraform, Kubernetes, and container definitions. It is
applicable only when the inventory contains at least one of those inputs.

## Run it

Review `adapters/checkov-v3.3.8.yaml` in the release or repository root, then
pre-pull the exact reviewed image:

```bash
docker pull docker.io/bridgecrew/checkov@sha256:c64ffb6d6fc8087c896341a2c697770a04a1cf558db04fa7b8129d8ca6bce336
```

Run the profile with an explicit local-verification capability grant:

```bash
./prc scan /path/to/project \
  --profile prc/iac \
  --mode verify-local \
  --adapter-manifest adapters/checkov-v3.3.8.yaml
```

The ordinary `prc scan` command does not launch Checkov or any other external
analyzer. The `verify-local` mode and exact manifest are both required.

## What is evaluated

The profile includes native checks for immutable container base identities,
non-root container users, Terraform dependency locks, Kubernetes non-root
workloads, privilege restrictions, Linux capabilities, seccomp, and container
resource requests and limits. `PRC-A-IAC-001` additionally consumes only the
`iac-policy` observations from the exact manifest digest bound into the catalog.

Checkov receives a scanner-owned list containing exactly the inventoried
Terraform, Kubernetes, and Dockerfile paths. A passing adapter result means that
the checks embedded in this pinned image reported no violation for supported
inputs. It does not prove deployed-state correctness, runtime behavior, cost,
drift, accessibility of external modules, or coverage by future policies.

## Isolation and failure behavior

The adapter:

- mounts a content-verified snapshot read-only and rechecks its digest before
  and after execution;
- runs as the invoking non-root identity with no network, dropped capabilities,
  a read-only container filesystem, bounded resources, and scanner-owned
  scratch;
- disables policy downloads, result upload, external modules, variable
  evaluation, and external checks;
- runs outside the target working directory so `.checkov.yml` cannot change the
  scanner-owned policy or hide passed results;
- rejects inline suppressions, parsing errors, output drift, unexpected online
  metadata, code disclosure, unbound paths, and unbounded records; and
- accepts bounded graph entity objects only long enough to validate the native
  report, then discards them so raw target configuration is not persisted in
  normalized artifacts.

A verified policy violation becomes a canonical finding in the detailed HTML
report. Unsupported files, ambiguous outcomes, parser failures, and safety
violations never become a pass. `prc scan` still does not modify Terraform,
Kubernetes, Dockerfiles, cloud resources, clusters, or any other target.
