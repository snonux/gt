// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"strings"
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
	for _, op := range []string{":=", "=: "} {
		trimmed := op
		if op == "=: " {
			trimmed = "=:"
		}
		if !reg.IsStandardOperator(trimmed) {
			t.Errorf("operator %q should be registered as standard operator", trimmed)
		}
	}
}

// TestAssignmentNotTriggeredByEqualEqual ensures that == is not misparsed as an assignment.
func TestAssignmentNotTriggeredByEqualEqual(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	// "a == b" should NOT be treated as an assignment
	_, err := rpn.ParseAndEvaluate("a == b")
	// It may error (unknown tokens) or produce a result, but it should NOT assign a variable
	// The key check: no variable "a" should be created
	if _, exists := vars.GetVariable("a"); exists {
		t.Error("'a == b' should not create variable a")
	}
	// The error from parsing == as an unknown operator is acceptable
	if err == nil {
		// If no error, that's also fine - it just means == was handled somehow
		t.Logf("'a == b' returned: %v", err)
	}
}

// TestAssignmentNotTriggeredByNotEqual ensures that != is not misparsed as an assignment.
func TestAssignmentNotTriggeredByNotEqual(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	// "a != b" should NOT be treated as an assignment
	_, err := rpn.ParseAndEvaluate("a != b")
	// The key check: no variable "a" should be created
	if _, exists := vars.GetVariable("a"); exists {
		t.Error("'a != b' should not create variable a")
	}
	_ = err // error is acceptable
}

// TestAssignmentAfterEqualEqual ensures "x 5 =" still works even when == exists in the expression.
func TestAssignmentAfterEqualEqual(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	// A standalone assignment should work regardless of == elsewhere
	result, err := rpn.ParseAndEvaluate("x 5 =")
	if err != nil {
		t.Fatalf("'x 5 =' returned error: %v", err)
	}
	if result != "x = 5" {
		t.Errorf("result = %q, want 'x = 5'", result)
	}

	val, exists := vars.GetVariable("x")
	if !exists || val != 5 {
		t.Errorf("x = %v (exists=%v), want 5", val, exists)
	}
}

// TestAssignRightOperator tests the := (right assignment) operator via ParseAndEvaluate.
func TestAssignRightOperator(t *testing.T) {
	tests := []struct {
		name  string
		expr  string
		want  string
	}{
		{"basic right assignment", "x 5 :=", "x = 5"},
		{"right assignment with decimal", "pi 3.14159 :=", "pi = 3.14159"},
		{"right assignment with negative value", "neg -10 :=", "neg = -10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRPN(NewVariables())
			result, err := r.ParseAndEvaluate(tt.expr)
			if err != nil {
				t.Fatalf("ParseAndEvaluate(%q) error = %v", tt.expr, err)
			}
			if result != tt.want {
				t.Errorf("ParseAndEvaluate(%q) = %q, want %q", tt.expr, result, tt.want)
			}
		})
	}
}

// TestAssignLeftOperator tests the =: (left assignment) operator via ParseAndEvaluate.
func TestAssignLeftOperator(t *testing.T) {
	tests := []struct {
		name  string
		expr  string
		want  string
	}{
		{"basic left assignment", "10 y =:", "y = 10"},
		{"left assignment with decimal", "2.71828 e =:", "e = 2.71828"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRPN(NewVariables())
			result, err := r.ParseAndEvaluate(tt.expr)
			if err != nil {
				t.Fatalf("ParseAndEvaluate(%q) error = %v", tt.expr, err)
			}
			if result != tt.want {
				t.Errorf("ParseAndEvaluate(%q) = %q, want %q", tt.expr, result, tt.want)
			}
		})
	}
}

// TestStandardAssignOperator tests the = (standard assignment) operator via ParseAndEvaluate.
func TestStandardAssignOperator(t *testing.T) {
	tests := []struct {
		name  string
		expr  string
		want  string
	}{
		{"infix assignment", "x = 5", "x = 5"},
		{"postfix assignment", "x 10 =", "x = 10"},
		{"assignment with expression continuation", "x 10 = x 5 +", "15"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRPN(NewVariables())
			result, err := r.ParseAndEvaluate(tt.expr)
			if err != nil {
				t.Fatalf("ParseAndEvaluate(%q) error = %v", tt.expr, err)
			}
			if result != tt.want {
				t.Errorf("ParseAndEvaluate(%q) = %q, want %q", tt.expr, result, tt.want)
			}
		})
	}
}

// TestVariableInExpression tests using variables in subsequent expressions.
func TestVariableInExpression(t *testing.T) {
	tests := []struct {
		name  string
		setup string
		expr  string
		want  string
	}{
		{"single variable reuse", "x 5 :=", "x x +", "10"},
		{"multi-variable expression", "a 10 := b 3 :=", "a b +", "13"},
		{"variable in complex expression", "x 5 :=", "x 2 + 3 *", "21"},
		{"variable reassignment", "x 5 :=", "x 10 := x", "10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRPN(NewVariables())
			_, err := r.ParseAndEvaluate(tt.setup)
			if err != nil {
				t.Fatalf("Setup failed for %q: %v", tt.setup, err)
			}
			result, err := r.ParseAndEvaluate(tt.expr)
			if err != nil {
				t.Fatalf("ParseAndEvaluate(%q) error = %v", tt.expr, err)
			}
			if result != tt.want {
				t.Errorf("ParseAndEvaluate(%q) = %q, want %q", tt.expr, result, tt.want)
			}
		})
	}
}

