# CI/CD Coverage Audit — 2026-06-24

End-to-end inventory of every job that runs on PR #66 (the PR that resolves
PR #55's open Coderabbit threads). Verifies the user's requirement of full
cross-platform coverage (Windows, Linux, macOS) plus the rest of the
quality/security/perf gate stack.

## Cross-platform builds (the headline requirement)

**Build job (`ci.yml`)** — runs on every PR against `main`:

| Runner            | Status on PR #66 | Notes |
|-------------------|------------------|-------|
| `ubuntu-latest`   | ✅ pass (21s)    | Native Linux build + smoke test (`./atheon --version && ./atheon list categories`) |
| `macos-latest`    | ✅ pass (13s)    | Native macOS build + smoke test |
| `windows-latest`  | ✅ pass (52s)    | Native Windows build + smoke test (`.\atheon.exe --version && .\atheon.exe list categories`) |

The matrix covers Linux/macOS/Windows as required. Each build artifact is
uploaded as `binaries-${{ matrix.os }}` so PRs can download and inspect them.

**Dev branch builds (`dev-testing.yml`)** — same matrix, relaxed gates:

| Runner            | Status                | Notes |
|-------------------|-----------------------|-------|
| `ubuntu-latest`   | ✅ pass (relaxed)     | 50% coverage threshold, non-blocking |
| `macos-latest`    | ✅ pass (relaxed)     | Same relaxed gates |
| `windows-latest`  | ✅ pass (relaxed)     | Same relaxed gates |

## Multi-version Go testing

`ci.yml`'s `test` job runs the suite against four Go versions on every PR:

| Go version | Status on PR #66 |
|------------|------------------|
| 1.21       | ✅ pass (30s)     |
| 1.22       | ✅ pass (28s)     |
| 1.23       | ✅ pass (35s)     |
| 1.24       | ✅ pass (33s)     |

The `-p 1` flag is MANDATORY (core has package-level state in `init()`),
enforced in both the workflow and the Copilot instructions.

## Complete check list — PR #66

| Check                                  | Source           | Type          | Status       |
|----------------------------------------|------------------|---------------|--------------|
| Test (Go 1.21)                         | ci.yml           | test          | ✅ pass      |
| Test (Go 1.22)                         | ci.yml           | test          | ✅ pass      |
| Test (Go 1.23)                         | ci.yml           | test          | ✅ pass      |
| Test (Go 1.24)                         | ci.yml           | test          | ✅ pass      |
| Lint (go vet + staticcheck + golangci-lint + gofmt + goimports + grep) | ci.yml | static analysis | ✅ pass |
| Build (ubuntu-latest)                  | ci.yml           | cross-platform | ✅ pass      |
| Build (macos-latest)                   | ci.yml           | cross-platform | ✅ pass      |
| Build (windows-latest)                 | ci.yml           | cross-platform | ✅ pass      |
| Integration Tests (pattern bundle ≥250, MCP startup, profile JSON validation) | ci.yml | integration | ✅ pass |
| Performance Benchmarks (3s benchtime, uploaded) | ci.yml | performance | ✅ pass |
| Documentation Check (branch strategy, profile doc) | ci.yml | docs | ✅ pass |
| Test Results & Coverage (JUnit + Codecov) | ci.yml        | reporting     | ✅ pass      |
| CodeQL (Go)                            | security.yml     | SAST          | ✅ pass      |
| Go Vulnerability Check (govulncheck)   | security.yml     | supply chain  | ✅ pass      |
| Self-Scan (secrets — blocking)         | security.yml     | dogfood       | ✅ pass      |
| Security Anti-Patterns                 | security.yml     | dogfood       | ✅ pass      |
| Self-Scan (code-quality — informational) | security.yml   | dogfood       | ✅ pass      |
| Auto-merge / Report Conflicts          | auto-merge.yml   | automation    | ✅ pass      |
| CodeRabbit                             | external         | code review   | ✅ pass      |
| Sourcery review                        | external         | code review   | ⏭️ skipping  |

**20 checks total, 19 PASS, 1 external skip.**

## Workflow inventory

| Workflow file                       | Purpose                                     | Triggers |
|-------------------------------------|---------------------------------------------|----------|
| `.github/workflows/ci.yml`          | Tests, lint, build (Win/Mac/Linux), integration, benchmarks, docs, reporting | `push:main`, `pull_request:main` |
| `.github/workflows/security.yml`    | CodeQL, govulncheck, self-scan (secrets/quality), anti-patterns | `push:main`, `schedule:weekly` |
| `.github/workflows/release.yml`     | Tag-driven GoReleaser + scheduled auto-versioned release | `push:tags:v*`, `schedule:monthly`, `workflow_dispatch` |
| `.github/workflows/sync.yml`        | Rebase `main` onto `upstream/HoraDomu:main` | `schedule:weekly`, `workflow_dispatch` |
| `.github/workflows/auto-merge.yml`  | Auto-merge Dependabot PRs once checks pass  | `pull_request:labeled:dependencies` |
| `.github/workflows/community-pattern-review.yml` | YAML schema validation + AI review (aliasfoxkde only) | `pull_request:paths:community/**` |
| `.github/workflows/dev-testing.yml` | Relaxed gates for `dev/testing` branch     | `push:dev/testing`, `pull_request:dev/testing` |

## CI/CD coverage gaps addressed in this fix (PR #66)

1. **CRIT-3 (code injection in `release.yml`)** — the consolidated
   `release.yml` inherited the `${{ inputs.version }}` injection bug from
   pre-consolidation `scheduled-release.yml`. Fixed by binding all
   `workflow_dispatch` inputs and step outputs to `env:` blocks so the
   operator-controlled value is treated as data, not shell.
2. **MAJ-1 (coverage threshold parameterization)** — `ci.yml`'s coverage
   gate now reads `vars.COVERAGE_THRESHOLD` with `'70'` as fallback, so
   raising the repo variable actually raises the gate (previously the
   threshold was hardcoded).
3. **CRIT-1/2/7 (code injection in `community-pattern-review.yml`)** —
   PR-controlled filenames are now passed via `env:` and never expanded
   into Python source or shell.

## Items deliberately not in scope

- **Windows self-scan / Windows govulncheck**: all security jobs run on
  `ubuntu-latest` for speed and consistency. GoReleaser validates Windows
  binaries via the build matrix, and govulncheck is OS-agnostic at the
  module level.
- **m1/m2 ARM runners**: macos-latest is currently Intel-on-ARM via
  Rosetta; ARM-specific runners can be added if a release regression
  appears on Apple Silicon.
- **Windows scheduled release smoke test**: smoke test for the release
  workflow is manual (`workflow_dispatch`) — the auto-merge job handles
  the gating.

## Conclusion

The CI/CD pipeline satisfies the cross-platform requirement (Linux, macOS,
Windows) at both the `main` and `dev/testing` levels. Multi-version Go
testing covers 1.21 → 1.24. Code scanning (CodeQL), supply-chain
(govulncheck), dogfood self-scans (blocking secrets + informational code
quality + anti-patterns), static analysis (vet + staticcheck +
golangci-lint + gofmt + goimports), performance benchmarks, integration
tests, documentation checks, and reporting (JUnit + Codecov) are all
present and green on PR #66.

PR #66 brings PR #55's review surface to a clean state (0 unresolved
threads) and removes the code-injection risk in the consolidated release
workflow.
