# Go HTTP timeout analysis

The core profile contains two deliberately narrow, deterministic Go checks for
HTTP timeout hazards in non-test `.go` files. Both inspect syntax without
running, building, importing, or type-checking target code. Files ending in
`_test.go` are excluded because production timeout gates should not classify
test-only helpers as deployed network behavior.

## Client package helpers

`PRC-A-GO-001` checks direct calls to the Go `net/http` package functions
`Get`, `Head`, `Post`, and `PostForm`. The Go documentation defines these
functions as wrappers around `http.DefaultClient`; the default is a usable
zero-value `http.Client`, and a zero `Client.Timeout` means no end-to-end
timeout. See the official
[`net/http` package documentation](https://pkg.go.dev/net/http) and
[`DefaultClient` source](https://go.dev/src/net/http/client.go).

The timeout decision must be visible at a bounded request or client boundary
instead of depending on mutable process-global client state. A typical
compliant call uses a reused `http.Client` with a positive `Timeout`:

```go
client := &http.Client{Timeout: 10 * time.Second}
response, err := client.Get(endpoint)
```

For APIs that need a different deadline per operation, construct a request with
`http.NewRequestWithContext` and use `Client.Do`. The scanner does not determine
whether a particular numeric duration is suitable for the service; that remains
a project-specific reliability decision.

## Server package helpers

`PRC-A-GO-002` checks direct calls to the `net/http` package functions
`ListenAndServe`, `ListenAndServeTLS`, `Serve`, and `ServeTLS`. The official Go
[`server.go` source](https://go.dev/src/net/http/server.go) shows that each
helper constructs a new `http.Server` with only an address and/or handler before
serving. The [`Server` documentation](https://pkg.go.dev/net/http#Server) states
that zero or negative read, header, write, and idle timeout values can mean
there is no timeout. The package's own
[`doc.go` guidance](https://go.dev/src/net/http/doc.go) recommends constructing
a custom `Server` when more control is required.

A typical compliant server makes its request timeout policy explicit on a
reused `http.Server` and calls a method on that value:

```go
server := &http.Server{
    Addr:              ":8080",
    Handler:           handler,
    ReadHeaderTimeout: 5 * time.Second,
    ReadTimeout:       30 * time.Second,
    WriteTimeout:      30 * time.Second,
    IdleTimeout:       60 * time.Second,
}
err := server.ListenAndServe()
```

This rule proves only the absence of the four package-level shortcuts. It does
not claim that every field is populated, positive, or operationally suitable.

## Syntax and evidence contract

The implementations parse inventoried non-test `.go` files with Go's syntax
tree parser. They recognize normal, aliased, and dot imports of `net/http`,
while excluding locally shadowed identifiers. They never type-check, build,
test, import, or execute the target. Scope evidence binds the content-addressed
inventory; a failure also records hashes for affected files, never source text. The result
summary reports a bounded file, line, column, and helper name; the structured
assertion result and canonical finding identify the exact affected callsite.

Each check returns Blocked/Unknown when more than 4,096 Go files or 256 MiB of
Go source would need inspection, or when one file exceeds the native 4 MiB
parser limit. A changed or unreadable inventoried file and invalid Go syntax
return an execution error. None of those states can become Pass.

## Deliberate limitations

These rules are not a whole-program proof that every network operation has an
appropriate timeout. In particular, they do not currently evaluate:

- direct calls through `http.DefaultClient` or another client variable;
- custom transports, generated clients, RPC libraries, databases, queues, or
  other network stacks;
- request contexts, data flow, build tags, generated code provenance, or whether
  a configured duration is operationally appropriate;
- whether a custom `http.Server` initializes every relevant timeout field,
  whether later assignments change those fields, or how proxies affect request
  bounds; or
- dependency source, binaries, deployed configuration, and runtime behavior.

Those broader claims require language-aware data flow and runtime evidence from
separately benchmarked, digest-pinned adapters. These native rules cover two
high-confidence hazards and must not be described as complete network-timeout
validation.
