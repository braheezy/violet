package app

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/braheezy/violet/pkg/vagrant"
	tea "github.com/charmbracelet/bubbletea"
)

type ecosystemMsg Ecosystem

type ecosystemErrMsg struct{ err error }

func (e ecosystemErrMsg) Error() string { return e.err.Error() }

func getInitialGlobalStatus() tea.Msg {
	client, err := vagrant.NewVagrantClient()
	if err != nil {
		log.Fatal(err)
	}
	ecosystem, err := createEcosystem(client)
	if err != nil {
		return ecosystemErrMsg{err}
	}
	return ecosystemMsg(ecosystem)
}
func createEcosystem(client *vagrant.VagrantClient) (Ecosystem, error) {
	result := client.GetGlobalStatus()
	var nilEcosystem Ecosystem
	results := vagrant.ParseVagrantOutput(result)
	if results == nil {
		return nilEcosystem, nil
	}
	// Create the VM struct
	var VMs []VM
	for _, machineInfo := range results {
		vm := VM{
			machineID: machineInfo.Fields["machine-id"],
			provider:  machineInfo.Fields["provider-name"],
			state:     machineInfo.Fields["state"],
			home:      machineInfo.Fields["machine-home"],
		}
		VMs = append(VMs, vm)
	}
	envGroups := make(map[string][]VM)
	for _, vm := range VMs {
		// TODO: Bug if two different paths have the same folder name
		envGroups[path.Base(vm.home)] = append(envGroups[path.Base(vm.home)], vm)
	}
	var environments []Environment
	for envName, vms := range envGroups {
		env := Environment{
			name:       envName,
			VMs:        vms,
			selectedVM: &vms[0],
		}
		environments = append(environments, env)
	}
	return Ecosystem{
		environments: environments,
		selectedEnv:  &environments[0],
	}, nil
}
func (v Violet) Init() tea.Cmd {
	return getInitialGlobalStatus
}

type statusMsg struct {
	identifier string
	status     vagrant.MachineInfo
}

type statusErrMsg struct{ err error }

func (e statusErrMsg) Error() string { return e.err.Error() }

func (v Violet) getVMStatus(identifier string) tea.Cmd {
	return func() tea.Msg {
		result, err := v.client.GetStatusForID(identifier)

		if err != nil {
			return statusErrMsg{err}
		}

		return statusMsg{
			identifier: identifier,
			status:     vagrant.ParseVagrantOutput(result)[0],
		}
	}
}

type streamMsg chan string

func (v Violet) streamCommandOnVM(identifier string, command string) tea.Cmd {
	return func() tea.Msg {
		output := make(chan string)
		go v.client.RunCommand(fmt.Sprintf("%v %v", command, identifier), output)
		return streamMsg(output)
	}
}

func Run() {
	if os.Getenv("VIOLET_DEBUG") != "" {
		if f, err := tea.LogToFile("violet-debug.log", "help"); err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		} else {
			defer f.Close()
		}
	}
	if _, err := tea.NewProgram(newViolet()).Run(); err != nil {
		log.Fatalf("Could not start program :(\n%v\n", err)
	}
}
