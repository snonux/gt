// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"cmp"
	"maps"
	"math"
	"slices"
	"sync"
)

// ConstantsReader defines the interface for reading constant values.
type ConstantsReader interface {
	GetConstant(name string) (float64, bool)
	ListConstants() []ConstantInfo
	Count() int
	HasConstant(name string) bool
}

// ConstantsWriter defines the interface for writing constant values.
type ConstantsWriter interface {
	SetConstant(name string, value float64) error
}

// ConstantsAdmin defines the interface for administrative constant operations.
type ConstantsAdmin interface {
	ClearConstants()
	ReloadBuiltInConstants()
}

// ConstantsProvider combines ConstantsReader, ConstantsWriter, and ConstantsAdmin
// for full constant storage access.
type ConstantsProvider interface {
	ConstantsReader
	ConstantsWriter
	ConstantsAdmin
}

// Ensure Constants implements all constant sub-interfaces at compile time.
var (
	_ ConstantsReader   = (*Constants)(nil)
	_ ConstantsWriter   = (*Constants)(nil)
	_ ConstantsAdmin    = (*Constants)(nil)
	_ ConstantsProvider = (*Constants)(nil)
)

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

// loadBuiltInConstants loads the standard mathematical constants and returns the set of keys it set.
func (c *Constants) loadBuiltInConstants() map[string]struct{} {
	builtIns := make(map[string]struct{})
	// Pi (π) - ratio of a circle's circumference to its diameter
	c.constants["pi"] = math.Pi
	builtIns["pi"] = struct{}{}
	c.constants["π"] = math.Pi
	builtIns["π"] = struct{}{}

	// Euler's number (e) - base of natural logarithm
	c.constants["e"] = math.E
	builtIns["e"] = struct{}{}
	c.constants["euler"] = math.E
	builtIns["euler"] = struct{}{}

	// Golden ratio (φ)
	c.constants["phi"] = 1.618033988749895
	builtIns["phi"] = struct{}{}
	c.constants["φ"] = 1.618033988749895
	builtIns["φ"] = struct{}{}

	// Square root of 2
	c.constants["sqrt2"] = 1.414213562373095
	builtIns["sqrt2"] = struct{}{}
	c.constants["√2"] = 1.414213562373095
	builtIns["√2"] = struct{}{}

	// Square root of 3
	c.constants["sqrt3"] = 1.732050807568877
	builtIns["sqrt3"] = struct{}{}
	c.constants["√3"] = 1.732050807568877
	builtIns["√3"] = struct{}{}

	// Square root of 5
	c.constants["sqrt5"] = 2.23606797749979
	builtIns["sqrt5"] = struct{}{}
	c.constants["√5"] = 2.23606797749979
	builtIns["√5"] = struct{}{}

	// Natural logarithm of 2
	c.constants["ln2"] = 0.693147180559945
	builtIns["ln2"] = struct{}{}
	c.constants["log2"] = 0.693147180559945
	builtIns["log2"] = struct{}{}

	// Natural logarithm of 10
	c.constants["ln10"] = 2.302585092994046
	builtIns["ln10"] = struct{}{}
	c.constants["log10"] = 2.302585092994046
	builtIns["log10"] = struct{}{}

	// Logarithm of e base 10
	c.constants["log_e"] = 0.434294481903252
	builtIns["log_e"] = struct{}{}
	c.constants["log_e10"] = 0.434294481903252
	builtIns["log_e10"] = struct{}{}

	// Tau (2π) - circle constant
	c.constants["tau"] = 2 * math.Pi
	builtIns["tau"] = struct{}{}
	c.constants["τ"] = 2 * math.Pi
	builtIns["τ"] = struct{}{}

	// Fraction 1/π
	c.constants["1/π"] = 1 / math.Pi
	builtIns["1/π"] = struct{}{}
	c.constants["inv_pi"] = 1 / math.Pi
	builtIns["inv_pi"] = struct{}{}

	// Fraction 1/e
	c.constants["1/e"] = 1 / math.E
	builtIns["1/e"] = struct{}{}
	c.constants["inv_e"] = 1 / math.E
	builtIns["inv_e"] = struct{}{}

	// Infinity
	c.constants["inf"] = math.Inf(1)
	builtIns["inf"] = struct{}{}
	c.constants["infinity"] = math.Inf(1)
	builtIns["infinity"] = struct{}{}

	// Negative infinity
	c.constants["-inf"] = math.Inf(-1)
	builtIns["-inf"] = struct{}{}
	c.constants["-infinity"] = math.Inf(-1)
	builtIns["-infinity"] = struct{}{}

	// NaN (Not a Number)
	c.constants["nan"] = math.NaN()
	builtIns["nan"] = struct{}{}

	return builtIns
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
	slices.SortFunc(infos, func(a, b ConstantInfo) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return infos
}

// ClearConstants removes all constants and reloads the built-in defaults.
// Note: This resets all constants, including any user-defined ones.
func (c *Constants) ClearConstants() {
	c.mu.Lock()
	defer c.mu.Unlock()

	clear(c.constants)
	c.loadBuiltInConstants()
}

// ReloadBuiltInConstants restores all built-in constants to their default values.
// It also removes any user-defined constants, effectively resetting the store.
func (c *Constants) ReloadBuiltInConstants() {
	c.mu.Lock()
	defer c.mu.Unlock()

	builtIns := c.loadBuiltInConstants()

	maps.DeleteFunc(c.constants, func(k string, _ float64) bool {
		_, ok := builtIns[k]
		return !ok
	})
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
