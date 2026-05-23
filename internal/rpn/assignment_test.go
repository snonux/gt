// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"testing"
)

// TestAssignmentStandard tests standard assignment 'x 5 ='
func TestAssignmentStandard(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	result, err := rpn.ParseAndEvaluate("x 5 =")
	if err != nil {
		t.Fatalf("assignment 'x 5 =' returned error: %v", err)
	}
	if result != "x = 5" {
		t.Errorf("assignment result = %q, want 'x = 5'", result)
	}

	val, exists := vars.GetVariable("x")
	if !exists {
		t.Fatal("variable x should exist after assignment")
	}
	if val != 5 {
		t.Errorf("variable x = %v, want 5", val)
	}
}

// TestAssignmentRight tests right assignment '5 x :='
func TestAssignmentRight(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	result, err := rpn.ParseAndEvaluate("5 x :=")
	if err != nil {
		t.Fatalf("right assignment '5 x :=' returned error: %v", err)
	}
	if result != "x = 5" {
		t.Errorf("right assignment result = %q, want 'x = 5'", result)
	}

	val, exists := vars.GetVariable("x")
	if !exists {
		t.Fatal("variable x should exist after right assignment")
	}
	if val != 5 {
		t.Errorf("variable x = %v, want 5", val)
	}
}

// TestAssignmentLeft tests left assignment '5 x =:'
func TestAssignmentLeft(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	result, err := rpn.ParseAndEvaluate("5 x =:")
	if err != nil {
		t.Fatalf("left assignment '5 x =:' returned error: %v", err)
	}
	if result != "x = 5" {
		t.Errorf("left assignment result = %q, want 'x = 5'", result)
	}

	val, exists := vars.GetVariable("x")
	if !exists {
		t.Fatal("variable x should exist after left assignment")
	}
	if val != 5 {
		t.Errorf("variable x = %v, want 5", val)
	}
}

// TestAssignmentReassignment tests variable reassignment
func TestAssignmentReassignment(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	_, err := rpn.ParseAndEvaluate("x 5 =")
	if err != nil {
		t.Fatalf("initial assignment failed: %v", err)
	}

	result, err := rpn.ParseAndEvaluate("x 10 =")
	if err != nil {
		t.Fatalf("reassignment failed: %v", err)
	}
	if result != "x = 10" {
		t.Errorf("reassignment result = %q, want 'x = 10'", result)
	}

	val, exists := vars.GetVariable("x")
	if !exists {
		t.Fatal("variable x should exist after reassignment")
	}
	if val != 10 {
		t.Errorf("variable x = %v, want 10 after reassignment", val)
	}
}

// TestAssignmentInExpression tests assignment used in an expression
func TestAssignmentInExpression(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	// 'x 5 = x 2 +' should assign x=5, then push x and 2, then add: 5 + 2 = 7
	result, err := rpn.ParseAndEvaluate("x 5 = x 2 +")
	if err != nil {
		t.Fatalf("assignment in expression failed: %v", err)
	}
	if result != "7" {
		t.Errorf("expression result = %q, want '7'", result)
	}
}

// TestAssignmentMultipleVariables tests multiple variable assignments
func TestAssignmentMultipleVariables(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	_, _ = rpn.ParseAndEvaluate("x 5 =")
	_, _ = rpn.ParseAndEvaluate("y 10 =")
	_, _ = rpn.ParseAndEvaluate("z 15 =")

	if vars.Count() != 3 {
		t.Errorf("vars.Count() = %d, want 3", vars.Count())
	}

	val, _ := vars.GetVariable("x")
	if val != 5 {
		t.Errorf("x = %v, want 5", val)
	}
	val, _ = vars.GetVariable("y")
	if val != 10 {
		t.Errorf("y = %v, want 10", val)
	}
	val, _ = vars.GetVariable("z")
	if val != 15 {
		t.Errorf("z = %v, want 15", val)
	}
}

