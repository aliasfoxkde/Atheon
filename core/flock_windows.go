//go:build windows

package core

// withFileLock is a no-op on Windows. syscall.Flock is not implemented on
// Windows in Go's syscall package; the underlying LockFileEx API behaves
// differently (mandatory vs advisory locking) and is not portable to
// cross-compiled targets. The atomicWriteFile rename still gives us
// per-process crash safety; cross-process concurrent saves on Windows are
// best-effort and rely on the user not running two writers simultaneously.
//
// If Windows users report lost preferences, the right fix is a portable
// LockFileEx wrapper — not forcing flock semantics through the Go runtime.
func withFileLock(path string, fn func() error) error {
	return fn()
}
