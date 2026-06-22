# Atheon-Enhanced: Comprehensive Improvement Plan

**Prepared for:** AI handoff — implement each section independently, branch per section, PR for review.  
**Repository:** https://github.com/aliasfoxkde/Atheon-Enhanced  
**Module:** `github.com/aliasfoxkde/Atheon`  
**Go version:** 1.21+ (CI tests 1.21–1.24)  
**Current state:** 223 patterns, 94.9% test coverage, 10 CI workflows  

---

## Environment Notes

- Platform: Windows 11, Git Bash
- Go binary: `/c/Program Files/Go/bin/go` (add to PATH)
- `goimports` and `staticcheck` are blocked by corporate proxy on this machine — embed them or install from a mirror; the unrestricted system should install them via `go install golang.org/x/tools/cmd/goimports@latest` and `go install honnef.co/go/tools/cmd/staticcheck@latest`
- Test command: `go test ./... -p 1 -timeout 15m` — the `-p 1` is **mandatory** (global bundle state corrupts under parallel package execution)
- Bundle rebuild: `go run ./bundler` (output: `core/patterns.bundle`, embedded via `//go:embed`)
- Hooks: `.githooks/pre-commit` and `.githooks/pre-push`; wired via `git config core.hooksPath .githooks`

---

## Section 1: CI/CD Enhancements

**Branch name:** `feat/ci-improvements`

### 1.1 Consolidate 10 workflows into ~4

Current 10 workflows have massive duplication. Consolidate to:
- `ci.yml` — test + lint + build (all PRs and pushes to main)
- `security.yml` — CodeQL + self-scan (merge of codeql.yml + security-scanning.yml)
- `release.yml` — publish + scheduled-release (only on tags/schedule)
- `sync.yml` — keep sync-stable-clean.yml as-is

