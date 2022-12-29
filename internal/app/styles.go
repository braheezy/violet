package app

import "github.com/charmbracelet/lipgloss"

var (
	focusedStyle = lipgloss.NewStyle().
			Bold(true).
			Italic(true)

	commandSelectStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("69"))
)
