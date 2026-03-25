// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"math"
)

// ArithmeticOperations provides arithmetic operator implementations.
type ArithmeticOperations struct {
	mode CalculationMode
}

// NewArithmeticOperations creates a new ArithmeticOperations instance.
func NewArithmeticOperations(mode CalculationMode) *ArithmeticOperations {
	return &ArithmeticOperations{mode: mode}
}

// Add pops two values from stack, adds them, and pushes result.
func (o *ArithmeticOperations) Add(stack *Stack) error {
	bVal, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for +: %w", err)
	}

	aVal, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for +: %w", err)
	}

	// Use the Number interface for arithmetic
	stack.Push(aVal.Add(bVal))
	return nil
}

// Subtract pops two values from stack, subtracts (a - b), and pushes result.
func (o *ArithmeticOperations) Subtract(stack *Stack) error {
	b, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for -: %w", err)
	}

	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for -: %w", err)
	}

	stack.Push(a.Sub(b))
	return nil
}

// Multiply pops two values from stack, multiplies them, and pushes result.
func (o *ArithmeticOperations) Multiply(stack *Stack) error {
	b, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for *: %w", err)
	}

	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for *: %w", err)
	}

	stack.Push(a.Mul(b))
	return nil
}

// Divide pops two values from stack, divides (a / b), and pushes result.
func (o *ArithmeticOperations) Divide(stack *Stack) error {
	b, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for /: %w", err)
	}

	if b.IsZero() {
		return fmt.Errorf("division by zero")
	}

	a, err2 := stack.Pop()
	if err2 != nil {
		return fmt.Errorf("insufficient operands for /: %w", err2)
	}

	result, err2 := a.Div(b)
	if err2 != nil {
		return fmt.Errorf("division error: %w", err2)
	}
	stack.Push(result)
	return nil
}

// Power pops two values from stack, raises first to power of second (a ^ b), and pushes result.
func (o *ArithmeticOperations) Power(stack *Stack) error {
	b, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for ^: %w", err)
	}

	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for ^: %w", err)
	}

	stack.Push(a.Pow(b))
	return nil
}

// Modulo pops two values from stack, computes modulo (a % b), and pushes result.
func (o *ArithmeticOperations) Modulo(stack *Stack) error {
	b, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for %%: %w", err)
	}

	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for %%: %w", err)
	}

	if b.IsZero() {
		return fmt.Errorf("modulo by zero")
	}

	result, err := a.Mod(b)
	if err != nil {
		return fmt.Errorf("modulo error: %w", err)
	}
	stack.Push(result)
	return nil
}

// Log2 pops one value from stack, computes log base 2 (log₂(a)), and pushes result.
func (o *ArithmeticOperations) Log2(stack *Stack) error {
	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for lg: %w", err)
	}

	// Check if value is zero or negative
	val := a.Float64()
	if val <= 0 {
		return fmt.Errorf("log2 undefined for non-positive numbers")
	}

	// Compute log2 using the number interface
	stack.Push(NewNumber(math.Log2(val), o.mode))
	return nil
}

// Log10 pops one value from stack, computes log base 10 (log₁₀(a)), and pushes result.
func (o *ArithmeticOperations) Log10(stack *Stack) error {
	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for log: %w", err)
	}

	// Check if value is zero or negative
	val := a.Float64()
	if val <= 0 {
		return fmt.Errorf("log10 undefined for non-positive numbers")
	}

	// Compute log10 using the number interface
	stack.Push(NewNumber(math.Log10(val), o.mode))
	return nil
}

// Ln pops one value from stack, computes natural log (ln(a)), and pushes result.
func (o *ArithmeticOperations) Ln(stack *Stack) error {
	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for ln: %w", err)
	}

	// Check if value is zero or negative
	val := a.Float64()
	if val <= 0 {
		return fmt.Errorf("ln undefined for non-positive numbers")
	}

	// Compute ln using the number interface
	stack.Push(NewNumber(math.Log(val), o.mode))
	return nil
}
