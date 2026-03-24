package rpn

import (
	"fmt"
	"math"
	"math/big"
)

// Number represents a number that can be used in RPN calculations.
// It can be either a float64 or a *big.Rat for precise rational calculations.
type Number interface {
	// String returns the string representation of the number.
	String() string
	// Float64 returns the float64 representation, or panics if not representable.
	Float64() float64
	// Add returns the sum of this number and another.
	Add(other Number) Number
	// Sub returns the difference of this number and another.
	Sub(other Number) Number
	// Mul returns the product of this number and another.
	Mul(other Number) Number
	// Div returns the quotient of this number and another.
	// Returns (nil, error) if division by zero.
	Div(other Number) (Number, error)
	// Pow returns this number raised to the power of another.
	Pow(other Number) Number
	// Mod returns the remainder of this number divided by another.
	// Returns (nil, error) if modulo by zero.
	Mod(other Number) (Number, error)
	// IsZero returns true if the number is zero.
	IsZero() bool
	// IsNegative returns true if the number is negative.
	IsNegative() bool
	// Compare returns -1, 0, or 1 if this number is less than, equal to, or greater than another.
	Compare(other Number) int
}

// NewNumber creates a Number from a float64 value.
// The actual type depends on the current calculation mode.
func NewNumber(value float64, mode CalculationMode) Number {
	if mode == RationalMode {
		return NewRat(value)
	}
	return &Float{n: value}
}

// Float is a Number implementation using float64.
type Float struct {
	n float64
}

// NewFloat creates a new Float number.
func NewFloat(n float64) *Float {
	return &Float{n: n}
}

// String returns the string representation of the float.
func (f *Float) String() string {
	return fmt.Sprintf("%.10g", f.n)
}

// Float64 returns the float64 value.
func (f *Float) Float64() float64 {
	return f.n
}

// Add returns the sum of two float numbers.
func (f *Float) Add(other Number) Number {
	return NewFloat(f.n + other.Float64())
}

// Sub returns the difference of two float numbers.
func (f *Float) Sub(other Number) Number {
	return NewFloat(f.n - other.Float64())
}

// Mul returns the product of two float numbers.
func (f *Float) Mul(other Number) Number {
	return NewFloat(f.n * other.Float64())
}

// Div returns the quotient of two float numbers.
func (f *Float) Div(other Number) (Number, error) {
	if other.IsZero() {
		return nil, fmt.Errorf("division by zero")
	}
	return NewFloat(f.n / other.Float64()), nil
}

// Pow returns this float raised to the power of another.
func (f *Float) Pow(other Number) Number {
	return NewFloat(math.Pow(f.n, other.Float64()))
}

// Mod returns the remainder of this float divided by another.
func (f *Float) Mod(other Number) (Number, error) {
	if other.IsZero() {
		return nil, fmt.Errorf("modulo by zero")
	}
	return NewFloat(math.Mod(f.n, other.Float64())), nil
}

// IsZero returns true if the float is zero.
func (f *Float) IsZero() bool {
	return f.n == 0
}

// IsNegative returns true if the float is negative.
func (f *Float) IsNegative() bool {
	return f.n < 0
}

// Compare returns -1, 0, or 1 if this float is less than, equal to, or greater than another.
func (f *Float) Compare(other Number) int {
	otherF := other.Float64()
	if f.n < otherF {
		return -1
	}
	if f.n > otherF {
		return 1
	}
	return 0
}

// Rat is a Number implementation using *big.Rat.
type Rat struct {
	n *big.Rat
}

// NewRat creates a new Rat number from a float64.
func NewRat(n float64) *Rat {
	r := &big.Rat{}
	r.SetFloat64(n)
	return &Rat{n: r}
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
	// Format as decimal for consistency with Float
	// Use a reasonable precision
	return r.n.FloatString(10)
}

// Float64 returns the float64 representation.
func (r *Rat) Float64() float64 {
	f, _ := r.n.Float64()
	return f
}

// Add returns the sum of two rational numbers.
func (r *Rat) Add(other Number) Number {
	result := &big.Rat{}
	result.Add(r.n, other.(*Rat).n)
	return &Rat{n: result}
}

// Sub returns the difference of two rational numbers.
func (r *Rat) Sub(other Number) Number {
	result := &big.Rat{}
	result.Sub(r.n, other.(*Rat).n)
	return &Rat{n: result}
}

// Mul returns the product of two rational numbers.
func (r *Rat) Mul(other Number) Number {
	result := &big.Rat{}
	result.Mul(r.n, other.(*Rat).n)
	return &Rat{n: result}
}

// Div returns the quotient of two rational numbers.
func (r *Rat) Div(other Number) (Number, error) {
	if other.IsZero() {
		return nil, fmt.Errorf("division by zero")
	}
	result := &big.Rat{}
	result.Quo(r.n, other.(*Rat).n)
	return &Rat{n: result}, nil
}

// Pow returns this rational raised to the power of another.
func (r *Rat) Pow(other Number) Number {
	// For rational powers, convert to float and back
	// This may lose precision but is necessary for non-integer exponents
	power := other.Float64()
	result := &big.Rat{}
	f, _ := r.n.Float64()
	result.SetFloat64(math.Pow(f, power))
	return &Rat{n: result}
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
	f2 := other.Float64()
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

// Compare returns -1, 0, or 1 if this rational is less than, equal to, or greater than another.
func (r *Rat) Compare(other Number) int {
	return r.n.Cmp(other.(*Rat).n)
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
	return n.Float64()
}
