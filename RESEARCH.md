# Research — Atheon-Enhanced

**Last Updated:** 2026-06-20

This file exists to satisfy the global harness's scaffolding
requirement. The substantive research, design, and planning
artifacts live in [`docs/`](./docs/). See in particular:

- [`docs/PLAN.md`](./docs/PLAN.md) — high-level plan and the
  four-phase roadmap
- [`docs/GOAL_ROADMAP.md`](./docs/GOAL_ROADMAP.md) — /goal
  production-quality sweep
- [`docs/DESIGN.md`](./docs/DESIGN.md) — design rationale
  (goals, non-goals, key decisions)
- [`docs/ARCHITECTURE.md`](./docs/ARCHITECTURE.md) — system
  architecture
- [`docs/STANDARDS.md`](./docs/STANDARDS.md) — engineering
  standards
- [`docs/ROADMAP.md`](./docs/ROADMAP.md) — long-term themes
- [`docs/audits/DEAD_CODE_AUDIT.md`](./docs/audits/DEAD_CODE_AUDIT.md) —
  the latest code-quality audit

## Project at a glance

- **What:** Local-first pattern-matching engine + MCP server for
  scanning code, text, and URLs for secrets, PII, security
  issues, code smells, and AI-tells.
- **Stack:** Go 1.21+, single binary (zero runtime deps), MCP
  over stdio, JSON-RPC.
- **Patterns:** 179 across 17 categories (validated metric —
  see `docs/memories/VALIDATED_METRICS.md`).
- **Coverage:** 98.7% (target ≥95%).

## Constraints

- Cross-platform binaries (Linux/macOS/Windows, Intel + ARM).
- Zero network calls for the local scanner.
- Bundle is embedded; downloaded override lives in
  `~/.atheon/`.
- MIT license with additional terms (mirrors upstream).
