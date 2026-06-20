# API Reference

The public Go API of the `atheon` module. This file is the
contract: anything documented here is supported; anything
not documented here is internal and may change without
notice. Pair this with [`DESIGN.md`](DESIGN.md) for the
rationale behind the shapes.

> **Stability:** the v1 API is stable. The planned v2 API
> is documented in [`MIGRATION.md`](MIGRATION.md).

## Package layout

```
github.com/aliasfoxkde/atheon
├── core/        — scanner engine, public API
├── bundler/     — pattern bundle builder (library)
└── internal/    — implementation details, not importable
```

The CLI and the MCP server live under `cmd/` and import
from `core/`; they are not part of the library API.

## `core` package

The scanner engine. All public functions take a
`context.Context` as their first argument and return a
`*Result` (or a partial result plus an error).

### Types

#### `Finding`

```go
type Finding struct {
    RuleID    string `json:"rule_id"`
    Category  string `json:"category"`
    Severity  Severity `json:"severity"`
    File      string `json:"file"`
    Line      int    `json:"line"`
    Column    int    `json:"column"`
    Match     string `json:"match"`
    Excerpt   string `json:"excerpt"`
    Entropy   float64 `json:"entropy,omitempty"`
    Masked    bool   `json:"masked"`
}
```

A single match. `File` is the path as supplied to the scan
function; callers that resolve symlinks or chdir should
re-base `File` to the canonical path. `Masked` is true if
the `Match` field has been redacted; the full value is not
exported.

#### `Severity`

```go
type Severity int

const (
    SeverityInfo Severity = iota
    SeverityLow
    SeverityMedium
    SeverityHigh
    SeverityCritical
)
```

#### `Stats`

```go
type Stats struct {
    FilesScanned   int           `json:"files_scanned"`
    FilesSkipped   int           `json:"files_skipped"`
    BytesScanned   int64         `json:"bytes_scanned"`
    FindingsTotal  int           `json:"findings_total"`
    Duration       time.Duration `json:"duration_ns"`
    PatternsLoaded int           `json:"patterns_loaded"`
}
```

Aggregate counters. `FilesSkipped` counts files the
scanner chose not to read (binary, too large, ignore
match, etc.) — not files that produced no findings.

#### `Result`

```go
type Result struct {
    Findings []Finding `json:"findings"`
    Stats    Stats     `json:"stats"`
}
```

The combined return value of every scan function. The
`Findings` slice is sorted by `(file, line, column)` and
is safe to range over without further locking.

#### `Options`

```go
type Options struct {
    IgnorePatterns  []string
    IncludeGlobs    []string
    ExcludeGlobs    []string
    SeverityMin     Severity
    MaxFileSize     int64
    MaxLineLength   int
    Redaction       bool
    Workers         int
    Patterns        []Pattern // overrides the default bundle
    ReportProgress  func(Progress)
}

func DefaultOptions() Options { /* … */ }
```

Optional knobs. `Options` is passed by value to the scan
functions; the zero value is valid and means "all
defaults". `DefaultOptions()` returns a fully-populated
`Options` with sensible values for an interactive CLI run.

`ReportProgress`, if non-nil, is called once per file. It
runs on a worker goroutine; implementations must be
thread-safe and must not block.

#### `Progress`

```go
type Progress struct {
    FilesScanned int
    FilesTotal   int // -1 if unknown
    CurrentFile  string
}
```

#### `Pattern`

```go
type Pattern struct {
    ID          string
    Category    string
    Severity    Severity
    Regex       *regexp.Regexp // compiled
    Keywords    []string
    MinEntropy  float64
    Description string
}
```

A single pattern. Constructing a `Pattern` directly is
discouraged; load a bundle with `LoadBundle` instead.

#### `Error`

```go
type Error struct {
    Op  string
    Err error
}

func (e *Error) Error() string
func (e *Error) Unwrap() error
```

Wrapped errors. Use `errors.Is` and `errors.As` to inspect
them. Common sentinels:

| Sentinel | Meaning |
|---|---|
| `core.ErrSecretInEnv` | A secret was found in an environment variable. |
| `core.ErrBinaryFile` | The scanner skipped a binary file. Not fatal; reported in `Stats.FilesSkipped`. |
| `core.ErrFileTooLarge` | The scanner skipped a file over `Options.MaxFileSize`. |
| `core.ErrLineTooLong` | A line exceeded `Options.MaxLineLength`; the line is reported as a `lines-too-long` finding. |
| `core.ErrBundleCorrupt` | The embedded bundle failed to parse. |
| `core.ErrContextCanceled` | The context was canceled mid-scan. |

