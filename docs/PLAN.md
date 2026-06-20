# Project Plan — Atheon-Enhanced

**Version:** 0.3.0 (pre-release; see `docs/planning/atheon-enhanced/14_RELEASE_AND_VERSIONING.md`)
**Last Updated:** 2026-06-20
**Status:** ACTIVE — supersedes prior template

---

## What this document is

This is the project's high-level plan. The detailed plan (18
documents, ~4,900 lines) lives in
[`docs/planning/atheon-enhanced/`](./planning/atheon-enhanced/00_README.md).

Read the [planning folder README](./planning/atheon-enhanced/00_README.md)
first; it has the table of contents and the ground-truth metrics.

---

## What the project is

**Atheon** is a 7 MB Go binary that scans code (or any text)
for secrets, PII, security issues, code smells, and AI tells
using **179 regex patterns across 17 categories**. It runs in
~50 microseconds per scan, is fully deterministic, and costs
zero per check.

**`atheon-mcp`** is the project's Model Context Protocol server.
It exposes `scan_string`, `scan_file`, `scan_dir`, and (the
fork's headline feature) `gate`. Wire it into Claude Code /
Cursor / Windsurf and the AI agent calls `gate` before
returning any output.

**Atheon-Enhanced** (this fork) is the upstream at
`github.com/HoraDomu/Atheon` plus experimental features
under test: the AI-detection patterns, the `gate` tool, the
RWMutex-guarded Engine, and the per-call pattern framework.

---

## The thesis (the killer sentence)

> Why pay an LLM $0.01 and 2 seconds to re-read its own output
> looking for leaked secrets, when a 7 MB binary can do it for
> free in 50 microseconds?

---

## The numbers (verified 2026-06-20)

- **179 patterns** across **17 categories**
- **6 AI-detection patterns** (unique to the fork)
- **98.7% test coverage**
- **~7 MB binary**, zero runtime dependencies
- **Cross-platform:** Linux, macOS (Intel + ARM), Windows
- **License:** MIT with Additional Terms (mirrors upstream)

Full audit: [`docs/planning/atheon-enhanced/03_CURRENT_STATE_AUDIT.md`](./planning/atheon-enhanced/03_CURRENT_STATE_AUDIT.md)
Validated metrics: [`docs/memories/VALIDATED_METRICS.md`](./memories/VALIDATED_METRICS.md)

---

## The four phases

### Phase 1 — this week (Tier 1)

Goal: unblock the GTM. Ship the `gate` tool. Fix the lying
config. Commit the working tree.

- Commit the 28 uncommitted working-tree changes (TD-01)
- Wire the Engine RWMutex into the public API (F-02)
- Strip the unimplemented fields from `config/profiles/*.json`
  (F-03)
- Fix the README pattern counts (F-04)
- Ship the `gate` MCP tool (F-05)
- Build `atheon-gate.pages.dev` (F-21)

### Phase 2 — this month (Tier 2)

Goal: credibility. Run benchmark v1. Write the white paper.
Open upstream PRs.

- Open 3+ PRs upstream (F-08)
- Implement `--profile` flag end-to-end (F-09)
- Add per-pattern `description:` field (F-10)
- Run benchmark v1, publish results table (F-11)
- Write white paper v1 (F-12)
- Add `findings` field to MCP responses (F-13)
- Add allow-list mechanism (F-14)

### Phase 3 — this quarter (Tier 3)

Goal: distribution. Ship `auto_fix`. Run benchmark v2. Do the
GTM push.

- Ship `auto_fix` MCP tool (F-17)
- Per-call pattern injection (F-18)
- Build benchmark v2 with human-eval rubric (F-19)
- Publish white paper v2 + blog post (F-20)
- Build `atheon-gate.pages.dev` + audit existing demos (F-21,
  F-22)
- Show HN (F-23)
- Submit to MCP directories (F-24)
- Per-language pattern files (F-25)
- Decouple pattern bundle from binary (F-26)

### Phase 4 — steady state

Goal: sustain. Quarterly benchmark releases, community-led
pattern review, monthly fork releases, every-other-Tuesday
release cadence.

---

## Where the detailed plan lives

All details are in
[`docs/planning/atheon-enhanced/`](./planning/atheon-enhanced/00_README.md).
The 18 documents are:

| # | File | Purpose |
|---|---|---|
| 00 | [README](./planning/atheon-enhanced/00_README.md) | Index + ground-truth metrics |
| 01 | [OVERVIEW](./planning/atheon-enhanced/01_OVERVIEW.md) | What Atheon is, fork relationship |
| 02 | [VISION_AND_PITCH](./planning/atheon-enhanced/02_VISION_AND_PITCH.md) | The pitch, three lengths |
| 03 | [CURRENT_STATE_AUDIT](./planning/atheon-enhanced/03_CURRENT_STATE_AUDIT.md) | Audit of engine, patterns, docs, CI |
| 04 | [TECHNICAL_DEBT](./planning/atheon-enhanced/04_TECHNICAL_DEBT.md) | 20 prioritized debt items (P0/P1/P2) |
| 05 | [FEATURE_ROADMAP](./planning/atheon-enhanced/05_FEATURE_ROADMAP.md) | 28 features across 3 tiers |
| 06 | [PATTERN_LIBRARY_STRATEGY](./planning/atheon-enhanced/06_PATTERN_LIBRARY_STRATEGY.md) | Pattern curation + AI-detection |
| 07 | [MCP_INTEGRATION_PLAN](./planning/atheon-enhanced/07_MCP_INTEGRATION_PLAN.md) | MCP tool surface + recipes |
| 08 | [TOKEN_ECONOMICS](./planning/atheon-enhanced/08_TOKEN_ECONOMICS.md) | The dollar numbers |
| 09 | [BENCHMARK_PLAN](./planning/atheon-enhanced/09_BENCHMARK_PLAN.md) | v1 and v2 design |
| 10 | [WHITE_PAPER_OUTLINE](./planning/atheon-enhanced/10_WHITE_PAPER_OUTLINE.md) | White paper structure |
| 11 | [GO_TO_MARKET](./planning/atheon-enhanced/11_GO_TO_MARKET.md) | Launch, distribution, growth |
| 12 | [DEMOS_AND_INFRA](./planning/atheon-enhanced/12_DEMOS_AND_INFRA.md) | Demos, CI, releases, registry |
| 13 | [UPSTREAM_CONTRIBUTION](./planning/atheon-enhanced/13_UPSTREAM_CONTRIBUTION.md) | The upstream playbook |
| 14 | [RELEASE_AND_VERSIONING](./planning/atheon-enhanced/14_RELEASE_AND_VERSIONING.md) | SemVer, cadence, goreleaser |
| 15 | [RISKS_AND_LIMITATIONS](./planning/atheon-enhanced/15_RISKS_AND_LIMITATIONS.md) | 12 risks, mitigations |
| 16 | [ACTION_PLANS](./planning/atheon-enhanced/16_ACTION_PLANS.md) | Day-by-day to-do lists |
| 17 | [QUICK_REFERENCE](./planning/atheon-enhanced/17_QUICK_REFERENCE.md) | Cheat sheet |

---

## Open questions

- **Should the fork renumber to 0.3.0?** Today it's at 1.2.0
  with a public API that is still moving. Renumbering to
  0.3.0 honestly signals "not yet 1.0." See
  [`14_RELEASE_AND_VERSIONING.md`](./planning/atheon-enhanced/14_RELEASE_AND_VERSIONING.md).
- **What's the right approach to AI-detection patterns?** The
  current 6 patterns have ~22% FP rate on benchmark v1 (TBD).
  Either tighten aggressively or deprecate the category. See
  [`06_PATTERN_LIBRARY_STRATEGY.md`](./planning/atheon-enhanced/06_PATTERN_LIBRARY_STRATEGY.md).
- **When to do Show HN, take 2?** After benchmark v2, ~end of
  quarter. See [`11_GO_TO_MARKET.md`](./planning/atheon-enhanced/11_GO_TO_MARKET.md).

---

## Dependencies

- **Upstream** `github.com/HoraDomu/Atheon` — rebased weekly
  via `sync-stable-clean.yml`.
- **Go 1.21+** toolchain.
- **Cloudflare Pages** for the demos.
- **GitHub Actions** for CI and release.

---

## What this plan does NOT cover

- **Long-term financial sustainability.** Volunteer-driven.
- **Multi-language scanner.** Engine is Go-only; patterns are
  language-agnostic.
- **Hosted SaaS version.** Local-first is the project's
  positioning.
- **Enterprise support contracts.** Out of scope.
- **Mobile or web UI for the scanner.** Out of scope.

These are deliberate omissions. The project is small and
local; adding any of them is a different project.