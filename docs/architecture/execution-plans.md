# Execution plans and capability policy

`prc plan` produces a content-addressed `prc.plan/v0.6` artifact before any
external adapter can run. The artifact records both what should be evaluated and
which capabilities the selected execution mode permits.

## Dependency graph

The `nodes` array is a deterministic topological ordering:

1. one inventory node binds the complete target inventory;
2. zero or more immutable adapter nodes depend on that inventory;
3. one assertion node per profile assertion depends on the inventory and any
   adapter evidence it requires; and
4. one gate node depends on every assertion node.

Every dependency must name an earlier node. Duplicate IDs, forward references,
missing assertion nodes, and unavailable implementation IDs fail closed. The
engine evaluates assertion nodes in this recorded order; it does not reconstruct
an implicit execution order at scan time.

The `implementations` array is a deduplicated registry view. It records the exact
implementation ID, implementation kind, all assertions that use it, required
capabilities, and whether the scanner has that exact version available. The
`adapters` array similarly records immutable adapter IDs and manifest digests,
observation kinds, capabilities, and plan-time authorization. Runtime and image
availability are checked separately before launch.

## Execution modes

The current vertical slice supports two modes:

| Mode | Workspace | Scratch | Process | Network | Secrets |
| --- | --- | --- | --- | --- | --- |
| `inspect` | Read-only | Denied | Denied | Denied | Denied |
| `verify-local` | Read-only | Isolated only | Rootless OCI only | Denied | Denied |

`inspect` is the default for `plan` and ordinary native scans. A live adapter
scan requires the operator to pass `--mode verify-local` explicitly. The adapter
must still be bound by an applicable assertion using its exact ID and manifest
SHA-256 digest, and its manifest must pass the stricter OCI runner validation.
Selecting a mode never authorizes an unbound adapter.

Unsupported execution modes are configuration errors. An adapter requirement
that exceeds the selected capability policy produces a blocked plan node and is
denied before the OCI runtime is invoked. The current modes never grant target
writes, general subprocess execution, network access, secret handles, production
access, or external mutation.

## Blocked work remains visible

A plan is still useful when some nodes are blocked. Examples include an
unregistered implementation, an applicable analysis assertion without an
immutable adapter binding, or an adapter that requires a higher mode. The gate
node is then marked blocked, and a scan records an explicit unknown/blocked
assertion result. Missing execution never becomes Pass or Not Applicable.

Version-specific plan v0.1 through v0.6 and run v0.1 through v0.8 schemas
remain frozen for archived consumers. Plans before v0.6 do not contain the
execution DAG or capability policy and therefore require conservative
re-evaluation when compared with current plans.
