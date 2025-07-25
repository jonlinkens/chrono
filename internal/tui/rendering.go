package tui

import (
	"fmt"
	"strings"

	"chrono/internal/colours"
	"chrono/internal/stats"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) calculateLeftContentWidth() int {
	maxWidth := 0

	cmd := m.formatCommandDisplay()
	if m.config.Phrase != "" {
		cmd += fmt.Sprintf("\nPhrase: \"%s\"", m.config.Phrase)
	}
	for line := range strings.SplitSeq(cmd, "\n") {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	var timeoutStr string
	if m.config.Timeout > 0 {
		timeoutStr = m.config.Timeout.String()
	} else {
		timeoutStr = "none"
	}
	configLines := []string{
		fmt.Sprintf("Warmups: %d", m.config.Warmups),
		fmt.Sprintf("Runs: %d", m.config.Runs),
		fmt.Sprintf("Timeout: %s", timeoutStr),
	}
	if !m.config.SkipCalibration {
		configLines = append(configLines, fmt.Sprintf("Shell overhead: %s", formatDuration(m.shellOverhead)))
	}
	for _, line := range configLines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	var statusLine string
	switch m.state {
	case StateCalibrating:
		statusLine = "Status: Calibrating"
	case StateWarmup:
		statusLine = fmt.Sprintf("Status: Warmup (%d/%d)", m.warmupProgress, m.config.Warmups)
	case StateBenchmarking:
		statusLine = fmt.Sprintf("Status: Benchmarking (%d/%d)", m.benchmarkProgress+1, m.config.Runs)
	case StateCompleted:
		statusLine = "Status: Completed"
	}
	if len(statusLine) > maxWidth {
		maxWidth = len(statusLine)
	}

	timingLines := []string{"Run Timings:"}
	for i := range m.warmupResults {
		line := fmt.Sprintf("  W%d: %s", i+1, formatDuration(m.warmupResults[i].Duration))
		if !m.warmupResults[i].Found {
			line = fmt.Sprintf("  W%d: timeout", i+1)
		}
		timingLines = append(timingLines, line)
	}
	for i := range m.benchmarkResults {
		line := fmt.Sprintf("  #%d: %s", i+1, formatDuration(m.benchmarkResults[i].Duration))
		if !m.benchmarkResults[i].Found {
			line = fmt.Sprintf("  #%d: timeout", i+1)
		}
		timingLines = append(timingLines, line)
	}

	if m.state == StateCompleted {
		summaryLines := []string{"Final Results:"}
		validResults, _ := m.filterValidResults()
		if len(validResults) > 1 {
			stats := stats.CalculateStatistics(validResults)
			summaryLines = append(summaryLines,
				fmt.Sprintf("  Mean: %s", formatDuration(stats.Mean)),
				fmt.Sprintf("  Min: %s", formatDuration(stats.Min)),
				fmt.Sprintf("  Max: %s", formatDuration(stats.Max)),
				fmt.Sprintf("  Range: %s", formatDuration(stats.Range)),
			)
		}
		timingLines = append(timingLines, summaryLines...)
	}

	for _, line := range timingLines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	return maxWidth + 4
}

func (m Model) renderLeftColumnContentText() string {
	var s strings.Builder

	if m.err != nil {
		errorText := lipgloss.NewStyle().
			Foreground(lipgloss.Color(colours.Red)).
			Bold(true).
			Render(fmt.Sprintf("Error: %v", m.err))
		s.WriteString(errorText)
		return s.String()
	}

	commandStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colours.Blue)).
		Bold(true)

	cmd := m.formatCommandDisplay()
	if m.config.Phrase != "" {
		cmd += fmt.Sprintf("\nPhrase: \"%s\"", m.config.Phrase)
	}
	s.WriteString(commandStyle.Render(cmd))
	s.WriteString("\n\n")

	configStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colours.Subtext0))

	var configInfo strings.Builder
	configInfo.WriteString(fmt.Sprintf("Warmups: %d\n", m.config.Warmups))
	configInfo.WriteString(fmt.Sprintf("Runs: %d\n", m.config.Runs))
	if m.config.Timeout > 0 {
		configInfo.WriteString(fmt.Sprintf("Timeout: %s\n", m.config.Timeout))
	} else {
		configInfo.WriteString("Timeout: none\n")
	}

	if !m.config.SkipCalibration {
		configInfo.WriteString(fmt.Sprintf("Shell overhead: %s\n", formatDuration(m.shellOverhead)))
	}

	s.WriteString(configStyle.Render(configInfo.String()))
	s.WriteString("\n")

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colours.Yellow)).
		Bold(true)

	var status string
	switch m.state {
	case StateCalibrating:
		status = "Status: Calibrating"
	case StateWarmup:
		status = fmt.Sprintf("Status: Warmup (%d/%d)", m.warmupProgress, m.config.Warmups)
	case StateBenchmarking:
		status = fmt.Sprintf("Status: Benchmarking (%d/%d)", m.benchmarkProgress+1, m.config.Runs)
	case StateCompleted:
		status = "Status: Completed"
	}
	s.WriteString(statusStyle.Render(status))
	s.WriteString("\n\n")

	runTimingsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colours.Blue)).
		Bold(true)
	s.WriteString(runTimingsStyle.Render("Run Timings:"))
	s.WriteString("\n")

	warmupStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colours.Lavender))

	for i, result := range m.warmupResults {
		if result.Found {
			s.WriteString(warmupStyle.Render(fmt.Sprintf("  W%d: %s", i+1, formatDuration(result.Duration))))
		} else {
			s.WriteString(warmupStyle.Render(fmt.Sprintf("  W%d: timeout", i+1)))
		}
		s.WriteString("\n")
	}

	if m.isRunning && m.state == StateWarmup {
		currentStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(colours.Green))
		s.WriteString(currentStyle.Render(fmt.Sprintf("  W%d: %s", m.warmupProgress+1, formatDuration(m.elapsedTime))))
		s.WriteString("\n")
	}

	timingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colours.Text))

	for i, result := range m.benchmarkResults {
		if result.Found {
			s.WriteString(timingStyle.Render(fmt.Sprintf("  #%d: %s", i+1, formatDuration(result.Duration))))
		} else {
			s.WriteString(timingStyle.Render(fmt.Sprintf("  #%d: timeout", i+1)))
		}
		s.WriteString("\n")
	}

	if m.isRunning && m.state == StateBenchmarking {
		currentStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(colours.Green))
		s.WriteString(currentStyle.Render(fmt.Sprintf("  #%d: %s", m.benchmarkProgress+1, formatDuration(m.elapsedTime))))
		s.WriteString("\n")
	}

	if m.state == StateCompleted {
		s.WriteString("\n")
		validResults, failedCount := m.filterValidResults()

		if len(validResults) == 0 {
			errorStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(colours.Red))
			if m.config.Phrase == "" {
				s.WriteString(errorStyle.Render("No successful runs - all commands timed out"))
			} else {
				s.WriteString(errorStyle.Render("No successful runs - phrase was not found"))
			}
		} else {
			summaryStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(colours.Green)).
				Bold(true)
			s.WriteString(summaryStyle.Render("Final Results:"))
			s.WriteString("\n")

			if failedCount > 0 {
				failedStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color(colours.Red))
				s.WriteString(failedStyle.Render(fmt.Sprintf("  Failed: %d/%d", failedCount, len(m.benchmarkResults))))
				s.WriteString("\n")
			}

			if len(validResults) == 1 {
				s.WriteString(fmt.Sprintf("  Time: %s", formatDuration(validResults[0])))
			} else {
				stats := stats.CalculateStatistics(validResults)
				s.WriteString(fmt.Sprintf("  Mean: %s\n", formatDuration(stats.Mean)))
				s.WriteString(fmt.Sprintf("  Min: %s\n", formatDuration(stats.Min)))
				s.WriteString(fmt.Sprintf("  Max: %s\n", formatDuration(stats.Max)))
				s.WriteString(fmt.Sprintf("  Range: %s", formatDuration(stats.Range)))
			}
		}
	}

	return s.String()
}

