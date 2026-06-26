package core

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestScanDirFollowsSymlinksByDefault pins the CLI's historical default
// (follow symlinks) so a future change to flip the default doesn't
// silently break scripts that depend on it. The MCP server overrides
// this — see TestScanDirNoFollowSymlinks — but the CLI keeps the old
// behaviour unless --no-follow-symlinks is passed.
func TestScanDirFollowsSymlinksByDefault(t *testing.T) {
	dir := t.TempDir()
	outside := t.TempDir()
	secret := filepath.Join(outside, "secret.txt")
	// Real, known-fireable content (AKIAIOSFODNN7EXAMPLE is the AWS
	// example key from community/secrets/aws.yaml).
	if err := os.WriteFile(secret, []byte("aws_access_key_id=AKIAIOSFODNN7EXAMPLE"), 0o644); err != nil {
		t.Fatalf("write secret: %v", err)
	}
	// Symlink inside the scan root pointing OUT to the secret file.
	// With the default (follow), scanLines reads through the link and
	// finds the secret. Pin this so we don't regress.
	if err := os.Symlink(secret, filepath.Join(dir, "leak.txt")); err != nil {
		t.Skipf("symlink unsupported on this fs: %v", err)
	}

	findings, _, err := ScanDir(context.Background(), dir, ScanOpts{})
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	var saw bool
	for _, f := range findings {
		if f.File == filepath.Join(dir, "leak.txt") && f.Pattern == "aws-access-key" {
			saw = true
		}
	}
	if !saw {
		t.Fatal("expected default ScanDir to follow symlink and report aws-access-key finding in the linked file")
	}
}

// TestScanDirNoFollowSymlinks is the security regression: a symlink
// inside the scan root pointing OUT to /etc/passwd (or any other file
// outside the tree) must NOT be read when NoFollowSymlinks is set.
// Without this guard, an attacker can plant `repo/leak -> /etc/passwd`
// in a PR and exfiltrate the file's contents into the findings
// stream — the canonical "lone symlink" escape.
func TestScanDirNoFollowSymlinks(t *testing.T) {
	dir := t.TempDir()
	outside := t.TempDir()
	secret := filepath.Join(outside, "secret.txt")
	if err := os.WriteFile(secret, []byte("aws_access_key_id=AKIAIOSFODNN7EXAMPLE"), 0o644); err != nil {
		t.Fatalf("write secret: %v", err)
	}
	link := filepath.Join(dir, "leak.txt")
	if err := os.Symlink(secret, link); err != nil {
		t.Skipf("symlink unsupported on this fs: %v", err)
	}

	findings, _, err := ScanDir(context.Background(), dir, ScanOpts{NoFollowSymlinks: true})
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	for _, f := range findings {
		if f.File == link {
			t.Fatalf("NoFollowSymlinks did not block escape: got findings for %s", link)
		}
	}
}

// TestScanDirNoFollowSymlinksLoop guards against an even nastier case:
// a symlink that points back to the scan root (or to one of its
// ancestors). Without NoFollowSymlinks, WalkDir itself would loop
// forever reading the same files. WalkDir actually catches this case
// (it doesn't follow symlinks by default), but the test exists as a
// belt-and-braces guard: if anyone later changes WalkDir's behaviour
// (e.g. switches to filepath.Walk which DOES follow), this test
// catches the regression.
func TestScanDirNoFollowSymlinksLoop(t *testing.T) {
	dir := t.TempDir()
	// loop -> dir's parent (..). Walking through `loop` would revisit
	// dir's parent, then everything underneath it again, ad infinitum.
	if err := os.Symlink(filepath.Join(dir, ".."), filepath.Join(dir, "loop")); err != nil {
		t.Skipf("symlink unsupported on this fs: %v", err)
	}
	// Plain text file so the scan has something to look at.
	if err := os.WriteFile(filepath.Join(dir, "normal.txt"), []byte("nothing of interest"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _, err := ScanDir(ctx, dir, ScanOpts{NoFollowSymlinks: true})
	if err != nil {
		// ctx.DeadlineExceeded is the expected failure mode if the loop
		// ever regresses — surface that explicitly so a reviewer
		// doesn't think the test is flaky.
		if ctx.Err() == context.DeadlineExceeded {
			t.Fatalf("ScanDir looped (deadline exceeded) — symlink guard regressed")
		}
		// Other errors (e.g. permission) are environment-specific and
		// not the regression we're guarding against.
		t.Logf("ScanDir returned non-fatal error (not a regression): %v", err)
	}
}

// TestScanDirNoFollowSymlinksDangling covers the "symlink to nowhere"
// case. WalkDir reports a symlink even if the target doesn't exist;
// without the guard, ReadFile on a dangling link returns an error that
// gets logged as a per-file failure in Stats.Errors. With the guard,
// the link is silently skipped — no error, no scan, no panic.
func TestScanDirNoFollowSymlinksDangling(t *testing.T) {
	dir := t.TempDir()
	dangling := filepath.Join(dir, "nowhere.txt")
	if err := os.Symlink(filepath.Join(dir, "does-not-exist"), dangling); err != nil {
		t.Skipf("symlink unsupported on this fs: %v", err)
	}

	_, stats, err := ScanDir(context.Background(), dir, ScanOpts{NoFollowSymlinks: true})
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	for _, e := range stats.Errors {
		// The dangling link must NOT appear as an error — it was
		// filtered out before the read attempt.
		if filepath.Base(e.Error()) == "nowhere.txt" {
			t.Fatalf("dangling symlink surfaced as error: %v", e)
		}
	}
}
