package app

import (
	"github.com/charmbracelet/bubbles/spinner"
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
)

type currentSpinner struct {
	spinner spinner.Model
	show    bool
	title   string
}

func newSpinner() currentSpinner {
	s := spinner.New()
	s.Spinner = spinners[0]
	return currentSpinner{
		spinner: s,
		show:    false,
		title:   "",
	}
}
