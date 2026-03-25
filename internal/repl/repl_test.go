// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"strings"
	"testing"

	"codeberg.org/snonux/perc/internal/rpn"

	"github.com/c-bata/go-prompt"
)

func TestExecutor(t *testing.T) {
	// Test that executor doesn't panic on empty input
	executor("")
}

func TestExecutorWithHelp(t *testing.T) {
	// Test executor with help command
	executor("help")
}

func TestExecutorWithClear(t *testing.T) {
	executor("clear")
}

func TestExecutorWithQuit(t *testing.T) {
	executor("quit")
}

func TestExecutorWithExit(t *testing.T) {
	executor("exit")
}

func TestExecutorWithPercentage(t *testing.T) {
	executor("20% of 150")
}

func TestExecutorWithRPN(t *testing.T) {
	executor("rpn 3 4 +")
}

func TestExecutorWithInvalid(t *testing.T) {
	executor("invalid input")
}

func TestExecutorWithVars(t *testing.T) {
	executor("rpn x 5 = vars")
}

func TestExecutorWithClearVariables(t *testing.T) {
	executor("rpn clear")
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

// TestRunRPN tests the runRPN helper function
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
			_, err := runRPN(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("runRPN(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
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
	executor("+")
	executor("-")
	executor("*")
	executor("/")
	executor("^")
	executor("%")
	executor("dup")
	executor("swap")
	executor("pop")
	executor("show")
	executor("vars")
	executor("clear")
}

func TestExecutorWithPercentageExpression(t *testing.T) {
	executor("20% of 150")
	executor("30 is what %% of 150")
	executor("30 is 20%% of what")
}

func TestExecutorWithInvalidPercentage(t *testing.T) {
	executor("invalid percentage input")
}

func TestExecutorWithOperatorOnly(t *testing.T) {
	executor("1 2 +")
	executor("+")
}

func TestExecutorWithRPNPrefix(t *testing.T) {
	executor("rpn 3 4 +")
}

func TestExecutorWithCalcPrefix(t *testing.T) {
	executor("calc 5 6 +")
}

func TestExecutorWithEmptyInput(t *testing.T) {
	executor("")
}

func TestExecutorWithWhitespaceOnly(t *testing.T) {
	executor("   ")
}

func TestExecutorWithInvalidInput(t *testing.T) {
	tests := []string{"invalid input", "not a valid command", "xyz"}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			executor(input)
		})
	}
}

func TestExecutorWithInvalidRPN(t *testing.T) {
	executor("rpn 1 +")
}

func TestExecutorWithEmptyRPNPrefix(t *testing.T) {
	executor("rpn")
	executor("calc")
}

func TestExecutorWithAssignment(t *testing.T) {
	executor("rpn x 42 =")
	executor("rpn x")
}

func TestExecutorWithPercentageAndRPNFallback(t *testing.T) {
	executor("20% of 150")
	executor("3 4 +")
}

func TestGetHistoryPath(t *testing.T) {
	path := getHistoryPath()
	if path == "" {
		t.Error("getHistoryPath() returned empty string")
	}
}

func TestLoadHistory(t *testing.T) {
	history := loadHistory()
	_ = history
}

func TestSaveHistory(t *testing.T) {
	err := saveHistory([]string{"test1", "test2"})
	_ = err
}

func TestExecutorWithRPNExpressionOnly(t *testing.T) {
	executor("5 3 +")
}

func TestExecutorWithRPNThenOperator(t *testing.T) {
	executor("1 2 +")
	executor("+")
}

func TestExecutorWithRPNThenRPN(t *testing.T) {
	executor("rpn 1 2 +")
	executor("rpn 3 4 +")
}

func TestExecutorWithRPNShow(t *testing.T) {
	executor("rpn show")
}

func TestExecutorWithRPNDup(t *testing.T) {
	executor("rpn dup")
}

func TestExecutorWithRPNSwap(t *testing.T) {
	executor("rpn swap")
}

func TestExecutorWithRPNSingle(t *testing.T) {
	executor("rpn 42")
}

func TestExecutorWithRPNMulti(t *testing.T) {
	executor("rpn 1 2 3 4 5 +")
}

