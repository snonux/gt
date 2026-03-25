// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
)

// StackOperations provides stack manipulation operator implementations.
type StackOperations struct {
}

// NewStackOperations creates a new StackOperations instance.
func NewStackOperations() *StackOperations {
	return &StackOperations{}
}

// Dup duplicates the top stack value.
func (o *StackOperations) Dup(stack *Stack) error {
	val, err := stack.Peek()
	if err != nil {
		return fmt.Errorf("insufficient operands for dup: %w", err)
	}
	stack.Push(val)
	return nil
}

// Swap swaps the top two stack values.
func (o *StackOperations) Swap(stack *Stack) error {
	b, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for swap: %w", err)
	}

	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for swap: %w", err)
	}

	// Push in swapped order
	stack.Push(b)
	stack.Push(a)
	return nil
}

// Pop removes the top stack value.
func (o *StackOperations) Pop(stack *Stack) error {
	_, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for pop: %w", err)
	}
	return nil
}

// Show returns the current stack state as a string without modifying it.
func (o *StackOperations) Show(stack *Stack) (string, error) {
	if stack.Len() == 0 {
		return "", fmt.Errorf("empty stack")
	}
	// For now, just return the top value as a string
	// In a full implementation, this would show the entire stack
	val, err := stack.Peek()
	if err != nil {
		return "", err
	}
	return val.String(), nil
}
