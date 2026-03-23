package rpn

import (
	"fmt"
	"strconv"
	"strings"
)

// RPN represents the RPN parser and evaluator.
type RPN struct {
	vars         VariableStore
	ops          Operator
	maxStack     int
	currentStack *Stack
}

// NewRPN creates a new RPN parser and evaluator with the given variable store.
func NewRPN(vars VariableStore) *RPN {
	return &RPN{
		vars:         vars,
		ops:          NewOperations(vars),
		maxStack:     1000, // Reasonable limit for RPN expressions
		currentStack: NewStack(),
	}
}

// ParseAndEvaluate parses and evaluates an RPN expression.
// Returns the result as a formatted string or an error.
func (r *RPN) ParseAndEvaluate(input string) (string, error) {
	// Validate input and initialize
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("empty expression")
	}
	if r.currentStack == nil {
		r.currentStack = NewStack()
	}

	// Handle assignment formats
	if assignmentResult, isAssignment, err := r.handleAssignment(input); err != nil {
		return "", err
	} else if isAssignment {
		return assignmentResult, nil
	}

	// Evaluate standard RPN expression
	tokens := Tokenize(input)
	if len(tokens) == 0 {
		return "", fmt.Errorf("no valid tokens found")
	}

	return r.evaluate(tokens)
}

// ResultStack returns the final stack state after evaluation.
// This is useful for commands that need to show the stack without consuming it.
func (r *RPN) ResultStack(tokens []string) (string, error) {
	stack := NewStack()

	for _, token := range tokens {
		// Check if it's a number
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			if stack.Len() >= r.maxStack {
				return "", fmt.Errorf("stack overflow")
			}
			stack.Push(num)
			continue
		}

		// Check for operators and special commands
		switch token {
		case "+":
			if err := r.ops.Add(stack); err != nil {
				return "", err
			}
		case "-":
			if err := r.ops.Subtract(stack); err != nil {
				return "", err
			}
		case "*":
			if err := r.ops.Multiply(stack); err != nil {
				return "", err
			}
		case "/":
			if err := r.ops.Divide(stack); err != nil {
				return "", err
			}
		case "^":
			if err := r.ops.Power(stack); err != nil {
				return "", err
			}
		case "%":
			if err := r.ops.Modulo(stack); err != nil {
				return "", err
			}
		case "dup":
			if err := r.ops.Dup(stack); err != nil {
				return "", err
			}
		case "swap":
			if err := r.ops.Swap(stack); err != nil {
				return "", err
			}
		case "pop":
			if err := r.ops.Pop(stack); err != nil {
				return "", err
			}
		case "show", "showstack", "print":
			return r.ops.Show(stack)
		case "vars":
			return r.ops.ListVariables()
		case "clear":
			r.ops.ClearVariables()
			return "All variables cleared", nil
		default:
			// Check if it's a variable reference (push its value)
			val, exists := r.vars.GetVariable(token)
			if exists {
				stack.Push(val)
			} else {
				return "", fmt.Errorf("unknown token '%s'", token)
			}
		}
	}

	return r.ops.Show(stack)
}

// EvalOperator evaluates a single operator on the current stack state.
// This allows incremental RPN operations like: "1 2 +" then "+".
func (r *RPN) EvalOperator(op string) (string, error) {
	if r.currentStack == nil {
		r.currentStack = NewStack()
	}

	switch op {
	case "+":
		if err := r.ops.Add(r.currentStack); err != nil {
			return "", fmt.Errorf("operator +: %w", err)
		}
	case "-":
		if err := r.ops.Subtract(r.currentStack); err != nil {
			return "", fmt.Errorf("operator -: %w", err)
		}
	case "*":
		if err := r.ops.Multiply(r.currentStack); err != nil {
			return "", fmt.Errorf("operator *: %w", err)
		}
	case "/":
		if err := r.ops.Divide(r.currentStack); err != nil {
			return "", fmt.Errorf("operator /: %w", err)
		}
	case "^":
		if err := r.ops.Power(r.currentStack); err != nil {
			return "", fmt.Errorf("operator ^: %w", err)
		}
	case "%":
		if err := r.ops.Modulo(r.currentStack); err != nil {
			return "", fmt.Errorf("operator %%: %w", err)
		}
	case "dup":
		if err := r.ops.Dup(r.currentStack); err != nil {
			return "", fmt.Errorf("dup: %w", err)
		}
	case "swap":
		if err := r.ops.Swap(r.currentStack); err != nil {
			return "", fmt.Errorf("swap: %w", err)
		}
	case "pop":
		if err := r.ops.Pop(r.currentStack); err != nil {
			return "", fmt.Errorf("pop: %w", err)
		}
	case "show", "showstack", "print":
		return r.ops.Show(r.currentStack)
	case "clear":
		r.ops.ClearVariables()
		return "All variables cleared", nil
	case "vars":
		return r.ops.ListVariables()
	default:
		return "", fmt.Errorf("unknown operator '%s'", op)
	}

	// Return the current stack state
	return r.ops.Show(r.currentStack)
}

