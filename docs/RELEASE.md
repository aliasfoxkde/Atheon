# Release Process

How Atheon ships. The goal is a release process that is
boring, repeatable, and produces artifacts that downstream
users can trust. Pair this with
[`GOVERNANCE.md`](GOVERNANCE.md) for who has merge rights
and [`SECURITY.md`](../.github/SECURITY.md) for how security
fixes are handled out-of-band.

## Cadence

- **Minor releases (1.x):** roughly every 6–8 weeks.
- **Patch releases (1.x.y):** as needed for bug fixes.
- **Security releases:** within 24–72 hours of disclosure,
  bypassing the normal cadence. See
  [`SECURITY.md`](../.github/SECURITY.md).
- **Major releases (2.0):** only when there is a documented
  breaking change. See [`MIGRATION.md`](MIGRATION.md).

There is no fixed calendar. "Ready" means: open issues and
PRs tagged for the release are closed, CI is green on
`stable`, and the release notes are drafted.

## Versioning

Atheon follows [Semantic Versioning 2.0.0](https://semver.org/):

- **MAJOR** — incompatible API changes.
- **MINOR** — new functionality, backwards-compatible.
- **PATCH** — backwards-compatible bug fixes.

The embedded pattern bundle, the public Go API, and the CLI
flags are all versioned together. A change to one is not a
hidden change to the others.

## Branch model

```
main      — tag source of truth; only fast-forwarded from stable
stable    — integration branch; PRs land here
fix/…     — bug fix branches
feat/…    — feature branches
```

- All PRs target `stable`.
- Releases are cut from `stable` once CI is green and the
  release notes are drafted.
- `main` is updated by the release script (`make release`)
  after a tag is pushed; it should never be targeted by a PR.

## Cutting a release

The release script in `scripts/release.sh` is the canonical
mechanism. It is idempotent and safe to re-run.

1. **Freeze `stable`.** No new PRs merge while the release
   is being cut. The freeze lifts as soon as the tag is
   pushed.
2. **Update `CHANGELOG.md`.** Move the `## [Unreleased]`
   section into a dated `## [1.x.y] - YYYY-MM-DD` section.
   Compare against the previous tag with
   `git log v1.x.y-1..HEAD --oneline`.
3. **Bump the version.** Update the `const Version = "…"`
   declaration in `cmd/atheon/main.go` and
   `cmd/mcp/main.go`.
4. **Run the release script.**
   ```bash
   ./scripts/release.sh 1.x.y
   ```
   The script:
   - Runs the full test suite with `-race`.
   - Runs `go vet ./...`, `staticcheck ./...`,
     `golangci-lint run`.
   - Builds binaries for all targets via GoReleaser.
   - Generates the embedded pattern bundle.
   - Creates the annotated tag `v1.x.y`.
5. **Push the tag.** `git push origin v1.x.y` triggers the
   GitHub Actions release workflow, which signs the binaries
   with cosign, attaches checksums, and publishes to the
   GitHub release.
6. **Update `main`.** `git push origin stable:main` (this is
   the one place a force-push of `main` is allowed; the
   `stable` history is fast-forwarded into `main`).
7. **Announce.** Post a short summary in the GitHub
   Discussions *Announcements* category and on the project
   Matrix/Discord channel (links in
   [`.github/SUPPORT.md`](../.github/SUPPORT.md)).

## Hotfix releases

For a security or data-loss fix that cannot wait for the
next scheduled release:

1. Branch from the most recent release tag, not from
   `stable`:
   ```bash
   git checkout -b fix/cve-XXXX v1.x.y
   ```
2. Land the fix, update `CHANGELOG.md` with a `## [1.x.y+1]`
   section, and tag.
3. Cherry-pick the fix into `stable` so the next regular
   release includes it.
4. Announce under *Security* in Discussions and follow the
   disclosure timeline in [`SECURITY.md`](../.github/SECURITY.md).

## Pre-release builds

Every commit on `stable` produces a pre-release build with
the tag `v1.x.y-pre.NNN` where `NNN` is the build counter.
These are **not** signed and are **not** announced. They
exist for maintainers and CI to smoke-test the release
script before the real tag.

## What never goes into a release

- **Untested code.** CI must be green on the exact commit
  being tagged.
- **Unreviewed PRs.** Use the standard review process even
  for "obvious" fixes; the discussion in the PR is part of
  the release record.
- **External links to the binaries in the git repo.** The
  binaries are in the GitHub release, not in the repo. The
  repo is source; the release is distribution.
- **Versioned docs in the repo.** `docs/MIGRATION.md` is the
  one place where we plan forward; everything else is written
  as the current truth.

## After the release

- **Close the milestone.** Any issues or PRs still in the
  milestone move to the next one.
- **Triage the backlog.** Anything tagged for the next
  release is reviewed and either merged, deferred with a
  reason, or closed.
- **Update `ROADMAP.md`** if the release changed the
  priorities recorded there.

## Release checklist

- [ ] `CHANGELOG.md` updated, dated, and complete.
- [ ] Version bumped in `cmd/atheon/main.go` and
      `cmd/mcp/main.go`.
- [ ] `./scripts/release.sh 1.x.y` runs clean.
- [ ] Tag pushed and GitHub release published.
- [ ] `main` updated from `stable`.
- [ ] Announcements posted.
- [ ] Milestone closed.
- [ ] `ROADMAP.md` reviewed.
