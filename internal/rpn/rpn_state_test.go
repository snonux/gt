// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"testing"
)

func TestRPNGetConstants(t *testing.T) {
	v := NewVariables()
	r := NewRPN(v)
	consts := r.GetConstants()
	if consts == nil {
		t.Fatal("GetConstants() returned nil")
	}
	// Verify we can actually use it
	val, ok := consts.GetConstant("pi")
	if !ok {
		t.Error("pi constant should exist")
	}
	if val <= 0 {
		t.Errorf("pi constant = %v, want > 0", val)
	}
}

func TestRPNGetModeAndSetMode(t *testing.T) {
	v := NewVariables()
	r := NewRPN(v)

	// Default should be FloatMode
	if r.GetMode() != FloatMode {
		t.Errorf("default mode = %v, want %v", r.GetMode(), FloatMode)
	}

	// Switch to RationalMode
	r.SetMode(RationalMode)
	if r.GetMode() != RationalMode {
		t.Errorf("mode after SetMode(RationalMode) = %v, want %v", r.GetMode(), RationalMode)
	}

	// Switch back
	r.SetMode(FloatMode)
	if r.GetMode() != FloatMode {
		t.Errorf("mode after SetMode(FloatMode) = %v, want %v", r.GetMode(), FloatMode)
	}
}

func TestRPNSetCurrentStack(t *testing.T) {
	v := NewVariables()
	r := NewRPN(v)

	// Start with empty stack
	if got := r.GetCurrentStack(); len(got) != 0 {
		t.Errorf("initial stack len = %d, want 0", len(got))
	}

	// Set a new stack
	values := []StackValue{
		NewNumber(1.0, FloatMode),
		NewNumber(2.0, FloatMode),
		NewNumber(3.0, FloatMode),
	}
	r.SetCurrentStack(values)

	got := r.GetCurrentStack()
	if len(got) != 3 {
		t.Fatalf("stack len = %d, want 3", len(got))
	}
	for i, sv := range got {
		if f := ToFloat(sv); f != float64(i+1) {
			t.Errorf("stack[%d] = %v, want %d", i, sv, i+1)
		}
	}
}

func TestRPNStack(t *testing.T) {
	v := NewVariables()
	r := NewRPN(v)

	values := []StackValue{
		NewNumber(42.0, FloatMode),
		NewNumber(99.0, FloatMode),
	}
	r.SetCurrentStack(values)

	// Stack() should return same as GetCurrentStack()
	stack := r.Stack()
	current := r.GetCurrentStack()
	if len(stack) != len(current) {
		t.Errorf("Stack() len = %d, GetCurrentStack() len = %d", len(stack), len(current))
	}
}

func TestRPNIsStandardOperator(t *testing.T) {
	v := NewVariables()
	r := NewRPN(v)

	tests := []struct {
		token    string
		expected bool
	}{
		{"+", true},
		{"-", true},
		{"*", true},
		{"/", true},
		{"^", true},
		{"dup", true},
		{"swap", true},
		{"pop", true},
		{"show", true},
		{"notanop", false},
		{"xyz", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			if got := r.IsStandardOperator(tt.token); got != tt.expected {
				t.Errorf("IsStandardOperator(%q) = %v, want %v", tt.token, got, tt.expected)
			}
		})
	}
}

func TestRPNIsHyperOperator(t *testing.T) {
	v := NewVariables()
	r := NewRPN(v)

	tests := []struct {
		token    string
		expected bool
	}{
		{"[+]", true},
		{"[-]", true},
		{"[*]", true},
		{"[/]", true},
		{"[^]", true},
		{"[lg]", true},
		{"[log]", true},
		{"[ln]", true},
		{"+", false},
		{"notanop", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			if got := r.IsHyperOperator(tt.token); got != tt.expected {
				t.Errorf("IsHyperOperator(%q) = %v, want %v", tt.token, got, tt.expected)
			}
		})
	}
}

func TestOperationsMetricRegistry(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)

	reg := o.MetricRegistry()
	if reg == nil {
		t.Fatal("MetricRegistry() returned nil")
	}

	// Should be able to look up a known metric
	metric, ok := reg.Find("Cool")
	if !ok {
		t.Fatal("Find(Cool) should succeed")
	}
	if metric.Name != "Cool" {
		t.Errorf("metric name = %q, want %q", metric.Name, "Cool")
	}
}

func TestRPNSetPrefixMode(t *testing.T) {
	v := NewVariables()
	r := NewRPN(v)

	// Default should be SI
	if r.GetPrefixMode() != SI {
		t.Errorf("default prefix mode = %v, want %v", r.GetPrefixMode(), SI)
	}

	r.SetPrefixMode(IEC)
	if r.GetPrefixMode() != IEC {
		t.Errorf("prefix mode after SetPrefixMode(IEC) = %v, want %v", r.GetPrefixMode(), IEC)
	}

	r.SetPrefixMode(SI)
	if r.GetPrefixMode() != SI {
		t.Errorf("prefix mode after SetPrefixMode(SI) = %v, want %v", r.GetPrefixMode(), SI)
	}
}
