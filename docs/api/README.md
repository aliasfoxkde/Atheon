# Atheon API Documentation

Complete API reference for integrating Atheon into Go applications and workflows.

## Package

```
github.com/aliasfoxkde/Atheon-Enhanced/core
```

Go 1.21+ required. All scan functions accept a `context.Context` for cancellation.

---

## Types

### Finding

```go
type Finding struct {
    Pattern string // pattern name that matched
    File    string // source path, or "env:KEY" for environment scans
    Line    int    // 1-indexed line number (0 for env scans)
    Content string // trimmed matching line or, for env scans, the matching value
}
```

### Stats

```go
type Stats struct {
    Files     int   // files whose contents were scanned (binary/skipped files excluded)
    Bytes     int64 // total bytes scanned
    ElapsedMs int64 // wall-clock duration in milliseconds
}
```

### PatternDef

```go
type PatternDef struct {
    Name     string
    Category string
    Pattern  string // regex string
    Enabled  bool
}
```

### Pattern (interface)

```go
type Pattern interface {
    Name() string
    Category() string
    Matches(line string) bool
    Enabled() bool
    SetEnabled(enabled bool)
}
```

---

## Scan Functions

### ScanFile

```go
func ScanFile(ctx context.Context, path string) ([]Finding, *Stats, error)
```

Scans a single file. Returns an empty slice (not nil) and zero-value Stats when the
file has no matches.

```go
findings, stats, err := core.ScanFile(ctx, "config.yaml")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%d finding(s) in %d ms\n", len(findings), stats.ElapsedMs)
```

### ScanDir

```go
func ScanDir(ctx context.Context, root string) ([]Finding, *Stats, error)
```

Recursively scans a directory. Respects `.atheonignore` files. Binary files are
skipped automatically.

```go
findings, stats, err := core.ScanDir(ctx, "/path/to/project")
```

### ScanString

```go
func ScanString(ctx context.Context, content, source string) []Finding
```

Scans an in-memory string. `source` is used as the `File` field in returned
findings (use any meaningful label, e.g., `"stdin"` or `"request-body"`).

```go
findings := core.ScanString(ctx, requestBody, "request-body")
```

### ScanEnv

```go
func ScanEnv(ctx context.Context) []Finding
```

Scans the current process environment variables. Each finding has `File` set to
`"env:VARIABLE_NAME"` and `Line` set to 0.

```go
findings := core.ScanEnv(ctx)
```

---

## Pattern Management

### Register

```go
func Register(p Pattern)
```

Registers a custom pattern. Custom patterns are included in all subsequent scans.

### All

```go
func All() []Pattern
```

Returns all registered patterns (bundle patterns + custom).

### Categories

```go
func Categories() []string
```

Returns the list of distinct category names from the loaded bundle.

### SetActiveCategories

```go
func SetActiveCategories(cats []string)
```

Restricts scanning to the specified categories. Pass `nil` to enable all categories.

```go
// Scan for secrets and PII only
core.SetActiveCategories([]string{"secrets", "pii"})
findings, _, _ := core.ScanDir(ctx, ".")

// Reset to all categories
core.SetActiveCategories(nil)
```

### EnablePattern / DisablePattern

```go
func EnablePattern(name string) bool
func DisablePattern(name string) bool
```

Enable or disable a pattern by name. Returns `true` if the pattern was found.

### SetPatternEnabled

```go
func SetPatternEnabled(name string, enabled bool) bool
```

Single function to enable or disable by name.

### ListEnabledPatterns / ListDisabledPatterns

```go
func ListEnabledPatterns() []string
func ListDisabledPatterns() []string
```

### EnableAllPatterns

```go
func EnableAllPatterns()
```

Enables every loaded pattern regardless of saved state.

### ValidatePattern

```go
func ValidatePattern(def PatternDef) error
```

Validates a pattern definition — checks that `Name` is non-empty and `Pattern` is
a valid Go regular expression.

---

## Bundle Management

### DownloadBundle

```go
func DownloadBundle(ctx context.Context) error
```

Downloads the latest pattern bundle from the release URL and reloads patterns.

### ReloadBundle

```go
func ReloadBundle()
```

Re-initializes patterns from the embedded bundle (undoes any runtime enable/disable
changes).

### InitializePatternState

```go
func InitializePatternState() error
```

Persists current enable/disable state to disk so it survives process restarts.

---

## Custom Pattern Example

```go
package main

import (
    "context"
    "fmt"
    "regexp"

    "github.com/aliasfoxkde/Atheon-Enhanced/core"
)

type myPattern struct {
    re *regexp.Regexp
}

func (p *myPattern) Name() string             { return "internal-api-token" }
func (p *myPattern) Category() string         { return "custom" }
func (p *myPattern) Matches(line string) bool { return p.re.MatchString(line) }
func (p *myPattern) Enabled() bool            { return true }
func (p *myPattern) SetEnabled(_ bool)        {}

func main() {
    core.Register(&myPattern{re: regexp.MustCompile(`INTERNAL_[A-Z0-9]{32}`)})

    findings, stats, err := core.ScanDir(context.Background(), ".")
    if err != nil {
        panic(err)
    }
    fmt.Printf("%d finding(s) across %d file(s) in %d ms\n",
        len(findings), stats.Files, stats.ElapsedMs)
}
```

---

## MCP Server

Build and register the MCP server so AI assistants (Claude, Cursor, etc.) can scan
files and strings on demand:

```bash
go build -o atheon-mcp ./cmd/mcp
```

Add to your MCP config (`~/.config/claude/mcp.json` or `.mcp.json`):

```json
{
  "mcpServers": {
    "atheon": {
      "command": "/path/to/atheon-mcp"
    }
  }
}
```

**Tools exposed by the MCP server:**

| Tool | Description |
|------|-------------|
| `scan_directory` | Scan a path, optional category filter |
| `scan_string` | Scan an in-memory string |
| `list_patterns` | List patterns with optional category filter |

---

## CLI Quick Reference

```bash
# Scan a directory
atheon .

# Specific categories
atheon --categories=secrets,pii .

# JSON output (--json must precede the path)
atheon --json --categories=secrets . > findings.json

# List loaded patterns and count
atheon list

# Enable/disable patterns persistently
atheon enable stripe-api-key
atheon disable todo-comments
```

Exit code is 0 when there are no findings, 1 when findings are detected.

---

*Package: `github.com/aliasfoxkde/Atheon-Enhanced/core` — Go 1.21+*
