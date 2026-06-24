# ADR 0002: Five-workflow CI surface

- **Status**: Accepted
- **Date**: 2026-06-23
- **Deciders**: aliasfoxkde
- **Supersedes**: prior 10-workflow configuration

## Context

The repository had 10 GitHub Actions workflows with substantial overlap
before 2026-06-23:

- Three workflows ran variants of "test + lint" (ci.yml,
  comprehensive-ci.yml, quality-assurance.yml).
- Two ran variants of CodeQL (codeql.yml, security-scanning.yml).
- Two ran variants of the Atheon self-scan (security-scanning.yml,
  self-scan.yml).
- Release split across two files (publish.yml + scheduled-release.yml).

Each duplicate triggered its own runner, its own checkout, its own
`go build`, and its own upload. The cumulative runner-minute cost was
~3× the actual work, and the noise of overlapping status checks made
PR review slower.

## Decision

**Consolidate to five single-responsibility workflows:**

| File | Trigger | Purpose |
|------|---------|---------|
| `ci.yml` | push to main, PR | test (4 Go versions), lint (vet/staticcheck/golangci-lint/gofmt), build (3 OS), integration, JUnit + Codecov |
| `security.yml` | push to main, PR, weekly Mon 06:00 UTC | CodeQL, Atheon self-scan (secrets blocking + quality informational), govulncheck |
| `release.yml` | tag push, schedule (10th/21st), workflow_dispatch | version derive + tag push + GoReleaser |
| `sync.yml` | workflow_dispatch only | merge upstream HoraDomu/Atheon main into `stable/clean` |
| `auto-merge.yml` | PR open/sync | enable auto-merge on green CI |

Removed (obsolete duplicates): codeql.yml, comprehensive-ci.yml,
publish.yml, quality-assurance.yml, scheduled-release.yml,
security-scanning.yml, self-scan.yml.

## Consequences

**Positive**

- Single source of truth per concern. Test failures point to one
  workflow; CodeQL alerts point to one workflow; release builds happen
  in one place.
- Runner-minute cost drops by ~60% (1 checkout per concern instead of 3).
- PR status checks collapse from ~12 status contexts to 5.
- Schedule-based runs (CodeQL weekly, release 10th/21st) are now
  co-located with their ad-hoc counterparts, easier to reason about.

**Negative**

- Larger individual workflows are harder to read end-to-end. We
  mitigate with banner-style section headers (`# ──`) so the job
  boundaries are visible.
- Removing `auto-merge.yml` was considered (GitHub's repo setting
  offers native auto-merge) but kept — the workflow version handles
  rebase-vs-squash selection and conflict reporting, which the setting
  alone does not.

**Neutral**

- All actions remain SHA-pinned. All shell commands still pass
  `-p 1` to `go test` (mandatory because `core/` has package-level
  state in `init()` that breaks under parallel package execution).