// TestChainedAssignments tests chaining multiple assignments in one expression.
func TestChainedAssignments(t *testing.T) {
	r := NewRPN(NewVariables())
	result, err := r.ParseAndEvaluate("a 10 := b 3 := c 2 :=")
	if err != nil {
		t.Fatalf("Chained assignment failed: %v", err)
	}
	// Chained assignments return empty result (side effects)
	if result != "" {
		t.Errorf("Chained assignment result = %q, want empty", result)
	}

	// Verify variables were set
	if val, exists := r.vars.GetVariable("a"); !exists || val != 10 {
		t.Errorf("Variable a should be 10, got %v (exists=%v)", val, exists)
	}
	if val, exists := r.vars.GetVariable("b"); !exists || val != 3 {
		t.Errorf("Variable b should be 3, got %v (exists=%v)", val, exists)
	}
	if val, exists := r.vars.GetVariable("c"); !exists || val != 2 {
		t.Errorf("Variable c should be 2, got %v (exists=%v)", val, exists)
	}

	// Use variables in expression
	result, err = r.ParseAndEvaluate("a b + c *")
	if err != nil {
		t.Fatalf("Expression with variables failed: %v", err)
	}
	if result != "26" {
		t.Errorf("a b + c * = %q, want %q", result, "26")
	}
}

// TestVarsCommand tests the vars command.
func TestVarsCommand(t *testing.T) {
	r := NewRPN(NewVariables())

	// Empty vars
	result, err := r.ParseAndEvaluate("vars")
	if err != nil {
		t.Fatalf("vars failed: %v", err)
	}
	if !strings.Contains(result, "No variables defined") {
		t.Errorf("Empty vars should say 'No variables defined', got: %q", result)
	}

	// With variables
	r.ParseAndEvaluate("z 3 :=")
	r.ParseAndEvaluate("a 1 :=")
	r.ParseAndEvaluate("m 2 :=")

	result, err = r.ParseAndEvaluate("vars")
	if err != nil {
		t.Fatalf("vars with data failed: %v", err)
	}

	// Should contain all variable names
	if !strings.Contains(result, "a") || !strings.Contains(result, "m") || !strings.Contains(result, "z") {
		t.Errorf("vars should contain a, m, z; got: %q", result)
	}
}

// TestClearCommand tests the clear command.
func TestClearCommand(t *testing.T) {
	r := NewRPN(NewVariables())
	r.ParseAndEvaluate("x 5 :=")
	r.ParseAndEvaluate("y 10 :=")

	result, err := r.ParseAndEvaluate("clear")
	if err != nil {
		t.Fatalf("clear failed: %v", err)
	}
	if !strings.Contains(result, "All variables cleared") {
		t.Errorf("clear result = %q, want to contain 'All variables cleared'", result)
	}

	// Verify variables are gone
	if r.vars.Count() != 0 {
		t.Errorf("After clear, count = %d, want 0", r.vars.Count())
	}
}

// TestDeleteVariableOperator tests the d (delete) operator.
func TestDeleteVariableOperator(t *testing.T) {
	r := NewRPN(NewVariables())

	// Create variable
	r.ParseAndEvaluate("x 5 :=")
	_, exists := r.vars.GetVariable("x")
	if !exists {
		t.Fatal("Variable x should exist before delete")
	}

	// Delete with symbol syntax — d is a side-effect operator that leaves
	// the stack empty, so ParseAndEvaluate returns "empty result". The delete
	// still succeeds; we verify by checking the variable store directly.
	_, _ = r.ParseAndEvaluate(":x d")

	// Verify it's gone
	_, exists = r.vars.GetVariable("x")
	if exists {
		t.Error("Variable x should not exist after delete")
	}
}

// TestDeleteNonExistentVariableOperator tests deleting a variable that doesn't exist.
func TestDeleteNonExistentVariableOperator(t *testing.T) {
	r := NewRPN(NewVariables())
	_, err := r.ParseAndEvaluate(":nonexistent d")
	if err == nil {
		t.Error("Deleting non-existent variable should return error")
	}
}

// TestVariablePersistenceAcrossExpressions tests that variables persist across
// multiple ParseAndEvaluate calls on the same RPN instance.
func TestVariablePersistenceAcrossExpressions(t *testing.T) {
	r := NewRPN(NewVariables())

	r.ParseAndEvaluate("a 10 :=")
	r.ParseAndEvaluate("b 3 :=")

	// Variables should persist
	result, err := r.ParseAndEvaluate("a b +")
	if err != nil {
		t.Fatalf("Expression failed: %v", err)
	}
	if result != "13" {
		t.Errorf("a b + = %q, want %q", result, "13")
	}

	// Verify vars shows both
	result, _ = r.ParseAndEvaluate("vars")
	if !strings.Contains(result, "a") || !strings.Contains(result, "b") {
		t.Errorf("vars should contain a and b, got: %q", result)
	}
}
