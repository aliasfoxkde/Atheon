#!/bin/bash
# Smart Exemption System for Documentation Validation
# Determines when to skip documentation checks based on change context

set -euo pipefail

# Source the categorization script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/categorize.sh"

# Exemption reasons
EXEMPT_TEST_ONLY="test-only"
EXEMPT_DOCS_UPDATED="docs-updated"
EXEMPT_INTERNAL_ONLY="internal-only"
EXEMPT_FORMATTING_ONLY="formatting-only"
EXEMPT_CHANGES_ONLY="changes-only"
EXEMPT_EXPLICIT="explicit"

# Check if only test files changed
is_test_only_changes() {
    local files=("$@")
    local only_tests=true

    for file in "${files[@]}"; do
        if [[ ! "$file" =~ _test\.go$ ]] && [[ ! "$file" =~ ^test/ ]] && [[ ! "$file" == test/* ]]; then
            only_tests=false
            break
        fi
    done

    if [ "$only_tests" = true ]; then
        echo "$EXEMPT_TEST_ONLY"
        return 0
    fi

    return 1
}

# Check if documentation files were updated
has_documentation_updates() {
    local files=("$@")
    local has_docs=false

    for file in "${files[@]}"; do
        if is_documentation_file "$file"; then
            has_docs=true
            break
        fi
    done

    if [ "$has_docs" = true ]; then
        echo "$EXEMPT_DOCS_UPDATED"
        return 0
    fi

    return 1
}

# Check if only internal files changed
is_internal_only_changes() {
    local files=("$@")
    local only_internal=true

    for file in "${files[@]}"; do
        local category=$(categorize_file "$file")
        if [ "$category" != "internal" ]; then
            only_internal=false
            break
        fi
    done

    if [ "$only_internal" = true ]; then
        echo "$EXEMPT_INTERNAL_ONLY"
        return 0
    fi

    return 1
}

# Check if changes are only formatting related
is_formatting_only_changes() {
    local files=("$@")

    # Check git diff for formatting-only changes
    for file in "${files[@]}"; do
        if [ ! -f "$file" ]; then
            continue
        fi

        # Get the diff
        local diff=$(git diff --cached "$file" 2>/dev/null || echo "")

        # Check if diff contains only whitespace changes
        if [ -n "$diff" ]; then
            # Remove whitespace-only lines and check if anything remains
            local content_changes=$(echo "$diff" | grep -v "^[+-]\s*$" | grep -v "^@" || echo "")

            if [ -z "$content_changes" ]; then
                echo "$EXEMPT_FORMATTING_ONLY"
                return 0
            fi
        fi
    done

    return 1
}

# Check if changelog was updated
has_changelog_updates() {
    local files=("$@")

    for file in "${files[@]}"; do
        if [[ "$file" =~ CHANGELOG ]] || [[ "$file" =~ changelog ]] || [[ "$file" =~ CHANGES ]]; then
            echo "$EXEMPT_CHANGES_ONLY"
            return 0
        fi
    done

    return 1
}

# Check for explicit exemption in commit message
has_exemption_intent() {
    local commit_msg=$(git log -1 --pretty=%B 2>/dev/null || echo "")

    if echo "$commit_msg" | grep -qiE "\[skip-docs\]|\[no-docs\]|skip documentation|no doc update"; then
        echo "$EXEMPT_EXPLICIT:intent"
        return 0
    fi

    if echo "$commit_msg" | grep -qiE "internal|refactor|perf|optimization"; then
        # Check if the message suggests internal changes
        if echo "$commit_msg" | grep -qiE "no user impact|internal only|implementation detail"; then
            echo "$EXEMPT_EXPLICIT:intent-internal"
            return 0
        fi
    fi

    return 1
}

# Determine overall exemption status
get_exemption_status() {
    local files=("$@")
    local exemptions=()

    # Check for various exemption conditions
    for file in "${files[@]}"; do
        if [ ! -f "$file" ] && [ ! -d "$file" ]; then
            continue
        fi
    done

    # Test-only changes
    if is_test_only_changes "${files[@]}"; then
        exemptions+=("$(is_test_only_changes)")
    fi

    # Documentation updates
    if has_documentation_updates "${files[@]}"; then
        exemptions+=("$(has_documentation_updates)")
    fi

    # Changelog updates
    if has_changelog_updates "${files[@]}"; then
        exemptions+=("$(has_changelog_updates)")
    fi

    # Internal-only changes
    if is_internal_only_changes "${files[@]}"; then
        exemptions+=("$(is_internal_only_changes)")
    fi

    # Explicit exemption
    if has_exemption_intent; then
        exemptions+=("$(has_exemption_intent)")
    fi

    # Return exemptions
    if [ ${#exemptions[@]} -gt 0 ]; then
        printf '%s\n' "${exemptions[@]}"
        return 0
    fi

    return 1
}

# Check if changes should be exempted
should_exempt_documentation() {
    local files=("$@")

    if exemptions=$(get_exemption_status "${files[@]}"); then
        echo "$exemptions"
        return 0
    fi

    return 1
}

# Format exemption message
format_exemption_message() {
    local exemption=$1

    case "$exemption" in
        "$EXEMPT_TEST_ONLY")
            echo "Only test files changed"
            ;;
        "$EXEMPT_DOCS_UPDATED")
            echo "Documentation files updated"
            ;;
        "$EXEMPT_INTERNAL_ONLY")
            echo "Only internal implementation changes"
            ;;
        "$EXEMPT_FORMATTING_ONLY")
            echo "Only formatting changes"
            ;;
        "$EXEMPT_CHANGES_ONLY")
            echo "Changelog updated"
            ;;
        "$EXEMPT_EXPLICIT"*)
            echo "Explicit exemption requested"
            ;;
        *)
            echo "Unknown exemption: $exemption"
            ;;
    esac
}

# Main execution
main() {
    local files=("$@")

    if [ ${#files[@]} -eq 0 ]; then
        # If no files specified, check git staged changes
        mapfile -t files < <(git diff --cached --name-only)
    fi

    if [ ${#files[@]} -eq 0 ]; then
        echo "No files to check"
        exit 0
    fi

    # Check for exemptions
    if should_exempt_documentation "${files[@]}"; then
        echo "true"
        exit 0
    else
        echo "false"
        exit 1
    fi
}

# Run main if executed directly
if [ "${BASH_SOURCE[0]}" == "${0}" ]; then
    main "$@"
fi