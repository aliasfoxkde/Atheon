# Pattern Authoring Guide

This guide explains how to create effective patterns for Atheon.

## Pattern File Format

Patterns are YAML files with the following structure:

```yaml
name: pattern-name
match: 'regex pattern'
enabled: true  # optional, defaults to true
```

## Naming Conventions

- **File name**: `category-pattern-name.yaml` (all lowercase, hyphenated)
- **Pattern name**: `category-pattern-name` (must be unique across all patterns)
- **Category**: directory name (e.g., `secrets`, `pii`, `code-quality`)

## Regex Requirements

Atheon uses RE2 regex (Go's standard library). Key restrictions:

### Allowed
- Character classes: `[a-z]`, `[0-9]`, `[A-Z]`
- Quantifiers: `*`, `+`, `?`, `{n}`, `{n,m}`
- Anchors: `^`, `$`, `\b`
- Groups: `(?:...)`, `(...)`
- Alternation: `|`
- Escape sequences: `\d`, `\w`, `\s`, `\.`, `\\`

### NOT Allowed (RE2 limitations)
- Lookahead: `(?=...)`, `(?!...)`
- Lookbehind: `(?<=...)`, `(?<!...)`
- Backreferences: `\1`, `\2`
- Word boundaries with unicode: use `\b` carefully

## Example Patterns

### Simple Secret Pattern
```yaml
name: my-api-key
match: '\bMY_API_KEY_[A-Z0-9]{32}\b'
```

### Complex Pattern with Groups
```yaml
name: aws-access-key
match: '\bAKIA[0-9A-Z]{16}\b'
```

### PII Pattern
```yaml
name: email-address
match: '\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b'
```

## Severity Levels (Future)

Pattern metadata will include severity levels:
- `critical` - Immediate security risk
- `high` - Significant risk
- `medium` - Moderate risk
- `low` - Minor issue
- `info` - Informational

## False Positive Minimization

1. **Use word boundaries** (`\b`) to avoid partial matches
2. **Be specific** with character classes (e.g., `[A-Z]` not `.`)
3. **Anchor patterns** where appropriate (`^` for start, `$` for end)
4. **Test against real code** to verify matches

## Testing Patterns

1. Create the YAML file in the appropriate `community/` subdirectory
2. Run `go run ./bundler` to rebuild the bundle
3. Test with `./atheon --file test-file.txt`
4. Verify no false positives on clean code

## Pattern Categories

Place patterns in the appropriate category directory:

| Category | Description |
|----------|-------------|
| `secrets` | API keys, tokens, passwords |
| `pii` | Personal identifying information |
| `code-quality` | Code smells, anti-patterns |
| `security-hardening` | Security misconfigurations |
| `cloud-native` | Cloud infrastructure issues |
| `compliance` | Regulatory compliance concerns |
| `git-hygiene` | Git workflow issues |

## Adding New Categories

1. Create the directory: `community/new-category/`
2. Add the category name to `core/pattern_test.go`'s `validCategories`
3. Create patterns following the naming conventions
4. Rebuild and test