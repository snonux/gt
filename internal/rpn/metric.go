// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"strings"
	"sync"
)

// Category identifies which domain a metric belongs to.
type Category int

const (
	// Universal is the default category — "Cool" unit.
	Universal Category = iota
	// DataRate covers bits/sec, bytes/sec, etc.
	DataRate
	// DataSize covers bits, bytes, KB, MB, GB, etc.
	DataSize
	// Time covers ms, s, min, hr, day.
	Time
	// Weight covers mg, g, kg, ton, lb, oz.
	Weight
	// Speed covers m/s, km/h, mph, knots.
	Speed
	// Distance covers m, km, mi, ft, in.
	Distance
	// Custom is for user-defined units.
	Custom
)

// String returns the human-readable name of the category.
func (c Category) String() string {
	switch c {
	case Universal:
		return "Universal"
	case DataRate:
		return "DataRate"
	case DataSize:
		return "DataSize"
	case Time:
		return "Time"
	case Weight:
		return "Weight"
	case Speed:
		return "Speed"
	case Distance:
		return "Distance"
	case Custom:
		return "Custom"
	default:
		return fmt.Sprintf("Category(%d)", c)
	}
}

// PrefixMode determines whether data size prefixes are SI (1000-based) or IEC (1024-based).
type PrefixMode int

const (
	// SI uses powers of 1000 for prefixes: K=1000, M=1000000, etc.
	SI PrefixMode = iota
	// IEC uses powers of 1024 for prefixes: Ki=1024, Mi=1048576, etc.
	IEC
)

// String returns the human-readable name of the prefix mode.
func (p PrefixMode) String() string {
	switch p {
	case SI:
		return "SI"
	case IEC:
		return "IEC"
	default:
		return fmt.Sprintf("PrefixMode(%d)", p)
	}
}

// Metric defines a unit of measurement with conversion to its base unit.
type Metric struct {
	// Name is the canonical name, e.g. "Mbps", "GB", "hr".
	Name string
	// Category identifies the domain this metric belongs to.
	Category Category
	// BaseUnit is the name of the base unit in this category.
	BaseUnit string
	// Factor returns the conversion factor to the base unit for the given prefix mode.
	Factor func(mode PrefixMode) float64
	// IsRate is true for rate units (per-second).
	IsRate bool
	// IsCustom is true for user-defined custom units.
	IsCustom bool
}

// String returns the canonical name of the metric.
func (m *Metric) String() string {
	return m.Name
}

// MetricRegistry provides thread-safe storage and lookup for metrics.
type MetricRegistry struct {
	mu      sync.RWMutex
	metrics map[string]*Metric
}

// Global registry instance.
var defaultRegistry *MetricRegistry
var registryOnce sync.Once

// GetMetricRegistry returns the global metric registry, initialized with built-in metrics.
func GetMetricRegistry() *MetricRegistry {
	registryOnce.Do(func() {
		defaultRegistry = &MetricRegistry{
			metrics: make(map[string]*Metric),
		}
		registerBuiltInMetrics(defaultRegistry)
	})
	return defaultRegistry
}

// NewMetricRegistry creates a new empty registry (no built-in metrics).
func NewMetricRegistry() *MetricRegistry {
	return &MetricRegistry{
		metrics: make(map[string]*Metric),
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

// registerBuiltInMetrics populates the registry with built-in metrics.
func registerBuiltInMetrics(r *MetricRegistry) {
	// Universal
	r.Register(&Metric{
		Name:     "Cool",
		Category: Universal,
		BaseUnit: "cool",
		Factor:   func(PrefixMode) float64 { return 1 },
		IsRate:   false,
	})

	// DataRate (base: bps)
	r.Register(&Metric{
		Name:     "bps",
		Category: DataRate,
		BaseUnit: "bps",
		Factor:   func(PrefixMode) float64 { return 1 },
		IsRate:   true,
	})
	r.Register(&Metric{
		Name:     "Kbps",
		Category: DataRate,
		BaseUnit: "bps",
		Factor:   func(PrefixMode) float64 { return 1e3 },
		IsRate:   true,
	})
	r.Register(&Metric{
		Name:     "Mbps",
		Category: DataRate,
		BaseUnit: "bps",
		Factor:   func(PrefixMode) float64 { return 1e6 },
		IsRate:   true,
	})
	r.Register(&Metric{
		Name:     "Gbps",
		Category: DataRate,
		BaseUnit: "bps",
		Factor:   func(PrefixMode) float64 { return 1e9 },
		IsRate:   true,
	})
	r.Register(&Metric{
		Name:     "Tbps",
		Category: DataRate,
		BaseUnit: "bps",
		Factor:   func(PrefixMode) float64 { return 1e12 },
		IsRate:   true,
	})

	// DataSize (base: bits)
	r.Register(&Metric{
		Name:     "bits",
		Category: DataSize,
		BaseUnit: "bits",
		Factor:   func(PrefixMode) float64 { return 1 },
		IsRate:   false,
	})
	r.Register(&Metric{
		Name:     "bytes",
		Category: DataSize,
		BaseUnit: "bits",
		Factor:   func(PrefixMode) float64 { return 8 },
		IsRate:   false,
	})

	// SI-prefixed data size (KB, MB, GB, TB, PB)
	siPrefixes := []struct {
		name     string
		multiple float64
	}{
		{"KB", 1e3},
		{"MB", 1e6},
		{"GB", 1e9},
		{"TB", 1e12},
		{"PB", 1e15},
	}
	for _, p := range siPrefixes {
		p := p
		r.Register(&Metric{
			Name:     p.name,
			Category: DataSize,
			BaseUnit: "bits",
			Factor:   func(PrefixMode) float64 { return 8 * p.multiple },
			IsRate:   false,
		})
	}

	// IEC-prefixed data size (KiB, MiB, GiB, TiB, PiB)
	iecPrefixes := []struct {
		name    string
		power10 int
	}{
		{"KiB", 10},
		{"MiB", 20},
		{"GiB", 30},
		{"TiB", 40},
		{"PiB", 50},
	}
	for _, p := range iecPrefixes {
		p := p
		r.Register(&Metric{
			Name:     p.name,
			Category: DataSize,
			BaseUnit: "bits",
			Factor:   func(PrefixMode) float64 { return 8 * float64(uint64(1)<<uint64(p.power10)) },
			IsRate:   false,
		})
	}

	// Time (base: seconds)
	r.Register(&Metric{
		Name:     "ms",
		Category: Time,
		BaseUnit: "seconds",
		Factor:   func(PrefixMode) float64 { return 0.001 },
		IsRate:   false,
	})
	r.Register(&Metric{
		Name:     "s",
		Category: Time,
		BaseUnit: "seconds",
		Factor:   func(PrefixMode) float64 { return 1 },
		IsRate:   false,
	})
	r.Register(&Metric{
		Name:     "min",
		Category: Time,
		BaseUnit: "seconds",
		Factor:   func(PrefixMode) float64 { return 60 },
		IsRate:   false,
	})
	r.Register(&Metric{
		Name:     "hr",
		Category: Time,
		BaseUnit: "seconds",
		Factor:   func(PrefixMode) float64 { return 3600 },
		IsRate:   false,
	})
	r.Register(&Metric{
		Name:     "day",
		Category: Time,
		BaseUnit: "seconds",
		Factor:   func(PrefixMode) float64 { return 86400 },
		IsRate:   false,
	})
}
