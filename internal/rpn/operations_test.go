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
	s.Push(NewNumberValue(1.0))
	s.Push(NewNumberValue(2.0))
	s.Push(NewNumberValue(3.0))

	if s.Len() != 3 {
		t.Errorf("Length after 3 pushes = %d, want 3", s.Len())
	}

	val, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop() returned error: %v", err)
	}
	if val.Number() != 3.0 {
		t.Errorf("Pop() = %v, want 3.0", val)
	}

	if s.Len() != 2 {
		t.Errorf("Length after pop = %d, want 2", s.Len())
	}
}

func TestStackPeek(t *testing.T) {
	s := NewStack()
	s.Push(NewNumberValue(5.0))

	val, err := s.Peek()
	if err != nil {
		t.Fatalf("Peek() returned error: %v", err)
	}
	if val.Number() != 5.0 {
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
	s.Push(NewNumberValue(1.0))
	s.Push(NewNumberValue(2.0))
	s.Push(NewNumberValue(3.0))

	vals := s.Values()
	if len(vals) != 3 {
		t.Errorf("Values() length = %d, want 3", len(vals))
	}

	// Values() returns values in storage order (bottom-to-top)
	// Push order: 1, 2, 3 so storage is [1, 2, 3] with 3 on top
	if vals[0].Number() != 1.0 || vals[1].Number() != 2.0 || vals[2].Number() != 3.0 {
		t.Errorf("Values() = %v, want [1 2 3] (bottom-to-top)", vals)
	}
}

func TestStackClear(t *testing.T) {
	s := NewStack()
	s.Push(NewNumberValue(1.0))
	s.Push(NewNumberValue(2.0))
	s.Push(NewNumberValue(3.0))

	s.Clear()

	if s.Len() != 0 {
		t.Errorf("Length after Clear() = %d, want 0", s.Len())
	}
}

func TestOperationsAdd(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumberValue(3.0))
	s.Push(NewNumberValue(4.0))

	err := o.Add(s)
	if err != nil {
		t.Fatalf("Add() returned error: %v", err)
	}

	val, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop() after Add() returned error: %v", err)
	}
	if val.Number() != 7.0 {
		t.Errorf("Add result = %v, want 7.0", val)
	}
}

func TestOperationsSubtract(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumberValue(10.0))
	s.Push(NewNumberValue(4.0))

	err := o.Subtract(s)
	if err != nil {
		t.Fatalf("Subtract() returned error: %v", err)
	}

	val, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop() after Subtract() returned error: %v", err)
	}
	if val.Number() != 6.0 {
		t.Errorf("Subtract result = %v, want 6.0 (10 - 4)", val)
	}
}

func TestOperationsMultiply(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumberValue(5.0))
	s.Push(NewNumberValue(3.0))

	err := o.Multiply(s)
	if err != nil {
		t.Fatalf("Multiply() returned error: %v", err)
	}

	val, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop() after Multiply() returned error: %v", err)
	}
	if val.Number() != 15.0 {
		t.Errorf("Multiply result = %v, want 15.0", val)
	}
}

func TestOperationsDivide(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumberValue(20.0))
	s.Push(NewNumberValue(4.0))

	err := o.Divide(s)
	if err != nil {
		t.Fatalf("Divide() returned error: %v", err)
	}

	val, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop() after Divide() returned error: %v", err)
	}
	if val.Number() != 5.0 {
		t.Errorf("Divide result = %v, want 5.0", val)
	}
}

func TestOperationsDivideByZero(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumberValue(10.0))
	s.Push(NewNumberValue(0.0))

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
	s.Push(NewNumberValue(2.0))
	s.Push(NewNumberValue(3.0))

	err := o.Power(s)
	if err != nil {
		t.Fatalf("Power() returned error: %v", err)
	}

	val, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop() after Power() returned error: %v", err)
	}
	if val.Number() != 8.0 {
		t.Errorf("Power result = %v, want 8.0 (2^3)", val)
	}
}

func TestOperationsModulo(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumberValue(10.0))
	s.Push(NewNumberValue(3.0))

	err := o.Modulo(s)
	if err != nil {
		t.Fatalf("Modulo() returned error: %v", err)
	}

	val, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop() after Modulo() returned error: %v", err)
	}
	if val.Number() != 1.0 {
		t.Errorf("Modulo result = %v, want 1.0 (10 %% 3)", val)
	}
}

func TestOperationsModuloByZero(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumberValue(10.0))
	s.Push(NewNumberValue(0.0))

	err := o.Modulo(s)
	if err == nil {
		t.Error("Modulo by zero should return error")
	}
}

func TestOperationsInsufficientOperands(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumberValue(5.0))

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
	s.Push(NewNumberValue(7.0))

	err := o.Dup(s)
	if err != nil {
		t.Fatalf("Dup() returned error: %v", err)
	}

	if s.Len() != 2 {
		t.Errorf("Length after Dup() = %d, want 2", s.Len())
	}

	val1, _ := s.Pop()
	val2, _ := s.Pop()
	if val1.Number() != 7.0 || val2.Number() != 7.0 {
		t.Errorf("Dup values = %v, %v, want both 7.0", val1, val2)
	}
}

