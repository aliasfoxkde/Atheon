# Atheon

![Go](https://img.shields.io/badge/Go-1.21%2B-00ADD8)
![License](https://img.shields.io/badge/license-MIT-green)

> **Status:** Feature complete. The engine is done. What grows from here are patterns and bug fixes — nothing more, nothing less.

---

**One tool. All patterns. Any input.**

Atheon is a pattern matching engine. You tell it what to look for. You point it at anything. It finds every match and tells you exactly where.

---

## Why this matters

Data ends up where it shouldn't. A hardcoded credential in a config file. A production secret in a log. A sensitive string committed into a repository by accident and now permanently in git history. These mistakes happen constantly — across every team, every stack, every domain.

The problem isn't that people are careless. The problem is there's no systematic way to catch what you can't see.

Atheon is that system. A pattern matching engine you define, run anywhere, and trust completely — because you wrote the rules.

---

## What pattern matching means

A pattern is a rule: "if a line looks like this, flag it." That rule can be a regex, a keyword check, a structural test — anything that returns true or false. Every pattern has a name. Every match tells you the file, the line, and what was found.

The engine itself is deliberately minimal. It doesn't know what a secret is, what compliance means, or what matters to your organization. You do. So you define it, and the engine enforces it — over files, directories, environment variables, or any stream of text piped through it.

Pattern matching is useful in any domain where text contains something that shouldn't be there, or something that must be there. Security. Compliance. Legal. Operations. Healthcare. Finance. If you can describe the rule, Atheon can run it.

---

## The scenario that makes this real

A developer wraps up a sprint and pushes a configuration file. Inside it, buried in a comment from a debugging session three weeks ago, is a production API key. The commit goes through. The pipeline passes. The key is now in git history, in the build artifact, and eventually in a production image. Someone rotates it two months later after a billing alert.

Atheon, wired into a pre-push hook:

```
$ atheon ./

[api-key] config/app.yaml:47  →  # debug key: sk-prod-a8f3c...
```

Exit code `1`. The push never happens. The key never leaves the machine.

That's it. That's the product.

---

## Install

Download the binary for your platform from [Releases](https://github.com/HoraDomu/Atheon/releases/latest). No install, no runtime, no dependencies. Drop it in your PATH and run it.

**Or build from source:**

```
go build -o atheon .
```

Cross-compile for any platform:

```
GOOS=windows GOARCH=amd64 go build -o atheon.exe .
GOOS=linux   GOARCH=amd64 go build -o atheon-linux .
GOOS=darwin  GOARCH=arm64 go build -o atheon-macos .
```

---

## Usage

```
atheon <path>          scan a directory
atheon --file <path>   scan a single file
atheon --env           scan environment variables
atheon list            list loaded patterns
```

Pipe support:

```
cat file.txt | atheon -
```

Exit code `0` = clean. Exit code `1` = findings. CI-friendly by default.

---

## Adding a pattern

One file. Two methods.

```go
package patterns

import (
    "atheon/core"
    "regexp"
)

func init() {
    core.Register(&myPattern{re: regexp.MustCompile(`your-regex-here`)})
}

type myPattern struct{ re *regexp.Regexp }

func (p *myPattern) Name() string             { return "my-pattern-name" }
func (p *myPattern) Matches(line string) bool { return p.re.MatchString(line) }
```

Drop the file in `patterns/`, rebuild. It appears in `atheon list` automatically.

The same two methods work for anything — credentials, PII, internal token formats, compliance markers, prohibited strings. If you can describe the rule, this is all the code it takes.

---

## Contributing

Atheon is not looking for new features. The engine is done.

What it will always accept:
- **Bug fixes** — if something behaves incorrectly, open an issue and it will be addressed
- **New patterns** — if you have a pattern worth adding, open an issue describing what it detects and why it matters

To contribute:
- Open an issue on [GitHub](https://github.com/HoraDomu/Atheon/issues) describing the bug or pattern
- Or email directly: [dommcpro@gmail.com](mailto:dommcpro@gmail.com)

Issues are reviewed and addressed by maintainers. The simpler and more focused the contribution, the faster it moves.

---

## License

MIT — Copyright © 2026 Dominick Yanez
