// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import "fmt"

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
func NewOperatorRegistry(op *Operations) *OperatorRegistry {
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
	registry.registerStandardOperator("=", func(stack *Stack) error { return op.AssignRight(stack) })
	registry.registerStandardOperator(":=", func(stack *Stack) error { return op.AssignRight(stack) })
	registry.registerStandardOperator("=:", func(stack *Stack) error { return op.AssignLeft(stack) })
	registry.registerStandardOperator("convert", func(stack *Stack) error { return op.Convert(stack) })
	registry.registerStandardOperator("dup", func(stack *Stack) error { return op.Dup(stack) })
	registry.registerStandardOperator("swap", func(stack *Stack) error { return op.Swap(stack) })
	registry.registerStandardOperator("pop", func(stack *Stack) error { return op.Pop(stack) })
	registry.registerStandardOperator("d", func(stack *Stack) error {
		val, err := popStack(stack, "d")
		if err != nil {
			return err
		}
		// Extract variable name from the value
		var name string
		switch v := val.(type) {
		case *Symbol:
			name = v.Name()
		case *StringNum:
			name = v.String()
		default:
			return fmt.Errorf("delete expects a variable name, got %T", val)
		}
		return op.DeleteVariable(name)
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
