// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"math"
)

// HyperOperations provides hyper operator implementations.
type HyperOperations struct {
	mode CalculationMode
}

// NewHyperOperations creates a new HyperOperations instance.
func NewHyperOperations(mode CalculationMode) *HyperOperations {
	return &HyperOperations{mode: mode}
}

// HyperAdd pops all values from stack, sums them, and pushes result.
func (o *HyperOperations) HyperAdd(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperadd: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Number
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return fmt.Errorf("hyperadd: %w", err)
		}
		values = append(values, val)
	}

	// Reverse to get left-to-right order (first pushed = first in)
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	// Process left-associative with Number interface
	sum := 0.0
	for i := 0; i < len(values); i++ {
		sum += values[i].Float64()
	}
	stack.Push(NewNumber(sum, o.mode))
	return nil
}

// HyperMultiply pops all values from stack, multiplies them left-associative, and pushes result.
func (o *HyperOperations) HyperMultiply(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hypermultiply: need at least 2 values")
	}

	product := 1.0
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return fmt.Errorf("hypermultiply: %w", err)
		}
		product *= val.Float64()
	}
	stack.Push(NewNumber(product, o.mode))
	return nil
}

// HyperSubtract pops all values from stack, subtracts them left-associative, and pushes result.
func (o *HyperOperations) HyperSubtract(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hypersubtract: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Number
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return fmt.Errorf("hypersubtract: %w", err)
		}
		values = append(values, val)
	}

	// Reverse to get left-to-right order (first pushed = first in)
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	// Process left-associative with Number interface
	result := values[0].Float64()
	for i := 1; i < len(values); i++ {
		result -= values[i].Float64()
	}
	stack.Push(NewNumber(result, o.mode))
	return nil
}

// HyperDivide pops all values from stack, divides them left-associative, and pushes result.
func (o *HyperOperations) HyperDivide(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperdivide: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Number
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return fmt.Errorf("hyperdivide: %w", err)
		}
		values = append(values, val)
	}

	// Reverse to get left-to-right order (first pushed = first in)
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	// Process left-associative with Number interface
	result := values[0].Float64()
	for i := 1; i < len(values); i++ {
		val := values[i].Float64()
		if val == 0 {
			return fmt.Errorf("division by zero")
		}
		result /= val
	}
	stack.Push(NewNumber(result, o.mode))
	return nil
}

// HyperPower pops all values from stack, raises to power left-associative, and pushes result.
func (o *HyperOperations) HyperPower(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperpower: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Number
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return fmt.Errorf("hyperpower: %w", err)
		}
		values = append(values, val)
	}

	// Reverse to get left-to-right order (first pushed = first in)
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	// Process left-associative with Number interface
	result := values[0].Float64()
	for i := 1; i < len(values); i++ {
		result = math.Pow(result, values[i].Float64())
	}
	stack.Push(NewNumber(result, o.mode))
	return nil
}

// HyperModulo pops all values from stack, computes modulo left-associative, and pushes result.
func (o *HyperOperations) HyperModulo(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hypermodulo: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Number
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return fmt.Errorf("hypermodulo: %w", err)
		}
		values = append(values, val)
	}

	// Reverse to get left-to-right order (first pushed = first in)
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	// Process left-associative with Number interface
	result := values[0].Float64()
	for i := 1; i < len(values); i++ {
		val := values[i].Float64()
		if val == 0 {
			return fmt.Errorf("modulo by zero")
		}
		result = math.Mod(result, val)
	}
	stack.Push(NewNumber(result, o.mode))
	return nil
}

// HyperLog2 pops all values from stack, computes sum of log2 for all values, and pushes result.
// This follows the same pattern as HyperAdd (sum) and HyperMultiply (product).
func (o *HyperOperations) HyperLog2(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperlog2: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Number
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return fmt.Errorf("hyperlog2: %w", err)
		}
		values = append(values, val)
	}

	// Reverse to get left-to-right order (first pushed = first in)
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	// Sum the log2 of all values with Number interface
	var result float64 = 0
	for i := 0; i < len(values); i++ {
		val := values[i].Float64()
		if val <= 0 {
			return fmt.Errorf("hyperlog2 undefined for non-positive numbers")
		}
		result += math.Log2(val)
	}

	// Push the result as a Number
	stack.Push(NewNumber(result, o.mode))
	return nil
}

// HyperLog10 pops all values from stack, computes sum of log10 for all values, and pushes result.
// This follows the same pattern as HyperAdd (sum) and HyperMultiply (product).
func (o *HyperOperations) HyperLog10(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperlog10: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Number
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return fmt.Errorf("hyperlog10: %w", err)
		}
		values = append(values, val)
	}

	// Reverse to get left-to-right order (first pushed = first in)
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	// Sum the log10 of all values
	var result float64 = 0
	for i := 0; i < len(values); i++ {
		val := values[i].Float64()
		if val <= 0 {
			return fmt.Errorf("hyperlog10 undefined for non-positive numbers")
		}
		result += math.Log10(val)
	}

	// Push the result as a Number
	stack.Push(NewNumber(result, o.mode))
	return nil
}

// HyperLn pops all values from stack, computes sum of natural log for all values, and pushes result.
// This follows the same pattern as HyperAdd (sum) and HyperMultiply (product).
func (o *HyperOperations) HyperLn(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperln: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Number
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return fmt.Errorf("hyperln: %w", err)
		}
		values = append(values, val)
	}

	// Reverse to get left-to-right order (first pushed = first in)
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	// Sum the natural log of all values with Number interface
	var result float64 = 0
	for i := 0; i < len(values); i++ {
		val := values[i].Float64()
		if val <= 0 {
			return fmt.Errorf("hyperln undefined for non-positive numbers")
		}
		result += math.Log(val)
	}
	stack.Push(NewNumber(result, o.mode))
	return nil
}
