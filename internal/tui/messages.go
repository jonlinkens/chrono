package tui

import (
	"chrono/internal/benchmark"
	"time"
)

type calibrationCompleteMsg struct {
	overhead time.Duration
}

type runStartMsg struct {
	isWarmup bool
}

type runCompleteMsg struct {
	result   benchmark.Result
	isWarmup bool
	output   []string
}

type outputLineMsg struct {
	line string
}

type tickMsg time.Time

type errorMsg struct {
	err error
}

type clipboardCopiedMsg struct{}

type clearClipboardFeedbackMsg struct{}
