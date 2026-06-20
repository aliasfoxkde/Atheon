# Installation Instructions for Enhanced Atheon v2.0.0

## 🚀 Go Install Method (Recommended)

```bash
# Install latest enhanced version
go install github.com/aliasfoxkde/Atheon@v2.0.0-enhanced

# Or install from main branch (latest commits)
go install github.com/aliasfoxkde/Atheon@latest

# Install MCP server
go install github.com/aliasfoxkde/Atheon/cmd/mcp@v2.0.0-enhanced
```

## 🛠️ Build from Source

```bash
# Clone repository
git clone https://github.com/aliasfoxkde/Atheon.git
cd Atheon

# Checkout specific version
git checkout v2.0.0-enhanced

# Build binaries
go build -o atheon ./cmd/atheon
go build -o atheon-mcp ./cmd/mcp

# Install to PATH
sudo mv atheon /usr/local/bin/
sudo mv atheon-mcp /usr/local/bin/
```

## ✅ Verify Installation

```bash
# Check version
atheon --version

# Check pattern count (should be 87)
atheon list | wc -l

# List categories
atheon list categories
```

## 🧪 Quick Test

```bash
# Test basic functionality
atheon README.md --categories=secrets

# List all patterns
atheon list

# Scan current directory
atheon . --categories=ai-detection,secrets
```

## 📚 Documentation

- **Complete Documentation**: https://github.com/aliasfoxkde/Atheon/tree/main/docs
- **API Reference**: https://github.com/aliasfoxkde/Atheon/blob/main/docs/api/README.md
- **Development Guide**: https://github.com/aliasfoxkde/Atheon/blob/main/docs/development/SETUP.md
- **Troubleshooting**: https://github.com/aliasfoxkde/Atheon/blob/main/docs/guides/TROUBLESHOOTING.md

## 📦 Package Information

- **Module**: github.com/aliasfoxkde/Atheon
- **Version**: v2.0.0-enhanced
- **Go Version**: 1.21+
- **Patterns**: 87 total patterns
- **Categories**: 10 categories
- **Platforms**: Linux, macOS, Windows

## 🔧 Platform-Specific Notes

### Linux
```bash
# Install Go dependencies
sudo apt install golang-go

# Build and install
go build -o atheon ./cmd/atheon
sudo mv atheon /usr/local/bin/
```

### macOS
```bash
# Install Go using Homebrew
brew install go

# Build and install
go build -o atheon ./cmd/atheon
sudo mv atheon /usr/local/bin/
```

### Windows
```bash
# Install Go from https://go.dev/dl/

# Build
go build -o atheon.exe

# Add to PATH or move to system directory
```

## ⚡ After Installation

### Remove Local Pattern Cache
If you have used Atheon before, remove the cached bundle:
```bash
rm -f ~/.atheon/patterns.bundle
```

### Verify Pattern Loading
```bash
# Should show 87 patterns
atheon list | wc -l

# Should show 10 categories
atheon list categories
```

## 🎯 Next Steps

1. **Configure Atheon** for your needs
2. **Explore patterns**: `atheon list` to see all 87 patterns
3. **Set up CI/CD integration**
4. **Contribute patterns** if you have custom needs

---

**Enhanced Atheon**: https://github.com/aliasfoxkde/Atheon
**Release**: v2.0.0-enhanced
**Date**: 2026-06-19
**Maintainer**: Micheal Kinney (aliasfoxkde)
