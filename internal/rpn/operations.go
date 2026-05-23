// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"sync"
)

// Helper functions to reduce error handling boilerplate in RPN operations

// popStack pops a value from the stack and returns a wrapped error if insufficient operands.
func popStack(stack *Stack, op string) (StackValue, error) {
	val, err := stack.Pop()
	if err != nil {
		return nil, fmt.Errorf("insufficient operands for %s: %w", op, err)
	}
	return val, nil
}

// popTwo pops two values from the stack for binary operations.
func popTwo(stack *Stack, op string) (StackValue, StackValue, error) {
	b, err := stack.Pop()
	if err != nil {
		return nil, nil, fmt.Errorf("insufficient operands for %s: %w", op, err)
	}

	a, err := stack.Pop()
	if err != nil {
		return nil, nil, fmt.Errorf("insufficient operands for %s: %w", op, err)
	}

	return a, b, nil
}

// toFloat64 converts a StackValue to float64 with proper error wrapping.
// Asserts the value to NumericValue; returns an error for non-numeric types.
func toFloat64(val StackValue, context string) (float64, error) {
	nv, ok := val.(NumericValue)
	if !ok {
		return 0, fmt.Errorf("%s: value %q is not numeric", context, val)
	}
	f, err := nv.Float64()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get float64 value: %w", context, err)
	}
	return f, nil
}

// ensureStackLength checks if the stack has at least min values and returns error if not.
func ensureStackLength(stack *Stack, min int, op string) error {
	if stack.Len() < min {
		return fmt.Errorf("insufficient operands for %s: need at least %d values", op, min)
	}
	return nil
}

// buildError wraps an error with context for the given operator.
func buildError(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}

// popAll pops all values from stack into a slice and reverses them for left-to-right processing.
// Returns values in order from bottom to top of stack (first pushed to last pushed).
func popAll(stack *Stack, op string) ([]StackValue, error) {
	if stack.Len() < 2 {
		return nil, fmt.Errorf("insufficient operands for %s: need at least 2 values", op)
	}

	var values []StackValue
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return nil, fmt.Errorf("%s: failed to pop: %w", op, err)
		}
		values = append(values, val)
	}

	// Reverse to get left-to-right order (first pushed = first in)
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	return values, nil
}

// ArithmeticOperator defines the interface for basic arithmetic operators.
type ArithmeticOperator interface {
	Add(stack *Stack) error
	Subtract(stack *Stack) error
	Multiply(stack *Stack) error
	Divide(stack *Stack) error
	Power(stack *Stack) error
	Modulo(stack *Stack) error
	Log2(stack *Stack) error
	Log10(stack *Stack) error
	Ln(stack *Stack) error
	Convert(stack *Stack) error
}

// BooleanOperator defines the interface for boolean comparison operators.
type BooleanOperator interface {
	GT(stack *Stack) error
	LT(stack *Stack) error
	GTE(stack *Stack) error
	LTE(stack *Stack) error
	EQ(stack *Stack) error
	NEQ(stack *Stack) error
}

// HyperOperator defines the interface for hyper operators.
type HyperOperator interface {
	HyperAdd(stack *Stack) error
	HyperSubtract(stack *Stack) error
	HyperMultiply(stack *Stack) error
	HyperDivide(stack *Stack) error
	HyperPower(stack *Stack) error
	HyperModulo(stack *Stack) error
	HyperLog2(stack *Stack) error
	HyperLog10(stack *Stack) error
	HyperLn(stack *Stack) error
}

// StackOperator defines the interface for stack manipulation operators.
type StackOperator interface {
	Dup(stack *Stack) error
	Swap(stack *Stack) error
	Pop(stack *Stack) error
	Show(stack *Stack) (string, error)
}

// VariableOperator defines the interface for variable operations.
type VariableOperator interface {
	ListVariables() (string, error)
	ClearVariables()
	AssignLeft(stack *Stack) error
	AssignRight(stack *Stack) error
}

// ConstantOperator defines the interface for constant operations.
type ConstantOperator interface {
	ListConstants() (string, error)
	ClearConstants()
}

// PowerIntOperator defines the interface for integer power operations (**) using binary exponentiation.
type PowerIntOperator interface {
	FastPower(stack *Stack) error
}

