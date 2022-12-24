package app

import (
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
	envArea += "\t[env1] [env2] [env3]"
	view += envArea
	view += "\n\n"

	// Show VMs for the selected environment
	vmArea := "VMs in [env2]:\n"
	vmArea += "\t[ ] vm1 (provider: virtualbox, state: running)\n"
	vmArea += "\t[x] vm2 (provider: vmware,     state: not created)\n"
	vmArea += "\t[ ] vm3 (provider: virtualbox, state: running)\n"
	view += vmArea
	view += "\n\n"

	// The available commands to run on selected VM
	commandView := "Commands:\n"
	supportedCommands := []string{"up", "halt", "provision"}
	commandView += strings.Join(supportedCommands, "\t")
	view += commandView
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
