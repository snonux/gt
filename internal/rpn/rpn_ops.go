// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"strings"
)

// Tokenize splits the input string into tokens (numbers, operators, variables).
// This is exported for testing purposes.
func Tokenize(input string) []string {
	// Standard RPN tokenization
	return strings.Fields(input)
}

// ResultStack returns the final stack state after evaluation.
// This is useful for commands that need to show the stack without consuming it.
func (r *RPN) ResultStack(tokens []string) (string, error) {
	r.mu.RLock()
	mode := r.ops.GetMode()
	r.mu.RUnlock()

	stack := NewStack()

	for _, token := range tokens {
		// Push literal values (numbers, booleans, metric values, @metric prefix)
		if pushed, err := r.pushLiteral(stack, token); err != nil {
			return "", err
		} else if pushed {
			continue
		}

		// Handle operator (common logic from executeOperator)
		if result, handled, err := r.executeOperator(stack, token); err != nil {
			// If not a recognized operator, fall through to check for variable
			if r.opRegistry.IsStandardOperator(token) || r.opRegistry.IsHyperOperator(token) {
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
			if mode == RationalMode {
				stack.Push(NewRat(val))
			} else {
				stack.Push(NewFloat(val))
			}
		} else {
			return "", fmt.Errorf("unknown token '%s'", token)
		}
	}

	return r.ops.Show(stack)
}

// EvalOperator evaluates a single operator on the current stack state.
// This allows incremental RPN operations, e.g., after ParseAndEvaluate("1 2 3 +") leaves
// the stack with [1, 5], calling EvalOperator("+") produces "6".
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

// executeOperator handles operator execution (standard or hyper) and returns (result string, handled bool, error error).
// This is a shared helper used by handleOperator, ResultStack, and EvalOperator.
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
