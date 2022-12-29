package app

import (
	"log"

	"github.com/braheezy/violet/pkg/vagrant"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
)

// **************************************************************************
//
//	Ecosystem Model
//
// **************************************************************************
// Ecosystem provides a single pane into all Environments
type Ecosystem struct {
	environments []Environment
	selectedEnv  *Environment
}

func (d *Ecosystem) UpdateEnvs(envs []Environment) {
	d.environments = envs
}

// **************************************************************************
//
//	Environment Model
//
// **************************************************************************
// Environment represents a single Vagrant project
type Environment struct {
	// A friendly name for the Environment
	name string
	// Environments have 0 or more VMs
	VMs        []VM
	selectedVM *VM
}

// **************************************************************************
//
//	VM Model
//
// **************************************************************************
// VM contains all the data and actions associated with a specific VM
type VM struct {
	name     string
	provider string
	state    string
	home     string
	// Other VM properties...
}

// **************************************************************************
//
//	Violet Model
//
// **************************************************************************
// focusState is used to track which model is focused
type focusState uint

// Enumerate available areas the user can focus
const (
	environmentView focusState = iota
	vmView
	commandView
)

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
	keys helpKeyMap
	// User input to Vagrant terminal
	textInput textinput.Model
	// The currently selected view
	focus focusState
	// A copy of the Vagrant client to share around
	client *vagrant.VagrantClient
	// The Vagrant commands that Violet can run
	supportedCommands []string
	// Currently selected Vagrant command to run
	selectedCommand string
	// The viewport to view Vagrant output
	vagrantOutputView outputViewport
}

func newViolet() Violet {
	client, err := vagrant.NewVagrantClient()
	if err != nil {
		log.Fatal(err)
	}

	textInput := textinput.New()
	textInput.Placeholder = "Send text to the terminal running Vagrant..."

	help := help.New()
	help.ShowAll = true

	vagrantOutputView := viewport.New(outputHeight, outputWidth)
	vagrantOutputView.Style = outputViewStyle

	return Violet{
		ecosystem: Ecosystem{
			environments: nil,
			selectedEnv:  nil,
		},
		keys:              keys,
		help:              help,
		textInput:         textInput,
		focus:             environmentView,
		client:            client,
		supportedCommands: []string{"up", "halt", "provision", "ssh"},
		selectedCommand:   "up",
		vagrantOutputView: outputViewport{vagrantOutputView},
	}
}
