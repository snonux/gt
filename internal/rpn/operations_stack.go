// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"strings"
)

// stack manipulation operators

// Dup duplicates the top stack value.
func (o *Operations) Dup(stack *Stack) error {
	val, err := stack.Peek()
	if err != nil {
		return buildError("dup", err)
	}
	stack.Push(val)
	return nil
}

// Swap swaps the top two stack values.
func (o *Operations) Swap(stack *Stack) error {
	if err := ensureStackLength(stack, 2, "swap"); err != nil {
		return err
	}

	// Get the values without popping
	vals := stack.Values()
	top := vals[len(vals)-1]
	second := vals[len(vals)-2]

	// Pop both values
	if _, err := stack.Pop(); err != nil {
		return buildError("swap", fmt.Errorf("failed to pop top value: %w", err))
	}
	if _, err := stack.Pop(); err != nil {
		return buildError("swap", fmt.Errorf("failed to pop second value: %w", err))
	}

	// Push in swapped order
	stack.Push(top)
	stack.Push(second)

	return nil
}

// Pop removes and discards the top stack value.
func (o *Operations) Pop(stack *Stack) error {
	if _, err := stack.Pop(); err != nil {
		return buildError("pop", err)
	}
	return nil
}

// Show returns the current stack as a formatted string.
// Each value uses its StackValue.String() method for display.
// Numeric values with non-Universal metrics append the metric suffix (e.g., "100Mbps").
// Values with the Universal metric (Cool) display without a suffix.
// Booleans display as "true"/"false". Symbols display with a leading ":" prefix.
func (o *Operations) Show(stack *Stack) (string, error) {
	if stack.Len() == 0 {
		return "Stack is empty", nil
	}

	vals := stack.Values()
	var sb strings.Builder
	for i, val := range vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		// Append metric suffix for non-Cool metrics
		m := val.Metric()
		if m != nil && m.Category != Universal {
			sb.WriteString(val.String())
			sb.WriteString(m.Name)
		} else {
			sb.WriteString(val.String())
		}
	}
	return sb.String(), nil
}
