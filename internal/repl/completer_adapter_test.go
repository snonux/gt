// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"testing"
)

func TestNewAutoCompleter(t *testing.T) {
	adapter := NewAutoCompleter()
	if adapter == nil {
		t.Fatal("NewAutoCompleter returned nil")
	}
	if adapter.commands == nil {
		t.Fatal("AutoCompleteAdapter.commands is nil")
	}
	expectedCommands := Commands()
	if len(adapter.commands) != len(expectedCommands) {
		t.Errorf("commands count = %d, want %d", len(adapter.commands), len(expectedCommands))
	}
}

func TestAutoCompleteAdapterDo(t *testing.T) {
	adapter := NewAutoCompleter()
	commands := Commands()

	tests := []struct {
		name        string
		line        []rune
		pos         int
		wantLen     int
		wantMinLen  int
		description string
	}{
		// Empty / whitespace input returns all commands
		{
			name:       "empty input returns all commands",
			line:       []rune(""),
			pos:        0,
			wantLen:    len(commands),
			wantMinLen: 0,
		},
		{
			name:       "whitespace only returns all commands",
			line:       []rune("   "),
			pos:        3,
			wantLen:    len(commands),
			wantMinLen: 0,
		},
		// Exact matches — no completion offered (readline would append the word)
		{
			name:       "exact match help",
			line:       []rune("help"),
			pos:        4,
			wantLen:    0,
			wantMinLen: 0,
		},
		{
			name:       "exact match clear",
			line:       []rune("clear"),
			pos:        5,
			wantLen:    0,
			wantMinLen: 0,
		},
		// Partial matches (single match, commonLen always 0 since minLen capped at len(lastWord))
		{
			name:       "partial match he",
			line:       []rune("he"),
			pos:        2,
			wantLen:    1,
			wantMinLen: 0,
		},
		{
			name:       "partial match cl",
			line:       []rune("cl"),
			pos:        2,
			wantLen:    1,
			wantMinLen: 0,
		},
		{
			name:       "partial match q",
			line:       []rune("q"),
			pos:        1,
			wantLen:    1,
			wantMinLen: 0,
		},
		{
			name:       "partial match rp",
			line:       []rune("rp"),
			pos:        2,
			wantLen:    1,
			wantMinLen: 0,
		},
		// Multiple matches
		{
			name:       "partial match c matches calc and clear",
			line:       []rune("c"),
			pos:        1,
			wantLen:    2,
			wantMinLen: 0,
		},
		// No matches
		{
			name:       "no match xyz",
			line:       []rune("xyz"),
			pos:        3,
			wantLen:    0,
			wantMinLen: 0,
		},
		{
			name:       "no match numbers",
			line:       []rune("123"),
			pos:        3,
			wantLen:    0,
			wantMinLen: 0,
		},
		{
			name:       "no match symbols",
			line:       []rune("!@#"),
			pos:        3,
			wantLen:    0,
			wantMinLen: 0,
		},
		{
			name:       "no match too long prefix",
			line:       []rune("heloooooo"),
			pos:        9,
			wantLen:    0,
			wantMinLen: 0,
		},
		// Case insensitive
		{
			name:       "uppercase HELP",
			line:       []rune("HELP"),
			pos:        4,
			wantLen:    1,
			wantMinLen: -4, // HELP vs help: no case match in prefix calc
		},
		{
			name:       "uppercase CLEAR",
			line:       []rune("CLEAR"),
			pos:        5,
			wantLen:    1,
			wantMinLen: -5,
		},
		{
			name:       "mixed case HeLp",
			line:       []rune("HeLp"),
			pos:        4,
			wantLen:    1,
			wantMinLen: -4,
		},
		{
			name:       "uppercase Q matches quit",
			line:       []rune("Q"),
			pos:        1,
			wantLen:    1,
			wantMinLen: -1,
		},
		{
			name:       "uppercase C matches calc and clear",
			line:       []rune("C"),
			pos:        1,
			wantLen:    2,
			wantMinLen: -1,
		},
		// Multi-word input completes last word (exact match returns nothing)
		{
			name:       "multi word input completes last word",
			line:       []rune("rpn help"),
			pos:        8,
			wantLen:    0,
			wantMinLen: 0,
		},
		// Cursor position matters
		{
			name:       "cursor in middle of word",
			line:       []rune("help"),
			pos:        2,
			wantLen:    1,
			wantMinLen: 0,
		},
		// Single character prefixes
		{
			name:       "single char h",
			line:       []rune("h"),
			pos:        1,
			wantLen:    1,
			wantMinLen: 0,
		},
		{
			name:       "single char e matches exit",
			line:       []rune("e"),
			pos:        1,
			wantLen:    1,
			wantMinLen: 0,
		},
		{
			name:       "single char r matches rpn and rat",
			line:       []rune("r"),
			pos:        1,
			wantLen:    2,
			wantMinLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches, commonLen := adapter.Do(tt.line, tt.pos)

			if len(matches) != tt.wantLen {
				t.Errorf("got %d matches, want %d. matches: %v",
					len(matches), tt.wantLen, runeSliceToStringSlice(matches))
			}

			if commonLen != tt.wantMinLen {
				t.Errorf("common prefix len = %d, want %d", commonLen, tt.wantMinLen)
			}

			// Verify all matches are actual commands
			for _, m := range matches {
				matchStr := string(m)
				found := false
				for _, cmd := range commands {
					if cmd == matchStr {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("match %q is not a valid command", matchStr)
				}
			}
		})
	}
}

func TestAutoCompleteAdapterDoPreserveCommandOrder(t *testing.T) {
	adapter := NewAutoCompleter()
	// 'c' matches calc and clear; they should appear in Commands() order
	matches, _ := adapter.Do([]rune("c"), 1)
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}

	cmds := Commands()
	calcIdx, clearIdx := -1, -1
	for i, cmd := range cmds {
		if cmd == "calc" {
			calcIdx = i
		}
		if cmd == "clear" {
			clearIdx = i
		}
	}

	match0 := string(matches[0])
	match1 := string(matches[1])

	if calcIdx < clearIdx {
		if match0 != "calc" || match1 != "clear" {
			t.Errorf("expected [calc, clear], got [%s, %s]", match0, match1)
		}
	} else {
		if match0 != "clear" || match1 != "calc" {
			t.Errorf("expected [clear, calc], got [%s, %s]", match0, match1)
		}
	}
}

func TestAutoCompleteAdapterDoMultilineInput(t *testing.T) {
	adapter := NewAutoCompleter()

	// Tab-separated words — exact match "help" returns no completions
	matches, _ := adapter.Do([]rune("rpn\thelp"), 8)
	if len(matches) != 0 {
		t.Errorf("tab-separated 'rpn\\thelp' (exact match) should return 0, got %d: %v",
			len(matches), runeSliceToStringSlice(matches))
	}
}

// Helper function to convert [][]rune to []string for error messages
func runeSliceToStringSlice(runes [][]rune) []string {
	result := make([]string, len(runes))
	for i, r := range runes {
		result[i] = string(r)
	}
	return result
}
