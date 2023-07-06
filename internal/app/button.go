package app

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	defaultButtonStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Background(theme.Bg()).
				Padding(1)
	activeButtonStyle = defaultButtonStyle.Copy().
				Foreground(textColor).
				Background(primaryColor).
				Bold(true)
	buttonGroupStyle = lipgloss.NewStyle().
				Padding(1)
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
	width   int
}

func newCommandButtons() buttonGroup {
	var buttons []button
	var longestCommand int
	// Create buttons based on supported commands
	for _, command := range supportedVagrantCommands {
		longestCommand = max(longestCommand, len(command))
		buttons = append(buttons, button{
			content: command,
			style:   defaultButtonStyle,
		})
	}
	// Set the button width based on longest command
	defaultButtonStyle.Width(longestCommand + defaultButtonStyle.GetHorizontalFrameSize())
	activeButtonStyle.Width(longestCommand + activeButtonStyle.GetHorizontalFrameSize())

	return buttonGroup{
		buttons: buttons,
		// This provides excellent space for each command
		width: longestCommand * 3,
	}
}

func (bg *buttonGroup) View(selectedCommand int) string {
	for i := range bg.buttons {
		if i == selectedCommand {
			bg.buttons[i].style = activeButtonStyle
		} else {
			bg.buttons[i].style = defaultButtonStyle
		}
	}

	// TODO: Hacky to hardcode the row items. Is there a better way?
	topRow := lipgloss.JoinHorizontal(lipgloss.Center, bg.buttons[0].View(), bg.buttons[1].View(), bg.buttons[2].View())
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Center, bg.buttons[3].View(), bg.buttons[4].View())

	grid := lipgloss.JoinVertical(lipgloss.Center, topRow, bottomRow)

	return buttonGroupStyle.Render(grid)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
