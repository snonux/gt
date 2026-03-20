package rpn

import (
	"fmt"
	"math"
)

// Operations provides operator implementations and stack manipulation.
type Operations struct {
	vars *Variables
}

// NewOperations creates a new Operations instance with the given variable store.
func NewOperations(vars *Variables) *Operations {
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
