package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// TestBundleWalkPatternsNoDir exercises the WalkDir error branch when the
// community directory doesn't exist.
func TestBundleWalkPatternsNoDir(t *testing.T) {
	_, err := walkPatterns("/this/dir/does/not/exist/anywhere")
	if err == nil {
		t.Error("expected error from walkPatterns with missing directory")
	}
}

// TestBundleWalkPatternsBadYAML exercises the yaml.Unmarshal error branch.
func TestBundleWalkPatternsBadYAML(t *testing.T) {
	tmp := t.TempDir()
	if err := writeFile(t, tmp+"/secrets/bad.yaml", "this is: not: valid: yaml: : :"); err != nil {
		t.Fatal(err)
	}
	_, err := walkPatterns(tmp)
	if err == nil {
		t.Error("expected error from walkPatterns with malformed YAML")
	}
}

// TestBundleWalkPatternsMissingFields exercises the missing name/match branch.
func TestBundleWalkPatternsMissingFields(t *testing.T) {
	tmp := t.TempDir()
	if err := writeFile(t, tmp+"/secrets/missing-name.yaml", "match: foo"); err != nil {
		t.Fatal(err)
	}
	_, err := walkPatterns(tmp)
	if err == nil {
		t.Error("expected error from walkPatterns with missing name field")
	}

	tmp2 := t.TempDir()
	if err := writeFile(t, tmp2+"/secrets/missing-match.yaml", "name: foo"); err != nil {
		t.Fatal(err)
	}
	_, err = walkPatterns(tmp2)
	if err == nil {
		t.Error("expected error from walkPatterns with missing match field")
	}
}

// TestBundleWalkPatternsEnabledFalse exercises the explicit enabled=false
// branch in walkPatterns.
func TestBundleWalkPatternsEnabledFalse(t *testing.T) {
	tmp := t.TempDir()
	if err := writeFile(t, tmp+"/secrets/x.yaml", `name: x
match: '\bX\b'
enabled: false
`); err != nil {
		t.Fatal(err)
	}
	defs, err := walkPatterns(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if len(defs) != 1 {
		t.Fatalf("expected 1 pattern, got %d", len(defs))
	}
	if defs[0].Enabled {
		t.Error("expected enabled=false to be parsed")
	}
}

// TestWriteGzippedHappy exercises the gzip write/close success path.
func TestWriteGzippedHappy(t *testing.T) {
	var buf bytes.Buffer
	data := []byte("hello, world")
	if err := writeGzipped(&buf, data); err != nil {
		t.Fatal(err)
	}
	// Verify by decompressing
	gr, err := gzip.NewReader(&buf)
	if err != nil {
		t.Fatal(err)
	}
	defer gr.Close()
	got, _ := io.ReadAll(gr)
	if string(got) != string(data) {
		t.Errorf("decompressed = %q, want %q", got, data)
	}
}

// TestWriteGzippedFailingWriter exercises the gzip write/close error branches
// by passing a writer that always errors.
func TestWriteGzippedFailingWriter(t *testing.T) {
	fw := &alwaysFailingWriter{}
	if err := writeGzipped(fw, []byte("hello")); err == nil {
		t.Error("expected error from gzip write to failing writer")
	}
}

// TestWriteGzippedCloseFailingWriter exercises the gz.Close error branch
// by passing a writer that lets the first Write succeed but fails the second.
func TestWriteGzippedCloseFailingWriter(t *testing.T) {
	fw := &closeFailingWriter{}
	if err := writeGzipped(fw, []byte("hello")); err == nil {
		t.Error("expected error from gzip close to failing writer")
	}
}

// TestBundleWithFailingGzip exercises bundle()'s writeGzipped error path
// via the bundleToWriter variant.
func TestBundleWithFailingGzip(t *testing.T) {
	tmp := t.TempDir()
	if err := writeFile(t, tmp+"/secrets/x.yaml", `name: x
match: '\bX\b'
`); err != nil {
		t.Fatal(err)
	}
	// bundleToWriter with a failing writer — covers the
	// `if err := writeGzipped(...); err != nil` branch in bundle().
	_, err := bundleToWriter(tmp, &alwaysFailingWriter{})
	if err == nil {
		t.Error("expected error from bundleToWriter with failing writer")
	}
}

// TestBundleToWriterHappy exercises bundleToWriter's success path.
func TestBundleToWriterHappy(t *testing.T) {
	tmp := t.TempDir()
	if err := writeFile(t, tmp+"/secrets/x.yaml", `name: x
match: '\bX\b'
`); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	n, err := bundleToWriter(tmp, &buf)
	if err != nil {
		t.Fatalf("bundleToWriter failed: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 pattern, got %d", n)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

// alwaysFailingWriter is an io.Writer that fails on every Write call.
type alwaysFailingWriter struct{}

var errAlwaysFail = fmt.Errorf("intentional write failure for test")

func (alwaysFailingWriter) Write(p []byte) (int, error) { return 0, errAlwaysFail }

// closeFailingWriter lets the first Write succeed (covering the gz.Write
// branch) but fails the second Write (covering the gz.Close branch).
type closeFailingWriter struct {
	count int
}

func (w *closeFailingWriter) Write(p []byte) (int, error) {
	w.count++
	if w.count >= 2 {
		return 0, errAlwaysFail
	}
	return len(p), nil
}

// TestBundleJSONMarshalError exercises the json.Marshal error branch inside
// bundle() by using a value that can't be marshaled. Since patternDef has
// only basic types, we replace writeGzipped or use an alternate path.
// Since patternDef marshals cleanly, we skip — json.Marshal errors on
// patternDef are unreachable. This test exists for documentation.
func TestBundleJSONMarshalError(t *testing.T) {
	t.Skip("patternDef has no unmarshalable fields; branch unreachable")
}

// TestBundleJsonMarshalRealBytes verifies the production JSON marshal path.
func TestBundleJsonMarshalRealBytes(t *testing.T) {
	defs := []patternDef{
		{Name: "x", Category: "secrets", Match: `\bX\b`, Enabled: true},
	}
	jb, err := json.Marshal(defs)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(jb, []byte(`"name":"x"`)) {
		t.Errorf("expected serialized name, got: %s", jb)
	}
}

// TestBundleWalkErrMethods exercises the Error and Unwrap methods on
// bundleWalkErr.
func TestBundleWalkErrMethods(t *testing.T) {
	inner := os.ErrNotExist
	werr := &bundleWalkErr{path: "/some/path", err: inner}
	if werr.Error() != "/some/path: file does not exist" {
		t.Errorf("unexpected error string: %s", werr.Error())
	}
	if werr.Unwrap() != inner {
		t.Error("Unwrap should return the inner error")
	}
}

// writeFile is a small helper to write a file and create parent dirs.
func writeFile(t *testing.T, path, content string) error {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}
