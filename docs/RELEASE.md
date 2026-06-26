# Release Runbook

**Last Updated**: 2026-06-25
**Audience**: Maintainers cutting a release. Contributors do not need to read this.

A release is: (1) cut a version tag, (2) ship Go binaries + the pattern bundle, (3) publish a GitHub release. The full sequence is below.

---

## Tag format

```
v0.YY.MM.DD[-rcN]
```

- **`0.YY.MM.DD`** — calver aligned with the date the release is cut. The leading `0` reserves the `1.x.x` slot for a future semver-breaking release if we ever decide to ship one.
- **`-rcN`** — release candidate. Use `-rc1` for the first pre-release; bump `N` for subsequent fixes. RCs publish binaries and the bundle, but the GitHub release is marked "Pre-release".

Examples:
- `v0.26.06.25` — full release cut on 2026-06-25.
- `v0.26.07.01-rc1` — first RC for the 2026-07-01 slot.
- `v0.26.07.01` — promote the RC after the soak period.

> The CHANGELOG uses the same numeric prefix without the leading `v`: `[0.6.0]` in CHANGELOG corresponds to tag `v0.26.06.25`.

---

## Pre-release checklist

Run these in order. Stop and investigate if any step fails.

```bash
# 1. Working tree clean, on main, up to date with origin.
git checkout main
git pull --rebase origin main
git status

# 2. No uncommitted patterns.bundle drift.
./scripts/pattern-count.sh --total
./atheon list | tail -1
# Counts MUST match. If they don't, regenerate:
go run ./bundler
git add core/patterns.bundle
git commit -m "chore(bundle): regenerate"

# 3. All tests pass with -race to catch any deferred concurrency issues
#    (e.g. the pattern_state mutex work, if not yet merged).
go test ./... -p 1 -race -timeout 15m

# 4. gofmt + goimports clean.
gofmt -l .   # must print nothing
goimports -l .   # must print nothing

# 5. CI is green on the current main.
gh workflow view ci --yaml >/dev/null
gh run list --workflow=ci --limit=1 --json conclusion --jq '.[0].conclusion'
# expect: "success"

# 6. Lint clean.
go vet ./...
go install honnef.co/go/tools/cmd/staticcheck@v0.6.1 && staticcheck ./...

# 7. Coverage threshold holds.
go test ./... -p 1 -coverprofile=/tmp/cov.out
go tool cover -func=/tmp/cov.out | grep total
# expect: ≥ vars.COVERAGE_THRESHOLD (default 70%)

# 8. Smoke scan a real repo to confirm runtime behavior.
git clone --depth 1 https://github.com/aliasfoxkde/Atheon-Enhanced /tmp/atheon-self-scan
./atheon --sarif /tmp/atheon-self-scan > /tmp/self-scan.sarif
# Findings are expected — this is a *scanner*, not the scanner's source.
# What matters is that the exit code reflects findings + errors:
echo "exit: $?"   # expect: 1
```

If all eight pass, proceed. If any fails, **do not** cut a tag — fix the regression first, merge the fix, then re-run the checklist from step 5.

---

## Cut the release

### Step 1 — Update CHANGELOG.md

Open `CHANGELOG.md` and:

1. Move the current `[Unreleased]` section's content into a new dated section with the version matching the tag, e.g. `[0.26.06.25] - 2026-06-25`.
2. Leave a fresh empty `[Unreleased]` section at the top for the next cycle.
3. Add a comparison link at the bottom if this is a notable release: `[0.26.06.25]: https://github.com/aliasfoxkde/Atheon-Enhanced/releases/tag/v0.26.06.25`.

Keep version section ordering: `[Unreleased]` → most recent dated → older dated → initial release.

### Step 2 — Commit the CHANGELOG update

```bash
git add CHANGELOG.md
git commit -m "chore(release): cut v0.26.06.25"
```

This commit goes directly on `main` — do not open a PR for the CHANGELOG alone.

### Step 3 — Tag

```bash
git tag -s v0.26.06.25 -m "Release v0.26.06.25"
```

`-s` (GPG-sign) is required for tag verification by `go install`. If your local GPG key isn't set up, use `-a` and document why in the PR description.

```bash
git push origin main
git push origin v0.26.06.25
```

### Step 4 — Let CI publish

`.github/workflows/release.yml` is triggered by the `v*` tag. It will:

1. Build binaries for linux/macos/windows × amd64/arm64.
2. Run GoReleaser to attach binaries + the generated `patterns.bundle` to a draft GitHub release.
3. Sign the artifacts with cosign if `COSIGN_PRIVATE_KEY` is configured in the `release` environment.

Watch the run:

```bash
gh run list --workflow=release --limit=1
gh run watch <run-id>
```

