// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"codeberg.org/snonux/gt/internal/rpn"

	"github.com/chzyer/readline"
)

// RPNState holds the state for RPN (Reverse Polish Notation) operations in the REPL.
// It maintains a variable store and calculator instance.
//
// Note: This struct should never be copied - use pointer receivers only.
type RPNState struct {
	vars         rpn.VariableStore
	calculator   Calculator
	varStoreFile string // Path to persistent variable store file
}

// NewRPNState creates a new RPNState with the given variable store and calculator.
// It also configures the variable store file path in the user's config directory.
//
// vars: the VariableStore instance to use
// calculator: the Calculator instance for RPN operations
// Returns a new RPNState instance with configured variable store file path
func NewRPNState(vars rpn.VariableStore, calculator Calculator) *RPNState {
	varStoreFile := getVarStoreFilePath()
	return &RPNState{
		vars:         vars,
		calculator:   calculator,
		varStoreFile: varStoreFile,
	}
}

// LoadVariables loads the variable store from the persistent file.
// Returns nil on success, or an error if loading fails (except when file doesn't exist).
func (r *RPNState) LoadVariables() error {
	if r.varStoreFile == "" {
		return nil
	}
	return r.vars.Load(r.varStoreFile)
}

// SaveVariables saves the variable store to the persistent file.
// Returns an error if saving fails.
func (r *RPNState) SaveVariables() error {
	if r.varStoreFile == "" {
		return nil
	}
	return r.vars.Save(r.varStoreFile)
}

// getVarStoreFilePath returns the path to the persistent variable store file.
// Variables are stored in ~/.local/state/gt/vars in JSON format (XDG spec).
//
// Returns the absolute path to the variable store file, or empty string on error
func getVarStoreFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".local", "state", "gt", "vars")
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
//   - logWriter: optional writer for session logging
type REPL struct {
	ttyChecker    *TTYChecker
	historyMgr    *HistoryManager
	signalHandler *SignalHandler
	prompt        *ReadlinePrompt
	commandChain  CommandHandler
	rpnState      *RPNState
	logWriter     io.WriteCloser
}

// ReadlinePrompt provides an interactive prompt using chzyer/readline.
// It supports:
//   - Ctrl+R for reverse history search
//   - Arrow keys for history navigation
//   - Tab completion
//   - Multi-line input
type ReadlinePrompt struct {
	instance  *readline.Instance
	executor  func(string)
}

// NewReadlinePrompt creates a new readline-based prompt instance.
func NewReadlinePrompt(prefix string, history []string, executor func(string), completer *AutoCompleteAdapter) (*ReadlinePrompt, error) {
	config := &readline.Config{
		Prompt:          prefix,
		HistoryFile:     "",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	}

	if completer != nil {
		config.AutoComplete = completer
	}

	rl, err := readline.NewEx(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create readline instance: %w", err)
	}

	// Load history - readline handles this automatically via HistoryFile
	// But we can pre-populate the history
	for _, entry := range history {
		rl.SaveHistory(entry)
	}

	return &ReadlinePrompt{
		instance: rl,
		executor: executor,
	}, nil
}

// Run starts the prompt loop and blocks until it exits.
func (p *ReadlinePrompt) Run() error {
	defer p.instance.Close()

	for {
		line, err := p.instance.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				fmt.Println("\nUse 'quit' or 'exit' to exit, or Ctrl+D")
				continue
			}
			return fmt.Errorf("readline error: %w", err)
		}

		input := strings.TrimSpace(line)
		if input == "" {
			continue
		}

		p.executor(input)
	}

	return nil
}

// Close closes the prompt instance.
func (p *ReadlinePrompt) Close() error {
	return p.instance.Close()
}

// NewREPL creates a new REPL instance with default components.
// If executor is nil, it uses defaultExecutor which processes input through commandChain.
// If completer is nil, it uses defaultCompleter which provides built-in command suggestions.
// If logWriter is non-nil, all REPL output is duplicated to the log writer.
//
// The executor function is called for each non-empty input line.
// The completer function provides tab-completion suggestions for the prompt.
func NewREPL(executor func(string), completer func() []string, logWriter io.WriteCloser) *REPL {
	// Initialize RPN state via dependency injection
	vars := rpn.NewVariables()
	// Load persisted variables from file (if any)
	if err := vars.Load(getVarStoreFilePath()); err != nil {
		fmt.Printf("Warning: Could not load saved variables: %v\n", err)
	}

	rpnCalc := rpn.NewRPN(vars)
	calculator := NewRPNCalculator(rpnCalc)
	rpnState := NewRPNState(vars, calculator)

	repl := &REPL{
		ttyChecker:    &TTYChecker{},
		historyMgr:    NewHistoryManager(".gt_history"),
		signalHandler: NewSignalHandler(),
		commandChain:  NewCommandChain(),
		rpnState:      rpnState,
		logWriter:     logWriter,
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
		completerFn = func() []string {
			return Commands()
		}
	}

	// Load history from file
	history := repl.historyMgr.Load()

	// Build the prompt
	var err error
	completerAdapter := NewAutoCompleter()
	repl.prompt, err = NewReadlinePrompt(
		"> ",
		history,
		execFn,
		completerAdapter,
	)
	if err != nil {
		fmt.Printf("Warning: Could not create prompt: %v\n", err)
	}

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
	if r.prompt != nil {
		if err := r.prompt.Run(); err != nil {
			return err
		}
	}

	// Save variables on exit
	_ = r.rpnState.SaveVariables()

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
// Returns a slice of strings for matching built-in commands
func defaultCompleter(r *REPL) []string {
	return Commands()
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
// This is a convenience wrapper around NewREPL(nil, nil, nil).Run().
// It's typically used when the standard REPL behavior is sufficient.
//
// Returns an error if the REPL cannot start (e.g., stdin is not a TTY)
func RunREPL() error {
	repl := NewREPL(nil, nil, nil)
	return repl.Run()
}

// RunREPLWithLog starts the interactive REPL with logging to the specified file.
// The log file receives input commands and output (each prefixed with '> ' for input).
//
// logFile: path to a file to append log output
// Returns an error if the REPL cannot start (e.g., stdin is not a TTY)
func RunREPLWithLog(logFile string) error {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file %q: %w", logFile, err)
	}
	defer file.Close()
	return NewREPL(nil, nil, file).Run()
}
