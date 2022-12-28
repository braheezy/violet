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
		name := machineInfo.Name
		if name == "" {
			result := client.GetStatusForID(machineInfo.Fields["machine-id"])
			status := vagrant.ParseVagrantOutput(result)
			name = status[0].Name
		}
		vm := VM{
			name:     name,
			provider: machineInfo.Fields["provider-name"],
			state:    machineInfo.Fields["state"],
			home:     machineInfo.Fields["machine-home"],
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
