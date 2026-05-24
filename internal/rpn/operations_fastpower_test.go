// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"math"
	"strings"
	"testing"
)

func TestFastPower(t *testing.T) {
	tests := []struct {
		name       string
		prepare    func(*Stack)
		wantErr    bool
		wantErrSub string
		want       float64
	}{
		{
			name:    "2 ** 10 = 1024",
			prepare: func(s *Stack) { s.Push(NewNumber(2.0, FloatMode)); s.Push(NewNumber(10.0, FloatMode)) },
			want:    1024.0,
		},
		{
			name:    "3 ** 4 = 81",
			prepare: func(s *Stack) { s.Push(NewNumber(3.0, FloatMode)); s.Push(NewNumber(4.0, FloatMode)) },
			want:    81.0,
		},
		{
			name:    "5 ** 0 = 1",
			prepare: func(s *Stack) { s.Push(NewNumber(5.0, FloatMode)); s.Push(NewNumber(0.0, FloatMode)) },
			want:    1.0,
		},
		{
			name:    "2 ** -3 = 0.125",
			prepare: func(s *Stack) { s.Push(NewNumber(2.0, FloatMode)); s.Push(NewNumber(-3.0, FloatMode)) },
			want:    0.125,
		},
		{
			name:    "10 ** -2 = 0.01",
			prepare: func(s *Stack) { s.Push(NewNumber(10.0, FloatMode)); s.Push(NewNumber(-2.0, FloatMode)) },
			want:    0.01,
		},
		{
			name:    "1.5 ** 3 = 3.375",
			prepare: func(s *Stack) { s.Push(NewNumber(1.5, FloatMode)); s.Push(NewNumber(3.0, FloatMode)) },
			want:    3.375,
		},
		{
			name:    "-2 ** 3 = -8",
			prepare: func(s *Stack) { s.Push(NewNumber(-2.0, FloatMode)); s.Push(NewNumber(3.0, FloatMode)) },
			want:    -8.0,
		},
		{
			name:       "non-integer exponent",
			prepare:    func(s *Stack) { s.Push(NewNumber(2.0, FloatMode)); s.Push(NewNumber(1.5, FloatMode)) },
			wantErr:    true,
			wantErrSub: "must be an integer",
		},
		{
			name:    "empty stack",
			prepare: func(s *Stack) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := NewStack()
			tt.prepare(stack)

			o := NewOperations(NewVariables(), nil)
			err := o.FastPower(stack)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.wantErrSub != "" && !strings.Contains(err.Error(), tt.wantErrSub) {
					t.Errorf("error = %q, want substring %q", err.Error(), tt.wantErrSub)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			val, err := stack.Pop()
			if err != nil {
				t.Fatalf("Pop() error: %v", err)
			}
			got, err := toFloat64(val, "check")
			if err != nil {
				t.Fatalf("toFloat64() error: %v", err)
			}
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("result = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBinaryExponentiationFloat(t *testing.T) {
	tests := []struct {
		base float64
		exp  int
		want float64
	}{
		{2, 10, 1024},
		{3, 4, 81},
		{5, 0, 1},
		{2, -3, 0.125},
		{10, 1, 10},
		{10, 2, 100},
		{2, 20, 1048576},
		{0, 5, 0},
		{0.5, 2, 0.25},
		{7, 3, 343},
		{1, 100, 1},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := binaryExponentiationFloat(tt.base, tt.exp)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("binaryExponentiationFloat(%v, %d) = %v, want %v", tt.base, tt.exp, got, tt.want)
			}
		})
	}
}
