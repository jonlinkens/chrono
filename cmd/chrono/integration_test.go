package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestBuildAndRun(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "test-benchmark", ".")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("test-benchmark")

	t.Run("help flag", func(t *testing.T) {
		cmd := exec.Command("./test-benchmark", "--help")
		output, _ := cmd.CombinedOutput()
		outputStr := string(output)
		if !strings.Contains(outputStr, "Usage of") {
			t.Errorf("Expected help text in output, got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "-phrase") {
			t.Errorf("Expected -phrase flag in help, got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "-runs") {
			t.Errorf("Expected -runs flag in help, got: %s", outputStr)
		}
	})

	t.Run("no arguments", func(t *testing.T) {
		cmd := exec.Command("./test-benchmark")
		output, err := cmd.CombinedOutput()
		if err == nil {
			t.Error("Expected non-zero exit for no arguments")
		}
		outputStr := string(output)
		if !strings.Contains(outputStr, "command to benchmark is required") {
			t.Errorf("Expected error message in output, got: %s", outputStr)
		}
	})

	t.Run("successful run with CLI flag", func(t *testing.T) {
		cmd := exec.Command("./test-benchmark", "--cli", "--runs", "1", "--skip-calibration", "echo", "test")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("Expected successful run, got error: %v, output: %s", err, string(output))
		}
		outputStr := string(output)
		if !strings.Contains(outputStr, "Running 1 benchmark runs") {
			t.Errorf("Expected benchmark output, got: %s", outputStr)
		}
	})

	t.Run("phrase detection", func(t *testing.T) {
		cmd := exec.Command("./test-benchmark", "--cli", "--phrase", "hello", "--runs", "1", "--skip-calibration", "echo", "hello world")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("Expected successful run, got error: %v, output: %s", err, string(output))
		}
		outputStr := string(output)
		if !strings.Contains(outputStr, "hello") {
			t.Errorf("Expected phrase in output, got: %s", outputStr)
		}
	})

	t.Run("calibration", func(t *testing.T) {
		cmd := exec.Command("./test-benchmark", "--cli", "--calibration", "2", "--runs", "1", "echo", "test")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("Expected successful run, got error: %v, output: %s", err, string(output))
		}
		outputStr := string(output)
		if !strings.Contains(outputStr, "Running 2 calibration runs") {
			t.Errorf("Expected calibration output, got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Shell overhead:") {
			t.Errorf("Expected shell overhead output, got: %s", outputStr)
		}
	})

	t.Run("warmup runs", func(t *testing.T) {
		cmd := exec.Command("./test-benchmark", "--cli", "--warmups", "2", "--runs", "1", "--skip-calibration", "echo", "test")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("Expected successful run, got error: %v, output: %s", err, string(output))
		}
		outputStr := string(output)
		if !strings.Contains(outputStr, "Running 2 warmup runs") {
			t.Errorf("Expected warmup output, got: %s", outputStr)
		}
	})

	t.Run("timeout flag", func(t *testing.T) {
		cmd := exec.Command("./test-benchmark", "--cli", "--timeout", "500ms", "--runs", "1", "--skip-calibration", "sleep", "2")
		start := time.Now()
		output, _ := cmd.CombinedOutput()
		duration := time.Since(start)

		if duration > 2*time.Second {
			t.Errorf("Expected timeout around 500ms, but took %v", duration)
		}
		outputStr := string(output)
		if !strings.Contains(outputStr, "timed out") && !strings.Contains(outputStr, "timeout") {
			t.Logf("Timeout test output: %s", outputStr)
		}
	})
	t.Run("invalid command", func(t *testing.T) {
		cmd := exec.Command("./test-benchmark", "--cli", "--runs", "1", "--skip-calibration", "nonexistent-command-12345")
		output, err := cmd.CombinedOutput()
		if err == nil {
			t.Error("Expected non-zero exit for invalid command")
		}
		outputStr := string(output)
		if !strings.Contains(outputStr, "not found") && !strings.Contains(outputStr, "executable file not found") {
			t.Errorf("Expected command not found error, got: %s", outputStr)
		}
	})

	t.Run("multiple runs summary", func(t *testing.T) {
		cmd := exec.Command("./test-benchmark", "--cli", "--runs", "3", "--skip-calibration", "echo", "test")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("Expected successful run, got error: %v, output: %s", err, string(output))
		}
		outputStr := string(output)
		if !strings.Contains(outputStr, "Running 3 benchmark runs") {
			t.Errorf("Expected 3 benchmark runs header, got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Mean:") {
			t.Errorf("Expected summary statistics with Mean, got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Min:") {
			t.Errorf("Expected summary statistics with Min, got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Max:") {
			t.Errorf("Expected summary statistics with Max, got: %s", outputStr)
		}
	})

	t.Run("zero runs behavior", func(t *testing.T) {
		cmd := exec.Command("./test-benchmark", "--cli", "--runs", "0", "--skip-calibration", "echo", "test")
		output, _ := cmd.CombinedOutput()
		outputStr := string(output)
		if !strings.Contains(outputStr, "Running 0 benchmark runs") {
			t.Errorf("Expected to see 0 runs header, got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "No successful runs") {
			t.Errorf("Expected 'No successful runs' message, got: %s", outputStr)
		}
	})

	t.Run("negative runs causes panic", func(t *testing.T) {
		cmd := exec.Command("./test-benchmark", "--cli", "--runs", "-1", "--skip-calibration", "echo", "test")
		output, err := cmd.CombinedOutput()
		if err == nil {
			t.Error("Expected non-zero exit for negative runs")
		}
		outputStr := string(output)
		if !strings.Contains(outputStr, "panic") {
			t.Errorf("Expected panic for negative runs, got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "makeslice: cap out of range") {
			t.Errorf("Expected makeslice error for negative runs, got: %s", outputStr)
		}
	})
}
