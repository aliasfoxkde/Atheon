# Contributing

Atheon grows through patterns. Every pattern is one file with two methods, so keep contributions small, focused, and easy to review.

## Pattern workflow

1. **Define what you are detecting**
   - What does it look like: a fixed prefix, structural shape, or known format?
   - Why does it matter: leaked credential, compliance violation, or prohibited string?

2. **Check it does not already exist**
   ```sh
   atheon list
   ```

3. **Create the pattern file**
   Add a new `.go` file in `patterns/`, named after what it detects.
   ```go
   package patterns

   import (
       "atheon/core"
       "regexp"
   )

   func init() { core.Register(&examplePattern{re: regexp.MustCompile(`your-regex-here`)}) }
   type examplePattern struct{ re *regexp.Regexp }
   func (p *examplePattern) Name() string             { return "my-pattern-name" }
   func (p *examplePattern) Matches(line string) bool { return p.re.MatchString(line) }
   ```
   Use lowercase hyphenated names, and be specific: `stripe-live-key`, not `stripe`.

4. **Build and confirm it loaded**
   ```sh
   go build -o atheon . || go build .
   atheon list
   ```

5. **Add a test case**
   Open `patterns/patterns_test.go` and add an entry for your pattern under the `cases` map:
   ```go
   "my-pattern-name": {
       matches:    []string{"line that must match"},
       nonMatches: []string{"line that must not match"},
   },
   ```
   The test suite enforces that every registered pattern has a case — `go test ./patterns` will fail without it.

6. **Verify manually**
   Run against a sample file or directory to confirm the output looks right:
   ```sh
   atheon --file <path>
   ```
   Every expected match should appear, with no unexpected matches.

7. **Submit the contribution**
   Open a pull request with what the pattern detects, why it matters, and the test cases you used.

Maintainers review for correctness, false positive rate, name clarity, and overlap with existing patterns.
