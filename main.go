package main

import (
	"atheon/core"
	_ "atheon/patterns"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "list":
		cmdList()

	case "--help", "help", "-h":
		printHelp()

	case "--env":
		findings := core.ScanEnv()
		printFindings(findings, nil)
		if len(findings) > 0 {
			os.Exit(1)
		}

	case "-", "--stdin":
		data, _ := io.ReadAll(os.Stdin)
		findings := core.ScanString(string(data), "stdin")
		printFindings(findings, nil)
		if len(findings) > 0 {
			os.Exit(1)
		}

	case "--file":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "error: --file requires a path")
			os.Exit(1)
		}
		findings, stats, err := core.ScanFile(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		printFindings(findings, stats)
		if len(findings) > 0 {
			os.Exit(1)
		}

	default:
		path := os.Args[1]
		info, err := os.Stat(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: path not found:", path)
			os.Exit(1)
		}
		var findings []core.Finding
		var stats *core.Stats
		if info.IsDir() {
			findings, stats, err = core.ScanDir(path)
		} else {
			findings, stats, err = core.ScanFile(path)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		printFindings(findings, stats)
		if len(findings) > 0 {
			os.Exit(1)
		}
	}
}

func printFindings(findings []core.Finding, stats *core.Stats) {
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

func cmdList() {
	for _, p := range core.All() {
		fmt.Println(p.Name())
	}
	fmt.Printf("\n%d pattern(s) loaded\n", len(core.All()))
}

func printHelp() {
	fmt.Print(`atheon - pattern matching engine

usage:
  atheon <path>          scan a directory
  atheon --file <path>   scan a single file
  atheon --env           scan environment variables
  atheon list            list loaded patterns
  atheon --help          show this message
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
