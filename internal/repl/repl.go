package repl

import (
	"fmt"
	"strings"
	"sync"

	"codeberg.org/snonux/perc/internal/rpn"

	"github.com/c-bata/go-prompt"
)

// REPL manages the interactive command-line interface.
type REPL struct {
	ttyChecker    *TTYChecker
	historyMgr    *HistoryManager
	signalHandler *SignalHandler
	prompt        *prompt.Prompt
	commandChain  CommandHandler
}

// NewREPL creates a new REPL instance with default components.
// If executor is nil, it uses a default executor.
// If completer is nil, it uses a default completer.
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

// defaultExecutor is the default executor function.
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

// defaultCompleter is the default completer function.
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

// defaultGetCommandDescription returns the description for a command.
func (r *REPL) defaultGetCommandDescription(cmd string) string {
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

// RunREPL starts the interactive REPL.
// This is a convenience wrapper around NewREPL().Run().
func RunREPL() error {
	repl := NewREPL(nil, nil)
	return repl.Run()
}

// executor runs a calculation command and returns the result.
// This is a package-level wrapper for backward compatibility.
// It creates a minimal REPL instance without a prompt for testing purposes.
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

// RPNState holds the state for RPN operations in REPL
// Note: This struct should never be copied - use pointer receivers only
type RPNState struct {
	vars    rpn.VariableStore
	rpnCalc *rpn.RPN
}

// rpnState holds the singleton RPN state for REPL operations
var rpnState *RPNState

// rpnStateOnce ensures rpnState is initialized exactly once
var rpnStateOnce sync.Once

// getRPNState returns or creates the RPN state
// Thread-safe implementation using sync.Once for simpler singleton initialization
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
// This is a REPL instance method for backward compatibility.
func (r *REPL) getRPNState() *RPNState {
	return getRPNState()
}

// runRPN parses and evaluates an RPN expression.
func runRPN(input string) (string, error) {
	state := getRPNState()
	return state.rpnCalc.ParseAndEvaluate(input)
}

// runRPN parses and evaluates an RPN expression.
// This is a REPL instance method for backward compatibility.
func (r *REPL) runRPN(input string) (string, error) {
	return runRPN(input)
}

// getHistoryPath returns the path to the history file.
// This is a package-level wrapper for backward compatibility.
func getHistoryPath() string {
	historyMgr := NewHistoryManager(".gt_history")
	return historyMgr.Path()
}

// loadHistory loads history from file.
// This is a package-level wrapper for backward compatibility.
func loadHistory() []string {
	historyMgr := NewHistoryManager(".gt_history")
	return historyMgr.Load()
}

// saveHistory saves history to file.
// This is a package-level wrapper for backward compatibility.
func saveHistory(history []string) error {
	historyMgr := NewHistoryManager(".gt_history")
	return historyMgr.Save(history)
}

// isBuiltinCommand checks if input starts with a built-in command
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
