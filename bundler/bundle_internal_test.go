package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBundleEmpty walks an empty community directory and verifies bundle()
// produces an empty-but-valid gzip+json output.
func TestBundleEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "out.bundle")

	n, err := bundle(tmpDir, outPath)
	if err != nil {
		t.Fatalf("bundle failed: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 patterns from empty dir, got %d", n)
	}

	// Verify output is valid gzip + JSON
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("could not read output: %v", err)
	}
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("output not gzip: %v", err)
	}
	defer gz.Close()

	var defs []patternDef
	if err := json.NewDecoder(gz).Decode(&defs); err != nil {
		t.Fatalf("could not decode bundle: %v", err)
	}
	if len(defs) != 0 {
		t.Errorf("expected empty defs, got %d", len(defs))
	}
}

// TestBundleOnePattern builds a community/secrets/test.yaml and verifies
// the bundle() function picks it up correctly.
func TestBundleOnePattern(t *testing.T) {
	tmpDir := t.TempDir()
	communityDir := filepath.Join(tmpDir, "community", "secrets")
	if err := os.MkdirAll(communityDir, 0o755); err != nil {
		t.Fatal(err)
	}
	patternPath := filepath.Join(communityDir, "test.yaml")
	if err := os.WriteFile(patternPath, []byte(`name: my-pattern
match: '\bMY_[A-Z0-9]{8}\b'
`), 0o644); err != nil {
		t.Fatal(err)
	}

	outPath := filepath.Join(tmpDir, "out.bundle")
	n, err := bundle(filepath.Join(tmpDir, "community"), outPath)
	if err != nil {
		t.Fatalf("bundle failed: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 pattern, got %d", n)
	}

	// Decode and verify
	data, _ := os.ReadFile(outPath)
	gz, _ := gzip.NewReader(bytes.NewReader(data))
	defer gz.Close()
	var defs []patternDef
	if err := json.NewDecoder(gz).Decode(&defs); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if len(defs) != 1 {
		t.Fatalf("expected 1 def, got %d", len(defs))
	}
	if defs[0].Name != "my-pattern" {
		t.Errorf("unexpected name: %s", defs[0].Name)
	}
	if defs[0].Category != "secrets" {
		t.Errorf("unexpected category: %s", defs[0].Category)
	}
	if !defs[0].Enabled {
		t.Error("expected enabled=true (default)")
	}
}

// TestBundleDisabledFlag verifies that an explicit enabled: false is honored.
func TestBundleDisabledFlag(t *testing.T) {
	tmpDir := t.TempDir()
	categoryDir := filepath.Join(tmpDir, "community", "secrets")
	if err := os.MkdirAll(categoryDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// enabled: false explicitly
	if err := os.WriteFile(filepath.Join(categoryDir, "a.yaml"), []byte(`name: a
match: '\bA\b'
enabled: false
`), 0o644); err != nil {
		t.Fatal(err)
	}

	outPath := filepath.Join(tmpDir, "out.bundle")
	if _, err := bundle(filepath.Join(tmpDir, "community"), outPath); err != nil {
		t.Fatalf("bundle failed: %v", err)
	}

	data, _ := os.ReadFile(outPath)
	gz, _ := gzip.NewReader(bytes.NewReader(data))
	defer gz.Close()
	var defs []patternDef
	_ = json.NewDecoder(gz).Decode(&defs)
	if len(defs) != 1 || defs[0].Enabled {
		t.Errorf("expected enabled=false, got %+v", defs)
	}
}

// TestBundleMissingFields returns an error when a YAML file is missing name
// or match fields.
func TestBundleMissingFields(t *testing.T) {
	tmpDir := t.TempDir()
	categoryDir := filepath.Join(tmpDir, "community", "secrets")
	if err := os.MkdirAll(categoryDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Missing match
	if err := os.WriteFile(filepath.Join(categoryDir, "bad.yaml"), []byte(`name: bad
`), 0o644); err != nil {
		t.Fatal(err)
	}

	outPath := filepath.Join(tmpDir, "out.bundle")
	_, err := bundle(filepath.Join(tmpDir, "community"), outPath)
	if err == nil {
		t.Error("expected error for missing fields, got nil")
	}
	if !strings.Contains(err.Error(), "missing name or match") {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestBundleBadYAML returns an error when YAML is malformed.
func TestBundleBadYAML(t *testing.T) {
	tmpDir := t.TempDir()
	categoryDir := filepath.Join(tmpDir, "community", "secrets")
	if err := os.MkdirAll(categoryDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(categoryDir, "broken.yaml"), []byte(`: : :
`), 0o644); err != nil {
		t.Fatal(err)
	}

	outPath := filepath.Join(tmpDir, "out.bundle")
	_, err := bundle(filepath.Join(tmpDir, "community"), outPath)
	if err == nil {
		t.Error("expected error for bad YAML, got nil")
	}
}

// TestBundleMissingCommunity returns an error when the community dir doesn't
// exist.
func TestBundleMissingCommunity(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "out.bundle")
	_, err := bundle(filepath.Join(tmpDir, "does-not-exist"), outPath)
	if err == nil {
		t.Error("expected error for missing community dir, got nil")
	}
}
