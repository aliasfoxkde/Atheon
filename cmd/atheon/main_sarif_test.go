package main

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aliasfoxkde/Atheon/core"
)

// captureStdout redirects os.Stdout to a pipe and returns the captured bytes.
// It restores os.Stdout regardless of what happens in f.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	done := make(chan string, 1)
	go func() {
		var sb strings.Builder
		io.Copy(&sb, r) //nolint:errcheck
		done <- sb.String()
	}()

	f()

	w.Close()
	os.Stdout = orig
	out := <-done
	r.Close()
	return out
}

// TestRunSARIF exercises the --sarif flag path through run().
func TestRunSARIF(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "clean.go")
	if err := os.WriteFile(tmp, []byte("package x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	code := run(context.Background(), []string{"--sarif", tmp})
	if code != 0 {
		t.Errorf("expected exit 0 for --sarif on clean file, got %d", code)
	}
}

// TestRunSARIFWithFindings exercises --sarif when findings are present (exit 1).
func TestRunSARIFWithFindings(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "secrets.txt")
	if err := os.WriteFile(tmp, []byte(`AKIAIOSFODNN7EXAMPLE`), 0o644); err != nil {
		t.Fatal(err)
	}
	// exit 1 because findings exist; just don't panic
	_ = run(context.Background(), []string{"--sarif", tmp})
}

// TestRunSARIFOutputIsValidJSON verifies --sarif emits parseable SARIF JSON.
func TestRunSARIFOutputIsValidJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "clean.go")
	if err := os.WriteFile(tmp, []byte("package x\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out := captureStdout(t, func() {
		run(context.Background(), []string{"--sarif", tmp}) //nolint:errcheck
	})

	var sarif map[string]any
	if err := json.Unmarshal([]byte(out), &sarif); err != nil {
		t.Fatalf("--sarif output is not valid JSON: %v\noutput: %s", err, out)
	}
	if sarif["version"] != "2.1.0" {
		t.Errorf("expected SARIF version 2.1.0, got %v", sarif["version"])
	}
	if _, ok := sarif["runs"]; !ok {
		t.Error("SARIF output missing 'runs' key")
	}
}

// TestRunFileWithFindings exercises the --file exit-1 branch (findings found).
func TestRunFileWithFindings(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "secrets.txt")
	// A string that looks like an AWS key to trigger a finding
	if err := os.WriteFile(tmp, []byte(`aws_key = "AKIAIOSFODNN7EXAMPLE"`), 0o644); err != nil {
		t.Fatal(err)
	}
	// May exit 0 or 1 depending on active patterns; just must not panic.
	_ = run(context.Background(), []string{"--file", tmp})
}

// TestRunDefaultPathWithFindings exercises the default-branch file scan exit-1 path.
func TestRunDefaultPathWithFindings(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "secrets.txt")
	if err := os.WriteFile(tmp, []byte(`AKIAIOSFODNN7EXAMPLE`), 0o644); err != nil {
		t.Fatal(err)
	}
	_ = run(context.Background(), []string{tmp})
}

// TestRunSARIFDir exercises --sarif on a directory path.
func TestRunSARIFDir(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "clean.go")
	if err := os.WriteFile(f, []byte("package x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	code := run(context.Background(), []string{"--sarif", dir})
	if code != 0 {
		t.Errorf("expected exit 0 for --sarif on clean dir, got %d", code)
	}
}

// TestPrintFindingsLineZero exercises the f.Line == 0 branch in printFindings
// where loc stays as f.File (no ":N" suffix appended).
func TestPrintFindingsLineZero(t *testing.T) {
	findings := []core.Finding{
		{Pattern: "test", File: "somefile.txt", Line: 0, Content: "secret content here!"},
	}
	out := captureStdout(t, func() {
		printFindings(findings, nil, false, false)
	})
	if !strings.Contains(out, "somefile.txt") {
		t.Errorf("expected file name in output, got: %s", out)
	}
	// Must NOT contain ":0" — the zero-line branch skips the colon-number suffix
	if strings.Contains(out, ":0") {
		t.Errorf("output should not contain ':0' for zero line number, got: %s", out)
	}
}

// TestPrintFindingsStatsZeroFiles exercises the stats.Files == 0 branch
// (stats block is suppressed when no files were scanned).
func TestPrintFindingsStatsZeroFiles(t *testing.T) {
	findings := []core.Finding{}
	stats := &core.Stats{Files: 0, Bytes: 0, ElapsedMs: 0}
	out := captureStdout(t, func() {
		printFindings(findings, stats, false, false)
	})
	// Stats line ("scanned N file(s)...") must not appear when Files == 0
	if strings.Contains(out, "scanned") {
		t.Errorf("expected no stats line for Files=0, got: %s", out)
	}
}

// TestBuildSARIFRulesEmpty exercises buildSARIFRules with no findings.
func TestBuildSARIFRulesEmpty(t *testing.T) {
	rules := buildSARIFRules(nil)
	if rules != nil {
		t.Errorf("expected nil rules for empty findings, got %v", rules)
	}
}
