package rpn

import (
	"fmt"
	"strconv"
	"strings"
)

// RPN represents the RPN parser and evaluator.
type RPN struct {
	vars    *Variables
	ops     *Operations
	maxStack int
}

// NewRPN creates a new RPN parser and evaluator with the given variable store.
func NewRPN(vars *Variables) *RPN {
	return &RPN{
		vars:     vars,
		ops:      NewOperations(vars),
		maxStack: 1000, // Reasonable limit for RPN expressions
	}
}

// ParseAndEvaluate parses and evaluates an RPN expression.
// Returns the result as a formatted string or an error.
func (r *RPN) ParseAndEvaluate(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("empty expression")
	}

	// First check for special assignment pattern: "name value ="
	// This is handled separately from RPN evaluation
	if strings.Contains(input, " = ") {
		parts := strings.SplitN(input, " = ", 2)
		if len(parts) == 2 {
			name := strings.TrimSpace(parts[0])
			valueStr := strings.TrimSpace(parts[1])
			// Validate name is a single word (variable name)
			nameFields := strings.Fields(name)
			if len(nameFields) == 1 {
				// Parse the value
				val, err := strconv.ParseFloat(valueStr, 64)
				if err != nil {
					return "", fmt.Errorf("invalid value '%s' for assignment: %w", valueStr, err)
				}
				// Assign the variable
				if err := r.vars.SetVariable(nameFields[0], val); err != nil {
					return "", err
				}
				return fmt.Sprintf("%s = %.10g", nameFields[0], val), nil
			}
		}
	}

	tokens := tokenize(input)
	if len(tokens) == 0 {
		return "", fmt.Errorf("no valid tokens found")
	}

	return r.evaluate(tokens)
}

// tokenize splits the input string into tokens (numbers, operators, variables).
func tokenize(input string) []string {
	// Standard RPN tokenization
	return strings.Fields(input)
}

// evaluate evaluates a list of tokens and returns the result.
func (r *RPN) evaluate(tokens []string) (string, error) {
	stack := NewStack()

	for i, token := range tokens {
		// Check for variable assignment: name value =
		if token == "=" {
			// Assignment requires: variable_name (as previous token), value (on stack)
			// But tokens are processed linearly, so we need special handling
			// For "name value =", we have tokens: [name, value, =]
			// When we see =, we need to have the value on stack and name from before
			// This approach won't work well, so we handle assignment at parse time
			return "", fmt.Errorf("invalid assignment: '=' must be used with 'name value =' syntax")
		}

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
			// For show, we return the stack state instead of continuing
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
				return "", fmt.Errorf("unknown token '%s' at position %d", token, i)
			}
		}
	}

	// Check final stack state
	if stack.Len() == 0 {
		return "", fmt.Errorf("empty result: expression evaluated to nothing")
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

// ResultStack returns the final stack state after evaluation.
// This is useful for commands that need to show the stack without consuming it.
func (r *RPN) ResultStack(tokens []string) (string, error) {
	stack := NewStack()

	for _, token := range tokens {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack.Push(num)
			continue
		}

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
