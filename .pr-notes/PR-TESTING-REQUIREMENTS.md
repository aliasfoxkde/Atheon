# Atheon Testing Requirements and PR Validation Standards

**Branch:** docs/pr-testing-requirements
**Purpose:** Define comprehensive testing requirements for all PR submissions

---

## Mandatory Testing Requirements

### Pre-Submission Testing Checklist ✅

All PRs MUST include the following testing before submission:

- [ ] **Unit Tests** - Add unit tests for new functionality
- [ ] **Integration Tests** - Test integration with existing systems
- [ ] **Compatibility Tests** - Test backward compatibility
- [ ] **Build Tests** - Verify project builds on all platforms
- [ ] **Pattern Tests** - Test new/modified patterns with test cases
- [ ] **Edge Cases** - Test boundary conditions and error scenarios
- [ ] **Performance Tests** - Ensure no performance regressions

### Breaking Changes Detection 🔍

**Automated Breaking Changes Detection:**
- Semantic versioning compliance checks
- API compatibility verification
- Pattern format validation
- Configuration migration testing

**Manual Breaking Changes Review:**
1. Pattern format changes require version bump
2. CLI flag changes require migration documentation
3. Core API changes require deprecation period
4. Output format changes require documentation update

---

## Platform Compatibility Testing

### Required Platform Testing

**All PRs must be tested on:**

**Operating Systems:**
- Linux (Ubuntu 20.04+, Debian 11+)
- macOS (12+ Monterey+)
- Windows (10+, Server 2019+)

**Architectures:**
- amd64 (x86_64)
- arm64 (aarch64)

**Go Versions:**
- Go 1.24+ (current stable)
- Go 1.23 (minimum supported)

### Compatibility Test Commands

```bash
# Test build on all platforms
GOOS=linux GOARCH=amd64 go build
GOOS=linux GOARCH=arm64 go build
GOOS=darwin GOARCH=amd64 go build
GOOS=darwin GOARCH=arm64 go build
GOOS=windows GOARCH=amd64 go build

# Test functionality
go test ./...
./atheon list categories
./atheon scan ./test_data/
```

---

## Pattern Testing Requirements

### New Pattern Submission Template

**Every new pattern MUST include:**

1. **Pattern Definition:**
```yaml
name: pattern-name
match: 'regex pattern'
# enabled: false  # optional
```

2. **Test Case in core/bundle_test.go:**
```go
{
    name: "pattern-name",
    matches:    []string{"example_match_1", "example_match_2"},
    nonMatches: []string{"should_not_match_1", "should_not_match_2"},
}
```

3. **Documentation:**
- Pattern purpose and use case
- False positive considerations
- Performance notes if complex regex

### Pattern Validation Process

**Before PR submission:**
1. Test pattern against real-world data
2. Verify regex correctness using regex tester
3. Check for catastrophic backtracking
4. Measure performance impact
5. Document edge cases

**Automated validation:**
```bash
# Run bundler to generate new bundle
go run ./bundler community core/patterns.bundle

# Run pattern tests
go test ./core -run TestRegisteredPatterns

# Verify no compilation errors
go build
```

---

## Integration Testing Standards

### MCP Server Testing

**atheon-mcp binary MUST be tested for:**

1. **Protocol Compliance:**
   - MCP JSON-RPC 2.0 specification
   - Tool definition format validation
   - Response format verification

2. **Tool Functionality:**
   - scan_string tool with various content types
   - scan_file tool with different file types
   - scan_dir tool with category filtering

3. **Error Handling:**
   - Invalid tool names
   - Missing parameters
   - File not found scenarios
   - Permission denied handling

### CLI Testing Requirements

**CLI commands MUST be tested for:**

1. **Basic Commands:**
   - `atheon list`
   - `atheon list categories`
   - `atheon scan <directory>`
   - `atheon --file <file>`
   - `atheon --env`

2. **Flag Combinations:**
   - `--json` with various commands
   - `--categories=<cat1,cat2>` filtering
   - `--all` flag usage

3. **Edge Cases:**
   - Empty directories
   - Files with no matches
   - Binary file handling
   - Permission errors

---

## Performance Testing Standards

### Performance Benchmarks

**Before PR submission, verify:**

1. **No Performance Regression:**
   ```bash
   cd /tmp/atheon-projects/Atheon-Benchmark
   go run bench.go --compare
   ```

2. **Memory Usage:**
   - Peak memory < 50MB for typical usage
   - No memory leaks in long-running processes
   - Efficient garbage collection

3. **File Processing Speed:**
   - < 100ms for 1000-line files
   - < 1s for 10,000-line files
   - Linear scaling for file size

### Performance Regression Detection

**Automated Checks:**
```bash
# Run performance comparison
go test -bench=. -benchmem > before.txt
# Make changes
go test -bench=. -benchmem > after.txt
# Compare results
benchcmp before.txt after.txt
```

---

## Breaking Changes Policy

### What Constitutes a Breaking Change

**Requires version bump (MAJOR version increment):**
- Pattern YAML format changes
- CLI interface changes
- API signature modifications
- Output format changes
- Configuration file format changes

**Does NOT require version bump (MINOR version):**
- Adding new patterns
- Adding new CLI flags
- New tools to MCP server
- Performance improvements
- Bug fixes

### Breaking Changes Process

1. **Documentation:**
   - Document breaking change clearly
   - Provide migration guide
   - Update version compatibility matrix

