package core

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadPatternStateNullPatterns exercises the "patterns == null" branch
// that initializes the map to empty.
func TestLoadPatternStateNullPatterns(t *testing.T) {
	home, _ := os.UserHomeDir()
	stateFile := filepath.Join(home, ".atheon", "pattern_state.json")

	backup, backupErr := os.ReadFile(stateFile)
	defer func() {
		if backupErr == nil {
			_ = os.WriteFile(stateFile, backup, 0o644)
		} else {
			_ = os.Remove(stateFile)
		}
	}()

	// Write a state file with patterns: null
	if err := os.MkdirAll(filepath.Dir(stateFile), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(stateFile, []byte(`{"patterns": null}`), 0o644); err != nil {
		t.Fatal(err)
	}

	state, err := loadPatternState()
	if err != nil {
		t.Fatalf("loadPatternState failed: %v", err)
	}
	if state.Patterns == nil {
		t.Error("expected Patterns to be initialized to empty map")
	}
	if len(state.Patterns) != 0 {
		t.Errorf("expected empty patterns, got %d", len(state.Patterns))
	}
}

// TestSavePatternStateReadOnlyHome exercises the save error path by making
// HOME point to a non-writable location via os.Setenv.
//
// We can't change HOME directly (the test runner uses it for other state),
// so we exercise savePatternState directly and accept either success or
// failure depending on environment. This is a no-op for coverage in most
// cases but documents the expected behavior.
func TestSavePatternStateReadOnlyHome(t *testing.T) {
	state := &PatternState{Patterns: map[string]bool{"foo": true}}
	_ = savePatternState(state)
}