// GetCurrentStack returns a copy of the current stack for inspection.
func (r *RPN) GetCurrentStack() []float64 {
	if r.currentStack == nil {
		return nil
	}
	return r.currentStack.Values()
}

// Tokenize splits the input string into tokens (numbers, operators, variables).
// This is exported for testing purposes.
func Tokenize(input string) []string {
	// Standard RPN tokenization
	return strings.Fields(input)
}

// evaluate evaluates a list of tokens and returns the result.
func (r *RPN) evaluate(tokens []string) (string, error) {
	// Use the current stack for evaluation to preserve state
	// This allows incremental operations in REPL mode
	if r.currentStack == nil {
		r.currentStack = NewStack()
	}
	stack := r.currentStack

	for i, token := range tokens {
		// Check for variable assignment: name value =
		if token == "=" {
			return "", fmt.Errorf("invalid assignment syntax at token %d: 'name value =' requires spaces around =", i)
		}

		// Check if it's a number
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			if stack.Len() >= r.maxStack {
				return "", fmt.Errorf("stack overflow")
			}
			stack.Push(num)
			continue
		}

		// Handle special operators and commands
		if result, err := r.handleOperator(stack, token, i); err != nil {
			return "", err
		} else if result != "" {
			return result, nil
		}
	}

	// Check final stack state
	if stack.Len() == 0 {
		return "", fmt.Errorf("empty result: expression evaluated to nothing")
	}

	// Save the current stack state for continued operations
	// Create a copy of the stack to preserve it
	r.currentStack = NewStack()
	for _, val := range stack.Values() {
		r.currentStack.Push(val)
	}

	// Get the final result
	if stack.Len() > 1 {
		// Multiple values on stack - show them all
		result, err := r.ops.Show(stack)
		if err != nil {
			return "", fmt.Errorf("final result: %w", err)
		}
		return result, nil
	}

	// Single value - return it
	val, _ := stack.Pop()
	return fmt.Sprintf("%.10g", val), nil
}

// handleOperator handles operators and special commands
func (r *RPN) handleOperator(stack *Stack, token string, tokenIndex int) (string, error) {
	// Handle hyperoperators
	if isHyperOperator(token) {
		if err := r.handleHyperOperator(stack, token); err != nil {
			return "", err
		}
		return "", nil
	}

	// Handle standard operators
	switch token {
	case "+":
		if err := r.ops.Add(stack); err != nil {
			return "", fmt.Errorf("operator +: %w", err)
		}
	case "-":
		if err := r.ops.Subtract(stack); err != nil {
			return "", fmt.Errorf("operator -: %w", err)
		}
	case "*":
		if err := r.ops.Multiply(stack); err != nil {
			return "", fmt.Errorf("operator *: %w", err)
		}
	case "/":
		if err := r.ops.Divide(stack); err != nil {
			return "", fmt.Errorf("operator /: %w", err)
		}
	case "^":
		if err := r.ops.Power(stack); err != nil {
			return "", fmt.Errorf("operator ^: %w", err)
		}
	case "%":
		if err := r.ops.Modulo(stack); err != nil {
			return "", fmt.Errorf("operator %%: %w", err)
		}
	case "dup":
		if err := r.ops.Dup(stack); err != nil {
			return "", fmt.Errorf("dup: %w", err)
		}
	case "swap":
		if err := r.ops.Swap(stack); err != nil {
			return "", fmt.Errorf("swap: %w", err)
		}
	case "pop":
		if err := r.ops.Pop(stack); err != nil {
			return "", fmt.Errorf("pop: %w", err)
		}
	case "show", "showstack", "print":
		result, err := r.ops.Show(stack)
		if err != nil {
			return "", fmt.Errorf("show: %w", err)
		}
		return result, nil
	case "vars":
		result, err := r.ops.ListVariables()
		if err != nil {
			return "", fmt.Errorf("vars: %w", err)
		}
		return result, nil
	case "clear":
		r.ops.ClearVariables()
		return "All variables cleared", nil
	case "d":
		return "", fmt.Errorf("'d' command not supported as standalone token")
	default:
		// Check if it's a variable reference (push its value)
		val, exists := r.vars.GetVariable(token)
		if exists {
			stack.Push(val)
		} else {
			return "", fmt.Errorf("unknown token '%s' at position %d", token, tokenIndex)
		}
	}
	return "", nil
}

