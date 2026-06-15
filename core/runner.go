package core

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var skipDirs = map[string]bool{
	".git": true, "node_modules": true, "vendor": true,
	".terraform": true, "dist": true, "build": true, "__pycache__": true,
}

var binaryExts = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
	".pdf": true, ".zip": true, ".tar": true, ".gz": true,
	".exe": true, ".bin": true, ".so": true, ".dylib": true,
	".class": true, ".jar": true,
}

func ScanFile(path string) ([]Finding, *Stats, error) {
	start := time.Now()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	findings := scanLines(string(data), path)
	return findings, &Stats{
		Files:     1,
		Bytes:     int64(len(data)),
		ElapsedMs: time.Since(start).Milliseconds(),
	}, nil
}

func ScanDir(root string) ([]Finding, *Stats, error) {
	start := time.Now()
	var paths []string

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if !binaryExts[ext] {
			paths = append(paths, path)
		}
		return nil
	})

	results := make([][]Finding, len(paths))
	sizes := make([]int64, len(paths))
	var wg sync.WaitGroup

	for i, p := range paths {
		wg.Add(1)
		go func(i int, p string) {
			defer wg.Done()
			data, err := os.ReadFile(p)
			if err != nil {
				return
			}
			results[i] = scanLines(string(data), p)
			sizes[i] = int64(len(data))
		}(i, p)
	}
	wg.Wait()

	var findings []Finding
	var totalBytes int64
	for i := range results {
		findings = append(findings, results[i]...)
		totalBytes += sizes[i]
	}

	return findings, &Stats{
		Files:     len(paths),
		Bytes:     totalBytes,
		ElapsedMs: time.Since(start).Milliseconds(),
	}, nil
}

func ScanEnv() []Finding {
	var findings []Finding
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, val := parts[0], parts[1]
		for _, p := range registry {
			if p.Matches(val) {
				findings = append(findings, Finding{
					Pattern: p.Name(),
					File:    "env:" + key,
					Content: val,
				})
			}
		}
	}
	return findings
}

func ScanString(content, source string) []Finding {
	return scanLines(content, source)
}

func scanLines(content, file string) []Finding {
	var findings []Finding
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		for _, p := range registry {
			if p.Matches(line) {
				findings = append(findings, Finding{
					Pattern: p.Name(),
					File:    file,
					Line:    i + 1,
					Content: strings.TrimSpace(line),
				})
			}
		}
	}
	return findings
}
