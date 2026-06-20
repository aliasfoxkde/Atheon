package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aliasfoxkde/Atheon/core"
)

// version is injected at build time via ldflags
var version = "dev"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := run(ctx, os.Args[1:])
	cancel() // explicit: os.Exit skips deferred cancel
	os.Exit(code)
}

// run executes the CLI with the given args and returns the exit code.
// This is separated from main() so tests can call it without os.Exit
// terminating the test process.
//
// The context flows through every Scan*/DownloadBundle call so callers
// (typically signal.NotifyContext from main) can cancel in-flight work.
func run(ctx context.Context, args []string) int {
	// Handle --version flag
	if len(args) > 0 && args[0] == "--version" {
		fmt.Printf("atheon %s\n", version)
		return 0
	}

	jsonOutput, args := extractJSONFlag(args)
	reportFormat, args := extractFormatFlag(args)
	if jsonOutput {
		reportFormat = core.FormatJSON
	}

	cats, args, enableAll := parseCategories(args)
	if enableAll {
		core.EnableAllPatterns()
	}
	core.SetActiveCategories(cats)

	if len(args) == 0 {
		printHelp()
		return 0
	}

	switch args[0] {
	case "update":
		fmt.Println("downloading patterns bundle...")
		if err := core.DownloadBundle(ctx); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return 1
		}
		return 0

	case "enable":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "error: enable requires a pattern name")
			return 1
		}
		if !core.EnablePattern(args[1]) {
			fmt.Fprintf(os.Stderr, "error: pattern '%s' not found\n", args[1])
			return 1
		}
		fmt.Printf("enabled pattern: %s\n", args[1])
		return 0

	case "disable":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "error: disable requires a pattern name")
			return 1
		}
		if !core.DisablePattern(args[1]) {
			fmt.Fprintf(os.Stderr, "error: pattern '%s' not found\n", args[1])
			return 1
		}
		fmt.Printf("disabled pattern: %s\n", args[1])
		return 0

	case "list":
		cmdList(args[1:])
		return 0

	case "--help", "help", "-h":
		printHelp()
		return 0

	case "--env":
		findings := core.ScanEnv(ctx)
		printFindings(findings, nil, reportFormat)
		if len(findings) > 0 {
			return 1
		}
		return 0

	case "-", "--stdin":
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: reading stdin:", err)
			return 1
		}
		findings := core.ScanString(ctx, string(data), "stdin")
		printFindings(findings, nil, reportFormat)
		if len(findings) > 0 {
			return 1
		}
		return 0

	case "--file":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "error: --file requires a path")
			return 1
		}
		findings, stats, err := core.ScanFile(ctx, args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return 1
		}
		printFindings(findings, stats, reportFormat)
		if len(findings) > 0 {
			return 1
		}
		return 0

	case "scan-url":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "error: scan-url requires a URL")
			return 1
		}
		findings, stats, err := core.ScanURL(ctx, args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return 1
		}
		printFindings(findings, stats, reportFormat)
		if len(findings) > 0 {
			return 1
		}
		return 0

	case "scan-git":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "error: scan-git requires a git remote URL")
			return 1
		}
		findings, stats, err := core.ScanGitRemote(ctx, args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return 1
		}
		printFindings(findings, stats, reportFormat)
		if len(findings) > 0 {
			return 1
		}
		return 0

	case "audit":
		root := "."
		if len(args) >= 2 {
			root = args[1]
		}
		report, err := core.Audit(ctx, root)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return 1
		}
		// Write to docs/audits/ with a timestamped subdirectory.
		ts := report.GeneratedAt.Format("2006-01-02-150405")
		dir := fmt.Sprintf("docs/audits/%s", ts)
		if err := core.WriteReport(report, dir); err != nil {
			fmt.Fprintln(os.Stderr, "error writing report:", err)
			return 1
		}
		fmt.Printf("audit complete: %s/REPORT.md\n", dir)
		return 0

	default:
		path := args[0]
		info, err := os.Stat(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: path not found:", path)
			return 1
		}
		var findings []core.Finding
		var stats *core.Stats
		if info.IsDir() {
			findings, stats, err = core.ScanDir(ctx, path)
		} else {
			findings, stats, err = core.ScanFile(ctx, path)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return 1
		}
		if stats != nil {
			for _, werr := range stats.WalkErrors {
				fmt.Fprintf(os.Stderr, "warning: skipped file: %v\n", werr)
			}
		}
		printFindings(findings, stats, reportFormat)
		if len(findings) > 0 {
			return 1
		}
		return 0
	}
}

