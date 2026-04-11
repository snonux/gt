// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"math"
	"strings"
	"testing"
)

func TestNewConstants(t *testing.T) {
	c := NewConstants()
	if c == nil {
		t.Fatal("NewConstants() returned nil")
	}
	if c.Count() == 0 {
		t.Error("NewConstants() should have built-in constants")
	}
}

func TestConstants_GetConstant(t *testing.T) {
	c := NewConstants()

	tests := []struct {
		name     string
		key      string
		expected float64
	}{
		{"pi", "pi", math.Pi},
		{"Pi with Greek letter", "π", math.Pi},
		{"Euler's number", "e", math.E},
		{"Euler with name", "euler", math.E},
		{"Golden ratio", "phi", 1.618033988749895},
		{"Golden ratio with Greek letter", "φ", 1.618033988749895},
		{"Square root of 2", "sqrt2", 1.414213562373095},
		{"Square root of 2 with symbol", "√2", 1.414213562373095},
		{"Infinity", "inf", math.Inf(1)},
		{"Infinity with name", "infinity", math.Inf(1)},
		{"NaN", "nan", math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, exists := c.GetConstant(tt.key)
			if !exists {
				t.Errorf("Constant %q should exist", tt.key)
				return
			}
			if tt.key == "nan" {
				if !math.IsNaN(val) {
					t.Errorf("Constant %q = %v, want NaN", tt.key, val)
				}
			} else if tt.key == "inf" || tt.key == "infinity" {
				if val != math.Inf(1) {
					t.Errorf("Constant %q = %v, want %v", tt.key, val, math.Inf(1))
				}
			} else {
				if val != tt.expected {
					t.Errorf("Constant %q = %v, want %v", tt.key, val, tt.expected)
				}
			}
		})
	}
}

func TestConstants_GetConstant_NonExistent(t *testing.T) {
	c := NewConstants()

	val, exists := c.GetConstant("nonexistent")
	if exists {
		t.Errorf("Non-existent constant should return exists=false")
	}
	if val != 0 {
		t.Errorf("Non-existent constant value should be 0, got %v", val)
	}
}

func TestConstants_Count(t *testing.T) {
	c := NewConstants()
	count := c.Count()

	// Should have at least: pi, e, phi, sqrt2, inf, nan = 6
	if count < 6 {
		t.Errorf("Count = %d, want at least 6 built-in constants", count)
	}
}

func TestConstants_HasConstant(t *testing.T) {
	c := NewConstants()

	if !c.HasConstant("pi") {
		t.Error("HasConstant(\"pi\") should return true")
	}
	if !c.HasConstant("e") {
		t.Error("HasConstant(\"e\") should return true")
	}
	if c.HasConstant("nonexistent") {
		t.Error("HasConstant(\"nonexistent\") should return false")
	}
}

func TestConstants_ListConstants(t *testing.T) {
	c := NewConstants()
	infos := c.ListConstants()

	if len(infos) == 0 {
		t.Error("ListConstants() should return at least one constant")
	}

	// Check that pi is in the list
	foundPi := false
	foundE := false
	for _, info := range infos {
		if info.Name == "pi" {
			foundPi = true
			if info.Value != math.Pi {
				t.Errorf("pi constant value = %v, want %v", info.Value, math.Pi)
			}
		}
		if info.Name == "e" {
			foundE = true
			if info.Value != math.E {
				t.Errorf("e constant value = %v, want %v", info.Value, math.E)
			}
		}
	}

	if !foundPi {
		t.Error("pi constant not found in ListConstants() output")
	}
	if !foundE {
		t.Error("e constant not found in ListConstants() output")
	}
}

func TestConstants_ListConstantsSorted(t *testing.T) {
	c := NewConstants()
	infos := c.ListConstants()

	// Verify the list is sorted alphabetically by name
	for i := 0; i < len(infos)-1; i++ {
		if infos[i].Name > infos[i+1].Name {
			t.Errorf("ListConstants() not sorted: %q > %q", infos[i].Name, infos[i+1].Name)
		}
	}
}

func TestConstants_SetConstant(t *testing.T) {
	c := NewConstants()

	// Test setting a custom constant
	err := c.SetConstant("custom", 42.0)
	if err != nil {
		t.Errorf("SetConstant() returned error: %v", err)
	}

	val, exists := c.GetConstant("custom")
	if !exists {
		t.Error("Custom constant should exist after SetConstant()")
	}
	if val != 42.0 {
		t.Errorf("Custom constant value = %v, want 42.0", val)
	}
}

func TestConstants_ClearConstants(t *testing.T) {
	c := NewConstants()

	// First add a user-defined constant
	c.SetConstant("custom", 42.0)

	// Clear the constants
	c.ClearConstants()

	// Built-in constants should still exist (pi, e, phi, sqrt2, inf, nan = 6 + 1 duplicate each = 11)
	// But custom should be removed
	count := c.Count()
	if count < 11 {
		t.Errorf("Count after ClearConstants() = %d, want at least 11 (built-in constants)", count)
	}

	// Verify custom constant is gone
	_, exists := c.GetConstant("custom")
	if exists {
		t.Error("Custom constant should be removed after ClearConstants()")
	}
}

