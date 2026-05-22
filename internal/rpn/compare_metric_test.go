// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"strings"
	"testing"
)

// TestMetricAwareComparisons tests all six comparison operators with metric operands.
func TestMetricAwareComparisons(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		expected    string
		prefixMode  PrefixMode
		description string
	}{
		{
			name:        "gt 1GB > 500MB SI",
			expression:  "1GB 500MB gt",
			expected:    "true",
			prefixMode:  SI,
			description: "1GB (1000MB SI) > 500MB",
		},
		{
			name:        "eq 1GB == 1000MB SI",
			expression:  "1GB 1000MB eq",
			expected:    "true",
			prefixMode:  SI,
			description: "1GB = 1000MB in SI mode",
		},
		{
			name:        "eq 1GB == 1024MB IEC",
			expression:  "1GB 1024MB eq",
			expected:    "true",
			prefixMode:  IEC,
			description: "1GB = 1024MB in IEC mode",
		},
		{
			name:        "neq 1GB != 1000MB IEC",
			expression:  "1GB 1000MB neq",
			expected:    "true",
			prefixMode:  IEC,
			description: "1GB != 1000MB in IEC mode",
		},
		{
			name:        "neq 1GB != 1024MB SI",
			expression:  "1GB 1024MB neq",
			expected:    "true",
			prefixMode:  SI,
			description: "1GB != 1024MB in SI mode",
		},
		{
			name:        "lt 30s < 1min",
			expression:  "30s 1min lt",
			expected:    "true",
			prefixMode:  0,
			description: "30s < 60s",
		},
		{
			name:        "gte 1hr >= 3600s",
			expression:  "1hr 3600s gte",
			expected:    "true",
			prefixMode:  0,
			description: "3600s >= 3600s",
		},
		{
			name:        "lte 2.2lb <= 1kg",
			expression:  "2.2lb 1kg lte",
			expected:    "true",
			prefixMode:  0,
			description: "2.2lb (0.998kg) <= 1kg",
		},
		{
			name:        "eq 1Gbps == 1000Mbps",
			expression:  "1Gbps 1000Mbps eq",
			expected:    "true",
			prefixMode:  0,
			description: "1Gbps = 1000Mbps",
		},
		{
			name:        "gt 1mi > 1600m",
			expression:  "1mi 1600m gt",
			expected:    "true",
			prefixMode:  0,
			description: "1mi (1609.344m) > 1600m",
		},
		{
			name:        "eq 1GiB == 1024MiB IEC",
			expression:  "1GiB 1024MiB eq",
			expected:    "true",
			prefixMode:  IEC,
			description: "1GiB = 1024MiB",
		},
		{
			name:        "lte 100Mbps <= 1Gbps",
			expression:  "100Mbps 1Gbps lte",
			expected:    "true",
			prefixMode:  0,
			description: "100Mbps <= 1000Mbps",
		},
		{
			name:        "neq plain numbers",
			expression:  "5 3 neq",
			expected:    "true",
			prefixMode:  0,
			description: "5 != 3",
		},
		{
			name:        "eq plain numbers",
			expression:  "5 5 eq",
			expected:    "true",
			prefixMode:  0,
			description: "5 == 5",
		},
		{
			name:        "gt 500MB > 0 SI",
			expression:  "500MB 0 gt",
			expected:    "true",
			prefixMode:  SI,
			description: "500MB > 0 (Cool absorbed)",
		},
		// False cases — exercises the false-result code path with metrics
		{
			name:        "gt 500MB > 1GB false SI",
			expression:  "500MB 1GB gt",
			expected:    "false",
			prefixMode:  SI,
			description: "500MB < 1GB, so 500MB > 1GB is false",
		},
		{
			name:        "lt 1GB < 500MB false",
			expression:  "1GB 500MB lt",
			expected:    "false",
			prefixMode:  0,
			description: "1GB > 500MB, so 1GB < 500MB is false",
		},
		{
			name:        "gte 1min >= 1hr false",
			expression:  "1min 1hr gte",
			expected:    "false",
			prefixMode:  0,
			description: "60s >= 3600s is false",
		},
		{
			name:        "lte 1hr <= 1min false",
			expression:  "1hr 1min lte",
			expected:    "false",
			prefixMode:  0,
			description: "3600s <= 60s is false",
		},
		{
			name:        "eq 1GB == 1GiB SI",
			expression:  "1GB 1GiB eq",
			expected:    "false",
			prefixMode:  SI,
			description: "1GB (SI=8e9 bits) != 1GiB (8*2^30 bits) in SI mode",
		},
		{
			name:        "neq 1GB != 1000MB false SI",
			expression:  "1GB 1000MB neq",
			expected:    "false",
			prefixMode:  SI,
			description: "1GB == 1000MB in SI, so != is false",
		},
		{
			name:        "neq 1hr != 3600s false",
			expression:  "1hr 3600s neq",
			expected:    "false",
			prefixMode:  0,
			description: "1hr == 3600s, so != is false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars := NewVariables()
			rpnCalc := NewRPN(vars)

			if tt.prefixMode == IEC {
				rpnCalc.SetPrefixMode(IEC)
			}

			result, err := rpnCalc.ParseAndEvaluate(tt.expression)

			if err != nil {
				t.Fatalf("Evaluate(%q) returned error: %v", tt.expression, err)
			}

			if result != tt.expected {
				t.Errorf("Evaluate(%q) = %q, want %q (%s)", tt.expression, result, tt.expected, tt.description)
			}
		})
	}
}