// extractJSONFlag scans args for a `--json` token in any position and
// returns (true, rest) if found, (false, args) otherwise. The returned
// slice is freshly allocated; the input is not mutated.
//
// Recognising the flag in any position matters for users who reach for
// shell aliases (`alias atheon='atheon --json'`) and for scripts that
// build the argv vector incrementally, both of which would otherwise
// see `--json` treated as a path and fail.
func extractJSONFlag(args []string) (bool, []string) {
	for i, a := range args {
		if a == "--json" {
			rest := make([]string, 0, len(args)-1)
			rest = append(rest, args[:i]...)
			rest = append(rest, args[i+1:]...)
			return true, rest
		}
	}
	return false, args
}

// extractFormatFlag scans args for a `--format=<value>` token and returns
// the corresponding core.Format. Supported values: text, json, sarif, html.
// If absent or unknown, returns core.FormatText. The flag is stripped from
// the returned args slice so downstream parsers don't see it.
func extractFormatFlag(args []string) (core.Format, []string) {
	for i, a := range args {
		if strings.HasPrefix(a, "--format=") {
			val := strings.TrimPrefix(a, "--format=")
			format := core.Format(val)
			switch format {
			case core.FormatJSON, core.FormatSARIF, core.FormatHTML:
				rest := make([]string, 0, len(args)-1)
				rest = append(rest, args[:i]...)
				rest = append(rest, args[i+1:]...)
				return format, rest
			}
		}
	}
	return core.FormatText, args
}

func parseCategories(args []string) (cats, rest []string, enableAll bool) {
	for _, a := range args {
		switch {
		case strings.HasPrefix(a, "--categories="):
			val := strings.TrimPrefix(a, "--categories=")
			for _, c := range strings.Split(val, ",") {
				if c = strings.TrimSpace(c); c != "" {
					cats = append(cats, c)
				}
			}
		case a == "--all":
			enableAll = true
		default:
			rest = append(rest, a)
		}
	}
	return
}

func printFindings(findings []core.Finding, stats *core.Stats, reportFormat core.Format) {
	var s core.Stats
	if stats != nil {
		s = *stats
	}
	rep := core.Report{
		Version:     version,
		GeneratedAt: time.Now(),
		Findings:    findings,
		Stats:       s,
	}
	fmt.Print(core.Render(rep, reportFormat))
}

func cmdList(args []string) {
	if len(args) > 0 && args[0] == "categories" {
		for _, c := range core.Categories() {
			fmt.Println(c)
		}
		return
	}

	var categoryFilter string
	showEnabled := false
	showDisabled := false
	for _, a := range args {
		switch {
		case strings.HasPrefix(a, "--category="):
			categoryFilter = strings.TrimPrefix(a, "--category=")
		case a == "--enabled":
			showEnabled = true
		case a == "--disabled":
			showDisabled = true
		}
	}

	var filtered []core.Pattern
	for _, p := range core.All() {
		if categoryFilter != "" && p.Category() != categoryFilter {
			continue
		}
		if showEnabled && !p.Enabled() {
			continue
		}
		if showDisabled && p.Enabled() {
			continue
		}
		filtered = append(filtered, p)
	}

	for _, p := range filtered {
		status := "enabled"
		if !p.Enabled() {
			status = "disabled"
		}
		fmt.Printf("%s [%s] [%s]\n", p.Name(), p.Category(), status)
	}
	fmt.Printf("\n%d pattern(s)\n", len(filtered))
}

func printHelp() {
	fmt.Print(`atheon - pattern matching engine

usage:
  atheon <path>                      scan a directory
  atheon --file <path>               scan a single file
  atheon --env                       scan environment variables
  atheon --json <path>               print findings as JSON (same as --format=json)
  atheon --format=<fmt> <path>       output format: text (default), json, sarif, html
  atheon --categories=<c1,c2> <path> scan specific categories
  atheon --all <path>                scan all patterns including disabled ones
  atheon list                        list all patterns with enabled/disabled status
  atheon list --enabled              list only enabled patterns
  atheon list --disabled             list only disabled patterns
  atheon scan-url <url>             scan a remote URL for secrets
  atheon scan-git <url>             scan a remote git repository for secrets
  atheon audit [path]              run audit checks and write REPORT.md + REPORT.json
  atheon list categories             list available categories
  atheon enable <pattern>            enable a pattern
  atheon disable <pattern>           disable a pattern
  atheon update                      download latest patterns bundle
  atheon --help                      show this message
`)
}

func redact(s string) string {
	if len(s) <= 8 {
		return "***"
	}
	return s[:4] + "****" + s[len(s)-4:]
}

func formatBytes(b int64) string {
	if b >= 1<<20 {
		return fmt.Sprintf("%.1f MB", float64(b)/(1<<20))
	}
	if b >= 1<<10 {
		return fmt.Sprintf("%.1f KB", float64(b)/(1<<10))
	}
	return fmt.Sprintf("%d B", b)
}
