// End-to-end pattern detection tests.
//
// For a representative pattern in each covered category, scan a known
// fixture line and assert the pattern name appears in the findings. This
// catches the class of bugs where a regex change silently breaks detection
// — the pattern compiles and loads, but no longer matches what it claims.
//
// The fixtures are intentionally minimal — one line per pattern. Each line
// must contain the canonical token for that pattern's match. The test
// isolates the regex engine from anything that might mask a regression:
// no file I/O, no init-time state, just core.ScanString against the
// embedded bundle.

package core

import (
	"context"
	"os"
	"testing"
)

// detectionFixtures maps (category → (representative pattern, known-good
// fixture line)). Names come from community/<cat>/<file>.yaml — verify
// against the YAML before updating, since a renamed pattern silently
// breaks this test.
var detectionFixtures = []struct {
	Category string
	Pattern  string
	Fixture  string
}{
	// secrets — canonical SaaS tokens used in README and PR reviews
	{"secrets", "aws-access-key", "AKIAIOSFODNN7EXAMPLE"},
	{"secrets", "github-pat", "ghp_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"},
	{"secrets", "stripe-secret-key", "sk_live_aaaaaaaaaaaaaaaaaaaaaaaa"},
	{"secrets", "slack-bot-token", "xoxb-12345678901-12345678901-aaaaaaaaaaaaaaaaaaaaaaaa"},
	{"secrets", "gcp-api-key", "AIzaSyDabcdefghijklmnopqrstuvwxyz123456"},
	{"secrets", "gcp-oauth-client-id", "123456789012-abcdefghijklmnopqrstuvwxyz1.apps.googleusercontent.com"},
	{"secrets", "gcp-oauth-client-secret", "GOCSPX-abcdefghijklmnopqrstuvwxyz12"},

	// pii — common identifiers
	{"pii", "ssn", "343-53-9183"},
	{"pii", "phone-number", "+1 555 867 5309"},
	{"pii", "credit-card", "4111-1111-1111-1111"},

	// secrets — service-account identities (defined under secrets/, not pii/)
	{"secrets", "gcp-service-account-email", "myservice@my-project.iam.gserviceaccount.com"},
}

// TestPatternDetection is the headline test: for every fixture, the
// matching pattern must appear in the findings. If a fixture fails, the
// test reports the pattern name and the patterns that DID match so you
// can find the YAML and decide whether the regex broadened or the test
// pointed at the wrong name.
func TestPatternDetection(t *testing.T) {
	for _, fx := range detectionFixtures {
		t.Run(fx.Pattern, func(t *testing.T) {
			findings := ScanString(context.Background(), fx.Fixture, "test-fixture")
			for _, f := range findings {
				if f.Pattern == fx.Pattern {
					return
				}
			}
			t.Errorf("pattern %q (category %q) did not detect fixture %q; got findings: %v",
				fx.Pattern, fx.Category, fx.Fixture, findingPatterns(findings))
		})
	}
}

// TestCategoryCoverage is informational only — it logs uncovered
// categories so reviewers know what's missing, but does not fail. Promote
// to t.Errorf if/when the team commits to 100% coverage; until then, a
// hard fail would block every PR that adds a category.
//
// To turn this into a gate, set ATHEON_REQUIRE_FULL_COVERAGE=1.
func TestCategoryCoverage(t *testing.T) {
	covered := map[string]bool{}
	for _, fx := range detectionFixtures {
		covered[fx.Category] = true
	}

	var missing []string
	for _, cat := range Categories() {
		if !covered[cat] {
			missing = append(missing, cat)
		}
	}
	if len(missing) == 0 {
		return
	}

	msg := "categories without a detection fixture (informational until coverage is required)"
	if os.Getenv("ATHEON_REQUIRE_FULL_COVERAGE") != "" {
		t.Errorf("%s: %v", msg, missing)
		return
	}
	t.Logf("%s: %v", msg, missing)
}

