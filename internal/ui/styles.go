package ui

import "charm.land/lipgloss/v2"

var (
	Purple = lipgloss.Color("#7D56F4")
	Green  = lipgloss.Color("#00FF00")
	Red    = lipgloss.Color("#FF0000")
	Gray   = lipgloss.Color("#444444")

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Purple).
			Padding(0, 1)

	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Purple).
			Padding(1, 2)

	ProfitStyle = lipgloss.NewStyle().Foreground(Green).Bold(true)
	LossStyle   = lipgloss.NewStyle().Foreground(Red).Bold(true)

	TableBorderStyle = lipgloss.NewStyle().Foreground(Gray)
)
