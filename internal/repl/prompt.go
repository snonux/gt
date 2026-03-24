package repl

import (
	"github.com/c-bata/go-prompt"
)

// PromptBuilder constructs a prompt instance with the given configuration.
type PromptBuilder struct {
	prefix     string
	title      string
	history    []string
	executor   func(string)
	completer  func(prompt.Document) []prompt.Suggest
	livePrefix func() (string, bool)
}

// NewPromptBuilder creates a new prompt builder.
func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{
		prefix:     "> ",
		title:      "gt - Percentage Calculator",
		executor:   func(string) {},
		completer:  func(prompt.Document) []prompt.Suggest { return nil },
		livePrefix: func() (string, bool) { return "> ", true },
	}
}

// SetPrefix sets the prompt prefix.
func (b *PromptBuilder) SetPrefix(prefix string) *PromptBuilder {
	b.prefix = prefix
	return b
}

// SetTitle sets the prompt title.
func (b *PromptBuilder) SetTitle(title string) *PromptBuilder {
	b.title = title
	return b
}

// SetHistory sets the history for the prompt.
func (b *PromptBuilder) SetHistory(history []string) *PromptBuilder {
	b.history = history
	return b
}

// SetExecutor sets the executor function for processing input.
func (b *PromptBuilder) SetExecutor(executor func(string)) *PromptBuilder {
	b.executor = executor
	return b
}

// SetCompleter sets the completer function for auto-completion.
func (b *PromptBuilder) SetCompleter(completer func(prompt.Document) []prompt.Suggest) *PromptBuilder {
	b.completer = completer
	return b
}

// SetLivePrefix sets the live prefix function.
func (b *PromptBuilder) SetLivePrefix(livePrefix func() (string, bool)) *PromptBuilder {
	b.livePrefix = livePrefix
	return b
}

// Build creates and returns a new prompt instance.
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
