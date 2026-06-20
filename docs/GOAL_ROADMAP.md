# /goal Implementation Roadmap — Atheon-Enhanced

**Created:** 2026-06-20
**Status:** ACTIVE
**Complements:** [`docs/PLAN.md`](../PLAN.md), [`docs/TASKS.md`](../TASKS.md), [`docs/PROGRESS.md`](../PROGRESS.md)

---

## Purpose

The /goal directive is a production-quality sweep across the entire
Atheon-Enhanced codebase. It is **orthogonal** to the existing
launch/GTM phases in `docs/PLAN.md` (which target `gate`, demos,
and Show HN). The /goal targets engineering quality and the
features the user has specifically called out.

## What the /goal covers

1. **Add an ultra-fast network scanner** — scan HTTP/HTTPS URLs and
   remote git repositories in addition to local files.
2. **Add structured reporting** — multi-format reports (findings,
   patterns, file types, severity breakdown, JSON/SARIF/HTML).
3. **Generate/improve technical documentation** — fill gaps,
   remove staleness, make docs production-grade.
4. **Audit for dead and redundant code** + **make the audit
   enforceable** in the harness and as a project quality gate.
5. **Audit upstream issues** at `github.com/HoraDomu/Atheon/issues`
   and fix them in the fork; create feature/fix branches for
   future PR submissions.
6. **Improve the audit pipeline** — make audits repeatable, fast,
   and CI-friendly.
7. **Refactor to production quality** — withstand external
   scrutiny.

## Why a separate document

The existing `PLAN.md` is heavily GTM-focused. The /goal work
overlaps some items (upstream issues) but the **dead-code
enforcement**, **network scanner**, and **structured reporting**
are new features that deserve a focused plan.

## Phasing

The /goal is divided into 7 phases. **Each phase ends with a
commit and (where applicable) a push to a feature branch.** The
existing 28 uncommitted changes are the *first* item; we do not
start new work until the working tree is clean.

### Phase 0 — Reconnaissance (DONE)

- [x] Read existing core/, cmd/, bundler/ source
- [x] Audit upstream `HoraDomu/Atheon` open issues (20 issues)
- [x] Map CLI/MCP surface and existing tests
- [x] Map harness system (`/home/mkinney/.claude/hooks/`)
- [x] Confirm pre-commit/pre-push/pre-pr hooks in `.githooks/`

### Phase 1 — Foundation: dead-code + harness enforcement

> **Why first:** the user explicitly called this out as needing
> to be "a pattern and part of the quality checks and tests ...
> added to our harness/backend system and enforced for best
> practices." Every subsequent phase benefits from this
> enforcement running in pre-commit.

- [ ] **G-1.1:** Audit `core/`, `cmd/`, `bundler/` for dead code
      - Run `go vet`, `staticcheck`, `unused`, `golangci-lint run`
      - Grep for unexported functions defined but never called
      - Identify `// nolint` comments that should be removed
      - Output: `docs/audits/DEAD_CODE_AUDIT.md`
