package main

import (
	"strings"
	"testing"
)

// BenchmarkRedact measures the redact() cost. redact is called once per
// finding's Content field during JSON output and is on the per-finding
// hot path for `--json` mode. A regression here shows up as slower CI
// scans over repos with many findings.
//
// Run with: go test -bench=BenchmarkRedact -benchmem ./cmd/atheon/
func BenchmarkRedact(b *testing.B) {
	cases := []string{
		"AKIAIOSFODNN7EXAMPLE",
		"ghp_1234567890abcdefghij",
		"sk-1234567890abcdefghij",
		"1234567890",
		"short",
		"a very long string that simulates a multi-line token captured in a finding",
	}
	// Build a large body so we get a meaningful sample size.
	var sb strings.Builder
	for i := 0; i < 1000; i++ {
		sb.WriteString(cases[i%len(cases)])
		sb.WriteByte('\n')
	}
	body := sb.String()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, c := range cases {
			_ = redact(c)
		}
		_ = redact(body)
	}
}
