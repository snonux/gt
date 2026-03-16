package main

import (
	"fmt"
	"os"
	"strings"

	"codeberg.org/snonux/perc/internal"
	"codeberg.org/snonux/perc/internal/calculator"
	"codeberg.org/snonux/perc/internal/repl"
	"github.com/mattn/go-isatty"
)

func main() {
	output, err := runCommand(os.Args)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println(output)
}

func runCommand(args []string) (string, error) {
	if len(args) < 2 {
		// No args provided - check if stdin is a TTY for REPL mode
		if isatty.IsTerminal(os.Stdin.Fd()) {
			repl.RunREPL()
			return "", nil
		}
		printUsage()
		return "", fmt.Errorf("no input provided")
	}

	if args[1] == "version" {
		return internal.Version, nil
	}

	// Check for --repl flag
	if args[1] == "--repl" || args[1] == "repl" {
		repl.RunREPL()
		return "", nil
	}

	input := strings.Join(args[1:], " ")
	result, err := calculator.Parse(input)
	if err != nil {
		return "", err
	}

	return result, nil
}

func printUsage() {
	fmt.Println("Usage: perc <calculation>")
	fmt.Println("       perc version")
	fmt.Println("       perc [--repl|repl]")
	fmt.Println("\nExamples:")
	fmt.Println("  perc 20% of 150")
	fmt.Println("  perc what is 20% of 150")
	fmt.Println("  perc 30 is what % of 150")
	fmt.Println("  perc 30 is 20% of what")
	fmt.Println("\nStart REPL mode interactively by running without arguments:")
	fmt.Println("  perc")
}