### Step 5 — Review and publish the draft

```bash
gh release view v0.26.06.25 --web
```

Verify in the draft:

- All six binary archives attached (3 OSes × 2 arches).
- `patterns.bundle` attached at the root.
- `SHA256SUMS` file present.
- Auto-generated notes match the CHANGELOG entry.

If anything is missing, **delete the draft and re-run the workflow** with a `-rcN` tag first; do not edit the draft by hand.

When it looks right: click **Publish release**. There is no approval gate beyond the `release` GitHub Environment (which is currently set to no required reviewers — see [TASKS.md](./TASKS.md) for related follow-ups).

---

## Hotfix workflow

A hotfix is a release that needs to land between scheduled releases, e.g. a false-positive regression reported by a user.

1. Branch from the latest released tag, **not** from `main`:
   ```bash
   git fetch --tags
   git checkout -b hotfix/short-desc v0.26.06.25
   ```
2. Make the fix. Add a regression test.
3. Open a PR targeting `main` (not `stable/clean`). The PR body must say `HOTFIX` so reviewers know it's urgent.
4. After approval + merge, immediately cut a new `-rc1` tag and follow the [Cut the release](#cut-the-release) flow from Step 3.
5. Skip the RC soak period for confirmed-hotfix releases — `-rc1` may be promoted to a full release within the same day if the fix is verified by the reporter.

The CHANGELOG entry for a hotfix is dated to the day it's published, not the day the fix was authored.

---

## Bundle regeneration

The bundle is rebuilt from `community/**/*.yaml` via the `bundler` package:

```bash
go run ./bundler
```

This rewrites `core/patterns.bundle`. **Always commit the regenerated bundle alongside any pattern YAML change.** The CI check "Check bundle freshness vs source" will block any PR that changes YAML but doesn't regenerate the bundle — that's by design.

To bump the bundle format version (when adding a new field to `PatternDef`):

1. Add the field to `PatternDef` in `core/bundle.go` with `omitempty` if it's additive.
2. Update the decoder in `loadBundle` to handle both old (missing field) and new (field present) bundles.
3. Update `bundler/main.go` to emit the new field.
4. Regenerate: `go run ./bundler`.
5. Document the version bump in CHANGELOG under `[Unreleased]`.

---

## ldflags

The CLI's `version` variable (`cmd/atheon/main.go:17`) is injected at build time:

```go
var version = "dev"
```

GoReleaser sets this via:

```yaml
builds:
  - ldflags:
      - -X main.version={{.Version}}
```

Local builds keep `version="dev"`. CI builds of a tag get `version="0.26.06.25"`. The `--version` flag prints it:

```
$ atheon --version
atheon 0.26.06.25
```

If you ever add another package that needs an ldflag variable, follow the same `-X package.path.VarName` convention and document it here.

---

## Troubleshooting

### `gh release view` says release is not found after publish

GH indexes new releases with a 1-2 minute delay. Wait and retry, or check `https://github.com/aliasfoxkde/Atheon-Enhanced/releases` directly.

### GoReleaser skipped a build target

Check `.goreleaser.yml`. The matrix is `goos: [linux, darwin, windows]` × `goarch: [amd64, arm64]`. A missing archive usually means an `archives:` block filter dropped it. Run `goreleaser build --snapshot --clean` locally to debug without publishing.

### Tag was created on the wrong commit

Tags cannot be moved without force-push, which is **forbidden** on this repo. The fix is to delete the tag locally + remotely and re-tag:

```bash
git tag -d v0.26.06.25
git push origin :refs/tags/v0.26.06.25
git tag -s v0.26.06.25   # on the correct commit
git push origin v0.26.06.25
```

If the release was already published, treat it as a hotfix: bump to `-rc1`, fix, then cut the corrected tag.

### Coverage dropped below threshold

The CI gate is `vars.COVERAGE_THRESHOLD` (default 70). If a deliberate drop is needed:

1. Update the var via `gh variable set COVERAGE_THRESHOLD <new>`.
2. Document the reason in the PR body that triggered the drop.
3. Plan a follow-up wave to restore coverage.

Do **not** disable the gate or lower the threshold in `ci.yml` to dodge a real coverage regression.

---

## References

- [PLAN.md](./PLAN.md) — overall project plan and hardening wave timeline
- [TASKS.md](./TASKS.md) — task ledger (open and completed work)
- [BRANCH_STRATEGY.md](./BRANCH_STRATEGY.md) — branch layout, protection rules
- [CHANGELOG.md](../CHANGELOG.md) — release notes
- [`.github/workflows/release.yml`](../.github/workflows/release.yml) — the workflow triggered by tag push
