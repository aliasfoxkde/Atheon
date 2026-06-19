# Atheon API Documentation

Complete API reference for integrating Atheon into your applications and workflows.

## Overview

Atheon provides both a CLI interface and a programmatic Go API for pattern matching and security scanning.

## Go Package API

### Core Package: `github.com/aliasfoxkde/Atheon/core`

#### Pattern Interface

```go
type Pattern interface {
    Name() string
    Category() string
    Matches(line string) bool
    Enabled() bool
    SetEnabled(enabled bool)
}
```

#### Main Functions

##### ScanFile
```go
func ScanFile(path string) ([]Finding, error)
```
Scans a single file and returns detected patterns.

**Parameters:**
- `path string`: Path to the file to scan

**Returns:**
- `[]Finding`: List of detected findings
- `error`: Error if file cannot be read

**Example:**
```go
findings, err := core.ScanFile("config.yaml")
if err != nil {
    log.Fatal(err)
}
for _, finding := range findings {
    fmt.Printf("%s: %s\n", finding.Pattern, finding.Line)
}
```

##### ScanDir
```go
func ScanDir(path string) ([]Finding, *Stats, error)
```
Recursively scans a directory and returns detected patterns with statistics.

**Parameters:**
- `path string`: Path to the directory to scan

**Returns:**
- `[]Finding`: List of detected findings
- `*Stats`: Scanning statistics (files scanned, bytes processed, etc.)
- `error`: Error if directory cannot be read

##### ScanEnv
```go
func ScanEnv() ([]Finding, error)
```
Scans environment variables for sensitive patterns.

##### ScanString
```go
func ScanString(content string) []Finding
```
Scans a string content for patterns.

##### SetActiveCategories
```go
func SetActiveCategories(categories []string)
```
Sets active categories for filtering. `nil` enables all categories.

##### Register
```go
func Register(p Pattern)
```
Registers a custom pattern programmatically.

### Finding Structure

```go
type Finding struct {
    Pattern  string    // Pattern name that matched
    File    string    // File where match occurred
    Line    int       // Line number
    Match   string    // The matched text
    Context string    // Surrounding context (if available)
}
```

### Stats Structure

```go
type Stats struct {
    Files      int64      // Number of files scanned
    Bytes      int64      // Total bytes processed
    Duration   time.Duration // Scanning duration
}
```

## MCP Server API

### Tool: `scan_directory`

Scans a directory and returns findings.

**Parameters:**
```json
{
  "path": "/path/to/directory",
  "categories": ["secrets", "pii"],
  "enabled_only": true
}
```

**Returns:**
```json
{
  "findings": [
    {
      "pattern": "aws-access-key",
      "file": "config.yaml",
      "line": 15,
      "match": "AKIAIOSFODNN7EXAMPLE"
    }
  ],
  "stats": {
    "files": 42,
    "bytes": 1024000,
    "duration_ms": 150
  }
}
```

### Tool: `scan_string`

Scans a string content for patterns.

**Parameters:**
```json
{
  "content": "string to scan",
  "categories": ["secrets"]
}
```

**Returns:**
```json
{
  "findings": [
    {
      "pattern": "api-key",
      "match": "sk_live_1234567890"
    }
  ]
}
```

### Tool: `list_patterns`

Lists available patterns with optional filtering.

**Parameters:**
```json
{
  "categories": ["secrets"],
  "enabled_only": false
}
```

**Returns:**
```json
{
  "patterns": [
    {
      "name": "aws-access-key",
      "category": "secrets",
      "enabled": true,
      "match": "\\b(?:AKIA|ASIA)[0-9A-Z]{16}\\b"
    }
  ]
}
```

## Configuration API

### Pattern State Management

#### Enable Pattern
```go
func EnablePattern(name string) bool
```
Enables a pattern by name. Returns `true` if successful.

#### Disable Pattern
```go
func DisablePattern(name string) bool
```
Disables a pattern by name. Returns `true` if successful.

