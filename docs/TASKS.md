# Task Ledger — Atheon Enhanced

**Last Updated**: 2026-06-26
**Status**: Active — Wave 9 in progress

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

## Wave 7: Concurrent pattern state

PR: #89

- [x] Add `sync.Mutex` around `core/pattern_state.go` writes
- [x] Add concurrent pattern state test harness
- [x] Document concurrent access contract

---

## Wave 8: Detection + CI + Patterns + MCP

PRs: #92-98

- [x] Add detection fixtures per category (`core/testdata/`)
- [x] Add `-race` to CI test gate
- [x] Fix float coverage comparison in CI
- [x] Remove shotgun regex patterns (circleci-token, heroku-api-key, etc.)
- [x] Add symlink guard to ScanDir
- [x] Add maxFileSize enforcement in ScanDir workers
- [x] Add NUL-byte binary content sniff
- [x] Add walk-error capture to Stats.Errors
- [x] Expand skip-dirs to include `.idea`, `.vscode`, etc.
- [x] Add SARIF uriBaseId, columns, redacted snippets, severity=none
- [x] Add panic recovery to MCP dispatchRequest
- [x] Move rate limiting to top of MCP run loop
- [x] Add MCP 64 MiB request cap and 30s timeout
- [x] Add MCP scan_string content cap and scan_env categories cap
- [x] Add JSON-RPC version validation and notification handling
- [x] Update CHANGELOG for all Wave 8 changes

---

## Wave 9: MCP Protocol + SARIF Ecosystem + Bundle Integrity

Three parallel Explore agents (2026-06-26) surfaced **62 findings** across MCP+DX, SARIF+ecosystem, and Security dimensions. Four PRs planned.

### PR #99/100: MCP error sanitization + cancel handler + stale-bundle detection

- [x] Add `$/cancelRequest` notification handler to MCP server (sync.Map per-request tracking)
- [x] Sanitize JSON-RPC error messages: map `os.IsNotExist` → `"file not found"`, `os.IsPermission` → `"permission denied"`, others → `"internal error"`
- [x] Add ETag-based stale-bundle detection with `force: bool` bypass parameter
- [x] Add `bundle.etag` and `bundle.lastChecked` to pattern state file
- [x] Add `cmd/mcp/mcp_cancel_test.go`
- [x] Add `cmd/mcp/mcp_error_sanitization_test.go`
- [x] Add `core/bundle_etag_test.go`
- [~] Progress notifications during bundle download (deferred to future PR)
- [~] Add `cmd/mcp/mcp_progress_test.go` (deferred with progress notifications)

### PR #100: SARIF rules[].relationships + output parity + community pattern triage

- [ ] Add CWE `relationships` to SARIF rules (secrets→CWE-798, web-security→CWE-79/601, etc.)
- [ ] Add `severity`, `column`, `fingerprint`, `category` fields to JSON output
- [ ] Scope `dummy-function`, `mock-stub`, `fake-data` patterns to test files only
- [ ] Scope `sleep-in-test` to `*_test.go` / `test_*.py` files
- [ ] Fix `skip-tests` overly-broad regex (anchor `skip` word boundary, remove `mvn.*skip.*test`)
- [ ] Lower `todo-comment`/`fixme-comment` severity to `info`
- [ ] Remove broken `helpUri` from SARIF rules (wiki/patterns#<name> does not exist)
- [ ] Add `cmd/atheon/main_json_output_test.go`
- [ ] Add `cmd/atheon/main_sarif_relationships_test.go`

### PR #101: Bundle hash verification + rate-limiter hardening + binary sniff

- [ ] Add SHA-256 verification for downloaded bundles (fetch checksums.txt first)
- [ ] Publish `checksums.txt` alongside GitHub releases (bundler computes at release time)
- [ ] Add concurrent request cap to MCP server (atomic.Int counter, maxConcurrent=50)
- [ ] Add extension-based binary heuristic for large `.log`/`.cfg`/`.conf`/`.ini` files
- [ ] Add UTF-16 BOM detection to binary sniff (`\xff\xfe` / `\xfe\xff`)
- [ ] Add `core/bundle_hash_test.go`
- [ ] Add `core/binary_sniff_test.go`
- [ ] Add `cmd/mcp/mcp_concurrency_test.go`

### PR #102: Help text + Go 1.25 prep + yaml.v3 deprecation

- [ ] Document `--all` and `--no-follow-symlinks` in `--help`
- [ ] Add Go 1.25 to CI matrix
- [ ] Add `yaml.v3` deprecation comment to `go.mod`
- [ ] Run `go vet ./...` and `golangci-lint` with latest version to catch new lints

---

## Deferred / Backlog

- [ ] Migrate `gopkg.in/yaml.v3` to `github.com/goccy/go-yaml` (breaking API change — requires careful review)
- [ ] Add per-tool MCP `isError` and `structuredContent` fields
- [ ] Branch protection ruleset consolidation
- [ ] Schema version for bundle format (`version: 2`)
- [ ] `update_bundle` force confirmation parameter

---

## Bug Tracker (active)

None open as of 2026-06-26.

Historical (all closed in their respective waves):
- Wave 4: severity dropped at SARIF boundary → closed in PR #84.
- Wave 5: scan silently dropped permission-denied files → closed in PR #85.
- Wave 6: legacy bundle flip indistinguishable from intentional all-disabled → closed in PR #86.
- Wave 6: `atheon --json --version` errored → closed in PR #86.
- Wave 6: MCP smoke test missed framing/parsing regressions → closed in PR #87.
- Wave 7: pattern_state race condition → closed in PR #89.
- Wave 8: shotgun regexes fired on every UUID/SHA1 → closed in PR #94.
- Wave 8: ScanDir unbounded file OOM → closed in PR #95.
- Wave 8: MCP panic killed entire server → closed in PR #97.
- Wave 8: SARIF missing uriBaseId/columns/snippets → closed in PR #96.
- Wave 8: rate limiter bypassable via initialize flood → closed in PR #97.

---

## Progress Summary

| Wave | Status | PRs | Theme |
|------|--------|-----|-------|
| 1 | [x] | #74-76, #79-80 | Scaffold + rename |
| 2 | [x] | #81 | CI/security plumbing |
| 3 | [x] | #83 | Fuzz + SBOM |
| 4 | [x] | #84 | Severity wiring |
| 5 | [x] | #85 | Audit findings |
| 6 | [x] | #86-88 | Audit-followup + docs |
| 7 | [x] | #89 | pattern_state mutex |
| 8 | [x] | #92-98 | Detection, CI, patterns, MCP hardening |
| 9 | [~] | #99-102 (planned) | MCP protocol, SARIF ecosystem, bundle integrity |

**Completed waves**: 8 / 8
**Total merged PRs**: 27 through Wave 8
**Wave 9**: In progress — 4 PRs planned
