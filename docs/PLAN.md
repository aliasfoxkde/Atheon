# AST Pattern Implementation Plan

## Status: Phase 1-4 Implemented ✅

## Overview

Add AST-based pattern analysis using Go's standard `go/ast` library to detect security issues that regex cannot find.

## Phases

### Phase 1: Core AST Scanner ✅
- [x] Create `core/ast_patterns.go`
- [x] Define `ASTFinding` and `ASTPattern` types
- [x] Implement `ScanFileAST()` and `ScanDirAST()`
- [x] Pattern registration (builtin patterns)
- [x] Added `ToFinding()` for integration with main Finding type

### Phase 2: Built-in Patterns ✅
- [x] `go-command-injection` - exec.Command with string concat
- [x] `go-shell-command` - Shell invocation with user input
- [x] `go-sql-injection` - String concat in query
- [x] `go-sql-template-query` - Query method with user input
- [x] `go-path-traversal` - os.Open with user input
- [x] `go-symlink-attack` - File open without O_NOFOLLOW
- [x] `go-unsafe-deserialization` - encoding.Unmarshal with user input
- [x] `go-gob-deserialization` - gob decoding with untrusted data
- [x] `go-ssrf` - HTTP request with user-controlled URL
- [x] `go-http-unvalidated-url` - http.Get with user URL
- [x] `go-template-injection` - Template with user input
- [x] `go-template-raw-html` - template.HTML bypasses escaping
- [x] `go-redos` - Regex with nested quantifiers
- [x] `go-regex-dynamic` - regexp.Compile with user pattern
- [x] `go-hardcoded-secret` - Credential with string literal
- [x] `go-private-key` - Embedded private key
- [x] `go-weak-crypto-md5` - Use of MD5
- [x] `go-weak-crypto-sha1` - Use of SHA-1
- [x] `go-insecure-random` - math/rand for security
- [x] `go-weak-cipher` - DES, RC4, ECB mode
- [x] `go-unchecked-error` - Unchecked error return
- [x] `go-silent-panic` - panic in production code
- [x] `go-ldap-injection` - LDAP query string concat
- [x] `go-xxe` - XML without external entity protection
- [x] `go-yaml-unsafe` - yaml.Unmarshal with untrusted data
- [x] `go-trust-boundary` - User input to internal state
- [x] `go-tls-skip-verify` - TLS verification disabled
- [x] `go-insecure-tls` - Weak TLS configuration

### Phase 3: CLI Integration ✅
- [x] Add `--ast` flag to `atheon scan`
- [x] Add `--ast-only` flag to skip regex
- [x] Update help text

### Phase 4: Core Integration ✅
- [x] Add `ScanOpts.EnableAST` and `ScanOpts.ASTOnly`
- [x] Add `ScanFileWithAST()` function
- [x] Add `ScanDirWithAST()` function
- [x] AST findings flow into main Finding type
- [x] SARIF/JSON output works automatically

### Phase 5: Documentation ✅
- [x] Extend PATTERN_FORMAT.md with AST section
- [x] Add `community/ast-security/` directory with README
- [x] Update `community/README.md`

## Files Created/Modified

| File | Change |
|------|--------|
| `core/ast_patterns.go` | Rewrite - 28 Go-specific AST patterns |
| `core/ast_patterns_test.go` | Rewrite - comprehensive tests |
| `core/runner.go` | Add ScanFileWithAST, ScanDirWithAST, ScanOpts |
| `cmd/atheon/main.go` | Add --ast, --ast-only flags |
| `docs/PATTERN_FORMAT.md` | Add AST documentation section |
| `community/ast-security/README.md` | NEW - category documentation |
| `community/README.md` | Add ast-security category |

## Testing

- [x] `go test ./core/ -run TestAST` - All pass
- [x] `go test ./... -p 1` - Full test suite passes
- [x] Manual testing with `atheon --ast`

## Rollout

1. CLI flag `--ast` enables AST scanning ✅
2. AST patterns in separate category (ast-security) ✅
3. Performance: AST scanning is fast for Go files ✅

## Future Work

- [ ] YAML-based AST pattern definition (Phase 6)
- [ ] Taint tracking for complex data flows
- [ ] Multi-language AST support (Python, JavaScript)
- [ ] AST pattern editor/visualizer
