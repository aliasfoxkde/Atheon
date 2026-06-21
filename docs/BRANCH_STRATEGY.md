# Branch Strategy

This document describes the branch strategy for Atheon-Enhanced.

## Main Branches

- `main` - Production-ready code, always deployable
- `stable/clean` - Stable integration branch synced with upstream HoraDomu/Atheon
- `dev/full-feature` - Development branch for major features

## Branch Naming Conventions

- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `test/` - Test improvements
- `refactor/` - Code refactoring
- `infra/` - Infrastructure changes

## Workflow

1. Create a feature branch from `main`
2. Make changes and commit
3. Push to origin and create a PR
4. After review, merge to `main`
5. `stable/clean` syncs with upstream HoraDomu/Atheon

## Protected Branches

The following branches are protected and require PRs:
- `main`
- `stable/clean`
- `dev/full-feature`
