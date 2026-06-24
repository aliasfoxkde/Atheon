# ADR 0001: YAML-on-disk as the canonical pattern format

- **Status**: Accepted
- **Date**: 2026-06-23
- **Deciders**: aliasfoxkde
- **Supersedes**: —

## Context

Atheon's pattern engine needs a representation that:

1. Survives version control without opaque encoding (so contributors can
   diff, review, and amend patterns in PRs).
2. Loads fast at startup (the engine reads ~250 patterns on every CLI
   invocation).
3. Tolerates the stdlib regex engine (RE2) — no lookahead, no
   backreferences, bounded memory.
4. Is authorable by humans without writing Go.

Two obvious formats:

- **YAML on disk + bundler to gzip+JSON**: text in git, fast at runtime.
- **Go source code**: requires recompile to add a pattern; high friction
  for the community-contribution model.

## Decision

**YAML files in `community/<category>/*.yaml` are the canonical format.
The bundler compiles them into a single `core/patterns.bundle`
(gzip-compressed JSON) that is `//go:embed`-ed into the binary.**

Each YAML file has exactly three fields:

```yaml
name: anthropic-api-key
match: '\bsk-ant-api[0-9]{2}-[A-Za-z0-9_\-]{93,}\b'
enabled: true
```

- `name`: kebab-case identifier; unique across the bundle (validated by
  the bundler — duplicate names are a hard error).
- `match`: RE2-compatible regex.
- `enabled`: optional bool, default true. Patterns that fire false
  positives in practice are shipped disabled (opt-in via
  `SetActiveCategories` or `--category`).

## Consequences

**Positive**

- Patterns can be added, removed, or amended by anyone who can open a
  PR — no Go knowledge required.
- The bundle is a single `//go:embed`-ed artifact; the binary is
  self-contained. No filesystem dance at startup beyond the optional
  `~/.atheon/patterns.bundle` override.
- The bundler is a fast, deterministic compile step (under 1s for 274
  patterns); it runs in CI and in `goreleaser`'s `before` hook.

**Negative**

- Two representations (YAML on disk, JSON in the bundle) means two
  sources of truth. We mitigate by treating the bundle as a *build
  artifact* — never edit it by hand, always regenerate via
  `go run ./bundler`.
- RE2 imposes a real cost on pattern authors: no lookahead, no
  backreferences, no atomic groups. Several patterns that would be
  trivial in PCRE have to be reformulated. This is documented in
  `community/README.md`.

**Neutral**

- The 255 → 274 pattern count is tracked separately by
  `scripts/pattern-count.sh` to avoid drift between docs and reality.