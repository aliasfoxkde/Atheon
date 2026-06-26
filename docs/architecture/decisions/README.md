# Architecture Decision Records

This directory contains the Architecture Decision Records (ADRs) for
Atheon-Enhanced. Each ADR documents one significant design choice:
the context that prompted it, the chosen path, and the consequences
(positive, negative, neutral) that follow from it.

ADRs are immutable once accepted. Superseded ADRs are kept in place
and pointed to by their replacement. New ADRs get the next sequential
number (`NNNN-short-slug.md`).

## Index

- [ADR 0001 — Pattern YAML format](0001-pattern-yaml-format.md):
  Why patterns live as YAML on disk and compile to a gzip+JSON bundle.
- [ADR 0002 — Five-workflow CI surface](0002-ci-workflow-consolidation.md):
  Why we consolidated 10 overlapping GitHub Actions workflows to 5.
- [ADR 0003 — MCP server design](0003-mcp-server-design.md):
  Why the MCP server speaks stdio JSON-RPC and how the seven tools
  are bounded.
- [ADR 0004 — RE2 regex choice](0004-re2-regex.md):
  Why the scanner uses Google's RE2 engine (linear time, no catastrophic
  backtracking) instead of PCRE.
- [ADR 0005 — Gzip+JSON bundle format](0005-gzip-bundle.md):
  Why patterns compile into a single gzip-compressed JSON blob embedded
  in the binary.
- [ADR 0006 — Parallel test requirement (-p 1)](0006-parallel-tests.md):
  Why `go test` must use `-p 1` — `core/` has package-level state in
  `init()` that breaks under parallel package execution.

## Conventions

- Use Markdown headings: `#`, `##`.
- File name: `NNNN-kebab-case-slug.md` with zero-padded sequence number.
- Required sections: Status, Date, Deciders, Context, Decision,
  Consequences.
- Status values: `Proposed`, `Accepted`, `Superseded`, `Deprecated`.