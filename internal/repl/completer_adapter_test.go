// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"strings"
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
		// Partial matches (single match, commonLen = shared prefix length)
		{
			name:       "partial match he",
			line:       []rune("he"),
			pos:        2,
			wantLen:    1,
			wantMinLen: 2, // "he" shared, suffix "lp"
		},
		{
			name:       "partial match cl",
			line:       []rune("cl"),
			pos:        2,
			wantLen:    1,
			wantMinLen: 2, // "cl" shared, suffix "ear"
		},
		{
			name:       "partial match q",
			line:       []rune("q"),
			pos:        1,
			wantLen:    1,
			wantMinLen: 1, // "q" shared, suffix "uit"
		},
		{
			name:       "partial match rp",
			line:       []rune("rp"),
			pos:        2,
			wantLen:    1,
			wantMinLen: 2, // "rp" shared, suffix "n"
		},
		// Multiple matches
		{
			name:       "partial match c matches calc and clear",
			line:       []rune("c"),
			pos:        1,
			wantLen:    2,
			wantMinLen: 1, // "c" shared, suffixes "alc" and "lear"
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
		// Case insensitive — case mismatch means no shared prefix
		{
			name:       "uppercase HELP",
			line:       []rune("HELP"),
			pos:        4,
			wantLen:    1,
			wantMinLen: 0, // H vs h: no shared prefix
		},
		{
			name:       "uppercase CLEAR",
			line:       []rune("CLEAR"),
			pos:        5,
			wantLen:    1,
			wantMinLen: 0, // C vs c: no shared prefix
		},
		{
			name:       "mixed case HeLp",
			line:       []rune("HeLp"),
			pos:        4,
			wantLen:    1,
			wantMinLen: 0, // H vs h: no shared prefix
		},
		{
			name:       "uppercase Q matches quit",
			line:       []rune("Q"),
			pos:        1,
			wantLen:    1,
			wantMinLen: 0, // Q vs q: no shared prefix
		},
		{
			name:       "uppercase C matches calc and clear",
			line:       []rune("C"),
			pos:        1,
			wantLen:    2,
			wantMinLen: 0, // C vs c: no shared prefix
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
			wantMinLen: 2, // "he" shared, suffix "lp"
		},
		{
			name:       "single char h",
			line:       []rune("h"),
			pos:        1,
			wantLen:    1,
			wantMinLen: 1, // "h" shared, suffix "elp"
		},
		{
			name:       "single char e matches exit",
			line:       []rune("e"),
			pos:        1,
			wantLen:    1,
			wantMinLen: 1, // "e" shared, suffix "xit"
		},
		{
			name:       "single char r matches rpn and rat",
			line:       []rune("r"),
			pos:        1,
			wantLen:    2,
			wantMinLen: 1, // "r" shared, suffixes "pn" and "at"
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

			// Verify all matches are actual commands (prefix + suffix)
			prefixWords := strings.Fields(string(tt.line[:tt.pos]))
			var prefix string
			if len(prefixWords) > 0 {
				prefix = prefixWords[len(prefixWords)-1]
			}
			for _, m := range matches {
				matchStr := prefix + string(m)
				// For case-insensitive matches, also check case-folded
				found := false
				for _, cmd := range commands {
					if cmd == matchStr || strings.EqualFold(cmd, matchStr) {
						found = true
						break
					}
					// Case-insensitive: the suffix itself (lowercased) should be a valid command
					if commonLen == 0 && strings.EqualFold(string(m), cmd) {
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

	// Suffixes: prefix is "c", so matches are "alc" and "lear"
	match0 := "c" + string(matches[0])
	match1 := "c" + string(matches[1])

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
