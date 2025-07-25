package benchmark

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	t.Run("command completion mode", func(t *testing.T) {
		config := Config{
			Phrase:  "",
			Timeout: 5 * time.Second,
			Command: []string{"echo", "hello"},
		}

		result := Run(config, 0)

		if !result.Found {
			t.Errorf("Expected result found, got %v", result.Found)
		}
		if result.Duration <= 0 {
			t.Errorf("Expected positive duration, got %v", result.Duration)
		}
	})

	t.Run("phrase detection mode success", func(t *testing.T) {
		config := Config{
			Phrase:  "hello",
			Timeout: 5 * time.Second,
			Command: []string{"echo", "hello world"},
		}

		result := Run(config, 0)

		if !result.Found {
			t.Errorf("Expected phrase found, got %v", result.Found)
		}
		if result.Duration <= 0 {
			t.Errorf("Expected positive duration, got %v", result.Duration)
		}
	})

	t.Run("phrase detection mode not found", func(t *testing.T) {
		config := Config{
			Phrase:  "notfound",
			Timeout: 100 * time.Millisecond,
			Command: []string{"sleep", "1"},
		}

		result := Run(config, 0)

		if result.Found {
			t.Errorf("Expected phrase not found, got %v", result.Found)
		}
	})

	t.Run("shell overhead subtraction", func(t *testing.T) {
		config := Config{
			Phrase:  "",
			Timeout: 5 * time.Second,
			Command: []string{"echo", "test"},
		}

		shellOverhead := 50 * time.Millisecond
		result := Run(config, shellOverhead)

		if !result.Found {
			t.Errorf("Expected result found, got %v", result.Found)
		}
	})

	t.Run("negative duration adjustment", func(t *testing.T) {
		config := Config{
			Phrase:  "",
			Timeout: 5 * time.Second,
			Command: []string{"echo", "test"},
		}

		shellOverhead := 10 * time.Second
		result := Run(config, shellOverhead)

		if !result.Found {
			t.Errorf("Expected result found, got %v", result.Found)
		}
		if result.Duration != 0 {
			t.Errorf("Expected duration adjusted to 0, got %v", result.Duration)
		}
	})
}

func TestCommandTimeout(t *testing.T) {
	config := Config{
		Phrase:  "",
		Timeout: 100 * time.Millisecond,
		Command: []string{"sleep", "1"},
	}

	result := Run(config, 0)

	if result.Found {
		t.Errorf("Expected timeout result (found=false), got found=%v", result.Found)
	}
}

func TestPhraseInStderr(t *testing.T) {
	config := Config{
		Phrase:  "error",
		Timeout: 5 * time.Second,
		Command: []string{"sh", "-c", "echo 'error message' >&2"},
	}

	result := Run(config, 0)

	if !result.Found {
		t.Errorf("Expected phrase found in stderr, got %v", result.Found)
	}
}

func TestKillProcess(t *testing.T) {
	cmd := exec.Command("sleep", "10")
	cmd.Start()

	if cmd.Process == nil {
		t.Fatal("Expected process to be started")
	}

	killProcess(cmd)

	err := cmd.Wait()
	if err == nil {
		t.Errorf("Expected process to be killed")
	}
	if !strings.Contains(err.Error(), "killed") && !strings.Contains(err.Error(), "signal") {
		t.Errorf("Expected kill-related error, got: %v", err)
	}
}

func TestScanOutput(t *testing.T) {
	t.Run("phrase found", func(t *testing.T) {
		reader := strings.NewReader("line1\nphrase here\nline3")
		done := make(chan time.Duration, 1)
		cancel := make(chan struct{})
		startTime := time.Now()

		go scanOutput(reader, "phrase", startTime, done, cancel)

		select {
		case duration := <-done:
			if duration <= 0 {
				t.Errorf("Expected positive duration, got %v", duration)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("Expected phrase to be found quickly")
		}
	})

	t.Run("phrase not found", func(t *testing.T) {
		reader := strings.NewReader("line1\nline2\nline3")
		done := make(chan time.Duration, 1)
		cancel := make(chan struct{})
		startTime := time.Now()

		go scanOutput(reader, "notfound", startTime, done, cancel)

		select {
		case <-done:
			t.Errorf("Did not expect phrase to be found")
		case <-time.After(100 * time.Millisecond):
		}
	})

	t.Run("cancelled", func(t *testing.T) {
		reader := strings.NewReader("line1\nphrase here\nline3")
		done := make(chan time.Duration, 1)
		cancel := make(chan struct{})
		startTime := time.Now()

		close(cancel)
		go scanOutput(reader, "phrase", startTime, done, cancel)

		select {
		case <-done:
			t.Errorf("Did not expect phrase to be found after cancellation")
		case <-time.After(100 * time.Millisecond):
		}
	})
}
