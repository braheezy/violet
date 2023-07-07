package app

import (
	"log"
	"strings"

	"github.com/braheezy/violet/pkg/vagrant"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

// Order matters here.
var supportedVagrantCommands = []string{"up", "halt", "ssh", "reload", "provision"}

// Environment represents a single Vagrant project
type Environment struct {
	// Friendly name for the Environment
	name string
	// Environments have 0 or more machines
	machines []Machine
	// The currently selected command to run on the machine.
	selectedCommand int
	home            string
	hasFocus        bool
}

// Machine contains all the data and actions associated with a specific Machine
type Machine struct {
	name      string
	provider  string
	state     string
	home      string
	machineID string
	// The currently selected command to run on the machine.
	selectedCommand int
}

func (m *Machine) View() string {
	displayName := m.name
	// If there's no name yet, at least show the machineID
	if displayName == "" {
		displayName = m.machineID
	}

	// Join the machine info for the card view
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		cardTitleStyle.Render(displayName),
		cardStatusStyle.Foreground(statusColors[m.state]).Render(m.state),
		cardProviderStyle.Render(m.provider),
	)

	return content
}

// Complete app state (i.e. the BubbleTea model)
type Violet struct {
	// Reference to the Ecosystem
	ecosystem Ecosystem
	// Fancy help bubble
	help help.Model
	// To support help
	keys helpKeyMap
	// Spinner to show while commands are running
	spinner currentSpinner
	// Current terminal size
	terminalWidth  int
	terminalHeight int
}

// Return the default Violet model
func newViolet() Violet {
	client, err := vagrant.NewVagrantClient()
	if err != nil {
		log.Fatal(err)
	}

	help := help.New()
	help.ShowAll = true

	return Violet{
		ecosystem: Ecosystem{
			environments: nil,
			client:       client,
		},
		keys:    keys,
		help:    help,
		spinner: newSpinner(),
	}
}

// Simple helper to get the specific machine the user is interacting with
func (e *Ecosystem) currentMachine() *Machine {
	return &e.environments[e.selectedEnv].machines[e.selectedMachine]
}

func (e *Ecosystem) currentEnv() *Environment {
	return &e.environments[e.selectedEnv]
}

func (e *Ecosystem) View() (result string) {
	// machineCards will be the set of machines to show for the selected env.
	// They are dealt with first so we know the size of content we need to
	// wrap in "tabs"
	machineCards := []string{}
	selectedEnv := e.environments[e.selectedEnv]
	for i, machine := range selectedEnv.machines {
		// "Viewing" a machine will get it's specific info
		machineView := machine.View()
		// Commands are the same for everyone so they are grabbed from the main model
		commands := e.commandButtons.View(machine.selectedCommand)
		cardInfo := lipgloss.JoinHorizontal(lipgloss.Center, machineView, commands)
		if !selectedEnv.hasFocus && i == e.selectedMachine {
			cardInfo = selectedCardStyle.Render(cardInfo)
		}
		machineCards = append(machineCards, cardInfo)

		// This card always exists and controls the top-level environment
		envTitle := envCardTitleStyle.Render(selectedEnv.name)
		envCommands := e.commandButtons.View(selectedEnv.selectedCommand)
		if selectedEnv.hasFocus {
			envTitle = selectedEnvCardStyle.Render(selectedEnv.name)
		}
		envCard := lipgloss.JoinHorizontal(lipgloss.Center, envTitle, envCommands)

		tabContent := envCard + "\n" + strings.Join(machineCards, "\n")

		// Now create the tab headers, one for each environment.
		var tabs []string
		for i, env := range e.environments {
			// Figure out which "tab" is selected and stylize accordingly
			var style lipgloss.Style
			isFirst, _, isActive := i == 0, i == len(e.environments)-1, i == e.selectedEnv
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
		// Create the window effect by creating a blank tab to fill the rest of the width.
		commandsWidth := e.commandButtons.width
		gap := tabGapStyle.Render(strings.Repeat(" ", commandsWidth))
		tabs = append(tabs, gap)
		tabHeader := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

		// Not rendering the top left corder of window border, account for it with magic 2 :(
		tabWindowStyle = tabWindowStyle.Width(lipgloss.Width(tabHeader) - 2)
		result = lipgloss.JoinVertical(lipgloss.Left, tabHeader, tabWindowStyle.Render(tabContent))
		result = lipgloss.NewStyle().Padding(0, 2).Render(result)
	}
	return result
}
