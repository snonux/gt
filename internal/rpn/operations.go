// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"math"
)

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
}

// Operator is the combined interface for all operator implementations.
// This allows RPN to depend on an abstraction instead of the concrete Operations type.
type Operator interface {
	ArithmeticOperator
	HyperOperator
	StackOperator
	VariableOperator
	// SetMode sets the calculation mode for number formatting
	SetMode(CalculationMode)
}

// Operations provides operator implementations and stack manipulation.
type Operations struct {
	vars VariableStore
	mode CalculationMode
}

// NewOperations creates a new Operations instance with the given variable store.
func NewOperations(vars VariableStore) *Operations {
	return &Operations{
		vars: vars,
		mode: FloatMode, // default
	}
}

// SetMode sets the calculation mode for the Operations instance.
func (o *Operations) SetMode(mode CalculationMode) {
	o.mode = mode
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
	registry.registerStandardOperator("%", func(stack *Stack) error { return op.Modulo(stack) })
	registry.registerStandardOperator("lg", func(stack *Stack) error { return op.Log2(stack) })
	registry.registerStandardOperator("log", func(stack *Stack) error { return op.Log10(stack) })
	registry.registerStandardOperator("ln", func(stack *Stack) error { return op.Ln(stack) })
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
	registry.registerCommandOperator("clear", func(stack *Stack) (string, error) { op.ClearVariables(); return "All variables cleared", nil })

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

// Add pops two values from stack, adds them (with boolean-to-number coercion), and pushes result.
func (o *Operations) Add(stack *Stack) error {
	bVal, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for +: %w", err)
	}

	aVal, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for +: %w", err)
	}

	// Use toNumber for automatic boolean-to-number coercion
	stack.Push(NewNumberValue(toNumber(aVal) + toNumber(bVal)))
	return nil
}

// Subtract pops two values from stack, subtracts (a - b), and pushes result.
func (o *Operations) Subtract(stack *Stack) error {
	b, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for -: %w", err)
	}

	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for -: %w", err)
	}

	stack.Push(NewNumberValue(toNumber(a) - toNumber(b)))
	return nil
}

// Multiply pops two values from stack, multiplies them, and pushes result.
func (o *Operations) Multiply(stack *Stack) error {
	b, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for *: %w", err)
	}

	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for *: %w", err)
	}

	stack.Push(NewNumberValue(toNumber(a) * toNumber(b)))
	return nil
}

// Divide pops two values from stack, divides (a / b), and pushes result.
func (o *Operations) Divide(stack *Stack) error {
	b, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for /: %w", err)
	}

	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for /: %w", err)
	}

	if toNumber(b) == 0 {
		return fmt.Errorf("division by zero")
	}

	stack.Push(NewNumberValue(toNumber(a) / toNumber(b)))
	return nil
}

// Power pops two values from stack, raises first to power of second (a ^ b), and pushes result.
func (o *Operations) Power(stack *Stack) error {
	b, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for ^: %w", err)
	}

	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for ^: %w", err)
	}

	stack.Push(NewNumberValue(math.Pow(toNumber(a), toNumber(b))))
	return nil
}

// Modulo pops two values from stack, computes modulo (a % b), and pushes result.
func (o *Operations) Modulo(stack *Stack) error {
	b, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for %%: %w", err)
	}

	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for %%: %w", err)
	}

	if toNumber(b) == 0 {
		return fmt.Errorf("modulo by zero")
	}

	stack.Push(NewNumberValue(math.Mod(toNumber(a), toNumber(b))))
	return nil
}

// Log2 pops one value from stack, computes log base 2 (log₂(a)), and pushes result.
func (o *Operations) Log2(stack *Stack) error {
	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for lg: %w", err)
	}

	if toNumber(a) <= 0 {
		return fmt.Errorf("log2 undefined for non-positive numbers")
	}

	stack.Push(NewNumberValue(math.Log2(toNumber(a))))
	return nil
}

// Log10 pops one value from stack, computes log base 10 (log₁₀(a)), and pushes result.
func (o *Operations) Log10(stack *Stack) error {
	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for log: %w", err)
	}

	if toNumber(a) <= 0 {
		return fmt.Errorf("log10 undefined for non-positive numbers")
	}

	stack.Push(NewNumberValue(math.Log10(toNumber(a))))
	return nil
}

// Ln pops one value from stack, computes natural log (ln(a)), and pushes result.
func (o *Operations) Ln(stack *Stack) error {
	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for ln: %w", err)
	}

	if toNumber(a) <= 0 {
		return fmt.Errorf("ln undefined for non-positive numbers")
	}

	stack.Push(NewNumberValue(math.Log(toNumber(a))))
	return nil
}

// Hyper operators - operate on all values on the stack

// HyperAdd pops all values from stack, adds them left-associative (with boolean-to-number coercion), and pushes result.
func (o *Operations) HyperAdd(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperadd: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Value
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return fmt.Errorf("hyperadd: %w", err)
		}
		values = append(values, val)
	}

	// Reverse to get left-to-right order (first pushed = first in)
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	// Process left-associative with toNumber coercion
	sum := 0.0
	for i := 0; i < len(values); i++ {
		sum += toNumber(values[i])
	}
	stack.Push(NewNumberValue(sum))
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
		product *= toNumber(val)
	}
	stack.Push(NewNumberValue(product))
	return nil
}

