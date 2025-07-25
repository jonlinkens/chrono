package stats

import (
	"testing"
	"time"
)

func TestCalculateStatistics(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		stats := CalculateStatistics([]time.Duration{})

		if stats.Mean != 0 {
			t.Errorf("Expected mean 0, got %v", stats.Mean)
		}
		if stats.Min != 0 {
			t.Errorf("Expected min 0, got %v", stats.Min)
		}
		if stats.Max != 0 {
			t.Errorf("Expected max 0, got %v", stats.Max)
		}
		if stats.Range != 0 {
			t.Errorf("Expected range 0, got %v", stats.Range)
		}
	})

	t.Run("single value", func(t *testing.T) {
		durations := []time.Duration{100 * time.Millisecond}
		stats := CalculateStatistics(durations)

		expected := 100 * time.Millisecond
		if stats.Mean != expected {
			t.Errorf("Expected mean %v, got %v", expected, stats.Mean)
		}
		if stats.Min != expected {
			t.Errorf("Expected min %v, got %v", expected, stats.Min)
		}
		if stats.Max != expected {
			t.Errorf("Expected max %v, got %v", expected, stats.Max)
		}
		if stats.Range != 0 {
			t.Errorf("Expected range 0, got %v", stats.Range)
		}
	})

	t.Run("multiple values", func(t *testing.T) {
		durations := []time.Duration{
			100 * time.Millisecond,
			200 * time.Millisecond,
			300 * time.Millisecond,
		}
		stats := CalculateStatistics(durations)

		expectedMean := 200 * time.Millisecond
		expectedMin := 100 * time.Millisecond
		expectedMax := 300 * time.Millisecond
		expectedRange := 200 * time.Millisecond

		if stats.Mean != expectedMean {
			t.Errorf("Expected mean %v, got %v", expectedMean, stats.Mean)
		}
		if stats.Min != expectedMin {
			t.Errorf("Expected min %v, got %v", expectedMin, stats.Min)
		}
		if stats.Max != expectedMax {
			t.Errorf("Expected max %v, got %v", expectedMax, stats.Max)
		}
		if stats.Range != expectedRange {
			t.Errorf("Expected range %v, got %v", expectedRange, stats.Range)
		}
	})

	t.Run("unsorted values", func(t *testing.T) {
		durations := []time.Duration{
			300 * time.Millisecond,
			100 * time.Millisecond,
			200 * time.Millisecond,
		}
		stats := CalculateStatistics(durations)

		expectedMin := 100 * time.Millisecond
		expectedMax := 300 * time.Millisecond

		if stats.Min != expectedMin {
			t.Errorf("Expected min %v, got %v", expectedMin, stats.Min)
		}
		if stats.Max != expectedMax {
			t.Errorf("Expected max %v, got %v", expectedMax, stats.Max)
		}
	})
}
