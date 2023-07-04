package app

import (
	"log"

	"github.com/braheezy/violet/pkg/vagrant"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

var supportedVagrantCommands = []string{"up", "halt", "reload", "provision"}

type Layout struct {
	// Spinner to show while commands are running
	spinner currentSpinner
	// Buttons to allow the user to run commands
	commandButtons buttonGroup
}

// Environment represents a single Vagrant project
type Environment struct {
	// Friendly name for the Environment
	name string
	// Environments have 0 or more VMs
	VMs []VM
}

// VM contains all the data and actions associated with a specific VM
type VM struct {
	name      string
	provider  string
	state     string
	home      string
	machineID string
	// The currently selected command to run on the VM.
	selectedCommand int
}

func (vm *VM) View() string {
	displayName := vm.name
	// If there's no name yet, at least show the machine-id
	if displayName == "" {
		displayName = vm.machineID
	}

	// Join the VM info for the card view
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		cardTitleStyle.Render(displayName),
		cardStatusStyle.Foreground(statusColors[vm.state]).Render(vm.state),
		cardProviderStyle.Render(vm.provider),
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
	// Indexes of the respective lists that are currently selected.
	selectedEnv int
	selectedVM  int
	// Current layout to use
	layout Layout
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
		keys:        keys,
		help:        help,
		selectedEnv: 0,
		selectedVM:  0,
		layout: Layout{
			spinner:        newSpinner(),
			commandButtons: newCommandButtons(),
		},
	}
}

// Simple helper to get the specific VM the user is interacting with
func (v *Violet) getCurrentVM() *VM {
	return &v.ecosystem.environments[v.selectedEnv].VMs[v.selectedVM]
}
