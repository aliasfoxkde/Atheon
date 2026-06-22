package core

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
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
}

func loadIgnorePatternsMatcher(root string) []*ignoreMatcher {
	var matchers []*ignoreMatcher
	for _, name := range []string{".atheonignore", ".gitignore"} {
		m, _ := compileIgnoreFile(filepath.Join(root, name))
		if m != nil {
			matchers = append(matchers, m)
		}
	}
	return matchers
}

func isIgnored(path string, matchers []*ignoreMatcher) bool {
	clean := filepath.ToSlash(path)
	for _, m := range matchers {
		if m.matchesPath(clean) {
			return true
		}
	}
	return false
}

// ScanFile reads a single file and reports every Finding produced by the
// currently active patterns. It honors .atheonignore and .gitignore when
// the file lives under the current working directory. The returned
// *Stats describes the read; findings are nil and stats are nil only when
// the file is ignored.
//
// The context controls read-side cancellation: if ctx is canceled before
// the read completes, ScanFile returns ctx.Err().
func ScanFile(ctx context.Context, path string) ([]Finding, *Stats, error) {
	start := time.Now()
	// Respect .atheonignore and .gitignore for files under the working directory,
	// so that `atheon file.go` and `atheon .` agree on what gets scanned.
	if absPath, err := filepath.Abs(path); err == nil {
		if root, err := os.Getwd(); err == nil {
			if rel, err := filepath.Rel(root, absPath); err == nil && !strings.HasPrefix(rel, "..") {
				if matchers := loadIgnorePatternsMatcher(root); isIgnored(filepath.ToSlash(rel), matchers) {
					return []Finding{}, nil, nil
				}
			}
		}
	}
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	findings := scanLines(ctx, string(data), path)
	return findings, &Stats{
		Files:     1,
		Bytes:     int64(len(data)),
		ElapsedMs: time.Since(start).Milliseconds(),
	}, nil
}

// ScanDir walks root and scans every non-binary, non-ignored file in
// parallel using one worker per CPU (up to a sensible cap). It honors
// .atheonignore and .gitignore at root and skips well-known noise
// directories (e.g. .git, node_modules, vendor). The returned *Stats
// counts only the files whose contents were actually scanned.
//
// The context controls worker cancellation: if ctx is canceled mid-walk
// the goroutines exit promptly and ScanDir returns ctx.Err() after
// WaitGroup drains.
func ScanDir(ctx context.Context, root string) ([]Finding, *Stats, error) {
	start := time.Now()
	ignoreMatcher := loadIgnorePatternsMatcher(root)
	var paths []string

	if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil //nolint:nilerr // skip unreadable entries during walk; reported via stats
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		if d.IsDir() {
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			// Don't SkipDir for user ignore rules — walk the dir and check files
			// individually so negation rules (e.g. !dist/keep.yaml) can un-ignore
			// specific files inside an otherwise-ignored directory.
			return nil
		}
		if isIgnored(rel, ignoreMatcher) {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if !binaryExts[ext] {
			paths = append(paths, path)
		}
		return nil
	}); err != nil {
		return nil, nil, err
	}

	results := make([][]Finding, len(paths))
	sizes := make([]int64, len(paths))
	scanned := make([]bool, len(paths))
	var wg sync.WaitGroup
	// I/O-bound file reads saturate well below CPU count; cap at 2× CPUs with
	// a minimum of 4 and a ceiling of 64 to avoid overwhelming shared runners.
	workers := min(max(runtime.NumCPU()*2, 4), 64)
	sem := make(chan struct{}, workers)

	for i, p := range paths {
		wg.Add(1)
		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			wg.Done()
			wg.Wait()
			return nil, nil, ctx.Err()
		}
		go func(i int, p string) {
			defer wg.Done()
			defer func() { <-sem }()
			if err := ctx.Err(); err != nil {
				return
			}
			data, err := os.ReadFile(p)
			if err != nil {
				return
			}
			results[i] = scanLines(ctx, string(data), p)
			sizes[i] = int64(len(data))
			scanned[i] = true
		}(i, p)
	}
	wg.Wait()

	var findings []Finding
	var totalBytes int64
	var filesScanned int
	for i := range results {
		if scanned[i] {
			filesScanned++
		}
		findings = append(findings, results[i]...)
		totalBytes += sizes[i]
	}

	return findings, &Stats{
		Files:     filesScanned,
		Bytes:     totalBytes,
		ElapsedMs: time.Since(start).Milliseconds(),
	}, nil
}

// ScanEnv scans the current process's environment variables for matches
// against the active patterns. Each finding uses "env:KEY" as its File
// and the matching value as its Content; Line is zero.
//
// The context is accepted for symmetry with the other Scan* entry
// points; the implementation checks ctx between iterations and returns
// early if canceled.
func ScanEnv(ctx context.Context) []Finding {
	return scanEnv(ctx, os.Environ())
}

// scanEnv is the inner implementation that accepts an explicit env list.
// Splitting this out lets tests exercise the len(parts) != 2 branch
// without having to mutate the real process environment.
func scanEnv(ctx context.Context, envs []string) []Finding {
	var findings []Finding
	for _, env := range envs {
		if err := ctx.Err(); err != nil {
			return findings
		}
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, val := parts[0], parts[1]
		for _, cs := range activeScanners {
			if !cs.combined.MatchString(val) {
				continue
			}
			for _, p := range cs.patterns {
				if p.Matches(val) {
					findings = append(findings, Finding{
						Pattern: p.Name(),
						File:    "env:" + key,
						Content: val,
					})
				}
			}
		}
	}
	return findings
}

// ScanString scans a string in memory and returns every Finding produced
// by the active patterns. source is recorded as the File on each
// Finding; lines are reported with their 1-indexed line number.
//
// The context is accepted for API symmetry; the scan is in-memory and
// the implementation only checks ctx for cancellation between lines.
func ScanString(ctx context.Context, content, source string) []Finding {
	return scanLines(ctx, content, source)
}

func scanLines(ctx context.Context, content, file string) []Finding {
	var findings []Finding
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if err := ctx.Err(); err != nil {
			return findings
		}
		if strings.Contains(line, "atheon:ignore") {
			continue
		}
		for _, cs := range activeScanners {
			if !cs.combined.MatchString(line) {
				continue
			}
			for _, p := range cs.patterns {
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
	}
	return findings
}
