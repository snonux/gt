// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestVariablesSaveLoadCorruptFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "corrupt.json")

	// Write invalid JSON
	if err := os.WriteFile(tmpFile, []byte("{invalid json"), 0644); err != nil {
		t.Fatalf("failed to create corrupt file: %v", err)
	}

	vars := NewVariables()
	err := vars.Load(tmpFile)
	if err == nil {
		t.Error("Load from corrupt JSON file should return an error")
	}
	if !strings.Contains(err.Error(), "unmarshal") {
		t.Errorf("error should mention unmarshal, got: %v", err)
	}
}

func TestVariablesLoadFiltersInvalidVariableNames(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "vars.json")

	// Write JSON with valid and invalid variable names
	jsonData := `[
		{"Name": "valid_name", "Value": 42},
		{"Name": "alsoValid123", "Value": 99.5},
		{"Name": "", "Value": 10},
		{"Name": "has space", "Value": 20},
		{"Name": "has@symbol", "Value": 30},
		{"Name": "valid_too", "Value": -5}
	]`
	if err := os.WriteFile(tmpFile, []byte(jsonData), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	vars := NewVariables()
	if err := vars.Load(tmpFile); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Valid variables should be loaded
	if vars.Count() != 3 {
		t.Errorf("should have 3 valid variables, got %d", vars.Count())
	}

	if val, exists := vars.GetVariable("valid_name"); !exists || val != 42 {
		t.Errorf("valid_name = %v (exists=%v), want 42", val, exists)
	}
	if val, exists := vars.GetVariable("alsoValid123"); !exists || val != 99.5 {
		t.Errorf("alsoValid123 = %v (exists=%v), want 99.5", val, exists)
	}
	if val, exists := vars.GetVariable("valid_too"); !exists || val != -5 {
		t.Errorf("valid_too = %v (exists=%v), want -5", val, exists)
	}

	// Invalid names should not be loaded
	if vars.HasVariable("") {
		t.Error("empty variable name should not be loaded")
	}
	if vars.HasVariable("has space") {
		t.Error("variable with space should not be loaded")
	}
	if vars.HasVariable("has@symbol") {
		t.Error("variable with symbol should not be loaded")
	}
}

func TestVariablesSaveLoadRoundTripDiverseValues(t *testing.T) {
	tests := []struct {
		name  string
		value float64
	}{
		{"zero", 0},
		{"negative", -42.5},
		{"positive", 42.5},
		{"veryLarge", 1.7976931348623157e+308}, // near max float64
		{"verySmall", 5e-324},                  // near min float64
		{"scientific", 1.23e-10},
		{"pi", 3.141592653589793},
		{"integer", 100},
		{"decimal", 0.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := filepath.Join(t.TempDir(), "vars.json")

			// Create and save
			vars := NewVariables()
			err := vars.SetVariable("v", tt.value)
			if err != nil {
				t.Fatalf("SetVariable failed: %v", err)
			}
			if err := vars.Save(tmpFile); err != nil {
				t.Fatalf("Save failed: %v", err)
			}

			// Load and verify
			newVars := NewVariables()
			if err := newVars.Load(tmpFile); err != nil {
				t.Fatalf("Load failed: %v", err)
			}

			loadedVal, exists := newVars.GetVariable("v")
			if !exists {
				t.Fatal("variable 'v' not found after load")
			}
			if loadedVal != tt.value {
				t.Errorf("value = %v, want %v", loadedVal, tt.value)
			}
		})
	}
}

func TestVariablesSaveFileFormat(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "vars.json")

	vars := NewVariables()
	_ = vars.SetVariable("z", 3.0)
	_ = vars.SetVariable("a", 1.0)

	if err := vars.Save(tmpFile); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify JSON structure
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	content := string(data)
	// Variables should be sorted by name, so "a" appears before "z"
	aIdx := strings.Index(content, `"a"`)
	zIdx := strings.Index(content, `"z"`)
	if aIdx >= zIdx {
		t.Errorf("variables should be sorted by name in JSON, got: %s", content)
	}
	// Should use indented JSON (MarshalIndent)
	if !strings.Contains(content, "\n  ") {
		t.Errorf("expected indented JSON, got: %s", content)
	}
}

func TestVariablesSaveLoadMultipleVariables(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "vars.json")

	// Create and save many variables
	vars := NewVariables()
	for i := 0; i < 50; i++ {
		name := string(rune('a'+i%26)) + string(rune('a'+(i/26)%26))
		_ = vars.SetVariable(name, float64(i)*1.5)
	}

	if err := vars.Save(tmpFile); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load and verify all
	newVars := NewVariables()
	if err := newVars.Load(tmpFile); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if newVars.Count() != 50 {
		t.Errorf("Count = %d, want 50", newVars.Count())
	}

	for i := 0; i < 50; i++ {
		name := string(rune('a'+i%26)) + string(rune('a'+(i/26)%26))
		want := float64(i) * 1.5
		got, exists := newVars.GetVariable(name)
		if !exists {
			t.Errorf("variable %q not found", name)
			continue
		}
		if got != want {
			t.Errorf("%s = %v, want %v", name, got, want)
		}
	}
}