// Operator is the combined interface for all operator implementations.
// This allows RPN to depend on an abstraction instead of the concrete Operations type.
//
// Design note: Operator intentionally mixes behavioral methods (arithmetic, stack,
// boolean ops) with configuration methods (SetMode, SetPrefixMode, GetPrefixMode)
// and metric command handlers. Per ISP this could be split, but RPN is the sole
// client and splitting would add indirection without practical benefit. The
// concrete *Operations type satisfies this interface exclusively.
type Operator interface {
	ArithmeticOperator
	BooleanOperator
	HyperOperator
	StackOperator
	VariableOperator
	ConstantOperator
	PowerIntOperator
	// SetMode sets the calculation mode for number formatting
	SetMode(CalculationMode)
	// SetPrefixMode sets the prefix mode for data size calculations
	SetPrefixMode(PrefixMode)
	// GetPrefixMode returns the current prefix mode
	GetPrefixMode() PrefixMode
	// Metric command handlers
	MetricShow(stack *Stack) (string, error)
	MetricList(stack *Stack) (string, error)
	MetricCategory(stack *Stack, categoryName string) (string, error)
	MetricCompatible(stack *Stack) (string, error)
	// Custom metric commands
	CustomList(stack *Stack) (string, error)
	CustomDefine(name string, factor float64, category string) error
	CustomUndefine(name string) error
}

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

// OperatorHandler represents a function that handles an operator.
// Returns (result string, handled bool, error error).
// result is non-empty only for commands that return immediately (like show, vars).
type OperatorHandler func(stack *Stack) (result string, handled bool, err error)

// OperatorRegistry maintains a registry of operators.
type OperatorRegistry struct {
	standardOperators map[string]OperatorHandler
	hyperOperators    map[string]OperatorHandler
}

// NewOperatorRegistry creates a new operator registry and registers all operators.
func NewOperatorRegistry(op Operator) *OperatorRegistry {
	registry := &OperatorRegistry{
		standardOperators: make(map[string]OperatorHandler),
		hyperOperators:    make(map[string]OperatorHandler),
	}

	// Register standard operators
	registry.registerStandardOperator("+", func(stack *Stack) error { return op.Add(stack) })
	registry.registerStandardOperator("-", func(stack *Stack) error { return op.Subtract(stack) })
	registry.registerStandardOperator("*", func(stack *Stack) error { return op.Multiply(stack) })
	registry.registerStandardOperator("/", func(stack *Stack) error { return op.Divide(stack) })
	registry.registerStandardOperator("^", func(stack *Stack) error { return op.Power(stack) })
	registry.registerStandardOperator("**", func(stack *Stack) error { return op.FastPower(stack) })
	registry.registerStandardOperator("%", func(stack *Stack) error { return op.Modulo(stack) })
	registry.registerStandardOperator("lg", func(stack *Stack) error { return op.Log2(stack) })
	registry.registerStandardOperator("log", func(stack *Stack) error { return op.Log10(stack) })
	registry.registerStandardOperator("ln", func(stack *Stack) error { return op.Ln(stack) })
	registry.registerStandardOperator("gt", func(stack *Stack) error { return op.GT(stack) })
	registry.registerStandardOperator("lt", func(stack *Stack) error { return op.LT(stack) })
	registry.registerStandardOperator("<", func(stack *Stack) error { return op.LT(stack) })
	registry.registerStandardOperator(">", func(stack *Stack) error { return op.GT(stack) })
	registry.registerStandardOperator("gte", func(stack *Stack) error { return op.GTE(stack) })
	registry.registerStandardOperator(">=", func(stack *Stack) error { return op.GTE(stack) })
	registry.registerStandardOperator("lte", func(stack *Stack) error { return op.LTE(stack) })
	registry.registerStandardOperator("<=", func(stack *Stack) error { return op.LTE(stack) })
	registry.registerStandardOperator("eq", func(stack *Stack) error { return op.EQ(stack) })
	registry.registerStandardOperator("==", func(stack *Stack) error { return op.EQ(stack) })
	registry.registerStandardOperator("neq", func(stack *Stack) error { return op.NEQ(stack) })
	registry.registerStandardOperator("!=", func(stack *Stack) error { return op.NEQ(stack) })
	registry.registerStandardOperator("=", func(stack *Stack) error { return op.AssignLeft(stack) })
	registry.registerStandardOperator(":=", func(stack *Stack) error { return op.AssignRight(stack) })
	registry.registerStandardOperator("=:", func(stack *Stack) error { return op.AssignLeft(stack) })
	registry.registerStandardOperator("convert", func(stack *Stack) error { return op.Convert(stack) })
	registry.registerStandardOperator("dup", func(stack *Stack) error { return op.Dup(stack) })
	registry.registerStandardOperator("swap", func(stack *Stack) error { return op.Swap(stack) })
	registry.registerStandardOperator("pop", func(stack *Stack) error { return op.Pop(stack) })
	registry.registerStandardOperator("d", func(stack *Stack) error {
		return fmt.Errorf("'d' command not supported as standalone token")
	})

	// Commands that return immediately
	registry.registerCommandOperator("show", func(stack *Stack) (string, error) { return op.Show(stack) })
	registry.registerCommandOperator("showstack", func(stack *Stack) (string, error) { return op.Show(stack) })
	registry.registerCommandOperator("print", func(stack *Stack) (string, error) { return op.Show(stack) })
	registry.registerCommandOperator("vars", func(stack *Stack) (string, error) { return op.ListVariables() })
	registry.registerCommandOperator("constants", func(stack *Stack) (string, error) { return op.ListConstants() })
	registry.registerCommandOperator("clear", func(stack *Stack) (string, error) { op.ClearVariables(); return "All variables cleared", nil })
	registry.registerCommandOperator("clearconstants", func(stack *Stack) (string, error) { op.ClearConstants(); return "All constants cleared", nil })

	// Register hyper operators
	registry.registerHyperOperator("[+]", func(stack *Stack) error { return op.HyperAdd(stack) })
	registry.registerHyperOperator("[-]", func(stack *Stack) error { return op.HyperSubtract(stack) })
	registry.registerHyperOperator("[*]", func(stack *Stack) error { return op.HyperMultiply(stack) })
	registry.registerHyperOperator("[/]", func(stack *Stack) error { return op.HyperDivide(stack) })
	registry.registerHyperOperator("[^]", func(stack *Stack) error { return op.HyperPower(stack) })
	registry.registerHyperOperator("[%]", func(stack *Stack) error { return op.HyperModulo(stack) })
	registry.registerHyperOperator("[lg]", func(stack *Stack) error { return op.HyperLog2(stack) })
	registry.registerHyperOperator("[log]", func(stack *Stack) error { return op.HyperLog10(stack) })
	registry.registerHyperOperator("[ln]", func(stack *Stack) error { return op.HyperLn(stack) })

	return registry
}

