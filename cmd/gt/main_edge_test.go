// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package main

import (
	"testing"
)

// TestRunCommandEmptyArgs tests running with no arguments (just the program name)
func TestRunCommandEmptyArgs(t *testing.T) {
	// With only program name, it would try REPL mode (needs TTY)
	// or read from stdin. In test context, stdin is not a TTY,
	// so it tries to read from stdin which returns empty.
	_, err := runCommand([]string{"gt"})
	if err == nil {
		t.Error("expected error when no input provided")
	}
}

// TestRunCommandVersionFlag tests version flag
func TestRunCommandVersionFlag(t *testing.T) {
	result, err := runCommand([]string{"gt", "version"})
	if err != nil {
		t.Fatalf("runCommand version returned error: %v", err)
	}
	if result == "" {
		t.Error("version output should not be empty")
	}
	// Version should be a semver-like string
	if result[0] != 'v' && result[0] != '0' && result[0] != '1' {
		t.Errorf("version output format unexpected: %q", result)
	}
}

// TestRunCommandLogFlagWithArgs tests --log flag with calculation
func TestRunCommandLogFlagWithArgs(t *testing.T) {
	result, err := runCommand([]string{"gt", "--log", "/tmp/test-gt-log", "3", "4", "+"})
	if err != nil {
		t.Fatalf("runCommand with --log returned error: %v", err)
	}
	if result != "7" {
		t.Errorf("runCommand with --log = %q, want '7'", result)
	}
}

// TestRunCommandPercentageCalculation tests percentage calculation
func TestRunCommandPercentageCalculation(t *testing.T) {
	result, err := runCommand([]string{"gt", "20% of 150"})
	if err != nil {
		t.Fatalf("runCommand percentage returned error: %v", err)
	}
	if result == "" {
		t.Error("percentage result should not be empty")
	}
}

// TestRunCommandInvalidPercentage tests invalid percentage expression error
func TestRunCommandInvalidPercentage(t *testing.T) {
	_, err := runCommand([]string{"gt", "invalid percentage stuff here"})
	if err == nil {
		t.Error("expected error for invalid percentage expression")
	}
}

// TestRunCommandMultipleArgsAsExpression tests multiple CLI args as one expression
func TestRunCommandMultipleArgsAsExpression(t *testing.T) {
	// Multiple args are joined with spaces
	result, err := runCommand([]string{"gt", "50% of", "200"})
	if err != nil {
		t.Fatalf("runCommand multi-arg percentage returned error: %v", err)
	}
	if result == "" {
		t.Error("multi-arg percentage result should not be empty")
	}
}

// TestRunCommandLogFlagOnly tests --log flag with no other args
func TestRunCommandLogFlagOnly(t *testing.T) {
	// --log alone with no other args would trigger REPL/stdin
	_, err := runCommand([]string{"gt", "--log", "/tmp/test-gt-log"})
	if err == nil {
		t.Log("log flag alone returned nil (might trigger REPL)")
	}
	// Either way, it shouldn't panic
}

// TestRunCommandRPNWithLog tests RPN calculation with --log flag
func TestRunCommandRPNWithLog(t *testing.T) {
	result, err := runCommand([]string{"gt", "--log", "/tmp/test-gt-log2", "10", "2", "/"})
	if err != nil {
		t.Fatalf("runCommand RPN with --log returned error: %v", err)
	}
	if result != "5" {
		t.Errorf("RPN with --log = %q, want '5'", result)
	}
}

// TestRunCommandUnknownRPNToken tests unknown RPN token error
func TestRunCommandUnknownRPNToken(t *testing.T) {
	_, err := runCommand([]string{"gt", "3", "unknown_token", "4", "+"})
	if err == nil {
		t.Error("expected error for unknown RPN token")
	}
}

// TestRunCommandMixedArgs tests mixed RPN and percentage args
func TestRunCommandMixedArgs(t *testing.T) {
	// This tests the fallback path: RPN fails, tries percentage
	result, err := runCommand([]string{"gt", "25% of", "400"})
	if err != nil {
		t.Fatalf("runCommand mixed args returned error: %v", err)
	}
	if result == "" {
		t.Error("mixed args result should not be empty")
	}
}

// TestRunCommandRPNPowerWithLog tests RPN power with --log
func TestRunCommandRPNPowerWithLog(t *testing.T) {
	result, err := runCommand([]string{"gt", "--log", "/tmp/test-gt-log3", "2", "3", "^"})
	if err != nil {
		t.Fatalf("runCommand RPN power with --log returned error: %v", err)
	}
	if result != "8" {
		t.Errorf("RPN power with --log = %q, want '8'", result)
	}
}

// TestRunCommandRPNAssignmentWithLog tests RPN assignment with --log
func TestRunCommandRPNAssignmentWithLog(t *testing.T) {
	result, err := runCommand([]string{"gt", "--log", "/tmp/test-gt-log4", "x", "5", "="})
	if err != nil {
		t.Fatalf("runCommand assignment with --log returned error: %v", err)
	}
	// Assignment returns "x = 5"
	if result != "x = 5" {
		t.Errorf("assignment with --log = %q, want 'x = 5'", result)
	}
}
