package app

import (
	"log"

	"github.com/braheezy/violet/pkg/vagrant"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

var supportedVagrantCommands = []string{"up", "halt", "reload", "provision"}

// **************************************************************************
//
//	Ecosystem Model
//
// **************************************************************************
// Ecosystem contains the total Vagrant world information
type Ecosystem struct {
	// Collection of all Vagrant environments
	environments []Environment
	// Reference to a Vagrant client to run commands with
	client *vagrant.VagrantClient
}

// **************************************************************************
//
//	Environment Model
//
// **************************************************************************
// Environment represents a single Vagrant project
type Environment struct {
	// Friendly name for the Environment
	name string
	// Environments have 0 or more VMs
	VMs []VM
}

// **************************************************************************
//
//	VM Model
//
// **************************************************************************
// VM contains all the data and actions associated with a specific VM
type VM struct {
	name            string
	provider        string
	state           string
	home            string
	machineID       string
	selectedCommand int
}

func (vm *VM) View() string {
	displayName := vm.name
	if displayName == "" {
		displayName = vm.machineID
	}

	// displayName = turnBig(displayName)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		cardTitleStyle.Render(displayName),
		cardStatusStyle.Foreground(statusColors[vm.state]).Render(vm.state),
		cardProviderStyle.Render(vm.provider),
	)

	return content

	// return strings.Join([]string{displayName, vm.provider, vm.state}, " ")
}

// **************************************************************************
//
//	Violet Model
//
// **************************************************************************
type outputViewport struct {
	viewport viewport.Model
}

func (o *outputViewport) hasContent() bool {
	currentSize := len(o.viewport.View())
	defaultSize := outputHeight * outputWidth
	return currentSize > defaultSize
}

// Complete app state (i.e. the BubbleTea model)
type Violet struct {
	// Reference to the Ecosystem
	ecosystem Ecosystem
	// Fancy help bubble
	help help.Model
	// To support help
	keys        helpKeyMap
	selectedEnv int
	selectedVM  int
	// The viewport to view Vagrant output
	vagrantOutputView outputViewport
	commandButtons    buttonGroup
	spinner           currentSpinner
}

func newViolet() Violet {
	client, err := vagrant.NewVagrantClient()
	if err != nil {
		log.Fatal(err)
	}

	help := help.New()
	help.ShowAll = true

	vagrantOutputView := viewport.New(outputWidth, outputHeight)
	vagrantOutputView.Style = outputViewStyle

	return Violet{
		ecosystem: Ecosystem{
			environments: nil,
			client:       client,
		},
		keys:              keys,
		help:              help,
		selectedEnv:       0,
		selectedVM:        0,
		vagrantOutputView: outputViewport{vagrantOutputView},
		commandButtons:    newCommandButtons(),
		spinner:           newSpinner(),
	}
}

func (v *Violet) getCurrentVM() *VM {
	return &v.ecosystem.environments[v.selectedEnv].VMs[v.selectedVM]
}

type currentSpinner struct {
	spinner spinner.Model
	show    bool
	title   string
}

func newSpinner() currentSpinner {
	s := spinner.New()
	s.Spinner = spinners[0]
	// s.Style = spinnerStyle
	return currentSpinner{
		spinner: s,
		show:    false,
		title:   "",
	}
}
