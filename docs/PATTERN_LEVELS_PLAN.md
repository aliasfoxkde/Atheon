# Pattern Levels and Systematic Expansion Plan

## Overview

This document outlines the plan for implementing severity levels and strictness modes for Atheon patterns, along with a systematic approach to expanding the pattern library.

## 1. Pattern Severity Levels

### Proposed Levels

- **Critical**: Security credentials and high-risk secrets
  - Examples: API keys, tokens, private keys, passwords
  - Action: Immediate remediation required

- **High**: Potential security vulnerabilities
  - Examples: SQL injection vectors, eval usage, weak crypto
  - Action: Review and fix required

- **Medium**: Code quality and maintainability issues
  - Examples: Dead code, deprecated functions, TODO comments
  - Action: Review recommended

- **Low**: Style and minor quality issues
  - Examples: Commented code, placeholder text, formatting
  - Action: Optional cleanup

- **Info**: Informational markers
  - Examples: Fake data, AI detection patterns, stubs
  - Action: Awareness only

### Implementation Approach

1. **Add `severity` field to pattern YAML**:
   ```yaml
   name: aws-access-key
   match: '\b(?:AKIA|ASIA)[0-9A-Z]{16}\b'
   severity: critical
   ```

2. **Migrate existing patterns**:
   - secrets/: → critical
   - code-quality/: → medium/low
   - devops/: → medium
   - ai-detection/: → info
   - pii/: → high

3. **Update bundle structure**:
   - Add severity to PatternDef struct
   - Include in JSON bundle output
   - Display in scan results

## 2. Strictness Modes

### Proposed Modes

- **Strict Mode**: All patterns enabled, including false positives
  - Use case: CI/CD security scans, audit preparation
  - Default: Include all patterns regardless of severity

- **Standard Mode**: Balanced approach (default)
  - Use case: Development workflows
  - Default: Critical + High + Medium patterns

- **Relaxed Mode**: High-confidence security patterns only
  - Use case: Large legacy codebases, initial scans
  - Default: Critical patterns only

- **Custom Mode**: User-specified severity threshold
  - Use case: Tailored workflows
  - CLI flag: `--severity-level high`

### Implementation Approach

1. **Add CLI flags**:
   ```bash
   --mode strict|standard|relaxed
   --severity-level critical|high|medium|low|info
   ```

2. **Filter patterns by severity**:
   - Update pattern loading logic
   - Add severity filtering functions

3. **Mode configuration**:
   - Add to config file support
   - Allow per-project mode setting

## 3. Systematic Pattern Expansion

### Current Pattern Count
- **Total**: 88 patterns
- **Secrets**: 12 patterns
- **Code Quality**: 24 patterns
- **DevOps**: 6 patterns
- **AI Detection**: 4 patterns
- **PII**: 10 patterns
- **Frameworks**: 8 patterns
- **Finance**: 6 patterns
- **Healthcare**: 4 patterns
- **Writing**: 0 patterns (planned)

### Expansion Targets

#### Phase 1: API Keys & Secrets (Target: +20 patterns)
- [ ] GitHub App keys
- [ ] GitLab tokens
- [ ] Bitbucket tokens
- [ ] Azure keys
- [ ] Twilio tokens
- [ ] SendGrid keys
- [ ] PagerDuty tokens
- [ ] Datadog keys
- [ ] New Relic keys
- [ ] Splunk tokens
- [ ] Shopify keys
- [ ] Square tokens
- [ ] PayPal tokens
- [ ] Auth0 tokens
- [ ] Okta tokens
- [ ] Firebase tokens
- [ ] Slack webhooks
- [ ] Discord tokens
- [ ] Telegram tokens
- [ ] Zoom tokens

#### Phase 2: Code Quality (Target: +30 patterns)
- [ ] Go-specific patterns (10 more)
  - Goroutine without context
  - Missing error returns
  - Interface not checked
  - Pointer vs value confusion
  - Race condition patterns
  - Memory leak patterns
  - Import side effects
  - Global variable usage
  - Unsafe usage
  - Reflection usage
- [ ] JavaScript/TypeScript patterns (10)
  - var vs const/let
  - Promise without catch
  - Missing await
  - Prototype pollution
  - XSS vectors
  - eval alternatives
  - this binding issues
  - Closure issues
  - Event listener leaks
  - DOM manipulation
- [ ] Python patterns (10)
  - Mutable default args
  - Import side effects
  - Exception handling
  - Generator exhaustion
  - Thread safety
  - Memory leaks
  - Global state
  - Dynamically created types
  - Unpickled data
  - Subprocess injection

#### Phase 3: DevOps & Infrastructure (Target: +15 patterns)
- [ ] Terraform patterns (5)
  - Hardcoded credentials
  - Missing tags
  - Resource leaks
  - State issues
  - Provider versions
