// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import "fmt"

// variables operations

// AssignVariable assigns a value from stack to a variable.
// This is a direct API method that takes the variable name as a parameter
// and pops the value from the stack. It is not the handler for the `=` operator;
// use AssignLeft (for `=:`) or AssignRight (for `=` and `:=`) instead.
func (o *Operations) AssignVariable(stack *Stack, name string) error {
	if name == "" {
		return fmt.Errorf("variable name cannot be empty")
	}

	if stack.Len() < 1 {
		return buildError("=", fmt.Errorf("insufficient operands: need value"))
	}

	val, err := popStack(stack, "=")
	if err != nil {
		return err
	}

	// Convert NumericValue to float64 for variable storage
	valF, err := toFloat64(val, "assigning variable")
	if err != nil {
		return err
	}
	return o.vars.SetVariable(name, valF)
}

// UseVariable pushes a variable's value onto the stack.
// This is a direct API method that takes the variable name as a parameter.
// It is not wired into the operator registry.
func (o *Operations) UseVariable(stack *Stack, name string) error {
	if name == "" {
		return fmt.Errorf("variable name cannot be empty")
	}

	val, exists := o.vars.GetVariable(name)
	if !exists {
		return fmt.Errorf("%w: %s", ErrVariableNotFound, name)
	}

	mode := o.GetMode()
	stack.Push(NewNumber(val, mode))
	return nil
}

// DeleteVariable removes a variable.
// Usage: `name d`
func (o *Operations) DeleteVariable(name string) error {
	if name == "" {
		return fmt.Errorf("variable name cannot be empty")
	}

	deleted := o.vars.DeleteVariable(name)
	if !deleted {
		return fmt.Errorf("%w: %s", ErrVariableNotFound, name)
	}
	return nil
}

// Delete pops a variable name from the stack and deletes that variable.
// Usage: `name d`
func (o *Operations) Delete(stack *Stack) error {
	val, err := popStack(stack, "d")
	if err != nil {
		return err
	}

	var name string
	switch v := val.(type) {
	case *Symbol:
		name = v.Name()
	case *StringNum:
		name = v.String()
	default:
		return fmt.Errorf("delete expects a variable name, got %T", val)
	}
	return o.DeleteVariable(name)
}

// ListVariables lists all variables.
// Usage: `vars`
func (o *Operations) ListVariables() (string, error) {
	return o.vars.FormatVariables(), nil
}

// ClearVariables removes all variables.
// Usage: `clear`
func (o *Operations) ClearVariables() {
	o.vars.ClearVariables()
}

// AssignLeft assigns a value to a variable (for =: operator).
// Stack order: value name =: (value on bottom, name on top).
// This function pops name first (top of stack), then value.
// Usage: `value name =:` (e.g., `5 x =:`)
func (o *Operations) AssignLeft(stack *Stack) error {
	name, err := popStack(stack, "=:")
	if err != nil {
		return err
	}

	val, err := popStack(stack, "=:")
	if err != nil {
		return err
	}

	varName := extractVarName(name)

	valF, err := toFloat64(val, "assigning variable")
	if err != nil {
		return err
	}
	return o.vars.SetVariable(varName, valF)
}

// AssignRight assigns a value to a variable (for = and := operators).
// Stack order: name value = (name on bottom, value on top).
// This function pops value first (top of stack), then name.
// Usage: `name value :=` (e.g., `x 5 :=`)
func (o *Operations) AssignRight(stack *Stack) error {
	val, err := popStack(stack, ":=")
	if err != nil {
		return err
	}

	name, err := popStack(stack, ":=")
	if err != nil {
		return err
	}

	varName := extractVarName(name)

	valF, err := toFloat64(val, "assigning variable")
	if err != nil {
		return err
	}
	return o.vars.SetVariable(varName, valF)
}