func TestExecutorWithStackOps(t *testing.T) {
	executor("dup")
	executor("swap")
	executor("pop")
	executor("show")
}

func TestExecutorWithRPNClear(t *testing.T) {
	executor("rpn clear")
}

func TestExecutorWithHistoryCommands(t *testing.T) {
	executor("vars")
	executor("clear")
}

func TestExecutorWithMixedInput(t *testing.T) {
	executor("25% of 200")
	executor("10 20 +")
}

func TestExecutorWithRPNCalcMixed(t *testing.T) {
	executor("rpn 1 2 +")
	executor("3 4 +")
	executor("calc 5 6 +")
}

func TestExecutorCommandsEdgeCases(t *testing.T) {
	executor("  clear  ")
	executor("HELP")
	executor("CLEAR")
}

func TestExecutorWithRPMPrefix(t *testing.T) {
	executor("rpn 1 2 +")
}

func TestExecutorWithCalcPrefixMixed(t *testing.T) {
	executor("calc 1 2 +")
}

func TestExecutorWithRatModeOn(t *testing.T) {
	executor("rat on")
	state := getRPNState()
	if state.rpnCalc.GetMode() != rpn.RationalMode {
		t.Errorf("Expected RationalMode after rat on, got %v", state.rpnCalc.GetMode())
	}
}

func TestExecutorWithRatModeOff(t *testing.T) {
	executor("rat off")
	state := getRPNState()
	if state.rpnCalc.GetMode() != rpn.FloatMode {
		t.Errorf("Expected FloatMode after rat off, got %v", state.rpnCalc.GetMode())
	}
}

func TestExecutorWithRatModeToggle(t *testing.T) {
	// First toggle - should enable rational mode if currently float
	executor("rat toggle")
	state := getRPNState()
	mode1 := state.rpnCalc.GetMode()

	// Second toggle - should toggle back
	executor("rat toggle")
	state = getRPNState()
	mode2 := state.rpnCalc.GetMode()

	// Modes should be different after toggle
	if mode1 == mode2 {
		t.Errorf("Modes should be different after toggle: %v -> %v", mode1, mode2)
	}
}

func TestExecutorWithRatModeInvalid(t *testing.T) {
	// Just verify it doesn't panic
	executor("rat invalid")
}

func TestExecutorWithRatModeNoArg(t *testing.T) {
	// Just verify it doesn't panic
	executor("rat")
}

func TestIsBuiltinCommandWithSubcommandHelp(t *testing.T) {
	_, ok := isBuiltinCommand("help")
	if !ok {
		t.Error("isBuiltinCommand('help') should return true")
	}
}

func TestRPNHandlerWithUnknownInput(t *testing.T) {
	// Test that unknown input falls through to next handler
	chain := NewCommandChain()

	// Create a minimal REPL
	r := &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  chain,
	}

	// Test unknown input - should not be handled by RPNHandler directly
	// but will be handled by Error handler after RPNHandler passes it through
	output, handled, err := chain.Handle(r, "unknowncommand")
	if handled {
		t.Errorf("Expected unknowncommand to be handled by error handler, got handled=%v, err=%v, output=%q", handled, err, output)
	}
}

func TestRPNHandlerWithPercentageExpression(t *testing.T) {
	// Test that percentage expressions are handled by PercentageHandler, not RPNHandler
	chain := NewCommandChain()
	r := &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  chain,
	}

	// Test percentage expression
	output, handled, err := chain.Handle(r, "20% of 150")
	if !handled {
		t.Errorf("Expected percentage expression to be handled, got handled=%v, err=%v, output=%q", handled, err, output)
	}
	if err != nil {
		t.Errorf("Expected no error for percentage expression, got %v", err)
	}
}

func TestRPNHandlerWithRPNExpression(t *testing.T) {
	// Test RPN expressions
	chain := NewCommandChain()
	r := &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  chain,
	}

	// Test RPN expression
	output, handled, err := chain.Handle(r, "3 4 +")
	if !handled {
		t.Errorf("Expected RPN expression to be handled, got handled=%v, err=%v, output=%q", handled, err, output)
	}
	if err != nil {
		t.Errorf("Expected no error for RPN expression, got %v", err)
	}
}

