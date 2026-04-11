// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
)

// VariableOperations provides variable management operator implementations.
type VariableOperations struct {
	vars VariableStore
}

// NewVariableOperations creates a new VariableOperations instance.
func NewVariableOperations(vars VariableStore) *VariableOperations {
	return &VariableOperations{vars: vars}
}

// AssignVariable assigns a value from the stack to a variable.
// Usage: `name value =`
func (o *VariableOperations) AssignVariable(stack *Stack, name string) error {
	val, err := stack.Pop()
	if err != nil {
		return err
	}

	// Convert Number to float64 for variable storage
	valF, err := val.Float64()
	if err != nil {
		return fmt.Errorf("failed to get float64 value for variable: %w", err)
	}
	return o.vars.SetVariable(name, valF)
}

// UseVariable pushes a variable's value onto the stack.
// Usage: `varname` (pushes stored value)
func (o *VariableOperations) UseVariable(stack *Stack, name string) error {
	if name == "" {
		return fmt.Errorf("variable name cannot be empty")
	}

	val, exists := o.vars.GetVariable(name)
	if !exists {
		return fmt.Errorf("%w: %s", ErrVariableNotFound, name)
	}

	stack.Push(NewNumber(val, FloatMode))
	return nil
}

// DeleteVariable removes a variable.
// Usage: `name d`
func (o *VariableOperations) DeleteVariable(name string) error {
	if name == "" {
		return fmt.Errorf("variable name cannot be empty")
	}

	deleted := o.vars.DeleteVariable(name)
	if !deleted {
		return fmt.Errorf("%w: %s", ErrVariableNotFound, name)
	}
	return nil
}

// ListVariables returns a string listing all variables.
// Usage: `vars`
func (o *VariableOperations) ListVariables() (string, error) {
	return o.vars.FormatVariables(), nil
}

// ClearVariables removes all variables.
// Usage: `clear`
func (o *VariableOperations) ClearVariables() {
	o.vars.ClearVariables()
}
