# Project Plan — Atheon-Enhanced

**Project**: aliasfoxkde/Atheon-Enhanced
**Module**: `github.com/aliasfoxkde/Atheon`
**Go version**: 1.21+ (CI matrix: 1.21, 1.22, 1.23, 1.24)
**Last Updated**: 2026-06-23
**Status**: APPROVED — see [ANALYSIS_REPORT.md](ANALYSIS_REPORT.md) for full audit
**Detail**: [ANALYSIS_REPORT.md](ANALYSIS_REPORT.md) is the canonical source for this plan

---

## Pre-Planning Research

### Research Questions

1. **What problem are we solving?**
   Detection of secrets, PII, code-quality issues, AI-generated anti-patterns, and modern-dev
   best-practice violations in source trees — without sending code to a remote service.

2. **Who are we solving it for?**
   - Go developers who want a CLI/MCP/library in their own toolchain
   - AI assistants (Claude, Cursor, Windsurf) that need a tool to call via MCP
   - CI/CD pipelines running pre-commit, pre-merge, and pre-deploy gates
   - Security teams needing offline / on-prem scanning with a rich pattern library

3. **What similar solutions exist?**
   - **gitleaks** — secrets only, Go, popular
   - **trufflehog** — secrets with active verification
   - **detect-secrets** (Yelp) — Python, baseline-based
   - **semgrep** — broader (security + quality), Python
   - **HoraDomu/Atheon** (upstream) — 57 patterns, stable, smaller scope
   - **gitleaks/trufflehog** — secrets-only

4. **What makes our solution unique?**
   - **252 patterns across 18 categories** vs upstream's 57/5 — broadest open pattern library
   - **MCP server** (`atheon-mcp`) — first-class AI-assistant integration
   - **Library + CLI + MCP** from one Go module
   - **SARIF output** — integrates with GitHub Security tab
   - **Pattern state persistence** — survives across runs
   - **Self-scanning** — the tool scans its own codebase in CI

### Research Findings

The [ANALYSIS_REPORT.md](ANALYSIS_REPORT.md) §2 documents the 30+ concrete gaps identified
during the 2026-06-23 audit. Key findings:

- Engine layer (`core/`, `cmd/`) is **mature and idiomatic** — context cancellation, sentinel
  errors, slog, SARIF, MCP rate limiting all present.
- Documentation layer is **stale and contradictory** — multiple docs disagree on pattern counts;
  PLAN/TASKS/PROGRESS files were unfilled templates until this update.
- CI/CD has **10 workflows with significant duplication** — CodeQL in 2, full test chain in 3.
- Pattern library has an **empty `frameworks/` category** despite being documented.

### Technology Options

