package app

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// helpKeyMap defines a set of keybindings.
type helpKeyMap struct {
	Up            key.Binding
	Down          key.Binding
	Left          key.Binding
	Right         key.Binding
	Switch        key.Binding
	Execute       key.Binding
	ChooseCommand key.Binding
	Help          key.Binding
	Quit          key.Binding
}

// Setup the keybinding and help text for each key
var keys = helpKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "scroll output up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "scroll output down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "select vm"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "select vm"),
	),
	Switch: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("⭾ tab", "switch env tab"),
	),
	Execute: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("⏎ enter", "run selected command"),
	),
	ChooseCommand: key.NewBinding(
		key.WithKeys("1", "2", "3", "4"),
		key.WithHelp("number", "select command by number"),
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

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k helpKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k helpKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},                        // first column
		{k.Switch, k.Execute, k.ChooseCommand, k.Help, k.Quit}, // second column
	}
}

var verbs = []string{"Running", "Executing", "Performing", "Invoking", "Launching", "Casting"}

func (v Violet) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Window was resized
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can it can gracefully truncate
		// its view as needed.
		v.help.Width = msg.Width

	// User pressed a key
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.ChooseCommand):
			selected, _ := strconv.Atoi(msg.String())
			v.getCurrentVM().selectedCommand = selected - 1
			// v.ecosystem.environments[v.selectedEnv].VMs[v.selectedVM].selectedCommand = selected - 1
		case key.Matches(msg, v.keys.Left):
			if v.selectedVM == 0 {
				v.selectedVM = len(v.ecosystem.environments[v.selectedEnv].VMs) - 1
			} else {
				v.selectedVM -= 1
			}
		case key.Matches(msg, v.keys.Right):
			if v.selectedVM == len(v.ecosystem.environments[v.selectedEnv].VMs)-1 {
				v.selectedVM = 0
			} else {
				v.selectedVM += 1
			}
		case key.Matches(msg, v.keys.Switch):
			if v.selectedEnv == len(v.ecosystem.environments)-1 {
				v.selectedEnv = 0
			} else {
				v.selectedEnv += 1
			}
			return v, nil
		case key.Matches(msg, v.keys.Execute):
			currentVM := v.getCurrentVM()
			v.spinner.show = true
			vagrantCmd := supportedVagrantCommands[currentVM.selectedCommand]
			v.spinner.title = fmt.Sprintf(
				"%v %v command on %v...",
				verbs[rand.Intn(len(verbs))],
				vagrantCmd,
				currentVM.name)
			// This must be sent for the spinner to spin
			tickCmd := v.spinner.spinner.Tick
			// Run the command async and stream result back
			streamCmd := v.streamCommandOnVM(
				vagrantCmd,
				currentVM.machineID,
			)
			return v, tea.Batch(tickCmd, streamCmd)
		case key.Matches(msg, v.keys.Help):
			v.help.ShowAll = !v.help.ShowAll
		case key.Matches(msg, v.keys.Quit):
			return v, tea.Quit
		}

	// New data from `global-status` has come in
	case ecosystemMsg:
		eco := Ecosystem(msg)
		var statusCmds []tea.Cmd
		// Don't have the VM names yet, just machine-ids.
		// Queue up a bunch of async calls to go get those names.
		for _, env := range eco.environments {
			for _, vm := range env.VMs {
				statusCmds = append(statusCmds, v.getVMStatus(vm.machineID))
			}
		}
		// Set the new ecosystem
		v.ecosystem = eco

		return v, tea.Batch(statusCmds...)

	// New data about a specific VM has come in
	case statusMsg:
		// Find the VM this message is about
		for i, env := range v.ecosystem.environments {
			for j, vm := range env.VMs {
				if msg.identifier == vm.machineID || msg.identifier == vm.name {
					// Found the VM this status message is about.
					// Status msgs don't return some info so retain existing info
					updatedVM := VM{
						machineID: vm.machineID,
						provider:  msg.status.Fields["provider-name"],
						state:     msg.status.Fields["state"],
						home:      vm.home,
						name:      msg.status.Name,
						// Reset the selected command
						selectedCommand: 0,
					}
					v.ecosystem.environments[i].VMs[j] = updatedVM
				}
			}
		}

	// Result from a command has been streamed in
	case streamMsg:
		// Put the content directly in the viewport
		v.vagrantOutputView.viewport.SetContent(string(msg))
		// Stop spinning
		v.spinner.show = false
		// Pick new spinner for next time
		v.spinner.spinner.Spinner = spinners[rand.Intn(len(spinners))]
		// Getting a streamMsg means something happened so run async task to get
		// new status on the VM the command was just run on.
		return v, v.getVMStatus(v.getCurrentVM().machineID)

	// Handle error messages (just throw them in the viewport)
	case ecosystemErrMsg:
		v.vagrantOutputView.viewport.SetContent(msg.Error())
	case statusErrMsg:
		v.vagrantOutputView.viewport.SetContent(msg.Error())
	}

	// Spinner needs spinCmd every update to know to keep spinning?
	if v.spinner.show {
		var spinCmd tea.Cmd
		v.spinner.spinner, spinCmd = v.spinner.spinner.Update(msg)
		return v, spinCmd
	} else {
		// Viewport needs vpCmd every update to know to handle user input?
		var vpCmd tea.Cmd
		v.vagrantOutputView.viewport, vpCmd = v.vagrantOutputView.viewport.Update(msg)
		return v, vpCmd
	}
}
