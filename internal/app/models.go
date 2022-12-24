package app

import (
	"log"

	"github.com/braheezy/violet/pkg/vagrant"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	}
}

// **************************************************************************
//
//	Help Model
//
// **************************************************************************
// helpKeyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type helpKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k helpKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k helpKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit},                // second column
	}
}

var keys = helpKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
