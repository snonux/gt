// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"strings"
	"testing"

	"codeberg.org/snonux/gt/internal/rpn"
)

// Helper to create a minimal REPL for testing without prompt (no TTY required)
func createTestREPL() *REPL {
	vars := rpn.NewVariables()
	return &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  NewCommandChain(),
		rpnState:      &RPNState{vars: vars, calculator: NewRPNCalculator(rpn.NewRPN(vars))},
	}
}

func TestExecutorWithHelp(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "help")
}

func TestExecutorWithClear(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "clear")
}

func TestExecutorWithQuit(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "quit")
}

func TestExecutorWithExit(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "exit")
}

func TestExecutorWithPercentage(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "20% of 150")
}

func TestExecutorWithRPN(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "rpn 3 4 +")
}

func TestExecutorWithInvalid(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "invalid input")
}

func TestExecutorWithVars(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "rpn x 5 = vars")
}

func TestExecutorWithClearVariables(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "rpn clear")
}

func TestIsBuiltinCommand(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"help", true},
		{"clear", true},
		{"quit", true},
		{"exit", true},
		{"rpn", true},
		{"calc", true},
		{"20% of 150", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, ok := isBuiltinCommand(tt.input)
			if ok != tt.expected {
				t.Errorf("isBuiltinCommand(%q) = %v, want %v", tt.input, ok, tt.expected)
			}
		})
	}
}

func TestIsBuiltinCommandWithSubcommand(t *testing.T) {
	_, ok := isBuiltinCommand("help clear")
	if !ok {
		t.Error("isBuiltinCommand('help clear') should return true")
	}
}

func TestIsBuiltinCommandEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty string", "", false},
		{"single space", " ", false},
		{"case insensitive - HELP", "HELP", true},
		{"case insensitive - HeLp", "HeLp", true},
		{"command with extra spaces", "  help  ", true},
		{"partial match - hel", "hel", false},
		{"partial match - cal", "cal", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := isBuiltinCommand(tt.input)
			if ok != tt.expected {
				t.Errorf("isBuiltinCommand(%q) = %v, want %v", tt.input, ok, tt.expected)
			}
		})
	}
}

func TestIsBuiltinCommandEdgeCasesWithAllCommands(t *testing.T) {
	allCommands := []string{"help", "clear", "quit", "exit", "rpn", "calc"}
	for _, cmd := range allCommands {
		t.Run(cmd, func(t *testing.T) {
			input, ok := isBuiltinCommand(cmd)
			if !ok {
				t.Errorf("isBuiltinCommand(%q) should return true for builtin command", cmd)
			}
			if input != cmd {
				t.Errorf("isBuiltinCommand(%q) returned %q, want %q", cmd, input, cmd)
			}
		})
	}
}

func TestIsBuiltinCommandWithMixedCase(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"HELP", true},
		{"Help", true},
		{"hElP", true},
		{"CLEAR", true},
		{"quit", true},
		{"QUIT", true},
		{"RPN", true},
		{"calc", true},
		{"CALC", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, ok := isBuiltinCommand(tt.input)
			if ok != tt.expected {
				t.Errorf("isBuiltinCommand(%q) = %v, want %v", tt.input, ok, tt.expected)
			}
		})
	}
}

// TestRunRPN tests inline RPN evaluation (like cmd/gt/main.go does)
func TestRunRPN(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"simple addition", "3 4 +", false},
		{"simple subtraction", "10 3 -", false},
		{"simple multiplication", "2 3 *", false},
		{"simple division", "10 2 /", false},
		{"power operation", "2 3 ^", false},
		{"modulo operation", "10 3 %", false},
		{"with variables", "x 5 = x x +", false},
		{"empty input", "", true},
		{"invalid input", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars := rpn.NewVariables()
			rpnCalc := rpn.NewRPN(vars)

			input := strings.TrimSpace(tt.input)
			if strings.HasPrefix(input, "rpn ") {
				input = strings.TrimPrefix(input, "rpn ")
			} else if strings.HasPrefix(input, "calc ") {
				input = strings.TrimPrefix(input, "calc ")
			}

			_, err := rpnCalc.ParseAndEvaluate(input)
			if (err != nil) != tt.wantErr {
				t.Errorf("RPN evaluation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetCommandDescription(t *testing.T) {
	tests := []struct {
		cmd        string
		wantPrefix string
	}{
		{"help", "Show help"},
		{"clear", "Clear"},
		{"quit", "Exit"},
		{"exit", "Exit"},
		{"rpn", "Evaluate an RPN"},
		{"calc", "Same as rpn"},
		{"unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.cmd, func(t *testing.T) {
			desc := getCommandDescription(tt.cmd)
			if tt.wantPrefix != "" && !strings.Contains(desc, tt.wantPrefix) {
				t.Errorf("getCommandDescription(%q) = %q, should contain %q", tt.cmd, desc, tt.wantPrefix)
			}
		})
	}
}

func TestGetCommandDescriptionForUnknownCommand(t *testing.T) {
	desc := getCommandDescription("unknown")
	if desc != "" {
		t.Errorf("getCommandDescription(%q) = %q, should be empty", "unknown", desc)
	}
}

func TestExecutorWithSingleOperator(t *testing.T) {
	repl := createTestREPL()
	for _, op := range []string{"+", "-", "*", "/", "^", "%", "dup", "swap", "pop", "show", "vars", "clear"} {
		t.Run(op, func(t *testing.T) {
			defaultExecutor(repl, op)
		})
	}
}

func TestExecutorWithPercentageExpression(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "20% of 150")
	defaultExecutor(repl, "30 is what %% of 150")
	defaultExecutor(repl, "30 is 20%% of what")
}

