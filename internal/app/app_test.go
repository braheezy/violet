package app

import (
	"testing"

	"github.com/braheezy/violet/pkg/vagrant"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/require"
)

/*
Test Goals:
   1. Model: Test the data structures and their properties, such as the Environment and VM structs in the Violet model.

   2. View: Test the output of the View method to ensure that it is correct for different input data.

   3. Update: Test the Update method to ensure that it updates the model correctly in response to different messages.

   4. Interactions with the VagrantClient: Test the interactions with the VagrantClient, such as the RunCommand function, to ensure that they are working correctly and returning the expected results.

   5. User input handling: Test the handling of user input, such as keyboard and mouse events, to ensure that they are correctly processed and result in the expected changes to the model.
*/

func Test_createEcosystem(t *testing.T) {
	testMachineInfos := []vagrant.MachineInfo{
		{
			Name: "vm1",
			Fields: map[string]string{
				"provider-name": "libvirt",
				"state":         "shutoff",
				"machine-home":  "/home/test/env1",
				"machine-id":    "c03b277",
			},
		},
		{
			Name: "vm2",
			Fields: map[string]string{
				"provider-name": "libvirt",
				"state":         "running",
				"machine-home":  "/home/test/env1",
				"machine-id":    "23d32r",
			},
		},
		{
			Name: "vm3",
			Fields: map[string]string{
				"provider-name": "virtualbox",
				"state":         "not created",
				"machine-home":  "/home/test/env2",
				"machine-id":    "34reef3",
			},
		},
	}
	testVMs := []VM{
		{
			name:     "vm1",
			provider: "libvirt",
			state:    "shutoff",
			home:     "/home/test/env1",
		},
		{
			name:     "vm2",
			provider: "libvirt",
			state:    "running",
			home:     "/home/test/env1",
		},
		{
			name:     "vm3",
			provider: "virtualbox",
			state:    "not created",
			home:     "/home/test/env2",
		},
	}
	testEnvs := []Environment{
		{
			name: "env1",
			VMs: []VM{
				testVMs[0],
			},
			selectedVM: &testVMs[0],
		},
		{
			name: "env1",
			VMs: []VM{
				testVMs[0],
				testVMs[1],
			},
			selectedVM: &testVMs[0],
		},
		{
			name: "env2",
			VMs: []VM{
				testVMs[2],
			},
			selectedVM: &testVMs[2],
		},
	}

	tests := []struct {
		name     string
		input    []vagrant.MachineInfo
		expected tea.Msg
	}{
		{
			name:     "global-status: empty",
			input:    nil,
			expected: nil,
		},
		{
			name: "global-status: One VM, One Env",
			input: []vagrant.MachineInfo{
				testMachineInfos[0],
			},
			expected: Ecosystem{
				environments: []Environment{
					testEnvs[0],
				},
				selectedEnv: &testEnvs[0],
			},
		},
		{
			name: "global-status: Multi VM, Multi Env",
			input: []vagrant.MachineInfo{
				testMachineInfos[0],
				testMachineInfos[1],
				testMachineInfos[2],
			},
			expected: Ecosystem{
				environments: []Environment{
					testEnvs[1],
					testEnvs[2],
				},
				selectedEnv: &testEnvs[1],
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, createEcosystem(test.input))
		})
	}
}
