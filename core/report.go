package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Format is the output format for a scan report.
type Format string

const (
	FormatText  Format = "text"
	FormatJSON  Format = "json"
	FormatSARIF Format = "sarif"
	FormatHTML  Format = "html"
)

// Report summarises a complete scan for rendering in any supported Format.
// It is produced by the Render function and is not part of the stable public
// API (the public API is the Render function itself).
type Report struct {
	// Version is the atheon binary version that produced the report.
	Version string `json:"version"`
	// GeneratedAt is when the scan finished.
	GeneratedAt time.Time `json:"generatedAt"`
	// ScanType describes what was scanned: "file", "dir", "string", or "env".
	ScanType string `json:"scanType"`
	// Target is the path, URL, or identifier that was scanned.
	Target string `json:"target,omitempty"`
	// Stats holds the scan statistics.
	Stats Stats `json:"stats"`
	// Findings is every match produced by the scan.
	Findings []Finding `json:"findings"`
	// Errors collects any errors that occurred during the scan (as
	// opposed to per-file WalkErrors, which live in Stats).
	Errors []error `json:"errors,omitempty"`
}

// Render produces a string representation of the given Report in the
// requested format. It is the public API entry point for structured output.
//
// The returned string is formatted for human readability (text, html) or
// machine parsing (json, sarif) depending on the requested Format.
func Render(r Report, format Format) string {
	switch format {
	case FormatJSON:
		return renderJSON(r)
	case FormatSARIF:
		return renderSARIF(r)
	case FormatHTML:
		return renderHTML(r)
	default:
		return renderText(r)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Text renderer
// ─────────────────────────────────────────────────────────────────────────────

func renderText(r Report) string {
	var b strings.Builder
	if len(r.Findings) == 0 {
		b.WriteString("no findings.\n")
	} else {
		for _, f := range r.Findings {
			loc := f.File
			if f.Line > 0 {
				loc = fmt.Sprintf("%s:%d", f.File, f.Line)
			}
			b.WriteString(f.Pattern)
			b.WriteString("  ")
			b.WriteString(loc)
			b.WriteString("\n")
			if f.Content != "" {
				b.WriteString(" ")
				b.WriteString(redact(f.Content))
				b.WriteString("\n")
			}
		}
		b.WriteString(fmt.Sprintf("\n%d finding(s)\n", len(r.Findings)))
	}
	if r.Stats.Files > 0 {
		b.WriteString(fmt.Sprintf("scanned %d file(s)  %s  %dms\n",
			r.Stats.Files, formatBytes(r.Stats.Bytes), r.Stats.ElapsedMs))
	}
	return b.String()
}

// redact redacts secret content for display, matching the CLI's redact logic.
func redact(s string) string {
	if len(s) <= 8 {
		return "***"
	}
	return s[:4] + "****" + s[len(s)-4:]
}

// formatBytes mirrors the CLI helper of the same name.
func formatBytes(b int64) string {
	if b >= 1<<20 {
		return fmt.Sprintf("%.1f MB", float64(b)/(1<<20))
	}
	if b >= 1<<10 {
		return fmt.Sprintf("%.1f KB", float64(b)/(1<<10))
	}
	return fmt.Sprintf("%d B", b)
}

// ─────────────────────────────────────────────────────────────────────────────
// JSON renderer
// ─────────────────────────────────────────────────────────────────────────────

func renderJSON(r Report) string {
	// Mirrors the existing CLI JSON output shape so existing consumers
	// that parse "pattern", "file", "line", "match" fields see no change.
	type jsonFinding struct {
		Pattern string `json:"pattern"`
		File    string `json:"file"`
		Line    int    `json:"line"`
		Match   string `json:"match"`
	}
	items := make([]jsonFinding, 0, len(r.Findings))
	for _, f := range r.Findings {
		items = append(items, jsonFinding{
			Pattern: f.Pattern,
			File:    f.File,
			Line:    f.Line,
			Match:   redact(f.Content),
		})
	}
	out, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return "{}\n"
	}
	return string(out) + "\n"
}

// ─────────────────────────────────────────────────────────────────────────────
// SARIF 2.1.0 renderer
// ─────────────────────────────────────────────────────────────────────────────

// sarifArtifact represents a file in a SARIF physicalLocation.
type sarifArtifact struct {
	URI        string   `json:"uri"`
	URIBaseID  string   `json:"uriBaseId,omitempty"`
	Properties struct{} `json:"properties,omitempty"`
}

// sarifRegion represents a line/column region in a SARIF artifact.
type sarifRegion struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn,omitempty"`
	EndLine     int `json:"endLine,omitempty"`
	EndColumn   int `json:"endColumn,omitempty"`
}

