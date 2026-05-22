// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"strconv"
	"testing"
)

func TestCategoryString(t *testing.T) {
	tests := []struct {
		cat    Category
		expect string
	}{
		{Universal, "Universal"},
		{DataRate, "DataRate"},
		{DataSize, "DataSize"},
		{Time, "Time"},
		{Weight, "Weight"},
		{Speed, "Speed"},
		{Distance, "Distance"},
		{Custom, "Custom"},
	}
	for _, tt := range tests {
		if got := tt.cat.String(); got != tt.expect {
			t.Errorf("Category(%d).String() = %q, want %q", tt.cat, got, tt.expect)
		}
	}
}

func TestPrefixModeString(t *testing.T) {
	if got := SI.String(); got != "SI" {
		t.Errorf("SI.String() = %q, want %q", got, "SI")
	}
	if got := IEC.String(); got != "IEC" {
		t.Errorf("IEC.String() = %q, want %q", got, "IEC")
	}
}

func TestMetricString(t *testing.T) {
	m := &Metric{Name: "Mbps"}
	if got := m.String(); got != "Mbps" {
		t.Errorf("Metric.String() = %q, want %q", got, "Mbps")
	}
}

func TestMetricRegistryRegisterAndFind(t *testing.T) {
	reg := NewMetricRegistry()
	m := &Metric{
		Name:     "test",
		Category: Custom,
		BaseUnit: "test",
		Factor:   func(PrefixMode) float64 { return 42 },
	}
	reg.Register(m)

	found, ok := reg.Find("test")
	if !ok {
		t.Fatal("Find did not return metric")
	}
	if found != m {
		t.Error("Find returned wrong metric")
	}

	_, ok = reg.Find("nope")
	if ok {
		t.Error("Find should not find unknown metric")
	}
}

func TestMetricRegistryFindCaseInsensitive(t *testing.T) {
	reg := NewMetricRegistry()
	m := &Metric{
		Name:     "MyUnit",
		Category: Custom,
		BaseUnit: "myunit",
		Factor:   func(PrefixMode) float64 { return 1 },
	}
	reg.Register(m)

	found, ok := reg.FindCaseInsensitive("myunit")
	if !ok {
		t.Fatal("FindCaseInsensitive did not find metric")
	}
	if found.Name != "MyUnit" {
		t.Errorf("FindCaseInsensitive returned name %q, want %q", found.Name, "MyUnit")
	}

	found, ok = reg.FindCaseInsensitive("MYUNIT")
	if !ok || found.Name != "MyUnit" {
		t.Error("FindCaseInsensitive should match any case")
	}
}

func TestMetricRegistryDuplicatePanics(t *testing.T) {
	reg := NewMetricRegistry()
	reg.Register(&Metric{
		Name:   "dup",
		Category: Custom,
		Factor: func(PrefixMode) float64 { return 1 },
	})
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on duplicate register")
		}
	}()
	reg.Register(&Metric{
		Name:   "dup",
		Category: Custom,
		Factor: func(PrefixMode) float64 { return 2 },
	})
}

func TestMetricRegistryList(t *testing.T) {
	reg := NewMetricRegistry()
	reg.Register(&Metric{Name: "a", Category: Custom, Factor: func(PrefixMode) float64 { return 1 }})
	reg.Register(&Metric{Name: "b", Category: Custom, Factor: func(PrefixMode) float64 { return 2 }})
	reg.Register(&Metric{Name: "c", Category: Time, Factor: func(PrefixMode) float64 { return 3 }})

	all := reg.List()
	if len(all) != 3 {
		t.Errorf("List() returned %d metrics, want 3", len(all))
	}

	timeMetrics := reg.ListByCategory(Time)
	if len(timeMetrics) != 1 || timeMetrics[0].Name != "c" {
		t.Errorf("ListByCategory(Time) = %v, want [c]", timeMetrics)
	}

	customMetrics := reg.ListByCategory(Custom)
	if len(customMetrics) != 2 {
		t.Errorf("ListByCategory(Custom) = %d metrics, want 2", len(customMetrics))
	}
}

