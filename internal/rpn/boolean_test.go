package rpn

import (
	"testing"
)

// TestBooleanOperators tests that boolean comparison operators work correctly.
func TestBooleanOperators(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		expected    string
		description string
	}{
		{
			name:        "gt true case",
			expression:  "5 3 gt",
			expected:    "true",
			description: "5 > 3 = true",
		},
		{
			name:        "gt false case",
			expression:  "3 5 gt",
			expected:    "false",
			description: "3 > 5 = false",
		},
		{
			name:        "gt equal case",
			expression:  "5 5 gt",
			expected:    "false",
			description: "5 > 5 = false",
		},
		{
			name:        "lt true case",
			expression:  "3 5 lt",
			expected:    "true",
			description: "3 < 5 = true",
		},
		{
			name:        "lt false case",
			expression:  "5 3 lt",
			expected:    "false",
			description: "5 < 3 = false",
		},
		{
			name:        "gte true case",
			expression:  "5 5 gte",
			expected:    "true",
			description: "5 >= 5 = true",
		},
		{
			name:        "gte false case",
			expression:  "3 5 gte",
			expected:    "false",
			description: "3 >= 5 = false",
		},
		{
			name:        "lte true case",
			expression:  "5 5 lte",
			expected:    "true",
			description: "5 <= 5 = true",
		},
		{
			name:        "lte false case",
			expression:  "5 3 lte",
			expected:    "false",
			description: "5 <= 3 = false",
		},
		{
			name:        "eq true case",
			expression:  "5 5 eq",
			expected:    "true",
			description: "5 == 5 = true",
		},
		{
			name:        "eq false case",
			expression:  "5 3 eq",
			expected:    "false",
			description: "5 == 3 = false",
		},
		{
			name:        "neq true case",
			expression:  "5 3 neq",
			expected:    "true",
			description: "5 != 3 = true",
		},
		{
			name:        "neq false case",
			expression:  "5 5 neq",
			expected:    "false",
			description: "5 != 5 = false",
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

// TestBooleanToNumberCoercion tests that boolean values are automatically
// coerced to numbers (true → 1, false → 0) in arithmetic operations.
func TestBooleanToNumberCoercion(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		expected    string
		description string
	}{
		{
			name:        "true in addition",
			expression:  "true 5 +",
			expected:    "6",
			description: "true (1) + 5 = 6",
		},
		{
			name:        "false in addition",
			expression:  "false 5 +",
			expected:    "5",
			description: "false (0) + 5 = 5",
		},
		{
			name:        "true in multiplication",
			expression:  "true 5 *",
			expected:    "5",
			description: "true (1) * 5 = 5",
		},
		{
			name:        "false in multiplication",
			expression:  "false 5 *",
			expected:    "0",
			description: "false (0) * 5 = 0",
		},
		{
			name:        "false in subtraction",
			expression:  "5 false -",
			expected:    "5",
			description: "5 - false (0) = 5",
		},
		{
			name:        "mixed boolean-numeric",
			expression:  "true false +",
			expected:    "1",
			description: "true (1) + false (0) = 1",
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

// TestMixedBooleanNumericArithmetic tests mixed boolean-numeric arithmetic.
func TestMixedBooleanNumericArithmetic(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		expected    string
		description string
	}{
		{
			name:        "5 3 gt 1 +",
			expression:  "5 3 gt 1 +",
			expected:    "2",
			description: "5 > 3 is true (1), 1 + 1 = 2",
		},
		{
			name:        "3 5 gt 1 +",
			expression:  "3 5 gt 1 +",
			expected:    "1",
			description: "3 > 5 is false (0), 0 + 1 = 1",
		},
		{
			name:        "true 2 *",
			expression:  "true 2 *",
			expected:    "2",
			description: "true (1) * 2 = 2",
		},
		{
			name:        "0 false +",
			expression:  "0 false +",
			expected:    "0",
			description: "0 + false (0) = 0",
		},
		{
			name:        "9 3 gt 4 5 lt +",
			expression:  "9 3 gt 4 5 lt +",
			expected:    "2",
			description: "9 > 3 is true (1), 4 < 5 is true (1), 1 + 1 = 2",
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

// TestBooleanShowFormat tests that Show command displays boolean values as true/false
func TestBooleanShowFormat(t *testing.T) {
	tests := []struct {
		name     string
		expression string
		expected string
	}{
		{
			name:     "show true",
			expression: "true show",
			expected: "true",
		},
		{
			name:     "show false",
			expression: "false show",
			expected: "false",
		},
		{
			name:     "show mixed stack",
			expression: "1 true 2 show",
			expected: "1 true 2",
		},
		{
			name:     "show comparison result",
			expression: "5 3 gt show",
			expected: "true",
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
				t.Errorf("Evaluate(%q) = %q, want %q", tt.expression, result, tt.expected)
			}
		})
	}
}
