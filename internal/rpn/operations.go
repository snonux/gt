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

// Ensure Operations implements Operator at compile time.
// This is an explicit interface satisfaction check that will fail to compile
// if Operations doesn't implement all methods required by the Operator interface.
var _ Operator = (*Operations)(nil)

// NewOperations creates a new Operations instance with the given variable store.
// If no registry is provided, defaults to the global MetricRegistry.
func NewOperations(vars VariableStore, reg ...*MetricRegistry) *Operations {
	consts := NewConstants()
	r := GetMetricRegistry()
	if len(reg) > 0 && reg[0] != nil {
		r = reg[0]
	}
	return &Operations{
		vars:           vars,
		consts:         consts,
		mode:           FloatMode, // default
		prefixMode:     SI,        // default
		metricRegistry: r,
	}
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