// registerStandardOperator registers a standard operator that returns empty result.
func (r *OperatorRegistry) registerStandardOperator(name string, handler func(*Stack) error) {
	r.standardOperators[name] = func(stack *Stack) (string, bool, error) {
		if err := handler(stack); err != nil {
			return "", false, fmt.Errorf("%s: %w", name, err)
		}
		return "", true, nil
	}
}

// registerCommandOperator registers a command operator that returns a result immediately.
func (r *OperatorRegistry) registerCommandOperator(name string, handler func(*Stack) (string, error)) {
	r.standardOperators[name] = func(stack *Stack) (string, bool, error) {
		result, err := handler(stack)
		if err != nil {
			return "", false, fmt.Errorf("%s: %w", name, err)
		}
		return result, true, nil
	}
}

// registerHyperOperator registers a hyper operator.
func (r *OperatorRegistry) registerHyperOperator(name string, handler func(*Stack) error) {
	r.hyperOperators[name] = func(stack *Stack) (string, bool, error) {
		if err := handler(stack); err != nil {
			return "", false, fmt.Errorf("%s: %w", name, err)
		}
		return "", true, nil
	}
}

// HandleStandardOperator handles a standard operator.
// Returns (result string, handled bool, error error).
func (r *OperatorRegistry) HandleStandardOperator(stack *Stack, token string) (string, bool, error) {
	if handler, exists := r.standardOperators[token]; exists {
		return handler(stack)
	}
	return "", false, fmt.Errorf("unknown token '%s'", token)
}

// HandleHyperOperator handles a hyper operator.
// Returns (result string, handled bool, error error).
func (r *OperatorRegistry) HandleHyperOperator(stack *Stack, token string) (string, bool, error) {
	if handler, exists := r.hyperOperators[token]; exists {
		return handler(stack)
	}
	return "", false, fmt.Errorf("unknown token '%s'", token)
}

// IsStandardOperator checks if a token is a standard operator.
func (r *OperatorRegistry) IsStandardOperator(token string) bool {
	_, exists := r.standardOperators[token]
	return exists
}

// IsHyperOperator checks if a token is a hyper operator.
func (r *OperatorRegistry) IsHyperOperator(token string) bool {
	_, exists := r.hyperOperators[token]
	return exists
}
