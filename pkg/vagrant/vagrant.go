package vagrant

import (
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
func (c *VagrantClient) RunCommand(command string) (string, error) {
	// Create the Vagrant command and capture its output
	cmd := exec.Command(c.ExecPath, append(strings.Split(command, " "), "--machine-readable")...)
	cmd.Env = c.Env

	// Run command and get stdout/stderr
	fmt.Print("Running command")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// Result represents the result of a Vagrant command for a single VM.
type VagrantOutputResult struct {
	// Name is the name of the VM.
	Name string
	// Fields is a map of field names to field values.
	Fields map[string]string
}

// Parse parses the output from a Vagrant command and returns the result.
func ParseVagrantOutput(output string) []VagrantOutputResult {
	// Compile regular expressions for each field.
	fields := make(map[string]*regexp.Regexp)
	supportedFields := []string{"metadata", "provider-name", "state", "state-human-long"}
	for _, field := range supportedFields {
		fields[field] = regexp.MustCompile(`^\d+,(\S+),` + field + `,(.+)$`)
	}
	var results []VagrantOutputResult
	var result VagrantOutputResult
	for _, line := range strings.Split(output, "\n") {
		// Check if the line matches any of the fields.
		matched := false
		for field, re := range fields {
			if m := re.FindStringSubmatch(line); m != nil {
				if field == "metadata" {
					// Save name if it's the first we've seen
					if result.Name == "" {
						result.Name = m[1]
						result.Fields = make(map[string]string)
					} else if result.Name != m[1] {
						// New VM found. Create new Result object
						results = append(results, result)
						result = VagrantOutputResult{Name: m[1], Fields: make(map[string]string)}
					}
				} else {
					// Update the result with the field value.
					result.Fields[field] = m[2]
				}
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
	if result.Name != "" {
		results = append(results, result)
	}
	return results
}
