# Changelog - Atheon Enhanced

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/aliasfoxkde/Atheon-Enhanced/compare/v0.4.0...HEAD
[0.4.0]: https://github.com/aliasfoxkde/Atheon-Enhanced/releases/tag/v0.4.0
[0.3.0]: https://github.com/aliasfoxkde/Atheon-Enhanced/releases/tag/v0.3.0
[0.2.0]: https://github.com/aliasfoxkde/Atheon-Enhanced/releases/tag/v0.2.0