func TestExecutorWithInvalidPercentage(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "invalid percentage input")
}

func TestExecutorWithOperatorOnly(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "1 2 +")
	defaultExecutor(repl, "+")
}

func TestExecutorWithRPNPrefix(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "rpn 3 4 +")
}

func TestExecutorWithCalcPrefix(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "calc 5 6 +")
}

func TestExecutorWithEmptyInput(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "")
}

func TestExecutorWithWhitespaceOnly(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "   ")
}

func TestExecutorWithInvalidInput(t *testing.T) {
	tests := []string{"invalid input", "not a valid command", "xyz"}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			repl := createTestREPL()
			defaultExecutor(repl, input)
		})
	}
}

func TestExecutorWithInvalidRPN(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "rpn 1 +")
}

func TestExecutorWithEmptyRPNPrefix(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "rpn")
	defaultExecutor(repl, "calc")
}

func TestExecutorWithAssignment(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "rpn x 42 =")
	defaultExecutor(repl, "rpn x")
}

func TestExecutorWithPercentageAndRPNFallback(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "20% of 150")
	defaultExecutor(repl, "3 4 +")
}

func TestExecutorWithRPNExpressionOnly(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "5 3 +")
}

func TestExecutorWithRPNThenOperator(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "1 2 +")
	defaultExecutor(repl, "+")
}

func TestExecutorWithRPNThenRPN(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "rpn 1 2 +")
	defaultExecutor(repl, "rpn 3 4 +")
}

func TestExecutorWithRPNShow(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "rpn show")
}

func TestExecutorWithRPNDup(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "rpn dup")
}

func TestExecutorWithRPNSwap(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "rpn swap")
}

func TestExecutorWithRPNSingle(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "rpn 42")
}

func TestExecutorWithRPNMulti(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "rpn 1 2 3 4 5 +")
}

func TestExecutorWithStackOps(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "dup")
	defaultExecutor(repl, "swap")
	defaultExecutor(repl, "pop")
	defaultExecutor(repl, "show")
}

func TestExecutorWithRPNClear(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "rpn clear")
}

func TestExecutorWithHistoryCommands(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "vars")
	defaultExecutor(repl, "clear")
}

func TestExecutorWithMixedInput(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "25% of 200")
	defaultExecutor(repl, "10 20 +")
}

func TestExecutorWithRPNCalcMixed(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "rpn 1 2 +")
	defaultExecutor(repl, "3 4 +")
	defaultExecutor(repl, "calc 5 6 +")
}

func TestExecutorCommandsEdgeCases(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "  clear  ")
	defaultExecutor(repl, "HELP")
	defaultExecutor(repl, "CLEAR")
}

func TestExecutorWithRPMPrefix(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "rpn 1 2 +")
}

func TestExecutorWithCalcPrefixMixed(t *testing.T) {
	repl := createTestREPL()
	defaultExecutor(repl, "calc 1 2 +")
}

// TestExecutorWithRatModeOn tests that rat on works with fresh REPL
func TestExecutorWithRatModeOn(t *testing.T) {
	vars := rpn.NewVariables()
	rpnCalc := rpn.NewRPN(vars)
	rpl := &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  NewCommandChain(),
		rpnState:      &RPNState{vars: vars, calculator: NewRPNCalculator(rpnCalc)},
	}
	defaultExecutor(rpl, "rat on")
	if rpnCalc.GetMode() != rpn.RationalMode {
		t.Errorf("Expected RationalMode after rat on, got %v", rpnCalc.GetMode())
	}
}

