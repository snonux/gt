// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"strings"
)

// completer provides auto-completion for built-in commands.
// It returns suggestions for commands that match the current word being typed.
// The matching is case-insensitive.
//
// This function is used by readline for tab completion.
//
// text: the current word being typed
// Returns a slice of strings for matching built-in commands
func completer(text string) []string {
	if text == "" {
		return nil
	}

	var suggestions []string
	for _, cmd := range Commands() {
		if strings.HasPrefix(strings.ToLower(cmd), strings.ToLower(text)) {
			suggestions = append(suggestions, cmd)
		}
	}
	return suggestions
}

// AutoCompleteAdapter adapts our completer function to the readline AutoCompleter interface
type AutoCompleteAdapter struct {
	commands []string
}

// NewAutoCompleter creates a readline auto-completer that uses the completer function.
func NewAutoCompleter() *AutoCompleteAdapter {
	return &AutoCompleteAdapter{
		commands: Commands(),
	}
}

// Do implements the readline.AutoCompleter interface.
// It returns matching command completions for the given line.
func (a *AutoCompleteAdapter) Do(line []rune, pos int) ([][]rune, int) {
	text := string(line[:pos])
	words := strings.Fields(text)
	if len(words) == 0 {
		var result [][]rune
		for _, cmd := range a.commands {
			result = append(result, []rune(cmd))
		}
		return result, 0
	}

	lastWord := words[len(words)-1]
	var matches [][]rune
	for _, cmd := range a.commands {
		if strings.HasPrefix(strings.ToLower(cmd), strings.ToLower(lastWord)) {
			matches = append(matches, []rune(cmd))
		}
	}

	// Find common prefix length
	minLen := len(lastWord)
	for _, m := range matches {
		compare := string(m)
		i := 0
		for i < len(lastWord) && i < len(compare) && lastWord[i] == compare[i] {
			i++
		}
		if i < minLen {
			minLen = i
		}
	}

	return matches, minLen - len(lastWord)
}
