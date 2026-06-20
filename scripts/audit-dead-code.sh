#!/usr/bin/env bash
# scripts/audit-dead-code.sh
#
# Find unexported helpers that have no callers anywhere in the
# codebase. A helper is "dead" when its name appears only on
# its own declaration line -- i.e. it is defined but never
# invoked (including from tests).
#
# Usage:
#   scripts/audit-dead-code.sh           scan every .go file under core/, bundler/, cmd/
#   scripts/audit-dead-code.sh --staged  scan only files staged for the next commit
#
# Exit code: 0 = clean, 1 = at least one dead helper found.
#
# Both Makefile and .githooks/pre-commit invoke this script so
# the rule stays in one place.

set -euo pipefail

MODE="${1:-all}"

case "$MODE" in
    --staged)
        FILES=$(git diff --cached --name-only --diff-filter=AM 2>/dev/null \
            | grep '\.go$' || true)
        ;;
    all)
        FILES=$(find core bundler cmd -name '*.go' -not -name '*_test.go' 2>/dev/null || true)
        ;;
    *)
        echo "usage: $0 [--staged]" >&2
        exit 2
        ;;
esac

if [ -z "$FILES" ]; then
    echo "OK (no files to scan)"
    exit 0
fi

FAILED=0
for f in $FILES; do
    [ -f "$f" ] || continue
    # For --staged, use the staged content (post-edit) so deletions
    # and additions reflect what will land.
    if [ "$MODE" = "--staged" ]; then
        CONTENT=$(git show ":$f" 2>/dev/null || true)
        if [ -z "$CONTENT" ]; then
            # file deleted in this commit -- skip
            continue
        fi
        HELPERS=$(printf '%s\n' "$CONTENT" | grep -E '^func [a-z]' \
            | sed -E 's/^func +([A-Za-z0-9_]+).*/\1/' || true)
    else
        HELPERS=$(grep -E '^func [a-z]' "$f" \
            | sed -E 's/^func +([A-Za-z0-9_]+).*/\1/' || true)
    fi

    for h in $HELPERS; do
        # TOTAL references across the whole tree (incl. the def).
        # If TOTAL == 1, only the declaration line matches -- dead.
        # Guard the grep: no matches returns exit 1, which under
        # pipefail would abort the script.
        TOTAL=$(grep -rE "\b$h\b" --include='*.go' . 2>/dev/null | wc -l || echo 0)
        if [ "${TOTAL:-0}" -le 1 ]; then
            echo "DEAD: $f: $h (no callers)"
            FAILED=1
        fi
    done
done

if [ "$FAILED" -ne 0 ]; then
    exit 1
fi
echo "OK"