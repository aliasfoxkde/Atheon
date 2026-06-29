package core

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ASTFinding represents a finding from AST-based analysis.
type ASTFinding struct {
	File     string
	Line     int
	Column   int
	Rule     string
	Message  string
	Severity string
	Suggest  string // Optional suggested fix
}

// ToFinding converts an ASTFinding to a core Finding for unified output.
func (f ASTFinding) ToFinding() Finding {
	return Finding{
		Pattern:     "ast:" + f.Rule,
		File:        f.File,
		Line:        f.Line,
		Column:      f.Column,
		Content:     f.Message,
		Severity:    f.Severity,
		Category:    "ast-security",
		Description: f.Message,
	}
}

// ASTPattern defines an AST-based pattern using Go AST traversal.
type ASTPattern struct {
	Name        string
	Description string
	Severity    string
	Suggest     string
	Func        func(fset *token.FileSet, file *ast.File) []ASTFinding
	// FileFilter restricts this pattern to files matching these extensions.
	// Empty means all files.
	FileFilter []string
}

// builtinASTPatterns contains built-in Go AST security patterns.
var builtinASTPatterns = []ASTPattern{
	// =========================================================================
	// COMMAND INJECTION
	// =========================================================================
	{
		Name:        "go-command-injection",
		Description: "exec.Command with string concatenation or user input in command arguments",
		Severity:    "critical",
		Suggest:     "Use exec.Command with separate arguments instead of shell string",
		FileFilter:  []string{".go"},
		Func:        detectGoCommandInjection,
	},
	{
		Name:        "go-shell-command",
		Description: "os/exec: shell invocation (bash -c, sh -c) with user input",
		Severity:    "critical",
		Suggest:     "Avoid shell invocation; use exec.Command with separate args",
		FileFilter:  []string{".go"},
		Func:        detectGoShellCommand,
	},

	// =========================================================================
	// SQL INJECTION
	// =========================================================================
	{
		Name:        "go-sql-injection",
		Description: "String concatenation or formatting in SQL query construction",
		Severity:    "critical",
		Suggest:     "Use parameterized queries; never concatenate user input into SQL",
		FileFilter:  []string{".go"},
		Func:        detectGoSQLInjection,
	},
	{
		Name:        "go-sql-template-query",
		Description: "Database query method called with concatenated arguments",
		Severity:    "high",
		Suggest:     "Use parameterized query methods (?, $1) instead of string formatting",
		FileFilter:  []string{".go"},
		Func:        detectGoSQLTemplateQuery,
	},

	// =========================================================================
	// PATH TRAVERSAL
	// =========================================================================
	{
		Name:        "go-path-traversal",
		Description: "File operation (os.Open, ioutil.ReadFile, etc.) with user-controlled path",
		Severity:    "high",
		Suggest:     "Validate and sanitize user input before using in file paths; use filepath.Clean and scope to allowed directories",
		FileFilter:  []string{".go"},
		Func:        detectGoPathTraversal,
	},
	{
		Name:        "go-symlink-attack",
		Description: "File operation may be susceptible to symlink attacks",
		Severity:    "medium",
		Suggest:     "Use O_NOFOLLOW or check symlink targets before accessing",
		FileFilter:  []string{".go"},
		Func:        detectGoSymlinkAttack,
	},

	// =========================================================================
	// DESERIALIZATION
	// =========================================================================
	{
		Name:        "go-unsafe-deserialization",
		Description: "encoding.BinaryUnmarshaler or encoding.TextUnmarshaler with user input",
		Severity:    "high",
		Suggest:     "Validate input source; prefer safe alternatives like encoding/json",
		FileFilter:  []string{".go"},
		Func:        detectGoUnsafeDeserialization,
	},
	{
		Name:        "go-gob-deserialization",
		Description: "gob.NewDecoder or gob.NewDecoder with untrusted data",
		Severity:    "high",
		Suggest:     "gob can execute arbitrary code; use encoding/json instead",
		FileFilter:  []string{".go"},
		Func:        detectGoGobDeserialization,
	},

	// =========================================================================
	// SSRF (Server-Side Request Forgery)
	// =========================================================================
	{
		Name:        "go-ssrf",
		Description: "HTTP request with user-controlled URL (potential SSRF)",
		Severity:    "high",
		Suggest:     "Validate URLs against an allowlist; never use user input directly in URLs",
		FileFilter:  []string{".go"},
		Func:        detectGoSSRF,
	},
	{
		Name:        "go-http-unvalidated-url",
		Description: "http.Get, http.Post, or http.Client Do with URL from user input",
		Severity:    "medium",
		Suggest:     "Validate and sanitize URL input; use URL allowlist",
		FileFilter:  []string{".go"},
		Func:        detectGoHTTPUnvalidatedURL,
	},

	// =========================================================================
	// TEMPLATE INJECTION
	// =========================================================================
	{
		Name:        "go-template-injection",
		Description: "html/template or text/template Execute with user-controlled template data",
		Severity:    "high",
		Suggest:     "Never pass user input as template source; use template execution only with trusted data",
		FileFilter:  []string{".go"},
		Func:        detectGoTemplateInjection,
	},
	{
		Name:        "go-template-raw-html",
		Description: "Template uses template.HTML with user input (potential XSS)",
		Severity:    "high",
		Suggest:     "Avoid template.HTML with user input; rely on auto-escaping",
		FileFilter:  []string{".go"},
		Func:        detectGoTemplateRawHTML,
	},

	// =========================================================================
	// REGEX DoS (ReDoS)
	// =========================================================================
	{
		Name:        "go-redos",
		Description: "Regular expression may be susceptible to ReDoS (catastrophic backtracking)",
		Severity:    "medium",
		Suggest:     "Use simple character classes; avoid nested quantifiers and alternations",
		FileFilter:  []string{".go"},
		Func:        detectGoReDoS,
	},
	{
		Name:        "go-regex-dynamic",
		Description: "regexp.Compile or regexp.MustCompile with user-controlled pattern",
		Severity:    "high",
		Suggest:     "Never compile user input as a regex; validate against allowlist",
		FileFilter:  []string{".go"},
		Func:        detectGoRegexDynamic,
	},

	// =========================================================================
	// HARDCODED CREDENTIALS & SECRETS
	// =========================================================================
	{
		Name:        "go-hardcoded-secret",
		Description: "Assignment to credential variable with string literal (not from env)",
		Severity:    "high",
		Suggest:     "Use os.Getenv or a secrets manager instead of hardcoded values",
		FileFilter:  []string{".go"},
		Func:        detectGoHardcodedSecret,
	},
	{
		Name:        "go-private-key",
		Description: "Private key or certificate data embedded as string literal",
		Severity:    "critical",
		Suggest:     "Load certificates/keys from files or environment variables",
		FileFilter:  []string{".go"},
		Func:        detectGoPrivateKey,
	},

	// =========================================================================
	// INSECURE CRYPTO
	// =========================================================================
	{
		Name:        "go-weak-crypto-md5",
		Description: "Use of MD5 hash function (cryptographically broken)",
		Severity:    "medium",
		Suggest:     "Use SHA-256 or SHA-3 instead; MD5 is broken for security purposes",
		FileFilter:  []string{".go"},
		Func:        detectGoWeakCrypto,
	},
	{
		Name:        "go-weak-crypto-sha1",
		Description: "Use of SHA-1 hash function (deprecated for security)",
		Severity:    "medium",
		Suggest:     "Use SHA-256 or SHA-3 instead; SHA-1 is deprecated",
		FileFilter:  []string{".go"},
		Func:        detectGoWeakCryptoSHA1,
	},
	{
		Name:        "go-insecure-random",
		Description: "Use of math/rand for security-sensitive random values",
		Severity:    "medium",
		Suggest:     "Use crypto/rand for security-sensitive randomness",
		FileFilter:  []string{".go"},
		Func:        detectGoInsecureRandom,
	},
	{
		Name:        "go-weak-cipher",
		Description: "Use of weak cipher (DES, RC4) or ECB mode",
		Severity:    "high",
		Suggest:     "Use AES-GCM or ChaCha20-Poly1305; avoid ECB mode",
		FileFilter:  []string{".go"},
		Func:        detectGoWeakCipher,
	},

	// =========================================================================
	// ERROR HANDLING
	// =========================================================================
	{
		Name:        "go-unchecked-error",
		Description: "Function returns error but return value is not checked",
		Severity:    "medium",
		Suggest:     "Always check error return values; handle or propagate errors",
		FileFilter:  []string{".go"},
		Func:        detectGoUncheckedError,
	},
	{
		Name:        "go-silent-panic",
		Description: "panic called in non-test production code",
		Severity:    "medium",
		Suggest:     "Return errors instead of panicking; reserve panic for unrecoverable states",
		FileFilter:  []string{".go"},
		Func:        detectGoSilentPanic,
	},

	// =========================================================================
	// LDAP INJECTION
	// =========================================================================
	{
		Name:        "go-ldap-injection",
		Description: "LDAP query constructed with string concatenation (potential LDAP injection)",
		Severity:    "high",
		Suggest:     "Use parameterized LDAP queries; validate and sanitize user input",
		FileFilter:  []string{".go"},
		Func:        detectGoLDAPInjection,
	},

	// =========================================================================
	// XML EXTERNAL ENTITY (XXE)
	// =========================================================================
	{
		Name:        "go-xxe",
		Description: "XML parsing without disabling external entity resolution (potential XXE)",
		Severity:    "high",
		Suggest:     "Set xmlparser.DisableExternalEntities or use a safe XML parser configuration",
		FileFilter:  []string{".go"},
		Func:        detectGoXXE,
	},

	// =========================================================================
	// YAML UNSAFE LOADING
	// =========================================================================
	{
		Name:        "go-yaml-unsafe",
		Description: "yaml.Unmarshal or yaml.NewDecoder.Decode with untrusted data (unsafe)",
		Severity:    "high",
		Suggest:     "Use yaml.Unmarshal with a known type; avoid yaml.TypeUnmarshaler on untrusted input",
		FileFilter:  []string{".go"},
		Func:        detectGoYAMLUnsafe,
	},

	// =========================================================================
	// TRUST BOUNDARY VIOLATION
	// =========================================================================
	{
		Name:        "go-trust-boundary",
		Description: "User-controlled data assigned to internal/global state without validation",
		Severity:    "medium",
		Suggest:     "Validate user input at trust boundaries; sanitize before storing",
		FileFilter:  []string{".go"},
		Func:        detectGoTrustBoundary,
	},

	// =========================================================================
	// HOSTNAME VERIFICATION BYPASS
	// =========================================================================
	{
		Name:        "go-tls-skip-verify",
		Description: "TLS client configured to skip certificate verification (insecure)",
		Severity:    "critical",
		Suggest:     "Always verify TLS certificates; use proper certificate pinning if needed",
		FileFilter:  []string{".go"},
		Func:        detectGoTLSSkipVerify,
	},
	{
		Name:        "go-insecure-tls",
		Description: "TLS config with InsecureSkipVerify set to true or MinVersion < TLS 1.2",
		Severity:    "high",
		Suggest:     "Set InsecureSkipVerify = false; use TLS 1.2 or higher",
		FileFilter:  []string{".go"},
		Func:        detectGoInsecureTLS,
	},
}

