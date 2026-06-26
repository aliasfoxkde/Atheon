package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func setupCommunity(t *testing.T, files map[string]string) (communityDir string) {
	t.Helper()
	tmp := t.TempDir()
	communityDir = filepath.Join(tmp, "community")
	if err := os.MkdirAll(communityDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for rel, content := range files {
		path := filepath.Join(communityDir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return communityDir
}

func readBundle(t *testing.T, path string) []patternDef {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read bundle: %v", err)
	}
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("gzip open: %v", err)
	}
	defer r.Close()
	var defs []patternDef
	if err := json.NewDecoder(r).Decode(&defs); err != nil {
		t.Fatalf("json decode: %v", err)
	}
	return defs
}

func TestBundleBasic(t *testing.T) {
	community := setupCommunity(t, map[string]string{
		"secrets/api-key.yaml": "name: test-api-key\nmatch: '\\bTEST_[A-Z0-9]{32}\\b'\n",
	})
	out := filepath.Join(t.TempDir(), "out.bundle")

	n, err := bundle(community, out)
	if err != nil {
		t.Fatalf("bundle: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 pattern, got %d", n)
	}

	defs := readBundle(t, out)
	if len(defs) != 1 {
		t.Fatalf("expected 1 def in bundle, got %d", len(defs))
	}
	d := defs[0]
	if d.Name != "test-api-key" {
		t.Errorf("name = %q, want test-api-key", d.Name)
	}
	if d.Category != "secrets" {
		t.Errorf("category = %q, want secrets", d.Category)
	}
	if !d.Enabled {
		t.Error("enabled should default to true")
	}
}

func TestBundleEnabledFalse(t *testing.T) {
	community := setupCommunity(t, map[string]string{
		"secrets/disabled.yaml": "name: disabled-key\nmatch: '\\bDIS_[A-Z0-9]{32}\\b'\nenabled: false\n",
	})
	out := filepath.Join(t.TempDir(), "out.bundle")

	n, err := bundle(community, out)
	if err != nil {
		t.Fatalf("bundle: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 pattern, got %d", n)
	}

	defs := readBundle(t, out)
	if defs[0].Enabled {
		t.Error("enabled should be false")
	}
}

func TestBundleMultipleCategories(t *testing.T) {
	community := setupCommunity(t, map[string]string{
		"secrets/key.yaml": "name: secret-key\nmatch: 'sk_[a-z0-9]{32}'\n",
		"tokens/jwt.yaml":  "name: jwt-token\nmatch: 'eyJ[A-Za-z0-9_-]+\\.eyJ'\n",
	})
	out := filepath.Join(t.TempDir(), "out.bundle")

	n, err := bundle(community, out)
	if err != nil {
		t.Fatalf("bundle: %v", err)
	}
	if n != 2 {
		t.Fatalf("expected 2 patterns, got %d", n)
	}
}

// TestBundleInvalidYAMLSkips — invalid YAML is logged + skipped, not an error.
func TestBundleInvalidYAMLSkips(t *testing.T) {
	community := setupCommunity(t, map[string]string{
		"secrets/bad.yaml":  "name: bad\nmatch: [unclosed\n",
		"secrets/good.yaml": "name: good\nmatch: 'x'\n",
	})
	out := filepath.Join(t.TempDir(), "out.bundle")

	n, err := bundle(community, out)
	if err != nil {
		t.Fatalf("bundle should not error on bad YAML, got: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 pattern (skipping the bad one), got %d", n)
	}
}

// TestBundleMissingMatchSkips: walkPatterns now skips files with missing
// fields (logging a warning) rather than aborting the build. This matches
// loadBundle's runtime tolerance. The pattern is silently dropped from the
// output bundle; n must reflect that.
func TestBundleMissingMatchSkips(t *testing.T) {
	community := setupCommunity(t, map[string]string{
		"secrets/incomplete.yaml": "name: incomplete-key\n",
		"secrets/good.yaml":       "name: good-key\nmatch: 'x'\n",
	})
	out := filepath.Join(t.TempDir(), "out.bundle")

	n, err := bundle(community, out)
	if err != nil {
		t.Fatalf("bundle: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 pattern (skipping the broken one), got %d", n)
	}
}

// TestBundleMissingNameSkips — same skip policy as above.
func TestBundleMissingNameSkips(t *testing.T) {
	community := setupCommunity(t, map[string]string{
		"secrets/noname.yaml": "match: 'something'\n",
		"secrets/good.yaml":   "name: good-key\nmatch: 'x'\n",
	})
	out := filepath.Join(t.TempDir(), "out.bundle")

	n, err := bundle(community, out)
	if err != nil {
		t.Fatalf("bundle: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 pattern (skipping the broken one), got %d", n)
	}
}

func TestBundleEmptyDirectory(t *testing.T) {
	community := setupCommunity(t, map[string]string{})
	out := filepath.Join(t.TempDir(), "out.bundle")

	n, err := bundle(community, out)
	if err != nil {
		t.Fatalf("bundle: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 patterns, got %d", n)
	}
}

func TestBundleBadOutputPath(t *testing.T) {
	community := setupCommunity(t, map[string]string{
		"secrets/key.yaml": "name: k\nmatch: 'x'\n",
	})
	_, err := bundle(community, "/nonexistent/dir/out.bundle")
	if err == nil {
		t.Error("expected error for bad output path, got nil")
	}
}

// TestBundleWhitespaceNameSkips — the broken-name file is dropped, the
// good file is kept. walkPatterns logs the skip reason to stderr.
func TestBundleWhitespaceNameSkips(t *testing.T) {
	community := setupCommunity(t, map[string]string{
		"secrets/whitespace.yaml": "name: 'key with space'\nmatch: 'x'\n",
		"secrets/good.yaml":       "name: good-key\nmatch: 'x'\n",
	})
	out := filepath.Join(t.TempDir(), "out.bundle")

	n, err := bundle(community, out)
	if err != nil {
		t.Fatalf("bundle: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 pattern (skipping the broken-name one), got %d", n)
	}
}

// TestBundleDuplicatePatternNameSkips — first writer wins, second is dropped.
func TestBundleDuplicatePatternNameSkips(t *testing.T) {
	community := setupCommunity(t, map[string]string{
		"secrets/key1.yaml": "name: duplicate-key\nmatch: 'test1'\n",
		"secrets/key2.yaml": "name: duplicate-key\nmatch: 'test2'\n",
	})
	out := filepath.Join(t.TempDir(), "out.bundle")

	n, err := bundle(community, out)
	if err != nil {
		t.Fatalf("bundle: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 pattern (first wins, second skipped), got %d", n)
	}
}

// TestBundleInvalidRegexSkips — invalid regex is logged + skipped.
func TestBundleInvalidRegexSkips(t *testing.T) {
	community := setupCommunity(t, map[string]string{
		"secrets/bad.yaml":  "name: bad-regex\nmatch: '[invalid'\n",
		"secrets/good.yaml": "name: good-key\nmatch: 'x'\n",
	})
	out := filepath.Join(t.TempDir(), "out.bundle")

	n, err := bundle(community, out)
	if err != nil {
		t.Fatalf("bundle: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 pattern (skipping the invalid-regex one), got %d", n)
	}
}

// TestBundleToWriterGzipFailure exercises the writeGzipped error path in bundleToWriter.
type failWriter struct{}

func (failWriter) Write([]byte) (int, error) { return 0, errWriteFail }

var errWriteFail = errors.New("write failed")

func TestBundleToWriterGzipFailure(t *testing.T) {
	community := setupCommunity(t, map[string]string{
		"secrets/key.yaml": "name: test-key\nmatch: 'AKIAIOSFODNN7EXAMPLE'\n",
	})

	_, err := bundleToWriter(community, failWriter{})
	if err == nil {
		t.Error("expected error from failing writer, got nil")
	} else if !errors.Is(err, errWriteFail) {
		t.Errorf("expected errWriteFail sentinel, got: %v", err)
	}
}

// TestWalkPatternsWhitespaceName — broken-name file is skipped, not an error.
func TestWalkPatternsWhitespaceName(t *testing.T) {
	community := setupCommunity(t, map[string]string{
		"secrets/bad.yaml":  "name: 'bad name'\nmatch: 'foo'\n",
		"secrets/good.yaml": "name: good\nmatch: 'foo'\n",
	})

	defs, err := walkPatterns(community)
	if err != nil {
		t.Fatalf("walkPatterns should not error on whitespace name, got: %v", err)
	}
	if len(defs) != 1 {
		t.Errorf("expected 1 pattern (skipping the bad-name one), got %d", len(defs))
	}
}

// TestWalkPatternsMissingFields — missing field is skipped, not an error.
func TestWalkPatternsMissingFields(t *testing.T) {
	community := setupCommunity(t, map[string]string{
		"secrets/empty.yaml": "name: ''\nmatch: 'foo'\n",
		"secrets/good.yaml":  "name: good\nmatch: 'foo'\n",
	})

	defs, err := walkPatterns(community)
	if err != nil {
		t.Fatalf("walkPatterns should not error on missing field, got: %v", err)
	}
	if len(defs) != 1 {
		t.Errorf("expected 1 pattern (skipping the empty-name one), got %d", len(defs))
	}
}
