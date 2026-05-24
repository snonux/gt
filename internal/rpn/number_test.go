// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"math"
	"testing"
)

// === Float tests ===

func TestFloatString(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  string
	}{
		{"integer", 42, "42"},
		{"negative", -5.5, "-5.5"},
		{"zero", 0, "0"},
		{"decimal", 0.1, "0.1"},
		{"large", 1e15, "1e+15"},
		{"small", 1.5e-10, "1.5e-10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFloat(tt.value)
			if f.String() != tt.want {
				t.Errorf("Float(%v).String() = %q, want %q", tt.value, f.String(), tt.want)
			}
		})
	}
}

func TestFloatBool(t *testing.T) {
	f := NewFloatFromBool(true)
	if f.String() != "true" {
		t.Errorf("FloatFromBool(true).String() = %q, want 'true'", f.String())
	}

	val, err := f.Bool()
	if err != nil {
		t.Fatalf("Bool() returned error: %v", err)
	}
	if !val {
		t.Error("Bool() should return true")
	}

	f64, _ := f.Float64()
	if f64 != 1 {
		t.Errorf("Float64() = %v, want 1", f64)
	}

	f2 := NewFloatFromBool(false)
	if f2.String() != "false" {
		t.Errorf("FloatFromBool(false).String() = %q, want 'false'", f2.String())
	}
}

func TestFloatBoolOnNonBool(t *testing.T) {
	f := NewFloat(42)
	_, err := f.Bool()
	if err == nil {
		t.Error("Bool() on non-bool should return error")
	}
}

func TestFloatIsZero(t *testing.T) {
	if !NewFloat(0).IsZero() {
		t.Error("Float(0).IsZero() should be true")
	}
	if NewFloat(1).IsZero() {
		t.Error("Float(1).IsZero() should be false")
	}
	if !NewFloatFromBool(false).IsZero() {
		t.Error("FloatFromBool(false).IsZero() should be true")
	}
	if NewFloatFromBool(true).IsZero() {
		t.Error("FloatFromBool(true).IsZero() should be false")
	}
}

func TestFloatIsNegative(t *testing.T) {
	if NewFloat(0).IsNegative() {
		t.Error("Float(0).IsNegative() should be false")
	}
	if !NewFloat(-1).IsNegative() {
		t.Error("Float(-1).IsNegative() should be true")
	}
	if NewFloat(1).IsNegative() {
		t.Error("Float(1).IsNegative() should be false")
	}
}


func TestFloatSetMetricCopy(t *testing.T) {
	reg := GetMetricRegistry()
	mbps, _ := reg.Find("Mbps")

	f := NewFloat(100)
	if f.Metric().Name != "Cool" {
		t.Errorf("initial metric = %q, want 'Cool'", f.Metric().Name)
	}

	f2 := f.SetMetric(mbps).(*Float)
	if f2.Metric().Name != "Mbps" {
		t.Errorf("SetMetric result = %q, want 'Mbps'", f2.Metric().Name)
	}
	// Original should be unchanged
	if f.Metric().Name != "Cool" {
		t.Errorf("original metric changed to %q", f.Metric().Name)
	}
}

func TestFloatCompare(t *testing.T) {
	tests := []struct {
		a    float64
		b    float64
		want int
	}{
		{1, 2, -1},
		{2, 1, 1},
		{1, 1, 0},
		{-1, 1, -1},
		{0, 0, 0},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result, err := NewFloat(tt.a).Compare(NewFloat(tt.b))
			if err != nil {
				t.Fatalf("Compare returned error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Float(%v).Compare(Float(%v)) = %d, want %d", tt.a, tt.b, result, tt.want)
			}
		})
	}
}

func TestFloatFloat64(t *testing.T) {
	f := NewFloat(42.5)
	val, err := f.Float64()
	if err != nil {
		t.Fatalf("Float64() returned error: %v", err)
	}
	if val != 42.5 {
		t.Errorf("Float64() = %v, want 42.5", val)
	}
}

// === Rat tests ===

func TestRatString(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  string
	}{
		{"integer", 42, "42.0000000000"},
		{"negative", -5.5, "-5.5000000000"},
		{"zero", 0, "0.0000000000"},
		{"decimal", 0.1, "0.1000000000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRat(tt.value)
			if r.String() != tt.want {
				t.Errorf("Rat(%v).String() = %q, want %q", tt.value, r.String(), tt.want)
			}
		})
	}
}

