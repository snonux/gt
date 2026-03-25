// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

// CalculationMode represents the mode for number calculations.
type CalculationMode int

const (
	// FloatMode uses float64 for calculations (default).
	FloatMode CalculationMode = iota
	// RationalMode uses *big.Rat for precise rational calculations.
	RationalMode
)
