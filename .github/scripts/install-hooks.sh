#!/bin/bash
# Script to install enhanced git hooks for the repository

echo "=== Installing Enhanced Atheon Git Hooks ==="

# Copy hooks to .git/hooks directory
cp .githooks/pre-commit .git/hooks/pre-commit
cp .githooks/pre-push .git/hooks/pre-push

# Make hooks executable
chmod +x .git/hooks/pre-commit
chmod +x .git/hooks/pre-push

echo "✅ Enhanced hooks installed:"
echo "  • Pre-commit: Comprehensive validation before commits"
echo "  • Pre-push: Additional validation before pushing"
echo "  • AI Integration: Available at .githooks/ai-integration"
echo ""
echo "To enable AI integration, add to your pre-commit hook:"
echo "  source .githooks/ai-integration"
echo ""
echo "✅ Installation complete!"
