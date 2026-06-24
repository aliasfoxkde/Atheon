# Atheon Enhanced Fork - System Architecture & Workflow

## 🏗️ Repository Structure

This document serves as the complete reference for the aliasfoxkde/Atheon enhanced fork architecture, workflow, and operational procedures.

### **🎯 Mission Statement**
Maintain an enhanced version of Atheon that provides production-ready features while maintaining full compatibility with upstream HoraDomu/Atheon.

## 📋 Branch Strategy

### **Core Branches**

#### **`stable/clean`** (Upstream Tracking)
- **Purpose**: Source of truth tracking upstream HoraDomu/Atheon:main
- **Update Strategy**: Automatic sync via GitHub Actions daily
- **Usage**: Baseline for all development, reference for upstream changes
- **Protection**: Protected branch, only maintainers can push
- **Workflow**:
  ```bash
  git checkout stable/clean
  git pull upstream main
  git push origin stable/clean
  ```

#### **`main`** (Production Build)
- **Purpose**: Production-ready build with all validated enhancements
- **Update Strategy**: Merge validated feature PRs + periodic stable/clean sync
- **Usage**: User-facing installation via `go install github.com/aliasfoxkde/Atheon`
- **Protection**: Protected branch, requires PR review and passing CI
- **Features**: Enhanced patterns, performance optimizations, MCP integration

#### **`dev/full-feature`** (Development Branch)
- **Purpose**: Development branch with ALL patterns enabled for comprehensive testing
- **Update Strategy**: Continuous integration of feature branches
- **Usage**: Internal testing, validation, pattern development
- **Features**: Full pattern suite, experimental features, comprehensive test suite

### **Feature Branches**

#### **Naming Convention**
- `feat/feature-name`: New features
- `perf/performance-improvement`: Performance optimizations
- `patterns/category-expansion`: Pattern library additions
- `docs/documentation-updates`: Documentation improvements
- `infra/infrastructure-updates`: CI/CD, tooling, infrastructure

## 🔄 Development Workflow

### **Feature Development Process**

```
1. Branch Creation
   └─ git checkout -b feat/feature-name stable/clean

2. Development & Testing
   ├─ Implementation
   ├─ Local testing (go test ./...)
   ├─ Pattern validation (atheon list --all)
   └─ Documentation updates

3. Create Pull Request
   ├─ Target: main branch
   ├─ Title: feat: descriptive feature name
   ├─ Description: Implementation details, testing results, impact analysis
   └─ Labels: appropriate category labels

4. CI/CD Validation
   ├─ All tests must pass (multi-version Go testing)
   ├─ Lint checks must pass (staticcheck, golangci-lint)
   ├─ Security scanning (Atheon self-scan)
   ├─ Performance benchmarks
   └─ Code coverage requirements (≥45%)

5. Code Review
   ├─ Peer review required
   ├─ Documentation review
   └─ Architecture impact assessment

6. Merge to main
   ├─ Auto-merge when all checks pass
   ├─ Version bump in main branch
   └─ Changelog update

7. Periodic sync
   └─ stable/clean → main merge to get upstream fixes
```

### **Upstream Integration Process**

```bash
# Sync with upstream Atheon
git checkout stable/clean
git fetch upstream
git merge upstream/main
# Resolve conflicts if any
git push origin stable/clean

# Update main with upstream changes
git checkout main
git merge stable/clean
# Test integration
git push origin main
```

## ⚙️ Configuration Profiles

### **Profile Structure**
Located in `config/profiles/`:

#### **`production.json`** (Default)
```json
{
  "enabled_categories": ["secrets", "pii", "security"],
  "strict_mode": "standard",
  "performance_mode": "optimized",
  "exit_on_findings": true,
  "max_file_size_mb": 10
}
```

#### **`pipeline.json`** (CI/CD Optimized)
```json
{
  "enabled_categories": ["secrets", "security", "code-quality"],
  "strict_mode": "strict",
  "performance_mode": "fast",
  "exit_on_findings": true,
  "json_output": true,
  "max_file_size_mb": 50
}
```

#### **`mcp-integration.json`** (MCP Server)
```json
{
  "enabled_categories": ["all"],
  "strict_mode": "standard",
  "performance_mode": "balanced",
  "exit_on_findings": false,
  "streaming_mode": true
}
```

