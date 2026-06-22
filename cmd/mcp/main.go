package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/aliasfoxkde/Atheon/core"
)

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

// JSON-RPC method names handled by the MCP server. Extracted as
// constants so goconst can verify they're not duplicated and so
// readers can see the protocol surface in one place.
const (
	methodInitialize = "initialize"
	methodToolsList  = "tools/list"
	methodToolsCall  = "tools/call"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := run(ctx, os.Stdin, os.Stdout)
	cancel() // explicit: os.Exit skips deferred cancel
	os.Exit(code)
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
			continue
		}
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
				"serverInfo":      map[string]any{"name": "atheon", "version": "1.0.0"},
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
	}
}

func handleCall(ctx context.Context, params json.RawMessage) (any, *rpcError) {
	var p struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, &rpcError{Code: -32602, Message: "invalid params"}
	}

	switch p.Name {
	case "scan_string":
		var args struct {
			Content    string   `json:"content"`
			Source     string   `json:"source"`
			Categories []string `json:"categories"`
		}
		if err := json.Unmarshal(p.Arguments, &args); err != nil {
			return nil, &rpcError{Code: -32602, Message: "invalid params"}
		}
		if args.Source == "" {
			args.Source = "stdin"
		}
		core.SetActiveCategories(args.Categories)
		return textResult(core.ScanString(ctx, args.Content, args.Source)), nil

	case "scan_file":
		var args struct {
			Path       string   `json:"path"`
			Categories []string `json:"categories"`
		}
		if err := json.Unmarshal(p.Arguments, &args); err != nil {
			return nil, &rpcError{Code: -32602, Message: "invalid params"}
		}
		core.SetActiveCategories(args.Categories)
		findings, _, err := core.ScanFile(ctx, args.Path)
		if err != nil {
			return nil, &rpcError{Code: -32603, Message: err.Error()}
		}
		return textResult(findings), nil

	case "scan_dir":
		var args struct {
			Path       string   `json:"path"`
			Categories []string `json:"categories"`
		}
		if err := json.Unmarshal(p.Arguments, &args); err != nil {
			return nil, &rpcError{Code: -32602, Message: "invalid params"}
		}
		core.SetActiveCategories(args.Categories)
		findings, _, err := core.ScanDir(ctx, args.Path)
		if err != nil {
			return nil, &rpcError{Code: -32603, Message: err.Error()}
		}
		return textResult(findings), nil

	default:
		return nil, &rpcError{Code: -32601, Message: "unknown tool: " + p.Name}
	}
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
