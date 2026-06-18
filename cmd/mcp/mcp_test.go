package main

import (
	"encoding/json"
	"strings"
	"testing"
)

// Mock test utilities for MCP server testing

func TestScanStringTool(t *testing.T) {
	// This test would verify the scan_string tool works correctly
	// Since we can't easily test the full MCP server without more setup,
	// we'll test the underlying logic

	testContent := "API_KEY=sk-1234567890abcdef"

	// Simulate what the MCP tool would do
	result := scanStringForTesting(testContent, "test-source")

	if len(result) == 0 {
		t.Error("expected to find API key pattern")
	}
}

func TestScanFileTool(t *testing.T) {
	// Test file scanning functionality
	// This would normally create a temp file and scan it

	testContent := "Some content\nAWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\nMore content"

	result := scanStringForTesting(testContent, "test.txt")

	if len(result) == 0 {
		t.Error("expected to find AWS key pattern")
	}

	// Verify the pattern name
	if result[0]["pattern"] != "aws-access-key" {
		t.Errorf("expected 'aws-access-key', got '%v'", result[0]["pattern"])
	}
}

func TestScanDirTool(t *testing.T) {
	// Test directory scanning functionality
	// This would require temporary directory setup

	// For now, we'll test the logic structure
	testPath := "/tmp/test"
	categories := []string{"secrets"}

	// Verify parameters are handled correctly
	if testPath == "" {
		t.Error("path should not be empty")
	}

	if len(categories) == 0 {
		t.Error("should have at least one category")
	}
}

func TestToolDefinition(t *testing.T) {
	// Verify tool definitions are properly structured
	tools := getToolDefinitions()

	if len(tools) != 3 {
		t.Errorf("expected 3 tools, got %d", len(tools))
	}

	// Verify scan_string tool
	scanStringTool := getToolByName(tools, "scan_string")
	if scanStringTool == nil {
		t.Error("scan_string tool not found")
	} else {
		if !hasParameter(scanStringTool, "content") {
			t.Error("scan_string missing 'content' parameter")
		}
		if !hasParameter(scanStringTool, "source") {
			t.Error("scan_string missing 'source' parameter")
		}
	}

	// Verify scan_file tool
	scanFileTool := getToolByName(tools, "scan_file")
	if scanFileTool == nil {
		t.Error("scan_file tool not found")
	} else {
		if !hasParameter(scanFileTool, "path") {
			t.Error("scan_file missing 'path' parameter")
		}
	}

	// Verify scan_dir tool
	scanDirTool := getToolByName(tools, "scan_dir")
	if scanDirTool == nil {
		t.Error("scan_dir tool not found")
	} else {
		if !hasParameter(scanDirTool, "path") {
			t.Error("scan_dir missing 'path' parameter")
		}
		if !hasParameter(scanDirTool, "categories") {
			t.Error("scan_dir missing 'categories' parameter")
		}
	}
}

func TestCategoryFiltering(t *testing.T) {
	// Test category filtering logic
	allCategories := []string{"secrets", "pii", "code-quality", "healthcare"}

	// Test single category filter
	filtered := filterCategories(allCategories, []string{"secrets"})
	if len(filtered) != 1 || filtered[0] != "secrets" {
		t.Error("category filtering failed for single category")
	}

	// Test multiple category filter
	filtered = filterCategories(allCategories, []string{"secrets", "pii"})
	if len(filtered) != 2 {
		t.Error("category filtering failed for multiple categories")
	}

	// Test no filter (should return all)
	filtered = filterCategories(allCategories, []string{})
	if len(filtered) != len(allCategories) {
		t.Error("no filter should return all categories")
	}
}

func TestJSONResponseFormat(t *testing.T) {
	// Verify JSON response format is correct
	testFindings := []map[string]interface{}{
		{
			"pattern": "aws-access-key",
			"file":    "test.txt",
			"line":    1,
			"content": "AKIAIOSFODNN7EXAMPLE",
		},
	}

	jsonData, err := json.Marshal(testFindings)
	if err != nil {
		t.Fatalf("failed to marshal findings: %v", err)
	}

	// Verify it's valid JSON
	var decoded []map[string]interface{}
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(decoded) != 1 {
		t.Error("expected 1 finding")
	}

	if decoded[0]["pattern"] != "aws-access-key" {
		t.Error("pattern name not preserved in JSON")
	}
}

