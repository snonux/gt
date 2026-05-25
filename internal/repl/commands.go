// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"fmt"
	"strings"
)

// builtinCommandsList is the list of built-in REPL commands.
// It's exposed as a variable to allow for dependency injection in tests.
// Commands: help, clear, quit, exit, rpn, calc, rat, stack
var builtinCommandsList = []string{"help", "clear", "quit", "exit", "rpn", "calc", "rat", "stack"}

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
	case "stack":
		return cmdStack(), nil
	default:
		return "", fmt.Errorf("unknown command: %s. Available commands: %s", args[0], strings.Join(builtinCommandsList, ", "))
	}
}

// cmdHelp returns help text for built-in commands.
// When called with no subcommands, it returns the general help overview.
// When called with a subcommand, it returns specific inline help for that topic.
//
// subCmds: optional slice of subcommand arguments (e.g., ["+"] for "help +")
// Returns the help text as a string
func cmdHelp(subCmds []string) string {
	if len(subCmds) == 0 {
		return GetHelp("")
	}
	return GetHelp(strings.ToLower(subCmds[0]))
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

// cmdStack provides a hint for viewing the RPN stack.
// It does not have access to RPN state; users should use 'rpn show' instead.
//
// Returns a hint string directing users to the 'rpn show' command
func cmdStack() string {
	return "Use 'rpn show' to view the current stack state"
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