// sarifPhysicalLocation maps a SARIF result location to a file+region.
type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifact `json:"artifactLocation"`
	Region           sarifRegion   `json:"region,omitempty"`
}

// sarifLocation is a SARIF result location.
type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

// sarifResult is a single SARIF result (one finding).
type sarifResult struct {
	RuleID    string                `json:"ruleId"`
	RuleIndex int                   `json:"ruleIndex"`
	Level     string                `json:"level"`
	Message   struct{ Text string } `json:"message"`
	Locations []sarifLocation       `json:"locations"`
}

// sarifToolDriver identifies the tool (Atheon).
type sarifToolDriver struct {
	Driver struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Rules   []struct {
			ID               string                `json:"id"`
			Name             string                `json:"name"`
			ShortDescription struct{ Text string } `json:"shortDescription,omitempty"`
		} `json:"rules,omitempty"`
	} `json:"driver"`
}

// sarifRun is a single SARIF run.
type sarifRun struct {
	Tool    sarifToolDriver `json:"tool"`
	Results []sarifResult   `json:"results,omitempty"`
}

// sarifLog is the root SARIF v2.1.0 document.
type sarifLog struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema,omitempty"`
	Runs    []sarifRun `json:"runs"`
}

func renderSARIF(r Report) string {
	// Collect unique patterns for the tool.driver.rules section.
	patternSeen := make(map[string]bool)
	var rules []struct {
		ID               string                `json:"id"`
		Name             string                `json:"name"`
		ShortDescription struct{ Text string } `json:"shortDescription,omitempty"`
	}
	ruleIndex := make(map[string]int)
	for _, f := range r.Findings {
		if !patternSeen[f.Pattern] {
			patternSeen[f.Pattern] = true
			ruleIndex[f.Pattern] = len(rules)
			rules = append(rules, struct {
				ID               string                `json:"id"`
				Name             string                `json:"name"`
				ShortDescription struct{ Text string } `json:"shortDescription,omitempty"`
			}{ID: f.Pattern, Name: f.Pattern})
		}
	}

	results := make([]sarifResult, 0, len(r.Findings))
	for _, f := range r.Findings {
		ri := ruleIndex[f.Pattern]
		msg := f.Content
		if f.Line > 0 {
			msg = fmt.Sprintf("%s (line %d)", f.Content, f.Line)
		}
		res := sarifResult{
			RuleID:    f.Pattern,
			RuleIndex: ri,
			Level:     "warning",
			Message:   struct{ Text string }{Text: msg},
			Locations: []sarifLocation{{
				PhysicalLocation: sarifPhysicalLocation{
					ArtifactLocation: sarifArtifact{URI: f.File},
					Region:           sarifRegion{StartLine: f.Line},
				},
			}},
		}
		results = append(results, res)
	}

	run := sarifRun{
		Tool: sarifToolDriver{
			Driver: struct {
				Name    string `json:"name"`
				Version string `json:"version"`
				Rules   []struct {
					ID               string                `json:"id"`
					Name             string                `json:"name"`
					ShortDescription struct{ Text string } `json:"shortDescription,omitempty"`
				} `json:"rules,omitempty"`
			}{
				Name:    "Atheon",
				Version: r.Version,
				Rules:   rules,
			},
		},
		Results: results,
	}

	log := sarifLog{
		Version: "2.1.0",
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Runs:    []sarifRun{run},
	}

	out, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return "{}\n"
	}
	return string(out) + "\n"
}

// ─────────────────────────────────────────────────────────────────────────────
// HTML renderer (single self-contained file)
// ─────────────────────────────────────────────────────────────────────────────

