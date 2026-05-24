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
	mu             sync.RWMutex
	vars           VariableStore
	consts         ConstantsProvider
	ops            Operator
	opRegistry     *OperatorRegistry
	assignHandler  *assignmentHandler
	maxStack       int
	currentStack   *Stack
}

// NewRPN creates a new RPN parser and evaluator with the given variable store.
// If no registry is provided, defaults to the global MetricRegistry.
func NewRPN(vars VariableStore, reg ...*MetricRegistry) *RPN {
	consts := NewConstants()
	ops := NewOperations(vars, reg...)
	ops.SetMode(FloatMode) // Set default mode
	ops.SetConstants(consts) // Share the same ConstantsProvider
	return &RPN{
		vars:          vars,
		consts:        consts,
		ops:           ops,
		opRegistry:    NewOperatorRegistry(ops),
		assignHandler: newAssignmentHandler(),
		maxStack:      1000, // Reasonable limit for RPN expressions
		currentStack:  NewStack(),
	}
}

// GetConstants returns the constants provider.
// This method is thread-safe for concurrent reads.
func (r *RPN) GetConstants() ConstantsProvider {
	return r.consts
}

// GetMode returns the current calculation mode.
// This method is thread-safe for reads.
func (r *RPN) GetMode() CalculationMode {
	return r.ops.GetMode()
}

// SetMode sets the calculation mode.
// This method is thread-safe for writes.
func (r *RPN) SetMode(mode CalculationMode) {
	r.ops.SetMode(mode)
}

// GetCurrentStack returns a copy of the current stack for inspection.
// Returns []StackValue to preserve value types (numbers, booleans, strings, symbols).
// This method is thread-safe for concurrent reads.
func (r *RPN) GetCurrentStack() []StackValue {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.currentStack == nil {
		return nil
	}
	return r.currentStack.Values()
}

// SetCurrentStack sets the current stack from a slice of StackValues.
// This is useful for restoring stack state.
// This method is thread-safe for writes.
func (r *RPN) SetCurrentStack(values []StackValue) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.currentStack = NewStack()
	for _, v := range values {
		r.currentStack.Push(v)
	}
}

// Stack returns the current stack as a slice of StackValues.
// This is a convenience wrapper around GetCurrentStack().
// Returns an empty slice if the stack has no values.
func (r *RPN) Stack() []StackValue {
	return r.GetCurrentStack()
}

// SetPrefixMode sets the prefix mode (SI or IEC).
// Delegates to the Operations instance.
// This method is thread-safe for writes.
func (r *RPN) SetPrefixMode(mode PrefixMode) {
	r.ops.SetPrefixMode(mode)
}

// GetPrefixMode returns the current prefix mode.
// Delegates to the Operations instance.
// This method is thread-safe for reads.
func (r *RPN) GetPrefixMode() PrefixMode {
	return r.ops.GetPrefixMode()
}

// IsStandardOperator checks if a token is a registered standard operator.
// Delegates to the operator registry.
func (r *RPN) IsStandardOperator(token string) bool {
	return r.opRegistry.IsStandardOperator(token)
}

// IsHyperOperator checks if a token is a registered hyper operator.
// Delegates to the operator registry.
func (r *RPN) IsHyperOperator(token string) bool {
	return r.opRegistry.IsHyperOperator(token)
}
