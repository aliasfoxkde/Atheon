# Tasks — Atheon-Enhanced

**Version:** 0.3.0 (planned)
**Last Updated:** 2026-06-20

---

## Status legend

- [ ] Pending
- [~] In Progress
- [x] Completed
- [!] Blocked
- [-] Cancelled

---

## Phase 0: Planning (COMPLETE)

- [x] Read upstream repo and fork repo fully
- [x] Audit current state of fork
- [x] Identify technical debt (20 items)
- [x] Design feature roadmap (3 tiers)
- [x] Write pattern library strategy
- [x] Write MCP integration plan
- [x] Calculate token economics
- [x] Design benchmark v1 and v2
- [x] Write white paper outline
- [x] Write go-to-market plan
- [x] Design demos and infrastructure
- [x] Write upstream contribution playbook
- [x] Write release and versioning plan
- [x] Identify risks and limitations
- [x] Write day-by-day action plans
- [x] Write quick reference cheat sheet
- [x] Fill in top-level `PLAN.md`, `PROGRESS.md`, `TASKS.md`

**Deliverables:**
- 18 documents in `docs/planning/atheon-enhanced/`
- ~4,900 lines of planning content
- Linked from `PLAN.md`, `PROGRESS.md`, `TASKS.md`

**See:** [`docs/planning/atheon-enhanced/`](./planning/atheon-enhanced/00_README.md)

---

## Phase 1: Tier 1 — This Week (PENDING)

### Hygiene (TD-01 through TD-04)

- [ ] **TD-01:** Commit the 28 uncommitted working-tree changes
      to feature branches
  - `feat/engine-rwmutex` (engine.go + context_cancel_test.go)
  - `fix/bundler-...` (bundler fixes)
  - One commit per logical change
- [ ] **F-02:** Wire Engine RWMutex into public API;
      `go test -race ./...` clean
- [ ] **F-03:** Strip `config/profiles/*.json` to 3 lines each
- [ ] **F-04:** Fix README pattern counts (57→179, 85→179,
      152→179)
- [ ] **F-15:** MCP server version via -ldflags
- [ ] **F-16:** Move fork binaries out of repo root
      (.gitignore)

### Ship the `gate` tool (F-05)

- [ ] Implement `gate(content, source, categories)` in
      `cmd/mcp/main.go`
- [ ] Add `gate` to `tools/list` output
- [ ] Tests:
  - [ ] `TestGateApproved` (clean content)
  - [ ] `TestGateRejected` (one finding)
  - [ ] `TestGateEmptyContent`
  - [ ] `TestGateCategories` (categories filter works)
  - [ ] `TestGateStructuredFindings` (JSON shape)
- [ ] Update README's MCP integration section

### Document the integration pattern (F-06)

- [ ] Create `docs/MCP_INTEGRATION.md`
- [ ] Add system-prompt snippet for Claude Code, Cursor,
      Windsurf, Continue.dev, Zed
- [ ] Link from README

### Build the demo (F-21, F-22)

- [ ] Audit `atheon-scanner.pages.dev`
- [ ] Audit `atheon-benchmark.pages.dev`
- [ ] Build `atheon-gate.pages.dev` (JS port of patterns)
- [ ] Link all three demos from README

### Tag a release

- [ ] Update CHANGELOG.md
- [ ] Tag `v0.3.0` (renumber from 1.2.0)
- [ ] Push tag, verify goreleaser runs (or manual release)

### Launch

- [ ] Write Show HN draft
- [ ] Post Show HN (Tue/Wed, 8-10am US Pacific)
- [ ] Cross-post to r/golang, r/programming,
      r/MachineLearning, dev.to
- [ ] Reply to every comment for 24 hours

---

## Phase 2: Tier 2 — This Month (PENDING)

### Upstream PRs (F-08)

- [ ] PR: Engine RWMutex (if not already)
- [ ] PR: Context cancellation tests
- [ ] PR: Bundler fixes
- [ ] PR: MCP server version via -ldflags

### Benchmark v1 (F-11)

- [ ] Build labelled dataset
  - [ ] 50 true positives per category
  - [ ] 50 true negatives per category
  - [ ] 50 AI-generated vs 50 human-written examples
  - [ ] Pin dataset hash in
        `benchmarks/dataset.SHA256`
- [ ] Write runner (`benchmarks/runner.go`)
- [ ] Run benchmark; produce `benchmarks/results/v1.{json,md}`
- [ ] Wire `make bench-v1` into CI
- [ ] PR comment on every PR with the diff vs main

