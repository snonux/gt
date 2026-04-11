// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"math"
	"strings"
	"sync"
)

// Helper functions to reduce error handling boilerplate in RPN operations

// popStack pops a value from the stack and returns a wrapped error if insufficient operands.
func popStack(stack *Stack, op string) (Number, error) {
	val, err := stack.Pop()
	if err != nil {
		return nil, fmt.Errorf("insufficient operands for %s: %w", op, err)
	}
	return val, nil
}

// popTwo pops two values from the stack for binary operations.
func popTwo(stack *Stack, op string) (Number, Number, error) {
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

// toFloat64 converts a Number to float64 with proper error wrapping.
func toFloat64(val Number, context string) (float64, error) {
	f, err := val.Float64()
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
func popAll(stack *Stack, op string) ([]Number, error) {
	if stack.Len() < 2 {
		return nil, fmt.Errorf("insufficient operands for %s: need at least 2 values", op)
	}

	var values []Number
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
	// AssignLeft assigns a value to a variable (for := operator)
	AssignLeft(stack *Stack) error
	// AssignRight assigns a value to a variable (for =: operator)
	AssignRight(stack *Stack) error
}

// Operations provides operator implementations and stack manipulation.
type Operations struct {
	vars    VariableStore
	consts  ConstantsProvider
	mode    CalculationMode
	mu      sync.RWMutex
}

// Ensure Operations implements Operator at compile time.
// This is an explicit interface satisfaction check that will fail to compile
// if Operations doesn't implement all methods required by the Operator interface.
var _ Operator = (*Operations)(nil)

// NewOperations creates a new Operations instance with the given variable store.
func NewOperations(vars VariableStore) *Operations {
	consts := NewConstants()
	return &Operations{
		vars:   vars,
		consts: consts,
		mode:   FloatMode, // default
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

// OperatorHandler represents a function that handles an operator.
// Returns (result string, handled bool, error error).
// result is non-empty only for commands that return immediately (like show, vars).
// handled indicates if the token was recognized.
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

// arithmetic operators

// Add pops two values from stack, adds them, and pushes result.
func (o *Operations) Add(stack *Stack) error {
	a, b, err := popTwo(stack, "+")
	if err != nil {
		return err
	}

	// Use the Number interface for arithmetic
	result, err := a.Add(b)
	if err != nil {
		return buildError("addition", err)
	}
	stack.Push(result)
	return nil
}

// Subtract pops two values from stack, subtracts (a - b), and pushes result.
func (o *Operations) Subtract(stack *Stack) error {
	a, b, err := popTwo(stack, "-")
	if err != nil {
		return err
	}

	result, err := a.Sub(b)
	if err != nil {
		return buildError("subtraction", err)
	}
	stack.Push(result)
	return nil
}

// Multiply pops two values from stack, multiplies them, and pushes result.
func (o *Operations) Multiply(stack *Stack) error {
	a, b, err := popTwo(stack, "*")
	if err != nil {
		return err
	}

	result, err := a.Mul(b)
	if err != nil {
		return buildError("multiplication", err)
	}
	stack.Push(result)
	return nil
}

// Divide pops two values from stack, divides (a / b), and pushes result.
func (o *Operations) Divide(stack *Stack) error {
	b, err := popStack(stack, "/")
	if err != nil {
		return err
	}

	if b.IsZero() {
		return buildError("/", fmt.Errorf("division by zero"))
	}

	a, err := popStack(stack, "/")
	if err != nil {
		return err
	}

	result, err := a.Div(b)
	if err != nil {
		return buildError("division", err)
	}
	stack.Push(result)
	return nil
}

// Power pops two values from stack, raises first to power of second (a ^ b), and pushes result.
func (o *Operations) Power(stack *Stack) error {
	a, b, err := popTwo(stack, "^")
	if err != nil {
		return err
	}

	result, err := a.Pow(b)
	if err != nil {
		return buildError("power", err)
	}
	stack.Push(result)
	return nil
}

// Modulo pops two values from stack, computes modulo (a % b), and pushes result.
func (o *Operations) Modulo(stack *Stack) error {
	a, b, err := popTwo(stack, "%")
	if err != nil {
		return err
	}

	// Check if operands are symbols (not supported for arithmetic)
	if sym, ok := a.(*Symbol); ok {
		return fmt.Errorf("symbol %s cannot be used with modulo operator", sym.Name())
	}
	if sym, ok := b.(*Symbol); ok {
		return fmt.Errorf("symbol %s cannot be used with modulo operator", sym.Name())
	}

	if b.IsZero() {
		return buildError("%", fmt.Errorf("modulo by zero"))
	}

	result, err := a.Mod(b)
	if err != nil {
		return buildError("%", err)
	}
	stack.Push(result)
	return nil
}

// FastPower pops two values from stack, raises first to integer power of second (a ** b), and pushes result.
// Uses binary exponentiation for efficiency with large integer exponents.
func (o *Operations) FastPower(stack *Stack) error {
	b, err := popStack(stack, "**")
	if err != nil {
		return err
	}

	a, err := popStack(stack, "**")
	if err != nil {
		return err
	}

	// Get the integer exponent from b
	bVal, err := b.Float64()
	if err != nil {
		return buildError("**", fmt.Errorf("exponent must be a number: %w", err))
	}

	exp := int(bVal)
	if float64(exp) != bVal {
		return buildError("**", fmt.Errorf("exponent must be an integer, got %v", bVal))
	}

	result, err := a.PowInt(exp)
	if err != nil {
		return buildError("**", err)
	}
	stack.Push(result)
	return nil
}

// Log2 pops one value from stack, computes log base 2 (log₂(a)), and pushes result.
func (o *Operations) Log2(stack *Stack) error {
	a, err := popStack(stack, "lg")
	if err != nil {
		return err
	}

	val, err := toFloat64(a, "log2")
	if err != nil {
		return err
	}
	if val <= 0 {
		return buildError("lg", fmt.Errorf("log2 undefined for non-positive numbers"))
	}

	// Compute log2 using the number interface
	mode := o.GetMode()
	stack.Push(NewNumber(math.Log2(val), mode))
	return nil
}

// Log10 pops one value from stack, computes log base 10 (log₁₀(a)), and pushes result.
func (o *Operations) Log10(stack *Stack) error {
	a, err := popStack(stack, "log")
	if err != nil {
		return err
	}

	val, err := toFloat64(a, "log10")
	if err != nil {
		return err
	}
	if val <= 0 {
		return buildError("log", fmt.Errorf("log10 undefined for non-positive numbers"))
	}

	// Compute log10 using the number interface
	mode := o.GetMode()
	stack.Push(NewNumber(math.Log10(val), mode))
	return nil
}

// Ln pops one value from stack, computes natural log (ln(a)), and pushes result.
func (o *Operations) Ln(stack *Stack) error {
	a, err := popStack(stack, "ln")
	if err != nil {
		return err
	}

	val, err := toFloat64(a, "ln")
	if err != nil {
		return err
	}
	if val <= 0 {
		return buildError("ln", fmt.Errorf("ln undefined for non-positive numbers"))
	}

	// Compute ln using the number interface
	mode := o.GetMode()
	stack.Push(NewNumber(math.Log(val), mode))
	return nil
}

// Hyper operators - operate on all values on the stack

// HyperAdd pops all values from stack, adds them left-associative (with boolean-to-number coercion), and pushes result.
func (o *Operations) HyperAdd(stack *Stack) error {
	values, err := popAll(stack, "hyperadd")
	if err != nil {
		return err
	}

	// Process left-associative with Number interface
	sum := 0.0
	for i := 0; i < len(values); i++ {
		val, err := values[i].Float64()
		if err != nil {
			return buildError("hyperadd", fmt.Errorf("failed to get float64 value: %w", err))
		}
		sum += val
	}
	mode := o.GetMode()
	stack.Push(NewNumber(sum, mode))
	return nil
}

// HyperMultiply pops all values from stack, multiplies them left-associative, and pushes result.
func (o *Operations) HyperMultiply(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hypermultiply: need at least 2 values")
	}

	product := 1.0
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return fmt.Errorf("hypermultiply: %w", err)
		}
		floatVal, err := val.Float64()
		if err != nil {
			return fmt.Errorf("hypermultiply: failed to get float64 value: %w", err)
		}
		product *= floatVal
	}
	mode := o.GetMode()
	stack.Push(NewNumber(product, mode))
	return nil
}

// HyperSubtract pops all values from stack, subtracts them left-associative, and pushes result.
func (o *Operations) HyperSubtract(stack *Stack) error {
	values, err := popAll(stack, "hypersubtract")
	if err != nil {
		return err
	}

	// Process left-associative with Number interface
	firstVal, err := values[0].Float64()
	if err != nil {
		return buildError("hypersubtract", fmt.Errorf("failed to get float64 value: %w", err))
	}
	result := firstVal
	for i := 1; i < len(values); i++ {
		val, err := values[i].Float64()
		if err != nil {
			return buildError("hypersubtract", fmt.Errorf("failed to get float64 value: %w", err))
		}
		result -= val
	}
	mode := o.GetMode()
	stack.Push(NewNumber(result, mode))
	return nil
}

// HyperDivide pops all values from stack, divides them left-associative, and pushes result.
func (o *Operations) HyperDivide(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperdivide: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Number
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return fmt.Errorf("hyperdivide: %w", err)
		}
		values = append(values, val)
	}

	// Reverse to get left-to-right order (first pushed = first in)
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	// Process left-associative with Number interface
	firstVal, err := values[0].Float64()
	if err != nil {
		return fmt.Errorf("hyperdivide: failed to get float64 value: %w", err)
	}
	result := firstVal
	for i := 1; i < len(values); i++ {
		val, err := values[i].Float64()
		if err != nil {
			return fmt.Errorf("hyperdivide: failed to get float64 value: %w", err)
		}
		if val == 0 {
			return fmt.Errorf("division by zero")
		}
		result /= val
	}
	mode := o.GetMode()
	stack.Push(NewNumber(result, mode))
	return nil
}

