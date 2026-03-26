// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

// Package gt provides a command-line percentage calculator with RPN support.
//
// gt is a versatile calculator that supports both percentage calculations and
// Reverse Polish Notation (RPN) expressions. It can be used in two modes:
//
// 1. Command-line mode: Pass calculations as arguments
//    gt 20% of 150         # Calculate 20% of 150
//    gt 3 4 +              # RPN expression: 3 + 4
//
// 2. Interactive REPL mode: Run without arguments to start an interactive session
//    gt                    # Start interactive REPL
//
// Percentage Calculations
//
// The calculator supports various percentage formats:
//   - Basic percentage: "20% of 150" → 30
//   - With prefix: "what is 20% of 150" → 30
//   - Reverse percentage: "30 is what % of 150" → 20%
//   - Find base: "30 is 20% of what" → 150
//
// RPN (Reverse Polish Notation) Support
//
// RPN expressions use postfix notation where operators follow operands:
//   - Basic operations: "3 4 +" (3 + 4), "5 2 -" (5 - 2)
//   - Complex expressions: "3 4 + 4 4 - *" ((3 + 4) * (4 - 4))
//   - Exponentiation: "2 3 ^" (2^3 = 8)
//   - Variable assignment: "x 5 = x x +" (assign x=5, then x + x)
//   - Stack operations: "dup swap pop show"
//
// Error Handling
//
// Errors from calculations or parsing are printed to stdout with exit code 1.
// Invalid RPN expressions and malformed percentage queries both return errors.
//
// Architecture
//
// The package uses a layered architecture:
//   - main.go: Entry point and command routing
//   - perc/: Handles percentage calculation parsing
//   - rpn/: Handles RPN expression parsing and evaluation
//   - repl/: Provides interactive Read-Eval-Print Loop mode
//
// See the cmd/gt/internal package for version information.
package main

import (
	"fmt"
	"os"
	"strings"

	"codeberg.org/snonux/gt/internal"
	"codeberg.org/snonux/gt/internal/perc"
	"codeberg.org/snonux/gt/internal/repl"
	"codeberg.org/snonux/gt/internal/rpn"
	"github.com/mattn/go-isatty"
)

// main is the entry point for the gt command-line calculator.
func main() {
	output, err := runCommand(os.Args)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println(output)
}

// runCommand processes command-line arguments and executes the appropriate action.
//
// It handles:
//   - No arguments: Start REPL mode if stdin is a TTY, otherwise read from stdin
//   - "version" argument: Return the version string
//   - Other arguments: Try RPN parsing first, then fall back to percentage calculation
func runCommand(args []string) (string, error) {
	if len(args) < 2 {
		// No args provided - check if stdin is a TTY for REPL mode
		if isatty.IsTerminal(os.Stdin.Fd()) {
			if err := runREPL(); err != nil {
				return "", err
			}
			return "", nil
		}
		// Read from stdin (pipe or redirect)
		input, err := readStdin()
		if err != nil {
			return "", fmt.Errorf("failed to read stdin: %w", err)
		}
		input = strings.TrimSpace(input)
		if input == "" {
			printUsage()
			return "", fmt.Errorf("no input provided")
		}
		// Try RPN parsing first
		rpnResult, rpnErr := runRPN(input)
		if rpnErr == nil {
			return rpnResult, nil
		}
		// Fall back to percentage calculation
		result, err := perc.Parse(input)
		if err != nil {
			return "", fmt.Errorf("rpn fallback failed for input %q: %w", input, err)
		}
		return result, nil
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
	result, err := perc.Parse(input)
	if err != nil {
		return "", fmt.Errorf("rpn fallback failed for input %q: %w", input, err)
	}

	return result, nil
}

// readStdin reads all input from stdin and returns it as a string.
func readStdin() (string, error) {
	data, err := os.ReadFile("/dev/stdin")
	if err != nil {
		// Fallback if /dev/stdin is not available
		buf := make([]byte, 4096)
		n, err := os.Stdin.Read(buf)
		if n > 0 {
			return string(buf[:n]), nil
		}
		return "", err
	}
	return string(data), nil
}

// runREPL starts the interactive REPL mode.
//
// It wraps repl.RunREPL() and returns an error if the REPL fails to start.
func runREPL() error {
	if err := repl.RunREPL(); err != nil {
		return fmt.Errorf("REPL error: %w", err)
	}
	return nil
}

// runRPN parses and evaluates an RPN (Reverse Polish Notation) expression.
//
// It creates a fresh RPN calculator with fresh variable store for each call,
// making it suitable for one-off calculations.
func runRPN(input string) (string, error) {
	vars := rpn.NewVariables()
	rpnCalc := rpn.NewRPN(vars)

	// Strip "rpn " or "calc " prefix if present
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "rpn ") {
		input = strings.TrimPrefix(input, "rpn ")
	} else if strings.HasPrefix(input, "calc ") {
		input = strings.TrimPrefix(input, "calc ")
	}

	return rpnCalc.ParseAndEvaluate(input)
}

// printUsage displays the command-line usage information and examples.
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
