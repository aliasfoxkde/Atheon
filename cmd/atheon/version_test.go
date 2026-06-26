package main

import (
	"os/exec"
	"strings"
	"testing"
)

// TestVersionFlag tests the --version flag via subprocess to avoid os.Exit side effects.
func TestVersionFlag(t *testing.T) {
	bin, cleanup := buildTestBinary(t)
	defer cleanup()

	out, err := exec.Command(bin, "--version").CombinedOutput()
	if err != nil {
		t.Fatalf("--version flag error: %v\noutput: %s", err, out)
	}

	if !strings.Contains(string(out), "atheon") {
		t.Errorf("Version output should contain 'atheon', got: %s", out)
	}

	if !strings.Contains(string(out), "dev") && !strings.Contains(string(out), "v") {
		t.Errorf("Version output should contain version number, got: %s", out)
	}
}

// TestDevVersion tests that dev version works correctly
func TestDevVersion(t *testing.T) {
	if version != "dev" {
		t.Logf("Note: version is '%s' (expected 'dev' for development builds)", version)
	}

	// Version should not be empty
	if version == "" {
		t.Error("Version should not be empty")
	}
}

// TestVersionFlagWithJSON verifies that the --json --version combination
// (and other flag orders) print the version cleanly. Before Wave 6, the
// --version check ran before the --json strip, so `atheon --json --version`
// fell into the default branch and errored with "path not found: --version".
func TestVersionFlagWithJSON(t *testing.T) {
	bin, cleanup := buildTestBinary(t)
	defer cleanup()

	for _, args := range [][]string{
		{"--version"},
		{"--json", "--version"},
		{"--sarif", "--version"},
	} {
		t.Run(args[0]+"+rest", func(t *testing.T) {
			out, err := exec.Command(bin, args...).CombinedOutput()
			if err != nil {
				t.Fatalf("%v flag combo error: %v\noutput: %s", args, err, out)
			}
			if !strings.Contains(string(out), "atheon") {
				t.Errorf("expected 'atheon' in output, got: %s", out)
			}
			if strings.Contains(string(out), "path not found") {
				t.Errorf("flag combo should not error with 'path not found', got: %s", out)
			}
		})
	}
}
