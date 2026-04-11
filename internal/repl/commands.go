// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"fmt"
	"strings"
)

// builtinCommandsList is the list of built-in REPL commands.
// It's exposed as a variable to allow for dependency injection in tests.
// Commands: help, clear, quit, exit, rpn, calc, rat
var builtinCommandsList = []string{"help", "clear", "quit", "exit", "rpn", "calc", "rat"}

// Commands returns the list of built-in command names supported by the REPL.
// This is a public function that exposes the built-in command list.
//
// Returns a slice of built-in command names (e.g., "help", "clear", "quit")
func Commands() []string {
	return builtinCommandsList
}

// ExecuteCommand runs a built-in command and returns its output or error.
// It dispatches to the appropriate command handler based on the command name.
//
// cmd: the full command string (e.g., "help", "clear", "rpn 3 4 +")
// Returns the command output string and an error if the command failed
func ExecuteCommand(cmd string) (string, error) {
	args := strings.Fields(cmd)
	if len(args) == 0 {
		return "", nil
	}

	switch strings.ToLower(args[0]) {
	case "help":
		return cmdHelp(args[1:]), nil
	case "clear":
		return "", cmdClear()
	case "quit", "exit":
		return "", cmdQuit()
	case "rpn", "calc":
		// rpn/calc commands are handled in executor(), not here
		return "", nil
	case "rat":
		// rat command is handled in executor() with access to RPN state
		return "", nil
	default:
		return "", fmt.Errorf("unknown command: %s. Available commands: %s", args[0], strings.Join(builtinCommandsList, ", "))
	}
}

// cmdHelp returns help text for built-in commands.
// When called with no subcommands, it returns comprehensive help for all commands.
// When called with a subcommand, it returns specific help for that command.
//
// subCmds: optional slice of subcommand arguments (e.g., ["help"] for "help help")
// Returns the help text as a string
func cmdHelp(subCmds []string) string {
	helpText := `PERC - Percentage Calculator REPL

Built-in Commands:
  help             Show this help message
  help <command>   Show help for a specific topic
  clear            Clear the screen
  quit / exit      Exit the REPL
  rpn / calc       Evaluate an RPN (postfix notation) expression
  rat on/off/toggle Switch between float64 and rational number modes

Usage Examples:
  20% of 150           Calculate 20% of 150
  what is 20% of 150   Same as above (what is prefix is optional)
  30 is what % of 150  Calculate what percentage 30 is of 150
  30 is 20% of what    Calculate what number 30 is 20% of

RPN (Reverse Polish Notation) Examples:
  rpn 3 4 +            3 + 4 = 7
  rpn 3 4 + 4 4 - *    (3 + 4) * (4 - 4) = 0
  rpn x 5 = x x +      Assign x=5, then x + x = 10
  rpn 2 3 ^            2^3 = 8
  rpn 1 2 swap         Swap top two stack values
  rpn 1 2 3 dup        Duplicate top value
  rpn show             Show current stack state

Keyboard Shortcuts (Emacs Mode - default):
  Ctrl+A         Go to beginning of line
  Ctrl+E         Go to end of line
  Ctrl+L         Clear the screen
  Ctrl+D         Delete character under cursor
  Ctrl+H         Delete character before cursor (Backspace)
  Ctrl+F         Forward one character
  Ctrl+B         Backward one character
  Ctrl+W         Cut word before cursor
  Ctrl+K         Cut line after cursor
  Ctrl+U         Cut line before cursor

History Navigation:
  Up Arrow         Previous command
  Down Arrow       Next command

Press Ctrl+D or type 'quit'/'exit' to exit.
`

	if len(subCmds) == 0 {
		return helpText
	}

	subCmd := strings.ToLower(subCmds[0])
	switch subCmd {
	case "help":
		return "help - Show this help message\nUsage: help [command]"
	case "clear":
		return "clear - Clear the screen\nUsage: clear"
	case "quit", "exit":
		return "quit / exit - Exit the REPL\nUsage: quit or exit"
	default:
		return fmt.Sprintf("No help available for: %s\nAvailable commands: help, clear, quit, exit, rpn, calc", subCmd)
	}
}

// cmdClear clears the terminal screen using ANSI escape sequences.
// It prints \033[2J\033[H to clear all content and move the cursor to (0,0).
//
// Returns nil on success
func cmdClear() error {
	// Clear screen using ANSI escape sequence
	fmt.Print("\033[2J\033[H")
	return nil
}

// cmdQuit displays a farewell message and signals REPL exit.
// It's called when the user enters "quit" or "exit" commands.
//
// Returns nil (exit is handled by the REPL itself)
func cmdQuit() error {
	fmt.Println("Goodbye!")
	return nil
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
	for _, builtin := range builtinCommandsList {
		if cmd == builtin {
			return input, true
		}
	}
	return "", false
}
