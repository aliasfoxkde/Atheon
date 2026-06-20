## Summary

<!-- One paragraph: what this PR does and why. Reference any
related issue using the closing keywords (`Closes #123`,
`Fixes #456`) so the PR merges with the issue. -->

## Changes

<!-- Bullet list of the user-visible changes. Mention any new flags,
defaults that changed, or breaking behaviour. -->

## Test plan

<!-- How did you verify this works? Checklist the appropriate boxes. -->

- [ ] I added or updated unit tests for the change.
- [ ] I added or updated integration tests where applicable.
- [ ] I ran `go test ./... -race` locally and it passed.
- [ ] I checked coverage with `go test ./... -cover` and it did not
      regress below the threshold in `ci.yml`.
- [ ] I ran `go vet ./...` and `gofmt -l .` (both clean).
- [ ] I updated the docs (`docs/`, `README.md`) if the change is
      user-visible.

## Quality gates

- [ ] `golangci-lint run ./...` is clean (or justified warnings).
- [ ] `staticcheck ./...` is clean.
- [ ] Pre-commit hook (`./scripts/install-hooks.sh` then commit)
      passes locally.

## Risk and rollout

<!-- Pick the level that matches the change and describe rollback
strategy if needed. -->

- **Risk:** <!-- low / medium / high -->
- **Rollback:** <!-- how to revert safely if something goes wrong -->
- **Feature flag needed?** <!-- yes / no — if yes, which one -->

## Documentation

<!-- Link or list the doc pages you touched. -->

- [ ] `docs/CHANGELOG.md` updated under "Unreleased"
- [ ] `docs/ROADMAP.md` updated if a roadmap item moved
- [ ] `docs/API.md` updated if the public API changed
- [ ] `README.md` updated if the change is user-visible

## Screenshots / output

<!-- Paste command output, JSON snippets, or screenshots that
demonstrate the change. Drop the section if not applicable. -->

```text

```

## Checklist

- [ ] I have read [CONTRIBUTING.md](../CONTRIBUTING.md).
- [ ] My commits follow the project's commit-message convention
      (`type(scope): subject`).
- [ ] I have split unrelated changes into separate PRs.
- [ ] I am the author of every commit, or the commit message
      attributes the original author (Co-authored-by:).
