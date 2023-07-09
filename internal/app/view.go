package app

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (v Violet) View() (view string) {
	// Title view area
	title := titleStyle.Render("Violet:")
	greeter := greeterStyle.Render("Pretty manager for Vagrant")
	titleGreeter := title + greeter
	view += lipgloss.NewStyle().Margin(marginVertical, marginHorizontal).Render(titleGreeter)

	help := v.help.View(v.keys)
	view += lipgloss.NewStyle().
		Margin(marginVertical, marginHorizontal).
		Render(help)
	view += "\n"

	// Show the current environments
	view += v.ecosystem.View()
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

	return view
}
