package shellcalibration

import (
	"testing"
)

func TestCalibrateShellOverhead(t *testing.T) {
	t.Run("shell calibration success", func(t *testing.T) {
		overhead := CalibrateShellOverhead(3)
		if overhead <= 0 {
			t.Errorf("Expected positive overhead, got %v", overhead)
		}
	})

	t.Run("minimal calibration runs", func(t *testing.T) {
		overhead := CalibrateShellOverhead(1)
		if overhead < 0 {
			t.Errorf("Expected non-negative overhead, got %v", overhead)
		}
	})
}
