// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import "fmt"

// resolveMetric returns the metric for a Number, defaulting to Cool if nil.
func resolveMetric(n Number) *Metric {
	m := n.Metric()
	if m == nil {
		return GetCoolMetric()
	}
	return m
}

// categoriesCompatible checks if two metrics are compatible for arithmetic.
// Cool (Universal) is compatible with anything. Same category is compatible.
func categoriesCompatible(a, b *Metric) bool {
	if a == nil || b == nil {
		return true
	}
	if a.Category == Universal || b.Category == Universal {
		return true
	}
	return a.Category == b.Category
}

// compatibleMetric returns the resulting metric for + and - operations.
// Cool absorbs: if either is Cool, result is the other's metric (or Cool if both).
// Same category: result uses left operand's metric.
func compatibleMetric(a, b *Metric) *Metric {
	if a == nil {
		a = GetCoolMetric()
	}
	if b == nil {
		b = GetCoolMetric()
	}
	if a.Category == Universal && b.Category == Universal {
		return a // Cool
	}
	if a.Category == Universal {
		return b
	}
	if b.Category == Universal {
		return a
	}
	// Same category: use left operand's metric
	return a
}

// convertToBase converts a Number's value to its metric's base unit.
// Returns the converted float64 value.
func convertToBase(n Number, mode PrefixMode) (float64, error) {
	m := resolveMetric(n)
	val, err := n.Float64()
	if err != nil {
		return 0, fmt.Errorf("convertToBase: %w", err)
	}
	return val * m.Factor(mode), nil
}

// convertFromBase converts a base-unit value back to the given metric.
func convertFromBase(baseVal float64, m *Metric, mode PrefixMode) float64 {
	if m == nil {
		m = GetCoolMetric()
	}
	return baseVal / m.Factor(mode)
}

// resultMetricForMul computes the resulting metric for multiplication.
// Cross-category inference rules:
//   - DataRate × Time → DataSize
//   - Time × DataRate → DataSize
//   - Speed × Time → Distance
//   - Time × Speed → Distance
//   - Universal × X → X (Cool absorbs)
//   - Otherwise → Cool (result is unitless product)
func resultMetricForMul(a, b *Metric) *Metric {
	if a == nil || a.Category == Universal {
		if b == nil {
			return GetCoolMetric()
		}
		return b
	}
	if b == nil || b.Category == Universal {
		return a
	}

	// Cross-category inference
	switch {
	case a.Category == DataRate && b.Category == Time:
		return findBaseMetric("bits")
	case a.Category == Time && b.Category == DataRate:
		return findBaseMetric("bits")
	case a.Category == Speed && b.Category == Time:
		return findBaseMetric("m")
	case a.Category == Time && b.Category == Speed:
		return findBaseMetric("m")
	default:
		return GetCoolMetric()
	}
}

// resultMetricForDiv computes the resulting metric for division.
// Cross-category inference rules:
//   - DataSize / Time → DataRate (base unit)
//   - Distance / Time → Speed (base unit)
//   - DataRate / DataRate → Cool (ratio)
//   - Speed / Speed → Cool (ratio)
//   - Universal / X → X
//   - X / Universal → X
//   - Otherwise → Cool
func resultMetricForDiv(a, b *Metric) *Metric {
	if a == nil && b == nil {
		return GetCoolMetric()
	}
	if b == nil || b.Category == Universal {
		if a == nil {
			return GetCoolMetric()
		}
		return a
	}
	if a == nil || a.Category == Universal {
		return b
	}

	// Cross-category inference
	switch {
	case a.Category == DataSize && b.Category == Time:
		return findBaseMetric("bps")
	case a.Category == Distance && b.Category == Time:
		return findBaseMetric("mps")
	default:
		return GetCoolMetric()
	}
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

// findBaseMetric looks up a base metric from the global registry.
// Panics if not found (indicates misconfiguration).
func findBaseMetric(name string) *Metric {
	m, ok := GetMetricRegistry().Find(name)
	if !ok {
		panic(fmt.Sprintf("metric registry missing base unit %q", name))
	}
	return m
}
