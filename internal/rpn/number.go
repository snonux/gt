// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"math/big"
)

// Number represents a number that can be used in RPN calculations.
// It can be either a float64 or a *big.Rat for precise rational calculations.
// Booleans are also supported through IsBool() and Bool() methods.
// Arithmetic is handled by Operations (metric-aware); this interface provides
// value inspection and comparison only.
type Number interface {
	// String returns the string representation of the number.
	String() string
	// Float64 returns the float64 representation.
	// Returns error if the number is not representable (e.g., StringNum, Symbol).
	Float64() (float64, error)
	// IsZero returns true if the number is zero.
	IsZero() bool
	// IsNegative returns true if the number is negative.
	IsNegative() bool
	// Compare returns -1, 0, or 1 if this number is less than, equal to, or greater than another.
	// Returns error if the operation is not supported (e.g., StringNum, Symbol).
	Compare(other Number) (int, error)
	// IsBool returns true if this number represents a boolean value.
	IsBool() bool
	// Bool returns the boolean value.
	// Returns error if the number is not a boolean.
	Bool() (bool, error)
	// IsString returns true if this number represents a string value.
	IsString() bool
	// IsSymbol returns true if this number represents a symbol.
	IsSymbol() bool
	// Metric returns the metric unit for this number. Returns the Cool default if none set.
	Metric() *Metric
	// SetMetric returns a copy of this Number with the given metric attached.
	SetMetric(m *Metric) Number
}

// NewNumber creates a Number from a float64 value with the given metric.
// The actual type depends on the current calculation mode.
// If metric is nil or omitted, defaults to Cool.
func NewNumber(value float64, mode CalculationMode, metric ...*Metric) Number {
	m := GetCoolMetric()
	if len(metric) > 0 && metric[0] != nil {
		m = metric[0]
	}
	if mode == RationalMode {
		return NewRatWithMetric(value, m)
	}
	return NewFloatWithMetric(value, m)
}

// GetCoolMetric returns the universal default metric.
func GetCoolMetric() *Metric {
	m, _ := GetMetricRegistry().Find("Cool")
	return m
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
func (f *Float) SetMetric(m *Metric) Number {
	return &Float{n: f.n, isBool: f.isBool, boolVal: f.boolVal, metric: m}
}

// Compare returns -1, 0, or 1 if this float is less than, equal to, or greater than another.
func (f *Float) Compare(other Number) (int, error) {
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
	return &Rat{n: rat}, nil
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
func (r *Rat) SetMetric(m *Metric) Number {
	n := &big.Rat{}
	n.Set(r.n)
	return &Rat{n: n, isBool: r.isBool, boolVal: r.boolVal, metric: m}
}

// Compare returns -1, 0, or 1 if this rational is less than, equal to, or greater than another.
func (r *Rat) Compare(other Number) (int, error) {
	otherRat, ok := other.(*Rat)
	if !ok {
		return 0, fmt.Errorf("cannot compare: operand is not a rational number")
	}
	return r.n.Cmp(otherRat.n), nil
}

// ToRat converts a Number to *big.Rat.
// Returns nil if the number is not a Rat.
func ToRat(n Number) *big.Rat {
	if r, ok := n.(*Rat); ok {
		return r.n
	}
	return nil
}

// ToFloat converts a Number to float64.
func ToFloat(n Number) float64 {
	val, _ := n.Float64()
	return val
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

// Float64 returns error for string numbers (not numeric).
func (s *StringNum) Float64() (float64, error) {
	return 0, fmt.Errorf("string not supported for Float64()")
}

// IsString returns true for StringNum.
func (s *StringNum) IsString() bool {
	return true
}

func (s *StringNum) IsZero() bool                     { return false }
func (s *StringNum) IsNegative() bool                 { return false }
func (s *StringNum) Compare(other Number) (int, error) { return 0, fmt.Errorf("string not supported for comparison") }
func (s *StringNum) Bool() (bool, error)              { return false, fmt.Errorf("string not supported for Bool()") }
func (s *StringNum) IsBool() bool                     { return false }
func (s *StringNum) IsSymbol() bool                   { return false }
func (s *StringNum) Metric() *Metric                  { return GetCoolMetric() }
func (s *StringNum) SetMetric(m *Metric) Number       { return s }

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

// Float64 returns error for symbols (not numeric).
func (s *Symbol) Float64() (float64, error) {
	return 0, fmt.Errorf("symbol not supported for Float64()")
}

// Name returns the symbol name.
func (s *Symbol) Name() string {
	return s.name
}

// IsSymbol returns true for Symbol.
func (s *Symbol) IsSymbol() bool {
	return true
}

func (s *Symbol) IsZero() bool {
	return false
}
func (s *Symbol) IsNegative() bool {
	return false
}
func (s *Symbol) IsString() bool {
	return false
}
func (s *Symbol) Compare(other Number) (int, error) {
	return 0, fmt.Errorf("symbol not supported for comparison")
}
func (s *Symbol) Bool() (bool, error) {
	return false, fmt.Errorf("symbol not supported for Bool()")
}
func (s *Symbol) IsBool() bool {
	return false
}
func (s *Symbol) Metric() *Metric {
	return GetCoolMetric()
}
func (s *Symbol) SetMetric(m *Metric) Number {
	return s
}
