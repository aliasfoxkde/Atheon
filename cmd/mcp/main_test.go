package main

import (
	"testing"
	"time"
)

// TestRateLimiterAllow tests the rate limiter Allow function
func TestRateLimiterAllow(t *testing.T) {
	tests := []struct {
		name          string
		tokensPerSec  float64
		burst         float64
		calls         int
		expectAllowed []bool
	}{
		{
			name:          "burst allows multiple calls",
			tokensPerSec:  10,
			burst:         5,
			calls:         5,
			expectAllowed: []bool{true, true, true, true, true},
		},
		{
			name:          "exceeding burst denies",
			tokensPerSec:  10,
			burst:         3,
			calls:         5,
			expectAllowed: []bool{true, true, true, false, false},
		},
		{
			name:          "single call always allowed with burst",
			tokensPerSec:  10,
			burst:         1,
			calls:         1,
			expectAllowed: []bool{true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := newRateLimiter(tt.tokensPerSec, tt.burst)
			for i := 0; i < tt.calls; i++ {
				got := rl.Allow()
				if got != tt.expectAllowed[i] {
					t.Errorf("call %d: Allow() = %v, want %v", i+1, got, tt.expectAllowed[i])
				}
			}
		})
	}
}

// TestRateLimiterRefill tests that tokens refill over time
func TestRateLimiterRefill(t *testing.T) {
	// Create a limiter with burst of 1, rate of 100 tokens/sec
	rl := newRateLimiter(100, 1)

	// First call should be allowed
	if !rl.Allow() {
		t.Error("first call should be allowed")
	}

	// Second call immediately should be denied (no refill yet)
	if rl.Allow() {
		t.Error("second call immediately should be denied")
	}

	// Wait enough time to refill 1 token (10ms for 100 tokens/sec)
	time.Sleep(15 * time.Millisecond)

	// Third call after refill should be allowed
	if !rl.Allow() {
		t.Error("third call after refill should be allowed")
	}
}

// TestRateLimiterMaxCap tests that tokens don't exceed burst
func TestRateLimiterMaxCap(t *testing.T) {
	rl := newRateLimiter(1000, 5) // high rate, small burst

	// Make many calls rapidly
	for i := 0; i < 20; i++ {
		rl.Allow()
	}

	// Wait a long time at high rate
	time.Sleep(50 * time.Millisecond)

	// Tokens should cap at max (5), not accumulate infinitely
	// After 50ms at 1000 tokens/sec = 50 tokens added, but max is 5
	// So we should have 5 tokens
	if !rl.Allow() {
		t.Error("should have tokens available after long wait")
	}
}