#### **`development.json`** (Full Feature Testing)
```json
{
  "enabled_patterns": "all",
  "strict_mode": "strict",
  "performance_mode": "comprehensive",
  "exit_on_findings": false,
  "debug_mode": true
}
```

## 🚨 Error Prevention & Quality Gates

### **Pre-commit Hooks**
- **Author Verification**: Ensures correct user attribution (aliasfoxkde)
- **Commit Format**: Enforces conventional commits
- **Build Validation**: Go build must succeed
- **Test Coverage**: Minimum 45% coverage threshold
- **Formatting**: gofmt and goimports checks

### **Pre-push Hooks**
- **Binary Validation**: Both atheon and atheon-mcp must build
- **Functionality Tests**: Basic command validation
- **Common Issues**: Trailing whitespace, debug prints, etc.

### **CI/CD Pipeline Checks**
- **Multi-version Testing**: Go 1.21, 1.22, 1.23, 1.24
- **Static Analysis**: golangci-lint, staticcheck, go vet
- **Security Scanning**: Atheon self-scan on codebase
- **Performance Benchmarks**: Regression detection
- **Pattern Validation**: All patterns must load and function
- **Documentation**: Docs must be updated for user-facing changes

## 🔧 Decision Making Framework

### **Workflow Decisions**

**Merge Strategy**: **Merge Commits** (not rebase)
- **Reasoning**: Preserves history, easier conflict resolution, clear audit trail
- **Exception**: stable/clean uses rebase for linear upstream tracking

**Release Cadence**: **Continuous Delivery with Monthly Tags**
- **Reasoning**: Users get features quickly, predictable release schedule
- **Process**: Automatic tagging when main reaches stable state

**Version Naming**: **Semantic Versioning with Enhancement Suffix**
- **Format**: `v1.2.3-enhanced`
- **Reasoning**: Clear differentiation from upstream, compatibility tracking

**CI/CD Priority**: **All Improvements Equal Priority**
- **Security**: Critical (blocking)
- **Performance**: Critical (regression blocking)
- **Coverage**: Important (minimum threshold)
- **Documentation**: Important (user-facing changes required)

## 📁 Repository Organization

### **Documentation Structure**
```
docs/
├── SYSTEM_ARCHITECTURE.md (this file)
├── FEATURE_COMPARISON.md (upstream vs enhanced)
├── DEVELOPMENT_WORKFLOW.md (detailed dev guide)
├── CONFIGURATION_GUIDE.md (profile usage)
├── BRANCH_STRATEGY.md (branch documentation)
├── TESTING_GUIDE.md (testing procedures)
└── CHANGELOG.md (version history)
```

### **Code Organization**
```
core/                    (Core scanning engine — no external deps)
├── patterns.bundle      (Embedded gzip+JSON pattern database; //go:embed)
├── pattern.go           (Pattern interface, registry, ValidatePattern, sentinel errors)
├── bundle.go            (Bundle load, enable/disable, DownloadBundle, SetActiveCategories)
├── runner.go            (ScanFile, ScanDir, ScanString, ScanEnv — the public surface)
├── ignore.go            (.atheonignore, .gitignore compilation and matching)
├── pattern_state.go     (Persisted enabled/disabled state in ~/.atheon/pattern_state.json)
└── finding.go           (Finding, Stats result types)

cmd/
├── atheon/              (CLI binary — flags, subcommands, JSON/SARIF/human output)
└── mcp/                 (MCP server — stdio JSON-RPC, rate-limited tool dispatch)

bundler/                 (YAML → gzip+JSON bundle compiler; invoked via `go run ./bundler`)

config/
└── profiles/            (User-facing configuration profiles: production, pipeline,
                          mcp-integration, development)

community/               (Pattern contributions; one YAML per pattern, category = directory)
├── secrets/             (58 patterns — credential/API-key/token detection)
├── code-quality/        (35 patterns — maintenance and quality issues)
├── accessibility/       (19 patterns — WCAG, ARIA, keyboard nav)
├── security-hardening/  (18 patterns — auth, crypto, CSRF, XSS, injection)
├── web-security/        (15 patterns — web-app vulnerabilities)
├── cloud-native/        (14 patterns — Docker, K8s, Terraform, serverless)
├── performance/         (12 patterns — N+1, caching, lazy-loading)
├── web-development/     (12 patterns — React/Next.js, TypeScript, bundling)
├── pii/                 (11 patterns — emails, IDs, health records)
├── ai-detection/        (9 patterns — AI-generated code markers)
├── api-integration/     (9 patterns — REST/GraphQL, auth, rate-limit, error-handling)
├── devops/              (9 patterns — CI/CD, Docker, GitHub workflows)
├── healthcare/          (7 patterns — HIPAA, NPI, PHI fields)
├── finance/             (6 patterns — payment, routing, account numbers)
├── pwa/                 (5 patterns — service workers, manifests, offline)
├── data-visualization/  (5 patterns — charts, color, mobile)
├── compliance/          (4 patterns — GDPR, HIPAA, PCI, retention)
├── git-hygiene/         (4 patterns — conflict markers, fixup commits)
└── frameworks/          (3 patterns across django/, nodejs/, react/ subdirs)

scripts/
├── pattern-count.sh     (Single source of truth for pattern counts)
├── install-hooks.sh     (Wire up pre-commit/pre-push hooks)
├── build.sh             (Local build helper)
├── coverage.sh          (Generate coverage report)
└── hooks/               (Pre-commit and pre-push hook scripts)
```

