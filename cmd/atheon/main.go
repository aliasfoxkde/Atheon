package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

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

	jsonOutput := len(args) > 0 && args[0] == "--json"
	if jsonOutput {
		args = args[1:]
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
		fmt.Println("patterns updated.")
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
		printFindings(findings, nil, jsonOutput)
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
		printFindings(findings, nil, jsonOutput)
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
		printFindings(findings, stats, jsonOutput)
		if len(findings) > 0 {
			return 1
		}
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
		printFindings(findings, stats, jsonOutput)
		if len(findings) > 0 {
			return 1
		}
		return 0
	}
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

func printFindings(findings []core.Finding, stats *core.Stats, jsonOutput bool) {
	if jsonOutput {
		printJSONFindings(findings)
		return
	}
	if len(findings) == 0 {
		fmt.Println("no findings.")
	} else {
		for _, f := range findings {
			loc := f.File
			if f.Line > 0 {
				loc = fmt.Sprintf("%s:%d", f.File, f.Line)
			}
			fmt.Printf("%s  %s\n", f.Pattern, loc)
			if f.Content != "" {
				fmt.Println(" ", redact(f.Content))
			}
		}
		fmt.Printf("\n%d finding(s)\n", len(findings))
	}
	if stats != nil && stats.Files > 0 {
		fmt.Printf("scanned %d file(s)  %s  %dms\n",
			stats.Files, formatBytes(stats.Bytes), stats.ElapsedMs)
	}
}

func printJSONFindings(findings []core.Finding) {
	items := make([]map[string]any, 0, len(findings))
	for _, f := range findings {
		items = append(items, map[string]any{"pattern": f.Pattern, "file": f.File, "line": f.Line, "match": redact(f.Content)})
	}
	if err := json.NewEncoder(os.Stdout).Encode(items); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	}
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
  atheon --json <path>               print findings as JSON
  atheon --categories=<c1,c2> <path> scan specific categories
  atheon --all <path>                scan all patterns including disabled ones
  atheon list                        list all patterns with enabled/disabled status
  atheon list --enabled              list only enabled patterns
  atheon list --disabled             list only disabled patterns
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
