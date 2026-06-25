// Scaffolding for the combined-regex compile-error logging that Wave 8
// PR #95 introduces in core/bundle.go. The bundler currently builds one
// big combined regex per category for fast pre-filtering; if it fails to
// compile, the failure is silently dropped (see bundle.go:259-262), so
// operators only learn about the broken category when matches stop
// appearing. PR #95 wires slog.Warn into the failure path.
//
// This file documents the expected behaviour so the test exists in
// `go test -list` and so PR #95 can drop in concrete assertions without
// renaming or restructuring. Skip-based placeholder keeps CI green in
// the meantime.
//
// The eventual test (PR #95) will:
//   - capture slog output via a custom slog.Handler
//   - inject a category whose patterns cannot be combined (e.g. nested
//     quantifier) by overriding one pattern's regex at test time
//   - call rebuildActiveScanners (or whatever the constructor becomes)
//   - assert a Warn-level record was emitted with the category name and
//     the underlying compile error

package core

import (
	"testing"
)

// TestCombinedRegexCompileError_Scaffold is a placeholder until PR #95
// wires slog.Warn into bundle.go's combined-regex fallback path. When
// the helper exists, replace this with a logger-capturing assertion
// (see file comment for the shape).
func TestCombinedRegexCompileError_Scaffold(t *testing.T) {
	t.Skip("combined-regex compile-error log ships in Wave 8 PR #95 — see plan:fix/wave8-runner-safety")
}
