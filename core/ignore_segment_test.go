package core

import (
	"strings"
	"testing"
)

// TestWriteIgnoreSegment exercises all branches of writeIgnoreSegment:
// plain literals, regex-special chars that need escaping, single-star (?),
// double-star (**), ?, and character classes including negation.
func TestWriteIgnoreSegment(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"plain literal", "abc", "abc"},
		{"dot escape", "a.b", `a\.b`},
		{"plus escape", "a+b", `a\+b`},
		{"caret escape", "^abc", `\^abc`},
		{"dollar escape", "abc$", `abc\$`},
		{"brace escape", "{a,b}", `\{a,b\}`},
		{"paren escape", "(a)", `\(a\)`},
		{"pipe escape", "a|b", `a\|b`},
		{"backslash escape", `a\b`, `a\\b`},
		{"single star", "a*b", `a[^/]*b`},
		{"double star", "a**b", `a.*b`},
		{"question mark", "a?b", `a[^/]b`},
		{"char class", "[abc]", `[abc]`},
		{"negated char class", "[!abc]", `[^abc]`},
		{"unterminated char class", "[abc", `\[abc`},
		{"mixed", "*.go", `[^/]*\.go`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var b strings.Builder
			writeIgnoreSegment(tc.in, &b)
			if b.String() != tc.want {
				t.Errorf("writeIgnoreSegment(%q) = %q, want %q", tc.in, b.String(), tc.want)
			}
		})
	}
}

// TestIgnorePatternToRegexpExhaustive exercises various gitignore-style
// patterns to drive coverage of the glob translation logic.
func TestIgnorePatternToRegexpExhaustive(t *testing.T) {
	cases := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"*.go", "main.go", true},
		{"*.go", "main.txt", false},
		{"**/foo", "foo", true},
		{"**/foo", "a/b/c/foo", true},
		{"docs/**", "docs/a/b", true},
		{"docs/**", "src/docs/a", false},
		{"build/", "build", true},
		{"build/", "build/x", true},
		{"build/", "buildfile", false},
		{"/root-only", "root-only", true},
		{"/root-only", "a/root-only", false},
		{"file?.txt", "file1.txt", true},
		{"file?.txt", "file12.txt", false},
	}

	for _, tc := range cases {
		t.Run(tc.pattern, func(t *testing.T) {
			re, err := ignorePatternToRegexp(tc.pattern)
			if err != nil {
				t.Fatalf("compile error: %v", err)
			}
			got := re.MatchString(tc.input)
			if got != tc.want {
				t.Errorf("pattern %q against %q = %v, want %v", tc.pattern, tc.input, got, tc.want)
			}
		})
	}
}

// TestIgnorePatternEmpty exercises the empty-pattern error branch.
func TestIgnorePatternEmpty(t *testing.T) {
	_, err := ignorePatternToRegexp("")
	if err == nil {
		t.Error("expected error for empty pattern")
	}
	_, err = ignorePatternToRegexp("/")
	if err == nil {
		t.Error("expected error for slash-only pattern")
	}
}