func TestRatFromBool(t *testing.T) {
	r := NewRatFromBool(true)
	if r.String() != "true" {
		t.Errorf("RatFromBool(true).String() = %q, want 'true'", r.String())
	}

	val, err := r.Bool()
	if err != nil {
		t.Fatalf("Bool() returned error: %v", err)
	}
	if !val {
		t.Error("Bool() should return true")
	}

	f64, _ := r.Float64()
	if f64 != 1 {
		t.Errorf("Float64() = %v, want 1", f64)
	}
}

func TestRatBoolOnNonBool(t *testing.T) {
	r := NewRat(42)
	_, err := r.Bool()
	if err == nil {
		t.Error("Bool() on non-bool Rat should return error")
	}
}

func TestRatIsZero(t *testing.T) {
	if !NewRat(0).IsZero() {
		t.Error("Rat(0).IsZero() should be true")
	}
	if NewRat(1).IsZero() {
		t.Error("Rat(1).IsZero() should be false")
	}
}

func TestRatIsNegative(t *testing.T) {
	if NewRat(0).IsNegative() {
		t.Error("Rat(0).IsNegative() should be false")
	}
	if !NewRat(-1).IsNegative() {
		t.Error("Rat(-1).IsNegative() should be true")
	}
}


func TestRatSetMetricCopy(t *testing.T) {
	reg := GetMetricRegistry()
	mbps, _ := reg.Find("Mbps")

	r := NewRat(100)
	if r.Metric().Name != "Cool" {
		t.Errorf("initial metric = %q, want 'Cool'", r.Metric().Name)
	}

	r2 := r.SetMetric(mbps).(*Rat)
	if r2.Metric().Name != "Mbps" {
		t.Errorf("SetMetric result = %q, want 'Mbps'", r2.Metric().Name)
	}
	// Original should be unchanged
	if r.Metric().Name != "Cool" {
		t.Errorf("original metric changed to %q", r.Metric().Name)
	}
}

func TestRatCompare(t *testing.T) {
	tests := []struct {
		a    float64
		b    float64
		want int
	}{
		{1, 2, -1},
		{2, 1, 1},
		{1, 1, 0},
		{-1, 1, -1},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result, err := NewRat(tt.a).Compare(NewRat(tt.b))
			if err != nil {
				t.Fatalf("Compare returned error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Rat(%v).Compare(Rat(%v)) = %d, want %d", tt.a, tt.b, result, tt.want)
			}
		})
	}
}

func TestRatCompareWithFloat(t *testing.T) {
	result, err := NewRat(3).Compare(NewFloat(2))
	if err != nil {
		t.Fatalf("Compare returned error: %v", err)
	}
	if result != 1 {
		t.Errorf("Rat(3).Compare(Float(2)) = %d, want 1", result)
	}
}

func TestNewRatFromString(t *testing.T) {
	tests := []struct {
		s     string
		valid bool
	}{
		{"1/2", true},
		{"3", true},
		{"-7/4", true},
		{"0", true},
		{"abc", false},
		{"", false},
		{"1/", false},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			r, err := NewRatFromString(tt.s)
			if tt.valid {
				if err != nil {
					t.Errorf("NewRatFromString(%q) error: %v", tt.s, err)
				}
				if r == nil {
					t.Error("NewRatFromString returned nil")
				}
			} else {
				if err == nil {
					t.Errorf("NewRatFromString(%q) should have returned error", tt.s)
				}
			}
		})
	}
}

// === StringNum tests ===

func TestStringNumString(t *testing.T) {
	s := NewStringNum("hello")
	if s.String() != "hello" {
		t.Errorf("StringNum('hello').String() = %q, want 'hello'", s.String())
	}
}


func TestStringNumMetricNil(t *testing.T) {
	s := NewStringNum("test")
	if s.Metric() == nil {
		t.Error("StringNum.Metric() should not be nil")
	}
}

// === Symbol tests ===

func TestSymbolString(t *testing.T) {
	s := NewSymbol("x")
	if s.String() != ":x" {
		t.Errorf("Symbol('x').String() = %q, want ':x'", s.String())
	}
}

func TestSymbolName(t *testing.T) {
	s := NewSymbol("counter")
	if s.Name() != "counter" {
		t.Errorf("Symbol('counter').Name() = %q, want 'counter'", s.Name())
	}
}


func TestSymbolMetricNil(t *testing.T) {
	s := NewSymbol("z")
	if s.Metric() == nil {
		t.Error("Symbol.Metric() should not be nil")
	}
}

