# Task List — Atheon-Enhanced

**Project**: aliasfoxkde/Atheon-Enhanced
**Last Updated**: 2026-06-23
**Detail**: [ANALYSIS_REPORT.md](ANALYSIS_REPORT.md) is the canonical source of truth for
all phases below.

---

## Task Status Legend

- [ ] Pending
- [~] In Progress
- [x] Completed
- [!] Blocked
- [-] Cancelled

---

## Phase A — Trust restoration [IN PROGRESS]

> Make every docs number match reality. Fix broken links. Fill the planning docs.

### A.1 Planning docs (this file's siblings)

- [x] A.1.1 Fill `docs/PLAN.md` with real Atheon-Enhanced content
- [x] A.1.2 Fill `docs/TASKS.md` (this file)
- [x] A.1.3 Fill `docs/PROGRESS.md`
- [x] A.1.4 Remove the duplicate template copies in `docs/planning/`

### A.2 Pattern-count consistency

- [ ] A.2.1 Update `README.md` to actual 252 patterns / 18 categories
- [ ] A.2.2 Update `docs/FAQ.md` (currently says 225 / 57)
- [ ] A.2.3 Update `docs/INSTALL.md` (currently says 190)
- [ ] A.2.4 Update `docs/development/SETUP.md` (currently says "Expected: 87")
- [ ] A.2.5 Update `docs/architecture/PATTERN_CATEGORIES.md` (currently says 225)
- [ ] A.2.6 Update README's pattern-distribution table (currently totals ~177; actual 252)

### A.3 Broken doc links

- [ ] A.3.1 Stub or remove references in `docs/README.md` to:
  - `docs/reports/MCP_INTEGRATION_ANALYSIS.md` (does not exist)
  - `docs/reports/SECURITY_TESTING.md` (does not exist)
  - `docs/reports/FINALIZATION_SUMMARY.md` (does not exist)
  - `docs/reports/TEST_COVERAGE.md` (does not exist)
  - `docs/tests/TEST_HARNESS_STATUS.md` (does not exist)
  - `docs/contributors.md` (does not exist)

### A.4 Doc hygiene

- [ ] A.4.1 Remove `--debug-level=2` example from `docs/development/SETUP.md` (no such flag)
- [ ] A.4.2 Pick one canonical `BRANCH_STRATEGY.md`; remove or redirect the other
- [ ] A.4.3 Update `docs/architecture/SYSTEM_ARCHITECTURE.md` "Code Organization" section to
      match actual files (no `streaming.go`, no `quality_enforcement.go`, no `config/defaults/`)

### A.5 Root-level contract docs (per user CLAUDE.md)

- [x] A.5.1 Create `CHANGELOG_RECENT.md` at `docs/`
- [x] A.5.2 Create `AGENTS.md` at repo root

---

## Phase B — CI consolidation [PENDING]

> Cut CI minutes, eliminate duplication across the 10 workflows.

- [ ] B.1 Consolidate 10 workflows → 4 (`ci.yml`, `security.yml`, `release.yml`, `sync.yml`)
- [ ] B.2 Drop `auto-merge.yml` if GitHub native auto-merge is enabled
- [ ] B.3 Fix `gofmt` check in `quality-assurance.yml` to use `gofmt -l`
- [ ] B.4 Verify all `go test` invocations have `-p 1`
- [ ] B.5 Drop `continue-on-error: true` from Codecov upload
- [ ] B.6 Fix `scheduled-release.yml` tag format (`0.4.YYMMDD` is technically valid but odd)

---

## Phase C — MCP completeness [PENDING]

> Make MCP server match the API docs.

- [ ] C.1 Add `list_patterns` MCP tool
- [ ] C.2 Add `list_categories` MCP tool
- [ ] C.3 Add `scan_env` MCP tool
- [ ] C.4 Add `update_bundle` MCP tool
- [ ] C.5 Inject version via ldflag into `serverInfo`
- [ ] C.6 Change rate-limit error code from -32600 to -32000

---

## Phase D — Pattern expansion [PENDING]

