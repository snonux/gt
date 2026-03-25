// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"fmt"
	"strings"
	"sync"

	"codeberg.org/snonux/perc/internal/rpn"

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

// executorREPL holds the REPL instance created by the executor function.
// This is used for backward compatibility with tests that need to access RPN state
// after calling executor(). It's not part of the main REPL architecture.
// Thread safety: Use executorREPLOnce for lazy initialization and executorREPLMu for access.
var executorREPL *REPL
var executorREPLOnce sync.Once
var executorREPLMu sync.Mutex

// ResetExecutorREPL resets the executorREPL for clean test isolation.
// This should be called between tests that use executor() and getRPNState()
// to ensure each test starts with a fresh RPN state.
//
// Note: This function is intended for test use only and should not be used
// in production code. For production use, create new REPL instances with NewREPL().
func ResetExecutorREPL() {
	executorREPLMu.Lock()
	defer executorREPLMu.Unlock()
	executorREPL = nil
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
	for _, cmd := range builtinCommands() {
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

// executor runs a calculation command and returns the result.
// This is a package-level wrapper for backward compatibility and testing.
// It creates a minimal REPL instance without building a prompt, allowing
// calculation execution in non-interactive contexts.
//
// input: the calculation or command string to execute
// The function processes the input through defaultExecutor, which handles
// commands via the chain of responsibility pattern, including percentage
// calculations, RPN expressions, and built-in commands.
func executor(input string) {
	// Initialize executorREPL only once using sync.Once for thread-safe lazy initialization
	executorREPLOnce.Do(func() {
		vars := rpn.NewVariables()
		rpnState := &RPNState{
			vars:    vars,
			rpnCalc: rpn.NewRPN(vars),
		}

		// Create a minimal REPL instance without building a prompt
		executorREPL = &REPL{
			ttyChecker:    &TTYChecker{},
			historyMgr:    NewHistoryManager(".gt_history"),
			signalHandler: NewSignalHandler(),
			commandChain:  NewCommandChain(),
			rpnState:      rpnState,
		}
	})

	// Use mutex to protect access to executorREPL during execution
	executorREPLMu.Lock()
	repl := executorREPL
	executorREPLMu.Unlock()

	defaultExecutor(repl, input)
}

// runRPN parses and evaluates an RPN (Reverse Polish Notation) expression.
// This is a package-level wrapper for backward compatibility that delegates to
// the executor's REPL runRPN method.
//
// input: the RPN expression to evaluate
// Returns the result string and an error if the expression is invalid
func runRPN(input string) (string, error) {
	executorREPLMu.Lock()
	defer executorREPLMu.Unlock()

	if executorREPL != nil {
		return executorREPL.rpnState.rpnCalc.ParseAndEvaluate(input)
	}
	return "", fmt.Errorf("no executor REPL available - call executor() first")
}

// getRPNState returns the RPN state from the executor's REPL.
// This is a package-level helper for backward compatibility with tests that need
// to access RPN state after calling executor(). It's not part of the main REPL
// architecture.
//
// Returns the RPNState instance from the last executor() call, or nil if executor() hasn't been called
func getRPNState() *RPNState {
	executorREPLMu.Lock()
	defer executorREPLMu.Unlock()

	if executorREPL != nil {
		return executorREPL.rpnState
	}
	return nil
}





// getHistoryPath returns the absolute path to the history file.
// This is a package-level wrapper for backward compatibility.
// The history file is stored in the user's home directory.
//
// Returns the full path to the history file, or empty string on error
func getHistoryPath() string {
	historyMgr := NewHistoryManager(".gt_history")
	return historyMgr.Path()
}

// loadHistory loads history from the history file.
// This is a package-level wrapper for backward compatibility that uses NewHistoryManager.
//
// Returns a slice of history entries, or nil if the file doesn't exist
func loadHistory() []string {
	historyMgr := NewHistoryManager(".gt_history")
	return historyMgr.Load()
}

// saveHistory saves history to the history file.
// This is a package-level wrapper for backward compatibility that uses NewHistoryManager.
//
// history: the slice of history entries to save
// Returns an error if the file cannot be written
func saveHistory(history []string) error {
	historyMgr := NewHistoryManager(".gt_history")
	return historyMgr.Save(history)
}

// isBuiltinCommand checks if input starts with a built-in command.
// It performs case-insensitive matching against known built-in commands.
//
// input: the command string to check
// Returns the input string and true if it starts with a built-in command,
// or empty string and false otherwise
func isBuiltinCommand(input string) (string, bool) {
	args := strings.Fields(input)
	if len(args) == 0 {
		return "", false
	}

	cmd := strings.ToLower(args[0])
	for _, builtin := range builtinCommands() {
		if cmd == builtin {
			return input, true
		}
	}
	return "", false
}
