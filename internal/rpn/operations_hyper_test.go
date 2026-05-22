// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"math"
	"strconv"
	"testing"
)

func TestHyperAddMetricAware(t *testing.T) {
	reg := GetMetricRegistry()
	tolerance := 0.001

	// Same category: [100Mbps 50Mbps 25Mbps [+] ] = 175Mbps
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("100Mbps 50Mbps 25Mbps [+]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal < 175-tolerance || resultVal > 175+tolerance {
		t.Errorf("result = %g, want 175", resultVal)
	}
	stack := rpn.GetCurrentStack()
	m := stack[len(stack)-1].Metric()
	mbps, _ := reg.Find("Mbps")
	if m != mbps {
		t.Errorf("metric = %v, want Mbps", m)
	}
}

func TestHyperAddCoolAbsorbing(t *testing.T) {
	reg := GetMetricRegistry()
	tolerance := 0.001

	// 5 (Cool) 100Mbps 10Mbps [+] ≈ 110.000005Mbps
	// Cool (factor 1.0) contributes 5 base units (bps) which is negligible
	// Result metric is Mbps (first non-Cool), value ≈ 110Mbps
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("5 100Mbps 10Mbps [+]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal < 110-tolerance || resultVal > 110+tolerance {
		t.Errorf("result = %g, want ~110", resultVal)
	}
	stack := rpn.GetCurrentStack()
	m := stack[len(stack)-1].Metric()
	mbps, _ := reg.Find("Mbps")
	if m != mbps {
		t.Errorf("metric = %v, want Mbps", m)
	}
}

func TestHyperAddIncompatible(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	_, err := rpn.ParseAndEvaluate("100Mbps 2hr [+]")
	if err == nil {
		t.Error("expected error for incompatible categories")
	}
}

func TestHyperAddMixedUnitsSameCategory(t *testing.T) {
	reg := GetMetricRegistry()
	tolerance := 0.001

	// 1km 500m 100m [+] = 1600m converted back to km = 1.6km
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("1km 500m 100m [+]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal < 1.6-tolerance || resultVal > 1.6+tolerance {
		t.Errorf("result = %g, want 1.6", resultVal)
	}
	stack := rpn.GetCurrentStack()
	m := stack[len(stack)-1].Metric()
	km, _ := reg.Find("km")
	if m != km {
		t.Errorf("metric = %v, want km", m)
	}
}

func TestHyperMultiplyMetricResultIsCool(t *testing.T) {
	reg := GetMetricRegistry()

	// 3 4 5 [*] = 60 with Cool metric
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("3 4 5 [*]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal != 60 {
		t.Errorf("result = %g, want 60", resultVal)
	}
	stack := rpn.GetCurrentStack()
	m := stack[len(stack)-1].Metric()
	cool, _ := reg.Find("Cool")
	if m != cool {
		t.Errorf("metric = %v, want Cool", m)
	}
}

func TestHyperMultiplyMixedMetricsIsCool(t *testing.T) {
	reg := GetMetricRegistry()

	// 100Mbps 2hr [*] = 200 (raw product, Cool metric)
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("100Mbps 2hr [*]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal != 200 {
		t.Errorf("result = %g, want 200", resultVal)
	}
	stack := rpn.GetCurrentStack()
	m := stack[len(stack)-1].Metric()
	cool, _ := reg.Find("Cool")
	if m != cool {
		t.Errorf("metric = %v, want Cool", m)
	}
}

func TestHyperSubtractMetricAware(t *testing.T) {
	reg := GetMetricRegistry()
	tolerance := 0.001

	// 1000Mbps 100Mbps 50Mbps [-] = 850Mbps
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("1000Mbps 100Mbps 50Mbps [-]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal < 850-tolerance || resultVal > 850+tolerance {
		t.Errorf("result = %g, want 850", resultVal)
	}
	stack := rpn.GetCurrentStack()
	m := stack[len(stack)-1].Metric()
	mbps, _ := reg.Find("Mbps")
	if m != mbps {
		t.Errorf("metric = %v, want Mbps", m)
	}
}

