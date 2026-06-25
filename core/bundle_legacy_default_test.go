package core

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

// TestLoadBundleLegacyDefaultFlip verifies that when a bundle decodes with
// NO patterns reporting enabled=true (the legacy old-format case), the
// loadBundle path:
//   1. flips all patterns to enabled, and
//   2. emits an slog.Info line so the legacy-compat flip is observable
//      instead of silently indistinguishable from an intentional
//      all-disabled bundle.
//
// Regression guard for Wave 6 gap Item 1.
func TestLoadBundleLegacyDefaultFlip(t *testing.T) {
	// Capture slog output via a custom handler.
	var buf bytes.Buffer
	oldHandler := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})))
	defer slog.SetDefault(oldHandler)

	// Build a tiny bundle where every pattern has enabled=false.
	defs := []PatternDef{
		{Name: "legacy-flip-a", Category: "secrets", Match: `AKIA[0-9A-Z]{16}`, Enabled: false},
		{Name: "legacy-flip-b", Category: "secrets", Match: `ghp_[a-zA-Z0-9]{36}`, Enabled: false},
	}
	jsonBytes, err := json.Marshal(defs)
	if err != nil {
		t.Fatal(err)
	}
	var gzBuf bytes.Buffer
	gz := gzip.NewWriter(&gzBuf)
	if _, err := gz.Write(jsonBytes); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}

	if err := loadBundle(gzBuf.Bytes()); err != nil {
		t.Fatalf("loadBundle: %v", err)
	}

	// All patterns should be enabled after the legacy flip.
	for _, p := range All() {
		if !p.Enabled() {
			t.Errorf("pattern %q should be enabled after legacy flip", p.Name())
		}
	}

	// The slog.Info line should appear in the captured output.
	logged := buf.String()
	if !strings.Contains(logged, "bundle had no enabled patterns") {
		t.Errorf("expected legacy-flip log line, got: %q", logged)
	}
	if !strings.Contains(logged, "legacy compatibility") {
		t.Errorf("expected legacy-compat hint in log, got: %q", logged)
	}
}

// TestLoadBundleNoFlipWhenAnyEnabled verifies the log does NOT fire when
// the bundle has at least one enabled pattern — i.e. the log is gated on
// the actual legacy case, not emitted on every load.
func TestLoadBundleNoFlipWhenAnyEnabled(t *testing.T) {
	var buf bytes.Buffer
	oldHandler := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})))
	defer slog.SetDefault(oldHandler)

	defs := []PatternDef{
		{Name: "normal-a", Category: "secrets", Match: `AKIA[0-9A-Z]{16}`, Enabled: true},
		{Name: "normal-b", Category: "secrets", Match: `ghp_[a-zA-Z0-9]{36}`, Enabled: false},
	}
	jsonBytes, err := json.Marshal(defs)
	if err != nil {
		t.Fatal(err)
	}
	var gzBuf bytes.Buffer
	gz := gzip.NewWriter(&gzBuf)
	if _, err := gz.Write(jsonBytes); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}

	if err := loadBundle(gzBuf.Bytes()); err != nil {
		t.Fatalf("loadBundle: %v", err)
	}

	if strings.Contains(buf.String(), "bundle had no enabled patterns") {
		t.Errorf("legacy-flip log should NOT fire when at least one pattern is enabled, got: %q", buf.String())
	}
}