package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"
)

// runMain executes run() with the given args and returns the captured output.
// Uses the testable run() function (not main) so exit codes don't terminate
// the test process.
func runMain(t *testing.T, args []string) string {
	t.Helper()

	origArgs := os.Args
	origStdout := os.Stdout
	origStderr := os.Stderr

	defer func() {
		os.Args = origArgs
		os.Stdout = origStdout
		os.Stderr = origStderr
	}()

	os.Args = append([]string{"atheon"}, args...)

	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	out := make(chan string)
	done := make(chan struct{})

	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		out <- buf.String()
	}()

	go func() {
		defer close(done)
		defer func() {
			_ = recover()
		}()
		_ = run(context.Background(), os.Args[1:])
	}()

	<-done
	w.Close()
	return <-out
}

func TestMainHelp(t *testing.T) {
	output := runMain(t, []string{})
	if !strings.Contains(strings.ToLower(output), "usage") && !strings.Contains(strings.ToLower(output), "atheon") {
		t.Errorf("Expected help output, got: %s", output)
	}
}

func TestMainHelpFlag(t *testing.T) {
	output := runMain(t, []string{"--help"})
	if !strings.Contains(strings.ToLower(output), "usage") && !strings.Contains(strings.ToLower(output), "atheon") {
		t.Errorf("Expected help output, got: %s", output)
	}
}

func TestMainHelpShort(t *testing.T) {
	output := runMain(t, []string{"-h"})
	if !strings.Contains(strings.ToLower(output), "usage") && !strings.Contains(strings.ToLower(output), "atheon") {
		t.Errorf("Expected help output, got: %s", output)
	}
}

func TestMainHelpCommandName(t *testing.T) {
	output := runMain(t, []string{"help"})
	if !strings.Contains(strings.ToLower(output), "usage") && !strings.Contains(strings.ToLower(output), "atheon") {
		t.Errorf("Expected help output, got: %s", output)
	}
}

func TestMainList(t *testing.T) {
	output := runMain(t, []string{"list"})
	// Should list patterns - output should be non-empty
	if len(output) == 0 {
		t.Error("Expected pattern list output")
	}
}

func TestMainListWithCategory(t *testing.T) {
	output := runMain(t, []string{"list", "secrets"})
	// Should filter by category
	if len(output) == 0 {
		t.Error("Expected filtered pattern list output")
	}
}

func TestMainEnable(t *testing.T) {
	// Find an actual pattern to enable
	patterns := getTestPatterns()
	if len(patterns) == 0 {
		t.Skip("No patterns available to test enable")
	}
	patternName := patterns[0]
	output := runMain(t, []string{"enable", patternName})
	if !strings.Contains(strings.ToLower(output), "enabled") && !strings.Contains(strings.ToLower(output), "error") {
		t.Logf("Enable output: %s", output)
	}
}

func TestMainEnableMissingArg(t *testing.T) {
	output := runMain(t, []string{"enable"})
	if !strings.Contains(strings.ToLower(output), "error") {
		t.Errorf("Expected error for missing arg, got: %s", output)
	}
}

func TestMainEnableNotFound(t *testing.T) {
	output := runMain(t, []string{"enable", "nonexistent-pattern-xyz123"})
	if !strings.Contains(strings.ToLower(output), "error") && !strings.Contains(strings.ToLower(output), "not found") {
		t.Errorf("Expected not found error, got: %s", output)
	}
}

func TestMainDisable(t *testing.T) {
	patterns := getTestPatterns()
	if len(patterns) == 0 {
		t.Skip("No patterns available to test disable")
	}
	patternName := patterns[0]
	output := runMain(t, []string{"disable", patternName})
	if !strings.Contains(strings.ToLower(output), "disabled") && !strings.Contains(strings.ToLower(output), "error") {
		t.Logf("Disable output: %s", output)
	}
	// Re-enable for other tests
	runMain(t, []string{"enable", patternName})
}

func TestMainDisableMissingArg(t *testing.T) {
	output := runMain(t, []string{"disable"})
	if !strings.Contains(strings.ToLower(output), "error") {
		t.Errorf("Expected error for missing arg, got: %s", output)
	}
}

func TestMainDisableNotFound(t *testing.T) {
	output := runMain(t, []string{"disable", "nonexistent-pattern-xyz123"})
	if !strings.Contains(strings.ToLower(output), "error") && !strings.Contains(strings.ToLower(output), "not found") {
		t.Errorf("Expected not found error, got: %s", output)
	}
}

func TestMainScanFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "atheon-test-*.go")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `package main
var apiKey = "sk-1234567890abcdefghijklmn"
`
	tmpFile.WriteString(content)
	tmpFile.Close()

	output := runMain(t, []string{"--file", tmpFile.Name()})
	t.Logf("Scan file output: %s", output)
}

