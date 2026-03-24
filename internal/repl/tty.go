package repl

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
)

// TTYChecker provides TTY detection functionality.
type TTYChecker struct{}

// IsTTY returns true if stdin is a terminal.
func (c *TTYChecker) IsTTY() bool {
	return isatty.IsTerminal(os.Stdin.Fd())
}

// EnsureTTY checks if stdin is a TTY and returns an error if not.
func (c *TTYChecker) EnsureTTY() error {
	if !c.IsTTY() {
		fmt.Fprintln(os.Stderr, "REPL mode requires a TTY. Use 'gt <calculation>' for non-interactive mode.")
		return fmt.Errorf("stdin is not a TTY")
	}
	return nil
}
