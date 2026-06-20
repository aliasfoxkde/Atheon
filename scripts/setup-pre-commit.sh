#!/bin/bash

# Pre-commit setup script for Atheon

echo "Setting up pre-commit hooks for Atheon..."

# Check if pre-commit is installed
if ! command -v pre-commit &> /dev/null; then
    echo "Installing pre-commit..."
    pip install pre-commit
fi

# Install pre-commit hooks
echo "Installing pre-commit hooks..."
pre-commit install

# Create .atheonignore if it doesn't exist
if [ ! -f .atheonignore ]; then
    echo "Creating .atheonignore..."
    cat > .atheonignore << 'EOF'
# Atheon ignore patterns
# Add patterns here to exclude files from scanning

# Example patterns:
test/
*.generated.go
.env
vendor/
node_modules/

# Temporary files
*.tmp
*.temp
*_test.go
coverage.out
EOF
fi

echo "Pre-commit hooks setup complete!"
echo ""
echo "Installed hooks:"
echo "  - Go tests"
echo "  - Go vet"
echo "  - Go fmt check"
echo "  - Go build verification"
echo "  - Go mod tidy"
echo "  - Trailing whitespace removal"
echo "  - YAML validation"
echo "  - Large file detection"
echo "  - Private key detection"
echo "  - GolangCI-Lint"
echo ""
echo "You can skip hooks with:"
echo "  git commit --no-verify -m \"message\""
echo ""
echo "Run hooks manually with:"
echo "  pre-commit run --all-files"