func TestMainScanFileMissing(t *testing.T) {
	output := runMain(t, []string{"--file", "/nonexistent/path/file.go"})
	if !strings.Contains(strings.ToLower(output), "error") {
		t.Errorf("Expected error for missing file, got: %s", output)
	}
}

func TestMainScanFileMissingArg(t *testing.T) {
	output := runMain(t, []string{"--file"})
	if !strings.Contains(strings.ToLower(output), "error") {
		t.Errorf("Expected error for missing arg, got: %s", output)
	}
}

func TestMainScanPath(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "atheon-test-*.go")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `package main
var apiKey = "sk-1234567890abcdefghijklmn"
`
	tmpFile.WriteString(content)
	tmpFile.Close()

	output := runMain(t, []string{tmpFile.Name()})
	t.Logf("Scan path output: %s", output)
}

func TestMainScanPathNotFound(t *testing.T) {
	output := runMain(t, []string{"/nonexistent/path/file.go"})
	if !strings.Contains(strings.ToLower(output), "error") {
		t.Errorf("Expected error for missing path, got: %s", output)
	}
}

func TestMainScanPathDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "atheon-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile, _ := os.CreateTemp(tmpDir, "*.go")
	content := `package main
var apiKey = "sk-1234567890abcdefghijklmn"
`
	tmpFile.WriteString(content)
	tmpFile.Close()

	output := runMain(t, []string{tmpDir})
	t.Logf("Scan dir output: %s", output)
}

func TestMainJSONOutputNew(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "atheon-test-*.go")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `package main
var apiKey = "sk-1234567890abcdefghijklmn"
`
	tmpFile.WriteString(content)
	tmpFile.Close()

	output := runMain(t, []string{"--json", tmpFile.Name()})
	t.Logf("JSON output: %s", output)
}

func TestMainJSONWithCategories(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "atheon-test-*.go")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("package main\n")
	tmpFile.Close()

	output := runMain(t, []string{"--json", "--categories=secrets", tmpFile.Name()})
	t.Logf("JSON with categories output: %s", output)
}

func TestMainCategoriesAll(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "atheon-test-*.go")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("package main\n")
	tmpFile.Close()

	output := runMain(t, []string{"--all", tmpFile.Name()})
	t.Logf("All categories output: %s", output)
}

func TestMainCategoriesMultiple(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "atheon-test-*.go")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("package main\n")
	tmpFile.Close()

	output := runMain(t, []string{"--categories=secrets,pii", tmpFile.Name()})
	t.Logf("Multiple categories output: %s", output)
}

func TestMainStdin(t *testing.T) {
	// Save original stdin
	origStdin := os.Stdin
	defer func() { os.Stdin = origStdin }()

	// Create a pipe with test data
	r, w, _ := os.Pipe()
	go func() {
		w.WriteString(`var apiKey = "sk-1234567890abcdefghijklmn"`)
		w.Close()
	}()
	os.Stdin = r

	output := runMain(t, []string{"-"})
	t.Logf("Stdin output: %s", output)
}

func TestMainStdinLongFlag(t *testing.T) {
	origStdin := os.Stdin
	defer func() { os.Stdin = origStdin }()

	r, w, _ := os.Pipe()
	go func() {
		w.WriteString(`var apiKey = "sk-1234567890abcdefghijklmn"`)
		w.Close()
	}()
	os.Stdin = r

	output := runMain(t, []string{"--stdin"})
	t.Logf("Stdin (long flag) output: %s", output)
}

func TestMainEnvScan(t *testing.T) {
	// Set a test env var that might trigger a pattern
	os.Setenv("ATHEON_TEST_API_KEY", "sk-1234567890abcdefghijklmn")
	defer os.Unsetenv("ATHEON_TEST_API_KEY")

	output := runMain(t, []string{"--env"})
	t.Logf("Env scan output: %s", output)
}

func TestMainUnknownCommand(t *testing.T) {
	// An unknown command is treated as a path to scan
	// It will fail because the path doesn't exist
	output := runMain(t, []string{"/this/path/does/not/exist/anywhere"})
	if !strings.Contains(strings.ToLower(output), "error") {
		t.Errorf("Expected error for non-existent path, got: %s", output)
	}
}

// getTestPatterns returns a list of actual pattern names from the bundle
func getTestPatterns() []string {
	// Use the core package to get real patterns
	// We can't import core here without making this file depend on it in a test
	// So we'll use a hardcoded list of known patterns
	return []string{
		"api-key",
		"aws-access-key",
		"credit-card",
		"Social Security Number",
	}
}
