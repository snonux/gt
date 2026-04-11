// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"strings"
	"testing"
)

// TestCompleterLogic tests the completer logic directly
func TestCompleterLogic(t *testing.T) {
	testCases := []struct {
		name  string
		text  string
		match bool
	}{
		{"h", "h", true},
		{"he", "he", true},
		{"hel", "hel", true},
		{"help", "help", true},
		{"c", "c", true},
		{"cl", "cl", true},
		{"cle", "cle", true},
		{"clear", "clear", true},
		{"ca", "ca", true},
		{"cal", "cal", true},
		{"calc", "calc", true},
		{"q", "q", true},
		{"qu", "qu", true},
		{"qui", "qui", true},
		{"quit", "quit", true},
		{"e", "e", true},
		{"ex", "ex", true},
		{"exi", "exi", true},
		{"exit", "exit", true},
		{"r", "r", true},
		{"rp", "rp", true},
		{"rpn", "rpn", true},
		{"x", "x", false},
		{"xyz", "xyz", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			suggestions := completer(tc.text)
			found := false
			for _, s := range suggestions {
				if strings.HasPrefix(strings.ToLower(s), strings.ToLower(tc.text)) {
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

// TestCompleterEmptyText tests completer with empty text
func TestCompleterEmptyText(t *testing.T) {
	suggestions := completer("")
	if suggestions != nil {
		t.Errorf("Expected nil for empty text, got %d suggestions", len(suggestions))
	}
}

// TestCompleterNoPrefix tests completer with no matching prefix
func TestCompleterNoPrefix(t *testing.T) {
	suggestions := completer("xyz ")
	if len(suggestions) > 0 {
		t.Errorf("Expected no suggestions for 'xyz', got %d", len(suggestions))
	}
}
