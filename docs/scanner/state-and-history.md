# Durable state and run history

When `--state-dir` is supplied, a scan writes immutable canonical JSON records
and transactionally indexes their metadata in `state.sqlite`:

```text
prc-state/
├── state.sqlite
├── evidence/
│   └── ab/<evidence-id>.json
└── runs/
    └── <run-id>.json
```

The JSON records are authoritative. SQLite is a query index that can be rebuilt
from those content-addressed records; it never replaces evidence or changes an
assessment. Indexing verifies every run and evidence identity against its
canonical record before one transaction updates runs, results, evidence
metadata, inventory files, inventory facts, relationships, and the audit event.

## Create private state

Scanner evidence can contain sensitive paths and findings. Keep state outside
the target repository on a local filesystem and restrict it to the current user:

```bash
mkdir -m 0700 /safe/local/path/prc-state

./prc scan \
  --target /path/to/project \
  --catalog-root /path/to/production-readiness-checklist \
  --state-dir /safe/local/path/prc-state \
  --format json \
  --exit-policy never
```

The scanner creates a missing state root with mode `0700` and the database with
mode `0600`. It rejects a pre-existing state root accessible by group or other
users. Windows relies on platform access controls because POSIX mode bits are not
available.

The initial store uses a full-synchronous rollback journal, a five-second busy
timeout, immediate write transactions, strict tables, and connection-level
foreign-key enforcement. WAL is deliberately not enabled: SQLite documents that
[WAL requires all users to be on one host and does not work on network
filesystems](https://www.sqlite.org/wal.html). A remote or shared filesystem is
not a supported state location even with the rollback journal.

## List indexed runs

```bash
./prc history list \
  --state-dir /safe/local/path/prc-state \
  --limit 20
```

Exact filters are available for `--target-name`, `--profile`, and
`--terminal-state`. JSON output conforms to `prc.history/v0.1`:

```bash
./prc history list \
  --state-dir /safe/local/path/prc-state \
  --target-name project \
  --format json > history.json
```

Counts preserve distinct states: Pass, Fail, and unresolved/blocked results are
never averaged into a score.

## Load a canonical run

```bash
./prc history show \
  --state-dir /safe/local/path/prc-state \
  --format json \
  <run-id>
```

`history show` obtains the record path from the index, rejects absolute,
traversing, or symlink-escaping paths, loads the immutable JSON, and recomputes
its run identity before returning it. A missing, modified, or mismatched record
is an error; the database is never treated as sufficient proof by itself.

The implementation enables foreign keys explicitly because SQLite does not
guarantee they are enabled by default, and its integrity audit uses both
[`PRAGMA integrity_check`](https://www.sqlite.org/pragma.html#pragma_integrity_check)
and `PRAGMA foreign_key_check`.

Run both checks and obtain indexed record counts with:

```bash
./prc history check \
  --state-dir /safe/local/path/prc-state \
  --format json
```

A successful JSON response conforms to `prc.state-check/v0.1`. Corruption or a
foreign-key violation is an error and never produces an `integrity: ok` report.
