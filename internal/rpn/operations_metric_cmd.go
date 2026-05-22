// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"sort"
	"strings"
)

// MetricShow shows metric info for top of stack.
func (o *Operations) MetricShow(stack *Stack) (string, error) {
	if stack.Len() < 1 {
		return "", fmt.Errorf("metric show: stack is empty")
	}
	val, err := stack.Peek()
	if err != nil {
		return "", buildError("metric show", err)
	}
	m := val.Metric()
	if m == nil || m.Category == Universal {
		return "Cool (Universal)", nil
	}
	factor := m.Factor(o.GetPrefixMode())
	return fmt.Sprintf("%s, %s, base: %s, factor: %.0g", m.Name, m.Category, m.BaseUnit, factor), nil
}

// MetricList returns all category names.
func (o *Operations) MetricList(stack *Stack) (string, error) {
	reg := o.metricRegistry
	all := reg.List()
	seen := make(map[Category]bool)
	var cats []string
	for _, m := range all {
		if !seen[m.Category] {
			seen[m.Category] = true
			cats = append(cats, m.Category.String())
		}
	}
	sort.Strings(cats)
	return strings.Join(cats, ", "), nil
}

// MetricCategory returns all metric names in the given category.
func (o *Operations) MetricCategory(stack *Stack, categoryName string) (string, error) {
	reg := o.metricRegistry
	cat, ok := parseCategory(categoryName)
	if !ok {
		return "", fmt.Errorf("metric: unknown category %q", categoryName)
	}
	metrics := reg.ListByCategory(cat)
	var names []string
	for _, m := range metrics {
		names = append(names, m.Name)
	}
	sort.Strings(names)
	return strings.Join(names, ", "), nil
}

// MetricCompatible checks if top two stack metrics are compatible.
func (o *Operations) MetricCompatible(stack *Stack) (string, error) {
	if stack.Len() < 2 {
		return "", fmt.Errorf("metric compatible: need at least 2 values on stack")
	}
	vals := stack.Values()
	top := vals[len(vals)-1]
	second := vals[len(vals)-2]
	mA := resolveMetric(o.metricRegistry, second)
	mB := resolveMetric(o.metricRegistry, top)
	compatible := categoriesCompatible(mA, mB)
	result := fmt.Sprintf("%s (%s) and %s (%s): %v",
		mA.Name, mA.Category, mB.Name, mB.Category, compatible)
	return result, nil
}

// parseCategory converts a category name string to a Category constant.
func parseCategory(name string) (Category, bool) {
	switch name {
	case "Universal":
		return Universal, true
	case "DataRate":
		return DataRate, true
	case "DataSize":
		return DataSize, true
	case "Time":
		return Time, true
	case "Weight":
		return Weight, true
	case "Speed":
		return Speed, true
	case "Distance":
		return Distance, true
	case "Custom":
		return Custom, true
	default:
		return 0, false
	}
}

// CustomList returns all custom metric names.
func (o *Operations) CustomList(stack *Stack) (string, error) {
	reg := o.metricRegistry
	metrics := reg.ListByCategory(Custom)
	var names []string
	for _, m := range metrics {
		names = append(names, m.Name)
	}
	sort.Strings(names)
	if len(names) == 0 {
		return "no custom metrics defined", nil
	}
	return strings.Join(names, ", "), nil
}

// CustomDefine registers a custom metric.
func (o *Operations) CustomDefine(name string, factor float64, category string) error {
	reg := o.metricRegistry
	// Check if already exists
	if _, ok := reg.Find(name); ok {
		return fmt.Errorf("metric %q already exists", name)
	}
	cat, ok := parseCategory(category)
	if !ok {
		return fmt.Errorf("unknown category %q", category)
	}
	m := &Metric{
		Name:     name,
		Category: cat,
		BaseUnit: cat.String() + "_base",
		Factor:   func(PrefixMode) float64 { return factor },
		IsCustom: true,
	}
	reg.Register(m)
	return nil
}

// CustomUndefine removes a custom metric.
func (o *Operations) CustomUndefine(name string) error {
	return o.metricRegistry.Unregister(name)
}
