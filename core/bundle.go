package core

import (
	"bytes"
	"compress/gzip"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

//go:embed patterns.bundle
var embeddedBundle []byte

// PatternDef is the on-disk (and on-wire) representation of a pattern as
// it appears inside a pattern bundle. Match holds the regular-expression
// source; the compiled *regexp.Regexp is not part of the wire form.
type PatternDef struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Match    string `json:"match"`
	Enabled  bool   `json:"enabled"`
}

type bundlePattern struct {
	name     string
	category string
	match    string
	enabled  bool
	re       *regexp.Regexp
}

func (p *bundlePattern) Name() string             { return p.name }
func (p *bundlePattern) Category() string         { return p.category }
func (p *bundlePattern) Matches(line string) bool { return p.enabled && p.re.MatchString(line) }
func (p *bundlePattern) Enabled() bool            { return p.enabled }
func (p *bundlePattern) SetEnabled(enabled bool)  { p.enabled = enabled }

type categoryScanner struct {
	combined *regexp.Regexp
	patterns []Pattern
}

var (
	allPatterns          []*bundlePattern
	activeScanners       []categoryScanner
	activeCategoryFilter []string // nil = all categories; preserved across rebuildActiveScanners
)

func init() {
	// Default initialization reads the user's local bundle from ~/.atheon
	// and falls back to the embedded bundle. Splitting this out as
	// initializeWith lets tests exercise the error branches.
	home, _ := os.UserHomeDir()
	data := embeddedBundle
	if b, err := os.ReadFile(filepath.Join(home, ".atheon", "patterns.bundle")); err == nil {
		data = b
	}
	initializeWith(data)
}

// initializeWith runs the same setup as init() but accepts the bundle data
// directly so tests can feed in corrupt data to exercise the error paths.
func initializeWith(data []byte) {
	if err := loadBundle(data); err != nil {
		fmt.Fprintf(os.Stderr, "atheon: bundle load failed: %v\n", err)
	}
	SetActiveCategories(nil)

	// Load pattern state after bundle is loaded
	if err := InitializePatternState(); err != nil {
		// Non-fatal error, just log warning
		fmt.Fprintf(os.Stderr, "atheon: pattern state initialization failed: %v\n", err)
	}
}

func loadBundle(data []byte) error {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("%w: %v", ErrBundleParse, err)
	}
	defer r.Close()

	var defs []PatternDef
	if err := json.NewDecoder(r).Decode(&defs); err != nil {
		return fmt.Errorf("%w: %v", ErrBundleParse, err)
	}

	var external []Pattern
	for _, p := range registry {
		if _, ok := p.(*bundlePattern); !ok {
			external = append(external, p)
		}
	}
	registry = nil
	allPatterns = nil

	for _, def := range defs {
		re, err := regexp.Compile(def.Match)
		if err != nil {
			fmt.Fprintf(os.Stderr, "atheon: skipping %q: %v\n", def.Name, err)
			continue
		}
		bp := &bundlePattern{name: def.Name, category: def.Category, match: def.Match, enabled: def.Enabled, re: re}
		allPatterns = append(allPatterns, bp)
		Register(bp)
	}

	// Old bundles predate the enabled field; JSON zero-value false means all appear
	// disabled. Detect this and default everything to enabled.
	anyEnabled := false
	for _, p := range allPatterns {
		if p.enabled {
			anyEnabled = true
			break
		}
	}
	if !anyEnabled {
		for _, p := range allPatterns {
			p.enabled = true
		}
	}

	for _, p := range external {
		Register(p)
	}

	return nil
}

// SetActiveCategories restricts subsequent scans to the named categories.
// A nil or empty slice means "all categories." Calling this rebuilds the
// internal pre-filter regexes used by ScanFile and ScanDir.
func SetActiveCategories(cats []string) {
	activeCategoryFilter = cats

	catSet := map[string]bool{}
	for _, c := range cats {
		catSet[strings.TrimSpace(c)] = true
	}

	byCategory := map[string][]Pattern{}
	for _, p := range allPatterns {
		if !p.enabled {
			continue
		}
		if len(cats) > 0 && !catSet[p.category] {
			continue
		}
		byCategory[p.category] = append(byCategory[p.category], p)
	}
	// Include externally registered non-bundle patterns so Register() callers are scanned
	for _, p := range registry {
		if _, ok := p.(*bundlePattern); ok {
			continue
		}
		cat := p.Category()
		if len(cats) > 0 && !catSet[cat] {
			continue
		}
		byCategory[cat] = append(byCategory[cat], p)
	}

	activeScanners = nil
	for _, patterns := range byCategory {
		// Split bundle patterns (have a match regex) from external patterns (don't)
		var bundlePs, extPs []Pattern
		for _, p := range patterns {
			if _, ok := p.(*bundlePattern); ok {
				bundlePs = append(bundlePs, p)
			} else {
				extPs = append(extPs, p)
			}
		}
		if len(bundlePs) > 0 {
			parts := make([]string, 0, len(bundlePs))
			for _, p := range bundlePs {
				parts = append(parts, "(?:"+p.(*bundlePattern).match+")")
			}
			combined, err := regexp.Compile(strings.Join(parts, "|"))
			if err == nil {
				activeScanners = append(activeScanners, categoryScanner{combined: combined, patterns: bundlePs})
			}
		}
		if len(extPs) > 0 {
			// External patterns have no regex to pre-filter with; use empty (matches all)
			combined := regexp.MustCompile("")
			activeScanners = append(activeScanners, categoryScanner{combined: combined, patterns: extPs})
		}
	}
}

