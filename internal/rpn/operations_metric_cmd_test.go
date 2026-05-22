// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"strings"
	"testing"
)

func TestMetricShow(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	result, err := rpn.ParseAndEvaluate("100Mbps metric show")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Mbps") {
		t.Errorf("expected result to contain 'Mbps', got: %s", result)
	}
	if !strings.Contains(result, "DataRate") {
		t.Errorf("expected result to contain 'DataRate', got: %s", result)
	}
}

func TestMetricShowUniversal(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	result, err := rpn.ParseAndEvaluate("42 metric show")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Cool") {
		t.Errorf("expected result to contain 'Cool', got: %s", result)
	}
	if !strings.Contains(result, "Universal") {
		t.Errorf("expected result to contain 'Universal', got: %s", result)
	}
}

func TestMetricShowEmptyStack(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	_, err := rpn.ParseAndEvaluate("metric show")
	if err == nil {
		t.Error("expected error for empty stack")
	}
}

func TestMetricList(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("metric list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"DataRate", "DataSize", "Distance", "Speed", "Time", "Universal", "Weight"}
	for _, cat := range expected {
		if !strings.Contains(result, cat) {
			t.Errorf("expected result to contain %q, got: %s", cat, result)
		}
	}
}

func TestMetricCategory(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("metric DataRate")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"bps", "Kbps", "Mbps", "Gbps", "Tbps"}
	for _, name := range expected {
		if !strings.Contains(result, name) {
			t.Errorf("expected result to contain %q, got: %s", name, result)
		}
	}
}

func TestMetricCategoryTime(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("metric Time")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"ms", "s", "min", "hr", "day"}
	for _, name := range expected {
		if !strings.Contains(result, name) {
			t.Errorf("expected result to contain %q, got: %s", name, result)
		}
	}
}

func TestMetricCategoryUnknown(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	_, err := rpn.ParseAndEvaluate("metric Nope")
	if err == nil {
		t.Error("expected error for unknown category")
	}
}

func TestMetricCategoryCustom(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	// Custom category exists but has no built-in metrics, so result should be empty
	result, err := rpn.ParseAndEvaluate("metric Custom")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Empty string is fine — no custom metrics registered by default
	_ = result
}

func TestMetricSetModeDecimal(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	result, err := rpn.ParseAndEvaluate("metric decimal set")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "SI") {
		t.Errorf("expected SI, got: %s", result)
	}
	if rpn.GetPrefixMode() != SI {
		t.Errorf("prefix mode = %v, want SI", rpn.GetPrefixMode())
	}
}

func TestMetricSetModeBinary(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	result, err := rpn.ParseAndEvaluate("metric binary set")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "IEC") {
		t.Errorf("expected IEC, got: %s", result)
	}
	if rpn.GetPrefixMode() != IEC {
		t.Errorf("prefix mode = %v, want IEC", rpn.GetPrefixMode())
	}
}

func TestMetricSetModeToggle(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)

	// Default is SI
	if rpn.GetPrefixMode() != SI {
		t.Errorf("default prefix mode = %v, want SI", rpn.GetPrefixMode())
	}

	// Switch to IEC
	rpn.ParseAndEvaluate("metric binary set")
	if rpn.GetPrefixMode() != IEC {
		t.Errorf("after binary set, prefix mode = %v, want IEC", rpn.GetPrefixMode())
	}

	// Switch back to SI
	rpn.ParseAndEvaluate("metric decimal set")
	if rpn.GetPrefixMode() != SI {
		t.Errorf("after decimal set, prefix mode = %v, want SI", rpn.GetPrefixMode())
	}
}

func TestMetricSetModeBinaryIncomplete(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	_, err := rpn.ParseAndEvaluate("metric binary")
	if err == nil {
		t.Error("expected error for incomplete 'metric binary'")
	}
}

func TestMetricSetModeDecimalIncomplete(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	_, err := rpn.ParseAndEvaluate("metric decimal")
	if err == nil {
		t.Error("expected error for incomplete 'metric decimal'")
	}
}

func TestMetricCompatibleSameCategory(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("100Mbps 1Gbps metric compatible")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "true") {
		t.Errorf("expected Mbps and Gbps to be compatible, got: %s", result)
	}
}

func TestMetricCompatibleDifferentCategory(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("100Mbps 2hr metric compatible")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Mbps") {
		t.Errorf("expected result to mention Mbps, got: %s", result)
	}
	if !strings.Contains(result, "hr") {
		t.Errorf("expected result to mention hr, got: %s", result)
	}
	// DataRate and Time are different categories, neither is Universal -> false
	if !strings.Contains(result, "false") {
		t.Errorf("expected DataRate and Time to be incompatible for +/-, got: %s", result)
	}
}

func TestMetricCompatibleWithUniversal(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("100Mbps 42 metric compatible")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Universal (Cool) is compatible with anything
	if !strings.Contains(result, "true") {
		t.Errorf("expected Cool and Mbps to be compatible, got: %s", result)
	}
}

func TestMetricCompatibleNotEnoughValues(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	_, err := rpn.ParseAndEvaluate("100 metric compatible")
	if err == nil {
		t.Error("expected error for insufficient stack values")
	}
}

func TestMetricCompatibleEmptyStack(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	_, err := rpn.ParseAndEvaluate("metric compatible")
	if err == nil {
		t.Error("expected error for empty stack")
	}
}