func TestDefaultRegistryHasBuiltIns(t *testing.T) {
	reg := GetMetricRegistry()

	// Check some expected built-in metrics exist
	expectedNames := []string{
		"Cool", "bps", "Mbps", "Gbps", "bits", "bytes", "KB", "GB", "KiB", "MiB",
		"ms", "s", "min", "hr", "day",
	}
	for _, name := range expectedNames {
		m, ok := reg.Find(name)
		if !ok {
			t.Errorf("expected built-in metric %q not found", name)
			continue
		}
		if m.Name != name {
			t.Errorf("metric name = %q, want %q", m.Name, name)
		}
	}
}

func TestBuiltInMetricFactors(t *testing.T) {
	reg := GetMetricRegistry()
	tolerance := 0.0001

	tests := []struct {
		name    string
		mode    PrefixMode
		expect  float64
	}{
		// Universal
		{"Cool", SI, 1},
		{"Cool", IEC, 1},
		// DataRate
		{"bps", SI, 1},
		{"Kbps", SI, 1e3},
		{"Mbps", SI, 1e6},
		{"Gbps", SI, 1e9},
		{"Tbps", SI, 1e12},
		// DataSize (base: bits)
		{"bits", SI, 1},
		{"bytes", SI, 8},
		{"KB", SI, 8000},    // 8 * 1000
		{"MB", SI, 8e6},
		{"GB", SI, 8e9},
		{"TB", SI, 8e12},
		{"PB", SI, 8e15},
		// IEC (base: bits)
		{"KiB", SI, 8192},       // 8 * 1024
		{"MiB", SI, 8 * 1048576},
		{"GiB", SI, 8 * 1073741824},
		{"TiB", SI, 8 * float64(uint64(1)<<40)},
		{"PiB", SI, 8 * float64(uint64(1)<<50)},
		// Time (base: seconds)
		{"ms", SI, 0.001},
		{"s", SI, 1},
		{"min", SI, 60},
		{"hr", SI, 3600},
		{"day", SI, 86400},
	}

	for _, tt := range tests {
		m, ok := reg.Find(tt.name)
		if !ok {
			t.Fatalf("metric %q not found", tt.name)
		}
		got := m.Factor(tt.mode)
		if got < tt.expect-tolerance || got > tt.expect+tolerance {
			t.Errorf("%s.Factor(%s) = %g, want %g (tolerance %g)", tt.name, tt.mode, got, tt.expect, tolerance)
		}
	}
}

func TestBuiltInMetricCategories(t *testing.T) {
	reg := GetMetricRegistry()

	tests := []struct {
		name     string
		category Category
	}{
		{"Cool", Universal},
		{"bps", DataRate},
		{"Mbps", DataRate},
		{"Gbps", DataRate},
		{"bits", DataSize},
		{"bytes", DataSize},
		{"KB", DataSize},
		{"GB", DataSize},
		{"KiB", DataSize},
		{"MiB", DataSize},
		{"ms", Time},
		{"s", Time},
		{"min", Time},
		{"hr", Time},
		{"day", Time},
	}

	for _, tt := range tests {
		m, ok := reg.Find(tt.name)
		if !ok {
			t.Fatalf("metric %q not found", tt.name)
		}
		if m.Category != tt.category {
			t.Errorf("%s.Category = %s, want %s", tt.name, m.Category, tt.category)
		}
	}
}

func TestBuiltInMetricBaseUnits(t *testing.T) {
	reg := GetMetricRegistry()

	tests := []struct {
		name     string
		baseUnit string
	}{
		{"Cool", "cool"},
		{"bps", "bps"},
		{"Mbps", "bps"},
		{"bits", "bits"},
		{"bytes", "bits"},
		{"KB", "bits"},
		{"GB", "bits"},
		{"KiB", "bits"},
		{"ms", "seconds"},
		{"s", "seconds"},
		{"min", "seconds"},
		{"hr", "seconds"},
		{"day", "seconds"},
	}

	for _, tt := range tests {
		m, ok := reg.Find(tt.name)
		if !ok {
			t.Fatalf("metric %q not found", tt.name)
		}
		if m.BaseUnit != tt.baseUnit {
			t.Errorf("%s.BaseUnit = %q, want %q", tt.name, m.BaseUnit, tt.baseUnit)
		}
	}
}

