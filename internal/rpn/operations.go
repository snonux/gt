// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"sync"
)

// Operations provides operator implementations and stack manipulation.
type Operations struct {
	vars           VariableStore
	consts         ConstantsProvider
	mode           CalculationMode
	prefixMode     PrefixMode
	metricRegistry *MetricRegistry
	mu             sync.RWMutex
}

// Ensure Operations implements all operator sub-interfaces at compile time.
var (
	_ ArithmeticOperator = (*Operations)(nil)
	_ BooleanOperator    = (*Operations)(nil)
	_ HyperOperator      = (*Operations)(nil)
	_ StackOperator      = (*Operations)(nil)
	_ VariableOperator   = (*Operations)(nil)
	_ ConstantOperator   = (*Operations)(nil)
	_ PowerIntOperator   = (*Operations)(nil)
)

// NewOperations creates a new Operations instance with the given variable store.
// Does not create a ConstantsProvider internally; caller must use SetConstants.
// If no registry is provided, defaults to the global MetricRegistry.
func NewOperations(vars VariableStore, reg ...*MetricRegistry) *Operations {
	r := GetMetricRegistry()
	if len(reg) > 0 && reg[0] != nil {
		r = reg[0]
	}
	return &Operations{
		vars:           vars,
		mode:           FloatMode, // default
		prefixMode:     SI,        // default
		metricRegistry: r,
	}
}

// SetConstants sets the constants provider for the Operations instance.
// This allows sharing a single ConstantsProvider between RPN and Operations.
func (o *Operations) SetConstants(c ConstantsProvider) {
	o.consts = c
}

// SetMode sets the calculation mode for the Operations instance.
// This method is thread-safe for writes.
func (o *Operations) SetMode(mode CalculationMode) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.mode = mode
}

// GetMode returns the current calculation mode.
// This method is thread-safe for reads.
func (o *Operations) GetMode() CalculationMode {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.mode
}

// GetPrefixMode returns the current prefix mode.
// This method is thread-safe for reads.
func (o *Operations) GetPrefixMode() PrefixMode {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.prefixMode
}

// SetPrefixMode sets the prefix mode for data size calculations.
// This method is thread-safe for writes.
func (o *Operations) SetPrefixMode(mode PrefixMode) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.prefixMode = mode
}

// MetricRegistry returns the metric registry used by this Operations instance.
func (o *Operations) MetricRegistry() *MetricRegistry {
	return o.metricRegistry
}
