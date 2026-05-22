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

	// Weight (base: kilograms)
	r.Register(&Metric{
		Name:     "mg",
		Category: Weight,
		BaseUnit: "kilograms",
		Factor:   func(PrefixMode) float64 { return 1e-6 },
		IsRate:   false,
	})
	r.Register(&Metric{
		Name:     "g",
		Category: Weight,
		BaseUnit: "kilograms",
		Factor:   func(PrefixMode) float64 { return 1e-3 },
		IsRate:   false,
	})
	r.Register(&Metric{
		Name:     "kg",
		Category: Weight,
		BaseUnit: "kilograms",
		Factor:   func(PrefixMode) float64 { return 1 },
		IsRate:   false,
	})
	r.Register(&Metric{
		Name:     "ton",
		Category: Weight,
		BaseUnit: "kilograms",
		Factor:   func(PrefixMode) float64 { return 1000 },
		IsRate:   false,
	})
	r.Register(&Metric{
		Name:     "lb",
		Category: Weight,
		BaseUnit: "kilograms",
		Factor:   func(PrefixMode) float64 { return 0.45359237 },
		IsRate:   false,
	})
	r.Register(&Metric{
		Name:     "oz",
		Category: Weight,
		BaseUnit: "kilograms",
		Factor:   func(PrefixMode) float64 { return 0.028349523125 },
		IsRate:   false,
	})

	// Speed (base: m/s)
	r.Register(&Metric{
		Name:     "mps",
		Category: Speed,
		BaseUnit: "m/s",
		Factor:   func(PrefixMode) float64 { return 1 },
		IsRate:   false,
	})
	r.Register(&Metric{
		Name:     "kmh",
		Category: Speed,
		BaseUnit: "m/s",
		Factor:   func(PrefixMode) float64 { return 1.0 / 3.6 },
		IsRate:   false,
	})
	r.Register(&Metric{
		Name:     "mph",
		Category: Speed,
		BaseUnit: "m/s",
		Factor:   func(PrefixMode) float64 { return 0.44704 },
		IsRate:   false,
	})
	r.Register(&Metric{
		Name:     "knots",
		Category: Speed,
		BaseUnit: "m/s",
		Factor:   func(PrefixMode) float64 { return 1852.0 / 3600 },
		IsRate:   false,
	})

	// Distance (base: meters)
	r.Register(&Metric{
		Name:     "m",
		Category: Distance,
		BaseUnit: "meters",
		Factor:   func(PrefixMode) float64 { return 1 },
		IsRate:   false,
	})
	r.Register(&Metric{
		Name:     "km",
		Category: Distance,
		BaseUnit: "meters",
		Factor:   func(PrefixMode) float64 { return 1000 },
		IsRate:   false,
	})
	r.Register(&Metric{
		Name:     "mi",
		Category: Distance,
		BaseUnit: "meters",
		Factor:   func(PrefixMode) float64 { return 1609.344 },
		IsRate:   false,
	})
	r.Register(&Metric{
		Name:     "ft",
		Category: Distance,
		BaseUnit: "meters",
		Factor:   func(PrefixMode) float64 { return 0.3048 },
		IsRate:   false,
	})
	r.Register(&Metric{
		Name:     "in",
		Category: Distance,
		BaseUnit: "meters",
		Factor:   func(PrefixMode) float64 { return 0.0254 },
		IsRate:   false,
	})
	r.Register(&Metric{
		Name:     "nm",
		Category: Distance,
		BaseUnit: "meters",
		Factor:   func(PrefixMode) float64 { return 1852 },
		IsRate:   false,
	})

	// Aliases (canonical name must exist first)
	// Note: Bps is intentionally NOT an alias for bps (capital B = bytes)
	
	// Data rate units are case-sensitive: b = bits, B = bytes
	r.MarkExactMatch("bps", "Kbps", "Mbps", "Gbps", "Tbps")
	
	r.RegisterAlias("bit/s", "bps")
	r.RegisterAlias("kbit/s", "Kbps")
	r.RegisterAlias("mbit/s", "Mbps")
	r.RegisterAlias("gbit/s", "Gbps")
	r.RegisterAlias("tbit/s", "Tbps")
	r.RegisterAlias("sec", "s")
	r.RegisterAlias("secs", "s")
	r.RegisterAlias("knot", "knots")
	r.RegisterAlias("mile", "mi")
	r.RegisterAlias("miles", "mi")
	r.RegisterAlias("foot", "ft")
	r.RegisterAlias("feet", "ft")
}