2. **Deprecation Period:**
   - Maintain old functionality for at least one version
   - Add deprecation warnings
   - Document sunset timeline

3. **Testing:**
   - Add tests for deprecated functionality
   - Ensure migration path works
   - Test both old and new interfaces

---

## Test Coverage Requirements

### Minimum Coverage Standards

**All PRs MUST maintain or improve:**

- **Overall Coverage:** ≥ 35.8%
- **Core Package:** ≥ 59%
- **New Code:** ≥ 80% coverage
- **Critical Paths:** 100% coverage

### Coverage Enforcement

**Automated checks:**
```bash
# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total:

# Check minimum coverage
go test ./... -cover | grep coverage: | awk '{if ($4 < 35.8) print "FAIL"; exit 1}'

# Check specific packages
go test ./core -cover | grep coverage: | awk '{if ($4 < 59) print "FAIL"; exit 1}'
```

---

## PR Testing Workflow

### Before Submitting PR

1. **Run Full Test Suite:**
   ```bash
   go test ./... -v
   go test ./... -race
   go test ./... -cover
   ```

2. **Build Verification:**
   ```bash
   go build -v
   go run ./bundler community core/patterns.bundle
   ```

3. **Manual Testing:**
   ```bash
   # Test basic functionality
   ./atheon list categories
   ./atheon scan ./test_data/
   ./atheon --file README.md
   ```

4. **Platform Testing:**
   ```bash
   # Test on multiple platforms if possible
   GOOS=linux go test ./...
   GOOS=windows go test ./...
   ```

5. **Performance Verification:**
   ```bash
   cd /tmp/atheon-projects/Atheon-Benchmark
   go run bench.go
   ```

### Automated PR Checks

**Pre-commit hooks will run:**
- Go build validation
- Test execution (must pass)
- Coverage threshold checks
- Pattern format validation
- Breaking change detection

**CI pipeline will run:**
- Full test suite on all platforms
- Integration tests
- Performance benchmarks
- Security scanning

---

## Test Documentation Requirements

### Test Case Documentation

**All tests MUST include:**

1. **Clear Purpose:** What is being tested and why
2. **Setup Requirements:** What test data is needed
3. **Expected Behavior:** What should happen
4. **Error Conditions:** How errors should be handled

### Example Test Documentation

```go
// TestScanFile tests the file scanning functionality
// It verifies that Atheon can scan a single file and detect patterns
// Expected: Returns findings array with pattern matches
// Edge cases: File not found, permission denied, binary files
func TestScanFile(t *testing.T) {
    // Setup
    tmpDir := t.TempDir()
    testFile := filepath.Join(tmpDir, "test.txt")
    content := "AKIAIOSFODNN7EXAMPLE"

    if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
        t.Fatal(err)
    }

    // Test
    findings, stats, err := core.ScanFile(testFile)
    if err != nil {
        t.Fatalf("ScanFile failed: %v", err)
    }

    // Verify
    if stats.Files != 1 {
        t.Errorf("expected 1 file, got %d", stats.Files)
    }

    // Edge case: empty file
    emptyFile := filepath.Join(tmpDir, "empty.txt")
    findings, stats, err = core.ScanFile(emptyFile)
    if err != nil {
        t.Errorf("empty file scan should not fail: %v", err)
    }
}
```

---

## Quality Gates

### Automated Quality Checks

**All PRs must pass:**

1. **Build Verification:** Project builds cleanly on all platforms
2. **Test Success:** All tests pass (no skips allowed)
3. **Coverage Threshold:** Minimum 35.8% overall coverage maintained
4. **Performance Baseline:** No significant performance regression
5. **Pattern Validation:** All patterns compile and load correctly
6. **Breaking Changes:** Properly documented and versioned

### Manual Review Checklist

**Reviewers must verify:**

- [ ] Test coverage is adequate for changes
- [ ] Breaking changes are properly documented
- [ ] Performance impact is acceptable
- [ ] Security implications are considered
- [ ] Documentation is updated
- [ ] Backward compatibility is maintained

---

## Continuous Testing

### Automated Testing Infrastructure

**CI/CD Pipeline:**
- Triggered on every PR submission
- Runs on Linux, macOS, Windows
- Executes full test suite
- Generates coverage reports
- Runs performance benchmarks

### Monitoring and Alerts

**Quality Metrics:**
- Test pass rate: Must be 100%
- Coverage trends: Must not decrease
- Performance baselines: Alert on regressions
- Build success rate: Monitor for platform-specific issues

---

## Enforcement

### Pre-Merge Requirements

**PR cannot be merged unless:**

1. All tests pass on all platforms
2. Coverage threshold is met
3. Performance benchmarks pass
4. Breaking changes are documented
5. Manual review checklist completed

### Post-Merge Monitoring

**After merge, monitor for:**

1. CI/CD pipeline failures
2. Coverage decreases
3. Performance regressions
4. User-reported issues
5. Security vulnerabilities

---

## Related Documentation

- [PR Documentation](.pr-notes/branches-summary.md) - PR submission guidelines
- [Best Practices](docs/pr-documentation-and-notes/PR_DOCUMENTATION.md) - Development standards
- [Testing Guide](.github/CONTRIBUTING.md) - Contribution requirements
- [Performance Targets](/tmp/atheon-projects/Atheon-Benchmark/RESULTS.md) - Performance benchmarks

---

**Last Updated:** 2025-06-18
**Enforcement:** Mandatory for all PR submissions
**Owner:** aliasfoxkde
**Status:** Active and enforced