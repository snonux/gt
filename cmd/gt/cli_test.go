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
		// Explicitly ignore error return from os.Remove during cleanup
		_ = os.Remove("/tmp/gt-test")
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

// TestCLIStdin tests that stdin input works correctly.
func TestCLIStdin(t *testing.T) {
	binaryPath := buildBinary(t)

	tests := []struct {
		name     string
		stdin    string
		expected string
	}{
		{
			name:     "RPN expression via stdin",
			stdin:    "3 4 +",
			expected: "7",
		},
		{
			name:     "Percentage expression via stdin",
			stdin:    "20% of 150",
			expected: "30",
		},
		{
			name:     "Variable assignment via stdin",
			stdin:    "x 5 = x x +",
			expected: "10",
		},
		{
			name:     "Boolean comparison via stdin",
			stdin:    "5 3 ==",
			expected: "false",
		},
		{
			name:     "Boolean to number coercion via stdin",
			stdin:    "true 2 *",
			expected: "2",
		},
		{
			name:     "Complex expression via stdin",
			stdin:    "2 3 + 4 *",
			expected: "20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath)
			cmd.Stdin = strings.NewReader(tt.stdin)
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

// TestCLIStdinWithPiping tests stdin via pipe (simulating `echo EXP | gt`).
func TestCLIStdinWithPiping(t *testing.T) {
	binaryPath := buildBinary(t)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "echo 3 4 +",
			input:    "3 4 +",
			expected: "7",
		},
		{
			name:     "echo 20%% of 150",
			input:    "20% of 150",
			expected: "30",
		},
		{
			name:     "echo x 5 = x x +",
			input:    "x 5 = x x +",
			expected: "10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate piping: echo INPUT | gt
			cmd := exec.Command(binaryPath)
			cmd.Stdin = strings.NewReader(tt.input)
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

// TestCLIStdinWithMultipleLines tests stdin with multiple lines (should use first line).
func TestCLIStdinWithMultipleLines(t *testing.T) {
	binaryPath := buildBinary(t)

	tests := []struct {
		name     string
		stdin    string
		expected string
	}{
		{
			name:     "Multiple lines - first line used",
			stdin:    "3 4 +\n5 6 +",
			expected: "7",
		},
		{
			name:     "Empty line then expression",
			stdin:    "\n3 4 +",
			expected: "7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath)
			cmd.Stdin = strings.NewReader(tt.stdin)
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

// TestCLIVariableAssignment tests all variable assignment syntaxes.
func TestCLIVariableAssignment(t *testing.T) {
	binaryPath := buildBinary(t)

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Standard assignment x = 2",
			args:     []string{"x", "2", "=", "x", "2", "+"},
			expected: "4",
		},
		{
			name:     "Right assignment x 2 := (value on stack, right)",
			args:     []string{"x", "2", ":=", "x", "2", "+"},
			expected: "4",
		},
		{
			name:     "Left assignment 2 x =: (value on stack, left)",
			args:     []string{"2", "x", "=: ", "x", "2", "+"},
			expected: "4",
		},
		{
			name:     "Stack variant with =: (value on stack, left)",
			args:     []string{"2", "x", "=: ", "x", "2", "+"},
			expected: "4",
		},
		{
			name:     "Assignment with existing variable",
			args:     []string{"x", "5", "=", "x", "3", "+"},
			expected: "8",
		},
		{
			name:     "Assignment with complex expression",
			args:     []string{"x", "2", "=", "x", "x", "*", "x", "+"},
			expected: "6",
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

// TestCLIVariableAssignmentSyntaxes tests each assignment syntax individually
// and verifies the variable can be used in subsequent expressions.
func TestCLIVariableAssignmentSyntaxes(t *testing.T) {
	binaryPath := buildBinary(t)

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Assignment: x = 2, then x 2 +",
			args:     []string{"x", "2", "=", "x", "2", "+"},
			expected: "4",
		},
		{
			name:     "Right assignment: x 2 :=, then x 2 +",
			args:     []string{"x", "2", ":=", "x", "2", "+"},
			expected: "4",
		},
		{
			name:     "Left assignment: 2 x =:, then x 2 +",
			args:     []string{"2", "x", "=: ", "x", "2", "+"},
			expected: "4",
		},
		{
			name:     "Assignment: y = 10, then y 5 *",
			args:     []string{"y", "10", "=", "y", "5", "*"},
			expected: "50",
		},
		{
			name:     "Assignment: pi = 3.14159, then pi 2 *",
			args:     []string{"pi", "3.14159", "=", "pi", "2", "*"},
			expected: "6.28318",
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

// TestCLIVariableAssignmentPreservation tests that variables persist within a single expression.
func TestCLIVariableAssignmentPreservation(t *testing.T) {
	binaryPath := buildBinary(t)

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Single assignment, multiple uses",
			args:     []string{"x", "5", "=", "x", "x", "+"},
			expected: "10", // x=5, then x+x=10
		},
		{
			name:     "Assignment with calculation as value (stack-variant)",
			args:     []string{"2", "3", "+", "x", "=: ", "x", "4", "*"},
			expected: "20", // 2+3=5, x=: assigns 5 to x, x*4=20
		},
		{
			name:     "Stack-variant with := (value on stack)",
			args:     []string{"x", "2", "3", "+", ":=", "x", "1", "+"},
			expected: "6", // x=2+3=5, x+1=6
		},
		{
			name:     "Stack-variant with =: (value first)",
			args:     []string{"2", "3", "+", "x", "=: ", "x", "1", "+"},
			expected: "6", // 2+3=5, x=: assigns 5 to x, x+1=6
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

// TestCLIVariableAssignmentRepetition tests that repeated assignment works correctly.
func TestCLIVariableAssignmentRepetition(t *testing.T) {
	binaryPath := buildBinary(t)

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Reassign same variable",
			args:     []string{"x", "5", ":=", "x", "10", ":=", "x", "2", "+"},
			expected: "12", // x=5, x=10, x+2=12
		},
		{
			name:     "Stack-variant reassignment",
			args:     []string{"x", "1", ":=", "x", "2", ":=", "x", "3", "+"},
			expected: "5", // x=1 then x=2, x+3=5
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
