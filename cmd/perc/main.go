package main

import (
	"fmt"
	"os"
	"strings"

	"codeberg.org/snonux/perc/internal"
	"codeberg.org/snonux/perc/internal/calculator"
	"codeberg.org/snonux/perc/internal/rpn"
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
			if err := runREPL(); err != nil {
				return "", err
			}
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
		// REPL command explicitly requested - run it (may fail if not a TTY)
		if err := runREPL(); err != nil {
			// If not a TTY, just return empty string (REPL can't run in non-interactive mode)
			if !isatty.IsTerminal(os.Stdin.Fd()) {
				return "", nil
			}
			return "", err
		}
		return "", nil
	}

	// Check for calc subcommand
	if args[1] == "calc" || args[1] == "rpn" {
		if len(args) < 3 {
			return "", fmt.Errorf("missing expression after '%s'", args[1])
		}
		input := strings.Join(args[2:], " ")
		result, err := runRPN(input)
		if err != nil {
			return "", err
		}
		return result, nil
	}

	input := strings.Join(args[1:], " ")
	
	// Try RPN parsing first (for bare RPN expressions like "3 4 +")
	rpnResult, rpnErr := runRPN(input)
	if rpnErr == nil {
		return rpnResult, nil
	}
	
	// Fall back to percentage calculation
	result, err := calculator.Parse(input)
	if err != nil {
		return "", err
	}

	return result, nil
}

// runREPL runs the REPL and handles errors
func runREPL() error {
	if err := repl.RunREPL(); err != nil {
		return fmt.Errorf("REPL error: %w", err)
	}
	return nil
}

// runRPN parses and evaluates an RPN expression
func runRPN(input string) (string, error) {
	vars := rpn.NewVariables()
	rpnCalc := rpn.NewRPN(vars)
	return rpnCalc.ParseAndEvaluate(input)
}

func printUsage() {
	fmt.Println("Usage: perc <calculation>")
	fmt.Println("       perc calc <rpn-expression>")
	fmt.Println("       perc rpn <rpn-expression>")
	fmt.Println("       perc version")
	fmt.Println("       perc [--repl|repl]")
	fmt.Println("\nPercentage calculator examples:")
	fmt.Println("  perc 20% of 150")
	fmt.Println("  perc what is 20% of 150")
	fmt.Println("  perc 30 is what % of 150")
	fmt.Println("  perc 30 is 20% of what")
	fmt.Println("\nRPN (postfix notation) examples:")
	fmt.Println("  perc calc 3 4 +")
	fmt.Println("  perc calc 3 4 + 4 4 - *")
	fmt.Println("  perc calc x = 5 x x +")
	fmt.Println("  perc calc 2 3 ^")
	fmt.Println("  perc calc dup swap pop show")
	fmt.Println("\nStart REPL mode interactively by running without arguments:")
	fmt.Println("  perc")
}