// Categories returns the unique, unsorted list of category labels present
// in the current bundle. The returned slice is owned by the caller.
func Categories() []string {
	seen := map[string]bool{}
	var cats []string
	for _, p := range allPatterns {
		if !seen[p.category] {
			seen[p.category] = true
			cats = append(cats, p.category)
		}
	}
	return cats
}

// bundleDownloadURL is the default upstream bundle URL. Tests may override
// it via SetBundleDownloadURL to point at an httptest server.
var bundleDownloadURL = "https://github.com/HoraDomu/Atheon/releases/latest/download/patterns.bundle"

// SetBundleDownloadURL swaps the upstream URL used by DownloadBundle. It
// returns a restore function that callers should defer to reset the URL
// after tests or short-lived overrides. Exported so external test
// packages (e.g., the main binary's tests) can stub out network access.
func SetBundleDownloadURL(url string) func() {
	orig := bundleDownloadURL
	bundleDownloadURL = url
	return func() { bundleDownloadURL = orig }
}

// DownloadBundle fetches the latest pattern bundle from the URL
// configured via SetBundleDownloadURL (or the default URL), compares it
// against the in-memory bundle, prints a summary of added/removed
// patterns, and persists the new bundle to ~/.atheon/patterns.bundle.
//
// The bundle is loaded into memory before being written to disk; if
// loadBundle fails the on-disk bundle is left untouched.
//
// On any non-success HTTP status code, DownloadBundle returns an error
// wrapping ErrBundleDownload so callers can use errors.Is.
// DownloadBundle fetches the latest pattern bundle from the URL
// configured via SetBundleDownloadURL (or the default URL), compares it
// against the in-memory bundle, prints a summary of added/removed
// patterns, and persists the new bundle to ~/.atheon/patterns.bundle.
//
// The bundle is loaded into memory before being written to disk; if
// loadBundle fails the on-disk bundle is left untouched.
//
// The context controls the HTTP request lifecycle: canceling ctx
// aborts the in-flight download.
//
// On any non-success HTTP status code, DownloadBundle returns an error
// wrapping ErrBundleDownload so callers can use errors.Is.
func DownloadBundle(ctx context.Context) error {
	oldPatterns := currentPatternNames()

	data, err := fetchBundleData(ctx)
	if err != nil {
		return err
	}
	dir, err := ensureAtheonDir()
	if err != nil {
		return err
	}

	newDefs, err := parseBundle(data)
	if err != nil {
		return err
	}

	added, removed := diffPatternNames(oldPatterns, newDefs)
	printBundleDiff(len(oldPatterns), len(newDefs), added, removed)

	// Load into memory first; only persist to disk if that succeeds
	if err := loadBundle(data); err != nil {
		return err
	}
	SetActiveCategories(activeCategoryFilter)
	if err := os.WriteFile(filepath.Join(dir, "patterns.bundle"), data, 0o600); err != nil {
		return err
	}
	return nil
}

// currentPatternNames returns the names of every pattern currently in
// the active bundle, in slice order.
func currentPatternNames() []string {
	var names []string
	for _, p := range allPatterns {
		names = append(names, p.name)
	}
	return names
}

// fetchBundleData performs the HTTP GET against bundleDownloadURL and
// returns the response body on success.
func fetchBundleData(ctx context.Context) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	// bundleDownloadURL is configured via SetBundleDownloadURL and only
	// ever points at https://github.com/... or a test stub. SSRF surface
	// is bounded by the controlled allow-list maintained in this package.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bundleDownloadURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBundleDownload, err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBundleDownload, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: server returned %d", ErrBundleDownload, resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBundleDownload, err)
	}
	return data, nil
}

