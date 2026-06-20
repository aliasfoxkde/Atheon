#!/usr/bin/env bash
# coverage.sh — run go tests with coverage, writing outputs to a temp dir.
#
# Usage:
#   ./scripts/coverage.sh                # default: ./tmp-cov/<timestamp>/
#   ./scripts/coverage.sh --html         # also write coverage HTML report
#   ./scripts/coverage.sh /custom/path   # custom output dir
#
# Output files (all under the chosen dir, never at repo root):
#   coverage.out      raw coverprofile
#   coverage.html     HTML report (only with --html)
#   summary.txt       `go tool cover -func` totals
#   test.log          full test output
#
# The temp dir is .gitignored via .tmp-cov/ and /tmp-cov/.

set -euo pipefail

# Pick output directory: arg > env > timestamped default under tmp-cov/
if [[ $# -ge 1 && "${1:-}" != "--html" ]]; then
    OUT_DIR="$1"
    shift
else
    OUT_DIR=".tmp-cov/$(date +%Y%m%d-%H%M%S)"
fi
WRITE_HTML=0
for arg in "$@"; do
    case "$arg" in
        --html) WRITE_HTML=1 ;;
    esac
done

mkdir -p "$OUT_DIR"

echo "→ writing coverage outputs to $OUT_DIR"
COV="$OUT_DIR/coverage.out"
LOG="$OUT_DIR/test.log"
SUM="$OUT_DIR/summary.txt"

# Run tests with coverage, tee log so the user sees progress and we capture output.
go test ./... -coverprofile="$COV" 2>&1 | tee "$LOG"

# Per-package summary
go tool cover -func="$COV" | tee "$SUM"
echo
echo "→ total:"
tail -1 "$SUM"

if [[ $WRITE_HTML -eq 1 ]]; then
    HTML="$OUT_DIR/coverage.html"
    go tool cover -html="$COV" -o "$HTML"
    echo "→ HTML: $HTML"
fi

echo "→ done. outputs in $OUT_DIR/"
