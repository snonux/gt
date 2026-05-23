// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"math/big"
	"testing"
)

func TestNewNumberValue(t *testing.T) {
	v := NewNumberValue(42.5)
	if !v.IsNumber() {
		t.Error("IsNumber() should be true for number value")
	}
	if v.IsBool() {
		t.Error("IsBool() should be false for number value")
	}
	if got := v.Float64(); got != 42.5 {
		t.Errorf("Float64() = %v, want 42.5", got)
	}
	if got := v.Number(); got != 42.5 {
		t.Errorf("Number() = %v, want 42.5", got)
	}
	if got := v.String(); got != "42.5" {
		t.Errorf("String() = %q, want %q", got, "42.5")
	}
}

func TestNewBoolValue(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		v := NewBoolValue(true)
		if !v.IsBool() {
			t.Error("IsBool() should be true")
		}
		if v.IsNumber() {
			t.Error("IsNumber() should be false")
		}
		if !v.Bool() {
			t.Error("Bool() should be true")
		}
		if got := v.Float64(); got != 1 {
			t.Errorf("Float64() = %v, want 1", got)
		}
		if got := v.Number(); got != 0 {
			t.Errorf("Number() = %v, want 0", got)
		}
		if got := v.String(); got != "true" {
			t.Errorf("String() = %q, want %q", got, "true")
		}
	})

	t.Run("false", func(t *testing.T) {
		v := NewBoolValue(false)
		if !v.IsBool() {
			t.Error("IsBool() should be true")
		}
		if v.IsNumber() {
			t.Error("IsNumber() should be false")
		}
		if v.Bool() {
			t.Error("Bool() should be false")
		}
		if got := v.Float64(); got != 0 {
			t.Errorf("Float64() = %v, want 0", got)
		}
		if got := v.String(); got != "false" {
			t.Errorf("String() = %q, want %q", got, "false")
		}
	})
}

func TestValueString(t *testing.T) {
	tests := []struct {
		name string
		v    Value
		want string
	}{
		{"number 42", NewNumberValue(42.0), "42"},
		{"number 3.14", NewNumberValue(3.14), "3.14"},
		{"bool true", NewBoolValue(true), "true"},
		{"bool false", NewBoolValue(false), "false"},
		{"number -5.5", NewNumberValue(-5.5), "-5.5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestToRat(t *testing.T) {
	t.Run("with Rat", func(t *testing.T) {
		r := NewRat(0.75)
		got := ToRat(r)
		if got == nil {
			t.Fatal("ToRat(Rat) returned nil")
		}
		// Verify it's approximately 3/4
		if got.Cmp(big.NewRat(3, 4)) != 0 {
			t.Errorf("ToRat(Rat) = %v, want 3/4", got)
		}
	})

	t.Run("with Float", func(t *testing.T) {
		f := NewFloat(3.14)
		got := ToRat(f)
		if got != nil {
			t.Errorf("ToRat(Float) = %v, want nil", got)
		}
	})

	t.Run("with nil", func(t *testing.T) {
		var sv StackValue
		got := ToRat(sv)
		if got != nil {
			t.Errorf("ToRat(nil) = %v, want nil", got)
		}
	})
}

func TestToFloatUtility(t *testing.T) {
	t.Run("with Float", func(t *testing.T) {
		f := NewFloat(42.5)
		got := ToFloat(f)
		if got != 42.5 {
			t.Errorf("ToFloat(Float(42.5)) = %v, want 42.5", got)
		}
	})

	t.Run("with Rat", func(t *testing.T) {
		r := NewRat(0.5)
		got := ToFloat(r)
		if got != 0.5 {
			t.Errorf("ToFloat(Rat(0.5)) = %v, want 0.5", got)
		}
	})

	t.Run("with nil", func(t *testing.T) {
		var sv StackValue
		got := ToFloat(sv)
		if got != 0 {
			t.Errorf("ToFloat(nil) = %v, want 0", got)
		}
	})
}
