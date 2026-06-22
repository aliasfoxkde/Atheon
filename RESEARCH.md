# Research: Upstream PR Submissions

**Date:** 2026-06-22
**Purpose:** Prepare clean PRs for HoraDomu/Atheon upstream

## Issues to Address

1. **#155** - Duplicate enabled check in SetActiveCategories (bug fix)
2. **#158** - Non-deterministic list output (improvement)
3. **#157** - ScanDir error swallowing (error handling)
4. **#156** - JSON flag position (CLI)

## Constraints

- Keep PRs small (1 issue per PR)
- Minimal files changed
- Squash commits for clean history
- Do NOT push to upstream directly
- Submit via PRs from clean feature branches

## Approach

Create feature branches prefixed with `pi/{issue-number}-{short-name}` from stable/clean
