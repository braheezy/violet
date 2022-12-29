package app

import "github.com/charmbracelet/lipgloss"

var (
	focusedStyle = lipgloss.NewStyle().
			Bold(true).
			Italic(true)

	commandSelectStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("69"))
	outputHeight    = 60
	outputWidth     = 10
	outputViewStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true).
			BorderForeground(lipgloss.Color("#cba6f7")).
			MarginLeft(4)
)