### Functions

#### `ScanFile`

```go
func ScanFile(ctx context.Context, path string, opts ...Options) (*Result, error)
```

Scan a single file. `path` is opened, read line by line,
and matched against every pattern. The function does not
follow symlinks; resolve them yourself if you need to.

```go
result, err := core.ScanFile(ctx, "/srv/app/.env")
if err != nil {
    return err
}
for _, f := range result.Findings {
    log.Printf("rule=%s file=%s line=%d", f.RuleID, f.File, f.Line)
}
```

#### `ScanDir`

```go
func ScanDir(ctx context.Context, root string, opts ...Options) (*Result, error)
```

Scan every file under `root` recursively. Symlinks are
**not** followed. Hidden directories (names starting with
`.`) are visited unless excluded by `Options.ExcludeGlobs`.

```go
result, err := core.ScanDir(ctx, "/srv/app",
    core.DefaultOptions(),
)
```

#### `ScanString`

```go
func ScanString(ctx context.Context, content, name string, opts ...Options) (*Result, error)
```

Scan an in-memory string. `name` is used as the `File`
field of any `Finding`; pick something meaningful for your
caller (e.g. the editor buffer name).

```go
result, err := core.ScanString(ctx, source, "main.go")
```

#### `ScanEnv`

```go
func ScanEnv(ctx context.Context, opts ...Options) (*Result, error)
```

Scan the current process environment. Returns
`core.ErrSecretInEnv` (wrapped) if any pattern matched.
The MCP server uses this to refuse to start in a
contaminated environment.

#### `ScanStdin`

```go
func ScanStdin(ctx context.Context, opts ...Options) (*Result, error)
```

Scan stdin as a single file. The `File` field of any
finding is set to `<stdin>`.

#### `LoadBundle`

```go
func LoadBundle(path string) ([]Pattern, error)
```

Load a pattern bundle from a file produced by the
`bundler` command. Used by callers that ship a custom
bundle alongside their binary.

#### `EmbeddedBundle`

```go
func EmbeddedBundle() ([]Pattern, error)
```

Return the bundle embedded in the `atheon` binary at
compile time. This is what every other `Scan*` function
uses by default.

### Error handling

All scan functions return a `*Result` and a non-nil
error when the scan could not complete. A partial result
may be returned alongside the error; inspect it before
deciding how to proceed.

```go
result, err := core.ScanDir(ctx, root)
if err != nil {
    if errors.Is(err, core.ErrContextCanceled) {
        return ctx.Err()
    }
    // result may still be usable
    log.Printf("scan incomplete: %v (findings so far: %d)",
        err, len(result.Findings))
}
```

The CLI and the MCP server follow this pattern: never
discard a partial result on error.

## `bundler` package

A library API for building pattern bundles. The CLI is a
thin wrapper around it.

### Types

#### `Config`

```go
type Config struct {
    SourceDir   string
    OutputPath  string
    MinPatterns int
    Workers     int
}
```

#### `Build`

```go
func Build(ctx context.Context, cfg Config) (*Stats, error)
```

Walk `SourceDir`, compile every `*.yaml` pattern, and
write a bundle to `OutputPath`. Returns the bundle
statistics (count by category, count by severity, total
size).

## Versioning

The API is versioned via the import path:

```go
import "github.com/aliasfoxkde/atheon/core"
```

Breaking changes to `core` or `bundler` ship in v2; see
[`MIGRATION.md`](MIGRATION.md) for the planned v1 → v2
delta. The current API is **v1** and is stable.

## Compatibility promises

- **Public types and functions in `core/` and `bundler/`**
  are stable within a major version. New fields may be
  added to `Options`, `Result`, and `Stats`; existing
  fields will not be renamed or removed.
- **Sentinel errors** (`ErrSecretInEnv` etc.) are stable.
  New sentinels may be added.
- **`Finding` JSON shape** is stable: any change is a
  major version bump. Schema versioning for the JSON
  output is tracked in `internal/report/schema.go`.

## What's not part of the API

- Anything under `internal/` is not importable and may
  change without notice.
- The CLI's flag set, exit codes, and output format are
  documented in [`CLI.md`](CLI.md) and are versioned
  independently from the library API.
- The MCP server's tool surface is documented in
  [`MCP.md`](MCP.md) and follows its own deprecation
  policy.
- The pattern YAML schema is documented in
  [`PATTERNS.md`](PATTERNS.md) and is stable across the
  v1 line.
