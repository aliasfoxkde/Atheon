# Contributing

Thanks for your interest in Atheon-Enhanced!

Please see [`.github/CONTRIBUTING.md`](.github/CONTRIBUTING.md) for the full
contribution guide, including:

- How to add a new pattern (YAML file in `community/<category>/`)
- Development setup and local testing
- PR review process and CI gates

## First-Time Setup

Run these two commands once after cloning:

```bash
make setup   # Install pre-commit hooks (gofmt, govet, tests, bundle rebuild)
make init    # Configure the commit message template
```

## Commit Message Format

All commits follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <short description>

[optional body]

[optional footer]
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`, `ci`, `build`, `perf`, `release`, `revert`

Rules:
- Subject line ≤ 72 characters
- Total message ≤ 800 characters
- Use imperative mood ("add feature" not "added feature")

The commit template (installed via `make init`) provides guidance automatically.

Quick links:

- [Pattern authoring guide](docs/patterns/contributing-patterns.md)
- [Pattern YAML format](docs/PATTERN_FORMAT.md)
- [Development setup](docs/development/SETUP.md)
- [Owner checklist](docs/OWNER_CHECKLIST.md)
- [Code of conduct](CODE_OF_CONDUCT.md)
