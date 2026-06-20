# Migration Guide

How to move from one major version of Atheon to the next. This
file is intentionally forward-looking: the v2 API is the only
breaking change we have planned, and it has not shipped yet.
The structure of the file mirrors the v1 → v2 migration; once
that is finalised the other sections will follow the same
shape.

> When v2 ships, this document will have a top-level section
> for each previous major version (e.g. "From v1 to v2"). Until
> then, this file documents the *planned* breaking changes so
> early adopters can prepare.

## From v1 to v2 (planned)

> Status: **planned, not implemented.** v2 lives on a feature
> branch; nothing here is final until the v2.0.0 release is
> tagged.

### Why a v2?

The v1 API is fine but has a few sharp edges that are easier to
fix with a clean break than with deprecation aliases:

1. **Library-only `context.Context` placement.** In v1, only the
   library functions take a context. The CLI flags that wrap them
   do not, so the CLI cannot pass a per-flag deadline through to
   the scanner. v2 makes the CLI honour context-aware flags too.
2. **Option structs.** Several library functions in v1 take
   long argument lists that grow with every release. v2 collapses
   them into option structs so adding a knob does not change
   every call site.
3. **Result struct.** v1 returns `([]Finding, *Stats, error)`.
   v2 returns a `*Result` that bundles the same fields and adds
   a few that v1 could not express.

### Planned breaking changes

| v1 | v2 |
|---|---|
| `core.ScanFile(ctx, path)` | `core.ScanFile(ctx, path, opts...)` |
| `core.ScanDir(ctx, root)` | `core.ScanDir(ctx, root, opts...)` |
| `core.ScanString(ctx, content, name)` | `core.ScanString(ctx, content, name, opts...)` |
| `core.ScanEnv(ctx)` | `core.ScanEnv(ctx, opts...)` |
| `core.ScanStdin(ctx)` | `core.ScanStdin(ctx, opts...)` |
| `findings, stats, err := core.ScanFile(...)` | `result, err := core.ScanFile(...)` |

### Migration steps (planned)

1. **Update import path.** v2 will move to
   `github.com/.../atheon/v2/core`. Update your `go.mod` and
   imports.
2. **Add the option struct.** Where v1 had positional arguments,
   v2 takes an `Options` struct. Start by passing
   `core.DefaultOptions` and only set the fields you care about.
3. **Switch to the result struct.** Replace
   `findings, stats, err := core.ScanFile(...)` with
   `result, err := core.ScanFile(...)` and read
   `result.Findings` and `result.Stats`.
4. **Run `go vet` and `go test ./... -race`.** v2's stricter
   context-placement rules will flag any code that loses a
   context along the way.

### What will NOT change

- **Pattern names and categories.** A pattern called
  `aws-access-key` in v1 is still `aws-access-key` in v2.
- **YAML format.** Patterns are data; the v1 YAML loads in v2.
- **Exit codes.** The CLI maps the same errors to the same exit
  codes in both versions.
- **Embedded bundle.** v2 binaries still ship with the bundle
  embedded; `atheon update` still works.

## From pre-1.0 (0.x) to v1

> Status: **already happened.** This section is kept for
> completeness for users who are still on a 0.x release.

The 0.x series had a different CLI surface and a much smaller
pattern bundle. The migration to 1.0 is largely a clean
install:

1. **Uninstall the 0.x binary.** `rm $(which atheon)` or
   `go clean -i github.com/.../atheon/cmd/atheon`.
2. **Install v1.** Follow [INSTALL.md](INSTALL.md).
3. **Migrate ignore rules.** v1 honours `.atheonignore` and
   `.gitignore`; 0.x used a different file called `.scannerignore`.
   Rename the file (or run with `--ignore-from=.atheonignore`
   during the transition).
4. **Re-evaluate exit codes.** v1 maps `ErrSecretInEnv` to
   exit `1`. 0.x used exit `2`. Update any pipeline that
   branches on the exit code.
5. **Re-test your patterns.** v1 ships with 179 patterns vs.
   the 0.x default of 57; expect fewer false negatives and a
   few new false positives to tune.

## How to ask for migration help

If you are mid-migration and the docs do not cover your case:

- Open a *Documentation* issue under
  `.github/ISSUE_TEMPLATE/documentation.md`.
- For v2 specifically, the maintainers will be running a
  "v2 readiness" thread on GitHub Discussions once the API
  is frozen.
- For commercial migrations, see [`GOVERNANCE.md`](GOVERNANCE.md)
  for contact channels.