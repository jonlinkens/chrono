package tui

import (
	"fmt"
	"time"

	"chrono/internal/benchmark"

	"github.com/charmbracelet/lipgloss"
)

func formatDuration(d time.Duration) string {
	return fmt.Sprintf("%.3fs", d.Seconds())
}

func (m Model) getMaxScrollOffset() int {
	contentHeight := m.height - 6
	availableLines := (contentHeight - 1) - 2

	if len(m.commandOutput) <= availableLines {
		return 0
	}
	return len(m.commandOutput) - availableLines
}

func (m Model) autoScrollToBottom() bool {
	maxScroll := m.getMaxScrollOffset()
	return m.scrollOffset >= maxScroll
}

func truncateToWidth(s string, width int) string {
	if lipgloss.Width(s) <= width {
		return s
	}

	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		candidate := string(runes[:i])
		if lipgloss.Width(candidate) <= width {
			return candidate
		}
	}
	return ""
}

func (m Model) formatCommandDisplay() string {
	cmd := fmt.Sprintf("Command: %s", m.config.Command[0])
	if len(m.config.Command) > 1 {
		cmd += fmt.Sprintf(" %v", m.config.Command[1:])
	}
	if m.config.Phrase != "" {
		cmd += fmt.Sprintf(" (phrase matched: \"%s\")", m.config.Phrase)
	}
	return cmd
}

func (m Model) filterValidResults() ([]time.Duration, int) {
	validResults := make([]time.Duration, 0, len(m.benchmarkResults))
	failedCount := 0

	for _, result := range m.benchmarkResults {
		if result.Found {
			validResults = append(validResults, result.Duration)
		} else {
			failedCount++
		}
	}

	return validResults, failedCount
}

func (m Model) createAdjustedResult(duration time.Duration, found bool) benchmark.Result {
	adjustedDuration := max(duration-m.shellOverhead, 0)
	return benchmark.Result{
		Duration: adjustedDuration,
		Found:    found,
	}
}

func (m Model) createTimeoutResult() benchmark.Result {
	return benchmark.Result{
		Duration: 0,
		Found:    false,
	}
}

func (m Model) calculateBoundedLeftWidth() int {
	leftContentWidth := m.calculateLeftContentWidth()
	maxLeftWidth := int(float64(m.width) * maxLeftWidthMultiplier)

	if leftContentWidth < minLeftWidth {
		leftContentWidth = minLeftWidth
	}
	if leftContentWidth > maxLeftWidth {
		leftContentWidth = maxLeftWidth
	}

	return leftContentWidth
}
