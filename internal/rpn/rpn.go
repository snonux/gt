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
	currentStack *Stack
}

// NewRPN creates a new RPN parser and evaluator with the given variable store.
func NewRPN(vars *Variables) *RPN {
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
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("empty expression")
	}

	// Reset the stack for fresh evaluation (each call should be independent)
	r.currentStack = NewStack()

	// Handle single assignment: "name value ="
	// This is when the entire input is just an assignment
	if strings.Contains(input, " = ") {
		parts := strings.SplitN(input, " = ", 2)
		if len(parts) == 2 {
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
						return "", fmt.Errorf("invalid value '%s' for assignment: %w", valueFields[0], err)
					}
					if err := r.vars.SetVariable(nameFields[0], val); err != nil {
						return "", err
					}
					return fmt.Sprintf("%s = %.10g", nameFields[0], val), nil
				}
			}
		}
	}

	// Handle assignment with expression: "name value = expression..."
	// Format: variable_name value = expression (where = comes after value)
	if strings.Contains(input, " = ") {
		// Check if the input matches pattern: "name value = expr..."
		// where name and value are single tokens, and = comes after value
		// For example: "x 5 = x x +" or "pi 3.14 = pi 2 *"
		
		// Find " = " position and split
		pos := strings.Index(input, " = ")
		if pos >= 0 {
			before := input[:pos]   // "name value"
			after := input[pos+3:]  // "expr..."
			
			beforeFields := strings.Fields(before)
			if len(beforeFields) == 2 {
				name := beforeFields[0]
				valueStr := beforeFields[1]
				
				// Try to parse value as a number
				val, err := strconv.ParseFloat(valueStr, 64)
				if err == nil {
					// Valid assignment pattern: "name value = expr..."
					if err := r.vars.SetVariable(name, val); err != nil {
						return "", err
					}
					
					// Evaluate the remaining expression
					return r.evaluate(strings.Fields(strings.TrimSpace(after)))
				}
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

	// Save the current stack state for continued operations
	// Create a copy of the stack to preserve it
	r.currentStack = NewStack()
	for _, val := range stack.Values() {
		r.currentStack.Push(val)
	}

	// Get the final result
	if stack.Len() > 1 {
		// Multiple values - show them all
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