// =========================================================================
// SCANNING ENTRY POINTS
// =========================================================================

// ScanFileAST performs AST-based pattern scanning on a single Go file.
func ScanFileAST(path string, patterns []ASTPattern) ([]ASTFinding, error) {
	// Only scan Go files
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".go" {
		return nil, nil
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	var findings []ASTFinding
	for _, p := range patterns {
		// Skip patterns that don't apply to this file type
		if len(p.FileFilter) > 0 {
			skip := true
			for _, ext := range p.FileFilter {
				if ext == ext {
					skip = false
					break
				}
			}
			if skip {
				continue
			}
		}

		pFindings := p.Func(fset, file)
		for i := range pFindings {
			pFindings[i].File = path
			pFindings[i].Rule = p.Name
			// Column is already set by the individual pattern detectors
			// using fset.Position(node.Pos()).Column
		}
		findings = append(findings, pFindings...)
	}

	return findings, nil
}

// ScanDirAST recursively scans all Go files in a directory with AST patterns.
func ScanDirAST(dir string, patterns []ASTPattern) ([]ASTFinding, error) {
	var allFindings []ASTFinding

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors, continue walking
		}
		if info.IsDir() {
			// Skip common non-source directories
			switch info.Name() {
			case ".git", "node_modules", "vendor", ".terraform", "dist", "build", "__pycache__":
				return filepath.SkipDir
			}
			return nil
		}
		// Only scan Go files
		if filepath.Ext(path) != ".go" {
			return nil
		}
		findings, err := ScanFileAST(path, patterns)
		if err != nil {
			return nil // Skip parse errors, continue
		}
		allFindings = append(allFindings, findings...)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return allFindings, nil
}

