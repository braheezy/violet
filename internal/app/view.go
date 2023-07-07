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

var verbs = []string{"Running", "Executing", "Performing", "Invoking", "Launching", "Casting"}

func (v Violet) View() (view string) {
	// Title view area
	title := titleStyle.Render("Violet: ")
	greeter := greeterStyle.Render("Pretty manager for Vagrant " + randomEmoji())
	titleGreeter := title + greeter
	view += lipgloss.NewStyle().Margin(1, 2).Render(titleGreeter)
	view += "\n"

	help := v.help.View(v.keys)
	view += lipgloss.NewStyle().Margin(0, 2).Render(help)
	view += "\n\n"

	// Show the current environments
	envArea := ""
	if v.ecosystem.environments == nil {
		envArea += "No environments found :("
	} else {
		// vmCards will be the set of VMs to show for the selected env.
		// They are dealt with first so we know the size of content we need to
		// wrap in "tabs"
		vmCards := []string{}
		selectedEnv := v.currentEnv()
		for i, vm := range selectedEnv.VMs {
			// "Viewing" a VM will get it's specific info
			vmInfo := vm.View()
			// Commands are the same for everyone so they are grabbed from the main model
			commands := v.layout.commandButtons.View(vm.selectedCommand)
			cardInfo := lipgloss.JoinHorizontal(lipgloss.Center, vmInfo, commands)
			if !selectedEnv.hasFocus && i == v.selectedVM {
				cardInfo = selectedCardStyle.Render(cardInfo)
			}
			vmCards = append(vmCards, cardInfo)
		}

		// This card always exists and controls the top-level environment
		envTitle := lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Render(selectedEnv.name)
		envCommands := v.layout.commandButtons.View(selectedEnv.selectedCommand)
		if selectedEnv.hasFocus {
			envTitle = lipgloss.NewStyle().
				Bold(true).
				Foreground(accentColor).
				Render(selectedEnv.name)
		}
		envCard := lipgloss.JoinHorizontal(lipgloss.Center, envTitle, envCommands)

		tabContent := envCard + "\n" + strings.Join(vmCards, "\n")

		var tabs []string
		for i, env := range v.ecosystem.environments {
			// Figure out which "tab" is selected and stylize accordingly
			var style lipgloss.Style
			isFirst, _, isActive := i == 0, i == len(v.ecosystem.environments)-1, i == v.selectedEnv
			if isActive {
				style = activeTabStyle.Copy()
			} else {
				style = inactiveTabStyle.Copy()
			}
			border, _, _, _, _ := style.GetBorder()
			// Override border edges for these edge cases
			if isFirst && isActive {
				border.BottomLeft = "│"
			} else if isFirst && !isActive {
				border.BottomLeft = "├"
			}
			style = style.Border(border)
			tabs = append(tabs, style.Render(env.name))
		}
		// This trick is how the "window" effect is realized: "empty tab" to fill the width.
		commandsWidth := v.layout.commandButtons.width
		gap := tabGapStyle.Render(strings.Repeat(" ", commandsWidth))
		tabs = append(tabs, gap)
		tabHeader := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

		// Not rendering the top left corder of window border, account for it with magic 2 :(
		tabWindowStyle = tabWindowStyle.Width(lipgloss.Width(tabHeader) - 2)
		envArea += lipgloss.JoinVertical(lipgloss.Left, tabHeader, tabWindowStyle.Render(tabContent))
		envArea = lipgloss.NewStyle().Padding(0, 2).Render(envArea)
	}
	view += envArea
	view += "\n\n"

	if v.layout.spinner.show {
		commandIndex := v.currentVM().selectedCommand
		targetName := v.currentVM().name
		if v.currentEnv().hasFocus {
			commandIndex = v.currentEnv().selectedCommand
			targetName = v.currentEnv().name
		}
		command := spinnerCommandStyle.Render(supportedVagrantCommands[commandIndex])

		title := spinnerStyle.Render(fmt.Sprintf(
			"%v: %v command %v",
			targetName,
			v.layout.spinner.verb,
			command,
		))

		progressView := fmt.Sprintf("%v %v %v\n\n", v.layout.spinner.spinner.View(), title, v.layout.spinner.spinner.View())
		view += lipgloss.NewStyle().Padding(0, 2).Render(progressView)
	}

	return view
}
