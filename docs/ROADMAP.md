# Roadmap

This is the long-term plan for Atheon. It is intentionally
high-level — concrete work lives in [TASKS.md](TASKS.md) and
[planning/](planning/), and progress is tracked in
[PROGRESS.md](PROGRESS.md).

## Mission

Make secret, PII, and quality-issue detection the boring,
fast, dependable step at the start of every CI pipeline. The
project optimizes for the engineer who wants answers, not for
the engineer who wants to configure a tool.

## Guiding principles

- **Boring defaults, powerful overrides.** The shipped pattern
  bundle catches the common cases. Every default can be tuned
  without forking.
- **Single binary, multiple surfaces.** One scanner that works
  on disk, on stdin, in env vars, and as a library function — no
  "lite" / "pro" / "enterprise" tier.
- **Patterns as data.** Adding a pattern is a YAML edit, not a
  code change. The match engine is the only piece that needs to
  stay in Go.
- **Honest numbers.** No inflated pattern counts, no
  hand-curated category breakdowns. Counts are pulled from the
  bundle at build time.

## Themes

### 1.0 → 1.x — **Hardening** (current)

- **Pattern coverage.** Close the highest-impact gaps in
  cloud-provider credentials (AWS, Azure, GCP), CI tokens
  (GitHub Actions, GitLab CI, CircleCI), and PII (passport,
  national-ID formats).
- **Performance.** Keep the 1.0 throughput baseline (≈1 GB/s on
  cold start on a single core) and add benchmark-driven
  regression coverage.
- **Docs parity.** Every public function has a godoc example.
  Every CLI flag is documented in `docs/API.md`.
- **Distribution.** Goreleaser-produced binaries for Linux,
  macOS, Windows on amd64 and arm64. Homebrew tap and a
  Scoop bucket.

### 1.x → 2.0 — **Expanding surfaces**

- **MCP server.** Already in `cmd/mcp/`. Promote from
  experimental to stable; document the JSON-RPC surface in
  `docs/API.md`.
- **Library API.** Stabilize the Go API: explicit version
  guarantees, deprecation policy, and a `v2` branch that
  cleans up `context.Context` placement.
- **Pre-commit and IDE plugins.** A Git hook (already shipped)
  and a VS Code extension that calls the same library.

### 2.0 → 3.0 — **Ecosystem**

- **Remote bundle distribution.** `atheon update` is already
  wired up; back it with a signed manifest and a
  mirror-fallback chain so it works in air-gapped CI runners.
- **Policy bundles.** Named profiles (`config/profiles/*.json`)
  that combine category selection, ignore rules, and exit-code
  policy. Goal: one of these can be passed to a third-party CI
  action.
- **Pattern authoring tools.** A linter that flags regexes
  prone to ReDoS, a corpus-based fuzz tester, and a
  coverage-by-category report.

## Non-goals

These are explicitly **not** on the roadmap. Saying no is as
important as saying yes.

- **A GUI.** Atheon stays a CLI and a library. Web UIs are
  better built on top of the library.
- **A cloud-hosted scan service.** The project ships no SaaS.
  A self-hosted scan server is feasible but is a community
  project, not a core deliverable.
- **Auto-remediation.** Atheon finds things; it does not open
  PRs, rewrite files, or rotate secrets. Those tools exist;
  build them on top of this one.
- **Vendor-specific SDKs.** No `pip install atheon`, no
  `npm install @atheon/core`. The library is Go; users who
  want bindings can generate them with the public API.

## Release cadence

- **Minor versions** (`1.0` → `1.1`) ship roughly every 6–8
  weeks. They add patterns and library features; they keep
  backward compatibility.
- **Patch versions** (`1.0.0` → `1.0.1`) ship as needed for
  bug fixes. They never add patterns.
- **Major versions** are reserved for breaking changes to the
  CLI flags or the public Go API. The next major (`v2`) is
  scoped to library ergonomics; see `docs/MIGRATION.md` once
  it is published.

## How to influence the roadmap

- Open a *Feature request* issue using the template under
  `.github/ISSUE_TEMPLATE/feature_request.md`. Mark
  "Willingness to contribute" if you can land it.
- Discussions under the *Ideas* category are the right place
  for "wouldn't it be cool if…" brainstorming.
- The maintainers meet informally in the GitHub Discussions
  threads tagged `roadmap` and produce a refreshed roadmap
  roughly every quarter.

## Tracking

| Theme | Document |
|---|---|
| Hardening (current) | `docs/TASKS.md`, `docs/PROGRESS.md` |
| Expanding surfaces | `docs/planning/` |
| Ecosystem | `docs/planning/PROJECTS.md` |
