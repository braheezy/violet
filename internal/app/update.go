package app

import (
	"reflect"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// helpKeyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type helpKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Switch key.Binding
	Help   key.Binding
	Quit   key.Binding
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
	Switch: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("⭾ tab", "switch views"),
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
			cursorInc = 1
		case key.Matches(msg, v.keys.Down):
			cursorInc = -1
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
		v.ecosystem = Ecosystem(msg)
	}

	return v, nil
}