- [ ] **G-1.2:** Fix the dead code
      - Remove unused helpers (e.g. `contains` per #159)
      - Remove duplicate checks (#155)
      - Replace `// nolint` with real fixes where possible
- [ ] **G-1.3:** Add `make audit` target (run all audits in one go)
- [ ] **G-1.4:** Wire dead-code check into `.githooks/pre-commit`
      - `staticcheck` + `unused` + `go vet`
      - Block commit if dead code is found (unless exempted)
- [ ] **G-1.5:** Add `dead-code-prevention-hook.py` to global harness
      - Path: `/home/mkinney/.claude/hooks/dead-code-prevention-hook.py`
      - PreToolUse: scan Write/Edit for new unused exports
      - Use `go build` to verify before allowing Write
- [ ] **G-1.6:** Commit to `chore/audit-dead-code-cleanup`
- [ ] **G-1.7:** Open upstream PR for #155, #159, #160, #161

### Phase 2 — Upstream issue fixes (one-branch-per-issue)

> **Why second:** most are 1-2 line fixes; they warm up the
> branch-and-PR workflow we will use for the network scanner
> and reporting.

- [ ] **G-2.1:** `fix/scandir-error-propagation` — #157
      - Make `ScanDir` collect and return per-file read errors
      - Update tests
- [ ] **G-2.2:** `fix/list-output-deterministic` — #158
      - `sort.Slice` patterns by name before printing
- [ ] **G-2.3:** `fix/json-flag-position-independent` — #156
      - Move `--json` check into the flag-parsing loop
- [ ] **G-2.4:** `feat/update-reports-changes` — #127
      - `atheon update` should print `added: N, removed: M,
        updated: X → Y patterns`
- [ ] **G-2.5:** Open PRs for each branch

### Phase 3 — Structured reporting

> **Why third:** once the foundation is clean, we can build new
> features without re-introducing dead code.

- [ ] **G-3.1:** Define report types in `core/report.go`
      - `FindingsReport` (current JSON output, normalized)
      - `PatternsReport` (list of patterns, by category, enabled/disabled)
      - `FileTypesReport` (extension → finding count, by category)
      - `SeverityReport` (count by severity)
      - `StatsReport` (files, bytes, elapsed, throughput)
- [ ] **G-3.2:** Add `core.Report` and `core.Render(format, w)` API
- [ ] **G-3.3:** Add CLI flags: `--report=<type>`, `--format=json|yaml|sarif|html`
- [ ] **G-3.4:** Implement SARIF 2.1.0 output (GitHub Code Scanning)
- [ ] **G-3.5:** Implement HTML output (single-file, no JS)
- [ ] **G-3.6:** Tests for each report type and format
- [ ] **G-3.7:** Document in `docs/REPORTING.md`
- [ ] **G-3.8:** Commit to `feat/structured-reporting`

### Phase 4 — Ultra-fast network scanner

> **Why fourth:** the scanner is the most independent feature;
> keeping it last in the major-feature queue means the report
> format from Phase 3 is available to use as the output.

- [ ] **G-4.1:** Design the network surface
      - `ScanURL(ctx, url string)` — fetch + scan response
      - `ScanGitRemote(ctx, url string)` — shallow clone + scan
      - `ScanAPI(ctx, base, paths []string)` — multi-endpoint
- [ ] **G-4.2:** Implement `core/scanner_net.go`
      - HTTP client with timeouts, redirect policy, robots.txt
        respect
      - Body size limit (default 5MB, configurable)
      - Content-Type aware (text/* scanned, application/* skipped
        unless `--scan-binary` set)
      - Streaming download (no `io.ReadAll` of unbounded bodies)
- [ ] **G-4.3:** Implement `core/scanner_git.go`
      - `go-git` for shallow clone (or `git clone --depth 1`
        fallback)
      - Use existing `ScanDir` on the cloned tree
- [ ] **G-4.4:** Add CLI: `atheon scan-url <url>` and `atheon
      scan-git <url>`
- [ ] **G-4.5:** Add MCP tools: `scan_url`, `scan_git`
- [ ] **G-4.6:** Use the Phase 3 reporter for output
- [ ] **G-4.7:** Tests with `httptest.Server`
- [ ] **G-4.8:** Document in `docs/NETWORK_SCANNER.md`
- [ ] **G-4.9:** Commit to `feat/network-scanner`

### Phase 5 — Audit pipeline improvements

> **Why fifth:** with reporting in place, the audit pipeline can
> produce real reports (vs ad-hoc output).

- [ ] **G-5.1:** New `core/audit.go` package
      - `Audit(ctx, root) AuditReport`
      - Runs all audits: dead-code, linter, security, coverage,
        docs staleness
- [ ] **G-5.2:** `make audit` runs `Audit` and writes
      `docs/audits/REPORT.md` + `.json`
- [ ] **G-5.3:** Wire audit into CI (`.github/workflows/audit.yml`)
- [ ] **G-5.4:** Document in `docs/AUDITING.md`

### Phase 6 — Documentation sweep

> **Why sixth:** by now the surface has new features (reporting,
> network) that need docs, and existing docs may have stale
> references.

- [ ] **G-6.1:** Re-audit every doc for staleness (use Phase 5
      pipeline)
- [ ] **G-6.2:** Update README to reflect 179 patterns,
      network scanner, and reporting
- [ ] **G-6.3:** Add `docs/REPORTING.md`, `docs/NETWORK_SCANNER.md`,
      `docs/AUDITING.md`
- [ ] **G-6.4:** Update `docs/API.md` to include new APIs
- [ ] **G-6.5:** Update `docs/ARCHITECTURE.md` to include new
      packages

### Phase 7 — Production hardening

- [ ] **G-7.1:** Run `go test -race -coverprofile=... ./...` and
      confirm ≥95% coverage
- [ ] **G-7.2:** Run `staticcheck`, `golangci-lint`, `gosec` and
      address all findings
- [ ] **G-7.3:** `govulncheck ./...` clean
- [ ] **G-7.4:** Tag `v0.4.0` (production-quality release)
- [ ] **G-7.5:** Push all branches; open remaining PRs upstream

---

## Branch strategy

| Branch | Phase | Target |
|---|---|---|
| `chore/audit-dead-code-cleanup` | 1 | merged to main |
| `fix/scandir-error-propagation` | 2 | upstream PR |
| `fix/list-output-deterministic` | 2 | upstream PR |
| `fix/json-flag-position-independent` | 2 | upstream PR |
| `feat/update-reports-changes` | 2 | upstream PR |
| `feat/structured-reporting` | 3 | merged to main |
| `feat/network-scanner` | 4 | merged to main |
| `feat/audit-pipeline` | 5 | merged to main |
| `docs/goal-sweep` | 6 | merged to main |
| `release/v0.4.0` | 7 | tag |

## Commits and pushes

- **One branch per logical feature**, even for 1-line changes.
- **One commit per logical change** within a branch.
- **Push to `origin` (aliasfoxkde)** at end of each phase.
- **Upstream PRs** opened at end of Phase 2; opened progressively
  as each upstream-PR-targeted branch lands in main.
- **No `--no-verify`.** If a hook blocks, fix the underlying issue
  and retry. This rule is enforced by the global harness.

## Memory hygiene

- Add a new memory at the end of each phase summarizing
  non-obvious findings, decisions, and gotchas.
- Update `MEMORY.md` index.

## What this plan does NOT cover

- The GTM work in `docs/PLAN.md` (Show HN, demos, gate tool).
  Those continue on their own schedule.
- Multi-language scanner support.
- Hosted SaaS version.
- The 6 AI-detection patterns (orthogonal, may be deprecated).
