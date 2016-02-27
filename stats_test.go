package honeybadger

import (
	"testing"
	"time"
)

func TestPercentileEmpty(t *testing.T) {
	values := []time.Duration{}
	percentile_90 := percentile(values, 90)
	if percentile_90 != 0.0 {
		t.Errorf("Expected percentile to return 0.0 for no values.")
	}
}

func TestPercentileOne(t *testing.T) {
	values := []time.Duration{time.Millisecond}
	percentile_90 := percentile(values, 90)
	if percentile_90 != 1.0 {
		t.Errorf("Expected percentile %v but %v was returned.", 1.0, percentile_90)
	}
}

func TestPercentile(t *testing.T) {
	values := make([]time.Duration, 100)
	// 0..100
	for i := 0; i < 100; i++ {
		values[i] = time.Duration(i) * time.Millisecond
	}
	percentile_90 := percentile(values, 90)
	if percentile_90 != 90 {
		t.Errorf("Expected percentile %v but %v was returned.", 90, percentile_90)
	}
}

func TestMean(t *testing.T) {
	values := make([]time.Duration, 6)
	// 1, 2, 3, 4, 5
	for i := 1; i < 7; i++ {
		values[i-1] = time.Duration(i) * time.Millisecond
	}
	result := mean(values)
	if result != 3.5 {
		t.Errorf("Expected mean %v but %v was returned.", 3.5, result)
	}
}

func TestMedian(t *testing.T) {
	values := make([]time.Duration, 7)
	// 1, 2, 3, 4, 5, 6, 7
	for i := 1; i < 8; i++ {
		values[i-1] = time.Duration(i) * time.Millisecond
	}
	result := median(values)
	if result != 4 {
		t.Errorf("Expected median %v but %v was returned.", 4, result)
	}
}

func TestStdDevEmpty(t *testing.T) {
	values := []time.Duration{}
	result := stdDev(values, mean(values))
	if result != 0.0 {
		t.Errorf("Expected stdDev to return 0.0 for no values.")
	}
}

func TestStdDevOne(t *testing.T) {
	values := []time.Duration{3 * time.Millisecond}
	result := stdDev(values, mean(values))
	if result != 0.0 {
		t.Errorf("Expected stdDev to return 0.0 for one value.")
	}
}

func TestStdDev(t *testing.T) {
	values := make([]time.Duration, 100)
	// 0..100
	for i := 0; i < 100; i++ {
		values[i] = time.Duration(i) * time.Millisecond
	}
	result := stdDev(values, mean(values))
	if result != 29.011491975882016 {
		t.Errorf("Expected stdDev %v but %v was returned.", 29.011491975882016, result)
	}
}

func TestMin(t *testing.T) {
	values := make([]time.Duration, 6)
	// 1, 2, 3, 4, 5, 6
	for i := 1; i < 7; i++ {
		values[i-1] = time.Duration(i) * time.Millisecond
	}
	result := min(values)
	if result != 1 {
		t.Errorf("Expected min 1 but %v was returned.", result)
	}
}

func TestMax(t *testing.T) {
	values := make([]time.Duration, 6)
	// 1, 2, 3, 4, 5, 6
	for i := 1; i < 7; i++ {
		values[i-1] = time.Duration(i) * time.Millisecond
	}
	result := max(values)
	if result != 6 {
		t.Errorf("Expected max 6 but %v was returned.", result)
	}
}