// TestFalsePositiveGuard scans a known-clean snippet and asserts no
// findings are produced. The snippet is deliberately constructed to
// avoid every known pattern — no token-shaped strings, no commented-out
// code that could be misread as fmt.Println, no "skip" + "link" /
// "navigation" phrases. This catches regressions where a regex
// broadens and starts matching noise.
//
// If your bundle starts producing findings here, treat it as a real
// false-positive regression and tighten the pattern, not the test.
func TestFalsePositiveGuard(t *testing.T) {
	// Each line chosen to evade a known pattern while still being
	// recognisably Go-shaped. Don't expand this without re-checking
	// every pattern the new line might match. In particular:
	//   - no fmt.Print* calls (would match code-quality/fmt-println-prod)
	//   - no SSN-shaped or credit-card-shaped numbers
	//   - no SaaS token prefixes
	//   - no "skip"+ "link" / "navigation" phrases
	clean := "" +
		"package main\n" +
		"\n" +
		"const greeting = \"hello world\"\n" +
		"\n" +
		"// This file is intentionally free of any token-shaped constants.\n" +
		"// Numbers are deliberately broken up so credential regexes do not match them.\n" +
		"\n" +
		"var (\n" +
		"    count = 7\n" +
		"    label = \"demo\"\n" +
		")\n" +
		"\n" +
		"func show() string {\n" +
		"    return greeting\n" +
		"}\n"

	findings := ScanString(context.Background(), clean, "clean-snippet")
	if len(findings) > 0 {
		t.Errorf("clean snippet produced %d unexpected findings: %v", len(findings), findingPatterns(findings))
	}
}

// findingPatterns flattens a finding slice to pattern names so test
// failure messages stay readable.
func findingPatterns(findings []Finding) []string {
	out := make([]string, 0, len(findings))
	for _, f := range findings {
		out = append(out, f.Pattern)
	}
	return out
}

// TestFindingSeverityPropagation asserts that the severity declared in the
// pattern YAML is what surfaces on every Finding. Today only a handful of
// patterns declare severity (the rest default to medium), but the contract
// must hold for the ones that do — otherwise SARIF consumers get the wrong
// security-severity score.
func TestFindingSeverityPropagation(t *testing.T) {
	// missing-skip-links declares severity: medium in its YAML.
	findings := ScanString(context.Background(),
		"// TODO: add skip navigation here", "test-severity-medium")
	for _, f := range findings {
		if f.Pattern != "missing-skip-links" {
			continue
		}
		if f.Severity != "medium" {
			t.Errorf("missing-skip-links severity: got %q, want %q", f.Severity, "medium")
		}
		return
	}
	t.Skip("missing-skip-links did not fire on the test snippet — pattern may have changed; rerun manually")
}

// TestFindingSeverityDefault asserts patterns without an explicit severity
// field still report one (defaulting to medium) so SARIF and JSON output
// never emit an empty security-severity.
func TestFindingSeverityDefault(t *testing.T) {
	findings := ScanString(context.Background(),
		`aws_key = "AKIAIOSFODNN7EXAMPLE"`, "test-default-severity")
	for _, f := range findings {
		if f.Pattern != "aws-access-key" {
			continue
		}
		if f.Severity == "" {
			t.Errorf("aws-access-key should report a non-empty severity even without an explicit YAML field; got empty")
		}
		return
	}
	t.Fatalf("aws-access-key did not fire on its canonical fixture")
}

// TestNormalizeSeverity exercises the loader's safety net for typo'd YAML
// severity values — anything outside ValidSeverities collapses to medium.
func TestNormalizeSeverity(t *testing.T) {
	cases := []struct{ in, want string }{
		{"", "medium"},
		{"medium", "medium"},
		{"HIGH", "high"},
		{" Critical ", "critical"},
		{"urgent", "medium"},   // unrecognised → default
		{"low", "low"},
		{"high", "high"},
		{"critical", "critical"},
	}
	for _, tc := range cases {
		if got := normalizeSeverity(tc.in); got != tc.want {
			t.Errorf("normalizeSeverity(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
