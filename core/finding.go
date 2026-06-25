package core

// Finding represents a single pattern match produced by ScanFile, ScanDir,
// ScanString, or ScanEnv. File is the source path (or "env:KEY" when
// scanning the process environment); Line is the 1-indexed line number
// within the source (0 for env scans); Content is the trimmed matching
// line or, for env scans, the matching value.
//
// Severity is the pattern's declared severity at the time of the match —
// one of "low", "medium", "high", "critical". It's copied off the Pattern
// at match time so toggling severity later doesn't rewrite history.
type Finding struct {
	Pattern  string
	File     string
	Line     int
	Content  string
	Severity string
}

// Stats summarizes the work performed by ScanFile or ScanDir. Files is the
// number of files whose contents were scanned (binary files and skipped
// directories are excluded); Bytes is the total number of bytes scanned;
// ElapsedMs is the wall-clock duration of the scan in milliseconds.
// Errors collects any per-file read errors encountered during a ScanDir
// walk so the caller can surface them instead of silently dropping them.
type Stats struct {
	Files     int
	Bytes     int64
	ElapsedMs int64
	Errors    []error
}
