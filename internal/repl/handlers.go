package repl

import (
	"fmt"
	"strconv"
	"strings"

	"codeberg.org/snonux/perc/internal/calculator"
	"codeberg.org/snonux/perc/internal/rpn"
)

// CommandHandler represents a handler in the chain of responsibility
// Each handler can process a command or pass it to the next handler
type CommandHandler interface {
	Handle(repl *REPL, input string) (output string, handled bool, err error)
	SetNext(next CommandHandler)
}

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	next CommandHandler
}

// SetNext sets the next handler in the chain
func (h *BaseHandler) SetNext(next CommandHandler) {
	h.next = next
}

// Next forwards the request to the next handler in the chain
func (h *BaseHandler) Next(repl *REPL, input string) (output string, handled bool, err error) {
	if h.next == nil {
		return "", false, nil
	}
	return h.next.Handle(repl, input)
}

// BuiltInCommandHandler handles built-in commands like help, clear, quit, exit
type BuiltInCommandHandler struct {
	BaseHandler
}

// Handle processes built-in commands
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
func handleRatCommand(repl *REPL, input string) (string, bool, error) {
	args := strings.Fields(input)
	if len(args) < 2 {
		return "rat command requires an argument: on, off, or toggle", true, nil
	}

	modeArg := strings.ToLower(args[1])
	rpnState := repl.getRPNState()

	switch modeArg {
	case "on":
		rpnState.rpnCalc.SetMode(rpn.RationalMode)
		return "Rational mode enabled", true, nil
	case "off":
		rpnState.rpnCalc.SetMode(rpn.FloatMode)
		return "Rational mode disabled (using float64)", true, nil
	case "toggle":
		if rpnState.rpnCalc.GetMode() == rpn.FloatMode {
			rpnState.rpnCalc.SetMode(rpn.RationalMode)
			return "Rational mode enabled", true, nil
		} else {
			rpnState.rpnCalc.SetMode(rpn.FloatMode)
			return "Rational mode disabled (using float64)", true, nil
		}
	default:
		return "Unknown rat mode: " + modeArg + ". Valid modes: on, off, toggle", true, nil
	}
}

// RPNHandler handles RPN expressions and RPN-related commands
type RPNHandler struct {
	BaseHandler
}

// Handle processes RPN commands and expressions
func (h *RPNHandler) Handle(repl *REPL, input string) (output string, handled bool, err error) {
	// Check for rpn/calc prefix
	lowerInput := strings.ToLower(input)
	if strings.HasPrefix(lowerInput, "rpn ") || strings.HasPrefix(lowerInput, "calc ") {
		// Extract the expression after rpn/calc
		rest := strings.TrimSpace(strings.TrimPrefix(input, strings.SplitN(input, " ", 2)[0]))
		result, err := repl.runRPN(rest)
		if err != nil {
			return "", true, err
		}
		return result, true, nil
	}

	// Try RPN parsing first (for bare RPN expressions like "3 4 +")
	if state := repl.getRPNState(); state != nil {
		// Check if input looks like RPN (contains spaces or is a single known operator)
		if strings.Contains(input, " ") {
			result, err := repl.runRPN(input)
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
				result, err := state.rpnCalc.EvalOperator(op)
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
				result, err := state.rpnCalc.ParseAndEvaluate(fields[0])
				if err != nil {
					return "", true, err
				}
				return result, true, nil
			}
		}
	}

	return h.Next(repl, input)
}

// PercentageHandler handles percentage calculations
type PercentageHandler struct {
	BaseHandler
}

// Handle processes percentage calculation expressions
func (h *PercentageHandler) Handle(repl *REPL, input string) (output string, handled bool, err error) {
	// Run the percentage calculation
	result, err := calculator.Parse(input)
	if err != nil {
		// Not a percentage expression, pass to next handler
		return h.Next(repl, input)
	}
	return result, true, nil
}

// ErrorHandler handles unknown commands
type ErrorHandler struct {
	BaseHandler
}

// Handle processes unknown commands by returning an error
func (h *ErrorHandler) Handle(repl *REPL, input string) (output string, handled bool, err error) {
	// Unknown command - return error
	return "", false, fmt.Errorf("unknown command or invalid expression: %s", input)
}

// NewCommandChain creates and returns the complete command handling chain
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