func TestHyperSubtractMixedUnits(t *testing.T) {
	reg := GetMetricRegistry()
	tolerance := 0.001

	// 2km 500m 100m [-] = (2000 - 500 - 100)m = 1400m = 1.4km
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("2km 500m 100m [-]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal < 1.4-tolerance || resultVal > 1.4+tolerance {
		t.Errorf("result = %g, want 1.4", resultVal)
	}
	stack := rpn.GetCurrentStack()
	m := stack[len(stack)-1].Metric()
	km, _ := reg.Find("km")
	if m != km {
		t.Errorf("metric = %v, want km", m)
	}
}

func TestHyperSubtractCoolAbsorbing(t *testing.T) {
	reg := GetMetricRegistry()

	// 100km 5 [-] = 100km - 5m = 99.995km
	// Cool (factor 1.0) contributes 5 base units (meters) = 0.005km
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("100km 5 [-]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal < 99.99 || resultVal > 99.996 {
		t.Errorf("result = %g, want ~99.995", resultVal)
	}
	stack := rpn.GetCurrentStack()
	m := stack[len(stack)-1].Metric()
	km, _ := reg.Find("km")
	if m != km {
		t.Errorf("metric = %v, want km", m)
	}
}

func TestHyperSubtractNegativeResult(t *testing.T) {
	reg := GetMetricRegistry()

	// 1km 2km [-] = -1km
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("1km 2km [-]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal < -1.1 || resultVal > -0.9 {
		t.Errorf("result = %g, want -1", resultVal)
	}
	stack := rpn.GetCurrentStack()
	m := stack[len(stack)-1].Metric()
	km, _ := reg.Find("km")
	if m != km {
		t.Errorf("metric = %v, want km", m)
	}
}

func TestHyperSubtractIncompatible(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	_, err := rpn.ParseAndEvaluate("100km 2hr [-]")
	if err == nil {
		t.Error("expected error for incompatible categories")
	}
}

func TestHyperDivideResultIsCool(t *testing.T) {
	reg := GetMetricRegistry()

	// 100 20 2 [/] = (100 / 20) / 2 = 2.5 with Cool metric
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("100 20 2 [/]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal != 2.5 {
		t.Errorf("result = %g, want 2.5", resultVal)
	}
	stack := rpn.GetCurrentStack()
	m := stack[len(stack)-1].Metric()
	cool, _ := reg.Find("Cool")
	if m != cool {
		t.Errorf("metric = %v, want Cool", m)
	}
}

func TestHyperPowerResultIsCool(t *testing.T) {
	reg := GetMetricRegistry()

	// 2 3 2 [^] = (2^3)^2 = 64 with Cool metric
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("2 3 2 [^]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal != 64 {
		t.Errorf("result = %g, want 64", resultVal)
	}
	stack := rpn.GetCurrentStack()
	m := stack[len(stack)-1].Metric()
	cool, _ := reg.Find("Cool")
	if m != cool {
		t.Errorf("metric = %v, want Cool", m)
	}
}

func TestHyperModuloMetricAware(t *testing.T) {
	reg := GetMetricRegistry()
	tolerance := 0.001

	// 10km 3km 2km [%] = ((10 % 3) % 2)km = (1 % 2)km = 1km
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("10km 3km 2km [%]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal < 1-tolerance || resultVal > 1+tolerance {
		t.Errorf("result = %g, want 1", resultVal)
	}
	stack := rpn.GetCurrentStack()
	m := stack[len(stack)-1].Metric()
	km, _ := reg.Find("km")
	if m != km {
		t.Errorf("metric = %v, want km", m)
	}
}

