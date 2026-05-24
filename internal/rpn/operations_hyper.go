// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"errors"
	"fmt"
	"math"
)

// Hyper operators - operate on all values on the stack

// HyperAdd pops all values from stack, adds them left-associative, and pushes result.
// Metric-aware: validates all operands share the same category (Cool absorbs), converts to base units for the
// computation, and pushes the result with the first non-Cool metric (or Cool if all are Cool).
func (o *Operations) HyperAdd(stack *Stack) error {
	values, err := popAll(stack, "[+]")
	if err != nil {
		return err
	}

	// Resolve metrics for all values
	metrics := make([]*Metric, len(values))
	for i, v := range values {
		metrics[i] = resolveMetric(o.metricRegistry, v)
	}

	// Validate all are compatible (all same category, or Cool absorbs)
	if err := validateCategories(metrics, "[+]"); err != nil {
		return err
	}

	// Result metric: first non-Cool metric (Cool absorbs), or Cool
	resultMetric := resultMetricForAdd(metrics)

	// Convert all to base units, sum, convert back
	pm := o.GetPrefixMode()
	var sum float64
	for i, v := range values {
		base, err := convertToBase(o.metricRegistry, v, pm, resultMetric)
		if err != nil {
			return buildError("[+]", fmt.Errorf("operand %d: %w", i, err))
		}
		sum += base
	}

	resultVal := convertFromBase(o.metricRegistry, sum, resultMetric, pm)

	stack.Push(NewNumberWithMetric(resultVal, o.GetMode(), resultMetric))
	return nil
}

// HyperMultiply pops all values from stack, multiplies them left-associative, and pushes result.
// No metric validation; uses raw float64 values. Result is always Cool (unitless).
func (o *Operations) HyperMultiply(stack *Stack) error {
	values, err := popAll(stack, "[*]")
	if err != nil {
		return err
	}

	var product float64 = 1
	for i, v := range values {
		val, err := toFloat64(v, "[*]")
		if err != nil {
			return buildError("[*]", fmt.Errorf("operand %d: %w", i, err))
		}
		if i == 0 {
			product = val
		} else {
			product *= val
		}
	}

	cool := coolMetric(o.metricRegistry)
	stack.Push(NewNumberWithMetric(product, o.GetMode(), cool))
	return nil
}

// HyperSubtract pops all values from stack, subtracts them left-associative, and pushes result.
// Metric-aware: validates all operands share the same category (Cool absorbs), converts to base units for the
// computation, and pushes the result with the first non-Cool metric (or Cool if all are Cool).
func (o *Operations) HyperSubtract(stack *Stack) error {
	values, err := popAll(stack, "[-]")
	if err != nil {
		return err
	}

	// Resolve metrics for all values
	metrics := make([]*Metric, len(values))
	for i, v := range values {
		metrics[i] = resolveMetric(o.metricRegistry, v)
	}

	// Validate all are compatible (all same category, or Cool absorbs)
	if err := validateCategories(metrics, "[-]"); err != nil {
		return err
	}

	// Result metric: first non-Cool metric (Cool absorbs), or Cool
	resultMetric := resultMetricForAdd(metrics)

	// Convert all to base units, subtract, convert back
	pm := o.GetPrefixMode()
	firstBase, err := convertToBase(o.metricRegistry, values[0], pm, resultMetric)
	if err != nil {
		return buildError("[-]", fmt.Errorf("operand 0: %w", err))
	}
	result := firstBase
	for i := 1; i < len(values); i++ {
		base, err := convertToBase(o.metricRegistry, values[i], pm, resultMetric)
		if err != nil {
			return buildError("[-]", fmt.Errorf("operand %d: %w", i, err))
		}
		result -= base
	}

	resultVal := convertFromBase(o.metricRegistry, result, resultMetric, pm)

	stack.Push(NewNumberWithMetric(resultVal, o.GetMode(), resultMetric))
	return nil
}

// HyperDivide pops all values from stack, divides them left-associative, and pushes result.
// No metric validation; uses raw float64 values. Result is always Cool (unitless).
func (o *Operations) HyperDivide(stack *Stack) error {
	values, err := popAll(stack, "[/]")
	if err != nil {
		return err
	}

	firstVal, err := toFloat64(values[0], "[/]")
	if err != nil {
		return buildError("[/]", fmt.Errorf("operand 0: %w", err))
	}
	result := firstVal
	for i := 1; i < len(values); i++ {
		val, err := toFloat64(values[i], "[/]")
		if err != nil {
			return buildError("[/]", fmt.Errorf("operand %d: %w", i, err))
		}
		if val == 0 {
			return buildError("[/]", fmt.Errorf("division by zero at operand %d", i))
		}
		result /= val
	}

	cool := coolMetric(o.metricRegistry)
	stack.Push(NewNumberWithMetric(result, o.GetMode(), cool))
	return nil
}

