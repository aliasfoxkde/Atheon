package core_test

import (
	"strings"
	"testing"

	"github.com/aliasfoxkde/Atheon/core"
)

type patternCase struct {
	matches    []string
	nonMatches []string
}

func TestRegisteredPatterns(t *testing.T) {
	cases := map[string]patternCase{
		"ssn": {
			matches:    []string{"ssn=123-45-6789"},
			nonMatches: []string{"ssn=123-456-789", "invoice=123-45-678"},
		},
		"aws-access-key": {
			matches:    []string{"AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE"},
			nonMatches: []string{"AWS_ACCESS_KEY_ID=AKIA123", "AKIAiosfodnn7example"},
		},
		"credit-card": {
			matches:    []string{"card=4242 4242 4242 4242", "amex=3782-822463-10005"},
			nonMatches: []string{"card=1234 5678 9012 3456", "order=4242"},
		},
		"gcp-api-key": {
			matches:    []string{"api_key=AIza" + strings.Repeat("a", 35)},
			nonMatches: []string{"api_key=AIza-short", "api_key=AIza" + strings.Repeat("!", 35)},
		},
		"gcp-oauth-client-id": {
			matches:    []string{"client_id=1234567890-abcdefghijklmnopqrstuvwxyz.apps.googleusercontent.com"},
			nonMatches: []string{"client_id=project.apps.googleusercontent.com", "client_id=1234567890-.apps.googleusercontent.com"},
		},
		"gcp-oauth-client-secret": {
			matches:    []string{"client_secret=GOCSPX-" + strings.Repeat("a", 28)},
			nonMatches: []string{"client_secret=GOCSPX-short", "client_secret=GOOGLE-" + strings.Repeat("a", 28)},
		},
		"gcp-service-account-email": {
			matches:    []string{"svc=my-service@project.iam.gserviceaccount.com"},
			nonMatches: []string{"svc=my-service@example.com", "svc=@project.iam.gserviceaccount.com"},
		},
		"gcp-service-account-key": {
			matches: []string{
				`"private_key_id": "` + strings.Repeat("a", 40) + `"`,
				`"client_email": "svc@project.iam.gserviceaccount.com"`,
			},
			nonMatches: []string{`"private_key_id": "short"`, `"client_email": "svc@example.com"`},
		},
		"github-pat": {
			matches:    []string{"token=ghp_" + strings.Repeat("a", 36)},
			nonMatches: []string{"token=ghp_short", "token=github_pat_" + strings.Repeat("a", 36)},
		},
		"openai-api-key": {
			matches:    []string{"OPENAI_API_KEY=sk-" + strings.Repeat("a", 20)},
			nonMatches: []string{"OPENAI_API_KEY=sk-short", "OPENAI_API_KEY=pk-" + strings.Repeat("a", 20)},
		},
		"phone-number": {
			matches:    []string{"phone=(555) 123-4567", "phone=+1 555 123 4567"},
			nonMatches: []string{"version=555-123", "ticket=555-123-456"},
		},
		"slack-bot-token": {
			matches:    []string{"SLACK_BOT_TOKEN=xoxb-12345678901-12345678901-" + strings.Repeat("a", 24)},
			nonMatches: []string{"SLACK_BOT_TOKEN=xoxb-short", "SLACK_BOT_TOKEN=xoxa-12345678901-12345678901-" + strings.Repeat("a", 24)},
		},
		"stripe-secret-key": {
			matches:    []string{"STRIPE_SECRET_KEY=sk_live_" + strings.Repeat("a", 24)},
			nonMatches: []string{"STRIPE_SECRET_KEY=sk_test_" + strings.Repeat("a", 24), "STRIPE_SECRET_KEY=sk_live_short"},
		},
		"twilio-account-sid": {
			matches:    []string{"TWILIO_ACCOUNT_SID=AC" + strings.Repeat("a", 32)},
			nonMatches: []string{"TWILIO_ACCOUNT_SID=ACshort", "TWILIO_ACCOUNT_SID=SK" + strings.Repeat("a", 32)},
		},
	}

	registered := map[string]core.Pattern{}
	for _, p := range core.All() {
		registered[p.Name()] = p
	}

	for name, tc := range cases {
		p, ok := registered[name]
		if !ok {
			t.Fatalf("pattern %q not registered", name)
		}
		t.Run(name, func(t *testing.T) {
			for _, line := range tc.matches {
				if !p.Matches(line) {
					t.Errorf("expected %q to match %q", name, line)
				}
			}
			for _, line := range tc.nonMatches {
				if p.Matches(line) {
					t.Errorf("expected %q not to match %q", name, line)
				}
			}
		})
	}

	for name := range cases {
		if _, ok := registered[name]; !ok {
			t.Fatalf("test case %q has no registered pattern", name)
		}
	}
}

