package app

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// helpKeyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
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

var verbs = []string{"Running", "Executing", "Performing", "Invoking", "Launching", "Casting"}

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

func (v Violet) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can it can gracefully truncate
		// its view as needed.
		v.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.ChooseCommand):
			selected, _ := strconv.Atoi(msg.String())
			v.ecosystem.environments[v.selectedEnv].VMs[v.selectedVM].selectedCommand = selected - 1
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
			tickCmd := v.spinner.spinner.Tick
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

	case statusMsg:
		for i, env := range v.ecosystem.environments {
			for j, vm := range env.VMs {
				if msg.identifier == vm.machineID || msg.identifier == vm.name {
					// Status msgs don't return some info so retain existing info
					updatedVM := VM{
						machineID: vm.machineID,
						provider:  msg.status.Fields["provider-name"],
						state:     msg.status.Fields["state"],
						home:      vm.home,
						name:      msg.status.Name,
					}
					v.ecosystem.environments[i].VMs[j] = updatedVM
				}
			}
		}

	case streamMsg:
		v.vagrantOutputView.viewport.SetContent(string(msg))
		v.spinner.show = false
		return v, v.getVMStatus(v.getCurrentVM().machineID)

	case ecosystemErrMsg:
		v.vagrantOutputView.viewport.SetContent(msg.Error())
	case statusErrMsg:
		v.vagrantOutputView.viewport.SetContent(msg.Error())
	}

	if v.spinner.show {
		var spinCmd tea.Cmd
		v.spinner.spinner, spinCmd = v.spinner.spinner.Update(msg)
		return v, spinCmd
	} else {
		var vpCmd tea.Cmd
		v.vagrantOutputView.viewport, vpCmd = v.vagrantOutputView.viewport.Update(msg)
		return v, vpCmd
	}
}

// Go `%` operator is not the same as Python, which gives you the actual modulo in return.
// This function does what `%` does in Python.
// https://stackoverflow.com/questions/43018206/modulo-of-negative-integers-in-go
func mod(a, b int) int {
	return (a%b + b) % b
}

func Select[T any](list *[]T, current *T, direction int) *T {
	for i := range *list {
		if reflect.DeepEqual(current, &(*list)[i]) {
			// Return the next item using modulo math
			return &((*list)[mod((i+direction), len(*list))])
		}
	}
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
