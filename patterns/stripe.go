package patterns

import (
	"atheon/core"
	"regexp"
)

func init() {
	core.Register(&stripePattern{re: regexp.MustCompile(`\bsk_live_[0-9a-zA-Z]{24}\b`)})
}

type stripePattern struct{ re *regexp.Regexp }

func (p *stripePattern) Name() string             { return "stripe-secret-key" }
func (p *stripePattern) Matches(line string) bool { return p.re.MatchString(line) }
