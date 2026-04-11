// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"fmt"
	"strconv"
	"strings"

	"codeberg.org/snonux/gt/internal/perc"
	"codeberg.org/snonux/gt/internal/rpn"
)

// CommandHandler represents a handler in the chain of responsibility pattern.
// Each handler can process a command or pass it to the next handler in the chain.
//
// Handlers implement the Handle method to process REPL commands and expressions.
// If a handler cannot process the input, it calls Next() to forward to the next handler.
type CommandHandler interface {
	Handle(repl *REPL, input string) (output string, handled bool, err error)
	SetNext(next CommandHandler)
}

// BaseHandler provides common functionality for all handlers in the chain.
// It stores a reference to the next handler and provides the Next() method
// for forwarding requests.
type BaseHandler struct {
	next CommandHandler
}

// SetNext sets the next handler in the chain.
// This enables building a chain of responsibility by linking handlers together.
//
// next: the next CommandHandler in the chain
func (h *BaseHandler) SetNext(next CommandHandler) {
	h.next = next
}

// Next forwards the request to the next handler in the chain.
// If there is no next handler, it returns (false, nil) indicating the request
// was not handled.
//
// repl: the REPL instance
// input: the command string to process
// Returns: (output string, handled bool, err error)
func (h *BaseHandler) Next(repl *REPL, input string) (output string, handled bool, err error) {
	if h.next == nil {
		return "", false, nil
	}
	return h.next.Handle(repl, input)
}

// BuiltInCommandHandler handles built-in commands like help, clear, quit, exit.
// It also handles special commands that require RPN state access (e.g., "rat").
// If the input doesn't match a built-in command, it forwards to the next handler.
type BuiltInCommandHandler struct {
	BaseHandler
}

// Handle processes built-in commands from the input string.
// It first checks if the input starts with a built-in command using isBuiltinCommand.
// Special handling is provided for the "rat" command which requires RPN state access.
// If the command is handled, it returns the output and sets handled=true.
// Otherwise, it forwards to the next handler in the chain.
//
// repl: the REPL instance
// input: the command string to process
// Returns: (output string, handled bool, err error)
func (h *BuiltInCommandHandler) Handle(repl *REPL, input string) (output string, handled bool, err error) {
	if cmd, ok := isBuiltinCommand(input); ok {
		args := strings.Fields(cmd)
		if len(args) > 0 {
			subCmd := strings.ToLower(args[0])
			// Handle rat command specially - needs RPN state access
			if subCmd == "rat" {
				return handleRatCommand(repl, input)
			}
		}
		output, err := ExecuteCommand(cmd)
		if err != nil {
			return "", true, err
		}
		return output, true, nil
	}
	return h.Next(repl, input)
}

// handleRatCommand handles the rat mode command with access to RPN state.
// It allows switching between float64 and rational number modes for RPN calculations.
//
// Valid modes:
//   - "on":  Enable rational number mode
//   - "off": Disable rational mode (use float64)
//   - "toggle": Switch between the two modes
//
// repl: the REPL instance (provides access to RPN state)
// input: the full command string (e.g., "rat on")
// Returns: (output string, handled bool, err error)
func handleRatCommand(repl *REPL, input string) (string, bool, error) {
	args := strings.Fields(input)
	if len(args) < 2 {
		return "rat command requires an argument: on, off, or toggle", true, nil
	}

	modeArg := strings.ToLower(args[1])
	rpnState := repl.rpnState
	calculator := rpnState.calculator

	switch modeArg {
	case "on":
		calculator.SetMode(rpn.RationalMode)
		return "Rational mode enabled", true, nil
	case "off":
		calculator.SetMode(rpn.FloatMode)
		return "Rational mode disabled (using float64)", true, nil
	case "toggle":
		if calculator.GetMode() == rpn.FloatMode {
			calculator.SetMode(rpn.RationalMode)
			return "Rational mode enabled", true, nil
		} else {
			calculator.SetMode(rpn.FloatMode)
			return "Rational mode disabled (using float64)", true, nil
		}
	default:
		return "Unknown rat mode: " + modeArg + ". Valid modes: on, off, toggle", true, nil
	}
}

// RPNHandler handles RPN (Reverse Polish Notation) expressions and RPN-related commands.
// It processes commands with "rpn" or "calc" prefixes, bare RPN expressions,
// and single RPN operators (e.g., "+", "dup", "swap", "show").
type RPNHandler struct {
	BaseHandler
}