func TestBuiltInMetricIsRate(t *testing.T) {
	reg := GetMetricRegistry()

	rateMetrics := []string{"bps", "Kbps", "Mbps", "Gbps", "Tbps"}
	for _, name := range rateMetrics {
		m, ok := reg.Find(name)
		if !ok {
			t.Fatalf("metric %q not found", name)
		}
		if !m.IsRate {
			t.Errorf("%s.IsRate = false, want true", name)
		}
	}

	nonRateMetrics := []string{"Cool", "bits", "bytes", "KB", "GB", "ms", "s", "min", "hr", "day"}
	for _, name := range nonRateMetrics {
		m, ok := reg.Find(name)
		if !ok {
			t.Fatalf("metric %q not found", name)
		}
		if m.IsRate {
			t.Errorf("%s.IsRate = true, want false", name)
		}
	}
}

func TestGetMetricRegistrySingleton(t *testing.T) {
	r1 := GetMetricRegistry()
	r2 := GetMetricRegistry()
	if r1 != r2 {
		t.Error("GetMetricRegistry should return the same instance")
	}
}

func TestNewMetricRegistryIsEmpty(t *testing.T) {
	reg := NewMetricRegistry()
	if got := len(reg.List()); got != 0 {
		t.Errorf("NewMetricRegistry returned %d metrics, want 0", got)
	}
	_, ok := reg.Find("Cool")
	if ok {
		t.Error("NewMetricRegistry should not have built-in metrics")
	}
}

func TestMetricRegistryListImmutability(t *testing.T) {
	reg := NewMetricRegistry()
	reg.Register(&Metric{Name: "a", Category: Custom, Factor: func(PrefixMode) float64 { return 1 }})

	list := reg.List()
	originalLen := len(list)
	list = append(list, &Metric{Name: "fake"}) // mutate the returned slice
	if got := len(reg.List()); got != originalLen {
		t.Errorf("modifying returned slice affected registry: len = %d, want %d", got, originalLen)
	}
}

func TestMetricRegistryConcurrent(t *testing.T) {
	reg := NewMetricRegistry()
	const n = 50
	done := make(chan bool, n*2)

	for i := 0; i < n; i++ {
		i := i
		go func() {
			reg.Register(&Metric{
				Name:     fmt.Sprintf("m%d", i),
				Category: Custom,
				Factor:   func(PrefixMode) float64 { return 1 },
			})
			done <- true
		}()
		go func() {
			reg.Find(fmt.Sprintf("m%d", i))
			done <- true
		}()
	}

	for i := 0; i < n*2; i++ {
		<-done
	}

	if got := len(reg.List()); got != n {
		t.Errorf("concurrent register: got %d metrics, want %d", got, n)
	}
}

func TestBuiltInWeightMetrics(t *testing.T) {
	reg := GetMetricRegistry()
	tolerance := 0.00001

	tests := []struct {
		name   string
		expect float64
	}{
		{"mg", 1e-6},
		{"g", 1e-3},
		{"kg", 1},
		{"ton", 1000},
		{"lb", 0.45359237},
		{"oz", 0.028349523125},
	}

	for _, tt := range tests {
		m, ok := reg.Find(tt.name)
		if !ok {
			t.Fatalf("metric %q not found", tt.name)
		}
		if m.Category != Weight {
			t.Errorf("%s.Category = %s, want Weight", tt.name, m.Category)
		}
		got := m.Factor(SI)
		if got < tt.expect-tolerance || got > tt.expect+tolerance {
			t.Errorf("%s.Factor(SI) = %g, want %g", tt.name, got, tt.expect)
		}
	}
}

