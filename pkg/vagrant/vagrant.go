// Package vagrant provides
package vagrant

import (
	"bufio"
	"errors"
	"io"
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
		// Something weird happened
		version = "N/A"
	}
	return version, nil
}

// Runs a Vagrant command and returns the result.
func (c *VagrantClient) RunCommand(command string, outputCh chan<- string) error {
	// Create the Vagrant command and capture its output
	cmd := exec.Command(c.ExecPath, append(strings.Split(command, " "), "--machine-readable")...)
	cmd.Env = c.Env

	// Create pipes for stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}
	// Read pipes as command runs
	scanner := bufio.NewScanner(io.MultiReader(stdoutPipe, stderrPipe))
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			// Send output back to caller via channel
			outputCh <- line
		}
		close(outputCh)
	}()

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

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
						result.Fields = make(map[string]string)
					} else if result.Name != m[1] {
						// New VM found. Create new Result
						results = append(results, result)
						result = MachineInfo{Name: m[1], Fields: make(map[string]string)}
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
	if result.Name != "" || len(result.Fields) != 0 {
		results = append(results, result)
	}
	return results
}
