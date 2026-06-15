package patterns

import (
	"atheon/core"
	"regexp"
)

func init() {
	core.Register(&slackPattern{re: regexp.MustCompile(`xoxb-[0-9]{11}-[0-9]{11}-[0-9a-zA-Z]{24}`)})
}

type slackPattern struct{ re *regexp.Regexp }

func (p *slackPattern) Name() string             { return "slack-bot-token" }
func (p *slackPattern) Matches(line string) bool { return p.re.MatchString(line) }
