package rpn

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Error variables for external error checking.
var (
	ErrVariableNotFound    = fmt.Errorf("variable not found")
	ErrInvalidVariableName = fmt.Errorf("invalid variable name")
)

// Stack represents a simple float64 stack for RPN calculations.
type Stack struct {
	values []float64
}

// NewStack creates a new empty stack.
func NewStack() *Stack {
	return &Stack{
		values: make([]float64, 0),
	}
}

// Push adds a value to the top of the stack.
func (s *Stack) Push(val float64) {
	s.values = append(s.values, val)
}

// Pop removes and returns the top value from the stack.
// Returns an error if the stack is empty.
func (s *Stack) Pop() (float64, error) {
	if len(s.values) == 0 {
		return 0, fmt.Errorf("stack is empty")
	}

	val := s.values[len(s.values)-1]
	s.values = s.values[:len(s.values)-1]
	return val, nil
}

// Peek returns the top value without removing it.
// Returns an error if the stack is empty.
func (s *Stack) Peek() (float64, error) {
	if len(s.values) == 0 {
		return 0, fmt.Errorf("stack is empty")
	}
	return s.values[len(s.values)-1], nil
}

// Len returns the number of values on the stack.
func (s *Stack) Len() int {
	return len(s.values)
}

// Values returns a copy of all stack values (top-to-bottom order).
func (s *Stack) Values() []float64 {
	vals := make([]float64, len(s.values))
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
}

// NewVariables creates and initializes a new Variables instance.
func NewVariables() *Variables {
	return &Variables{
		variables: make(map[string]float64),
	}
}

// isValidVariableName checks if a variable name is valid.
// Variable names must be non-empty and contain only alphanumeric characters and underscores.
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