// isHyperOperator checks if the token is a hyperoperator
func isHyperOperator(token string) bool {
	switch token {
	case "[+]", "[-]", "[*]", "[/]", "[^]", "[%]":
		return true
	default:
		return false
	}
}

// handleHyperOperator handles hyperoperators
func (r *RPN) handleHyperOperator(stack *Stack, token string) error {
	switch token {
	case "[+]":
		if err := r.ops.HyperAdd(stack); err != nil {
			return fmt.Errorf("hyperoperator [+]: %w", err)
		}
	case "[-]":
		if err := r.ops.HyperSubtract(stack); err != nil {
			return fmt.Errorf("hyperoperator [-]: %w", err)
		}
	case "[*]":
		if err := r.ops.HyperMultiply(stack); err != nil {
			return fmt.Errorf("hyperoperator [*]: %w", err)
		}
	case "[/]":
		if err := r.ops.HyperDivide(stack); err != nil {
			return fmt.Errorf("hyperoperator [/]: %w", err)
		}
	case "[^]":
		if err := r.ops.HyperPower(stack); err != nil {
			return fmt.Errorf("hyperoperator [^]: %w", err)
		}
	case "[%]":
		if err := r.ops.HyperModulo(stack); err != nil {
			return fmt.Errorf("hyperoperator [%%]: %w", err)
		}
	}
	return nil
}

// handleAssignment checks if the input is an assignment format and handles it.
// Returns (result string, isAssignment bool, error error).
func (r *RPN) handleAssignment(input string) (string, bool, error) {
	// Check for assignment format (name = value or name value = expression)
	// We look for either " = " (with trailing space) or " =" (just space before equals)
	hasAssignment := strings.Contains(input, " = ") || strings.Contains(input, " =")
	if !hasAssignment {
		return "", false, nil
	}

	// Handle single assignment: "name = value"
	if parts := strings.SplitN(input, " = ", 2); len(parts) == 2 {
		name := strings.TrimSpace(parts[0])
		valueStr := strings.TrimSpace(parts[1])

		// Validate name is a single word (variable name)
		nameFields := strings.Fields(name)
		if len(nameFields) == 1 {
			// Validate value is a single number
			valueFields := strings.Fields(valueStr)
			if len(valueFields) == 1 {
				val, err := strconv.ParseFloat(valueFields[0], 64)
				if err != nil {
					return "", false, fmt.Errorf("invalid value '%s' for assignment: %w", valueFields[0], err)
				}
				if err := r.vars.SetVariable(nameFields[0], val); err != nil {
					return "", false, err
				}
				return fmt.Sprintf("%s = %.10g", nameFields[0], val), true, nil
			}
		}
	}

	// Handle assignment with expression: "name value = expression..."
	// Use " =" (space before equals) to find the boundary
	pos := strings.Index(input, " =")
	if pos >= 0 {
		// Extract content before the assignment
		before := strings.TrimSpace(input[:pos])
		// Extract content after " =" (may be empty or contain expression)
		after := strings.TrimSpace(input[pos+2:])

		beforeFields := strings.Fields(before)
		if len(beforeFields) == 2 {
			name := beforeFields[0]
			valueStr := beforeFields[1]

			// Try to parse value as a number
			val, err := strconv.ParseFloat(valueStr, 64)
			if err == nil {
				// Valid assignment pattern: "name value = expr..." or "name value ="
				if err := r.vars.SetVariable(name, val); err != nil {
					return "", false, err
				}

				// If no expression after assignment, just return assignment info
				if after == "" {
					return fmt.Sprintf("%s = %.10g", name, val), true, nil
				}
				result, err := r.evaluate(strings.Fields(after))
				return result, true, err
			}
		}
	}

	return "", false, nil
}
