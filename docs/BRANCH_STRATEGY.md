# Branch Strategy

This document describes the branching strategy for Atheon.

## Main Branches

### main
The main production branch. Contains stable, production-ready code.

### stable/clean
A cleaned stable branch synced from upstream HoraDomu/Atheon. Used for stable releases.

### dev/full-feature
Development branch for comprehensive feature work and integration testing.

## Branch Naming Conventions

- `main` - Production branch
- `stable/clean` - Stable sync branch
- `dev/full-feature` - Feature development branch
- `infra/**` - Infrastructure updates
- `feat/*` - New features
- `fix/*` - Bug fixes
- `docs/*` - Documentation updates
- `refactor/*` - Code refactoring
- `test/*` - Test improvements

## Workflow

1. Feature work happens on `feat/*` or `dev/full-feature` branches
2. Changes are submitted via pull requests
3. PRs require review before merging to `main`
4. `stable/clean` syncs with upstream HoraDomu/Atheon main branch
