# Leakr Roadmap

Leakr is a finished product.

The CLI is done. The library API is done. The scanner engine, parallel scanning, JSON output, exit codes, environment variable scanning, stdin support — all of it is done. The core will not change.

**The only thing Leakr will ever need is more scanners.**

---

## How to contribute

A new scanner is a single `.java` file in `src/leakr/scanners/`. Implement the `Scanner` interface, rebuild, and it works. See the README for the full walkthrough.

---

## Scanner additions 

Scanners to add over time, roughly by priority:

- `openai-api-key` — OpenAI API keys (`sk-...`)
- `private-key` — PEM private key blocks (`-----BEGIN RSA PRIVATE KEY-----`, EC, etc.)
- `jwt-token` — JSON Web Tokens (`eyJ...`)
- `google-service-account` — GCP service account JSON credential fragments
- `azure-storage-key` — Azure storage connection strings
- `azure-subscription-key` — Azure Cognitive Services keys
- `databricks-token` — Databricks personal access tokens (`dapi...`)
- `vault-token` — HashiCorp Vault tokens (`hvs.`)
- `npm-token` — npm publish tokens (`npm_...`)
- `pypi-token` — PyPI API tokens (`pypi-...`)
- `sendgrid-api-key` — SendGrid API keys (`SG.`)
- `mailgun-api-key` — Mailgun API keys
- `shopify-access-token` — Shopify Admin API access tokens
- `okta-api-token` — Okta API tokens
- `datadog-api-key` — Datadog API keys
- `heroku-api-key` — Heroku API keys
- `netlify-access-token` — Netlify personal access tokens
- `digitalocean-token` — DigitalOcean personal access tokens
- `telegram-bot-token` — Telegram bot tokens
- `discord-bot-token` — Discord bot tokens
- `firebase-server-key` — Firebase server keys
- `square-access-token` — Square OAuth access tokens
- `twitter-bearer-token` — Twitter/X bearer tokens
- `cloudflare-api-key` — Cloudflare API keys
- `basic-auth-url` — Credentials embedded in URLs (`https://user:pass@...`)

---

## What is NOT on the roadmap

These will never be added:

- Git history scanning — use [gitleaks](https://github.com/gitleaks/gitleaks)
- Pre-commit hook installer
- SARIF output format
- GitHub Actions workflow files
- Plugin or rule DSL system
- Server mode or daemon
- Config file format for defining patterns

Leakr scans content for secrets. That is all it does, and it does it well.
