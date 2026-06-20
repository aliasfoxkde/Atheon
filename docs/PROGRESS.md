# Project Progress — Atheon-Enhanced

**Last Updated:** 2026-06-20
**Current Phase:** Phase 1 — Tier 1 (this week)
**Overall Progress:** ~5% (planning done; Tier 1 work not yet
started)

---

## Where the detailed plan lives

All planning detail is in
[`docs/planning/atheon-enhanced/`](./planning/atheon-enhanced/00_README.md).
This file is the high-level progress tracker.

---

## Progress summary

| Phase | Tier | Status | Target |
|---|---|---|---|
| Planning | — | ✅ DONE (18 docs) | 2026-06-20 |
| Phase 1: launch | Tier 1 | ⏳ NOT STARTED | 2026-06-27 |
| Phase 2: credibility | Tier 2 | ⏳ NOT STARTED | 2026-07-20 |
| Phase 3: distribution | Tier 3 | ⏳ NOT STARTED | 2026-09-20 |
| Phase 4: sustained | — | ⏳ NOT STARTED | 2026-12-20 |

---

## Phase 1 — Tier 1 (this week, target 2026-06-27)

### 1.1 Hygiene (Day 1)

- [ ] Commit `core/engine.go` + `core/context_cancel_test.go`
      to `feat/engine-rwmutex` branch
- [ ] Wire Engine into public API; `go test -race ./...` clean
- [ ] Strip `config/profiles/*.json` to 3 lines each
- [ ] Fix README pattern counts (57→179, 85→179, 152→179)
- [ ] Audit existing demos: `atheon-scanner.pages.dev`,
      `atheon-benchmark.pages.dev`

### 1.2 Ship the `gate` tool (Day 2-4)

- [ ] Implement `gate(content, source, categories)` in
      `cmd/mcp/main.go`
- [ ] Tests: `TestGateApproved`, `TestGateRejected`,
      `TestGateEmptyContent`, `TestGateCategories`
- [ ] Update `tools/list` to include `gate`
- [ ] Update `docs/MCP_INTEGRATION.md` with system-prompt
      snippet

### 1.3 Ship the demo (Day 5)

- [ ] Build `atheon-gate.pages.dev` (JS port of patterns;
      1-page paste-text-see-findings)
- [ ] Add demo link to README
- [ ] Tag `v0.3.0` (renumber from 1.2.0 to be honest about
      pre-1.0)
- [ ] Set GitHub repo topics: `mcp`, `ai`, `code-review`,
      `secrets`, `developer-tools`, `model-context-protocol`

### 1.4 Launch (Day 5)

- [ ] Post Show HN (Tue/Wed, 8-10am US Pacific)
- [ ] Cross-post to r/golang, r/programming,
      r/MachineLearning, dev.to
- [ ] Reply to every comment for 24 hours

See [`16_ACTION_PLANS.md`](./planning/atheon-enhanced/16_ACTION_PLANS.md)
for the full day-by-day breakdown.

---

## Phase 2 — Tier 2 (this month, target 2026-07-20)

- [ ] Open 3+ PRs upstream (engine RWMutex, context
      cancellation, bundler fixes)
- [ ] Build labelled dataset (`benchmarks/dataset.jsonl`)
- [ ] Write runner (`benchmarks/runner.go`)
- [ ] Run benchmark v1; produce `benchmarks/results/v1.md`
- [ ] Write white paper v1 (`docs/WHITEPAPER.md`)
- [ ] Implement `--profile` flag end-to-end
- [ ] Add per-pattern `description:` field
- [ ] Add `findings` field to MCP tool responses
- [ ] Add `pattern:allow <name>` to `.atheonignore`

---

## Phase 3 — Tier 3 (this quarter, target 2026-09-20)

- [ ] Ship `auto_fix` MCP tool
- [ ] Per-call pattern injection
- [ ] Build benchmark v2 (human-eval rubric)
- [ ] White paper v2
- [ ] Blog post on dev.to
- [ ] Show HN (take 2)
- [ ] Per-language pattern files
- [ ] Decouple pattern bundle from binary

---

## Recent activity

### 2026-06-20 — planning complete

- 18 planning documents created (~4,900 lines) in
  `docs/planning/atheon-enhanced/`
- 3 top-level docs (`PLAN.md`, `PROGRESS.md`, `TASKS.md`)
  filled in
- All planning tasks in the task list marked completed

---

## Open blockers

None at the planning phase. Phase 1 work begins once the user
reviews the plan and confirms priorities.

---

## Risks (top 3)

1. **R-05: Maintainer burnout.** The plan is too big for one
   person. Cut Tier 3 if Tier 2 isn't done by month 2.
2. **R-01: Drift from upstream.** The fork has 28 uncommitted
   changes; needs immediate hygiene.
3. **R-03: AI-detection discredited.** Patterns have high FP
   rate; default to opt-in, don't lead with this category.

Full risk register:
[`15_RISKS_AND_LIMITATIONS.md`](./planning/atheon-enhanced/15_RISKS_AND_LIMITATIONS.md).

---

## Metrics to track

- GitHub stars (target: 50+ post-launch, 300+ by month 3)
- Weekly downloads (target: 100+ post-launch, 500+ by month 3)
- MCP directory listings (target: 2+ post-launch, 5+ by month 3)
- Upstream PRs merged (target: 3+ in 2026-Q3, 6+ in 2026-Q4)
- Test coverage (maintain 98.7%+)