| Option | Pros | Cons | Recommendation |
|--------|------|------|----------------|
| Go (current) | Single static binary, stdlib regex (RE2), great MCP story | None for this use case | **Keep** |
| RE2 regex engine | Bounded memory, linear time, no backtracking | No lookahead/lookbehind | **Keep** — matches all current patterns |
| Gzip+JSON bundle | Small (~50KB), embeddable, easy to parse | Not human-readable | **Keep** |
| MCP (stdio JSON-RPC) | Native to Claude/Cursor/Windsurf | stdio-only limits scaling | **Keep** — primary AI integration |

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        User Surface                         │
├──────────────┬──────────────────┬───────────────────────────┤
│  atheon CLI  │  atheon-mcp      │  Go library               │
│  (cmd/atheon)│  (cmd/mcp)       │  (import core)            │
└──────┬───────┴────────┬─────────┴──────────┬───────────────┘
       │                │                    │
       └────────────────┴────────────────────┘
                        │
              ┌─────────▼──────────┐
              │   core/            │
              │  ├─ pattern.go     │  registry, ValidatePattern
              │  ├─ bundle.go      │  load, enable/disable, download
              │  ├─ runner.go      │  ScanFile/ScanDir/ScanString/ScanEnv
              │  ├─ ignore.go      │  .atheonignore, .gitignore
              │  ├─ pattern_state  │  persisted enabled/disabled
              │  └─ finding.go     │  result type
              └─────────┬──────────┘
                        │ reads at init()
              ┌─────────▼──────────┐
              │  community/*.yaml  │  252 patterns, 18 categories
              │       │            │
              │       ▼            │
              │  bundler/          │  compiles YAML → core/patterns.bundle
              └────────────────────┘
```

### Technology Stack

- **Language**: Go 1.21+ (single static binary, no runtime deps)
- **Engine**: RE2 regex (stdlib `regexp`) — bounded-memory, no catastrophic backtracking
- **Bundle format**: gzip + JSON, `//go:embed` into the binary
- **MCP transport**: stdio JSON-RPC 2.0
- **CI**: GitHub Actions (10 workflows; consolidation planned)
- **Lint**: golangci-lint v1.64.8 (18 linters)
- **Vulnerability**: govulncheck in CI
- **Coverage**: Codecov v5 with project+patch thresholds
- **Output formats**: human-readable, JSON, SARIF 2.1.0

### Justification

The chosen stack is constrained by:

1. **MCP requirement** — stdio JSON-RPC is what Claude/Cursor/Windsurf expect. A Go binary with
   stdlib `encoding/json` covers this with zero dependencies.
2. **Static-binary distribution** — users want `go install` or a downloaded binary, not a
   runtime. CGO_ENABLED=0 across all builds.
3. **RE2 over PCRE** — Go's stdlib regex engine prevents ReDoS by design. The pattern library
   is large enough that a backtracking engine would be a foot-gun.
4. **No external dependencies** — `go.mod` has zero non-stdlib requires. This means no supply
   chain risk and no version-conflict noise for users who import the library.

---

## Development Approach

### Methodology

Lightweight Agile — single maintainer, weekly release cadence, PR-driven.

1. **Sprint Length**: 1 week (Mon → Sun)
2. **Planning**: Continuous — PR title describes the change
3. **Testing**: `go test ./... -p 1 -coverprofile=coverage.out` (the `-p 1` is **mandatory**
   because `core` package has package-level state in `init()`)
4. **Deployment**: GoReleaser on tag; tags cut automatically on the 10th and 21st of each month
   via `scheduled-release.yml`

### Quality Standards

| Standard | Target | Enforced by |
|----------|--------|-------------|
| Code coverage (project) | 90% | `codecov.yml` |
| Code coverage (patch) | 80% | `codecov.yml` |
| Coverage floor (CI gate) | 70% | `ci.yml` threshold check |
| Lint warnings | 0 | `golangci-lint v1.64.8`, 18 linters |
| Vulnerabilities | 0 high/critical | `govulncheck` in CI |
| Conventional commits | required | README + commit template |
| PR review | 1 approval from CODEOWNERS | GitHub branch protection |
| Branch protection on `main` | required | GitHub settings (verify §9 of IMPROVEMENT_PLAN) |

---

## Implementation Phases

### Phase A — Trust restoration (CURRENT — week of 2026-06-23)

**Goal**: Make every docs number match reality. Fix broken links.

See [ANALYSIS_REPORT.md](ANALYSIS_REPORT.md) §3 Phase A for the 8-item checklist.

**Deliverables**:
- `docs/PLAN.md`, `docs/TASKS.md`, `docs/PROGRESS.md` filled with real content
- Pattern count consistent across `README.md`, `FAQ.md`, `INSTALL.md`, `SETUP.md`,
  `PATTERN_CATEGORIES.md`
- `docs/README.md` broken links either stubbed or removed
- One canonical `BRANCH_STRATEGY.md`

### Phase B — CI consolidation (next 1–2 weeks)

**Goal**: Cut CI minutes, eliminate duplication across the 10 workflows.

See [ANALYSIS_REPORT.md](ANALYSIS_REPORT.md) §3 Phase B.

**Deliverables**:
- 4 workflows total (down from 10): `ci.yml`, `security.yml`, `release.yml`, `sync.yml`
- All `go test` invocations guarded with `-p 1`
- `gofmt -l` check replaces broken `git diff --name-only` check

### Phase C — MCP completeness (parallel with B)

**Goal**: Make MCP server match the API docs.

See [ANALYSIS_REPORT.md](ANALYSIS_REPORT.md) §3 Phase C.

**Deliverables**:
- `list_patterns`, `list_categories`, `scan_env`, `update_bundle` tools added
- Version injected via ldflag

### Phase D — Pattern expansion (ongoing weekly)

**Goal**: Reach 300+ patterns, fix `frameworks/` empty category.

Three pattern batches per the IMPROVEMENT_PLAN.md §3 model.

**Deliverables**:
- 10 SaaS-token patterns (Anthropic, OpenAI org-scoped, Supabase, Vault, etc.)
- 5 PII patterns (email, IP literals, more passport formats)
- Restore 3 frameworks patterns OR remove category

### Phase E — Architecture hygiene (1 week, after B)

**Goal**: Match the SYSTEM_ARCHITECTURE.md to reality.

**Deliverables**:
- 3 ADRs in `docs/architecture/decisions/`
- Either create `core/streaming.go` or fix the doc

### Phase F — Future features (post-release)

- Pattern metadata (severity, description, references) — wire-format change
- `--baseline` filter for incremental scans
- SBOM generation in release workflow
- LSP mode for IDE integration

---

## Risk Management

| Risk | Impact | Mitigation |
|------|--------|------------|
| User trusts outdated docs | High | Phase A fixes numbers; generated counts prevent recurrence |
| CI minutes cost balloons | Medium | Phase B consolidation cuts ~40–60% of runs |
| MCP users hit missing tools | Medium | Phase C adds the four advertised tools |
| False positive rate grows with patterns | Medium | Pattern metadata + confidence scoring (Phase F) |
| Bundle corruption | Low | Gzip + JSON format with init-time validation; corrupt bundle falls back to embedded |
| Path traversal via symlinks | Low | Use `filepath.EvalSymlinks` in `ScanDir` (planned) |
| API breakage for library users | High | All public APIs context-aware since 0.3.0; future changes require ADR |

---

## Success Metrics

### Technical Metrics (measured quarterly)

| Metric | Current (2026-06-23) | Target Q3 2026 |
|--------|---------------------|----------------|
| Pattern count | 252 | 300+ |
| Categories | 18 | 20 |
| Test coverage (project) | 97%+ | ≥95% sustained |
| Test coverage (patch) | 80% gate | 85% |
| CI minutes / week | ~X | -50% post-consolidation |
| Lint warnings | 0 | 0 |
| Critical vulnerabilities | 0 | 0 |
| Workflow count | 10 | 4 |
| Docs files with stale numbers | ~6 | 0 |

### Community Metrics

| Metric | Current | Target |
|--------|---------|--------|
| GitHub stars (this fork) | low triple digits | n/a (not the focus) |
| Pattern contributions / quarter | varies | ≥5 from non-maintainer |
| Open issues closed < 30 days | varies | ≥80% |
| Releases / quarter | 6 (10th & 21st schedule) | 6 maintained |

---

## Timeline

```
2026-06-23  ━ Phase A (docs trust)        ← YOU ARE HERE
2026-06-30  ━ Phase B (CI consolidation)  ━ Phase C (MCP) in parallel
2026-07-07  ━ Phase D batch 1 (SaaS secrets)
2026-07-14  ━ Phase D batch 2 (PII)
2026-07-21  ━ Phase D batch 3 (frameworks)
2026-07-28  ━ Phase E (ADRs + arch hygiene)
2026-08-XX  ━ Phase F begins (pattern metadata, --baseline)
```

---

## Dependencies

### External Dependencies (zero non-stdlib)
- **go**: 1.21+
- **git**: for hooks and ignore semantics
- **jq**: used in CI scripts (not a binary dep — only CI tooling)
- **golangci-lint v1.64.8**: CI linter
- **go-junit-report**: CI test reporter
- **govulncheck**: CI vulnerability scanner
- **Codecov**: coverage reporting
- **GoReleaser**: release automation

### Internal Dependencies
- `core` ← `cmd/atheon`, `cmd/mcp`, `bundler`
- `bundler` reads `community/**/*.yaml` → produces `core/patterns.bundle`
- `core/patterns.bundle` ← `//go:embed` in `core/bundle.go`

