# Structured Reporting

Atheon can produce scan findings in multiple formats: plain text, JSON, SARIF, and HTML.
Use `--format=<fmt>` or `--json` to select the output format.

---

## Output Formats

### `text` (default)

Human-readable summary of findings with file, line, and matched content:

```
$ atheon --format=text scanme.txt
secrets: openai-api-key     scanme.txt:3   sk-abcdefghijklmnopqrstuvwxyz
pii: aws-access-key         scanme.txt:7   AKIAIOSFODNN7EXAMPLE
```

### `json`

Pretty-printed JSON for machine consumption:

```json
{
  "version": "0.4.0",
  "generatedAt": "2026-06-20T14:30:00Z",
  "scanType": "file",
  "target": "scanme.txt",
  "stats": {
    "files": 1,
    "bytes": 128,
    "elapsedMs": 4
  },
  "findings": [
    {
      "pattern": "openai-api-key",
      "file": "scanme.txt",
      "line": 3,
      "content": "sk-abcdefghijklmnopqrstuvwxyz"
    }
  ]
}
```

The `--json` flag is an alias for `--format=json` and is accepted in any position in the argument list, even after the scan target.

### `sarif`

[SARIF 2.1.0](https://docs.github.com/en/code-security/code-scanning/integrating-with-code-scanning/sarif-support-for-code-scanning) format for GitHub Code Scanning and other SARIF-compatible tools:

```
$ atheon --format=sarif scanme.txt
{
  "version": "2.1.0",
  "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
  "runs": [
    {
      "tool": {
        "driver": {
          "name": "atheon",
          "version": "0.4.0",
          "informationUri": "https://github.com/aliasfoxkde/Atheon",
          "rules": [
            {
              "id": "openai-api-key",
              "name": "OpenAI API Key",
              "shortDescription": { "text": "Detects OpenAI API keys" },
              "properties": { "tags": ["secret", "openai"] }
            }
          ]
        }
      },
      "results": [
        {
          "ruleIndex": 0,
          "message": { "text": "OpenAI API key found: sk-abcdefghijklmnopqrstuvwxyz" },
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": { "uri": "scanme.txt" },
                "region": { "startLine": 3 }
              }
            }
          ]
        }
      ]
    }
  ]
}
```

Each finding maps to a SARIF `result` referencing a `rule` in `tool.driver.rules`.
Rules are deduplicated by pattern name.

### `html`

Self-contained single HTML file with inline CSS. Open the output in any browser — no external dependencies:

```html
$ atheon --format=html scanme.txt > findings.html
$ open findings.html
```

The HTML report includes a findings table with pattern, file, line, and matched content columns, and a summary header with the scan version and timestamp.

---

## CLI Usage

```
atheon --format=<fmt> <path>   output format: text (default), json, sarif, html
atheon --json <path>           shorthand for --format=json; accepted in any position
```

`--format=<fmt>` and `--json` may appear anywhere in the argument list (before or after the scan target). The flag is stripped before the scan path is resolved.

| Flag | Format | Notes |
|------|--------|-------|
| `--format=text` | text | Default; human-readable |
| `--format=json` | JSON | Pretty-printed with 2-space indent |
| `--format=sarif` | SARIF 2.1.0 | GitHub Code Scanning compatible |
| `--format=html` | HTML | Self-contained, no external dependencies |
| `--json` | JSON | Shorthand; accepted in any argument position |

---

## Report Structure

All four formats are derived from a single `core.Report` struct:

```go
type Report struct {
    Version     string           // atheon version string
    GeneratedAt time.Time        // when the scan ran
    ScanType    string           // "file", "dir", "env", "stdin"
    Target      string           // path or source that was scanned
    Stats       Stats            // files scanned, bytes, elapsed time
    Findings    []Finding        // individual findings
    Errors      []error          // any errors encountered (JSON/SARIF only)
}
```

`core.Render(report, format)` dispatches to the appropriate private renderer:

- `renderText` — line-by-line human output
- `renderJSON` — `json.MarshalIndent` with `"  "` spacing
- `renderSARIF` — SARIF 2.1.0 document with rule indexing and deduplication
- `renderHTML` — self-contained HTML with inline CSS

---

## Exit Codes

| Situation | Exit Code |
|-----------|-----------|
| No findings | 0 |
| One or more findings | 1 |
| Scan error (path not found, read error, etc.) | 1 |
| `--help` or `--version` | 0 |
| `update` command | 0 on success, 1 on failure |

---

## Examples

```bash
# Text output (default)
atheon ./src

# JSON output for automation
atheon --json ./src > findings.json

# SARIF for GitHub Code Scanning upload
atheon --format=sarif ./src > atheon-results.sarif

# HTML report to open in a browser
atheon --format=html ./src > findings.html

# --json after the path (same result as before it)
atheon ./src --json

# --format in any position
atheon --format=sarif --categories=secrets ./src
```