### White paper v1 (F-12)

- [ ] Assemble benchmark v1 numbers into paper structure
- [ ] Write prose
- [ ] Produce `docs/WHITEPAPER.md` (markdown)
- [ ] Produce `docs/WHITEPAPER.pdf` (pandoc)
- [ ] Link from README

### Profiles (F-09)

- [ ] Implement `Profile struct` in `core/`
- [ ] Implement `LoadProfile(path)`
- [ ] Add `--profile <path>` CLI flag
- [ ] Update `config/profiles/*.json` with implemented fields
- [ ] Add "Profiles" section to README

### Per-pattern descriptions (F-10)

- [ ] Add `description:` field to YAML schema
- [ ] Migrate most-trafficked 30 patterns to v2
- [ ] Update `atheon list` to show descriptions
- [ ] Display descriptions in MCP tool responses

### Structured findings (F-13)

- [ ] Add `findings` field to MCP `tools/call` responses
- [ ] Update README's MCP integration section

### Allow-list (F-14)

- [ ] Add `pattern:allow <name>` to `.atheonignore`
- [ ] Test: allow-list prevents pattern from firing
- [ ] Document in README

---

## Phase 3: Tier 3 — This Quarter (PENDING)

### `auto_fix` tool (F-17)

- [ ] Design tool schema
- [ ] Implement in `cmd/mcp/main.go`
- [ ] Implement redaction in `core/redact.go`
- [ ] Tests
- [ ] Update README

### Per-call pattern injection (F-18)

- [ ] Add `pattern:` parameter to `scan_string`
- [ ] Implement one-off regex in `core/runner.go`
- [ ] Tests

### Benchmark v2 (F-19)

- [ ] Design 40 tasks
- [ ] Run control + treatment (40 tasks × 2 conditions)
- [ ] 4 human raters, blind to condition
- [ ] Statistical analysis (Wilcoxon, Cohen's d)
- [ ] Produce `benchmarks/results/v2.md`

### White paper v2 (F-12)

- [ ] Update with v2 benchmark results
- [ ] Add agent-vs-agent comparison section
- [ ] Re-publish PDF + markdown

### Distribution (F-23, F-24, F-21)

- [ ] Show HN, take 2
- [ ] Submit to MCP directories (mcp.so, glama.ai/mcp,
      awesome-mcp)
- [ ] Build `atheon-docs.pages.dev` (Docusaurus)
- [ ] Cross-post to golangweekly.com, dev.to, Hashnode
- [ ] VS Code extension (separate repo)
- [ ] GitHub Action (`.github/actions/atheon`)
- [ ] Pre-commit hook

### Per-language patterns (F-25)

- [ ] Add `languages:` field to YAML schema
- [ ] Update bundler to validate
- [ ] Update engine to filter by file language
- [ ] Migrate web-security + frameworks patterns to v2

### Bundle decoupling (F-26)

- [ ] Update `.goreleaser.yml` to publish
      `atheon-patterns-vX.Y.Z.bundle` as separate artifact
- [ ] Document install-into-upstream-binary flow

---

## Phase 4: Steady State (PENDING)

### Quarterly cadence

- [ ] Quarterly benchmark releases
- [ ] Year-in-review blog post (December)
- [ ] Conference talk submission (GopherCon, AI Engineer
      Summit)

### Community

- [ ] Enable GitHub Discussions
- [ ] Add `good first pattern` label
- [ ] Pin "What we're working on" issue
- [ ] Triage issues daily for first 2 weeks, then weekly
- [ ] Find a co-maintainer (upstream maintainer first
      choice)

### Releases

- [ ] Every-other-Tuesday cadence
- [ ] Patch as needed for security
- [ ] Major with rc process (2-week rc)

---

## Progress summary

- **Total tasks planned:** ~80
- **Completed:** 16 (Phase 0: planning)
- **In progress:** 0
- **Pending:** ~64
- **Completion:** 20%

---

## Notes

- Tasks are ordered by tier and by file path.
- Each task should result in a commit on a feature branch.
- Each tier should end with a tagged release.
- Tier 3 is conditional on Tier 2 being 50%+ complete by
  month 2.
- See [`16_ACTION_PLANS.md`](./planning/atheon-enhanced/16_ACTION_PLANS.md)
  for day-by-day execution.