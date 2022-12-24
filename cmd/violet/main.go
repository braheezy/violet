package main

import (
	"os"

	"github.com/braheezy/violet/internal/app"
	"github.com/muesli/termenv"
)

func main() {
	// Clear screen before running app
	output := termenv.NewOutput(os.Stdout)
	output.ClearScreen()
	app.Run()
}