> **Verified 2026-06-23**: 255 patterns across 19 categories. Run
> `./scripts/pattern-count.sh` to regenerate the table above; the engine and bundler
> tolerate this drift without recompilation of `core/` because patterns are loaded at
> init time.

## 🎯 Success Metrics

### **Quality Indicators**
- **Test Coverage**: Maintain ≥45% (target: 60%+)
- **Pattern Count**: 105+ patterns (target: 200+)
- **False Positive Rate**: <5% (target: <2%)
- **CI/CD Pass Rate**: >95% (target: 99%+)
- **Upstream Sync Lag**: <24 hours (target: <6 hours)

### **User Satisfaction**
- **Installation Success**: >95% first-time installs
- **Pattern Accuracy**: >90% actionable findings
- **Documentation Clarity**: >80% find answers in docs
- **Performance**: <1s for 1000-file repositories

## 🚨 Common Issues & Solutions

### **Merge Conflicts**
- **Prevention**: Keep feature branches focused, update from stable/clean regularly
- **Resolution**: Use theirs strategy for upstream changes, ours for features
- **Escalation**: Document unresolvable conflicts in SYSTEM_ISSUES.md

### **CI/CD Failures**
- **Debug**: Check logs for specific failure points
- **Local Validation**: Run `go test ./...` and `go vet ./...` locally
- **Pattern Issues**: Validate patterns with `atheon list --all`

### **Performance Regressions**
- **Detection**: Benchmark suite alerts on >10% degradation
- **Investigation**: Profile with `go test -bench`
- **Resolution**: Optimize before merge to main

### **Documentation Gaps**
- **Prevention**: Pre-commit hook checks for user-facing changes
- **Detection**: User feedback and issue tracking
- **Resolution**: Update docs before merge

## 🔄 Continuous Improvement

### **Weekly Reviews**
- Branch status and cleanup
- CI/CD performance and optimization
- Documentation currency and accuracy
- Upstream changes assessment

### **Monthly Assessments**
- Architecture and workflow review
- Technology stack evaluation
- Security audit and enhancement
- User feedback integration

### **Quarterly Planning**
- Feature roadmap alignment
- Performance optimization targets
- Documentation restructuring
- Tool and process improvements

---

## 📞 Contact & Contribution

### **For Issues**
- **Technical Problems**: Open GitHub issue
- **Documentation Issues**: Submit PR with improvements
- **Feature Requests**: Discuss in issues first

### **Contribution Guidelines**
- Follow branch strategy documentation
- Respect workflow and quality gates
- Update documentation for all changes
- Test thoroughly before submission

### **Maintainer Responsibilities**
- Review and merge PRs promptly
- Monitor CI/CD health and performance
- Keep documentation current and accurate
- Maintain upstream compatibility

---

**This document is the single source of truth for all architectural and operational decisions. Any deviations must be documented and approved.**