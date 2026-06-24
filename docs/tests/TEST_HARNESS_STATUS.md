# Test Harness Status

> **Stub** — This document is referenced from `docs/README.md` but the test harness is
> standard `go test` with no custom framework.
>
> For how to run the suite, see [../development/SETUP.md §Testing Setup](../development/SETUP.md).

## Current Setup

- **Runner**: `go test ./... -p 1 -timeout 15m -coverprofile=coverage.out`
  - `-p 1` is **mandatory**: `core` has package-level state in `init()` that is not
    concurrency-safe under parallel package execution.
- **Coverage gate**: 70% minimum (CI fails below)
- **Reporting**: JUnit XML via `go-junit-report` + Codecov
- **Lint**: golangci-lint v1.64.8 (18 linters)
- **Vuln**: govulncheck
- **Pre-commit hooks**: in `scripts/hooks/`

## Test Layout

| Package | Coverage target | Notes |
|---------|----------------|-------|
| `core` | 90%+ | Sentinel errors, bundle loading, scanning |
| `cmd/atheon` | 90%+ | CLI subcommands, output formats |
| `cmd/mcp` | 90%+ | JSON-RPC, rate limiter |
| `bundler` | 90%+ | YAML → JSON+gzip |

## Status

Standard `go test` harness. No custom runner needed.