// =========================================================================
// COMMAND INJECTION
// =========================================================================

func detectGoCommandInjection(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if !isExecCommand(call) {
			return true
		}

		for _, arg := range call.Args {
			if containsStringConcatOrFormat(arg) || containsUserInput(arg) {
				findings = append(findings, ASTFinding{
					Line:     fset.Position(arg.Pos()).Line,
					Message:  "Potential command injection: string concat or user input in exec.Command argument",
					Severity: "critical",
					Suggest:  "Use exec.Command with separate, pre-validated arguments instead of shell string concatenation",
				})
				break
			}
		}
		return true
	})

	return findings
}

func detectGoShellCommand(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if !isExecCommand(call) {
			return true
		}

		// Check if using shell -c pattern
		for i, arg := range call.Args {
			if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				val := unquote(lit.Value)
				if strings.Contains(val, " -c ") || strings.Contains(val, " -c\n") {
					if i > 0 {
						// The next arg is the command string - check if it contains user input
						if i+1 < len(call.Args) && containsUserInput(call.Args[i+1]) {
							findings = append(findings, ASTFinding{
								Line:     fset.Position(arg.Pos()).Line,
								Message:  "Shell invocation with user input: exec.Command(\"sh\", \"-c\", userInput)",
								Severity: "critical",
								Suggest:  "Avoid shell invocation; use exec.Command with separate arguments",
							})
						}
					}
				}
			}
		}
		return true
	})

	return findings
}

// =========================================================================
// SQL INJECTION
// =========================================================================

func detectGoSQLInjection(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if !isDatabaseMethod(call) {
			return true
		}

		for _, arg := range call.Args {
			if containsStringConcatOrFormat(arg) {
				findings = append(findings, ASTFinding{
					Line:     fset.Position(arg.Pos()).Line,
					Message:  "Potential SQL injection: string concatenation or formatting in query",
					Severity: "critical",
					Suggest:  "Use parameterized queries with ? or $1 placeholders; never concatenate user input",
				})
				break
			}
		}
		return true
	})

	return findings
}

func detectGoSQLTemplateQuery(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if !isDatabaseQueryMethod(call) {
			return true
		}

		for _, arg := range call.Args {
			if containsUserInput(arg) && !isSafeQueryArg(arg) {
				findings = append(findings, ASTFinding{
					Line:     fset.Position(arg.Pos()).Line,
					Message:  "Potential SQL injection: user input passed to query method without parameterization",
					Severity: "high",
					Suggest:  "Use parameterized queries; pass user values as parameters, not in query string",
				})
				break
			}
		}
		return true
	})

	return findings
}

// =========================================================================
// PATH TRAVERSAL
// =========================================================================

