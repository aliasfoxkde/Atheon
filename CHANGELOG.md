# Changelog - Atheon Enhanced

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `docs/RELEASE.md` (maintainer runbook): tag format (`v0.YY.MM.DD[-rcN]`),
  pre-release checklist, GoReleaser publishing, hotfix workflow, ldflags,
  bundle regeneration, troubleshooting.
- `docs/PLAN.md` filled with project-specific content reflecting the
  multi-wave hardening cycle (274 patterns, 19 categories, MCP server,
  wave-by-wave hardening). Replaces the prior `{{...}}` template.
- `docs/TASKS.md` filled with the actual task ledger (waves 1–6 marked
  completed, deferred items including the `pattern_state` mutex work).
- `docs/BRANCH_STRATEGY.md` consolidated with the richer Quick Start,
  Decision Tree, Branch Comparison table, and per-branch configuration
  sections from the previous `docs/reports/BRANCH_STRATEGY.md`.

### Removed
- `docs/reports/BRANCH_STRATEGY.md` (duplicate of the canonical
  `docs/BRANCH_STRATEGY.md`). Its unique content was merged into the
  canonical copy and the duplicate was deleted to satisfy the project's
  "no duplicate implementations" rule.

## [0.6.0] - 2026-06-25

### Added
- Pattern severity wired end-to-end: each `community/*.yaml` declares
  `severity: low|medium|high|critical`. Defaults applied by category
  (secrets/pii/web-security/compliance/security-hardening = high;
  code-quality/performance/accessibility/web-development = low; all others = medium).
  Severity flows through `Pattern → bundlePattern → Finding → SARIF`.
- SARIF severity mapping: CVSS-like scores (9.5/7.5/5.0/2.5) and levels
  (error/warning/note) derived per-pattern, plus a `security-severity-label`
  property for human readability. Replaces hard-coded
  `security-severity: "High"` / `level: "error"` for every result.
- `Pattern.Severity()` interface method and `normalizeSeverity()` helper that
  coerces empty/typo'd values to `medium` so downstream code stays safe.
- Hot-path benchmarks: `BenchmarkLoadBundle` (gzip decode + 274-pattern
  compile), `BenchmarkCompileIgnoreFile` and `BenchmarkIgnoreMatcherMatch`
  (recursive regex on `.atheonignore`), and `BenchmarkRedact` (per-finding
  redact on the `--json` output path). Run with `go test -bench=. ./core/`
  and `./cmd/atheon/`.
- `scanErrorsPresent()` helper: bumps exit code when a scan silently dropped
  files (permission denied, unreadable). Closes a data-loss gap where a
  partial failure was reported as success.
- Pre-commit hook now surfaces bundler warnings (per-file skip reasons) on
  the same line as the success message, so contributors see when a pattern
  was dropped without a manual re-run.
- `atheon list --category=<bogus>` now errors with the known-category list,
  instead of silently filtering to zero matches.
- `slog.Info` line emitted when `loadBundle` flips all patterns to enabled
  (legacy compatibility path). Without this log, a contributor bundle
  that's accidentally all-`enabled: false` looks identical at runtime —
  every pattern silently comes on. Surface the path so it stays observable.
- Regression tests `TestLoadBundleLegacyDefaultFlip` and
  `TestLoadBundleNoFlipWhenAnyEnabled` in `core/bundle_legacy_default_test.go`
  guard the legacy-flip behavior and the log-on/only-on gate.
- `TestVersionFlagWithJSON` subtests for `["--version"]`,
  `["--json", "--version"]`, and `["--sarif", "--version"]` flag orderings.
- CI JSON-RPC roundtrip integration test in `.github/workflows/ci.yml`:
  replaces the prior smoke step (which only verified clean exit on empty
  stdin) with a real `initialize` + `tools/list` roundtrip and `jq`
  assertions on `protocolVersion`, `capabilities`, tool count, and tool
  names. Catches framing and discovery regressions that the smoke test
  missed.

### Changed
- Bundler (`go run ./bundler`) no longer aborts on broken pattern files.
  Malformed YAML, missing fields, whitespace in pattern names, duplicate
  names, and invalid regex are logged to stderr and the file is skipped.
  This mirrors `loadBundle`'s runtime tolerance.
- 67 community patterns had pre-existing regex corruption (severity text
  embedded in the `match:` value); these were repaired so all 274 patterns
  ship cleanly.
- `atheon --json --version` (and `--sarif --version`) now print the version
  cleanly. Previously the `--version` check ran before the `--json`/`--sarif`
  strip, so `atheon --json --version` fell into the default branch and
  errored with `path not found: --version`. Flag order is now forgiving.

### Fixed
- `community-pattern-review` workflow SIGPIPE: `git diff ... | head -10`
  exited 141 under `set -euo pipefail` when `head` closed the pipe early.
  Disabled pipefail around that pipeline only — captured files unchanged.
- `gofmt` alignment in `cmd/atheon/main.go` (the SARIF map literal had a
  misaligned key after the severity wiring change).

### Removed
- Dead scripts: `scripts/doc-validate.sh`, `scripts/doc-exemptions.sh`,
  `scripts/doc-categorize.sh`. No callers existed anywhere in the repo;
  deleting them removes a code-maintenance violation.

