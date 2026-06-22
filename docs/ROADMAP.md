# Roadmap

Atheon's trajectory is simple: a larger, more authoritative pattern library across more domains. The engine is stable. What grows from here are the patterns.

---

## Current state (June 2026)

- **58 patterns** across 5 categories: secrets, pii, code-quality, healthcare, finance
- CLI with category filtering, enable/disable with persistent state, JSON output, stdin piping
- MCP server (`atheon-mcp`) for Claude, Cursor, and Windsurf integration
- Git hook support (pre-commit, pre-push)
- CI/CD integration with native binaries for Windows, macOS, and Linux
- Automated releases on the 10th and 21st of each month via GoReleaser

---

## Near term

**More patterns in existing categories**

The existing categories have obvious gaps. Priority additions:

- `secrets` — more SaaS and cloud provider API key formats, OAuth tokens, JWT secrets
- `pii` — email addresses, passport numbers, driver's license numbers, IP addresses
- `code-quality` — hardcoded passwords, magic numbers, commented-out code blocks
- `finance` — credit card CVV, more regional routing number formats
- `healthcare` — NPI numbers, DEA numbers, HIPAA-relevant identifiers

**New categories from the community**

Domains actively accepting pattern contributions:

- `legal` — prohibited contract terms, restricted clause patterns
- `operations` — log anomaly signatures, error codes, pipeline failure markers
- `gaming` — anti-cheat string signatures, profanity filters, moderation patterns

---

## Medium term

**Pattern quality improvements**

As the library grows, false positive rate matters more. Plans include:

- Confidence metadata on patterns (`confidence: high/medium/low`) so users can filter by signal strength
- Context-aware matching — anchor patterns to specific file types or variable name prefixes
- Pattern deprecation workflow for retiring patterns that have been superseded

**Tooling**

- `atheon check <pattern-file>` — validate a YAML file before bundling: regex safety, naming, overlap check
- Improved `atheon update` diff — show pattern names and categories for new additions

---

## Long term

**Platform vision**

The long-term goal is Atheon as a community platform: a searchable, browseable library of patterns contributed by practitioners across every domain. Think of it as a package registry for detection rules — where a security engineer, a HIPAA compliance officer, and a game moderation team can each find patterns built by people who understand their domain.

**What "platform" means in practice:**

- Web-browseable pattern index so users can discover patterns without installing the CLI
- Pattern search by category, keyword, and domain
- Contributor profiles and attribution
- Community voting on pattern quality and false positive rate

---

## What's out of scope

Atheon is a pattern matching engine. It will not become:

- A SIEM or alerting platform
- A managed scanning service or SaaS product
- A dependency scanner or SBOM tool
- A secrets rotation or remediation tool

The engine stays minimal. The patterns are the product.

---

## Contributing

Every pattern on the roadmap above is open for contribution right now. See [CONTRIBUTING.md](../.github/CONTRIBUTING.md) to add yours.