func renderHTML(r Report) string {
	var b strings.Builder
	b.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Atheon Scan Report</title>
<style>
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;margin:40px;background:#f8f9fa;color:#212529}
h1{color:#1a1a2e;border-bottom:2px solid #dee2e6;padding-bottom:.5rem}
.summary{background:#e9ecef;border-radius:6px;padding:1rem 1.25rem;margin-bottom:1.5rem;display:flex;gap:2rem;flex-wrap:wrap}
.summary-item{display:flex;flex-direction:column}
.summary-label{font-size:.75rem;text-transform:uppercase;color:#6c757d;letter-spacing:.05em}
.summary-value{font-size:1.25rem;font-weight:600}
table{width:100%;border-collapse:collapse;background:#fff;border-radius:6px;overflow:hidden;box-shadow:0 1px 3px rgba(0,0,0,.1)}
th{background:#1a1a2e;color:#fff;text-align:left;padding:.75rem 1rem;font-size:.875rem}
td{padding:.75rem 1rem;border-bottom:1px solid #dee2e6;vertical-align:top}
tr:last-child td{border-bottom:none}
.pattern{font-weight:600;color:#1a1a2e}
.loc{color:#6c757d;font-size:.875rem}
.content{font-family:ui-monospace,'Cascadia Code','Fira Code',monospace;font-size:.8125rem;background:#f8f9fa;padding:.25rem .5rem;border-radius:4px;word-break:break-all}
.no-findings{background:#d1e7dd;color:#0f5132;padding:1rem;border-radius:6px}
footer{text-align:center;color:#6c757d;font-size:.875rem;margin-top:2rem}
</style>
</head>
<body>
<h1>Atheon Scan Report</h1>
`)
	// Summary
	b.WriteString(`<div class="summary">
<div class="summary-item"><span class="summary-label">Scan Type</span><span class="summary-value">`)
	b.WriteString(htmlEscape(r.ScanType))
	b.WriteString(`</span></div>
<div class="summary-item"><span class="summary-label">Target</span><span class="summary-value">`)
	b.WriteString(htmlEscape(r.Target))
	b.WriteString(`</span></div>
<div class="summary-item"><span class="summary-label">Findings</span><span class="summary-value">`)
	b.WriteString(fmt.Sprintf("%d", len(r.Findings)))
	b.WriteString(`</span></div>
<div class="summary-item"><span class="summary-label">Files Scanned</span><span class="summary-value">`)
	b.WriteString(fmt.Sprintf("%d", r.Stats.Files))
	b.WriteString(`</span></div>
<div class="summary-item"><span class="summary-label">Duration</span><span class="summary-value">`)
	b.WriteString(fmt.Sprintf("%dms", r.Stats.ElapsedMs))
	b.WriteString(`</span></div>
<div class="summary-item"><span class="summary-label">Generated</span><span class="summary-value">`)
	b.WriteString(htmlEscape(r.GeneratedAt.Format(time.RFC3339)))
	b.WriteString(`</span></div>
</div>
`)
	if len(r.Findings) == 0 {
		b.WriteString(`<div class="no-findings">No secrets detected.</div>
`)
	} else {
		b.WriteString(`<table>
<thead><tr><th>Pattern</th><th>Location</th><th>Match</th></tr></thead>
<tbody>
`)
		for _, f := range r.Findings {
			loc := f.File
			if f.Line > 0 {
				loc = fmt.Sprintf("%s:%d", f.File, f.Line)
			}
			b.WriteString("<tr>")
			b.WriteString(`<td class="pattern">`)
			b.WriteString(htmlEscape(f.Pattern))
			b.WriteString("</td>")
			b.WriteString(`<td class="loc">`)
			b.WriteString(htmlEscape(loc))
			b.WriteString("</td>")
			b.WriteString(`<td class="content">`)
			b.WriteString(htmlEscape(redact(f.Content)))
			b.WriteString("</td>")
			b.WriteString("</tr>\n")
		}
		b.WriteString("</tbody>\n</table>\n")
	}
	b.WriteString(`<footer>Generated by Atheon `)
	b.WriteString(htmlEscape(r.Version))
	b.WriteString(`</footer>
</body>
</html>
`)
	return b.String()
}

func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}
