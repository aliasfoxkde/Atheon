// Concurrency regression test for the pattern_state RWMutex (Wave 7 PR #89).
//
// Without the RWMutex, concurrent EnablePattern/DisablePattern/SetPatternEnabled
// calls race against each other and against the Scan* goroutines that read
// activeScanners — the race detector flags torn *regexp.Regexp reads and,
// occasionally, slice out-of-bounds panics. Run with `go test -race` to
// verify both that this test catches regressions and that the production
// mutexes are actually doing their job.
//
// To prove the test would have caught the original bug, comment out the
// patternMu.RLock()/RUnlock() lines in core/runner.go:scanLines and rerun
// `go test -race ./core/` — the race detector reports "WARNING: DATA RACE"
// on activeScanners.
package core

import (
	"context"
	"sync"
	"testing"
)

// TestConcurrentEnableDisable exercises every write path against the read
// paths in ScanString. It runs under -race; without the patternMu guards
// added in Wave 7 PR #89, it reports a data race on activeScanners and
// (depending on timing) on bundlePattern.enabled.
//
// The iteration count is intentionally modest: every SetPatternEnabled
// rebuilds the combined-regex scanner via regexp.Compile, which is the
// expensive side. With 4 writers + 4 readers at 50 iters the whole test
// completes in <2s under -race; larger values risk the 30s timeout under
// contention on slow CI runners.
func TestConcurrentEnableDisable(t *testing.T) {
	patterns := All()
	if len(patterns) == 0 {
		t.Skip("no patterns loaded; can't exercise concurrent toggle")
	}
	target := ""
	for _, p := range patterns {
		target = p.Name()
		break
	}
	if target == "" {
		t.Fatal("no pattern name available")
	}

	const writers = 4
	const readers = 4
	const iters = 50

	var wg sync.WaitGroup
	wg.Add(writers + readers)

	for w := 0; w < writers; w++ {
		go func(seed int) {
			defer wg.Done()
			for i := 0; i < iters; i++ {
				if (seed+i)%2 == 0 {
					SetPatternEnabled(target, true)
				} else {
					SetPatternEnabled(target, false)
				}
			}
		}(w)
	}

	// Readers exercise scanLines via ScanString. Two distinct inputs make
	// sure the per-line inner pattern loop runs (no short-circuit on
	// combined-regex pre-filter when the line has no hits).
	for r := 0; r < readers; r++ {
		go func() {
			defer wg.Done()
			for i := 0; i < iters; i++ {
				_ = ScanString(context.Background(), "AKIA1234567890ABCDEF", "race-test")
				_ = ScanString(context.Background(), "ghp_aBcDeFgHiJkLmNoPqRsTuVwXyZ0123456789", "race-test")
			}
		}()
	}

	wg.Wait()

	// Final state must reflect one of the toggle calls (any value), not a
	// torn bool. We don't assert the specific value — under contention the
	// order is non-deterministic. We assert the API doesn't deadlock and
	// the pattern is still findable in the registry.
	if !patternExists(target) {
		t.Errorf("pattern %q disappeared from All() after concurrent toggles", target)
	}
}

func patternExists(name string) bool {
	for _, p := range All() {
		if p.Name() == name {
			return true
		}
	}
	return false
}
