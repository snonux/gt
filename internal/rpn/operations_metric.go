// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import "fmt"

// resolveMetric returns the metric for a Number, defaulting to Cool if nil.
func resolveMetric(reg *MetricRegistry, n Number) *Metric {
	m := n.Metric()
	if m == nil {
		return coolMetric(reg)
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
func compatibleMetric(reg *MetricRegistry, a, b *Metric) *Metric {
	if a == nil {
		a = coolMetric(reg)
	}
	if b == nil {
		b = coolMetric(reg)
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
	m := resolveMetric(GetMetricRegistry(), n)
	val, err := n.Float64()
	if err != nil {
		return 0, fmt.Errorf("convertToBase: %w", err)
	}
	return val * m.Factor(mode), nil
}

// convertFromBase converts a base-unit value back to the given metric.
func convertFromBase(baseVal float64, m *Metric, mode PrefixMode) float64 {
	if m == nil {
		m = coolMetric(GetMetricRegistry())
	}
	return baseVal / m.Factor(mode)
}

// resultMetricForMul computes the resulting metric for multiplication.
func resultMetricForMul(reg *MetricRegistry, a, b *Metric) *Metric {
	if a == nil || a.Category == Universal {
		if b == nil {
			return coolMetric(reg)
		}
		return b
	}
	if b == nil || b.Category == Universal {
		return a
	}

	// Cross-category inference
	switch {
	case a.Category == DataRate && b.Category == Time:
		return baseMetric(reg, "bits")
	case a.Category == Time && b.Category == DataRate:
		return baseMetric(reg, "bits")
	case a.Category == Speed && b.Category == Time:
		return baseMetric(reg, "m")
	case a.Category == Time && b.Category == Speed:
		return baseMetric(reg, "m")
	default:
		return coolMetric(reg)
	}
}

// resultMetricForDiv computes the resulting metric for division.
func resultMetricForDiv(reg *MetricRegistry, a, b *Metric) *Metric {
	if a == nil && b == nil {
		return coolMetric(reg)
	}
	if b == nil || b.Category == Universal {
		if a == nil {
			return coolMetric(reg)
		}
		return a
	}
	if a == nil || a.Category == Universal {
		return b
	}

	// Cross-category inference
	switch {
	case a.Category == DataSize && b.Category == Time:
		return baseMetric(reg, "bps")
	case a.Category == Distance && b.Category == Time:
		return baseMetric(reg, "mps")
	default:
		return coolMetric(reg)
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

// coolMetric returns the Cool metric from the registry.
func coolMetric(reg *MetricRegistry) *Metric {
	m, ok := reg.Find("Cool")
	if !ok {
		panic("metric registry missing Cool metric")
	}
	return m
}

// baseMetric looks up a base metric from the registry.
func baseMetric(reg *MetricRegistry, name string) *Metric {
	m, ok := reg.Find(name)
	if !ok {
		panic(fmt.Sprintf("metric registry missing base unit %q", name))
	}
	return m
}
