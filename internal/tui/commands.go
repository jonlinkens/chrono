package tui

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
	"time"

	"chrono/internal/shellcalibration"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) tickCmd() tea.Cmd {
	return tea.Tick(TickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) runCalibration() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		overhead := shellcalibration.CalibrateShellOverhead(m.config.CalibrationRuns)
		return calibrationCompleteMsg{overhead: overhead}
	})
}

func (m Model) startWarmup() tea.Cmd {
	m.state = StateWarmup
	return tea.Batch(
		func() tea.Msg { return runStartMsg{isWarmup: true} },
		m.runWithOutput(true),
	)
}

func (m Model) startBenchmark() tea.Cmd {
	m.state = StateBenchmarking
	return tea.Batch(
		func() tea.Msg { return runStartMsg{isWarmup: false} },
		m.runWithOutput(false),
	)
}

func (m Model) runWithOutput(isWarmup bool) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command(m.config.Command[0], m.config.Command[1:]...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return errorMsg{err: err}
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return errorMsg{err: err}
		}

		startTime := time.Now()
		if err := cmd.Start(); err != nil {
			return errorMsg{err: err}
		}

		if m.config.Phrase == "" {
			done := make(chan time.Duration, 1)
			outputLines := make(chan string, OutputChannelBuffer)

			go m.captureOutput(stdout, outputLines, false)
			go m.captureOutput(stderr, outputLines, true)

			go func() {
				cmd.Wait()
				close(outputLines)
				duration := time.Since(startTime)
				done <- duration
			}()

			return startStreamingMsg{
				cmd:         cmd,
				stdout:      stdout,
				stderr:      stderr,
				startTime:   startTime,
				isWarmup:    isWarmup,
				phrase:      m.config.Phrase,
				done:        done,
				cancel:      nil,
				outputLines: outputLines,
			}
		} else {
			done := make(chan time.Duration, 1)
			cancel := make(chan struct{})
			outputLines := make(chan string, OutputChannelBuffer)

			go m.scanOutputWithStreaming(stdout, m.config.Phrase, startTime, done, cancel, outputLines, false)
			go m.scanOutputWithStreaming(stderr, m.config.Phrase, startTime, done, cancel, outputLines, true)

			go func() {
				cmd.Wait()
				close(outputLines)
			}()

			return startStreamingMsg{
				cmd:         cmd,
				stdout:      stdout,
				stderr:      stderr,
				startTime:   startTime,
				isWarmup:    isWarmup,
				phrase:      m.config.Phrase,
				done:        done,
				cancel:      cancel,
				outputLines: outputLines,
			}
		}
	}
}

func (m Model) captureOutput(reader io.ReadCloser, outputLines chan string, isStderr bool) {
	defer func() {
		reader.Close()
	}()

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		if isStderr {
			line = "stderr: " + line
		}

		select {
		case outputLines <- line:
		default:
		}
	}
}

func (m Model) scanOutputWithStreaming(reader io.ReadCloser, phrase string, startTime time.Time, done chan time.Duration, cancel chan struct{}, outputLines chan string, isStderr bool) {
	defer func() {
		reader.Close()
	}()

	scanner := bufio.NewScanner(reader)
	phraseFound := false

	for scanner.Scan() {
		select {
		case <-cancel:
			return
		default:
		}

		line := scanner.Text()
		originalLine := line
		if isStderr {
			line = "stderr: " + line
		}

		select {
		case outputLines <- line:
		case <-cancel:
			return
		}

		if phrase != "" && !phraseFound && strings.Contains(originalLine, phrase) {
			phraseFound = true

			select {
			case outputLines <- "Match found!":
			case <-cancel:
				return
			}

			select {
			case done <- time.Since(startTime):
			case <-cancel:
				return
			default:
			}
			return
		}
	}

	if scanner.Err() != nil {
		return
	}
}
