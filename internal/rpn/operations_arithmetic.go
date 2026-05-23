// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"math"
)

// binaryMetricOp executes a metric-aware binary operation.
// It pops two operands, validates metric compatibility, converts to base units,
// applies the computation, converts back, and pushes the result.
//
// Parameters:
//   - op: operator name for error messages
//   - compatCheck: metric compatibility validator (nil to skip, e.g., for Multiply)
//   - compute: the arithmetic function to apply on base-unit values
//   - resultMetricFn: function to determine the result metric from input metrics
func (o *Operations) binaryMetricOp(
	stack *Stack,
	op string,
	compatCheck func(*Metric, *Metric) error,
	compute func(float64, float64) float64,
	resultMetricFn func(*MetricRegistry, *Metric, *Metric) (*Metric, error),
) error {
	a, b, err := popTwo(stack, op)
	if err != nil {
		return err
	}

	aM, err := resolveMetric(o.metricRegistry, a)
	if err != nil {
		return buildError(op, err)
	}
	bM, err := resolveMetric(o.metricRegistry, b)
	if err != nil {
		return buildError(op, err)
	}
	if compatCheck != nil {
		if err := compatCheck(aM, bM); err != nil {
			return err
		}
	}

	pm := o.GetPrefixMode()
	resultMetric, err := resultMetricFn(o.metricRegistry, aM, bM)
	if err != nil {
		return buildError(op, err)
	}
	aBase, err := convertToBase(o.metricRegistry, a, pm, resultMetric)
	if err != nil {
		return buildError("arithmetic", err)
	}
	bBase, err := convertToBase(o.metricRegistry, b, pm, resultMetric)
	if err != nil {
		return buildError("arithmetic", err)
	}
	resultVal, err := convertFromBase(o.metricRegistry, compute(aBase, bBase), resultMetric, pm)
	if err != nil {
		return buildError(op, err)
	}

	stack.Push(NewNumberWithMetric(resultVal, o.GetMode(), resultMetric))
	return nil
}

// Add pops two values from stack, adds them, and pushes result.
func (o *Operations) Add(stack *Stack) error {
	return o.binaryMetricOp(
		stack, "+",
		func(aM, bM *Metric) error {
			if !categoriesCompatible(aM, bM) {
				return metricError("+", aM, bM)
			}
			return nil
		},
		func(a, b float64) float64 { return a + b },
		compatibleMetric,
	)
}

// Subtract pops two values from stack, subtracts (a - b), and pushes result.
func (o *Operations) Subtract(stack *Stack) error {
	return o.binaryMetricOp(
		stack, "-",
		func(aM, bM *Metric) error {
			if !categoriesCompatible(aM, bM) {
				return metricError("-", aM, bM)
			}
			return nil
		},
		func(a, b float64) float64 { return a - b },
		compatibleMetric,
	)
}

// Multiply pops two values from stack, multiplies them, and pushes result.
func (o *Operations) Multiply(stack *Stack) error {
	return o.binaryMetricOp(
		stack, "*",
		nil, // no compatibility check for multiplication
		func(a, b float64) float64 { return a * b },
		resultMetricForMul,
	)
}

// Divide pops two values from stack, divides (a / b), and pushes result.
func (o *Operations) Divide(stack *Stack) error {
	b, err := popStack(stack, "/")
	if err != nil {
		return err
	}

	bF, err := toFloat64(b, "/")
	if err != nil {
		return err
	}
	if bF == 0 {
		return buildError("/", fmt.Errorf("division by zero"))
	}

	a, err := popStack(stack, "/")
	if err != nil {
		return err
	}

	aM, err := resolveMetric(o.metricRegistry, a)
	if err != nil {
		return buildError("/", err)
	}
	bM, err := resolveMetric(o.metricRegistry, b)
	if err != nil {
		return buildError("/", err)
	}

	pm := o.GetPrefixMode()
	resultMetric, err := resultMetricForDiv(o.metricRegistry, aM, bM)
	if err != nil {
		return buildError("/", err)
	}
	aBase, err := convertToBase(o.metricRegistry, a, pm, resultMetric)
	if err != nil {
		return buildError("division", err)
	}
	bBase, err := convertToBase(o.metricRegistry, b, pm, resultMetric)
	if err != nil {
		return buildError("division", err)
	}
	resultVal, err := convertFromBase(o.metricRegistry, aBase/bBase, resultMetric, pm)
	if err != nil {
		return buildError("/", err)
	}

	stack.Push(NewNumberWithMetric(resultVal, o.GetMode(), resultMetric))
	return nil
}

