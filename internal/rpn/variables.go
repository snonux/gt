// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

// Error variables for external error checking.
var (
	ErrVariableNotFound    = fmt.Errorf("variable not found")
	ErrInvalidVariableName = fmt.Errorf("invalid variable name")
)

// Value represents a variant type that can hold either a number (float64) or a boolean.
//
// When used in arithmetic operations, boolean values are automatically coerced:
//   - true → 1
//   - false → 0
//
// This allows boolean results from comparison operations to be used directly in
// arithmetic expressions (e.g., "5 3 == 1 +" where "5 3 ==" produces false=0,
// and "0 + 1" produces 1).
type Value struct {
	isBool  bool
	boolVal bool
	numVal  float64
}

// NewNumberValue creates a new Value containing a float64 number.
func NewNumberValue(n float64) Value {
	return Value{isBool: false, numVal: n}
}

// NewBoolValue creates a new Value containing a boolean.
func NewBoolValue(b bool) Value {
	return Value{isBool: true, boolVal: b}
}

// IsBool returns true if the value is a boolean.
func (v Value) IsBool() bool {
	return v.isBool
}

// IsNumber returns true if the value is a number.
func (v Value) IsNumber() bool {
	return !v.isBool
}

// Bool returns the boolean value, or false if the value is not a boolean.
func (v Value) Bool() bool {
	return v.boolVal
}

// Float64 returns the float64 value.
// If the value is a boolean, true returns 1 and false returns 0.
// If the value is a number, it returns the numeric value directly.
func (v Value) Float64() float64 {
	if v.isBool {
		if v.boolVal {
			return 1
		}
		return 0
	}
	return v.numVal
}

// Number returns the float64 value (deprecated, use Float64 instead).
// If the value is a boolean, this returns 0 (the numeric value is not used for booleans).
func (v Value) Number() float64 {
	return v.numVal
}

// String returns the string representation of the value.
// For booleans, it returns "true" or "false".
// For numbers, it returns the formatted float64 value.
func (v Value) String() string {
	if v.isBool {
		if v.boolVal {
			return "true"
		}
		return "false"
	}
	return fmt.Sprintf("%.10g", v.numVal)
}

// Stack represents a stack for RPN calculations using the Number interface.
// It can hold both number and boolean values through the Number interface.
type Stack struct {
	values []Number
}

// NewStack creates a new empty stack.
func NewStack() *Stack {
	return &Stack{
		values: make([]Number, 0),
	}
}

// Push adds a value to the top of the stack.
func (s *Stack) Push(val Number) {
	s.values = append(s.values, val)
}

// Pop removes and returns the top value from the stack.
// Returns an error if the stack is empty.
func (s *Stack) Pop() (Number, error) {
	if len(s.values) == 0 {
		return nil, fmt.Errorf("stack is empty")
	}

	val := s.values[len(s.values)-1]
	s.values = s.values[:len(s.values)-1]
	return val, nil
}

// Peek returns the top value without removing it.
// Returns an error if the stack is empty.
func (s *Stack) Peek() (Number, error) {
	if len(s.values) == 0 {
		return nil, fmt.Errorf("stack is empty")
	}
	return s.values[len(s.values)-1], nil
}

// Len returns the number of values on the stack.
func (s *Stack) Len() int {
	return len(s.values)
}

// Values returns a copy of all stack values (top-to-bottom order).
func (s *Stack) Values() []Number {
	vals := make([]Number, len(s.values))
	copy(vals, s.values)
	return vals
}

// Clear removes all values from the stack.
// Note: This resets the slice to nil, releasing the underlying memory.
func (s *Stack) Clear() {
	s.values = nil
}

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

// VariableStore combines VariableReader and VariableWriter for full variable storage access.
type VariableStore interface {
	VariableReader
	VariableWriter
	// Save writes the variable store to a file in JSON format.
	Save(path string) error
	// Load reads the variable store from a file in JSON format.
	// Existing variables are overwritten; new variables are added.
	Load(path string) error
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
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name < infos[j].Name
	})

	return infos
}

// ClearVariables removes all variables from storage.
// Usage: `clear` removes all variables.
func (v *Variables) ClearVariables() {
	v.mu.Lock()
	defer v.mu.Unlock()

	for k := range v.variables {
		delete(v.variables, k)
	}
}

// formatVariablesUnsafe returns a list of variable info without acquiring a lock.
// The caller must hold a read lock.
func (v *Variables) formatVariablesUnsafe() string {
	var infos []VariableInfo
	for name, value := range v.variables {
		infos = append(infos, VariableInfo{Name: name, Value: value})
	}

	// Sort by name for consistent output
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name < infos[j].Name
	})

	if len(infos) == 0 {
		return "No variables defined"
	}

	var sb strings.Builder
	for i, info := range infos {
		if i > 0 {
			sb.WriteString("\n")
		}
		// Use Number interface for consistent formatting
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
// This method is thread-safe for concurrent writes.
func (v *Variables) Save(path string) error {
	v.mu.RLock()
	defer v.mu.RUnlock()

	// Convert variables to JSON-compatible format
	var infos []VariableInfo
	for name, value := range v.variables {
		infos = append(infos, VariableInfo{Name: name, Value: value})
	}

	// Sort by name for consistent output
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name < infos[j].Name
	})

	return saveVariables(path, infos)
}

// Load reads the variable store from a file in JSON format.
// Existing variables are overwritten; new variables are added.
// This method is thread-safe for concurrent reads.
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
