package colours

import "github.com/charmbracelet/lipgloss"

const (
	Surface0 = "#313244"
	Surface1 = "#585b70"
	Surface2 = "#6c7086"

	Text     = "#cdd6f4"
	Subtext0 = "#a6adc8"

	Lavender = "#cba6f7"
	Blue     = "#89b4fa"
	Green    = "#a6e3a1"
	Red      = "#f38ba8"
	Yellow   = "#f9e2af"
	Cyan     = "#94e2d5"
)

var (
	RedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(Red))
	GreenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(Green))
	YellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(Yellow))
	BlueStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(Blue))
	PurpleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(Lavender))
	CyanStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(Cyan))
	GrayStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(Surface2))
	BoldStyle   = lipgloss.NewStyle().Bold(true)
)
