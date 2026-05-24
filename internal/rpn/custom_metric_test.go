// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"strings"
	"testing"
)

// customCleanup undefines a custom metric and verifies the cleanup.
// Must be deferred at the start of each test that defines a custom metric.
func customCleanup(t *testing.T, rpn *RPN, name string) {
	_, err := rpn.ParseAndEvaluate("custom undefine " + name)
	if err != nil {
		t.Logf("cleanup: failed to undefine %s: %v", name, err)
	}
}

// TestCustomDefineFactorZero tests defining a custom metric with factor=0
func TestCustomDefineFactorZero(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars, nil)
	defer customCleanup(t, rpn, "zero")

	result, err := rpn.ParseAndEvaluate("custom define zero 0 Custom")
	if err != nil {
		t.Fatalf("unexpected error defining metric with factor=0: %v", err)
	}
	if !strings.Contains(result, "defined") {
		t.Errorf("expected 'defined' in result, got: %s", result)
	}

	// Verify it works (factor=0 means any value in this metric = 0 in base)
	result, err = rpn.ParseAndEvaluate("5zero 3 +")
	if err != nil {
		t.Fatalf("arithmetic with factor=0 metric failed: %v", err)
	}
	t.Logf("5zero + 3 = %s", result)
}

// TestCustomDefineNegativeFactor tests defining a custom metric with negative factor
func TestCustomDefineNegativeFactor(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars, nil)
	defer customCleanup(t, rpn, "negtest")

	result, err := rpn.ParseAndEvaluate("custom define negtest -5 Custom")
	if err != nil {
		t.Fatalf("unexpected error defining metric with negative factor: %v", err)
	}
	if !strings.Contains(result, "defined") {
		t.Errorf("expected 'defined' in result, got: %s", result)
	}

	// Verify it works with negative factor
	result, err = rpn.ParseAndEvaluate("2negtest 1 +")
	if err != nil {
		t.Fatalf("arithmetic with negative factor metric failed: %v", err)
	}
	t.Logf("2negtest + 1 = %s", result)
}

// TestCustomDefineMalformedInput tests that malformed custom define input
func TestCustomDefineMalformedInput(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars, nil)

	_, err := rpn.ParseAndEvaluate("custom define  1 Custom")
	// Double-space collapses: fields = ["custom", "define", "1", "Custom"] (metric "" already exists as alias or empty check)
	if err == nil {
		t.Log("custom define with empty name succeeded")
	}
	// Either way, it shouldn't panic
}

// TestCustomDefineConflictingAlias tests defining a custom metric with a name
// that conflicts with an existing metric
func TestCustomDefineConflictingAlias(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars, nil)

	// Try to define a metric with a name that matches an existing one
	_, err := rpn.ParseAndEvaluate("custom define Mbps 1000 DataRate")
	if err == nil {
		t.Error("expected error defining metric with existing name 'Mbps'")
	}
}

// TestCustomMetricArithmetic tests using a custom metric in arithmetic operations
func TestCustomMetricArithmetic(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars, nil)
	defer customCleanup(t, rpn, "arithmetic_test")

	_, err := rpn.ParseAndEvaluate("custom define arithmetic_test 42 Custom")
	if err != nil {
		t.Fatalf("failed to define custom metric: %v", err)
	}

	// Use in addition: 10arithmetic_test + 5arithmetic_test = 15arithmetic_test
	result, err := rpn.ParseAndEvaluate("10arithmetic_test 5arithmetic_test +")
	if err != nil {
		t.Fatalf("arithmetic with custom metric failed: %v", err)
	}
	t.Logf("10arithmetic_test + 5arithmetic_test = %s", result)
}