func detectGoPathTraversal(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if !isFileOperation(call) {
			return true
		}

		for _, arg := range call.Args {
			if containsUserInput(arg) {
				findings = append(findings, ASTFinding{
					Line:     fset.Position(arg.Pos()).Line,
					Message:  "Potential path traversal: user input in file operation path",
					Severity: "high",
					Suggest:  "Validate user input; use filepath.Clean, check against allowed directory list, resolve symlinks",
				})
				break
			}
		}
		return true
	})

	return findings
}

func detectGoSymlinkAttack(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	// Look for os.Open without O_NOFOLLOW
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isFileOpen(call) {
			// Check if O_NOFOLLOW is used
			if !hasOpenFlag(call, "O_NOFOLLOW") && !hasOpenFlag(call, "O_RDONLY") {
				// Heuristic: file open without O_NOFOLLOW on user-controlled path
				for _, arg := range call.Args {
					if containsUserInput(arg) {
						findings = append(findings, ASTFinding{
							Line:     fset.Position(call.Pos()).Line,
							Message:  "File operation may be susceptible to symlink attacks (no O_NOFOLLOW)",
							Severity: "medium",
							Suggest:  "Use O_NOFOLLOW flag to prevent following symlinks",
						})
						break
					}
				}
			}
		}
		return true
	})

	return findings
}

// =========================================================================
// DESERIALIZATION
// =========================================================================

func detectGoUnsafeDeserialization(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isUnmarshalCall(call) {
			for _, arg := range call.Args {
				if containsUserInput(arg) {
					findings = append(findings, ASTFinding{
						Line:     fset.Position(arg.Pos()).Line,
						Message:  "Potential unsafe deserialization: unmarshal with user-controlled data",
						Severity: "high",
						Suggest:  "Validate data source; prefer encoding/json over binary unmarshalers for untrusted data",
					})
					break
				}
			}
		}
		return true
	})

	return findings
}

func detectGoGobDeserialization(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isGobCall(call) {
			for _, arg := range call.Args {
				if containsUserInput(arg) {
					findings = append(findings, ASTFinding{
						Line:     fset.Position(arg.Pos()).Line,
						Message:  "Potential code execution via gob deserialization: gob can execute type methods during decode",
						Severity: "high",
						Suggest:  "Use encoding/json instead of encoding/gob for untrusted data; gob can execute arbitrary code",
					})
					break
				}
			}
		}
		return true
	})

	return findings
}

// =========================================================================
// SSRF
// =========================================================================

func detectGoSSRF(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isHTTPClientCall(call) {
			for _, arg := range call.Args {
				if containsUserInput(arg) {
					findings = append(findings, ASTFinding{
						Line:     fset.Position(arg.Pos()).Line,
						Message:  "Potential SSRF: user-controlled URL in HTTP request",
						Severity: "high",
						Suggest:  "Validate URL against allowlist; check scheme (http/https only), host, and path",
					})
					break
				}
			}
		}
		return true
	})

	return findings
}

func detectGoHTTPUnvalidatedURL(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	// Check for http.Get, http.Post, http.PostForm, client.Do etc. with URL variable
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isHTTPGetPost(call) {
			if len(call.Args) > 0 && containsUserInput(call.Args[0]) {
				findings = append(findings, ASTFinding{
					Line:     fset.Position(call.Pos()).Line,
					Message:  "HTTP request with user-controlled URL (potential SSRF)",
					Severity: "medium",
					Suggest:  "Validate URL before use; restrict to known-good hosts and schemes",
				})
			}
		}
		return true
	})

	return findings
}

// =========================================================================
// TEMPLATE INJECTION
// =========================================================================

func detectGoTemplateInjection(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isTemplateExecute(call) {
			for _, arg := range call.Args {
				if containsUserInput(arg) {
					findings = append(findings, ASTFinding{
						Line:     fset.Position(arg.Pos()).Line,
						Message:  "Potential template injection: user input passed to template execution",
						Severity: "high",
						Suggest:  "Never pass user input as template source; use template.Clone and validate input",
					})
					break
				}
			}
		}
		return true
	})

	return findings
}

func detectGoTemplateRawHTML(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check for template.HTML wrapping user input
		if isTemplateHTMLCall(call) {
			for _, arg := range call.Args {
				if containsUserInput(arg) {
					findings = append(findings, ASTFinding{
						Line:     fset.Position(arg.Pos()).Line,
						Message:  "Potential XSS: template.HTML with user-controlled data bypasses auto-escaping",
						Severity: "high",
						Suggest:  "Avoid template.HTML with user input; rely on template auto-escaping for user data",
					})
					break
				}
			}
		}
		return true
	})

	return findings
}

// =========================================================================
// REGEX DOS
// =========================================================================

var dangerousRegexPatterns = regexp.MustCompile(`(\.\*|\.\+|\.\?|.\*\*|.\+\+).*(\.\*|\.\+|\.\?|.\*\*|.\+\+)`)

