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

	aVal, err := a.Float64()
	if err != nil {
		return fmt.Errorf("failed to get float64 value for a: %w", err)
	}
	bVal, err := b.Float64()
	if err != nil {
		return fmt.Errorf("failed to get float64 value for b: %w", err)
	}

	stack.Push(NewFloatFromBool(aVal > bVal))
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

	aVal, err := a.Float64()
	if err != nil {
		return fmt.Errorf("failed to get float64 value for a: %w", err)
	}
	bVal, err := b.Float64()
	if err != nil {
		return fmt.Errorf("failed to get float64 value for b: %w", err)
	}

	stack.Push(NewFloatFromBool(aVal < bVal))
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

	aVal, err := a.Float64()
	if err != nil {
		return fmt.Errorf("failed to get float64 value for a: %w", err)
	}
	bVal, err := b.Float64()
	if err != nil {
		return fmt.Errorf("failed to get float64 value for b: %w", err)
	}

	stack.Push(NewFloatFromBool(aVal >= bVal))
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

	aVal, err := a.Float64()
	if err != nil {
		return fmt.Errorf("failed to get float64 value for a: %w", err)
	}
	bVal, err := b.Float64()
	if err != nil {
		return fmt.Errorf("failed to get float64 value for b: %w", err)
	}

	stack.Push(NewFloatFromBool(aVal <= bVal))
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

	aVal, err := a.Float64()
	if err != nil {
		return fmt.Errorf("failed to get float64 value for a: %w", err)
	}
	bVal, err := b.Float64()
	if err != nil {
		return fmt.Errorf("failed to get float64 value for b: %w", err)
	}

	stack.Push(NewFloatFromBool(aVal == bVal))
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

	aVal, err := a.Float64()
	if err != nil {
		return fmt.Errorf("failed to get float64 value for a: %w", err)
	}
	bVal, err := b.Float64()
	if err != nil {
		return fmt.Errorf("failed to get float64 value for b: %w", err)
	}

	stack.Push(NewFloatFromBool(aVal != bVal))
	return nil
}
