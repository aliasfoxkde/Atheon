# Pattern Expansion: Issue #149

## Context
Expand community library based on real-world repo benchmarking via Atheon-GitHub-Scanner and Atheon-Benchmark projects.

## Reference Projects

### Atheon-GitHub-Scanner (`/nas/Temp/repos/Atheon-GitHub-Scanner`)
Scans real GitHub repositories to find security issues and identify pattern gaps.

### Atheon-Benchmark (`/nas/Temp/repos/Atheon-Benchmark`)
Benchmarks AI code generation using Atheon's pattern scanning as quality gates:
- Uses **185+ patterns** from Atheon bundle
- Validates code across **8 categories**
- Tests AI outputs for security and quality issues
- Reference: `/nas/Temp/repos/Atheon-Benchmark/dashboard/lib/atheon/quality-gates.ts`

## Benchmark Results Summary

### From Atheon-GitHub-Scanner (5 repos, 55 findings)
- 2 patterns validated with accuracy >85%:
  - **API Key Exposure in Configuration Files** (CWE-798, accuracy 85%)
  - **SQL Injection via String Concatenation** (CWE-89, accuracy 88%)

### From Atheon-Benchmark Quality Gates
The benchmark's quality gates validate code against 185+ patterns across 8 categories.
Patterns from the validated fallback corpus not yet in community library were identified:
- `test-skip` - Detects skipped tests (describe.skip, it.skip, etc.)
- `ai-generated-content` - Detects AI-generated content markers

## Patterns Added

### From Atheon-GitHub-Scanner Benchmark
1. `community/code-quality/sql-injection-string-concat.yaml`
   - Detects SQL injection via string concatenation
   - CWE-89, severity: high
   - Source: awesome-python-project benchmark

2. `community/secrets/generic-api-key-config.yaml`
   - Detects hardcoded API keys, secrets, tokens in config
   - CWE-798, severity: critical
   - Source: popular-javascript-lib benchmark

### From Atheon-Benchmark Quality Gates
3. `community/code-quality/test-skip.yaml`
   - Detects skipped tests (describe.skip, it.skip, test.skip, context.skip)
   - Category: code-quality
   - Source: Atheon-Benchmark fallback patterns

4. `community/code-quality/ai-generated-content.yaml`
   - Detects AI-generated content markers
   - Category: code-quality
   - Source: Atheon-Benchmark fallback patterns

## Validation

All patterns validated against:
- Go `regexp.Compile` syntax check
- Bundle loading and registration
- Pattern matching on test inputs
- `go test ./...` passes

## Pattern Count

| Source | Patterns |
|--------|----------|
| Community library (before) | 58 |
| Added from GitHub-Scanner | 2 |
| Added from Benchmark | 2 |
| **Total** | **62** |

## Source Data References

- Atheon-GitHub-Scanner pipeline: `/nas/Temp/repos/Atheon-GitHub-Scanner/pipeline_results.json`
- Combined scan results: `/nas/Temp/repos/Atheon-GitHub-Scanner/data/combined_scan_results.json`
- Atheon-Benchmark quality gates: `/nas/Temp/repos/Atheon-Benchmark/dashboard/lib/atheon/quality-gates.ts`
- Atheon-Benchmark fallback patterns: `/nas/Temp/repos/Atheon-Benchmark/dashboard/lib/claude/atheon-integration.ts`

---

**Date:** 2026-06-22
**Branch:** `pr/149-patterns-expansion`
**Sources:**
- Atheon-GitHub-Scanner (5 repos, 55 findings, 2 validated patterns)
- Atheon-Benchmark (185+ patterns, 8 categories, 2 fallback patterns added)