// HyperPower pops all values from stack, raises to power left-associative, and pushes result.
func (o *Operations) HyperPower(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperpower: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Number
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return fmt.Errorf("hyperpower: %w", err)
		}
		values = append(values, val)
	}

	// Reverse to get left-to-right order (first pushed = first in)
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	// Process left-associative with Number interface
	firstVal, err := values[0].Float64()
	if err != nil {
		return fmt.Errorf("hyperpower: failed to get float64 value: %w", err)
	}
	result := firstVal
	for i := 1; i < len(values); i++ {
		val, err := values[i].Float64()
		if err != nil {
			return fmt.Errorf("hyperpower: failed to get float64 value: %w", err)
		}
		result = math.Pow(result, val)
	}
	mode := o.GetMode()
	stack.Push(NewNumber(result, mode))
	return nil
}

// HyperModulo pops all values from stack, computes modulo left-associative, and pushes result.
func (o *Operations) HyperModulo(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hypermodulo: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Number
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return fmt.Errorf("hypermodulo: %w", err)
		}
		values = append(values, val)
	}

	// Reverse to get left-to-right order (first pushed = first in)
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	// Process left-associative with Number interface
	firstVal, err := values[0].Float64()
	if err != nil {
		return fmt.Errorf("hypermodulo: failed to get float64 value: %w", err)
	}
	result := firstVal
	for i := 1; i < len(values); i++ {
		val, err := values[i].Float64()
		if err != nil {
			return fmt.Errorf("hypermodulo: failed to get float64 value: %w", err)
		}
		if val == 0 {
			return fmt.Errorf("modulo by zero")
		}
		result = math.Mod(result, val)
	}
	mode := o.GetMode()
	stack.Push(NewNumber(result, mode))
	return nil
}

