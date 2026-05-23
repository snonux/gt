// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"strings"
	"testing"
)

func TestTTYCheckerEnsureTTYNotATerminal(t *testing.T) {
	// In test context, stdin is never a terminal
	checker := &TTYChecker{}
	err := checker.EnsureTTY()
	if err == nil {
		t.Error("EnsureTTY should return an error when stdin is not a TTY")
	}
}

func TestTTYCheckerEnsureTTYErrorMessage(t *testing.T) {
	checker := &TTYChecker{}
	err := checker.EnsureTTY()
	if err == nil {
		t.Skip("stdin appears to be a TTY; skipping error message check")
	}
	if !strings.Contains(err.Error(), "TTY") {
		t.Errorf("error message should mention TTY, got: %q", err.Error())
	}
}

func TestTTYCheckerIsTTYReturnsFalseInTests(t *testing.T) {
	// In test context, stdin is not a terminal
	checker := &TTYChecker{}
	if checker.IsTTY() {
		t.Skip("stdin appears to be a TTY (e.g., interactive terminal); skipping")
	}
	// Good — IsTTY correctly returned false
}

func TestTTYCheckerIsTTYConsistentWithEnsureTTY(t *testing.T) {
	checker := &TTYChecker{}
	isTTY := checker.IsTTY()
	err := checker.EnsureTTY()

	if isTTY && err != nil {
		t.Error("IsTTY returned true but EnsureTTY returned an error")
	}
	if !isTTY && err == nil {
		t.Error("IsTTY returned false but EnsureTTY returned nil")
	}
}

func TestTTYCheckerMultipleCalls(t *testing.T) {
	checker := &TTYChecker{}
	// Multiple calls should be consistent
	result1 := checker.IsTTY()
	result2 := checker.IsTTY()
	if result1 != result2 {
		t.Errorf("IsTTY returned inconsistent results: %v, %v", result1, result2)
	}

	err1 := checker.EnsureTTY()
	err2 := checker.EnsureTTY()
	// Both errors should be the same type (both nil or both non-nil)
	bothNil := err1 == nil && err2 == nil
	bothNonNil := err1 != nil && err2 != nil
	if !bothNil && !bothNonNil {
		t.Errorf("EnsureTTY returned inconsistent results: %v, %v", err1, err2)
	}
}
