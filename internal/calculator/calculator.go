// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package calculator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// CalculationType represents the type of calculation performed.
type CalculationType int

const (
	// PercentOfY: "X% of Y" → "X.00% of Y.00 = Z.00"
	PercentOfY CalculationType = iota
	// IsWhatPercentOfY: "X is what % of Y" → "X.00 is P.00% of Y.00"
	IsWhatPercentOfY
	// IsYPercentOfWhat: "X is Y% of what" → "X.00 is Y.00% of W.00"
	IsYPercentOfWhat
)

// Calculation represents the result of a percentage calculation.
type Calculation struct {
	Type    CalculationType
	Percent float64
	Base    float64
	Result  float64
	Steps   string
}

// Format returns the formatted calculation result.
func (c *Calculation) Format() string {
	var baseStr string
	switch c.Type {
	case PercentOfY:
		baseStr = fmt.Sprintf("%.2f%% of %.2f = %.2f", c.Percent, c.Base, c.Result)
	case IsWhatPercentOfY:
		// percent is the result, base is the "whole"
		baseStr = fmt.Sprintf("%.2f is %.2f%% of %.2f", c.Result, c.Percent, c.Base)
	case IsYPercentOfWhat:
		// percent is the known value, base is the "what"
		baseStr = fmt.Sprintf("%.2f is %.2f%% of %.2f", c.Result, c.Percent, c.Base)
	}
	if c.Steps != "" {
		return baseStr + "\n  Steps: " + c.Steps
	}
	return baseStr
}

// ParsingStrategy represents a parsing function that attempts to parse input.
// Returns a Calculation if handled, or error if not.
type ParsingStrategy func(input string) (*Calculation, bool, error)

// strategyRegistry maintains a registry of parsing strategies.
type strategyRegistry struct {
	strategies []ParsingStrategy
}

// newStrategyRegistry creates a new strategy registry.
func newStrategyRegistry() *strategyRegistry {
	return &strategyRegistry{
		strategies: make([]ParsingStrategy, 0),
	}
}

// register adds a parsing strategy to the registry.
func (r *strategyRegistry) register(strategy ParsingStrategy) {
	r.strategies = append(r.strategies, strategy)
}

// parse attempts to parse input using registered strategies in order.
func (r *strategyRegistry) parse(input string) (*Calculation, bool, error) {
	for _, strategy := range r.strategies {
		if result, handled, err := strategy(input); handled {
			return result, true, err
		}
	}
	return nil, false, nil
}

// Parse parses a percentage calculation input string and returns the result as a formatted string.
// It handles formats like "20% of 150", "30 is what % of 150", and "30 is 20% of what".
// Note: This function only handles percentage calculations, not RPN expressions.
func Parse(input string) (string, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	input = strings.ReplaceAll(input, "what is ", "")
	input = strings.TrimSpace(input)

	// Create registry and register percentage parsing strategies
	registry := newStrategyRegistry()
	registry.register(parseXPercentOfY)
	registry.register(parseXIsWhatPercentOfY)
	registry.register(parseXIsYPercentOfWhat)

	calc, ok, err := registry.parse(input)
	if ok {
		return calc.Format(), nil
	}
	if err != nil {
		return "", fmt.Errorf("calculator: unable to parse input %q: %w", input, err)
	}

	return "", fmt.Errorf("calculator: unable to parse input %q: unknown error", input)
}

// ParseCalculation parses a percentage calculation input string and returns the Calculation object.
// It handles formats like "20% of 150", "30 is what % of 150", and "30 is 20% of what".
// This provides callers with more flexibility to access raw values and formatting options.
func ParseCalculation(input string) (*Calculation, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	input = strings.ReplaceAll(input, "what is ", "")
	input = strings.TrimSpace(input)

	// Create registry and register percentage parsing strategies
	registry := newStrategyRegistry()
	registry.register(parseXPercentOfY)
	registry.register(parseXIsWhatPercentOfY)
	registry.register(parseXIsYPercentOfWhat)

	calc, ok, err := registry.parse(input)
	if ok {
		return calc, nil
	}
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("calculator: unable to parse input %q. See usage for examples", input)
}

// parseXPercentOfY calculates "X% of Y" and returns a Calculation.
func parseXPercentOfY(input string) (*Calculation, bool, error) {
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*%\s*(?:of\s+)?(\d+(?:\.\d+)?)$`)
	matches := re.FindStringSubmatch(input)

	if matches == nil {
		return nil, false, nil
	}

	percent, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return nil, false, err
	}
	base, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return nil, false, err
	}

	result := (percent / 100.0) * base

	calc := &Calculation{
		Type:    PercentOfY,
		Percent: percent,
		Base:    base,
		Result:  result,
		Steps:   fmt.Sprintf("(%.2f / 100) * %.2f = %.2f * %.2f = %.2f", percent, base, percent/100.0, base, result),
	}

	return calc, true, nil
}

// parseXIsWhatPercentOfY calculates "X is what % of Y" and returns a Calculation.
func parseXIsWhatPercentOfY(input string) (*Calculation, bool, error) {
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s+is\s+what\s*%\s*(?:of\s+)?(\d+(?:\.\d+)?)$`)
	matches := re.FindStringSubmatch(input)

	if matches == nil {
		return nil, false, nil
	}

	part, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return nil, false, err
	}
	whole, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return nil, false, err
	}

	if whole == 0 {
		return nil, false, fmt.Errorf("division by zero")
	}

	percent := (part / whole) * 100.0

	calc := &Calculation{
		Type:    IsWhatPercentOfY,
		Percent: percent,
		Base:    whole,
		Result:  part,
		Steps:   fmt.Sprintf("(%.2f / %.2f) * 100 = %.2f * 100 = %.2f%%", part, whole, part/whole, percent),
	}

	return calc, true, nil
}

// parseXIsYPercentOfWhat calculates "X is Y% of what" and returns a Calculation.
func parseXIsYPercentOfWhat(input string) (*Calculation, bool, error) {
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s+is\s+(\d+(?:\.\d+)?)\s*%\s*(?:of\s+)?what$`)
	matches := re.FindStringSubmatch(input)

	if matches == nil {
		return nil, false, nil
	}

	part, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return nil, false, err
	}
	percent, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return nil, false, err
	}

	if percent == 0 {
		return nil, false, fmt.Errorf("division by zero")
	}

	whole := (part / percent) * 100.0

	calc := &Calculation{
		Type:    IsYPercentOfWhat,
		Percent: percent,
		Base:    whole,
		Result:  part,
		Steps:   fmt.Sprintf("(%.2f / %.2f) * 100 = %.2f * 100 = %.2f", part, percent, part/percent, whole),
	}

	return calc, true, nil
}
