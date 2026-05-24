// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"math/big"
)

// StackValue represents a value on the RPN stack.
// All stack items implement this interface, including numeric values,
// strings, and symbols.
type StackValue interface {
	// String returns the string representation of the value.
	String() string
	// IsBool returns true if this value represents a boolean value.
	IsBool() bool
	// IsString returns true if this value represents a string value.
	IsString() bool
	// IsSymbol returns true if this value represents a symbol.
	IsSymbol() bool
	// Metric returns the metric unit for this value.
	Metric() *Metric
}

// NumericValue represents a numeric value that supports arithmetic operations.
// Float and Rat implement this; StringNum and Symbol do not.
type NumericValue interface {
	StackValue
	// Float64 returns the float64 representation.
	Float64() (float64, error)
	// IsZero returns true if the number is zero.
	IsZero() bool
	// IsNegative returns true if the number is negative.
	IsNegative() bool
	// Compare returns -1, 0, or 1 if this number is less than, equal to, or
	// greater than another.
	Compare(other NumericValue) (int, error)
	// Bool returns the boolean value.
	// Returns error if the number is not a boolean.
	Bool() (bool, error)
	// SetMetric returns a copy of this NumericValue with the given metric attached.
	SetMetric(m *Metric) NumericValue
}

// NumericValue is the interface for numeric values (float64 or *big.Rat).
// It extends StackValue with metric operations.
//
// Note: type Number was a legacy alias for NumericValue, removed to avoid
// confusion between the alias and the concrete type.

// Compile-time interface satisfaction checks.
var _ StackValue = (*Float)(nil)
var _ NumericValue = (*Float)(nil)
var _ StackValue = (*Rat)(nil)
var _ NumericValue = (*Rat)(nil)
var _ StackValue = (*StringNum)(nil)
var _ StackValue = (*Symbol)(nil)

// NewNumber creates a NumericValue from a float64 value with the given mode.
// The actual type depends on the current calculation mode (Float or Rat).
// The metric defaults to Cool (unitless).
func NewNumber(value float64, mode CalculationMode) NumericValue {
	if mode == RationalMode {
		return NewRat(value)
	}
	return NewFloat(value)
}

// NewNumberWithMetric creates a NumericValue from a float64 value with an explicit metric.
// The actual type depends on the current calculation mode (Float or Rat).
// If metric is nil, defaults to Cool.
func NewNumberWithMetric(value float64, mode CalculationMode, metric *Metric) NumericValue {
	if metric == nil {
		metric = GetCoolMetric()
	}
	if mode == RationalMode {
		return NewRatWithMetric(value, metric)
	}
	return NewFloatWithMetric(value, metric)
}

// GetCoolMetric returns the universal default metric.
// Ensures the metric registry is initialized first, then returns the
// cached Cool metric pointer. Safe for concurrent use.
func GetCoolMetric() *Metric {
	// Ensure registry (and cachedCoolMetric) are initialized.
	// registryOnce.Do() in GetMetricRegistry() sets cachedCoolMetric.
	GetMetricRegistry()
	return cachedCoolMetric
}

// Float is a Number implementation using float64.
// It can also represent boolean values (true=1, false=0).
type Float struct {
	n       float64
	isBool  bool
	boolVal bool
	metric  *Metric
}

// NewFloat creates a new Float number with the Cool metric.
func NewFloat(n float64) *Float {
	return &Float{n: n, isBool: false, boolVal: false, metric: GetCoolMetric()}
}

// NewFloatWithMetric creates a new Float number with the given metric.
func NewFloatWithMetric(n float64, metric *Metric) *Float {
	return &Float{n: n, isBool: false, boolVal: false, metric: metric}
}

// NewFloatFromBool creates a new Float representing a boolean.
func NewFloatFromBool(b bool) *Float {
	return &Float{n: 0, isBool: true, boolVal: b, metric: GetCoolMetric()}
}

// String returns the string representation of the float.
func (f *Float) String() string {
	if f.isBool {
		if f.boolVal {
			return "true"
		}
		return "false"
	}
	return fmt.Sprintf("%.10g", f.n)
}

