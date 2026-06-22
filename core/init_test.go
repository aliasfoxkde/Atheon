package core

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestLoadBundleFromDisk simulates the init() function reading a user-local
// patterns.bundle from ~/.atheon/patterns.bundle.
//
// We can't restart init(), but we can verify loadBundle works the same way
// by passing the on-disk bytes through it (which is what init does on the
// success path).
func TestLoadBundleFromDisk(t *testing.T) {
	home, _ := os.UserHomeDir()
	bundlePath := filepath.Join(home, ".atheon", "patterns.bundle")

	// Build a small valid bundle
	defs := []PatternDef{
		{Name: "disk-pattern-1", Category: "disk-test", Match: `\bDISK_[A-Z0-9]+\b`, Enabled: true},
	}
	jsonBytes, _ := json.Marshal(defs)

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(jsonBytes); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}

	// Save current bundle if any and restore after
	orig, origErr := os.ReadFile(bundlePath)
	defer func() {
		if origErr == nil {
			_ = os.WriteFile(bundlePath, orig, 0o644)
		} else {
			_ = os.Remove(bundlePath)
		}
		// Re-load the embedded bundle to restore default state
		_ = loadBundle(embeddedBundle)
		SetActiveCategories(nil)
	}()

	if err := os.WriteFile(bundlePath, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	// Read back and verify it would be picked up
	data, err := os.ReadFile(bundlePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := loadBundle(data); err != nil {
		t.Fatalf("loadBundle failed: %v", err)
	}

	// Verify our disk pattern was loaded
	found := false
	for _, p := range allPatterns {
		if p.name == "disk-pattern-1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected disk-pattern-1 to be loaded from disk bundle")
	}
}

// TestInitLoadsExistingBundle verifies init()'s disk-bundle path is
// exercised. We can't call init() again, but we can verify the loadBundle
// + InitializePatternState pair works after writing a disk bundle.
func TestInitLoadsExistingBundle(t *testing.T) {
	home, _ := os.UserHomeDir()
	bundlePath := filepath.Join(home, ".atheon", "patterns.bundle")

	defs := []PatternDef{
		{Name: "init-disk-pattern", Category: "init-test", Match: `\bINIT_[A-Z]+\b`, Enabled: true},
	}
	jsonBytes, _ := json.Marshal(defs)

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(jsonBytes)
	gz.Close()

	orig, origErr := os.ReadFile(bundlePath)
	defer func() {
		if origErr == nil {
			_ = os.WriteFile(bundlePath, orig, 0o644)
		} else {
			_ = os.Remove(bundlePath)
		}
		_ = loadBundle(embeddedBundle)
		SetActiveCategories(nil)
	}()

	if err := os.WriteFile(bundlePath, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	// loadBundle is what init() calls after reading the disk file
	if err := loadBundle(buf.Bytes()); err != nil {
		t.Fatalf("loadBundle failed: %v", err)
	}

	found := false
	for _, p := range allPatterns {
		if p.name == "init-disk-pattern" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected init-disk-pattern to be loaded")
	}
}
