package patterns

import (
	"atheon/core"
	"regexp"
)

func init() {
	core.Register(&githubPattern{re: regexp.MustCompile(`ghp_[0-9a-zA-Z]{36}`)})
}

type githubPattern struct{ re *regexp.Regexp }

func (p *githubPattern) Name() string             { return "github-pat" }
func (p *githubPattern) Matches(line string) bool { return p.re.MatchString(line) }
