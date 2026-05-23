// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package perc

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

// TestParseCalculationPercentOfY tests ParseCalculation for "X% of Y" type
func TestParseCalculationPercentOfY(t *testing.T) {
	calc, err := ParseCalculation("20% of 150")
	if err != nil {
		t.Fatalf("ParseCalculation returned error: %v", err)
	}
	if calc == nil {
		t.Fatal("ParseCalculation returned nil")
	}
	if calc.Type != PercentOfY {
		t.Errorf("Type = %v, want PercentOfY", calc.Type)
	}
	if calc.Percent != 20 {
		t.Errorf("Percent = %v, want 20", calc.Percent)
	}
	if calc.Base != 150 {
		t.Errorf("Base = %v, want 150", calc.Base)
	}
	if calc.Result != 30 {
		t.Errorf("Result = %v, want 30", calc.Result)
	}
	if calc.Steps == "" {
		t.Error("Steps should not be empty")
	}
}

// TestParseCalculationIsWhatPercentOfY tests ParseCalculation for "X is what % of Y" type
func TestParseCalculationIsWhatPercentOfY(t *testing.T) {
	calc, err := ParseCalculation("30 is what % of 150")
	if err != nil {
		t.Fatalf("ParseCalculation returned error: %v", err)
	}
	if calc == nil {
		t.Fatal("ParseCalculation returned nil")
	}
	if calc.Type != IsWhatPercentOfY {
		t.Errorf("Type = %v, want IsWhatPercentOfY", calc.Type)
	}
	if calc.Percent != 20 {
		t.Errorf("Percent = %v, want 20", calc.Percent)
	}
	if calc.Base != 150 {
		t.Errorf("Base = %v, want 150", calc.Base)
	}
	if calc.Result != 30 {
		t.Errorf("Result = %v, want 30", calc.Result)
	}
}

// TestParseCalculationIsYPercentOfWhat tests ParseCalculation for "X is Y% of what" type
func TestParseCalculationIsYPercentOfWhat(t *testing.T) {
	calc, err := ParseCalculation("30 is 20% of what")
	if err != nil {
		t.Fatalf("ParseCalculation returned error: %v", err)
	}
	if calc == nil {
		t.Fatal("ParseCalculation returned nil")
	}
	if calc.Type != IsYPercentOfWhat {
		t.Errorf("Type = %v, want IsYPercentOfWhat", calc.Type)
	}
	if calc.Percent != 20 {
		t.Errorf("Percent = %v, want 20", calc.Percent)
	}
	if calc.Base != 150 {
		t.Errorf("Base = %v, want 150", calc.Base)
	}
	if calc.Result != 30 {
		t.Errorf("Result = %v, want 30", calc.Result)
	}
}

// TestParseCalculationWhatIsPrefix tests "what is" prefix stripping
func TestParseCalculationWhatIsPrefix(t *testing.T) {
	calc, err := ParseCalculation("what is 20% of 150")
	if err != nil {
		t.Fatalf("ParseCalculation returned error: %v", err)
	}
	if calc.Percent != 20 || calc.Base != 150 || calc.Result != 30 {
		t.Errorf("ParseCalculation with 'what is' prefix: Percent=%v, Base=%v, Result=%v",
			calc.Percent, calc.Base, calc.Result)
	}
}

// TestParseCalculationCaseInsensitive tests case insensitivity
func TestParseCalculationCaseInsensitive(t *testing.T) {
	tests := []struct {
		input    string
		wantType CalculationType
	}{
		{"20% OF 150", PercentOfY},
		{"WHAT IS 20% OF 150", PercentOfY},
		{"30 IS WHAT % OF 150", IsWhatPercentOfY},
		{"30 IS 20% OF WHAT", IsYPercentOfWhat},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			calc, err := ParseCalculation(tt.input)
			if err != nil {
				t.Fatalf("ParseCalculation(%q) returned error: %v", tt.input, err)
			}
			if calc.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", calc.Type, tt.wantType)
			}
		})
	}
}

