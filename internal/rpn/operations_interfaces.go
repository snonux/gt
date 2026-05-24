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
}

// LogarithmicOperator defines the interface for logarithmic operators.
type LogarithmicOperator interface {
	Log2(stack *Stack) error
	Log10(stack *Stack) error
	Ln(stack *Stack) error
}

// MetricOperator defines the interface for metric unit conversion.
type MetricOperator interface {
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
	Delete(stack *Stack) error
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

// ModeController defines the interface for calculation mode and prefix mode control.
type ModeController interface {
	SetMode(CalculationMode)
	GetMode() CalculationMode
	SetPrefixMode(PrefixMode)
	GetPrefixMode() PrefixMode
}

// MetricCommander defines the interface for metric query commands.
type MetricCommander interface {
	MetricRegistry() MetricReader
	MetricShow(stack *Stack) (string, error)
	MetricList(stack *Stack) (string, error)
	MetricCategory(stack *Stack, categoryName string) (string, error)
	MetricCompatible(stack *Stack) (string, error)
}

// CustomMetricManager defines the interface for custom metric operations.
type CustomMetricManager interface {
	CustomShow(stack *Stack, name string) (string, error)
	CustomList(stack *Stack) (string, error)
	CustomDefine(name string, factor float64, category string) error
	CustomUndefine(name string) error
}

// OperatorProvider combines all operator interfaces needed by the OperatorRegistry.
// This satisfies DIP by decoupling the registry from the concrete *Operations type.
type OperatorProvider interface {
	ArithmeticOperator
	PowerIntOperator
	LogarithmicOperator
	BooleanOperator
	StackOperator
	VariableOperator
	ConstantOperator
	MetricOperator
	HyperOperator
}

// Registration interfaces — focused composites for each operator group.
// Adding a new operator category only requires creating a new registration
// interface and helper, without editing OperatorProvider or existing helpers.
type (
	// ArithmeticOpProvider covers standard arithmetic, fast power, and log operators.
	ArithmeticOpProvider interface {
		ArithmeticOperator
		PowerIntOperator
		LogarithmicOperator
	}
	// ComparisonOpProvider covers boolean comparison operators.
	ComparisonOpProvider interface {
		BooleanOperator
	}
	// VariableOpProvider covers variable assignment and metric conversion.
	VariableOpProvider interface {
		VariableOperator
		MetricOperator
	}
	// CommandOpProvider covers command operators that return immediate results.
	CommandOpProvider interface {
		StackOperator
		VariableOperator
		ConstantOperator
	}
)

// OperationsProvider combines all interfaces that RPN needs from Operations.
// RPN depends on this interface (DIP) rather than the concrete *Operations type.
type OperationsProvider interface {
	ModeController
	StackOperator
	MetricCommander
	CustomMetricManager
}

// Operator implementations are split across focused sub-interfaces
// (ArithmeticOperator, LogarithmicOperator, MetricOperator, BooleanOperator,
// HyperOperator, StackOperator, VariableOperator, ConstantOperator,
// PowerIntOperator) for clarity.
//
// The OperationsProvider interface is the combined interface that RPN depends
// on, satisfying DIP by decoupling RPN from the concrete *Operations type.
//
// Each sub-interface and OperationsProvider is satisfied by *Operations,
// verified by compile-time checks in operations.go.
