package app

import (
	"fmt"
	"strings"

	"github.com/braheezy/violet/pkg/vagrant"
	tea "github.com/charmbracelet/bubbletea"
)

func (v *Violet) RunCommandInProject(command string, dir string) (output string, err error) {
	v.ecosystem.client.WorkingDir = dir
	output, _ = v.ecosystem.client.RunCommand(command)
	v.ecosystem.client.WorkingDir = ""
	return output, nil
}

type runMsg struct {
	content string
	err     error
}

func (v *Violet) getRunCommandOnVM(command string, identifier string) tea.Cmd {
	return func() tea.Msg {
		content, _ := v.ecosystem.client.RunCommand(fmt.Sprintf("%v %v", command, identifier))
		return runMsg{content: content}
	}
}

func (v *Violet) getRunCommandInVagrantProject(command string, dir string) tea.Cmd {
	return func() tea.Msg {
		content, _ := v.RunCommandInProject(command, dir)

		return runMsg{content: content}
	}
}

// Status messages get emitted on a per-VM basis
type statusMsg struct {
	// identifier is the name or machine-id for this status info
	identifier string
	// Resultant status about machine
	status vagrant.MachineInfo
}

type statusErrMsg struct{ err error }

func (e statusErrMsg) Error() string { return e.err.Error() }

func (v *Violet) getVMStatus(identifier string) tea.Cmd {
	return func() tea.Msg {
		result, err := v.ecosystem.client.GetStatusForID(identifier)

		if err != nil {
			return statusErrMsg{err}
		}

		vmStatus := vagrant.ParseVagrantOutput(result)[0]
		vmStatus.Fields["state"] = strings.Replace(vmStatus.Fields["state"], "_", " ", -1)

		return statusMsg{
			identifier: identifier,
			status:     vmStatus,
		}
	}
}

type envStatusMsg struct {
	name   string
	status []vagrant.MachineInfo
}

func (v *Violet) getEnvStatus(env *Environment) tea.Cmd {
	return func() tea.Msg {
		result, _ := v.RunCommandInProject("status --machine-readable", env.home)

		newStatus := vagrant.ParseVagrantOutput(result)
		return envStatusMsg{
			name:   env.name,
			status: newStatus,
		}
	}
}
