package tui

import (
	"time"

	"chrono/internal/colours"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.width < MinTerminalWidth || m.height < MinTerminalHeight {
		return "Terminal too small (minimum 80x20)"
	}

	leftContentWidth := m.calculateBoundedLeftWidth()

	rightContentWidth := m.width - leftContentWidth

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colours.Surface1))

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colours.Lavender)).
		Bold(true).
		Padding(0, 1).
		Align(lipgloss.Left)

	outputTitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colours.Blue)).
		Bold(true).
		Padding(0, 1).
		Align(lipgloss.Center)

	contentStyle := lipgloss.NewStyle().
		Padding(1)
	contentHeight := m.height - 6

	leftTitle := titleStyle.
		Width(leftContentWidth - 4).
		Render("chrono")

	rightTitle := outputTitleStyle.
		Width(rightContentWidth - 4).
		Render("Output")

	leftContent := contentStyle.
		Width(leftContentWidth - 4).
		Height(contentHeight).
		Render(m.renderLeftColumnContentText())

	rightContent := contentStyle.
		Width(rightContentWidth - 4).
		Height(contentHeight - 1).
		Render(m.renderRightColumnContentText(rightContentWidth-4, contentHeight-1))

	leftPanel := lipgloss.JoinVertical(
		lipgloss.Top,
		leftTitle,
		leftContent,
	)

	rightPanel := lipgloss.JoinVertical(
		lipgloss.Top,
		rightTitle,
		rightContent,
	)

	leftBorderedPanel := borderStyle.
		Width(leftContentWidth - 2).
		Render(leftPanel)

	rightBorderedPanel := borderStyle.
		Width(rightContentWidth - 2).
		Render(rightPanel)

	columns := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftBorderedPanel,
		rightBorderedPanel,
	)

	shortcuts := m.renderShortcuts()
	bottomText := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colours.Surface2)).
		Width(m.width).
		Align(lipgloss.Center).
		Render(shortcuts)

	layout := lipgloss.JoinVertical(
		lipgloss.Top,
		columns,
		"",
		bottomText,
	)

	return layout
}

func (m Model) renderShortcuts() string {
	if m.clipboardFeedback != "" && time.Since(m.clipboardFeedbackTime) < ClipboardFeedbackTime {
		return m.clipboardFeedback
	}

	if m.state == StateCompleted {
		return "↑/↓ j/k: scroll • Esc: top/bottom • y: copy results • q/Ctrl+C: quit"
	} else {
		return "↑/↓ j/k: scroll • Esc: top/bottom • q/Ctrl+C: quit"
	}
}
