package output

import (
	"fmt"
	"time"

	"chrono/internal/benchmark"
	"chrono/internal/stats"
	"github.com/charmbracelet/lipgloss"
)

var (
	redStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8"))
	greenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1"))
	yellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f9e2af"))
	blueStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#89b4fa"))
	purpleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cba6f7"))
	cyanStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#94e2d5"))
	grayStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086"))
	boldStyle   = lipgloss.NewStyle().Bold(true)
)

func FormatDuration(d time.Duration) string {
	return fmt.Sprintf("%.3fs", d.Seconds())
}

func PrintCalibration(runs int) {
	fmt.Printf("%s\n", purpleStyle.Render(fmt.Sprintf("Running %d calibration runs to measure shell startup overhead...", runs)))
}

func PrintShellOverhead(overhead time.Duration) {
	fmt.Printf("%s %s\n\n",
		purpleStyle.Render("Shell overhead:"),
		boldStyle.Render(FormatDuration(overhead)))
}

func PrintWarmupHeader(warmups int) {
	fmt.Printf("%s\n", yellowStyle.Render(fmt.Sprintf("Running %d warmup runs...", warmups)))
}

func PrintWarmupResult(run int, result benchmark.Result) {
	if !result.Found {
		fmt.Printf("%s\n", grayStyle.Render(fmt.Sprintf("Warmup %d: phrase not found within timeout", run)))
	} else {
		fmt.Printf("%s\n", grayStyle.Render(fmt.Sprintf("Warmup %d: %s", run, FormatDuration(result.Duration))))
	}
}

func PrintBenchmarkHeader(runs int) {
	fmt.Printf("%s\n", blueStyle.Render(fmt.Sprintf("Running %d benchmark runs...", runs)))
}

func PrintBenchmarkResult(run int, result benchmark.Result) {
	headerStyle := boldStyle.Foreground(lipgloss.Color("#94e2d5"))
	fmt.Printf("\n%s\n", headerStyle.Render(fmt.Sprintf("--- Benchmark Run %d ---", run)))

	if !result.Found {
		fmt.Printf("%s\n", redStyle.Render(fmt.Sprintf("Run %d: phrase not found within timeout", run)))
	} else {
		fmt.Printf("%s\n", greenStyle.Render(fmt.Sprintf("Run %d: %s", run, boldStyle.Render(FormatDuration(result.Duration)))))
	}
}
func PrintSummary(results []benchmark.Result, config benchmark.Config, shellOverhead time.Duration) {
	fmt.Println()

	validResults := make([]time.Duration, 0, len(results))
	failedCount := 0

	for _, result := range results {
		if result.Found {
			validResults = append(validResults, result.Duration)
		} else {
			failedCount++
		}
	}

	warmupInfo := ""
	if config.Warmups > 0 {
		warmupInfo = fmt.Sprintf(" (%d warmups)", config.Warmups)
	}

	calibrationInfo := ""
	if !config.SkipCalibration {
		calibrationInfo = fmt.Sprintf(" (-%s shell overhead)", FormatDuration(shellOverhead))
	}

	if config.Phrase == "" {
		fmt.Printf("%s%s\n",
			cyanStyle.Render("Mode:"),
			fmt.Sprintf(" Command completion timing%s%s", warmupInfo, calibrationInfo))
	} else {
		fmt.Printf("%s%s\n",
			cyanStyle.Render("Phrase:"),
			fmt.Sprintf(" \"%s\"%s%s", boldStyle.Render(config.Phrase), warmupInfo, calibrationInfo))
	}

	if len(validResults) == 0 {
		if config.Phrase == "" {
			fmt.Printf("%s\n", redStyle.Render("No successful runs - all commands timed out"))
		} else {
			fmt.Printf("%s\n", redStyle.Render("No successful runs - phrase was not found in any execution"))
		}
		return
	}

	if failedCount > 0 {
		fmt.Printf("%s  ", redStyle.Render(fmt.Sprintf("Failed: %d/%d", failedCount, len(results))))
	}

	if len(validResults) == 1 {
		fmt.Printf("%s %s\n",
			cyanStyle.Render("Time:"),
			boldStyle.Render(FormatDuration(validResults[0])))
		return
	}

	stats := stats.CalculateStatistics(validResults)
	if failedCount == 0 {
		fmt.Printf("%s  ", greenStyle.Render(fmt.Sprintf("Runs: %d", len(validResults))))
	}
	fmt.Printf("%s %s  %s %s  %s %s  %s %s\n",
		cyanStyle.Render("Mean:"), boldStyle.Render(FormatDuration(stats.Mean)),
		greenStyle.Render("Min:"), FormatDuration(stats.Min),
		redStyle.Render("Max:"), FormatDuration(stats.Max),
		yellowStyle.Render("Range:"), FormatDuration(stats.Range))
}
