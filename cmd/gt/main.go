package main

import (
	"fmt"
	"os"
	"strings"

	"codeberg.org/snonux/perc/internal"
	"codeberg.org/snonux/perc/internal/calculator"
	"codeberg.org/snonux/perc/internal/repl"
	"codeberg.org/snonux/perc/internal/rpn"
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

	input := strings.Join(args[1:], " ")

	// Try RPN parsing first (for bare RPN expressions like "3 4 +")
	rpnResult, rpnErr := runRPN(input)
	if rpnErr == nil {
		return rpnResult, nil
	}

	// Fall back to percentage calculation
	result, err := calculator.Parse(input)
	if err != nil {
		return "", fmt.Errorf("rpn fallback failed for input %q: %w", input, err)
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
	fmt.Println("Usage: gt <calculation>")
	fmt.Println("       gt version")
	fmt.Println("\nPercentage calculator examples:")
	fmt.Println("  gt 20% of 150")
	fmt.Println("  gt what is 20% of 150")
	fmt.Println("  gt 30 is what % of 150")
	fmt.Println("  gt 30 is 20% of what")
	fmt.Println("\nRPN (postfix notation) examples:")
	fmt.Println("  gt 3 4 +")
	fmt.Println("  gt 3 4 + 4 4 - *")
	fmt.Println("  gt x 5 = x x +")
	fmt.Println("  gt 2 3 ^")
	fmt.Println("  gt dup swap pop show")
	fmt.Println("\nStart REPL mode interactively by running without arguments:")
	fmt.Println("  gt")
}
