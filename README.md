# Atheon

![Go](https://img.shields.io/badge/Go-1.21%2B-00ADD8)
![License](https://img.shields.io/badge/license-MIT-green)

**A pattern matching engine. Define what you're looking for. Point it at anything.**

---

## Download

Grab the binary for your platform from [Releases](https://github.com/HoraDomu/Atheon/releases/latest) — no install, no runtime required.

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

---

## Build

```
go build -o atheon .
```

Cross-compile:

```
GOOS=windows GOARCH=amd64 go build -o atheon.exe .
GOOS=linux   GOARCH=amd64 go build -o atheon-linux .
GOOS=darwin  GOARCH=arm64 go build -o atheon-macos .
```

---

## License

MIT — Copyright © 2026 Dominick Yanez