#### List Patterns
```go
func All() []Pattern
```
Returns all registered patterns.

#### Get Categories
```go
func Categories() []string
```
Returns all available categories.

## Advanced Usage

### Custom Pattern Registration

```go
package main

import (
    "github.com/aliasfoxkde/Atheon/core"
    "regexp"
)

type CustomPattern struct {
    name     string
    category string
    re       *regexp.Regexp
    enabled  bool
}

func (p *CustomPattern) Name() string        { return p.name }
func (p *CustomPattern) Category() string    { return p.category }
func (p *CustomPattern) Matches(line string) bool {
    return p.enabled && p.re.MatchString(line)
}
func (p *CustomPattern) Enabled() bool       { return p.enabled }
func (p *CustomPattern) SetEnabled(e bool)  { p.enabled = e }

func main() {
    custom := &CustomPattern{
        name:     "my-custom-pattern",
        category: "custom",
        re:       regexp.MustCompile(`MY_SECRET_KEY`),
        enabled:  true,
    }

    core.Register(custom)
    // Now scan with your custom pattern included
}
```

### Category Filtering

```go
// Only scan for secrets and PII
core.SetActiveCategories([]string{"secrets", "pii"})
findings, stats, _ := core.ScanDir("/path/to/scan")

// Reset to all categories
core.SetActiveCategories(nil)
```

### JSON Output

```go
findings, _, err := core.ScanDir("/path")
if err != nil {
    log.Fatal(err)
}

jsonBytes, err := json.MarshalIndent(findings, "", "  ")
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(jsonBytes))
```

## Error Handling

### Common Errors

**File Access Errors:**
```go
findings, err := core.ScanFile("/protected/file")
if err != nil {
    if os.IsPermission(err) {
        log.Printf("Permission denied: %v", err)
    } else if os.IsNotExist(err) {
        log.Printf("File not found: %v", err)
    } else {
        log.Printf("Error scanning file: %v", err)
    }
}
```

**Pattern Compilation Errors:**
Pattern compilation errors are logged but don't stop scanning:
```
atheon: skipping "bad-pattern": invalid regex: error parsing regexp: missing closing bracket
```

## Performance Considerations

### Memory Efficiency
Atheon uses streaming and chunked processing for large files:
- **Large files**: Processed in chunks to minimize memory
- **Parallel scanning**: Multiple files processed concurrently
- **Efficient regex**: Pre-compiled patterns for fast matching

### Speed Optimization
- **Category filtering**: Reduces patterns to check
- **File extension filtering**: Skips binary files
- **Directory exclusion**: Avoids unnecessary scanning

### Best Practices
1. **Use category filtering** when you only need specific pattern types
2. **Enable appropriate patterns** for your use case
3. **Configure ignore patterns** to avoid false positives
4. **Use parallel processing** for large codebases

## Integration Examples

### CI/CD Integration

```go
package main

import (
    "github.com/aliasfoxkde/Atheon/core"
    "os"
)

func main() {
    // Scan current directory
    findings, stats, err := core.ScanDir(".")
    if err != nil {
        os.Exit(1)
    }

    // Fail if critical findings detected
    critical := 0
    for _, f := range findings {
        if isCriticalPattern(f.Pattern) {
            critical++
        }
    }

    if critical > 0 {
        os.Exit(1)
    }
}
```

### Web Service Integration

```go
package main

import (
    "encoding/json"
    "net/http"
    "github.com/aliasfoxkde/Atheon/core"
)

type ScanRequest struct {
    Content string   `json:"content"`
    Categories []string `json:"categories"`
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
    var req ScanRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    // Set categories if provided
    if len(req.Categories) > 0 {
        core.SetActiveCategories(req.Categories)
    }

    findings := core.ScanString(req.Content)
    json.NewEncoder(w).Encode(findings)
}
```

---

**API Version**: 1.0.0
**Go Version**: 1.21+
**Last Updated**: 2026-06-19