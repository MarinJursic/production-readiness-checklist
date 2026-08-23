# Go HTTP client timeout analysis

`PRC-A-GO-001` is a deterministic, no-execution source check for direct calls
to the Go `net/http` package functions `Get`, `Head`, `Post`, and `PostForm`.
The Go documentation defines these functions as wrappers around
`http.DefaultClient`; the default is a usable zero-value `http.Client`, and a
zero `Client.Timeout` means no end-to-end timeout. See the official
[`net/http` package documentation](https://pkg.go.dev/net/http) and
[`DefaultClient` source](https://go.dev/src/net/http/client.go).

The rule requires the timeout decision to be visible at a bounded request or
client boundary instead of depending on mutable process-global client state.
A typical compliant call uses a reused `http.Client` with a positive `Timeout`:

```go
client := &http.Client{Timeout: 10 * time.Second}
response, err := client.Get(endpoint)
```

For APIs that need a different deadline per operation, construct a request with
`http.NewRequestWithContext` and use `Client.Do`. The scanner does not determine
whether a particular numeric duration is suitable for the service; that remains
a project-specific reliability decision.

## Syntax and evidence contract

The implementation parses inventoried `.go` files with Go's syntax tree parser.
It recognizes normal, aliased, and dot imports of `net/http`, while excluding
locally shadowed identifiers. It never type-checks, builds, tests, imports, or
executes the target. Scope evidence binds the content-addressed inventory; a
failure also records hashes for affected files, never source text. The result
summary reports a bounded file, line, column, and helper name; the structured
finding location identifies the affected file.

The check returns Blocked/Unknown when more than 4,096 Go files or 256 MiB of Go
source would need inspection, or when one file exceeds the native 4 MiB parser
limit. A changed or unreadable inventoried file and invalid Go syntax return an
execution error. None of those states can become Pass.

## Deliberate limitations

This is not a whole-program proof that every outbound call has an appropriate
timeout. In particular, it does not currently evaluate:

- direct calls through `http.DefaultClient` or another client variable;
- custom transports, generated clients, RPC libraries, databases, queues, or
  other network stacks;
- request contexts, data flow, build tags, generated code provenance, or whether
  a configured duration is operationally appropriate; or
- dependency source, binaries, deployed configuration, and runtime behavior.

Those broader claims require language-aware data flow and runtime evidence from
separately benchmarked, digest-pinned adapters. This native rule covers one
high-confidence hazard and must not be described as complete outbound-call
timeout validation.
