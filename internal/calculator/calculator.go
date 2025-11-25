package calculator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

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

	return "", fmt.Errorf("unable to parse input. See usage for examples")
}

func parseXPercentOfY(input string) (string, bool) {
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*%\s*(?:of\s+)?(\d+(?:\.\d+)?)$`)
	matches := re.FindStringSubmatch(input)

	if matches == nil {
		return "", false
	}

	percent, _ := strconv.ParseFloat(matches[1], 64)
	base, _ := strconv.ParseFloat(matches[2], 64)

	result := (percent / 100.0) * base
	return fmt.Sprintf("%.2f%% of %.2f = %.2f", percent, base, result), true
}

func parseXIsWhatPercentOfY(input string) (string, bool) {
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s+is\s+what\s*%\s*(?:of\s+)?(\d+(?:\.\d+)?)$`)
	matches := re.FindStringSubmatch(input)

	if matches == nil {
		return "", false
	}

	part, _ := strconv.ParseFloat(matches[1], 64)
	whole, _ := strconv.ParseFloat(matches[2], 64)

	if whole == 0 {
		return "", false
	}

	percent := (part / whole) * 100.0
	return fmt.Sprintf("%.2f is %.2f%% of %.2f", part, percent, whole), true
}

func parseXIsYPercentOfWhat(input string) (string, bool) {
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s+is\s+(\d+(?:\.\d+)?)\s*%\s*(?:of\s+)?what$`)
	matches := re.FindStringSubmatch(input)

	if matches == nil {
		return "", false
	}

	part, _ := strconv.ParseFloat(matches[1], 64)
	percent, _ := strconv.ParseFloat(matches[2], 64)

	if percent == 0 {
		return "", false
	}

	whole := (part / percent) * 100.0
	return fmt.Sprintf("%.2f is %.2f%% of %.2f", part, percent, whole), true
}
