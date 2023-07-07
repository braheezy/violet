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
	// Indexes of the respective lists that are currently selected.
	selectedEnv     int
	selectedMachine int
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

	var machines []Machine
	for _, machineInfo := range results {
		machine := Machine{
			machineID: machineInfo.MachineID,
			provider:  machineInfo.Fields["provider-name"],
			state:     strings.Replace(machineInfo.Fields["state"], "_", " ", -1),
			home:      machineInfo.Fields["machine-home"],
		}
		machines = append(machines, machine)
	}
	// Create different envs by grouping machines based on machine-home
	envGroups := make(map[string][]Machine)
	for _, machine := range machines {
		// TODO: Bug if two different paths have the same folder name e.g. /foo/env1 and /bar/env1 will incorrectly be treated the same
		envGroups[path.Base(machine.home)] = append(envGroups[path.Base(machine.home)], machine)
	}
	var environments []Environment
	for envName, machines := range envGroups {
		env := Environment{
			name:     envName,
			machines: machines,
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
