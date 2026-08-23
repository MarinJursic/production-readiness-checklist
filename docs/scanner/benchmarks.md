# Scanner benchmarks and quality budgets

The benchmark runner measures scanner behavior on small, labeled fixture
repositories. It is a release-regression tool, not a readiness score and not a
claim that an assertion is accurate outside the represented fixture classes.

The checked-in `core-native` suite currently covers deterministic examples of:

- Pass and Fail results;
- Not Applicable planning;
- Blocked adapter evidence;
- Manual Review evidence requirements; and
- malformed-input execution errors.

Each case binds its target inventory digest and is scanned twice at one fixed
evaluation time. The report fails its quality gate if expected assessment or
execution states drift, repeated runs differ, precision or recall falls below
the suite budget, or the false-positive rate exceeds it. The content-addressed
corpus digest combines the normalized suite definition with every fixture
inventory digest.

Run the suite locally:

```bash
prc benchmark run \
  --catalog-root . \
  --suite fixtures/benchmarks/core-native/suite.yaml \
  --format human
```

Emit the versioned report used by CI:

```bash
prc benchmark run \
  --catalog-root . \
  --suite fixtures/benchmarks/core-native/suite.yaml \
  --evaluated-at 2026-08-23T12:00:00Z \
  --format json > benchmark-report.json
```

`--evaluated-at` makes the evidence timestamps and run identities reproducible
for a controlled comparison. It does not override target content, catalog, or
fixture identity.

## Interpreting the metrics

For benchmark metrics, an expected `fail` is the positive class. Precision is
the fraction of reported failures that were labeled failures; recall is the
fraction of labeled failures detected; the false-positive rate is measured
against all labeled non-failure outcomes. Exact-state matching remains stricter
than the binary metrics: confusing Blocked, Error, Manual Review, Not
Applicable, or Pass still fails the suite even when it does not change the
failure-class confusion matrix.

The initial suite is deliberately small. Its perfect budget protects known
behavior but does not establish broad real-world accuracy. New assertions need
representative Pass, Fail, unsupported or Not Applicable, and execution-error
fixtures before they can support a default release gate. Cross-platform,
adversarial, performance, and tool-update differential suites remain required
as coverage grows.

## Validated packs

A pack is a versioned distribution claim over a measured subset of the
catalog. The `core-foundation` pack binds each included assertion to its exact
catalog implementation and declares only the outcomes present in the pinned
benchmark suite. Validation fails on catalog drift, suite digest drift,
undeclared profile membership, or overstated outcome coverage:

```bash
prc pack validate \
  --catalog-root . \
  --file packs/core-foundation.yaml \
  --format human
```

The smaller `core-foundation` pack contains three assertions and remains a
fast contract test. The `core-native` pack binds all 30 assertions in the core
profile to a 21-case, 103-expectation fixture corpus:

```bash
prc benchmark run \
  --catalog-root . \
  --suite fixtures/benchmarks/core-native/suite-comprehensive.yaml \
  --evaluated-at 2026-08-23T12:00:00Z
prc pack validate \
  --catalog-root . \
  --file packs/core-native.yaml
```

The v0.2 benchmark contract can materialize five bounded fixture conditions in
a temporary copy: exact text-token replacement, final-newline removal, file
truncation, permission-mode changes, and a synthetic immutable Git HEAD. It
never executes target code, follows symlinks, or modifies the checked-in
fixture. This gives each deterministic native assertion both its relevant
positive and negative outcome; manual and adapter-backed assertions remain
measured only in their fail-closed states. All limitations are part of the pack
manifest and therefore its digest.

Pack membership does not authorize adapter execution and does not change gate
semantics. Detached, scoped Ed25519 verification is supported through an
explicit [publisher trust store](../architecture/publisher-trust.md). The
repository does not yet publish an official release trust store or signed pack,
so checked-in packs remain local integrity contracts rather than a remote trust
channel.

## Adapter release fixtures

Repository benchmarks measure assertion decisions over project inventories.
Recorded adapter fixtures instead measure the untrusted-tool boundary: parser
limits, protocol rejection, manifest-declared output kinds, incomplete statuses,
and determinism. Run the separate gate with `prc adapter fixture-validate`.
Passing it is necessary for an adapter release but never substitutes for a
tool-specific detection benchmark. See the
[sandboxed adapter protocol](../architecture/adapters.md#inspect-and-validate)
for the exact contract and limitations.
