package rpn

import (
	"fmt"
	"strconv"
	"strings"
)

// CalculationMode represents the mode for number calculations.
type CalculationMode int

const (
	// FloatMode uses float64 for calculations (default).
	FloatMode CalculationMode = iota
	// RationalMode uses *big.Rat for precise rational calculations.
	RationalMode
)

// RPN represents the RPN parser and evaluator.
type RPN struct {
	vars         VariableStore
	ops          Operator
	opRegistry   *OperatorRegistry
	maxStack     int
	currentStack *Stack
	mode         CalculationMode
}

// NewRPN creates a new RPN parser and evaluator with the given variable store.
func NewRPN(vars VariableStore) *RPN {
	ops := NewOperations(vars)
	ops.SetMode(FloatMode) // Set default mode
	return &RPN{
		vars:         vars,
		ops:          ops,
		opRegistry:   NewOperatorRegistry(ops),
		maxStack:     1000, // Reasonable limit for RPN expressions
		currentStack: NewStack(),
		mode:         FloatMode, // Default mode
	}
}

// GetMode returns the current calculation mode.
func (r *RPN) GetMode() CalculationMode {
	return r.mode
}

// SetMode sets the calculation mode.
func (r *RPN) SetMode(mode CalculationMode) {
	r.mode = mode
	r.ops.SetMode(mode)
}

// ParseAndEvaluate parses and evaluates an RPN expression.
// Returns the result as a formatted string or an error.
func (r *RPN) ParseAndEvaluate(input string) (string, error) {
	// Validate input and initialize
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("rpn: empty expression")
	}
	if r.currentStack == nil {
		r.currentStack = NewStack()
	}

	// Handle assignment formats
	if assignmentResult, isAssignment, err := r.handleAssignment(input); err != nil {
		return "", fmt.Errorf("rpn: failed to handle assignment: %w", err)
	} else if isAssignment {
		return assignmentResult, nil
	}

	// Evaluate standard RPN expression
	tokens := Tokenize(input)
	if len(tokens) == 0 {
		return "", fmt.Errorf("rpn: no valid tokens found in input: %q", input)
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
			stack.Push(NewNumberValue(num))
			continue
		}

		// Handle operator (common logic from executeOperator)
		if result, handled, err := r.executeOperator(stack, token); err != nil {
			// If the error is not "unknown token", return it
			// Otherwise, fall through to check for variable
			if !strings.Contains(err.Error(), "unknown token") {
				return "", err
			}
		} else if handled {
			if result != "" {
				return result, nil
			}
			continue
		}

		// Check if it's a variable reference (push its value)
		val, exists := r.vars.GetVariable(token)
		if exists {
			stack.Push(NewNumberValue(val))
		} else {
			return "", fmt.Errorf("unknown token '%s'", token)
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

	// Handle operator (common logic from executeOperator)
	if result, handled, err := r.executeOperator(r.currentStack, op); err != nil {
		return "", err
	} else if handled {
		if result != "" {
			return result, nil
		}
		// For EvalOperator, show the stack after operation
		stackShow, err := r.ops.Show(r.currentStack)
		if err != nil {
			return "", fmt.Errorf("show stack: %w", err)
		}
		return stackShow, nil
	}

	return "", fmt.Errorf("unknown operator '%s'", op)
}

// GetCurrentStack returns a copy of the current stack for inspection.
func (r *RPN) GetCurrentStack() []float64 {
	if r.currentStack == nil {
		return nil
	}
	values := r.currentStack.Values()
	result := make([]float64, len(values))
	for i, v := range values {
		result[i] = v.Number()
	}
	return result
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
			return "", fmt.Errorf("rpn: invalid assignment syntax at token %d: 'name value =' requires spaces around =", i)
		}

		// Check if it's a number
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			if stack.Len() >= r.maxStack {
				return "", fmt.Errorf("stack overflow")
			}
			stack.Push(NewNumberValue(num))
			continue
		}

		// Handle special operators and commands
		if result, err := r.handleOperator(stack, token, i); err != nil {
			return "", fmt.Errorf("rpn: failed to handle operator '%s' at position %d: %w", token, i, err)
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
		r.currentStack.Push(NewNumberValue(toNumber(val)))
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
	return fmt.Sprintf("%.10g", toNumber(val)), nil
}

// handleOperator handles operators and special commands using the operator registry.
func (r *RPN) handleOperator(stack *Stack, token string, tokenIndex int) (string, error) {
	// Check if it's a number first
	if _, err := strconv.ParseFloat(token, 64); err == nil {
		return "", nil
	}

	// Check if it's a variable reference first (before operators)
	if val, exists := r.vars.GetVariable(token); exists {
		stack.Push(NewNumberValue(val))
		return "", nil
	}

	// Handle standard operators (common logic extracted for DRY)
	if result, handled, err := r.executeOperator(stack, token); err != nil {
		return "", err
	} else if handled {
		return result, nil
	}

	return "", fmt.Errorf("unknown token '%s'", token)
}

// executeOperator handles operator execution (standard or hyper) and returns (result string, handled bool, error error).
// This is a helper to avoid code duplication between handleOperator and EvalOperator.
func (r *RPN) executeOperator(stack *Stack, token string) (string, bool, error) {
	// Check for hyperoperators first
	if r.opRegistry.IsHyperOperator(token) {
		result, handled, err := r.opRegistry.HandleHyperOperator(stack, token)
		return result, handled, err
	}

	// Then check standard operators
	result, handled, err := r.opRegistry.HandleStandardOperator(stack, token)
	return result, handled, err
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
