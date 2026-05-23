// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
)

// Helper functions to reduce error handling boilerplate in RPN operations

// popStack pops a value from the stack and returns a wrapped error if insufficient operands.
func popStack(stack *Stack, op string) (StackValue, error) {
	val, err := stack.Pop()
	if err != nil {
		return nil, fmt.Errorf("insufficient operands for %s: %w", op, err)
	}
	return val, nil
}

// popTwo pops two values from the stack for binary operations.
func popTwo(stack *Stack, op string) (StackValue, StackValue, error) {
	b, err := stack.Pop()
	if err != nil {
		return nil, nil, fmt.Errorf("insufficient operands for %s: %w", op, err)
	}

	a, err := stack.Pop()
	if err != nil {
		return nil, nil, fmt.Errorf("insufficient operands for %s: %w", op, err)
	}

	return a, b, nil
}

// toFloat64 converts a StackValue to float64 with proper error wrapping.
// Asserts the value to NumericValue; returns an error for non-numeric types.
func toFloat64(val StackValue, context string) (float64, error) {
	nv, ok := val.(NumericValue)
	if !ok {
		return 0, fmt.Errorf("%s: value %q is not numeric", context, val)
	}
	f, err := nv.Float64()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get float64 value: %w", context, err)
	}
	return f, nil
}

// ensureStackLength checks if the stack has at least min values and returns error if not.
func ensureStackLength(stack *Stack, min int, op string) error {
	if stack.Len() < min {
		return fmt.Errorf("insufficient operands for %s: need at least %d values", op, min)
	}
	return nil
}

// buildError wraps an error with context for the given operator.
func buildError(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}

// popAll pops all values from stack into a slice and reverses them for left-to-right processing.
// Returns values in order from bottom to top of stack (first pushed to last pushed).
func popAll(stack *Stack, op string) ([]StackValue, error) {
	if stack.Len() < 2 {
		return nil, fmt.Errorf("insufficient operands for %s: need at least 2 values", op)
	}

	var values []StackValue
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return nil, fmt.Errorf("%s: failed to pop: %w", op, err)
		}
		values = append(values, val)
	}

	// Reverse to get left-to-right order (first pushed = first in)
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	return values, nil
}
