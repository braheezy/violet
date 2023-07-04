package app

import (
	"fmt"

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
				MarginLeft(4).
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
		content := fmt.Sprintf("%v", command)
		longestCommand = max(longestCommand, len(content))
		buttons = append(buttons, button{
			content: content,
			style:   defaultButtonStyle,
		})
	}
	// Set the button width based on longest command
	defaultButtonStyle.Width(longestCommand + defaultButtonStyle.GetHorizontalFrameSize())
	activeButtonStyle.Width(longestCommand + activeButtonStyle.GetHorizontalFrameSize())

	return buttonGroup{
		buttons: buttons,
		// This provides excellent space for each command
		width: longestCommand * 2,
	}
}

func (b *buttonGroup) View(selectedCommand int) string {
	for i := range b.buttons {
		if i == selectedCommand {
			b.buttons[i].style = activeButtonStyle
		} else {
			b.buttons[i].style = defaultButtonStyle
		}
	}

	leftSide := lipgloss.JoinVertical(lipgloss.Center, b.buttons[0].View(), b.buttons[2].View())
	rightSide := lipgloss.JoinVertical(lipgloss.Center, b.buttons[1].View(), b.buttons[3].View())

	group := lipgloss.JoinHorizontal(lipgloss.Center, leftSide, rightSide)

	return buttonGroupStyle.Render(group)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
