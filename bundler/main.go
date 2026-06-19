package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type patternFile struct {
	Name    string `yaml:"name"`
	Match   string `yaml:"match"`
	Enabled *bool  `yaml:"enabled,omitempty"`
}

type patternDef struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Match    string `json:"match"`
	Enabled  bool   `json:"enabled"`
}

// bundleWalkErr is a sentinel for WalkDir errors so the loop above still
// returns errors but tests can introspect.
type bundleWalkErr struct {
	path string
	err  error
}

func (e *bundleWalkErr) Error() string { return e.path + ": " + e.err.Error() }
func (e *bundleWalkErr) Unwrap() error { return e.err }

// bundleToWriter bundles the community directory and writes the result to out.
// Splitting this out from bundle() lets tests pass a failing writer.
func bundleToWriter(communityDir string, out io.Writer) (int, error) {
	defs, err := walkPatterns(communityDir)
	if err != nil {
		return 0, err
	}

	jsonBytes, err := json.Marshal(defs)
	if err != nil {
		return 0, err
	}

	if err := writeGzipped(out, jsonBytes); err != nil {
		return 0, err
	}
	return len(defs), nil
}

func bundle(communityDir, outPath string) (int, error) {
	var buf bytes.Buffer
	n, err := bundleToWriter(communityDir, &buf)
	if err != nil {
		return 0, err
	}
	if err := os.WriteFile(outPath, buf.Bytes(), 0o644); err != nil {
		return 0, err
	}
	return n, nil
}

// walkPatterns walks communityDir and returns all parsed pattern definitions.
func walkPatterns(communityDir string) ([]patternDef, error) {
	var defs []patternDef
	err := filepath.WalkDir(communityDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".yaml") {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return &bundleWalkErr{path: path, err: err}
		}
		var pf patternFile
		if err := yaml.Unmarshal(data, &pf); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if pf.Name == "" || pf.Match == "" {
			return fmt.Errorf("%s: missing name or match", path)
		}
		category := filepath.Base(filepath.Dir(path))
		enabled := true
		if pf.Enabled != nil {
			enabled = *pf.Enabled
		}
		defs = append(defs, patternDef{
			Name:     pf.Name,
			Category: category,
			Match:    pf.Match,
			Enabled:  enabled,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return defs, nil
}

// writeGzipped gzip-encodes jsonBytes into out. Accepts io.Writer so tests
// can substitute a failing writer to exercise the gzip error paths.
func writeGzipped(out io.Writer, jsonBytes []byte) error {
	gz := gzip.NewWriter(out)
	if _, err := gz.Write(jsonBytes); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	return nil
}

func main() {
	os.Exit(run(os.Args[1:]))
}

// run executes the bundler with the given args and returns the exit code.
// This is separated from main() so tests can call it without os.Exit
// terminating the test process.
func run(args []string) int {
	communityDir := "community"
	outPath := filepath.Join("core", "patterns.bundle")
	if len(args) > 0 {
		communityDir = args[0]
	}
	if len(args) > 1 {
		outPath = args[1]
	}

	n, err := bundle(communityDir, outPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	fmt.Printf("bundled %d patterns → %s\n", n, outPath)
	return 0
}