// TestShorthandComparisonOperators tests the shorthand operators <, >, >=, <=, ==, !=
// that delegate to the same comparison methods.
func TestShorthandComparisonOperators(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		expected    string
		description string
	}{
		{
			name:        "shorthand > gt",
			expression:  "5 3 >",
			expected:    "true",
			description: "5 > 3 via shorthand",
		},
		{
			name:        "shorthand < lt",
			expression:  "3 5 <",
			expected:    "true",
			description: "3 < 5 via shorthand",
		},
		{
			name:        "shorthand >= gte",
			expression:  "5 5 >=",
			expected:    "true",
			description: "5 >= 5 via shorthand",
		},
		{
			name:        "shorthand <= lte",
			expression:  "5 5 <=",
			expected:    "true",
			description: "5 <= 5 via shorthand",
		},
		{
			name:        "shorthand == eq",
			expression:  "5 5 ==",
			expected:    "true",
			description: "5 == 5 via shorthand",
		},
		{
			name:        "shorthand != neq",
			expression:  "5 3 !=",
			expected:    "true",
			description: "5 != 3 via shorthand",
		},
		{
			name:        "shorthand > metric",
			expression:  "1GB 500MB >",
			expected:    "true",
			description: "1GB > 500MB via shorthand >",
		},
		{
			name:        "shorthand < metric",
			expression:  "500MB 1GB <",
			expected:    "true",
			description: "500MB < 1GB via shorthand <",
		},
		{
			name:        "shorthand == metric",
			expression:  "1hr 3600s ==",
			expected:    "true",
			description: "1hr == 3600s via shorthand ==",
		},
		{
			name:        "shorthand != metric",
			expression:  "1hr 1min !=",
			expected:    "true",
			description: "1hr != 1min via shorthand !=",
		},
		{
			name:        "shorthand >= metric",
			expression:  "1km 1000m >=",
			expected:    "true",
			description: "1km >= 1000m via shorthand >=",
		},
		{
			name:        "shorthand <= metric",
			expression:  "1mi 1609.344m <=",
			expected:    "true",
			description: "1mi <= 1609.344m via shorthand <=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars := NewVariables()
			rpnCalc := NewRPN(vars)

			result, err := rpnCalc.ParseAndEvaluate(tt.expression)

			if err != nil {
				t.Fatalf("Evaluate(%q) returned error: %v", tt.expression, err)
			}

			if result != tt.expected {
				t.Errorf("Evaluate(%q) = %q, want %q (%s)", tt.expression, result, tt.expected, tt.description)
			}
		})
	}
}

