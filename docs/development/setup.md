# Development Setup Guide

Complete guide for setting up a development environment for contributing to Enhanced Atheon.

## Prerequisites

### Required Software

- **Go**: 1.21 or later (supports 1.21, 1.22, 1.23, 1.24)
- **Git**: Latest stable version
- **Make**: Optional, for build automation
- **Docker**: Optional, for containerized development

### Recommended Tools

- **VS Code** or **GoLand**: IDE with Go support
- **golangci-lint**: Comprehensive Go linting
- **staticcheck**: Advanced static analysis
- **jq**: JSON processing for CI/CD workflows

## 🛠️ Environment Setup

### 1. Clone Repository

```bash
# Clone your fork
git clone https://github.com/aliasfoxkde/Atheon.git
cd Atheon

# Add upstream remote
git remote add upstream https://github.com/HoraDomu/Atheon.git
```

### 2. Go Module Setup

```bash
# Verify Go installation
go version
# Expected: go version go1.21+ linux/amd64

# Download dependencies
go mod download

# Verify module integrity
go mod verify
```

### 3. Build from Source

```bash
# Build main binary
go build -o atheon .

# Build MCP server
go build -o atheon-mcp ./core/mcp

# Build bundler
go build -o bundler/bundler ./bundler

# Verify installation
./atheon --version
./atheon list | wc -l
# Expected: 87 patterns
```

### 4. Development Tools

```bash
# Install linting tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install honnef.co/go/tools/cmd/staticcheck@latest

# Install Go tools
go install golang.org/x/tools/cmd/goimports@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
```

## 🧪 Testing Setup

### Run Test Suite

```bash
# Run all tests
go test ./... -v

# Run with coverage
go test ./... -coverprofile=coverage.out

# View coverage report
go tool cover -html=coverage.out
```

### Test Specific Packages

```bash
# Test core package
go test ./core -v

# Test CLI
go test . -v

# Test MCP server
go test ./core/mcp -v
```

### Pre-commit Hooks

The repository includes comprehensive pre-commit validation:

```bash
# Pre-commit hooks are automatically installed in .git/hooks/
# They run:
# - Go formatting checks
# - go vet static analysis
# - Test execution with 54.4% coverage requirement
# - Author validation (Micheal Kinney)
```

## 📝 Development Workflow

### 1. Create Feature Branch

```bash
# Start from stable baseline
git checkout stable/clean
git checkout -b feat/my-feature

# Or from main for enhancements
git checkout main
git checkout -b feat/my-feature
```

### 2. Development Process

```bash
# Make your changes
# ... code changes ...

# Run tests
go test ./... -v

# Check formatting
go fmt ./...

# Run linters
golangci-lint run
staticcheck ./...

# Build to verify
go build -o atheon .
```

### 3. Pattern Development

If adding new patterns:

```bash
# Create pattern YAML
cat > community/new-category/pattern-name.yaml <<'EOF'
name: pattern-name
match: 'regex_pattern_here'
EOF

# Rebuild bundle
./bundler/bundler

# Rebuild binary
go build -o atheon .

# Test new pattern
./atheon list | grep pattern-name
./atheon --categories=new-category .
```

### 4. Commit Changes

```bash
# Stage changes
git add .

# Commit with conventional format
git commit -m "feat: add new pattern category for X"

# Push to your fork
git push origin feat/my-feature
```

### 5. Create Pull Request

```bash
# Create PR via GitHub CLI
gh pr create --base main --head feat/my-feature
```

## 🔧 Development Configuration

### IDE Setup

**VS Code:**
```json
{
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "package",
  "go.testFlags": ["-v"],
  "go.testTimeout": "30s",
  "go.coverOnSave": false,
  "go.coverOnTestPackage": true
}
```

**GoLand:**
- Enable Go Modules integration
- Configure golangci-lint as external tool
- Set up test coverage visualization

### Git Configuration

Required author configuration:
```bash
git config user.name "Micheal Kinney"
git config user.email "micheal.l.c.kinney@gmail.com"
```

## 🐛 Debugging Setup

### Debug Build

