package patterns

import (
    "atheon/core"
    "regexp"
)

func init() { core.Register(&ssnPattern{re: regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`)}) }
type ssnPattern struct{ re *regexp.Regexp }
func (p *ssnPattern) Name() string             { return "Social Security Number" }
func (p *ssnPattern) Matches(line string) bool { return p.re.MatchString(line) }
