# Owner Checklist & Recommendations

Comprehensive list of setup tasks, pending work, and recommendations for `aliasfoxkde/Atheon-Enhanced`.
Last updated: 2026-06-23.

---

## IMMEDIATE: GitHub Actions Secrets

No secrets are currently configured. Several CI features are silently degraded until these are added.

**How:** GitHub → Settings → Secrets and variables → Actions → New repository secret

| Secret | Value | Effect if missing |
|--------|-------|-------------------|
| `CODECOV_TOKEN` | From codecov.io (see below) | Coverage uploads use tokenless mode — unreliable, no PR comments |

---

## IMMEDIATE: Codecov Setup

1. Go to [codecov.io](https://codecov.io) and sign in with GitHub
2. Click **+ Add a repository** → select `aliasfoxkde/Atheon-Enhanced`
3. Codecov will display an upload token — copy it
4. Add it as `CODECOV_TOKEN` in GitHub Actions secrets (above)
5. Install the [Codecov GitHub App](https://github.com/apps/codecov) on your account so the `codecov-commenter` bot can post PR comments

**After setup:** The live coverage badge in README will populate, and future PRs will get automatic patch coverage diff comments.

> **Recommendation:** Keep the `codecov.yml` patch target at 80% and project target at 90%. These are ambitious but achievable — the codebase is already at 97%. Lowering them later is easy; letting coverage drift is hard to reverse.

---

## IMMEDIATE: Open Pull Requests to Merge

Both PRs are CI-tested and ready.

| PR | Branch | What it does |
|----|--------|-------------|
| [#53](https://github.com/aliasfoxkde/Atheon-Enhanced/pull/53) | `feat/codecov-integration` | Codecov v5.4.2, live README badge, `codecov.yml` config |
| [#54](https://github.com/aliasfoxkde/Atheon-Enhanced/pull/54) | `fix/funding-yml-github-sponsors` | Removes unenrolled `github: aliasfoxkde` from FUNDING.yml |

> **Recommendation:** Merge #54 first (1-line fix, unblocks the Sponsor button error). Then merge #53 after adding `CODECOV_TOKEN` so the badge works immediately on merge.

---

## GitHub Repository Settings

### Pull Requests (Settings → General)

| Setting | Recommended | Why |
|---------|------------|-----|
| Automatically delete head branches | **Enable** | Prevents branch accumulation after merge; `pr/*` branches are excluded from auto-delete because they're not PR head branches on this repo |
| Allow squash merging | Enable (default) | Clean linear history |
| Allow merge commits | Disable (optional) | Keeps history linear |
| Allow rebase merging | Enable (optional) | Useful for simple fixes |

### Branch Protection (Settings → Branches)

**`stable/clean`** is currently protected and cannot be deleted. This branch has no active purpose.

To remove it:
1. Settings → Branches → find the `stable/clean` protection rule → delete it
2. Then: `git push origin --delete stable/clean`

> **Recommendation:** Add a branch protection rule for `main` if one doesn't exist: require PR reviews, require status checks (CI, self-scan/secrets), and disable force-push. This prevents accidental direct pushes.

### GitHub Sponsors (Settings → Sponsor this project)

`aliasfoxkde` is not yet enrolled in GitHub Sponsors, which causes a FUNDING.yml parse error. To add the GitHub Sponsors button:

1. Go to [github.com/sponsors](https://github.com/sponsors) and apply
2. Once approved, add `github: aliasfoxkde` back to `.github/FUNDING.yml`

The Ko-fi and Patreon buttons are working now (after PR #54 merges).

---

## Upstream PRs (HoraDomu/Atheon)

These are PRs submitted from this repo upstream. The `pr/*` branches **must be kept** — deleting them auto-closes the upstream PRs.

| Upstream PR | Branch | Status |
|-------------|--------|--------|
| [#176](https://github.com/HoraDomu/Atheon/pull/176) | `pr/158-deterministic-list` | Open |
| [#175](https://github.com/HoraDomu/Atheon/pull/175) | `pr/156-json-flag` | Open |
| [#173](https://github.com/HoraDomu/Atheon/pull/173) | `pr/146-147-148-perf` | Open |
| [#172](https://github.com/HoraDomu/Atheon/pull/172) | `pr/149-patterns-expansion` | Open |

Also present (check status): `pr/155-157-scan-errors`, `pr/177-fix-build-and-ci-lint-schema-timeout`, `pr/177-v2-clean`.

> **Recommendation:** Periodically check upstream PR status. When upstream merges one, sync it back: `git fetch upstream && git merge upstream/main`. This keeps the repos aligned and avoids growing divergence.

---

## Code Quality: main.go Coverage Gap

A previous Codecov report flagged `main.go` at 5.88% patch coverage (16 missing lines). This is from the CLI-layer functions (`printSARIFFindings`, `buildSARIFRules`, `buildSARIFResults`, `cmdList`, `printHelp`, `redact`, `formatBytes`) which are exercised by integration tests but not unit tests with coverage instrumentation.

**To fix:** Add a `cmd/atheon/main_test.go` that tests the CLI output functions directly — invoke `printFindings`, `printJSONFindings`, `printSARIFFindings` with captured stdout. This is straightforward table-driven test work.

> **Recommendation:** Target this in the next sprint. Patch coverage below 80% will block PRs once Codecov is properly configured. Addressing it proactively avoids pressure later.

---

## Stale Branch Cleanup

After enabling "Automatically delete head branches," future merged PRs clean themselves up. For existing stale branches:

**Branches safe to delete** (no open upstream PR association):
- `stable/clean` — protected; remove protection first (see above), then delete
- `pr/177-v2-clean` — verify no open upstream PR, then delete
- `pr/177-fix-build-and-ci-lint-schema-timeout` — verify no open upstream PR, then delete

**To verify before deleting any `pr/*` branch:**
```bash
gh api "repos/HoraDomu/Atheon/pulls?state=open&per_page=50" -q '.[] | .head.ref'
```
Cross-reference against any branch you plan to delete. If the branch name appears, do not delete it.

---

## Automation Gaps

### Dependabot

No Dependabot configuration exists. GitHub Actions and Go modules are not automatically updated.

Create `.github/dependabot.yml`:

```yaml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 5

  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10
```

> **Recommendation:** Add this. Actions pinned to hashes (as this repo does) won't auto-update without Dependabot — security patches to actions like `checkout` and `setup-go` will be missed.

### goreleaser

Binary releases are currently manual (the scheduled release workflow creates GitHub releases but does not attach compiled binaries). Adding goreleaser would produce signed, checksummed binaries for linux/mac/windows/arm automatically on each release.

Create `.goreleaser.yml` at the repo root and update the release workflow to call `goreleaser release`. The goreleaser action is free for open-source projects.

> **Recommendation:** Medium priority. Users can build from source with `go build`, but binary releases significantly lower the barrier for non-Go users.

### SECURITY.md

No responsible disclosure policy exists. GitHub recommends this file in the repo root (or `docs/` or `.github/`).

> **Recommendation:** Add a short `SECURITY.md` with: supported versions, how to report a vulnerability (private disclosure email or GitHub private advisory), and expected response time. This is a one-time 15-minute task.

---

## Pattern Library Roadmap

Current count: 255 patterns across all categories.

| Category | Estimated count | Growth opportunity |
|----------|-----------------|--------------------|
| secrets | ~80 | Cloud provider tokens, new SaaS APIs |
| pii | ~40 | Regional ID formats (EU, AU, CA) |
| code-quality | ~60 | Language-specific anti-patterns |
| devops | ~50 | Kubernetes, Terraform, GitHub Actions misconfigs |
| ai-detection | ~25 | Emerging AI prompt injection patterns |

> **Recommendation:** Use real-world project scans (Atheon on your own codebases) as the primary signal for which patterns have the most false positives. The self-scan CI loop will surface these automatically. Prioritize tightening over expanding — a pattern that fires accurately on 5 cases is more valuable than one that fires noisily on 50.

---

## Documentation Gaps

| Missing doc | Priority | Notes |
|-------------|---------|-------|
| `docs/patterns/contributing-patterns.md` | High | End-to-end guide: YAML schema, test cases, FP risk, PR process |
| `SECURITY.md` | High | Responsible disclosure policy |
| `docs/integrations/pre-commit.md` | Medium | How to install the pre-commit hook from `.hooks/` |
| `docs/integrations/vscode.md` | Low | Using Atheon with VS Code tasks or the MCP server |
| `docs/api/mcp.md` | Medium | MCP server endpoints, message format, client examples |

---

## Recurring Maintenance

| Task | Cadence | How |
|------|---------|-----|
| Sync upstream (`HoraDomu/Atheon`) | Weekly or when upstream merges land | `git fetch upstream && git merge upstream/main` |
| Check upstream PRs still open | Before deleting any `pr/*` branch | `gh api "repos/HoraDomu/Atheon/pulls?state=open" -q '.[].head.ref'` |
| Review Codecov patch coverage on each PR | Per PR | Codecov bot will post automatically once configured |
| Pattern quality review | Monthly | Run `bash .github/scripts/self-scan.sh --quality` against a diverse set of real-world repos |
| Dependabot PRs | Weekly | Merge Go module + Actions updates promptly |