```bash
# Build with debug symbols
go build -gcflags="all=-N -l" -o atheon-debug .

# Run with debug output
./atheon-debug --debug-level=2 scan .
```

### Test Debugging

```bash
# Run specific test with verbose output
go test -v -run TestScanFile

# Run with race detector
go test -race -v

# Test with coverage profiling
go test -coverprofile=coverage.out -covermode=count
go tool cover -func=coverage.out
```

### Performance Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof
go tool pprof mem.prof
```

## 📊 Quality Checks

### Pre-commit Validation

The pre-commit hook ensures:

1. **Author Validation**: Correct git user configuration
2. **Code Formatting**: Proper Go formatting
3. **Static Analysis**: go vet passes
4. **Test Coverage**: Minimum 54.4% coverage
5. **Documentation**: Docs updated for user-facing changes

### Manual Quality Checks

```bash
# Full test suite
go test ./... -race -coverprofile=coverage.out

# Linting
golangci-lint run --timeout=5m

# Security scanning
./atheon . --all

# Bundle validation
./bundler/bundler
```

## 🚀 CI/CD Integration

### Local CI Simulation

```bash
# Run comprehensive testing
go test ./... -v -race -coverprofile=coverage.out

# Check coverage threshold
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
echo "Coverage: $COVERAGE%"
if [ "$COVERAGE" -lt 54 ]; then
    echo "❌ Coverage below threshold"
    exit 1
fi
```

### Matrix Testing

Test across multiple Go versions:
```bash
# Test with Go 1.21
go test ./... -v

# Test with Go 1.24
go test ./... -v

# Or use go version matrix
for version in "1.21" "1.22" "1.23" "1.24"; do
    echo "Testing with Go $version"
    go test ./... -v
done
```

## 📖 Development Resources

### Core Documentation

- [System Architecture](SYSTEM_ARCHITECTURE.md) - Technical architecture
- [Branch Strategy](BRANCH_STRATEGY.md) - Git workflow
- [Pattern Development](patterns/development.md) - Pattern creation guide
- [API Documentation](api/README.md) - Programmatic API

### External Resources

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Modules Reference](https://go.dev/ref/mod)
- [Atheon Upstream](https://github.com/HoraDomu/Atheon)

## 🎯 Development Best Practices

### Code Quality

1. **Follow Go Conventions**
   - Use standard Go formatting
   - Implement proper error handling
   - Write clear, descriptive comments
   - Use meaningful variable names

2. **Testing Requirements**
   - Maintain 54.4%+ coverage
   - Write comprehensive unit tests
   - Include integration tests for features
   - Test edge cases and error conditions

3. **Documentation Standards**
   - Document exported functions
   - Provide usage examples
   - Update user guides for user-facing changes
   - Maintain inline comments for complex logic

### Pattern Development

1. **Pattern Quality**
   - Low false positive rate
   - Efficient regex compilation
   - Clear category assignment
   - Comprehensive testing

2. **Pattern Documentation**
   - Clear pattern purpose
   - Usage examples
   - Known limitations
   - Performance characteristics

### Security Considerations

1. **Safe File Handling**
   - Validate file paths
   - Handle permission errors
   - Prevent directory traversal
   - Manage resource limits

2. **Input Validation**
   - Sanitize user input
   - Validate regex patterns
   - Handle malformed data
   - Prevent injection attacks

## 🔍 Troubleshooting

### Common Issues

**Build Failures:**
```bash
# Clean Go cache
go clean -cache -mod

# Rebuild dependencies
go mod tidy
go mod download
```

**Test Failures:**
```bash
# Update go modules
go mod tidy

# Run specific failing test
go test -v -run TestName

# Check for race conditions
go test -race -v
```

**Pattern Loading Issues:**
```bash
# Rebuild bundle
./bundler/bundler

# Remove local bundle cache
rm ~/.atheon/patterns.bundle

# Rebuild binary
go build -o atheon .
```

---

**Setup Guide Version**: 1.0.0
**Last Updated**: 2026-06-19
**Maintainer**: Micheal Kinney (aliasfoxkde)