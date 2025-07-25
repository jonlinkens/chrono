package tui

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"chrono/internal/stats"

	"github.com/aymanbagabas/go-osc52/v2"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) copyResultsToClipboard() tea.Cmd {
	return func() tea.Msg {
		if m.state != StateCompleted {
			return nil
		}

		results := m.buildResultsString()

		sequence := osc52.New(results)

		if os.Getenv("TMUX") != "" {
			sequence = sequence.Tmux()
		} else if os.Getenv("TERM") == "screen" || strings.Contains(os.Getenv("TERM"), "screen") {
			sequence = sequence.Screen()
		}

		fmt.Fprint(os.Stderr, sequence)
		return clipboardCopiedMsg{}
	}
}

func (m Model) buildResultsString() string {

	var results strings.Builder

	cmd := m.formatCommandDisplay()
	results.WriteString(cmd)
	results.WriteString("\n")

	validResults, failedCount := m.filterValidResults()

	if len(validResults) == 0 {
		if m.config.Phrase == "" {
			results.WriteString("No successful runs - all executions timed out")
		} else {
			results.WriteString("No successful runs - phrase not found")
		}
	} else if len(validResults) == 1 {
		results.WriteString(fmt.Sprintf("Time: %s", formatDuration(validResults[0])))
		if failedCount > 0 {
			results.WriteString(fmt.Sprintf(" (%d failed)", failedCount))
		}
	} else {
		stats := stats.CalculateStatistics(validResults)

		var variance float64
		meanSeconds := stats.Mean.Seconds()
		for _, d := range validResults {
			diff := d.Seconds() - meanSeconds
			variance += diff * diff
		}
		variance /= float64(len(validResults))
		stdDev := time.Duration(math.Sqrt(variance) * float64(time.Second))

		results.WriteString(fmt.Sprintf("Mean: %s ± %s",
			formatDuration(stats.Mean), formatDuration(stdDev)))

		if failedCount > 0 {
			results.WriteString(fmt.Sprintf(" (%d/%d completed)", len(validResults), len(m.benchmarkResults)))
		}

		results.WriteString(fmt.Sprintf("\nRange: %s … %s",
			formatDuration(stats.Min), formatDuration(stats.Max)))
	}

	return results.String()
}
