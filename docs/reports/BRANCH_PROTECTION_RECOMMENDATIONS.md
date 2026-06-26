# Branch Protection Recommendations — 2026-06-24

## Background

Audit confirmed that PRs **#57** and **#58** (the foundational refactors that
introduced the consolidated CI/CD pipeline) were merged to `main` with **zero**
GitHub Actions workflow runs. The repository's `main` branch had no branch
protection rules, so contributors were able to bypass all required status checks.

This file captures the recommended `main` branch protection configuration so the
findings can be re-applied as a single admin action.

## Required Status Checks

After the next green run of the consolidated CI, configure the following as
required status checks on `main`:

| Check name                              | Workflow                |
| --------------------------------------- | ----------------------- |
| `CI / Test (Go 1.21)`                   | ci.yml                  |
| `CI / Test (Go 1.22)`                   | ci.yml                  |
| `CI / Test (Go 1.23)`                   | ci.yml                  |
| `CI / Test (Go 1.24)`                   | ci.yml                  |
| `CI / Lint`                             | ci.yml                  |
| `CI / Build (ubuntu-latest)`            | ci.yml                  |
| `CI / Build (macos-latest)`             | ci.yml                  |
| `CI / Build (windows-latest)`           | ci.yml                  |
| `CI / Integration Tests`                | ci.yml                  |
| `CI / Performance Benchmarks`           | ci.yml                  |
| `CI / Documentation Check`              | ci.yml                  |
| `CI / Test Results & Coverage`          | ci.yml                  |
| `Security / CodeQL (Go)`                | security.yml            |
| `Security / Self-Scan (secrets — blocking)` | security.yml        |
| `Security / Security Anti-Patterns`     | security.yml            |
| `Security / Self-Scan (code-quality — informational)` | security.yml |
| `Security / Go Vulnerability Check`     | security.yml            |

## API Payload (preview — do not apply until CI is green)

```json
{
  "required_status_checks": {
    "strict": true,
    "contexts": [
      "CI / Test (Go 1.21)",
      "CI / Test (Go 1.22)",
      "CI / Test (Go 1.23)",
      "CI / Test (Go 1.24)",
      "CI / Lint",
      "CI / Build (ubuntu-latest)",
      "CI / Build (macos-latest)",
      "CI / Build (windows-latest)",
      "CI / Integration Tests",
      "CI / Performance Benchmarks",
      "CI / Documentation Check",
      "CI / Test Results & Coverage",
      "Security / CodeQL (Go)",
      "Security / Self-Scan (secrets — blocking)",
      "Security / Security Anti-Patterns",
      "Security / Self-Scan (code-quality — informational)",
      "Security / Go Vulnerability Check"
    ]
  },
  "enforce_admins": true,
  "required_pull_request_reviews": {
    "dismissal_restrictions": {},
    "dismiss_stale_reviews": true,
    "require_code_owner_reviews": false,
    "required_approving_review_count": 1
  },
  "restrictions": null,
  "required_linear_history": true,
  "allow_force_pushes": false,
  "allow_deletions": false,
  "block_creations": false,
  "required_conversation_resolution": true
}
```

## How to Apply

1. Wait for the consolidated CI to produce a green run on the latest `main` HEAD.
2. Open `https://github.com/aliasfoxkde/Atheon-Enhanced/settings/branches` (requires admin).
3. Edit the `main` branch protection rule and paste the contexts above.
4. Enable:
   - "Require status checks to pass before merging"
   - "Require branches to be up to date before merging" (matches `strict: true`)
   - "Require conversation resolution before merging"
   - "Include administrators" (`enforce_admins: true`)

Or apply via the GitHub API:

```bash
gh api \
  --method PUT \
  -H "Accept: application/vnd.github+json" \
  /repos/aliasfoxkde/Atheon-Enhanced/branches/main/protection \
  --input branch-protection.json
```

## What This Audit Does Not Cover

- **CODEOWNERS file** is not yet defined for the repository. Adding a `CODEOWNERS`
  file at `.github/CODEOWNERS` would unlock per-path required reviewers and is a
  natural follow-up.
- **Secret scanning push protection** is not yet enabled. Recommended for
  repositories with workflow secrets (`GITHUB_TOKEN`, `CODECOV_TOKEN`, `GH_PAT`).
- **Dependabot** is not configured. Recommended for `go.mod` and GitHub Actions
  version bumps.
