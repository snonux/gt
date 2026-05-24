// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import "testing"

// TestParseCategoryKnown tests that known category names resolve correctly.
func TestParseCategoryKnown(t *testing.T) {
	tests := []struct {
		name string
		want Category
	}{
		{"Universal", Universal},
		{"DataRate", DataRate},
		{"DataSize", DataSize},
		{"Time", Time},
		{"Weight", Weight},
		{"Speed", Speed},
		{"Distance", Distance},
		{"Custom", Custom},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseCategory(tt.name)
			if !ok {
				t.Errorf("parseCategory(%q) = !ok, want ok", tt.name)
			}
			if got != tt.want {
				t.Errorf("parseCategory(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

// TestParseCategoryUnknown does not return Universal.
// If parseCategory returned Category(0) on failure, callers ignoring the
// bool return would silently get Universal — this test prevents that regression.
func TestParseCategoryUnknown(t *testing.T) {
	got, ok := parseCategory("DoesNotExist")
	if ok {
		t.Fatal("parseCategory(\"DoesNotExist\") = ok, want !ok")
	}
	if got == Universal {
		t.Errorf("parseCategory(\"DoesNotExist\") returned Universal (Category(0)); "+
			"use a sentinel like Category(-1) so callers can't mistake it for a valid category")
	}
	if got != invalidCategory {
		t.Errorf("parseCategory(\"DoesNotExist\") = %v, want %v (invalidCategory)", got, invalidCategory)
	}
}

// TestParseCategoryEmpty tests that an empty string is not treated as a valid category.
func TestParseCategoryEmpty(t *testing.T) {
	got, ok := parseCategory("")
	if ok {
		t.Fatal("parseCategory(\"\") = ok, want !ok")
	}
	if got == Universal {
		t.Error("parseCategory(\"\") must not return Universal")
	}
}

// TestInvalidCategorySentinel verifies the sentinel is not a valid category.
func TestInvalidCategorySentinel(t *testing.T) {
	if invalidCategory >= 0 && invalidCategory < _sentinel {
		t.Errorf("invalidCategory (%d) overlaps with valid range [0, %d)", invalidCategory, _sentinel)
	}
	if invalidCategory == Universal {
		t.Error("invalidCategory must never equal Universal")
	}
}