// HyperSubtract pops all values from stack, subtracts them left-associative, and pushes result.
func (o *Operations) HyperSubtract(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hypersubtract: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Value
	for stack.Len() > 0 {
		val, err := stack.Pop()
		if err != nil {
			return fmt.Errorf("hypersubtract: %w", err)
		}
		values = append(values, val)
	}

	// Reverse to get left-to-right order (first pushed = first in)
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	// Process left-associative with toNumber coercion
	result := toNumber(values[0])
	for i := 1; i < len(values); i++ {
		result -= toNumber(values[i])
	}
	stack.Push(NewNumberValue(result))
	return nil
}

// HyperDivide pops all values from stack, divides them left-associative, and pushes result.
func (o *Operations) HyperDivide(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperdivide: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Value
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

	// Process left-associative with toNumber coercion
	result := toNumber(values[0])
	for i := 1; i < len(values); i++ {
		if toNumber(values[i]) == 0 {
			return fmt.Errorf("division by zero")
		}
		result /= toNumber(values[i])
	}
	stack.Push(NewNumberValue(result))
	return nil
}

// HyperPower pops all values from stack, raises to power left-associative, and pushes result.
func (o *Operations) HyperPower(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperpower: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Value
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

	// Process left-associative with toNumber coercion
	result := toNumber(values[0])
	for i := 1; i < len(values); i++ {
		result = math.Pow(result, toNumber(values[i]))
	}
	stack.Push(NewNumberValue(result))
	return nil
}

// HyperModulo pops all values from stack, computes modulo left-associative, and pushes result.
func (o *Operations) HyperModulo(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hypermodulo: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Value
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

	// Process left-associative with toNumber coercion
	result := toNumber(values[0])
	for i := 1; i < len(values); i++ {
		if toNumber(values[i]) == 0 {
			return fmt.Errorf("modulo by zero")
		}
		result = math.Mod(result, toNumber(values[i]))
	}
	stack.Push(NewNumberValue(result))
	return nil
}

// HyperLog2 pops all values from stack, computes sum of log2 for all values, and pushes result.
// This follows the same pattern as HyperAdd (sum) and HyperMultiply (product).
func (o *Operations) HyperLog2(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperlog2: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Value
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

	// Sum the log2 of all values with toNumber coercion
	var result float64 = 0
	for i := 0; i < len(values); i++ {
		if toNumber(values[i]) <= 0 {
			return fmt.Errorf("hyperlog2 undefined for non-positive numbers")
		}
		result += math.Log2(toNumber(values[i]))
	}
	stack.Push(NewNumberValue(result))
	return nil
}

// HyperLog10 pops all values from stack, computes sum of log10 for all values, and pushes result.
// This follows the same pattern as HyperAdd (sum) and HyperMultiply (product).
func (o *Operations) HyperLog10(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperlog10: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Value
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

	// Sum the log10 of all values with toNumber coercion
	var result float64 = 0
	for i := 0; i < len(values); i++ {
		if toNumber(values[i]) <= 0 {
			return fmt.Errorf("hyperlog10 undefined for non-positive numbers")
		}
		result += math.Log10(toNumber(values[i]))
	}
	stack.Push(NewNumberValue(result))
	return nil
}

// HyperLn pops all values from stack, computes sum of natural log for all values, and pushes result.
// This follows the same pattern as HyperAdd (sum) and HyperMultiply (product).
func (o *Operations) HyperLn(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperln: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []Value
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

	// Sum the natural log of all values with toNumber coercion
	var result float64 = 0
	for i := 0; i < len(values); i++ {
		if toNumber(values[i]) <= 0 {
			return fmt.Errorf("hyperln undefined for non-positive numbers")
		}
		result += math.Log(toNumber(values[i]))
	}
	stack.Push(NewNumberValue(result))
	return nil
}

// stack manipulation operators

// Dup duplicates the top stack value.
func (o *Operations) Dup(stack *Stack) error {
	val, err := stack.Peek()
	if err != nil {
		return fmt.Errorf("insufficient operands for dup: %w", err)
	}
	stack.Push(val)
	return nil
}

// Swap swaps the top two stack values.
func (o *Operations) Swap(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for swap: need at least 2 values")
	}

	// Get the values without popping
	vals := stack.Values()
	top := vals[len(vals)-1]
	second := vals[len(vals)-2]

	// Pop both values - we know this won't fail because we checked stack.Len() >= 2 above
	if _, err := stack.Pop(); err != nil {
		return fmt.Errorf("swap: failed to pop top value: %w", err)
	}
	if _, err := stack.Pop(); err != nil {
		return fmt.Errorf("swap: failed to pop second value: %w", err)
	}

	// Push in swapped order
	stack.Push(top)
	stack.Push(second)

	return nil
}

// Pop removes and discards the top stack value.
func (o *Operations) Pop(stack *Stack) error {
	if _, err := stack.Pop(); err != nil {
		return fmt.Errorf("insufficient operands for pop: %w", err)
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
		// Use Number interface for consistent formatting with the current mode
		num := NewNumber(toNumber(val), o.mode)
		result += num.String()
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
		return fmt.Errorf("insufficient operands for assignment: need value")
	}

	val, err := stack.Pop()
	if err != nil {
		return err
	}

	// Convert Value to float64 for variable storage
	o.vars.SetVariable(name, toNumber(val))
	return nil
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

	stack.Push(NewNumberValue(val))
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
