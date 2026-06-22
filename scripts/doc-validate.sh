#!/bin/bash
# Documentation Validator - Checks if documentation is current for code changes
# Uses intelligent file mapping and age-based validation

set -euo pipefail

# Source the categorization script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/doc-categorize.sh"

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Documentation mapping
declare -A DOC_MAPPING
DOC_MAPPING["main.go"]="README.md docs/user-guide.md docs/cli-reference.md"
DOC_MAPPING["cmd/"]="README.md docs/user-guide.md"
DOC_MAPPING["community/*.yaml"]="docs/patterns/*.md README.md"
DOC_MAPPING["core/pattern"]="docs/pattern-development.md"
DOC_MAPPING["config/"]="docs/configuration.md INSTALL.md"

# Grace periods in seconds
GRACE_PERIOD_IMMEDIATE=3600    # 1 hour
GRACE_PERIOD_SHORT=86400      # 24 hours
GRACE_PERIOD_LONG=604800       # 7 days

# Get last modification time from git
get_file_time() {
    local file=$1
    if [ -f "$file" ]; then
        git log -1 --format=%ct -- "$file" 2>/dev/null || echo "0"
    else
        echo "0"
    fi
}

# Calculate age difference in seconds
calculate_age_diff() {
    local code_time=$1
    local doc_time=$2
    echo $((doc_time - code_time))
}

