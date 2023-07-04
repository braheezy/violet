package app

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Available spinners
	spinners = []spinner.Spinner{
		spinner.MiniDot,
		spinner.Dot,
		spinner.Line,
		spinner.Jump,
		spinner.Pulse,
		spinner.Points,
		spinner.Globe,
		spinner.Moon,
		spinner.Monkey,
	}

	spinnerStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Italic(true)
	spinnerCommandStyle = spinnerStyle.Copy().
				Bold(true).
				Foreground(accentColor)
)

type currentSpinner struct {
	spinner spinner.Model
	show    bool
	verb    string
}

func newSpinner() currentSpinner {
	s := spinner.New()
	s.Spinner = spinners[0]
	return currentSpinner{
		spinner: s,
	}
}
