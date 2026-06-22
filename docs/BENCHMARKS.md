# Benchmarks

Atheon ships with a small set of in-package benchmarks that measure the
hot paths of the scan engine. These are intended to catch regressions
in the regex-prefilter pipeline, the parallel directory walker, and the
per-line matching cost.

For large-scale cross-corpus measurement (e.g. running Atheon against
the top 1000 GitHub repos), see the companion project
[`Atheon-Benchmark`](https://github.com/aliasfoxkde/Atheon-Benchmark).

## Running the in-tree benchmarks

```bash
# All benchmarks in the core package, one iteration each.
go test ./core/ -bench=. -benchtime=1x -run=^$

# Benchmarks with 3 seconds of measurement per case, plus memory stats.
go test ./core/ -bench=. -benchtime=3s -benchmem -run=^$

# A single benchmark (e.g. ScanDir only).
go test ./core/ -bench=BenchmarkScanDir -benchtime=5s -run=^$
```

The `-run=^$` flag skips the regular test suite so only benchmarks run.
Use `go test -bench=. ./...` to include benchmarks from `bundler`,
`cmd/atheon`, and `cmd/mcp`.

## What is benchmarked

| Function | What it measures | Why it matters |
|----------|------------------|----------------|
| `BenchmarkScanString` | in-memory scan over a 100KB synthetic file with ~3% matches | The MCP server's hot path; most user input is small strings. |
| `BenchmarkScanFile` | read + scan of a 1MB file | Disk I/O + scan combined; regressing here means the reader or scanner slowed down. |
| `BenchmarkScanDir` | parallel walk + scan of 200 small files | Validates that the worker-pool scales; use `-cpu=1,2,4,8` to compare scaling. |
| `BenchmarkMatchPattern` | single combined-regex match against a candidate line | Isolates the regex engine cost from line iteration. |
| `BenchmarkScanStringEmpty` | scan with no matches | Separates iteration cost from matching cost. |

## Interpreting results

When comparing runs, focus on **allocations per op** as well as ns/op.
A regression that adds allocs without changing ns/op is still a problem
(it pressures GC). Use `go test -benchmem` to see both.

For CPU-scaling analysis:

```bash
go test ./core/ -bench=BenchmarkScanDir -cpu=1,2,4,8 -benchtime=3s -run=^$
```

`BenchmarkScanDir` should scale nearly linearly with GOMAXPROCS up to
the number of physical cores.

## Adding new benchmarks

Add a `BenchmarkX` function to `core/bench_test.go` (or the package
whose behavior you want to measure). Follow these conventions:

- Use `b.ResetTimer()` after setup that allocates large buffers.
- Use `b.ReportAllocs()` so allocation counts are tracked.
- Use `t.TempDir()` / `b.TempDir()` for filesystem fixtures so they
  clean up automatically.
- Call the production code path with realistic input sizes.

## Reference data

A snapshot of typical numbers on an AMD Ryzen 7 5700U (8 cores, 16
threads) running Go 1.24:

```
BenchmarkScanString-16           ~20 ms   (100KB input, ~3% matches)
BenchmarkScanFile-16             ~840 ms  (1MB file, single match)
BenchmarkScanDir-16              ~10 ms   (200 small files)
BenchmarkMatchPattern-16         ~19 µs   (single combined-regex match)
BenchmarkScanStringEmpty-16      ~70 ms   (1000 non-matching lines)
```

These numbers will drift with hardware and Go version; the goal of
in-tree benchmarks is to detect regressions, not to be a substitute for
`Atheon-Benchmark`'s cross-corpus measurement.

## See also

- [Atheon-Benchmark](https://github.com/aliasfoxkde/Atheon-Benchmark) —
  measures Atheon across many public repos, tracks per-pattern precision
  and recall over time.
- [Atheon-GitHub-Scanner](https://github.com/aliasfoxkde/Atheon-GitHub-Scanner) —
  feeds Atheon the contents of every file in a target repo and reports
  aggregate findings.
- [PERFORMANCE.md](docs/architecture/PERFORMANCE.md) — notes on the
  scanner's algorithmic complexity and design choices.