func TestBuiltInSpeedMetrics(t *testing.T) {
	reg := GetMetricRegistry()
	tolerance := 0.001

	tests := []struct {
		name   string
		expect float64
	}{
		{"mps", 1},
		{"kmh", 1.0 / 3.6},
		{"mph", 0.44704},
		{"knots", 1852.0 / 3600},
	}

	for _, tt := range tests {
		m, ok := reg.Find(tt.name)
		if !ok {
			t.Fatalf("metric %q not found", tt.name)
		}
		if m.Category != Speed {
			t.Errorf("%s.Category = %s, want Speed", tt.name, m.Category)
		}
		got := m.Factor(SI)
		if got < tt.expect-tolerance || got > tt.expect+tolerance {
			t.Errorf("%s.Factor(SI) = %g, want %g", tt.name, got, tt.expect)
		}
	}
}

func TestBuiltInDistanceMetrics(t *testing.T) {
	reg := GetMetricRegistry()

	tests := []struct {
		name   string
		expect float64
	}{
		{"m", 1},
		{"km", 1000},
		{"mi", 1609.344},
		{"ft", 0.3048},
		{"in", 0.0254},
		{"nm", 1852},
	}

	for _, tt := range tests {
		m, ok := reg.Find(tt.name)
		if !ok {
			t.Fatalf("metric %q not found", tt.name)
		}
		if m.Category != Distance {
			t.Errorf("%s.Category = %s, want Distance", tt.name, m.Category)
		}
		got := m.Factor(SI)
		if got != tt.expect {
			t.Errorf("%s.Factor(SI) = %g, want %g", tt.name, got, tt.expect)
		}
	}
}

func TestMetricRegistryAliases(t *testing.T) {
	reg := GetMetricRegistry()

	tests := []struct {
		alias    string
		canonical string
	}{
		{"bit/s", "bps"},
		{"kbit/s", "Kbps"},
		{"mbit/s", "Mbps"},
		{"gbit/s", "Gbps"},
		{"tbit/s", "Tbps"},
		{"sec", "s"},
		{"secs", "s"},
		{"knot", "knots"},
		{"mile", "mi"},
		{"miles", "mi"},
		{"foot", "ft"},
		{"feet", "ft"},
	}

	for _, tt := range tests {
		m, ok := reg.FindWithAliases(tt.alias)
		if !ok {
			t.Fatalf("alias %q not found", tt.alias)
		}
		if m.Name != tt.canonical {
			t.Errorf("alias %q resolved to %q, want %q", tt.alias, m.Name, tt.canonical)
		}
	}

	// Bps should NOT resolve to bps (capital B = bytes)
	if m, ok := reg.FindWithAliases("Bps"); ok {
		t.Errorf("Bps should not resolve, but got %q", m.Name)
	}

	// Other case variations of exact-match names should also not resolve
	for _, bad := range []string{"BPS", "KBPS", "MBPS", "GBPS", "TBPS"} {
		if m, ok := reg.FindWithAliases(bad); ok {
			t.Errorf("%s should not resolve (exact-match guard), got %q", bad, m.Name)
		}
	}
}

func TestFindWithAliasesPriority(t *testing.T) {
	reg := NewMetricRegistry()
	reg.Register(&Metric{
		Name:   "Foo",
		Category: Custom,
		Factor: func(PrefixMode) float64 { return 1 },
	})
	reg.RegisterAlias("foo", "Foo")

	// Exact match on canonical name
	m, ok := reg.FindWithAliases("Foo")
	if !ok || m.Name != "Foo" {
		t.Errorf("exact match failed")
	}

	// Alias match
	m, ok = reg.FindWithAliases("foo")
	if !ok || m.Name != "Foo" {
		t.Errorf("alias match failed")
	}
}

func TestRegisterAliasPanicOnMissingCanonical(t *testing.T) {
	reg := NewMetricRegistry()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on alias to nonexistent metric")
		}
	}()
	reg.RegisterAlias("nope", "doesNotExist")
}

