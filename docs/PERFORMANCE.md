# Performance Improvements: Issues #146, #147, #148

## Issue #146: Profile Scanning on Large Repositories

Add profiling support to identify bottlenecks during scanning.

**Suggested Implementation:**
- Add flag: `--profile <file>`
- Output pprof-compatible profile data

## Issue #147: Stream Findings Instead of Buffering

Instead of buffering all results in memory, stream findings as they are found.

**Suggested Implementation:**
```go
func ScanStreaming(root string, onFinding func(Finding)) error
// or
func ScanChannel(root string) <-chan Finding
```

## Issue #148: Chunk Large Files

Currently files are loaded entirely into memory with `os.ReadFile()`. For large files this causes memory issues.

**Suggested Implementation:**
- Use bufio.Reader with chunked reading for files > 10MB
- Process file in chunks instead of loading entirely

---

**Date:** 2026-06-22
**Branch:** `pi/146-147-148-perf-v2`
**Base:** upstream/main
