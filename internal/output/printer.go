package output

import (
	"fmt"
	"time"

	"chrono/internal/benchmark"
	"chrono/internal/colours"
	"chrono/internal/stats"
	"github.com/charmbracelet/lipgloss"
)

func FormatDuration(d time.Duration) string {
	return fmt.Sprintf("%.3fs", d.Seconds())
}

func PrintCalibration(runs int) {
	fmt.Printf("%s\n", colours.PurpleStyle.Render(fmt.Sprintf("Running %d calibration runs to measure shell startup overhead...", runs)))
}

func PrintShellOverhead(overhead time.Duration) {
	fmt.Printf("%s %s\n\n",
		colours.PurpleStyle.Render("Shell overhead:"),
		colours.BoldStyle.Render(FormatDuration(overhead)))
}

func PrintWarmupHeader(warmups int) {
	fmt.Printf("%s\n", colours.YellowStyle.Render(fmt.Sprintf("Running %d warmup runs...", warmups)))
}

func PrintWarmupResult(run int, result benchmark.Result) {
	if !result.Found {
		fmt.Printf("%s\n", colours.GrayStyle.Render(fmt.Sprintf("Warmup %d: phrase not found within timeout", run)))
	} else {
		fmt.Printf("%s\n", colours.GrayStyle.Render(fmt.Sprintf("Warmup %d: %s", run, FormatDuration(result.Duration))))
	}
}

func PrintBenchmarkHeader(runs int) {
	fmt.Printf("%s\n", colours.BlueStyle.Render(fmt.Sprintf("Running %d benchmark runs...", runs)))
}

func PrintBenchmarkResult(run int, result benchmark.Result) {
	headerStyle := colours.BoldStyle.Foreground(lipgloss.Color(colours.Cyan))
	fmt.Printf("\n%s\n", headerStyle.Render(fmt.Sprintf("--- Benchmark Run %d ---", run)))

	if !result.Found {
		fmt.Printf("%s\n", colours.RedStyle.Render(fmt.Sprintf("Run %d: phrase not found within timeout", run)))
	} else {
		fmt.Printf("%s\n", colours.GreenStyle.Render(fmt.Sprintf("Run %d: %s", run, colours.BoldStyle.Render(FormatDuration(result.Duration)))))
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
			colours.CyanStyle.Render("Mode:"),
			fmt.Sprintf(" Command completion timing%s%s", warmupInfo, calibrationInfo))
	} else {
		fmt.Printf("%s%s\n",
			colours.CyanStyle.Render("Phrase:"),
			fmt.Sprintf(" \"%s\"%s%s", colours.BoldStyle.Render(config.Phrase), warmupInfo, calibrationInfo))
	}

	if len(validResults) == 0 {
		if config.Phrase == "" {
			fmt.Printf("%s\n", colours.RedStyle.Render("No successful runs - all commands timed out"))
		} else {
			fmt.Printf("%s\n", colours.RedStyle.Render("No successful runs - phrase was not found in any execution"))
		}
		return
	}

	if failedCount > 0 {
		fmt.Printf("%s  ", colours.RedStyle.Render(fmt.Sprintf("Failed: %d/%d", failedCount, len(results))))
	}

	if len(validResults) == 1 {
		fmt.Printf("%s %s\n",
			colours.CyanStyle.Render("Time:"),
			colours.BoldStyle.Render(FormatDuration(validResults[0])))
		return
	}

	stats := stats.CalculateStatistics(validResults)
	if failedCount == 0 {
		fmt.Printf("%s  ", colours.GreenStyle.Render(fmt.Sprintf("Runs: %d", len(validResults))))
	}
	fmt.Printf("%s %s  %s %s  %s %s  %s %s\n",
		colours.CyanStyle.Render("Mean:"), colours.BoldStyle.Render(FormatDuration(stats.Mean)),
		colours.GreenStyle.Render("Min:"), FormatDuration(stats.Min),
		colours.RedStyle.Render("Max:"), FormatDuration(stats.Max),
		colours.YellowStyle.Render("Range:"), FormatDuration(stats.Range))
}
