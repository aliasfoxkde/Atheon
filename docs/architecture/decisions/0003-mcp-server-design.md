# ADR 0003: MCP server over stdio JSON-RPC

- **Status**: Accepted
- **Date**: 2026-06-23
- **Deciders**: aliasfoxkde
- **Supersedes**: —

## Context

Modern IDE clients (Cursor, Claude Desktop, Zed, etc.) integrate
third-party tools via the [Model Context Protocol][mcp], a JSON-RPC
2.0 dialect running over stdin/stdout. To make Atheon invokable from
those environments without the user pasting CLI output into a chat
manually, we need a server that:

- Speaks JSON-RPC 2.0 over stdio (no port binding, no auth dance).
- Exposes scan operations as named *tools* the host can discover via
  `tools/list` and invoke via `tools/call`.
- Is rate-limited so a runaway host loop cannot DoS the local scanner.
- Carries version info in `serverInfo` so hosts can display it.

[mcp]: https://modelcontextprotocol.io/

## Decision

**Ship `atheon-mcp`, a single-binary stdio JSON-RPC server written in
the same Go module, with seven tools:**

| Tool | Purpose |
|------|---------|
| `scan_string` | Scan a literal string for pattern matches |
| `scan_file`   | Scan a file on disk |
| `scan_dir`    | Recursively scan a directory |
| `scan_env`    | Scan the process environment variables |
| `list_patterns` | Enumerate loaded patterns (optionally filtered by category) |
| `list_categories` | Enumerate known categories |
| `update_bundle` | Re-download the latest pattern bundle |

The server is rate-limited via a stdlib-only token bucket (10 req/sec,
burst 20) that returns JSON-RPC code **-32000** (Server Error) when
exceeded. (Earlier drafts used -32600 "Invalid Request", which is the
wrong semantic — -32600 implies the *client* sent a malformed request,
not that the *server* is throttling.)

The version string in `serverInfo.version` is read from a package-level
`var version = "dev"`, injected at build time via
`-ldflags "-X main.version=<v>"`. Both `atheon` and `atheon-mcp`
builds in `.goreleaser.yml` apply this ldflag.

## Consequences

**Positive**

- One transport, no extra dependencies: `bufio.Scanner` on stdin,
  `json.NewEncoder` on stdout. The Go stdlib is the entire transport.
- A runaway host cannot OOM the host machine: the rate limiter
  caps at burst 20 even with infinite request rates.
- Tool discovery is standardised — hosts can present a UI of all
  seven tools without Atheon-specific code paths.
- Version reporting is honest: hosts see "dev" for `go run` and the
  real tag for releases.

**Negative**

- A token-bucket rate limiter is global to the process, which makes
  test suites fragile (a 25-call test suite can exhaust a 20-burst
  bucket). We mitigate by adding `TestMain` to install a
  generously-sized limiter (10000/10000) for the test binary; the
  rate-limit denial test still swaps in a zero-token limiter via defer.
- Stdout is consumed exclusively by the JSON-RPC stream — any
  incidental `fmt.Println` would corrupt the protocol. The codebase
  uses `slog` (stderr) for diagnostics and `fmt.Fprintf(os.Stderr, ...)`
  for the one-off "malformed request" warning.

**Neutral**

- Seven tools is the right number for now. Adding more should follow
  the principle: each new tool maps to a `core/` API or a thin wrapper
  over one. Tools that exist only to combine other tools should be
  composed by the host, not added to the server.