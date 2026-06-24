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

## Conventions

- Use Markdown headings: `#`, `##`.
- File name: `NNNN-kebab-case-slug.md` with zero-padded sequence number.
- Required sections: Status, Date, Deciders, Context, Decision,
  Consequences.
- Status values: `Proposed`, `Accepted`, `Superseded`, `Deprecated`.