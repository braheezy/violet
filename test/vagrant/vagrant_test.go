package vagrant_test

import (
	"testing"

	"violet/pkg/vagrant"

	"github.com/stretchr/testify/assert"
)

// Confirm the Vagrant Client can connect be successfully created
func TestNewVagrantClient(t *testing.T) {
	// Test client when binary is available and accessible
	client, err := vagrant.NewVagrantClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)
	// Confirm the client can run
	version, err := client.GetVersion()
	assert.NoError(t, err)
	assert.NotNil(t, version)

	// Test client when binary is not available
	client.ExecPath = "/fake/path/vagrant"
	version, err = client.GetVersion()

	assert.Empty(t, version)
	assert.ErrorContains(t, err, "unable to run vagrant binary")
}

func TestRunCommand_Basic(t *testing.T) {
	client, _ := vagrant.NewVagrantClient()

	// Verify a valid Vagrant command
	out, err := client.RunCommand("global-status")
	assert.NoError(t, err)
	assert.NotEmpty(t, out)

	// Verify an invalid command
	out, err = client.RunCommand("invalid-command")
	assert.Error(t, err)
	assert.Empty(t, out)
}

func TestParseVagrantOutput_Status(t *testing.T) {
	// Test output from the "vagrant status" command.
	output := `1,builder-f35,metadata,provider,libvirt
2,builder-f35,provider-name,libvirt
3,builder-f35,state,shutoff
4,builder-f35,state-human-short,shutoff
5,builder-f35,state-human-long,The Libvirt domain is not running. Run 'vagrant up' to start it.
5,,ui,info,Current machine states:\n\nbuilder-f35               shutoff (libvirt)\n\nThe Libvirt domain is not running. Run 'vagrant up' to start it.
6,$spe_Cat4,metadata,provider,virtualbox
7,$spe_Cat4,provider-name,virtualbox
8,$spe_Cat4,state,running
9,$spe_Cat4,state-human-short,running
10,$spe_Cat4,state-human-long,`

	expected := []vagrant.VagrantOutputResult{
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
	}
	results := vagrant.ParseVagrantOutput(output)

	assert.EqualValues(t, expected, results)
}