// TestIncompatibleCategoryComparison verifies that comparing incompatible metric
// categories returns an error.
func TestIncompatibleCategoryComparison(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		description string
	}{
		{
			name:        "dataRate vs time",
			expression:  "100Mbps 2hr eq",
			description: "100Mbps vs 2hr are incompatible",
		},
		{
			name:        "dataRate vs time gt",
			expression:  "1Gbps 1hr gt",
			description: "1Gbps vs 1hr are incompatible",
		},
		{
			name:        "time vs weight",
			expression:  "1hr 5kg lt",
			description: "1hr vs 5kg are incompatible",
		},
		{
			name:        "weight vs distance",
			expression:  "10lb 5mi neq",
			description: "10lb vs 5mi are incompatible",
		},
		{
			name:        "dataSize vs speed",
			expression:  "1GB 100mph gte",
			description: "1GB vs 100mph are incompatible",
		},
		{
			name:        "distance vs dataRate",
			expression:  "1km 10Mbps lte",
			description: "1km vs 10Mbps are incompatible",
		},
		// Shorthand operators with incompatible categories
		{
			name:        "shorthand > incompatible",
			expression:  "1Gbps 1hr >",
			description: "1Gbps > 1hr should error (shorthand >)",
		},
		{
			name:        "shorthand != incompatible",
			expression:  "5kg 10mi !=",
			description: "5kg != 10mi should error (shorthand !=)",
		},
		{
			name:        "shorthand == incompatible",
			expression:  "100mph 1GB ==",
			description: "100mph == 1GB should error (shorthand ==)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars := NewVariables()
			rpnCalc := NewRPN(vars)

			_, err := rpnCalc.ParseAndEvaluate(tt.expression)

			if err == nil {
				t.Errorf("Evaluate(%q) expected error for incompatible categories (%s), got nil", tt.expression, tt.description)
				return
			}

			if !strings.Contains(err.Error(), "incompatible") {
				t.Errorf("Evaluate(%q) expected error mentioning 'incompatible', got: %v", tt.expression, err)
			}
		})
	}
}

// TestMetricComparisonEdgeCases tests edge cases for metric-aware comparisons.
func TestMetricComparisonEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		expected    string
		prefixMode  PrefixMode
		description string
	}{
		{
			name:        "zero values 0GB == 0MB",
			expression:  "0GB 0MB eq",
			expected:    "true",
			prefixMode:  0,
			description: "0GB == 0MB",
		},
		{
			name:        "negative -1GB == -1000MB SI",
			expression:  "-1GB -1000MB eq",
			expected:    "true",
			prefixMode:  SI,
			description: "-1GB == -1000MB in SI mode",
		},
		{
			name:        "1000bps == 1Kbps",
			expression:  "1000bps 1Kbps eq",
			expected:    "true",
			prefixMode:  0,
			description: "1000bps == 1Kbps",
		},
		{
			name:        "cool gt 5 > 3",
			expression:  "5 3 gt",
			expected:    "true",
			prefixMode:  0,
			description: "5 > 3 (both Cool)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars := NewVariables()
			rpnCalc := NewRPN(vars)

			if tt.prefixMode == IEC {
				rpnCalc.SetPrefixMode(IEC)
			}

			result, err := rpnCalc.ParseAndEvaluate(tt.expression)

			if err != nil {
				t.Fatalf("Evaluate(%q) returned error: %v", tt.expression, err)
			}

			if result != tt.expected {
				t.Errorf("Evaluate(%q) = %q, want %q (%s)", tt.expression, result, tt.expected, tt.description)
			}
		})
	}
}
