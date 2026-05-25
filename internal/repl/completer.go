// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"strings"

	"github.com/chzyer/readline"
)

// completer provides auto-completion for built-in commands.
// It returns suggestions for commands that match the current word being typed.
// The matching is case-insensitive.
//
// The function is used in tests; readline tab completion uses AutoCompleteAdapter instead.
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

// Ensure AutoCompleteAdapter implements readline.AutoCompleter at compile time.
var _ readline.AutoCompleter = (*AutoCompleteAdapter)(nil)

// AutoCompleteAdapter implements the readline AutoCompleter interface,
// providing tab-completion suggestions for built-in commands.
// When the first word is "help", it completes with help topics instead.
type AutoCompleteAdapter struct {
	commands []string
}

// NewAutoCompleter creates a new AutoCompleteAdapter with the current list of built-in commands.
func NewAutoCompleter() *AutoCompleteAdapter {
	return &AutoCompleteAdapter{
		commands: Commands(),
	}
}

// Do implements the readline.AutoCompleter interface.
// It returns matching command completions for the given line.
// When the first word is "help" followed by a space, it offers help topic completions.
func (a *AutoCompleteAdapter) Do(line []rune, pos int) ([][]rune, int) {
	text := string(line[:pos])
	words := strings.Fields(text)
	if len(words) == 0 {
		return a.completeCommands("")
	}

	// If first word is "help" and user typed a space after it (or more words),
	// complete help topics. If first word is just "help" with no trailing space,
	// complete the command "help".
	if strings.ToLower(words[0]) == "help" && len(words) == 1 {
		// Check if there's a trailing space — means user wants topics
		if len(text) > 0 && text[len(text)-1] == ' ' {
			return a.completeHelpTopics("")
		}
		// No trailing space — just completing the command "help"
	}
	if strings.ToLower(words[0]) == "help" && len(words) > 1 {
		return a.completeHelpTopics(words[len(words)-1])
	}

	lastWord := words[len(words)-1]
	return a.completeCommands(lastWord)
}

// completeCommands returns matching command completions and prefix length.
func (a *AutoCompleteAdapter) completeCommands(lastWord string) ([][]rune, int) {
	var matches [][]rune
	for _, cmd := range a.commands {
		if strings.HasPrefix(strings.ToLower(cmd), strings.ToLower(lastWord)) {
			matches = append(matches, []rune(cmd))
		}
	}
	// If the word already exactly matches one command, don't offer it as a
	// completion — readline would append it, giving e.g. "helphelp".
	if len(matches) == 1 && string(matches[0]) == lastWord {
		return nil, 0
	}
	return a.withCommonPrefix(matches, lastWord)
}

// completeHelpTopics returns matching help topic completions and prefix length.
func (a *AutoCompleteAdapter) completeHelpTopics(lastWord string) ([][]rune, int) {
	var matches [][]rune
	topics := GetCompletionTopics()
	for _, topic := range topics {
		if strings.HasPrefix(strings.ToLower(topic), strings.ToLower(lastWord)) {
			matches = append(matches, []rune(topic))
		}
	}
	return a.withCommonPrefix(matches, lastWord)
}

// withCommonPrefix calculates the common prefix and returns suffixes
// as required by the readline AutoCompleter interface.
// Per the readline docs, candidates should be the characters AFTER the
// common prefix, and length is the common prefix length.
// Example: input "he" matches "help" => return [["lp"]], 2
func (a *AutoCompleteAdapter) withCommonPrefix(matches [][]rune, lastWord string) ([][]rune, int) {
	if len(matches) == 0 {
		return matches, 0
	}

	// Find the common prefix length across all matches.
	minLen := len(matches[0])
	for _, m := range matches[1:] {
		if len(m) < minLen {
			minLen = len(m)
		}
	}
	for i := 0; i < minLen; i++ {
		ch := matches[0][i]
		for _, m := range matches[1:] {
			if m[i] != ch {
				minLen = i
				break
			}
		}
		if minLen == i {
			break
		}
	}

	// Also cap by the number of characters that match the typed input.
	// Case-insensitive matches (e.g. "HELP" -> "help") have 0 shared chars.
	inputWord := []rune(lastWord)
	sharedWithInput := 0
	for i := 0; i < minLen && i < len(inputWord); i++ {
		// Check against first match (all matches share this prefix)
		if i < len(matches[0]) && matches[0][i] == inputWord[i] {
			sharedWithInput++
		} else {
			break
		}
	}
	if sharedWithInput < minLen {
		minLen = sharedWithInput
	}

	// Return only the suffixes (characters after the common prefix).
	suffixes := make([][]rune, len(matches))
	for i, m := range matches {
		suffixes[i] = m[minLen:]
	}

	return suffixes, minLen
}
