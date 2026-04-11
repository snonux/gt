// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"codeberg.org/snonux/gt/internal/rpn"
)

// Calculator defines the interface for RPN calculation operations.
// This interface abstracts the RPN engine to decouple the REPL from specific
// RPN implementation details.
type Calculator interface {
	// ParseAndEvaluate parses and evaluates an RPN expression.
	// Returns the result string and any error encountered.
	ParseAndEvaluate(input string) (string, error)

	// EvalOperator evaluates a single RPN operator on the current stack.
	// Returns the result string and any error encountered.
	EvalOperator(op string) (string, error)

	// GetMode returns the current calculation mode.
	GetMode() rpn.CalculationMode

	// SetMode sets the calculation mode.
	SetMode(mode rpn.CalculationMode)
}

// RPNCalculator is an adapter that wraps an rpn.RPN instance to implement Calculator.
type RPNCalculator struct {
	rpnCalc *rpn.RPN
}

// NewRPNCalculator creates a new RPNCalculator that wraps the given RPN instance.
func NewRPNCalculator(rpnCalc *rpn.RPN) *RPNCalculator {
	return &RPNCalculator{rpnCalc: rpnCalc}
}

// ParseAndEvaluate parses and evaluates an RPN expression.
// Implements Calculator interface.
func (c *RPNCalculator) ParseAndEvaluate(input string) (string, error) {
	return c.rpnCalc.ParseAndEvaluate(input)
}

// EvalOperator evaluates a single RPN operator on the current stack.
// Implements Calculator interface.
func (c *RPNCalculator) EvalOperator(op string) (string, error) {
	return c.rpnCalc.EvalOperator(op)
}

// GetMode returns the current calculation mode.
// Implements Calculator interface.
func (c *RPNCalculator) GetMode() rpn.CalculationMode {
	return c.rpnCalc.GetMode()
}

// SetMode sets the calculation mode.
// Implements Calculator interface.
func (c *RPNCalculator) SetMode(mode rpn.CalculationMode) {
	c.rpnCalc.SetMode(mode)
}