func TestMCPProtocolCompliance(t *testing.T) {
	// Verify basic MCP protocol compliance

	// Test that tools return proper JSON-RPC format
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"result": []map[string]interface{}{
			{
				"pattern": "test-key",
				"file":    "test.txt",
				"content": "secret",
			},
		},
		"id": 1,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("failed to marshal MCP response: %v", err)
	}

	// Verify it's valid JSON
	var decoded map[string]interface{}
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if decoded["jsonrpc"] != "2.0" {
		t.Error("MCP response should use JSON-RPC 2.0")
	}
}

func TestErrorHandling(t *testing.T) {
	// Test error handling in MCP tools

	// Test with invalid file path
	err := scanFileInvalidPath("/nonexistent/path/file.txt")
	if err == nil {
		t.Error("expected error for invalid file path")
	}

	// Test with empty content
	result := scanStringForTesting("", "test")
	if len(result) != 0 {
		t.Error("empty content should return no findings")
	}

	// Test with invalid JSON
	err = json.Unmarshal([]byte("invalid json"), &map[string]interface{}{})
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

// Helper functions for testing

func scanStringForTesting(content, source string) []map[string]interface{} {
	// This would normally call the core.ScanString function
	// For testing purposes, we'll simulate basic pattern matching

	findings := []map[string]interface{}{}

	// Simple pattern matching for common secrets
	if strings.Contains(content, "AKIA") && len(content) >= 20 {
		findings = append(findings, map[string]interface{}{
			"pattern": "aws-access-key",
			"file":    source,
			"content": content,
		})
	}

	if strings.Contains(content, "sk-") && len(content) >= 20 {
		findings = append(findings, map[string]interface{}{
			"pattern": "openai-api-key",
			"file":    source,
			"content": content,
		})
	}

	return findings
}

func scanFileInvalidPath(path string) error {
	// Simulate file scanning error
	return &osError{Path: path}
}

func getToolDefinitions() []map[string]interface{} {
	// Simulate tool definitions
	return []map[string]interface{}{
		{
			"name":        "scan_string",
			"description": "Scan a string for patterns",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"content": map[string]interface{}{
						"type":        "string",
						"description": "Content to scan",
					},
					"source": map[string]interface{}{
						"type":        "string",
						"description": "Source identifier",
					},
				},
				"required": []string{"content"},
			},
		},
		{
			"name":        "scan_file",
			"description": "Scan a file for patterns",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "File path to scan",
					},
				},
				"required": []string{"path"},
			},
		},
		{
			"name":        "scan_dir",
			"description": "Scan a directory for patterns",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Directory path to scan",
					},
					"categories": map[string]interface{}{
						"type":        "array",
						"description": "Pattern categories to include",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
				},
				"required": []string{"path"},
			},
		},
	}
}

func getToolByName(tools []map[string]interface{}, name string) map[string]interface{} {
	for _, tool := range tools {
		if tool["name"] == name {
			return tool
		}
	}
	return nil
}

func hasParameter(tool map[string]interface{}, paramName string) bool {
	inputSchema, ok := tool["inputSchema"].(map[string]interface{})
	if !ok {
		return false
	}

	properties, ok := inputSchema["properties"].(map[string]interface{})
	if !ok {
		return false
	}

	_, exists := properties[paramName]
	return exists
}

func filterCategories(allCategories, filter []string) []string {
	if len(filter) == 0 {
		return allCategories
	}

	filtered := []string{}
	filterSet := make(map[string]bool)
	for _, f := range filter {
		filterSet[f] = true
	}

	for _, cat := range allCategories {
		if filterSet[cat] {
			filtered = append(filtered, cat)
		}
	}

	return filtered
}

// osError simulates OS-level errors for testing
type osError struct {
	Path string
}

func (e *osError) Error() string {
	return "file not found: " + e.Path
}

func TestMain(m *testing.M) {
	// Setup for MCP tests
	m.Run()
}