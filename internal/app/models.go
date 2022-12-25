package app

import (
	"log"

	"github.com/braheezy/violet/pkg/vagrant"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
)

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
	VMs []VM
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
//	Dashboard Model
//
// **************************************************************************
// Dashboard provides a single pane into all Environments. Includes methods for managing Environments (and VMs).
type Dashboard struct {
	environments []Environment
	selected     *Environment
}

// **************************************************************************
//
//	Violet Model
//
// **************************************************************************
// sessionState is used to track which model is focused
type sessionState uint

// Enumerate available states
const (
	environmentView sessionState = iota
	vmView
	commandView
)

// Complete app state (i.e. the BubbleTea model)
type Violet struct {
	// Reference to the VagrantClient to use for all calls.
	vagrantClient *vagrant.VagrantClient
	// Reference to the Dashboard
	dashboard Dashboard
	// Fancy help bubble
	help help.Model
	// To support help
	keys helpKeyMap
	// User input to Vagrant terminal
	textInput textinput.Model
	// The currently selected view
	state sessionState
}

func newViolet() Violet {
	vagrantClient, err := vagrant.NewVagrantClient()
	if err != nil {
		log.Fatal(err)
	}

	textInput := textinput.New()
	textInput.Placeholder = "Send text to the terminal running Vagrant..."

	return Violet{
		vagrantClient: vagrantClient,
		dashboard: Dashboard{
			environments: nil,
			selected:     nil,
		},
		keys:      keys,
		help:      help.New(),
		textInput: textInput,
		state:     environmentView,
	}
}
