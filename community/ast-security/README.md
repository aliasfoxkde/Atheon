# AST Security Patterns

This directory contains **built-in AST-based security patterns** for Go code analysis.

## Overview

Unlike regex-based patterns in other categories, AST patterns analyze the **structure** of Go source code using the `go/ast` package. This enables detection of complex vulnerabilities that text-based regex cannot find, such as:

- Command injection via `exec.Command` with string concatenation
- SQL injection through AST analysis of query construction
- Path traversal detection in file operations
- Template injection vulnerabilities
- ReDoS (Regular Expression Denial of Service)
- Weak cryptography usage

## Built-in Patterns

All AST security patterns are **built into `core/ast_patterns.go`** and are automatically enabled with the `--ast` flag. They are not loaded from YAML files.

To enable AST scanning:

```bash
atheon ./path/to/code --ast
```

## Pattern List

| Pattern | Severity | Description |
|---------|----------|-------------|
| `go-command-injection` | CRITICAL | exec.Command with user input concatenation |
| `go-shell-command` | CRITICAL | Shell invocation with user input |
| `go-sql-injection` | CRITICAL | String concatenation in SQL query |
| `go-sql-template-query` | HIGH | Query method with user-controlled argument |
| `go-path-traversal` | HIGH | File operation with user-controlled path |
| `go-symlink-attack` | MEDIUM | File open without O_NOFOLLOW flag |
| `go-unsafe-deserialization` | HIGH | Binary unmarshal with untrusted data |
| `go-gob-deserialization` | HIGH | gob decoding with untrusted data |
| `go-ssrf` | HIGH | HTTP request to user-controlled URL |
| `go-http-unvalidated-url` | MEDIUM | http.Get/Post with user-provided URL |
| `go-template-injection` | HIGH | Template execution with user data |
| `go-template-raw-html` | HIGH | template.HTML bypasses auto-escaping |
| `go-redos` | MEDIUM | Regex with nested quantifiers (ReDoS) |
| `go-regex-dynamic` | HIGH | Dynamic regex from user input |
| `go-hardcoded-secret` | HIGH | Credential as string literal |
| `go-private-key` | CRITICAL | Embedded private key/certificate |
| `go-weak-crypto-md5` | MEDIUM | Use of MD5 (broken for security) |
| `go-weak-crypto-sha1` | MEDIUM | Use of SHA-1 (deprecated) |
| `go-insecure-random` | MEDIUM | math/rand for security randomness |
| `go-weak-cipher` | HIGH | DES, RC4, or ECB mode |
| `go-unchecked-error` | MEDIUM | Error return not checked |
| `go-silent-panic` | MEDIUM | panic() in production code |
| `go-ldap-injection` | HIGH | LDAP query string concatenation |
| `go-xxe` | HIGH | XML external entity enabled |
| `go-yaml-unsafe` | HIGH | yaml.Unmarshal on untrusted data |
| `go-trust-boundary` | MEDIUM | User input to internal state |
| `go-tls-skip-verify` | CRITICAL | TLS verification disabled |
| `go-insecure-tls` | HIGH | Weak TLS configuration |

## Future Plans

Future versions will allow custom AST patterns via YAML files in this directory, enabling:

- Custom AST-based detections
- Language-specific analyzers
- Taint tracking configurations
