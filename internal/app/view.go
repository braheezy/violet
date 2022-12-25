package app

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func randomEmoji() string {
	emojis := []rune("🍦🧋🍡🤠👾😭🦊🐯🦆🥨🎏🍔🍒🍥🎮📦🦁🐶🐸🍕🥐🧲🚒🥇🏆🌽")
	return string(emojis[rand.Intn(len(emojis))])
}

func (v Violet) View() string {
	// Build up final view
	view := ""

	// Title view area
	title := "Violet"
	greeter := "A splash of color for vagrant " + randomEmoji()
	view += lipgloss.JoinVertical(lipgloss.Center, title, greeter)
	view += "\n\n"

	// Show the current environments
	envArea := "Environments:\n"
	envArea += "\t"
	for _, env := range v.dashboard.environments {
		envArea += fmt.Sprintf("[%v]", env.name)
	}
	if v.state == environmentView {
		view += focusedStyle.Render(envArea)
	} else {
		view += envArea
	}
	view += "\n\n"

	// Show VMs for the selected environment
	vmArea := "VMs in [env2]:\n"
	vmArea += "\t[ ] vm1 (provider: virtualbox, state: running)\n"
	vmArea += "\t[x] vm2 (provider: vmware,     state: not created)\n"
	vmArea += "\t[ ] vm3 (provider: virtualbox, state: running)\n"
	if v.state == vmView {
		view += focusedStyle.Render(vmArea)
	} else {
		view += vmArea
	}
	view += "\n\n"

	// The available commands to run on selected VM
	commandArea := "Commands:\n"
	supportedCommands := []string{"up", "halt", "provision"}
	commandArea += strings.Join(supportedCommands, "\t")
	if v.state == commandView {
		view += focusedStyle.Render(commandArea)
	} else {
		view += commandArea
	}
	view += "\n\n"

	// Area to view output from Vagrant commands
	outputView := "Vagrant Output:\n"
	outputView += "\n\n\n\n\n\n\n"
	view += outputView
	view += "\n\n"

	// Area to let user type things to Vagrant terminal
	inputView := v.textInput.View()
	view += inputView
	view += "\n\n"

	help := v.help.View(v.keys)
	view += help
	view += "\n"

	return view
}
