// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewVariables(t *testing.T) {
	v := NewVariables()
	if v == nil {
		t.Fatal("NewVariables() returned nil")
	}
	if v.Count() != 0 {
		t.Errorf("NewVariables() count = %d, want 0", v.Count())
	}
}

func TestSetVariable(t *testing.T) {
	v := NewVariables()

	tests := []struct {
		name  string
		value float64
	}{
		{"x", 5.0},
		{"pi", 3.14159},
		{"result", 42.0},
		{"alpha", -10.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.SetVariable(tt.name, tt.value)
			if err != nil {
				t.Fatalf("SetVariable(%q, %v) returned error: %v", tt.name, tt.value, err)
			}

			value, exists := v.GetVariable(tt.name)
			if !exists {
				t.Errorf("Variable %q should exist after SetVariable", tt.name)
			}
			if value != tt.value {
				t.Errorf("GetVariable(%q) = %v, want %v", tt.name, value, tt.value)
			}
		})
	}
}

func TestGetVariable(t *testing.T) {
	v := NewVariables()
	_ = v.SetVariable("test", 100.0)

	tests := []struct {
		name     string
		value    float64
		expected bool
	}{
		{"test", 100.0, true},
		{"nonexistent", 0, false},
		{"", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, exists := v.GetVariable(tt.name)
			if exists != tt.expected {
				t.Errorf("GetVariable(%q) exists = %v, want %v", tt.name, exists, tt.expected)
			}
			if exists && value != tt.value {
				t.Errorf("GetVariable(%q) = %v, want %v", tt.name, value, tt.value)
			}
		})
	}
}

func TestDeleteVariable(t *testing.T) {
	v := NewVariables()
	_ = v.SetVariable("temp", 50.0)

	// Delete existing variable
	deleted := v.DeleteVariable("temp")
	if !deleted {
		t.Error("DeleteVariable(\"temp\") should return true for existing variable")
	}
	if v.Count() != 0 {
		t.Errorf("Count after delete = %d, want 0", v.Count())
	}

	// Delete non-existent variable
	deleted = v.DeleteVariable("nonexistent")
	if deleted {
		t.Error("DeleteVariable(\"nonexistent\") should return false for non-existent variable")
	}
}

func TestListVariables(t *testing.T) {
	v := NewVariables()
	_ = v.SetVariable("z", 3.0)
	_ = v.SetVariable("a", 1.0)
	_ = v.SetVariable("m", 2.0)

	infos := v.ListVariables()

	if len(infos) != 3 {
		t.Errorf("ListVariables() returned %d variables, want 3", len(infos))
	}

	// Verify sorted order
	expectedOrder := []string{"a", "m", "z"}
	for i, info := range infos {
		if info.Name != expectedOrder[i] {
			t.Errorf("ListVariables()[%d] name = %q, want %q", i, info.Name, expectedOrder[i])
		}
	}
}

func TestClearVariables(t *testing.T) {
	v := NewVariables()
	_ = v.SetVariable("x", 1.0)
	_ = v.SetVariable("y", 2.0)
	_ = v.SetVariable("z", 3.0)

	v.ClearVariables()

	if v.Count() != 0 {
		t.Errorf("Count after ClearVariables() = %d, want 0", v.Count())
	}

	// Verify all variables are gone
	_, exists := v.GetVariable("x")
	if exists {
		t.Error("Variable x should not exist after ClearVariables()")
	}
}

func TestFormatVariables(t *testing.T) {
	v := NewVariables()
	_ = v.SetVariable("pi", 3.14159)
	_ = v.SetVariable("e", 2.71828)

	formatted := v.FormatVariables()

	if strings.Contains(formatted, "No variables defined") {
		t.Error("Formatted variables should not show 'No variables defined' when variables exist")
	}

	if !strings.Contains(formatted, "pi") || !strings.Contains(formatted, "e") {
		t.Errorf("Formatted variables should contain all variable names, got: %s", formatted)
	}
}

func TestFormatVariablesEmpty(t *testing.T) {
	v := NewVariables()
	formatted := v.FormatVariables()

	if !strings.Contains(formatted, "No variables defined") {
		t.Errorf("Empty variables should show 'No variables defined', got: %s", formatted)
	}
}

func TestHasVariable(t *testing.T) {
	v := NewVariables()
	_ = v.SetVariable("exists", 1.0)

	tests := []struct {
		name     string
		expected bool
	}{
		{"exists", true},
		{"nonexistent", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if v.HasVariable(tt.name) != tt.expected {
				t.Errorf("HasVariable(%q) = %v, want %v", tt.name, v.HasVariable(tt.name), tt.expected)
			}
		})
	}
}