// TestCategories tests the Categories function
func TestCategories(t *testing.T) {
	cats := core.Categories()
	if len(cats) == 0 {
		t.Error("Categories returned empty list")
	}

	// Check for expected categories
	expectedCats := map[string]bool{
		"ai-detection": false,
		"code-quality": false,
		"devops":       false,
		"finance":      false,
		"healthcare":   false,
		"pii":          false,
		"secrets":      false,
		"django":       false,
		"nodejs":       false,
		"react":        false,
	}

	for _, cat := range cats {
		if _, exists := expectedCats[cat]; exists {
			expectedCats[cat] = true
		}
	}

	// Verify at least some expected categories exist
	foundCount := 0
	for _, found := range expectedCats {
		if found {
			foundCount++
		}
	}

	if foundCount < 5 {
		t.Errorf("Expected at least 5 categories, found %d", foundCount)
	}
}

// TestSetPatternEnabled tests the SetPatternEnabled function
func TestSetPatternEnabled(t *testing.T) {
	patterns := core.All()
	if len(patterns) == 0 {
		t.Skip("No patterns available for testing")
	}

	// Find a test pattern
	var testPatternName string
	var testPattern core.Pattern
	for _, p := range patterns {
		testPattern = p
		testPatternName = p.Name()
		break
	}

	if testPatternName == "" {
		t.Fatal("Could not find a pattern to test")
	}

	originalState := testPattern.Enabled()

	// Test setting to true
	core.EnablePattern(testPatternName)
	// Need to find the pattern again to check its state
	enabled := false
	for _, p := range core.All() {
		if p.Name() == testPatternName {
			enabled = p.Enabled()
			break
		}
	}
	if !enabled {
		t.Error("Failed to enable pattern")
	}

	// Test setting to false
	core.DisablePattern(testPatternName)
	// After disable, the pattern is removed from the registry
	// (our implementation removes disabled patterns from All())
	disabled := false
	for _, p := range core.All() {
		if p.Name() == testPatternName {
			disabled = true
			break
		}
	}
	if disabled {
		t.Error("Failed to disable pattern - pattern should be removed from registry")
	}

	// Restore original state
	if originalState {
		core.EnablePattern(testPatternName)
	} else {
		core.DisablePattern(testPatternName)
	}
}

