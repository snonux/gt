package repl

import (
	"strings"
	"testing"

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

func TestCompleter(t *testing.T) {
	suggestions := completer(prompt.Document{})
	_ = suggestions
}

func TestCompleterWithPartialMatch(t *testing.T) {
	// Use trailing space to ensure GetWordBeforeCursor() returns non-empty
	suggestions := completer(prompt.Document{Text: "h "})
	_ = suggestions
}

func TestCompleterWithClearPrefix(t *testing.T) {
	suggestions := completer(prompt.Document{Text: "cl "})
	_ = suggestions
}

func TestCompleterWithEmptyText(t *testing.T) {
	suggestions := completer(prompt.Document{Text: ""})
	if suggestions != nil {
		t.Errorf("completer with empty text should return nil, got %d suggestions", len(suggestions))
	}
}

func TestCompleterWithAllCommands(t *testing.T) {
	allCommands := []string{"help", "clear", "quit", "exit", "rpn", "calc"}
	for _, cmd := range allCommands {
		t.Run(cmd, func(t *testing.T) {
			suggestions := completer(prompt.Document{Text: cmd + " "})
			_ = suggestions
		})
	}
}

func TestCompleterWithNoPrefix(t *testing.T) {
	suggestions := completer(prompt.Document{Text: "xyz "})
	if len(suggestions) > 0 {
		t.Errorf("completer(%q) should return no suggestions, got %d", "xyz", len(suggestions))
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

func TestCompleterEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		doc  prompt.Document
	}{
		{"single character q", prompt.Document{Text: "q"}},
		{"single character e", prompt.Document{Text: "e"}},
		{"single character r", prompt.Document{Text: "r"}},
		{"partial help", prompt.Document{Text: "he"}},
		{"partial quit", prompt.Document{Text: "qui"}},
		{"partial exit", prompt.Document{Text: "ex"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := completer(tt.doc)
			_ = suggestions
		})
	}
}

func TestIsBuiltinCommandWithSubcommandHelp(t *testing.T) {
	_, ok := isBuiltinCommand("help")
	if !ok {
		t.Error("isBuiltinCommand('help') should return true")
	}
}
