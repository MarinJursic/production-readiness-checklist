# Kubernetes workload analysis

The focused `prc/kubernetes@0.1` profile performs bounded, no-execution
repository analysis of detected Kubernetes workload manifests. Its six
assertions cover non-root execution, resource requests and limits, privileged
host access, privilege escalation, Linux capabilities, and seccomp profiles.
The same assertions are included in `prc/core-repository@1.0`.

Run only the Kubernetes assertions with:

```bash
prc scan \
  --target PATH \
  --catalog-root PATH_TO_RELEASE \
  --profile prc/kubernetes
```

## Security properties

The native checks inspect `Pod`, `Deployment`, `StatefulSet`, `DaemonSet`,
`ReplicaSet`, `Job`, and `CronJob` pod specifications. For ordinary, init, and
ephemeral containers where the field is defined by Kubernetes, they verify:

- `runAsNonRoot` is affirmed at pod or container scope and `runAsUser: 0` is
  absent;
- host networking, PID, and IPC namespaces remain unset or false, Windows
  HostProcess remains unset or false, `hostPath` volumes are absent, and
  containers are not privileged;
- Linux containers explicitly set `allowPrivilegeEscalation: false`;
- Linux containers drop `ALL` capabilities and add back at most
  `NET_BIND_SERVICE`;
- Linux containers inherit or select a `RuntimeDefault` or `Localhost` seccomp
  profile; and
- ordinary and init containers declare nonempty resource requests and limits.

The security behavior follows the relevant Baseline and Restricted fields in
the authoritative
[Kubernetes Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/).
The Linux-only privilege-escalation, capability, and seccomp requirements are
not applied when `spec.os.name` is explicitly `windows`, matching the standard's
operating-system distinction. The
[Kubernetes security-context documentation](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)
explains how those fields affect container isolation.

## Evidence and failure behavior

Each result binds content-addressed evidence to the inventoried manifest. A
failure reports bounded file, line, and column locations without executing
`kubectl`, Helm, Kustomize, manifest hooks, or target code. A file over the
native parser limit, malformed YAML, a duplicate mapping key, an alias, a
mutation after inventory, or another unprovable parse condition returns an
execution error and `unknown`, never a pass.

## Deliberate limitations

This profile is not a complete Pod Security Admission implementation. It does
not render templates or overlays, inspect custom workload resources, resolve
external values, validate every Kubernetes API field, prove that admission
policy is enabled, inspect an image, contact a cluster, or prove what is
deployed. It also does not infer a risk exception for workloads that
intentionally need host or elevated access. Those cases require explicit,
separately authorized evidence or a signed time-bounded exception; the native
repository result remains a failure.
