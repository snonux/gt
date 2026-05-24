// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import "fmt"

// resolveMetric returns the metric for a StackValue, defaulting to Cool if nil.
func resolveMetric(reg MetricReader, n StackValue) (*Metric, error) {
	m := n.Metric()
	if m == nil {
		var ok bool
		m, ok = reg.Find("Cool")
		if !ok {
			return nil, fmt.Errorf("metric registry missing Cool metric")
		}
	}
	return m, nil
}

// validateCategories checks that all metrics belong to the same category.
// Cool (Universal) absorbs — it is compatible with any single non-Universal category.
// Returns nil if compatible, error if metrics span multiple non-Universal categories.
// Works for both binary (2 metrics) and N-ary (slice) cases.
func validateCategories(metrics []*Metric, opName string) error {
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

// resultMetricForAdd finds the appropriate result metric for add/subtract/modulo operations.
// When Cool absorbs a non-Cool category, use the first non-Cool metric.
// When all are Cool, use Cool.
func resultMetricForAdd(metrics []*Metric) *Metric {
	for _, m := range metrics {
		if m != nil && m.Category != Universal {
			return m
		}
	}
	if len(metrics) > 0 && metrics[0] != nil {
		return metrics[0]
	}
	return GetCoolMetric()
}

// categoriesCompatible checks if two metrics are compatible for arithmetic.
func categoriesCompatible(a, b *Metric) bool {
	return validateCategories([]*Metric{a, b}, "") == nil
}

// compatibleMetric returns the resulting metric for + and - operations.
func compatibleMetric(reg MetricReader, a, b *Metric) *Metric {
	return resultMetricForAdd([]*Metric{a, b})
}

// convertToBase converts a StackValue's value to its metric's base unit.
//
// Cool absorbing (important): When the value's metric is Cool (Universal)
// and resultMetric is non-Cool, the Cool value is treated as units of the
// result metric. For example, '5 100Mbps +' converts the Cool 5 as 5Mbps,
// producing 105Mbps. This allows seamless mixing of unitless scalars with
// metric values. The Cool value is multiplied by resultMetric.Factor(),
// NOT by 1 (base units).
//
// Returns the converted float64 value in base units.
func convertToBase(reg MetricReader, n StackValue, mode PrefixMode, resultMetric *Metric) (float64, error) {
	nv, ok := n.(NumericValue)
	if !ok {
		return 0, fmt.Errorf("convertToBase: value %q is not numeric", n)
	}
	m, err := resolveMetric(reg, n)
	if err != nil {
		return 0, err
	}
	val, err := nv.Float64()
	if err != nil {
		return 0, fmt.Errorf("convertToBase: %w", err)
	}
	// Cool absorbing: if operand is Cool but result metric is not,
	// treat the Cool value as units of the result metric.
	if m.Category == Universal && resultMetric != nil && resultMetric.Category != Universal {
		return val * resultMetric.Factor(mode), nil
	}
	return val * m.Factor(mode), nil
}

// convertFromBase converts a base-unit value back to the given metric.
func convertFromBase(reg MetricReader, baseVal float64, m *Metric, mode PrefixMode) (float64, error) {
	if m == nil {
		var err error
		m, err = baseMetric(reg, "Cool")
		if err != nil {
			return 0, err
		}
	}
	return baseVal / m.Factor(mode), nil
}

// multiplicationInference maps ordered category pairs to the result base-unit name.
// Both orderings are stored so that map lookup replaces the switch for commutative mul.
var multiplicationInference = map[[2]Category]string{
	{DataRate, Time}: "bits",
	{Time, DataRate}: "bits",
	{Speed, Time}:    "m",
	{Time, Speed}:    "m",
}

// resultMetricForMul computes the resulting metric for multiplication.
func resultMetricForMul(reg MetricReader, a, b *Metric) (*Metric, error) {
	if a == nil || a.Category == Universal {
		if b == nil {
			return baseMetric(reg, "Cool")
		}
		return b, nil
	}
	if b == nil || b.Category == Universal {
		return a, nil
	}

	if name, ok := multiplicationInference[[2]Category{a.Category, b.Category}]; ok {
		return baseMetric(reg, name)
	}
	return baseMetric(reg, "Cool")
}

// divisionInference maps (dividend, divisor) category pairs to the result base-unit name.
// Unlike multiplication, division is not commutative so order matters.
var divisionInference = map[[2]Category]string{
	{DataSize, Time}: "bps",
	{Distance, Time}: "mps",
}

// resultMetricForDiv computes the resulting metric for division.
func resultMetricForDiv(reg MetricReader, a, b *Metric) (*Metric, error) {
	if a == nil && b == nil {
		return baseMetric(reg, "Cool")
	}
	if b == nil || b.Category == Universal {
		if a == nil {
			return baseMetric(reg, "Cool")
		}
		return a, nil
	}
	// When dividend is Cool (unitless) and divisor has a metric,
	// result should be Cool. E.g., 5 / 10km → 0.5 (Cool, not km).
	// Cool-absorbing is designed for addition, not division.
	if a == nil || a.Category == Universal {
		return baseMetric(reg, "Cool")
	}

	if name, ok := divisionInference[[2]Category{a.Category, b.Category}]; ok {
		return baseMetric(reg, name)
	}
	return baseMetric(reg, "Cool")
}

// metricError returns a descriptive error for incompatible metric operations.
func metricError(op string, a, b *Metric) error {
	aName, aCat := "Cool", "Universal"
	if a != nil {
		aName, aCat = a.Name, a.Category.String()
	}
	bName, bCat := "Cool", "Universal"
	if b != nil {
		bName, bCat = b.Name, b.Category.String()
	}
	return fmt.Errorf("%s: incompatible metrics %s (%s) and %s (%s)",
		op, aName, aCat, bName, bCat)
}

// coolMetric returns the Cool metric from the registry.
func coolMetric(reg MetricReader) (*Metric, error) {
	m, ok := reg.Find("Cool")
	if !ok {
		return nil, fmt.Errorf("metric registry missing Cool metric")
	}
	return m, nil
}

// baseMetric looks up a base metric from the registry.
func baseMetric(reg MetricReader, name string) (*Metric, error) {
	m, ok := reg.Find(name)
	if !ok {
		return nil, fmt.Errorf("metric registry missing base unit %q", name)
	}
	return m, nil
}

// Convert converts a value from its current metric to a target metric.
// Pops target metric (from @X syntax), then pops value to convert.
// Validates category compatibility, converts through base unit,
// and pushes the result with the target metric.
func (o *Operations) Convert(stack *Stack) error {
	// 1. Pop target (from @X syntax)
	target, err := popStack(stack, "convert")
	if err != nil {
		return err
	}
	// 2. Pop value to convert
	value, err := popStack(stack, "convert")
	if err != nil {
		return err
	}
	// 3. Get metrics
	targetMetric := target.Metric()
	valueMetric, err := resolveMetric(o.metricRegistry, value)
	if err != nil {
		return buildError("convert", err)
	}
	// 4. Validate same category (or Cool absorbing)
	if !categoriesCompatible(valueMetric, targetMetric) {
		return metricError("convert", valueMetric, targetMetric)
	}
	// 5. Convert through base unit: value → base → target
	pm := o.GetPrefixMode()
	baseVal, err := convertToBase(o.metricRegistry, value, pm, targetMetric)
	if err != nil {
		return buildError("convert", err)
	}
	resultVal, err := convertFromBase(o.metricRegistry, baseVal, targetMetric, pm)
	if err != nil {
		return buildError("convert", err)
	}
	// 6. Push result with target metric
	stack.Push(NewNumberWithMetric(resultVal, o.GetMode(), targetMetric))
	return nil
}
