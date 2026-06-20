# Auditing

Atheon includes an audit pipeline that runs code-quality checks against the codebase and produces structured reports in both JSON and Markdown formats.

---

## Running Audits

### CLI

```bash
atheon audit              # audit the current directory
atheon audit ./path/to/repo  # audit a specific directory
```

The CLI writes its report to `docs/audits/<timestamp>/REPORT.md` and `docs/audits/<timestamp>/REPORT.json`, then prints the path to the Markdown report.

### Make

```bash
make audit                # run every audit gate in sequence
make audit-dead-code      # dead code check only (staticcheck)
make audit-nolint         # enumerate //nolint annotations
make audit-fixmes         # enumerate TODO/FIXME/XXX markers
make audit-sentinels      # verify exported sentinel errors have callers
```

---

## Audit Checks

The `core.Audit()` function runs five checks:

| Check | What it does | Failure severity |
|-------|-------------|-----------------|
| `dead-code` | Stub (always passes); real enforcement via `staticcheck` in pre-commit | N/A |
| `nolint` | Enumerates every `//nolint` annotation in the codebase | warning |
| `todo-fixme` | Enumerates every `// TODO`, `// FIXME`, `// XXX` marker | info |
| `go-vet` | Runs `go vet ./...` across the entire project | error |
| `sentinel-errors` | Verifies every exported `var ErrFoo = errors.New(...)` has at least one caller | warning |

The `dead-code` check is a placeholder; the actual dead-code enforcement is handled by:
- **Pre-commit hook** (`.githooks/pre-commit`): runs `staticcheck` and blocks commits with dead code
- **Make target** (`make audit-dead-code`): runs `scripts/audit-dead-code.sh`

---

## Report Structure

The JSON report follows this schema:

```json
{
  "version": "1.0",
  "generatedAt": "2026-06-20T15:00:00Z",
  "root": "/path/to/repo",
  "elapsedMs": 420,
  "results": [
    {
      "check": "go-vet",
      "passed": true,
      "findings": []
    },
    {
      "check": "nolint",
      "passed": false,
      "findings": [
        { "file": "core/bundle.go", "line": 42, "message": "ineffectual assignment", "severity": "warning" }
      ]
    }
  ],
  "summary": {
    "total": 5,
    "passed": 4,
    "failed": 1
  }
}
```

---

## CI Integration

The audit pipeline is wired into GitHub Actions via `.github/workflows/quality-assurance.yml`, which runs `make audit` on every push to `main` and on pull requests.

To run the audit in CI independently:

```yaml
- name: Run audit
  run: make audit
```

---

## Adding a New Audit Check

1. Implement a function matching `func(root string) AuditResult`.
2. Add it to the `checks` slice in `core/audit.go`.
3. If it produces findings, set `Result.Passed = false` and append `AuditFinding` structs.
4. If it requires a new tool, add a corresponding `make audit-<name>` target and wire it into the umbrella `audit` target.

---

## Historical Reports

Past audit reports are stored in `docs/audits/` with timestamped subdirectories:

```
docs/audits/
├── 2026-06-20-150405/
│   ├── REPORT.json
│   └── REPORT.md
└── 2026-06-19-090000/
    ├── REPORT.json
    └── REPORT.md
```
