package app

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	defaultLargeButtonStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Padding(1)
	activeLargeButtonStyle = defaultLargeButtonStyle.
				Foreground(secondaryColor).
				Bold(true)
	buttonLargeGroupStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder(), true).
				BorderForeground(primaryColor).
				Margin(0)

	defaultSmallButtonStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Margin(0, 1)
	activeSmallButtonStyle = defaultSmallButtonStyle.
				Foreground(secondaryColor).
				Bold(true)
	buttonSmallGroupStyle = lipgloss.NewStyle().
				Margin(marginVertical, marginHorizontal, 0).
				Border(lipgloss.RoundedBorder(), true).
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
			bg.buttons[i].style = activeSmallButtonStyle.Padding(0)
		} else {
			bg.buttons[i].style = defaultSmallButtonStyle.Padding(0)
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
	// Create buttons based on supported commands
	for _, command := range supportedVagrantCommands {
		cont := symbols[command] + " " + command
		buttons = append(buttons, button{
			content: cont,
			style:   defaultLargeButtonStyle,
		})
	}

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
