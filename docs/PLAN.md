# Project Plan — Atheon-Enhanced

**Version**: 0.6.0 (post-Wave 6)
**Last Updated**: 2026-06-25
**Status**: Active — multi-wave hardening cycle in progress

---

## Problem & Audience

**Problem.** Repositories leak secrets, PII, and security-sensitive code patterns. Off-the-shelf scanners (TruffleHog, GitLeaks, Semgrep) each cover a slice. We want a single fast scanner that handles **all three** (secrets + PII + code-quality) with a **transparent, community-editable pattern catalog** and an **MCP server** so IDEs can query it.

**Audience.** Open-source maintainers, security-conscious teams, AI-coding-assistant users who want inline pattern feedback inside their editor.

**Differentiators.**
- 274 patterns across 19 categories (secrets, PII, web-security, compliance, security-hardening, cloud-native, devops, healthcare, finance, ai-detection, code-quality, performance, accessibility, web-development, data-visualization, api-integration, git-hygiene, pwa, frameworks).
- Pattern catalog is plain YAML in `community/` — anyone can add a pattern with a PR. No proprietary DSL.
- Per-pattern severity (low/medium/high/critical) flows end-to-end into SARIF for IDE integration.
- MCP server (`atheon-mcp`) exposes the same scanner to Claude/Cursor/etc. via JSON-RPC.

## Architecture

```
                ┌──────────────────────┐
                │ community/*.yaml     │  274 patterns, declarative
                │ (severity, regex)    │
                └──────────┬───────────┘
                           │  go run ./bundler
                           ▼
                ┌──────────────────────┐
                │ core/patterns.bundle │  gzipped JSON, embedded
                └──────────┬───────────┘
                           │  go:embed
                           ▼
   ┌──────────────┐   ┌──────────┐   ┌──────────────┐
   │ cmd/atheon   │   │ core/    │   │ cmd/mcp      │
   │ (CLI scan)   ├──▶│ scanner  │◀──┤ (JSON-RPC)   │
   └──────┬───────┘   └──────┬───┘   └──────────────┘
          │                  │
          ▼                  ▼
       text/JSON/SARIF    Finding{Severity,File,Line,Content}
```

- **Language**: Go (RE2 regex via stdlib `regexp`, guaranteed linear time).
- **Embed**: `//go:embed patterns.bundle` so the binary is self-contained.
- **Concurrency**: parallel per-file scan via goroutines; package-level `init()` loads the bundle exactly once.

## Development Approach

**Methodology**: wave-based hardening. Each wave is one merged PR cluster, scoped by a fresh gap-analysis subagent. Waves do not have fixed sprint lengths — they end when the gap list is exhausted or risk outweighs benefit.

**Iteration cycle per wave**:
1. **Plan** (subagent or direct): enumerate gaps with risk × effort ranking.
2. **Confirm** with user via AskUserQuestion (scope, defer choices).
3. **Implement** on a `feature/wave-N-*` branch off `main`.
4. **Verify**: `go test ./...`, `gofmt -l .`, manual smoke.
5. **Open PR** via REST API (gh CLI's GraphQL is flaky on this repo).
6. **Resolve** any CodeRabbit threads before merge.
7. **Squash-merge** via REST.
8. **Cleanup**: `git branch -d`, update `MEMORY.md` with wave summary.

**Quality gates**:
- Code coverage: ≥70% on touched packages (enforced in `scripts/hooks/pre-commit`).
- `gofmt -l .` must be empty (enforced in CI Lint job).
- All commits end with `Co-Authored-By: Claude <noreply@anthropic.com>`.
- No `git stash`, `git reset --hard`, `git push --force`, `git rebase -i` — per standing repo rules.

## Hardening Waves

| Wave | PR | Theme | Status |
|------|-----|-------|--------|
| 1 | #74-76, #79-80 | Initial scaffold, rename cleanup, gap-analysis docs | Merged |
| 2 | #81 | CI/security: dependabot groups, govulncheck pin, PR template | Merged |
| 3 | #83 | Fuzz tests, detection coverage, -trimpath + SPDX SBOM | Merged |
| 4 | #84 | Pattern severity wired end-to-end (Pattern → Finding → SARIF) | Merged |
| 5 | #85 | Stats.Errors surfacing, list --category validation, hot-path benchmarks, dead-script removal | Merged |
| 6 | #86, #87, #88 | Legacy-flip log, --json --version, MCP JSON-RPC roundtrip, docs fill | Merged |

## Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Pattern regexes accidentally match common code | High (false positives) | `TestFalsePositiveGuard` covers known clean snippets; PR review required for `community/*.yaml` |
| Bundle corruption on bad YAML | Medium (silent drop) | Bundler now logs and skips; pre-commit hook surfaces stderr |
| SARIF severity mismatch consumer expectations | Medium | CVSS-like scores (9.5/7.5/5.0/2.5) + level (error/warning/note) per spec |
| `core/` package-level state races between CLI + MCP | Medium | Documented; mutex work deferred to Wave 7 |
| Coverage drop slips past CI | Low | Codecov status check; `require_ci_to_pass: false` intentional (per ci.yml comment) |

## Success Metrics

**Adoption** (proxy): bundle download count from GitHub releases (visible in Insights).
**Quality**: 
- Pattern false-positive rate (target: <5% on `TestFalsePositiveGuard` corpus)
- Time-to-fix for a known-bad regex (target: <1 wave)
**Reliability**: zero `panic` in CI across all four Go versions (1.21, 1.22, 1.23, 1.24).

## Timeline

No fixed timeline. Each wave lands when ready. Velocity is gated by review throughput, not calendar.

## Dependencies

**External**:
- Go 1.21+ (1.22, 1.23, 1.24 also tested in CI matrix)
- `gopkg.in/yaml.v3` (only direct dep)
- GitHub Actions runners (ubuntu, windows, macos)

**Internal**:
- `core/` is the only package with embedded data; everything else imports it.
- `bundler/` is a separate `main` package that produces the embedded bundle.

## Open Questions

- Should the pattern_state race (gap #2 from Wave 5) ship as Wave 7, or stay deferred until someone reports a real bug?
- Bundle format versioning — should we add a `version: 2` field now that severity is wired through, or wait until the next breaking wire change?
- Should `dev/full-feature` branch (referenced in CI grep) actually exist, or should the canonical `docs/BRANCH_STRATEGY.md` no longer mention it?