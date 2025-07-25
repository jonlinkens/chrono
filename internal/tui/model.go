package tui

import (
	"time"

	"chrono/internal/benchmark"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	StateCalibrating int = iota
	StateWarmup
	StateBenchmarking
	StateCompleted
)

type Model struct {
	config        benchmark.Config
	state         int
	shellOverhead time.Duration

	warmupProgress    int
	benchmarkProgress int

	warmupResults    []benchmark.Result
	benchmarkResults []benchmark.Result

	currentRun int
	totalRuns  int

	currentRunStartTime time.Time
	elapsedTime         time.Duration
	isRunning           bool

	commandOutput []string
	scrollOffset  int

	width  int
	height int

	clipboardFeedback     string
	clipboardFeedbackTime time.Time

	err error
}

func NewModel(config benchmark.Config) Model {
	totalRuns := config.Warmups + config.Runs
	return Model{
		config:           config,
		state:            StateCalibrating,
		warmupResults:    make([]benchmark.Result, 0, config.Warmups),
		benchmarkResults: make([]benchmark.Result, 0, config.Runs),
		totalRuns:        totalRuns,
		commandOutput:    make([]string, 0),
		width:            DefaultWidth,
		height:           DefaultHeight,
	}
}

func (m Model) Init() tea.Cmd {
	if m.config.SkipCalibration {
		if m.config.Warmups > 0 {
			return tea.Batch(
				m.startWarmup(),
				m.tickCmd(),
			)
		}
		return tea.Batch(
			m.startBenchmark(),
			m.tickCmd(),
		)
	}
	return tea.Batch(
		m.runCalibration(),
		m.tickCmd(),
	)
}

func Run(config benchmark.Config) error {
	model := NewModel(config)

	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())

	_, err := p.Run()
	return err
}
