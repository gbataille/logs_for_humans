package main

import (
	"github.com/gbataille/logs_for_humans/consolefmt"
	"github.com/pterm/pterm"
)

func main() {
	// Enable debug messages.
	pterm.EnableDebugMessages()

	formatter := consolefmt.ConsoleFormatter{}
	formatter.Run()
}
