#!/usr/bin/env bash
# build.sh — canonical build for Atheon.
#
# Produces (only) the two binaries the project ships:
#   ./atheon        — main CLI
#   ./atheon-mcp    — MCP server
#   ./bundler/bundler — pattern bundler tool
#
# Always pass -o so the binary name is explicit (otherwise `go build .` in
# this repo names the binary "Atheon" — capitalized after the directory —
# which is wrong and is a build artifact, not part of the repo).
#
# Usage:
#   ./scripts/build.sh              # build all three
#   ./scripts/build.sh cli          # only ./atheon
#   ./scripts/build.sh mcp          # only ./atheon-mcp
#   ./scripts/build.sh bundler      # only ./bundler/bundler
#   ./scripts/build.sh clean        # remove any stray build artifacts

set -euo pipefail

build_cli() {
    echo "→ building atheon"
    go build -o atheon .
}

build_mcp() {
    echo "→ building atheon-mcp"
    go build -o atheon-mcp ./cmd/mcp
}

build_bundler() {
    echo "→ building bundler/bundler"
    go build -o bundler/bundler ./bundler
}

clean_strays() {
    # Remove any stray build artifact (e.g. Atheon produced by `go build .`
    # without -o, which names the binary after the directory). These are
    # build artifacts, not source.
    shopt -s nullglob nocaseglob
    for f in Atheon Atheon.exe atheon.exe atheon-mcp.exe; do
        if [[ -f "$f" ]]; then
            echo "→ removing stray build artifact: $f"
            rm -f -- "$f"
        fi
    done
}

case "${1:-all}" in
    cli)     build_cli ;;
    mcp)     build_mcp ;;
    bundler) build_bundler ;;
    clean)   clean_strays ;;
    all|"")  build_cli; build_mcp; build_bundler; clean_strays ;;
    *)
        echo "usage: $0 [all|cli|mcp|bundler|clean]" >&2
        exit 2
        ;;
esac
