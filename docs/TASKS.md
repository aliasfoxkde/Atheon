# Task Ledger — Atheon Enhanced

**Last Updated**: 2026-06-27
**Status**: Active — Wave 11 complete, release v1.3.1-enhanced deployed (2026-06-27)

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

- [x] Add CWE `relationships` to SARIF rules (secrets→CWE-798, web-security→CWE-79/601, etc.)
- [x] Add `severity`, `column`, `fingerprint`, `category` fields to JSON output
- [x] Scope `dummy-function`, `mock-stub`, `fake-data` patterns to test files only
- [x] Scope `sleep-in-test` to `*_test.go` / `test_*.py` files
- [x] Fix `skip-tests` overly-broad regex (anchor `skip` word boundary, remove `mvn.*skip.*test`)
- [x] Lower `todo-comment`/`fixme-comment` severity to `info`
- [x] Remove broken `helpUri` from SARIF rules (wiki/patterns#<name> does not exist)
- [x] Add `cmd/atheon/main_json_output_test.go`
- [x] Add `cmd/atheon/main_sarif_relationships_test.go`

> **Note**: All PR #100 items landed in PR #102 (Wave 10) which supersedes PR #100.

### PR #101: Bundle hash verification + rate-limiter hardening + binary sniff

- [x] Add SHA-256 verification for downloaded bundles (fetch checksums.txt first)
- [~] Publish `checksums.txt` alongside GitHub releases (bundler computes at release time) — release process step, tracked separately
> **Note**: `release.yml` generates checksums.txt at release time via goreleaser; `SetBundleDownloadURL` now enforces HTTPS-only URLs (SSRF guard); hash mismatch is fatal.
- [x] Add concurrent request cap to MCP server (atomic.Int counter, maxConcurrent=50)
- [x] Add extension-based binary heuristic for large `.log`/`.cfg`/`.conf`/`.ini` files
- [x] Add UTF-16 BOM detection to binary sniff (`\xff\xfe` / `\xfe\xff`)
- [x] Add `core/bundle_hash_test.go`
- [x] Add `core/binary_sniff_test.go`
- [x] Add `cmd/mcp/mcp_concurrency_test.go`

### PR #102: Help text + Go 1.25 prep + yaml.v3 deprecation

- [x] Document `--all` and `--no-follow-symlinks` in `--help`
- [x] Add Go 1.25 to CI matrix
- [x] Add `yaml.v3` deprecation comment to `go.mod`
- [x] Run `go vet ./...` and `golangci-lint` with latest version to catch new lints

---

## Wave 10: Post-Wave 9 Hardening (PR #102)

> PR #102 supersedes PR #100. All items below are implemented in PR #102 (11 commits, 26 files, +1105/-123).

### MCP path traversal fix (`cmd/mcp/main.go`)
- [x] Add `sandboxPath(path)` helper: `filepath.Clean` + `EvalSymlinks` on relative paths before dispatch
- [x] Block `../../etc/passwd` and relative symlink escapes (`cwd/subdir -> /etc`)
- [x] Absolute paths pass through unchanged (explicit user intent)
- [x] `handleScanFile` and `handleScanDir` call `sandboxPath` before dispatch
- [x] `cmd/mcp/mcp_sandbox_test.go`: 5 test cases

### Bundle download hardening (`core/bundle.go`)
- [x] `io.LimitedReader` caps bundle downloads at `maxBundleDownloadBytes` (100 MiB)
- [x] `Content-Length` header vs actual-bytes validation in `fetchBundleData`
- [x] `SetBundleDownloadURL` rejects non-HTTP(S) schemes (`file://`, `ftp://`, etc.) — SSRF prevention
- [x] `verifyBundleHash` failure now propagates as error (was: warn-and-proceed)

### TOCTOU fix (`core/runner.go`)
- [x] `readFileCapped` calls `filepath.EvalSymlinks` before `os.Stat`
- [x] Symlinks to huge files sized by resolved target, not symlink itself

### JSON-RPC error `data` field (`cmd/mcp/main.go`)
- [x] `rpcError` struct gains `Data any` field per JSON-RPC 2.0 spec
- [x] `rate_limit` → `Data: "rate_limit"`, `concurrent_limit` → `Data: "concurrent_limit"`, `invalid_params` → `Data: "invalid_params"`

### CI/Release fixes
- [x] `release.yml`: add `-race` to pre-release test gate
- [x] `release.yml`: pin goreleaser-action `version: '7.2.2'`
- [x] `release.yml`: add `--prov` flag for SLSA provenance attestation
- [x] `community-pattern-review.yml`: add `--max-time 30` to curl call

### Test infrastructure (required by the fatal hash check)
- [x] `bundle_download_test.go`: `serveBundle` and all mock servers serve `checksums.txt`
- [x] `state_io_errors_test.go`: mock servers updated
- [x] `cli_test.go`, `main_run_test.go`: mock servers updated, trailing-slash URL fix
- [x] `bundle_hash_test.go`: `TestVerifyBundleHashMismatch` updated for fatal error expectation

---

## Deferred / Backlog

- [x] ~~Migrate `gopkg.in/yaml.v3` to `github.com/goccy/go-yaml`~~ — DONE PR #111
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
| 9 | [x] | #99-100 | MCP protocol, SARIF ecosystem, bundle integrity |
| 10 | [x] | #102 | Post-Wave 9 hardening: SSRF, TOCTOU, fatal hash, path sandbox |
| 11 | [x] | #109-111 | Test fix, CHANGELOG fix, yaml.v3 → goccy/go-yaml |
| 12 | [x] | #113-125 | SDLC: commitlint, stale cleanup, PR labeler, goreleaser fixes, release v1.3.1 |
| 13 | [x] | #129 | Comprehensive audit: CI/CD fixes (Go 1.21 EOL, goreleaser pin, ascend-again, MODEL env, jitter), labeler v6.1.0 fix, lint fix |

**Completed waves**: 13 / 13
**Total merged PRs**: 44 through Wave 13