// Float64 returns the float64 value.
func (f *Float) Float64() (float64, error) {
	if f.isBool {
		if f.boolVal {
			return 1, nil
		}
		return 0, nil
	}
	return f.n, nil
}

// IsBool returns true if this number represents a boolean value.
func (f *Float) IsBool() bool {
	return f.isBool
}

// Bool returns the boolean value.
// Returns error if the number is not a boolean.
func (f *Float) Bool() (bool, error) {
	if !f.isBool {
		return false, fmt.Errorf("not a boolean")
	}
	return f.boolVal, nil
}

// IsZero returns true if the float is zero.
// For boolean values, false (0) is zero, true (1) is not zero.
func (f *Float) IsZero() bool {
	val, _ := f.Float64()
	return val == 0
}

// IsNegative returns true if the float is negative.
func (f *Float) IsNegative() bool {
	return f.n < 0
}

// IsString returns true if this number represents a string value.
func (f *Float) IsString() bool {
	return false
}

// IsSymbol returns true if this number represents a symbol.
func (f *Float) IsSymbol() bool {
	return false
}

// Metric returns the metric for this number.
func (f *Float) Metric() *Metric {
	return f.metric
}

// SetMetric returns a copy of this Float with the given metric.
func (f *Float) SetMetric(m *Metric) NumericValue {
	return &Float{n: f.n, isBool: f.isBool, boolVal: f.boolVal, metric: m}
}

// Compare returns -1, 0, or 1 if this float is less than, equal to, or greater than another.
func (f *Float) Compare(other NumericValue) (int, error) {
	otherF, err := other.Float64()
	if err != nil {
		return 0, fmt.Errorf("cannot compare: %w", err)
	}
	if f.n < otherF {
		return -1, nil
	}
	if f.n > otherF {
		return 1, nil
	}
	return 0, nil
}

// Rat is a Number implementation using *big.Rat.
// It can also represent boolean values (true=1, false=0).
type Rat struct {
	n       *big.Rat
	isBool  bool
	boolVal bool
	metric  *Metric
}

// NewRat creates a new Rat number from a float64 with the Cool metric.
func NewRat(n float64) *Rat {
	r := &big.Rat{}
	r.SetFloat64(n)
	return &Rat{n: r, isBool: false, boolVal: false, metric: GetCoolMetric()}
}

// NewRatWithMetric creates a new Rat number from a float64 with the given metric.
func NewRatWithMetric(n float64, metric *Metric) *Rat {
	r := &big.Rat{}
	r.SetFloat64(n)
	return &Rat{n: r, isBool: false, boolVal: false, metric: metric}
}

// NewRatFromBool creates a new Rat representing a boolean.
func NewRatFromBool(b bool) *Rat {
	r := &big.Rat{}
	if b {
		r.SetInt64(1)
	} else {
		r.SetInt64(0)
	}
	return &Rat{n: r, isBool: true, boolVal: b, metric: GetCoolMetric()}
}

// NewRatFromString creates a new Rat number from a string representation.
func NewRatFromString(s string) (*Rat, error) {
	r := &big.Rat{}
	rat, ok := r.SetString(s)
	if !ok || rat == nil {
		return nil, fmt.Errorf("invalid rational number: %s", s)
	}
	return &Rat{n: rat, metric: GetCoolMetric()}, nil
}

// String returns the string representation of the rational number.
func (r *Rat) String() string {
	if r.isBool {
		if r.boolVal {
			return "true"
		}
		return "false"
	}
	// Format as decimal for consistency with Float
	// Use a reasonable precision
	return r.n.FloatString(10)
}

// Float64 returns the float64 representation.
func (r *Rat) Float64() (float64, error) {
	if r.isBool {
		if r.boolVal {
			return 1, nil
		}
		return 0, nil
	}
	f, ok := r.n.Float64()
	if !ok {
		return 0, fmt.Errorf("cannot convert rational number to float64")
	}
	return f, nil
}

