package app

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/braheezy/violet/pkg/vagrant"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func Run() {
	if os.Getenv("VIOLET_DEBUG") != "" {
		if f, err := tea.LogToFile("violet-debug.log", "help"); err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		} else {
			defer f.Close()
		}
	}
	// Set the color palette for the application.
	if lipgloss.HasDarkBackground() {
		theme = defaultDarkTheme
	} else {
		theme = defaultLightTheme
	}
	if _, err := tea.NewProgram(newViolet()).Run(); err != nil {
		log.Fatalf("Could not start program :(\n%v\n", err)
	}
}

func (v Violet) Init() tea.Cmd {
	return getInitialGlobalStatus
}

// Runs on boot to get current Vagrant status on host.
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

// Updates for the entire ecosystem. Usually with results from `global-status`
type ecosystemMsg Ecosystem

type ecosystemErrMsg struct{ err error }

func (e ecosystemErrMsg) Error() string { return e.err.Error() }

// Call `global-status` and translate result into a new Ecosystem
func createEcosystem(client *vagrant.VagrantClient) (Ecosystem, error) {
	// Fetch (not stream) the current global status
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
			state:     strings.Replace(machineInfo.Fields["state"], "_", " ", -1),
			home:      machineInfo.Fields["machine-home"],
		}
		VMs = append(VMs, vm)
	}
	// Create different envs by grouping VMs based on machine-home
	envGroups := make(map[string][]VM)
	for _, vm := range VMs {
		// TODO: Bug if two different paths have the same folder name e.g. /foo/env1 and /bar/env1 will incorrectly be treated the same
		envGroups[path.Base(vm.home)] = append(envGroups[path.Base(vm.home)], vm)
	}
	var environments []Environment
	for envName, vms := range envGroups {
		env := Environment{
			name: envName,
			VMs:  vms,
		}
		environments = append(environments, env)
	}
	return Ecosystem{
		environments: environments,
		client:       client,
	}, nil
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

func (v Violet) getVMStatus(identifier string) tea.Cmd {
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

type runMsg string

func (v Violet) runCommandOnVM(command string, identifier string) tea.Cmd {
	return func() tea.Msg {
		output := make(chan string)
		go v.ecosystem.client.RunCommand(fmt.Sprintf("%v %v", command, identifier), output)
		var content string
		for value := range output {
			content += string(value) + "\n"
		}
		return runMsg(content)
	}
}
