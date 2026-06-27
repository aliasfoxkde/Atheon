# Wave 13 — Comprehensive Audit & Gap Closure

**Created**: 2026-06-27
**Status**: In Progress

---

## Executive Summary

Comprehensive audit identified **28 issues** across CI/CD (9), hooks/scripts (7), code quality (5), and documentation (7). This wave targets all HIGH priority items and MEDIUM items that are quick wins.

---

## CI/CD Issues

### HIGH Priority

| # | Issue | File | Fix |
|---|-------|------|-----|
| 1 | Go 1.21 EOL but in test matrix | ci.yml | Remove 1.21 |
| 2 | GoReleaser pinned to `latest` | release.yml | Pin to specific version |
| 3 | Hardcoded gpt-4o-mini model | community-pattern-review.yml | Add MODEL env var |
| 4 | Silent API failure fallback | community-pattern-review.yml | Emit warning, fail on API errors |
| 5 | Build blocked by all 5 test variants | ci.yml | Remove `needs: test` from build |
| 6 | No CI failure notifications | ci.yml | Add notification step |
| 7 | Duplicate security-events:write | security.yml | Remove job-level permissions |

### MEDIUM Priority

| # | Issue | File | Fix |
|---|-------|------|-----|
| 8 | go-junit-report uses @latest | ci.yml | Pin version |
| 9 | ascend-again: true can re-stale | stale.yml | Set to false |
| 10 | 60-day stale may be too aggressive | stale.yml | Consider 90 days |
| 11 | Auto-merge polling no jitter | auto-merge.yml | Add RANDOM jitter |
| 12 | report job if:always() fragile | ci.yml | Check file existence |

---

## Code Quality Issues

### HIGH Priority

| # | Issue | File | Fix |
|---|-------|------|-----|
| 1 | Semaphore leak in ScanDir on ctx cancel | core/runner.go | Restructure defer |

### MEDIUM Priority

| # | Issue | File | Fix |
|---|-------|------|-----|
| 2 | Bundle freshness hardcoded (24h) | core/bundle.go | Add env var |
| 3 | HTTP timeouts hardcoded | core/bundle.go | Add env vars |
| 4 | Worker count hardcoded | core/runner.go | Add env var |

---

## Documentation Issues

### HIGH Priority

| # | Issue | File | Fix |
|---|-------|------|-----|
| 1 | ARCHITECTURE.md refs yaml.v3 (deprecated) | docs/ARCHITECTURE.md | Update to goccy/go-yaml |

### MEDIUM Priority

| # | Issue | File | Fix |
|---|-------|------|-----|
| 2 | PATTERN_CATEGORIES.md says 152 (stale) | docs/ | Update to 272 |
| 3 | External wiki links may be orphaned | .github/wiki/ | Update or remove |

---

## Hooks/Scripts Issues

### MEDIUM Priority

| # | Issue | File | Fix |
|---|-------|------|-----|
| 1 | goimports error swallowed with `\|\| true` | scripts/hooks/pre-commit | Remove `\|\| true` |
| 2 | mktemp -d no failure check | scripts/hooks/pre-commit | Add check |
| 3 | install-hooks uses @latest | scripts/install-hooks.sh | Pin versions |

---

## Execution Plan

### PR 1: CI/CD Fixes (part 1)
- Remove Go 1.21 from matrix
- Remove `needs: test` from build
- Add notification step
- Pin go-junit-report version

### PR 2: CI/CD Fixes (part 2) + Community Review Fixes
- Pin GoReleaser version
- Fix community-pattern-review.yml (MODEL env, warning on API fail)
- Remove duplicate permissions
- Add jitter to auto-merge

### PR 3: Code Quality Fixes
- Fix semaphore leak in runner.go
- Add ATHEON_BUNDLE_FRESHNESS_HOURS env var
- Add ATHEON_HTTP_TIMEOUT env vars
- Add ATHEON_MAX_WORKERS env var

### PR 4: Documentation Fixes
- Update ARCHITECTURE.md yaml.v3 → goccy/go-yaml
- Update PATTERN_CATEGORIES.md count to 272
- Review/update .github/wiki/ links

### PR 5: Hooks/Scripts Fixes
- Fix goimports error handling
- Fix mktemp check
- Pin tool versions in install-hooks

---

## Verification

All PRs must pass:
```bash
go vet ./... && go build ./... && go test ./... -q
cargo fmt --check  # if Rust code changed
```
