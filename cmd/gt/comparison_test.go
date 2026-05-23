// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package main

import (
	"testing"
)

// TestRunCommandComparisonGreater tests gt operator
func TestRunCommandComparisonGreater(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"3 gt 5 = false", []string{"gt", "3", "5", "gt"}, "false"},
		{"5 gt 3 = true", []string{"gt", "5", "3", "gt"}, "true"},
		{"5 gt 5 = false", []string{"gt", "5", "5", "gt"}, "false"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runCommand(tt.args)
			if err != nil {
				t.Fatalf("runCommand(%v) returned error: %v", tt.args, err)
			}
			if result != tt.want {
				t.Errorf("runCommand(%v) = %q, want %q", tt.args, result, tt.want)
			}
		})
	}
}

// TestRunCommandComparisonLess tests lt operator
func TestRunCommandComparisonLess(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"3 lt 5 = true", []string{"gt", "3", "5", "lt"}, "true"},
		{"5 lt 3 = false", []string{"gt", "5", "3", "lt"}, "false"},
		{"5 lt 5 = false", []string{"gt", "5", "5", "lt"}, "false"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runCommand(tt.args)
			if err != nil {
				t.Fatalf("runCommand(%v) returned error: %v", tt.args, err)
			}
			if result != tt.want {
				t.Errorf("runCommand(%v) = %q, want %q", tt.args, result, tt.want)
			}
		})
	}
}

// TestRunCommandComparisonEqual tests eq and == operators
func TestRunCommandComparisonEqual(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"5 eq 5 = true", []string{"gt", "5", "5", "eq"}, "true"},
		{"5 eq 3 = false", []string{"gt", "5", "3", "eq"}, "false"},
		{"5 == 5 = true", []string{"gt", "5", "5", "=="}, "true"},
		{"5 == 3 = false", []string{"gt", "5", "3", "=="}, "false"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runCommand(tt.args)
			if err != nil {
				t.Fatalf("runCommand(%v) returned error: %v", tt.args, err)
			}
			if result != tt.want {
				t.Errorf("runCommand(%v) = %q, want %q", tt.args, result, tt.want)
			}
		})
	}
}

// TestRunCommandComparisonNotEqual tests neq and != operators
func TestRunCommandComparisonNotEqual(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"5 neq 3 = true", []string{"gt", "5", "3", "neq"}, "true"},
		{"5 neq 5 = false", []string{"gt", "5", "5", "neq"}, "false"},
		{"5 != 3 = true", []string{"gt", "5", "3", "!="}, "true"},
		{"5 != 5 = false", []string{"gt", "5", "5", "!="}, "false"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runCommand(tt.args)
			if err != nil {
				t.Fatalf("runCommand(%v) returned error: %v", tt.args, err)
			}
			if result != tt.want {
				t.Errorf("runCommand(%v) = %q, want %q", tt.args, result, tt.want)
			}
		})
	}
}

// TestRunCommandComparisonGteLte tests gte, lte, >=, <= operators
func TestRunCommandComparisonGteLte(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"5 gte 3 = true", []string{"gt", "5", "3", "gte"}, "true"},
		{"3 gte 3 = true", []string{"gt", "3", "3", "gte"}, "true"},
		{"3 gte 5 = false", []string{"gt", "3", "5", "gte"}, "false"},
		{"5 lte 3 = false", []string{"gt", "5", "3", "lte"}, "false"},
		{"3 lte 3 = true", []string{"gt", "3", "3", "lte"}, "true"},
		{"3 lte 5 = true", []string{"gt", "3", "5", "lte"}, "true"},
		{"5 >= 3 = true", []string{"gt", "5", "3", ">="}, "true"},
		{"3 >= 3 = true", []string{"gt", "3", "3", ">="}, "true"},
		{"3 >= 5 = false", []string{"gt", "3", "5", ">="}, "false"},
		{"5 <= 3 = false", []string{"gt", "5", "3", "<="}, "false"},
		{"3 <= 3 = true", []string{"gt", "3", "3", "<="}, "true"},
		{"3 <= 5 = true", []string{"gt", "3", "5", "<="}, "true"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runCommand(tt.args)
			if err != nil {
				t.Fatalf("runCommand(%v) returned error: %v", tt.args, err)
			}
			if result != tt.want {
				t.Errorf("runCommand(%v) = %q, want %q", tt.args, result, tt.want)
			}
		})
	}
}

// TestRunCommandComparisonWithMetrics tests metric-aware comparisons
func TestRunCommandComparisonWithMetrics(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"100Mbps gt 50Mbps = true", []string{"gt", "100Mbps", "50Mbps", "gt"}, "true"},
		{"50Mbps gt 100Mbps = false", []string{"gt", "50Mbps", "100Mbps", "gt"}, "false"},
		{"100Mbps eq 100Mbps = true", []string{"gt", "100Mbps", "100Mbps", "eq"}, "true"},
		{"1km gt 500m = true", []string{"gt", "1km", "500m", "gt"}, "true"},
		{"500m lt 1km = true", []string{"gt", "500m", "1km", "lt"}, "true"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runCommand(tt.args)
			if err != nil {
				t.Fatalf("runCommand(%v) returned error: %v", tt.args, err)
			}
			if result != tt.want {
				t.Errorf("runCommand(%v) = %q, want %q", tt.args, result, tt.want)
			}
		})
	}
}

// TestRunCommandComparisonCrossCategoryError tests that mixed categories
// in comparisons produce an error (the exact error message depends on
// the parsing pipeline, so we just check for any error)
func TestRunCommandComparisonCrossCategoryError(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"DataRate vs Time", []string{"gt", "100Mbps", "2hr", "gt"}},
		{"Speed vs Distance", []string{"gt", "100mph", "5mi", "eq"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := runCommand(tt.args)
			if err == nil {
				t.Errorf("runCommand(%v) should return error for cross-category comparison", tt.args)
			}
		})
	}
}

// TestRunCommandComparisonSymbolError tests that comparisons with symbols produce an error
func TestRunCommandComparisonSymbolError(t *testing.T) {
	_, err := runCommand([]string{"gt", "5", ":x", "gt"})
	if err == nil {
		t.Error("comparing with a symbol should produce an error")
	}
}
