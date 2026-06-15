package patterns

import (
	"atheon/core"
	"regexp"
)

func init() {
	core.Register(&awsPattern{re: regexp.MustCompile(`\b(?:AKIA|ASIA)[0-9A-Z]{16}\b`)})
}

type awsPattern struct{ re *regexp.Regexp }

func (p *awsPattern) Name() string             { return "aws-access-key" }
func (p *awsPattern) Matches(line string) bool { return p.re.MatchString(line) }
