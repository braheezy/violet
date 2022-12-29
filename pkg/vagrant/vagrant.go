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

// A Vagrant client.
type VagrantClient struct {
	// The path to the Vagrant executable.
	ExecPath string
	// Environment variables used when running Vagrant commands
	Env []string
}

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

// Get the version of Vagrant.
// NB: Good way to check things are working
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
		// Did they change the format to the version?
		version = "N/A"
	}
	return version, nil
}

func ReadChanToString(channel chan string) (result string) {
	for value := range channel {
		result += string(value) + "\n"
	}
	return result
}

func (c *VagrantClient) GetGlobalStatus() string {
	output := make(chan string)
	go c.RunCommand("global-status --machine-readable", output)

	result := ReadChanToString(output)
	return result
}

func (c *VagrantClient) GetStatusForID(machineID string) (string, error) {
	output := make(chan string)
	go c.RunCommand(fmt.Sprintf("status %v --machine-readable", machineID), output)

	result := ReadChanToString(output)

	if strings.Contains(result, "Error") {
		return "", errors.New(result)
	}
	return result, nil
}

// Runs a Vagrant command and stream the result back to caller over channel
func (c *VagrantClient) RunCommand(command string, outputCh chan string) {
	defer close(outputCh)
	// Create the Vagrant command and capture its output
	cmd := exec.Command(c.ExecPath, strings.Split(command, " ")...)
	cmd.Env = c.Env

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		outputCh <- string(fmt.Sprintf("Error getting stdout pipe: %v", err))
	}
	cmd.Stderr = cmd.Stdout

	scanner := bufio.NewScanner(stdout)

	done := make(chan struct{})

	err = cmd.Start()
	if err != nil {
		outputCh <- string(fmt.Sprintf("Error executing: %v", err))
	}

	go func() {
		for scanner.Scan() {
			outputCh <- scanner.Text()
		}
		done <- struct{}{}
	}()

	<-done

	err = cmd.Wait()
	if err != nil {
		outputCh <- string(fmt.Sprintf("Error waiting for the script to complete: %v", err))
	}
}

// **************************************************************************
//
//	Extras. Are these opinionated and don't belong in a public package?
//
// **************************************************************************
// Represents the result of a Vagrant command under the context of a single VM.
type MachineInfo struct {
	// Name is the name of the machine.
	Name string
	// Fields is a map of field names to field values.
	Fields map[string]string
}

// Generically parses the output from a Vagrant command and returns the result.
// Multi-machine environments may result in a list of Result objects
func ParseVagrantOutput(output string) []MachineInfo {
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
				// metadata lines de-lineate VMs and are a good place to grab the name.
				if field == "metadata" {
					// Save name if it's the first we've seen
					if result.Name == "" {
						result.Name = m[1]
					} else if result.Name != m[1] {
						// New VM found. Create new Result
						results = append(results, result)
						result = MachineInfo{Name: m[1], Fields: make(map[string]string)}
					}
				} else if field == "machine-id" && result.Fields["machine-id"] != "" {
					// New VM found. Create new Result
					results = append(results, result)
					result = MachineInfo{Name: m[1], Fields: make(map[string]string)}
					result.Fields[field] = m[2]
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
	if len(result.Fields) != 0 {
		results = append(results, result)
	}
	return results
}
