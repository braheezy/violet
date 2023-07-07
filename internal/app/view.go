package app

import (
	"fmt"
	"math/rand"

	"github.com/charmbracelet/lipgloss"
)

func randomEmoji() string {
	emojis := []rune("🍦🧋🍡🤠👾😭🦊🐯🦆🥨🎏🍔🍒🍥🎮📦🦁🐶🐸🍕🥐🧲🚒🥇🏆🌽")
	return string(emojis[rand.Intn(len(emojis))])
}

var verbs = []string{"Running", "Executing", "Performing", "Invoking", "Launching", "Casting"}

func (v Violet) View() (view string) {
	// Title view area
	title := titleStyle.Render("Violet: ")
	greeter := greeterStyle.Render("Pretty manager for Vagrant " + randomEmoji())
	titleGreeter := title + greeter
	view += lipgloss.NewStyle().Margin(marginVertical, marginHorizontal).Render(titleGreeter)

	help := v.help.View(v.keys)
	view += lipgloss.NewStyle().Margin(marginVertical, marginHorizontal).Render(help)
	view += "\n\n"

	// Show the current environments
	envArea := ""
	if v.ecosystem.environments == nil {
		envArea += "No environments found :("
	} else {
		envArea += v.ecosystem.View()
	}
	view += envArea
	view += "\n\n"

	if v.spinner.show {
		commandIndex := v.ecosystem.currentMachine().selectedCommand
		targetName := v.ecosystem.currentMachine().name
		if v.ecosystem.currentEnv().hasFocus {
			commandIndex = v.ecosystem.currentEnv().selectedCommand
			targetName = v.ecosystem.currentEnv().name
		}
		command := spinnerCommandStyle.Render(supportedVagrantCommands[commandIndex])

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