// HyperLog2 pops all values from stack, computes sum of log2 for all values, and pushes result.
// This follows the same pattern as HyperAdd (sum) and HyperMultiply (product).
func (o *Operations) HyperLog2(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperlog2: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Number
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return fmt.Errorf("hyperlog2: %w", err)
		}
		values = append(values, val)
	}

	// Reverse to get left-to-right order (first pushed = first in)
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	// Sum the log2 of all values using Float64() for value conversion:
	// - true → 1, false → 0
	var result float64 = 0
	for i := 0; i < len(values); i++ {
		val, err := values[i].Float64()
		if err != nil {
			return fmt.Errorf("hyperlog2: failed to get float64 value: %w", err)
		}
		if val <= 0 {
			return fmt.Errorf("hyperlog2 undefined for non-positive numbers")
		}
		result += math.Log2(val)
	}

	// Push the result as a Number
	mode := o.GetMode()
	stack.Push(NewNumber(result, mode))
	return nil
}

// HyperLog10 pops all values from stack, computes sum of log10 for all values, and pushes result.
// This follows the same pattern as HyperAdd (sum) and HyperMultiply (product).
func (o *Operations) HyperLog10(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperlog10: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Number
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return fmt.Errorf("hyperlog10: %w", err)
		}
		values = append(values, val)
	}

	// Reverse to get left-to-right order (first pushed = first in)
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	// Sum the log10 of all values using Float64() for value conversion:
	// - true → 1, false → 0
	var result float64 = 0
	for i := 0; i < len(values); i++ {
		val, err := values[i].Float64()
		if err != nil {
			return fmt.Errorf("hyperlog10: failed to get float64 value: %w", err)
		}
		if val <= 0 {
			return fmt.Errorf("hyperlog10 undefined for non-positive numbers")
		}
		result += math.Log10(val)
	}

	// Push the result as a Number
	mode := o.GetMode()
	stack.Push(NewNumber(result, mode))
	return nil
}