func TestOperationsSwap(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumberValue(1.0))
	s.Push(NewNumberValue(2.0))

	err := o.Swap(s)
	if err != nil {
		t.Fatalf("Swap() returned error: %v", err)
	}

	val1, _ := s.Pop()
	val2, _ := s.Pop()
	if val1.Number() != 1.0 || val2.Number() != 2.0 {
		t.Errorf("After Swap, values = %v, %v, want 1.0, 2.0 (swapped)", val1, val2)
	}
}

func TestOperationsSwapInsufficient(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumberValue(5.0))

	err := o.Swap(s)
	if err == nil {
		t.Error("Swap with insufficient operands should return error")
	}
}

func TestOperationsPop(t *testing.T) {
	v := NewVariables()
	o := NewOperations(v)
	s := NewStack()
	s.Push(NewNumberValue(1.0))
	s.Push(NewNumberValue(2.0))
	s.Push(NewNumberValue(3.0))

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
	s.Push(NewNumberValue(1.0))
	s.Push(NewNumberValue(2.0))
	s.Push(NewNumberValue(3.0))

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
	s.Push(NewNumberValue(5.0))

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
	if val.Number() != 3.14159 {
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
			s.Push(NewNumberValue(float64(id)))
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
	stack.Push(NewNumberValue(8))
	err := o.Log2(stack)
	if err != nil {
		t.Errorf("Log2() returned error: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Number() != 3.0 {
		t.Errorf("Log2(8) = %f, want 3.0)", val.Number())
	}

	// Test log₂(1) = 0
	stack.Push(NewNumberValue(1))
	err = o.Log2(stack)
	if err != nil {
		t.Errorf("Log2(1) returned error: %v", err)
	}
	val, err = stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Number() != 0.0 {
		t.Errorf("Log2(1) = %f, want 0.0)", val.Number())
	}

	// Test log₂(0) should error
	stack.Push(NewNumberValue(0))
	err = o.Log2(stack)
	if err == nil {
		t.Errorf("Log2(0) should return error, got nil")
	}
}

func TestLog10(t *testing.T) {
	o := NewOperations(NewVariables())
	stack := NewStack()

	// Test log₁₀(100) = 2
	stack.Push(NewNumberValue(100))
	err := o.Log10(stack)
	if err != nil {
		t.Errorf("Log10() returned error: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Number() != 2.0 {
		t.Errorf("Log10(100) = %f, want 2.0)", val.Number())
	}

	// Test log₁₀(1) = 0
	stack.Push(NewNumberValue(1))
	err = o.Log10(stack)
	if err != nil {
		t.Errorf("Log10(1) returned error: %v", err)
	}
	val, err = stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Number() != 0.0 {
		t.Errorf("Log10(1) = %f, want 0.0)", val.Number())
	}
}

func TestLn(t *testing.T) {
	o := NewOperations(NewVariables())
	stack := NewStack()

	// Test ln(e) ≈ 1
	stack.Push(NewNumberValue(math.E))
	err := o.Ln(stack)
	if err != nil {
		t.Errorf("Ln() returned error: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if math.Abs(val.Number()-1.0) > 0.0001 {
		t.Errorf("ln(e) = %f, want ~1.0", val.Number())
	}

	// Test ln(1) = 0
	stack.Push(NewNumberValue(1))
	err = o.Ln(stack)
	if err != nil {
		t.Errorf("Ln(1) returned error: %v", err)
	}
	val, err = stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Number() != 0.0 {
		t.Errorf("Ln(1) = %f, want 0.0)", val.Number())
	}
}

func TestHyperLog2(t *testing.T) {
	o := NewOperations(NewVariables())
	stack := NewStack()

	// Test hyperlog₂(4, 16) = log₂(4) + log₂(16) = 2 + 4 = 6
	stack.Push(NewNumberValue(4))
	stack.Push(NewNumberValue(16))
	err := o.HyperLog2(stack)
	if err != nil {
		t.Errorf("HyperLog2() returned error: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Number() != 6.0 {
		t.Errorf("HyperLog2(4, 16) = %f, want 6.0)", val.Number())
	}

	// Test with single value (should error, like other hyper operators)
	stack.Push(NewNumberValue(8))
	err = o.HyperLog2(stack)
	if err == nil {
		t.Errorf("HyperLog2 with single value should return error, got nil")
	}
}

func TestHyperLog10(t *testing.T) {
	o := NewOperations(NewVariables())
	stack := NewStack()

	// Test hyperlog₁₀(10, 100) = log₁₀(10) + log₁₀(100) = 1 + 2 = 3
	stack.Push(NewNumberValue(10))
	stack.Push(NewNumberValue(100))
	err := o.HyperLog10(stack)
	if err != nil {
		t.Errorf("HyperLog10() returned error: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val.Number() != 3.0 {
		t.Errorf("HyperLog10(10, 100) = %f, want 3.0)", val.Number())
	}
}

func TestHyperLn(t *testing.T) {
	o := NewOperations(NewVariables())
	stack := NewStack()

	// Test hyperln(e, e²) = ln(e) + ln(e²) = 1 + 2 = 3
	stack.Push(NewNumberValue(math.E))
	stack.Push(NewNumberValue(math.E * math.E))
	err := o.HyperLn(stack)
	if err != nil {
		t.Errorf("HyperLn() returned error: %v", err)
	}
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if math.Abs(val.Number()-3.0) > 0.0001 {
		t.Errorf("HyperLn(e, e²) = %f, want ~3.0", val.Number())
	}
}
