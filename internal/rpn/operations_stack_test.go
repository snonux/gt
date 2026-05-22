// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"strings"
	"testing"
)

func TestShowWithMetrics(t *testing.T) {
	reg := GetMetricRegistry()
	vars := NewVariables()
	ops := NewOperations(vars)
	stack := NewStack()

	mbps, _ := reg.Find("Mbps")
	cool, _ := reg.Find("Cool")

	stack.Push(NewFloatWithMetric(100, mbps))
	stack.Push(NewFloatWithMetric(42, cool))
	stack.Push(NewFloatWithMetric(5.5, mbps))

	result, err := ops.Show(stack)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Non-Cool shows metric suffix, Cool doesn't
	if !strings.Contains(result, "100Mbps") {
		t.Errorf("expected '100Mbps' in result, got: %s", result)
	}
	if !strings.Contains(result, "42") {
		t.Errorf("expected '42' (plain) in result, got: %s", result)
	}
	if !strings.Contains(result, "5.5Mbps") {
		t.Errorf("expected '5.5Mbps' in result, got: %s", result)
	}
}