// HyperLn pops all values from stack, computes sum of natural log for all values, and pushes result.
// This follows the same pattern as HyperAdd (sum) and HyperMultiply (product).
func (o *Operations) HyperLn(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperln: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Number
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return fmt.Errorf("hyperln: %w", err)
		}
		values = append(values, val)
	}

	// Reverse to get left-to-right order (first pushed = first in)
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	// Sum the natural log of all values using Float64() for value conversion:
	// - true → 1, false → 0
	var result float64 = 0
	for i := 0; i < len(values); i++ {
		val, err := values[i].Float64()
		if err != nil {
			return fmt.Errorf("hyperln: failed to get float64 value: %w", err)
		}
		if val <= 0 {
			return fmt.Errorf("hyperln undefined for non-positive numbers")
		}
		result += math.Log(val)
	}
	mode := o.GetMode()
	stack.Push(NewNumber(result, mode))
	return nil
}

// GT pops two values from stack, compares (a > b), and pushes a boolean result.
func (o *Operations) GT(stack *Stack) error {
	a, b, err := popTwo(stack, "gt")
	if err != nil {
		return err
	}

	aVal, err := toFloat64(a, "gt comparison for a")
	if err != nil {
		return err
	}
	bVal, err := toFloat64(b, "gt comparison for b")
	if err != nil {
		return err
	}

	stack.Push(NewFloatFromBool(aVal > bVal))
	return nil
}

// LT pops two values from stack, compares (a < b), and pushes a boolean result.
func (o *Operations) LT(stack *Stack) error {
	a, b, err := popTwo(stack, "lt")
	if err != nil {
		return err
	}

	aVal, err := toFloat64(a, "lt comparison for a")
	if err != nil {
		return err
	}
	bVal, err := toFloat64(b, "lt comparison for b")
	if err != nil {
		return err
	}

	stack.Push(NewFloatFromBool(aVal < bVal))
	return nil
}

