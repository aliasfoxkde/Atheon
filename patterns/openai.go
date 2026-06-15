package patterns

import (
	"atheon/core"
	"regexp"
)

func init() {
	core.Register(&openaiPattern{re: regexp.MustCompile(`\bsk-[A-Za-z0-9_\-]{20,}\b`)})
}

type openaiPattern struct{ re *regexp.Regexp }

func (p *openaiPattern) Name() string             { return "openai-api-key" }
func (p *openaiPattern) Matches(line string) bool { return p.re.MatchString(line) }
