package patterns

import (
	"atheon/core"
	"regexp"
)

func init() {
	core.Register(&gcpServiceAccountKeyPattern{
		privateKeyID: regexp.MustCompile(`"private_key_id"\s*:\s*"[0-9a-f]{40}"`),
		clientEmail:  regexp.MustCompile(`"client_email"\s*:\s*"[^"]+@[^"]+\.iam\.gserviceaccount\.com"`),
	})
	core.Register(&gcpPattern{re: regexp.MustCompile(`\bAIza[0-9A-Za-z\-_]{35}\b`), name: "gcp-api-key"})
	core.Register(&gcpPattern{re: regexp.MustCompile(`\b[a-z0-9][a-z0-9\-\.]*@[a-z0-9\-]+\.iam\.gserviceaccount\.com\b`), name: "gcp-service-account-email"})
	core.Register(&gcpPattern{re: regexp.MustCompile(`\b[0-9]+-[0-9a-z]+\.apps\.googleusercontent\.com\b`), name: "gcp-oauth-client-id"})
	core.Register(&gcpPattern{re: regexp.MustCompile(`\bGOCSPX-[A-Za-z0-9_\-]{28}\b`), name: "gcp-oauth-client-secret"})
}

type gcpServiceAccountKeyPattern struct {
	privateKeyID *regexp.Regexp
	clientEmail  *regexp.Regexp
}

func (p *gcpServiceAccountKeyPattern) Name() string { return "gcp-service-account-key" }
func (p *gcpServiceAccountKeyPattern) Matches(line string) bool {
	return p.privateKeyID.MatchString(line) || p.clientEmail.MatchString(line)
}

type gcpPattern struct {
	re   *regexp.Regexp
	name string
}

func (p *gcpPattern) Name() string             { return p.name }
func (p *gcpPattern) Matches(line string) bool { return p.re.MatchString(line) }
