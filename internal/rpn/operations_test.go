// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"testing"
)

func TestStackNewStack(t *testing.T) {
	s := NewStack()
	if s == nil {
		t.Fatal("NewStack() returned nil")
	}
	if s.Len() != 0 {
		t.Errorf("NewStack() length = %d, want 0", s.Len())
	}
}

func TestStackPushPop(t *testing.T) {
	s := NewStack()
	s.Push(NewNumber(1.0, FloatMode))
	s.Push(NewNumber(2.0, FloatMode))
	s.Push(NewNumber(3.0, FloatMode))

	if s.Len() != 3 {
		t.Errorf("Length after 3 pushes = %d, want 3", s.Len())
	}

	val, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop() returned error: %v", err)
	}
	if val.Float64() != 3.0 {
		t.Errorf("Pop() = %v, want 3.0", val)
	}

	if s.Len() != 2 {
		t.Errorf("Length after pop = %d, want 2", s.Len())
	}
}

func TestStackPeek(t *testing.T) {
	s := NewStack()
	s.Push(NewNumber(5.0, FloatMode))

	val, err := s.Peek()
	if err != nil {
		t.Fatalf("Peek() returned error: %v", err)
	}
	if val.Float64() != 5.0 {
		t.Errorf("Peek() = %v, want 5.0", val)
	}

	// Peek should not remove the value
	if s.Len() != 1 {
		t.Errorf("Length after Peek() = %d, want 1", s.Len())
	}
}

func TestStackPeekEmpty(t *testing.T) {
	s := NewStack()
	_, err := s.Peek()
	if err == nil {
		t.Error("Peek() on empty stack should return error")
	}
	if !strings.Contains(err.Error(), "stack is empty") {
		t.Errorf("Peek() error = %v, should contain 'stack is empty'", err)
	}
}

func TestStackPopEmpty(t *testing.T) {
	s := NewStack()
	_, err := s.Pop()
	if err == nil {
		t.Error("Pop() on empty stack should return error")
	}
}

func TestStackValues(t *testing.T) {
	s := NewStack()
	s.Push(NewNumber(1.0, FloatMode))
	s.Push(NewNumber(2.0, FloatMode))
	s.Push(NewNumber(3.0, FloatMode))

	vals := s.Values()
	if len(vals) != 3 {
		t.Errorf("Values() length = %d, want 3", len(vals))
	}

	// Values() returns values in storage order (bottom-to-top)
	// Push order: 1, 2, 3 so storage is [1, 2, 3] with 3 on top
	if vals[0].Float64() != 1.0 || vals[1].Float64() != 2.0 || vals[2].Float64() != 3.0 {
		t.Errorf("Values() = %v, want [1 2 3] (bottom-to-top)", vals)
	}
}

func TestStackClear(t *testing.T) {
	s := NewStack()
	s.Push(NewNumber(1.0, FloatMode))
	s.Push(NewNumber(2.0, FloatMode))
	s.Push(NewNumber(3.0, FloatMode))

	s.Clear()

	if s.Len() != 0 {
		t.Errorf("Length after Clear() = %d, want 0", s.Len())
	}
}

func TestOperationsAdd(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumber(3.0, FloatMode))
	s.Push(NewNumber(4.0, FloatMode))

	err := o.Add(s)
	if err != nil {
		t.Fatalf("Add() returned error: %v", err)
	}

	val, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop() after Add() returned error: %v", err)
	}
	if val.Float64() != 7.0 {
		t.Errorf("Add result = %v, want 7.0", val)
	}
}

func TestOperationsSubtract(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumber(10.0, FloatMode))
	s.Push(NewNumber(4.0, FloatMode))

	err := o.Subtract(s)
	if err != nil {
		t.Fatalf("Subtract() returned error: %v", err)
	}

	val, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop() after Subtract() returned error: %v", err)
	}
	if val.Float64() != 6.0 {
		t.Errorf("Subtract result = %v, want 6.0 (10 - 4)", val)
	}
}

