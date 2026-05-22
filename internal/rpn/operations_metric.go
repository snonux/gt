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
	if a.Category == Universal {
		return b
	}
	if b.Category == Universal {
		return a
	}

	// Cross-category inference
	switch {
	case a.Category == DataRate && b.Category == Time:
		// Rate × Time = Size — use base unit of DataSize
		return dataSizeBaseMetric()
	case a.Category == Time && b.Category == DataRate:
		return dataSizeBaseMetric()
	case a.Category == Speed && b.Category == Time:
		// Speed × Time = Distance — use base unit of Distance
		return distanceBaseMetric()
	case a.Category == Time && b.Category == Speed:
		return distanceBaseMetric()
	default:
		// Same category or unknown combination → Cool (unitless)
		if a.Category == b.Category {
			return GetCoolMetric()
		}
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
	if a.Category == Universal && b.Category == Universal {
		return a
	}
	if b.Category == Universal {
		return a
	}
	if a.Category == Universal {
		return b
	}

	// Cross-category inference
	switch {
	case a.Category == DataSize && b.Category == Time:
		return dataRateBaseMetric()
	case a.Category == Distance && b.Category == Time:
		return speedBaseMetric()
	case a.Category == DataRate && b.Category == Time:
		// Rate / Time = Size per time² → Cool
		return GetCoolMetric()
	default:
		// Same category ratio → Cool
		if a.Category == b.Category {
			return GetCoolMetric()
		}
		return GetCoolMetric()
	}
}

// metricError returns a descriptive error for incompatible metric operations.
func metricError(op string, a, b *Metric) error {
	return fmt.Errorf("%s: incompatible metrics %s (%s) and %s (%s)",
		op, a.Name, a.Category, b.Name, b.Category)
}

// dataSizeBaseMetric returns the DataSize base unit (bits).
func dataSizeBaseMetric() *Metric {
	m, ok := GetMetricRegistry().Find("bits")
	if !ok {
		panic("metric registry missing base unit 'bits'")
	}
	return m
}

// dataRateBaseMetric returns the DataRate base unit (bps).
func dataRateBaseMetric() *Metric {
	m, ok := GetMetricRegistry().Find("bps")
	if !ok {
		panic("metric registry missing base unit 'bps'")
	}
	return m
}

// distanceBaseMetric returns the Distance base unit (meters).
func distanceBaseMetric() *Metric {
	m, ok := GetMetricRegistry().Find("m")
	if !ok {
		panic("metric registry missing base unit 'm'")
	}
	return m
}

// speedBaseMetric returns the Speed base unit (mps).
func speedBaseMetric() *Metric {
	m, ok := GetMetricRegistry().Find("mps")
	if !ok {
		panic("metric registry missing base unit 'mps'")
	}
	return m
}
