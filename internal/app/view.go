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
	if v.ecosystem.environments == nil {
		envArea += "No environments found :("
	} else {
		for _, env := range v.ecosystem.environments {
			envArea += fmt.Sprintf("[%v]", env.name)
		}
	}

	if v.focus == environmentView {
		view += focusedStyle.Render(envArea)
	} else {
		view += envArea
	}
	view += "\n\n"

	// Show VMs for the selected environment
	vmArea := ""
	if v.ecosystem.selectedEnv == nil {
		vmArea = "\n"
	} else {
		vmArea = fmt.Sprintf("VMs in [%v]:\n", v.ecosystem.selectedEnv.name)
		VMs := [3]VM{
			{"vm1", "virtualbox", "running", "/vm/home/runners"},
			{"vm2", "vmware", "not created", "/vm/home/runners"},
			{"vm3", "virtualbox", "running", "/vm/home/runners"},
		}
		for _, vm := range VMs {
			if vm.name == v.ecosystem.selectedEnv.selectedVM.name {
				vmArea += "\t[x] "
			} else {
				vmArea += "\t[ ] "
			}
			vmArea += strings.Join([]string{vm.name, vm.provider, vm.state}, " ")
			vmArea += "\n"
		}
	}

	if v.focus == vmView {
		view += focusedStyle.Render(vmArea)
	} else {
		view += vmArea
	}
	view += "\n\n"

	// The available commands to run on selected VM
	commandArea := "Commands:\n"
	supportedCommands := []string{"up", "halt", "provision"}
	commandArea += strings.Join(supportedCommands, "\t")
	if v.focus == commandView {
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
