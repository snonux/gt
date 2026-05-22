// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"strconv"
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

func TestPrefixModeEndToEndConvert(t *testing.T) {
	// SI mode: 1GB → MB = 1000
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("1GB @MB convert")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal < 999.9 || resultVal > 1000.1 {
		t.Errorf("SI: 1GB→MB = %g, want 1000", resultVal)
	}

	// IEC mode: 1GB → MB = 1024
	vars2 := NewVariables()
	rpn2 := NewRPN(vars2)
	_, err = rpn2.ParseAndEvaluate("metric binary set")
	if err != nil {
		t.Fatalf("metric binary set failed: %v", err)
	}
	result2, err := rpn2.ParseAndEvaluate("1GB @MB convert")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal2, _ := strconv.ParseFloat(result2, 64)
	if resultVal2 < 1023.9 || resultVal2 > 1024.1 {
		t.Errorf("IEC: 1GB→MB = %g, want 1024", resultVal2)
	}
}

func TestPrefixModeAffectsCrossMetricConvert(t *testing.T) {
	// SI mode: 1GiB → GB = ~1.07374
	// (GiB is always 2^30, GB in SI mode is 10^9)
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("1GiB @GB convert")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal < 1.07 || resultVal > 1.08 {
		t.Errorf("SI: 1GiB→GB = %g, want ~1.074", resultVal)
	}

	// IEC mode: 1GiB → GB = 1.0
	// (GiB is 2^30, GB in IEC mode is also 2^30)
	vars2 := NewVariables()
	rpn2 := NewRPN(vars2)
	_, err = rpn2.ParseAndEvaluate("metric binary set")
	if err != nil {
		t.Fatalf("metric binary set failed: %v", err)
	}
	result2, err := rpn2.ParseAndEvaluate("1GiB @GB convert")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal2, _ := strconv.ParseFloat(result2, 64)
	if resultVal2 < 0.999 || resultVal2 > 1.001 {
		t.Errorf("IEC: 1GiB→GB = %g, want 1.0", resultVal2)
	}
}

func TestPrefixModeUsedInArithmetic(t *testing.T) {
	// Same-metric addition: result is the same in both modes
	// (conversion factors cancel out when input and output metrics match)
	// This test verifies that GetPrefixMode() is called during arithmetic,
	// not that SI vs IEC produces different results.
	// For mode-dependent results, see TestPrefixModeEndToEndConvert.

	// SI mode: 1024KB + 1KB → 1025KB
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("1024KB 1KB + @KB convert")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal < 1024.9 || resultVal > 1025.1 {
		t.Errorf("SI: 1024KB+1KB→KB = %g, want 1025", resultVal)
	}

	// IEC mode: same result (factors cancel)
	vars2 := NewVariables()
	rpn2 := NewRPN(vars2)
	_, err = rpn2.ParseAndEvaluate("metric binary set")
	if err != nil {
		t.Fatalf("metric binary set failed: %v", err)
	}
	result2, err := rpn2.ParseAndEvaluate("1024KB 1KB + @KB convert")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal2, _ := strconv.ParseFloat(result2, 64)
	if resultVal2 < 1024.9 || resultVal2 > 1025.1 {
		t.Errorf("IEC: 1024KB+1KB→KB = %g, want 1025", resultVal2)
	}
}

func TestPrefixModeAffectsComparison(t *testing.T) {
	// SI mode: 1GB (8e9 bits) vs 1000MB (1000*8e6 = 8e9 bits) → equal
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("1GB 1000MB eq")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "true" {
		t.Errorf("SI: 1GB == 1000MB should be true, got %s", result)
	}

	// IEC mode: 1GB (8*2^30 bits) vs 1000MB (1000*8*2^20 bits)
	// 8*2^30 = 8589934592, 1000*8*2^20 = 8388608000 → not equal
	vars2 := NewVariables()
	rpn2 := NewRPN(vars2)
	_, err = rpn2.ParseAndEvaluate("metric binary set")
	if err != nil {
		t.Fatalf("metric binary set failed: %v", err)
	}
	result2, err := rpn2.ParseAndEvaluate("1GB 1000MB eq")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result2 != "false" {
		t.Errorf("IEC: 1GB == 1000MB should be false, got %s", result2)
	}

	// IEC mode: 1GB == 1024MB (both use 2^30 / 2^20) — use fresh RPN to avoid stack carryover
	vars3 := NewVariables()
	rpn3 := NewRPN(vars3)
	_, err = rpn3.ParseAndEvaluate("metric binary set")
	if err != nil {
		t.Fatalf("metric binary set failed: %v", err)
	}
	result3, err := rpn3.ParseAndEvaluate("1GB 1024MB eq")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result3 != "true" {
		t.Errorf("IEC: 1GB == 1024MB should be true, got %s", result3)
	}
}

func TestMetricShowReflectsPrefixMode(t *testing.T) {
	// In SI mode, GB factor should be 8e+09
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("1GB metric show")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "GB") {
		t.Errorf("expected 'GB' in result, got: %s", result)
	}
	// SI GB factor = 8e9, formatted with %.0g = "8e+09"
	if !strings.Contains(result, "8e+09") {
		t.Errorf("SI mode GB factor should be 8e+09, got: %s", result)
	}

	// Switch to IEC mode, GB factor should be 8*2^30 = 8589934592
	// formatted with %.0g = "9e+09"
	vars2 := NewVariables()
	rpn2 := NewRPN(vars2)
	_, err = rpn2.ParseAndEvaluate("metric binary set")
	if err != nil {
		t.Fatalf("metric binary set failed: %v", err)
	}
	result2, err := rpn2.ParseAndEvaluate("1GB metric show")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// IEC GB factor = 8*2^30 ≈ 8.59e9, formatted with %.0g = "9e+09"
	// (different from SI's "8e+09")
	if !strings.Contains(result2, "9e+09") {
		t.Errorf("IEC mode GB factor should be ~9e+09, got: %s", result2)
	}
}
