package app

import (
	"fmt"
	"math/rand"
	"reflect"
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
			if env.name == v.ecosystem.selectedEnv.name {
				envArea += focusedStyle.Render(fmt.Sprintf("[%v]\t", env.name))
			} else {
				envArea += fmt.Sprintf("[%v]\t", env.name)
			}
		}
	}
	view += envArea
	view += "\n\n"

	// Show VMs for the selected environment
	vmArea := ""
	if v.ecosystem.selectedEnv == nil {
		vmArea = "\n"
	} else {
		vmArea = fmt.Sprintf("VMs in %v environment:\n", v.ecosystem.selectedEnv.name)
		for _, vm := range v.ecosystem.selectedEnv.VMs {
			if reflect.DeepEqual(&vm, v.ecosystem.selectedEnv.selectedVM) {
				vmArea += "\t[x] "
			} else {
				vmArea += "\t[ ] "
			}
			displayName := vm.name
			if displayName == "" {
				displayName = vm.machineID
			}
			vmArea += strings.Join([]string{displayName, vm.provider, vm.state}, " ")
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
	commandArea := "Commands:\n\t"
	for _, cmd := range v.supportedCommands {
		if v.selectedCommand == cmd {
			commandArea += commandSelectStyle.Render(cmd + "\t")
		} else {
			commandArea += cmd + "\t"
		}
	}
	if v.focus == commandView {
		view += focusedStyle.Render(commandArea)
	} else {
		view += commandArea
	}
	view += "\n\n"

	// Area to view output from Vagrant commands
	outputView := "Vagrant Output:\n"
	if v.vagrantOutputView.hasContent() {
		outputView += v.vagrantOutputView.viewport.View()
	} else {
		// Reserve the whitespace anyway
		outputView += strings.Repeat("\n", outputWidth)
	}
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