> Reach 300+ patterns, fix `frameworks/` empty category.

### D.1 Top-tier SaaS secrets (10 patterns)

- [ ] D.1.1 `anthropic-api-key` — `sk-ant-*`
- [ ] D.1.2 `openai-project-key` — `sk-proj-*`
- [ ] D.1.3 `github-fine-grained-pat` — `github_pat_*`
- [ ] D.1.4 `supabase-service-key` — `eyJ...` JWT form
- [ ] D.1.5 `vercel-edge-config-token`
- [ ] D.1.6 `cloudflare-r2-token`
- [ ] D.1.7 `slack-workflow-token` — `xoxw-*`
- [ ] D.1.8 `bitbucket-app-password` — `ATBB*`
- [ ] D.1.9 `hashicorp-vault-token` — `hvs.*` / `hvb.*`
- [ ] D.1.10 `onepassword-service-account`

### D.2 PII patterns (5)

- [ ] D.2.1 `email-address`
- [ ] D.2.2 `ipv4-literal`
- [ ] D.2.3 `ipv6-literal`
- [ ] D.2.4 `passport-number` (international variants)
- [ ] D.2.5 `ssn-strong` (with format validation)

### D.3 Frameworks restore

- [ ] D.3.1 Either restore 3 patterns (`django`, `nodejs`, `react`) OR delete `community/frameworks/`

---

## Phase E — Architecture hygiene [PENDING]

- [ ] E.1 Create `docs/architecture/decisions/` directory
- [ ] E.2 Write `ADR-001-re2-regex.md`
- [ ] E.3 Write `ADR-002-gzip-bundle.md`
- [ ] E.4 Write `ADR-003-parallel-tests-must-be-1.md`
- [ ] E.5 Either create `core/streaming.go` (extract from runner.go) or fix SYSTEM_ARCHITECTURE.md

---

## Phase F — Future features [PENDING]

- [ ] F.1 Pattern metadata (severity, description, references) — wire-format change
- [ ] F.2 `--baseline` filter for incremental scans
- [ ] F.3 SBOM generation in release workflow
- [ ] F.4 LSP mode for IDE integration (VS Code / Cursor)
- [ ] F.5 `--lsp-mode` JSON-RPC bridge

---

## Bug Fixes & Improvements (open)

### Bugs (engine layer — none currently known)
- [ ] (no open bugs reported in the audited code)

### Documentation bugs (see Phase A)
- [x] Unfilled template placeholders in `docs/PLAN.md`, `docs/TASKS.md`, `docs/PROGRESS.md`
- [ ] Stale pattern counts in 5+ files
- [ ] 6 broken doc links in `docs/README.md`

### Improvements (housekeeping)
- [ ] Add `scripts/pattern-count.sh` to generate a single source of truth for counts
- [ ] Wire the `.github/wiki/` markdown files into a Wiki publish workflow (currently 3 files
      sitting unused in `.github/wiki/`)
- [ ] Add `docs/RELEASE.md` runbook for cutting a release

---

## Notes

This task list supersedes the prior template copy in `docs/planning/TASKS.md`. That file should
be removed (tracked as A.1.4).

The IMPROVEMENT_PLAN.md at `docs/reports/IMPROVEMENT_PLAN.md` was a partial-execution plan from
PR #45. This `TASKS.md` is the live task list going forward.

---

## Progress Summary

| Phase | Tasks | Completed | In Progress | Pending |
|-------|-------|-----------|-------------|---------|
| Phase A (docs trust) | 15 | 5 | 0 | 10 |
| Phase B (CI consolidation) | 6 | 0 | 0 | 6 |
| Phase C (MCP completeness) | 6 | 0 | 0 | 6 |
| Phase D (patterns) | 18 | 0 | 0 | 18 |
| Phase E (architecture) | 5 | 0 | 0 | 5 |
| Phase F (future features) | 5 | 0 | 0 | 5 |
| **Total** | **55** | **5** | **0** | **50** |

**Completion**: 9% (5/55) as of 2026-06-23.