// IsBool returns true if this number represents a boolean value.
func (r *Rat) IsBool() bool {
	return r.isBool
}

// Bool returns the boolean value.
// Returns error if the number is not a boolean.
func (r *Rat) Bool() (bool, error) {
	if !r.isBool {
		return false, fmt.Errorf("not a boolean")
	}
	return r.boolVal, nil
}

// IsZero returns true if the rational number is zero.
func (r *Rat) IsZero() bool {
	return r.n.Sign() == 0
}

// IsNegative returns true if the rational number is negative.
func (r *Rat) IsNegative() bool {
	return r.n.Sign() < 0
}

// IsString returns true if this number represents a string value.
func (r *Rat) IsString() bool {
	return false
}

// IsSymbol returns true if this number represents a symbol.
func (r *Rat) IsSymbol() bool {
	return false
}

// Metric returns the metric for this number.
func (r *Rat) Metric() *Metric {
	return r.metric
}

// SetMetric returns a copy of this Rat with the given metric.
func (r *Rat) SetMetric(m *Metric) NumericValue {
	n := &big.Rat{}
	n.Set(r.n)
	return &Rat{n: n, isBool: r.isBool, boolVal: r.boolVal, metric: m}
}

// Compare returns -1, 0, or 1 if this rational is less than, equal to, or greater than another.
// For Rat-to-Rat comparison, uses direct big.Rat.Cmp() for exact precision.
// For comparison with non-Rat types, converts via float64 (acceptable precision loss).
func (r *Rat) Compare(other NumericValue) (int, error) {
	// Rat-to-Rat: direct comparison with full precision
	if otherRat, ok := other.(*Rat); ok {
		return r.n.Cmp(otherRat.n), nil
	}
	// Rat-to-Float or other: convert via float64 (acceptable precision loss)
	otherF, err := other.Float64()
	if err != nil {
		return 0, fmt.Errorf("cannot compare: %w", err)
	}
	otherRat := &big.Rat{}
	otherRat.SetFloat64(otherF)
	return r.n.Cmp(otherRat), nil
}

// ToRat converts a StackValue to *big.Rat.
// Returns nil if the value is not a Rat.
func ToRat(n StackValue) *big.Rat {
	if r, ok := n.(*Rat); ok {
		return r.n
	}
	return nil
}

// ToFloat converts a StackValue to float64.
// Returns 0 if the value is not numeric.
func ToFloat(n StackValue) float64 {
	if nv, ok := n.(NumericValue); ok {
		val, _ := nv.Float64()
		return val
	}
	return 0
}

// StringNum represents a string value on the stack for variable names.
type StringNum struct {
	value string
}

// NewStringNum creates a new StringNum from a string.
func NewStringNum(s string) *StringNum {
	return &StringNum{value: s}
}

// String returns the string representation.
func (s *StringNum) String() string {
	return s.value
}

// IsString returns true for StringNum.
func (s *StringNum) IsString() bool {
	return true
}

func (s *StringNum) IsBool() bool    { return false }
func (s *StringNum) IsSymbol() bool  { return false }
func (s *StringNum) Metric() *Metric { return GetCoolMetric() }

// Symbol represents a variable symbol on the stack.
// Symbols are created when:
// - The user enters :x syntax (explicit symbol)
// - A bare identifier x is used but the variable is unbound
// When printed, symbols are prefixed with : (e.g., :x) to distinguish them from values.
type Symbol struct {
	name string
}

// NewSymbol creates a new Symbol from a name.
func NewSymbol(name string) *Symbol {
	return &Symbol{name: name}
}

// String returns the string representation of the symbol, prefixed with :.
func (s *Symbol) String() string {
	return ":" + s.name
}

// Name returns the symbol name.
func (s *Symbol) Name() string {
	return s.name
}

// IsSymbol returns true for Symbol.
func (s *Symbol) IsSymbol() bool {
	return true
}

func (s *Symbol) IsBool() bool    { return false }
func (s *Symbol) IsString() bool  { return false }
func (s *Symbol) Metric() *Metric { return GetCoolMetric() }