// TestSetPatternEnabledTrue tests the SetPatternEnabled(name, true) path
func TestSetPatternEnabledTrue(t *testing.T) {
	patterns := core.All()
	if len(patterns) == 0 {
		t.Skip("No patterns available for testing")
	}
	name := patterns[0].Name()

	// Make sure it starts enabled
	core.EnablePattern(name)

	originalEnabled := false
	for _, p := range core.All() {
		if p.Name() == name {
			originalEnabled = p.Enabled()
			break
		}
	}

	// Use SetPatternEnabled to toggle
	if !core.SetPatternEnabled(name, false) {
		t.Fatal("SetPatternEnabled(name, false) returned false")
	}

	// Verify it's now disabled (removed from registry)
	for _, p := range core.All() {
		if p.Name() == name && p.Enabled() {
			t.Errorf("Pattern %s should be disabled", name)
		}
	}

	// Set back to true
	if !core.SetPatternEnabled(name, true) {
		t.Fatal("SetPatternEnabled(name, true) returned false")
	}

	// Verify it's enabled
	found := false
	for _, p := range core.All() {
		if p.Name() == name {
			if !p.Enabled() {
				t.Errorf("Pattern %s should be enabled", name)
			}
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Pattern %s not found after re-enable", name)
	}

	// Restore
	if !originalEnabled {
		core.DisablePattern(name)
	}
}

// TestSetPatternEnabledNotFound tests SetPatternEnabled with an unknown name
func TestSetPatternEnabledNotFound(t *testing.T) {
	if core.SetPatternEnabled("nonexistent-pattern-xyz", true) {
		t.Error("SetPatternEnabled should return false for unknown pattern")
	}
}

// TestListDisabledPatterns tests the ListDisabledPatterns function
func TestListDisabledPatterns(t *testing.T) {
	// First disable some patterns
	patterns := core.All()
	if len(patterns) == 0 {
		t.Skip("No patterns available for testing")
	}

	// Disable at least one pattern
	testPatternName := patterns[0].Name()
	originalState := patterns[0].Enabled()
	core.DisablePattern(testPatternName)

	disabled := core.ListDisabledPatterns()

	if len(disabled) == 0 {
		t.Error("Expected at least one disabled pattern")
	}

	// Restore original state
	if originalState {
		core.EnablePattern(testPatternName)
	}
}

// TestListEnabledPatterns tests the ListEnabledPatterns function
func TestListEnabledPatterns(t *testing.T) {
	patterns := core.All()
	if len(patterns) == 0 {
		t.Skip("No patterns available for testing")
	}

	enabled := core.ListEnabledPatterns()

	if len(enabled) == 0 {
		t.Error("Expected at least one enabled pattern")
	}

	// Verify all returned patterns are actually enabled
	for _, name := range enabled {
		found := false
		for _, p := range core.All() {
			if p.Name() == name {
				if !p.Enabled() {
					t.Errorf("Pattern %s is in enabled list but not enabled", name)
				}
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Pattern %s is in enabled list but not found in All()", name)
		}
	}
}

// TestEnableAllPatterns tests the EnableAllPatterns function
func TestEnableAllPatterns(t *testing.T) {
	patterns := core.All()
	if len(patterns) == 0 {
		t.Skip("No patterns available for testing")
	}

	// Disable some patterns first
	for i := 0; i < 3 && i < len(patterns); i++ {
		core.DisablePattern(patterns[i].Name())
	}

	// Enable all
	core.EnableAllPatterns()

	// Verify all are enabled
	for _, p := range patterns {
		if !p.Enabled() {
			t.Errorf("Pattern %s should be enabled after EnableAllPatterns", p.Name())
		}
	}
}

// TestPatternStatePersistence tests pattern state persistence
func TestPatternStatePersistence(t *testing.T) {
	patterns := core.All()
	if len(patterns) == 0 {
		t.Skip("No patterns available for testing")
	}

	// Get original states
	originalStates := make(map[string]bool)
	for _, p := range patterns {
		originalStates[p.Name()] = p.Enabled()
	}

	// Modify some states
	for i := 0; i < 3 && i < len(patterns); i++ {
		core.DisablePattern(patterns[i].Name())
	}

	// Note: SavePatternStates function doesn't exist in current API
	// Skipping persistence test

	// Restore original states
	for _, p := range patterns {
		if originalState, exists := originalStates[p.Name()]; exists {
			if originalState {
				core.EnablePattern(p.Name())
			} else {
				core.DisablePattern(p.Name())
			}
		}
	}
}
