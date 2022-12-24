package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (v Violet) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can it can gracefully truncate
		// its view as needed.
		v.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, v.keys.Up):
			fallthrough
		case key.Matches(msg, v.keys.Down):
			fallthrough
		case key.Matches(msg, v.keys.Left):
			fallthrough
		case key.Matches(msg, v.keys.Right):
			fallthrough
		case key.Matches(msg, v.keys.Help):
			v.help.ShowAll = !v.help.ShowAll
		case key.Matches(msg, v.keys.Quit):
			return v, tea.Quit
		}

	case globalStatusMsg:
		v.dashboard = Dashboard(msg)
	}

	return v, nil
}
