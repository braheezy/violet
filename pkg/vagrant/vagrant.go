// Package vagrant provides
package vagrant

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// VagrantClient know how to runs Vagrant commands
type VagrantClient struct {
	// The path to the Vagrant executable.
	ExecPath string
	// Environment variables used when running Vagrant commands
	Env []string
	// The working directory for Vagrant commands
	WorkingDir string
}

// NewVagrantClient returns a new VagrantClient ready to run commands.
// Return error if there's issues getting Vagrant binary.
func NewVagrantClient() (*VagrantClient, error) {
	execPath, err := exec.LookPath("vagrant")
	if err != nil {
		return nil, errors.New("vagrant binary not found in PATH")
	}

	return &VagrantClient{
		ExecPath: execPath,
		Env:      os.Environ(),
	}, nil
}

// Return the version of Vagrant.
// NB: Good way to check VagrantClient is working
func (c *VagrantClient) GetVersion() (string, error) {
	cmd := exec.Command(c.ExecPath, "--version")
	result, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New("unable to run vagrant binary")
	}
	// Parse out version string
	version := string(result)
	r := regexp.MustCompile(`Vagrant (\d+.\d+.\d+)`)
	matches := r.FindStringSubmatch(version)

	if len(matches) > 0 {
		version = matches[1]
	} else {
		// Did they change the format of the version?
		version = "N/A"
	}
	return version, nil
}

func (c *VagrantClient) GetGlobalStatus() (result string) {
	result, _ = c.RunCommand("global-status --machine-readable")

	return result
}

func (c *VagrantClient) GetStatusForID(machineID string) (result string, err error) {
	result, err = c.RunCommand(fmt.Sprintf("status %v --machine-readable", machineID))

	if err != nil {
		return "", err
	}
	return result, nil
}

// Run a Vagrant command and stream the result back to caller over outputCh
func (c *VagrantClient) RunCommand(command string) (output string, err error) {
	cmd := exec.Command(c.ExecPath, strings.Split(command, " ")...)
	cmd.Env = c.Env
	cmd.Dir = c.WorkingDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", errors.New("Error creating stdout pipe: " + err.Error())
	}
	cmd.Stderr = cmd.Stdout

	scanner := bufio.NewScanner(stdout)

	err = cmd.Start()
	if err != nil {
		return "", errors.New("Error executing: " + err.Error())
	}

	go func() {
		for scanner.Scan() {
			output += scanner.Text() + "\n"
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return "", errors.New("Error waiting for the command to complete:" + err.Error())
	}
	return output, nil
}

// Represents the result of a Vagrant command under the context of a single VM.
type MachineInfo struct {
	// Name is the name of the machine.
	Name string
	// MachineID is the unique ID of the machine.
	MachineID string
	// Fields is a map of field names to field values.
	Fields map[string]string
}

// Returns true if item is in slice.
func Contains(slice []MachineInfo, item MachineInfo) bool {
	for _, s := range slice {
		if len(s.MachineID) > 0 && s.MachineID == item.MachineID {
			return true
		}
		if len(s.Name) > 0 && s.Name == item.Name {
			return true
		}
	}
	return false
}

// Generically parses the output from a Vagrant command and returns the result.
// Multi-machine environments may result in a list of Result objects
func ParseVagrantOutput(output string) []MachineInfo {
	/*
		This function operators on the output of --machine-readable Vagrant commands: https://developer.hashicorp.com/vagrant/docs/cli/machine-readable.

		The format is:
			timestamp,target,type,data...

			timestamp is a Unix timestamp in UTC of when the message was printed.

			target is the target of the following output. This is empty if the message is related to Vagrant globally. Otherwise, this is generally a machine name so you can relate output to a specific machine when multi-VM is in use.

			type is the type of machine-readable message being outputted. There are a set of standard types which are covered later.

			data is zero or more comma-separated values associated with the prior type. The exact amount and meaning of this data is type-dependent, so you must read the documentation associated with the type to understand fully.
	*/
	// Compile regular expressions for each field.
	fields := make(map[string]*regexp.Regexp)
	// DEV: Add more fields here as needed.
	supportedFields := []string{"metadata", "machine-id", "provider-name", "state", "state-human-long", "machine-home"}
	for _, field := range supportedFields {
		fields[field] = regexp.MustCompile(`^\s*\d+,(.*),` + field + `,(.+)$`)
	}
	var results []MachineInfo
	var result MachineInfo
	result.Fields = make(map[string]string)
	for _, line := range strings.Split(output, "\n") {
		// Check if the line matches any of the fields.
		matched := false
		for field, re := range fields {
			if m := re.FindStringSubmatch(line); m != nil {
				target := m[1]
				// This loops manages findings for machines we might have already seen.
				for _, r := range results {
					// When the name matches (and isn't empty), it's a machine seen before.
					if r.Name == target && r.Name != "" {
						// We're about to switch the result to a different machine so make sure
						// it's in the results list before we "drop" it.
						if !Contains(results, result) {
							results = append(results, result)
						}
						// Update this machine instead of whatever we we're tracking previously.
						result = r
					}
				}
				// Now, fill in machine data, keeping an eye out for new Names or MachineIDs, the clear indication a new machine has been found.
				// metadata lines de-lineate VMs and are a good place to grab the name.
				if field == "metadata" {
					// Save name if it's the first we've seen
					if result.Name == "" {
						// NB: target might be "" too, that's okay.
						result.Name = target
					} else if result.Name != target {
						// New VM found. Create new Result
						results = append(results, result)
						result = MachineInfo{Name: target, Fields: make(map[string]string)}
					}
				} else if field == "machine-id" {
					if result.MachineID != "" {
						// New VM found. Create new Result
						results = append(results, result)
						result = MachineInfo{Name: target, MachineID: m[2], Fields: make(map[string]string)}
					} else {
						result.MachineID = m[2]
					}
				} else {
					// Update the result with the field value.
					result.Fields[field] = m[2]
				}
				// Found a field for this line, move on to next line
				matched = true
				break
			}
		}
		if !matched {
			// No field was matched, move on to the next line.
			continue
		}
	}
	// Add the last result
	if len(result.Fields) != 0 && !Contains(results, result) {
		results = append(results, result)
	}
	return results
}
