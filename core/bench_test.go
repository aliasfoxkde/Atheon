package core

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// BenchmarkScanString measures the throughput of the in-memory line
// scanner against a synthetic 100KB file of mostly-non-matching lines
// with periodic matches. This is the hot path for the MCP server.
//
// Run with: go test -bench=BenchmarkScanString -benchmem ./core/
func BenchmarkScanString(b *testing.B) {
	// 100KB of mostly-non-matching content with one match per KB.
	var sb strings.Builder
	for i := 0; i < 100; i++ {
		sb.WriteString("this is line one of a thousand fictional secret-less lines\n")
		sb.WriteString("and another line\n")
		sb.WriteString("AKIAIOSFODNN7EXAMPLE\n") // one match per 30 lines
	}
	content := sb.String()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ScanString(context.Background(), content, "bench.txt")
	}
}

// BenchmarkScanFile measures the file-read-plus-scan path against a
// fixed 1MB temp file. Useful for detecting regressions in the read
// pipeline.
func BenchmarkScanFile(b *testing.B) {
	tmp := filepath.Join(b.TempDir(), "bench.txt")
	var sb strings.Builder
	for i := 0; i < 10000; i++ {
		sb.WriteString("the quick brown fox jumps over the lazy dog\n")
	}
	if err := os.WriteFile(tmp, []byte(sb.String()), 0o600); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _, err := ScanFile(context.Background(), tmp)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkScanDir exercises the parallel walk-and-scan path against a
// synthetic directory of 200 small files. Wall-clock time scales with
// GOMAXPROCS.
func BenchmarkScanDir(b *testing.B) {
	root := b.TempDir()
	for i := 0; i < 200; i++ {
		p := filepath.Join(root, "f_"+itoa(i)+".go")
		if err := os.WriteFile(p, []byte("package x\nvar apiKey = \"AKIAIOSFODNN7EXAMPLE\"\n"), 0o600); err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _, err := ScanDir(context.Background(), root)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMatchPattern measures the inner per-line, per-pattern cost.
// The combined regex pre-filter narrows the candidate set so most lines
// only hit one regex match per category.
func BenchmarkMatchPattern(b *testing.B) {
	// Pre-resolve the combined regex from the first active scanner.
	if len(activeScanners) == 0 {
		b.Skip("no active scanners; run with the embedded bundle")
	}
	cs := activeScanners[0]
	line := "this line contains AKIAIOSFODNN7EXAMPLE somewhere"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = cs.combined.MatchString(line)
	}
}

// BenchmarkScanStringEmpty measures the all-non-matching path — useful
// for separating the cost of line iteration from regex matching.
func BenchmarkScanStringEmpty(b *testing.B) {
	var sb strings.Builder
	for i := 0; i < 1000; i++ {
		sb.WriteString("nothing here, no matches, no nothing\n")
	}
	content := sb.String()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ScanString(context.Background(), content, "bench.txt")
	}
}

// itoa is a tiny strconv.Itoa shim that keeps this file import-free.
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	neg := i < 0
	if neg {
		i = -i
	}
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
