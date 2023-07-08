package app

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	defaultLargeButtonStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Background(theme.Bg()).
				Padding(1)
	activeLargeButtonStyle = defaultLargeButtonStyle.Copy().
				Foreground(textColor).
				Background(primaryColor).
				Bold(true)
	buttonLargeGroupStyle = lipgloss.NewStyle().
				Margin(0, 1)

	defaultSmallButtonStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Margin(0, 1)
	activeSmallButtonStyle = defaultSmallButtonStyle.Copy().
				Foreground(textColor).
				Bold(true)
	buttonSmallGroupStyle = lipgloss.NewStyle().
				Margin(marginVertical, marginHorizontal).
				Border(lipgloss.NormalBorder(), true).
				BorderForeground(primaryColor)
)

type button struct {
	content string
	style   lipgloss.Style
}

func (b *button) View() string {
	return b.style.Render(b.content)
}

type buttonGroup struct {
	buttons []button
}

type machineCommandButtons buttonGroup

func newMachineCommandButtons(supportedVagrantCommands []string) machineCommandButtons {
	var buttons []button

	for _, command := range supportedVagrantCommands {
		buttons = append(buttons, button{
			content: symbols[command],
			style:   defaultSmallButtonStyle,
		})
	}

	return machineCommandButtons{
		buttons: buttons,
	}
}

func (bg *machineCommandButtons) View(selectedCommand int, hasFocus bool) string {
	for i := range bg.buttons {
		if i == selectedCommand && hasFocus {
			bg.buttons[i].style = activeSmallButtonStyle.Copy().Padding(0)
		} else {
			bg.buttons[i].style = defaultSmallButtonStyle.Copy().Padding(0)
		}
	}

	var row []string
	for _, button := range bg.buttons {
		row = append(row, button.View())
	}

	grid := lipgloss.JoinHorizontal(lipgloss.Center, row...)

	return buttonSmallGroupStyle.Render(grid)
}

type envCommandButtons buttonGroup

func newEnvCommandButtons(supportedVagrantCommands []string) envCommandButtons {
	var buttons []button
	var longestContent int
	// Create buttons based on supported commands
	for _, command := range supportedVagrantCommands {
		cont := symbols[command] + " " + command
		longestContent = max(longestContent, len(cont))
		buttons = append(buttons, button{
			content: cont,
			style:   defaultLargeButtonStyle,
		})
	}
	// Set the button width based on longest command
	defaultLargeButtonStyle.Width(longestContent)
	activeLargeButtonStyle.Width(longestContent)

	return envCommandButtons{
		buttons: buttons,
	}
}

func (bg *envCommandButtons) View(selectedCommand int, hasFocus bool) string {
	for i := range bg.buttons {
		if i == selectedCommand && hasFocus {
			bg.buttons[i].style = activeLargeButtonStyle
		} else {
			bg.buttons[i].style = defaultLargeButtonStyle
		}
	}

	var row []string
	for _, button := range bg.buttons {
		row = append(row, button.View())
	}

	grid := lipgloss.JoinHorizontal(lipgloss.Center, row...)

	return buttonLargeGroupStyle.Render(grid)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
