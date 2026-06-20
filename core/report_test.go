package core

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestRenderText_NoFindings(t *testing.T) {
	r := Report{Version: "test", Findings: nil, Stats: Stats{}}
	out := Render(r, FormatText)
	if out == "" {
		t.Error("expected non-empty text output")
	}
	if !strings.Contains(out, "no findings") {
		t.Errorf("expected 'no findings' in output: %s", out)
	}
}

func TestRenderText_WithFindings(t *testing.T) {
	r := Report{
		Version:  "test",
		ScanType: "dir",
		Target:   "/src",
		Findings: []Finding{
			{Pattern: "aws-access-key", File: "config.txt", Line: 3, Content: "AKIAIOSFODNN7EXAMPLE"},
		},
		Stats: Stats{Files: 10, Bytes: 1024, ElapsedMs: 50},
	}
	out := Render(r, FormatText)
	if !strings.Contains(out, "aws-access-key") {
		t.Errorf("expected pattern in output: %s", out)
	}
	if !strings.Contains(out, "config.txt:3") {
		t.Errorf("expected location in output: %s", out)
	}
	if !strings.Contains(out, "1 finding(s)") {
		t.Errorf("expected summary in output: %s", out)
	}
	if !strings.Contains(out, "10 file(s)") {
		t.Errorf("expected file count in output: %s", out)
	}
}

func TestRenderJSON_Shape(t *testing.T) {
	r := Report{
		Version:  "test",
		Findings: []Finding{{Pattern: "aws-access-key", File: "f", Line: 1, Content: "AKIAIOSFODNN7EXAMPLE"}},
	}
	out := Render(r, FormatJSON)
	var items []map[string]any
	if err := json.Unmarshal([]byte(out), &items); err != nil {
		t.Fatalf("output is not valid JSON: %s", out)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
	if items[0]["pattern"] != "aws-access-key" {
		t.Errorf("unexpected pattern field: %v", items[0]["pattern"])
	}
	// Content must be redacted
	if items[0]["match"] == "AKIAIOSFODNN7EXAMPLE" {
		t.Error("match should be redacted, got raw value")
	}
}

func TestRenderJSON_Empty(t *testing.T) {
	r := Report{Version: "test", Findings: []Finding{}}
	out := Render(r, FormatJSON)
	var items []map[string]any
	if err := json.Unmarshal([]byte(out), &items); err != nil {
		t.Fatalf("output is not valid JSON: %s", out)
	}
	if len(items) != 0 {
		t.Errorf("expected empty array, got %d items", len(items))
	}
}

func TestRenderSARIF_Structure(t *testing.T) {
	r := Report{
		Version:  "test",
		ScanType: "dir",
		Target:   "/src",
		Findings: []Finding{
			{Pattern: "aws-access-key", File: "config.txt", Line: 5, Content: "AKIAIOSFODNN7EXAMPLE"},
		},
		Stats: Stats{Files: 1, Bytes: 100, ElapsedMs: 10},
	}
	out := Render(r, FormatSARIF)

	var log sarifLog
	if err := json.Unmarshal([]byte(out), &log); err != nil {
		t.Fatalf("output is not valid SARIF JSON: %s", out)
	}

	if log.Version != "2.1.0" {
		t.Errorf("expected SARIF version 2.1.0, got %s", log.Version)
	}
	if len(log.Runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(log.Runs))
	}
	run := log.Runs[0]
	if run.Tool.Driver.Name != "Atheon" {
		t.Errorf("expected tool name 'Atheon', got %s", run.Tool.Driver.Name)
	}
	if run.Tool.Driver.Version != "test" {
		t.Errorf("expected tool version 'test', got %s", run.Tool.Driver.Version)
	}
	if len(run.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(run.Results))
	}
	res := run.Results[0]
	if res.RuleID != "aws-access-key" {
		t.Errorf("expected ruleId 'aws-access-key', got %s", res.RuleID)
	}
	if res.Level != "warning" {
		t.Errorf("expected level 'warning', got %s", res.Level)
	}
	if len(res.Locations) != 1 {
		t.Errorf("expected 1 location, got %d", len(res.Locations))
	}
	loc := res.Locations[0]
	if loc.PhysicalLocation.ArtifactLocation.URI != "config.txt" {
		t.Errorf("expected URI 'config.txt', got %s", loc.PhysicalLocation.ArtifactLocation.URI)
	}
	if loc.PhysicalLocation.Region.StartLine != 5 {
		t.Errorf("expected startLine 5, got %d", loc.PhysicalLocation.Region.StartLine)
	}
}

func TestRenderSARIF_RuleDeduplication(t *testing.T) {
	// Two findings with the same pattern should produce one rule entry
	// and both results should reference the same ruleIndex.
	r := Report{
		Version: "test",
		Findings: []Finding{
			{Pattern: "aws-access-key", File: "a.txt", Line: 1, Content: "key1"},
			{Pattern: "aws-access-key", File: "b.txt", Line: 2, Content: "key2"},
			{Pattern: "openai-api-key", File: "c.txt", Line: 3, Content: "key3"},
		},
	}
	out := Render(r, FormatSARIF)
	var log sarifLog
	if err := json.Unmarshal([]byte(out), &log); err != nil {
		t.Fatalf("not valid JSON: %s", out)
	}
	run := log.Runs[0]
	if len(run.Tool.Driver.Rules) != 2 {
		t.Errorf("expected 2 unique rules, got %d", len(run.Tool.Driver.Rules))
	}
	if len(run.Results) != 3 {
		t.Errorf("expected 3 results, got %d", len(run.Results))
	}
	// Both aws-access-key results should have the same ruleIndex.
	var awsIdx int
	for i, res := range run.Results {
		if res.RuleID == "aws-access-key" {
			awsIdx = res.RuleIndex
			if i > 0 && run.Results[i-1].RuleID == "aws-access-key" {
				if run.Results[i-1].RuleIndex != awsIdx {
					t.Error("duplicate pattern should share ruleIndex")
				}
			}
		}
	}
}

func TestRenderHTML_Structure(t *testing.T) {
	r := Report{
		Version:     "v99.0.0",
		GeneratedAt: time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC),
		ScanType:    "dir",
		Target:      "/src",
		Findings: []Finding{
			{Pattern: "aws-access-key", File: "config.txt", Line: 5, Content: "AKIAIOSFODNN7EXAMPLE"},
		},
		Stats: Stats{Files: 10, Bytes: 2048, ElapsedMs: 123},
	}
	out := Render(r, FormatHTML)
	for _, needle := range []string{
		"<!DOCTYPE html>",
		"<title>Atheon Scan Report</title>",
		`<span class="summary-value">dir</span>`,
		`<span class="summary-value">/src</span>`,
		`<span class="summary-value">1</span>`,     // findings count
		`<span class="summary-value">10</span>`,    // files scanned
		`<span class="summary-value">123ms</span>`, // elapsed
		`Generated by Atheon v99.0.0`,
		"aws-access-key",
		"config.txt:5",
	} {
		if !strings.Contains(out, needle) {
			t.Errorf("expected %q in HTML output", needle)
		}
	}
	// No findings case
	r2 := Report{Version: "test", Findings: nil, GeneratedAt: time.Now()}
	out2 := Render(r2, FormatHTML)
	if !strings.Contains(out2, "No secrets detected") {
		t.Error("expected 'No secrets detected' in empty HTML")
	}
}