func TestOperationsMultiply(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumber(5.0, FloatMode))
	s.Push(NewNumber(3.0, FloatMode))

	err := o.Multiply(s)
	if err != nil {
		t.Fatalf("Multiply() returned error: %v", err)
	}

	val, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop() after Multiply() returned error: %v", err)
	}
	if val.Float64() != 15.0 {
		t.Errorf("Multiply result = %v, want 15.0", val)
	}
}

func TestOperationsDivide(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumber(20.0, FloatMode))
	s.Push(NewNumber(4.0, FloatMode))

	err := o.Divide(s)
	if err != nil {
		t.Fatalf("Divide() returned error: %v", err)
	}

	val, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop() after Divide() returned error: %v", err)
	}
	if val.Float64() != 5.0 {
		t.Errorf("Divide result = %v, want 5.0", val)
	}
}

func TestOperationsDivideByZero(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumber(10.0, FloatMode))
	s.Push(NewNumber(0.0, FloatMode))

	err := o.Divide(s)
	if err == nil {
		t.Error("Divide by zero should return error")
	}
	if !strings.Contains(err.Error(), "division by zero") {
		t.Errorf("Divide by zero error = %v, should contain 'division by zero'", err)
	}
}

func TestOperationsPower(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumber(2.0, FloatMode))
	s.Push(NewNumber(3.0, FloatMode))

	err := o.Power(s)
	if err != nil {
		t.Fatalf("Power() returned error: %v", err)
	}

	val, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop() after Power() returned error: %v", err)
	}
	if val.Float64() != 8.0 {
		t.Errorf("Power result = %v, want 8.0 (2^3)", val)
	}
}

func TestOperationsModulo(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumber(10.0, FloatMode))
	s.Push(NewNumber(3.0, FloatMode))

	err := o.Modulo(s)
	if err != nil {
		t.Fatalf("Modulo() returned error: %v", err)
	}

	val, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop() after Modulo() returned error: %v", err)
	}
	if val.Float64() != 1.0 {
		t.Errorf("Modulo result = %v, want 1.0 (10 %% 3)", val)
	}
}

func TestOperationsModuloByZero(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumber(10.0, FloatMode))
	s.Push(NewNumber(0.0, FloatMode))

	err := o.Modulo(s)
	if err == nil {
		t.Error("Modulo by zero should return error")
	}
}

func TestOperationsInsufficientOperands(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumber(5.0, FloatMode))

	// Try to add with only one operand
	err := o.Add(s)
	if err == nil {
		t.Error("Add with insufficient operands should return error")
	}
}

func TestOperationsDup(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumber(7.0, FloatMode))

	err := o.Dup(s)
	if err != nil {
		t.Fatalf("Dup() returned error: %v", err)
	}

	if s.Len() != 2 {
		t.Errorf("Length after Dup() = %d, want 2", s.Len())
	}

	val1, _ := s.Pop()
	val2, _ := s.Pop()
	if val1.Float64() != 7.0 || val2.Float64() != 7.0 {
		t.Errorf("Dup values = %v, %v, want both 7.0", val1, val2)
	}
}

func TestOperationsSwap(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumber(1.0, FloatMode))
	s.Push(NewNumber(2.0, FloatMode))

	err := o.Swap(s)
	if err != nil {
		t.Fatalf("Swap() returned error: %v", err)
	}

	val1, _ := s.Pop()
	val2, _ := s.Pop()
	if val1.Float64() != 1.0 || val2.Float64() != 2.0 {
		t.Errorf("After Swap, values = %v, %v, want 1.0, 2.0 (swapped)", val1, val2)
	}
}

func TestOperationsSwapInsufficient(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumber(5.0, FloatMode))

	err := o.Swap(s)
	if err == nil {
		t.Error("Swap with insufficient operands should return error")
	}
}

func TestOperationsPop(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumber(1.0, FloatMode))
	s.Push(NewNumber(2.0, FloatMode))
	s.Push(NewNumber(3.0, FloatMode))

	err := o.Pop(s)
	if err != nil {
		t.Fatalf("Pop() returned error: %v", err)
	}

	if s.Len() != 2 {
		t.Errorf("Length after Pop() = %d, want 2", s.Len())
	}
}

func TestOperationsPopEmpty(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()

	err := o.Pop(s)
	if err == nil {
		t.Error("Pop on empty stack should return error")
	}
}

