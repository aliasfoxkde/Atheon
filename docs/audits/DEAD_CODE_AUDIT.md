# Dead-Code Audit — Atheon-Enhanced

**Date:** 2026-06-20
**Scope:** `core/`, `cmd/atheon/`, `cmd/mcp/`, `bundler/`
**Tools:** `staticcheck 2025.1.1 (0.6.1)`, `go vet`, manual grep
**Auditor:** Phase 1.1 of [`docs/GOAL_ROADMAP.md`](../GOAL_ROADMAP.md)

---

## Summary

| Tool | Findings |
|---|---|
| `staticcheck ./...` | 0 issues |
| `go vet ./...` | 0 issues |
| Manual grep (unexported funcs with no callers) | 1 (see below) |
| `grep "//nolint"` | 3 (all justified, see below) |
| `grep "// FIXME\|// XXX"` | 0 |

The codebase is in **excellent shape**. One dead helper remains.

---

## Findings

### Finding 1: `contains(slice, item)` in `core/bundle.go` is dead code

**File:** `core/bundle.go:495-503`

```go
// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

**Evidence:**
- Zero callers in production code (`grep "contains(" core/ cmd/ bundler/` returns only the definition).
- One test caller: `core/internal_helpers_test.go:31` (`TestContains`).
- Upstream has the same dead helper (HoraDomu/Atheon `core/bundle.go:380`).
- **Upstream issue: #159** "Remove unused contains helper in core/bundle.go" (2026-06-20, labeled `good first issue`).
- Upstream PR #160 already opens a fix; **the fork should match**.

**Severity:** Low (no runtime impact, no security impact).

**Recommendation:** Remove the helper and its test. Replace any future callers with `slices.Contains` (stdlib, Go 1.21+).

---

## Justified `//nolint` / `//nolint:gosec` comments

| File:Line | Justification |
|---|---|
| `core/runner.go:98` | `return nil //nolint:nilerr // skip unreadable entries during walk; reported via stats`. The skip-on-walk-error path is intentional; the error is propagated via `Stats.WalkErrors` (or silently skipped for non-fatal symlink/permission races during a parallel directory walk). Acceptable. |
| `cmd/atheon/main_test.go:13` | `if err := os.WriteFile(file, []byte("token=sk-abcdefghijklmnopqrstuvwxyz\n"), 0o644); err != nil { //nolint atheon:ignore`. Test fixture writes a fake secret; the `atheon:ignore` comment is the documented escape hatch. Acceptable. |
| `core/runner_test.go:185` | `t.Errorf("isIgnored(%q, nil) = %v, want %v", tt.path, result, tt.expected) //nolint`. Standard `t.Errorf` with no `t.Fatal`; the test continues to other cases after reporting the failure. Comment is informational, not a suppression. |

**Recommendation:** Keep all three.

---

## Already-removed dead code (historical)

These items were removed in prior commits; documented here so future
auditors know not to "re-add" them:

| Removed in commit | What was removed | Why |
|---|---|---|
| `72fa65c` (fix: remove unused helper functions) | `isPatternInRegistry`, `removePatternFromRegistry` from `core/bundle.go` | Unused after pattern state refactor |
| `593b973` (fix(ci): repair workflow build paths, remove unused code) | Various unused helpers across the codebase | CI lint failures |

---

## Sentinel errors — verified all used

| Sentinel | Used by | Status |
|---|---|---|
| `ErrPatternNotFound` | `core/bundle.go` (EnablePattern/DisablePattern) | ✓ |
| `ErrBundleDownload` | `core/bundle.go` (DownloadBundle error wrapping) | ✓ |
| `ErrBundleParse` | `core/bundle.go` (loadBundle error wrapping) | ✓ |
| `ErrInvalidPattern` | exported for external test packages | ✓ (intentional public API) |

---

## Recommendations

1. **Remove `contains` and its test** (Phase 1.2 of GOAL_ROADMAP).
2. **No other dead code in production.** Existing `//nolint` comments are all justified.
3. **Add a `make audit` target** that runs staticcheck + go vet + this manual grep (Phase 1.3 of GOAL_ROADMAP).
4. **Wire staticcheck into `.githooks/pre-commit`** so new dead code is caught at commit time (Phase 1.4 of GOAL_ROADMAP).
5. **Add a PreToolUse hook in the global harness** that runs `go build ./...` after every Write/Edit on Go files to catch dead-code-introducing changes immediately (Phase 1.5 of GOAL_ROADMAP).

---

## Cross-reference with upstream

- Upstream #155 (duplicate `!p.enabled` check): **already fixed** in the fork in commit `2c45e7a`.
- Upstream #159 (`contains` helper): **NOT yet fixed** in the fork; matches upstream exactly. Will be fixed in Phase 1.2.
- Upstream #160 (PR for #159): exists on the upstream tracker; the fork fix will be its own PR.
