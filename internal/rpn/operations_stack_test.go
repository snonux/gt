// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"strings"
	"testing"
)

func TestShowWithMetrics(t *testing.T) {
	reg := GetMetricRegistry()
	vars := NewVariables()
	ops := NewOperations(vars, nil)
	stack := NewStack()

	mbps, _ := reg.Find("Mbps")
	cool, _ := reg.Find("Cool")

	stack.Push(NewFloatWithMetric(100, mbps))
	stack.Push(NewFloatWithMetric(42, cool))
	stack.Push(NewFloatWithMetric(5.5, mbps))

	result, err := ops.Show(stack)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Non-Cool shows metric suffix, Cool doesn't
	if !strings.Contains(result, "100Mbps") {
		t.Errorf("expected '100Mbps' in result, got: %s", result)
	}
	if !strings.Contains(result, "42") {
		t.Errorf("expected '42' (plain) in result, got: %s", result)
	}
	if !strings.Contains(result, "5.5Mbps") {
		t.Errorf("expected '5.5Mbps' in result, got: %s", result)
	}
}

func TestShowEmptyStack(t *testing.T) {
	vars := NewVariables()
	ops := NewOperations(vars, nil)
	stack := NewStack()

	result, err := ops.Show(stack)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Stack is empty" {
		t.Errorf("Show(empty) = %q, want 'Stack is empty'", result)
	}
}

func TestShowWithBooleanValues(t *testing.T) {
	vars := NewVariables()
	ops := NewOperations(vars, nil)
	stack := NewStack()

	stack.Push(NewFloatFromBool(true))
	stack.Push(NewFloatFromBool(false))

	result, err := ops.Show(stack)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "true") {
		t.Errorf("expected 'true' in result, got: %s", result)
	}
	if !strings.Contains(result, "false") {
		t.Errorf("expected 'false' in result, got: %s", result)
	}
}

func TestShowWithSymbols(t *testing.T) {
	vars := NewVariables()
	ops := NewOperations(vars, nil)
	stack := NewStack()

	stack.Push(NewSymbol("x"))
	stack.Push(NewSymbol("counter"))
	stack.Push(NewFloat(42))

	result, err := ops.Show(stack)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, ":x") {
		t.Errorf("expected ':x' in result, got: %s", result)
	}
	if !strings.Contains(result, ":counter") {
		t.Errorf("expected ':counter' in result, got: %s", result)
	}
	if !strings.Contains(result, "42") {
		t.Errorf("expected '42' in result, got: %s", result)
	}
}

func TestShowWithStringNum(t *testing.T) {
	vars := NewVariables()
	ops := NewOperations(vars, nil)
	stack := NewStack()

	stack.Push(NewStringNum("hello"))
	stack.Push(NewFloat(10))

	result, err := ops.Show(stack)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "hello") {
		t.Errorf("expected 'hello' in result, got: %s", result)
	}
	if !strings.Contains(result, "10") {
		t.Errorf("expected '10' in result, got: %s", result)
	}
}

func TestShowWithMixedTypes(t *testing.T) {
	vars := NewVariables()
	ops := NewOperations(vars, nil)
	stack := NewStack()
	reg := GetMetricRegistry()

	mbps, _ := reg.Find("Mbps")

	// Push various types: number, bool, symbol, string, metric number
	stack.Push(NewFloat(42))
	stack.Push(NewFloatFromBool(true))
	stack.Push(NewSymbol("x"))
	stack.Push(NewStringNum("hello"))
	stack.Push(NewFloatWithMetric(100, mbps))

	result, err := ops.Show(stack)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify all types appear in output
	if !strings.Contains(result, "42") {
		t.Errorf("expected '42' in result, got: %s", result)
	}
	if !strings.Contains(result, "true") {
		t.Errorf("expected 'true' in result, got: %s", result)
	}
	if !strings.Contains(result, ":x") {
		t.Errorf("expected ':x' in result, got: %s", result)
	}
	if !strings.Contains(result, "hello") {
		t.Errorf("expected 'hello' in result, got: %s", result)
	}
	if !strings.Contains(result, "100Mbps") {
		t.Errorf("expected '100Mbps' in result, got: %s", result)
	}
}

func TestShowWithMultipleMetrics(t *testing.T) {
	vars := NewVariables()
	ops := NewOperations(vars, nil)
	stack := NewStack()
	reg := GetMetricRegistry()

	mbps, _ := reg.Find("Mbps")
	hrs, _ := reg.Find("hr")

	stack.Push(NewFloatWithMetric(100, mbps))
	stack.Push(NewFloatWithMetric(50, mbps))
	stack.Push(NewFloatWithMetric(2.5, hrs))

	result, err := ops.Show(stack)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "100Mbps") {
		t.Errorf("expected '100Mbps' in result, got: %s", result)
	}
	if !strings.Contains(result, "50Mbps") {
		t.Errorf("expected '50Mbps' in result, got: %s", result)
	}
	if !strings.Contains(result, "2.5hr") {
		t.Errorf("expected '2.5hr' in result, got: %s", result)
	}
}

func TestShowWithRat(t *testing.T) {
	vars := NewVariables()
	ops := NewOperations(vars, nil)
	stack := NewStack()

	// Rational numbers
	rat := NewRat(0.5)
	stack.Push(rat)
	stack.Push(NewFloat(10))

	result, err := ops.Show(stack)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "0.5") {
		t.Errorf("expected '0.5' in result, got: %s", result)
	}
	if !strings.Contains(result, "10") {
		t.Errorf("expected '10' in result, got: %s", result)
	}
}

func TestShowWithRatFromBool(t *testing.T) {
	vars := NewVariables()
	ops := NewOperations(vars, nil)
	stack := NewStack()

	stack.Push(NewRatFromBool(true))
	stack.Push(NewRatFromBool(false))

	result, err := ops.Show(stack)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "true") {
		t.Errorf("expected 'true' in result, got: %s", result)
	}
	if !strings.Contains(result, "false") {
		t.Errorf("expected 'false' in result, got: %s", result)
	}
}

func TestShowValuesOrder(t *testing.T) {
	vars := NewVariables()
	ops := NewOperations(vars, nil)
	stack := NewStack()

	// Push in order: 1, 2, 3
	stack.Push(NewFloat(1))
	stack.Push(NewFloat(2))
	stack.Push(NewFloat(3))

	result, err := ops.Show(stack)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Show uses stack.Values() which is bottom-to-top, so "1 2 3"
	expected := "1 2 3"
	if result != expected {
		t.Errorf("Show = %q, want %q", result, expected)
	}
}
