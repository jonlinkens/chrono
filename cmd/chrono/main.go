package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"chrono/internal/benchmark"
	"chrono/internal/output"
	"chrono/internal/shellcalibration"
	"chrono/internal/tui"
)

func main() {
	config := parseFlags()

	if !config.UseCli {
		if err := tui.Run(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	var shellOverhead time.Duration
	if !config.SkipCalibration {
		output.PrintCalibration(config.CalibrationRuns)

		shellOverhead = shellcalibration.CalibrateShellOverhead(config.CalibrationRuns)

		output.PrintShellOverhead(shellOverhead)
	}

	if config.Warmups > 0 {
		output.PrintWarmupHeader(config.Warmups)
		for i := range config.Warmups {
			result := benchmark.Run(config, shellOverhead)
			output.PrintWarmupResult(i+1, result)
		}
		fmt.Println()
	}

	output.PrintBenchmarkHeader(config.Runs)
	results := make([]benchmark.Result, 0, config.Runs)

	for i := range config.Runs {
		result := benchmark.Run(config, shellOverhead)
		results = append(results, result)
		output.PrintBenchmarkResult(i+1, result)
	}

	output.PrintSummary(results, config, shellOverhead)
}

func parseCommandString(cmd string) ([]string, error) {
	var args []string
	var current strings.Builder
	var inQuotes bool
	var quoteChar rune

	runes := []rune(cmd)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		switch {
		case r == '\\' && i+1 < len(runes):

			i++
			current.WriteRune(runes[i])

		case !inQuotes && (r == '"' || r == '\''):
			inQuotes = true
			quoteChar = r

		case inQuotes && r == quoteChar:
			inQuotes = false

		case !inQuotes && unicode.IsSpace(r):
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}

		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	if inQuotes {
		return nil, fmt.Errorf("unclosed quote in command string")
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("empty command string")
	}

	return args, nil
}

func parseFlags() benchmark.Config {
	var (
		phrase          = flag.String("phrase", "", "Phrase to search for in command output (if not specified, measures until command completion)")
		warmups         = flag.Int("warmups", 0, "Number of warmup runs before benchmarking")
		runs            = flag.Int("runs", 1, "Number of benchmark runs")
		timeout         = flag.Duration("timeout", 0, "Maximum time to wait for phrase or command completion (default: no timeout)")
		calibrationRuns = flag.Int("calibration", 5, "Number of calibration runs to measure shell startup overhead")
		skipCalibration = flag.Bool("skip-calibration", false, "Skip calibration and don't subtract shell overhead")
		useCLI          = flag.Bool("cli", false, "Use CLI output instead of terminal UI")
		commandStr      = flag.String("command", "", "Command to benchmark as a quoted string (alternative to positional arguments)")
	)
	flag.Parse()

	var command []string
	if *commandStr != "" {
		var err error
		command, err = parseCommandString(*commandStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing command string: %v\n", err)
			os.Exit(1)
		}
		if flag.NArg() > 0 {
			fmt.Fprintf(os.Stderr, "Error: cannot specify both --command and positional arguments\n")
			flag.Usage()
			os.Exit(1)
		}
	} else {
		if flag.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "Error: command to benchmark is required\n")
			fmt.Fprintf(os.Stderr, "Use either:\n")
			fmt.Fprintf(os.Stderr, "  --command \"cmd with args\"  (command as quoted string)\n")
			fmt.Fprintf(os.Stderr, "  -- cmd with args           (command after -- separator)\n")
			fmt.Fprintf(os.Stderr, "  cmd with args              (command as positional arguments)\n")
			flag.Usage()
			os.Exit(1)
		}
		command = flag.Args()
	}

	config := benchmark.Config{
		Phrase:          *phrase,
		Warmups:         *warmups,
		Runs:            *runs,
		Timeout:         *timeout,
		CalibrationRuns: *calibrationRuns,
		SkipCalibration: *skipCalibration,
		Command:         command,
		UseCli:          *useCLI,
	}

	return config
}