- [ ] Ansible patterns (5)
  - No-op changes
  - Undefined variables
  - Vault usage
  -become usage
  - Loop issues
- [ ] CI/CD patterns (5)
  - Artifactory tokens
  - Jenkins credentials
  - GitLab CI variables
  - CircleCI env vars
  - Bitbucket pipelines

#### Phase 4: Security Vulnerabilities (Target: +20 patterns)
- [ ] OWASP Top 10 coverage (10)
  - Command injection
  - LDAP injection
  - XML injection
  - Path traversal
  - insecure deserialization
  - SSRF
  - CSRF
  - RCE patterns
  - File inclusion
  - Open redirect
- [ ] Cryptography issues (5)
  - Hardcoded IV
  - ECB mode
  - Padding oracle
  - Key reuse
  - Random number generation
- [ ] Authentication issues (5)
  - Hardcoded sessions
  - JWT without verification
  - Basic auth
  - Digest auth weaknesses
  - OAuth implementation issues

#### Phase 5: AI Detection Enhancement (Target: +10 patterns)
- [ ] Code structure patterns
  - Over-commented code
  - Generic variable names
  - Excessive abstraction
  - Template-like structure
- [ ] Textual patterns
  - Transitional phrases
  - Passive voice overuse
  - Hedging language
  - Boilerplate introductions
  - Generic conclusions

#### Phase 6: Writing & Documentation Quality (Target: +25 patterns)
- [ ] Documentation issues (10)
  - Outdated version references
  - Broken link patterns
  - Missing examples
  - Placeholder content
  - Inconsistent terminology
  - Unclear instructions
  - Missing context
  - Vague descriptions
  - Incomplete steps
  - Dead content
- [ ] Writing style issues (10)
  - Passive voice overuse
  - Wordy constructions
  - Jargon without explanation
  - Inconsistent tone
  - Unclear antecedents
  - Run-on sentences
  - Ambiguous language
  - Weasel words
  - Redundant phrases
  - Cliché usage
- [ ] Code comment quality (5)
  - Misleading comments
  - Obvious comments
  - Commented-out code
  - Outdated comments
  - Comment vs code mismatch

## 4. Pattern Taxonomy & Best Practices

### Pattern Naming Convention
- Use kebab-case: `aws-access-key`
- Category prefix for specificity: `go-goroutine-leak`
- Avoid generic names: `token` → `stripe-api-token`

### Pattern Quality Standards
1. **Testability**: Pattern must have test cases
2. **Low false positives**: Minimize noise
3. **Clear remediation**: Explain why pattern is problematic
4. **Severity accuracy**: Match severity to actual risk
5. **Documentation**: Include usage examples

### Pattern Metadata Structure
```yaml
name: pattern-name
match: 'regex pattern'
severity: critical|high|medium|low|info
confidence: high|medium|low  # New: confidence level
category: category-name
description: Human-readable description
remediation: How to fix the issue
references:
  - https://example.com/docs
tags:
  - security
  - api-key
```

## 5. Migration Strategy

### Phase 1: Add Severity Field (Week 1)
1. Update pattern YAML schema to include severity
2. Migrate existing patterns with default severity
3. Update bundler to handle severity field
4. Update core loading logic

### Phase 2: Implement Modes (Week 2)
1. Add CLI flags for mode/severity
2. Implement pattern filtering
3. Update scan output to show severity
4. Add tests for mode switching

### Phase 3: Pattern Expansion (Weeks 3-8)
1. Systematically add patterns by category
2. Focus on high-value patterns first
3. Test each pattern against real codebases
4. Document patterns as they're added

### Phase 4: Documentation & Polish (Week 9-10)
1. Create pattern contribution guide
2. Add severity usage examples
3. Document mode selection guidelines
4. Create pattern review checklist

## 6. Success Metrics

### Quantitative
- Pattern count: 88 → 200+ patterns
- Category coverage: 8 → 15+ categories
- Test coverage: Maintain 50%+

### Qualitative
- False positive rate: < 5%
- Pattern remediation clarity: 90%+ actionable
- User satisfaction: Feedback-based improvement

## 7. Open Questions

1. Should AI detection patterns be opt-in by default?
2. How to handle patterns with different confidence levels?
3. Should patterns have regional/industry variants?
4. How to balance pattern granularity vs. noise?
5. Should we support custom pattern repositories?

## 8. Next Steps

1. [ ] Review and approve this plan
2. [ ] Create tasks for each phase
3. [ ] Start with Phase 1 (severity field implementation)
4. [ ] Begin pattern expansion with high-value targets
5. [ ] Establish pattern review process