func detectGoReDoS(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isRegexCompile(call) {
			for _, arg := range call.Args {
				if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
					pattern := unquote(lit.Value)
					if dangerousRegexPatterns.MatchString(pattern) {
						findings = append(findings, ASTFinding{
							Line:     fset.Position(arg.Pos()).Line,
							Message:  fmt.Sprintf("Potential ReDoS: regex pattern contains nested quantifiers: %s", truncate(pattern, 40)),
							Severity: "medium",
							Suggest:  "Simplify regex; avoid nested quantifiers; use possessive quantifiers or atomic groups",
						})
					}
				}
			}
		}
		return true
	})

	return findings
}

func detectGoRegexDynamic(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isRegexCompile(call) {
			for _, arg := range call.Args {
				if containsUserInput(arg) {
					findings = append(findings, ASTFinding{
						Line:     fset.Position(arg.Pos()).Line,
						Message:  "Dynamic regex compilation with user-controlled pattern (ReDoS risk)",
						Severity: "high",
						Suggest:  "Never compile user input as regex; validate against a known pattern allowlist",
					})
					break
				}
			}
		}
		return true
	})

	return findings
}

// =========================================================================
// HARDCODED SECRETS
// =========================================================================

var credentialPattern = regexp.MustCompile(`(?i)^(password|passwd|pwd|secret|api_?key|token|private_?key|auth_?token|access_?token|client_?secret|encryption_?key|ssh_?key)$`)

func detectGoHardcodedSecret(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		for i, lhs := range assign.Lhs {
			if ident, ok := lhs.(*ast.Ident); ok {
				if credentialPattern.MatchString(ident.Name) {
					if i < len(assign.Rhs) {
						rhs := assign.Rhs[i]
						if isNonEmptyString(rhs) && !containsEnvVar(rhs) && !isFunctionCall(rhs) {
							findings = append(findings, ASTFinding{
								Line:     fset.Position(assign.Pos()).Line,
								Message:  fmt.Sprintf("Hardcoded credential: variable '%s' assigned string literal", ident.Name),
								Severity: "high",
								Suggest:  "Use os.Getenv or a secrets manager (Vault, AWS Secrets Manager, etc.)",
							})
						}
					}
				}
			}
		}
		return true
	})

	return findings
}

func detectGoPrivateKey(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	// Look for string literals containing PEM-encoded private key markers
	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}

		val := unquote(lit.Value)
		if (strings.Contains(val, "-----BEGIN PRIVATE KEY-----") ||
			strings.Contains(val, "-----BEGIN RSA PRIVATE KEY-----") ||
			strings.Contains(val, "-----BEGIN EC PRIVATE KEY-----") ||
			strings.Contains(val, "-----BEGIN OPENSSH PRIVATE KEY-----")) &&
			!containsEnvVarFromLit(lit) {
			findings = append(findings, ASTFinding{
				Line:     fset.Position(lit.Pos()).Line,
				Message:  "Private key or certificate embedded as string literal",
				Severity: "critical",
				Suggest:  "Load private keys from files or environment variables; never embed in source",
			})
		}
		return true
	})

	return findings
}

// =========================================================================
// INSECURE CRYPTO
// =========================================================================

func detectGoWeakCrypto(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isMD5Hash(call) {
			findings = append(findings, ASTFinding{
				Line:     fset.Position(call.Pos()).Line,
				Message:  "Use of MD5 hash function (cryptographically broken for security)",
				Severity: "medium",
				Suggest:  "Use crypto/sha256 or crypto/sha3 instead of crypto/md5",
			})
		}
		return true
	})

	return findings
}

func detectGoWeakCryptoSHA1(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isSHA1Hash(call) {
			findings = append(findings, ASTFinding{
				Line:     fset.Position(call.Pos()).Line,
				Message:  "Use of SHA-1 hash function (deprecated for security purposes)",
				Severity: "medium",
				Suggest:  "Use crypto/sha256 or crypto/sha3 instead of crypto/sha1",
			})
		}
		return true
	})

	return findings
}

func detectGoInsecureRandom(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isMathRandCall(call) {
			findings = append(findings, ASTFinding{
				Line:     fset.Position(call.Pos()).Line,
				Message:  "Use of math/rand for security-sensitive randomness (predictable)",
				Severity: "medium",
				Suggest:  "Use crypto/rand for security-sensitive random values",
			})
		}
		return true
	})

	return findings
}

func detectGoWeakCipher(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isWeakCipher(call) {
			findings = append(findings, ASTFinding{
				Line:     fset.Position(call.Pos()).Line,
				Message:  "Use of weak cipher (DES, 3DES, RC4) or ECB mode",
				Severity: "high",
				Suggest:  "Use AES-GCM or ChaCha20-Poly1305; avoid ECB mode for block ciphers",
			})
		}
		return true
	})

	return findings
}

// =========================================================================
// ERROR HANDLING
// =========================================================================

func detectGoUncheckedError(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	// Skip test files for this pattern
	// Note: we'd need file path here; for simplicity, caller should filter

	ast.Inspect(file, func(n ast.Node) bool {
		exprStmt, ok := n.(*ast.ExprStmt)
		if !ok {
			return true
		}

		call, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isErrorReturningCall(call) {
			findings = append(findings, ASTFinding{
				Line:     fset.Position(exprStmt.Pos()).Line,
				Message:  "Error return value not checked: function returns error but caller ignores it",
				Severity: "medium",
				Suggest:  "Check error return value; handle or propagate errors explicitly",
			})
		}
		return true
	})

	return findings
}

