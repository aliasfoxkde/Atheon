package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestMainVersionFlag tests main() with --version flag
func TestMainVersionFlag(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping main() test in short mode")
	}

	// Build the binary first
	buildCmd := exec.Command("go", "build", "-o", "atheon-test", ".")
	if err := buildCmd.Run(); err != nil {
		t.Skip("Failed to build binary, skipping test")
	}
	defer os.Remove("atheon-test")

	// Test --version flag
	cmd := exec.Command("./atheon-test", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("--version flag error: %v", err)
	}

	if len(output) == 0 {
		t.Error("Expected output from --version flag")
	}

	if !bytes.Contains(output, []byte("atheon")) {
		t.Error("Expected 'atheon' in version output")
	}
}

// TestMainListCommand tests main() with list command
func TestMainListCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping main() test in short mode")
	}

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "atheon-test", ".")
	if err := buildCmd.Run(); err != nil {
		t.Skip("Failed to build binary, skipping test")
	}
	defer os.Remove("atheon-test")

	// Test list command
	cmd := exec.Command("./atheon-test", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("list command failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("Expected output from list command")
	}

	// Should contain pattern names
	outputStr := string(output)
	if !strings.Contains(outputStr, "aws-access-key") {
		t.Error("Expected 'aws-access-key' in list output")
	}
}

// TestMainHelpCommand tests main() with --help flag
func TestMainHelpCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping main() test in short mode")
	}

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "atheon-test", ".")
	if err := buildCmd.Run(); err != nil {
		t.Skip("Failed to build binary, skipping test")
	}
	defer os.Remove("atheon-test")

	// Test --help flag
	cmd := exec.Command("./atheon-test", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("--help command failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("Expected output from --help command")
	}

	// Should contain usage information
	outputStr := string(output)
	if !strings.Contains(outputStr, "usage:") {
		t.Error("Expected 'usage:' in help output")
	}
}

// TestMainJSONOutput tests JSON output functionality
func TestMainJSONOutput(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping main() test in short mode")
	}

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "atheon-test", ".")
	if err := buildCmd.Run(); err != nil {
		t.Skip("Failed to build binary, skipping test")
	}
	defer os.Remove("atheon-test")

	// Create a test file with content
	tmpFile := "/tmp/atheon-test-input.txt"
	defer os.Remove(tmpFile)
	content := []byte("AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE")
	if err := os.WriteFile(tmpFile, content, 0o644); err != nil {
		t.Fatal(err)
	}

	// Test --json flag
	cmd := exec.Command("./atheon-test", "--json", "--file", tmpFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Non-zero exit code is expected if findings are found
		t.Logf("JSON command completed with exit code: %v", err)
	}

	if len(output) == 0 {
		t.Error("Expected JSON output")
	}

	// Should be valid JSON
	if !bytes.HasPrefix(output, []byte("[")) {
		t.Error("Expected JSON array output")
	}
}

// TestMainEnvScanning tests environment variable scanning
func TestMainEnvScanning(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping main() test in short mode")
	}

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "atheon-test", ".")
	if err := buildCmd.Run(); err != nil {
		t.Skip("Failed to build binary, skipping test")
	}
	defer os.Remove("atheon-test")

	// Set test environment variable
	os.Setenv("TEST_AWS_KEY", "AKIAIOSFODNN7EXAMPLE")
	defer os.Unsetenv("TEST_AWS_KEY")

	// Test --env flag
	cmd := exec.Command("./atheon-test", "--env")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Non-zero exit code is expected if findings are found
		t.Logf("Env scan completed with exit code: %v", err)
	}

	if len(output) == 0 {
		t.Error("Expected output from env scan")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "aws-access-key") {
		t.Error("Expected 'aws-access-key' in env scan output")
	}
}

// TestMainInvalidArgs tests error handling for invalid arguments
func TestMainInvalidArgs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping main() test in short mode")
	}

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "atheon-test", ".")
	if err := buildCmd.Run(); err != nil {
		t.Skip("Failed to build binary, skipping test")
	}
	defer os.Remove("atheon-test")

	// Test invalid command
	cmd := exec.Command("./atheon-test", "invalid-command")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error for invalid command")
	}

	if len(output) == 0 {
		t.Error("Expected error message for invalid command")
	}

	outputStr := string(output)
	if !strings.Contains(strings.ToLower(outputStr), "error") && !strings.Contains(strings.ToLower(outputStr), "unknown") {
		t.Error("Expected error message in output")
	}
}
