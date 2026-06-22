# Community Patterns

Pattern files for Atheon, organized by category. Each `.yaml` file defines one pattern.

## Categories

| Category | Patterns | Description |
|----------|----------|-------------|
| [accessibility](accessibility/) | 19 | WCAG compliance, ARIA, screen reader issues |
| [ai-detection](ai-detection/) | 9 | AI-generated code markers, LLM disclaimers |
| [api-integration](api-integration/) | 9 | API keys, webhook secrets, integration tokens |
| [cloud-native](cloud-native/) | 9 | Kubernetes, Terraform, Docker secrets |
| [code-quality](code-quality/) | 29 | Debug artifacts, hardcoded values, dead code |
| [data-visualization](data-visualization/) | 5 | Chart/graph library patterns |
| [devops](devops/) | 9 | CI/CD bypass markers, pipeline secrets |
| [finance](finance/) | 6 | Payment identifiers, financial data |
| [healthcare](healthcare/) | 7 | PHI, HIPAA-relevant field patterns |
| [performance](performance/) | 12 | Blocking calls, synchronous patterns |
| [pii](pii/) | 7 | Personally identifiable information |
| [pwa](pwa/) | 5 | Progressive Web App patterns |
| [secrets](secrets/) | 49 | API keys, tokens, credentials, private keys |
| [security-hardening](security-hardening/) | 18 | Insecure configs, weak crypto, unsafe calls |
| [web-development](web-development/) | 12 | Frontend anti-patterns |
| [web-security](web-security/) | 15 | XSS, SQLi, CORS, injection risks |

**Total: 225 patterns across 16 active categories**

## Pattern File Format

```yaml
name: pattern-name          # lowercase, hyphenated, unique
category: secrets           # must match directory name
match: "regex-pattern"      # RE2-compatible (no lookahead/lookbehind)
enabled: true               # false = opt-in only
description: "What it detects"
```

> **RE2 constraint**: Go uses the RE2 engine. No lookahead (`(?=...)`), lookbehind (`(?<=...)`), or backreferences (`\1`). Use prefix/suffix context in the pattern itself instead.

## Adding a Pattern

1. Create `community/<category>/<name>.yaml`
2. Run `go run ./bundler` to rebuild `core/patterns.bundle`
3. Run `go test ./... -p 1` to verify all tests pass
4. Open a PR — CI validates the bundle automatically

See [docs/guides/PATTERN_AUTHORING.md](../docs/guides/PATTERN_AUTHORING.md) if it exists, or the [CONTRIBUTING guide](../.github/CONTRIBUTING.md).