// GTE pops two values from stack, compares (a >= b), and pushes a boolean result.
func (o *Operations) GTE(stack *Stack) error {
	a, b, err := popTwo(stack, "gte")
	if err != nil {
		return err
	}

	aVal, err := toFloat64(a, "gte comparison for a")
	if err != nil {
		return err
	}
	bVal, err := toFloat64(b, "gte comparison for b")
	if err != nil {
		return err
	}

	stack.Push(NewFloatFromBool(aVal >= bVal))
	return nil
}

// LTE pops two values from stack, compares (a <= b), and pushes a boolean result.
func (o *Operations) LTE(stack *Stack) error {
	a, b, err := popTwo(stack, "lte")
	if err != nil {
		return err
	}

	aVal, err := toFloat64(a, "lte comparison for a")
	if err != nil {
		return err
	}
	bVal, err := toFloat64(b, "lte comparison for b")
	if err != nil {
		return err
	}

	stack.Push(NewFloatFromBool(aVal <= bVal))
	return nil
}

// EQ pops two values from stack, compares (a == b), and pushes a boolean result.
func (o *Operations) EQ(stack *Stack) error {
	a, b, err := popTwo(stack, "eq")
	if err != nil {
		return err
	}

	aVal, err := toFloat64(a, "eq comparison for a")
	if err != nil {
		return err
	}
	bVal, err := toFloat64(b, "eq comparison for b")
	if err != nil {
		return err
	}

	stack.Push(NewFloatFromBool(aVal == bVal))
	return nil
}

// NEQ pops two values from stack, compares (a != b), and pushes a boolean result.
func (o *Operations) NEQ(stack *Stack) error {
	a, b, err := popTwo(stack, "neq")
	if err != nil {
		return err
	}

	aVal, err := toFloat64(a, "neq comparison for a")
	if err != nil {
		return err
	}
	bVal, err := toFloat64(b, "neq comparison for b")
	if err != nil {
		return err
	}

	stack.Push(NewFloatFromBool(aVal != bVal))
	return nil
}

// stack manipulation operators

// Dup duplicates the top stack value.
func (o *Operations) Dup(stack *Stack) error {
	val, err := stack.Peek()
	if err != nil {
		return buildError("dup", err)
	}
	stack.Push(val)
	return nil
}

// Swap swaps the top two stack values.
func (o *Operations) Swap(stack *Stack) error {
	if err := ensureStackLength(stack, 2, "swap"); err != nil {
		return err
	}

	// Get the values without popping
	vals := stack.Values()
	top := vals[len(vals)-1]
	second := vals[len(vals)-2]

	// Pop both values
	if _, err := stack.Pop(); err != nil {
		return buildError("swap", fmt.Errorf("failed to pop top value: %w", err))
	}
	if _, err := stack.Pop(); err != nil {
		return buildError("swap", fmt.Errorf("failed to pop second value: %w", err))
	}

	// Push in swapped order
	stack.Push(top)
	stack.Push(second)

	return nil
}

// Pop removes and discards the top stack value.
func (o *Operations) Pop(stack *Stack) error {
	if _, err := stack.Pop(); err != nil {
		return buildError("pop", err)
	}
	return nil
}

// Show returns the current stack as a formatted string using the Number interface.
func (o *Operations) Show(stack *Stack) (string, error) {
	if stack.Len() == 0 {
		return "Stack is empty", nil
	}

	vals := stack.Values()
	var result string
	for i, val := range vals {
		if i > 0 {
			result += " "
		}
		// Use val.String() to format values correctly:
		// - Boolean values show as "true"/"false"
		// - Number values show with appropriate precision
		result += val.String()
	}
	return result, nil
}

// variables operations

