#!/bin/bash
# File Categorization Engine for Documentation Validation
# Determines if a change requires documentation based on file type and change nature

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# File categories
USER_FACING="user-facing"
INTERNAL="internal"
DOCUMENTATION="documentation"
CONFIG="config"

# Categorize a single file
categorize_file() {
    local file=$1

    # Skip if file doesn't exist
    if [ ! -f "$file" ]; then
        echo "$INTERNAL"
        return
    fi

    # Documentation files
    if is_documentation_file "$file"; then
        echo "$DOCUMENTATION"
        return
    fi

    # Configuration files
    if is_config_file "$file"; then
        echo "$CONFIG"
        return
    fi

    # User-facing changes
    if is_user_facing_file "$file"; then
        echo "$USER_FACING"
        return
    fi

    # Default to internal
    echo "$INTERNAL"
}

# Check if file is documentation
is_documentation_file() {
    local file=$1
    case "$file" in
        *.md|*.txt|docs/*|README*|CHANGELOG*|LICENSE*)
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

# Check if file is configuration
is_config_file() {
    local file=$1
    case "$file" in
        *.yaml|*.yml|*.json|*.toml|*.conf|config/*|*.example)
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

# Check if file is user-facing
is_user_facing_file() {
    local file=$1
    case "$file" in
        main.go|cmd/*|cli*)
            return 0
            ;;
        community/*.yaml)
            return 0
            ;;
        core/*.go)
            # Only if it contains exported functions
            if grep -q "^func [A-Z]" "$file" 2>/dev/null; then
                return 0
            fi
            return 1
            ;;
        *)
            return 1
            ;;
    esac
}

# Check if file is internal
is_internal_file() {
    local file=$1
    case "$file" in
        *_test.go|test/*|internal/*|*.test)
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

# Categorize multiple files
categorize_files() {
    local files=("$@")
    local user_facing=()
    local internal=()
    local documentation=()
    local config=()

    for file in "${files[@]}"; do
        local category=$(categorize_file "$file")
        case "$category" in
            "$USER_FACING")
                user_facing+=("$file")
                ;;
            "$INTERNAL")
                internal+=("$file")
                ;;
            "$DOCUMENTATION")
                documentation+=("$file")
                ;;
            "$CONFIG")
                config+=("$file")
                ;;
        esac
    done

    # Return results as JSON
    echo "{\"user_facing\":[\"$(IFS='"","'; echo "${user_facing[*]}" | sed 's/ /" "/g')\"],\"internal\":[\"$(IFS='"","'; echo "${internal[*]}" | sed 's/ /" "/g')\"],\"documentation\":[\"$(IFS='"","'; echo "${documentation[*]}" | sed 's/ /" "/g')\"],\"config\":[\"$(IFS='"","'; echo "${config[*]}" | sed 's/ /" "/g')\"]}"
}

# Main execution
main() {
    if [ $# -eq 0 ]; then
        echo "Usage: $0 <file1> [file2] ..."
        exit 1
    fi

    categorize_files "$@"
}

# Run main if executed directly
if [ "${BASH_SOURCE[0]}" == "${0}" ]; then
    main "$@"
fi