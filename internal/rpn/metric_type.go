// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import "fmt"

// Category identifies which domain a metric belongs to.
type Category int

const (
	// Universal is the default category — "Cool" unit.
	Universal Category = iota
	// DataRate covers bits/sec, bytes/sec, etc.
	DataRate
	// DataSize covers bits, bytes, KB, MB, GB, etc.
	DataSize
	// Time covers ms, s, min, hr, day.
	Time
	// Weight covers mg, g, kg, ton, lb, oz.
	Weight
	// Speed covers m/s, km/h, mph, knots.
	Speed
	// Distance covers m, km, mi, ft, in.
	Distance
	// Custom is for user-defined units.
	Custom
	// _sentinel marks the upper bound for bounded Category iteration.
	// Adding a new Category: insert before _sentinel, and it will be
	// automatically picked up by parseCategory().
	_sentinel
)

// invalidCategory is a sentinel value returned by parseCategory for unknown names.
// It must not equal any valid Category (which start at 0), so callers that
// ignore the bool return cannot accidentally treat it as Universal.
const invalidCategory Category = -1

// categoryNames maps each Category value to its human-readable name.
// Adding a new Category only requires appending here; _sentinel bounds the range.
var categoryNames = []string{
	"Universal",
	"DataRate",
	"DataSize",
	"Time",
	"Weight",
	"Speed",
	"Distance",
	"Custom",
}

// String returns the human-readable name of the category.
func (c Category) String() string {
	if c >= 0 && c < Category(len(categoryNames)) {
		return categoryNames[c]
	}
	return fmt.Sprintf("Category(%d)", c)
}

// PrefixMode determines whether data size prefixes are SI (1000-based) or IEC (1024-based).
type PrefixMode int

const (
	// SI uses powers of 1000 for prefixes: K=1000, M=1000000, etc.
	SI PrefixMode = iota
	// IEC uses powers of 1024 for prefixes: Ki=1024, Mi=1048576, etc.
	IEC
)

// String returns the human-readable name of the prefix mode.
func (p PrefixMode) String() string {
	switch p {
	case SI:
		return "SI"
	case IEC:
		return "IEC"
	default:
		return fmt.Sprintf("PrefixMode(%d)", p)
	}
}

// Metric defines a unit of measurement with conversion to its base unit.
type Metric struct {
	// Name is the canonical name, e.g. "Mbps", "GB", "hr".
	Name string
	// Category identifies the domain this metric belongs to.
	Category Category
	// BaseUnit is the name of the base unit in this category.
	BaseUnit string
	// Factor returns the conversion factor to the base unit for the given prefix mode.
	Factor func(mode PrefixMode) float64
	// IsRate is true for rate units (per-second).
	IsRate bool
	// IsCustom is true for user-defined custom units.
	IsCustom bool
}

// String returns the canonical name of the metric.
func (m *Metric) String() string {
	return m.Name
}
