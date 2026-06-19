#!/bin/bash
# Enhanced Git Hooks Setup Script for Atheon Repository

set -e

echo "=== 🔧 Setting Up Enhanced Atheon Git Hooks ==="

# Ensure we're in the Atheon repository
if [ ! -f "go.mod" ] || ! grep -q "github.com/aliasfoxkde/Atheon" go.mod; then
    echo "❌ Error: Not in Atheon repository root"
    echo "Run this script from the Atheon repository root directory"
    exit 1
fi

# Create hooks directory if it doesn't exist
mkdir -p .git/hooks

# Copy enhanced hooks
echo "Installing enhanced pre-commit hook..."
cp .githooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit

echo "Installing enhanced pre-push hook..."
cp .githooks/pre-push .git/hooks/pre-push
chmod +x .git/hooks/pre-push

echo "✅ Core hooks installed"

# Optional AI integration hook
echo ""
echo "🤖 AI Integration Hook (Optional)"
echo "The AI integration hook can catch AI-generated patterns early."
echo "To enable it, add this line to .git/hooks/pre-commit:"
echo "  source .githooks/ai-integration"
echo ""
echo "This will check for:"
echo "  • AI-generated code patterns (templates, buzzwords, emojis)"
echo "  • Quality enforcement patterns (skip-tests, force-push, etc.)"
echo "  • Incomplete code patterns (placeholders, TODOs)"
echo ""

echo "=== ✅ Enhanced Hooks Setup Complete ==="
echo ""
echo "📋 Installed Hooks:"
echo "  • Pre-commit: Comprehensive validation (author, formatting, tests, coverage)"
echo "  • Pre-push: Additional checks (static analysis, security, performance)"
echo "  • AI Integration: Optional AI-generated pattern detection"
echo ""
echo "🎯 Next Steps:"
echo "  1. Commit some code to test the hooks"
echo "  2. Try pushing to test pre-push validation"
echo "  3. (Optional) Enable AI integration for AI-assisted development"
echo ""
echo "📚 Documentation:"
echo "  • See .githooks/ for hook source code"
echo "  • See docs/development/setup.md for development guide"