func TestRegisterAliasRejectsExactMatchConflict(t *testing.T) {
	reg := NewMetricRegistry()
	reg.Register(&Metric{
		Name:     "bps",
		Category: Custom,
		BaseUnit: "bps",
		Factor:   func(PrefixMode) float64 { return 1 },
	})
	reg.MarkExactMatch("bps")

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on alias conflicting with exact-match name")
		}
	}()
	reg.RegisterAlias("Bps", "bps") // should panic
}

func TestAliasWinsOverCaseInsensitive(t *testing.T) {
	reg := NewMetricRegistry()
	reg.Register(&Metric{Name: "Foo", Category: Custom, Factor: func(PrefixMode) float64 { return 1 }})
	reg.Register(&Metric{Name: "Bar", Category: Custom, Factor: func(PrefixMode) float64 { return 2 }})
	reg.RegisterAlias("bar", "Foo") // alias "bar" -> "Foo"

	m, ok := reg.FindWithAliases("bar")
	if !ok || m.Name != "Foo" {
		t.Errorf("alias 'bar' should resolve to 'Foo', got %v (ok=%v)", m, ok)
	}
}

func TestNumberDefaultMetricIsCool(t *testing.T) {
	cool, _ := GetMetricRegistry().Find("Cool")

	// Float defaults to Cool
	f := NewFloat(42)
	if f.Metric() != cool {
		t.Errorf("NewFloat().Metric() = %v, want Cool", f.Metric())
	}

	// Rat defaults to Cool
	r := NewRat(42)
	if r.Metric() != cool {
		t.Errorf("NewRat().Metric() = %v, want Cool", r.Metric())
	}

	// NewNumber defaults to Cool
	n := NewNumber(42, FloatMode)
	if n.Metric() != cool {
		t.Errorf("NewNumber().Metric() = %v, want Cool", n.Metric())
	}
}

func TestNewFloatWithMetric(t *testing.T) {
	reg := GetMetricRegistry()
	mbps, _ := reg.Find("Mbps")

	f := NewFloatWithMetric(100, mbps)
	if f.Metric() != mbps {
		t.Errorf("NewFloatWithMetric.Metric() = %v, want Mbps", f.Metric())
	}
}

func TestNewRatWithMetric(t *testing.T) {
	reg := GetMetricRegistry()
	gb, _ := reg.Find("GB")

	r := NewRatWithMetric(5, gb)
	if r.Metric() != gb {
		t.Errorf("NewRatWithMetric.Metric() = %v, want GB", r.Metric())
	}
}

func TestNewNumberWithMetric(t *testing.T) {
	reg := GetMetricRegistry()
	hr, _ := reg.Find("hr")

	n := NewNumber(2, FloatMode, hr)
	if n.Metric() != hr {
		t.Errorf("NewNumber(metric).Metric() = %v, want hr", n.Metric())
	}
}

func TestFloatSetMetric(t *testing.T) {
	reg := GetMetricRegistry()
	mbps, _ := reg.Find("Mbps")
	gbps, _ := reg.Find("Gbps")

	f := NewFloatWithMetric(100, mbps)
	g := f.SetMetric(gbps)

	// Original unchanged
	if f.Metric() != mbps {
		t.Error("SetMetric should not modify original")
	}

	// New number has the new metric
	if g.Metric() != gbps {
		t.Errorf("SetMetric result = %v, want Gbps", g.Metric())
	}

	// Value preserved
	val, _ := g.Float64()
	if val != 100 {
		t.Errorf("SetMetric result value = %g, want 100", val)
	}
}

func TestRatSetMetric(t *testing.T) {
	reg := GetMetricRegistry()
	kg, _ := reg.Find("kg")
	lb, _ := reg.Find("lb")

	r := NewRatWithMetric(5, kg)
	s := r.SetMetric(lb)

	if r.Metric() != kg {
		t.Error("SetMetric should not modify original")
	}
	if s.Metric() != lb {
		t.Errorf("SetMetric result = %v, want lb", s.Metric())
	}

	val, _ := s.Float64()
	if val != 5 {
		t.Errorf("SetMetric result value = %g, want 5", val)
	}
}

