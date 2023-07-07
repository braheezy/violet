package app

import (
	"path"
	"strings"

	"github.com/braheezy/violet/pkg/vagrant"
)

// Ecosystem contains the total Vagrant world information
type Ecosystem struct {
	// Collection of all Vagrant environments
	environments []Environment
	// Reference to a Vagrant client to run commands with
	client *vagrant.VagrantClient
	// Buttons to allow the user to run commands
	commandButtons buttonGroup
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
			machineID: machineInfo.MachineID,
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
			name:     envName,
			VMs:      vms,
			home:     envGroups[envName][0].home,
			hasFocus: true,
		}
		environments = append(environments, env)
	}
	return Ecosystem{
		environments:   environments,
		client:         client,
		commandButtons: newCommandButtons(),
	}, nil
}
