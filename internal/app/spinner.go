package app

import (
	"math/rand"

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
	spinnerCommandStyle = spinnerStyle.
				Bold(true).
				Foreground(secondaryColor)
)

type currentSpinner struct {
	spinner spinner.Model
	show    bool
	verb    string
}

func newSpinner() currentSpinner {
	s := spinner.New()
	s.Spinner = spinners[rand.Intn(len(spinners))]
	return currentSpinner{
		spinner: s,
		verb:    verbs[rand.Intn(len(verbs))],
	}
}
