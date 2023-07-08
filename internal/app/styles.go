package app

import (
	"github.com/charmbracelet/lipgloss"
	tint "github.com/lrstanley/bubbletint"
)

// Tweak these to quickly change app feels
var defaultDarkTheme = tint.TintKonsolas
var defaultLightTheme = tint.TintCatppuccinLatte
var theme = defaultDarkTheme

// Tweak these for a different palette
var (
	primaryColor   = theme.Purple()
	secondaryColor = theme.Yellow()
	accentColor    = theme.Cyan()
	textColor      = theme.Fg()

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

	marginVertical   = 1
	marginHorizontal = 2

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Italic(true).
			Foreground(primaryColor)
	greeterStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	highlightColor   = primaryColor
	inactiveTabStyle = lipgloss.NewStyle().
				Border(inactiveTabBorder, true).
				BorderForeground(highlightColor).
				Padding(0, 1).
				Foreground(textColor)
	activeTabStyle = inactiveTabStyle.Copy().
			Border(activeTabBorder, true).
			Foreground(accentColor)
	tabGapStyle = inactiveTabStyle.Copy().
			Border(gapBorder)
	tabWindowStyle = lipgloss.NewStyle().
			BorderForeground(highlightColor).
			Padding(1, 1).
			Border(lipgloss.RoundedBorder()).
			UnsetBorderTop()
	errorTitleStyle = lipgloss.NewStyle().
			Foreground(theme.Red()).
			MarginLeft(marginHorizontal)
	errorStyle = lipgloss.NewStyle().
			MarginLeft(marginHorizontal).
			Foreground(theme.BrightRed()).
			Bold(true)
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

// ------ Card Style ---------
var (
	cardTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Margin(0, marginHorizontal)
	cardStatusStyle = lipgloss.NewStyle().
			MarginLeft(marginHorizontal * 2)
	statusColors = map[string]lipgloss.TerminalColor{
		"running":     theme.Green(),
		"shutoff":     theme.Red(),
		"stopped":     theme.Red(),
		"not started": theme.Black(),
	}
	cardProviderStyle = lipgloss.NewStyle().
				Faint(true).
				Italic(true).
				MarginLeft(marginHorizontal * 2).
				Foreground(textColor)
	defaultCardStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder(), false, false, false, false).
				MarginLeft(1)
	selectedCardStyle = defaultCardStyle.Copy().
				BorderLeft(true).
				Border(lipgloss.RoundedBorder(), false, false, false, true).
				BorderForeground(accentColor).
				MarginLeft(0)

	envCardTitleStyle = cardTitleStyle.Copy().
				MarginLeft(1)
	selectedEnvCardStyle = envCardTitleStyle.Copy().
				Bold(true).
				Foreground(accentColor)
)
