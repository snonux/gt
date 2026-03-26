// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"sync"
)

// RPN represents the RPN parser and evaluator with state management.
// It is thread-safe for concurrent read operations, but write operations
// on the stack or mode should be synchronized externally or use the provided methods.
type RPN struct {
	mu            sync.RWMutex
	vars          VariableStore
	ops           Operator
	opRegistry    *OperatorRegistry
	assignHandler *assignmentHandler
	maxStack      int
	currentStack  *Stack
	mode          CalculationMode
}

// NewRPN creates a new RPN parser and evaluator with the given variable store.
func NewRPN(vars VariableStore) *RPN {
	ops := NewOperations(vars)
	ops.SetMode(FloatMode) // Set default mode
	return &RPN{
		vars:          vars,
		ops:           ops,
		opRegistry:    NewOperatorRegistry(ops),
		assignHandler: newAssignmentHandler(),
		maxStack:      1000, // Reasonable limit for RPN expressions
		currentStack:  NewStack(),
		mode:          FloatMode, // Default mode
	}
}

// GetMode returns the current calculation mode.
// This method is thread-safe for concurrent reads.
func (r *RPN) GetMode() CalculationMode {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.mode
}

// SetMode sets the calculation mode.
// This method is thread-safe for writes.
func (r *RPN) SetMode(mode CalculationMode) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.mode = mode
	r.ops.SetMode(mode)
}

// GetCurrentStack returns a copy of the current stack for inspection.
// Returns []Number to preserve value types (numbers and booleans).
// This method is thread-safe for concurrent reads.
func (r *RPN) GetCurrentStack() []Number {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.currentStack == nil {
		return nil
	}
	return r.currentStack.Values()
}

// SetCurrentStack sets the current stack from a slice of numbers.
// This is useful for restoring stack state.
// This method is thread-safe for writes.
func (r *RPN) SetCurrentStack(values []Number) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.currentStack = NewStack()
	for _, v := range values {
		r.currentStack.Push(v)
	}
}