func TestConstants_ThreadSafety(t *testing.T) {
	c := NewConstants()

	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			name := "thread" + string(rune(id+'0'))
			c.SetConstant(name, float64(id*10))
			val, exists := c.GetConstant(name)
			if !exists {
				t.Errorf("Thread %d: constant %q should exist", id, name)
			}
			if val != float64(id*10) {
				t.Errorf("Thread %d: constant %q = %v, want %v", id, name, val, float64(id*10))
			}
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestConstants_RetrieveInRPN(t *testing.T) {
	v := NewVariables()
	r := NewRPN(v)

	// Test pi constant
	result, err := r.ParseAndEvaluate("pi")
	if err != nil {
		t.Fatalf("ParseAndEvaluate(\"pi\") returned error: %v", err)
	}
	if result != "3.141592654" {
		t.Errorf("ParseAndEvaluate(\"pi\") = %q, want \"3.141592654\"", result)
	}

	// Create a new RPN instance for each test to avoid stack state conflicts
	v2 := NewVariables()
	r2 := NewRPN(v2)

	// Test e constant
	result, err = r2.ParseAndEvaluate("e")
	if err != nil {
		t.Fatalf("ParseAndEvaluate(\"e\") returned error: %v", err)
	}
	if result != "2.718281828" {
		t.Errorf("ParseAndEvaluate(\"e\") = %q, want \"2.718281828\"", result)
	}

	// Create a new RPN instance for the next test
	v3 := NewVariables()
	r3 := NewRPN(v3)

	// Test pi in expression
	result, err = r3.ParseAndEvaluate("pi 2 *")
	if err != nil {
		t.Fatalf("ParseAndEvaluate(\"pi 2 *\") returned error: %v", err)
	}
	if result != "6.283185307" {
		t.Errorf("ParseAndEvaluate(\"pi 2 *\") = %q, want \"6.283185307\"", result)
	}

	// Create a new RPN instance for phi test
	v4 := NewVariables()
	r4 := NewRPN(v4)

	// Test phi constant
	result, err = r4.ParseAndEvaluate("phi")
	if err != nil {
		t.Fatalf("ParseAndEvaluate(\"phi\") returned error: %v", err)
	}
	if result != "1.618033989" {
		t.Errorf("ParseAndEvaluate(\"phi\") = %q, want \"1.618033989\"", result)
	}
}

func TestConstants_RetrieveWithGreekLetters(t *testing.T) {
	v := NewVariables()
	r := NewRPN(v)

	// Test pi with Greek letter
	result, err := r.ParseAndEvaluate("π")
	if err != nil {
		t.Fatalf("ParseAndEvaluate(\"π\") returned error: %v", err)
	}
	if result != "3.141592654" {
		t.Errorf("ParseAndEvaluate(\"π\") = %q, want \"3.141592654\"", result)
	}

	// Create a new RPN instance for phi test
	v2 := NewVariables()
	r2 := NewRPN(v2)

	// Test phi with Greek letter
	result, err = r2.ParseAndEvaluate("φ")
	if err != nil {
		t.Fatalf("ParseAndEvaluate(\"φ\") returned error: %v", err)
	}
	if result != "1.618033989" {
		t.Errorf("ParseAndEvaluate(\"φ\") = %q, want \"1.618033989\"", result)
	}
}

func TestConstants_ConflictWithVariables(t *testing.T) {
	v := NewVariables()
	r := NewRPN(v)

	// First set a variable named pi
	result, err := r.ParseAndEvaluate("pi = 3.0")
	if err != nil {
		t.Fatalf("ParseAndEvaluate(\"pi = 3.0\") returned error: %v", err)
	}
	if result != "pi = 3" {
		t.Errorf("ParseAndEvaluate(\"pi = 3.0\") = %q, want \"pi = 3\"", result)
	}

	// Now using pi should get the variable value, not the constant
	result, err = r.ParseAndEvaluate("pi")
	if err != nil {
		t.Fatalf("ParseAndEvaluate(\"pi\") after variable set returned error: %v", err)
	}
	if result != "3" {
		t.Errorf("ParseAndEvaluate(\"pi\") after variable set = %q, want \"3\"", result)
	}
}

func TestConstantsCommand(t *testing.T) {
	v := NewVariables()
	r := NewRPN(v)
	result, err := r.ParseAndEvaluate("constants")
	if err != nil {
		t.Fatalf("ParseAndEvaluate(\"constants\") returned error: %v", err)
	}
	if !strings.Contains(result, "pi") || !strings.Contains(result, "e") {
		t.Errorf("constants output should contain pi and e, got: %s", result)
	}
}

func TestClearConstantsCommand(t *testing.T) {
	v := NewVariables()
	r := NewRPN(v)

	// First add a custom constant
	result, err := r.ParseAndEvaluate("custom 42 =")
	if err != nil {
		t.Fatalf("Failed to set custom constant: %v", err)
	}

	// Clear constants
	result, err = r.ParseAndEvaluate("clearconstants")
	if err != nil {
		t.Fatalf("ParseAndEvaluate(\"clearconstants\") returned error: %v", err)
	}
	if result != "All constants cleared" {
		t.Errorf("clearconstants result = %q, want \"All constants cleared\"", result)
	}

	// Built-in constants should still exist
	result, err = r.ParseAndEvaluate("pi")
	if err != nil {
		t.Fatalf("pi constant should still work after clearconstants: %v", err)
	}
	if result != "3.141592654" {
		t.Errorf("pi after clearconstants = %q, want \"3.141592654\"", result)
	}
}