// TestAssignmentArithmeticWithVariable tests arithmetic using assigned variable
func TestAssignmentArithmeticWithVariable(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	_, _ = rpn.ParseAndEvaluate("x 5 =")

	// x x + should be 5 + 5 = 10
	result, err := rpn.ParseAndEvaluate("x x +")
	if err != nil {
		t.Fatalf("arithmetic with variable failed: %v", err)
	}
	if result != "10" {
		t.Errorf("x x + = %q, want '10'", result)
	}
}

// TestAssignmentNegativeValue tests assignment with negative value
func TestAssignmentNegativeValue(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	_, err := rpn.ParseAndEvaluate("x -5 =")
	if err != nil {
		t.Fatalf("negative assignment failed: %v", err)
	}

	val, exists := vars.GetVariable("x")
	if !exists {
		t.Fatal("variable x should exist")
	}
	if val != -5 {
		t.Errorf("x = %v, want -5", val)
	}
}

// TestAssignmentDecimalValue tests assignment with decimal value
func TestAssignmentDecimalValue(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	_, err := rpn.ParseAndEvaluate("x 3.14 =")
	if err != nil {
		t.Fatalf("decimal assignment failed: %v", err)
	}

	val, exists := vars.GetVariable("x")
	if !exists {
		t.Fatal("variable x should exist")
	}
	if val != 3.14 {
		t.Errorf("x = %v, want 3.14", val)
	}
}

// TestAssignmentUnderscoreVariable tests assignment with underscore variable
func TestAssignmentUnderscoreVariable(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	_, err := rpn.ParseAndEvaluate("my_var 42 =")
	if err != nil {
		t.Fatalf("underscore variable assignment failed: %v", err)
	}

	val, exists := vars.GetVariable("my_var")
	if !exists {
		t.Fatal("variable my_var should exist")
	}
	if val != 42 {
		t.Errorf("my_var = %v, want 42", val)
	}
}

// TestAssignmentWithCalculationResult tests assigning result of calculation
func TestAssignmentWithCalculationResult(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	// 3 4 + x =: should push 3, 4, add (=7), then assign 7 to x
	_, err := rpn.ParseAndEvaluate("3 4 + x =:")
	if err != nil {
		t.Fatalf("assignment with calculation failed: %v", err)
	}

	val, exists := vars.GetVariable("x")
	if !exists {
		t.Fatal("variable x should exist")
	}
	if val != 7 {
		t.Errorf("x = %v, want 7", val)
	}
}

// TestAssignmentAllOperatorsInSequence tests all three assignment operators
func TestAssignmentAllOperatorsInSequence(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	_, _ = rpn.ParseAndEvaluate("a 1 =")  // standard
	_, _ = rpn.ParseAndEvaluate("2 b :=") // right
	_, _ = rpn.ParseAndEvaluate("3 c =:") // left

	if val, exists := vars.GetVariable("a"); !exists || val != 1 {
		t.Errorf("a = %v (exists=%v), want 1", val, exists)
	}
	if val, exists := vars.GetVariable("b"); !exists || val != 2 {
		t.Errorf("b = %v (exists=%v), want 2", val, exists)
	}
	if val, exists := vars.GetVariable("c"); !exists || val != 3 {
		t.Errorf("c = %v (exists=%v), want 3", val, exists)
	}
}

// TestAssignmentOperatorRegistry verifies all assignment operators are registered
func TestAssignmentOperatorRegistry(t *testing.T) {
	vars := NewVariables()
	ops := NewOperations(vars)
	reg := NewOperatorRegistry(ops)

	// Verify assignment operators are registered as standard operators
	for _, op := range []string{"=", ":=", "=: "} {
		trimmed := op
		if op == "=: " {
			trimmed = "=:"
		}
		if !reg.IsStandardOperator(trimmed) {
			t.Errorf("operator %q should be registered as standard operator", trimmed)
		}
	}
}
