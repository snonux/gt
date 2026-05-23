// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import "strconv"

// parseNumberWithMetric attempts to split a token into a number and a metric suffix.
// E.g., "100Mbps" -> (100, Mbps), "5.5GB" -> (5.5, GB), "2hr" -> (2, hr).
// The metric suffix is looked up in the given registry (exact, alias, then case-insensitive).
// Returns (num, metric, true) if successful, or (0, nil, false) if the token
// does not contain a number+metric combination.
func parseNumberWithMetric(token string, reg *MetricRegistry) (float64, *Metric, bool) {
	if len(token) == 0 {
		return 0, nil, false
	}

	// Skip leading sign
	i := 0
	if token[i] == '-' || token[i] == '+' {
		i++
	}

	// Must have at least one digit after optional sign
	if i >= len(token) {
		return 0, nil, false
	}

	// Scan the numeric portion (digits, decimal point, e/E for scientific notation)
	start := i
	for i < len(token) {
		c := token[i]
		if c >= '0' && c <= '9' {
			i++
		} else if c == '.' {
			i++
		} else if c == 'e' || c == 'E' {
			i++
			if i < len(token) && (token[i] == '+' || token[i] == '-') {
				i++
			}
		} else {
			break
		}
	}

	// Must have consumed at least one character in the numeric portion;
	// actual digit validation is deferred to ParseFloat.
	if i == start {
		return 0, nil, false
	}

	// Must have a non-empty suffix
	if i >= len(token) {
		return 0, nil, false
	}

	numStr := token[:i]
	metricName := token[i:]

	// Early pre-check: skip unlikely metric suffixes to avoid
	// unnecessary registry lookups for tokens like "10x".
	// All built-in metric names are >= 2 chars, except for three
	// single-char metrics: s (seconds), m (meters), g (grams).
	// Require suffix length >= 2, or be one of those known singles.
	if len(metricName) < 2 {
		if metricName != "s" && metricName != "m" && metricName != "g" {
			return 0, nil, false
		}
	}

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, nil, false
	}

	metric, ok := reg.FindWithAliases(metricName)
	if !ok {
		return 0, nil, false
	}

	return num, metric, true
}
