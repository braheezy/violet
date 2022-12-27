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
	// Define reusable VM and Env struct, and to get their addresses
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
				{
					Name: "vm1",
					Fields: map[string]string{
						"provider-name": "libvirt",
						"state":         "shutoff",
						"machine-home":  "/home/test/env1",
						"machine-id":    "c03b277",
					},
				},
			},
			expected: Ecosystem{
				environments: []Environment{
					testEnvs[0],
				},
				selectedEnv: &testEnvs[0],
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, createEcosystem(test.input))
		})
	}
}
