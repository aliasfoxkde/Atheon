package core

import (
	"errors"
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
