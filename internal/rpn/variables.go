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
func (s *Stack) Clear() {
	s.values = s.values[:0]
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

// VariableStore defines the interface for variable storage operations.
type VariableStore interface {
	SetVariable(name string, value float64) error
	GetVariable(name string) (float64, bool)
	DeleteVariable(name string) bool
	ListVariables() []VariableInfo
	ClearVariables()
	Count() int
	HasVariable(name string) bool
	FormatVariables() string
}

// NewVariables creates and initializes a new Variables instance.
func NewVariables() VariableStore {
	return &Variables{
		variables: make(map[string]float64),
	}
}

// SetVariable assigns a value to a variable name.
// Usage: `name value =` stores value in variable.
func (v *Variables) SetVariable(name string, value float64) error {
	if name == "" {
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

	v.variables = make(map[string]float64)
}

// FormatVariables formats all variables for display.
func (v *Variables) FormatVariables() string {
	infos := v.ListVariables()
	if len(infos) == 0 {
		return "No variables defined"
	}

	var sb strings.Builder
	for i, info := range infos {
		if i > 0 {
			sb.WriteString("\n")
		}
		fmt.Fprintf(&sb, "%s = %.10g", info.Name, info.Value)
	}
	return sb.String()
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
