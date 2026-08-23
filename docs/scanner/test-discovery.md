# Test discovery evidence

`PRC-A-CORE-010` uses `prc.native.test-suite@0.2` to distinguish a
collectable test declaration from a file that merely has a test-looking name.
The check reads only inventory-bound repository bytes and never imports code or
runs a test command.

## What passes

The assertion passes when the scanner verifies either:

- a conventional test path containing a supported language or framework test
  declaration; or
- a `package.json` `scripts.test` value that is a nonempty, non-placeholder
  command.

The native recognizer covers conventional Go `TestXxx(*testing.T)` functions,
pytest-style `test_` functions, JavaScript and TypeScript `test(...)` or
`it(...)` declarations, and bounded declaration patterns for common Rust,
Ruby, JVM, .NET, PHP, Swift, C, C++, and Scala frameworks. A malformed file,
empty file, comment-only file, skipped JavaScript declaration, or the default
`no test specified` npm placeholder does not pass.

The conventions follow the [Go testing contract](https://go.dev/doc/code#Testing)
and [pytest collection rules](https://docs.pytest.org/en/stable/explanation/goodpractices.html#conventions-for-python-test-discovery).
Package scripts are treated only as declarations; their contents are not run.

## R2 proposal boundary

The scanner-owned missing-test task is narrower than general recognition. It
currently plans new tests only for Go, Python, JavaScript, and TypeScript source
files. Before creating a candidate workspace, the anti-gaming audit reconstructs
the proposed test from the validated patch and requires:

1. a conventional path;
2. a collectable test declaration;
3. a recognized behavioral assertion or failure check;
4. no suppression, skip, constant assertion, or modification to an existing
   test or specification file.

This prevents an agent from closing the discovery finding with an empty,
comment-only, invocation-only, skipped, or constant-assertion test. The
candidate is then rescanned from fresh bytes and must close the exact finding
without regressing a baseline pass.

## Deliberate limitation

Static discovery does not establish that the test compiles, executes, is
deterministic, covers important behavior, or would fail for the defect it is
intended to catch. The current R2 path therefore remains experimental and does
not edit production code. Scanner-owned sandboxed project commands,
test-first mutation evidence, and independent behavioral verification are
required before broader autonomous remediation can rely on this assertion.

The comprehensive core benchmark includes positive and negative cases for a
real declaration, a test-looking file without a declaration, a meaningful
package test command, and the default placeholder command.
