package app

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func randomEmoji() string {
	emojis := []rune("🍦🧋🍡🤠👾😭🦊🐯🦆🥨🎏🍔🍒🍥🎮📦🦁🐶🐸🍕🥐🧲🚒🥇🏆🌽")
	return string(emojis[rand.Intn(len(emojis))])
}

func (v Violet) View() string {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	// Build up final view
	view := ""

	// Title view area
	title := `
	██╗░░░██╗██╗░█████╗░██╗░░░░░███████╗████████╗
	██║░░░██║██║██╔══██╗██║░░░░░██╔════╝╚══██╔══╝
	╚██╗░██╔╝██║██║░░██║██║░░░░░█████╗░░░░░██║░░░
	░╚████╔╝░██║██║░░██║██║░░░░░██╔══╝░░░░░██║░░░
	░░╚██╔╝░░██║╚█████╔╝███████╗███████╗░░░██║░░░
	░░░╚═╝░░░╚═╝░╚════╝░╚══════╝╚══════╝░░░╚═╝░░░​​​​​`

	greeter := "Pretty manager for Vagrant " + randomEmoji()
	view += lipgloss.JoinVertical(lipgloss.Center, title, greeter)
	view = titleStyle.Render(view)
	view += "\n\n"

	// Show the current environments
	envArea := "Environments:\n"
	if v.ecosystem.environments == nil {
		envArea += "No environments found :("
	} else {
		// Create tabs
		var titleTabs []string
		for i, env := range v.ecosystem.environments {
			var style lipgloss.Style
			isFirst, _, isActive := i == 0, i == len(v.ecosystem.environments)-1, i == v.selectedEnv
			if isActive {
				style = activeTabStyle.Copy()
			} else {
				style = inactiveTabStyle.Copy()
			}
			border, _, _, _, _ := style.GetBorder()
			if isFirst && isActive {
				border.BottomLeft = "│"
			} else if isFirst && !isActive {
				border.BottomLeft = "├"
			}
			style = style.Border(border)
			titleTabs = append(titleTabs, style.Render(env.name))
		}
		tabTitleRow := lipgloss.JoinHorizontal(lipgloss.Top, titleTabs...)
		gap := tabGapStyle.Render(strings.Repeat(" ", max(0, physicalWidth-lipgloss.Width(tabTitleRow))-7))
		tabTitle := lipgloss.JoinHorizontal(lipgloss.Bottom, tabTitleRow, gap)

		vmCards := []string{}
		for i, vm := range v.ecosystem.environments[v.selectedEnv].VMs {
			vmInfo := vm.View()
			commands := v.commandButtons.View(vm.selectedCommand)
			cardInfo := lipgloss.JoinHorizontal(lipgloss.Center, vmInfo, commands)
			if i == v.selectedVM {
				cardInfo = selectedCardStyle.Render(cardInfo)
			}
			vmCards = append(vmCards, cardInfo)
		}

		tabContent := strings.Join(vmCards, "\n")

		// Not rendering the top left corder of window border, account for it with magic 2
		tabWindowStyle = tabWindowStyle.Width(lipgloss.Width(tabTitle) - 2)
		envArea += lipgloss.JoinVertical(lipgloss.Left, tabTitle, tabWindowStyle.Render(tabContent))
	}
	view += envArea
	view += "\n\n"

	// Show VMs for the selected environment
	// vmArea := ""
	// if v.ecosystem.selectedEnv == nil {
	// 	vmArea = "\n"
	// } else {
	// 	vmArea = fmt.Sprintf("VMs in %v environment:\n", v.ecosystem.selectedEnv.name)
	// 	for _, vm := range v.ecosystem.selectedEnv.VMs {
	// 		if reflect.DeepEqual(&vm, v.ecosystem.selectedEnv.selectedVM) {
	// 			vmArea += "\t[x] "
	// 		} else {
	// 			vmArea += "\t[ ] "
	// 		}
	// 		displayName := vm.name
	// 		if displayName == "" {
	// 			displayName = vm.machineID
	// 		}
	// 		vmArea += strings.Join([]string{displayName, vm.provider, vm.state}, " ")
	// 		vmArea += "\n"
	// 	}
	// }

	// if v.focus == vmView {
	// 	view += focusedStyle.Render(vmArea)
	// } else {
	// 	view += vmArea
	// }
	// view += "\n\n"

	// The available commands to run on selected VM
	// commandArea := "Commands:\n\t"
	// for _, cmd := range v.supportedCommands {
	// 	if v.selectedCommand == cmd {
	// 		commandArea += commandSelectStyle.Render(cmd + "\t")
	// 	} else {
	// 		commandArea += cmd + "\t"
	// 	}
	// }
	// view += commandArea
	// view += "\n\n"

	// Area to view output from Vagrant commands
	outputView := "Vagrant Output:\n"
	if v.spinner.show {
		outputView += fmt.Sprintf("%v %v", v.spinner.title, v.spinner.spinner.View())
	} else if v.vagrantOutputView.hasContent() {
		outputView += v.vagrantOutputView.viewport.View()
	} else {
		// Reserve the whitespace anyway
		outputView += strings.Repeat("\n", outputHeight)
	}

	view += outputView
	view += "\n\n"

	help := v.help.View(v.keys)
	view += help
	view += "\n"

	return view
}
