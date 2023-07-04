package app

/*
The main entry point for the application.
*/

import (
	"fmt"
	"log"
	"os"

	"github.com/braheezy/violet/pkg/vagrant"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func Run() {
	if os.Getenv("VIOLET_DEBUG") != "" {
		if f, err := tea.LogToFile("violet-debug.log", "help"); err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		} else {
			defer f.Close()
		}
	}
	// Set the color palette for the application.
	if lipgloss.HasDarkBackground() {
		theme = defaultDarkTheme
	} else {
		theme = defaultLightTheme
	}
	p := tea.NewProgram(newViolet(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Could not start program :(\n%v\n", err)
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