// TestExecutorWithRatModeOff tests that rat off works with fresh REPL
func TestExecutorWithRatModeOff(t *testing.T) {
	vars := rpn.NewVariables()
	rpnCalc := rpn.NewRPN(vars)
	rpl := &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  NewCommandChain(),
		rpnState:      &RPNState{vars: vars, calculator: NewRPNCalculator(rpnCalc)},
	}
	defaultExecutor(rpl, "rat off")
	if rpnCalc.GetMode() != rpn.FloatMode {
		t.Errorf("Expected FloatMode after rat off, got %v", rpnCalc.GetMode())
	}
}

// TestExecutorWithRatModeToggle tests that rat toggle works with fresh REPL
func TestExecutorWithRatModeToggle(t *testing.T) {
	vars := rpn.NewVariables()
	rpnCalc := rpn.NewRPN(vars)
	rpl := &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  NewCommandChain(),
		rpnState:      &RPNState{vars: vars, calculator: NewRPNCalculator(rpnCalc)},
	}
	
	// First toggle
	defaultExecutor(rpl, "rat toggle")
	mode1 := rpnCalc.GetMode()
	
	// Second toggle
	defaultExecutor(rpl, "rat toggle")
	mode2 := rpnCalc.GetMode()
	
	// Modes should be different after toggle
	if mode1 == mode2 {
		t.Errorf("Modes should be different after toggle: %v -> %v", mode1, mode2)
	}
}

func TestExecutorWithRatModeInvalid(t *testing.T) {
	repl := createTestREPL()
	// Just verify it doesn't panic
	defaultExecutor(repl, "rat invalid")
}

func TestExecutorWithRatModeNoArg(t *testing.T) {
	repl := createTestREPL()
	// Just verify it doesn't panic
	defaultExecutor(repl, "rat")
}

func TestIsBuiltinCommandWithSubcommandHelp(t *testing.T) {
	_, ok := isBuiltinCommand("help")
	if !ok {
		t.Error("isBuiltinCommand('help') should return true")
	}
}

// TestExecutorWithAssignmentRight tests := and =: operators
func TestExecutorWithAssignmentRight(t *testing.T) {
	vars := rpn.NewVariables()
	rpnCalc := rpn.NewRPN(vars)
	rpl := &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  NewCommandChain(),
		rpnState:      &RPNState{vars: vars, calculator: NewRPNCalculator(rpnCalc)},
	}
	
	// Test := operator
	defaultExecutor(rpl, "5 x :=")
	val, exists := vars.GetVariable("x")
	if !exists {
		t.Errorf("Variable x should exist after x :=")
	}
	if val != 5 {
		t.Errorf("Variable x = %v, want 5", val)
	}
	
	// Test =: operator
	defaultExecutor(rpl, "y 3 =:")
	val, exists = vars.GetVariable("y")
	if !exists {
		t.Errorf("Variable y should exist after y =:")
	}
	if val != 3 {
		t.Errorf("Variable y = %v, want 3", val)
	}
}

// TestExecutorWithAssignmentAfterCalculation tests assignment after a calculation
func TestExecutorWithAssignmentAfterCalculation(t *testing.T) {
	vars := rpn.NewVariables()
	rpnCalc := rpn.NewRPN(vars)
	rpl := &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  NewCommandChain(),
		rpnState:      &RPNState{vars: vars, calculator: NewRPNCalculator(rpnCalc)},
	}
	
	// Test that assignment works after a calculation
	defaultExecutor(rpl, "1 2 + z =:")
	val, exists := vars.GetVariable("z")
	if !exists {
		t.Errorf("Variable z should exist")
	}
	if val != 3 {
		t.Errorf("Variable z = %v, want 3", val)
	}
}

// TestExecutorWithIncrementalAssignment tests that assignment works after a calculation with separate commands
func TestExecutorWithIncrementalAssignment(t *testing.T) {
	vars := rpn.NewVariables()
	rpnCalc := rpn.NewRPN(vars)
	rpl := &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  NewCommandChain(),
		rpnState:      &RPNState{vars: vars, calculator: NewRPNCalculator(rpnCalc)},
	}
	
	// Test that assignment works after a calculation
	defaultExecutor(rpl, "1 2 +")
	
	// Now use z =: to assign the top of stack (3) to variable z
	defaultExecutor(rpl, "z =:")
	
	val, exists := vars.GetVariable("z")
	if !exists {
		t.Errorf("Variable z should exist after z =:")
	}
	if val != 3 {
		t.Errorf("Variable z = %v, want 3", val)
	}
}

