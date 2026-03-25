package repl

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
)

// TTYChecker provides TTY detection functionality.
// It uses the go-isatty package to determine if stdin is a terminal.
type TTYChecker struct{}

// IsTTY returns true if stdin is a terminal.
// This is useful for determining whether to run in interactive REPL mode.
//
// Returns true if stdin is a TTY, false otherwise
func (c *TTYChecker) IsTTY() bool {
	return isatty.IsTerminal(os.Stdin.Fd())
}

// EnsureTTY checks if stdin is a TTY and returns an error if not.
// This is used to prevent running the REPL in non-interactive contexts
// (e.g., when stdin is piped from a file or another command).
//
// Returns nil if stdin is a TTY, or an error describing the issue otherwise
func (c *TTYChecker) EnsureTTY() error {
	if !c.IsTTY() {
		fmt.Fprintln(os.Stderr, "REPL mode requires a TTY. Use 'gt <calculation>' for non-interactive mode.")
		return fmt.Errorf("stdin is not a TTY")
	}
	return nil
}