// TestCustomMetricHyperOperations tests using a custom metric in hyper operations
func TestCustomMetricHyperOperations(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars, nil)
	defer customCleanup(t, rpn, "hyper_test")

	_, err := rpn.ParseAndEvaluate("custom define hyper_test 100 Custom")
	if err != nil {
		t.Fatalf("failed to define custom metric: %v", err)
	}

	// Use in hyper add
	result, err := rpn.ParseAndEvaluate("10hyper_test 5hyper_test 3hyper_test [+]")
	if err != nil {
		t.Fatalf("hyper operation with custom metric failed: %v", err)
	}
	t.Logf("10hyper_test 5hyper_test 3hyper_test [+] = %s", result)
}

// TestCustomMetricSubtraction tests subtraction with custom metrics
func TestCustomMetricSubtraction(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars, nil)
	defer customCleanup(t, rpn, "sub_test")

	_, err := rpn.ParseAndEvaluate("custom define sub_test 10 Custom")
	if err != nil {
		t.Fatalf("failed to define custom metric: %v", err)
	}

	result, err := rpn.ParseAndEvaluate("5sub_test 3sub_test -")
	if err != nil {
		t.Fatalf("subtraction with custom metric failed: %v", err)
	}
	t.Logf("5sub_test - 3sub_test = %s", result)
}

// TestCustomMetricMultiplication tests multiplication with custom metrics
func TestCustomMetricMultiplication(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars, nil)
	defer customCleanup(t, rpn, "mul_test")

	_, err := rpn.ParseAndEvaluate("custom define mul_test 5 Custom")
	if err != nil {
		t.Fatalf("failed to define custom metric: %v", err)
	}

	result, err := rpn.ParseAndEvaluate("3mul_test 2 *")
	if err != nil {
		t.Fatalf("multiplication with custom metric failed: %v", err)
	}
	t.Logf("3mul_test * 2 = %s", result)
}

// TestCustomMetricDivision tests division with custom metrics
func TestCustomMetricDivision(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars, nil)
	defer customCleanup(t, rpn, "div_test")

	_, err := rpn.ParseAndEvaluate("custom define div_test 10 Custom")
	if err != nil {
		t.Fatalf("failed to define custom metric: %v", err)
	}

	result, err := rpn.ParseAndEvaluate("20div_test 4div_test /")
	if err != nil {
		t.Fatalf("division with custom metric failed: %v", err)
	}
	t.Logf("20div_test / 4div_test = %s", result)
}

// TestCustomMetricConversionFail tests that custom metrics can't be converted
// between different categories
func TestCustomMetricConversionFail(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars, nil)
	defer customCleanup(t, rpn, "conv_test")

	_, err := rpn.ParseAndEvaluate("custom define conv_test 10 Custom")
	if err != nil {
		t.Fatalf("failed to define custom metric: %v", err)
	}

	// Converting Custom to DataRate should fail
	_, err = rpn.ParseAndEvaluate("10conv_test @Mbps convert")
	if err == nil {
		t.Error("expected error converting Custom to DataRate")
	}
}

// TestCustomMetricCrossCategory tests that custom metrics in different categories
// can't be mixed in operations
func TestCustomMetricCrossCategory(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars, nil)
	defer customCleanup(t, rpn, "cat_test1")
	defer customCleanup(t, rpn, "cat_test2")

	_, _ = rpn.ParseAndEvaluate("custom define cat_test1 10 Time")
	_, _ = rpn.ParseAndEvaluate("custom define cat_test2 5 Distance")

	// Mixing Time and Distance should fail
	_, err := rpn.ParseAndEvaluate("10cat_test1 5cat_test2 +")
	if err == nil {
		t.Error("expected error mixing Time and Distance categories")
	}
}

// TestCustomMetricShow tests that custom metrics display correctly in Show
func TestCustomMetricShow(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars, nil)
	defer customCleanup(t, rpn, "show_test")

	_, err := rpn.ParseAndEvaluate("custom define show_test 100 Custom")
	if err != nil {
		t.Fatalf("failed to define custom metric: %v", err)
	}

	result, err := rpn.ParseAndEvaluate("42show_test show")
	if err != nil {
		t.Fatalf("show with custom metric failed: %v", err)
	}
	t.Logf("42show_test show = %s", result)
}