// === NewNumber constructor tests ===

func TestNewNumberFloatMode(t *testing.T) {
	n := NewNumber(42.5, FloatMode)
	if _, ok := n.(*Float); !ok {
		t.Error("NewNumber(FloatMode) should return *Float")
	}
	val, _ := n.Float64()
	if val != 42.5 {
		t.Errorf("Float64() = %v, want 42.5", val)
	}
}

func TestNewNumberRationalMode(t *testing.T) {
	n := NewNumber(42.5, RationalMode)
	if _, ok := n.(*Rat); !ok {
		t.Error("NewNumber(RationalMode) should return *Rat")
	}
	val, _ := n.Float64()
	if val != 42.5 {
		t.Errorf("Float64() = %v, want 42.5", val)
	}
}

func TestNewNumberWithExplicitMetric(t *testing.T) {
	reg := GetMetricRegistry()
	mbps, _ := reg.Find("Mbps")

	n := NewNumberWithMetric(100, FloatMode, mbps)
	val, _ := n.Float64()
	if val != 100 {
		t.Errorf("Float64() = %v, want 100", val)
	}
	if n.Metric().Name != "Mbps" {
		t.Errorf("Metric = %q, want 'Mbps'", n.Metric().Name)
	}
}

func TestNewNumberDefaultsToCool(t *testing.T) {
	n := NewNumber(10, FloatMode)
	if n.Metric().Name != "Cool" {
		t.Errorf("default Metric = %q, want 'Cool'", n.Metric().Name)
	}
}

// === Interface satisfaction tests ===

func TestFloatImplementsStackValue(t *testing.T) {
	var sv StackValue = NewFloat(42)
	if sv.String() != "42" {
		t.Errorf("StackValue interface broken: %q", sv.String())
	}
}

func TestRatImplementsStackValue(t *testing.T) {
	var sv StackValue = NewRat(42)
	// Rat.String() uses FloatString(10) which includes trailing zeros
	if sv.String() != "42.0000000000" {
		t.Errorf("StackValue interface broken: %q", sv.String())
	}
}

func TestStringNumImplementsStackValue(t *testing.T) {
	var sv StackValue = NewStringNum("hello")
	if sv.String() != "hello" {
		t.Errorf("StackValue interface broken: %q", sv.String())
	}
}

func TestSymbolImplementsStackValue(t *testing.T) {
	var sv StackValue = NewSymbol("x")
	if sv.String() != ":x" {
		t.Errorf("StackValue interface broken: %q", sv.String())
	}
}

func TestFloatImplementsNumericValue(t *testing.T) {
	var nv NumericValue = NewFloat(42)
	val, err := nv.Float64()
	if err != nil {
		t.Fatalf("NumericValue interface broken: %v", err)
	}
	if val != 42 {
		t.Errorf("Float64() = %v, want 42", val)
	}
}

func TestRatImplementsNumericValue(t *testing.T) {
	var nv NumericValue = NewRat(42)
	val, err := nv.Float64()
	if err != nil {
		t.Fatalf("NumericValue interface broken: %v", err)
	}
	if val != 42 {
		t.Errorf("Float64() = %v, want 42", val)
	}
}

// === Helper function tests ===

func TestToFloat(t *testing.T) {
	if ToFloat(NewFloat(42.5)) != 42.5 {
		t.Error("ToFloat(Float(42.5)) should return 42.5")
	}
	if ToFloat(NewRat(42.5)) != 42.5 {
		t.Error("ToFloat(Rat(42.5)) should return 42.5")
	}
	if ToFloat(NewSymbol("x")) != 0 {
		t.Error("ToFloat(Symbol) should return 0")
	}
}

// === Edge cases ===

func TestFloatStringNaN(t *testing.T) {
	f := NewFloat(math.NaN())
	s := f.String()
	// NaN formatted as "NaN"
	if s != "NaN" {
		t.Errorf("Float(NaN).String() = %q, expected 'NaN'", s)
	}
}

func TestFloatStringInf(t *testing.T) {
	f := NewFloat(math.Inf(1))
	if f.String() != "+Inf" {
		t.Errorf("Float(+Inf).String() = %q, want '+Inf'", f.String())
	}
	f2 := NewFloat(math.Inf(-1))
	if f2.String() != "-Inf" {
		t.Errorf("Float(-Inf).String() = %q, want '-Inf'", f2.String())
	}
}
