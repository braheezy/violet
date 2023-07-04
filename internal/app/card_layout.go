package app

import (
	"fmt"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type cardLayout struct {
	// Spinner to show while commands are running
	spinner currentSpinner
	// Buttons to allow the user to run commands
	commandButtons buttonGroup
}

var verbs = []string{"Running", "Executing", "Performing", "Invoking", "Launching", "Casting"}

func (c *cardLayout) UpdatePreExec(cmd string, name string) tea.Cmd {
	c.spinner.show = true
	c.spinner.title = fmt.Sprintf(
		"%v %v command on %v...",
		verbs[rand.Intn(len(verbs))],
		cmd,
		name)
	// This must be sent for the spinner to spin
	tickCmd := c.spinner.spinner.Tick
	return tickCmd
}

func (c *cardLayout) UpdatePostExec() {
	c.spinner.show = false
	c.spinner.spinner.Spinner = spinners[rand.Intn(len(spinners))]
}

func (c *cardLayout) View(v *Violet) (view string) {
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
		// vmCards will be the set of VMs to show for the selected env.
		// They are dealt with first so we know the size of content we need to
		// wrap in "tabs"
		vmCards := []string{}
		commandsWidth := 0
		for i, vm := range v.ecosystem.environments[v.selectedEnv].VMs {
			// "Viewing" a VM will get it's specific info
			vmInfo := vm.View()
			// Commands are the same for everyone so they are grabbed from the main model
			commands := c.commandButtons.View(vm.selectedCommand)
			commandsWidth = c.commandButtons.width
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
		// This trick is how they "window" effect is realized: "empty tab" to fill the width.
		gap := tabGapStyle.Render(strings.Repeat(" ", commandsWidth*2))
		tabs = append(tabs, gap)
		tabHeader := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

		// Not rendering the top left corder of window border, account for it with magic 2 :(
		tabWindowStyle = tabWindowStyle.Width(lipgloss.Width(tabHeader) - 2)
		envArea += lipgloss.JoinVertical(lipgloss.Left, tabHeader, tabWindowStyle.Render(tabContent))
	}
	view += envArea
	view += "\n\n"

	// Area to view output from Vagrant commands
	outputView := ""
	if c.spinner.show {
		outputView = "Vagrant Output:\n\n"
		outputView += fmt.Sprintf("%v %v", c.spinner.title, c.spinner.spinner.View())
		// Maintain whitespace to keep help view from jumping around
		outputView += strings.Repeat("\n", outputHeight-1)
	} else {
		// Reserve the whitespace anyway
		outputView += strings.Repeat("\n", outputHeight)
	}
	view += outputView

	help := v.help.View(v.keys)
	view += help
	view += "\n"

	return view
}

func (c cardLayout) UpdateAlways(msg tea.Msg) tea.Cmd {
	// Spinner needs spinCmd every update to know to keep spinning?
	if c.spinner.show {
		var spinCmd tea.Cmd
		c.spinner.spinner, spinCmd = c.spinner.spinner.Update(msg)
		return spinCmd
	}
	return nil
}
