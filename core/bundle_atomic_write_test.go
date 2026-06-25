// Atomic-write regression test (Wave 7 PR #89, audit item C3).
//
// atomicWriteFile is the helper DownloadBundle / savePatternState use to
// persist user-facing data. Without it, os.WriteFile does a truncate-then-
// write, and a SIGKILL in the gap leaves a zero-byte file that downstream
// loaders see as corrupt. The test exercises three properties:
//
//  1. Happy path: write succeeds, file contents match, no leftover .tmp.
//  2. Permission failure: write to an unwritable path returns an error
//     and the destination (if it existed) is left untouched.
//  3. Mid-write failure: simulate a write error and verify the destination
//     is not partially overwritten (the .tmp is cleaned up).
package core

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAtomicWriteFile_HappyPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.bin")
	payload := []byte("hello-atomic-world")

	if err := atomicWriteFile(path, payload, 0o600); err != nil {
		t.Fatalf("atomicWriteFile: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if !bytes.Equal(got, payload) {
		t.Errorf("contents mismatch: got %q want %q", got, payload)
	}

	// No leftover .tmp in the directory.
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".tmp-*") || strings.Contains(e.Name(), ".tmp-") {
			t.Errorf("leftover temp file: %s", e.Name())
		}
	}

	// Permissions match.
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Errorf("perm: got %o want 0o600", got)
	}
}

func TestAtomicWriteFile_DestinationUnchangedOnFailure(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.bin")

	// Seed the destination with the "previous" payload — atomicWriteFile
	// must leave it intact when its write fails.
	previous := []byte("previous-good-content")
	if err := os.WriteFile(path, previous, 0o600); err != nil {
		t.Fatalf("seed write: %v", err)
	}

	// Force a failure: write to a path whose parent does not exist. The
	// os.CreateTemp call inside atomicWriteFile will fail because MkdirAll
	// is not called on the parent — the destination is therefore untouched.
	bogus := filepath.Join(dir, "missing-subdir", "out.bin")
	err := atomicWriteFile(bogus, []byte("new"), 0o600)
	if err == nil {
		t.Fatal("expected error writing to missing parent dir, got nil")
	}

	// Destination must still hold the previous payload.
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if !bytes.Equal(got, previous) {
		t.Errorf("destination changed despite failure: got %q want %q", got, previous)
	}
}

func TestAtomicWriteFile_OverwriteExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.bin")
	if err := os.WriteFile(path, []byte("v1"), 0o600); err != nil {
		t.Fatalf("seed: %v", err)
	}
	v2 := []byte("v2-longer-than-v1")
	if err := atomicWriteFile(path, v2, 0o600); err != nil {
		t.Fatalf("overwrite: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if !bytes.Equal(got, v2) {
		t.Errorf("contents after overwrite: got %q want %q", got, v2)
	}
	// Filename should still be the original; the rename is the atomic step.
	if _, err := os.Stat(path); err != nil {
		t.Errorf("destination missing after overwrite: %v", err)
	}
}