func TestErrVariableNotFound(t *testing.T) {
	v := NewVariables()
	_, exists := v.GetVariable("nonexistent")
	if exists {
		t.Error("GetVariable on non-existent variable should return exists=false")
	}
}

func TestErrInvalidVariableName(t *testing.T) {
	v := NewVariables()
	err := v.SetVariable("", 1.0)
	if err != ErrInvalidVariableName {
		t.Errorf("SetVariable with empty name should return ErrInvalidVariableName, got %v", err)
	}
}

func TestVariableConcurrency(t *testing.T) {
	v := NewVariables()
	done := make(chan bool, 100)

	// Concurrent writes
	for i := 0; i < 50; i++ {
		go func(id int) {
			name := fmt.Sprintf("var%d", id)
			_ = v.SetVariable(name, float64(id))
			done <- true
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 50; i++ {
		go func(id int) {
			name := fmt.Sprintf("var%d", id)
			v.GetVariable(name)
			done <- true
		}(i)
	}

	// Wait for all operations
	for i := 0; i < 100; i++ {
		<-done
	}

	// Verify final count
	if v.Count() != 50 {
		t.Errorf("Final count = %d, want 50", v.Count())
	}
}

// TestVariablesSaveLoad tests saving and loading variables to/from a file
func TestVariablesSaveLoad(t *testing.T) {
	tmpFile := t.TempDir() + "/vars.json"

	// Create and populate a variable store
	vars := NewVariables()
	vars.SetVariable("x", 42)
	vars.SetVariable("y", 100)

	// Save to file
	if err := vars.Save(tmpFile); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Create a new variable store and load from file
	newVars := NewVariables()
	if err := newVars.Load(tmpFile); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify loaded values
	val, exists := newVars.GetVariable("x")
	if !exists || val != 42 {
		t.Errorf("Expected x=42 after load, got %v (exists=%v)", val, exists)
	}

	val, exists = newVars.GetVariable("y")
	if !exists || val != 100 {
		t.Errorf("Expected y=100 after load, got %v (exists=%v)", val, exists)
	}

	// Verify that variables from original store are NOT in the new one
	val, exists = newVars.GetVariable("z")
	if exists {
		t.Errorf("Unexpected variable z in loaded store: %v", val)
	}
}

// TestVariablesSaveLoadNonExistentFile tests that Load handles non-existent files gracefully
func TestVariablesSaveLoadNonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/nonexistent_vars.json"

	vars := NewVariables()
	// Load from non-existent file should not error
	if err := vars.Load(tmpFile); err != nil {
		t.Errorf("Load from non-existent file should not error: %v", err)
	}

	// Verify the store is empty
	if vars.Count() != 0 {
		t.Errorf("Expected empty store after loading non-existent file, got %d variables", vars.Count())
	}
}

// TestVariablesSaveLoadOverwrite tests that Load overwrites existing variables
func TestVariablesSaveLoadOverwrite(t *testing.T) {
	tmpFile := t.TempDir() + "/vars.json"

	// Create and save variables
	vars := NewVariables()
	vars.SetVariable("x", 10)
	vars.SetVariable("y", 20)

	if err := vars.Save(tmpFile); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Modify variables in original store
	vars.SetVariable("x", 100)
	vars.SetVariable("z", 30)

	// Load from file (should overwrite)
	if err := vars.Load(tmpFile); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify values are from file
	val, exists := vars.GetVariable("x")
	if !exists || val != 10 {
		t.Errorf("Expected x=10 after load, got %v", val)
	}

	val, exists = vars.GetVariable("y")
	if !exists || val != 20 {
		t.Errorf("Expected y=20 after load, got %v", val)
	}

	// z should not exist anymore
	val, exists = vars.GetVariable("z")
	if exists {
		t.Errorf("Expected z to be removed, got %v", val)
	}
}

// TestVariablesSaveLoadEmptyStore tests saving and loading an empty variable store
func TestVariablesSaveLoadEmptyStore(t *testing.T) {
	tmpFile := t.TempDir() + "/empty_vars.json"

	// Create empty variable store
	vars := NewVariables()

	// Save empty store
	if err := vars.Save(tmpFile); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load into new store
	newVars := NewVariables()
	if err := newVars.Load(tmpFile); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify empty
	if newVars.Count() != 0 {
		t.Errorf("Expected empty store, got %d variables", newVars.Count())
	}
}
