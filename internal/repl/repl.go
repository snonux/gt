// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"fmt"
	"strings"

	"codeberg.org/snonux/gt/internal/rpn"

	"github.com/c-bata/go-prompt"
)

// RPNState holds the state for RPN (Reverse Polish Notation) operations in the REPL.
// It maintains a variable store and RPN calculator instance.
//
// Note: This struct should never be copied - use pointer receivers only.
type RPNState struct {
	vars    rpn.VariableStore
	rpnCalc *rpn.RPN
}

// REPL manages the interactive command-line interface for the percentage calculator.
// It provides an interactive prompt with history, tab-completion, signal handling,
// and command processing through a chain of responsibility pattern.
//
// The REPL integrates various components:
//   - TTYChecker: validates stdin is a terminal
//   - HistoryManager: manages command history persistence
//   - SignalHandler: handles SIGINT (Ctrl+C)
//   - commandChain: processes commands via chain of responsibility
//   - rpnState: provides RPN state for calculations
type REPL struct {
	ttyChecker    *TTYChecker
	historyMgr    *HistoryManager
	signalHandler *SignalHandler
	prompt        *prompt.Prompt
	commandChain  CommandHandler
	rpnState      *RPNState
}

// NewREPL creates a new REPL instance with default components.
// If executor is nil, it uses defaultExecutor which processes input through commandChain.
// If completer is nil, it uses defaultCompleter which provides built-in command suggestions.
//
// The executor function is called for each non-empty input line.
// The completer function provides tab-completion suggestions for the prompt.
func NewREPL(executor func(string), completer func(prompt.Document) []prompt.Suggest) *REPL {
	// Initialize RPN state via dependency injection
	vars := rpn.NewVariables()
	rpnState := &RPNState{
		vars:    vars,
		rpnCalc: rpn.NewRPN(vars),
	}

	repl := &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  NewCommandChain(),
		rpnState:      rpnState,
	}

	// Set up executor - if nil, use default
	execFn := executor
	if execFn == nil {
		execFn = func(input string) {
			defaultExecutor(repl, input)
		}
	}

	// Set up completer - if nil, use default
	completerFn := completer
	if completerFn == nil {
		completerFn = func(d prompt.Document) []prompt.Suggest {
			return defaultCompleter(repl, d)
		}
	}

	// Load history from file
	history := repl.historyMgr.Load()

	// Build the prompt
	repl.prompt = NewPromptBuilder().
		SetTitle("gt - Percentage Calculator").
		SetPrefix("> ").
		SetLivePrefix(func() (string, bool) { return "> ", true }).
		SetExecutor(execFn).
		SetCompleter(completerFn).
		SetHistory(history).
		Build()

	return repl
}

// Run starts the REPL and blocks until it exits.
func (r *REPL) Run() error {
	// Check if stdin is a TTY
	if err := r.ttyChecker.EnsureTTY(); err != nil {
		return err
	}

	// Start signal handler
	r.signalHandler.Start(func() {
		fmt.Println("\nUse 'quit' or 'exit' to exit, or Ctrl+D")
	})

	// Run the prompt
	r.prompt.Run()

	return nil
}

// defaultExecutor is the default executor function used when no custom executor is provided.
// It processes input through the command chain of responsibility pattern.
// It includes panic recovery to gracefully handle unexpected errors during command execution.
//
// Input processing:
//   - Trims whitespace from input
//   - Skips empty input
//   - Routes to commandChain for processing
//   - Displays output and errors appropriately
//   - Adds handled commands to history
func defaultExecutor(r *REPL, input string) {
	// Add panic recovery for better resilience
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Printf("Error: Unexpected error occurred: %v\n", rec)
			fmt.Println("Please try a different expression or command.")
		}
	}()

	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	// Use chain of responsibility pattern to handle the command
	output, handled, err := r.commandChain.Handle(r, input)

	if handled {
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		if output != "" {
			fmt.Println(output)
		}
		// Don't add handled commands to history
		return
	}

	// Not handled by any handler in the chain
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// defaultCompleter is the default completer function used when no custom completer is provided.
// It provides tab-completion suggestions for built-in REPL commands.
// Suggestions are case-insensitive and include descriptions.
//
// d: the current prompt.Document containing cursor position and text
// Returns a slice of prompt.Suggest for matching built-in commands
func defaultCompleter(r *REPL, d prompt.Document) []prompt.Suggest {
	text := d.GetWordBeforeCursor()
	if text == "" {
		return nil
	}

	var suggestions []prompt.Suggest
	for _, cmd := range Commands() {
		if strings.HasPrefix(strings.ToLower(cmd), strings.ToLower(text)) {
			suggestions = append(suggestions, prompt.Suggest{
				Text:        cmd,
				Description: r.defaultGetCommandDescription(cmd),
			})
		}
	}
	return suggestions
}

// defaultGetCommandDescription returns the description for a built-in command.
// It's used by the default completer to provide helpful descriptions during tab-completion.
//
// cmd: the built-in command name (e.g., "help", "clear", "quit")
// Returns the description string for the command, or empty string if not found
func (r *REPL) defaultGetCommandDescription(cmd string) string {
	return getCommandDescription(cmd)
}

// getCommandDescription returns the description for a built-in command.
// This is a package-level function that provides a single source of truth
// for command descriptions, used by defaultGetCommandDescription.
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

// RunREPL starts the interactive REPL with default components.
// This is a convenience wrapper around NewREPL(nil, nil).Run().
// It's typically used when the standard REPL behavior is sufficient.
//
// Returns an error if the REPL cannot start (e.g., stdin is not a TTY)
func RunREPL() error {
	repl := NewREPL(nil, nil)
	return repl.Run()
}