// Power pops two values from stack, raises first to power of second (a ^ b), and pushes result.
//
// Result is always Cool (unitless). This is intentional: x^n has different physical
// units than x, so retaining the input metric would be misleading. Examples:
//
//	2hr 3 ^  → 8 (Cool, not 8hr³)
//	100Mbps 2 ^  → 10000 (Cool, not 10000Mbps)
//	3 4 ^  → 81 (Cool)
//
// Scalar exponents (metric=Cool) also produce Cool results. There is no concept of
// "preserving the base metric" for power operations.
func (o *Operations) Power(stack *Stack) error {
	a, b, err := popTwo(stack, "^")
	if err != nil {
		return err
	}

	aF, err := toFloat64(a, "power")
	if err != nil {
		return buildError("power", err)
	}
	bF, err := toFloat64(b, "power")
	if err != nil {
		return buildError("power", err)
	}

	stack.Push(NewNumberWithMetric(math.Pow(aF, bF), o.GetMode(), GetCoolMetric()))
	return nil
}

// Modulo pops two values from stack, computes modulo (a % b), and pushes result.
func (o *Operations) Modulo(stack *Stack) error {
	a, b, err := popTwo(stack, "%")
	if err != nil {
		return err
	}

	bF, err := toFloat64(b, "%")
	if err != nil {
		return err
	}
	if bF == 0 {
		return buildError("%", fmt.Errorf("modulo by zero"))
	}

	aM, err := resolveMetric(o.metricRegistry, a)
	if err != nil {
		return buildError("%", err)
	}
	bM, err := resolveMetric(o.metricRegistry, b)
	if err != nil {
		return buildError("%", err)
	}
	if !categoriesCompatible(aM, bM) {
		return metricError("%", aM, bM)
	}

	pm := o.GetPrefixMode()
	resultMetric, err := compatibleMetric(o.metricRegistry, aM, bM)
	if err != nil {
		return buildError("%", err)
	}
	aBase, err := convertToBase(o.metricRegistry, a, pm, resultMetric)
	if err != nil {
		return buildError("modulo", err)
	}
	bBase, err := convertToBase(o.metricRegistry, b, pm, resultMetric)
	if err != nil {
		return buildError("modulo", err)
	}
	resultVal, err := convertFromBase(o.metricRegistry, math.Mod(aBase, bBase), resultMetric, pm)
	if err != nil {
		return buildError("%", err)
	}

	stack.Push(NewNumberWithMetric(resultVal, o.GetMode(), resultMetric))
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

	bVal, err := toFloat64(b, "**")
	if err != nil {
		return buildError("**", fmt.Errorf("exponent must be a number: %w", err))
	}

	exp := int(bVal)
	if float64(exp) != bVal {
		return buildError("**", fmt.Errorf("exponent must be an integer, got %v", bVal))
	}

	aF, err := toFloat64(a, "**")
	if err != nil {
		return buildError("**", err)
	}

	// Result is unitless (Cool metric)
	if exp == 0 {
		stack.Push(NewNumberWithMetric(1, o.GetMode(), GetCoolMetric()))
		return nil
	}
	resultVal := binaryExponentiationFloat(aF, exp)
	stack.Push(NewNumberWithMetric(resultVal, o.GetMode(), GetCoolMetric()))
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
