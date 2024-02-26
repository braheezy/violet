package app

import (
	"errors"
	"path"
	"strings"

	"github.com/braheezy/violet/pkg/vagrant"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
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
	// Helper to paginate the list of environments
	envPager environmentPager
}

// Updates for the entire ecosystem. Usually with results from `global-status`
type ecosystemMsg Ecosystem

type ecosystemErrMsg struct{ err error }

func (e ecosystemErrMsg) Error() string { return e.err.Error() }

type environmentPager struct {
	pg             paginator.Model
	moreIsSelected bool
	backIsSelected bool
}

func (ep *environmentPager) hasMultiplePages() bool {
	return ep.pg.TotalPages > 1
}

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
	type EnvironmentGroup struct {
		Name     string
		Machines []Machine
	}
	var envGroups []EnvironmentGroup
	for _, machine := range machines {
		found := false
		for i, env := range envGroups {
			// TODO: Bug if two different paths have the same folder name e.g. /foo/env1 and /bar/env1 will incorrectly be treated the same
			if env.Name == path.Base(machine.home) {
				envGroups[i].Machines = append(envGroups[i].Machines, machine)
				found = true
				break
			}
		}
		if !found {
			env := EnvironmentGroup{
				Name:     path.Base(machine.home),
				Machines: []Machine{machine},
			}
			envGroups = append(envGroups, env)
		}
	}

	var environments []Environment
	for _, envGroup := range envGroups {
		if len(envGroup.Machines) > 0 {
			env := Environment{
				name:     envGroup.Name,
				machines: envGroup.Machines,
				home:     envGroup.Machines[0].home,
				hasFocus: true,
			}
			environments = append(environments, env)
		}
	}

	pager := paginator.New()
	pager.PerPage = 5
	pager.SetTotalPages(len(environments))

	return Ecosystem{
		environments:    environments,
		client:          client,
		machineCommands: newMachineCommandButtons(supportedMachineCommands),
		envCommands:     newEnvCommandButtons(supportedEnvCommands),
		envPager:        environmentPager{pg: pager},
	}, nil
}

// Simple helper to get the specific machine the user is interacting with
func (e *Ecosystem) currentMachine() (*Machine, error) {
	if e.selectedEnv >= len(e.environments) {
		return nil, errors.New("tried to access environment outside of ecosystem")
	} else if e.selectedMachine >= len(e.environments[e.selectedEnv].machines) {
		return nil, errors.New("tried to access machine outside of ecosystem")
	} else {
		return &e.environments[e.selectedEnv].machines[e.selectedMachine], nil
	}
}

func (e *Ecosystem) currentEnv() *Environment {
	return &e.environments[e.selectedEnv]
}

func (e *Ecosystem) View() (result string) {
	if e.environments == nil {
		return lipgloss.NewStyle().Foreground(textColor).Italic(true).Faint(true).Render("Still looking for environments...")
	}

	// Create the tab headers, one for each environment.
	var tabs []string
	start, end := e.envPager.pg.GetSliceBounds(len(e.environments))
	for i, env := range e.environments[start:end] {
		// Figure out which "tab" is selected and stylize accordingly
		var style lipgloss.Style
		idx := i
		if e.envPager.pg.Page > 0 {
			idx = i + e.envPager.pg.PerPage
		}
		isFirst, _, isActive := idx == start, idx == len(e.environments)-1, idx == e.selectedEnv
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		// Override border edges for these edge cases
		if e.envPager.pg.Page == 0 {
			if isFirst && isActive {
				border.BottomLeft = "│"
			} else if isFirst && !isActive {
				border.BottomLeft = "├"
			}
		}

		style = style.Border(border)
		tabs = append(tabs, zone.Mark(env.name, style.Render(env.name)))
	}

	var tabContent string

	// If there's paged environments, show a tab with a paged indicator
	if e.envPager.hasMultiplePages() {
		// Show a More button if there's additional pages
		var moreTab string
		var backTab string
		if e.envPager.pg.Page < e.envPager.pg.TotalPages-1 {
			moreTab = zone.Mark("more", "⮕ ")
			if e.envPager.moreIsSelected {
				moreTab = activeTabStyle.Render(moreTab)
			} else {
				moreTab = inactiveTabStyle.Render(moreTab)
			}
			tabs = append(tabs, moreTab)
		}
		// Show a Back button if there's previous pages
		if e.envPager.pg.Page > 0 {
			backTab = zone.Mark("back", "⬅ ")
			if e.envPager.backIsSelected {
				border, _, _, _, _ := activeTabStyle.GetBorder()
				border.BottomLeft = "│"
				style := activeTabStyle.Copy().Border(border)
				backTab = style.Render(backTab)
			} else {
				border, _, _, _, _ := inactiveTabStyle.GetBorder()
				border.BottomLeft = "├"
				style := inactiveTabStyle.Copy().Border(border)
				backTab = style.Render(backTab)
			}
			tabs = append([]string{backTab}, tabs...)
		}
	}

	if e.envPager.moreIsSelected {
		// Show More tab content
		tabContent = "There's more stuff ova there ->\nHit ENTER"
	} else if e.envPager.backIsSelected {
		tabContent = "<- There's stuff bak there\nHit ENTER"
	} else {
		// machineCards will be the set of machines to show for the selected env.
		// They are dealt with first so we know the size of content we need to
		// wrap in "tabs"
		machineCards := []string{}
		selectedEnv := e.environments[e.selectedEnv]
		for i, machine := range selectedEnv.machines {
			// "Viewing" a machine will get it's specific info
			machineView := machine.View()
			commands := e.machineCommands.View(machine.selectedCommand, !selectedEnv.hasFocus)
			cardInfo := lipgloss.JoinHorizontal(lipgloss.Center, machineView, commands)
			if !selectedEnv.hasFocus && i == e.selectedMachine {
				cardInfo = selectedCardStyle.Render(cardInfo)
			} else {
				cardInfo = defaultCardStyle.Render(cardInfo)
			}
			machineCards = append(machineCards, cardInfo)
		}

		// This card always exists and controls the top-level environment
		envTitle := envCardTitleStyle.Render(selectedEnv.name)
		envCommands := newEnvCommandButtons(supportedEnvCommands)
		if selectedEnv.hasFocus {
			envTitle = selectedEnvCardStyle.Render(selectedEnv.name)
		}
		envCard := lipgloss.JoinHorizontal(lipgloss.Center, envTitle, envCommands.View(selectedEnv.selectedCommand, selectedEnv.hasFocus))

		tabContent = envCard + "\n" + strings.Join(machineCards, "\n")
	}

	tabHeader := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	// Create the window effect by creating a blank tab to fill the rest of the width.
	gapWidth := lipgloss.Width(tabContent) - lipgloss.Width(tabHeader)
	if gapWidth < 0 {
		// There's more tabs than the standard width of a tab, so add padding
		tabContent = lipgloss.NewStyle().MarginRight(gapWidth * -1).Render(tabContent)
		gapWidth = 0
	}
	gap := tabGapStyle.Render(strings.Repeat(" ", gapWidth))
	tabHeader = lipgloss.JoinHorizontal(lipgloss.Top, tabHeader, gap)

	result = lipgloss.JoinVertical(lipgloss.Center, tabHeader, tabWindowStyle.Render(tabContent))

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
		lipgloss.Right,
		cardTitleStyle.Render(displayName),
		cardStatusStyle.Foreground(statusColors[m.state]).Render(m.state),
		cardProviderStyle.Render(m.provider),
	)

	return content
}
