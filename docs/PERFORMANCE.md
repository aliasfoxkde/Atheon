# Performance

How fast Atheon is, what it costs, and where the time goes.
Pair this with [`BENCHMARKS.md`](BENCHMARKS.md) for the raw
benchmark numbers and [`ARCHITECTURE.md`](ARCHITECTURE.md)
for the structure that produces them.

## Headline numbers

Measured on an AWS `c7a.large` (2 vCPU, 4 GiB RAM), scanning
the Linux kernel tree (`linux/` at v6.7, ~73k files, ~1.1 GB
of source) with the default pattern bundle:

| Metric | Value |
|---|---|
| Wall time | 4.8 s |
| Throughput | 15.2k files/s |
| Peak RSS | 220 MB |
| CPU utilisation | 195 % (≈2 cores saturated) |

The same scan on an Apple M2 (4 perf cores) finishes in
**3.1 s**. The single-core regression is real and expected;
the scanner is a parallel workload.

## Where the time goes

A profile of the kernel tree scan, in descending order of
flat CPU time:

| Function | Flat % | Cumulative % |
|---|---|---|
| `(*Engine).Match` | 38 % | 38 % |
| `readFile` (kernel) | 22 % | 60 % |
| `(*Engine).scanLine` | 14 % | 74 % |
| `filepath.WalkDir` | 9 % | 83 % |
| `(*Pattern).Compile` (cached) | 6 % | 89 % |
| reporting / JSON marshal | 4 % | 93 % |
| everything else | 7 % | 100 % |

Three observations:

1. **Matching dominates.** The regex engine is the
   bottleneck. Optimisations to matching (e.g. Aho-Corasick
   for the literal-prefix stage) have the highest
   leverage.
2. **`readFile` is mostly syscalls.** On Linux, the kernel
   page cache absorbs most of this; on cold storage the
   number doubles. Atheon does not buffer ahead — it scans
   one file at a time as `WalkDir` yields it. Parallel I/O
   is a known follow-up (see `ROADMAP.md`).
3. **Pattern compilation is cheap.** Patterns are compiled
   once at startup and reused for every file. This is why
   the per-file cost is dominated by matching, not by
   compilation.

## Memory

The scanner's memory profile is dominated by:

- **Pattern bundle** — the compiled regexes, kept in a
  single read-only `[]Pattern`. For the default bundle:
  ~28 MB resident.
- **Per-file line slice** — the current file, fully
  materialised as `[]string` by `bufio.Scanner`. For the
  99th-percentile file in the kernel tree this is 12 KB;
  for pathological minified JS it is ~50 MB. Atheon caps
  the line slice at 1 MB; files with longer lines are
  skipped with a `lines-too-long` finding.
- **Finding slice** — one `Finding` per match. A typical
  scan produces < 1000 findings; the slice stays under
  1 MB.

Peak RSS for the kernel scan is 220 MB. For 95 % of
real-world scans it is under 150 MB. The 220 MB figure
includes Go's heap arenas, the regex cache, and the page
cache, all of which are paid once per process.

## Concurrency model

Atheon fans out a fixed number of workers over the file
list. The default is `runtime.NumCPU()`. The model is:

```
WalkDir ──► job channel ──► [worker] ──► result channel ──► reporter
              N workers, N = runtime.NumCPU()
```

- **No shared mutable state between workers.** Each worker
  has its own `bufio.Scanner` and `[]byte` scratch buffer.
- **The reporter is single-threaded.** Findings are appended
  to a single `[]Finding` protected by a mutex. The cost of
  the mutex is negligible because findings are rare
  (≪ 1 per file).
- **Backpressure is implicit.** `WalkDir` blocks until a
  worker reads from the job channel, which caps memory at
  `N * (file size)`. We rely on `WalkDir`'s pacing to
  prevent I/O queue saturation.

The `-j N` flag overrides the worker count. For SSDs,
`N = runtime.NumCPU()` is optimal. For network mounts,
lower `N` (1–2) avoids thrashing.

## Cold vs warm

The numbers above are **warm**: the file system is in the
page cache, the binary is already running, the bundle is
parsed. The **cold** cost on a freshly-booted machine is
dominated by:

1. Process start (~25 ms).
2. Bundle parse + regex compile (~180 ms for the default
   bundle; scales linearly with pattern count).
3. First `readFile` per file (~1–3 ms each on a cold
   page cache).

For a CI job that boots a fresh container and scans a fresh
checkout, the cold overhead is ~5–10 % of the wall time.
This is why Atheon's Docker image keeps the bundle pre-parsed
in `/usr/local/share/atheon/bundle.bin`.

## What we optimise for

Atheon is a developer tool that runs in editors, pre-commit
hooks, and CI. The two constraints that matter most:

1. **Editor responsiveness.** A single file scan must
   complete in < 50 ms so the user does not see a stutter.
   This is the only "hard" latency requirement. The
   per-file path is optimised for it: no allocations on
   the happy path, no locks, no logger.
2. **CI throughput.** A repo scan should saturate the
   CPU. We measure against `runtime.NumCPU()` and report
   the per-core throughput in `BENCHMARKS.md`.

We **do not** optimise for:

- **Streaming mode.** Atheon does not currently stream
  findings as they are found; the reporter batches them.
  A future "watch" mode is on `ROADMAP.md`.
- **Sub-millisecond scans.** The minimum useful unit is
  ~10 ms; below that the editor/IDE cannot show feedback
  anyway. We do not chase sub-ms.

## Profiling

To reproduce the profile numbers above:

```bash
# CPU profile
go test -bench=BenchmarkScanDir -benchtime=10x \
  -cpuprofile=cpu.out ./core
go tool pprof -top -cum cpu.out

# Memory profile
go test -bench=BenchmarkScanDir -benchtime=10x \
  -memprofile=mem.out ./core
go tool pprof -top mem.out

# Allocation profile
go test -bench=BenchmarkScanDir -benchtime=10x \
  -benchmem ./core
```

For flame graphs:

```bash
go tool pprof -http=:8080 cpu.out
```

## Known regressions

These are accepted performance costs; if you find a way to
recover them, file an issue.

- **`.atheonignore` parsing cost.** A large
  `.atheonignore` is parsed once per scan. For repos with
  > 5k ignore lines this is ~15 ms. The fix is a compiled
  ignore cache; tracked on `ROADMAP.md`.
- **Windows file I/O.** `os.ReadFile` on Windows is
  ~30 % slower than `mmap + read` on Linux. The Go runtime
  is the bottleneck; nothing to do in user space.
- **macOS quarantine.** On macOS, the first `exec` of a
  downloaded binary triggers a Gatekeeper check that adds
  ~200 ms. This is an OS-level cost, not a scanner cost.

## When to file a performance issue

Open an issue if:

- A scan that used to finish in < N seconds now takes
  > 2N seconds on the same hardware.
- A single file scan takes > 100 ms on a representative
  file.
- RSS exceeds 500 MB on a default scan.

Please include: command, machine, file count, time, RSS,
and the version of Atheon. A `go test -bench` output is
helpful but not required.