# Format age for display
format_age() {
    local seconds=$1
    if [ "$seconds" -lt 0 ]; then
        local abs_seconds=${seconds#-}
        if [ "$abs_seconds" -lt 3600 ]; then
            echo "$((abs_seconds / 60)) minutes ago"
        elif [ "$abs_seconds" -lt 86400 ]; then
            echo "$((abs_seconds / 3600)) hours ago"
        else
            echo "$((abs_seconds / 86400)) days ago"
        fi
    else
        echo "in the future"
    fi
}

# Check if documentation needs update based on age
check_documentation_age() {
    local code_file=$1
    local doc_file=$2

    local code_time=$(get_file_time "$code_file")
    local doc_time=$(get_file_time "$doc_file")

    if [ "$code_time" -eq 0 ] || [ "$doc_time" -eq 0 ]; then
        echo "skip"  # Can't determine age
        return
    fi

    local age_diff=$(calculate_age_diff "$code_time" "$doc_time")
    echo "$age_diff"
}

# Get relevant documentation for a code file
get_relevant_docs() {
    local code_file=$1
    local relevant_docs=()

    # Direct mappings
    for key in "${!DOC_MAPPING[@]}"; do
        if [[ "$code_file" == $key ]]; then
            IFS=' ' read -ra docs <<< "${DOC_MAPPING[$key]}"
            for doc in "${docs[@]}"; do
                if [ -f "$doc" ]; then
                    relevant_docs+=("$doc")
                fi
            done
        fi
    done

    # Pattern-based mappings
    if [[ "$code_file" == community/*.yaml ]]; then
        relevant_docs+=("README.md")
        if [ -d "docs/patterns" ]; then
            relevant_docs+=("docs/patterns/"*)
        fi
    fi

    if [[ "$code_file" == core/*.go ]] && grep -q "^func [A-Z]" "$code_file"; then
        relevant_docs+=("docs/api.md")
    fi

    echo "${relevant_docs[@]}"
}

# Validate documentation for a single file
validate_file_documentation() {
    local code_file=$1
    local category=$(categorize_file "$code_file")

    # Skip documentation and internal files
    if [ "$category" == "documentation" ] || [ "$category" == "internal" ]; then
        echo "skip"
        return
    fi

    local relevant_docs=($(get_relevant_docs "$code_file"))

    if [ ${#relevant_docs[@]} -eq 0 ]; then
        echo "no-docs"
        return
    fi

    local status="current"
    local stale_docs=()

    for doc_file in "${relevant_docs[@]}"; do
        local age=$(check_documentation_age "$code_file" "$doc_file")

        case "$age" in
            skip)
                continue
                ;;
            *)
                if [ "$age" -lt "-$GRACE_PERIOD_LONG" ]; then
                    status="stale"
                    stale_docs+=("$doc_file")
                elif [ "$age" -lt "-$GRACE_PERIOD_SHORT" ]; then
                    status="warning"
                    stale_docs+=("$doc_file")
                fi
                ;;
        esac
    done

    echo "$status:${stale_docs[*]}"
}

# Main validation function
validate_documentation() {
    local files=("$@")
    local results=()
    local stale_count=0
    local warning_count=0

    echo -e "${BLUE}=== 📚 Documentation Validation ===${NC}"
    echo ""

    for file in "${files[@]}"; do
        if [ ! -f "$file" ]; then
            continue
        fi

        local result=$(validate_file_documentation "$file")
        local status="${result%%:*}"
        local docs="${result#*:}"

        case "$status" in
            stale)
                echo -e "${RED}❌ $file${NC}"
                echo -e "   ${RED}Stale documentation: $docs${NC}"
                ((stale_count++))
                ;;
            warning)
                echo -e "${YELLOW}⚠️  $file${NC}"
                echo -e "   ${YELLOW}Documentation may need update: $docs${NC}"
                ((warning_count++))
                ;;
            current)
                echo -e "${GREEN}✓ $file${NC} - Documentation current"
                ;;
            no-docs)
                echo -e "${BLUE}ℹ️  $file${NC} - No documentation required"
                ;;
            skip)
                # Skipped (documentation/internal files)
                ;;
        esac
    done

    echo ""
    echo "Summary: $stale_count stale, $warning_count warnings"

    if [ "$stale_count" -gt 0 ]; then
        return 1
    elif [ "$warning_count" -gt 0 ]; then
        return 0  # Warning only
    else
        return 0
    fi
}

# Check if changes should be exempt from documentation validation
should_validate_docs() {
    local files=("$@")

    # Check if only test files changed
    local only_tests=true
    for file in "${files[@]}"; do
        if [[ ! "$file" =~ _test\.go$ ]] && [[ ! "$file" =~ test/ ]]; then
            only_tests=false
            break
        fi
    done

    if [ "$only_tests" = true ]; then
        echo "test-only"
        return
    fi

    # Check if documentation files were updated
    local has_docs=false
    for file in "${files[@]}"; do
        if is_documentation_file "$file"; then
            has_docs=true
            break
        fi
    done

    if [ "$has_docs" = true ]; then
        echo "has-docs"
        return
    fi

    echo "validate"
}

# Main execution
main() {
    if [ $# -eq 0 ]; then
        # If no files specified, check git staged changes
        mapfile -t files < <(git diff --cached --name-only)
    else
        files=("$@")
    fi

    if [ ${#files[@]} -eq 0 ]; then
        echo "No files to validate"
        exit 0
    fi

    # Check if we should validate
    local validation=$(should_validate_docs "${files[@]}")

    case "$validation" in
        test-only)
            echo -e "${BLUE}=== 📚 Documentation Validation ===${NC}"
            echo -e "${GREEN}✓ Only test files changed, skipping documentation check${NC}"
            exit 0
            ;;
        has-docs)
            echo -e "${BLUE}=== 📚 Documentation Validation ===${NC}"
            echo -e "${GREEN}✓ Documentation files updated, validation passed${NC}"
            exit 0
            ;;
    esac

    # Perform validation
    if ! validate_documentation "${files[@]}"; then
        echo ""
        echo -e "${RED}=== ❌ Documentation Validation Failed ===${NC}"
        echo "Please update relevant documentation before committing"
        echo "Or use --force to skip this check (not recommended)"
        exit 1
    fi

    echo ""
    echo -e "${GREEN}=== ✅ Documentation Validation Passed ===${NC}"
}

# Run main if executed directly
if [ "${BASH_SOURCE[0]}" == "${0}" ]; then
    main "$@"
fi