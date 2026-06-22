#!/usr/bin/env bash
# Install Atheon development hooks and tools.
# Run once after cloning: ./scripts/install-hooks.sh

set -euo pipefail

echo "=== Atheon: Installing development hooks ==="

# Wire git to use the repo's hooks/ directory
git config core.hooksPath hooks
echo "  ✓ git core.hooksPath set to 'hooks'"

# Install pre-commit framework (optional but recommended)
if command -v pre-commit &>/dev/null; then
    pre-commit install
    echo "  ✓ pre-commit framework installed"
else
    echo "  ℹ pre-commit not found — install with: pip install pre-commit"
    echo "    Then re-run this script to activate YAML/whitespace checks"
fi

# Install optional Go tools
if command -v go &>/dev/null; then
    echo "  Installing optional Go tools..."
    go install honnef.co/go/tools/cmd/staticcheck@2024.1 2>/dev/null && echo "  ✓ staticcheck" || echo "  ⚠ staticcheck unavailable (non-blocking)"
    go install golang.org/x/tools/cmd/goimports@latest 2>/dev/null && echo "  ✓ goimports" || echo "  ⚠ goimports unavailable (non-blocking)"
else
    echo "  ⚠ Go not found — install from https://go.dev/dl/"
fi

echo ""
echo "=== Setup complete ==="
echo "Hooks will run automatically on 'git commit' and 'git push'."
echo "To skip a hook once: git commit --no-verify (use sparingly)"