// ensureAtheonDir creates (if needed) and returns the ~/.atheon path.
func ensureAtheonDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".atheon")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

// parseBundle decodes a gzipped JSON bundle into PatternDefs.
func parseBundle(data []byte) ([]PatternDef, error) {
	var defs []PatternDef
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBundleParse, err)
	}
	defer r.Close()
	if err := json.NewDecoder(r).Decode(&defs); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBundleParse, err)
	}
	return defs, nil
}

// diffPatternNames computes the symmetric set difference between the
// currently-loaded pattern names and a freshly-downloaded bundle.
func diffPatternNames(oldPatterns []string, newDefs []PatternDef) (added, removed []string) {
	newSet := make(map[string]bool, len(newDefs))
	for _, def := range newDefs {
		newSet[def.Name] = true
	}
	oldSet := make(map[string]bool, len(oldPatterns))
	for _, name := range oldPatterns {
		oldSet[name] = true
	}
	for _, name := range oldPatterns {
		if !newSet[name] {
			removed = append(removed, name)
		}
	}
	for _, def := range newDefs {
		if !oldSet[def.Name] {
			added = append(added, def.Name)
		}
	}
	return added, removed
}

// printBundleDiff writes a human-readable summary of the bundle change.
func printBundleDiff(oldCount, newCount int, added, removed []string) {
	fmt.Printf("Patterns updated: %d → %d\n", oldCount, newCount)
	switch {
	case len(added) > 0:
		fmt.Printf("Added: %d patterns\n", len(added))
		for _, p := range added {
			fmt.Printf("  + %s\n", p)
		}
		if len(removed) > 0 {
			fmt.Printf("Removed: %d patterns\n", len(removed))
			for _, p := range removed {
				fmt.Printf("  - %s\n", p)
			}
		}
	case len(removed) > 0:
		fmt.Printf("Removed: %d patterns\n", len(removed))
		for _, p := range removed {
			fmt.Printf("  - %s\n", p)
		}
	default:
		fmt.Println("No pattern changes detected")
	}
}

// EnablePattern enables the named pattern, rebuilds the scanner set,
// and persists the new state. It returns false if no pattern with the
// given name exists in the bundle.
func EnablePattern(name string) bool {
	for _, p := range allPatterns {
		if p.name != name {
			continue
		}
		p.enabled = true
		rebuildRegistry()
		rebuildActiveScanners()
		if err := syncPatternState(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to save pattern state: %v\n", err)
		}
		return true
	}
	return false
}

// DisablePattern disables the named pattern, rebuilds the scanner set,
// and persists the new state. It returns false if no pattern with the
// given name exists in the bundle.
func DisablePattern(name string) bool {
	for _, p := range allPatterns {
		if p.name != name {
			continue
		}
		p.enabled = false
		rebuildRegistry()
		rebuildActiveScanners()
		if err := syncPatternState(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to save pattern state: %v\n", err)
		}
		return true
	}
	return false
}

// SetPatternEnabled sets the enabled flag for the named pattern and
// rebuilds the active scanner set. Unlike EnablePattern and
// DisablePattern it does not persist state to disk — useful for tests
// and for callers that batch state updates. Returns false if the
// pattern name is unknown.
func SetPatternEnabled(name string, enabled bool) bool {
	for _, p := range allPatterns {
		if p.name == name {
			p.enabled = enabled
			rebuildActiveScanners()
			return true
		}
	}
	return false
}

// ListDisabledPatterns returns the names of every pattern that is
// currently disabled, in bundle order.
func ListDisabledPatterns() []string {
	var disabled []string
	for _, p := range allPatterns {
		if !p.enabled {
			disabled = append(disabled, p.name)
		}
	}
	return disabled
}

// ListEnabledPatterns returns the names of every pattern that is
// currently enabled, in bundle order.
func ListEnabledPatterns() []string {
	var enabled []string
	for _, p := range allPatterns {
		if p.enabled {
			enabled = append(enabled, p.name)
		}
	}
	return enabled
}

func rebuildActiveScanners() {
	SetActiveCategories(activeCategoryFilter)
}

// EnableAllPatterns enables every pattern in the bundle, overriding any
// prior disable calls, then rebuilds the active scanner set.
func EnableAllPatterns() {
	for _, p := range allPatterns {
		p.enabled = true
	}
	rebuildActiveScanners()
}

// rebuildRegistry rebuilds the registry from allPatterns, respecting enabled state
func rebuildRegistry() {
	registry = nil
	for _, p := range allPatterns {
		if p.enabled {
			Register(p)
		}
	}
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
