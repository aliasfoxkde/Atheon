package core

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
)

// Sentinel errors returned by the core API. Callers should use errors.Is
// to compare against these rather than string-matching error messages.
var (
	// ErrPatternNotFound is returned when a lookup by name fails to find
	// a pattern in the active bundle.
	ErrPatternNotFound = errors.New("pattern not found")

	// ErrBundleDownload is returned by DownloadBundle when the upstream
	// bundle cannot be fetched or the response body cannot be persisted.
	ErrBundleDownload = errors.New("bundle download failed")

	// ErrBundleParse is returned by loadBundle and DownloadBundle when
	// the bundle payload cannot be decompressed or decoded.
	ErrBundleParse = errors.New("bundle parse failed")

	// ErrInvalidPattern is returned when a pattern definition fails
	// validation (e.g. a malformed regex).
	ErrInvalidPattern = errors.New("invalid pattern")
)

// Pattern is the interface implemented by all scanner patterns, both those
// loaded from the embedded bundle and those registered programmatically
// via Register.
type Pattern interface {
	// Name returns the stable, human-readable pattern identifier.
	Name() string
	// Category returns the grouping label (e.g. "secrets", "pii").
	Category() string
	// Enabled reports whether the pattern is currently active.
	Enabled() bool
	// Severity returns the pattern's severity: one of "low", "medium", "high",
	// "critical". Implementations that don't track severity return "medium".
	Severity() string
	// Matches reports whether the pattern's regex matches the given line.
	Matches(line string) bool
}

// registry holds every registered Pattern. It is package-private; callers
// mutate it via Register and read it via All.
var registry []Pattern

// Register adds p to the active registry. It is safe to call before
// loadBundle (external patterns survive bundle loads) and after (they
// appear alongside bundle patterns in subsequent All() calls).
func Register(p Pattern) {
	registry = append(registry, p)
}

// All returns a sorted snapshot of the registered patterns ordered by
// Name. The returned slice is owned by the caller and may be mutated.
func All() []Pattern {
	sorted := make([]Pattern, len(registry))
	copy(sorted, registry)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name() < sorted[j].Name()
	})
	return sorted
}

// ValidatePattern checks that def has a valid regex and required fields.
// It returns nil if the pattern is valid, or an error describing the issue.
func ValidatePattern(def PatternDef) error {
	if def.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidPattern)
	}
	if def.Match == "" {
		return fmt.Errorf("%w: match regex is required", ErrInvalidPattern)
	}
	if _, err := regexp.Compile(def.Match); err != nil {
		return fmt.Errorf("%w: %q: %v", ErrInvalidPattern, def.Name, err)
	}
	return nil
}
