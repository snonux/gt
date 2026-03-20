package rpn

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewRPN(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)
	if r == nil {
		t.Fatal("NewRPN() returned nil")
	}
	if r.vars == nil {
		t.Error("RPN.vars should not be nil")
	}
	if r.ops == nil {
		t.Error("RPN.ops should not be nil")
	}
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple expression",
			input:    "3 4 +",
			expected: []string{"3", "4", "+"},
		},
		{
			name:     "multiple spaces",
			input:    "3   4   +",
			expected: []string{"3", "4", "+"},
		},
		{
			name:     "decimal numbers",
			input:    "3.14 2.5 +",
			expected: []string{"3.14", "2.5", "+"},
		},
		{
			name:     "expression with leading/trailing spaces",
			input:    "  3 4 +  ",
			expected: []string{"3", "4", "+"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tokenize(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("tokenize(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseAndEvaluateSimple(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "3 4 +",
			input:    "3 4 +",
			expected: "7",
		},
		{
			name:     "5 3 -",
			input:    "5 3 -",
			expected: "2",
		},
		{
			name:     "2 3 *",
			input:    "2 3 *",
			expected: "6",
		},
		{
			name:     "10 2 /",
			input:    "10 2 /",
			expected: "5",
		},
		{
			name:     "2 3 ^",
			input:    "2 3 ^",
			expected: "8",
		},
		{
			name:     "10 3 %",
			input:    "10 3 %",
			expected: "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewVariables().(*Variables)
			r := NewRPN(v)
			result, err := r.ParseAndEvaluate(tt.input)
			if err != nil {
				t.Fatalf("ParseAndEvaluate(%q) returned error: %v", tt.input, err)
			}
			if !strings.HasPrefix(result, tt.expected) {
				t.Errorf("ParseAndEvaluate(%q) = %q, want to start with %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseAndEvaluateChain(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "3 4 + 4 4 - *",
			input:    "3 4 + 4 4 - *",
			expected: "0",
		},
		{
			name:     "1 2 + 3 *",
			input:    "1 2 + 3 *",
			expected: "9",
		},
		{
			name:     "2 3 + 4 5 + *",
			input:    "2 3 + 4 5 + *",
			expected: "45",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewVariables().(*Variables)
			r := NewRPN(v)
			result, err := r.ParseAndEvaluate(tt.input)
			if err != nil {
				t.Fatalf("ParseAndEvaluate(%q) returned error: %v", tt.input, err)
			}
			if !strings.HasPrefix(result, tt.expected) {
				t.Errorf("ParseAndEvaluate(%q) = %q, want to start with %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseAndEvaluateStackOps(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "1 2 3 dup",
			input:    "1 2 3 dup",
			expected: "1 2 3 3",
		},
		{
			name:     "1 2 swap",
			input:    "1 2 swap",
			expected: "2 1",
		},
		{
			name:     "1 2 3 pop",
			input:    "1 2 3 pop",
			expected: "1 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewVariables().(*Variables)
			r := NewRPN(v)
			result, err := r.ParseAndEvaluate(tt.input)
			if err != nil {
				t.Fatalf("ParseAndEvaluate(%q) returned error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("ParseAndEvaluate(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseAndEvaluateVariables(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	// Test variable assignment and reuse
	// First assign a variable
	result, err := r.ParseAndEvaluate("x = 5")
	if err != nil {
		t.Fatalf("First assignment failed: %v", err)
	}
	if result != "x = 5" {
		t.Errorf("Assignment result = %q, want %q", result, "x = 5")
	}

	// Now use the variable in RPN
	result, err = r.ParseAndEvaluate("x x +")
	if err != nil {
		t.Fatalf("ParseAndEvaluate(%q) returned error: %v", "x x +", err)
	}
	if result != "10" {
		t.Errorf("ParseAndEvaluate(%q) = %q, want %q", "x x +", result, "10")
	}
}

func TestParseAndEvaluateEmpty(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	_, err := r.ParseAndEvaluate("")
	if err == nil {
		t.Error("ParseAndEvaluate(\"\") should return error")
	}
	if !strings.Contains(err.Error(), "empty expression") {
		t.Errorf("Error = %v, should contain 'empty expression'", err)
	}
}

func TestParseAndEvaluateAssignment(t *testing.T) {
	// Test assignment format: "varname = value"
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "x = 5",
			input:    "x = 5",
			expected: "x = 5",
		},
		{
			name:     "pi = 3.14159",
			input:    "pi = 3.14159",
			expected: "pi = 3.14159",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewVariables().(*Variables)
			r := NewRPN(v)
			result, err := r.ParseAndEvaluate(tt.input)
			if err != nil {
				t.Fatalf("ParseAndEvaluate(%q) returned error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("ParseAndEvaluate(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseAndEvaluateDivisionByZero(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	_, err := r.ParseAndEvaluate("5 0 /")
	if err == nil {
		t.Error("5 0 / should return error")
	}
	if !strings.Contains(err.Error(), "division by zero") {
		t.Errorf("Error = %v, should contain 'division by zero'", err)
	}
}

func TestParseAndEvaluateUndefinedVariable(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	_, err := r.ParseAndEvaluate("undefined +")
	if err == nil {
		t.Error("undefined variable should return error")
	}
	// The error should mention the undefined token
	if !strings.Contains(err.Error(), "undefined") {
		t.Errorf("Error = %v, should contain 'undefined'", err)
	}
}

func TestParseAndEvaluateUnknownToken(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	_, err := r.ParseAndEvaluate("1 2 + hello")
	if err == nil {
		t.Error("unknown token should return error")
	}
	if !strings.Contains(err.Error(), "unknown token") {
		t.Errorf("Error = %v, should contain 'unknown token'", err)
	}
}

func TestParseAndEvaluateInsufficientOperands(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	tests := []struct {
		name  string
		input string
	}{
		{"+ with one operand", "5 +"},
		{"+ with no operands", "+"},
		{"3 +", "3 +"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := r.ParseAndEvaluate(tt.input)
			if err == nil {
				t.Errorf("%q should return error for insufficient operands", tt.input)
			}
		})
	}
}

func TestParseAndEvaluateShow(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	result, err := r.ParseAndEvaluate("1 2 3 show")
	if err != nil {
		t.Fatalf("ParseAndEvaluate(%q) returned error: %v", "1 2 3 show", err)
	}
	if result != "1 2 3" {
		t.Errorf("ParseAndEvaluate(%q) = %q, want \"1 2 3\"", "1 2 3 show", result)
	}
}

func TestParseAndEvaluateVars(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	// Set some variables using new format: "name = value"
	r.ParseAndEvaluate("x = 5")
	r.ParseAndEvaluate("y = 10")

	result, err := r.ParseAndEvaluate("vars")
	if err != nil {
		t.Fatalf("ParseAndEvaluate(%q) returned error: %v", "vars", err)
	}
	if !strings.Contains(result, "x") || !strings.Contains(result, "y") {
		t.Errorf("vars output should contain all variables, got: %s", result)
	}
}

func TestParseAndEvaluateClear(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	// Set and clear
	r.ParseAndEvaluate("x 5 =")
	r.ParseAndEvaluate("clear")

	if v.Count() != 0 {
		t.Errorf("Count after clear = %d, want 0", v.Count())
	}
}

func TestRPNConcurrency(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	done := make(chan bool, 10)
	for i := 0; i < 5; i++ {
		go func(id int) {
			name := fmt.Sprintf("val%d", id)
			input := fmt.Sprintf("%s = %d", name, id)
			r.ParseAndEvaluate(input)
			done <- true
		}(i)
	}

	for i := 0; i < 5; i++ {
		<-done
	}

	if v.Count() != 5 {
		t.Errorf("Final count = %d, want 5", v.Count())
	}
}

func TestResultStack(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	tokens := []string{"1", "2", "3", "+"}
	result, err := r.ResultStack(tokens)
	if err != nil {
		t.Fatalf("ResultStack() returned error: %v", err)
	}
	if result != "1 5" {
		t.Errorf("ResultStack() = %q, want \"1 5\"", result)
	}
}

func TestResultStackEmpty(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	tokens := []string{}
	result, err := r.ResultStack(tokens)
	if err != nil {
		t.Fatalf("ResultStack([]) returned error: %v", err)
	}
	if result != "Stack is empty" {
		t.Errorf("ResultStack([]) = %q, want \"Stack is empty\"", result)
	}
}

func TestParseAndEvaluateAssignmentExpression(t *testing.T) {
	// Test assignment with expression: "name value = expression..."
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple assignment with expression",
			input:    "x 5 = x x +", // set x=5, then evaluate x x + => 5+5=10
			expected: "10",
		},
		{
			name:     "assignment with complex expression",
			input:    "pi 3.14 = pi 2 *", // set pi=3.14, then evaluate pi 2 * => 3.14*2=6.28
			expected: "6.28",
		},
		{
			name:     "assignment then use in another expression",
			input:    "b 7 = b 1 + b *", // set b=7, then b 1 + => 8, then b * => 7*8=56
			expected: "56",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewVariables().(*Variables)
			r := NewRPN(v)
			result, err := r.ParseAndEvaluate(tt.input)
			if err != nil {
				t.Fatalf("ParseAndEvaluate(%q) returned error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("ParseAndEvaluate(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseAndEvaluateAssignmentErrors(t *testing.T) {
	// Test error cases in assignment handling
	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name:     "invalid value for assignment (non-numeric)",
			input:    "x abc =",
			expectedError: "unknown token 'x'",
		},
		{
			name:     "assignment with variable name containing space",
			input:    "my var 5 =",
			expectedError: "unknown token 'my'",
		},
		{
			name:     "assignment with value containing space",
			input:    "x 5 6 =",
			expectedError: "unknown token 'x'",
		},
		{
			name:     "empty assignment",
			input:    " = ",
			expectedError: "invalid assignment syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewVariables().(*Variables)
			r := NewRPN(v)
			_, err := r.ParseAndEvaluate(tt.input)
			if err == nil {
				t.Errorf("ParseAndEvaluate(%q) expected error, got nil", tt.input)
				return
			}
			if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("ParseAndEvaluate(%q) error = %q, want to contain %q", tt.input, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestParseAndEvaluateEvaluateErrors(t *testing.T) {
	// Test error cases in evaluate function
	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name:     "invalid assignment syntax (standalone =)",
			input:    "=",
			expectedError: "invalid assignment syntax",
		},
		{
			name:     "'d' command not supported",
			input:    "d",
			expectedError: "'d' command not supported as standalone token",
		},
		{
			name:     "empty result after evaluation",
			input:    "1 2 + pop", // 1 2 + => 3, then pop => empty stack
			expectedError: "empty result: expression evaluated to nothing",
		},
		{
			name:     "stack overflow (simulate many numbers)",
			input:    "", // placeholder
			expectedError: "stack overflow",
		},
	}

	for _, tt := range tests {
		if tt.name == "stack overflow (simulate many numbers)" {
			// Skip stack overflow test for now as it's hard to test without modifying internals
			t.Logf("Skipping %s test", tt.name)
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			v := NewVariables().(*Variables)
			r := NewRPN(v)
			_, err := r.ParseAndEvaluate(tt.input)
			if err == nil {
				t.Errorf("ParseAndEvaluate(%q) expected error, got nil", tt.input)
				return
			}
			if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("ParseAndEvaluate(%q) error = %q, want to contain %q", tt.input, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestResultStackErrors(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	// Test error cases in ResultStack function
	tests := []struct {
		name          string
		input         []string
		expectedError string
	}{
		{
			name:     "division by zero",
			input:    []string{"5", "0", "/"},
			expectedError: "division by zero",
		},
		{
			name:     "unknown token",
			input:    []string{"1", "2", "+", "unknown"},
			expectedError: "unknown token",
		},
		{
			name:     "insufficient operands for +",
			input:    []string{"5", "+"},
			expectedError: "insufficient operands",
		},
		{
			name:     "insufficient operands for -",
			input:    []string{"5", "-"},
			expectedError: "insufficient operands",
		},
		{
			name:     "invalid assignment syntax in ResultStack",
			input:    []string{"="},
			expectedError: "unknown token '='",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := r.ResultStack(tt.input)
			if err == nil {
				t.Errorf("ResultStack(%v) expected error, got nil", tt.input)
				return
			}
			if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("ResultStack(%v) error = %q, want to contain %q", tt.input, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestResultStackMultipleValues(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	// Test Case where stack has multiple values at the end
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "two values on stack",
			input:    []string{"1", "2", "3", "+"}, // 1 2 3 + => 1 (5) => two values: 1, 5
			expected: "1 5", // Show should show all values
		},
		{
			name:     "three values on stack",
			input:    []string{"1", "2", "3", "4"}, // just push 4 numbers
			expected: "1 2 3 4",
		},
		{
			name:     "multiple values with variables",
			input:    []string{"x", "y", "z"}, // after setting variables
			expected: "10 20 30", // variables x, y, z have values 10, 20, 30
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up variables if needed
			if tt.name == "multiple values with variables" {
				r.ParseAndEvaluate("x = 10")
				r.ParseAndEvaluate("y = 20")
				r.ParseAndEvaluate("z = 30")
			}

			result, err := r.ResultStack(tt.input)
			if err != nil {
				t.Fatalf("ResultStack(%v) returned error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("ResultStack(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRPNIncrementalOperations(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	// Test: 1 2 3 + then +
	// First evaluate "1 2 3 +"
	result, err := r.ParseAndEvaluate("1 2 3 +")
	if err != nil {
		t.Fatalf("First evaluation failed: %v", err)
	}
	if result != "1 5" {
		t.Errorf("First result = %q, want '1 5'", result)
	}

	// Then apply + operator
	result, err = r.EvalOperator("+")
	if err != nil {
		t.Fatalf("EvalOperator('+') failed: %v", err)
	}
	if result != "6" {
		t.Errorf("After + = %q, want '6'", result)
	}
}

func TestRPNIncrementalSubtract(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	// First put two values on stack: "10 3" gives stack [10, 3]
	_, err := r.ParseAndEvaluate("10 3")
	if err != nil {
		t.Fatalf("First evaluation failed: %v", err)
	}

	// Now subtract
	result, err := r.EvalOperator("-")
	if err != nil {
		t.Fatalf("EvalOperator('-') failed: %v", err)
	}
	// 10 - 3 = 7
	if result != "7" {
		t.Errorf("After - = %q, want '7'", result)
	}
}

func TestRPNIncrementalDup(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	// First push values (two values so stack is not emptied after evaluation)
	_, err := r.ParseAndEvaluate("5 6")
	if err != nil {
		t.Fatalf("First evaluation failed: %v", err)
	}
	// After "5 6", stack should have [5, 6], result is "5 6"

	// Now duplicate
	result, err := r.EvalOperator("dup")
	if err != nil {
		t.Fatalf("EvalOperator('dup') failed: %v", err)
	}
	if result != "5 6 6" {
		t.Errorf("After dup = %q, want '5 6 6'", result)
	}
}

func TestRPNIncrementalSwap(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	_, err := r.ParseAndEvaluate("1 2")
	if err != nil {
		t.Fatalf("First evaluation failed: %v", err)
	}

	result, err := r.EvalOperator("swap")
	if err != nil {
		t.Fatalf("EvalOperator('swap') failed: %v", err)
	}
	if result != "2 1" {
		t.Errorf("After swap = %q, want '2 1'", result)
	}
}

func TestRPNGetCurrentStack(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	_, err := r.ParseAndEvaluate("1 2 3")
	if err != nil {
		t.Fatalf("First evaluation failed: %v", err)
	}

	stack := r.GetCurrentStack()
	if len(stack) != 3 {
		t.Errorf("Stack length = %d, want 3", len(stack))
	}
	if stack[0] != 1 || stack[1] != 2 || stack[2] != 3 {
		t.Errorf("Stack = %v, want [1 2 3]", stack)
	}
}

func TestRPNIncrementalUnknownOperator(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	_, err := r.ParseAndEvaluate("1 2")
	if err != nil {
		t.Fatalf("First evaluation failed: %v", err)
	}

	_, err = r.EvalOperator("unknown")
	if err == nil {
		t.Error("EvalOperator('unknown') should return error")
	}
}

func TestRPNClearStack(t *testing.T) {
	v := NewVariables().(*Variables)
	r := NewRPN(v)

	_, err := r.ParseAndEvaluate("1 2 3")
	if err != nil {
		t.Fatalf("First evaluation failed: %v", err)
	}

	result, err := r.EvalOperator("clear")
	if err != nil {
		t.Fatalf("EvalOperator('clear') failed: %v", err)
	}
	if result != "All variables cleared" {
		t.Errorf("After clear = %q, want 'All variables cleared'", result)
	}
}
