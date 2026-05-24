// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import "strings"

// ListConstants lists all constants.
// Usage: `constants`
func (o *Operations) ListConstants() (string, error) {
	infos := o.consts.ListConstants()
	if len(infos) == 0 {
		return "No constants defined", nil
	}
	var sb strings.Builder
	for i, info := range infos {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(info.Name)
		sb.WriteString(" = ")
		// Use NumericValue interface for consistent formatting
		num := NewNumber(info.Value, FloatMode)
		sb.WriteString(num.String())
	}
	return sb.String(), nil
}

// ClearConstants removes all user-defined constants and resets built-in constants to their default values.
// Usage: `clearconstants`
func (o *Operations) ClearConstants() {
	o.consts.ReloadBuiltInConstants()
}
