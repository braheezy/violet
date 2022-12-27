package app

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
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
}

func newViolet() Violet {

	textInput := textinput.New()
	textInput.Placeholder = "Send text to the terminal running Vagrant..."

	return Violet{
		ecosystem: Ecosystem{
			environments: nil,
			selectedEnv:  nil,
		},
		keys:      keys,
		help:      help.New(),
		textInput: textInput,
		focus:     environmentView,
	}
}