Delete: `comprehensive-ci.yml`, `quality-assurance.yml` (duplicate of ci.yml), `auto-merge.yml` (use GitHub's native auto-merge setting instead).

### 1.2 SHA-pin all GitHub Actions

All `uses:` lines currently use mutable tags (`@v4`, `@v5`). Replace every action with a SHA-pinned reference and a human-readable comment:

```yaml
# actions/checkout@v4 (2024-10-22)
uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683
```

Actions to pin:
- `actions/checkout@v5` → latest SHA for v4 (use v4; v5 unconfirmed stable)
- `actions/setup-go@v6` → latest SHA for v5
- `actions/upload-artifact@v4`
- `codecov/codecov-action@v4`
- `golangci/golangci-lint-action@v6`
- `github/codeql-action/init@v4`
- `github/codeql-action/analyze@v4`
- `goreleaser/goreleaser-action@v6`

Use `pinact` or `renovatebot` to automate: `go install github.com/suzuki-shunsuke/pinact/cmd/pinact@latest`

### 1.3 Add dependency update bot

Create `.github/renovate.json`:
```json
{
  "$schema": "https://docs.renovatebot.com/renovate-schema.json",
  "extends": ["config:recommended"],
  "packageRules": [
    { "matchManagers": ["gomod"], "automerge": true, "matchUpdateTypes": ["patch"] },
    { "matchManagers": ["github-actions"], "automerge": false }
  ]
}
```

### 1.4 Add test result reporting

After `go test`, generate a JUnit XML report and upload it:

```yaml
- name: Run tests with JUnit output
  run: |
    go install github.com/jstemmer/go-junit-report/v2@latest
    go test ./... -p 1 -v -race -coverprofile=coverage.out 2>&1 | go-junit-report -set-exit-code > report.xml

- name: Publish test results
  uses: EnricoMi/publish-unit-test-result-action@v2
  if: always()
  with:
    files: report.xml
```

### 1.5 Add Go vulnerability check

```yaml
- name: Check for vulnerabilities
  run: |
    go install golang.org/x/vuln/cmd/govulncheck@latest
    govulncheck ./...
```

### 1.6 Fix -p 1 in ALL test commands

Every `go test ./...` in every workflow must have `-p 1`. Current gaps:
- `ci.yml:47` — `go test ./... -race -coverprofile=coverage.out` → add `-p 1`
- `security-scanning.yml:220-221` — multiple `go test ./... -run Test*` → add `-p 1` to each

### 1.7 Add golangci-lint config

Create `.golangci.yml` at repo root:
```yaml
linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - gofmt
    - goimports
    - misspell
    - bodyclose
    - noctx
    - exportloopref

linters-settings:
  goimports:
    local-prefixes: github.com/aliasfoxkde/Atheon

issues:
  exclude-rules:
    - path: _test\.go
      linters: [errcheck]
```

---

## Section 2: Test Coverage Improvements

**Branch name:** `feat/coverage-improvements`  
**Current:** bundler 90.9%, cmd/atheon 96.0%, cmd/mcp 92.4%, core 95.5%, total 94.9%  
**Target:** 97%+ total

### 2.1 core/ignore.go — writeIgnoreSegment (25% coverage)

File: `core/ignore.go` function `writeIgnoreSegment` at line 112.  
Add test in `core/ignore_test.go`:
```go
func TestWriteIgnoreSegment(t *testing.T) {
    // Create a temp .atheonignore and write segments to it
    // Test both header write and pattern write paths
}
```

### 2.2 core/pattern_state.go — InitializePatternState (66.7%)

File: `core/pattern_state.go`. Add tests covering:
- Load from valid state file
- Load from corrupted state file (JSON decode error)  
- Load from missing file (graceful fallback)
- Save state and reload

### 2.3 core/bundle.go — loadBundle error paths

Add tests for:
- `loadBundle` with invalid gzip (already partial — cover the `json.Unmarshal` error path)
- `init()` with corrupt `~/.atheon/patterns.bundle` (should fall back to embedded)
- `EnablePattern`/`DisablePattern` with non-existent name

### 2.4 main.go — printFindings, printJSONFindings

Currently at 86.7% and 80%. Add test cases via `run()` for:
- Large findings output (many findings, stats display)
- JSON mode with findings (currently only clean-file tested)
- `--json` with `--categories` that filters to zero findings

### 2.5 cmd/mcp — handleCall remaining branches

Add tests for:
- `scan_env` tool
- `list_patterns` tool  
- `list_categories` tool
- Invalid tool name (returns -32601 method not found)

---

## Section 3: Pattern Additions

**Branch name:** `feat/patterns-batch-2`

### 3.1 PII patterns to add

```
community/pii/national-id.yaml      — Generic national ID (country-agnostic numeric pattern)
community/pii/dob-format.yaml       — Date of birth in YYYY-MM-DD, MM/DD/YYYY formats
community/pii/gender-field.yaml     — Explicit gender field names in forms/schemas
community/pii/health-record-id.yaml — MRN patterns (Medical Record Numbers)
community/pii/tax-id-ein.yaml       — US EIN: XX-XXXXXXX format
```

### 3.2 Secrets patterns to add

```
community/secrets/cloudflare-api-token.yaml  — CF_ prefix 40-char token
community/secrets/okta-api-token.yaml        — 00x format 40-char Okta token  
community/secrets/pagerduty-api-key.yaml     — 20-char alphanumeric PD key
community/secrets/heroku-api-key.yaml        — UUID-format Heroku key
community/secrets/travis-ci-token.yaml       — Travis CI token pattern
community/secrets/circleci-token.yaml        — CircleCI v2 token
community/secrets/sonarqube-token.yaml       — sqp_ or squ_ prefixed tokens
community/secrets/artifactory-token.yaml     — AKC or AP prefix tokens
community/secrets/firebase-api-key.yaml      — Firebase config apiKey field
community/secrets/vercel-token.yaml          — vercel token pattern
```

### 3.3 Cloud-native patterns to add

```
community/cloud-native/aws-arn.yaml              — AWS ARN format
community/cloud-native/gcp-project-id.yaml       — GCP project ID in URLs
community/cloud-native/azure-connection-string.yaml — DefaultEndpointsProtocol pattern
community/cloud-native/k8s-imagepullsecret.yaml  — imagePullSecrets reference
community/cloud-native/helm-secret-value.yaml    — Helm values with .secret. paths
```

### 3.4 Code Quality patterns to add

```
community/code-quality/sleep-in-test.yaml        — time.Sleep in test files
community/code-quality/fmt-println-prod.yaml     — fmt.Println in non-test Go files
community/code-quality/panic-in-handler.yaml     — panic() calls in HTTP handlers
community/code-quality/direct-sql-query.yaml     — Raw SQL without parameterization context
community/code-quality/global-variable.yaml      — Package-level var declarations
community/code-quality/unused-import-comment.yaml — Blank import with no explanation
```

### 3.5 New category: `compliance`

```
community/compliance/gdpr-personal-data-comment.yaml — Comments mentioning personal data storage
community/compliance/hipaa-phi-field.yaml            — Field names matching PHI categories
community/compliance/pci-cardholder-data.yaml        — PAN/CVV/expiry in combination
community/compliance/data-retention-comment.yaml     — "never delete", "keep forever" comments
```

### 3.6 New category: `git-hygiene`

```
community/git-hygiene/merge-conflict-marker.yaml    — <<<<<<<, =======, >>>>>>> markers
community/git-hygiene/fixup-commit-message.yaml     — "fixup!" or "squash!" commit markers in code
community/git-hygiene/rebase-todo-leftover.yaml     — "pick", "reword", "edit" rebase instructions
community/git-hygiene/git-rerere-conflict.yaml      — rerere conflict markers
```

---

## Section 4: Code Quality & Architecture

**Branch name:** `refactor/code-quality`

### 4.1 Add `context.Context` throughout

`core/runner.go` — `ScanFile` and `ScanDir` accept `ctx context.Context` but don't check `ctx.Err()` during the directory walk. Add cancellation checks:
```go
if err := ctx.Err(); err != nil {
    return nil, nil, err
}
```
inside the `filepath.WalkDir` callback.

### 4.2 Extract pattern validation helper

`bundler/main.go:bundle()` validates regex inline. Extract to `core/pattern.go` as `ValidatePattern(def PatternDef) error` so both the bundler and tests can use the same validation.

### 4.3 Add structured logging

Replace `fmt.Fprintf(os.Stderr, ...)` calls with a simple leveled logger. Use `log/slog` (stdlib, Go 1.21+):
```go
slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn})))
```

### 4.4 Pattern metadata enrichment

Add optional fields to `PatternDef`:
```go
type PatternDef struct {
    Name        string   `json:"name"`
    Category    string   `json:"category"`
    Match       string   `json:"match"`
    Enabled     bool     `json:"enabled"`
    Description string   `json:"description,omitempty"`  // NEW
    Severity    string   `json:"severity,omitempty"`      // NEW: critical/high/medium/low/info
    References  []string `json:"references,omitempty"`    // NEW: CWE, CVE, docs URLs
    Tags        []string `json:"tags,omitempty"`          // NEW: freeform tags
}
```

Update YAML schema, bundler, and CLI `list` output to display severity and description.

### 4.5 Fix Finding output consistency

`core/runner.go` — `Finding.Line` is 1-indexed but sometimes reported as 0 when scanning in-memory strings. Add a guard:
```go
if f.Line == 0 { f.Line = 1 }
```

### 4.6 Add `atheon validate` command

New CLI subcommand that validates a YAML pattern file:
```
atheon validate community/secrets/my-pattern.yaml
```
Checks: valid RE2 regex, required fields present, name uniqueness against loaded bundle.

---

## Section 5: Documentation

**Branch name:** `docs/comprehensive`

### 5.1 API documentation

`docs/api/README.md` is referenced but may be sparse. Add:
- Full `core` package API reference (all exported functions with signatures and examples)
- `Finding` struct fields documented
- `PatternDef` wire format documented
- Bundle download/update flow documented

### 5.2 Pattern authoring guide

Create `docs/guides/PATTERN_AUTHORING.md`:
- RE2 regex restrictions (no lookahead/lookbehind, no backreferences)
- Naming conventions: `category-specific-name` (all lowercase hyphenated)
- Test case requirements
- Severity levels definition
- False positive minimization techniques
- The `enabled: false` use case (opt-in patterns)

### 5.3 Development setup guide

`docs/development/SETUP.md` — add:
- On this machine: proxy blocks `go install` for external tools; workaround via GONOSUMCHECK or GOPROXY direct
- Tool installation script: `./scripts/install-hooks.sh`
- How to add tools embedded in the repo (vendor or embed in scripts/)

### 5.4 Architecture decision records (ADRs)

Create `docs/architecture/decisions/` with:
- `ADR-001-re2-regex.md` — Why RE2 (Go stdlib) rather than PCRE
- `ADR-002-gzip-bundle.md` — Why gzip+JSON bundle rather than embedded YAML
- `ADR-003-parallel-tests.md` — Why -p 1 is required (global state in bundle init)

### 5.5 CHANGELOG.md

Create `CHANGELOG.md` at repo root following Keep a Changelog format. Document:
- Current version enhancements
- Pattern additions
- CI/CD improvements

---

## Section 6: Security Hardening

**Branch name:** `feat/security-hardening`

### 6.1 Add SARIF output format

For GitHub Security tab integration, add `--sarif` output flag that emits SARIF 2.1.0 JSON. Upload in CI:
```yaml
- name: Upload SARIF results
  uses: github/codeql-action/upload-sarif@v4
  with:
    sarif_file: results.sarif
```

### 6.2 Rate limiting in MCP server

`cmd/mcp/main.go` — no rate limiting on scan requests. Add a simple token bucket:
```go
var limiter = rate.NewLimiter(rate.Every(time.Second), 10) // 10 req/sec
```

### 6.3 Path traversal prevention

`core/runner.go:ScanDir` — validate that the resolved path doesn't escape the provided root via symlink traversal. Use `filepath.EvalSymlinks` and check prefix.

### 6.4 Input size limits

Add configurable max file size (default 10MB) to prevent memory exhaustion on large binary files accidentally included in a scan.

---

## Section 7: Performance

**Branch name:** `feat/performance`

### 7.1 Parallel pattern matching

Currently `scanLines` applies each pattern sequentially per line. For large files with many patterns (223+), this is O(lines × patterns). Parallelize pattern matching per line using a worker pool — safe because patterns are read-only after `init()`.

### 7.2 Pattern caching per process

Pre-compile all patterns on first use (currently done) but add a fast-path for disabled patterns:
```go
if !p.enabled { continue }  // already exists — verify it's before regex.FindString
```

### 7.3 Streaming API

`ScanFile` currently reads the entire file into memory. For large files, add a streaming line-by-line reader that processes chunks:
```go
scanner := bufio.NewScanner(f)
scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
for scanner.Scan() { ... }
```
(Note: this may already exist in the codebase as a `chunkedScan` function — verify before implementing.)

### 7.4 Bundle load time

Current bundle parse is sequential. For 223+ patterns, add parallel regex compilation:
```go
var wg sync.WaitGroup
for i := range defs {
    wg.Add(1)
    go func(def PatternDef) {
        defer wg.Done()
        // compile regex
    }(defs[i])
}
wg.Wait()
```

---

## Section 8: Tooling Infrastructure

**Branch name:** `feat/tooling`

### 8.1 Embed goimports and staticcheck

Since the corporate proxy blocks `go install`, embed pre-built binaries or use a Makefile target that downloads from GitHub releases directly:

```makefile
GOIMPORTS_VERSION = v0.22.0
STATICCHECK_VERSION = 2024.1.1

tools/goimports:
	GOBIN=$(PWD)/tools go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)

tools/staticcheck:
	GOBIN=$(PWD)/tools go install honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION)
```

Add `tools/` to `.gitignore` and update `.githooks/pre-commit` to check `./tools/goimports` before `goimports` in PATH.

### 8.2 Add Makefile

Create `Makefile` with targets:
```makefile
.PHONY: build test lint bundle setup clean

build:
	go build -o bin/atheon ./cmd/atheon
	go build -o bin/atheon-mcp ./cmd/mcp

test:
	go test ./... -p 1 -timeout 15m -coverprofile=coverage.out
	go tool cover -func=coverage.out | grep total:

lint:
	go vet ./...
	./tools/staticcheck ./... 2>/dev/null || staticcheck ./... 2>/dev/null || true
	gofmt -l . | xargs -r false

bundle:
	go run ./bundler

setup:
	git config core.hooksPath .githooks
	mkdir -p tools
	GOBIN=$(PWD)/tools go install golang.org/x/tools/cmd/goimports@latest || true
	GOBIN=$(PWD)/tools go install honnef.co/go/tools/cmd/staticcheck@latest || true

clean:
	rm -rf bin/ coverage.out
```

### 8.3 Update install-hooks.sh to use local tools

Modify `scripts/install-hooks.sh` to install tools to `./tools/` (repo-local) rather than global GOPATH:
```bash
TOOLS_DIR="$(git rev-parse --show-toplevel)/tools"
mkdir -p "$TOOLS_DIR"
GOBIN="$TOOLS_DIR" go install golang.org/x/tools/cmd/goimports@latest 2>/dev/null && echo "  ✓ goimports → tools/" || echo "  ⚠ goimports blocked (try: make setup)"
GOBIN="$TOOLS_DIR" go install honnef.co/go/tools/cmd/staticcheck@latest 2>/dev/null && echo "  ✓ staticcheck → tools/" || echo "  ⚠ staticcheck blocked"
```

---

## Section 9: Branch Protection & CODEOWNERS

**Branch name:** `docs/governance`

### 9.1 Verify CODEOWNERS is active

`.github/CODEOWNERS` exists with `* @aliasfoxkde`. Verify by checking GitHub repo Settings → Code and Automation → Branches → Branch protection rules for `main`:
- ✅ Require a pull request before merging
- ✅ Require approvals: 1
- ✅ Require review from Code Owners
- ✅ Require status checks to pass before merging: ci/test, ci/build, ci/lint
- ✅ Require branches to be up to date before merging
- ✅ Do not allow bypassing the above settings
- ✅ Restrict force pushes: disabled
- ✅ Restrict deletions: enabled

### 9.2 Document branch protection setup

Create `docs/governance/BRANCH_PROTECTION.md` with:
- Exact settings to configure
- Why each rule exists
- How the AI should create PRs (never force-push, never --no-verify)
- The review/approval flow

---

## Implementation Order (Recommended)

1. **Section 8** (Tooling — Makefile, tool install) — unblocks everything else
2. **Section 1.6 + Section 2** (CI -p 1 fixes + Coverage) — ensures tests stay green
3. **Section 3** (Patterns) — high value, low risk
4. **Section 1.1–1.5** (CI consolidation + SHA pinning) — reduces maintenance debt
5. **Section 4** (Code quality) — architectural improvements
6. **Section 5** (Docs) — parallel with any section
7. **Section 6** (Security) — SARIF + rate limiting
8. **Section 7** (Performance) — measure before optimizing
9. **Section 9** (Governance docs) — documentation-only

---

## Success Criteria (A+ in every category)

| Category | Current | Target | Key Actions |
|----------|---------|--------|-------------|
| Test Coverage | 94.9% | 97%+ | Section 2 |
| Pattern Count | 223 | 250+ | Section 3 |
| CI/CD Quality | B | A+ | Sections 1, 6 |
| Documentation | B+ | A+ | Sections 5, 9 |
| Security | B | A+ | Sections 1.2, 6 |
| Code Quality | A | A+ | Sections 4, 8 |
| Performance | unknown | measured | Section 7 |

---

*Generated: 2026-06-22 | Branch: fix/audit-comprehensive-v2 | For: AI implementation handoff*
