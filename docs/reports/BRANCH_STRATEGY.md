# Branch Strategy — Deprecated Location

> **This document is deprecated.** The canonical branch strategy has moved to
> [../BRANCH_STRATEGY.md](../BRANCH_STRATEGY.md).
>
> Older, longer-form branch-strategy content that previously lived here is preserved below
> for historical reference only. Do not edit — make changes in
> [../BRANCH_STRATEGY.md](../BRANCH_STRATEGY.md) instead.

---

# Branch Strategy Documentation (Historical)

## 🎯 Overview

The aliasfoxkde/Atheon repository uses a systematic branch strategy designed to maintain high-quality releases while staying synchronized with upstream HoraDomu/Atheon.

## 📋 Core Branches

### **`stable/clean`** (Upstream Tracking Branch)
**Purpose**: Source of truth that tracks upstream HoraDomu/Atheon:main
- **Update Strategy**: Automatic sync via GitHub Actions daily
- **Usage**: Baseline for all development, reference for upstream changes
- **Protection**: Protected branch, only maintainers can push

### **`main`** (Production Build)
**Purpose**: Production-ready build with all validated enhancements
- **Update Strategy**: Merge validated feature PRs + periodic stable/clean sync
- **Usage**: User-facing installation via `go install github.com/aliasfoxkde/Atheon`
- **Protection**: Protected branch, requires PR review and passing CI

### **`dev/full-feature`** (Development Branch)
**Purpose**: Development branch with ALL patterns enabled for comprehensive testing
- **Update Strategy**: Continuous integration of feature branches
- **Usage**: Internal testing, validation, pattern development

## 🔄 Development Workflow

(See canonical doc at [../BRANCH_STRATEGY.md](../BRANCH_STRATEGY.md).)

## Configuration Profiles

Configuration profiles live at `config/profiles/`:

- `development.json` — Development environment configuration
- `mcp-integration.json` — MCP server integration settings
- `pipeline.json` — CI/CD pipeline configuration
- `production.json` — Production environment settings

For the up-to-date branch strategy and development workflow, see
[../BRANCH_STRATEGY.md](../BRANCH_STRATEGY.md).
