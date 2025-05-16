package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	PrimaryColor   = lipgloss.Color("#00ff00")
	SecondaryColor = lipgloss.Color("#ff00ff")
	ErrorColor     = lipgloss.Color("#ff0000")
	InfoColor      = lipgloss.Color("#00ffff")

	// Base styles
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryColor).
		MarginBottom(1)

	MenuItem = lipgloss.NewStyle().
			PaddingLeft(4)

	SelectedMenuItem = MenuItem.Copy().
				Foreground(SecondaryColor).
				Bold(true)

	// Status styles
	Success = lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true)

	Error = lipgloss.NewStyle().
		Foreground(ErrorColor).
		Bold(true)

	Info = lipgloss.NewStyle().
		Foreground(InfoColor)

	// Layout styles
	Container = lipgloss.NewStyle().
			Padding(1, 2)

	Section = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(0, 1)
)
