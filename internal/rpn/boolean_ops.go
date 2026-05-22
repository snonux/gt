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

	aF, err := toFloat64(a, "gt")
	if err != nil {
		return err
	}
	bF, err := toFloat64(b, "gt")
	if err != nil {
		return err
	}

	stack.Push(NewFloatFromBool(aF > bF))
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

	aF, err := toFloat64(a, "lt")
	if err != nil {
		return err
	}
	bF, err := toFloat64(b, "lt")
	if err != nil {
		return err
	}

	stack.Push(NewFloatFromBool(aF < bF))
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

	aF, err := toFloat64(a, "gte")
	if err != nil {
		return err
	}
	bF, err := toFloat64(b, "gte")
	if err != nil {
		return err
	}

	stack.Push(NewFloatFromBool(aF >= bF))
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

	aF, err := toFloat64(a, "lte")
	if err != nil {
		return err
	}
	bF, err := toFloat64(b, "lte")
	if err != nil {
		return err
	}

	stack.Push(NewFloatFromBool(aF <= bF))
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

	aF, err := toFloat64(a, "eq")
	if err != nil {
		return err
	}
	bF, err := toFloat64(b, "eq")
	if err != nil {
		return err
	}

	stack.Push(NewFloatFromBool(aF == bF))
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

	aF, err := toFloat64(a, "neq")
	if err != nil {
		return err
	}
	bF, err := toFloat64(b, "neq")
	if err != nil {
		return err
	}

	stack.Push(NewFloatFromBool(aF != bF))
	return nil
}
