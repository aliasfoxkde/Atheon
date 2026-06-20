package core

// Finding represents a single pattern match produced by ScanFile, ScanDir,
// ScanString, or ScanEnv. File is the source path (or "env:KEY" when
// scanning the process environment); Line is the 1-indexed line number
// within the source (0 for env scans); Content is the trimmed matching
// line or, for env scans, the matching value.
type Finding struct {
	Pattern string
	File    string
	Line    int
	Content string
}

// Stats summarizes the work performed by ScanFile or ScanDir. Files is the
// number of files whose contents were scanned (binary files and skipped
// directories are excluded); Bytes is the total number of bytes scanned;
// ElapsedMs is the wall-clock duration of the scan in milliseconds.
//
// WalkErrors, when non-nil, lists per-file read errors collected by
// ScanDir -- files that the directory walk enumerated but whose contents
// could not be read (typically because of a permission change, a
// symlink whose target disappeared, or a TOCTOU race between WalkDir and
// ReadFile). ScanDir returns nil for its error in this case so callers
// that only want findings are unaffected; callers that care about every
// skipped file should inspect Stats.WalkErrors.
type Stats struct {
	Files      int
	Bytes      int64
	ElapsedMs  int64
	WalkErrors []error
}
