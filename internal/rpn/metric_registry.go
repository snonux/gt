// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"strings"
	"sync"
)

// MetricRegistry provides thread-safe storage and lookup for metrics.
type MetricRegistry struct {
	mu              sync.RWMutex
	metrics         map[string]*Metric
	aliases         map[string]string // alias name -> canonical metric name
	exactMatchNames map[string]bool   // names that must match exactly (no case-insensitive fallback)
}

// Global registry instance.
var defaultRegistry *MetricRegistry
var registryOnce sync.Once

// GetMetricRegistry returns the global metric registry, initialized with built-in metrics.
func GetMetricRegistry() *MetricRegistry {
	registryOnce.Do(func() {
		defaultRegistry = NewMetricRegistry()
		registerBuiltInMetrics(defaultRegistry)
	})
	return defaultRegistry
}

// NewMetricRegistry creates a new empty registry (no built-in metrics).
func NewMetricRegistry() *MetricRegistry {
	return &MetricRegistry{
		metrics:         make(map[string]*Metric),
		aliases:         make(map[string]string),
		exactMatchNames: make(map[string]bool),
	}
}

// Register adds a metric to the registry.
// Panics if a metric with the same name already exists.
func (r *MetricRegistry) Register(m *Metric) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.metrics[m.Name]; exists {
		panic(fmt.Sprintf("metric %q already registered", m.Name))
	}
	r.metrics[m.Name] = m
}

// Find looks up a metric by name (case-sensitive).
func (r *MetricRegistry) Find(name string) (*Metric, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.metrics[name]
	return m, ok
}

// FindCaseInsensitive looks up a metric by name ignoring case.
// Returns the canonical name match.
func (r *MetricRegistry) FindCaseInsensitive(name string) (*Metric, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	lower := strings.ToLower(name)
	for _, m := range r.metrics {
		if strings.ToLower(m.Name) == lower {
			return m, true
		}
	}
	return nil, false
}

// RegisterAlias maps an alternative name to an existing metric.
// The canonical metric must already be registered.
// Panics if the alias case-insensitively matches a MarkExactMatch name.
func (r *MetricRegistry) RegisterAlias(alias, canonicalName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.metrics[canonicalName]; !ok {
		panic(fmt.Sprintf("cannot alias %q to %q: canonical metric not found", alias, canonicalName))
	}
	aliasLower := strings.ToLower(alias)
	for emName := range r.exactMatchNames {
		if strings.ToLower(emName) == aliasLower {
			panic(fmt.Sprintf("cannot alias %q: conflicts with exact-match name %q", alias, emName))
		}
	}
	r.aliases[alias] = canonicalName
}

// MarkExactMatch marks metric names that require exact case matching,
// disabling case-insensitive fallback. Used for units where case matters
// (e.g., bps vs Bps where lowercase b = bits, uppercase B = bytes).
func (r *MetricRegistry) MarkExactMatch(names ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, name := range names {
		r.exactMatchNames[name] = true
	}
}

// FindWithAliases looks up a metric by name, resolving aliases.
// Checks exact match first, then aliases, then case-insensitive (unless exact-match).
func (r *MetricRegistry) FindWithAliases(name string) (*Metric, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Exact match
	if m, ok := r.metrics[name]; ok {
		return m, true
	}

	// Alias
	if canonical, ok := r.aliases[name]; ok {
		return r.metrics[canonical], true
	}

	// Case-insensitive (skip if exact-match guard matches)
	lower := strings.ToLower(name)
	if len(r.exactMatchNames) > 0 {
		for emName := range r.exactMatchNames {
			if strings.ToLower(emName) == lower {
				return nil, false
			}
		}
	}

	for _, m := range r.metrics {
		if strings.ToLower(m.Name) == lower {
			return m, true
		}
	}
	return nil, false
}

// List returns all registered metrics.
func (r *MetricRegistry) List() []*Metric {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Metric, 0, len(r.metrics))
	for _, m := range r.metrics {
		result = append(result, m)
	}
	return result
}

// ListByCategory returns all metrics in the given category.
func (r *MetricRegistry) ListByCategory(cat Category) []*Metric {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Metric
	for _, m := range r.metrics {
		if m.Category == cat {
			result = append(result, m)
		}
	}
	return result
}