// AssignVariable assigns a value from stack to a variable.
// Usage: `name value =`
func (o *Operations) AssignVariable(stack *Stack, name string) error {
	if name == "" {
		return fmt.Errorf("variable name cannot be empty")
	}

	if stack.Len() < 1 {
		return buildError("=", fmt.Errorf("insufficient operands: need value"))
	}

	val, err := popStack(stack, "=")
	if err != nil {
		return err
	}

	// Convert Number to float64 for variable storage
	valF, err := toFloat64(val, "assigning variable")
	if err != nil {
		return err
	}
	return o.vars.SetVariable(name, valF)
}

// UseVariable pushes a variable's value onto the stack.
// Usage: `varname` (pushes stored value)
func (o *Operations) UseVariable(stack *Stack, name string) error {
	if name == "" {
		return fmt.Errorf("variable name cannot be empty")
	}

	val, exists := o.vars.GetVariable(name)
	if !exists {
		return fmt.Errorf("%w: %s", ErrVariableNotFound, name)
	}

	mode := o.GetMode()
	stack.Push(NewNumber(val, mode))
	return nil
}

// DeleteVariable removes a variable.
// Usage: `name d`
func (o *Operations) DeleteVariable(name string) error {
	if name == "" {
		return fmt.Errorf("variable name cannot be empty")
	}

	deleted := o.vars.DeleteVariable(name)
	if !deleted {
		return fmt.Errorf("%w: %s", ErrVariableNotFound, name)
	}
	return nil
}

// ListVariables lists all variables.
// Usage: `vars`
func (o *Operations) ListVariables() (string, error) {
	return o.vars.FormatVariables(), nil
}

// ClearVariables removes all variables.
// Usage: `clear`
func (o *Operations) ClearVariables() {
	o.vars.ClearVariables()
}

// ListConstants lists all constants.
// Usage: `constants`
func (o *Operations) ListConstants() (string, error) {
	infos := o.consts.ListConstants()
	if len(infos) == 0 {
		return "No constants defined", nil
	}
	var sb strings.Builder
	for i, info := range infos {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(info.Name)
		sb.WriteString(" = ")
		// Use Number interface for consistent formatting
		num := NewNumber(info.Value, FloatMode)
		sb.WriteString(num.String())
	}
	return sb.String(), nil
}

// ClearConstants removes all constants from storage.
// Note: This clears only user-defined constants; built-in constants are preserved.
// Usage: `clearconstants`
func (o *Operations) ClearConstants() {
	o.consts.ReloadBuiltInConstants()
}

// AssignLeft assigns a value to a variable (for =: operator).
// Stack order: value name =: (value on bottom, name on top).
// This function pops name first (top of stack), then value.
// Usage: `value name =:` (e.g., `5 x =:`)
func (o *Operations) AssignLeft(stack *Stack) error {
	name, err := popStack(stack, "=:")
	if err != nil {
		return err
	}

	val, err := popStack(stack, "=:")
	if err != nil {
		return err
	}

	// Get the variable name - handle Symbol, StringNum, or convert to string
	varName := ""
	switch v := name.(type) {
	case *Symbol:
		varName = v.Name()
	case *StringNum:
		varName = v.String()
	default:
		varName = name.String()
	}

	valF, err := toFloat64(val, "assigning variable")
	if err != nil {
		return err
	}
	return o.vars.SetVariable(varName, valF)
}

// AssignRight assigns a value to a variable (for := operator).
// Stack order: name value := (name on bottom, value on top).
// This function pops value first (top of stack), then name.
// Usage: `name value :=` (e.g., `x 5 :=`)
func (o *Operations) AssignRight(stack *Stack) error {
	val, err := popStack(stack, ":=")
	if err != nil {
		return err
	}

	name, err := popStack(stack, ":=")
	if err != nil {
		return err
	}

	// Get the variable name - handle Symbol, StringNum, or convert to string
	varName := ""
	switch v := name.(type) {
	case *Symbol:
		varName = v.Name()
	case *StringNum:
		varName = v.String()
	default:
		varName = name.String()
	}

	valF, err := toFloat64(val, "assigning variable")
	if err != nil {
		return err
	}
	return o.vars.SetVariable(varName, valF)
}
