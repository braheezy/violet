package vagrant

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Confirm the Vagrant Client can be successfully created
func TestNewVagrantClient(t *testing.T) {
	t.Run("Verify client when binary is available and accessible", func(t *testing.T) {
		client, err := NewVagrantClient()
		require.NoError(t, err)
		require.NotNil(t, client)
		// Confirm the client can run
		version, err := client.GetVersion()
		require.NoError(t, err)
		require.NotNil(t, version)
	})
	t.Run("Verify client when binary is not installed", func(t *testing.T) {
		// Save the current PATH so we can restore it after
		originalPathEnv := os.Getenv("PATH")
		// Destroy PATH for next call
		modifiedPathEnv := "/nothing"
		os.Setenv("PATH", modifiedPathEnv)
		client, err := NewVagrantClient()

		assert.Nil(t, client)
		require.ErrorContains(t, err, "vagrant binary not found")

		// Restore PATH
		os.Setenv("PATH", originalPathEnv)
	})
	t.Run("Verify client when binary is not available", func(t *testing.T) {
		client, _ := NewVagrantClient()
		client.ExecPath = "/fake/path/to/vagrant"

		_, err := client.GetVersion()

		require.ErrorContains(t, err, "unable to run vagrant binary")
	})

}

func TestRunCommand(t *testing.T) {
	client, _ := NewVagrantClient()

	t.Run("Verify a valid Vagrant command", func(t *testing.T) {
		output := make(chan string)
		go client.RunCommand("global-status", output)

		result := ReadChanToString(output)

		require.Greater(t, len(result), 40)
	})
}

func TestRunCommandInDir(t *testing.T) {
	client, _ := NewVagrantClient()

	t.Run("Verify a valid Vagrant command in a directory", func(t *testing.T) {
		output := make(chan string)
		go client.RunCommandInDir("global-status", "/tmp", output)

		result := ReadChanToString(output)

		require.Greater(t, len(result), 40)
	})
}

func TestGetGlobalStatus(t *testing.T) {
	client, _ := NewVagrantClient()

	result := client.GetGlobalStatus()

	require.NotEmpty(t, result)
}

func TestGetStatusForID(t *testing.T) {
	client, _ := NewVagrantClient()

	tests := []struct {
		name      string
		input     string
		expected  string
		wantError bool
	}{
		{
			name:      "Test bad ID",
			input:     "fake",
			expected:  "",
			wantError: true,
		},
		{
			name:     "Test good ID",
			input:    "12deee0",
			expected: "node1,metadata,provider,libvirt",
		},
	}
	for _, test := range tests {
		result, err := client.GetStatusForID(test.input)
		require.Contains(t, result, test.expected)
		if test.wantError {
			require.Error(t, err)
		}
	}
}
func TestParseVagrantOutput_Status(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []MachineInfo
	}{
		{
			name: "Test output from the 'vagrant status' command.",
			input: `1,builder-f35,metadata,provider,libvirt
			2,builder-f35,provider-name,libvirt
			3,builder-f35,state,shutoff
			4,builder-f35,state-human-short,shutoff
			5,builder-f35,state-human-long,The Libvirt domain is not running. Run 'vagrant up' to start it.
			5,,ui,info,Current machine states:\n\nbuilder-f35               shutoff (libvirt)\n\nThe Libvirt domain is not running. Run 'vagrant up' to start it.
			6,$spe_Cat4,metadata,provider,virtualbox
			7,$spe_Cat4,provider-name,virtualbox
			8,$spe_Cat4,state,running
			9,$spe_Cat4,state-human-short,running
			10,$spe_Cat4,state-human-long,`,
			expected: []MachineInfo{
				{
					Name: "builder-f35",
					Fields: map[string]string{
						"provider-name":    "libvirt",
						"state":            "shutoff",
						"state-human-long": "The Libvirt domain is not running. Run 'vagrant up' to start it.",
					},
				},
				{
					Name: "$spe_Cat4",
					Fields: map[string]string{
						"provider-name": "virtualbox",
						"state":         "running",
					},
				},
			},
		},
		{
			name: "Test error status",
			input: `1671329290,,ui,error,A Vagrant environment or target machine is required to run this\ncommand. Run 'vagrant init' to create a new Vagrant environment. Or%!(VAGRANT_COMMA)\nget an ID of a target machine from 'vagrant global-status' to run\nthis command on. A final option is to change to a directory with a\nVagrantfile and to try again.
			1671329290,,error-exit,Vagrant::Errors::NoEnvironmentError,A Vagrant environment or target machine is required to run this\ncommand. Run 'vagrant init' to create a new Vagrant environment. Or%!(VAGRANT_COMMA)\nget an ID of a target machine from 'vagrant global-status' to run\nthis command on. A final option is to change to a directory with a\nVagrantfile and to try again.`,
			expected: nil,
		},
		{
			name:     "Test empty status",
			input:    ``,
			expected: nil,
		},
	}

	for _, test := range tests {
		results := ParseVagrantOutput(test.input)
		assert.EqualValues(t, test.expected, results)
	}

}

