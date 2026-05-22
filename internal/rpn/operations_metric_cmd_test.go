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

type prefixModeTestCase struct {
	name         string
	setMode      string // "si", "iec", or "" (default)
	expr         string
	wantNum      float64
	tol          float64
	wantContains string
}

func TestPrefixMode(t *testing.T) {
	cases := []prefixModeTestCase{
		// End-to-end conversion: SI mode
		{
			name:    "SI:1GB-to-MB",
			setMode: "si",
			expr:    "1GB @MB convert",
			wantNum: 1000,
			tol:     0.1,
		},
		// End-to-end conversion: IEC mode
		{
			name:    "IEC:1GB-to-MB",
			setMode: "iec",
			expr:    "1GB @MB convert",
			wantNum: 1024,
			tol:     0.1,
		},
		// Cross-metric conversion: SI mode (GiB is 2^30, GB is 10^9)
		{
			name:    "SI:1GiB-to-GB",
			setMode: "si",
			expr:    "1GiB @GB convert",
			wantNum: 1.07374,
			tol:     0.01,
		},
		// Cross-metric conversion: IEC mode (both use powers of 2)
		{
			name:    "IEC:1GiB-to-GB",
			setMode: "iec",
			expr:    "1GiB @GB convert",
			wantNum: 1.0,
			tol:     0.001,
		},
		// Arithmetic: SI mode (1024KB + 1KB = 1025KB)
		{
			name:    "SI:1024KB-plus-1KB",
			setMode: "si",
			expr:    "1024KB 1KB + @KB convert",
			wantNum: 1025,
			tol:     0.1,
		},
		// Arithmetic: IEC mode (same result, factors cancel)
		{
			name:    "IEC:1024KB-plus-1KB",
			setMode: "iec",
			expr:    "1024KB 1KB + @KB convert",
			wantNum: 1025,
			tol:     0.1,
		},
		// Comparison: SI mode — 1GB == 1000MB (8e9 == 1000*8e6)
		{
			name:         "SI:1GB-eq-1000MB",
			setMode:      "si",
			expr:         "1GB 1000MB eq",
			wantContains: "true",
		},
		// Comparison: IEC mode — 1GB != 1000MB (different base)
		{
			name:         "IEC:1GB-neq-1000MB",
			setMode:      "iec",
			expr:         "1GB 1000MB eq",
			wantContains: "false",
		},
		// Comparison: IEC mode — 1GB == 1024MB (powers of 2)
		{
			name:         "IEC:1GB-eq-1024MB",
			setMode:      "iec",
			expr:         "1GB 1024MB eq",
			wantContains: "true",
		},
		// Metric show: SI mode — GB factor is 8e+09
		{
			name:         "SI:metric-show-GB-factor",
			setMode:      "si",
			expr:         "1GB metric show",
			wantContains: "8e+09",
		},
		// Metric show: IEC mode — GB factor is ~9e+09 (8*2^30 ≈ 8.59e9)
		{
			name:         "IEC:metric-show-GB-factor",
			setMode:      "iec",
			expr:         "1GB metric show",
			wantContains: "9e+09",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			vars := NewVariables()
			rpn := NewRPN(vars)

			if tc.setMode == "iec" {
				if _, err := rpn.ParseAndEvaluate("metric binary set"); err != nil {
					t.Fatalf("metric binary set failed: %v", err)
				}
			}

			result, err := rpn.ParseAndEvaluate(tc.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.wantNum > 0 || tc.tol > 0 {
				got, err := strconv.ParseFloat(result, 64)
				if err != nil {
					t.Fatalf("cannot parse result %q as float: %v", result, err)
				}
				if got < tc.wantNum-tc.tol || got > tc.wantNum+tc.tol {
					t.Errorf("got %g, want %g±%g", got, tc.wantNum, tc.tol)
				}
			}

			if tc.wantContains != "" {
				if !strings.Contains(result, tc.wantContains) {
					t.Errorf("expected result to contain %q, got: %s", tc.wantContains, result)
				}
			}
		})
	}
}
