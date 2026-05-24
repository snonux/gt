// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"errors"
	"fmt"
	"math"
)

// nAryMetricOp handles metric-aware n-ary operations.
// Resolves metrics, validates categories, converts to base units,
// applies fn left-associatively, converts back, and pushes the result.
func (o *Operations) nAryMetricOp(stack *Stack, opName string, fn func(float64, float64) float64) error {
	values, err := popAll(stack, opName)
	if err != nil {
		return err
	}
	if len(values) == 0 {
		return buildError(opName, fmt.Errorf("need at least one operand"))
	}

	metrics := make([]*Metric, len(values))
	for i, v := range values {
		m, err := resolveMetric(o.metricRegistry, v)
		if err != nil {
			return buildError(opName, err)
		}
		metrics[i] = m
	}

	if err := validateCategories(metrics, opName); err != nil {
		return err
	}

	resultMetric := resultMetricForAdd(metrics)
	pm := o.GetPrefixMode()

	firstBase, err := convertToBase(o.metricRegistry, values[0], pm, resultMetric)
	if err != nil {
		return buildError(opName, fmt.Errorf("operand 0: %w", err))
	}
	result := firstBase
	for i := 1; i < len(values); i++ {
		base, err := convertToBase(o.metricRegistry, values[i], pm, resultMetric)
		if err != nil {
			return buildError(opName, fmt.Errorf("operand %d: %w", i, err))
		}
		result = fn(result, base)
	}

	resultVal, err := convertFromBase(o.metricRegistry, result, resultMetric, pm)
	if err != nil {
		return buildError(opName, err)
	}

	stack.Push(NewNumberWithMetric(resultVal, o.GetMode(), resultMetric))
	return nil
}

// nAryScalarOp applies a binary function to pre-converted float64 values
// and pushes the result with Cool metric.
func (o *Operations) nAryScalarOp(stack *Stack, opName string, floats []float64, fn func(float64, float64) float64) error {
	if len(floats) == 0 {
		return buildError(opName, fmt.Errorf("need at least one operand"))
	}
	result := floats[0]
	for i := 1; i < len(floats); i++ {
		result = fn(result, floats[i])
	}

	cool, err := coolMetric(o.metricRegistry)
	if err != nil {
		return buildError(opName, err)
	}
	stack.Push(NewNumberWithMetric(result, o.GetMode(), cool))
	return nil
}

// HyperAdd pops all values from stack, adds them left-associative, and pushes result.
// Metric-aware: validates all operands share the same category (Cool absorbs), converts to base units for the
// computation, and pushes the result with the first non-Cool metric (or Cool if all are Cool).
func (o *Operations) HyperAdd(stack *Stack) error {
	return o.nAryMetricOp(stack, "[+]", func(a, b float64) float64 { return a + b })
}

// HyperMultiply pops all values from stack, multiplies them left-associative, and pushes result.
// No metric validation; uses raw float64 values. Result is always Cool (unitless).
func (o *Operations) HyperMultiply(stack *Stack) error {
	values, err := popAll(stack, "[*]")
	if err != nil {
		return err
	}

	floats := make([]float64, len(values))
	for i, v := range values {
		val, err := toFloat64(v, "[*]")
		if err != nil {
			return buildError("[*]", fmt.Errorf("operand %d: %w", i, err))
		}
		floats[i] = val
	}

	var product float64 = 1
	for i, val := range floats {
		if i == 0 {
			product = val
		} else {
			product *= val
		}
	}

	cool, err := coolMetric(o.metricRegistry)
	if err != nil {
		return buildError("[*]", err)
	}
	stack.Push(NewNumberWithMetric(product, o.GetMode(), cool))
	return nil
}

// HyperSubtract pops all values from stack, subtracts them left-associative, and pushes result.
// Metric-aware: validates all operands share the same category (Cool absorbs), converts to base units for the
// computation, and pushes the result with the first non-Cool metric (or Cool if all are Cool).
func (o *Operations) HyperSubtract(stack *Stack) error {
	return o.nAryMetricOp(stack, "[-]", func(a, b float64) float64 { return a - b })
}

// HyperDivide pops all values from stack, divides them left-associative, and pushes result.
// No metric validation; uses raw float64 values. Result is always Cool (unitless).
func (o *Operations) HyperDivide(stack *Stack) error {
	values, err := popAll(stack, "[/]")
	if err != nil {
		return err
	}

	floats := make([]float64, len(values))
	for i, v := range values {
		val, err := toFloat64(v, "[/]")
		if err != nil {
			return buildError("[/]", fmt.Errorf("operand %d: %w", i, err))
		}
		if i > 0 && val == 0 {
			return buildError("[/]", fmt.Errorf("division by zero at operand %d", i))
		}
		floats[i] = val
	}

	return o.nAryScalarOp(stack, "[/]", floats, func(a, b float64) float64 { return a / b })
}

// HyperPower pops all values from stack, raises to power left-associative, and pushes result.
// No metric validation; uses raw float64 values. Result is always Cool (unitless).
func (o *Operations) HyperPower(stack *Stack) error {
	values, err := popAll(stack, "[^]")
	if err != nil {
		return err
	}

	floats := make([]float64, len(values))
	for i, v := range values {
		val, err := toFloat64(v, "[^]")
		if err != nil {
			return buildError("[^]", fmt.Errorf("operand %d: %w", i, err))
		}
		floats[i] = val
	}

	return o.nAryScalarOp(stack, "[^]", floats, math.Pow)
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
		m, err := resolveMetric(o.metricRegistry, v)
		if err != nil {
			return buildError("[%]", err)
		}
		metrics[i] = m
	}

	if err := validateCategories(metrics, "[%]"); err != nil {
		return err
	}

	resultMetric := resultMetricForAdd(metrics)
	pm := o.GetPrefixMode()

	bases := make([]float64, len(values))
	for i, v := range values {
		base, err := convertToBase(o.metricRegistry, v, pm, resultMetric)
		if err != nil {
			return buildError("[%]", fmt.Errorf("operand %d: %w", i, err))
		}
		if i > 0 && base == 0 {
			return buildError("[%]", fmt.Errorf("modulo by zero at operand %d", i))
		}
		bases[i] = base
	}

	result := bases[0]
	for i := 1; i < len(bases); i++ {
		result = math.Mod(result, bases[i])
	}

	resultVal, err := convertFromBase(o.metricRegistry, result, resultMetric, pm)
	if err != nil {
		return buildError("[%]", err)
	}

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

	cool, err := coolMetric(o.metricRegistry)
	if err != nil {
		return buildError(opName, err)
	}
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
