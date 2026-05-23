// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

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
	DeleteVariable(name string) error
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
// boolean ops) with configuration methods (SetMode, SetPrefixMode, GetPrefixMode),
// metric command handlers, and custom metric commands. Per ISP this could be split,
// but RPN is the sole client and splitting would add indirection without practical
// benefit. The concrete *Operations type satisfies this interface exclusively.
type Operator interface {
	ArithmeticOperator
	BooleanOperator
	HyperOperator
	StackOperator
	VariableOperator
	ConstantOperator
	PowerIntOperator
	// SetMode sets the calculation mode (e.g., FloatMode, RationalMode).
	SetMode(CalculationMode)
	// GetMode returns the current calculation mode.
	GetMode() CalculationMode
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
