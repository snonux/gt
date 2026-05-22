// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"math"
)

// Hyper operators - operate on all values on the stack

// HyperAdd pops all values from stack, adds them left-associative (with boolean-to-number coercion), and pushes result.
// Metric-aware: validates all operands share the same category (Cool absorbs), converts to base units for the
// computation, and pushes the result with the first operand's metric.
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
	if err := validateSameCategory(metrics, "[+]"); err != nil {
		return err
	}

	// Result metric: first non-Cool metric (Cool absorbs), or Cool
	resultMetric := resultMetricForHyperAdd(metrics)

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

	stack.Push(NewNumber(resultVal, o.GetMode(), resultMetric))
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
		val, err := v.Float64()
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
	stack.Push(NewNumber(product, o.GetMode(), cool))
	return nil
}

// HyperSubtract pops all values from stack, subtracts them left-associative, and pushes result.
// Metric-aware: validates all operands share the same category (Cool absorbs), converts to base units for the
// computation, and pushes the result with the first operand's metric.
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
	if err := validateSameCategory(metrics, "[-]"); err != nil {
		return err
	}

	// Result metric: first non-Cool metric (Cool absorbs), or Cool
	resultMetric := resultMetricForHyperAdd(metrics)

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

	stack.Push(NewNumber(resultVal, o.GetMode(), resultMetric))
	return nil
}

// HyperDivide pops all values from stack, divides them left-associative, and pushes result.
// No metric validation; uses raw float64 values. Result is always Cool (unitless).
func (o *Operations) HyperDivide(stack *Stack) error {
	values, err := popAll(stack, "[/]")
	if err != nil {
		return err
	}

	firstVal, err := values[0].Float64()
	if err != nil {
		return buildError("[/]", fmt.Errorf("operand 0: %w", err))
	}
	result := firstVal
	for i := 1; i < len(values); i++ {
		val, err := values[i].Float64()
		if err != nil {
			return buildError("[/]", fmt.Errorf("operand %d: %w", i, err))
		}
		if val == 0 {
			return buildError("[/]", fmt.Errorf("division by zero at operand %d", i))
		}
		result /= val
	}

	cool := coolMetric(o.metricRegistry)
	stack.Push(NewNumber(result, o.GetMode(), cool))
	return nil
}

// HyperPower pops all values from stack, raises to power left-associative, and pushes result.
// No metric validation; uses raw float64 values. Result is always Cool (unitless).
func (o *Operations) HyperPower(stack *Stack) error {
	values, err := popAll(stack, "[^]")
	if err != nil {
		return err
	}

	firstVal, err := values[0].Float64()
	if err != nil {
		return buildError("[^]", fmt.Errorf("operand 0: %w", err))
	}
	result := firstVal
	for i := 1; i < len(values); i++ {
		val, err := values[i].Float64()
		if err != nil {
			return buildError("[^]", fmt.Errorf("operand %d: %w", i, err))
		}
		result = math.Pow(result, val)
	}

	cool := coolMetric(o.metricRegistry)
	stack.Push(NewNumber(result, o.GetMode(), cool))
	return nil
}

// HyperModulo pops all values from stack, computes modulo left-associative, and pushes result.
// Metric-aware: validates all operands share the same category (Cool absorbs), converts to base units for the
// computation, and pushes the result with the first operand's metric.
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
	if err := validateSameCategory(metrics, "[%]"); err != nil {
		return err
	}

	// Result metric: first non-Cool metric (Cool absorbs), or Cool
	resultMetric := resultMetricForHyperAdd(metrics)

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

	stack.Push(NewNumber(resultVal, o.GetMode(), resultMetric))
	return nil
}

// HyperLog2 pops all values from stack, computes sum of log2 for all values, and pushes result.
// No metric validation; uses raw float64 values. Result is always Cool (unitless).
func (o *Operations) HyperLog2(stack *Stack) error {
	values, err := popAll(stack, "[lg]")
	if err != nil {
		return err
	}

	var result float64 = 0
	for i := 0; i < len(values); i++ {
		val, err := values[i].Float64()
		if err != nil {
			return buildError("[lg]", fmt.Errorf("operand %d: %w", i, err))
		}
		if val <= 0 {
			return buildError("[lg]", fmt.Errorf("log2 undefined for non-positive numbers"))
		}
		result += math.Log2(val)
	}

	cool := coolMetric(o.metricRegistry)
	stack.Push(NewNumber(result, o.GetMode(), cool))
	return nil
}

// HyperLog10 pops all values from stack, computes sum of log10 for all values, and pushes result.
// No metric validation; uses raw float64 values. Result is always Cool (unitless).
func (o *Operations) HyperLog10(stack *Stack) error {
	values, err := popAll(stack, "[log]")
	if err != nil {
		return err
	}

	var result float64 = 0
	for i := 0; i < len(values); i++ {
		val, err := values[i].Float64()
		if err != nil {
			return buildError("[log]", fmt.Errorf("operand %d: %w", i, err))
		}
		if val <= 0 {
			return buildError("[log]", fmt.Errorf("log10 undefined for non-positive numbers"))
		}
		result += math.Log10(val)
	}

	cool := coolMetric(o.metricRegistry)
	stack.Push(NewNumber(result, o.GetMode(), cool))
	return nil
}

// HyperLn pops all values from stack, computes sum of natural log for all values, and pushes result.
// No metric validation; uses raw float64 values. Result is always Cool (unitless).
func (o *Operations) HyperLn(stack *Stack) error {
	values, err := popAll(stack, "[ln]")
	if err != nil {
		return err
	}

	var result float64 = 0
	for i := 0; i < len(values); i++ {
		val, err := values[i].Float64()
		if err != nil {
			return buildError("[ln]", fmt.Errorf("operand %d: %w", i, err))
		}
		if val <= 0 {
			return buildError("[ln]", fmt.Errorf("ln undefined for non-positive numbers"))
		}
		result += math.Log(val)
	}

	cool := coolMetric(o.metricRegistry)
	stack.Push(NewNumber(result, o.GetMode(), cool))
	return nil
}

// resultMetricForHyperAdd finds the appropriate result metric for add/subtract/modulo hyper operations.
// When Cool absorbs a non-Cool category, use the first non-Cool metric.
// When all are Cool, use Cool.
func resultMetricForHyperAdd(metrics []*Metric) *Metric {
	for _, m := range metrics {
		if m != nil && m.Category != Universal {
			return m
		}
	}
	// All Universal (Cool) or empty slice — default to Cool.
	// In practice, metrics is never empty (popAll enforces >= 2 operands),
	// but we handle it defensively.
	if len(metrics) > 0 && metrics[0] != nil {
		return metrics[0]
	}
	m, _ := GetMetricRegistry().Find("Cool")
	return m
}

// validateSameCategory checks that all metrics belong to the same category.
// Cool (Universal) absorbs — it is compatible with any single non-Universal category.
// Returns an error if metrics span multiple non-Universal categories.
func validateSameCategory(metrics []*Metric, opName string) error {
	var dominantCat Category = Universal
	for _, m := range metrics {
		if m == nil || m.Category == Universal {
			continue
		}
		if dominantCat == Universal {
			dominantCat = m.Category
		} else if m.Category != dominantCat {
			return fmt.Errorf("%s: incompatible metrics: mixed %s and %s categories",
				opName, dominantCat, m.Category)
		}
	}
	return nil
}