func TestVariablesSaveLoadConcurrency(t *testing.T) {
	vars := NewVariables()
	for i := 0; i < 10; i++ {
		_ = vars.SetVariable(string(rune('a'+i)), float64(i))
	}

	// Save to a stable file first
	tmpFile := filepath.Join(t.TempDir(), "vars.json")
	if err := vars.Save(tmpFile); err != nil {
		t.Fatalf("initial Save failed: %v", err)
	}

	errChan := make(chan error, 10)

	// Concurrent loads from the same file (reads are safe)
	for i := 0; i < 10; i++ {
		go func() {
			v := NewVariables()
			errChan <- v.Load(tmpFile)
		}()
	}

	for i := 0; i < 10; i++ {
		if err := <-errChan; err != nil {
			t.Errorf("concurrent load failed: %v", err)
		}
	}

	// Test concurrent SetVariable/GetVariable on same instance
	vars2 := NewVariables()
	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_ = vars2.SetVariable(string(rune('x'+id%26)), float64(id))
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			vars2.GetVariable(string(rune('x' + id%26)))
		}(i)
	}

	// Concurrent count/list
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = vars2.Count()
			_ = vars2.ListVariables()
		}()
	}

	wg.Wait()

	// Verify final state is consistent
	if vars2.Count() > 26 {
		t.Errorf("Count = %d, expected at most 26 (x0-x25)", vars2.Count())
	}
}

func TestVariablesLoadEmptyJSON(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "empty.json")

	// Write empty JSON array
	if err := os.WriteFile(tmpFile, []byte("[]"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	vars := NewVariables()
	if err := vars.Load(tmpFile); err != nil {
		t.Fatalf("Load of empty JSON array failed: %v", err)
	}

	if vars.Count() != 0 {
		t.Errorf("Count = %d, want 0", vars.Count())
	}
}

func TestVariablesLoadPreservesExistingAfterError(t *testing.T) {
	vars := NewVariables()
	if err := vars.SetVariable("existing", 42.0); err != nil {
		t.Fatalf("SetVariable failed: %v", err)
	}

	// Load from corrupt file
	tmpFile := filepath.Join(t.TempDir(), "corrupt.json")
	if err := os.WriteFile(tmpFile, []byte("not json"), 0644); err != nil {
		t.Fatalf("failed to create corrupt file: %v", err)
	}

	err := vars.Load(tmpFile)
	if err == nil {
		t.Fatal("expected error loading corrupt file")
	}

	// Load() calls loadVariables() before acquiring the write lock.
	// When loadVariables() errors, Load() returns the error without
	// modifying the map, so existing state is preserved.
	val, exists := vars.GetVariable("existing")
	if !exists || val != 42.0 {
		t.Errorf("existing variable should be preserved after failed Load: got %v (exists=%v)", val, exists)
	}
}

func TestVariablesSaveInfReturnsError(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "inf.json")

	vars := NewVariables()
	// Inf values cannot be marshaled to JSON, so Save should fail
	if err := vars.SetVariable("inf", math.Inf(1)); err != nil {
		t.Fatalf("SetVariable failed: %v", err)
	}
	err := vars.Save(tmpFile)
	if err == nil {
		t.Error("Save with Inf value should return an error (JSON cannot represent Inf)")
	}
	if !strings.Contains(err.Error(), "unsupported value") {
		t.Errorf("error should mention unsupported value, got: %v", err)
	}
}

func TestVariablesLoadDuplicateNames(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "dupes.json")

	// JSON with duplicate names — last one wins in Go maps
	jsonData := `[{"Name":"x","Value":1},{"Name":"x","Value":2}]`
	if err := os.WriteFile(tmpFile, []byte(jsonData), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	vars := NewVariables()
	if err := vars.Load(tmpFile); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should have exactly 1 variable (x), with value 2 (last wins)
	if vars.Count() != 1 {
		t.Errorf("Count = %d, want 1", vars.Count())
	}
	val, exists := vars.GetVariable("x")
	if !exists || val != 2 {
		t.Errorf("x = %v (exists=%v), want 2 (last wins)", val, exists)
	}
}

func TestVariablesLoadEmptyFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "empty.json")

	// Write empty file (not even valid JSON)
	if err := os.WriteFile(tmpFile, []byte{}, 0644); err != nil {
		t.Fatalf("failed to create empty file: %v", err)
	}

	vars := NewVariables()
	err := vars.Load(tmpFile)
	if err == nil {
		t.Error("Load from empty file should return an error (invalid JSON)")
	}
}