// HyperPower pops all values from stack, raises to power left-associative, and pushes result.
// No metric validation; uses raw float64 values. Result is always Cool (unitless).
func (o *Operations) HyperPower(stack *Stack) error {
	values, err := popAll(stack, "[^]")
	if err != nil {
		return err
	}

	firstVal, err := toFloat64(values[0], "[^]")
	if err != nil {
		return buildError("[^]", fmt.Errorf("operand 0: %w", err))
	}
	result := firstVal
	for i := 1; i < len(values); i++ {
		val, err := toFloat64(values[i], "[^]")
		if err != nil {
			return buildError("[^]", fmt.Errorf("operand %d: %w", i, err))
		}
		result = math.Pow(result, val)
	}

	cool := coolMetric(o.metricRegistry)
	stack.Push(NewNumberWithMetric(result, o.GetMode(), cool))
	return nil
}

// HyperModulo pops all values from stack, computes modulo left-associative, and pushes result.
// Metric-aware: validates all operands share the same category (Cool absorbs), converts to base units for the
// computation, and pushes the result with the first non-Cool metric (or Cool if all are Cool).
func (o *Operations) HyperModulo(stack *Stack) error {
	values, err := popAll(stack, "[%]")
	if err != nil {
		return err
	}

	// Resolve metrics for all values
	metrics := make([]*Metric, len(values))
	for i, v := range values {
		metrics[i] = resolveMetric(o.metricRegistry, v)
	}

	// Validate all are compatible (all same category, or Cool absorbs)
	if err := validateCategories(metrics, "[%]"); err != nil {
		return err
	}

	// Result metric: first non-Cool metric (Cool absorbs), or Cool
	resultMetric := resultMetricForAdd(metrics)

	// Convert all to base units, compute modulo, convert back
	pm := o.GetPrefixMode()
	firstBase, err := convertToBase(o.metricRegistry, values[0], pm, resultMetric)
	if err != nil {
		return buildError("[%]", fmt.Errorf("operand 0: %w", err))
	}
	result := firstBase
	for i := 1; i < len(values); i++ {
		base, err := convertToBase(o.metricRegistry, values[i], pm, resultMetric)
		if err != nil {
			return buildError("[%]", fmt.Errorf("operand %d: %w", i, err))
		}
		if base == 0 {
			return buildError("[%]", fmt.Errorf("modulo by zero at operand %d", i))
		}
		result = math.Mod(result, base)
	}

	resultVal := convertFromBase(o.metricRegistry, result, resultMetric, pm)

	stack.Push(NewNumberWithMetric(resultVal, o.GetMode(), resultMetric))
	return nil
}

// hyperLog computes the sum of a log function over all stack values.
// Each value must be positive. Result is pushed with Cool metric.
func (o *Operations) hyperLog(stack *Stack, opName string, logFn func(float64) float64, errMsg string) error {
	values, err := popAll(stack, opName)
	if err != nil {
		return err
	}

	var result float64
	for i, v := range values {
		val, err := toFloat64(v, opName)
		if err != nil {
			return buildError(opName, fmt.Errorf("operand %d: %w", i, err))
		}
		if val <= 0 {
			return buildError(opName, errors.New(errMsg))
		}
		result += logFn(val)
	}

	cool := coolMetric(o.metricRegistry)
	stack.Push(NewNumberWithMetric(result, o.GetMode(), cool))
	return nil
}

// HyperLog2 pops all values from stack, computes sum of log2 for all values, and pushes result.
// No metric validation; uses raw float64 values. Result is always Cool (unitless).
func (o *Operations) HyperLog2(stack *Stack) error {
	return o.hyperLog(stack, "[lg]", math.Log2, "log2 undefined for non-positive numbers")
}

// HyperLog10 pops all values from stack, computes sum of log10 for all values, and pushes result.
// No metric validation; uses raw float64 values. Result is always Cool (unitless).
func (o *Operations) HyperLog10(stack *Stack) error {
	return o.hyperLog(stack, "[log]", math.Log10, "log10 undefined for non-positive numbers")
}

// HyperLn pops all values from stack, computes sum of natural log for all values, and pushes result.
// No metric validation; uses raw float64 values. Result is always Cool (unitless).
func (o *Operations) HyperLn(stack *Stack) error {
	return o.hyperLog(stack, "[ln]", math.Log, "ln undefined for non-positive numbers")
}
