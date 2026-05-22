// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

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
