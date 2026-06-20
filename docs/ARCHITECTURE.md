# Architecture

Atheon is a small, layered system. This document describes the
moving parts, how data flows through them, and the rules each
layer plays by. For a deep dive into the runtime internals see
[`architecture/SYSTEM_ARCHITECTURE.md`](architecture/SYSTEM_ARCHITECTURE.md);
for the rationale behind the layering see [`DESIGN.md`](DESIGN.md).

## Bird's-eye view

```
┌─────────────────────────────────────────────────────────────┐
│  CLI         cmd/atheon/main.go   cmd/atheon/cli.go         │
├─────────────────────────────────────────────────────────────┤
│  MCP server  cmd/mcp/                                       │
├─────────────────────────────────────────────────────────────┤
│  Library API  core/runner.go   core/scanner.go              │
├─────────────────────────────────────────────────────────────┤
│  Engines     core/scan*.go (string, file, dir, env, stdin)  │
├─────────────────────────────────────────────────────────────┤
│  Matchers    core/match.go   core/patterns.bundle (embed)   │
├─────────────────────────────────────────────────────────────┤
│  Bundler     bundler/                                       │
└─────────────────────────────────────────────────────────────┘
```

The four layers below the CLI are all in the `core/` Go package
and form the **library API**. The CLI and the MCP server are thin
front-ends that translate their respective inputs into calls into
that library.

## The four layers

### 1. CLI front-ends (`cmd/atheon/`, `cmd/mcp/`)

Each front-end:

- Parses arguments with the standard library's `flag` package.
- Resolves the input source (`--file`, positional path, `--env`,
  `-` for stdin, or `--all` for the bundle).
- Calls into `core/` and writes the result in the requested
  format (text, JSON, or quiet).

The CLI is intentionally thin. All matching logic lives in
`core/` so the library and the CLI cannot drift.

### 2. Library API (`core/runner.go`, `core/scanner.go`)

Exposes the public surface:

```go
findings, stats, err := core.ScanFile(ctx, "path/to/file")
findings, stats, err := core.ScanDir(ctx, "path/to/dir")
findings, stats, err := core.ScanString(ctx, content, "name")
findings, stats, err := core.ScanEnv(ctx)
findings, stats, err := core.ScanStdin(ctx)
```

Every public function takes a `context.Context` first. The
context controls cancellation and (optionally) a deadline.
The function returns `([]Finding, *Stats, error)` so the caller
knows what was found, how much work the scan did, and whether
something went wrong.

The library never exits, never writes to `os.Stdout`, never reads
from `os.Stdin`. It is a pure Go package and it can be embedded
in any host.

### 3. Engines (`core/scan*.go`)

The engines are the workers that turn one input source into
`[]Finding`. Each one implements the same shape:

1. Resolve the input (file path, dir path, env map, string, …).
2. Apply the working-directory `.atheonignore` and `.gitignore`
   rules.
3. Skip files that look binary by extension (the binary-extension
   list lives in `core/binary.go`).
4. Stream the bytes line-by-line through the matcher.

The engines run in parallel where it helps. `ScanDir` uses one
goroutine per CPU (capped) with a semaphore; `ScanString` is
single-threaded because the input is already in memory.

### 4. Matchers and the bundle (`core/match.go`, `core/patterns.bundle`)

The matcher is a hand-written Aho-Corasick-style scanner that
pre-computes a per-pattern decision tree from the bundle. The
bundle is a gzip-compressed JSON blob embedded into the binary
at build time via `//go:embed`.

Patterns are data, not code. Adding a pattern means adding a
YAML file in `community/`, running `bundler/bundler`, and
rebuilding. The match engine itself never changes for a new
pattern.

## Data flow

```
input source
   │
   ▼
resolve + ignore + binary-filter
   │
   ▼
line stream
   │
   ▼
matcher (per-line, anchored or unbounded)
   │
   ▼
[]Finding
```

Each `Finding` is a struct with the pattern name, category,
file/line, and a redacted snippet:

```go
type Finding struct {
    Pattern  string
    Category string
    File     string
    Line     int
    Column   int
    Snippet  string  // redacted; never the full secret
}
```

The snippet is redacted by the matcher — Atheon never logs the
secret in the clear.

## Concurrency model

- **Within a single scan.** `ScanDir` spawns a worker pool. Each
  worker reads one file at a time, scans it, and writes the
  result into a per-worker slice. After the pool drains the
  slices are concatenated. Workers honour `ctx.Done()` between
  files so a cancellation propagates quickly.
- **Across calls.** Each call to a `core.Scan*` function is
  independent. There is no shared mutable state in the library
  except the read-only pattern bundle.
- **Process-wide state.** The `enable`/`disable` commands mutate
  the on-disk `~/.atheon/state.json`. This file is read once at
  start-up by `LoadState` and not re-read; if you `enable`
  something you need to either restart the process or call
  `LoadState` again yourself.

## Failure modes

Atheon follows the Go convention of returning errors and never
panicking in production code paths. The library defines a small
set of sentinel errors (`core/errors.go`) that callers can match
with `errors.Is`:

| Sentinel | Meaning |
|---|---|
| `ErrBinaryFile` | The file's extension marks it as binary and was skipped. |
| `ErrCancelled` | The context was cancelled before the scan completed. |
| `ErrInvalidPath` | The path could not be resolved or does not exist. |
| `ErrTooLarge` | The input exceeds the configured size limit. |
| `ErrSecretInEnv` | The env scan found something — not an error, but a typed signal so callers can branch. |
| `ErrSecretInStdin` | Same as `ErrSecretInEnv`, for stdin. |

The CLI maps these to exit codes: `ErrSecretInEnv` →
exit `1`, everything else → exit `2`. Cancellation (`Ctrl-C`,
context deadline) → exit `130`.

## Boundaries

What lives in each layer and what does not:

- **`cmd/`** owns argument parsing, output formatting, and
  process exit codes. It does not own pattern matching.
- **`core/`** owns pattern matching, ignore rules, and the
  public Go API. It does not own argument parsing or output
  formatting.
- **`bundler/`** owns the YAML → bundle pipeline. It is a
  separate `main` package that builds into its own binary so
  the pattern-bundling tool is not a runtime dependency of the
  scanner.

## See also

- [`DESIGN.md`](DESIGN.md) — the why behind the architecture.
- [`architecture/SYSTEM_ARCHITECTURE.md`](architecture/SYSTEM_ARCHITECTURE.md) —
  detailed walkthrough with diagrams of the matcher and bundler.
- [`architecture/PATTERN_CATEGORIES.md`](architecture/PATTERN_CATEGORIES.md) —
  what each category means and how patterns are organised.
- [`API.md`](API.md) — the public Go API reference.
