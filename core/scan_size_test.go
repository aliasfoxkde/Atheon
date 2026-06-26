package core

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeFileWithSize creates a file of exactly `size` bytes at the given
// path. Uses Truncate so it works for size==0 (Seek(size-1) would fail for
// size==0 because there's no negative offset to seek to in a fresh file).
// Sparse-friendly: no actual data is written for non-zero sizes.
func writeFileWithSize(t *testing.T, path string, size int64) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create %s: %v", path, err)
	}
	defer f.Close()
	if err := f.Truncate(size); err != nil {
		t.Fatalf("truncate %s to %d: %v", path, size, err)
	}
}

// TestReadFileCappedUnderCap asserts that a file smaller than the cap
// is read in full. This is the common case (every well-behaved source
// file) and a regression here would silently truncate matches.
func TestReadFileCappedUnderCap(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "small.txt")
	want := []byte("hello world\n")
	if err := os.WriteFile(path, want, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, err := readFileCapped(path, int64(len(want)+1))
	if err != nil {
		t.Fatalf("readFileCapped: %v", err)
	}
	if string(got) != string(want) {
		t.Fatalf("contents: got %q, want %q", got, want)
	}
}

// TestReadFileCappedBoundary asserts that a file of EXACTLY maxBytes
// is read in full (boundary inclusive — the cap is "less than or equal",
// not "strictly less than"). Off-by-one here would either truncate the
// last byte of every file at the cap, or fail to enforce the cap on
// files of exactly maxBytes+1.
func TestReadFileCappedBoundary(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "exact.txt")
	writeFileWithSize(t, path, 1024)
	got, err := readFileCapped(path, 1024)
	if err != nil {
		t.Fatalf("readFileCapped at boundary: %v", err)
	}
	if int64(len(got)) != 1024 {
		t.Fatalf("boundary read length: got %d, want 1024", len(got))
	}
}

// TestReadFileCappedOverCap asserts the over-cap path: a single byte
// over the limit must produce ErrFileTooLarge (not a generic I/O
// error and not a partial read). This is the property that bounds
// memory — without it, a 10 GiB log could OOM the scanner.
func TestReadFileCappedOverCap(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "big.txt")
	writeFileWithSize(t, path, 1025)
	_, err := readFileCapped(path, 1024)
	if err == nil {
		t.Fatal("expected ErrFileTooLarge, got nil")
	}
	if !errors.Is(err, ErrFileTooLarge) {
		t.Fatalf("expected ErrFileTooLarge, got %v", err)
	}
	// The error message should identify the file path so an operator
	// scanning 100k files can find the offender without re-running
	// with strace.
	if !strings.Contains(err.Error(), path) {
		t.Errorf("error %q should mention path %q", err.Error(), path)
	}
}

// TestReadFileCappedZeroBytes asserts that a zero-byte file is read
// successfully (zero is a valid size, not an error). Filesystem tools
// (touch, truncate) frequently create zero-byte placeholders; rejecting
// them would surface false-positive scan errors.
func TestReadFileCappedZeroBytes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")
	writeFileWithSize(t, path, 0)
	got, err := readFileCapped(path, 1024)
	if err != nil {
		t.Fatalf("readFileCapped zero bytes: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("zero-byte file: got %d bytes, want 0", len(got))
	}
}

// TestReadFileCappedUnreadable asserts that a permission-denied file
// produces a regular os.ReadFile error (NOT ErrFileTooLarge — the cap
// only triggers on size, not on permission). Callers need to distinguish
// the two: a size skip is benign, a permission error is an environment
// issue the operator should know about.
func TestReadFileCappedUnreadable(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root bypasses file permissions; cannot test perm-denied path")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "noperm.txt")
	if err := os.WriteFile(path, []byte("secret"), 0o000); err != nil {
		t.Fatalf("write: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(path, 0o644) })

	_, err := readFileCapped(path, 1024)
	if err == nil {
		t.Fatal("expected error for unreadable file, got nil")
	}
	if errors.Is(err, ErrFileTooLarge) {
		t.Fatalf("perm error should not be ErrFileTooLarge: %v", err)
	}
}

// TestScanDirSizeCapSurfacesError asserts that ScanDir populates
// stats.Errors with the ErrFileTooLarge sentinel when a worker hits
// the cap, so the CLI can surface a non-zero exit and the MCP server
// can return the count to the caller. Before PR #95, ScanDir's
// per-file goroutines called os.ReadFile directly and bypassed the
// ScanFile size check entirely — a 10 GiB file would just OOM.
func TestScanDirSizeCapSurfacesError(t *testing.T) {
	dir := t.TempDir()
	big := filepath.Join(dir, "big.txt")
	writeFileWithSize(t, big, 2048)

	// Force a tiny cap via opts so we don't have to allocate a 10 MiB
	// fixture. The helper respects opts.MaxFileSize, which is what
	// we're actually testing.
	_, stats, err := ScanDir(context.Background(), dir, ScanOpts{MaxFileSize: 1024})
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	if stats == nil {
		t.Fatal("stats nil")
	}
	if len(stats.Errors) == 0 {
		t.Fatal("expected at least one error for over-cap file")
	}
	var sawTooLarge bool
	for _, e := range stats.Errors {
		if errors.Is(e, ErrFileTooLarge) {
			sawTooLarge = true
			break
		}
	}
	if !sawTooLarge {
		t.Fatalf("expected ErrFileTooLarge in stats.Errors, got: %v", stats.Errors)
	}
}
