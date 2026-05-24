// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
)

// Error variables for external error checking.
var (
	ErrVariableNotFound    = fmt.Errorf("variable not found")
	ErrInvalidVariableName = fmt.Errorf("invalid variable name")
)

// VariableInfo represents a single variable with its name and value.
type VariableInfo struct {
	Name  string
	Value float64
}

// Variables stores variable name-value pairs for RPN calculations.
// It provides thread-safe access to variable storage.
type Variables struct {
	mu        sync.RWMutex
	variables map[string]float64
}

// VariableReader defines the interface for reading variable storage.
type VariableReader interface {
	GetVariable(name string) (float64, bool)
	ListVariables() []VariableInfo
	FormatVariables() string
	Count() int
	HasVariable(name string) bool
}

// VariableWriter defines the interface for writing to variable storage.
type VariableWriter interface {
	SetVariable(name string, value float64) error
	DeleteVariable(name string) bool
	ClearVariables()
}

// VariablePersistence defines the interface for persisting variable storage to disk.
type VariablePersistence interface {
	// Save writes the variable store to a file in JSON format.
	Save(path string) error
	// Load reads the variable store from a file in JSON format.
	// All existing variables are replaced with the loaded values.
	Load(path string) error
}

// VariableStore combines VariableReader, VariableWriter, and VariablePersistence
// for full variable storage access.
type VariableStore interface {
	VariableReader
	VariableWriter
	VariablePersistence
}

// NewVariables creates and initializes a new Variables instance.
func NewVariables() *Variables {
	return &Variables{
		variables: make(map[string]float64),
	}
}

// isValidVariableName checks if a variable name is valid.
// Variable names must be non-empty and contain only alphanumeric characters and underscores.
//
// name: the variable name to validate
// Returns true if the name is valid, false otherwise
func isValidVariableName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		// Check if character is NOT alphanumeric or underscore
		// Apply De Morgan's law: !(P || Q || R || S) == !P && !Q && !R && !S
		// where P = 'a' <= r && r <= 'z' (lowercase)
		// Q = 'A' <= r && r <= 'Z' (uppercase)
		// R = '0' <= r && r <= '9' (digit)
		// S = r == '_' (underscore)
		// !P = r < 'a' || r > 'z'
		// !Q = r < 'A' || r > 'Z'
		// !R = r < '0' || r > '9'
		// !S = r != '_'
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '_' {
			return false
		}
	}
	return true
}

// SetVariable assigns a value to a variable name.
// Usage: `name value =` stores value in variable.
func (v *Variables) SetVariable(name string, value float64) error {
	if !isValidVariableName(name) {
		return ErrInvalidVariableName
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	v.variables[name] = value
	return nil
}

// GetVariable retrieves the value of a variable.
// Returns the value and true if found, or 0 and false if not found.
func (v *Variables) GetVariable(name string) (float64, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	value, exists := v.variables[name]
	return value, exists
}

// DeleteVariable removes a variable from storage.
// Usage: `name d` removes the variable.
func (v *Variables) DeleteVariable(name string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	_, exists := v.variables[name]
	if exists {
		delete(v.variables, name)
	}
	return exists
}

// ListVariables returns a sorted list of all variable names and their values.
func (v *Variables) ListVariables() []VariableInfo {
	v.mu.RLock()
	defer v.mu.RUnlock()

	var infos []VariableInfo
	for name, value := range v.variables {
		infos = append(infos, VariableInfo{Name: name, Value: value})
	}

	// Sort by name for consistent output
	slices.SortFunc(infos, func(a, b VariableInfo) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return infos
}

// ClearVariables removes all variables from storage.
// Usage: `clear` removes all variables.
func (v *Variables) ClearVariables() {
	v.mu.Lock()
	defer v.mu.Unlock()

	clear(v.variables)
}

// formatVariablesUnsafe returns a list of variable info without acquiring a lock.
// The caller must hold a read lock.
func (v *Variables) formatVariablesUnsafe() string {
	var infos []VariableInfo
	for name, value := range v.variables {
		infos = append(infos, VariableInfo{Name: name, Value: value})
	}

	// Sort by name for consistent output
	slices.SortFunc(infos, func(a, b VariableInfo) int {
		return cmp.Compare(a.Name, b.Name)
	})

	if len(infos) == 0 {
		return "No variables defined"
	}

	var sb strings.Builder
	for i, info := range infos {
		if i > 0 {
			sb.WriteString("\n")
		}
		// Use NumericValue interface for consistent formatting
		num := NewNumber(info.Value, FloatMode)
		sb.WriteString(info.Name)
		sb.WriteString(" = ")
		sb.WriteString(num.String())
	}
	return sb.String()
}

// FormatVariables formats all variables for display.
func (v *Variables) FormatVariables() string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return v.formatVariablesUnsafe()
}

// Count returns the number of defined variables.
func (v *Variables) Count() int {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return len(v.variables)
}

// HasVariable checks if a variable exists.
func (v *Variables) HasVariable(name string) bool {
	v.mu.RLock()
	defer v.mu.RUnlock()

	_, exists := v.variables[name]
	return exists
}

// Save writes the variable store to a file in JSON format.
// The file path should be an absolute path.
// This method acquires a read lock and does not block concurrent readers.
func (v *Variables) Save(path string) error {
	v.mu.RLock()
	defer v.mu.RUnlock()

	// Convert variables to JSON-compatible format
	var infos []VariableInfo
	for name, value := range v.variables {
		infos = append(infos, VariableInfo{Name: name, Value: value})
	}

	// Sort by name for consistent output
	slices.SortFunc(infos, func(a, b VariableInfo) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return saveVariables(path, infos)
}

// Load reads the variable store from a file in JSON format.
// All existing variables are replaced with the loaded values.
// This method is thread-safe.
func (v *Variables) Load(path string) error {
	infos, err := loadVariables(path)
	if err != nil {
		return err
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	// Clear existing variables and load from file
	v.variables = make(map[string]float64)
	for _, info := range infos {
		if isValidVariableName(info.Name) {
			v.variables[info.Name] = info.Value
		}
	}

	return nil
}

// saveVariables saves variable info to a file in JSON format.
// This is a helper function that does not acquire locks.
func saveVariables(path string, infos []VariableInfo) error {
	data, err := json.MarshalIndent(infos, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal variables: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// loadVariables loads variable info from a file in JSON format.
// Returns an empty slice if the file doesn't exist.
// This is a helper function that does not acquire locks.
func loadVariables(path string) ([]VariableInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []VariableInfo{}, nil
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var infos []VariableInfo
	if err := json.Unmarshal(data, &infos); err != nil {
		return nil, fmt.Errorf("failed to unmarshal variables: %w", err)
	}

	return infos, nil
}