// Handle processes RPN commands and expressions.
// It handles:
//   - Commands with "rpn" or "calc" prefix
//   - Bare RPN expressions (e.g., "3 4 +")
//   - Single RPN operators on the current stack
//   - Single numbers (push onto stack)
//
// If the input doesn't match any RPN pattern, it forwards to the next handler.
//
// repl: the REPL instance (provides access to RPN state)
// input: the command string to process
// Returns: (output string, handled bool, err error)
func (h *RPNHandler) Handle(repl *REPL, input string) (output string, handled bool, err error) {
	// Check for rpn/calc prefix
	lowerInput := strings.ToLower(input)
	if strings.HasPrefix(lowerInput, "rpn ") || strings.HasPrefix(lowerInput, "calc ") {
		// Extract the expression after rpn/calc
		rest := strings.TrimSpace(strings.TrimPrefix(input, strings.SplitN(input, " ", 2)[0]))
		result, err := repl.rpnState.calculator.ParseAndEvaluate(rest)
		if err != nil {
			return "", true, err
		}
		return result, true, nil
	}

	// Try RPN parsing first (for bare RPN expressions like "3 4 +")
	if state := repl.rpnState; state != nil {
		calculator := state.calculator
		// Check if input looks like RPN (contains spaces or is a single known operator)
		if strings.Contains(input, " ") {
			result, err := calculator.ParseAndEvaluate(input)
			if err == nil {
				return result, true, nil
			}
		}

		// Try evaluating as a single operator on the current RPN stack
		fields := strings.Fields(input)
		if len(fields) == 1 {
			op := strings.ToLower(fields[0])
			// Check if it's a known operator (standard or hyper)
			isStandardOp := op == "+" || op == "-" || op == "*" || op == "/" || op == "^" || op == "%" ||
				op == "dup" || op == "swap" || op == "pop" || op == "show" || op == "clear" || op == "vars" ||
				op == "lg" || op == "log" || op == "ln"
			isHyperOp := op == "[+]" || op == "[-]" || op == "[*]" || op == "[/]" || op == "[^]" || op == "[%]" ||
				op == "[lg]" || op == "[log]" || op == "[ln]"

			if isStandardOp || isHyperOp {
				result, err := calculator.EvalOperator(op)
				if err != nil {
					return "", true, err
				}
				return result, true, nil
			}
		}

		// Check if input is a single number (valid RPN - pushes number onto stack)
		if len(fields) == 1 {
			if _, err := strconv.ParseFloat(fields[0], 64); err == nil {
				// Push the number onto the RPN stack using ParseAndEvaluate
				// This maintains the RPN state across multiple inputs in REPL mode
				result, err := calculator.ParseAndEvaluate(fields[0])
				if err != nil {
					return "", true, err
				}
				return result, true, nil
			}
		}
		
		// Check if input is a symbol syntax (:x) - valid RPN that pushes a symbol
		if len(fields) == 1 {
			token := fields[0]
			if len(token) > 0 && token[0] == ':' {
				// This is a symbol syntax like :x
				result, err := calculator.ParseAndEvaluate(token)
				if err != nil {
					return "", true, err
				}
				return result, true, nil
			}
		}
	}

	return h.Next(repl, input)
}

// PercentageHandler handles percentage calculation expressions.
// It uses the perc.Parse function to evaluate expressions like:
//   - "20% of 150"
//   - "what is 20% of 150"
//   - "30 is what % of 150"
//   - "30 is 20% of what"
type PercentageHandler struct {
	BaseHandler
}

// Handle processes percentage calculation expressions.
// If the input matches a percentage expression pattern, it evaluates and returns
// the result. Otherwise, it forwards to the next handler.
//
// repl: the REPL instance
// input: the command string to process
// Returns: (output string, handled bool, err error)
func (h *PercentageHandler) Handle(repl *REPL, input string) (output string, handled bool, err error) {
	// Run the percentage calculation
	result, err := perc.Parse(input)
	if err != nil {
		// Not a percentage expression, pass to next handler
		return h.Next(repl, input)
	}
	return result, true, nil
}

// ErrorHandler handles unknown commands and invalid expressions.
// It returns an error indicating that the input was not recognized.
type ErrorHandler struct {
	BaseHandler
}

// Handle processes unknown commands by returning an error.
// This is typically the last handler in the chain.
//
// repl: the REPL instance
// input: the command string that was not handled by previous handlers
// Returns: ("", false, error) with an error describing the unknown command
func (h *ErrorHandler) Handle(repl *REPL, input string) (output string, handled bool, err error) {
	// Unknown command - return error
	return "", false, fmt.Errorf("unknown command or invalid expression: %s", input)
}

// NewCommandChain creates and returns the complete command handling chain.
// The chain is built in the following order:
//  1. BuiltInCommandHandler: handles built-in commands (help, clear, quit, exit, rat)
//  2. RPNHandler: handles RPN expressions and operators
//  3. PercentageHandler: handles percentage calculations
//  4. ErrorHandler: handles unknown commands (returns error)
//
// Returns a CommandHandler representing the first handler in the chain
func NewCommandChain() CommandHandler {
	// Create handlers
	builtInHandler := &BuiltInCommandHandler{}
	rpnHandler := &RPNHandler{}
	percentageHandler := &PercentageHandler{}
	errorHandler := &ErrorHandler{}

	// Build the chain: BuiltIn -> RPN -> Percentage -> Error
	builtInHandler.SetNext(rpnHandler)
	rpnHandler.SetNext(percentageHandler)
	percentageHandler.SetNext(errorHandler)

	return builtInHandler
}