// TestExecutorWithSimpleIncrementalAssignment tests x =: after 2 in REPL
func TestExecutorWithSimpleIncrementalAssignment(t *testing.T) {
	vars := rpn.NewVariables()
	rpnCalc := rpn.NewRPN(vars)
	rpl := &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  NewCommandChain(),
		rpnState:      &RPNState{vars: vars, calculator: NewRPNCalculator(rpnCalc)},
	}
	
	// First execute 2 to put it on the stack
	defaultExecutor(rpl, "2")
	
	// Then use x =: to assign the top of stack to variable x
	defaultExecutor(rpl, "x =:")
	val, exists := vars.GetVariable("x")
	if !exists {
		t.Errorf("Variable x should exist after x =:")
	}
	if val != 2 {
		t.Errorf("Variable x = %v, want 2", val)
	}
}

// TestExecutorWithExactUserScenario tests the exact user scenario: 2 then x =:
func TestExecutorWithExactUserScenario(t *testing.T) {
	vars := rpn.NewVariables()
	rpnCalc := rpn.NewRPN(vars)
	rpl := &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  NewCommandChain(),
		rpnState:      &RPNState{vars: vars, calculator: NewRPNCalculator(rpnCalc)},
	}
	
	// This test replicates the exact user interaction:
	// > 2
	// > x =:
	// The variable should be assigned the value 2
	
	defaultExecutor(rpl, "2")
	
	// Verify stack has 2
	// (can't directly check stack without exposing it, but next command will fail if stack is empty)
	
	defaultExecutor(rpl, "x =:")
	val, exists := vars.GetVariable("x")
	if !exists {
		t.Errorf("Variable x should exist after x =:")
	}
	if val != 2 {
		t.Errorf("Variable x = %v, want 2", val)
	}
}

// TestExecutorWithExactUserScenarioWithOutput tests that x =: assigns and shows result
func TestExecutorWithExactUserScenarioWithOutput(t *testing.T) {
	vars := rpn.NewVariables()
	rpnCalc := rpn.NewRPN(vars)
	rpl := &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  NewCommandChain(),
		rpnState:      &RPNState{vars: vars, calculator: NewRPNCalculator(rpnCalc)},
	}
	
	// Clear any previous state
	defaultExecutor(rpl, "rpn clear")
	
	// Put 2 on stack
	defaultExecutor(rpl, "2")
	_, _ = rpnCalc.ResultStack([]string{})
	
	// Assign to x =:
	result, err := rpnCalc.ParseAndEvaluate("x =:")
	t.Logf("ParseAndEvaluate('x =:') returned result=%q, err=%v", result, err)
	
	val, exists := vars.GetVariable("x")
	if !exists {
		t.Errorf("Variable x should exist after x =:")
	}
	if val != 2 {
		t.Errorf("Variable x = %v, want 2", val)
	}
}

// TestExecutorWithExactUserScenarioDirect simulates REPL input flow
func TestExecutorWithExactUserScenarioDirect(t *testing.T) {
	vars := rpn.NewVariables()
	rpnCalc := rpn.NewRPN(vars)
	rpl := &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  NewCommandChain(),
		rpnState:      &RPNState{vars: vars, calculator: NewRPNCalculator(rpnCalc)},
	}
	
	// Clear any previous state
	defaultExecutor(rpl, "rpn clear")
	
	// Simulate typing "2" in REPL
	defaultExecutor(rpl, "2")
	
	// Simulate typing "x =:" in REPL
	defaultExecutor(rpl, "x =:")
	
	// Verify variable was set
	val, exists := vars.GetVariable("x")
	if !exists {
		t.Errorf("Variable x should exist after x =:")
	}
	if val != 2 {
		t.Errorf("Variable x = %v, want 2", val)
	}
}

func TestExecutorWithUnknownCommand(t *testing.T) {
	repl := createTestREPL()
	// Test that unknown commands are handled by the error handler
	defaultExecutor(repl, "completelyunknowncommand123")
}

func TestDefaultExecutorCodePaths(t *testing.T) {
	// Test all code paths in defaultExecutor
	repl := createTestREPL()
	
	// Path 1: Empty input
	defaultExecutor(repl, "")
	
	// Path 2: Built-in command with error (clear should not error but let's verify)
	defaultExecutor(repl, "clear")
	
	// Path 3: Built-in command with output (help returns help text)
	defaultExecutor(repl, "help")
	
	// Path 4: Unknown command (error handler returns handled=false, err!=nil)
	defaultExecutor(repl, "completelyunknowncommand123")
	
	// Path 5: Whitespace only (trimmed to empty, returns early)
	defaultExecutor(repl, "   ")
}
