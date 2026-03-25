// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
)

// BooleanOperations provides boolean comparison operator implementations.
type BooleanOperations struct {
}

// NewBooleanOperations creates a new BooleanOperations instance.
func NewBooleanOperations() *BooleanOperations {
	return &BooleanOperations{}
}

// GT pops two values from stack, compares (a > b), and pushes a boolean result.
func (o *BooleanOperations) GT(stack *Stack) error {
	b, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for gt: %w", err)
	}

	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for gt: %w", err)
	}

	stack.Push(NewFloatFromBool(a.Float64() > b.Float64()))
	return nil
}

// LT pops two values from stack, compares (a < b), and pushes a boolean result.
func (o *BooleanOperations) LT(stack *Stack) error {
	b, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for lt: %w", err)
	}

	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for lt: %w", err)
	}

	stack.Push(NewFloatFromBool(a.Float64() < b.Float64()))
	return nil
}

// GTE pops two values from stack, compares (a >= b), and pushes a boolean result.
func (o *BooleanOperations) GTE(stack *Stack) error {
	b, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for gte: %w", err)
	}

	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for gte: %w", err)
	}

	stack.Push(NewFloatFromBool(a.Float64() >= b.Float64()))
	return nil
}

// LTE pops two values from stack, compares (a <= b), and pushes a boolean result.
func (o *BooleanOperations) LTE(stack *Stack) error {
	b, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for lte: %w", err)
	}

	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for lte: %w", err)
	}

	stack.Push(NewFloatFromBool(a.Float64() <= b.Float64()))
	return nil
}

// EQ pops two values from stack, compares (a == b), and pushes a boolean result.
func (o *BooleanOperations) EQ(stack *Stack) error {
	b, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for eq: %w", err)
	}

	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for eq: %w", err)
	}

	stack.Push(NewFloatFromBool(a.Float64() == b.Float64()))
	return nil
}

// NEQ pops two values from stack, compares (a != b), and pushes a boolean result.
func (o *BooleanOperations) NEQ(stack *Stack) error {
	b, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for neq: %w", err)
	}

	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for neq: %w", err)
	}

	stack.Push(NewFloatFromBool(a.Float64() != b.Float64()))
	return nil
}
