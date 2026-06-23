# ADR-001: RE2 Regex Choice

**Status**: Accepted
**Date**: 2026-06-17

## Context

We needed to choose a regex engine for pattern matching in Atheon. Options considered:

- **PCRE/PCRE2**: Powerful but not available in Go stdlib, adds C dependency
- **RE2**: Go stdlib, safe (no backtracking), consistent performance
- **Custom**: Too much work, reinventing the wheel

## Decision

We use RE2 (via `regexp` package in Go's stdlib).

## Rationale

1. **No external dependencies**: RE2 is in Go's standard library
2. **Guaranteed linear time**: No ReDoS attacks possible
3. **Consistent performance**: O(n) regardless of input
4. **Thread-safe**: Patterns can be shared across goroutines
5. **Well-maintained**: Part of Go itself

## Consequences

### Positive
- Single dependency (Go itself)
- Safe patterns by default
- Fast compilation and matching

### Negative
- No lookahead/lookbehind assertions
- No backreferences
- Some complex patterns not possible

## Alternatives Considered

### PCRE2
- Would allow more complex patterns
- Requires CGO or external dependency
- Potential ReDoS vulnerability
- **Rejected**: Complexity and security concerns

### RE2 (C library)
- Same issues as PCRE2
- **Rejected**: Go's implementation is sufficient

---

**References**:
- [RE2 syntax](https://github.com/google/re2/wiki/Syntax)
- [Go regexp package](https://pkg.go.dev/regexp)