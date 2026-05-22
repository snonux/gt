// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

// comparison operators

// GT pops two values from stack, compares (a > b), and pushes a boolean result.
func (o *Operations) GT(stack *Stack) error {
	a, b, err := popTwo(stack, "gt")
	if err != nil {
		return err
	}

	aVal, err := toFloat64(a, "gt comparison for a")
	if err != nil {
		return err
	}
	bVal, err := toFloat64(b, "gt comparison for b")
	if err != nil {
		return err
	}

	stack.Push(NewFloatFromBool(aVal > bVal))
	return nil
}

// LT pops two values from stack, compares (a < b), and pushes a boolean result.
func (o *Operations) LT(stack *Stack) error {
	a, b, err := popTwo(stack, "lt")
	if err != nil {
		return err
	}

	aVal, err := toFloat64(a, "lt comparison for a")
	if err != nil {
		return err
	}
	bVal, err := toFloat64(b, "lt comparison for b")
	if err != nil {
		return err
	}

	stack.Push(NewFloatFromBool(aVal < bVal))
	return nil
}

// GTE pops two values from stack, compares (a >= b), and pushes a boolean result.
func (o *Operations) GTE(stack *Stack) error {
	a, b, err := popTwo(stack, "gte")
	if err != nil {
		return err
	}

	aVal, err := toFloat64(a, "gte comparison for a")
	if err != nil {
		return err
	}
	bVal, err := toFloat64(b, "gte comparison for b")
	if err != nil {
		return err
	}

	stack.Push(NewFloatFromBool(aVal >= bVal))
	return nil
}

// LTE pops two values from stack, compares (a <= b), and pushes a boolean result.
func (o *Operations) LTE(stack *Stack) error {
	a, b, err := popTwo(stack, "lte")
	if err != nil {
		return err
	}

	aVal, err := toFloat64(a, "lte comparison for a")
	if err != nil {
		return err
	}
	bVal, err := toFloat64(b, "lte comparison for b")
	if err != nil {
		return err
	}

	stack.Push(NewFloatFromBool(aVal <= bVal))
	return nil
}

// EQ pops two values from stack, compares (a == b), and pushes a boolean result.
func (o *Operations) EQ(stack *Stack) error {
	a, b, err := popTwo(stack, "eq")
	if err != nil {
		return err
	}

	aVal, err := toFloat64(a, "eq comparison for a")
	if err != nil {
		return err
	}
	bVal, err := toFloat64(b, "eq comparison for b")
	if err != nil {
		return err
	}

	stack.Push(NewFloatFromBool(aVal == bVal))
	return nil
}

// NEQ pops two values from stack, compares (a != b), and pushes a boolean result.
func (o *Operations) NEQ(stack *Stack) error {
	a, b, err := popTwo(stack, "neq")
	if err != nil {
		return err
	}

	aVal, err := toFloat64(a, "neq comparison for a")
	if err != nil {
		return err
	}
	bVal, err := toFloat64(b, "neq comparison for b")
	if err != nil {
		return err
	}

	stack.Push(NewFloatFromBool(aVal != bVal))
	return nil
}
