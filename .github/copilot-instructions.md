# Copilot Instructions for Atheon-Enhanced

## What this project does

Atheon is a Go pattern-matching engine for detecting secrets, PII, and code quality issues. It scans files, directories, environment variables, and strings against a library of YAML-defined regex patterns. It outputs plain text, JSON, or SARIF format.

Key binaries:
- `atheon` — CLI scanner (`cmd/atheon/`)
- `atheon-mcp` — MCP server for AI assistant integration (`cmd/mcp/`)

## Repository layout

```
core/           — pattern loading, scanning engine, matcher
cmd/atheon/     — CLI entrypoint, flag parsing, output formatting
cmd/mcp/        — MCP server (stdio JSON-RPC)
community/      — YAML pattern library, organized by category:
  secrets/      — API keys, tokens, credentials
  pii/          — Personal identifiable information
  code-quality/ — Anti-patterns, insecure coding practices
  devops/       — Infrastructure misconfigs
  ai-detection/ — AI prompt injection patterns
bundler/        — Compiles community/ YAML into core/patterns.bundle
docs/           — Documentation (published via GitHub Pages)
```

## Pattern YAML format

Each pattern file follows this schema:

```yaml
name: my-pattern-name
description: What this pattern detects
severity: low|medium|high|critical
category: secrets|pii|code-quality|devops|ai-detection
enabled: true
pattern: "regex here"
keywords: ["optional", "context", "keywords"]
true_positives:
  - "example that should match"
false_positives:
  - "example that should NOT match"
```

## Common issue types

**False positive reports:** The user is reporting that Atheon flagged something it shouldn't have. Look up the pattern name in `community/<category>/<name>.yaml`. Explain why the regex matched and whether narrowing with `keywords` would help.

**False negative reports:** Atheon missed something it should have caught. Check if an existing pattern covers the case or if a new pattern is needed.

**Build failures:** Run `go build ./...` from the repo root. The bundler (`go run ./bundler`) must be run before building to regenerate `core/patterns.bundle` if community patterns changed.

**Test failures:** Run `go test ./... -timeout 15m -p 1`. Tests use real pattern files — a failing test usually means a pattern regex changed without updating test fixtures.

## CI workflows

- `ci.yml` — main build, test, coverage gate (70%)
- `quality-assurance.yml` — self-scan, pattern count gate (250+), lint
- `scheduled-release.yml` — auto-increments patch version on 10th/21st, requires `release` environment approval
- `publish.yml` — triggered by `v*` tags, runs goreleaser to build/publish binaries
- `dev-testing.yml` — relaxed gates for `dev/testing` branch (50% coverage, non-blocking)

## Contribution guidelines

- Pattern PRs go to `community/` — must include `true_positives` and `false_positives`
- Code PRs must pass `go vet ./...` and `go test ./...`
- Coverage must not drop below 70% on `main`
- See `docs/patterns/contributing-patterns.md` for the full pattern contribution guide
