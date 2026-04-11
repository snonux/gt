// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"math"
	"math/big"
)

// Number represents a number that can be used in RPN calculations.
// It can be either a float64 or a *big.Rat for precise rational calculations.
// Booleans are also supported through IsBool() and Bool() methods.
type Number interface {
	// String returns the string representation of the number.
	String() string
	// Float64 returns the float64 representation.
	// Returns error if the number is not representable (e.g., StringNum, Symbol).
	Float64() (float64, error)
	// Add returns the sum of this number and another.
	// Returns error if the operation is not supported (e.g., StringNum, Symbol).
	Add(other Number) (Number, error)
	// Sub returns the difference of this number and another.
	// Returns error if the operation is not supported (e.g., StringNum, Symbol).
	Sub(other Number) (Number, error)
	// Mul returns the product of this number and another.
	// Returns error if the operation is not supported (e.g., StringNum, Symbol).
	Mul(other Number) (Number, error)
	// Div returns the quotient of this number and another.
	// Returns (nil, error) if division by zero or operation not supported.
	Div(other Number) (Number, error)
	// Pow returns this number raised to the power of another.
	// Returns error if the operation is not supported (e.g., StringNum, Symbol).
	Pow(other Number) (Number, error)
	// Mod returns the remainder of this number divided by another.
	// Returns (nil, error) if modulo by zero or operation not supported.
	Mod(other Number) (Number, error)
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
}

// NewNumber creates a Number from a float64 value.
// The actual type depends on the current calculation mode.
func NewNumber(value float64, mode CalculationMode) Number {
	if mode == RationalMode {
		return NewRat(value)
	}
	return NewFloat(value)
}

// Float is a Number implementation using float64.
// It can also represent boolean values (true=1, false=0).
type Float struct {
	n       float64
	isBool  bool
	boolVal bool
}

// NewFloat creates a new Float number.
func NewFloat(n float64) *Float {
	return &Float{n: n, isBool: false, boolVal: false}
}

