package main

import (
	"os"
	"strings"
	"testing"
)

func TestRunCommandVersion(t *testing.T) {
	args := []string{"perc", "version"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand(['perc', 'version']) returned error: %v", err)
	}
	if result != "dev" && !strings.HasPrefix(result, "v") {
		t.Errorf("runCommand(['perc', 'version']) = %q, expected version string", result)
	}
}

func TestRunCommandCalc(t *testing.T) {
	args := []string{"perc", "calc", "3", "4", "+"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand(['perc', 'calc', '3', '4', '+']) returned error: %v", err)
	}
	if result != "7" {
		t.Errorf("runCommand(['perc', 'calc', '3', '4', '+']) = %q, want '7'", result)
	}
}

func TestRunCommandRPN(t *testing.T) {
	args := []string{"perc", "rpn", "3", "4", "+"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand(['perc', 'rpn', '3', '4', '+']) returned error: %v", err)
	}
	if result != "7" {
		t.Errorf("runCommand(['perc', 'rpn', '3', '4', '+']) = %q, want '7'", result)
	}
}

func TestRunCommandRPNWithAssignment(t *testing.T) {
	args := []string{"perc", "rpn", "x", "5", "=", "x", "x", "+"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand with assignment returned error: %v", err)
	}
	if result != "10" {
		t.Errorf("runCommand with assignment = %q, want '10'", result)
	}
}

func TestRunCommandPercentage(t *testing.T) {
	args := []string{"perc", "20% of 150"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand(['perc', '20%% of 150']) returned error: %v", err)
	}
	if !strings.Contains(result, "30") {
		t.Errorf("runCommand(['perc', '20%% of 150']) = %q, should contain '30'", result)
	}
}

func TestRunCommandMissingExpression(t *testing.T) {
	args := []string{"perc", "calc"}
	_, err := runCommand(args)
	if err == nil {
		t.Error("runCommand(['perc', 'calc']) should return error for missing expression")
	}
	if !strings.Contains(err.Error(), "missing expression") {
		t.Errorf("Error = %v, should contain 'missing expression'", err)
	}
}

func TestRunCommandInvalidRPN(t *testing.T) {
	args := []string{"perc", "rpn", "5", "0", "/"}
	_, err := runCommand(args)
	if err == nil {
		t.Error("runCommand(['perc', 'rpn', '5', '0', '/']) should return error for division by zero")
	}
}

func TestRunCommandUnknownToken(t *testing.T) {
	args := []string{"perc", "rpn", "unknown +"}
	_, err := runCommand(args)
	if err == nil {
		t.Error("runCommand(['perc', 'rpn', 'unknown +']) should return error")
	}
}

func TestPrintUsage(t *testing.T) {
	// Just verify the function doesn't panic
	// We can't easily test the output since it goes to stdout
	printUsage()
}

func TestRunCommandUnknownSubcommand(t *testing.T) {
	args := []string{"perc", "unknown", "3", "4", "+"}
	// This will fall through to calculator.Parse which will fail
	_, err := runCommand(args)
	if err == nil {
		t.Error("runCommand with unknown subcommand should return error")
	}
}

func TestMain(t *testing.T) {
	// Save original os.Args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test with version command
	os.Args = []string{"perc", "version"}
	// Note: we can't actually call main() in tests because it calls os.Exit()
	// Instead we test via runCommand which is what main() calls
	result, err := runCommand(os.Args)
	if err != nil {
		t.Fatalf("runCommand(['perc', 'version']) returned error: %v", err)
	}
	if result != "dev" && !strings.HasPrefix(result, "v") {
		t.Errorf("runCommand(['perc', 'version']) = %q, expected version string", result)
	}
}

func TestRunCommandCalcChain(t *testing.T) {
	args := []string{"perc", "calc", "3", "4", "+", "4", "4", "-", "*"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand with chain returned error: %v", err)
	}
	if result != "0" {
		t.Errorf("runCommand with chain = %q, want '0'", result)
	}
}

func TestRunCommandRPNPower(t *testing.T) {
	args := []string{"perc", "rpn", "2", "3", "^"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand with power returned error: %v", err)
	}
	if result != "8" {
		t.Errorf("runCommand with power = %q, want '8'", result)
	}
}

func TestRunCommandRPNModulo(t *testing.T) {
	args := []string{"perc", "rpn", "10", "3", "%"}
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
	args := []string{"perc"}
	_, err := runCommand(args)
	if err == nil {
		t.Error("runCommand with no args should return error")
	}
	if !strings.Contains(err.Error(), "no input provided") {
		t.Errorf("Error = %v, should contain 'no input provided'", err)
	}
}

func TestRunCommandRepl(t *testing.T) {
	// Test repl command (returns empty string, doesn't start REPL in tests)
	args := []string{"perc", "repl"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand(['perc', 'repl']) returned error: %v", err)
	}
	if result != "" {
		t.Errorf("runCommand(['perc', 'repl']) = %q, want empty string", result)
	}
}

func TestRunCommandReplFlag(t *testing.T) {
	// Test --repl flag
	args := []string{"perc", "--repl"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand(['perc', '--repl']) returned error: %v", err)
	}
	if result != "" {
		t.Errorf("runCommand(['perc', '--repl']) = %q, want empty string", result)
	}
}

func TestRunCommandRPNWithVariables(t *testing.T) {
	// Test rpn with single variable assignment and usage
	args := []string{"perc", "rpn", "x", "5", "=", "x", "x", "+"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand with variables returned error: %v", err)
	}
	if result != "10" {
		t.Errorf("runCommand with variables = %q, want '10'", result)
	}
}

func TestRunCommandCalcWithShow(t *testing.T) {
	// Test calc with show command
	args := []string{"perc", "calc", "1", "2", "3", "show"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand with show returned error: %v", err)
	}
	if result != "1 2 3" {
		t.Errorf("runCommand with show = %q, want '1 2 3'", result)
	}
}

func TestRunCommandCalcWithVars(t *testing.T) {
	// Test calc with vars command
	args := []string{"perc", "calc", "x", "5", "=", "vars"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand with vars returned error: %v", err)
	}
	if !strings.Contains(result, "x") {
		t.Errorf("runCommand with vars = %q, should contain 'x'", result)
	}
}

func TestRunCommandCalcWithClear(t *testing.T) {
	// Test calc with clear command
	args := []string{"perc", "calc", "x", "5", "=", "clear"}
	result, err := runCommand(args)
	if err != nil {
		t.Fatalf("runCommand with clear returned error: %v", err)
	}
	if !strings.Contains(result, "All variables cleared") {
		t.Errorf("runCommand with clear = %q, should contain 'All variables cleared'", result)
	}
}
