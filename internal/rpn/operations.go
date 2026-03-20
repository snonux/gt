package rpn

import (
	"fmt"
	"math"
)

// Operator defines the interface for operator implementations and stack manipulation.
// This allows RPN to depend on an abstraction instead of the concrete Operations type.
type Operator interface {
	// Arithmetic operators
	Add(stack *Stack) error
	Subtract(stack *Stack) error
	Multiply(stack *Stack) error
	Divide(stack *Stack) error
	Power(stack *Stack) error
	Modulo(stack *Stack) error

	// Hyper operators
	HyperAdd(stack *Stack) error
	HyperSubtract(stack *Stack) error
	HyperMultiply(stack *Stack) error
	HyperDivide(stack *Stack) error
	HyperPower(stack *Stack) error
	HyperModulo(stack *Stack) error

	// Stack manipulation operators
	Dup(stack *Stack) error
	Swap(stack *Stack) error
	Pop(stack *Stack) error
	Show(stack *Stack) (string, error)

	// Variable operations
	ListVariables() (string, error)
	ClearVariables()
}

// Operations provides operator implementations and stack manipulation.
type Operations struct {
	vars VariableStore
}

// NewOperations creates a new Operations instance with the given variable store.
func NewOperations(vars VariableStore) *Operations {
	return &Operations{
		vars: vars,
	}
}

// arithmetic operators

// Add pops two values from stack, adds them, and pushes result.
func (o *Operations) Add(stack *Stack) error {
	b, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for +: %w", err)
	}

	a, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for +: %w", err)
	}

	stack.Push(a + b)
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

	stack.Push(a - b)
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

	stack.Push(a * b)
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

	if b == 0 {
		return fmt.Errorf("division by zero")
	}

	stack.Push(a / b)
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

	stack.Push(math.Pow(a, b))
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

	if b == 0 {
		return fmt.Errorf("modulo by zero")
	}

	stack.Push(math.Mod(a, b))
	return nil
}

// Hyper operators - operate on all values on the stack

// HyperAdd pops all values from stack, adds them left-associative, and pushes result.
func (o *Operations) HyperAdd(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperadd: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []float64
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

	// Process left-associative
	sum := 0.0
	for i := 0; i < len(values); i++ {
		sum += values[i]
	}
	stack.Push(sum)
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
		product *= val
	}
	stack.Push(product)
	return nil
}

// HyperSubtract pops all values from stack, subtracts them left-associative, and pushes result.
func (o *Operations) HyperSubtract(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hypersubtract: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []float64
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

	// Process left-associative
	result := values[0]
	for i := 1; i < len(values); i++ {
		result -= values[i]
	}
	stack.Push(result)
	return nil
}

// HyperDivide pops all values from stack, divides them left-associative, and pushes result.
func (o *Operations) HyperDivide(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperdivide: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []float64
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

	// Process left-associative
	result := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] == 0 {
			return fmt.Errorf("division by zero")
		}
		result /= values[i]
	}
	stack.Push(result)
	return nil
}

// HyperPower pops all values from stack, raises to power left-associative, and pushes result.
func (o *Operations) HyperPower(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hyperpower: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []float64
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

	// Process left-associative
	result := values[0]
	for i := 1; i < len(values); i++ {
		result = math.Pow(result, values[i])
	}
	stack.Push(result)
	return nil
}

// HyperModulo pops all values from stack, computes modulo left-associative, and pushes result.
func (o *Operations) HyperModulo(stack *Stack) error {
	if stack.Len() < 2 {
		return fmt.Errorf("insufficient operands for hypermodulo: need at least 2 values")
	}

	// Pop all values into a slice (in reverse order - top first)
	var values []float64
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

	// Process left-associative
	result := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] == 0 {
			return fmt.Errorf("modulo by zero")
		}
		result = math.Mod(result, values[i])
	}
	stack.Push(result)
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

	// Pop both
	stack.Pop()
	stack.Pop()

	// Push in swapped order
	stack.Push(top)
	stack.Push(second)

	return nil
}

// Pop removes and discards the top stack value.
func (o *Operations) Pop(stack *Stack) error {
	_, err := stack.Pop()
	if err != nil {
		return fmt.Errorf("insufficient operands for pop: %w", err)
	}
	return nil
}

// Show returns the current stack as a formatted string.
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
		result += fmt.Sprintf("%.10g", val)
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

	return o.vars.SetVariable(name, val)
}

// UseVariable pushes a variable's value onto the stack.
// Usage: `varname` (pushes stored value)
func (o *Operations) UseVariable(stack *Stack, name string) error {
	if name == "" {
		return fmt.Errorf("variable name cannot be empty")
	}

	val, exists := o.vars.GetVariable(name)
	if !exists {
		return fmt.Errorf("undefined variable: %s", name)
	}

	stack.Push(val)
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
		return fmt.Errorf("undefined variable: %s", name)
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
