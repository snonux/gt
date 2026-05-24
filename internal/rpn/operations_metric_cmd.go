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
// Iterates over all valid Category values using a for loop bounded by _sentinel,
// so adding a new Category constant (between Universal and _sentinel) automatically
// makes it available here without modifying this function.
func parseCategory(name string) (Category, bool) {
	for cat := Category(0); cat < _sentinel; cat++ {
		if cat.String() == name {
			return cat, true
		}
	}
	return 0, false
}

// CustomShow returns detailed info for custom metrics.
// If name is empty, shows all custom metrics; otherwise shows the named one.
func (o *Operations) CustomShow(stack *Stack, name string) (string, error) {
	reg := o.metricRegistry
	if name != "" {
		m, ok := reg.Find(name)
		if !ok {
			return "", fmt.Errorf("unknown custom metric %q", name)
		}
		if !m.IsCustom {
			return "", fmt.Errorf("metric %q is not a custom metric", name)
		}
		factor := m.Factor(o.GetPrefixMode())
		return fmt.Sprintf("%s, category: %s, base: %s, factor: %.10g", m.Name, m.Category, m.BaseUnit, factor), nil
	}
	metrics := reg.ListByCategory(Custom)
	if len(metrics) == 0 {
		return "no custom metrics defined", nil
	}
	var lines []string
	for _, m := range metrics {
		factor := m.Factor(o.GetPrefixMode())
		lines = append(lines, fmt.Sprintf("  %s, category: %s, base: %s, factor: %.10g", m.Name, m.Category, m.BaseUnit, factor))
	}
	return strings.Join(lines, "\n"), nil
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
