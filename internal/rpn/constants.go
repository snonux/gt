// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"math"
	"sort"
	"sync"
)

// ConstantsProvider defines the interface for reading constant values.
type ConstantsProvider interface {
	GetConstant(name string) (float64, bool)
	ListConstants() []ConstantInfo
	Count() int
	HasConstant(name string) bool
	SetConstant(name string, value float64) error
	ClearConstants()
	ReloadBuiltInConstants()
}

// ConstantInfo represents a single constant with its name and value.
type ConstantInfo struct {
	Name  string
	Value float64
}

// Constants stores constant name-value pairs for RPN calculations.
// It provides thread-safe access to constant storage.
type Constants struct {
	mu        sync.RWMutex
	constants map[string]float64
}

// NewConstants creates and initializes a new Constants instance with built-in constants.
func NewConstants() *Constants {
	c := &Constants{
		constants: make(map[string]float64),
	}
	c.loadBuiltInConstants()
	return c
}

// loadBuiltInConstants loads the standard mathematical constants.
func (c *Constants) loadBuiltInConstants() {
	// Pi (π) - ratio of a circle's circumference to its diameter
	c.constants["pi"] = math.Pi
	c.constants["π"] = math.Pi

	// Euler's number (e) - base of natural logarithm
	c.constants["e"] = math.E
	c.constants["euler"] = math.E

	// Golden ratio (φ)
	c.constants["phi"] = 1.618033988749895
	c.constants["φ"] = 1.618033988749895

	// Square root of 2
	c.constants["sqrt2"] = 1.414213562373095
	c.constants["√2"] = 1.414213562373095

	// Infinity
	c.constants["inf"] = math.Inf(1)
	c.constants["infinity"] = math.Inf(1)

	// NaN (Not a Number)
	c.constants["nan"] = math.NaN()
}

// SetConstant assigns a value to a constant name.
func (c *Constants) SetConstant(name string, value float64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.constants[name] = value
	return nil
}

// GetConstant retrieves the value of a constant.
// Returns the value and true if found, or 0 and false if not found.
func (c *Constants) GetConstant(name string) (float64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.constants[name]
	return value, exists
}

// ListConstants returns a sorted list of all constant names and their values.
func (c *Constants) ListConstants() []ConstantInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var infos []ConstantInfo
	for name, value := range c.constants {
		infos = append(infos, ConstantInfo{Name: name, Value: value})
	}

	// Sort by name for consistent output
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name < infos[j].Name
	})

	return infos
}

// ClearConstants removes all constants from storage.
// Note: This clears only user-defined constants; built-in constants are preserved.
func (c *Constants) ClearConstants() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Remove only user-defined constants (not built-in ones)
	builtIns := map[string]bool{
		"pi": true, "π": true,
		"e": true, "euler": true,
		"phi": true, "φ": true,
		"sqrt2": true, "√2": true,
		"inf": true, "infinity": true,
		"nan": true,
	}
	for k := range c.constants {
		if !builtIns[k] {
			delete(c.constants, k)
		}
	}
}

// ReloadBuiltInConstants restores all built-in constants.
// This is called internally when ClearConstants is used to ensure
// built-in constants are preserved.
func (c *Constants) ReloadBuiltInConstants() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// First remove only user-defined constants
	builtIns := map[string]bool{
		"pi": true, "π": true,
		"e": true, "euler": true,
		"phi": true, "φ": true,
		"sqrt2": true, "√2": true,
		"inf": true, "infinity": true,
		"nan": true,
	}
	for k := range c.constants {
		if !builtIns[k] {
			delete(c.constants, k)
		}
	}
	// Then reload built-in constants
	c.loadBuiltInConstants()
}

// Count returns the number of defined constants.
func (c *Constants) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.constants)
}

// HasConstant checks if a constant exists.
func (c *Constants) HasConstant(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.constants[name]
	return exists
}
