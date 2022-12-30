package app

import (
	"github.com/charmbracelet/lipgloss"
	tint "github.com/lrstanley/bubbletint"
)

var defaultDarkTheme = tint.TintMaterialDark
var defaultLightTheme = tint.TintMaterial
var theme = defaultDarkTheme

var (
	focusedStyle = lipgloss.NewStyle().
			Bold(true).
			Italic(true)

	commandSelectStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(theme.BrightBlue())

	outputHeight    = 10
	outputWidth     = 70
	outputViewStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true).
			BorderForeground(theme.Fg()).
			MarginLeft(2)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Italic(true).
			Foreground(theme.BrightPurple())
)