func TestParseVagrantOutput_GlobalStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []MachineInfo
	}{
		// {
		// 	name: "Verify empty status",
		// 	input: `1,,metadata,machine-count,0
		// 	2,,ui,info,id
		// 	3,,ui,info,name
		// 	4,,ui,info,provider
		// 	4,,ui,info,state
		// 	5,,ui,info,directory
		// 	6,,ui,info,
		// 	7,,ui,info,--------------------------------------------------------------------
		// 	8,,ui,info,There are no active Vagrant environments on this computer! Or%!(VAGRANT_COMMA)\nyou haven't destroyed and recreated Vagrant environments that were\nstarted with an older version of Vagrant.`,
		// 	expected: nil,
		// },
		// {
		// 	name: "Verify single status",
		// 	input: `1671330325,,metadata,machine-count,1
		// 	1671330325,,machine-id,c03b277
		// 	1671330325,,provider-name,libvirt
		// 	1671330325,,machine-home,/home/braheezy/prettybox/runners
		// 	1671330325,,state,shutoff
		// 	1671330325,,ui,info,id
		// 	1671330325,,ui,info,name
		// 	1671330325,,ui,info,provider
		// 	1671330325,,ui,info,state
		// 	1671330325,,ui,info,directory
		// 	1671330325,,ui,info,
		// 	1671330325,,ui,info,--------------------------------------------------------------------------
		// 	1671330325,,ui,info,c03b277
		// 	1671330325,,ui,info,builder-f35
		// 	1671330325,,ui,info,libvirt
		// 	1671330325,,ui,info,shutoff
		// 	1671330325,,ui,info,/home/braheezy/prettybox/runners
		// 	1671330325,,ui,info,
		// 	1671330325,,ui,info, \nThe above shows information about all known Vagrant environments\non this machine. This data is cached and may not be completely\nup-to-date (use "vagrant global-status --prune" to prune invalid\nentries). To interact with any of the machines%!(VAGRANT_COMMA) you can go to that\ndirectory and run Vagrant%!(VAGRANT_COMMA) or you can use the ID directly with\nVagrant commands from any directory. For example:\n"vagrant destroy 1a2b3c4d"`,
		// 	expected: []MachineInfo{
		// 		{
		// 			Name: "",
		// 			Fields: map[string]string{
		// 				"provider-name": "libvirt",
		// 				"state":         "shutoff",
		// 				"machine-home":  "/home/braheezy/prettybox/runners",
		// 				"machine-id":    "c03b277",
		// 			},
		// 		},
		// 	},
		// },
		{
			name: "Verify multi status",
			input: `1672263560,,metadata,machine-count,3
			1672263560,,machine-id,12deee0
			1672263560,,provider-name,libvirt
			1672263560,,machine-home,/home/braheezy/vagrant-envs/violet-test/env1
			1672263560,,state,running
			1672263560,,machine-id,15b6a07
			1672263560,,provider-name,libvirt
			1672263560,,machine-home,/home/braheezy/vagrant-envs/violet-test/env1
			1672263560,,state,running
			1672263560,,machine-id,200d64a
			1672263560,,provider-name,libvirt
			1672263560,,machine-home,/home/braheezy/vagrant-envs/violet-test/env2
			1672263560,,state,running
			1672263560,,ui,info,id
			1672263560,,ui,info,name
			1672263560,,ui,info,provider
			1672263560,,ui,info,state
			1672263560,,ui,info,directory
			1672263560,,ui,info,
			1672263560,,ui,info,--------------------------------------------------------------------------------
			1672263560,,ui,info,12deee0
			1672263560,,ui,info,node1
			1672263560,,ui,info,libvirt
			1672263560,,ui,info,running
			1672263560,,ui,info,/home/braheezy/vagrant-envs/violet-test/env1
			1672263560,,ui,info,
			1672263560,,ui,info,15b6a07
			1672263560,,ui,info,node2
			1672263560,,ui,info,libvirt
			1672263560,,ui,info,running
			1672263560,,ui,info,/home/braheezy/vagrant-envs/violet-test/env1
			1672263560,,ui,info,
			1672263560,,ui,info,200d64a
			1672263560,,ui,info,server-3
			1672263560,,ui,info,libvirt
			1672263560,,ui,info,running
			1672263560,,ui,info,/home/braheezy/vagrant-envs/violet-test/env2
			1672263560,,ui,info,`,
			expected: []MachineInfo{
				{
					Name: "",
					Fields: map[string]string{
						"provider-name": "libvirt",
						"state":         "running",
						"machine-home":  "/home/braheezy/vagrant-envs/violet-test/env1",
						"machine-id":    "12deee0",
					},
				},
				{
					Name: "",
					Fields: map[string]string{
						"provider-name": "libvirt",
						"state":         "running",
						"machine-home":  "/home/braheezy/vagrant-envs/violet-test/env1",
						"machine-id":    "15b6a07",
					},
				},
				{
					Name: "",
					Fields: map[string]string{
						"provider-name": "libvirt",
						"state":         "running",
						"machine-home":  "/home/braheezy/vagrant-envs/violet-test/env2",
						"machine-id":    "200d64a",
					},
				},
			},
		},
	}

	for _, test := range tests {
		results := ParseVagrantOutput(test.input)
		assert.EqualValues(t, test.expected, results)
	}
}
