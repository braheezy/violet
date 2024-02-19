package app

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

func (v Violet) View() (view string) {
	// Title view area
	title := titleStyle.Render("Violet:")
	greeter := greeterStyle.Render("Pretty manager for Vagrant")
	titleGreeter := lipgloss.NewStyle().Margin(marginVertical, marginHorizontal).Render(title + greeter)
	view += lipgloss.PlaceHorizontal(v.terminalWidth, lipgloss.Center, titleGreeter)

	help := v.help.View(v.keys)
	helpText := lipgloss.NewStyle().
		Margin(marginVertical, marginHorizontal).
		Render(help)
	view += lipgloss.PlaceHorizontal(v.terminalWidth, lipgloss.Center, helpText)
	view += "\n"

	// Show the current environments
	ecosystemView := v.ecosystem.View()
	view += lipgloss.PlaceHorizontal(v.terminalWidth, lipgloss.Center, ecosystemView)
	view += "\n\n"

	if len(v.errorMessage) > 0 {
		view += errorTitleStyle.Render("Violet ran into an error: ")
		view += "\n"
		view += errorStyle.Render(v.errorMessage)
	} else if v.spinner.show {
		currentMachine, _ := v.ecosystem.currentMachine()
		commandIndex := currentMachine.selectedCommand
		targetName := currentMachine.name
		if v.ecosystem.currentEnv().hasFocus {
			commandIndex = v.ecosystem.currentEnv().selectedCommand
			targetName = v.ecosystem.currentEnv().name
		}
		command := spinnerCommandStyle.Render(supportedMachineCommands[commandIndex])

		title := spinnerStyle.Render(fmt.Sprintf(
			"%v: %v command %v",
			targetName,
			v.spinner.verb,
			command,
		))

		progressView := fmt.Sprintf("%v %v %v\n\n", v.spinner.spinner.View(), title, v.spinner.spinner.View())
		view += lipgloss.NewStyle().Margin(marginVertical, marginHorizontal).Render(progressView)
	}

	// Monitor mouse zones and strip injected ANSI sequences
	view = zone.Scan(view)

	return view
}
