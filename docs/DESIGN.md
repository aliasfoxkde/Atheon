# Design

The why behind Atheon's architecture and API choices. Pair this
with [`ARCHITECTURE.md`](ARCHITECTURE.md), which describes the
what.

## The problem we're solving

Engineers commit secrets. It happens. The question is whether
the secret sits in the repository for hours or whether it
gets caught before the push.

The dominant tools in this space are either:

- **Server-side scanners** that only see what reaches the
  remote. They catch leaks after the fact, when rotating the
  secret is already expensive.
- **Heavy local linters** that take minutes to run, write
  findings to a custom format, and require a config file
  before the first invocation.

Atheon is a **fast local scanner that ships sane defaults**.
Run it once with no arguments and it does the right thing.
Run it with `--json` and you get machine-readable output for
a pre-commit hook or a CI step.

## Goals

1. **No configuration to start.** `atheon .` should catch the
   common cases out of the box.
2. **Fast.** A first scan of a small repo should finish in
   under a second. A large monorepo should finish in tens of
   seconds, not minutes.
3. **Composable.** Embedding Atheon as a library must be as
   easy as `import "github.com/.../core"` and calling one
   function. No global init, no flag-parsing side effects.
4. **Honest.** No inflated pattern counts, no marketing
   numbers. The README says `179 patterns, 19 categories`
   because that's what `atheon list` reports, full stop.
5. **Predictable failure modes.** Errors are typed. Cancellation
   propagates. The CLI never panics.

## Non-goals

- **Auto-fix.** We do not rewrite files, open PRs, or rotate
  secrets. We find; someone else fixes.
- **A central server.** Atheon is a single binary and a Go
  library. There is no hosted service.
- **A graphical interface.** CLIs and libraries first; UIs
  are community projects.

## Key decisions

### 1. Patterns as data, not code

Patterns live in YAML files in `community/`. The match engine
is a single Go file (`core/match.go`) that consumes a bundle.
Adding a pattern does not touch the engine.

This decision is the single most important one in the project.
It means:

- Non-Go-contributors can ship patterns.
- The engine can be tested once and frozen.
- Pattern bugs do not take down the binary.
- The bundle can be regenerated, signed, and distributed
  without re-releasing the scanner.

### 2. One Go package for the library

The public API is `github.com/.../core`. There is no
`atheonlib`, no `atheon-sdk`, no v2/next/legacy split. Adding a
package is cheap; maintaining a split is expensive, and we
have no users who would benefit from the split yet.

### 3. `context.Context` everywhere

Every public function takes a context first. This is the
idiomatic Go choice and it gives us cancellation, deadlines,
and tracing for free. The CLI translates `Ctrl-C` into a
context cancel; the MCP server translates the JSON-RPC
`notifications/cancelled` into a context cancel.

### 4. Redaction at the matcher, not the formatter

The matcher is the only place that ever sees the secret. It
emits a redacted snippet, never the original. This means:

- The CLI is safe to pipe to `tee`.
- The JSON output is safe to upload to a CI artifact store.
- The library is safe to embed in a process that may log
  findings to a third-party service.

### 5. Embedded bundle, not network-loaded by default

The pattern bundle is `//go:embed`-ed into the binary. The
`atheon update` command exists for users who want a fresher
bundle than the one shipped with their binary, but the default
is offline and self-contained.

This decision has trade-offs:

- **Pro.** The binary works in air-gapped CI runners without
  any network access.
- **Pro.** There is no startup-time DNS or HTTP request that
  could fail.
- **Con.** A binary that is six months old has a six-month-old
  pattern set. Users must run `atheon update` to get the
  latest.

The trade-off favours the boring, predictable path. Users who
want freshness can run `update`.

### 6. Single-binary bundler

The `bundler/` directory builds into its own binary. The
bundler is not a runtime dependency of the scanner — it is a
developer tool. Keeping it separate means:

- The scanner binary is small.
- The bundler can be rewritten in any language without
  affecting users.
- A contributor who wants to ship a pattern does not need
  to compile the scanner.

### 7. Errors as values

The library defines sentinel errors and returns them from
public functions. Callers branch on `errors.Is(err,
core.ErrBinaryFile)` rather than string-matching. This is
standard Go and it lets us evolve the error messages without
breaking callers.

### 8. The MCP server is a peer of the CLI

The MCP server (`cmd/mcp/`) is not a wrapper around the CLI;
it is a separate entry point that calls directly into the
library. This means the MCP server has the same performance
characteristics as the CLI, and we can change the CLI's
argument syntax without breaking MCP clients.

## What we'd do differently

Hindsight, with the project at 1.0:

- **Bundle signing.** We should have shipped signed bundles
  from day one. It is in the roadmap for the 2.x series.
- **Streaming JSON output.** The CLI currently buffers the full
  result before printing. A streaming JSON encoder would let
  the first finding appear instantly, which matters for
  pre-commit hooks where the user is waiting.
- **A `v2` Go API earlier.** The 1.0 API is fine but a clean
  break to put `ctx` first and use option structs would have
  been easier at release-time than now.

## See also

- [`ARCHITECTURE.md`](ARCHITECTURE.md) — the structural
  description.
- [`ROADMAP.md`](ROADMAP.md) — where these design choices
  point next.
- [`MIGRATION.md`](MIGRATION.md) — when the v2 Go API lands,
  how to get there.
