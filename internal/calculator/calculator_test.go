// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package calculator

import (
	"strings"
	"testing"
)

// commonTestCases contains common test case patterns used across multiple tests
var commonTestCases = []struct {
	name     string
	input    string
	expected string
}{
	{
		name:     "20% of 150",
		input:    "20% of 150",
		expected: "20.00% of 150.00 = 30.00",
	},
	{
		name:     "what is 20% of 150",
		input:    "what is 20% of 150",
		expected: "20.00% of 150.00 = 30.00",
	},
	{
		name:     "50% of 200",
		input:    "50% of 200",
		expected: "50.00% of 200.00 = 100.00",
	},
	{
		name:     "decimal percent",
		input:    "12.5% of 80",
		expected: "12.50% of 80.00 = 10.00",
	},
	{
		name:     "decimal base",
		input:    "20% of 75.5",
		expected: "20.00% of 75.50 = 15.10",
	},
	{
		name:     "without 'of'",
		input:    "25% 400",
		expected: "25.00% of 400.00 = 100.00",
	},
	{
		name:     "30 is what % of 150",
		input:    "30 is what % of 150",
		expected: "30.00 is 20.00% of 150.00",
	},
	{
		name:     "50 is what % of 200",
		input:    "50 is what % of 200",
		expected: "50.00 is 25.00% of 200.00",
	},
	{
		name:     "decimal values",
		input:    "12.5 is what % of 50",
		expected: "12.50 is 25.00% of 50.00",
	},
	{
		name:     "without spaces around %",
		input:    "75 is what% of 300",
		expected: "75.00 is 25.00% of 300.00",
	},
	{
		name:     "without 'of'",
		input:    "100 is what % 400",
		expected: "100.00 is 25.00% of 400.00",
	},
	{
		name:     "30 is 20% of what",
		input:    "30 is 20% of what",
		expected: "30.00 is 20.00% of 150.00",
	},
	{
		name:     "50 is 25% of what",
		input:    "50 is 25% of what",
		expected: "50.00 is 25.00% of 200.00",
	},
	{
		name:     "decimal values",
		input:    "15 is 30% of what",
		expected: "15.00 is 30.00% of 50.00",
	},
	{
		name:     "without spaces around %",
		input:    "75 is 25% of what",
		expected: "75.00 is 25.00% of 300.00",
	},
	{
		name:     "without 'of'",
		input:    "40 is 20% what",
		expected: "40.00 is 20.00% of 200.00",
	},
}

// runParseTest runs a parse test with common validation logic
func runParseTest(t *testing.T, tests []struct {
	name     string
	input    string
	expected string
}) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.input, err)
			}
			if !strings.HasPrefix(result, tt.expected) {
				t.Errorf("Parse(%q) = %q, expected to start with %q", tt.input, result, tt.expected)
			}
			if !strings.Contains(result, "Steps:") {
				t.Errorf("Parse(%q) = %q, expected to contain calculation steps", tt.input, result)
			}
		})
	}
}
func TestParseXPercentOfY(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		commonTestCases[0], // "20% of 150"
		commonTestCases[1], // "what is 20% of 150"
		commonTestCases[2], // "50% of 200"
		commonTestCases[3], // "decimal percent"
		commonTestCases[4], // "decimal base"
		commonTestCases[5], // "without 'of'"
	}
	runParseTest(t, tests)
}

func TestParseXIsWhatPercentOfY(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		commonTestCases[6],  // "30 is what % of 150"
		commonTestCases[7],  // "50 is what % of 200"
		commonTestCases[8],  // "decimal values"
		commonTestCases[9],  // "without spaces around %"
		commonTestCases[10], // "without 'of'"
	}
	runParseTest(t, tests)
}

func TestParseXIsYPercentOfWhat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		commonTestCases[11], // "30 is 20% of what"
		commonTestCases[12], // "50 is 25% of what"
		commonTestCases[13], // "decimal values"
		commonTestCases[14], // "without spaces around %"
		commonTestCases[15], // "without 'of'"
	}
	runParseTest(t, tests)
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid input",
			input: "hello world",
		},
		{
			name:  "incomplete input",
			input: "20%",
		},
		{
			name:  "missing numbers",
			input: "% of",
		},
		{
			name:  "random text",
			input: "calculate percentage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err == nil {
				t.Errorf("Parse(%q) expected error, got nil", tt.input)
			}
		})
	}
}

func TestParseCaseInsensitive(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "uppercase WHAT IS",
			input: "WHAT IS 20% OF 150",
		},
		{
			name:  "mixed case What Is",
			input: "What Is 20% Of 150",
		},
		{
			name:  "uppercase IS WHAT",
			input: "30 IS WHAT % OF 150",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse(%q) should be case-insensitive, got error: %v", tt.input, err)
			}
		})
	}
}

func TestParseDivisionByZero(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "X is what % of 0",
			input: "30 is what % of 0",
		},
		{
			name:  "X is 0% of what",
			input: "30 is 0% of what",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err == nil {
				t.Errorf("Parse(%q) should handle division by zero, expected error", tt.input)
			}
		})
	}
}

func TestParseWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "extra spaces",
			input:    "  20%   of   150  ",
			expected: "20.00% of 150.00 = 30.00",
		},
		{
			name:     "tabs and spaces",
			input:    "30  is  what  %  of  150",
			expected: "30.00 is 20.00% of 150.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.input, err)
			}
			if !strings.Contains(result, "of") {
				t.Errorf("Parse(%q) should handle whitespace properly, got %q", tt.input, result)
			}
		})
	}
}
