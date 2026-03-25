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

// rpnState holds the singleton RPN state for REPL operations.
// It is initialized lazily using sync.Once to ensure thread-safe initialization.
var rpnState *RPNState

// rpnStateOnce ensures rpnState is initialized exactly once.
// It's used by getRPNState to guarantee lazy singleton initialization.
var rpnStateOnce sync.Once

// REPL manages the interactive command-line interface for the percentage calculator.
// It provides an interactive prompt with history, tab-completion, signal handling,
// and command processing through a chain of responsibility pattern.
//
// The REPL integrates various components:
//   - TTYChecker: validates stdin is a terminal
//   - HistoryManager: manages command history persistence
//   - SignalHandler: handles SIGINT (Ctrl+C)
//   - commandChain: processes commands via chain of responsibility
type REPL struct {
	ttyChecker    *TTYChecker
	historyMgr    *HistoryManager
	signalHandler *SignalHandler
	prompt        *prompt.Prompt
	commandChain  CommandHandler
}

// NewREPL creates a new REPL instance with default components.
// If executor is nil, it uses defaultExecutor which processes input through commandChain.
// If completer is nil, it uses defaultCompleter which provides built-in command suggestions.
//
// The executor function is called for each non-empty input line.
// The completer function provides tab-completion suggestions for the prompt.
func NewREPL(executor func(string), completer func(prompt.Document) []prompt.Suggest) *REPL {
	repl := &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  NewCommandChain(),
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
	// Create a minimal REPL instance without building a prompt
	r := &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  NewCommandChain(),
	}
	defaultExecutor(r, input)
}

// getRPNState returns or creates the RPN state using lazy initialization.
// It's thread-safe using sync.Once to ensure the RPN state is initialized exactly once.
// The RPN state is shared across all REPL instances.
//
// Returns the RPNState instance for performing RPN calculations
func getRPNState() *RPNState {
	rpnStateOnce.Do(func() {
		vars := rpn.NewVariables()
		rpnState = &RPNState{
			vars:    vars,
			rpnCalc: rpn.NewRPN(vars),
		}
	})
	return rpnState
}

// getRPNState returns the RPN state.
// This is a REPL instance method for backward compatibility that delegates to the package-level getRPNState.
//
// Returns the RPNState instance for performing RPN calculations
func (r *REPL) getRPNState() *RPNState {
	return getRPNState()
}

// runRPN parses and evaluates an RPN (Reverse Polish Notation) expression.
// It uses the shared RPN state to maintain stack state across multiple calls.
//
// input: the RPN expression to evaluate (e.g., "3 4 +" or "x 5 = x x +")
// Returns the result string and an error if the expression is invalid
func runRPN(input string) (string, error) {
	state := getRPNState()
	return state.rpnCalc.ParseAndEvaluate(input)
}

// runRPN parses and evaluates an RPN (Reverse Polish Notation) expression.
// This is a REPL instance method for backward compatibility that delegates to the package-level runRPN.
//
// input: the RPN expression to evaluate
// Returns the result string and an error if the expression is invalid
func (r *REPL) runRPN(input string) (string, error) {
	return runRPN(input)
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
