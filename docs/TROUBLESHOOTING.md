# Troubleshooting

When something doesn't work, start here. Each section lists the
symptom, the most likely cause, and the smallest fix.

> Looking for the legacy "Enhanced Atheon" troubleshooting content?
> It moved to [`reports/TROUBLESHOOTING_LEGACY.md`](reports/TROUBLESHOOTING_LEGACY.md)
> when this consolidated guide was promoted to the project root.

## Quick reference

| Symptom | Jump to |
|---|---|
| `command not found: atheon` | [Install failures](#install-failures) |
| Pattern count seems low or wrong | [Pattern bundle problems](#pattern-bundle-problems) |
| A scan hangs forever | [Scans that hang](#scans-that-hang) |
| False positive on test data | [False positives](#false-positives) |
| `go build` fails after upgrade | [Build failures](#build-failures) |
| Tests fail with `-race` | [Race conditions](#race-conditions) |
| CI fails on Windows but works locally | [Cross-platform issues](#cross-platform-issues) |
| MCP server won't start | [MCP problems](#mcp-problems) |

## Install failures

### `command not found: atheon`

The binary isn't on `$PATH`.

- **`go install` path:** `$(go env GOPATH)/bin` must be on
  `$PATH`. `go env GOPATH` prints the directory; add
  `${GOPATH}/bin` to your shell profile.
- **Downloaded binary:** extract to a directory on `$PATH`
  (`~/.local/bin`, `/usr/local/bin`, etc.) or call it with an
  absolute path.

### Permission denied on Linux/macOS

The binary is not executable.

```sh
chmod +x atheon
```

If you downloaded via `curl` and got "cannot execute binary
file", your CPU architecture probably does not match. Check
`uname -m` and pick the right archive.

## Pattern bundle problems

### Pattern count looks low

The shipped binary embeds a fixed bundle. To verify what is
loaded:

```sh
atheon list --all | wc -l
```

If the number is far below the
[README badge](../README.md), you may be running a stale
binary. Update with:

```sh
go install github.com/.../atheon/cmd/atheon@latest
```

Or download the latest release binary.

### `atheon update` fails

`update` requires outbound HTTPS to the bundle host. Behind a
corporate proxy:

```sh
HTTPS_PROXY=https://proxy.example.com atheon update
```

In an air-gapped CI runner, skip `update` entirely — the
embedded bundle is the source of truth.

## Scans that hang

The scanner respects `context.Context`. A hanging scan usually
means the context has no deadline and the scan is walking a
filesystem you didn't expect. Two fixes:

1. **Pass a deadline** when embedding the library:

   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   core.ScanDir(ctx, root)
   ```

2. **Inspect the input.** If you passed `.` from a directory
   that contains `node_modules` or `.git`, exclude them via
   `.atheonignore`:

   ```
   node_modules/
   .git/
   ```

## False positives

### Test data trips the scanner

Test fixtures often contain example credentials. Two options:

- **Ignore the file** via `.atheonignore`:

  ```
  testdata/
  ```

- **Disable the pattern** for the run:

  ```sh
  atheon . --exclude=aws-access-key
  ```

### Specific regex too broad

File an issue with the
[pattern submission template](../.github/ISSUE_TEMPLATE/pattern_submission.md)
and include the false-positive sample. The maintainers will
either tighten the regex or document when to disable the
pattern.

## Build failures

### `cannot find package` after `go mod tidy`

Run `go mod tidy` from the repository root, not from a
subdirectory. Go modules work at the module root; tidy
elsewhere will leave `go.sum` inconsistent.

### `undefined: core.ScanString` (older API)

You are looking at code written against a pre-1.0 API. See
[`MIGRATION.md`](MIGRATION.md) for the symbol renames.

## Race conditions

### `data race` under `go test -race`

If you wrote a test that calls `os.Chdir`, switch to
`t.Chdir`:

```go
t.Chdir(tmpDir)  // restores via t.Cleanup
```

`t.Chdir` exists in Go 1.24+. If you must support older Go,
keep the manual `defer os.Chdir(orig)` pattern but make sure
your test does not run in parallel with any other test that
relies on the package CWD.

## Cross-platform issues

### Tests fail on Windows only

Two common causes:

- **chmod-based tests.** `os.Chmod(path, 0o000)` does not
  restrict the file owner on Windows; only the read-only bit
  is honored. Skip those tests on Windows with
  `if runtime.GOOS == "windows" { t.Skip(...) }`.
- **Path separators.** A path containing `\` will fail on
  Linux/macOS. Always build paths with `filepath.Join`.

### macOS: `executable file not found in $PATH` from integration tests

Caused by a `go test` run that changed CWD mid-flight. Update
to the 1.0 integration-test pattern (build once in `TestMain`,
share via a package-level variable) and the issue goes away.

## MCP problems

### `connection refused` on `127.0.0.1:8080`

The MCP server defaults to a stdio transport, not a TCP port.
Make sure your client is configured for stdio, not HTTP. If
you do want TCP, see the `--listen` flag in
[`API.md`](API.md).

### JSON-RPC parse errors

The server expects newline-delimited JSON. If your client is
sending length-prefixed framing, switch to `\n`-delimited or
add a framing shim.

## Still stuck?

1. Re-read the relevant section of [`FAQ.md`](FAQ.md).
2. Search [existing issues](../../issues) for the error message.
3. If nothing matches, open an issue with the *Bug report*
   template under `.github/ISSUE_TEMPLATE/`. Include the full
   command and its output.