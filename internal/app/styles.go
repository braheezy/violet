package app

import (
	"github.com/charmbracelet/lipgloss"
	tint "github.com/lrstanley/bubbletint"
)

// Tweak these to quickly change app feels
var defaultDarkTheme = tint.TintKonsolas
var defaultLightTheme = tint.TintMaterial
var theme = defaultDarkTheme

// Tweak these for a different palette
var (
	primaryColor   = theme.Purple()
	secondaryColor = theme.BrightBlue()
	accentColor    = theme.Yellow()
	textColor      = theme.Fg()
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Italic(true).
			Foreground(primaryColor)

	greeterStyle = lipgloss.NewStyle().
			Foreground(accentColor)
)

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	gapBorder         = lipgloss.Border{
		Top:         " ",
		Bottom:      "─",
		Left:        " ",
		Right:       " ",
		TopLeft:     " ",
		TopRight:    " ",
		BottomLeft:  "─",
		BottomRight: "╮",
	}

	highlightColor   = primaryColor
	inactiveTabStyle = lipgloss.NewStyle().
				Border(inactiveTabBorder, true).
				BorderForeground(highlightColor).
				Padding(0, 1)
	activeTabStyle = inactiveTabStyle.Copy().
			Border(activeTabBorder, true)
	tabGapStyle = inactiveTabStyle.Copy().
			Border(gapBorder)
	tabWindowStyle = lipgloss.NewStyle().
			BorderForeground(highlightColor).
			Padding(1, 1).
			Border(lipgloss.RoundedBorder()).
			UnsetBorderTop()
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
