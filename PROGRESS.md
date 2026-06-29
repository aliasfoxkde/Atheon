# Progress

## 2026-06-29

### AST Pattern Enhancements (feature/ast-patterns-v2)

**Major Enhancements:**

1. **Fixed and Enhanced Built-in AST Patterns** (`core/ast_patterns.go`)
   - Removed Python-specific patterns (getattr, exec, eval) that were incorrectly mixed with Go patterns
   - Expanded from 9 to 28 Go-specific AST security patterns
   - Fixed detection logic (token comparison bug with `+` operator)
   - Added patterns for: SSRF, template injection, LDAP injection, XXE, ReDoS, weak crypto, and more

2. **CLI Integration** (`cmd/atheon/main.go`)
   - Added `--ast` flag to enable AST scanning
   - Added `--ast-only` flag for AST-only scanning (skip regex)
   - Updated help text

3. **Core Integration** (`core/runner.go`)
   - Added `ScanOpts.EnableAST` and `ScanOpts.ASTOnly` options
   - Added `ScanFileWithAST()` function
   - Added `ScanDirWithAST()` function

4. **Documentation**
   - Extended `docs/PATTERN_FORMAT.md` with AST pattern documentation
   - Created `community/ast-security/README.md`
   - Updated `community/README.md` with ast-security category

5. **Tests** (`core/ast_patterns_test.go`)
   - Comprehensive tests for all AST patterns
   - Tests for command injection, SQL injection, path traversal, SSRF, weak crypto, etc.

**Pattern Count: 28 built-in AST patterns**

### See Also
- [docs/PLAN.md](docs/PLAN.md) - Implementation plan
- [docs/PATTERN_FORMAT.md](docs/PATTERN_FORMAT.md) - Pattern format documentation