---

## Appendix

### Assumptions

- Single maintainer (`@aliasfoxkde`) — changes go through CODEOWNERS review
- No corporate proxy blocks `go install` in CI (documented in IMPROVEMENT_PLAN §8.1 for dev
  machines that do)
- The 18 categories are stable; new categories require an ADR
- The fork's relationship to upstream HoraDomu/Atheon is "test ground, propose upstream" — not
  a competing product

### Constraints

- Go 1.21+ compatibility (no `min`, `max`, `clear` from 1.21 — they ARE available, so the
  constraint is the floor not the ceiling)
- Single static binary — no runtime deps
- Zero non-stdlib dependencies in `core/` (matches `.golangci.yml` discipline)
- RE2 regex only (no PCRE features)

### Open Questions

- Should `frameworks/` be deleted or restored? (Pending Phase D batch 3)
- Should `--baseline` filter live in `core` or as a CLI-only feature? (Phase F design)
- Pattern metadata schema — what's the minimum useful set? (Phase F design)

### Reference

- **Canonical analysis**: [docs/ANALYSIS_REPORT.md](ANALYSIS_REPORT.md)
- **Improvement plan (older, partially executed)**: [docs/reports/IMPROVEMENT_PLAN.md](reports/IMPROVEMENT_PLAN.md)
- **Completion report (older)**: [docs/reports/COMPLETION_REPORT.md](reports/COMPLETION_REPORT.md)
- **Roadmap**: [docs/reports/ROADMAP.md](reports/ROADMAP.md)
- **Branch strategy**: [docs/BRANCH_STRATEGY.md](BRANCH_STRATEGY.md)
- **System architecture**: [docs/architecture/SYSTEM_ARCHITECTURE.md](architecture/SYSTEM_ARCHITECTURE.md)
- **Pattern categories**: [docs/architecture/PATTERN_CATEGORIES.md](architecture/PATTERN_CATEGORIES.md)
- **API reference**: [docs/api/README.md](api/README.md)
- **Setup guide**: [docs/development/SETUP.md](development/SETUP.md)
- **Troubleshooting**: [docs/guides/TROUBLESHOOTING.md](guides/TROUBLESHOOTING.md)
- **Changelog**: [../CHANGELOG.md](../CHANGELOG.md)
- **Recent changelog excerpt**: [CHANGELOG_RECENT.md](CHANGELOG_RECENT.md)
