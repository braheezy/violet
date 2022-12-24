package app

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/braheezy/violet/pkg/vagrant"
	tea "github.com/charmbracelet/bubbletea"
)

type globalStatusMsg Dashboard

type globalStatusErrMsg struct{ err error }

func (e globalStatusErrMsg) Error() string { return e.err.Error() }

func initGlobalStatus() tea.Msg {
	results := []vagrant.MachineInfo{
		{
			Name: "",
			Fields: map[string]string{
				"provider-name": "libvirt",
				"state":         "shutoff",
				"machine-home":  "/home/braheezy/prettybox/runners",
				"machine-id":    "c03b277",
			},
		},
	}

	// environments := make(Environment, len(results))
	// Create the VM struct
	VMs := make([]VM, len(results))
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
	for envName, vms := range envGroups {
		env := Environment{
			name: envName,
			VMs:  vms,
		}
		environments = append(environments, env)
	}

	return globalStatusMsg{
		environments: environments,
		selected:     nil,
	}
}

func (v Violet) Init() tea.Cmd {
	return initGlobalStatus
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
