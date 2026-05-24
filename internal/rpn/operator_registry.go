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

	registry.registerArithmeticOperators(op)
	registry.registerComparisonOperators(op)
	registry.registerStackOperators(op)
	registry.registerVariableOperators(op)
	registry.registerCommandOperators(op)
	registry.registerHyperOperators(op)

	return registry
}

// registerArithmeticOperators registers math and numeric operators.
func (r *OperatorRegistry) registerArithmeticOperators(op *Operations) {
	r.registerStandardOperator("+", func(stack *Stack) error { return op.Add(stack) })
	r.registerStandardOperator("-", func(stack *Stack) error { return op.Subtract(stack) })
	r.registerStandardOperator("*", func(stack *Stack) error { return op.Multiply(stack) })
	r.registerStandardOperator("/", func(stack *Stack) error { return op.Divide(stack) })
	r.registerStandardOperator("^", func(stack *Stack) error { return op.Power(stack) })
	r.registerStandardOperator("**", func(stack *Stack) error { return op.FastPower(stack) })
	r.registerStandardOperator("%", func(stack *Stack) error { return op.Modulo(stack) })
	r.registerStandardOperator("lg", func(stack *Stack) error { return op.Log2(stack) })
	r.registerStandardOperator("log", func(stack *Stack) error { return op.Log10(stack) })
	r.registerStandardOperator("ln", func(stack *Stack) error { return op.Ln(stack) })
}

// registerComparisonOperators registers equality and ordering operators.
func (r *OperatorRegistry) registerComparisonOperators(op *Operations) {
	r.registerStandardOperator("gt", func(stack *Stack) error { return op.GT(stack) })
	r.registerStandardOperator("lt", func(stack *Stack) error { return op.LT(stack) })
	r.registerStandardOperator("<", func(stack *Stack) error { return op.LT(stack) })
	r.registerStandardOperator(">", func(stack *Stack) error { return op.GT(stack) })
	r.registerStandardOperator("gte", func(stack *Stack) error { return op.GTE(stack) })
	r.registerStandardOperator(">=", func(stack *Stack) error { return op.GTE(stack) })
	r.registerStandardOperator("lte", func(stack *Stack) error { return op.LTE(stack) })
	r.registerStandardOperator("<=", func(stack *Stack) error { return op.LTE(stack) })
	r.registerStandardOperator("eq", func(stack *Stack) error { return op.EQ(stack) })
	r.registerStandardOperator("==", func(stack *Stack) error { return op.EQ(stack) })
	r.registerStandardOperator("neq", func(stack *Stack) error { return op.NEQ(stack) })
	r.registerStandardOperator("!=", func(stack *Stack) error { return op.NEQ(stack) })
}

// registerStackOperators registers stack manipulation operators.
func (r *OperatorRegistry) registerStackOperators(op *Operations) {
	r.registerStandardOperator("dup", func(stack *Stack) error { return op.Dup(stack) })
	r.registerStandardOperator("swap", func(stack *Stack) error { return op.Swap(stack) })
	r.registerStandardOperator("pop", func(stack *Stack) error { return op.Pop(stack) })
	r.registerStandardOperator("d", func(stack *Stack) error {
		val, err := popStack(stack, "d")
		if err != nil {
			return err
		}
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
}

// registerVariableOperators registers assignment and conversion operators.
func (r *OperatorRegistry) registerVariableOperators(op *Operations) {
	r.registerStandardOperator("=", func(stack *Stack) error { return op.AssignRight(stack) })
	r.registerStandardOperator(":=", func(stack *Stack) error { return op.AssignRight(stack) })
	r.registerStandardOperator("=:", func(stack *Stack) error { return op.AssignLeft(stack) })
	r.registerStandardOperator("convert", func(stack *Stack) error { return op.Convert(stack) })
}

// registerCommandOperators registers operators that return a result immediately.
func (r *OperatorRegistry) registerCommandOperators(op *Operations) {
	r.registerCommandOperator("show", func(stack *Stack) (string, error) { return op.Show(stack) })
	r.registerCommandOperator("showstack", func(stack *Stack) (string, error) { return op.Show(stack) })
	r.registerCommandOperator("print", func(stack *Stack) (string, error) { return op.Show(stack) })
	r.registerCommandOperator("vars", func(stack *Stack) (string, error) { return op.ListVariables() })
	r.registerCommandOperator("constants", func(stack *Stack) (string, error) { return op.ListConstants() })
	r.registerCommandOperator("clear", func(stack *Stack) (string, error) { op.ClearVariables(); return "All variables cleared", nil })
	r.registerCommandOperator("clearconstants", func(stack *Stack) (string, error) { op.ClearConstants(); return "All constants cleared", nil })
}

// registerHyperOperators registers hyper (vectorized) operators.
func (r *OperatorRegistry) registerHyperOperators(op *Operations) {
	r.registerHyperOperator("[+]", func(stack *Stack) error { return op.HyperAdd(stack) })
	r.registerHyperOperator("[-]", func(stack *Stack) error { return op.HyperSubtract(stack) })
	r.registerHyperOperator("[*]", func(stack *Stack) error { return op.HyperMultiply(stack) })
	r.registerHyperOperator("[/]", func(stack *Stack) error { return op.HyperDivide(stack) })
	r.registerHyperOperator("[^]", func(stack *Stack) error { return op.HyperPower(stack) })
	r.registerHyperOperator("[%]", func(stack *Stack) error { return op.HyperModulo(stack) })
	r.registerHyperOperator("[lg]", func(stack *Stack) error { return op.HyperLog2(stack) })
	r.registerHyperOperator("[log]", func(stack *Stack) error { return op.HyperLog10(stack) })
	r.registerHyperOperator("[ln]", func(stack *Stack) error { return op.HyperLn(stack) })
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
