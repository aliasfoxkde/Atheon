# Upgrade Guide

How to upgrade an existing Atheon installation or an existing
project that uses Atheon as a library.

## Picking a version

The project follows [Semantic Versioning](https://semver.org/):

- **Major** (1.0 → 2.0) — breaking changes to the CLI flags
  or the public Go API. Read [`MIGRATION.md`](MIGRATION.md)
  before upgrading.
- **Minor** (1.0 → 1.1) — backward-compatible additions: new
  patterns, new flags, new library functions. Safe to upgrade.
- **Patch** (1.0.0 → 1.0.1) — bug fixes only. Always safe.

## Upgrading the binary

### From a downloaded release

```sh
# 1. Note the version you're replacing
atheon --version

# 2. Download the new release for your OS/arch
curl -L https://github.com/.../releases/latest/download/atheon_$(uname -s)_$(uname -m).tar.gz | tar xz

# 3. Replace the installed binary
sudo mv atheon /usr/local/bin/atheon

# 4. Verify
atheon --version
```

### From `go install`

```sh
go install github.com/.../atheon/cmd/atheon@latest
```

The `@latest` resolves to the newest tagged release. To pin to a
specific minor:

```sh
go install github.com/.../atheon/cmd/atheon@v1.1
```

### Refresh the pattern bundle

After upgrading the binary, refresh the on-disk bundle so you
have the latest patterns:

```sh
atheon update
```

This step is optional — the binary ships with the bundle it was
built with — but recommended when the binary is more than a
month old.

## Upgrading a project that uses `core/` as a library

In `go.mod`:

```diff
-require github.com/.../atheon v1.0.0
+require github.com/.../atheon v1.1.0
```

Then:

```sh
go mod tidy
go test ./...
```

Read the [`CHANGELOG.md`](CHANGELOG.md) section for the new
version. Pay attention to the **Breaking** tag.

## Upgrading CI

If CI installs Atheon via `go install`, bump the version:

```yaml
- name: Install Atheon
  run: go install github.com/.../atheon/cmd/atheon@v1.1
```

If CI downloads a binary, update the URL:

```yaml
- name: Download Atheon
  run: |
    curl -L -o atheon https://github.com/.../releases/download/v1.1.0/atheon_$(uname -s)_$(uname -m)
    chmod +x atheon
```

## Rolling back

If the new version breaks your pipeline, downgrade:

```sh
go install github.com/.../atheon/cmd/atheon@v1.0.0
```

Or restore the previous binary from your backup. There is no
in-place migration that would need to be undone — upgrades are
in-place overwrites, downgrades are the same operation in
reverse.

## Pre-flight checklist

Before upgrading a production pipeline, verify on a staging
copy:

- [ ] `atheon --version` reports the new version.
- [ ] `atheon .` returns the same exit code as before on a
      representative input.
- [ ] The pre-commit hook still passes.
- [ ] If you embed `core/`, `go test ./... -race` is clean on
      the consumer project.

## See also

- [`CHANGELOG.md`](CHANGELOG.md) — what changed.
- [`MIGRATION.md`](MIGRATION.md) — how to adapt to a major
  version bump.
- [`RELEASE.md`](RELEASE.md) — release cadence and process.