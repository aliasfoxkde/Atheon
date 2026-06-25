# ADR-002: Gzip+JSON Bundle Format

**Status**: Accepted
**Date**: 2026-06-17

## Context

We needed a format to store and distribute pattern definitions. Options considered:

- **Embedded YAML**: Simple but large file size
- **JSON**: Compact but no compression
- **Gzip+JSON**: Best of both - compact and fast to parse
- **Protocol Buffers**: Faster parsing but complex schema

## Decision

We use gzip-compressed JSON (`patterns.bundle`).

## Rationale

1. **Small size**: Gzip compresses ~90% on typical pattern bundles
2. **Fast parsing**: Single decompress + JSON decode
3. **Human-readable**: Can be inspected with `zcat` or `gunzip`
4. **No schema**: JSON with dynamic parsing
5. **Go stdlib**: `compress/gzip` and `encoding/json` are built-in

## Bundle Structure

```json
[
  {
    "name": "pattern-name",
    "category": "category-dir",
    "match": "regex-pattern",
    "enabled": true
  },
  ...
]
```

Then gzip-compressed into `patterns.bundle`.

## Consequences

### Positive
- Small download size (~10KB vs ~100KB uncompressed)
- Fast to parse (one decompress pass)
- Easy to generate and debug

### Negative
- Not directly human-readable (need to decompress)
- Need bundler tool to create bundles

---

**References**:
- [Go compress/gzip](https://pkg.go.dev/compress/gzip)
- [Go encoding/json](https://pkg.go.dev/encoding/json)