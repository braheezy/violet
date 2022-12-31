package app

import "github.com/charmbracelet/lipgloss"

var (
	cardTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)
	cardStatusStyle = lipgloss.NewStyle().
			PaddingLeft(1)
	statusColors = map[string]lipgloss.TerminalColor{
		"running":     theme.BrightGreen(),
		"shutoff":     theme.BrightRed(),
		"not started": theme.Black(),
	}
	cardProviderStyle = lipgloss.NewStyle().
				Faint(true).
				Italic(true).
				PaddingLeft(1)
	selectedCardStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(accentColor).
				PaddingLeft(1)
)
