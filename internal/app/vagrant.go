package app

import (
	"fmt"
	"log"
	"strings"

	"github.com/braheezy/violet/pkg/vagrant"
	tea "github.com/charmbracelet/bubbletea"
)

// Order matters here.
var supportedMachineCommands = []string{"up", "halt", "ssh", "reload", "provision"}
var supportedEnvCommands = []string{"up", "halt", "reload", "provision"}
var symbols = map[string]string{
	"up":        "▶",
	"halt":      "■",
	"ssh":       "＞＿ssh",
	"reload":    "↺",
	"provision": "🛠",
}

// runMsg is emitted after a command is run.
type runMsg struct {
	content string
}
type runErrMsg string

// Create the tea.Cmd that will run command on the machine specified by identifier.
func (v *Violet) createMachineRunCmd(command string, identifier string) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Running %v on %v", command, identifier)
		content, err := v.ecosystem.client.RunCommand(fmt.Sprintf("%v %v", command, identifier))

		if err != nil {
			return runErrMsg(vagrant.ParseVagrantError(err.Error()))
		}

		return runMsg{content: content}
	}
}

// Create the tea.Cmd that will run command in the directory.
func (v *Violet) createEnvRunCmd(command string, dir string) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Running %v in %v", command, dir)
		content, err := v.ecosystem.client.RunCommandInDirectory(command, dir)

		if err != nil {
			return runErrMsg(vagrant.ParseVagrantError(err.Error()))
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
		log.Printf("Getting status for %v", identifier)
		result, err := v.ecosystem.client.GetStatusForID(identifier)

		if err != nil {
			return statusErrMsg{err}
		}

		machineStatus := vagrant.ParseVagrantOutput(result)[0]
		machineStatus.Fields["state"] = strings.Replace(machineStatus.Fields["state"], "_", " ", -1)

		return machineStatusMsg{
			identifier: identifier,
			status:     machineStatus,
		}
	}
}

type nameStatusMsg struct {
	machineID string
	name      string
}

type nameStatusErrMsg struct{ err error }

func (e nameStatusErrMsg) Error() string { return e.err.Error() }

// Create the tea.Cmd that will get name of a machine.
func (v *Violet) createNameStatusCmd(identifier string) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Getting status for %v", identifier)
		result, err := v.ecosystem.client.GetStatusForID(identifier)

		if err != nil {
			return nameStatusErrMsg{err}
		}

		machineStatus := vagrant.ParseVagrantOutput(result)[0]

		return nameStatusMsg{
			machineID: identifier,
			name:      machineStatus.Name,
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
		log.Printf("Getting status in %v", env.home)
		result, err := v.ecosystem.client.RunCommandInDirectory("status --machine-readable", env.home)

		if err != nil {
			return statusErrMsg{err}
		}

		newStatus := vagrant.ParseVagrantOutput(result)
		return envStatusMsg{
			name:   env.name,
			status: newStatus,
		}
	}
}
