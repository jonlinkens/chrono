package shellcalibration

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const fallbackShell = "/bin/sh"

func CalibrateShellOverhead(runs int) time.Duration {
	durations := make([]time.Duration, 0, runs)

	calibrationShell := fallbackShell
	calibrationArgs := []string{"-c", "true"}

	if userShell := os.Getenv("SHELL"); userShell != "" {
		calibrationShell = userShell
	}

	warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f9e2af"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8"))

	for i := range runs {
		start := time.Now()
		cmd := exec.Command(calibrationShell, calibrationArgs...)
		err := cmd.Run()
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("%s\n", warningStyle.Render(fmt.Sprintf("Calibration run %d failed with %s: %v", i+1, calibrationShell, err)))

			if calibrationShell != fallbackShell {
				start = time.Now()
				cmd = exec.Command(fallbackShell, "-c", "true")
				err = cmd.Run()
				duration = time.Since(start)
				if err != nil {
					fmt.Printf("%s\n", errorStyle.Render(fmt.Sprintf("Calibration run %d failed completely: %v", i+1, err)))
					continue
				}
			} else {
				continue
			}
		}

		durations = append(durations, duration)
	}

	if len(durations) == 0 {
		fmt.Printf("%s\n", warningStyle.Render("All calibration runs failed, using 0 overhead"))
		return 0
	}

	var total time.Duration
	for _, d := range durations {
		total += d
	}

	return total / time.Duration(len(durations))
}