func detectGoSilentPanic(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isPanicCall(call) && !isInTestFile(file) {
			findings = append(findings, ASTFinding{
				Line:     fset.Position(call.Pos()).Line,
				Message:  "panic called in production code",
				Severity: "medium",
				Suggest:  "Return errors instead of panicking; reserve panic for unrecoverable states",
			})
		}
		return true
	})

	return findings
}

// =========================================================================
// LDAP INJECTION
// =========================================================================

func detectGoLDAPInjection(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isLDAPMethod(call) {
			for _, arg := range call.Args {
				if containsStringConcatOrFormat(arg) || containsUserInput(arg) {
					findings = append(findings, ASTFinding{
						Line:     fset.Position(arg.Pos()).Line,
						Message:  "Potential LDAP injection: string concatenation in LDAP query",
						Severity: "high",
						Suggest:  "Use parameterized LDAP queries; validate and sanitize user input",
					})
					break
				}
			}
		}
		return true
	})

	return findings
}

// =========================================================================
// XXE
// =========================================================================

func detectGoXXE(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isXMLParserCall(call) && !hasXXEPrevention(call) {
			findings = append(findings, ASTFinding{
				Line:     fset.Position(call.Pos()).Line,
				Message:  "XML parsing without external entity protection (potential XXE)",
				Severity: "high",
				Suggest:  "Disable external entities: xmlparser.DisableExternalEntities(true) or use safe defaults",
			})
		}
		return true
	})

	return findings
}

// =========================================================================
// YAML UNSAFE
// =========================================================================

func detectGoYAMLUnsafe(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isYAMLUnmarshal(call) {
			for _, arg := range call.Args {
				if containsUserInput(arg) {
					findings = append(findings, ASTFinding{
						Line:     fset.Position(arg.Pos()).Line,
						Message:  "yaml.Unmarshal with untrusted data (can trigger arbitrary code in some configurations)",
						Severity: "high",
						Suggest:  "Validate input; be aware that yaml can invoke UnmarshalYAML methods",
					})
					break
				}
			}
		}
		return true
	})

	return findings
}

// =========================================================================
// TRUST BOUNDARY
// =========================================================================

func detectGoTrustBoundary(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	// Heuristic: assignment to package-level or global variable from function param
	ast.Inspect(file, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		for _, lhs := range assign.Lhs {
			// Check if assigning to an exported package-level identifier
			if ident, ok := lhs.(*ast.Ident); ok && ast.IsExported(ident.Name) {
				for _, rhs := range assign.Rhs {
					if call, ok := rhs.(*ast.CallExpr); ok {
						// Check if the function param comes from user input
						if isUserInputFunction(call) {
							findings = append(findings, ASTFinding{
								Line:     fset.Position(assign.Pos()).Line,
								Message:  fmt.Sprintf("Trust boundary violation: exported var '%s' assigned from user input", ident.Name),
								Severity: "medium",
								Suggest:  "Validate user input at trust boundaries before assigning to shared state",
							})
						}
					}
				}
			}
		}
		return true
	})

	return findings
}

// =========================================================================
// TLS / SSL
// =========================================================================

func detectGoTLSSkipVerify(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isTLSSkipVerify(call) {
			findings = append(findings, ASTFinding{
				Line:     fset.Position(call.Pos()).Line,
				Message:  "TLS certificate verification disabled (InsecureSkipVerify = true)",
				Severity: "critical",
				Suggest:  "Set InsecureSkipVerify = false; use proper TLS verification",
			})
		}
		return true
	})

	return findings
}

func detectGoInsecureTLS(fset *token.FileSet, file *ast.File) []ASTFinding {
	var findings []ASTFinding

	ast.Inspect(file, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		for i, lhs := range assign.Lhs {
			if sel, ok := lhs.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "InsecureSkipVerify" || sel.Sel.Name == "MinVersion" {
					if i < len(assign.Rhs) {
						rhs := assign.Rhs[i]
						if sel.Sel.Name == "InsecureSkipVerify" && isTrue(rhs) {
							findings = append(findings, ASTFinding{
								Line:     fset.Position(assign.Pos()).Line,
								Message:  "TLS InsecureSkipVerify set to true (disables certificate verification)",
								Severity: "critical",
								Suggest:  "Set InsecureSkipVerify = false for production TLS connections",
							})
						}
					}
				}
			}
		}
		return true
	})

	return findings
}

// =========================================================================
// HELPER FUNCTIONS - CALL IDENTIFICATION
// =========================================================================

func isExecCommand(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return ident.Name == "exec" && sel.Sel.Name == "Command"
		}
	}
	return false
}

func isDatabaseMethod(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		name := sel.Sel.Name
		return name == "Query" || name == "Exec" || name == "QueryRow" ||
			name == "Execute" || name == "ExecContext" || name == "QueryContext" ||
			name == "QueryRowContext" || name == "rawQuery" || name == "Raw" ||
			name == "Scan"
	}
	if ident, ok := call.Fun.(*ast.Ident); ok {
		name := ident.Name
		return name == "Query" || name == "Exec" || name == "QueryRow" ||
			name == "Execute" || name == "rawQuery"
	}
	return false
}

func isDatabaseQueryMethod(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		name := sel.Sel.Name
		return name == "Query" || name == "QueryRow" || name == "rawQuery" ||
			name == "Select" || name == "Get" || name == "Find"
	}
	return false
}

