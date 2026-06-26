package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/aliasfoxkde/Atheon/core"
)

// version is the server version, set at build time via:
//
//	-ldflags "-X main.version=1.2.3"
//
// Defaults to "dev" so `go run ./cmd/mcp` is usable without a build script.
var version = "dev"

type request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type response struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      any       `json:"id"`
	Result  any       `json:"result,omitempty"`
	Error   *rpcError `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// rateLimitCode is the JSON-RPC error code returned when a request is
// denied by the rate limiter. JSON-RPC reserves -32000..-32099 for
// implementation-defined server errors; -32600 is "Invalid Request",
// which is the wrong code for a throttling response.
const rateLimitCode = -32000

// rateLimiter implements a simple token bucket rate limiter.
// Uses stdlib only to avoid external dependencies.
type rateLimiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTime time.Time
}

// newRateLimiter creates a rate limiter allowing maxTokens per second, up to burst.
func newRateLimiter(tokensPerSecond, burst float64) *rateLimiter {
	return &rateLimiter{
		tokens:   burst,
		max:      burst,
		rate:     tokensPerSecond,
		lastTime: time.Now(),
	}
}

// Allow checks if a request is permitted under the rate limit.
// Returns true if allowed, false if rate limited.
func (rl *rateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastTime).Seconds()
	rl.lastTime = now

	// Add tokens based on elapsed time
	rl.tokens += elapsed * rl.rate
	if rl.tokens > rl.max {
		rl.tokens = rl.max
	}

	if rl.tokens < 1 {
		return false
	}
	rl.tokens--
	return true
}

// mcpRateLimiter is the global rate limiter for MCP requests.
// Allows 10 requests per second with a burst of 20.
var mcpRateLimiter = newRateLimiter(10, 20)

// JSON-RPC method names handled by the MCP server. Extracted as
// constants so goconst can verify they're not duplicated and so
// readers can see the protocol surface in one place.
const (
	methodInitialize = "initialize"
	methodToolsList  = "tools/list"
	methodToolsCall  = "tools/call"
)

func main() {
	configureLogging()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := run(ctx, os.Stdin, os.Stdout)
	cancel() // explicit: os.Exit skips deferred cancel
	os.Exit(code)
}

// configureLogging mirrors cmd/atheon's setup so MCP server logs are
// configurable via the same env vars (ATHEON_LOG_FORMAT, ATHEON_LOG_LEVEL).
// Without this, slog's default text handler is used and downstream
// aggregators have to parse key=value pairs from a non-deterministic
// format.
func configureLogging() {
	var level slog.Level
	switch strings.ToLower(os.Getenv("ATHEON_LOG_LEVEL")) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if strings.EqualFold(os.Getenv("ATHEON_LOG_FORMAT"), "json") {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, opts)
	}
	slog.SetDefault(slog.New(handler))
}

// run executes the JSON-RPC loop reading from r and writing to w, returning
// the exit code. Separated from main() so tests can call it without os.Exit
// terminating the test process.
//
// The context is forwarded into the core scan helpers so a SIGTERM
// received mid-scan aborts cleanly.
func run(ctx context.Context, r io.Reader, w io.Writer) int {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	enc := json.NewEncoder(w)

	for sc.Scan() {
		var req request
		if err := json.Unmarshal(sc.Bytes(), &req); err != nil {
			// JSON-RPC: malformed requests have no ID, so we cannot send an
			// error response. Log to stderr for debuggability.
			fmt.Fprintf(os.Stderr, "atheon-mcp: malformed JSON-RPC request: %v\n", err)
			continue
		}
		// Per-request structured log so MCP traffic is observable in ELK /
		// Loki / Datadog. Gated at Debug so the default Info level stays
		// quiet for the common case. Use fmt.Sprintf for ID since the
		// field is `any` and JSON-encoding a nil ID emits "null" which is
		// technically correct but harder to grep than "<notif>".
		idStr := "<notif>"
		if req.ID != nil {
			idStr = fmt.Sprintf("%v", req.ID)
		}
		slog.Debug("mcp request", "method", req.Method, "id", idStr)
		if req.Method == "initialized" {
			continue
		}

		var result any
		var rerr *rpcError

		switch req.Method {
		case methodInitialize:
			result = map[string]any{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]any{"tools": map[string]any{}},
				"serverInfo":      map[string]any{"name": "atheon", "version": version},
			}
		case methodToolsList:
			result = map[string]any{"tools": toolList()}
		case methodToolsCall:
			result, rerr = handleCall(ctx, req.Params)
		default:
			rerr = &rpcError{Code: -32601, Message: "method not found"}
		}

		resp := response{JSONRPC: "2.0", ID: req.ID}
		if rerr != nil {
			resp.Error = rerr
		} else {
			resp.Result = result
		}
		_ = enc.Encode(resp)
	}
	return 0
}

// toolList returns the MCP tool registry. The schema helper wraps a Go
// property bag into the JSON Schema shape MCP expects.
func toolList() []map[string]any {
	schema := func(required []string, props map[string]any) map[string]any {
		return map[string]any{"type": "object", "properties": props, "required": required}
	}
	str := map[string]any{"type": "string"}
	cats := map[string]any{"type": "array", "items": str, "description": "categories to scan (omit for all)"}

	return []map[string]any{
		{
			"name":        "scan_string",
			"description": "Scan a string for pattern matches",
			"inputSchema": schema([]string{"content"}, map[string]any{
				"content":    map[string]any{"type": "string"},
				"source":     str,
				"categories": cats,
			}),
		},
		{
			"name":        "scan_file",
			"description": "Scan a file for pattern matches",
			"inputSchema": schema([]string{"path"}, map[string]any{
				"path":       map[string]any{"type": "string"},
				"categories": cats,
			}),
		},
		{
			"name":        "scan_dir",
			"description": "Scan a directory for pattern matches",
			"inputSchema": schema([]string{"path"}, map[string]any{
				"path":       map[string]any{"type": "string"},
				"categories": cats,
			}),
		},
		{
			"name":        "scan_env",
			"description": "Scan process environment variables for pattern matches",
			"inputSchema": schema([]string{}, map[string]any{
				"categories": cats,
			}),
		},
		{
			"name":        "list_patterns",
			"description": "List all loaded patterns (name, category, enabled)",
			"inputSchema": schema([]string{}, map[string]any{
				"category": map[string]any{
					"type":        "string",
					"description": "filter to a single category (omit for all)",
				},
			}),
		},
		{
			"name":        "list_categories",
			"description": "List all pattern categories available in the bundle",
			"inputSchema": schema([]string{}, map[string]any{}),
		},
		{
			"name":        "update_bundle",
			"description": "Download the latest pattern bundle from the configured URL",
			"inputSchema": schema([]string{}, map[string]any{}),
		},
	}
}

// toolHandler is the signature every per-tool dispatcher implements.
// Extracted so handleCall stays under the lint funlen limit and each
// tool's parse-and-execute logic is independently testable.
type toolHandler func(ctx context.Context, args json.RawMessage) (any, *rpcError)

// toolHandlers maps each tool name to its handler. Lookups via map keep
// handleCall flat instead of an ever-growing switch.
var toolHandlers = map[string]toolHandler{
	"scan_string":     handleScanString,
	"scan_file":       handleScanFile,
	"scan_dir":        handleScanDir,
	"scan_env":        handleScanEnv,
	"list_patterns":   handleListPatterns,
	"list_categories": handleListCategories,
	"update_bundle":   handleUpdateBundle,
}

// handleCall parses the JSON-RPC params envelope, looks up the tool
// handler, and dispatches. Rate-limit and envelope validation live here
// so every tool inherits them.
func handleCall(ctx context.Context, params json.RawMessage) (any, *rpcError) {
	if !mcpRateLimiter.Allow() {
		slog.Warn("rate limit exceeded for MCP request")
		return nil, &rpcError{Code: rateLimitCode, Message: "rate limit exceeded"}
	}
	var p struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, &rpcError{Code: -32602, Message: "invalid params"}
	}
	h, ok := toolHandlers[p.Name]
	if !ok {
		return nil, &rpcError{Code: -32601, Message: "unknown tool: " + p.Name}
	}
	return h(ctx, p.Arguments)
}

// invalidParams is a small helper so every tool handler returns the
// same JSON-RPC error shape on argument parse failure.
func invalidParams(err error) *rpcError {
	return &rpcError{Code: -32602, Message: "invalid params: " + err.Error()}
}

func handleScanString(ctx context.Context, raw json.RawMessage) (any, *rpcError) {
	var args struct {
		Content    string   `json:"content"`
		Source     string   `json:"source"`
		Categories []string `json:"categories"`
	}
	if err := json.Unmarshal(raw, &args); err != nil {
		return nil, invalidParams(err)
	}
	if args.Source == "" {
		args.Source = "stdin"
	}
	core.SetActiveCategories(args.Categories)
	return textResult(core.ScanString(ctx, args.Content, args.Source)), nil
}

func handleScanFile(ctx context.Context, raw json.RawMessage) (any, *rpcError) {
	var args struct {
		Path       string   `json:"path"`
		Categories []string `json:"categories"`
	}
	if err := json.Unmarshal(raw, &args); err != nil {
		return nil, invalidParams(err)
	}
	core.SetActiveCategories(args.Categories)
	findings, _, err := core.ScanFile(ctx, args.Path)
	if err != nil {
		return nil, &rpcError{Code: -32603, Message: err.Error()}
	}
	return textResult(findings), nil
}

func handleScanDir(ctx context.Context, raw json.RawMessage) (any, *rpcError) {
	var args struct {
		Path       string   `json:"path"`
		Categories []string `json:"categories"`
	}
	if err := json.Unmarshal(raw, &args); err != nil {
		return nil, invalidParams(err)
	}
	core.SetActiveCategories(args.Categories)
	// MCP defaults to the safe symlink policy. Agents invoking scan_dir
	// are typically operating on untrusted trees (third-party repos,
	// generated code, scratch dirs), and a symlink escape would let a
	// crafted repo leak /etc/passwd or ~/.aws/credentials into the
	// findings without the operator ever noticing. The CLI keeps the
	// historical follow-symlinks behaviour behind an opt-in flag.
	findings, _, err := core.ScanDir(ctx, args.Path, core.ScanOpts{NoFollowSymlinks: true})
	if err != nil {
		return nil, &rpcError{Code: -32603, Message: err.Error()}
	}
	return textResult(findings), nil
}

func handleScanEnv(ctx context.Context, raw json.RawMessage) (any, *rpcError) {
	var args struct {
		Categories []string `json:"categories"`
	}
	if err := json.Unmarshal(raw, &args); err != nil {
		return nil, invalidParams(err)
	}
	core.SetActiveCategories(args.Categories)
	return textResult(core.ScanEnv(ctx)), nil
}

func handleListPatterns(_ context.Context, raw json.RawMessage) (any, *rpcError) {
	var args struct {
		Category string `json:"category"`
	}
	if err := json.Unmarshal(raw, &args); err != nil {
		return nil, invalidParams(err)
	}
	return patternsResult(core.All(), args.Category), nil
}

func handleListCategories(_ context.Context, _ json.RawMessage) (any, *rpcError) {
	return categoriesResult(core.Categories()), nil
}

func handleUpdateBundle(ctx context.Context, _ json.RawMessage) (any, *rpcError) {
	if err := core.DownloadBundle(ctx); err != nil {
		return nil, &rpcError{Code: -32603, Message: err.Error()}
	}
	return map[string]any{
		"content": []map[string]any{{
			"type": "text",
			"text": "bundle updated successfully",
		}},
	}, nil
}

func textResult(findings []core.Finding) map[string]any {
	var sb strings.Builder
	if len(findings) == 0 {
		sb.WriteString("no findings")
	} else {
		for _, f := range findings {
			fmt.Fprintf(&sb, "%s  %s:%d\n", f.Pattern, f.File, f.Line)
		}
		fmt.Fprintf(&sb, "\n%d finding(s)", len(findings))
	}
	return map[string]any{
		"content": []map[string]any{{"type": "text", "text": sb.String()}},
	}
}

// patternsResult renders the pattern list as a markdown table inside the
// MCP text-content wrapper. If category is non-empty, only patterns whose
// Category() matches are returned.
func patternsResult(patterns []core.Pattern, category string) map[string]any {
	var sb strings.Builder
	fmt.Fprintf(&sb, "| Name | Category | Enabled |\n")
	fmt.Fprintf(&sb, "|------|----------|---------|\n")
	count := 0
	for _, p := range patterns {
		if category != "" && p.Category() != category {
			continue
		}
		enabled := "no"
		if p.Enabled() {
			enabled = "yes"
		}
		fmt.Fprintf(&sb, "| %s | %s | %s |\n", p.Name(), p.Category(), enabled)
		count++
	}
	if count == 0 {
		if category != "" {
			fmt.Fprintf(&sb, "\n(no patterns in category %q)", category)
		} else {
			sb.WriteString("\n(no patterns loaded)")
		}
	} else {
		fmt.Fprintf(&sb, "\n%d pattern(s)", count)
	}
	return map[string]any{
		"content": []map[string]any{{"type": "text", "text": sb.String()}},
	}
}

// categoriesResult renders the category list as a simple comma-separated
// text block, since categories are short labels.
func categoriesResult(cats []string) map[string]any {
	var sb strings.Builder
	if len(cats) == 0 {
		sb.WriteString("(no categories loaded)")
	} else {
		fmt.Fprintf(&sb, "%s\n\n%d categories", strings.Join(cats, ", "), len(cats))
	}
	return map[string]any{
		"content": []map[string]any{{"type": "text", "text": sb.String()}},
	}
}
