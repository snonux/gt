// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"math"
)

// arithmetic operators

// Add pops two values from stack, adds them, and pushes result.
func (o *Operations) Add(stack *Stack) error {
	a, b, err := popTwo(stack, "+")
	if err != nil {
		return err
	}

	aM, bM := resolveMetric(o.metricRegistry, a), resolveMetric(o.metricRegistry, b)
	if !categoriesCompatible(aM, bM) {
		return metricError("+", aM, bM)
	}

		pm := o.GetPrefixMode()
		resultMetric := compatibleMetric(o.metricRegistry, aM, bM)
		// Convert both to base units, add, convert back to result metric
		aBase, err := convertToBase(o.metricRegistry, a, pm, resultMetric)
		if err != nil {
			return buildError("addition", err)
		}
		bBase, err := convertToBase(o.metricRegistry, b, pm, resultMetric)
		if err != nil {
			return buildError("addition", err)
		}
		resultVal := convertFromBase(o.metricRegistry, aBase+bBase, resultMetric, pm)

	stack.Push(NewNumber(resultVal, o.GetMode(), resultMetric))
	return nil
}

// Subtract pops two values from stack, subtracts (a - b), and pushes result.
func (o *Operations) Subtract(stack *Stack) error {
	a, b, err := popTwo(stack, "-")
	if err != nil {
		return err
	}

	aM, bM := resolveMetric(o.metricRegistry, a), resolveMetric(o.metricRegistry, b)
	if !categoriesCompatible(aM, bM) {
		return metricError("-", aM, bM)
	}

	pm := o.GetPrefixMode()
	resultMetric := compatibleMetric(o.metricRegistry, aM, bM)
	aBase, err := convertToBase(o.metricRegistry, a, pm, resultMetric)
	if err != nil {
		return buildError("subtraction", err)
	}
	bBase, err := convertToBase(o.metricRegistry, b, pm, resultMetric)
	if err != nil {
		return buildError("subtraction", err)
	}
	resultVal := convertFromBase(o.metricRegistry, aBase-bBase, resultMetric, pm)

	stack.Push(NewNumber(resultVal, o.GetMode(), resultMetric))
	return nil
}

// Multiply pops two values from stack, multiplies them, and pushes result.
func (o *Operations) Multiply(stack *Stack) error {
	a, b, err := popTwo(stack, "*")
	if err != nil {
		return err
	}

	aM, bM := resolveMetric(o.metricRegistry, a), resolveMetric(o.metricRegistry, b)

	pm := o.GetPrefixMode()
	resultMetric := resultMetricForMul(o.metricRegistry, aM, bM)
	// Convert both to base units, multiply, convert back to result metric
	aBase, err := convertToBase(o.metricRegistry, a, pm, resultMetric)
	if err != nil {
		return buildError("multiplication", err)
	}
	bBase, err := convertToBase(o.metricRegistry, b, pm, resultMetric)
	if err != nil {
		return buildError("multiplication", err)
	}
	resultVal := convertFromBase(o.metricRegistry, aBase*bBase, resultMetric, pm)

	stack.Push(NewNumber(resultVal, o.GetMode(), resultMetric))
	return nil
}

// Divide pops two values from stack, divides (a / b), and pushes result.
func (o *Operations) Divide(stack *Stack) error {
	b, err := popStack(stack, "/")
	if err != nil {
		return err
	}

	if b.IsZero() {
		return buildError("/", fmt.Errorf("division by zero"))
	}

	a, err := popStack(stack, "/")
	if err != nil {
		return err
	}

	aM, bM := resolveMetric(o.metricRegistry, a), resolveMetric(o.metricRegistry, b)

	pm := o.GetPrefixMode()
	resultMetric := resultMetricForDiv(o.metricRegistry, aM, bM)
	aBase, err := convertToBase(o.metricRegistry, a, pm, resultMetric)
	if err != nil {
		return buildError("division", err)
	}
	bBase, err := convertToBase(o.metricRegistry, b, pm, resultMetric)
	if err != nil {
		return buildError("division", err)
	}
	resultVal := convertFromBase(o.metricRegistry, aBase/bBase, resultMetric, pm)

	stack.Push(NewNumber(resultVal, o.GetMode(), resultMetric))
	return nil
}

// Power pops two values from stack, raises first to power of second (a ^ b), and pushes result.
// Result is unitless (Cool metric).
func (o *Operations) Power(stack *Stack) error {
	a, b, err := popTwo(stack, "^")
	if err != nil {
		return err
	}

	aF, err := a.Float64()
	if err != nil {
		return buildError("power", err)
	}
	bF, err := b.Float64()
	if err != nil {
		return buildError("power", err)
	}

	stack.Push(NewNumber(math.Pow(aF, bF), o.GetMode(), GetCoolMetric()))
	return nil
}

