package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "y":
			if m.state == StateCompleted {
				return m, m.copyResultsToClipboard()
			}
			return m, nil
		case "up", "k":
			if m.scrollOffset > 0 {
				m.scrollOffset--
			}
			return m, nil
		case "down", "j":
			maxScroll := m.getMaxScrollOffset()
			if m.scrollOffset < maxScroll {
				m.scrollOffset++
			}
			return m, nil
		case "esc":
			m.scrollOffset = m.getMaxScrollOffset()
			return m, nil
		}

	case tea.MouseMsg:
		leftContentWidth := m.calculateBoundedLeftWidth()

		rightColumnStart := leftContentWidth

		if msg.X >= rightColumnStart {
			switch msg.Button {
			case tea.MouseButtonWheelUp:
				if m.scrollOffset > 0 {
					m.scrollOffset = max(0, m.scrollOffset-ScrollSteps)
				}
				return m, nil
			case tea.MouseButtonWheelDown:
				maxScroll := m.getMaxScrollOffset()
				if m.scrollOffset < maxScroll {
					m.scrollOffset = min(maxScroll, m.scrollOffset+ScrollSteps)
				}
				return m, nil
			}
		}

	case tickMsg:
		if m.isRunning {
			m.elapsedTime = time.Since(m.currentRunStartTime)
		}
		return m, m.tickCmd()

	case calibrationCompleteMsg:
		m.shellOverhead = msg.overhead
		if m.config.Warmups > 0 {
			m.state = StateWarmup
			return m, m.startWarmup()
		}
		m.state = StateBenchmarking
		return m, m.startBenchmark()

	case runStartMsg:
		m.currentRunStartTime = time.Now()
		m.isRunning = true
		m.elapsedTime = 0

		if !msg.isWarmup {
			runNumber := m.benchmarkProgress + 1

			maxScroll := m.getMaxScrollOffset()
			shouldAutoScroll := m.scrollOffset >= maxScroll-1

			m.commandOutput = append(m.commandOutput, "")
			separator := fmt.Sprintf("--- Benchmark Run %d ---", runNumber)
			m.commandOutput = append(m.commandOutput, separator)
			m.commandOutput = append(m.commandOutput, "")

			if shouldAutoScroll {
				m.scrollOffset = m.getMaxScrollOffset()
			}
		}

		return m, nil
	case startStreamingMsg:
		return m, m.startStreaming(msg)

	case streamNextMsg:
		return m, m.handleStreamNext(msg)

	case newOutputLineMsg:
		shouldAutoScroll := m.autoScrollToBottom()
		m.commandOutput = append(m.commandOutput, msg.line)
		if shouldAutoScroll {
			m.scrollOffset = m.getMaxScrollOffset()
		}

		if msg.streamNext.matchFoundReceived && msg.streamNext.phraseMatchResult != nil {
			return m, tea.Cmd(func() tea.Msg {
				return runCompleteMsg{
					result:   *msg.streamNext.phraseMatchResult,
					isWarmup: msg.streamNext.isWarmup,
					output:   []string{},
				}
			})
		}

		if msg.streamNext.commandCompleted && msg.streamNext.completionResult != nil && msg.streamNext.phrase == "" {
			return m, tea.Cmd(func() tea.Msg {
				select {
				case line, ok := <-msg.streamNext.outputLines:
					if ok {
						return newOutputLineMsg{
							line:       line,
							streamNext: msg.streamNext,
						}
					} else {
						return runCompleteMsg{
							result:   *msg.streamNext.completionResult,
							isWarmup: msg.streamNext.isWarmup,
							output:   []string{},
						}
					}
				default:
					return msg.streamNext
				}
			})
		}

		return m, m.handleStreamNext(msg.streamNext)

	case runCompleteMsg:
		m.isRunning = false

		if msg.isWarmup {
			m.warmupResults = append(m.warmupResults, msg.result)
			m.warmupProgress++
			m.currentRun++

			if m.warmupProgress < m.config.Warmups {
				return m, m.startWarmup()
			}

			m.state = StateBenchmarking
			return m, m.startBenchmark()
		} else {
			m.benchmarkResults = append(m.benchmarkResults, msg.result)
			m.benchmarkProgress++
			m.currentRun++

			if m.benchmarkProgress < m.config.Runs {
				return m, m.startBenchmark()
			}

			m.state = StateCompleted
			return m, nil
		}

	case outputLineMsg:
		shouldAutoScroll := m.autoScrollToBottom()
		m.commandOutput = append(m.commandOutput, msg.line)
		if shouldAutoScroll {
			m.scrollOffset = m.getMaxScrollOffset()
		}
		return m, nil

	case errorMsg:
		m.err = msg.err
		m.isRunning = false
		return m, nil

	case clipboardCopiedMsg:
		m.clipboardFeedback = "Results copied to clipboard!"
		m.clipboardFeedbackTime = time.Now()
		return m, tea.Tick(ClipboardFeedbackTime, func(t time.Time) tea.Msg {
			return clearClipboardFeedbackMsg{}
		})

	case clearClipboardFeedbackMsg:
		if time.Since(m.clipboardFeedbackTime) >= ClipboardFeedbackTime {
			m.clipboardFeedback = ""
		}
		return m, nil
	}

	return m, nil
}
