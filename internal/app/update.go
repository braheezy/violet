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
	SelectVM      key.Binding
	Space         key.Binding
	Help          key.Binding
	Quit          key.Binding
}

// Setup the keybinding and help text for each key
var keys = helpKeyMap{
	SelectVM: key.NewBinding(
		key.WithKeys("up", "k", "down", "j"),
		key.WithHelp("↑/k ↓/j", "select vm"),
	),
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
		{k.SelectVM, k.SelectCommand, k.Switch}, // first column
		{k.Space, k.Execute, k.Help, k.Quit},    // second column
	}
}

func (v Violet) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Window was resized
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can it can gracefully truncate
		// its view as needed.
		v.help.Width = msg.Width
		v.terminalWidth = msg.Width
		v.terminalHeight = msg.Height

	// User pressed a key
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, v.keys.Left):
			if v.currentEnv().hasFocus {
				if v.currentEnv().selectedCommand == 0 {
					v.currentEnv().selectedCommand = len(supportedVagrantCommands) - 1
				} else {
					v.currentEnv().selectedCommand--
				}
			} else {
				if v.currentVM().selectedCommand == 0 {
					v.currentVM().selectedCommand = len(supportedVagrantCommands) - 1
				} else {
					v.currentVM().selectedCommand--
				}
			}
		case key.Matches(msg, v.keys.Right):
			if v.currentEnv().hasFocus {
				if v.currentEnv().selectedCommand == len(supportedVagrantCommands)-1 {
					v.currentEnv().selectedCommand = 0
				} else {
					v.currentEnv().selectedCommand++
				}
			} else {
				if v.currentVM().selectedCommand == len(supportedVagrantCommands)-1 {
					v.currentVM().selectedCommand = 0
				} else {
					v.currentVM().selectedCommand++
				}
			}
		case key.Matches(msg, v.keys.Up):
			if v.currentEnv().hasFocus {
				break
			}
			if v.selectedVM == 0 {
				v.selectedVM = len(v.currentEnv().VMs) - 1
			} else {
				v.selectedVM -= 1
			}
		case key.Matches(msg, v.keys.Down):
			if v.currentEnv().hasFocus {
				break
			}
			if v.selectedVM == len(v.currentEnv().VMs)-1 {
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
		case key.Matches(msg, v.keys.ShiftTab):
			if v.selectedEnv == 0 {
				v.selectedEnv = len(v.ecosystem.environments) - 1
			} else {
				v.selectedEnv -= 1
			}
			return v, nil
		case key.Matches(msg, v.keys.Space):
			v.currentEnv().hasFocus = !v.currentEnv().hasFocus
			return v, nil
		case key.Matches(msg, v.keys.Execute):
			if v.currentEnv().hasFocus {
				vagrantCommand := supportedVagrantCommands[v.currentEnv().selectedCommand]
				runCommand := v.createEnvRunCmd(vagrantCommand, v.currentEnv().home)
				v.spinner.show = true
				// This must be sent for the spinner to spin
				tickCmd := v.spinner.spinner.Tick
				return v, tea.Batch(runCommand, tickCmd)
			} else {
				currentVM := v.currentVM()
				vagrantCommand := supportedVagrantCommands[currentVM.selectedCommand]
				/*
					TODO: This doesn't support running commands in a desktop-less environment that doesn't have an external terminal to put commands on. One approach is to use `screen` to create virtual screen.

					Create a virtual screen:
						screen -dmS <session name> <command>
					Connect to it:
						screen -r <session name>
				*/

				if vagrantCommand == "ssh" {
					c := exec.Command("vagrant", "ssh", currentVM.machineID)
					if currentVM.provider == "docker" {
						c = exec.Command("vagrant", "docker-exec", currentVM.name, "-it", "--", "/bin/sh")
						c.Dir = currentVM.home
					}
					runCommand := tea.ExecProcess(c, func(err error) tea.Msg {
						return runMsg{content: "", err: err}
					})
					return v, runCommand
				} else {
					// Run the command async and stream result back
					runCommand := v.createMachineRunCmd(
						vagrantCommand,
						currentVM.machineID,
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
		// Don't have the VM names yet, just machineIDs.
		// Queue up a bunch of async calls to go get those names.
		for _, env := range eco.environments {
			for _, vm := range env.VMs {
				statusCmds = append(statusCmds, v.createMachineStatusCmd(vm.machineID))
			}
		}
		// Set the new ecosystem
		v.ecosystem = eco

		return v, tea.Batch(statusCmds...)

	// New data about a specific VM has come in
	case machineStatusMsg:
		v.spinner.show = false
		v.spinner.verb = verbs[rand.Intn(len(verbs))]
		v.spinner.spinner.Spinner = spinners[rand.Intn(len(spinners))]
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

	case envStatusMsg:
		v.spinner.show = false
		v.spinner.verb = verbs[rand.Intn(len(verbs))]
		v.spinner.spinner.Spinner = spinners[rand.Intn(len(spinners))]

		// Find the env this message is about
		for i, env := range v.ecosystem.environments {
			if msg.name == env.name {
				selectedEnv := &v.ecosystem.environments[i]
				newVMs := make([]VM, 0)
				for _, vmStatus := range msg.status {
					newVM := VM{
						provider: vmStatus.Fields["provider-name"],
						state:    vmStatus.Fields["state"],
						home:     selectedEnv.home,
						name:     vmStatus.Name,
						// Reset the selected command
						selectedCommand: 0,
					}
					newVMs = append(newVMs, newVM)
				}
				selectedEnv.VMs = newVMs
				break
			}
		}
		return v, nil

	// Result from a command has been streamed in
	case runMsg:
		if v.currentEnv().hasFocus {
			return v, v.createEnvStatusCmd(v.currentEnv())
		} else {
			// Getting a runMsg means something happened so run async task to get
			// new status on the VM the command was just run on.
			return v, v.createMachineStatusCmd(v.currentVM().machineID)
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