func (m Model) renderRightColumnContentText(maxWidth, maxHeight int) string {
	var s strings.Builder

	if len(m.commandOutput) == 0 && !m.isRunning {
		s.WriteString("No output yet...")
		return s.String()
	}

	separatorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colours.Yellow)).
		Bold(true)

	matchStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colours.Green)).
		Background(lipgloss.Color(colours.Surface0)).
		Bold(true)

	regularStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colours.Text))

	availableLines := maxHeight - 2

	startIdx := m.scrollOffset
	endIdx := min(len(m.commandOutput), startIdx+availableLines)

	linesRendered := 0
	for i := startIdx; i < endIdx && linesRendered < availableLines; i++ {
		line := m.commandOutput[i]

		if lipgloss.Width(line) > maxWidth {
			line = truncateToWidth(line, maxWidth)
		}

		if strings.HasPrefix(line, "--- Benchmark Run") {
			s.WriteString(separatorStyle.Render(line))
		} else if strings.Contains(line, "Match found!") {
			s.WriteString(matchStyle.Render(line))
		} else {
			s.WriteString(regularStyle.Render(line))
		}
		s.WriteString("\n")
		linesRendered++
	}

	for linesRendered < availableLines {
		s.WriteString("\n")
		linesRendered++
	}

	return s.String()
}
