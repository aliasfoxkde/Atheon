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

---

## License

MIT — Copyright © 2026 Dominick Yanez
