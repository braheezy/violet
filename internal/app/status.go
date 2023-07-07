package app

import (
	"strings"

	"github.com/braheezy/violet/pkg/vagrant"
	tea "github.com/charmbracelet/bubbletea"
)

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
		output := make(chan string)
		go v.RunCommandInProject("status --machine-readable", env.home, output)

		var result string
		for line := range output {
			result += line + "\n"
		}

		newStatus := vagrant.ParseVagrantOutput(result)
		return envStatusMsg{
			name:   env.name,
			status: newStatus,
		}
	}
}
