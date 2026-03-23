package calculator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"codeberg.org/snonux/perc/internal/rpn"
)

// Parse parses a percentage calculation input string and returns the result.
// It handles formats like "20% of 150", "30 is what % of 150", and "30 is 20% of what".
func Parse(input string) (string, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	input = strings.ReplaceAll(input, "what is ", "")
	input = strings.TrimSpace(input)

	if result, ok := parseXPercentOfY(input); ok {
		return result, nil
	}

	if result, ok := parseXIsWhatPercentOfY(input); ok {
		return result, nil
	}

	if result, ok := parseXIsYPercentOfWhat(input); ok {
		return result, nil
	}

	// Try RPN as a fallback
	result, err := ParseRPN(input)
	if err == nil {
		return result, nil
	}

	return "", fmt.Errorf("calculator: unable to parse input %q. See usage for examples", input)
}

// ParseRPN parses and evaluates an RPN (Reverse Polish Notation) expression.
// It handles formats like "3 4 +", "3 4 + 4 4 - *", "x 5 = x x +", etc.
func ParseRPN(input string) (string, error) {
	vars := rpn.NewVariables()
	rpnCalc := rpn.NewRPN(vars)
	return rpnCalc.ParseAndEvaluate(input)
}

func parseXPercentOfY(input string) (string, bool) {
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*%\s*(?:of\s+)?(\d+(?:\.\d+)?)$`)
	matches := re.FindStringSubmatch(input)

	if matches == nil {
		return "", false
	}

	percent, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return "", false
	}
	base, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return "", false
	}

	result := (percent / 100.0) * base

	output := fmt.Sprintf("%.2f%% of %.2f = %.2f\n", percent, base, result)
	output += fmt.Sprintf("  Steps: (%.2f / 100) * %.2f = %.2f * %.2f = %.2f", percent, base, percent/100.0, base, result)

	return output, true
}

func parseXIsWhatPercentOfY(input string) (string, bool) {
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s+is\s+what\s*%\s*(?:of\s+)?(\d+(?:\.\d+)?)$`)
	matches := re.FindStringSubmatch(input)

	if matches == nil {
		return "", false
	}

	part, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return "", false
	}
	whole, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return "", false
	}

	if whole == 0 {
		return "", false
	}

	percent := (part / whole) * 100.0

	output := fmt.Sprintf("%.2f is %.2f%% of %.2f\n", part, percent, whole)
	output += fmt.Sprintf("  Steps: (%.2f / %.2f) * 100 = %.2f * 100 = %.2f%%", part, whole, part/whole, percent)

	return output, true
}

func parseXIsYPercentOfWhat(input string) (string, bool) {
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s+is\s+(\d+(?:\.\d+)?)\s*%\s*(?:of\s+)?what$`)
	matches := re.FindStringSubmatch(input)

	if matches == nil {
		return "", false
	}

	part, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return "", false
	}
	percent, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return "", false
	}

	if percent == 0 {
		return "", false
	}

	whole := (part / percent) * 100.0

	output := fmt.Sprintf("%.2f is %.2f%% of %.2f\n", part, percent, whole)
	output += fmt.Sprintf("  Steps: (%.2f / %.2f) * 100 = %.2f * 100 = %.2f", part, percent, part/percent, whole)

	return output, true
}