func TestStringNumMetric(t *testing.T) {
	cool, _ := GetMetricRegistry().Find("Cool")
	s := NewStringNum("hello")
	if s.Metric() != cool {
		t.Errorf("StringNum.Metric() = %v, want Cool", s.Metric())
	}
}

func TestSymbolMetric(t *testing.T) {
	cool, _ := GetMetricRegistry().Find("Cool")
	s := NewSymbol("x")
	if s.Metric() != cool {
		t.Errorf("Symbol.Metric() = %v, want Cool", s.Metric())
	}
}

func TestGetCoolMetric(t *testing.T) {
	cool := GetCoolMetric()
	if cool == nil {
		t.Fatal("GetCoolMetric returned nil")
	}
	if cool.Name != "Cool" {
		t.Errorf("GetCoolMetric().Name = %q, want %q", cool.Name, "Cool")
	}
	if cool.Category != Universal {
		t.Errorf("GetCoolMetric().Category = %s, want Universal", cool.Category)
	}
}

func TestParseNumberWithMetric(t *testing.T) {
	reg := GetMetricRegistry()

	tests := []struct {
		token      string
		wantNum    float64
		wantMetric string
		wantOK     bool
	}{
		{"100Mbps", 100, "Mbps", true},
		{"5.5GB", 5.5, "GB", true},
		{"2hr", 2, "hr", true},
		{"3km", 3, "km", true},
		{"10lb", 10, "lb", true},
		{"-2.5kB", -2.5, "KB", true},
		{"1e6bps", 1e6, "bps", true},
		{"1E9Gbps", 1e9, "Gbps", true},
		// Plain number — no metric
		{"42", 0, "", false},
		// Plain number with decimal
		{"3.14", 0, "", false},
		// Unknown suffix
		{"100xyz", 0, "", false},
		// Empty
		{"", 0, "", false},
		// Just a letter
		{"Mbps", 0, "", false},
		// Negative
		{"-5MB", -5, "MB", true},
		// Positive sign
		{"+3hr", 3, "hr", true},
	}

	for _, tt := range tests {
		num, metric, ok := parseNumberWithMetric(tt.token)
		if ok != tt.wantOK {
			t.Errorf("parseNumberWithMetric(%q) ok = %v, want %v", tt.token, ok, tt.wantOK)
			continue
		}
		if ok {
			if num != tt.wantNum {
				t.Errorf("parseNumberWithMetric(%q) num = %g, want %g", tt.token, num, tt.wantNum)
			}
			if metric == nil {
				t.Errorf("parseNumberWithMetric(%q) metric = nil", tt.token)
				continue
			}
			if metric.Name != tt.wantMetric {
				// Check if it's an alias that resolved to the expected name
				expected, _ := reg.Find(tt.wantMetric)
				if expected == nil || metric.Name != expected.Name {
					t.Errorf("parseNumberWithMetric(%q) metric = %q, want %q", tt.token, metric.Name, tt.wantMetric)
				}
			}
		}
	}
}

func TestParseNumberWithMetricAliases(t *testing.T) {
	reg := GetMetricRegistry()

	tests := []struct {
		token      string
		wantMetric string
	}{
		{"1bit/s", "bps"},
		{"1kbit/s", "Kbps"},
		{"1mbit/s", "Mbps"},
		{"2sec", "s"},
		{"3secs", "s"},
		{"5knot", "knots"},
		{"1mile", "mi"},
		{"2miles", "mi"},
		{"1foot", "ft"},
		{"2feet", "ft"},
	}

	for _, tt := range tests {
		_, metric, ok := parseNumberWithMetric(tt.token)
		if !ok {
			t.Fatalf("parseNumberWithMetric(%q) = false, want true", tt.token)
		}
		expected, _ := reg.Find(tt.wantMetric)
		if metric != expected {
			t.Errorf("parseNumberWithMetric(%q) = %q, want %q", tt.token, metric.Name, tt.wantMetric)
		}
	}
}