func TestOperationsShow(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumber(1.0, FloatMode))
	s.Push(NewNumber(2.0, FloatMode))
	s.Push(NewNumber(3.0, FloatMode))

	result, err := o.Show(s)
	if err != nil {
		t.Fatalf("Show() returned error: %v", err)
	}

	if result != "1 2 3" {
		t.Errorf("Show() = %q, want \"1 2 3\"", result)
	}
}

func TestOperationsShowEmpty(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()

	result, err := o.Show(s)
	if err != nil {
		t.Fatalf("Show() on empty stack returned error: %v", err)
	}

	if !strings.Contains(result, "Stack is empty") {
		t.Errorf("Show() on empty stack = %q, should contain 'Stack is empty'", result)
	}
}

func TestOperationsAssignVariable(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumber(5.0, FloatMode))

	err := o.AssignVariable(s, "x")
	if err != nil {
		t.Fatalf("AssignVariable() returned error: %v", err)
	}

	val, exists := v.GetVariable("x")
	if !exists {
		t.Error("Variable x should exist after assignment")
	}
	if val != 5.0 {
		t.Errorf("Variable x value = %v, want 5.0", val)
	}

	// Verify value was popped from stack
	if s.Len() != 0 {
		t.Errorf("Stack length after assignment = %d, want 0", s.Len())
	}
}

func TestOperationsAssignVariableEmptyName(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()

	err := o.AssignVariable(s, "")
	if err == nil {
		t.Error("AssignVariable with empty name should return error")
	}
}

func TestOperationsUseVariable(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()

	if err := v.SetVariable("pi", 3.14159); err != nil {
		t.Fatalf("SetVariable() returned error: %v", err)
	}

	err := o.UseVariable(s, "pi")
	if err != nil {
		t.Fatalf("UseVariable() returned error: %v", err)
	}

	val, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop() after UseVariable() returned error: %v", err)
	}
	if val.Float64() != 3.14159 {
		t.Errorf("Variable value pushed to stack = %v, want 3.14159", val)
	}
}

func TestOperationsUseVariableUndefined(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()

	err := o.UseVariable(s, "undefined")
	if err == nil {
		t.Error("UseVariable for undefined variable should return error")
	}
	if !errors.Is(err, ErrVariableNotFound) {
		t.Errorf("UseVariable error = %v, should be ErrVariableNotFound", err)
	}
}

func TestOperationsDeleteVariable(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)

	if err := v.SetVariable("temp", 100.0); err != nil {
		t.Fatalf("SetVariable() returned error: %v", err)
	}

	err := o.DeleteVariable("temp")
	if err != nil {
		t.Fatalf("DeleteVariable() returned error: %v", err)
	}

	_, exists := v.GetVariable("temp")
	if exists {
		t.Error("Variable should not exist after deletion")
	}
}

func TestOperationsDeleteVariableUndefined(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)

	err := o.DeleteVariable("nonexistent")
	if err == nil {
		t.Error("DeleteVariable for undefined variable should return error")
	}
}

func TestOperationsListVariables(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)

	if err := v.SetVariable("x", 1.0); err != nil {
		t.Fatalf("SetVariable() returned error: %v", err)
	}
	if err := v.SetVariable("y", 2.0); err != nil {
		t.Fatalf("SetVariable() returned error: %v", err)
	}

	result, err := o.ListVariables()
	if err != nil {
		t.Fatalf("ListVariables() returned error: %v", err)
	}

	if strings.Contains(result, "No variables defined") {
		t.Error("ListVariables should show variables, not 'No variables defined'")
	}
	if !strings.Contains(result, "x") || !strings.Contains(result, "y") {
		t.Errorf("ListVariables output should contain all variable names, got: %s", result)
	}
}

func TestOperationsClearVariables(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)

	_, _ = v.SetVariable("x", 1.0), v.SetVariable("y", 2.0)

	o.ClearVariables()

	if v.Count() != 0 {
		t.Errorf("Count after ClearVariables() = %d, want 0", v.Count())
	}
}

