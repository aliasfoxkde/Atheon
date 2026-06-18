# Atheon Pattern Format Specification

## Pattern Definition

Every pattern in Atheon is defined as a YAML file in the `community/` directory.

## File Structure

```
community/
├── secrets/
│   ├── aws.yaml
│   └── github.yaml
├── pii/
│   ├── creditcard.yaml
│   └── ssn.yaml
└── code-quality/
    ├── todo-comment.yaml
    └── debug-println.yaml
```

## Pattern Specification

### Required Fields

**name** (string)
- Pattern identifier
- Must be unique across all patterns
- Should be lowercase with hyphens
- Should be descriptive and specific

**category** (string)
- Determined by directory name
- Valid categories: secrets, pii, code-quality, healthcare, finance
- Automatically extracted from directory path

**match** (string)
- Valid RE2 regex pattern
- Use single quotes to avoid backslash escaping
- Should match the specific pattern you're detecting

### Optional Fields

**enabled** (boolean)
- Default: true
- Set to false to disable pattern by default
- Can be enabled at runtime via CLI

## Pattern Examples

### Basic Pattern

```yaml
# community/secrets/aws.yaml
name: aws-access-key-id
match: '\bAKIA[0-9A-Z]{16}\b'
```

### Disabled Pattern

```yaml
# community/secrets/internal-key.yaml
name: internal-api-key
match: '\bINTERNAL_[A-Z0-9]{32}\b'
enabled: false
```

### Complex Pattern

```yaml
# community/secrets/gcp-service-account.yaml
name: gcp-service-account-key
match: '"private_key_id": "[A-Za-z0-9]{40}"'
```

## Category Guidelines

### secrets
- API keys, tokens, credentials
- Service principal credentials
- Database connection strings
- CI/CD tokens
- Container registry credentials

### pii
- Personal identifiable information
- Credit card numbers
- Social Security numbers
- Phone numbers
- Email addresses

### code-quality
- Debug statements
- TODO/FIXME comments
- Deprecated functions
- Hardcoded credentials
- Code smells

### healthcare
- Patient identifiers
- Medical record numbers
- Prescription numbers
- Insurance numbers
- Medical licenses

### finance
- IBANs
- ABA routing numbers
- SWIFT/BIC codes
- Financial account numbers

## Naming Conventions

### Pattern Names
- Use lowercase with hyphens
- Be specific and descriptive
- Examples:
  - `stripe-live-key` (good)
  - `stripe` (too vague)
  - `StripeKey` (wrong case)

### Category Names
- Use lowercase with hyphens
- Plural when appropriate
- Examples:
  - `code-quality` (good)
  - `CodeQuality` (wrong case)

## Regex Guidelines

### RE2 Syntax
- Use RE2 regex syntax (not PCRE)
- Avoid backreferences and lookarounds
- Use word boundaries (`\b`) when appropriate
- Be specific to avoid false positives

### Performance
- Avoid catastrophic backtracking
- Use character classes `[a-z]` instead of `(?:a|b|c)`
- Avoid nested quantifiers
- Test performance for complex patterns

### False Positives
- Include test cases for both matches and non-matches
- Consider edge cases
- Document expected formats
- Balance sensitivity vs specificity

## Pattern Testing

### Test Cases

Every pattern must include test cases in `core/bundle_test.go`:

```go
{
    name: "pattern-name",
    matches:    []string{"valid_match_1", "valid_match_2"},
    nonMatches: []string{"invalid_1", "invalid_2"},
}
```

### Testing Guidelines

1. **Positive Cases**: Include 2+ examples that should match
2. **Negative Cases**: Include 2+ examples that should NOT match
3. **Edge Cases**: Test boundary conditions
4. **Real Data**: Test with real-world examples
5. **Performance**: Ensure regex is efficient

## Common Patterns

### API Keys
```yaml
name: example-api-key
match: '\bEXAMPLE_[A-Za-z0-9]{32}\b'
```

### UUIDs
```yaml
name: generic-uuid
match: '\b[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\b'
```

### Base64
```yaml
name: base64-string
match: '(?:[A-Za-z0-9+/]{4}){2,}(?:[A-Za-z0-9+/]{2}==)?'
```

### Email
```yaml
name: generic-email
match: '[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}'
```

## Validation

Before committing a pattern:

1. **Test locally**: `go test ./...`
2. **Manual testing**: `atheon --file <test-file>`
3. **Bundle generation**: `go run ./bundler community core/patterns.bundle`
4. **Pattern validation**: Verify pattern is loaded correctly
5. **False positive check**: Test against real code

## Best Practices

1. **Start Specific**: Make patterns as specific as possible
2. **Test Thoroughly**: Include comprehensive test cases
3. **Document Edge Cases**: Note any limitations or false positives
4. **Consider Performance**: Avoid expensive regex operations
5. **Review Existing Patterns**: Check for overlap before creating new patterns
6. **Use Appropriate Categories**: Choose the right category for your pattern

## Pattern Review Process

Maintainers review patterns for:
- Correctness (regex accuracy)
- False positive rate
- Name clarity
- Overlap with existing patterns
- Performance impact
- Security considerations
- Test coverage