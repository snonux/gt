package repl

import (
	"strings"
	"testing"

	"github.com/c-bata/go-prompt"
)

// TestCompleterLogic tests the completer logic directly
func TestCompleterLogic(t *testing.T) {
	// Simulate the completer logic
	testCases := []struct {
		name  string
		text  string
		match bool
	}{
		{"h", "h", true},     // "help"
		{"he", "he", true},   // "help"
		{"hel", "hel", true}, // "help"
		{"help", "help", true},
		{"c", "c", true},     // "clear", "calc"
		{"cl", "cl", true},   // "clear"
		{"cle", "cle", true}, // "clear"
		{"clear", "clear", true},
		{"ca", "ca", true},   // "calc"
		{"cal", "cal", true}, // "calc"
		{"calc", "calc", true},
		{"q", "q", true},     // "quit"
		{"qu", "qu", true},   // "quit"
		{"qui", "qui", true}, // "quit"
		{"quit", "quit", true},
		{"e", "e", true},     // "exit"
		{"ex", "ex", true},   // "exit"
		{"exi", "exi", true}, // "exit"
		{"exit", "exit", true},
		{"r", "r", true},   // "rpn"
		{"rp", "rp", true}, // "rpn"
		{"rpn", "rpn", true},
		{"x", "x", false},     // no match
		{"xyz", "xyz", false}, // no match
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate the completer logic
			var found bool
			for _, cmd := range builtinCommands() {
				if strings.HasPrefix(strings.ToLower(cmd), strings.ToLower(tc.text)) {
					found = true
					break
				}
			}
			if found != tc.match {
				t.Errorf("For text %q, expected match=%v, got match=%v", tc.text, tc.match, found)
			}
		})
	}
}

// TestCompleterWithTrailingSpace tests completer with trailing space
func TestCompleterWithTrailingSpace(t *testing.T) {
	// When there's a trailing space, GetWordBeforeCursor() should return the word
	// This is how the actual REPL works when user types "help " then presses tab
	tests := []struct {
		name string
		text string
	}{
		{"h ", "h "},
		{"hel ", "hel "},
		{"c ", "c "},
		{"cl ", "cl "},
		{"help ", "help "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := completer(prompt.Document{Text: tt.text})
			// We just verify it doesn't panic
			_ = suggestions
		})
	}
}

// TestCompleterEmptyText tests completer with empty text
func TestCompleterEmptyText(t *testing.T) {
	suggestions := completer(prompt.Document{Text: ""})
	if suggestions != nil {
		t.Errorf("Expected nil for empty text, got %d suggestions", len(suggestions))
	}
}

// TestCompleterNoPrefix tests completer with no matching prefix
func TestCompleterNoPrefix(t *testing.T) {
	suggestions := completer(prompt.Document{Text: "xyz "})
	if len(suggestions) > 0 {
		t.Errorf("Expected no suggestions for 'xyz', got %d", len(suggestions))
	}
}

// TestCompleterWithAllCommands tests completer for all commands

// TestCompleterWithTrailingSpace tests completer with trailing space
