// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"slices"
)

// Stack represents an RPN stack of values.
type Stack struct {
	values []StackValue
}

// NewStack creates a new empty stack.
func NewStack() *Stack {
	return &Stack{
		values: make([]StackValue, 0),
	}
}

// Push adds a value to the top of the stack.
func (s *Stack) Push(val StackValue) {
	s.values = append(s.values, val)
}

// Pop removes and returns the top value from the stack.
// Returns an error if the stack is empty.
func (s *Stack) Pop() (StackValue, error) {
	if len(s.values) == 0 {
		return nil, fmt.Errorf("stack is empty")
	}

	val := s.values[len(s.values)-1]
	s.values = s.values[:len(s.values)-1]
	return val, nil
}

// Peek returns the top value without removing it.
// Returns an error if the stack is empty.
func (s *Stack) Peek() (StackValue, error) {
	if len(s.values) == 0 {
		return nil, fmt.Errorf("stack is empty")
	}
	return s.values[len(s.values)-1], nil
}

// Len returns the number of values on the stack.
func (s *Stack) Len() int {
	return len(s.values)
}

// Values returns a copy of all stack values (top-to-bottom order).
func (s *Stack) Values() []StackValue {
	return slices.Clone(s.values)
}

// Clear removes all values from the stack.
// Note: This resets the slice to nil, releasing the underlying memory.
func (s *Stack) Clear() {
	s.values = nil
}
