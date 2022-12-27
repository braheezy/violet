package app

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/braheezy/violet/pkg/vagrant"
	tea "github.com/charmbracelet/bubbletea"
)

func readChanToString(input chan string) (result string) {
	for line := range input {
		result += line + "\n"
	}
	return result
}

type ecosystemMsg Ecosystem

type ecosystemErrMsg struct{ err error }

func (e ecosystemErrMsg) Error() string { return e.err.Error() }

func getGlobalStatus() tea.Msg {
	client, err := vagrant.NewVagrantClient()
	if err != nil {
		log.Fatal(err)
	}
	output := make(chan string)
	go func() {
		err = client.RunCommand("global-status", output)
	}()
	if err != nil {
		return ecosystemErrMsg{err}
	}
	result := readChanToString(output)
	results := vagrant.ParseVagrantOutput(result)
	return createEcosystem(results)
}
func createEcosystem(results []vagrant.MachineInfo) tea.Msg {
	if results == nil {
		return nil
	}
	// Create the VM struct
	var VMs []VM
	for _, machineInfo := range results {
		name := machineInfo.Name
		if name == "" {
			name = machineInfo.Fields["machine-id"]
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
		envGroups[path.Base(vm.home)] = append(envGroups[path.Base(vm.home)], vm)
	}
	environments := make([]Environment, len(envGroups))
	i := 0
	for envName, vms := range envGroups {
		env := Environment{
			name:       envName,
			VMs:        vms,
			selectedVM: &vms[0],
		}
		environments[i] = env
		i += 1
	}
	return ecosystemMsg{
		environments: environments,
		selectedEnv:  &environments[0],
	}
}

func (v Violet) Init() tea.Cmd {
	return getGlobalStatus
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
