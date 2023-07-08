package app

import (
	"path"
	"strings"

	"github.com/braheezy/violet/pkg/vagrant"
	"github.com/charmbracelet/lipgloss"
)

// Ecosystem contains the total Vagrant world information
type Ecosystem struct {
	// Collection of all Vagrant environments
	environments []Environment
	// Reference to a Vagrant client to run commands with
	client *vagrant.VagrantClient
	// Buttons to allow the user to run commands
	machineCommands machineCommandButtons
	envCommands     envCommandButtons
	// Indexes of the respective lists that are currently selected.
	selectedEnv     int
	selectedMachine int
}

// Updates for the entire ecosystem. Usually with results from `global-status`
type ecosystemMsg Ecosystem

type ecosystemErrMsg struct{ err error }

func (e ecosystemErrMsg) Error() string { return e.err.Error() }

// Call `global-status` and translate result into a new Ecosystem
func createEcosystem(client *vagrant.VagrantClient) (Ecosystem, error) {
	// Fetch (not stream) the current global status
	result, err := client.GetGlobalStatus()
	var nilEcosystem Ecosystem

	if err != nil {
		return nilEcosystem, ecosystemErrMsg{err}
	}

	results := vagrant.ParseVagrantOutput(result)
	if results == nil {
		return nilEcosystem, nil
	}

	var machines []Machine
	for _, machineInfo := range results {
		machine := Machine{
			machineID: machineInfo.MachineID,
			provider:  machineInfo.Fields["provider-name"],
			state:     strings.Replace(machineInfo.Fields["state"], "_", " ", -1),
			home:      machineInfo.Fields["machine-home"],
		}
		machines = append(machines, machine)
	}
	// Create different envs by grouping machines based on machine-home
	envGroups := make(map[string][]Machine)
	for _, machine := range machines {
		// TODO: Bug if two different paths have the same folder name e.g. /foo/env1 and /bar/env1 will incorrectly be treated the same
		envGroups[path.Base(machine.home)] = append(envGroups[path.Base(machine.home)], machine)
	}
	var environments []Environment
	for envName, machines := range envGroups {
		env := Environment{
			name:     envName,
			machines: machines,
			home:     envGroups[envName][0].home,
			hasFocus: true,
		}
		environments = append(environments, env)
	}
	return Ecosystem{
		environments:    environments,
		client:          client,
		machineCommands: newMachineCommandButtons(supportedMachineCommands),
		envCommands:     newEnvCommandButtons(supportedEnvCommands),
	}, nil
}

// Simple helper to get the specific machine the user is interacting with
func (e *Ecosystem) currentMachine() *Machine {
	return &e.environments[e.selectedEnv].machines[e.selectedMachine]
}

func (e *Ecosystem) currentEnv() *Environment {
	return &e.environments[e.selectedEnv]
}

func (e *Ecosystem) View() (result string) {
	if e.environments == nil {
		return "No environments found :("
	}

	// machineCards will be the set of machines to show for the selected env.
	// They are dealt with first so we know the size of content we need to
	// wrap in "tabs"
	machineCards := []string{}
	selectedEnv := e.environments[e.selectedEnv]
	for i, machine := range selectedEnv.machines {
		// "Viewing" a machine will get it's specific info
		machineView := machine.View()
		commands := e.machineCommands.View(machine.selectedCommand)
		cardInfo := lipgloss.JoinHorizontal(lipgloss.Top, machineView, commands)
		if !selectedEnv.hasFocus && i == e.selectedMachine {
			cardInfo = selectedCardStyle.Render(cardInfo)
		} else {
			cardInfo = defaultCardStyle.Render(cardInfo)
		}
		machineCards = append(machineCards, cardInfo)

		// This card always exists and controls the top-level environment
		envTitle := envCardTitleStyle.Render(selectedEnv.name)
		envCommands := newEnvCommandButtons(supportedEnvCommands)
		// envCommands = e.commandButtons.View(selectedEnv.selectedCommand)
		if selectedEnv.hasFocus {
			envTitle = selectedEnvCardStyle.Render(selectedEnv.name)
		}
		envCard := lipgloss.JoinHorizontal(lipgloss.Center, envTitle, envCommands.View(selectedEnv.selectedCommand))

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

		tabHeader := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
		// Create the window effect by creating a blank tab to fill the rest of the width.
		gapWidth := lipgloss.Width(tabContent) - lipgloss.Width(tabHeader)
		gap := tabGapStyle.Render(strings.Repeat(" ", gapWidth))
		tabHeader = lipgloss.JoinHorizontal(lipgloss.Top, tabHeader, gap)

		result = lipgloss.JoinVertical(lipgloss.Left, tabHeader, tabWindowStyle.Render(tabContent))

	}
	return result
}

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
