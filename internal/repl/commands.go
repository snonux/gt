package repl

import (
	"fmt"
	"strings"
)

// builtinCommands defines the built-in REPL commands
var builtinCommands = []string{"help", "clear", "quit", "exit", "rpn", "calc"}

// Commands returns the list of built-in command names
func Commands() []string {
	return builtinCommands
}

// ExecuteCommand runs a built-in command and returns its output or error
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
	default:
		return "", fmt.Errorf("unknown command: %s. Available commands: help, clear, quit, exit, rpn, calc", args[0])
	}
}

func cmdHelp(subCmds []string) string {
	helpText := `PERC - Percentage Calculator REPL

Built-in Commands:
  help             Show this help message
  help <command>   Show help for a specific topic
  clear            Clear the screen
  quit / exit      Exit the REPL
  rpn / calc       Evaluate an RPN (postfix notation) expression

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

func cmdClear() error {
	// Clear screen using ANSI escape sequence
	fmt.Print("\033[2J\033[H")
	return nil
}

func cmdQuit() error {
	fmt.Println("Goodbye!")
	return nil
}
