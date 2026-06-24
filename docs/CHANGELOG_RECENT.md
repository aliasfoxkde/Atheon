# Changelog — Recent Version

This file tracks the most recent released version plus the in-progress unreleased section.
For the full history, see [../CHANGELOG.md](../CHANGELOG.md).

---

## [Unreleased]

### Added
- **Planning docs** filled with real Atheon-Enhanced content (`docs/PLAN.md`, `docs/TASKS.md`,
  `docs/PROGRESS.md`)
- **`docs/ANALYSIS_REPORT.md`** — deep audit cataloguing 30+ gaps across docs, CI, source code,
  patterns, MCP, and ops
- **`AGENTS.md`** at repo root — guidance for AI agents working on this codebase
- **`docs/CHANGELOG_RECENT.md`** — this file
- PII patterns: national-id, dob-format, gender-field, health-record-id, tax-id-ein
- Secrets patterns: cloudflare-token, okta-api-token, pagerduty-api-key, heroku-api-key,
  travis-ci-token, circleci-token, sonarqube-token, artifactory-token, firebase-api-key,
  vercel-token
- Cloud-native patterns: aws-arn, gcp-project-id, azure-connection-string,
  k8s-imagepullsecret, helm-secret-value
- Code-quality patterns: sleep-in-test, fmt-println-prod, panic-in-handler, direct-sql-query,
  global-variable, unused-import-comment
- New `compliance` category: gdpr-personal-data-comment, hipaa-phi-field,
  pci-cardholder-data, data-retention-comment
- New `git-hygiene` category: merge-conflict-marker, fixup-commit-message,
  rebase-todo-leftover, git-rerere-conflict

### Changed
- Structured logging via `log/slog` for consistency and flexibility
- `ValidatePattern()` helper in core for reusable pattern validation
- Bundle wiring: `patterns.bundle` is now 252 patterns / 18 non-empty categories

### Fixed
- `Finding.Line` guard ensures 1-indexed line numbers (0 becomes 1)
- `gofmt -l` check used in CI where applicable
- SHA-pinned all GitHub Actions (PR #43)

### Security
- govulncheck runs in CI (PR #43)
- JUnit test reporting in CI (PR #43)
- Atheon self-scan integrated into CI (PR #49)

---

## [0.4.0] - 2026-06-22

### Added
- 223+ patterns across 19 categories (rolled forward to 252 since)
- `atheon update` command for downloading latest pattern bundle
- `atheon list --enabled/--disabled/--category=` filtering
- `atheon --json` JSON output mode
- `atheon --env` for scanning environment variables
- `atheon --stdin` for scanning piped content
- MCP server (`atheon-mcp`) for IDE integration
- SARIF output (`--sarif`) for GitHub Security tab integration
- Rate limiter in MCP server (10 req/sec, burst 20)
- Codecov v5 integration (PR #53)

### Changed
- Bundle format: gzip-compressed JSON for smaller size and faster loading
- Pattern enable/disable persists across runs via `~/.atheon/pattern_state.json`
- All scan entry points now accept `context.Context`

### Security
- SHA-pinned GitHub Actions
- govulncheck in CI
- JUnit test reporting in CI
- Atheon self-scan gating secrets/PII in production source

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

For full history, see [../CHANGELOG.md](../CHANGELOG.md).