func isFileOperation(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		name := sel.Sel.Name
		return name == "Open" || name == "ReadFile" || name == "WriteFile" ||
			name == "Create" || name == "Stat" || name == "Rename" ||
			name == "Remove" || name == "MkdirAll" || name == "TempFile" ||
			name == "TempDir" || name == "ReadDir"
	}
	return false
}

func isFileOpen(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		return sel.Sel.Name == "Open" || sel.Sel.Name == "OpenFile"
	}
	return false
}

func hasOpenFlag(call *ast.CallExpr, flagName string) bool {
	for _, arg := range call.Args {
		if sel, ok := arg.(*ast.SelectorExpr); ok {
			if sel.Sel.Name == flagName {
				return true
			}
		}
	}
	return false
}

func isUnmarshalCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		name := sel.Sel.Name
		return name == "Unmarshal" || name == "Decode" || name == "NewDecoder" ||
			name == "Deserialize"
	}
	if ident, ok := call.Fun.(*ast.Ident); ok {
		name := ident.Name
		return name == "Unmarshal" || name == "Decode"
	}
	return false
}

func isGobCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return ident.Name == "gob" && (sel.Sel.Name == "NewDecoder" || sel.Sel.Name == "Decode")
		}
	}
	return false
}

func isHTTPClientCall(call *ast.CallExpr) bool {
	if ident, ok := call.Fun.(*ast.Ident); ok {
		return ident.Name == "Get" || ident.Name == "Post" || ident.Name == "PostForm" ||
			ident.Name == "Head" || ident.Name == "Do"
	}
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return (ident.Name == "client" || ident.Name == "Client") &&
				(sel.Sel.Name == "Get" || sel.Sel.Name == "Post" || sel.Sel.Name == "Do" ||
					sel.Sel.Name == "Head")
		}
	}
	return false
}

func isHTTPGetPost(call *ast.CallExpr) bool {
	if ident, ok := call.Fun.(*ast.Ident); ok {
		return ident.Name == "Get" || ident.Name == "Post" || ident.Name == "PostForm" ||
			ident.Name == "Head"
	}
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		return sel.Sel.Name == "Get" || sel.Sel.Name == "Post" || sel.Sel.Name == "Do"
	}
	return false
}

func isTemplateExecute(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		return sel.Sel.Name == "Execute" || sel.Sel.Name == "ExecuteTemplate"
	}
	return false
}

func isTemplateHTMLCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return ident.Name == "template" && sel.Sel.Name == "HTML"
		}
	}
	return false
}

func isRegexCompile(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return ident.Name == "regexp" &&
				(sel.Sel.Name == "Compile" || sel.Sel.Name == "MustCompile")
		}
	}
	return false
}

func isErrorReturningCall(call *ast.CallExpr) bool {
	if ident, ok := call.Fun.(*ast.Ident); ok {
		name := ident.Name
		// Common functions that return errors
		return name == "Read" || name == "Write" || name == "Close" ||
			name == "Scan" || name == "Next" || name == "Decode" ||
			name == "Open" || name == "Create" || name == "Stat"
	}
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		return sel.Sel.Name == "Do" || sel.Sel.Name == "Scan" ||
			sel.Sel.Name == "Next" || sel.Sel.Name == "Decode"
	}
	return false
}

func isPanicCall(call *ast.CallExpr) bool {
	if ident, ok := call.Fun.(*ast.Ident); ok {
		return ident.Name == "panic"
	}
	return false
}

func isInTestFile(file *ast.File) bool {
	// This would need the filename passed in; heuristic check via package name
	for _, decl := range file.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok {
			if gen.Tok == token.TYPE && gen.Specs != nil {
				for _, spec := range gen.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						if strings.HasSuffix(ts.Name.Name, "_test") {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func isLDAPMethod(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		name := sel.Sel.Name
		return name == "Search" || name == "SearchFunc" || name == "Bind" ||
			name == "Modify" || name == "Add" || name == "Delete"
	}
	return false
}

func isXMLParserCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return (ident.Name == "xml" || ident.Name == "xmlparser" || ident.Name == "encoding/xml") &&
				(sel.Sel.Name == "NewParser" || sel.Sel.Name == "Parse" || sel.Sel.Name == "Decode")
		}
	}
	return false
}

func hasXXEPrevention(call *ast.CallExpr) bool {
	// Check if DisableExternalEntities or similar was called earlier in the same file
	// This is a simplification; real implementation would track state
	return false
}

func isYAMLUnmarshal(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return ident.Name == "yaml" && sel.Sel.Name == "Unmarshal"
		}
	}
	return false
}

func isUserInputFunction(call *ast.CallExpr) bool {
	if ident, ok := call.Fun.(*ast.Ident); ok {
		name := ident.Name
		return name == "ReadAll" || name == "ioutil.ReadAll" ||
			strings.HasPrefix(name, "Read") || name == "ParseForm" ||
			name == "ParseMultipartForm"
	}
	return false
}

func isTLSSkipVerify(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return ident.Name == "tls" && sel.Sel.Name == "Config"
		}
	}
	return false
}

