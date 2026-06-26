# ADR 0006: Parallel Test Requirement (-p 1)

**Status**: Accepted
**Date**: 2026-06-17

## Context

Go's `go test` runs package tests in parallel by default. We discovered this causes test failures in Atheon due to global state in the `core` package. Options considered:

- **Remove global state**: Significant refactoring
- **Use mutexes/locks**: Performance impact, complex
- **Run tests serially (-p 1)**: Simple workaround
- **Use test suffixes**: Fragile, easy to forget

## Decision

Require `-p 1` flag for all test runs.

## Root Cause

The `core` package has global variables:
- `allPatterns` - slice of all loaded patterns
- `activeScanners` - compiled combined regex per category
- `activeCategoryFilter` - current category filter

When `init()` runs, it loads the bundle and sets up state. In parallel test execution:
1. Package A's init loads bundle with patterns P1,P2,P3
2. Package B's init loads bundle with patterns P1,P2,P3,P4
3. One test disables pattern P1
4. Another test expects P1 to exist
5. Race condition causes flaky tests or failures

## Rationale

1. **Simple**: Just add `-p 1` to test commands
2. **Effective**: Serializes access to global state
3. **No code changes**: Doesn't require refactoring core
4. **CI consistency**: Same behavior locally and in CI

## Consequences

### Positive
- Tests are reliable and repeatable
- No race conditions in bundle loading
- Simple to understand and document

### Negative
- Slightly slower test execution (but reliability > speed)
- Easy to forget if not in CI config

## Enforcement

All CI workflows must use `-p 1`:
```yaml
- run: go test ./... -p 1 -timeout 15m
```

Makefile targets include `-p 1`:
```makefile
test:
    go test ./... -p 1 -timeout 15m -coverprofile=coverage.out
```

---

**References**:
- [Go test flags](https://pkg.go.dev/cmd/go#hdr-Test_flags)
- [Go concurrency bugs in tests](https://github.com/golang/go/wiki/TestFlags)