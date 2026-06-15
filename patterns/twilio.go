package patterns

import (
	"atheon/core"
	"regexp"
)

func init() {
	core.Register(&twilioPattern{re: regexp.MustCompile(`AC[0-9a-fA-F]{32}`)})
}

type twilioPattern struct{ re *regexp.Regexp }

func (p *twilioPattern) Name() string             { return "twilio-account-sid" }
func (p *twilioPattern) Matches(line string) bool { return p.re.MatchString(line) }
