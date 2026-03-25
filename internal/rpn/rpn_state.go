// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

// RPN represents the RPN parser and evaluator with state management.
type RPN struct {
	vars         VariableStore
	ops          Operator
	opRegistry   *OperatorRegistry
	maxStack     int
	currentStack *Stack
	mode         CalculationMode
}

// NewRPN creates a new RPN parser and evaluator with the given variable store.
func NewRPN(vars VariableStore) *RPN {
	ops := NewOperations(vars)
	ops.SetMode(FloatMode) // Set default mode
	return &RPN{
		vars:         vars,
		ops:          ops,
		opRegistry:   NewOperatorRegistry(ops),
		maxStack:     1000, // Reasonable limit for RPN expressions
		currentStack: NewStack(),
		mode:         FloatMode, // Default mode
	}
}

// GetMode returns the current calculation mode.
func (r *RPN) GetMode() CalculationMode {
	return r.mode
}

// SetMode sets the calculation mode.
func (r *RPN) SetMode(mode CalculationMode) {
	r.mode = mode
	r.ops.SetMode(mode)
}

// GetCurrentStack returns a copy of the current stack for inspection.
// Returns []Value to preserve value types (numbers and booleans).
func (r *RPN) GetCurrentStack() []Value {
	if r.currentStack == nil {
		return nil
	}
	return r.currentStack.Values()
}

// SetCurrentStack sets the current stack from a slice of values.
// This is useful for restoring stack state.
func (r *RPN) SetCurrentStack(values []Value) {
	r.currentStack = NewStack()
	for _, v := range values {
		r.currentStack.Push(v)
	}
}
