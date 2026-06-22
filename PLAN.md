# Plan: Upstream PR Submissions

**Date:** 2026-06-22

## Goal
Create clean PRs for upstream HoraDomu/Atheon issues using `pi/{number}-{name}` branch naming

## PRs to Create

1. **pi/155-dup-enabled-check** - Remove duplicate enabled check in SetActiveCategories
2. **pi/158-deterministic-list** - Sort list output for deterministic results
3. **pi/157-scandir-error-propagation** - Propagate per-file read errors
4. **pi/156-json-flag-position** - Allow `--json` in any argument position

## Steps

1. Create root docs (RESEARCH.md, PLAN.md, TASKS.md, PROGRESS.md) ✓ (in progress)
2. Create branches from stable/clean
3. Apply fixes
4. Push to aliasfoxkde/Atheon-Enhanced fork
5. Open PRs

## Constraints
- NEVER push to upstream HoraDomu/Atheon
- Only push to aliasfoxkde/Atheon-Enhanced
- One issue per PR
