# Engineering Standards

The conventions that the Atheon codebase follows, the
reasoning behind them, and the exceptions we accept. Pair
this with [`DESIGN.md`](DESIGN.md) for the architectural
decisions and [`CONTRIBUTING.md`](../.github/CONTRIBUTING.md)
for the day-to-day workflow.

> **TL;DR:** Go's standard library, `gofmt`, and `go vet`.
> Anything more is in this file.

## Language and runtime

- **Go version:** the version pinned in `go.mod`. CI
  verifies the codebase against this version and the two
  previous minor releases.
- **Module path:** `github.com/aliasfoxkde/atheon`.
  Sub-packages use the same root.
- **Build tags:** none required. We do not use
  `//go:build` to fork the codebase on OS; runtime checks
  (`runtime.GOOS`) are reserved for cases where there is
  no portable alternative.

## Code style

- **`gofmt` is the law.** The CI runs `gofmt -l` and fails
  on any diff. There is no "house style" layered on top
  of `gofmt`; if the formatter is happy, we are happy.
- **`goimports` for grouping.** Imports are grouped
  three ways: stdlib, third-party, internal. A single
  blank line separates each group. CI runs
  `goimports -l`.
- **Line length:** no hard cap. The default Go formatter
  wraps as needed, and we follow its lead. If a line is
  awkwardly long, refactor the expression rather than
  breaking it across multiple lines.
- **Comments:** every exported identifier has a doc
  comment that starts with the identifier name and reads
  as a sentence. Example:
  ```go
  // ScanFile scans a single file and returns its findings.
  func ScanFile(...) { ... }
  ```
  Unexported identifiers are commented when the
  implementation is non-obvious.

## Naming

- **Packages:** short, lowercase, single-word. Avoid
  stutter: `core.Scan`, not `core.CoreScan`.
- **Types:** `PascalCase`. Acronyms are all-caps
  (`HTTPClient`, not `HttpClient`).
- **Functions and methods:** `PascalCase` for exported,
  `camelCase` for unexported. Predicate functions start
  with `Is`, `Has`, or `Should` (`IsValid`, `HasPattern`).
- **Variables:** `camelCase`. Loop indices are `i`, `j`,
  `k`. Short-lived variables are short
  (`for _, f := range findings`).
- **Constants:** `PascalCase` for exported, `camelCase`
  for unexported. The `kFoo` / `eFoo` prefixes some code
  bases use are not used here.
- **Receivers:** one or two letters, derived from the
  type. Consistent across the type. Never `this` or
  `self`.

## Errors

- **Wrap, don't replace.** Use `fmt.Errorf("...: %w", err)`
  to add context, or a typed `*Error` from
  `core/errors.go` for richer context.
- **Sentinel errors** are exported for cases that callers
  must check with `errors.Is`. New sentinels are added in
  `core/errors.go`.
- **No `panic` in the library.** The CLI may `panic` on
  programmer error, but the library always returns an
  error. The one exception is `init` functions, which may
  `panic` if a required asset is missing — that is
  considered a build-time error, not a runtime error.
- **Error messages** are lowercase and do not end with a
  period. They read as a sentence fragment because they
  are wrapped by callers: `"read config: open file: no
  such file"`, not `"Read config: open file: no such
  file or directory."`.

## Concurrency

- **Share by communicating.** Channels and `sync` types
  are both fine; the choice is driven by clarity, not
  dogma. A mutex around a counter is usually clearer than
  a channel for the same purpose.
- **Context first.** Every function that may block takes
  a `context.Context` as its first argument. The context
  is never stored in a struct.
- **No goroutine leaks.** Every `go func` has a documented
  termination condition. The scanner's worker pool is the
  only place we spawn goroutines; the pool is closed in
  `defer` and the `WaitGroup` is waited on.
- **Race detector in CI.** `go test -race` runs on every
  PR and on `stable`. A race that only reproduces under
  the detector is treated as a bug.

## Testing

