// Package errors provides shared error utilities for Atheon.
package errors

import (
	"os"
)

// SafeError returns a user-safe error message that does not leak filesystem
// paths or internal details. Used by both the CLI and MCP server to map
// OS-level errors to human-readable strings.
func SafeError(err error) string {
	if err == nil {
		return "unknown error"
	}
	switch {
	case os.IsNotExist(err):
		return "file not found"
	case os.IsPermission(err):
		return "permission denied"
	default:
		return "internal error"
	}
}
