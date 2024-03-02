package app

/*
The main entry point for the application.
*/

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/braheezy/violet/pkg/vagrant"
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

func Run() {
	if os.Getenv("VIOLET_DEBUG") != "" {
		if f, err := tea.LogToFile("violet-debug.log", "debug"); err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		} else {
			defer f.Close()
		}
	} else {
		// Set up a dummy logger that discards log output
		log.SetOutput(io.Discard)
	}
	// Set the color palette for the application.
	if lipgloss.HasDarkBackground() {
		theme = defaultDarkTheme
	} else {
		theme = defaultLightTheme
	}

	// Setup mouse tracking
	zone.NewGlobal()

	p := tea.NewProgram(newViolet(), tea.WithAltScreen(), tea.WithMouseAllMotion())
	p.SetWindowTitle("♡♡ violet ♡♡")
	if _, err := p.Run(); err != nil {
		log.Fatalf("Could not start program :(\n%v\n", err)
	}
}

// Complete app state (i.e. the BubbleTea model)
type Violet struct {
	// Reference to the Ecosystem
	ecosystem Ecosystem
	// Fancy help bubble
	help help.Model
	// To support help
	keys helpKeyMap
	// Spinner to show while commands are running
	spinner currentSpinner
	// Current terminal size
	terminalWidth  int
	terminalHeight int
	errorMessage   string
}

func (v *Violet) setErrorMessage(message string) {
	v.errorMessage = message
}

// Return the default Violet model
func newViolet() Violet {
	client, err := vagrant.NewVagrantClient()
	if err != nil {
		log.Fatal(err)
	}

	help := help.New()
	help.ShowAll = true

	return Violet{
		ecosystem: Ecosystem{
			environments: nil,
			client:       client,
		},
		keys:    keys,
		help:    help,
		spinner: newSpinner(),
	}
}

func (v Violet) Init() tea.Cmd {
	return getInitialGlobalStatus

}

// Runs on boot to get current Vagrant status on host.
func getInitialGlobalStatus() tea.Msg {
	client, err := vagrant.NewVagrantClient()
	if err != nil {
		log.Fatal(err)
	}
	ecosystem, err := createEcosystem(client)
	if err != nil {
		return ecosystemErrMsg{err}
	}
	return ecosystemMsg(ecosystem)
}