func TestOperationsConcurrent(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)

	// Test concurrent variable access
	// Each goroutine uses its own stack to avoid race conditions
	done := make(chan bool, 10)
	for i := 0; i < 5; i++ {
		go func(id int) {
			name := fmt.Sprintf("concurrent%d", id)
			s := NewStack()
			s.Push(NewNumber(float64(id), FloatMode))
			if err := o.AssignVariable(s, name); err != nil {
				t.Errorf("AssignVariable() returned error: %v", err)
			}
			done <- true
		}(i)
	}

	for i := 0; i < 5; i++ {
		<-done
	}

	if v.Count() != 5 {
		t.Errorf("Final count = %d, want 5", v.Count())
	}
}

func TestLog2(t *testing.T) {
	o := NewOperations(NewVariables())
	stack := NewStack()

	// Test log₂(8) = 3
	stack.Push(NewNumber(8, FloatMode))
	err := o.Log2(stack)
	if err != nil {
		t.Errorf("Log2() returned error: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Float64() != 3.0 {
		t.Errorf("Log2(8) = %f, want 3.0)", val.Float64())
	}

	// Test log₂(1) = 0
	stack.Push(NewNumber(1.0, FloatMode))
	err = o.Log2(stack)
	if err != nil {
		t.Errorf("Log2(1) returned error: %v", err)
	}
	val, err = stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Float64() != 0.0 {
		t.Errorf("Log2(1) = %f, want 0.0)", val.Float64())
	}

	// Test log₂(0) should error
	stack.Push(NewNumber(0.0, FloatMode))
	err = o.Log2(stack)
	if err == nil {
		t.Errorf("Log2(0) should return error, got nil")
	}
}

func TestLog10(t *testing.T) {
	o := NewOperations(NewVariables())
	stack := NewStack()

	// Test log₁₀(100) = 2
	stack.Push(NewNumber(100.0, FloatMode))
	err := o.Log10(stack)
	if err != nil {
		t.Errorf("Log10() returned error: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Float64() != 2.0 {
		t.Errorf("Log10(100) = %f, want 2.0)", val.Float64())
	}

	// Test log₁₀(1) = 0
	stack.Push(NewNumber(1.0, FloatMode))
	err = o.Log10(stack)
	if err != nil {
		t.Errorf("Log10(1) returned error: %v", err)
	}
	val, err = stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Float64() != 0.0 {
		t.Errorf("Log10(1) = %f, want 0.0)", val.Float64())
	}
}

func TestLn(t *testing.T) {
	o := NewOperations(NewVariables())
	stack := NewStack()

	// Test ln(e) ≈ 1
	stack.Push(NewNumber(math.E, FloatMode))
	err := o.Ln(stack)
	if err != nil {
		t.Errorf("Ln() returned error: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if math.Abs(val.Float64()-1.0) > 0.0001 {
		t.Errorf("ln(e) = %f, want ~1.0", val.Float64())
	}

	// Test ln(1) = 0
	stack.Push(NewNumber(1.0, FloatMode))
	err = o.Ln(stack)
	if err != nil {
		t.Errorf("Ln(1) returned error: %v", err)
	}
	val, err = stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Float64() != 0.0 {
		t.Errorf("Ln(1) = %f, want 0.0)", val.Float64())
	}
}

func TestLog2WithBoolean(t *testing.T) {
	o := NewOperations(NewVariables())
	stack := NewStack()

	// Test with boolean true (should be converted to 1, log₂(1) = 0)
	stack.Push(NewFloatFromBool(true))
	err := o.Log2(stack)
	if err != nil {
		t.Errorf("Log2(true) returned error: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Float64() != 0.0 {
		t.Errorf("Log2(true) = %f, want 0.0 (log₂(1) = 0)", val.Float64())
	}

	// Test with boolean false (should be converted to 0, log₂(0) should error)
	stack.Push(NewFloatFromBool(false))
	err = o.Log2(stack)
	if err == nil {
		t.Errorf("Log2(false) should return error for log₂(0), got nil")
	}
}

func TestLog10WithBoolean(t *testing.T) {
	o := NewOperations(NewVariables())
	stack := NewStack()

	// Test with boolean true (should be converted to 1, log₁₀(1) = 0)
	stack.Push(NewFloatFromBool(true))
	err := o.Log10(stack)
	if err != nil {
		t.Errorf("Log10(true) returned error: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Float64() != 0.0 {
		t.Errorf("Log10(true) = %f, want 0.0 (log₁₀(1) = 0)", val.Float64())
	}

	// Test with boolean false (should be converted to 0, log₁₀(0) should error)
	stack.Push(NewFloatFromBool(false))
	err = o.Log10(stack)
	if err == nil {
		t.Errorf("Log10(false) should return error for log₁₀(0), got nil")
	}
}

func TestLnWithBoolean(t *testing.T) {
	o := NewOperations(NewVariables())
	stack := NewStack()

	// Test with boolean true (should be converted to 1, ln(1) = 0)
	stack.Push(NewFloatFromBool(true))
	err := o.Ln(stack)
	if err != nil {
		t.Errorf("Ln(true) returned error: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Float64() != 0.0 {
		t.Errorf("Ln(true) = %f, want 0.0 (ln(1) = 0)", val.Float64())
	}

	// Test with boolean false (should be converted to 0, ln(0) should error)
	stack.Push(NewFloatFromBool(false))
	err = o.Ln(stack)
	if err == nil {
		t.Errorf("Ln(false) should return error for ln(0), got nil")
	}
}

func TestLnEdgeCases(t *testing.T) {
	o := NewOperations(NewVariables())
	stack := NewStack()

	// Test ln(negative) should error
	stack.Push(NewNumber(-1.0, FloatMode))
	err := o.Ln(stack)
	if err == nil {
		t.Errorf("Ln(-1) should return error, got nil")
	}

	// Test ln(0) should error
	stack.Push(NewNumber(0.0, FloatMode))
	err = o.Ln(stack)
	if err == nil {
		t.Errorf("Ln(0) should return error, got nil")
	}

	// Test ln(very small positive) should work
	stack.Push(NewNumber(0.001, FloatMode))
	err = o.Ln(stack)
	if err != nil {
		t.Errorf("Ln(0.001) should not return error, got: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Float64() > -6.0 || val.Float64() < -7.0 {
		t.Errorf("Ln(0.001) = %f, want ~-6.9 (ln(0.001))", val.Float64())
	}
}

func TestHyperLog2WithBoolean(t *testing.T) {
	o := NewOperations(NewVariables())
	stack := NewStack()

	// Test hyperlog₂(4, true) = log₂(4) + log₂(1) = 2 + 0 = 2
	// true should be converted to 1
	stack.Push(NewNumber(4.0, FloatMode))
	stack.Push(NewFloatFromBool(true))
	err := o.HyperLog2(stack)
	if err != nil {
		t.Errorf("HyperLog2(4, true) returned error: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Float64() != 2.0 {
		t.Errorf("HyperLog2(4, true) = %f, want 2.0 (log₂(4) + log₂(1) = 2 + 0)", val.Float64())
	}

	// Test hyperlog₂(4, false) = log₂(4) + log₂(0) should error
	// false should be converted to 0, which is undefined for log₂
	stack.Push(NewNumber(4.0, FloatMode))
	stack.Push(NewFloatFromBool(false))
	err = o.HyperLog2(stack)
	if err == nil {
		t.Errorf("HyperLog2(4, false) should return error for log₂(0), got nil")
	}
}

func TestHyperLog10WithBoolean(t *testing.T) {
	o := NewOperations(NewVariables())
	stack := NewStack()

	// Test hyperlog₁₀(10, true) = log₁₀(10) + log₁₀(1) = 1 + 0 = 1
	// true should be converted to 1
	stack.Push(NewNumber(10.0, FloatMode))
	stack.Push(NewFloatFromBool(true))
	err := o.HyperLog10(stack)
	if err != nil {
		t.Errorf("HyperLog10(10, true) returned error: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Float64() != 1.0 {
		t.Errorf("HyperLog10(10, true) = %f, want 1.0 (log₁₀(10) + log₁₀(1) = 1 + 0)", val.Float64())
	}

	// Test hyperlog₁₀(10, false) = log₁₀(10) + log₁₀(0) should error
	stack.Push(NewNumber(10.0, FloatMode))
	stack.Push(NewFloatFromBool(false))
	err = o.HyperLog10(stack)
	if err == nil {
		t.Errorf("HyperLog10(10, false) should return error for log₁₀(0), got nil")
	}
}

func TestHyperLnWithBoolean(t *testing.T) {
	o := NewOperations(NewVariables())
	stack := NewStack()

	// Test hyperln(e, true) = ln(e) + ln(1) = 1 + 0 = 1
	// true should be converted to 1
	stack.Push(NewNumber(math.E, FloatMode))
	stack.Push(NewFloatFromBool(true))
	err := o.HyperLn(stack)
	if err != nil {
		t.Errorf("HyperLn(e, true) returned error: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if math.Abs(val.Float64()-1.0) > 0.0001 {
		t.Errorf("HyperLn(e, true) = %f, want ~1.0 (ln(e) + ln(1) = 1 + 0)", val.Float64())
	}

	// Test hyperln(e, false) = ln(e) + ln(0) should error
	stack.Push(NewNumber(math.E, FloatMode))
	stack.Push(NewFloatFromBool(false))
	err = o.HyperLn(stack)
	if err == nil {
		t.Errorf("HyperLn(e, false) should return error for ln(0), got nil")
	}
}

func TestHyperLog2(t *testing.T) {
	o := NewOperations(NewVariables())
	stack := NewStack()

	// Test hyperlog₂(4, 16) = log₂(4) + log₂(16) = 2 + 4 = 6
	stack.Push(NewNumber(4.0, FloatMode))
	stack.Push(NewNumber(16, FloatMode))
	err := o.HyperLog2(stack)
	if err != nil {
		t.Errorf("HyperLog2() returned error: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Float64() != 6.0 {
		t.Errorf("HyperLog2(4, 16) = %f, want 6.0)", val.Float64())
	}

	// Test with single value (should error, like other hyper operators)
	stack.Push(NewNumber(8, FloatMode))
	err = o.HyperLog2(stack)
	if err == nil {
		t.Errorf("HyperLog2 with single value should return error, got nil")
	}
}

func TestHyperLog10(t *testing.T) {
	o := NewOperations(NewVariables())
	stack := NewStack()

	// Test hyperlog₁₀(10, 100) = log₁₀(10) + log₁₀(100) = 1 + 2 = 3
	stack.Push(NewNumber(10.0, FloatMode))
	stack.Push(NewNumber(100.0, FloatMode))
	err := o.HyperLog10(stack)
	if err != nil {
		t.Errorf("HyperLog10() returned error: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Float64() != 3.0 {
		t.Errorf("HyperLog10(10, 100) = %f, want 3.0)", val.Float64())
	}
}

func TestHyperLn(t *testing.T) {
	o := NewOperations(NewVariables())
	stack := NewStack()

	// Test hyperln(e, e²) = ln(e) + ln(e²) = 1 + 2 = 3
	stack.Push(NewNumber(math.E, FloatMode))
	stack.Push(NewNumber(math.E*math.E, FloatMode))
	err := o.HyperLn(stack)
	if err != nil {
		t.Errorf("HyperLn() returned error: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if math.Abs(val.Float64()-3.0) > 0.0001 {
		t.Errorf("HyperLn(e, e²) = %f, want ~3.0", val.Float64())
	}
}

func TestOperatorRegistry(t *testing.T) {
	o := NewOperations(NewVariables())
	registry := NewOperatorRegistry(o)

	// Test IsStandardOperator with valid operators
	validOperators := []string{"+", "-", "*", "/", "^", "%", "lg", "log", "ln", "gt", "lt", "gte", "lte", "eq", "neq", "dup", "swap", "pop", "show", "showstack", "print", "vars", "clear"}
	for _, op := range validOperators {
		if !registry.IsStandardOperator(op) {
			t.Errorf("IsStandardOperator(%q) = false, want true", op)
		}
	}

	// Test IsStandardOperator with invalid operators
	invalidOperators := []string{"invalid", "xyz", "123"}
	for _, op := range invalidOperators {
		if registry.IsStandardOperator(op) {
			t.Errorf("IsStandardOperator(%q) = true, want false", op)
		}
	}

	// Test IsHyperOperator with valid operators
	hyperOperators := []string{"[+]", "[-]", "[*]", "[/]", "[^]", "[%]", "[lg]", "[log]", "[ln]"}
	for _, op := range hyperOperators {
		if !registry.IsHyperOperator(op) {
			t.Errorf("IsHyperOperator(%q) = false, want true", op)
		}
	}

	// Test IsHyperOperator with invalid operators
	for _, op := range invalidOperators {
		if registry.IsHyperOperator(op) {
			t.Errorf("IsHyperOperator(%q) = true, want false", op)
		}
	}
}

func TestOperatorRegistryHandleStandardOperator(t *testing.T) {
	o := NewOperations(NewVariables())
	registry := NewOperatorRegistry(o)
	stack := NewStack()

	// Test standard operator handling
	testCases := []struct {
		name     string
		token    string
		prepare  func()
		expected float64
	}{
		{"Addition", "+", func() { stack.Push(NewNumber(3.0, FloatMode)); stack.Push(NewNumber(4.0, FloatMode)) }, 7.0},
		{"Subtraction", "-", func() { stack.Push(NewNumber(10.0, FloatMode)); stack.Push(NewNumber(4.0, FloatMode)) }, 6.0},
		{"Multiplication", "*", func() { stack.Push(NewNumber(5.0, FloatMode)); stack.Push(NewNumber(3.0, FloatMode)) }, 15.0},
		{"Division", "/", func() { stack.Push(NewNumber(20.0, FloatMode)); stack.Push(NewNumber(4.0, FloatMode)) }, 5.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare()
			result, handled, err := registry.HandleStandardOperator(stack, tc.token)
			if err != nil {
				t.Errorf("HandleStandardOperator(%q) returned error: %v", tc.token, err)
			}
			if !handled {
				t.Errorf("HandleStandardOperator(%q) = false, want true", tc.token)
			}
			if result != "" {
				t.Errorf("HandleStandardOperator(%q) returned non-empty result: %q", tc.token, result)
			}
			val, err := stack.Pop()
			if err != nil {
				t.Errorf("Pop() returned error: %v", err)
			}
			if val.Float64() != tc.expected {
				t.Errorf("Result = %f, want %f", val.Float64(), tc.expected)
			}
		})
	}
}


func TestAssignLeft(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()

	// For "5 x =:":
	// AssignLeft pops name first, then value
	// Push value first (will be popped second), then name (will be popped first)
	s.Push(NewNumber(5, FloatMode))  // value (will be popped second)
	s.Push(NewStringNum("x"))  // name (will be popped first)

	err := o.AssignLeft(s)
	if err != nil {
		t.Errorf("AssignLeft() error = %v", err)
	}

	// Check that x = 5
	val, exists := v.GetVariable("x")
	if !exists {
		t.Errorf("Variable x should exist after assignment")
	}
	if val != 5 {
		t.Errorf("Variable x = %v, want 5", val)
	}

	// Stack should be empty
	if s.Len() != 0 {
		t.Errorf("Stack length = %d, want 0", s.Len())
	}
}

func TestAssignRight(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()

	// For "x 5 :=":
	// AssignRight pops value first, then name
	// Push name first (will be popped second), then value (will be popped first)
	s.Push(NewStringNum("x"))  // name (will be popped second)
	s.Push(NewNumber(5, FloatMode))  // value (will be popped first)

	err := o.AssignRight(s)
	if err != nil {
		t.Errorf("AssignRight() error = %v", err)
	}

	// Check that x = 5
	val, exists := v.GetVariable("x")
	if !exists {
		t.Errorf("Variable x should exist after assignment")
	}
	if val != 5 {
		t.Errorf("Variable x = %v, want 5", val)
	}

	// Stack should be empty
	if s.Len() != 0 {
		t.Errorf("Stack length = %d, want 0", s.Len())
	}
}

func TestAssignLeftErrorCases(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()

	err := o.AssignLeft(s)
	if err == nil {
		t.Error("AssignLeft() should return error when stack is empty")
	}
}

func TestAssignRightErrorCases(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()

	err := o.AssignRight(s)
	if err == nil {
		t.Error("AssignRight() should return error when stack is empty")
	}
}