// Modulo pops two values from stack, computes modulo (a % b), and pushes result.
func (o *Operations) Modulo(stack *Stack) error {
	a, b, err := popTwo(stack, "%")
	if err != nil {
		return err
	}

	if sym, ok := a.(*Symbol); ok {
		return fmt.Errorf("symbol %s cannot be used with modulo operator", sym.Name())
	}
	if sym, ok := b.(*Symbol); ok {
		return fmt.Errorf("symbol %s cannot be used with modulo operator", sym.Name())
	}

	if b.IsZero() {
		return buildError("%", fmt.Errorf("modulo by zero"))
	}

	aM, bM := resolveMetric(o.metricRegistry, a), resolveMetric(o.metricRegistry, b)
	if !categoriesCompatible(aM, bM) {
		return metricError("%", aM, bM)
	}

	pm := o.GetPrefixMode()
	resultMetric := compatibleMetric(o.metricRegistry, aM, bM)
	aBase, err := convertToBase(o.metricRegistry, a, pm, resultMetric)
	if err != nil {
		return buildError("modulo", err)
	}
	bBase, err := convertToBase(o.metricRegistry, b, pm, resultMetric)
	if err != nil {
		return buildError("modulo", err)
	}
	resultVal := convertFromBase(o.metricRegistry, math.Mod(aBase, bBase), resultMetric, pm)

	stack.Push(NewNumber(resultVal, o.GetMode(), resultMetric))
	return nil
}

// FastPower pops two values from stack, raises first to integer power of second (a ** b), and pushes result.
// Uses binary exponentiation for efficiency with large integer exponents.
func (o *Operations) FastPower(stack *Stack) error {
	b, err := popStack(stack, "**")
	if err != nil {
		return err
	}

	a, err := popStack(stack, "**")
	if err != nil {
		return err
	}

	bVal, err := b.Float64()
	if err != nil {
		return buildError("**", fmt.Errorf("exponent must be a number: %w", err))
	}

	exp := int(bVal)
	if float64(exp) != bVal {
		return buildError("**", fmt.Errorf("exponent must be an integer, got %v", bVal))
	}

	aF, err := a.Float64()
	if err != nil {
		return buildError("**", err)
	}

	// Result is unitless (Cool metric)
	if exp == 0 {
		stack.Push(NewNumber(1, o.GetMode(), GetCoolMetric()))
		return nil
	}
	resultVal := binaryExponentiationFloat(aF, exp)
	stack.Push(NewNumber(resultVal, o.GetMode(), GetCoolMetric()))
	return nil
}

// binaryExponentiationFloat computes base^exp using the square-and-multiply algorithm.
// Time Complexity: O(log exp)
// Space Complexity: O(1)
func binaryExponentiationFloat(base float64, exp int) float64 {
	if exp == 0 {
		return 1.0
	}

	// Handle negative exponents: base^-exp = 1 / (base^exp)
	if exp < 0 {
		return 1.0 / binaryExponentiationFloat(base, -exp)
	}

	res := 1.0
	for exp > 0 {
		// If exponent is odd, multiply result by current base
		if exp%2 == 1 {
			res *= base
		}
		// Square the base and divide exponent by 2
		base *= base
		exp /= 2
	}
	return res
}

// Log2 pops one value from stack, computes log base 2 (log₂(a)), and pushes result.
func (o *Operations) Log2(stack *Stack) error {
	a, err := popStack(stack, "lg")
	if err != nil {
		return err
	}

	val, err := toFloat64(a, "log2")
	if err != nil {
		return err
	}
	if val <= 0 {
		return buildError("lg", fmt.Errorf("log2 undefined for non-positive numbers"))
	}

	// Compute log2 using the number interface
	mode := o.GetMode()
	stack.Push(NewNumber(math.Log2(val), mode))
	return nil
}

// Log10 pops one value from stack, computes log base 10 (log₁₀(a)), and pushes result.
func (o *Operations) Log10(stack *Stack) error {
	a, err := popStack(stack, "log")
	if err != nil {
		return err
	}

	val, err := toFloat64(a, "log10")
	if err != nil {
		return err
	}
	if val <= 0 {
		return buildError("log", fmt.Errorf("log10 undefined for non-positive numbers"))
	}

	// Compute log10 using the number interface
	mode := o.GetMode()
	stack.Push(NewNumber(math.Log10(val), mode))
	return nil
}

// Ln pops one value from stack, computes natural log (ln(a)), and pushes result.
func (o *Operations) Ln(stack *Stack) error {
	a, err := popStack(stack, "ln")
	if err != nil {
		return err
	}

	val, err := toFloat64(a, "ln")
	if err != nil {
		return err
	}
	if val <= 0 {
		return buildError("ln", fmt.Errorf("ln undefined for non-positive numbers"))
	}

	// Compute ln using the number interface
	mode := o.GetMode()
	stack.Push(NewNumber(math.Log(val), mode))
	return nil
}
