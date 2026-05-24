// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"math"
	"strconv"
	"testing"
)

// TestHyperMetricAwareOperations tests [+-%] ops that preserve metric categories.
func TestHyperMetricAwareOperations(t *testing.T) {
	reg := GetMetricRegistry()

	tests := []struct {
		name    string
		expr    string
		wantNum float64
		wantMet string // empty skips metric check
		tol     float64
		wantErr bool
	}{
		// Addition
		{
			name: "add same category", expr: "100Mbps 50Mbps 25Mbps [+]",
			wantNum: 175, wantMet: "Mbps", tol: 0.001,
		},
		{
			name: "add mixed units", expr: "1km 500m 100m [+]",
			wantNum: 1.6, wantMet: "km", tol: 0.001,
		},
		{
			name: "add cool absorbing", expr: "5 100Mbps 10Mbps [+]",
			wantNum: 115, wantMet: "Mbps", tol: 0.001,
		},
		{
			name: "add all cool", expr: "1 2 3 [+]",
			wantNum: 6, wantMet: "Cool", tol: 0.001,
		},
		{
			name: "add incompatible", expr: "100Mbps 2hr [+]",
			wantErr: true,
		},
		// Subtraction
		{
			name: "sub same category", expr: "1000Mbps 100Mbps 50Mbps [-]",
			wantNum: 850, wantMet: "Mbps", tol: 0.001,
		},
		{
			name: "sub mixed units", expr: "2km 500m 100m [-]",
			wantNum: 1.4, wantMet: "km", tol: 0.001,
		},
		{
			name: "sub cool absorbing", expr: "100km 5 [-]",
			wantNum: 95, wantMet: "km", tol: 1,
		},
		{
			name: "sub negative", expr: "1km 2km [-]",
			wantNum: -1, wantMet: "km", tol: 0.1,
		},
		{
			name: "sub incompatible", expr: "100km 2hr [-]",
			wantErr: true,
		},
		// Modulo
		{
			name: "mod same category", expr: "10km 3km 2km [%]",
			wantNum: 1, wantMet: "km", tol: 0.001,
		},
		{
			name: "mod mixed units", expr: "1000m 300m 200m [%]",
			wantNum: 100, tol: 0.001,
		},
		{
			name: "mod incompatible", expr: "100km 2hr [%]",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vars := NewVariables()
			rpn := NewRPN(vars, nil)
			result, err := rpn.ParseAndEvaluate(tc.expr)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got, _ := strconv.ParseFloat(result, 64)
			tol := tc.tol
			if tol == 0 {
				tol = 0.001
			}
			if got < tc.wantNum-tol || got > tc.wantNum+tol {
				t.Errorf("result = %g, want %g (tol %g)", got, tc.wantNum, tol)
			}

			if tc.wantMet != "" {
				stack := rpn.GetCurrentStack()
				m := stack[len(stack)-1].Metric()
				want, _ := reg.Find(tc.wantMet)
				if m != want {
					t.Errorf("metric = %v, want %s", m, tc.wantMet)
				}
			}
		})
	}
}

// TestHyperCoolResultOperations tests [*][/][^][lg][log][ln] ops that always
// produce a Cool metric result.
func TestHyperCoolResultOperations(t *testing.T) {
	reg := GetMetricRegistry()

	tests := []struct {
		name    string
		expr    string
		wantNum float64
		tol     float64
	}{
		{
			name: "multiply basic", expr: "3 4 5 [*]",
			wantNum: 60, tol: 0.001,
		},
		{
			name: "multiply mixed metrics", expr: "100Mbps 2hr [*]",
			wantNum: 200, tol: 0.001,
		},
		{
			name: "divide basic", expr: "100 20 2 [/]",
			wantNum: 2.5, tol: 0.001,
		},
		{
			name: "power basic", expr: "2 3 2 [^]",
			wantNum: 64, tol: 0.001,
		},
		{
			name: "log2 basic", expr: "2 4 8 [lg]",
			wantNum: 6, tol: 0.001,
		},
		{
			name: "log10 basic", expr: "10 100 [log]",
			wantNum: 3, tol: 0.001,
		},
		{
			name: "ln basic", expr: "2.718281828459045 7.38905609893065 [ln]",
			wantNum: 3, tol: 0.001,
		},
		{
			name: "log2 with metrics", expr: "100Mbps 1000Mbps [lg]",
			wantNum: math.Log2(100) + math.Log2(1000), tol: 0.01,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vars := NewVariables()
			rpn := NewRPN(vars, nil)
			result, err := rpn.ParseAndEvaluate(tc.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got, _ := strconv.ParseFloat(result, 64)
			if got < tc.wantNum-tc.tol || got > tc.wantNum+tc.tol {
				t.Errorf("result = %g, want %g (tol %g)", got, tc.wantNum, tc.tol)
			}

			// All results must be Cool
			stack := rpn.GetCurrentStack()
			m := stack[len(stack)-1].Metric()
			cool, _ := reg.Find("Cool")
			if m != cool {
				t.Errorf("metric = %v, want Cool", m)
			}
		})
	}
}

// TestHyperErrorCases covers expressions that must fail evaluation.
func TestHyperErrorCases(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"modulo by zero", "10km 0km [%]"},
		{"divide by zero", "10 0 [/]"},
		{"log2 non-positive", "0 5 [lg]"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vars := NewVariables()
			rpn := NewRPN(vars, nil)
			_, err := rpn.ParseAndEvaluate(tc.expr)
			if err == nil {
				t.Fatal("expected error, got none")
			}
		})
	}
}