## [0.5.0] - 2026-06-25

### Added
- New patterns in PII category: national-id, dob-format, gender-field, health-record-id, tax-id-ein
- New patterns in Secrets category: cloudflare-token, okta-api-token, pagerduty-api-key, heroku-api-key, travis-ci-token, circleci-token, sonarqube-token, artifactory-token, firebase-api-key, vercel-token
- New patterns in Cloud-native category: aws-arn, gcp-project-id, azure-connection-string, k8s-imagepullsecret, helm-secret-value
- New patterns in Code-quality category: sleep-in-test, fmt-println-prod, panic-in-handler, direct-sql-query, global-variable, unused-import-comment
- New `compliance` category: gdpr-personal-data-comment, hipaa-phi-field, pci-cardholder-data, data-retention-comment
- New `git-hygiene` category: merge-conflict-marker, fixup-commit-message, rebase-todo-leftover, git-rerere-conflict
- `scripts/pattern-count.sh` — single source of truth for pattern counts (replaces
  hardcoded numbers scattered across docs). Supports `--json`, `--total`, `--table`,
  `--help`. Confirmed: **274 patterns / 19 categories**.
- `docs/architecture/decisions/` directory for Architecture Decision Records (ADRs)
  (planned)

### Changed
- Structured logging via `log/slog` for consistency and flexibility
- `ValidatePattern()` helper in core for reusable pattern validation
- **CI consolidation**: 10 GitHub Actions workflows → 5 (ci, security, release,
  sync, auto-merge). Removed duplicate test/lint/build, self-scan, and CodeQL
  workflows. Consolidated into a coherent set with single-responsibility jobs.
- `gofmt` check now uses the standard `gofmt -l .` idiom (replaced a non-standard
  `--debug-level` pattern from docs)
- Codecov upload no longer has `continue-on-error` (was silently masking failures)
- Scheduled release tag format changed from `0.4.YYMMDD` to `v0.YY.MM.DD` for
  consistency with manual tags
- All documented `go test ./...` invocations now include `-p 1` (the flag is
  mandatory because `core/` has package-level state in `init()` that breaks under
  parallel package execution). Updated: `.pre-commit-config.yaml`,
  `scripts/coverage.sh`, `.github/wiki/TROUBLESHOOTING.md`,
  `docs/guides/TROUBLESHOOTING.md`, `docs/PATTERN_FORMAT.md`,
  `docs/reports/REPOSITORY_RENAME_PLAN.md`, `docs/reports/BRANCH_STRATEGY.md`,
  `docs/reports/FEATURE_COMPARISON.md`, `docs/architecture/SYSTEM_ARCHITECTURE.md`,
  `docs/self-scan.md`. README CI badge updated to `ci.yml`.
- `docs/reports/BRANCH_STRATEGY.md` is now a redirect stub pointing to canonical
  `docs/BRANCH_STRATEGY.md`
- `docs/architecture/SYSTEM_ARCHITECTURE.md` "Code Organization" section rewritten
  to reflect actual file inventory (removed phantom `core/streaming.go`,
  `core/quality_enforcement.go`, `config/defaults/` entries)
- Pattern counts in 6 doc files updated to reflect actual 274 patterns / 19
  categories (was stale: 177 / 225 / 255 / 190 depending on file)

### Fixed
- Finding.Line guard ensures 1-indexed line numbers (0 becomes 1)
- README CI badge pointed at non-existent `comprehensive-ci.yml` — now points at
  the consolidated `ci.yml`

## [0.4.0] - 2026-06-22

### Added
- 223+ patterns across 19 categories
- `atheon update` command for downloading latest pattern bundle
- `atheon list --enabled/--disabled/--category=` filtering
- `atheon --json` JSON output mode
- `atheon --env` for scanning environment variables
- `atheon --stdin` for scanning piped content
- MCP server (`atheon-mcp`) for IDE integration

### Changed
- Bundle format: gzip-compressed JSON for smaller size and faster loading
- Pattern enable/disable persists across runs via `~/.atheon/pattern_state.json`

### Security
- SHA-pinned GitHub Actions
- govulncheck in CI
- JUnit test reporting in CI

---

## [0.3.0] - 2026-06-20

### Added
- `.atheonignore` file support
- Context cancellation support for all scan operations
- `atehon:ignore` inline directive

### Changed
- Improved performance with combined regex per category

---

## [0.2.0] - 2026-06-17

### Added
- Initial release with core pattern categories
- Secrets and PII detection
- Multiple output formats (text, JSON)

---

[Unreleased]: https://github.com/aliasfoxkde/Atheon-Enhanced/compare/v0.5.0...HEAD
[0.5.0]: https://github.com/aliasfoxkde/Atheon-Enhanced/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/aliasfoxkde/Atheon-Enhanced/releases/tag/v0.4.0
[0.3.0]: https://github.com/aliasfoxkde/Atheon-Enhanced/releases/tag/v0.3.0
[0.2.0]: https://github.com/aliasfoxkde/Atheon-Enhanced/releases/tag/v0.2.0