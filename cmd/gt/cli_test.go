// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// buildBinary builds the gt binary to a temporary location.
func buildBinary(t *testing.T) string {
	t.Helper()
	
	// Get the project root directory
	projectRoot := os.Getenv("GITHUB_WORKSPACE")
	if projectRoot == "" {
		projectRoot = "/home/paul/git/gt"
	}

	buildCmd := exec.Command("go", "build", "-o", "/tmp/gt-test", "./cmd/gt")
	buildCmd.Dir = projectRoot
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("build failed: %v", err)
	}
	
	t.Cleanup(func() {
		os.Remove("/tmp/gt-test")
	})
	
	return "/tmp/gt-test"
}

// TestCLIVersion tests that the version command works correctly.
func TestCLIVersion(t *testing.T) {
	binaryPath := buildBinary(t)
	
	cmd := exec.Command(binaryPath, "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version command failed: %v\nOutput: %s", err, string(output))
	}

	versionOutput := strings.TrimSpace(string(output))
	if !strings.HasPrefix(versionOutput, "v") {
		t.Errorf("version output should start with 'v', got: %s", versionOutput)
	}
}

// TestCLIVersionOnly tests that the version command works correctly.
func TestCLIVersionOnly(t *testing.T) {
	binaryPath := buildBinary(t)
	
	cmd := exec.Command(binaryPath, "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version command failed: %v\nOutput: %s", err, string(output))
	}

	versionOutput := strings.TrimSpace(string(output))
	if !strings.HasPrefix(versionOutput, "v") {
		t.Errorf("version output should start with 'v', got: %s", versionOutput)
	}
}

// TestCLIPercentageCalculation tests percentage calculation commands.

// TestCLIPercentageCalculation tests percentage calculation commands.
func TestCLIPercentageCalculation(t *testing.T) {
	binaryPath := buildBinary(t)

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "20% of 150",
			args:     []string{"20% of 150"},
			expected: "30",
		},
		{
			name:     "what is 20% of 150",
			args:     []string{"what is 20% of 150"},
			expected: "30",
		},
		{
			name:     "30 is what % of 150",
			args:     []string{"30 is what % of 150"},
			expected: "20",
		},
		{
			name:     "30 is 20% of what",
			args:     []string{"30 is 20% of what"},
			expected: "150",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("command failed: %v\nOutput: %s", err, string(output))
			}

			outputStr := strings.TrimSpace(string(output))
			if !strings.Contains(outputStr, tt.expected) {
				t.Errorf("output should contain '%s', got: %s", tt.expected, outputStr)
			}
		})
	}
}

// TestCLIRPNCalculation tests RPN calculation commands.
func TestCLIRPNCalculation(t *testing.T) {
	binaryPath := buildBinary(t)

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "3 4 +",
			args:     []string{"3 4 +"},
			expected: "7",
		},
		{
			name:     "5 6 *",
			args:     []string{"5 6 *"},
			expected: "30",
		},
		{
			name:     "10 2 /",
			args:     []string{"10 2 /"},
			expected: "5",
		},
		{
			name:     "2 3 ^",
			args:     []string{"2 3 ^"},
			expected: "8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("command failed: %v\nOutput: %s", err, string(output))
			}

			outputStr := strings.TrimSpace(string(output))
			if !strings.Contains(outputStr, tt.expected) {
				t.Errorf("output should contain '%s', got: %s", tt.expected, outputStr)
			}
		})
	}
}

// TestCLIExitCode tests that the CLI returns proper exit codes.
func TestCLIExitCode(t *testing.T) {
	binaryPath := buildBinary(t)

	// Test successful command returns 0
	cmd := exec.Command(binaryPath, "version")
	if err := cmd.Run(); err != nil {
		t.Errorf("version command should succeed, got error: %v", err)
	}

	// Test invalid command returns non-zero
	cmd = exec.Command(binaryPath, "invalidcommand")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("invalid command should fail, but succeeded")
	}
	if len(output) > 0 && !strings.Contains(string(output), "Error:") {
		t.Errorf("error output should contain 'Error:', got: %s", string(output))
	}
}

// TestCLIInvalidPercentage tests invalid percentage calculation.
func TestCLIInvalidPercentage(t *testing.T) {
	binaryPath := buildBinary(t)

	cmd := exec.Command(binaryPath, "invalid percentage")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("invalid percentage should fail, but succeeded")
	}
	if len(output) > 0 && !strings.Contains(string(output), "Error:") {
		t.Errorf("error output should contain 'Error:', got: %s", string(output))
	}
}

// TestCLIInvalidRPN tests invalid RPN expression.
func TestCLIInvalidRPN(t *testing.T) {
	binaryPath := buildBinary(t)

	cmd := exec.Command(binaryPath, "3 +")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("invalid RPN should fail, but succeeded")
	}
	if len(output) > 0 && !strings.Contains(string(output), "Error:") {
		t.Errorf("error output should contain 'Error:', got: %s", string(output))
	}
}
