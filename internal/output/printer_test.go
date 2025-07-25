package output

import (
	"bytes"
	"chrono/internal/benchmark"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{0, "0.000s"},
		{100 * time.Millisecond, "0.100s"},
		{1*time.Second + 500*time.Millisecond, "1.500s"},
		{2 * time.Minute, "120.000s"},
	}

	for _, test := range tests {
		result := FormatDuration(test.duration)
		if result != test.expected {
			t.Errorf("FormatDuration(%v) = %s, expected %s", test.duration, result, test.expected)
		}
	}
}

func TestPrintCalibration(t *testing.T) {
	output := captureOutput(func() {
		PrintCalibration(5)
	})

	expected := "Running 5 calibration runs"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain '%s', got '%s'", expected, output)
	}
}

func TestPrintShellOverhead(t *testing.T) {
	output := captureOutput(func() {
		PrintShellOverhead(50 * time.Millisecond)
	})

	expected := "Shell overhead:"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain '%s', got '%s'", expected, output)
	}
	if !strings.Contains(output, "0.050s") {
		t.Errorf("Expected output to contain duration '0.050s', got '%s'", output)
	}
}

func TestPrintWarmupHeader(t *testing.T) {
	output := captureOutput(func() {
		PrintWarmupHeader(3)
	})

	expected := "Running 3 warmup runs"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain '%s', got '%s'", expected, output)
	}
}

func TestPrintWarmupResult(t *testing.T) {
	t.Run("successful result", func(t *testing.T) {
		result := benchmark.Result{
			Duration: 100 * time.Millisecond,
			Found:    true,
		}

		output := captureOutput(func() {
			PrintWarmupResult(1, result)
		})

		if !strings.Contains(output, "Warmup 1:") {
			t.Errorf("Expected output to contain 'Warmup 1:', got '%s'", output)
		}
		if !strings.Contains(output, "0.100s") {
			t.Errorf("Expected output to contain duration, got '%s'", output)
		}
	})

	t.Run("failed result", func(t *testing.T) {
		result := benchmark.Result{
			Duration: 0,
			Found:    false,
		}

		output := captureOutput(func() {
			PrintWarmupResult(2, result)
		})

		if !strings.Contains(output, "Warmup 2:") {
			t.Errorf("Expected output to contain 'Warmup 2:', got '%s'", output)
		}
		if !strings.Contains(output, "phrase not found") {
			t.Errorf("Expected output to contain 'phrase not found', got '%s'", output)
		}
	})
}

func TestPrintBenchmarkHeader(t *testing.T) {
	output := captureOutput(func() {
		PrintBenchmarkHeader(10)
	})

	expected := "Running 10 benchmark runs"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain '%s', got '%s'", expected, output)
	}
}

func TestPrintBenchmarkResult(t *testing.T) {
	t.Run("successful result", func(t *testing.T) {
		result := benchmark.Result{
			Duration: 250 * time.Millisecond,
			Found:    true,
		}

		output := captureOutput(func() {
			PrintBenchmarkResult(3, result)
		})

		if !strings.Contains(output, "Run 3:") {
			t.Errorf("Expected output to contain 'Run 3:', got '%s'", output)
		}
		if !strings.Contains(output, "0.250s") {
			t.Errorf("Expected output to contain duration, got '%s'", output)
		}
	})

	t.Run("failed result", func(t *testing.T) {
		result := benchmark.Result{
			Duration: 0,
			Found:    false,
		}

		output := captureOutput(func() {
			PrintBenchmarkResult(4, result)
		})

		if !strings.Contains(output, "Run 4:") {
			t.Errorf("Expected output to contain 'Run 4:', got '%s'", output)
		}
		if !strings.Contains(output, "phrase not found") {
			t.Errorf("Expected output to contain 'phrase not found', got '%s'", output)
		}
	})
}

func TestPrintSummary(t *testing.T) {
	t.Run("command completion mode", func(t *testing.T) {
		config := benchmark.Config{
			Phrase:          "",
			Warmups:         2,
			SkipCalibration: false,
		}
		results := []benchmark.Result{
			{Duration: 100 * time.Millisecond, Found: true},
			{Duration: 200 * time.Millisecond, Found: true},
		}
		shellOverhead := 10 * time.Millisecond

		output := captureOutput(func() {
			PrintSummary(results, config, shellOverhead)
		})

		if !strings.Contains(output, "Command completion timing") {
			t.Errorf("Expected output to contain 'Command completion timing', got '%s'", output)
		}
		if !strings.Contains(output, "(2 warmups)") {
			t.Errorf("Expected output to contain warmup info, got '%s'", output)
		}
		if !strings.Contains(output, "shell overhead") {
			t.Errorf("Expected output to contain shell overhead info, got '%s'", output)
		}
		if !strings.Contains(output, "Mean:") {
			t.Errorf("Expected output to contain statistics, got '%s'", output)
		}
	})

	t.Run("phrase mode", func(t *testing.T) {
		config := benchmark.Config{
			Phrase:          "test phrase",
			Warmups:         0,
			SkipCalibration: true,
		}
		results := []benchmark.Result{
			{Duration: 150 * time.Millisecond, Found: true},
		}

		output := captureOutput(func() {
			PrintSummary(results, config, 0)
		})

		if !strings.Contains(output, "test phrase") {
			t.Errorf("Expected output to contain phrase, got '%s'", output)
		}
		if !strings.Contains(output, "Time:") {
			t.Errorf("Expected single result format, got '%s'", output)
		}
	})

	t.Run("no successful runs - command mode", func(t *testing.T) {
		config := benchmark.Config{
			Phrase: "",
		}
		results := []benchmark.Result{
			{Duration: 0, Found: false},
		}

		output := captureOutput(func() {
			PrintSummary(results, config, 0)
		})

		if !strings.Contains(output, "all commands timed out") {
			t.Errorf("Expected timeout message for command mode, got '%s'", output)
		}
	})

	t.Run("no successful runs - phrase mode", func(t *testing.T) {
		config := benchmark.Config{
			Phrase: "notfound",
		}
		results := []benchmark.Result{
			{Duration: 0, Found: false},
		}

		output := captureOutput(func() {
			PrintSummary(results, config, 0)
		})

		if !strings.Contains(output, "phrase was not found") {
			t.Errorf("Expected phrase not found message, got '%s'", output)
		}
	})

	t.Run("mixed results", func(t *testing.T) {
		config := benchmark.Config{
			Phrase: "test",
		}
		results := []benchmark.Result{
			{Duration: 100 * time.Millisecond, Found: true},
			{Duration: 0, Found: false},
			{Duration: 200 * time.Millisecond, Found: true},
		}

		output := captureOutput(func() {
			PrintSummary(results, config, 0)
		})

		if !strings.Contains(output, "Failed:") {
			t.Errorf("Expected failure indication, got '%s'", output)
		}
		if !strings.Contains(output, "1/3") {
			t.Errorf("Expected failure count ratio, got '%s'", output)
		}
		if !strings.Contains(output, "Mean:") {
			t.Errorf("Expected statistics for multiple results, got '%s'", output)
		}
	})

	t.Run("all failed results", func(t *testing.T) {
		config := benchmark.Config{
			Phrase: "test",
		}
		results := []benchmark.Result{
			{Duration: 0, Found: false},
			{Duration: 0, Found: false},
		}

		output := captureOutput(func() {
			PrintSummary(results, config, 0)
		})

		if !strings.Contains(output, "No successful runs") {
			t.Errorf("Expected no successful runs message, got '%s'", output)
		}
	})
}