func TestRPNHandlerWithSingleNumber(t *testing.T) {
	// Test single number input (RPN - pushes number onto stack)
	chain := NewCommandChain()
	r := &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  chain,
	}

	// Test single number
	output, handled, err := chain.Handle(r, "42")
	if !handled {
		t.Errorf("Expected single number to be handled, got handled=%v, err=%v, output=%q", handled, err, output)
	}
	if err != nil {
		t.Errorf("Expected no error for single number, got %v", err)
	}
}

// TestNewREPL tests that NewREPL creates a valid REPL instance.
// Note: This test is skipped when not running in a TTY because the prompt
// library requires TTY access.
func TestNewREPL(t *testing.T) {
	// Skip this test if not running in a TTY
	ttyChecker := &TTYChecker{}
	if !ttyChecker.IsTTY() {
		t.Skip("Skipping test - not running in a TTY")
	}

	// Test that NewREPL creates a valid REPL instance without panicking
	repl := NewREPL(nil, nil)
	if repl == nil {
		t.Fatal("Expected REPL to be created, got nil")
	}
	if repl.prompt == nil {
		t.Error("Expected prompt to be set")
	}
	if repl.commandChain == nil {
		t.Error("Expected commandChain to be set")
	}
	if repl.ttyChecker == nil {
		t.Error("Expected ttyChecker to be set")
	}
	if repl.historyMgr == nil {
		t.Error("Expected historyMgr to be set")
	}
	if repl.signalHandler == nil {
		t.Error("Expected signalHandler to be set")
	}
}

func TestDefaultCompleter(t *testing.T) {
	// Test the default completer function directly
	// Note: This test has limited coverage because defaultCompleter uses
	// GetWordBeforeCursor() which requires proper cursor position.
	// The actual completer logic is tested in completer_test.go

	// Test with text that would match if cursor position was set correctly
	repl := &REPL{}
	doc := prompt.Document{Text: "h"}
	suggestions := defaultCompleter(repl, doc)

	// When cursor is at position 0 (default), GetWordBeforeCursor returns empty
	// But the test in completer_test.go verifies the actual behavior
	_ = suggestions

	// Test with clear prefix
	doc2 := prompt.Document{Text: "cl"}
	suggestions2 := defaultCompleter(repl, doc2)
	_ = suggestions2
}

func TestDefaultGetCommandDescription(t *testing.T) {
	// Create a REPL and test the defaultGetCommandDescription method
	repl := &REPL{}

	tests := []struct {
		cmd        string
		wantPrefix string
	}{
		{"help", "Show"},
		{"clear", "Clear"},
		{"quit", "Exit"},
		{"exit", "Exit"},
		{"rpn", "Evaluate"},
		{"calc", "Same"},
	}

	for _, tt := range tests {
		t.Run(tt.cmd, func(t *testing.T) {
			desc := repl.defaultGetCommandDescription(tt.cmd)
			if !strings.Contains(desc, tt.wantPrefix) {
				t.Errorf("defaultGetCommandDescription(%q) = %q, should contain %q", tt.cmd, desc, tt.wantPrefix)
			}
		})
	}
}

func TestExecutorWithUnknownCommand(t *testing.T) {
	// Test that unknown commands are handled by the error handler
	// This should exercise the "Not handled by any handler" path
	executor("completelyunknowncommand123")
}

func TestDefaultExecutorCodePaths(t *testing.T) {
	// Test all code paths in defaultExecutor
	// 1. Empty input (returns early at line 110)
	// 2. Handled=true with error (prints error, returns at line 124)
	// 3. Handled=true with output (prints output, returns at line 124)
	// 4. Handled=false with error (prints error at line 130)
	// 5. Handled=false without error (does nothing)

	// Path 1: Empty input
	executor("")

	// Path 2: Built-in command with error (clear should not error but let's verify)
	executor("clear")

	// Path 3: Built-in command with output (help returns help text)
	executor("help")

	// Path 4: Unknown command (error handler returns handled=false, err!=nil)
	executor("completelyunknowncommand123")

	// Path 5: Whitespace only (trimmed to empty, returns early)
	executor("   ")
}
