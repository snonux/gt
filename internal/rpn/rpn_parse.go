// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseAndEvaluate parses and evaluates an RPN expression.
// Returns the result as a formatted string or an error.
// This method is thread-safe for concurrent execution.
func (r *RPN) ParseAndEvaluate(input string) (string, error) {
	// Validate input and initialize
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("rpn: empty expression")
	}

	// Lock for write operations on currentStack
	r.mu.Lock()
	if r.currentStack == nil {
		r.currentStack = NewStack()
	}
	r.mu.Unlock()

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

	return r.evaluate(input, tokens)
}

// evaluate evaluates a list of tokens and returns the result.
// This method is thread-safe for concurrent execution.
func (r *RPN) evaluate(input string, tokens []string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
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

		// Check if it's a boolean literal
		if token == "true" {
			stack.Push(NewFloatFromBool(true))
			continue
		}
		if token == "false" {
			stack.Push(NewFloatFromBool(false))
			continue
		}

		// Check if it's a number
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			if stack.Len() >= r.maxStack {
				return "", fmt.Errorf("stack overflow")
			}
			stack.Push(NewNumber(num, r.mode))
			continue
		}

		// Check if this is a variable name for assignment (:= or =:)
		// For := (right assignment): name value := - first token is always a variable name
		// For =: (left assignment): value name =: - token before =: is a variable name
		shouldPushName := false
		if i+1 < len(tokens) {
			nextToken := tokens[i+1]
			if nextToken == ":=" || nextToken == "=:" {
				// This token is a variable name
				// Only push as StringNum if it's not a number
				if _, err := strconv.ParseFloat(token, 64); err != nil {
					shouldPushName = true
				}
			}
		}
		
		// Special case: first token in := expression (e.g., "x 5 :=")
		if i == 0 && len(tokens) >= 3 && tokens[len(tokens)-1] == ":=" {
			shouldPushName = true
		}
		
		if shouldPushName {
			// This token is a variable name, push as StringNum
			stack.Push(NewStringNum(token))
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
		// Empty stack might be valid for assignment operators (:= or =:)
		// Check if the input was an assignment expression
		if strings.Contains(input, ":=") || strings.Contains(input, "=:") {
			// Assignment expression - empty stack is valid (side effect is variable assignment)
			return "", nil
		}
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
	return val.String(), nil
}

// handleOperator handles operators and special commands using the operator registry.
func (r *RPN) handleOperator(stack *Stack, token string, tokenIndex int) (string, error) {
	// Check if it's a number first
	if _, err := strconv.ParseFloat(token, 64); err == nil {
		return "", nil
	}

	// Check if it's a variable reference first (before operators)
	if val, exists := r.vars.GetVariable(token); exists {
		stack.Push(NewNumber(val, r.mode))
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

// handleAssignment checks if the input is an assignment format and handles it.
// Returns (result string, isAssignment bool, error error).
func (r *RPN) handleAssignment(input string) (string, bool, error) {
	// Handle := operator
	// Format 1: value name := (value on bottom, name on top)
	// Format 2: name value := (name on bottom, value on top) - for compatibility
	if strings.Contains(input, ":=") {
		pos := strings.Index(input, ":=")
		if pos >= 0 {
			before := strings.TrimSpace(input[:pos])
			after := strings.TrimSpace(input[pos+2:])

			beforeFields := strings.Fields(before)
			if len(beforeFields) == 2 {
				// Try value name := format first
				name := beforeFields[1]
				valueStr := beforeFields[0]

				val, err := strconv.ParseFloat(valueStr, 64)
				if err == nil {
					if err := r.vars.SetVariable(name, val); err != nil {
						return "", false, err
					}
					if after == "" {
						return fmt.Sprintf("%s = %.10g", name, val), true, nil
					}
					result, err := r.evaluate(input, strings.Fields(after))
					return result, true, err
				}

				// Try name value := format (for backward compatibility)
				name = beforeFields[0]
				valueStr = beforeFields[1]

				val, err = strconv.ParseFloat(valueStr, 64)
				if err == nil {
					if err := r.vars.SetVariable(name, val); err != nil {
						return "", false, err
					}
					if after == "" {
						return fmt.Sprintf("%s = %.10g", name, val), true, nil
					}
					result, err := r.evaluate(input, strings.Fields(after))
					return result, true, err
				}
			}
		}
	}

	// Handle =: operator
	// Format 1: value name =: (value on bottom, name on top)
	// Format 2: name value =: (name on bottom, value on top) - for compatibility
	if strings.Contains(input, "=:") {
		pos := strings.Index(input, "=:")
		if pos >= 0 {
			before := strings.TrimSpace(input[:pos])
			after := strings.TrimSpace(input[pos+2:])

			beforeFields := strings.Fields(before)
			if len(beforeFields) == 2 {
				// Try value name =: format first
				name := beforeFields[1]
				valueStr := beforeFields[0]

				val, err := strconv.ParseFloat(valueStr, 64)
				if err == nil {
					if err := r.vars.SetVariable(name, val); err != nil {
						return "", false, err
					}
					if after == "" {
						return fmt.Sprintf("%s = %.10g", name, val), true, nil
					}
					result, err := r.evaluate(input, strings.Fields(after))
					return result, true, err
				}

				// Try name value =: format (for backward compatibility)
				name = beforeFields[0]
				valueStr = beforeFields[1]

				val, err = strconv.ParseFloat(valueStr, 64)
				if err == nil {
					if err := r.vars.SetVariable(name, val); err != nil {
						return "", false, err
					}
					if after == "" {
						return fmt.Sprintf("%s = %.10g", name, val), true, nil
					}
					result, err := r.evaluate(input, strings.Fields(after))
					return result, true, err
				}
			}
		}
	}

	// Check for standard assignment format (name = value or name value = expression)
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
				result, err := r.evaluate(input, strings.Fields(after))
				return result, true, err
			}
		}
	}

	return "", false, nil
}
