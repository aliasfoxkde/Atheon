package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
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

func bundle(communityDir, outPath string) (int, error) {
	var defs []patternDef
	err := filepath.WalkDir(communityDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".yaml") {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
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
		return 0, err
	}

	jsonBytes, err := json.Marshal(defs)
	if err != nil {
		return 0, err
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(jsonBytes); err != nil {
		return 0, err
	}
	if err := gz.Close(); err != nil {
		return 0, err
	}

	if err := os.WriteFile(outPath, buf.Bytes(), 0o644); err != nil {
		return 0, err
	}
	return len(defs), nil
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
