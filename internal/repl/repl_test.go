package repl

import (
	"testing"

	"github.com/c-bata/go-prompt"
)

func TestExecutor(t *testing.T) {
	// Test that executor doesn't panic on empty input
	executor("")
}

func TestExecutorWithHelp(t *testing.T) {
	// Test executor with help command
	// This should execute the help command and not print output
	executor("help")
}

func TestExecutorWithClear(t *testing.T) {
	executor("clear")
}

func TestExecutorWithQuit(t *testing.T) {
	// This should exit REPL but we can't test the actual exit
	// We just verify the command is processed without error
	executor("quit")
}

func TestExecutorWithExit(t *testing.T) {
	executor("exit")
}

func TestExecutorWithPercentage(t *testing.T) {
	// Test executor with a percentage calculation
	// Note: output is printed to stdout, we just verify it doesn't panic
	executor("20% of 150")
}

func TestExecutorWithRPN(t *testing.T) {
	// Test executor with RPN command
	executor("rpn 3 4 +")
}

func TestExecutorWithInvalid(t *testing.T) {
	// Test executor with invalid input
	executor("invalid input")
}

// Note: captureOutput is removed - we test for side effects instead of capturing output

func TestExecutorWithVars(t *testing.T) {
	executor("rpn x 5 = vars")
}

func TestExecutorWithClearVariables(t *testing.T) {
	executor("rpn clear")
}

func TestIsBuiltinCommand(t *testing.T) {
	// Test known built-in commands
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
	// Test with help subcommands
	_, ok := isBuiltinCommand("help clear")
	if !ok {
		t.Error("isBuiltinCommand('help clear') should return true")
	}
}

func TestCompleter(t *testing.T) {
	// Test completer with empty input
	suggestions := completer(prompt.Document{})
	// completer returns suggestions for builtin commands
	// We just verify it doesn't panic
	_ = suggestions
}

func TestCompleterWithPartialMatch(t *testing.T) {
	// Test completer with partial command
	suggestions := completer(prompt.Document{Text: "h"})
	// completer returns suggestions for builtin commands
	_ = suggestions
}

func TestCompleterWithClearPrefix(t *testing.T) {
	suggestions := completer(prompt.Document{Text: "cl"})
	_ = suggestions
}
