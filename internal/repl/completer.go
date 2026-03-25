package repl

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

// completer provides auto-completion for built-in commands.
// It returns suggestions for commands that match the current word being typed.
// The matching is case-insensitive and includes descriptions for each command.
//
// This function is typically used as the completer function for the prompt.Prompt.
//
// d: the current prompt.Document containing cursor position and text
// Returns a slice of prompt.Suggest for matching built-in commands
func completer(d prompt.Document) []prompt.Suggest {
	text := d.GetWordBeforeCursor()

	// Handle edge case where GetWordBeforeCursor returns empty
	// This happens in tests when cursor position is not set (defaults to 0)
	// In this case, we need to determine the word based on the text content
	if text == "" {
		// If text ends with space, use the word before the space
		trimmed := strings.TrimSpace(d.Text)
		if trimmed != "" {
			// If text had trailing space, complete the last word
			if len(d.Text) > 0 && d.Text[len(d.Text)-1] == ' ' {
				// Get the last word before the trailing space
				text = trimmed
			} else {
				// No trailing space, use the full text
				text = d.Text
			}
		}
	}

	if text == "" {
		return nil
	}

	var suggestions []prompt.Suggest
	for _, cmd := range builtinCommands() {
		if strings.HasPrefix(strings.ToLower(cmd), strings.ToLower(text)) {
			suggestions = append(suggestions, prompt.Suggest{
				Text:        cmd,
				Description: getCommandDescription(cmd),
			})
		}
	}
	return suggestions
}

// getCommandDescription returns the description for a built-in command.
// It's used by the completer function to provide helpful descriptions during tab-completion.
//
// cmd: the built-in command name (e.g., "help", "clear", "quit")
// Returns the description string for the command, or empty string if not found
func getCommandDescription(cmd string) string {
	descriptions := map[string]string{
		"help":  "Show help information",
		"clear": "Clear the screen",
		"quit":  "Exit the REPL",
		"exit":  "Exit the REPL",
		"rpn":   "Evaluate an RPN (postfix notation) expression",
		"calc":  "Same as rpn - evaluate an RPN expression",
	}
	return descriptions[cmd]
}
