package app

import (
	"math/rand"
	"os/exec"

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
	ShiftTab      key.Binding
	Execute       key.Binding
	SelectCommand key.Binding
	Space         key.Binding
	Help          key.Binding
	Quit          key.Binding
	// These are defined to assist with help text.
	SelectMachine key.Binding
}

// Setup the keybinding and help text for each key
var keys = helpKeyMap{
	SelectCommand: key.NewBinding(
		key.WithKeys("left", "h", "right", "l"),
		key.WithHelp("←/h →/l", "select command"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
	),
	Switch: key.NewBinding(
		key.WithKeys("tab", "shift+tab"),
		key.WithHelp("⭾ tab/⇧+⭾", "switch env tab"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
	),
	Execute: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("⏎ enter", "run selected command"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle env/vm"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	SelectMachine: key.NewBinding(
		key.WithKeys("up", "k", "down", "j"),
		key.WithHelp("↑/k ↓/j", "select vm"),
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
		{k.SelectMachine, k.SelectCommand, k.Switch}, // first column
		{k.Space, k.Execute, k.Help, k.Quit},         // second column
	}
}

func (v Violet) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Window was resized
	case tea.WindowSizeMsg:
		// During development, there were horrific UI bugs when the screen was resized. Things would wrap
		// to the next line. This approach repaints the screen when the happens and it seems to handle all
		// cases. Hopefully this check is good enough to not spam ClearScreen commands.
		needsRepaint := false

		if msg.Width < lipgloss.Width(v.ecosystem.View()) {
			needsRepaint = true
		}

		// If we set a width on the help menu it can it can gracefully truncate
		// its view as needed.
		v.help.Width = msg.Width
		v.terminalWidth = msg.Width
		v.terminalHeight = msg.Height

		if needsRepaint {
			return v, tea.ClearScreen
		}

	// User pressed a key
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, v.keys.Left):
			currentEnv := v.ecosystem.currentEnv()
			currentMachine := v.ecosystem.currentMachine()
			if currentEnv.hasFocus {
				if currentEnv.selectedCommand == 0 {
					currentEnv.selectedCommand = len(supportedVagrantCommands) - 1
				} else {
					currentEnv.selectedCommand--
				}
			} else {
				if currentMachine.selectedCommand == 0 {
					currentMachine.selectedCommand = len(supportedVagrantCommands) - 1
				} else {
					currentMachine.selectedCommand--
				}
			}
		case key.Matches(msg, v.keys.Right):
			currentEnv := v.ecosystem.currentEnv()
			currentMachine := v.ecosystem.currentMachine()
			if currentEnv.hasFocus {
				if currentEnv.selectedCommand == len(supportedVagrantCommands)-1 {
					currentEnv.selectedCommand = 0
				} else {
					currentEnv.selectedCommand++
				}
			} else {
				if currentMachine.selectedCommand == len(supportedVagrantCommands)-1 {
					currentMachine.selectedCommand = 0
				} else {
					currentMachine.selectedCommand++
				}
			}
		case key.Matches(msg, v.keys.Up):
			if v.ecosystem.currentEnv().hasFocus {
				break
			}
			if v.ecosystem.selectedMachine == 0 {
				v.ecosystem.selectedMachine = len(v.ecosystem.currentEnv().machines) - 1
			} else {
				v.ecosystem.selectedMachine -= 1
			}
		case key.Matches(msg, v.keys.Down):
			if v.ecosystem.currentEnv().hasFocus {
				break
			}
			if v.ecosystem.selectedMachine == len(v.ecosystem.currentEnv().machines)-1 {
				v.ecosystem.selectedMachine = 0
			} else {
				v.ecosystem.selectedMachine += 1
			}
		case key.Matches(msg, v.keys.Switch):
			if v.ecosystem.selectedEnv == len(v.ecosystem.environments)-1 {
				v.ecosystem.selectedEnv = 0
			} else {
				v.ecosystem.selectedEnv += 1
			}
			return v, nil
		case key.Matches(msg, v.keys.ShiftTab):
			if v.ecosystem.selectedEnv == 0 {
				v.ecosystem.selectedEnv = len(v.ecosystem.environments) - 1
			} else {
				v.ecosystem.selectedEnv -= 1
			}
			return v, nil
		case key.Matches(msg, v.keys.Space):
			v.ecosystem.currentEnv().hasFocus = !v.ecosystem.currentEnv().hasFocus
			return v, nil
		case key.Matches(msg, v.keys.Execute):
			if v.ecosystem.currentEnv().hasFocus {
				vagrantCommand := supportedVagrantCommands[v.ecosystem.currentEnv().selectedCommand]
				runCommand := v.createEnvRunCmd(vagrantCommand, v.ecosystem.currentEnv().home)
				v.spinner.show = true
				// This must be sent for the spinner to spin
				tickCmd := v.spinner.spinner.Tick
				return v, tea.Batch(runCommand, tickCmd)
			} else {
				currentMachine := v.ecosystem.currentMachine()
				vagrantCommand := supportedVagrantCommands[currentMachine.selectedCommand]
				/*
					TODO: This doesn't support running commands in a desktop-less environment that doesn't have an external terminal to put commands on. One approach is to use `screen` to create virtual screen.

					Create a virtual screen:
						screen -dmS <session name> <command>
					Connect to it:
						screen -r <session name>
				*/

				if vagrantCommand == "ssh" {
					c := exec.Command("vagrant", "ssh", currentMachine.machineID)
					if currentMachine.provider == "docker" {
						c = exec.Command("vagrant", "docker-exec", currentMachine.name, "-it", "--", "/bin/sh")
						c.Dir = currentMachine.home
					}
					runCommand := tea.ExecProcess(c, func(err error) tea.Msg {
						return runMsg{content: "", err: err}
					})
					return v, runCommand
				} else {
					// Run the command async and stream result back
					runCommand := v.createMachineRunCmd(
						vagrantCommand,
						currentMachine.machineID,
					)
					v.spinner.show = true
					// This must be sent for the spinner to spin
					tickCmd := v.spinner.spinner.Tick
					return v, tea.Batch(runCommand, tickCmd)
				}
			}
		case key.Matches(msg, v.keys.Help):
			v.help.ShowAll = !v.help.ShowAll
		case key.Matches(msg, v.keys.Quit):
			return v, tea.Quit
		}

	// New data from `global-status` has come in
	case ecosystemMsg:
		eco := Ecosystem(msg)
		var statusCmds []tea.Cmd
		// Don't have the machine names yet, just machineIDs.
		// Queue up a bunch of async calls to go get those names.
		for _, env := range eco.environments {
			for _, machine := range env.machines {
				statusCmds = append(statusCmds, v.createMachineStatusCmd(machine.machineID))
			}
		}
		// Set the new ecosystem
		v.ecosystem = eco

		return v, tea.Batch(statusCmds...)

	// New data about a specific machine has come in
	case machineStatusMsg:
		v.spinner.show = false
		v.spinner.verb = verbs[rand.Intn(len(verbs))]
		v.spinner.spinner.Spinner = spinners[rand.Intn(len(spinners))]
		// Find the machine this message is about
		for i, env := range v.ecosystem.environments {
			for j, machine := range env.machines {
				if msg.identifier == machine.machineID || msg.identifier == machine.name {
					// Found the machine this status message is about.
					// Status msgs don't return some info so retain existing info
					updateMachine := Machine{
						machineID: machine.machineID,
						provider:  msg.status.Fields["provider-name"],
						state:     msg.status.Fields["state"],
						home:      machine.home,
						name:      msg.status.Name,
						// Reset the selected command
						selectedCommand: 0,
					}
					v.ecosystem.environments[i].machines[j] = updateMachine
				}
			}
		}

	case envStatusMsg:
		v.spinner.show = false
		v.spinner.verb = verbs[rand.Intn(len(verbs))]
		v.spinner.spinner.Spinner = spinners[rand.Intn(len(spinners))]

		// Find the env this message is about
		for i, env := range v.ecosystem.environments {
			if msg.name == env.name {
				selectedEnv := &v.ecosystem.environments[i]
				newMachines := make([]Machine, 0)
				for _, machineStatus := range msg.status {
					newMachine := Machine{
						provider: machineStatus.Fields["provider-name"],
						state:    machineStatus.Fields["state"],
						home:     selectedEnv.home,
						name:     machineStatus.Name,
						// Reset the selected command
						selectedCommand: 0,
					}
					newMachines = append(newMachines, newMachine)
				}
				selectedEnv.machines = newMachines
				break
			}
		}
		return v, nil

	// Result from a command has been streamed in
	case runMsg:
		if v.ecosystem.currentEnv().hasFocus {
			return v, v.createEnvStatusCmd(v.ecosystem.currentEnv())
		} else {
			// Getting a runMsg means something happened so run async task to get
			// new status on the machine the command was just run on.
			return v, v.createMachineStatusCmd(v.ecosystem.currentMachine().machineID)
		}

	// TODO: Handle error messages (just throw them in the viewport)
	case ecosystemErrMsg:
	case statusErrMsg:
	}

	if v.spinner.show {
		var spinCmd tea.Cmd
		v.spinner.spinner, spinCmd = v.spinner.spinner.Update(msg)
		return v, spinCmd
	}

	return v, nil
}
