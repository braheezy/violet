package app

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	tint "github.com/lrstanley/bubbletint"
)

var defaultDarkTheme = tint.TintMaterialDark
var defaultLightTheme = tint.TintMaterial
var theme = defaultDarkTheme

var (
	primaryColor   = theme.BrightPurple()
	secondaryColor = theme.Purple()
	accentColor    = theme.BrightYellow()
	textColor      = theme.Fg()
)

var (
	outputHeight    = 10
	outputWidth     = 70
	outputViewStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true).
			BorderForeground(accentColor).
			MarginLeft(2)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Italic(true).
			Foreground(theme.BrightPurple())
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
			Padding(1, 2).
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

var (
	// Available spinners
	spinners = []spinner.Spinner{
		spinner.MiniDot,
		spinner.Dot,
		spinner.Line,
		spinner.Jump,
		spinner.Pulse,
		spinner.Points,
		spinner.Globe,
		spinner.Moon,
		spinner.Monkey,
	}

	spinnerStyle = lipgloss.NewStyle().Foreground(secondaryColor)
)
