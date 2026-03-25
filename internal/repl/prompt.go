package repl

import (
	"github.com/c-bata/go-prompt"
)

// PromptBuilder constructs a prompt instance with the given configuration.
// It uses the builder pattern to configure all aspects of the prompt before calling Build.
type PromptBuilder struct {
	prefix     string
	title      string
	history    []string
	executor   func(string)
	completer  func(prompt.Document) []prompt.Suggest
	livePrefix func() (string, bool)
}

// NewPromptBuilder creates a new prompt builder with default values.
// Default values:
//   - prefix: "> "
//   - title: "gt - Percentage Calculator"
//   - executor: empty function
//   - completer: function that returns nil
//   - livePrefix: function that returns ("> ", true)
//
// Returns a new PromptBuilder instance
func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{
		prefix:     "> ",
		title:      "gt - Percentage Calculator",
		executor:   func(string) {},
		completer:  func(prompt.Document) []prompt.Suggest { return nil },
		livePrefix: func() (string, bool) { return "> ", true },
	}
}

// SetPrefix sets the prompt prefix string.
// This is the string displayed before each input line (default: "> ").
//
// prefix: the prefix string to display
// Returns the builder for method chaining
func (b *PromptBuilder) SetPrefix(prefix string) *PromptBuilder {
	b.prefix = prefix
	return b
}

// SetTitle sets the prompt title.
// This title is displayed in the terminal window/tab title.
//
// title: the title string to set
// Returns the builder for method chaining
func (b *PromptBuilder) SetTitle(title string) *PromptBuilder {
	b.title = title
	return b
}

// SetHistory sets the history for the prompt.
// The history is a slice of strings representing previously entered commands.
// This allows users to navigate through their command history using arrow keys.
//
// history: the slice of history entries
// Returns the builder for method chaining
func (b *PromptBuilder) SetHistory(history []string) *PromptBuilder {
	b.history = history
	return b
}

// SetExecutor sets the executor function for processing input.
// The executor is called for each non-empty input line after the user presses Enter.
//
// executor: the function to call with each input line
// Returns the builder for method chaining
func (b *PromptBuilder) SetExecutor(executor func(string)) *PromptBuilder {
	b.executor = executor
	return b
}

// SetCompleter sets the completer function for auto-completion.
// The completer is called when the user presses Tab to get suggestions.
//
// completer: the function to call for tab-completion suggestions
// Returns the builder for method chaining
func (b *PromptBuilder) SetCompleter(completer func(prompt.Document) []prompt.Suggest) *PromptBuilder {
	b.completer = completer
	return b
}

// SetLivePrefix sets the live prefix function.
// The live prefix is displayed on the left side of the current input line
// and can be used to show context-dependent information (e.g., multi-line input).
//
// livePrefix: the function that returns the current prefix string
// Returns the builder for method chaining
func (b *PromptBuilder) SetLivePrefix(livePrefix func() (string, bool)) *PromptBuilder {
	b.livePrefix = livePrefix
	return b
}

// Build creates and returns a new prompt instance with the configured options.
// After calling Build, the PromptBuilder should not be modified.
//
// Returns a new prompt.Prompt instance ready to use with prompt.Run()
func (b *PromptBuilder) Build() *prompt.Prompt {
	return prompt.New(
		b.executor,
		b.completer,
		prompt.OptionTitle(b.title),
		prompt.OptionPrefix(b.prefix),
		prompt.OptionLivePrefix(b.livePrefix),
		prompt.OptionHistory(b.history),
	)
}