func TestHyperModuloMixedUnits(t *testing.T) {
	// 1000m 300m 200m [%] = ((1000 % 300) % 200)m = (100 % 200)m = 100m = 0.1km
	// Result uses first operand's metric (m)
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("1000m 300m 200m [%]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	tolerance := 0.001
	if resultVal < 100-tolerance || resultVal > 100+tolerance {
		t.Errorf("result = %g, want 100", resultVal)
	}
}

func TestHyperModuloIncompatible(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	_, err := rpn.ParseAndEvaluate("100km 2hr [%]")
	if err == nil {
		t.Error("expected error for incompatible categories")
	}
}

func TestHyperLog2CoolResult(t *testing.T) {
	reg := GetMetricRegistry()
	tolerance := 0.001

	// [lg] on all values → sum of log2: log2(2) + log2(4) + log2(8) = 1 + 2 + 3 = 6 → Cool
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("2 4 8 [lg]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal < 6-tolerance || resultVal > 6+tolerance {
		t.Errorf("result = %g, want 6", resultVal)
	}
	stack := rpn.GetCurrentStack()
	m := stack[len(stack)-1].Metric()
	cool, _ := reg.Find("Cool")
	if m != cool {
		t.Errorf("metric = %v, want Cool", m)
	}
}

func TestHyperLog10CoolResult(t *testing.T) {
	reg := GetMetricRegistry()
	tolerance := 0.001

	// [log] on all values → sum of log10: log10(10) + log10(100) = 1 + 2 = 3 → Cool
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("10 100 [log]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal < 3-tolerance || resultVal > 3+tolerance {
		t.Errorf("result = %g, want 3", resultVal)
	}
	stack := rpn.GetCurrentStack()
	m := stack[len(stack)-1].Metric()
	cool, _ := reg.Find("Cool")
	if m != cool {
		t.Errorf("metric = %v, want Cool", m)
	}
}

func TestHyperLnCoolResult(t *testing.T) {
	reg := GetMetricRegistry()
	tolerance := 0.001

	// [ln] on all values → sum of ln: ln(e) + ln(e*e) = 1 + 2 = 3 → Cool
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("2.718281828459045 7.38905609893065 [ln]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal < 3-tolerance || resultVal > 3+tolerance {
		t.Errorf("result = %g, want 3", resultVal)
	}
	stack := rpn.GetCurrentStack()
	m := stack[len(stack)-1].Metric()
	cool, _ := reg.Find("Cool")
	if m != cool {
		t.Errorf("metric = %v, want Cool", m)
	}
}

func TestHyperLogWithMetrics(t *testing.T) {
	reg := GetMetricRegistry()

	// [lg] with metrics still uses raw values and returns Cool
	// log2(100) + log2(1000) ≈ 6.64 + 9.97 ≈ 16.61
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("100Mbps 1000Mbps [lg]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	expected := math.Log2(100) + math.Log2(1000)
	tolerance := 0.01
	if resultVal < expected-tolerance || resultVal > expected+tolerance {
		t.Errorf("result = %g, want ~%g", resultVal, expected)
	}
	stack := rpn.GetCurrentStack()
	m := stack[len(stack)-1].Metric()
	cool, _ := reg.Find("Cool")
	if m != cool {
		t.Errorf("metric = %v, want Cool", m)
	}
}

func TestHyperAddAllCool(t *testing.T) {
	reg := GetMetricRegistry()

	// 1 2 3 [+] with no metrics = 6 Cool
	vars := NewVariables()
	rpn := NewRPN(vars)
	result, err := rpn.ParseAndEvaluate("1 2 3 [+]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultVal, _ := strconv.ParseFloat(result, 64)
	if resultVal != 6 {
		t.Errorf("result = %g, want 6", resultVal)
	}
	stack := rpn.GetCurrentStack()
	m := stack[len(stack)-1].Metric()
	cool, _ := reg.Find("Cool")
	if m != cool {
		t.Errorf("metric = %v, want Cool", m)
	}
}

func TestHyperModuloByZero(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	_, err := rpn.ParseAndEvaluate("10km 0km [%]")
	if err == nil {
		t.Error("expected error for modulo by zero")
	}
}

func TestHyperDivideByZero(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	_, err := rpn.ParseAndEvaluate("10 0 [/]")
	if err == nil {
		t.Error("expected error for division by zero")
	}
}

func TestHyperLogNonPositive(t *testing.T) {
	vars := NewVariables()
	rpn := NewRPN(vars)
	_, err := rpn.ParseAndEvaluate("0 5 [lg]")
	if err == nil {
		t.Error("expected error for log2 of non-positive number")
	}
}