func isMD5Hash(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return ident.Name == "md5" && (sel.Sel.Name == "New" || sel.Sel.Name == "Sum")
		}
	}
	return false
}

func isSHA1Hash(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return ident.Name == "sha1" && (sel.Sel.Name == "New" || sel.Sel.Name == "Sum")
		}
	}
	return false
}

func isMathRandCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			// Check for any math/rand function calls
			if ident.Name == "rand" {
				method := sel.Sel.Name
				// Intn, Int31n, Int63n, Float64, etc. are all insecure for security
				return method == "New" || method == "Int" || method == "Intn" ||
					method == "Int31" || method == "Int31n" || method == "Int63" ||
					method == "Int63n" || method == "Float64" || method == "Uint64" ||
					method == "Shuffle"
			}
		}
	}
	return false
}

func isWeakCipher(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		name := sel.Sel.Name
		return name == "NewCipher" || name == "NewCBCEncrypter" || name == "NewCFBEncrypter"
	}
	return false
}

// =========================================================================
// HELPER FUNCTIONS - PATTERN MATCHING
// =========================================================================

func containsStringConcatOrFormat(n ast.Node) bool {
	var found bool
	ast.Inspect(n, func(node ast.Node) bool {
		if found {
			return false
		}
		if bin, ok := node.(*ast.BinaryExpr); ok {
			if bin.Op == token.ADD {
				found = true
				return false
			}
		}
		// Also check for fmt.Sprintf, fmt.Fprintf, strings.Join, etc.
		if call, ok := node.(*ast.CallExpr); ok {
			if ident, ok := call.Fun.(*ast.Ident); ok {
				if ident.Name == "Sprintf" || ident.Name == "Sprint" || ident.Name == "Errorf" {
					found = true
					return false
				}
			}
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok {
					if ident.Name == "fmt" && (sel.Sel.Name == "Sprintf" || sel.Sel.Name == "Sprint" || sel.Sel.Name == "Errorf") {
						found = true
						return false
					}
					if ident.Name == "strings" && sel.Sel.Name == "Join" {
						found = true
						return false
					}
				}
			}
		}
		return true
	})
	return found
}

func containsUserInput(n ast.Node) bool {
	var found bool
	ast.Inspect(n, func(node ast.Node) bool {
		if found {
			return false
		}
		if ident, ok := node.(*ast.Ident); ok {
			name := strings.ToLower(ident.Name)
			// Common user input variable names - expanded to catch more patterns
			if name == "req" || name == "request" || name == "body" ||
				name == "input" || name == "params" || name == "query" ||
				name == "form" || name == "ctx" || name == "r" ||
				name == "w" || name == "response" || name == "data" ||
				name == "args" || name == "argv" || name == "env" ||
				name == "filename" || name == "filepath" || name == "path" ||
				name == "name" || name == "url" || name == "uri" ||
				strings.HasPrefix(name, "user") || strings.HasPrefix(name, "http") ||
				strings.HasPrefix(name, "post") || strings.HasPrefix(name, "get") ||
				strings.HasPrefix(name, "cookie") || strings.HasPrefix(name, "header") ||
				strings.HasPrefix(name, "file") || strings.HasPrefix(name, "path") {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

func isSafeQueryArg(n ast.Node) bool {
	// Check if argument is a parameterized placeholder or a constant
	if ident, ok := n.(*ast.Ident); ok {
		name := strings.ToLower(ident.Name)
		return name == "_" || name == "args" || name == "params"
	}
	return false
}

func containsEnvVar(n ast.Node) bool {
	if call, ok := n.(*ast.CallExpr); ok {
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok {
				return ident.Name == "os" && sel.Sel.Name == "Getenv"
			}
		}
	}
	return false
}

func containsEnvVarFromLit(lit *ast.BasicLit) bool {
	// Check if the file contains os.Getenv near this literal
	return false // Simplified; real implementation would search surrounding context
}

func isNonEmptyString(n ast.Node) bool {
	if lit, ok := n.(*ast.BasicLit); ok {
		if lit.Kind == token.STRING {
			val := unquote(lit.Value)
			return len(val) > 0
		}
	}
	return false
}

func isFunctionCall(n ast.Node) bool {
	_, ok := n.(*ast.CallExpr)
	return ok
}

func isTrue(n ast.Node) bool {
	ident, ok := n.(*ast.Ident)
	return ok && ident.Name == "true"
}

func isNonEmptyStringOrIdent(n ast.Node) bool {
	if isNonEmptyString(n) {
		return true
	}
	if ident, ok := n.(*ast.Ident); ok {
		return ident.Name != "_" && ident.Name != "" // Allow underscore (blank identifier)
	}
	return false
}

// =========================================================================
// UTILITIES
// =========================================================================

func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '`' && s[len(s)-1] == '`') {
			// Simple unquote - doesn't handle escape sequences
			return s[1 : len(s)-1]
		}
		if s[0] == '\'' && s[len(s)-1] == '\'' {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// StringContains reports whether substr is found in the AST node's string value.
func stringContains(n ast.Node, substr string) bool {
	var found bool
	ast.Inspect(n, func(node ast.Node) bool {
		if found {
			return false
		}
		if lit, ok := node.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			if strings.Contains(unquote(lit.Value), substr) {
				found = true
				return false
			}
		}
		return true
	})
	return found
}
