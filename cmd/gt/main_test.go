package main

import (
	"os"
	"strings"
	"testing"
)

func TestRunCommandVersion(t *testing.T) {
	args := []string{"gt", "version"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand(['gt', 'version']) returned error: %v", err)
	}
	if result != "dev" && !strings.HasPrefix(result, "v") {
		t.Errorf("runCommand(['gt', 'version']) = %q, expected version string", result)
	}
}

func TestRunCommandCalc(t *testing.T) {
	// RPN expressions are now parsed directly without 'calc' prefix
	args := []string{"gt", "3", "4", "+"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand(['gt', '3', '4', '+']) returned error: %v", err)
	}
	if result != "7" {
		t.Errorf("runCommand(['gt', '3', '4', '+']) = %q, want '7'", result)
	}
}

func TestRunCommandRPN(t *testing.T) {
	// RPN expressions are now parsed directly without 'rpn' prefix
	args := []string{"gt", "3", "4", "+"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand(['gt', '3', '4', '+']) returned error: %v", err)
	}
	if result != "7" {
		t.Errorf("runCommand(['gt', '3', '4', '+']) = %q, want '7'", result)
	}
}

func TestRunCommandRPNWithAssignment(t *testing.T) {
	// RPN expressions with assignment are now parsed directly
	args := []string{"gt", "x", "5", "=", "x", "x", "+"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand with assignment returned error: %v", err)
	}
	if result != "10" {
		t.Errorf("runCommand with assignment = %q, want '10'", result)
	}
}

func TestRunCommandPercentage(t *testing.T) {
	args := []string{"gt", "20% of 150"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand(['gt', '20%% of 150']) returned error: %v", err)
	}
	if !strings.Contains(result, "30") {
		t.Errorf("runCommand(['gt', '20%% of 150']) = %q, should contain '30'", result)
	}
}

func TestRunCommandInvalidRPN(t *testing.T) {
	args := []string{"gt", "5", "0", "/"}
	_, err := runCommand(args)
	if err == nil {
		t.Error("runCommand(['gt', '5', '0', '/']) should return error for division by zero")
	}
}

func TestRunCommandUnknownToken(t *testing.T) {
	// Unknown token in RPN expression should fail
	args := []string{"gt", "unknown"}
	_, err := runCommand(args)
	if err == nil {
		t.Error("runCommand(['gt', 'unknown']) should return error")
	}
}

func TestPrintUsage(t *testing.T) {
	// Just verify the function doesn't panic
	// We can't easily test the output since it goes to stdout
	printUsage()
}

func TestRunCommandUnknownInput(t *testing.T) {
	// Unknown input should fail
	args := []string{"gt", "unknown 3 4 +"}
	_, err := runCommand(args)
	if err == nil {
		t.Error("runCommand with unknown input should return error")
	}
}

func TestMain(t *testing.T) {
	// Save original os.Args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test with version command
	os.Args = []string{"gt", "version"}
	// Note: we can't actually call main() in tests because it calls os.Exit()
	// Instead we test via runCommand which is what main() calls
	result, err := runCommand(os.Args)
	if err != nil {
		t.Fatalf("runCommand(['gt', 'version']) returned error: %v", err)
	}
	if result != "dev" && !strings.HasPrefix(result, "v") {
		t.Errorf("runCommand(['gt', 'version']) = %q, expected version string", result)
	}
}

func TestRunCommandCalcChain(t *testing.T) {
	// RPN expression chain without 'calc' prefix
	args := []string{"gt", "3", "4", "+", "4", "4", "-", "*"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand with chain returned error: %v", err)
	}
	if result != "0" {
		t.Errorf("runCommand with chain = %q, want '0'", result)
	}
}

func TestRunCommandRPNPower(t *testing.T) {
	// RPN power without 'rpn' prefix
	args := []string{"gt", "2", "3", "^"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand with power returned error: %v", err)
	}
	if result != "8" {
		t.Errorf("runCommand with power = %q, want '8'", result)
	}
}

func TestRunCommandRPNModulo(t *testing.T) {
	// RPN modulo without 'rpn' prefix
	args := []string{"gt", "10", "3", "%"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand with modulo returned error: %v", err)
	}
	if result != "1" {
		t.Errorf("runCommand with modulo = %q, want '1'", result)
	}
}

func TestRunCommandNoArgs(t *testing.T) {
	// Test with no arguments (simulating stdin not being TTY)
	args := []string{"gt"}
	_, err := runCommand(args)
	if err == nil {
		t.Error("runCommand with no args should return error")
	}
	if !strings.Contains(err.Error(), "no input provided") {
		t.Errorf("Error = %v, should contain 'no input provided'", err)
	}
}

// The following tests were removed because they tested subcommand handling
// which has been removed:
// - TestRunCommandRepl (repl command)
// - TestRunCommandReplFlag (--repl flag)
// - TestRunCommandCalcWithShow (calc with show)
// - TestRunCommandCalcWithVars (calc with vars)
// - TestRunCommandCalcWithClear (calc with clear)

// These commands are now only available in REPL mode, not in command-line mode.