// NewFloatFromBool creates a new Float representing a boolean.
func NewFloatFromBool(b bool) *Float {
	return &Float{n: 0, isBool: true, boolVal: b}
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

// Add returns the sum of two float numbers.
func (f *Float) Add(other Number) (Number, error) {
	otherF, err := other.Float64()
	if err != nil {
		return nil, fmt.Errorf("cannot add: %w", err)
	}
	fF, err := f.Float64()
	if err != nil {
		return nil, fmt.Errorf("cannot add: %w", err)
	}
	// Use Float64() to handle both regular numbers and boolean values
	return NewFloat(fF + otherF), nil
}

// Sub returns the difference of two float numbers.
func (f *Float) Sub(other Number) (Number, error) {
	otherF, err := other.Float64()
	if err != nil {
		return nil, fmt.Errorf("cannot subtract: %w", err)
	}
	fF, err := f.Float64()
	if err != nil {
		return nil, fmt.Errorf("cannot subtract: %w", err)
	}
	// Use Float64() to handle both regular numbers and boolean values
	return NewFloat(fF - otherF), nil
}

// Mul returns the product of two float numbers.
func (f *Float) Mul(other Number) (Number, error) {
	otherF, err := other.Float64()
	if err != nil {
		return nil, fmt.Errorf("cannot multiply: %w", err)
	}
	fF, err := f.Float64()
	if err != nil {
		return nil, fmt.Errorf("cannot multiply: %w", err)
	}
	// Use Float64() to handle both regular numbers and boolean values
	return NewFloat(fF * otherF), nil
}

// Div returns the quotient of two float numbers.
func (f *Float) Div(other Number) (Number, error) {
	otherF, err := other.Float64()
	if err != nil {
		return nil, fmt.Errorf("cannot divide: %w", err)
	}
	fF, err := f.Float64()
	if err != nil {
		return nil, fmt.Errorf("cannot divide: %w", err)
	}
	if other.IsZero() {
		return nil, fmt.Errorf("division by zero")
	}
	// Use Float64() to handle both regular numbers and boolean values
	return NewFloat(fF / otherF), nil
}

// Pow returns this float raised to the power of another.
func (f *Float) Pow(other Number) (Number, error) {
	otherF, err := other.Float64()
	if err != nil {
		return nil, fmt.Errorf("cannot power: %w", err)
	}
	fF, err := f.Float64()
	if err != nil {
		return nil, fmt.Errorf("cannot power: %w", err)
	}
	// Use Float64() to handle both regular numbers and boolean values
	return NewFloat(math.Pow(fF, otherF)), nil
}

// Mod returns the remainder of this float divided by another.
func (f *Float) Mod(other Number) (Number, error) {
	otherF, err := other.Float64()
	if err != nil {
		return nil, fmt.Errorf("cannot modulo: %w", err)
	}
	fF, err := f.Float64()
	if err != nil {
		return nil, fmt.Errorf("cannot modulo: %w", err)
	}
	if other.IsZero() {
		return nil, fmt.Errorf("modulo by zero")
	}
	// Use Float64() to handle both regular numbers and boolean values
	return NewFloat(math.Mod(fF, otherF)), nil
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
}

// NewRat creates a new Rat number from a float64.
func NewRat(n float64) *Rat {
	r := &big.Rat{}
	r.SetFloat64(n)
	return &Rat{n: r, isBool: false, boolVal: false}
}

// NewRatFromBool creates a new Rat representing a boolean.
func NewRatFromBool(b bool) *Rat {
	r := &big.Rat{}
	if b {
		r.SetInt64(1)
	} else {
		r.SetInt64(0)
	}
	return &Rat{n: r, isBool: true, boolVal: b}
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

// Add returns the sum of two rational numbers.
func (r *Rat) Add(other Number) (Number, error) {
	otherRat, ok := other.(*Rat)
	if !ok {
		return nil, fmt.Errorf("cannot add: operand is not a rational number")
	}
	result := &big.Rat{}
	result.Add(r.n, otherRat.n)
	return &Rat{n: result}, nil
}

// Sub returns the difference of two rational numbers.
func (r *Rat) Sub(other Number) (Number, error) {
	otherRat, ok := other.(*Rat)
	if !ok {
		return nil, fmt.Errorf("cannot subtract: operand is not a rational number")
	}
	result := &big.Rat{}
	result.Sub(r.n, otherRat.n)
	return &Rat{n: result}, nil
}

// Mul returns the product of two rational numbers.
func (r *Rat) Mul(other Number) (Number, error) {
	otherRat, ok := other.(*Rat)
	if !ok {
		return nil, fmt.Errorf("cannot multiply: operand is not a rational number")
	}
	result := &big.Rat{}
	result.Mul(r.n, otherRat.n)
	return &Rat{n: result}, nil
}

// Div returns the quotient of two rational numbers.
func (r *Rat) Div(other Number) (Number, error) {
	if other.IsZero() {
		return nil, fmt.Errorf("division by zero")
	}
	otherRat, ok := other.(*Rat)
	if !ok {
		return nil, fmt.Errorf("cannot divide: operand is not a rational number")
	}
	result := &big.Rat{}
	result.Quo(r.n, otherRat.n)
	return &Rat{n: result}, nil
}

// Pow returns this rational raised to the power of another.
func (r *Rat) Pow(other Number) (Number, error) {
	otherF, err := other.Float64()
	if err != nil {
		return nil, fmt.Errorf("cannot power: %w", err)
	}
	// For rational powers, convert to float and back
	// This may lose precision but is necessary for non-integer exponents
	power := otherF
	result := &big.Rat{}
	f, _ := r.n.Float64()
	result.SetFloat64(math.Pow(f, power))
	return &Rat{n: result}, nil
}

// Mod returns the remainder of this rational divided by another.
func (r *Rat) Mod(other Number) (Number, error) {
	if other.IsZero() {
		return nil, fmt.Errorf("modulo by zero")
	}
	// For rational modulo, use float64 conversion
	// This may lose precision but is necessary for non-integer moduli
	result := &big.Rat{}
	f1, _ := r.n.Float64()
	f2, err := other.Float64()
	if err != nil {
		return nil, fmt.Errorf("cannot modulo: %w", err)
	}
	result.SetFloat64(math.Mod(f1, f2))
	return &Rat{n: result}, nil
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

// Other methods return errors for strings
func (s *StringNum) Add(other Number) (Number, error) { return nil, fmt.Errorf("string not supported for addition") }
func (s *StringNum) Sub(other Number) (Number, error) { return nil, fmt.Errorf("string not supported for subtraction") }
func (s *StringNum) Mul(other Number) (Number, error) { return nil, fmt.Errorf("string not supported for multiplication") }
func (s *StringNum) Div(other Number) (Number, error) { return nil, fmt.Errorf("string not supported for division") }
func (s *StringNum) Pow(other Number) (Number, error) { return nil, fmt.Errorf("string not supported for power") }
func (s *StringNum) Mod(other Number) (Number, error) { return nil, fmt.Errorf("string not supported for modulo") }
func (s *StringNum) IsZero() bool                     { return false }
func (s *StringNum) IsNegative() bool                 { return false }
func (s *StringNum) Compare(other Number) (int, error) { return 0, fmt.Errorf("string not supported for comparison") }
func (s *StringNum) Bool() (bool, error)              { return false, fmt.Errorf("string not supported for Bool()") }
func (s *StringNum) IsBool() bool                     { return false }
func (s *StringNum) IsSymbol() bool                   { return false }

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

// Other methods return errors for symbols
func (s *Symbol) Add(other Number) (Number, error) {
	return nil, fmt.Errorf("symbol not supported for addition")
}
func (s *Symbol) Sub(other Number) (Number, error) {
	return nil, fmt.Errorf("symbol not supported for subtraction")
}
func (s *Symbol) Mul(other Number) (Number, error) {
	return nil, fmt.Errorf("symbol not supported for multiplication")
}
func (s *Symbol) Div(other Number) (Number, error) {
	return nil, fmt.Errorf("symbol not supported for division")
}
func (s *Symbol) Pow(other Number) (Number, error) {
	return nil, fmt.Errorf("symbol not supported for power")
}
func (s *Symbol) Mod(other Number) (Number, error) {
	return nil, fmt.Errorf("symbol not supported for modulo")
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
