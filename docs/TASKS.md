# Task List — Atheon-Enhanced

**Version**: 0.6.0 (post-Wave 6)
**Last Updated**: 2026-06-25
**Format**: One section per hardening wave. Each wave ends with a merged PR cluster; new waves begin when a gap-analysis subagent surfaces enough new work to justify a PR.

---

## Status Legend

- `[x]` Merged
- `[~]` In Progress
- `[ ]` Open
- `[!]` Blocked
- `[-]` Cancelled / Deferred

---

## Wave 1: Initial scaffold + rename cleanup

PRs: #74, #75, #76, #79, #80

- [x] Move repo from `agravels` namespace to `aliasfoxkde` org
- [x] Add LICENSE (MIT), CODEOWNERS, dependabot config
- [x] Initial README, AGENTS.md, CLAUDE.md
- [x] Baseline CI: go test, gofmt, go vet
- [x] First community/*.yaml pattern submission flow

---

## Wave 2: CI/security plumbing

PR: #81

- [x] Dependabot groups (GitHub Actions grouped, Go minor/patch split)
- [x] govulncheck pinned to v1.1.4 (workaround for setup-go v6 toolchain)
- [x] PR template (`.github/PULL_REQUEST_TEMPLATE.md`)
- [x] Release environment (`github.event.repository.default_branch`)
- [x] `GO_VERSION` repo variable with `1.23` default
- [x] Pin actions by SHA + version comment convention

---

## Wave 3: Fuzz + coverage + SBOM

PR: #83

- [x] Fuzz harness for `core/loadBundle` (gzip round-trip)
- [x] Fuzz harness for `core.ScanString` (RE2 no-panic guarantee)
- [x] Detection-coverage test against known-bad fixture corpus
- [x] `-trimpath` for reproducible builds
- [x] SPDX SBOM generation in release workflow
- [x] Codecov integration with `CODECOV_TOKEN`

---

## Wave 4: Severity wiring

PR: #84

- [x] Add `severity` field to `PatternDef` (low/medium/high/critical)
- [x] Wire through `core.Pattern` → `core.Finding` → SARIF
- [x] CVSS-like score mapping (9.5/7.5/5.0/2.5)
- [x] SARIF `level` enum (error/warning/note)
- [x] Bundle regenerated; 274 patterns now carry severity
- [x] SARIF rule descriptors for code-scanning UI

---

## Wave 5: Closing Wave 4 audit findings

PR: #85

- [x] Surface `Stats.Errors` to stderr + bump exit code (silent data-loss fix)
- [x] Validate `--category=` against `core.Categories()` (typo guard)
- [x] Hot-path benchmarks uploaded as CI artifact
- [x] Remove dead `scripts/` (no-op maintenance)
- [x] Add false-positive guard test corpus

---

## Wave 6: Audit-followup hardening

PRs: #86, #87, #88

- [x] PR #86: `slog.Info` in legacy-default-flip branch (`core/bundle.go`)
- [x] PR #86: reorder `--version` after `--json`/`--sarif` strip
- [x] PR #86: `TestVersionFlagWithJSON` subtests for flag orderings
- [x] PR #87: CI JSON-RPC roundtrip replaces smoke-step on `atheon-mcp`
- [x] PR #88: fill `docs/PLAN.md` and `docs/TASKS.md` (this file)
- [x] PR #88: add `docs/RELEASE.md` maintainer runbook
- [x] PR #88: consolidate `docs/BRANCH_STRATEGY.md`; delete `docs/reports/BRANCH_STRATEGY.md`

---

## Deferred / Open

These came out of the Wave 5 subagent gap report and are intentionally not in Wave 6. The defer was an explicit user choice.

- [ ] **Item 2 — `pattern_state` race fix.** Add `sync.Mutex` around `core/pattern_state.go` writes; switch bundle temp-file write to tempfile+rename. Needs `-race` validation and a focused concurrent test harness. Best done as a single PR with its own review. **Candidate for Wave 7.**
- [ ] **Item 3+** — minor items (TODO comments in non-core, lint cleanup). Fold into the next wave or a standalone refactor PR.
- [ ] **Bundle format `version: 2` field.** Should land before the next breaking wire change. Tracked separately; not blocking.
- [ ] **`dev/full-feature` branch.** Mentioned in CI `docs/BRANCH_STRATEGY.md` check but does not exist. Decide whether to create it or remove the reference.
- [ ] **Doc tasks #171, #172, #173.** Resolved by PR #88 (the `{{...}}` placeholders and the duplicate `BRANCH_STRATEGY.md` are now gone).

---

## Bug Tracker (active)

None open as of 2026-06-25.

Historical (all closed in their respective waves):
- Wave 4: severity dropped at SARIF boundary → closed in PR #84.
- Wave 5: scan silently dropped permission-denied files → closed in PR #85.
- Wave 6: legacy bundle flip indistinguishable from intentional all-disabled → closed in PR #86.
- Wave 6: `atheon --json --version` errored → closed in PR #86.
- Wave 6: MCP smoke test missed framing/parsing regressions → closed in PR #87.

---

## Notes

- **Subagent gap analysis** is the project's standing input channel. After each merged wave, an Explore agent enumerates new gaps with file/line citations. The output becomes the next wave's input.
- **Each wave is one or more PRs** (typically 2-3) for reviewability. Never land >5 PRs in a single wave without an explicit reason.
- **Coverage threshold is a guardrail, not a goal.** Coverage has held ≥70% across all shipped waves. If a deliberate drop is needed (e.g. to add a new package), adjust `vars.COVERAGE_THRESHOLD` in repo settings and call it out in the PR.
- **CI gates that may surprise a contributor**: the `gofmt -l .` check, the `gofmt -d .` (prints diff on failure), the no-TODO / no-debug grep, the `goimports -l .` check, and the bundle-freshness check (`SOURCE_COUNT == BUNDLE_COUNT`). The bundle check is the one most likely to bite a first-time contributor — the fix is `go run ./bundler && git add core/patterns.bundle && git commit --amend`.

---

## Progress Summary

| Wave | Status | PRs | Theme |
|------|--------|-----|-------|
| 1 | [x] | 4 (#74-76, #79-80) | Scaffold + rename |
| 2 | [x] | 1 (#81) | CI/security plumbing |
| 3 | [x] | 1 (#83) | Fuzz + SBOM |
| 4 | [x] | 1 (#84) | Severity wiring |
| 5 | [x] | 1 (#85) | Audit findings |
| 6 | [x] | 3 (#86-88) | Audit-followup + docs |
| 7 | [ ] | TBD | pattern_state mutex (deferred) |

**Completed waves**: 6 / 6 in flight.
**Total merged PRs** (through Wave 6): 11.
**Open work**: 1 (Item 2 mutex) + 1 (bundle v2 field) + 1 (`dev/full-feature` decision).