// TestParseCalculationDivisionByZero tests division by zero edge cases
func TestParseCalculationDivisionByZero(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"X is what % of 0", "30 is what % of 0"},
		{"X is 0% of what", "30 is 0% of what"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCalculation(tt.input)
			if err == nil {
				t.Errorf("ParseCalculation(%q) should return error for division by zero", tt.input)
			}
		})
	}
}

// TestParseCalculationInvalidInput tests invalid inputs return errors
func TestParseCalculationInvalidInput(t *testing.T) {
	tests := []string{
		"hello world",
		"20%",
		"% of",
		"calculate percentage",
		"",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := ParseCalculation(input)
			if err == nil {
				t.Errorf("ParseCalculation(%q) should return error", input)
			}
		})
	}
}

// TestCalculationFormatPercentOfY tests Format() for PercentOfY type
func TestCalculationFormatPercentOfY(t *testing.T) {
	calc := &Calculation{
		Type:    PercentOfY,
		Percent: 20,
		Base:    150,
		Result:  30,
		Steps:   "(20 / 100) * 150 = 0.20 * 150 = 30",
	}

	formatted := calc.Format()
	if !strings.HasPrefix(formatted, "20.00% of 150.00 = 30.00") {
		t.Errorf("Format() = %q, expected to start with '20.00%% of 150.00 = 30.00'", formatted)
	}
	if !strings.Contains(formatted, "Steps:") {
		t.Error("Format() should contain Steps")
	}
}

// TestCalculationFormatIsWhatPercentOfY tests Format() for IsWhatPercentOfY type
func TestCalculationFormatIsWhatPercentOfY(t *testing.T) {
	calc := &Calculation{
		Type:    IsWhatPercentOfY,
		Percent: 20,
		Base:    150,
		Result:  30,
		Steps:   "",
	}

	formatted := calc.Format()
	if !strings.HasPrefix(formatted, "30.00 is 20.00% of 150.00") {
		t.Errorf("Format() = %q, expected to start with '30.00 is 20.00%% of 150.00'", formatted)
	}
}

// TestCalculationFormatIsYPercentOfWhat tests Format() for IsYPercentOfWhat type
func TestCalculationFormatIsYPercentOfWhat(t *testing.T) {
	calc := &Calculation{
		Type:    IsYPercentOfWhat,
		Percent: 20,
		Base:    150,
		Result:  30,
		Steps:   "",
	}

	formatted := calc.Format()
	if !strings.HasPrefix(formatted, "30.00 is 20.00% of 150.00") {
		t.Errorf("Format() = %q, expected to start with '30.00 is 20.00%% of 150.00'", formatted)
	}
}

// TestCalculationFormatWithoutSteps tests Format() without Steps field
func TestCalculationFormatWithoutSteps(t *testing.T) {
	calc := &Calculation{
		Type:    PercentOfY,
		Percent: 25,
		Base:    100,
		Result:  25,
		Steps:   "",
	}

	formatted := calc.Format()
	if formatted != "25.00% of 100.00 = 25.00" {
		t.Errorf("Format() without steps = %q, want '25.00%% of 100.00 = 25.00'", formatted)
	}
	if strings.Contains(formatted, "Steps") {
		t.Error("Format() should not contain 'Steps' when Steps is empty")
	}
}

// TestParseCalculationDecimalValues tests decimal values in ParseCalculation
func TestParseCalculationDecimalValues(t *testing.T) {
	calc, err := ParseCalculation("12.5% of 80")
	if err != nil {
		t.Fatalf("ParseCalculation returned error: %v", err)
	}
	if calc.Percent != 12.5 {
		t.Errorf("Percent = %v, want 12.5", calc.Percent)
	}
	if calc.Result != 10 {
		t.Errorf("Result = %v, want 10", calc.Result)
	}
}

// TestParseCalculationLargeValues tests large values in ParseCalculation
func TestParseCalculationLargeValues(t *testing.T) {
	calc, err := ParseCalculation("50% of 1000000")
	if err != nil {
		t.Fatalf("ParseCalculation returned error: %v", err)
	}
	if calc.Percent != 50 {
		t.Errorf("Percent = %v, want 50", calc.Percent)
	}
	if calc.Result != 500000 {
		t.Errorf("Result = %v, want 500000", calc.Result)
	}
}
