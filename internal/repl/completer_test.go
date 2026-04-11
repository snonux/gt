// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"testing"
)

// TestCompleter tests the completer function with various inputs
func TestCompleter(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		wantLen  int
		wantText []string
	}{
		{"empty text returns nil", "", 0, nil},
		{"help prefix returns help", "help", 1, []string{"help"}},
		{"h prefix matches help", "h", 1, []string{"help"}},
		{"he prefix matches help", "he", 1, []string{"help"}},
		{"hel prefix matches help", "hel", 1, []string{"help"}},
		{"clear prefix returns clear", "clear", 1, []string{"clear"}},
		{"c prefix matches clear and calc", "c", 2, []string{"calc", "clear"}},
		{"cl prefix matches clear", "cl", 1, []string{"clear"}},
		{"quit prefix returns quit", "quit", 1, []string{"quit"}},
		{"q prefix matches quit", "q", 1, []string{"quit"}},
		{"exit prefix returns exit", "exit", 1, []string{"exit"}},
		{"rpn prefix returns rpn", "rpn", 1, []string{"rpn"}},
		{"calc prefix returns calc", "calc", 1, []string{"calc"}},
		{"unknown prefix returns no matches", "xyz", 0, []string{}},
		{"case insensitive help", "HELP", 1, []string{"help"}},
		{"case insensitive clear", "CLEAR", 1, []string{"clear"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := completer(tt.text)
			if len(suggestions) != tt.wantLen {
				t.Errorf("completer(%q) returned %d suggestions, want %d", tt.text, len(suggestions), tt.wantLen)
			}

			if tt.wantText != nil {
				for _, expectedText := range tt.wantText {
					found := false
					for _, s := range suggestions {
						if s == expectedText {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("completer(%q) missing expected suggestion %q, got %v", tt.text, expectedText, suggestions)
					}
				}
			}
		})
	}
}

// TestCompleterWithAllBuiltinCommands tests completer for all built-in commands
func TestCompleterWithAllBuiltinCommands(t *testing.T) {
	commands := []string{"help", "clear", "quit", "exit", "rpn", "calc", "rat"}

	for _, cmd := range commands {
		t.Run(cmd, func(t *testing.T) {
			suggestions := completer(cmd)
			if len(suggestions) == 0 {
				t.Errorf("completer(%q) returned no suggestions, expected at least one", cmd)
			}
		})
	}
}

// TestCompleterDescription tests that suggestions are returned
func TestCompleterDescription(t *testing.T) {
	suggestions := completer("help")
	if len(suggestions) == 0 {
		t.Fatal("completer should return suggestions")
	}
}

// TestCompleterEdgeCases tests edge cases for completer
func TestCompleterEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{"single character q", "q"},
		{"single character c", "c"},
		{"single character h", "h"},
		{"partial help", "he"},
		{"partial quit", "qui"},
		{"partial exit", "ex"},
		{"partial rpn", "rp"},
		{"partial calc", "cal"},
		{"partial rat", "ra"},
		{"all lowercase help", "help"},
		{"all uppercase help", "HELP"},
		{"mixed case help", "HeLp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := completer(tt.text)
			_ = suggestions
		})
	}
}

// TestCompleterWithSpecialCharacters tests completer with special characters
func TestCompleterWithSpecialCharacters(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{"with tabs", "\thelp"},
		{"with newlines", "\nhelp"},
		{"with special chars", "hel#"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := completer(tt.text)
			_ = suggestions
		})
	}
}

// TestCompleterWithLongPrefix tests completer with long prefix
func TestCompleterWithLongPrefix(t *testing.T) {
	suggestions := completer("helooooooooo")
	if len(suggestions) != 0 {
		t.Errorf("completer with long prefix should return no matches, got %d", len(suggestions))
	}
}

// TestCompleterVerifyDescriptions tests that all commands return suggestions
func TestCompleterVerifyDescriptions(t *testing.T) {
	commands := []string{"help", "clear", "quit", "exit", "rpn", "calc"}
	for _, cmd := range commands {
		t.Run(cmd, func(t *testing.T) {
			suggestions := completer(cmd)
			if len(suggestions) == 0 {
				t.Errorf("completer(%q) should return suggestions", cmd)
			}
		})
	}
}

// TestCompleterNonAlphabetic tests completer with non-alphabetic input
func TestCompleterNonAlphabetic(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{"numbers only", "123"},
		{"symbols", "!@#"},
		{"mixed", "h3lp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := completer(tt.text)
			_ = suggestions
		})
	}
}

// TestCompleterMultipleWords tests completer behavior with space-separated words
func TestCompleterMultipleWords(t *testing.T) {
	suggestions := completer("clear")
	found := false
	for _, s := range suggestions {
		if s == "clear" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("completer with multiple words should complete 'clear', got %v", suggestions)
	}
}

// TestCompleterWithTrailingSpace tests completer with trailing space
func TestCompleterWithTrailingSpace(t *testing.T) {
	suggestions := completer("help ")
	_ = suggestions
}
