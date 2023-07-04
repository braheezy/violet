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
	view += help
	view += "\n\n"

	// Show the current environments
	envTitle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Padding(0, 1).
		Render("Environments:")
	envArea := "\n\n"
	if v.ecosystem.environments == nil {
		envArea += "No environments found :("
	} else {
		// vmCards will be the set of VMs to show for the selected env.
		// They are dealt with first so we know the size of content we need to
		// wrap in "tabs"
		vmCards := []string{}
		commandsWidth := 0
		for i, vm := range v.ecosystem.environments[v.selectedEnv].VMs {
			// "Viewing" a VM will get it's specific info
			vmInfo := vm.View()
			// Commands are the same for everyone so they are grabbed from the main model
			commands := v.layout.commandButtons.View(vm.selectedCommand)
			commandsWidth = v.layout.commandButtons.width
			cardInfo := lipgloss.JoinHorizontal(lipgloss.Center, vmInfo, commands)
			if i == v.selectedVM {
				cardInfo = selectedCardStyle.Render(cardInfo)
			}
			vmCards = append(vmCards, cardInfo)
		}

		tabContent := strings.Join(vmCards, "\n")

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
		gap := tabGapStyle.Render(strings.Repeat(" ", commandsWidth))
		tabs = append(tabs, gap)
		tabHeader := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

		// Not rendering the top left corder of window border, account for it with magic 2 :(
		tabWindowStyle = tabWindowStyle.Width(lipgloss.Width(tabHeader) - 2)
		envArea += lipgloss.JoinVertical(lipgloss.Left, tabHeader, tabWindowStyle.Render(tabContent))
		envArea = lipgloss.NewStyle().Padding(0, 2).Render(envArea)
	}
	view += envTitle + envArea
	view += "\n\n"

	outputView := ""
	if v.layout.spinner.show {
		outputView = fmt.Sprintf("%v %v\n\n", v.layout.spinner.title, v.layout.spinner.spinner.View())
	}
	view += outputView

	return view
}
