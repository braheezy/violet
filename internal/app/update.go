package app

import (
	"reflect"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// helpKeyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type helpKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Switch  key.Binding
	Execute key.Binding
	Help    key.Binding
	Quit    key.Binding
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
		key.WithHelp("←/h", "move selection left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move selection right"),
	),
	Switch: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("⭾ tab", "switch views"),
	),
	Execute: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("⏎ enter/return", "run command"),
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
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Switch, k.Help, k.Quit},      // second column
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

func (v Violet) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can it can gracefully truncate
		// its view as needed.
		v.help.Width = msg.Width

	case tea.KeyMsg:
		cursorInc := 0
		switch {
		case key.Matches(msg, v.keys.Up):
			break
		case key.Matches(msg, v.keys.Down):
			break
		case key.Matches(msg, v.keys.Left):
			cursorInc = -1
		case key.Matches(msg, v.keys.Right):
			cursorInc = 1
		case key.Matches(msg, v.keys.Switch):
			switch v.focus {
			case environmentView:
				v.focus = vmView
			case vmView:
				v.focus = commandView
			case commandView:
				v.focus = environmentView
			}
		case key.Matches(msg, v.keys.Execute):
			return v, v.streamCommandOnVM(
				v.selectedCommand,
				v.ecosystem.selectedEnv.selectedVM.name,
				v.ecosystem.selectedEnv.selectedVM.home,
			)
		case key.Matches(msg, v.keys.Help):
			v.help.ShowAll = !v.help.ShowAll
		case key.Matches(msg, v.keys.Quit):
			return v, tea.Quit
		}

		if cursorInc != 0 {
			switch v.focus {
			case environmentView:
				v.ecosystem.selectedEnv = Select(&(v.ecosystem.environments), v.ecosystem.selectedEnv, cursorInc)
			case vmView:
				v.ecosystem.selectedEnv.selectedVM = Select(&(v.ecosystem.selectedEnv.VMs), v.ecosystem.selectedEnv.selectedVM, cursorInc)
			case commandView:
				v.selectedCommand = *(Select(&(v.supportedCommands), &(v.selectedCommand), cursorInc))
			}
		}

	case ecosystemMsg:
		eco := Ecosystem(msg)
		var statusCmds []tea.Cmd
		for _, env := range eco.environments {
			for _, vm := range env.VMs {
				statusCmds = append(statusCmds, v.getVMStatus(vm.machineID))
			}
		}
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
		var content string
		for value := range msg {
			content += string(value) + "\n"
		}
		v.vagrantOutputView.viewport.SetContent(content)

	case ecosystemErrMsg:
		v.vagrantOutputView.viewport.SetContent(msg.Error())
	case statusErrMsg:
		v.vagrantOutputView.viewport.SetContent(msg.Error())
	}

	var vpCmd tea.Cmd
	v.vagrantOutputView.viewport, vpCmd = v.vagrantOutputView.viewport.Update(msg)

	return v, vpCmd
}