- **Table-driven.** Tests are written as a slice of test
  cases, each with a name. A failing case is identified by
  its name, not by its index. See `core/coverage_test.go`
  for the canonical shape.
- **Subtests via `t.Run`.** The outer test is the function
  (`TestScanFile`), the inner is the case
  (`t.Run("empty file", ...)`). Test names read as paths
  in the verbose output.
- **No `init` in test files.** Test setup is explicit
  inside the test function or in `TestMain`. Mocks are
  constructed in the test that uses them.
- **Coverage is reported but not gated.** The CI prints
  per-package coverage. New code is expected to maintain
  or improve the package's coverage. We do not block
  merges on a coverage number; we do block merges on
  *unjustified* coverage drops.
- **Test data in `testdata/`.** Go's tooling ignores
  `testdata/` directories, so they are the right place
  for fixtures, golden files, and pattern samples.

## Dependencies

- **The standard library is the first choice.** A
  third-party dependency is added only when the stdlib
  does not have a workable answer. The dependency list is
  reviewed at every minor release.
- **Pinning.** Dependencies are pinned to a specific
  version in `go.mod` and updated via Dependabot. We do
  not use `^` or `~` ranges.
- **`go mod tidy` is committed.** `go.mod` and
  `go.sum` are the canonical record; no untracked
  dependencies are allowed.
- **`go mod verify` in CI.** The CI verifies the checksum
  database before running tests.

## Performance

- **Measure first.** Optimisation is gated on a benchmark
  in `BENCHMARKS.md` that demonstrates the cost. Without
  a measurement, the change is a "cleanup", not a
  performance improvement.
- **No premature allocation.** A `make([]T, 0, n)` is
  preferred over repeated `append` only when the size is
  known up front and is large.
- **Hot paths are documented.** A function that runs
  per-line or per-byte has a comment that calls this out.
  See `core/engine.go` for examples.

## Documentation

- **Public API is documented in the source.** Doc
  comments on exported identifiers are required; CI
  runs `go doc ./...` and fails on missing comments for
  the public surface (`core/`, `bundler/`).
- **Behavioural docs in `docs/`.** Long-form explanations
  (how to migrate, how to write a pattern, how the
  scanner works) live in `docs/`, not in code comments.
  This keeps the source readable and the docs linkable.
- **Examples in `Example*` functions.** The `go test
  -run=Example` invocation picks up `// Output:`
  comments and verifies that the example produces the
  claimed output. Use this for runnable doc examples; it
  is the only test that doubles as documentation.

## Versioning and API stability

- **SemVer 2.0.0** for the Go module, the CLI, and the
  embedded bundle. A change to any one is a version bump
  in all three.
- **The v1 API is stable** within the 1.x line. The
  planned v2 API is documented in
  [`MIGRATION.md`](MIGRATION.md).
- **Deprecations are signalled in the doc comment** with
  a `// Deprecated: …` line. Deprecated symbols are
  retained for at least one minor release before
  removal.

## Security

- **Secrets in test fixtures** are obviously fake
  (`AKIAIOSFODNN7EXAMPLE` for AWS keys, etc.). A scanner
  that finds a real secret in a test fixture is a
  security incident, not a "funny coincidence".
- **No credentials in the repo.** Even `.env.example`
  files use placeholder values. CI scans the repo for
  known secret patterns on every push.
- **`SECURITY.md`** is the disclosure channel. Private
  issues for security reports; public issues for
  everything else.

## Exceptions

Every rule has an exception, and every exception is in a
PR description or a code comment. The exception list is
short:

- **Windows file I/O is slow.** We accept it; the
  alternative (raw syscalls) is not worth the portability
  cost. See [`PERFORMANCE.md`](PERFORMANCE.md).
- **The CLI uses `os.Exit` directly** rather than
  returning an error. The library never does; the CLI is
  the boundary that turns errors into exit codes.
- **The bundler has a long-running command** that uses
  signals. The scanner does not, because the scanner is
  embeddable and the host owns the signal handling.

If you find yourself adding a new exception, write it
down. Unrecorded exceptions are technical debt.
