#!/usr/bin/env bash
# scripts/pattern-count.sh — single source of truth for pattern counts.
#
# Usage:
#   ./scripts/pattern-count.sh           # human-readable summary
#   ./scripts/pattern-count.sh --json    # machine-readable JSON
#   ./scripts/pattern-count.sh --total   # total count only
#   ./scripts/pattern-count.sh --table   # markdown table for README inclusion
#
# This script reads community/**/*.yaml and produces counts. Use it to keep docs
# in sync — never hardcode numbers in README, FAQ, or INSTALL again.
#
# Output is deterministic; safe to diff in CI.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

# Collect all YAML files (excluding hidden dirs).
mapfile -t FILES < <(find community -mindepth 2 -type f \( -name '*.yaml' -o -name '*.yml' \) | sort)
TOTAL=${#FILES[@]}

# Per-category counts.
declare -A COUNTS
for f in "${FILES[@]}"; do
  cat="${f#community/}"
  cat="${cat%%/*}"
  COUNTS[$cat]=$((${COUNTS[$cat]:-0} + 1))
done

CATEGORIES=${#COUNTS[@]}
NON_EMPTY=0
for c in "${!COUNTS[@]}"; do
  if [[ "${COUNTS[$c]}" -gt 0 ]]; then
    NON_EMPTY=$((NON_EMPTY + 1))
  fi
done

case "${1:-}" in
  --json)
    {
      printf '{"total":%d,"categories":%d,"non_empty_categories":%d,"per_category":{' "$TOTAL" "$CATEGORIES" "$NON_EMPTY"
      first=1
      for c in $(printf '%s\n' "${!COUNTS[@]}" | sort); do
        if [[ $first -eq 0 ]]; then printf ','; fi
        printf '"%s":%d' "$c" "${COUNTS[$c]}"
        first=0
      done
      printf '}}\n'
    }
    ;;
  --total)
    printf '%d\n' "$TOTAL"
    ;;
  --table)
    printf '| Category | Patterns |\n|----------|----------|\n'
    for c in $(printf '%s\n' "${!COUNTS[@]}" | sort); do
      printf '| %s | %d |\n' "$c" "${COUNTS[$c]}"
    done
    printf '| **TOTAL** | **%d** |\n' "$TOTAL"
    ;;
  --help|-h)
    sed -n '2,17p' "$0"
    ;;
  *)
    printf 'Atheon pattern catalog\n'
    printf '  Total patterns:    %d\n' "$TOTAL"
    printf '  Categories:        %d (with content: %d)\n' "$CATEGORIES" "$NON_EMPTY"
    printf '\n'
    printf '%-25s %s\n' CATEGORY COUNT
    printf '%-25s %s\n' -------- -----
    for c in $(printf '%s\n' "${!COUNTS[@]}" | sort); do
      printf '%-25s %d\n' "$c" "${COUNTS[$c]}"
    done
    ;;
esac
