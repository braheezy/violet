package app

import (
	"fmt"
	"strings"

	"github.com/braheezy/violet/pkg/vagrant"
	tea "github.com/charmbracelet/bubbletea"
)

// runMsg is emitted after a command is run.
type runMsg struct {
	content string
	err     error
}
type runErrMsg struct{ err error }

func (e runErrMsg) Error() string { return e.err.Error() }

// Create the tea.Cmd that will run command on the machine specified by identifier.
func (v *Violet) createMachineRunCmd(command string, identifier string) tea.Cmd {
	return func() tea.Msg {
		content, err := v.ecosystem.client.RunCommand(fmt.Sprintf("%v %v", command, identifier))

		if err != nil {
			return runErrMsg{err}
		}

		return runMsg{content: content}
	}
}

// Create the tea.Cmd that will run command in the directory.
func (v *Violet) createEnvRunCmd(command string, dir string) tea.Cmd {
	return func() tea.Msg {
		content, err := v.ecosystem.client.RunCommandInDirectory(command, dir)

		if err != nil {
			return runErrMsg{err}
		}

		return runMsg{content: content}
	}
}

// machineStatusMsg is emitted when status on a machine is received.
type machineStatusMsg struct {
	// identifier is the name or machine-id for this status info
	identifier string
	// Resultant status about machine
	status vagrant.MachineInfo
}

type statusErrMsg struct{ err error }

func (e statusErrMsg) Error() string { return e.err.Error() }

// Create the tea.Cmd that will get status on a machine.
func (v *Violet) createMachineStatusCmd(identifier string) tea.Cmd {
	return func() tea.Msg {
		result, err := v.ecosystem.client.GetStatusForID(identifier)

		if err != nil {
			return statusErrMsg{err}
		}

		vmStatus := vagrant.ParseVagrantOutput(result)[0]
		vmStatus.Fields["state"] = strings.Replace(vmStatus.Fields["state"], "_", " ", -1)

		return machineStatusMsg{
			identifier: identifier,
			status:     vmStatus,
		}
	}
}

// envStatusMsg is emitted when status on an environment is received.
type envStatusMsg struct {
	name   string
	status []vagrant.MachineInfo
}

// Create the tea.Cmd that will get status on an environment.
func (v *Violet) createEnvStatusCmd(env *Environment) tea.Cmd {
	return func() tea.Msg {
		result, _ := v.ecosystem.client.RunCommandInDirectory("status --machine-readable", env.home)

		newStatus := vagrant.ParseVagrantOutput(result)
		return envStatusMsg{
			name:   env.name,
			status: newStatus,
		}
	}
}
