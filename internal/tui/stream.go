package tui

import (
	"chrono/internal/benchmark"
	"io"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type startStreamingMsg struct {
	cmd         *exec.Cmd
	stdout      io.ReadCloser
	stderr      io.ReadCloser
	startTime   time.Time
	isWarmup    bool
	phrase      string
	done        <-chan time.Duration
	cancel      chan struct{}
	outputLines <-chan string
}

type streamNextMsg struct {
	cmd                *exec.Cmd
	startTime          time.Time
	isWarmup           bool
	phrase             string
	done               <-chan time.Duration
	cancel             chan struct{}
	outputLines        <-chan string
	matchFoundReceived bool
	phraseMatchResult  *benchmark.Result
	commandCompleted   bool
	completionResult   *benchmark.Result
}

type newOutputLineMsg struct {
	line       string
	streamNext streamNextMsg
}

func (m Model) startStreaming(msg startStreamingMsg) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		return streamNextMsg{
			cmd:                msg.cmd,
			startTime:          msg.startTime,
			isWarmup:           msg.isWarmup,
			phrase:             msg.phrase,
			done:               msg.done,
			cancel:             msg.cancel,
			outputLines:        msg.outputLines,
			matchFoundReceived: false,
			phraseMatchResult:  nil,
			commandCompleted:   false,
			completionResult:   nil,
		}
	})
}

func (m Model) handleStreamNext(msg streamNextMsg) tea.Cmd {
	return func() tea.Msg {
		if m.config.Timeout > 0 {
			select {
			case duration := <-msg.done:
				if msg.cancel != nil {
					result := m.createAdjustedResult(duration, true)

					updatedMsg := msg
					updatedMsg.phraseMatchResult = &result

					return tea.Tick(StreamTickInterval, func(t time.Time) tea.Msg {
						return updatedMsg
					})()
				} else {
					result := m.createAdjustedResult(duration, true)

					updatedMsg := msg
					updatedMsg.commandCompleted = true
					updatedMsg.completionResult = &result

					return tea.Tick(StreamTickInterval, func(t time.Time) tea.Msg {
						return updatedMsg
					})()
				}

			case <-time.After(m.config.Timeout):
				if msg.cancel != nil {
					close(msg.cancel)
				}
				if msg.cmd.Process != nil {
					msg.cmd.Process.Kill()
				}

				result := m.createTimeoutResult()

				return runCompleteMsg{
					result:   result,
					isWarmup: msg.isWarmup,
					output:   []string{},
				}

			case line, ok := <-msg.outputLines:
				if ok {
					if line == "Match found!" && msg.phraseMatchResult != nil {
						if msg.cancel != nil {
							close(msg.cancel)
							if msg.cmd.Process != nil {
								msg.cmd.Process.Kill()
							}
						}

						updatedMsg := msg
						updatedMsg.matchFoundReceived = true
						return newOutputLineMsg{
							line:       line,
							streamNext: updatedMsg,
						}
					}

					return newOutputLineMsg{
						line:       line,
						streamNext: msg,
					}
				} else {
					if msg.commandCompleted && msg.completionResult != nil {
						return runCompleteMsg{
							result:   *msg.completionResult,
							isWarmup: msg.isWarmup,
							output:   []string{},
						}
					}

					if msg.phrase != "" {
						result := benchmark.Result{
							Duration: 0,
							Found:    false,
						}

						return runCompleteMsg{
							result:   result,
							isWarmup: msg.isWarmup,
							output:   []string{},
						}
					}

					return tea.Tick(StreamTickInterval, func(t time.Time) tea.Msg {
						return msg
					})()
				}

			default:
				if msg.phraseMatchResult != nil && msg.matchFoundReceived {
					return runCompleteMsg{
						result:   *msg.phraseMatchResult,
						isWarmup: msg.isWarmup,
						output:   []string{},
					}
				}

				return tea.Tick(StreamTickInterval, func(t time.Time) tea.Msg {
					return msg
				})()
			}
		} else {
			select {
			case duration := <-msg.done:
				if msg.cancel != nil {
					result := m.createAdjustedResult(duration, true)

					updatedMsg := msg
					updatedMsg.phraseMatchResult = &result

					return tea.Tick(StreamTickInterval, func(t time.Time) tea.Msg {
						return updatedMsg
					})()
				} else {
					result := m.createAdjustedResult(duration, true)

					updatedMsg := msg
					updatedMsg.commandCompleted = true
					updatedMsg.completionResult = &result

					return tea.Tick(StreamTickInterval, func(t time.Time) tea.Msg {
						return updatedMsg
					})()
				}

			case line, ok := <-msg.outputLines:
				if ok {
					if line == "Match found!" && msg.phraseMatchResult != nil {
						if msg.cancel != nil {
							close(msg.cancel)
							if msg.cmd.Process != nil {
								msg.cmd.Process.Kill()
							}
						}

						updatedMsg := msg
						updatedMsg.matchFoundReceived = true
						return newOutputLineMsg{
							line:       line,
							streamNext: updatedMsg,
						}
					}

					return newOutputLineMsg{
						line:       line,
						streamNext: msg,
					}
				} else {
					if msg.commandCompleted && msg.completionResult != nil {
						return runCompleteMsg{
							result:   *msg.completionResult,
							isWarmup: msg.isWarmup,
							output:   []string{},
						}
					}

					if msg.phrase != "" {
						result := benchmark.Result{
							Duration: 0,
							Found:    false,
						}

						return runCompleteMsg{
							result:   result,
							isWarmup: msg.isWarmup,
							output:   []string{},
						}
					}

					return tea.Tick(StreamTickInterval, func(t time.Time) tea.Msg {
						return msg
					})()
				}

			default:
				if msg.phraseMatchResult != nil && msg.matchFoundReceived {
					return runCompleteMsg{
						result:   *msg.phraseMatchResult,
						isWarmup: msg.isWarmup,
						output:   []string{},
					}
				}

				return tea.Tick(StreamTickInterval, func(t time.Time) tea.Msg {
					return msg
				})()
			}
		}
	}
}
