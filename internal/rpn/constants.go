// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"cmp"

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

// builtInConstants is the immutable set of standard mathematical constants.
var builtInConstants = map[string]float64{
	// Pi (π) - ratio of a circle's circumference to its diameter
	"pi":    math.Pi,
	"π":     math.Pi,
	// Euler's number (e) - base of natural logarithm
	"e":     math.E,
	"euler": math.E,
	// Golden ratio (φ)
	"phi": 1.618033988749895,
	"φ":   1.618033988749895,
	// Square root of 2
	"sqrt2": 1.414213562373095,
	"√2":    1.414213562373095,
	// Square root of 3
	"sqrt3": 1.732050807568877,
	"√3":    1.732050807568877,
	// Square root of 5
	"sqrt5": 2.23606797749979,
	"√5":    2.23606797749979,
	// Natural logarithm of 2
	"ln2":  0.693147180559945,
	"log2": 0.693147180559945,
	// Natural logarithm of 10
	"ln10":  2.302585092994046,
	"log10": 2.302585092994046,
	// Logarithm of e base 10
	"log_e":   0.434294481903252,
	"log_e10": 0.434294481903252,
	// Tau (2π) - circle constant
	"tau": 2 * math.Pi,
	"τ":   2 * math.Pi,
	// Fraction 1/π
	"1/π":    1 / math.Pi,
	"inv_pi": 1 / math.Pi,
	// Fraction 1/e
	"1/e":    1 / math.E,
	"inv_e":  1 / math.E,
	// Infinity
	"inf":      math.Inf(1),
	"infinity": math.Inf(1),
	// Negative infinity
	"-inf":      math.Inf(-1),
	"-infinity": math.Inf(-1),
	// NaN (Not a Number)
	"nan": math.NaN(),
}

// NewConstants creates and initializes a new Constants instance with built-in constants.
func NewConstants() *Constants {
	c := &Constants{
		constants: make(map[string]float64),
	}
	c.loadBuiltInConstants()
	return c
}

// loadBuiltInConstants copies the built-in constants into the constants map.
func (c *Constants) loadBuiltInConstants() {
	for name, value := range builtInConstants {
		c.constants[name] = value
	}
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

	clear(c.constants)
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
