// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

// comparison operators

// compareValues pops two values, checks metric compatibility,
// converts to base units, and compares using the given function.
func compareValues(o *Operations, stack *Stack, op string, cmp func(float64, float64) bool) error {
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
	if !categoriesCompatible(aM, bM) {
		return metricError(op, aM, bM)
	}

	pm := o.GetPrefixMode()
	resultMetric, err := compatibleMetric(o.metricRegistry, aM, bM)
	if err != nil {
		return buildError(op, err)
	}
	aBase, err := convertToBase(o.metricRegistry, a, pm, resultMetric)
	if err != nil {
		return buildError(op, err)
	}
	bBase, err := convertToBase(o.metricRegistry, b, pm, resultMetric)
	if err != nil {
		return buildError(op, err)
	}

	stack.Push(NewFloatFromBool(cmp(aBase, bBase)))
	return nil
}

// GT pops two values from stack, compares (a > b), and pushes a boolean result.
// Converts to base units for metric-aware comparison.
func (o *Operations) GT(stack *Stack) error {
	return compareValues(o, stack, "gt", func(a, b float64) bool { return a > b })
}

// LT pops two values from stack, compares (a < b), and pushes a boolean result.
// Converts to base units for metric-aware comparison.
func (o *Operations) LT(stack *Stack) error {
	return compareValues(o, stack, "lt", func(a, b float64) bool { return a < b })
}

// GTE pops two values from stack, compares (a >= b), and pushes a boolean result.
// Converts to base units for metric-aware comparison.
func (o *Operations) GTE(stack *Stack) error {
	return compareValues(o, stack, "gte", func(a, b float64) bool { return a >= b })
}

// LTE pops two values from stack, compares (a <= b), and pushes a boolean result.
// Converts to base units for metric-aware comparison.
func (o *Operations) LTE(stack *Stack) error {
	return compareValues(o, stack, "lte", func(a, b float64) bool { return a <= b })
}

// EQ pops two values from stack, compares (a == b), and pushes a boolean result.
// Converts to base units for metric-aware comparison.
func (o *Operations) EQ(stack *Stack) error {
	return compareValues(o, stack, "eq", func(a, b float64) bool { return a == b })
}

// NEQ pops two values from stack, compares (a != b), and pushes a boolean result.
// Converts to base units for metric-aware comparison.
func (o *Operations) NEQ(stack *Stack) error {
	return compareValues(o, stack, "neq", func(a, b float64) bool { return a != b })
}