func TestParseNumberWithMetricExactMatch(t *testing.T) {
	// Bps (capital B) should NOT resolve to bps
	_, _, ok := parseNumberWithMetric("100Bps")
	if ok {
		t.Error("parseNumberWithMetric(100Bps) should fail (B = bytes)")
	}
}

func TestAtPrefixMetricParsing(t *testing.T) {
	reg := GetMetricRegistry()

	tests := []struct {
		token      string
		wantMetric string
		wantOK     bool
	}{
		{"@GB", "GB", true},
		{"@Mbps", "Mbps", true},
		{"@hr", "hr", true},
		{"@sec", "s", true},  // alias
		{"@foot", "ft", true}, // alias
		{"@nope", "", false},
		{"@", "", false},
	}

	for _, tt := range tests {
		// Simulate the @ prefix logic from rpn_parse.go
		if len(tt.token) <= 1 || tt.token[0] != '@' {
			continue
		}
		metricName := tt.token[1:]
		metric, ok := reg.FindWithAliases(metricName)
		if ok != tt.wantOK {
			t.Errorf("@%s ok = %v, want %v", metricName, ok, tt.wantOK)
			continue
		}
		if ok && tt.wantMetric != "" {
			expected, _ := reg.Find(tt.wantMetric)
			if metric != expected {
				t.Errorf("@%s = %q, want %q", metricName, metric.Name, tt.wantMetric)
			}
		}
	}
}

func TestAtPrefixIntegration(t *testing.T) {
	// Test that @GB parses correctly through the full RPN pipeline
	vars := NewVariables()
	rpn := NewRPN(vars)

	// Parse a standalone @ metric
	result, err := rpn.ParseAndEvaluate("@GB")
	if err != nil {
		t.Fatalf("ParseAndEvaluate(@GB) failed: %v", err)
	}
	// Should push 1 with GB metric; display is "1"
	if result != "1" {
		t.Errorf("@GB result = %q, want %q", result, "1")
	}

	// The stack should have a number with the GB metric
	stack := rpn.GetCurrentStack()
	if len(stack) != 1 {
		t.Fatalf("expected 1 item on stack, got %d", len(stack))
	}
	m := stack[0].Metric()
	if m == nil || m.Name != "GB" {
		t.Errorf("stack[0].Metric() = %v, want GB", m)
	}
}

func TestMetricAwareArithmetic(t *testing.T) {
	reg := GetMetricRegistry()
	tolerance := 0.001

	tests := []struct {
		expr    string
		wantNum float64
		wantMet string
		wantErr bool
	}{
		{"100Mbps 50Mbps +", 150, "Mbps", false},
		{"1 100hr +", 100, "hr", false},
		{"3 4 +", 7, "Cool", false},
		{"100Mbps 1hr *", 360000000000, "bits", false},
		{"1000000000bits 1s /", 1000000000, "bps", false},
		{"100kmh 1hr *", 100000, "m", false},
		{"1km 1s /", 1000, "mps", false},
		{"1Gbps 1000Mbps /", 1, "Cool", false},
		{"100Mbps 2 ^", 10000, "Cool", false},
		{"100Mbps 1hr +", 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			vars := NewVariables()
			rpn := NewRPN(vars)
			result, err := rpn.ParseAndEvaluate(tt.expr)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got %q", result)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			resultVal, err := strconv.ParseFloat(result, 64)
			if err != nil {
				t.Fatalf("failed to parse result %q: %v", result, err)
			}
			if resultVal < tt.wantNum-tolerance || resultVal > tt.wantNum+tolerance {
				t.Errorf("result = %g, want %g (tolerance %g)", resultVal, tt.wantNum, tolerance)
			}
			stack := rpn.GetCurrentStack()
			if len(stack) > 0 {
				m := stack[0].Metric()
				expected, _ := reg.Find(tt.wantMet)
				if m != expected {
					t.Errorf("metric = %v, want %v", m, expected)
				}
			}
		})
	}
}
